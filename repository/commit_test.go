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

func TestInMemInserter(t *testing.T) {
	var inserter InMemInserter
	err := inserter.Insert(context.Background(), "hey!")
	require.NoError(t, err)
	require.Equal(t, 1, len(inserter.Items))
}

func TestInmemInsertWithTx(t *testing.T) {
	var inserter InMemInserter
	fn := WithTx(insertSingleBv)
	err := fn(context.Background(), &inserter)
	require.NoError(t, err)
	require.Equal(t, 1, len(inserter.Items))
}

func insertSingleBv(ctx context.Context, inserter Inserter) error {
	var bv models.BaseVoltage
	return inserter.Insert(ctx, &bv)
}

func insertTwoBv(ctx context.Context, inserter Inserter) error {
	bvs := make([]models.BaseVoltage, 2)
	for i := range bvs {
		bvs[i].Mrid = uuid.New()
	}
	return inserter.Insert(ctx, &bvs)
}

func TestBunInserter(t *testing.T) {
	sql, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sql, sqlitedialect.New())

	ctx := context.Background()
	_, err = db.NewCreateTable().Model((*models.BaseVoltage)(nil)).Exec(ctx)
	require.NoError(t, err)

	bunInsert := BunInserter{Db: db}

	fn := WithTx(insertSingleBv)
	err = fn(ctx, &bunInsert)
	require.NoError(t, err)

	var bvResult []models.BaseVoltage
	err = db.NewSelect().Model(&bvResult).Scan(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(bvResult))

	fn = WithTx(insertTwoBv)
	err = fn(ctx, &bunInsert)
	require.NoError(t, err)

	err = db.NewSelect().Model(&bvResult).Scan(ctx)
	require.NoError(t, err)
	require.Equal(t, 3, len(bvResult))
}
