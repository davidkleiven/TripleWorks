package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(createCurrentLimitTable, revertCurrentLimit)
}

func createCurrentLimitTable(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model((*models.CurrentLimit)(nil)).
		IfNotExists().
		WithForeignKeys().
		Exec(ctx)
	return err
}

func revertCurrentLimit(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.CurrentLimit)(nil)).IfExists().Exec(ctx)
	return err
}
