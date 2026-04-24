package migrations

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addBidzoneView, revertAddBidzoneView)
}

func addBidzoneView(ctx context.Context, db *bun.DB) error {
	substation := MustGetQuery("substation_region.sql")
	crossBorder := MustGetQuery("cross_region_lines.sql")
	_, err1 := db.ExecContext(ctx, substation)
	_, err2 := db.ExecContext(ctx, crossBorder)
	return errors.Join(err1, err2)
}

func revertAddBidzoneView(ctx context.Context, db *bun.DB) error {
	_, err1 := db.ExecContext(ctx, "DROP VIEW v_substation_bidzones_latest")
	_, err2 := db.ExecContext(ctx, "DROP VIEW v_cross_region_lines_latest")
	return errors.Join(err1, err2)
}
