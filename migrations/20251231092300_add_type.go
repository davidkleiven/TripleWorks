package migrations

import (
	"context"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

func init() {
	migrations.MustRegister(addTypeToEntities, revertAddTypeToEntities)
}

func addTypeToEntities(ctx context.Context, db *bun.DB) error {
	dbDialect := db.Dialect().Name()
	if dbDialect == dialect.PG {
		query := "ALTER TABLE entities ADD COLUMN IF NOT EXISTS entity_type TEXT DEFAULT 'unknown'"
		_, err := db.Exec(query)
		return err
	}

	// SQLite
	query := "ALTER TABLE entities ADD COLUMN entity_type"
	_, err := db.Exec(query)
	if err != nil && !isSQLiteDuplicateColumn(err) {
		return err
	}

	backfill := "UPDATE entities SET entity_type = 'unknown' WHERE entity_type IS NULL"
	_, err = db.Exec(backfill)
	return err
}

func revertAddTypeToEntities(ctx context.Context, db *bun.DB) error {
	dbDialect := db.Dialect().Name()
	if dbDialect == dialect.PG {
		query := "ALTER TABLE entities DROP COLUMN IF EXISTS entity_type"
		_, err := db.Exec(query)
		return err
	}

	// SQLite ignore error
	db.Exec("ALTER TABLE entities DROP COLUMN entity_type")
	return nil
}

func isSQLiteDuplicateColumn(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate column") || strings.Contains(msg, "already exists")
}
