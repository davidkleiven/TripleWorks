package repository

import (
	"context"
	"database/sql"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func TestUserRepoTest(t *testing.T) {
	sqldb, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())

	user := models.User{Email: "alex@grid.no"}
	ctx := context.Background()
	_, err = db.NewCreateTable().Model(&user).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&user).Exec(ctx)
	require.NoError(t, err)

	repo := BunUserRepository{
		Db: db,
	}

	result, err := repo.GetByEmail(ctx, "john@example.com")
	require.ErrorIs(t, err, ErrUserNotFound)

	result, err = repo.GetByEmail(ctx, "alex@grid.no")
	require.NoError(t, err)
	require.Equal(t, result.Email, "alex@grid.no")
}
