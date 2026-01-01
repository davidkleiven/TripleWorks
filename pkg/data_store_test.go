package pkg

import (
	"context"
	"fmt"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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
	db := NewTestConfig().DatabaseConnection()
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
