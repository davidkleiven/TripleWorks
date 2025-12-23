package migrations

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/migrate"
)

func TestAllMigrations(t *testing.T) {
	ctx := context.Background()

	sqldb, err := sql.Open("sqlite3", "file::memory:?cache-shared")
	assert.Nil(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())
	applied, err := RunUp(ctx, db)
	t.Logf("Applied migrations: %d", len(applied.Migrations))

	assert.Nil(t, err)

	var tableCount int
	err = db.NewSelect().ColumnExpr("COUNT(*)").
		Table("sqlite_master").
		Where("type='table' AND name='commit'").
		Scan(ctx, &tableCount)
	assert.Nil(t, err)
	assert.Equal(t, tableCount, 1)

	migrator := migrate.NewMigrator(db, nil)

	rolledBack, err := migrator.Rollback(ctx)
	t.Logf("Rolled back %d", len(rolledBack.Migrations))
	assert.Nil(t, err)
	assert.Equal(t, len(rolledBack.Migrations), len(applied.Migrations))
}
