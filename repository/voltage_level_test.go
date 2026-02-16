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

func TestInSubstation(t *testing.T) {
	vls := make([]models.VoltageLevel, 110)
	for i := range vls {
		vls[i].Mrid = uuid.New()
		vls[i].SubstationMrid = uuid.New()
	}

	ctx := context.Background()

	t.Run("in-memory", func(t *testing.T) {
		var store InMemVoltageLevelReadRepository
		store.Items = vls
		result, err := store.InSubstation(ctx, vls[0].SubstationMrid.String())
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
	})

	t.Run("sqlite", func(t *testing.T) {
		db, err := sql.Open("sqlite3", ":memory:")
		require.NoError(t, err)

		bunDb := bun.NewDB(db, sqlitedialect.New())
		_, err = bunDb.NewCreateTable().Model((*models.VoltageLevel)(nil)).Exec(ctx)
		require.NoError(t, err)

		_, err = bunDb.NewInsert().Model(&vls).Exec(ctx)
		require.NoError(t, err)

		store := NewBunVoltageLevelReadRepository(bunDb)
		result, err := store.InSubstation(ctx, vls[0].SubstationMrid.String())
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
	})
}
