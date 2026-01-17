package pkg

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestOnlyLatest(t *testing.T) {
	uuid1 := Must(uuid.NewUUID())
	uuid2 := Must(uuid.NewUUID())
	baseVoltages := []models.BaseVoltage{
		{
			IdentifiedObject: models.IdentifiedObject{Mrid: uuid1, BaseEntity: models.BaseEntity{CommitId: 0}},
		},
		{
			IdentifiedObject: models.IdentifiedObject{Mrid: uuid2, BaseEntity: models.BaseEntity{CommitId: 1}},
		},
		{
			IdentifiedObject: models.IdentifiedObject{Mrid: uuid1, BaseEntity: models.BaseEntity{CommitId: 1}},
		},
	}

	latest := OnlyLatestVersion(baseVoltages)
	require.Equal(t, 2, len(latest))
}

func TestFindAll(t *testing.T) {
	ctx := context.Background()
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	uuid, err := uuid.NewUUID()
	require.NoError(t, err)

	bv := models.BaseVoltage{
		IdentifiedObject: models.IdentifiedObject{Mrid: uuid, Name: "420kV"},
	}

	_, err = db.NewInsert().Model(&bv).Exec(ctx)
	require.NoError(t, err)

	items, err := FindAll[models.BaseVoltage](db, ctx, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(items))

	nameAndMrid, err := FindNameAndMrid[models.BaseVoltage](db, ctx, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(nameAndMrid))
}

func TestFailedToFetchAllError(t *testing.T) {
	ctx := context.Background()
	db := NewTestConfig().DatabaseConnection()
	_, err := FindNameAndMrid[models.BaseVoltage](db, ctx, 0)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no such table")

}

func TestFinders(t *testing.T) {
	ctx := context.Background()
	db := NewTestConfig().DatabaseConnection()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)
	for k, finder := range Finders {
		_, err := finder(ctx, db, 0)
		require.NoError(t, err, fmt.Sprintf("Failed for %s", k))
	}

}

func TestAllEnumFinders(t *testing.T) {
	ctx := context.Background()
	db := NewTestConfig().DatabaseConnection()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)
	for k, finder := range EnumFinders {
		items, err := finder(ctx, db)
		errorMsg := fmt.Sprintf("Failed for %s", k)
		require.NoError(t, err, errorMsg)
		require.Greater(t, len(items), 0, errorMsg)
	}
}

func TestFilteredFinder(t *testing.T) {
	ctx := context.Background()
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	obj := models.IdentifiedObject{
		Mrid: uuid.New(),
		Name: "My object",
	}

	_, err = db.NewInsert().Model(&obj).Exec(ctx)
	require.NoError(t, err)

	// Add more than 100 objects to test that max 100 is returned
	vls := make([]models.VoltageLevel, 120)
	for i := range len(vls) {
		vls[i].Mrid = uuid.New()
	}

	_, err = db.NewInsert().Model(&vls).Exec(ctx)
	require.NoError(t, err)

	t.Run("error on unkown key", func(t *testing.T) {
		_, err := GetFinder("my random object", "", "")
		require.Error(t, err)
		require.ErrorContains(t, err, "find a finder")
	})

	t.Run("IdentifiedObject", func(t *testing.T) {
		finder, err := GetFinder("IdentifiedObject", "", "")
		require.NoError(t, err)
		result, err := finder(ctx, db, 0)
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
	})

	t.Run("all unfiltered", func(t *testing.T) {
		finder, err := GetFinder("all", "", "")
		require.NoError(t, err)
		result, err := finder(ctx, db, 0)
		require.NoError(t, err)
		require.Equal(t, 100, len(result))
	})

	t.Run("all name filtered", func(t *testing.T) {
		finder, err := GetFinder("all", "my", "")
		require.NoError(t, err)
		result, err := finder(ctx, db, 0)
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
	})

	t.Run("all name filtered no result", func(t *testing.T) {
		finder, err := GetFinder("all", "base voltage", "")
		require.NoError(t, err)
		result, err := finder(ctx, db, 0)
		require.NoError(t, err)
		require.Equal(t, 0, len(result))
	})

	t.Run("all type filtered", func(t *testing.T) {
		finder, err := GetFinder("all", "", "ident")
		require.NoError(t, err)
		result, err := finder(ctx, db, 0)
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
	})

	t.Run("all type filtered no result", func(t *testing.T) {
		finder, err := GetFinder("all", "", "BaseVoltage")
		require.NoError(t, err)
		result, err := finder(ctx, db, 0)
		require.NoError(t, err)
		require.Equal(t, 0, len(result))
	})

	t.Run("return immediatly on failing finder", func(t *testing.T) {
		finder := func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
			return nil, errors.New("something went wrong")
		}

		candidates := map[string]Finder{"failing": finder}
		_, err := AllFinder(ctx, db, 0, candidates, NoOpNameFilter)
		require.Error(t, err)
		require.ErrorContains(t, err, "went wrong")
	})
}

func TestCommitInsertFailsOnNoTables(t *testing.T) {
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	items := func(yield func(v any) bool) {}
	err := InsertAll(context.Background(), db, "insert commit", items, NoOpOnInsert)
	require.Error(t, err)
	require.ErrorContains(t, err, "Failed to insert commit")
}
