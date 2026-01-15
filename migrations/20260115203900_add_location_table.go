package migrations

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addLocations, revertAddLocations)
}

var locTables = []any{
	(*models.Location)(nil),
	(*models.PositionPoint)(nil),
	(*models.CoordinateSystem)(nil),
}

func addLocations(ctx context.Context, db *bun.DB) error {
	for i, table := range locTables {
		_, err := db.NewCreateTable().
			Model(table).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("addLocations: Failed to create table no. %d: %w", i, err)
		}
	}
	return nil
}

func revertAddLocations(ctx context.Context, db *bun.DB) error {
	for i := len(locTables); i > 0; i-- {
		_, err := db.NewDropTable().Model(locTables[i-1]).IfExists().Exec(ctx)
		if err != nil {
			return fmt.Errorf("addLocations: Failed to drop table no. %d: %w", i, err)
		}
	}
	return nil
}
