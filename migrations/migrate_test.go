package migrations

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/migrate"
)

func setupSqliteTestDb(t *testing.T) *bun.DB {
	sqldb, err := sql.Open("sqlite3", "file::memory:?cache-shared")
	require.Nil(t, err)
	return bun.NewDB(sqldb, sqlitedialect.New())
}

func TestAllMigrations(t *testing.T) {
	ctx := context.Background()

	sqldb, err := sql.Open("sqlite3", "file::memory:?cache-shared")
	assert.Nil(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())
	applied, err := RunUp(ctx, db)
	t.Logf("Applied migrations: %d", len(applied.Migrations))

	assert.Nil(t, err)

	migrator := migrate.NewMigrator(db, migrations)

	rolledBack, err := migrator.Rollback(ctx)
	t.Logf("Rolled back %d", len(rolledBack.Migrations))
	assert.Nil(t, err)
	assert.Equal(t, len(rolledBack.Migrations), len(applied.Migrations))
}

func TestAddCommitTable(t *testing.T) {
	db := setupSqliteTestDb(t)
	ctx := context.Background()
	err := createCommitTable(ctx, db)
	require.Nil(t, err)

	var tableCount int
	err = db.NewSelect().ColumnExpr("COUNT(*)").
		Table("sqlite_master").
		Where("type='table' AND name='commits'").
		Scan(ctx, &tableCount)

	require.Nil(t, err)
	require.Equal(t, 1, tableCount)
}

func TestAddModelTable(t *testing.T) {
	db := setupSqliteTestDb(t)
	ctx := context.Background()
	err := createModelTable(ctx, db)
	require.Nil(t, err)

	var tableCount int
	err = db.NewSelect().ColumnExpr("COUNT(*)").
		Table("sqlite_master").
		Where("type='table' AND name='models'").
		Scan(ctx, &tableCount)

	require.Nil(t, err)
	require.Equal(t, 1, tableCount)
}

func TestRunUp_CancelledContext(t *testing.T) {
	ctx := context.Background()

	// Create a cancelled context
	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately

	sqldb, err := sql.Open("sqlite3", "file::memory:?cache-shared")
	assert.Nil(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())

	// This should fail due to cancelled context
	applied, err := RunUp(cancelledCtx, db)

	// Verify error path is taken
	assert.Error(t, err)
	assert.Nil(t, applied)
	assert.Contains(t, err.Error(), "Failed to initialize migrator")
}
