package pkg

import (
	"context"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func setupDb(t *testing.T) (*bun.DB, int) {
	ctx := context.Background()

	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)

	_, err = migrations.RunUp(ctx, db)
	require.NoError(t, err)

	model := models.Model{Name: "test model"}
	_, err = db.NewInsert().Model(&model).Exec(ctx)
	require.NoError(t, err)
	return db, model.Id
}

func TestSubstation(t *testing.T) {
	substation := SubstationLight{
		Name:   "Sub A",
		Region: "NO1",
		X:      10.4495,
		Y:      63.2674,
	}

	ctx := context.Background()
	db, modelId := setupDb(t)
	itemIterator := substation.CimItems(modelId)

	err := InsertAll(ctx, db, "Create sub A", itemIterator, NoOpOnInsert)
	require.NoError(t, err)
}

func TestLine(t *testing.T) {
	line := LineLight{
		FromSubstation: "Sub A",
		ToSubstation:   "Sub B",
		Length:         32.0,
		Voltage:        300.0,
	}

	db, modelId := setupDb(t)
	itemIterator := line.CimItems(modelId)
	err := InsertAll(context.Background(), db, "Create line A-B", itemIterator, NoOpOnInsert)
	require.NoError(t, err)
}

func TestGenerator(t *testing.T) {
	gen := GeneratorLight{
		Kind:       "hydro",
		Substation: "Sub A",
		Num:        1,
		MaxP:       10.0,
		Voltage:    300,
	}

	db, modelId := setupDb(t)
	itemIterator := gen.CimItems(modelId)

	var substationEntity models.Entity
	substationEntity.Mrid = substationMrid(gen.Substation)
	substationEntity.EntityType = "Substation"
	substationEntity.ModelId = modelId

	var commit models.Commit
	_, err := db.NewInsert().Model(&commit).Exec(context.Background())
	require.NoError(t, err)

	substationEntity.CommitId = int(commit.Id)
	_, err = db.NewInsert().Model(&substationEntity).Exec(context.Background())
	require.NoError(t, err)

	err = InsertAll(context.Background(), db, "Create generator in substation A", itemIterator, NoOpOnInsert)
	require.NoError(t, err)

	gen.Kind = "wind"
	itemIterator = gen.CimItems(modelId)
	err = InsertAll(context.Background(), db, "Create wind generator in substation A", itemIterator, NoOpOnInsert)
	require.NoError(t, err)

	gen.Kind = "thermal"
	itemIterator = gen.CimItems(modelId)
	err = InsertAll(context.Background(), db, "Create thermal generator in substation A", itemIterator, NoOpOnInsert)
	require.NoError(t, err)
}

func TestBreakInYieldMany(t *testing.T) {
	iterator := func(yield func(v any) bool) {
		yieldMany(yield, 1, 2, 3, 4)
	}

	total := 0
	for integer := range iterator {
		total += integer.(int)
		if total > 5 {
			break
		}
	}

	require.Equal(t, 6, total)
}

func TestConformLoad(t *testing.T) {
	load := LoadLight{
		Substation: "Sub A",
		Num:        1,
		Voltage:    300,
		NominalP:   10.0,
	}

	db, modelId := setupDb(t)
	itemIterator := load.CimItems(modelId)

	var substationEntity models.Entity
	substationEntity.Mrid = substationMrid(load.Substation)
	substationEntity.EntityType = "Substation"
	substationEntity.ModelId = modelId

	var commit models.Commit
	_, err := db.NewInsert().Model(&commit).Exec(context.Background())
	require.NoError(t, err)

	substationEntity.CommitId = int(commit.Id)

	_, err = db.NewInsert().Model(&substationEntity).Exec(context.Background())
	require.NoError(t, err)

	err = InsertAll(context.Background(), db, "Create load in substation A", itemIterator, NoOpOnInsert)
	require.NoError(t, err)
}
