package repository

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func baseVoltageStore() *InMemReadRepository[models.BaseVoltage] {
	baseVoltages := make([]models.BaseVoltage, 5)
	for i := range baseVoltages {
		baseVoltages[i].Mrid = uuid.New()
		baseVoltages[i].CommitId = i
	}
	return &InMemReadRepository[models.BaseVoltage]{Items: baseVoltages}
}

func TestInMemGetById(t *testing.T) {
	repo := baseVoltageStore()
	mrid := repo.Items[0].Mrid
	bv, err := repo.GetByMrid(context.Background(), mrid.String())
	require.NoError(t, err)
	require.Equal(t, mrid, bv.Mrid)
}

func TestInMemErrorOnNotFound(t *testing.T) {
	repo := baseVoltageStore()
	_, err := repo.GetByMrid(context.Background(), "0000-0000")
	require.Error(t, err)
}

func TestInMemList(t *testing.T) {
	repo := baseVoltageStore()
	bvs, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Equal(t, 5, len(bvs))
}

func TestInMemListByMrid(t *testing.T) {
	repo := baseVoltageStore()
	mrids := []string{repo.Items[1].Mrid.String(), repo.Items[4].Mrid.String()}
	bvs, err := repo.ListByMrids(context.Background(), slices.Values(mrids))
	require.NoError(t, err)
	require.Equal(t, 2, len(bvs))
}

func dbWithBaseVoltageTable(t *testing.T) *bun.DB {
	dburl := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	sqldb, err := sql.Open("sqlite3", dburl)
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	_, err = db.NewCreateTable().Model((*models.BaseVoltage)(nil)).Exec(context.Background())
	require.NoError(t, err)
	return db
}

func TestBunReadRepository(t *testing.T) {
	db := dbWithBaseVoltageTable(t)
	repo := BunReadRepository[models.BaseVoltage]{Db: db}
	bvs := make([]models.BaseVoltage, 3)

	bvs[0].Mrid = uuid.New()
	bvs[0].CommitId = 1
	bvs[1].Mrid = uuid.New()
	bvs[2].Mrid = bvs[0].Mrid
	ctx := context.Background()
	_, err := db.NewInsert().Model(&bvs).Exec(ctx)
	require.NoError(t, err)

	t.Run("get by mrid", func(t *testing.T) {
		result, err := repo.GetByMrid(ctx, bvs[0].Mrid.String())
		require.NoError(t, err)
		require.Equal(t, 1, result.CommitId)
	})

	t.Run("list", func(t *testing.T) {
		result, err := repo.List(ctx)
		require.NoError(t, err)
		require.Equal(t, 3, len(result), "Should list all history")
	})

	t.Run("list by mrid", func(t *testing.T) {
		iter := func(yield func(string) bool) { yield(bvs[0].Mrid.String()) }
		result, err := repo.ListByMrids(ctx, iter)
		require.NoError(t, err)
		require.Equal(t, 2, len(result), "There are two version of mrid 0")
	})

}

func TestFailingRepo(t *testing.T) {
	ctx := context.Background()
	repo := FailingReadRepo[models.BaseVoltage]{}
	t.Run("test by mrid", func(t *testing.T) {
		_, err := repo.GetByMrid(ctx, "0000-000")
		require.Error(t, err)
	})

	t.Run("list has error", func(t *testing.T) {
		_, err := repo.List(ctx)
		require.Error(t, err)
	})

	t.Run("list by mrid has error", func(t *testing.T) {
		iter := func(yield func(v string) bool) {}
		_, err := repo.ListByMrids(ctx, iter)
		require.Error(t, err)
	})
}
