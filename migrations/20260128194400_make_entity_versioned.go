package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

func init() {
	migrations.MustRegister(addCommitInfoToEntities, revertAddCommitInfoToEntities)
}

func addCommitInfoToEntities(ctx context.Context, db *bun.DB) error {
	queryTemplate := "ALTER TABLE entities ADD COLUMN %s commit_id %s DEFAULT %s NOT NULL"
	intType := "INTEGER"
	ifNotExistClause := ""

	if db.Dialect().Name() == dialect.PG {
		intType = "INT"
		ifNotExistClause = "IF NOT EXISTS"

	}
	query := fmt.Sprintf(queryTemplate, ifNotExistClause, intType, "FALSE")
	_, err := db.ExecContext(ctx, query)
	if err != nil && db.Dialect().Name() == dialect.PG {
		return err
	}
	return nil
}

func revertAddCommitInfoToEntities(ctx context.Context, db *bun.DB) error {
	return nil
}
