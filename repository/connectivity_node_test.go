package repository

import (
	"context"
	"database/sql"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func TestInContainer(t *testing.T) {
	cns := make([]models.ConnectivityNode, 10)
	for i := range cns {
		cns[i].Mrid = uuid.New()
		cns[i].ConnectivityNodeContainerMrid = uuid.New()
	}

	ctx := context.Background()
	t.Run("in memory", func(t *testing.T) {
		var store InMemConnectivityNodeReadRepository
		store.Items = cns
		result, err := store.InContainer(ctx, cns[3].ConnectivityNodeContainerMrid.String())
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
	})

	t.Run("sqlite", func(t *testing.T) {
		sql, err := sql.Open("sqlite3", ":memory:")
		require.NoError(t, err)
		db := bun.NewDB(sql, sqlitedialect.New())

		_, err = db.NewCreateTable().Model((*models.ConnectivityNode)(nil)).Exec(ctx)
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&cns).Exec(ctx)
		require.NoError(t, err)

		var store BunConnectivityNodeReadRepository
		store.Db = db

		result, err := store.InContainer(ctx, cns[3].ConnectivityNodeContainerMrid.String())
		require.NoError(t, err)
		require.Equal(t, 1, len(result))

	})
}
