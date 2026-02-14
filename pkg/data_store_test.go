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

func TestCommitInsertAllStopWhenCallbackFails(t *testing.T) {
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	_, err := migrations.RunUp(context.Background(), db)
	require.NoError(t, err)

	items := func(yield func(v any) bool) { yield(&models.Substation{}) }

	cb := func(v any) error {
		return errors.New("something went wrong")
	}
	err = InsertAll(context.Background(), db, "insert commit", items, cb)
	require.Error(t, err)
	require.ErrorContains(t, err, "went wrong")
}

func TestMridIfPossible(t *testing.T) {
	var substation models.Substation
	substation.Mrid = uuid.New()

	require.Equal(t, mridIfPossible(1), uuid.UUID{})
	require.Equal(t, mridIfPossible(&substation), substation.Mrid)
}

func TestExistingMrids(t *testing.T) {
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	_, err := migrations.RunUp(context.Background(), db)
	require.NoError(t, err)

	entities := []models.Entity{
		{ModelEntity: models.ModelEntity{ModelId: 1}, Mrid: uuid.New()},
		{ModelEntity: models.ModelEntity{ModelId: 1}, Mrid: uuid.New()},
	}
	_, err = db.NewInsert().Model(&entities).Exec(context.Background())
	require.NoError(t, err)

	existing, err := ExistingMrids(context.Background(), db, 1)
	require.NoError(t, err)
	require.Equal(t, 2, len(existing))
}

func TestOnlyNewItems(t *testing.T) {
	existing := uuid.New()
	iterator := func(yield func(v any) bool) {
		var subst models.Substation
		subst.Mrid = existing

		var substNew models.Substation
		substNew.Mrid = uuid.New()

		if !yield(&subst) {
			return
		}

		if !yield(&substNew) {
			return
		}
	}

	existingMrids := map[uuid.UUID]struct{}{existing: {}}
	filteredIterator := OnlyNewItems(existingMrids, iterator)
	num := 0
	for range filteredIterator {
		num++
	}
	require.Equal(t, 1, num)
}

func TestLineConnectedToSubstationByName(t *testing.T) {
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	_, err := migrations.RunUp(context.Background(), db)
	require.NoError(t, err)
	lines := make([]models.ACLineSegment, 20)
	substations := make([]models.Substation, 20)
	locations := make([]models.Location, 20)
	points := make([]models.PositionPoint, 20)
	for i := range len(lines) {
		lines[i].Mrid = uuid.New()
		lines[i].Name = fmt.Sprintf("Sub%d - Sub%d", i, (i+1)%len(substations))

		substations[i].Mrid = uuid.New()
		substations[i].Name = fmt.Sprintf("Sub%d", i)
		substations[i].LocationMrid = uuid.New()

		locations[i].Mrid = substations[i].LocationMrid
		points[i].LocationMrid = locations[i].Mrid
		points[i].XPosition = float64(i)
		points[i].YPosition = float64(i)
	}

	ctx := context.Background()
	_, err = db.NewInsert().Model(&lines).Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&substations).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&locations).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&points).Exec(ctx)
	require.NoError(t, err)

	target := substations[8]

	t.Run("two lines connected to Sub8", func(t *testing.T) {
		linesConnected, err := LinesConnectedToSubstationByName(context.Background(), db, &target)
		require.NoError(t, err)
		connectedNames := make(map[string]struct{})
		for _, line := range linesConnected {
			connectedNames[line.Name] = struct{}{}
		}

		require.Equal(t, 2, len(linesConnected))

		_, ok := connectedNames["Sub7 - Sub8"]
		require.True(t, ok)
	})

	t.Run("error on cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := LinesConnectedToSubstationByName(cancelledCtx, db, &target)
		require.Error(t, err)
	})

	target = substations[1]
	t.Run("two lines connected to Sub1", func(t *testing.T) {
		// Check that we don't get any match for Sub10 which also contains Sub1x
		linesConnected, err := LinesConnectedToSubstationByName(context.Background(), db, &target)
		require.NoError(t, err)

		connectedNames := make(map[string]struct{})
		for _, line := range linesConnected {
			connectedNames[line.Name] = struct{}{}
		}

		require.Equal(t, 2, len(linesConnected))
		require.Equal(t, Set("Sub0 - Sub1", "Sub1 - Sub2"), connectedNames)
	})
}
