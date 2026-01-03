package migrations

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

func init() {
	migrations.MustRegister(addOpsColsToEnergyConsumer, revertAddOpsColsToEnergyConsumer)
}

var newColumns = []string{"qfixed", "qfixed_pct", "pfixed", "pfixed_pct"}

func addOpsColsToEnergyConsumer(ctx context.Context, db *bun.DB) error {
	dbDialect := db.Dialect().Name()
	for i, col := range newColumns {
		var err error
		switch dbDialect {
		case dialect.PG:
			query := fmt.Sprintf("ALTER TABLE energy_consumers ADD COLUMN IF NOT EXISTS %s DOUBLE PRECISION DEFAULT 0.0", col)
			_, err = db.ExecContext(ctx, query)
		case dialect.SQLite:
			query := fmt.Sprintf("ALTER TABLE energy_consumers ADD COLUMN %s REAL DEFAULT 0.0", col)
			_, err = db.ExecContext(ctx, query)

			if err != nil && !isSQLiteDuplicateColumn(err) {
				return fmt.Errorf("Query %d failed: %w", i, err)
			}
			backfill := fmt.Sprintf("UPDATE energy_consumers SET %s = 0.0 WHERE %s IS NULL", col, col)
			_, err = db.Exec(backfill)
		}

		if err != nil {
			return fmt.Errorf("Query: %d failed: %w", i, err)
		}
	}
	return nil
}

func revertAddOpsColsToEnergyConsumer(ctx context.Context, db *bun.DB) error {
	for i, col := range newColumns {
		query := fmt.Sprintf("ALTER TABLE energy_consumers DROP COLUMN %s", col)
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("Drop column %d failed: %w", i, err)
		}

	}
	_, err := db.NewDelete().Model(&models.TransformerControlMode{}).Where("1=1").Exec(ctx)
	return err
}
