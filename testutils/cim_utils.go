package testutils

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ValidTerminal struct {
	Entities []models.Entity
	Terminal *models.Terminal
	Model    *models.Model
	Commit   *models.Commit
}

func CreateValidTerminal() *ValidTerminal {
	commit := models.Commit{}
	model := models.Model{Name: "testmodel"}
	terminalEntity := models.Entity{
		Mrid:        uuid.New(),
		ModelEntity: models.ModelEntity{ModelId: 1},
		EntityType:  pkg.StructName(models.Terminal{}),
	}

	condEquipment := models.Entity{
		Mrid:        uuid.New(),
		ModelEntity: models.ModelEntity{ModelId: 1},
		EntityType:  pkg.StructName(models.ConductingEquipment{}),
	}

	busNameMarker := models.Entity{
		Mrid:        uuid.New(),
		ModelEntity: models.ModelEntity{ModelId: 1},
		EntityType:  pkg.StructName(models.BusNameMarker{}),
	}

	conNode := models.Entity{
		Mrid:        uuid.New(),
		ModelEntity: models.ModelEntity{ModelId: 1},
		EntityType:  pkg.StructName(models.ConnectivityNode{}),
	}

	return &ValidTerminal{
		Entities: []models.Entity{terminalEntity, condEquipment, busNameMarker, conNode},
		Terminal: &models.Terminal{
			ConductingEquipmentMrid: condEquipment.Mrid,
			ACDCTerminal: models.ACDCTerminal{
				BusNameMarkerMrid: busNameMarker.Mrid,
				SequenceNumber:    1,
				IdentifiedObject:  models.IdentifiedObject{BaseEntity: models.BaseEntity{CommitId: 1}},
			},
			PhasesId:           1,
			TerminalOperations: models.TerminalOperations{ConnectivityNodeMrid: conNode.Mrid},
		},
		Model:  &model,
		Commit: &commit,
	}
}

func InsertTerminalFactory(data *ValidTerminal) func(context.Context, bun.Tx) error {
	return func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(data.Model).Exec(ctx); err != nil {
			return fmt.Errorf("Failed to insert model: %w", err)
		}
		if _, err := tx.NewInsert().Model(data.Commit).Exec(ctx); err != nil {
			return fmt.Errorf("Failed to insert commit: %w", err)
		}

		for i, entity := range data.Entities {
			entity.ModelId = data.Model.Id
			if _, err := tx.NewInsert().Model(&entity).Exec(ctx); err != nil {
				return fmt.Errorf("Failed to insert entity %d: %w", i, err)
			}
		}

		data.Terminal.CommitId = int(data.Commit.Id)
		if _, err := tx.NewInsert().Model(data.Terminal).Exec(ctx); err != nil {
			return fmt.Errorf("Failed to insert terminal: %w", err)
		}
		return nil
	}
}
