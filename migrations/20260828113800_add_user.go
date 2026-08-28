package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addUserTable, removeUserTable)
}

func addUserTable(ctx context.Context, db *bun.DB) error {
	model := struct {
		Id    int    `bun:"id,pk,autoincrement"`
		Email string `bun:"email,unique,notnull"`
	}{}
	_, err := db.NewCreateTable().Model(&model).ModelTableExpr("users").IfNotExists().Exec(ctx)
	return err
}

func removeUserTable(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Table("users").Exec(ctx)
	return err
}
