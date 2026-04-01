package pkg

import (
	"database/sql"
	"embed"
	"log/slog"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"gopkg.in/yaml.v3"
)

//go:embed profiles/*
var configProfiles embed.FS

func MustGetPredfinedProfile(name string) *Config {
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
	Port    int           `yaml:"port"`
	DbUrl   string        `yaml:"dbUrl"`
	Timeout time.Duration `yaml:"timeout"`
}

func (c *Config) DatabaseConnection() *bun.DB {
	if strings.Contains(c.DbUrl, "postgres") {
		slog.Info("Connecting to postgres database", "url", c.DbUrl)
		sqldb := Must(sql.Open("pgx", c.DbUrl))
		return bun.NewDB(sqldb, pgdialect.New())
	}

	slog.Info("Connecting to sqlite database", "url", c.DbUrl)
	sqldb := Must(sql.Open("sqlite3", c.DbUrl))
	return bun.NewDB(sqldb, sqlitedialect.New())
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
	case "local_pg":
		return MustGetPredfinedProfile(name)
	default:
		return NewDefaultConfig()
	}
}
