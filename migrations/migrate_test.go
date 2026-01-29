package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/testutils"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func setupSqliteTestDb(t *testing.T) *bun.DB {
	sqldb, err := sql.Open("sqlite3", "file::memory:?cache-shared")
	require.Nil(t, err)
	return bun.NewDB(sqldb, sqlitedialect.New())
}

func setupPostgresTestDb(t *testing.T) *bun.DB {
	sqldb, err := sql.Open("pgx", pgTestUrl())
	require.NoError(t, err)
	return bun.NewDB(sqldb, pgdialect.New())
}

func pgTestUrl() string {
	url, ok := os.LookupEnv("POSTGRES_TESTDB")
	if ok {
		return url
	}
	return "postgres://test:test@localhost:5432/testdb?sslmode=disable"
}

func skipLocallyIfNoConnection(t *testing.T) {
	_, isCi := os.LookupEnv("CI")
	sqldb, err := sql.Open("pgx", pgTestUrl())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = sqldb.PingContext(ctx)
	if err != nil && isCi {
		t.Fatalf("Could not context postgres db: %s", err)
	} else if err != nil {
		t.Skip("Could not contact postgres db")
	}
}

func clearPgDataBase(ctx context.Context, db *bun.DB) error {
	rows, err := db.QueryContext(ctx, "SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		return fmt.Errorf("Failed to create query: %w", err)
	}
	defer rows.Close()

	tabNo := 0
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("Failed for table %d: %w", tabNo, err)
		}
		deleteQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		_, err = db.ExecContext(ctx, deleteQuery)
		if err != nil {
			return fmt.Errorf("Failed to delete table %s (%d): %w", table, tabNo, err)
		}
		tabNo++
	}
	return nil
}

func TestAllMigrations(t *testing.T) {
	for _, test := range []struct {
		db   *bun.DB
		name string
	}{
		{
			db:   setupSqliteTestDb(t),
			name: "Sqlite",
		},
		{
			db:   setupPostgresTestDb(t),
			name: "Postgres",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "Postgres" {
				skipLocallyIfNoConnection(t)
				defer clearPgDataBase(context.Background(), test.db)
			}
			ctx := context.Background()
			applied, err := RunUp(ctx, test.db)
			t.Logf("Applied migrations: %d", len(applied.Migrations))

			assert.NoError(t, err)

			rolledBack, err := RunDown(ctx, test.db)
			require.NoError(t, err)
			require.Equal(t, len(applied.Migrations), len(rolledBack.Migrations))
		})
	}

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

type SqliteForeignKey struct {
	Id       int    `bun:"id"`
	Seq      int    `bun:"seq"`
	Table    string `bun:"table"`
	From     string `bun:"from"`
	To       string `bun:"to"`
	OnUpdate string `bun:"on_update"`
	OnDelete string `bun:"on_delete"`
	Match    string `bun:"match"`
}

func TestCreateCim16TablesSqlite(t *testing.T) {
	ctx := context.Background()
	sqldb := setupSqliteTestDb(t)
	err := createCim16Tables(ctx, sqldb)
	require.NoError(t, err)

	_, err = sqldb.NewRaw("PRAGMA foreign_keys = ON").Exec(ctx)
	require.NoError(t, err)

	var fks []SqliteForeignKey
	err = sqldb.NewRaw("PRAGMA foreign_key_list(terminals)").Scan(ctx, &fks)
	require.NoError(t, err)

	has_equipment := false
	for _, fk := range fks {
		if fk.From == "conducting_equipment_mrid" {
			has_equipment = true
		}
	}
	require.True(t, has_equipment)
}

func invalidTerminal() *models.Terminal {
	return &models.Terminal{
		PhasesId:                1,
		ConductingEquipmentMrid: uuid.New(),
	}
}

func TestInvalidInsertFails(t *testing.T) {
	ctx := context.Background()
	sqldb := setupSqliteTestDb(t)
	err := createCim16Tables(ctx, sqldb)
	require.NoError(t, err)

	_, err = sqldb.NewRaw("PRAGMA foreign_keys = ON").Exec(ctx)
	require.NoError(t, err)
	terminal := invalidTerminal()
	_, err = sqldb.NewInsert().Model(terminal).Exec(ctx)
	require.Error(t, err)
}

func TestValidInsertOkSqlite(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		db   *bun.DB
		name string
	}{
		{
			db:   setupSqliteTestDb(t),
			name: "Sqlite",
		},
		{
			db:   setupPostgresTestDb(t),
			name: "Postgres",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "Postgres" {
				skipLocallyIfNoConnection(t)
				defer clearPgDataBase(ctx, test.db)
			}
			_, err := RunUp(ctx, test.db)

			require.NoError(t, err)

			data := testutils.CreateValidTerminal()
			err = test.db.RunInTx(ctx, nil, testutils.InsertTerminalFactory(data))
			assert.NoError(t, err)

			// Must delete terminal to ensure rollback works
			_, err = test.db.NewDelete().Table("terminals").Where("1=1").Exec(ctx)
			assert.NoError(t, err)
			_, err = RunDown(ctx, test.db)
			assert.NoError(t, err)
		})
	}
}

