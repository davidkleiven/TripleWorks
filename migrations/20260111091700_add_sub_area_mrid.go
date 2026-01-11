package migrations

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addOpsColsToLoadGroupSubTypes, revertAddOpsColsToLoadGroupSubTypes)
}

var loadGroupTables = []any{(*models.LoadGroup)(nil), (*models.ConformLoadGroup)(nil), (*models.NonConformLoadGroup)(nil)}

func addOpsColsToLoadGroupSubTypes(ctx context.Context, db *bun.DB) error {
	for _, table := range loadGroupTables {
		_, err := db.NewAddColumn().Model(table).ColumnExpr("sub_load_area_mrid UUID").Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to add sub_load_area_mrid column: %w", err)
		}
	}
	return nil
}

func revertAddOpsColsToLoadGroupSubTypes(ctx context.Context, db *bun.DB) error {
	for _, table := range loadGroupTables {
		_, err := db.NewDropColumn().Model(table).Column("sub_load_area_mrid").Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to drop sub_load_area_mrid column: %w", err)
		}
	}
	return nil
}
