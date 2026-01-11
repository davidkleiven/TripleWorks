package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

func init() {
	migrations.MustRegister(addOpsColsToEnergyConsumerSubTypes, revertAddOpsColsToEnergyConsumerSubTypes)
}

var newTables = []string{"conform_loads", "non_conform_loads"}

func addOpsColsToEnergyConsumerSubTypes(ctx context.Context, db *bun.DB) error {
	dbDialect := db.Dialect().Name()
	for _, table := range newTables {
		for i, col := range newColumns {
			var err error
			switch dbDialect {
			case dialect.PG:
				query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s DOUBLE PRECISION DEFAULT 0.0", table, col)
				_, err = db.ExecContext(ctx, query)
			case dialect.SQLite:
				query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s REAL DEFAULT 0.0", table, col)
				_, err = db.ExecContext(ctx, query)

				if err != nil && !isSQLiteDuplicateColumn(err) {
					return fmt.Errorf("Query %d failed: %w", i, err)
				}
				backfill := fmt.Sprintf("UPDATE %s SET %s = 0.0 WHERE %s IS NULL", table, col, col)
				_, err = db.Exec(backfill)
			}

			if err != nil {
				return fmt.Errorf("Query: %d failed: %w", i, err)
			}
		}
	}
	return nil
}

func revertAddOpsColsToEnergyConsumerSubTypes(ctx context.Context, db *bun.DB) error {
	for _, table := range newTables {
		for i, col := range newColumns {
			query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table, col)
			_, err := db.ExecContext(ctx, query)
			if err != nil {
				return fmt.Errorf("Drop column %d failed: %w", i, err)
			}

		}
	}
	return nil
}
