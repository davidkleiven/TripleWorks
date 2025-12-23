package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(createCommitTable, revertCreateCommitTable)
}

func createCommitTable(ctx context.Context, db *bun.DB) error {
	var commit models.Commit
	_, err := db.NewCreateTable().
		Model(commit).
		IfNotExists().
		Exec(ctx)
	return err
}

func revertCreateCommitTable(ctx context.Context, db *bun.DB) error {
	var commit models.Commit
	_, err := db.NewDropTable().Model(commit).Exec(ctx)
	return err
}
