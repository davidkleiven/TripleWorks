package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations()

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
