package pkg

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/caarlos0/env"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"gopkg.in/yaml.v3"
)

//go:embed profiles/*
var configProfiles embed.FS

func MustGetPredefinedProfile(name string) *Config {
	reader := Must(configProfiles.Open("profiles/" + name + ".yaml"))
	config := NewDefaultConfig()
	PanicOnErr(yaml.NewDecoder(reader).Decode(&config))
	return config
}

func ConfigFromExternalFile(name string) *Config {
	config := NewDefaultConfig()
	f, err := os.Open(name)
	if err != nil {
		slog.Error("Could not open file. Using default config", "error", err, "file", name)
		return config
	}

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		slog.Error("Could not decode yaml file. Using default", "error", err)
		return config
	}
	return config
}

type Config struct {
	DbUrl                           string        `yaml:"dbUrl"`
	LoadflowServiceEndpoint         string        `yaml:"load_flow_service" env:"TRIPLEWORKS_LOAD_FLOW_SERVICE"`
	LocalPtdfFolder                 string        `yaml:"local_parquest_folder" env:"TRIPLEWORKS_LOCAL_PTDF_FOLDER"`
	PtdfBucket                      string        `yaml:"parquet_bucket" env:"TRIPLEWORKS_PTDF_BUCKET"`
	Port                            int           `yaml:"port" env:"TRIPLEWORKS_PORT"`
	Timeout                         time.Duration `yaml:"timeout" env:"TRIPLEWORKS_TIMEOUT"`
	WithTailscaleUserIdentification bool          `yaml:"withTailscaleUserIdentification" env:"WITH_TAILSCALE_USER_IDENTIFICATION"`
}

func (c *Config) DatabaseConnection() *bun.DB {
	if strings.Contains(c.DbUrl, "postgres") {
		slog.Info("Connecting to postgres database")
		sqldb := Must(sql.Open("pgx", c.DbUrl))
		return bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	}

	slog.Info("Connecting to sqlite database", "url", c.DbUrl)
	sqldb := Must(sql.Open("sqlite3", c.DbUrl))
	return bun.NewDB(sqldb, sqlitedialect.New(), bun.WithDiscardUnknownColumns())
}

// SafeString returns a loggable (e.g. no secrets) string representation of the config object
func (c *Config) SafeString() string {
	var builder strings.Builder
	builder.WriteString("port=")
	builder.WriteString(strconv.Itoa(c.Port))
	builder.WriteString(", timeout=")
	builder.WriteString(c.Timeout.String())
	builder.WriteString(", withTailscaleUserIdentification=")
	builder.WriteString(strconv.FormatBool(c.WithTailscaleUserIdentification))
	builder.WriteString("ptdfBucket=")
	builder.WriteString(c.PtdfBucket)
	builder.WriteString("localPtdfFolder=")
	builder.WriteString(c.LocalPtdfFolder)
	builder.WriteString("loadflowServiceEndpoint=")
	builder.WriteString(c.LoadflowServiceEndpoint)
	return builder.String()
}

func (c *Config) PtdfWriterFactory() *MultiWriterFactory {
	var factory MultiWriterFactory
	if c.LocalPtdfFolder != "" {
		factory.Factories = append(factory.Factories, &LocalWriterFactory{})
	}
	if c.PtdfBucket != "" {
		client := Must(storage.NewClient(context.Background()))
		factory.Factories = append(factory.Factories, &GcsWriterFactory{Client: client})
	}
	return &factory
}

func NewDefaultConfig() *Config {
	return &Config{
		Port:    36000,
		DbUrl:   "tripleworks.db",
		Timeout: 10 * time.Minute,
	}
}

func NewTestConfig(opts ...func(c *Config)) *Config {
	config := NewDefaultConfig()
	config.DbUrl = "file:memdb?mode=memory&cache=shared"
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func NewEnvParsedConfig() *Config {
	config := NewDefaultConfig()
	env.Parse(config)
	return config
}

func WithDbName(name string) func(c *Config) {
	return func(c *Config) {
		c.DbUrl = strings.ReplaceAll(c.DbUrl, "memdb", name)
	}
}

func GetConfig(name string) *Config {
	if strings.HasSuffix(name, ".yaml") {
		return ConfigFromExternalFile(name)
	}

	switch name {
	case "test":
		return NewTestConfig()
	case "local_pg", "e2e_sqlite":
		return MustGetPredefinedProfile(name)
	case "pg_env":
		return PgEnv(&FsOpener{})
	default:
		return NewDefaultConfig()
	}
}

func PgEnv(opener Opener) *Config {
	config := NewEnvParsedConfig()
	prefix := "TRIPLEWORKS_DB"
	user := os.Getenv(prefix + "_USER")
	port := os.Getenv(prefix + "_PORT")
	host := os.Getenv(prefix + "_HOST")
	db := os.Getenv(prefix + "_DATABASE")
	passwordFile := os.Getenv(prefix + "_PW_FILE")
	f, err := opener.Open(passwordFile)
	if err != nil {
		slog.Error("Could not open file", "error", err)
		return config
	}
	defer f.Close()

	passwordBytes, err := io.ReadAll(f)
	if err != nil {
		slog.Error("Could not read file content", "error", err)
		return config
	}
	password := strings.TrimSpace(string(passwordBytes))
	slog.Info("Loaded postgres config from env", "user", user, "port", port, "host", host, "db", db, "password", loggablePassword(password))
	config.DbUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)
	return config
}

func loggablePassword(password string) string {
	if len(password) < 2 {
		return password
	}
	return password[:2] + "*******"
}
