package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addSubstationsPosition, revertSubstations)
}

func addSubstationsPosition(ctx context.Context, db *bun.DB) error {
	query := MustGetQuery("subtations_with_coordinates.sql")
	_, err := db.ExecContext(ctx, query)
	return err
}

func revertSubstations(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, "DROP VIEW substations_geo")
	return err
}
