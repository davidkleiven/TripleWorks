package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addIndexOnCommit, removeIndexOnCommit)
}

func addIndexOnCommit(ctx context.Context, db *bun.DB) error {
	sql := "CREATE INDEX IF NOT EXISTS idx_commits_id_created_at ON commits (id, created_at DESC)"
	_, err := db.ExecContext(ctx, sql)
	return err
}

func removeIndexOnCommit(ctx context.Context, db *bun.DB) error {
	sql := "DROP INDEX IF EXISTS idx_commits_id_created_at"
	_, err := db.ExecContext(ctx, sql)
	return err
}