func TestInvalidInsertPostgres(t *testing.T) {
	skipLocallyIfNoConnection(t)
	ctx := context.Background()
	sqldb := setupPostgresTestDb(t)
	_, err := RunUp(ctx, sqldb)
	defer clearPgDataBase(ctx, sqldb)
	require.NoError(t, err)
	terminal := invalidTerminal()
	_, err = sqldb.NewInsert().Model(terminal).Exec(ctx)
	require.Error(t, err)
}

func TestCreateCim16Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel to trigger error
	sqldb := setupSqliteTestDb(t)
	err := createCim16Tables(ctx, sqldb)
	require.Error(t, err)
	require.ErrorContains(t, err, "no. 0")
}

func TestCreateCim16ErrorDown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel to trigger error
	sqldb := setupSqliteTestDb(t)
	err := revertCreateCim16Tables(ctx, sqldb)
	require.Error(t, err)
	require.ErrorContains(t, err, "no. 143")
}

func TestRunDownErrorCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel to trigger error
	sqldb := setupSqliteTestDb(t)
	_, err := RunDown(ctx, sqldb)
	require.Error(t, err)
	require.ErrorContains(t, err, "Failed to initialize")
}

func TestPopulateEnumWithErrorContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel to trigger error
	sqldb := setupSqliteTestDb(t)
	err := populateEnumTables(ctx, sqldb)
	require.Error(t, err)
	require.ErrorContains(t, err, "Failed to insert enum")

	err = revertPopulateEnumTables(ctx, sqldb)
	require.Error(t, err)
	require.ErrorContains(t, err, "Failed to clear table")
}

func TestEnumsArePopulated(t *testing.T) {
	skipLocallyIfNoConnection(t)
	ctx := context.Background()
	sqldb := setupPostgresTestDb(t)
	_, err := RunUp(ctx, sqldb)
	defer clearPgDataBase(ctx, sqldb)
	require.NoError(t, err)

	var phaseCodes []models.PhaseCode
	err = sqldb.NewSelect().Model(&phaseCodes).Scan(ctx)
	require.NoError(t, err)
	require.Greater(t, len(phaseCodes), 0)

	var windUnitKind []models.WindGenUnitKind
	err = sqldb.NewSelect().Model(&windUnitKind).Scan(ctx)
	require.Greater(t, len(windUnitKind), 0)
}

func TestIsSqliteDuplicateColumnError(t *testing.T) {
	sqldb := setupSqliteTestDb(t)
	create := "CREATE TABLE duplication (name TEXT)"
	alter := "ALTER TABLE duplication ADD COLUMN name TEXT"
	_, err := sqldb.Exec(create)
	require.NoError(t, err)
	_, err = sqldb.Exec(alter)
	require.Error(t, err)
	require.True(t, isSQLiteDuplicateColumn(err))
}

func TestIsSqliteErrorNoErrorIsFalse(t *testing.T) {
	require.False(t, isSQLiteDuplicateColumn(nil))
}

func TestErrorDuringSqlMigrationOfAddColumn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sqldb := setupSqliteTestDb(t)
	err := addTypeToEntities(ctx, sqldb)
	require.Error(t, err)
}

func TestAddConNodeExistingTerminals(t *testing.T) {
	sqldb := setupSqliteTestDb(t)
	terminal := models.Terminal{}

	ctx := context.Background()
	_, err := sqldb.NewCreateTable().Model((*models.Terminal)(nil)).Exec(ctx)
	require.NoError(t, err)

	_, err = sqldb.NewCreateTable().Model((*models.ConnectivityNode)(nil)).Exec(ctx)
	require.NoError(t, err)

	_, err = sqldb.NewInsert().Model(&terminal).Exec(context.Background())
	require.NoError(t, err)

	addConNodeToTerminals(ctx, sqldb)

	num, err := sqldb.NewSelect().Model((*models.ConnectivityNode)(nil)).Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, num)

}

func TestAddNodeErrorOnCancelledContext(t *testing.T) {
	sqldb := setupSqliteTestDb(t)
	ctx := context.Background()
	err := addConNodeToTerminals(ctx, sqldb)
	require.Error(t, err)
	require.ErrorContains(t, err, "no such table")
}

func TestCanNotInsertZeroSequenceNumber(t *testing.T) {
	db := setupPostgresTestDb(t)
	ctx := context.Background()
	defer clearPgDataBase(ctx, db)
	_, err := RunUp(ctx, db)
	require.NoError(t, err)
	var terminal models.Terminal
	terminal.Mrid = uuid.New()
	_, err = db.NewInsert().Model(&terminal).Exec(ctx)
	require.Error(t, err)
	require.ErrorContains(t, err, "sequence_number_positive")
}

func TestCanNotInsertZeroNominalVoltage(t *testing.T) {
	db := setupPostgresTestDb(t)
	ctx := context.Background()
	defer clearPgDataBase(ctx, db)
	_, err := RunUp(ctx, db)
	require.NoError(t, err)

	model := models.Model{Name: "test model"}
	_, err = db.NewInsert().Model(&model).Exec(ctx)
	require.NoError(t, err)

	var commit models.Commit
	commit.Message = "add base voltage"
	_, err = db.NewInsert().Model(&commit).Exec(ctx)
	require.NoError(t, err)

	var bv models.BaseVoltage
	bv.Mrid = uuid.New()
	bv.CommitId = int(commit.Id)

	_, err = db.NewInsert().Model(&bv).Exec(ctx)
	require.Error(t, err)
	require.ErrorContains(t, err, "nominal_voltage_positive")
}
