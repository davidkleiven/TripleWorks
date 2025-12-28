package migrations

import (
	"context"
	"embed"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations()

//go:embed sql/*
var sqlQueries embed.FS

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
