package migrations

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addLocationToSubstations, revertAddLocationToSubstations)
}

func addLocationToSubstations(ctx context.Context, db *bun.DB) error {
	table := (*models.Substation)(nil)
	var count int
	err := db.NewSelect().Model(table).ColumnExpr("COUNT(location_mrid)").Where("1=1").Scan(ctx, &count)
	if err == nil {
		return nil
	}
	_, err = db.NewAddColumn().Model(table).ColumnExpr("location_mrid UUID").Exec(ctx)
	if err != nil {
		return fmt.Errorf("Failed to add location_mrid column: %w", err)
	}
	return nil
}

func revertAddLocationToSubstations(ctx context.Context, db *bun.DB) error {
	return nil
}
