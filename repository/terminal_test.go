package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func TestTerminalByConNode(t *testing.T) {
	terminals := make([]models.Terminal, 10)
	for i := range terminals {
		terminals[i].Mrid = uuid.New()
		terminals[i].ConnectivityNodeMrid = uuid.New()
	}

	store := InMemTerminalReadRepository{InMemReadRepository: InMemReadRepository[models.Terminal]{Items: terminals}}
	mridIter := func(yield func(v string) bool) {
		if !yield(terminals[1].ConnectivityNodeMrid.String()) {
			return
		}

		if !yield(terminals[7].ConnectivityNodeMrid.String()) {
			return
		}
	}

	ctx := context.Background()
	t.Run("inmem store", func(t *testing.T) {
		result, err := store.WithConnectivityNode(ctx, mridIter)
		require.NoError(t, err)
		require.Equal(t, 2, len(result))
	})

	t.Run("sqlite store", func(t *testing.T) {
		dbUrl := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
		sqldb, err := sql.Open("sqlite3", dbUrl)
		require.NoError(t, err)
		db := bun.NewDB(sqldb, sqlitedialect.New())

		_, err = db.NewCreateTable().Model((*models.Terminal)(nil)).Exec(ctx)
		require.NoError(t, err)

		_, err = db.NewInsert().Model(&terminals).Exec(ctx)
		require.NoError(t, err)

		repo := BunTerminalReadRepository{BunReadRepository: BunReadRepository[models.Terminal]{Db: db}}

		result, err := repo.WithConnectivityNode(ctx, mridIter)
		require.NoError(t, err)
		require.Equal(t, 2, len(result))
	})
}
