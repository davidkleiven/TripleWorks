package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addLatestSubstationView, revertAddLatestSubstationView)
}

func addLatestSubstationView(ctx context.Context, db *bun.DB) error {
	sql := MustGetViewSql("substations")
	_, err := db.ExecContext(ctx, sql)
	return err
}

func revertAddLatestSubstationView(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, "DROP VIEW v_substations_latest")
	return err
}
