package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(createModelTable, revertCreateModelTable)
}

func createModelTable(ctx context.Context, db *bun.DB) error {
	var model models.Model
	_, err := db.NewCreateTable().
		Model(&model).
		IfNotExists().
		Exec(ctx)
	return err
}

func revertCreateModelTable(ctx context.Context, db *bun.DB) error {
	var model models.Model
	_, err := db.NewDropTable().Model(&model).Exec(ctx)
	return err
}
