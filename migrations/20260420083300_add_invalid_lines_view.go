package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addInvalidLines, removeInvalidLines)
}

func addInvalidLines(ctx context.Context, db *bun.DB) error {
	query := MustGetQuery("invalid_lines.sql")
	_, err := db.ExecContext(ctx, query)
	return err
}

func removeInvalidLines(ctx context.Context, db *bun.DB) error {
	query := "DROP VIEW IF EXISTS v_invalid_lines"
	_, err := db.ExecContext(ctx, query)
	return err
}
