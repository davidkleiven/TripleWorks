package pkg

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
)

func IdentifiedLine(mrid uuid.UUID, name string) *models.ACLineSegment {
	var line models.ACLineSegment
	line.Mrid = mrid
	line.Name = name
	return &line
}

func IdentifiedSubstation(mrid uuid.UUID, name string) *models.Substation {
	var substation models.Substation
	substation.Mrid = mrid
	substation.Name = name
	return &substation
}

type E2EData struct {
	Model  models.Model
	Commit models.Commit
	Data   []any
}

func MakeE2eData() *E2EData {
	model := models.Model{Name: "e2e-model", Id: 1}

	concreteKinds := []models.MridGetter{
		// e2e test where Sub A is connected to Sub B
		IdentifiedLine(uuid.MustParse("ce8e57c7-8f6c-42c3-8b8e-e06aa39f0da3"), "Unconnected line"),
		IdentifiedSubstation(uuid.MustParse("fed4f58f-199c-43c7-95f1-b353f55ae12c"), "Substation A"),
		IdentifiedSubstation(uuid.MustParse("8fbd0382-e14c-491b-b4d1-7b2b13be27fb"), "Substation B"),
	}

	var entities []*models.Entity
	for _, item := range concreteKinds {
		entities = append(entities, &models.Entity{
			Mrid:        item.GetMrid(),
			ModelEntity: models.ModelEntity{ModelId: 1},
			EntityType:  StructName(item),
		})
	}

	commit := models.Commit{
		Id:        1,
		Message:   "Populate database with e2e data",
		Author:    "TripleWorks",
		CreatedAt: time.Now(),
	}

	var data []any
	for _, e := range entities {
		data = append(data, e)
	}

	for _, c := range concreteKinds {
		data = append(data, c)
	}

	return &E2EData{
		Model:  model,
		Commit: commit,
		Data:   data,
	}
}

func InsertE2eData(data *E2EData, inserter repository.Inserter) {
	slog.Info("Inserting E2E data")
	err1 := inserter.Insert(context.Background(), &data.Model)

	var numInserted int

	onInsert := func(item any) error {
		numInserted++
		return nil
	}

	err := InsertAllInserter(context.Background(), inserter, data.Commit, slices.Values(data.Data), onInsert)
	slog.Info("Inserted e2e data", "numRecords", numInserted, "errModelInsert", err1, "errEntityInsert", err)
}
