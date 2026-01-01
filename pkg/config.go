package pkg

import (
	"database/sql"
	"log/slog"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

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
		Timeout: 10 * time.Second,
	}
}

func NewTestConfig() *Config {
	config := NewDefaultConfig()
	config.DbUrl = "file::memory:?cache-shared"
	return config
}
