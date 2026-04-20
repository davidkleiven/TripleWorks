package migrations

import (
	"context"
	"embed"
	"fmt"
	"io"
	"path/filepath"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

//go:embed sql/*
var sqlQueries embed.FS
var migrations = migrate.NewMigrations()

func MustGetQuery(name string) string {
	file, err := sqlQueries.Open(filepath.Join("sql", name))
	PanicOnErr(err)
	defer file.Close()

	query, err := io.ReadAll(file)
	PanicOnErr(err)
	return string(query)
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func RunUp(ctx context.Context, db *bun.DB) (*migrate.MigrationGroup, error) {
	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return nil, fmt.Errorf("Failed to initialize migrator: %w", err)
	}
	return migrator.Migrate(ctx)
}

func RunDown(ctx context.Context, db *bun.DB) (*migrate.MigrationGroup, error) {
	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return nil, fmt.Errorf("Failed to initialize migrator: %w", err)
	}
	return migrator.Rollback(ctx)
}
