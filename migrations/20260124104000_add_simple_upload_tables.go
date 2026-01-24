package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(createSimpleUploadTable, revertCreateSimpleUploadTable)
}

func createSimpleUploadTable(ctx context.Context, db *bun.DB) error {
	var simpelUpload models.SimpleUpload
	_, err := db.NewCreateTable().
		Model(&simpelUpload).
		IfNotExists().
		Exec(ctx)
	return err
}

func revertCreateSimpleUploadTable(ctx context.Context, db *bun.DB) error {
	var simpleUpload models.SimpleUpload
	_, err := db.NewDropTable().Model(&simpleUpload).Exec(ctx)
	return err
}
