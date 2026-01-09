package pkg

import (
	"context"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSingleLineModel(t *testing.T) {
	data := NewVoltageLevelEquipment(WithLines([]models.ACLineSegment{{}}))
	data.LineTerminalNumbers[uuid.UUID{}] = 2
	result := CreateFullyConnectedVoltageLevel(data)
	require.Equal(t, 1, len(result.Switches), "Switches")
	require.Equal(t, 3, len(result.Terminals), "Terminals")
	require.Equal(t, 2, len(result.ConnectivityNodes), "ConnectivityNodes")
	require.Equal(t, result.Terminals[0].SequenceNumber, 2)
}

func TestGeneratorAndLine(t *testing.T) {
	data := NewVoltageLevelEquipment(
		WithLines([]models.ACLineSegment{{}}),
		WithGenerators([]models.SynchronousMachine{{}}),
	)
	result := CreateFullyConnectedVoltageLevel(data)
	require.Equal(t, 2, len(result.Switches), "Switches")
	require.Equal(t, 6, len(result.Terminals), "Terminals")
	require.Equal(t, 3, len(result.ConnectivityNodes), "ConnectivityNodes")
}

func TestLoadAndLine(t *testing.T) {
	data := NewVoltageLevelEquipment(
		WithLines([]models.ACLineSegment{{}}),
		WithConformLoads([]models.ConformLoad{{}}),
	)
	result := CreateFullyConnectedVoltageLevel(data)
	require.Equal(t, 2, len(result.Switches), "Switches")
	require.Equal(t, 6, len(result.Terminals), "Terminals")
	require.Equal(t, 3, len(result.ConnectivityNodes), "ConnectivityNodes")
}

func TestGenAndLoadNotConnectedBySwitch(t *testing.T) {
	data := NewVoltageLevelEquipment(
		WithConformLoads([]models.ConformLoad{{}}),
		WithGenerators([]models.SynchronousMachine{{}}),
	)
	result := CreateFullyConnectedVoltageLevel(data)
	require.Equal(t, 0, len(result.Switches), "Switches")
	require.Equal(t, 2, len(result.Terminals), "Terminals")
	require.Equal(t, 2, len(result.ConnectivityNodes), "ConnectivityNodes")
}

func TestLineIsConnectedToLine(t *testing.T) {
	data := NewVoltageLevelEquipment(WithLines([]models.ACLineSegment{{}, {}}))
	result := CreateFullyConnectedVoltageLevel(data)
	require.Equal(t, 3, len(result.Switches), "Switches")
	require.Equal(t, 8, len(result.Terminals), "Terminals")
	require.Equal(t, 4, len(result.ConnectivityNodes), "ConnectivityNodes")
}

func TestInsertVoltageModelToDb(t *testing.T) {
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	ctx := context.Background()
	migrations, err := migrations.RunUp(ctx, db)
	t.Logf("Running %d migrations", len(migrations.Migrations))
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, "PRAGMA foreign_keys=ON;")
	require.NoError(t, err)

	model := models.Model{Name: "Test model"}
	_, err = db.NewInsert().Model(&model).Exec(ctx)
	require.NoError(t, err, "Model")

	zeroMrid := models.Entity{ModelEntity: models.ModelEntity{ModelId: model.Id}}
	_, err = db.NewInsert().Model(&zeroMrid).Exec(ctx)
	require.NoError(t, err)

	var entities []models.Entity
	err = db.NewSelect().Model(&entities).Scan(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(entities))
	require.Equal(t, uuid.UUID{}, entities[0].Mrid)

	data := NewVoltageLevelEquipment(
		WithConformLoads([]models.ConformLoad{{}, {}, {}}),
		WithGenerators([]models.SynchronousMachine{{}, {}, {}}),
		WithLines([]models.ACLineSegment{{}, {}, {}}),
	)
	result := CreateFullyConnectedVoltageLevel(data)
	err = result.Write(ctx, db, model.Id, "Add voltage level")
	require.NoError(t, err)
}

func TestLineMrids(t *testing.T) {
	lines := []models.ACLineSegment{{}}
	mrids := LineMrids(lines)
	require.Equal(t, 1, len(mrids))
	require.Equal(t, uuid.UUID{}, mrids[0])
}

func TestGetTargetTerminalSequenceNumber(t *testing.T) {
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	ctx := context.Background()
	migrations, err := migrations.RunUp(ctx, db)
	t.Logf("Running %d migrations", len(migrations.Migrations))
	require.NoError(t, err)

	terminals := make([]models.Terminal, 10)
	lines := make([]uuid.UUID, len(terminals))
	for i := range terminals {
		terminals[i].Mrid = uuid.New()
		terminals[i].ConductingEquipmentMrid = uuid.New()
		terminals[i].SequenceNumber = i%2 + 1
		lines[i] = terminals[i].ConductingEquipmentMrid
	}

	_, err = db.NewInsert().Model(&terminals).Exec(ctx)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		result, err := GetTargetTerminalSequenceNumber(ctx, db, lines)
		require.NoError(t, err)

		// Terminal 1 exists, thus terminal 2 should be created
		require.Equal(t, 2, result[terminals[0].ConductingEquipmentMrid])

		// Terminal 2 exists, thus terminal 1 should be created
		require.Equal(t, 1, result[terminals[1].ConductingEquipmentMrid])
	})

	t.Run("error on multi terminal line", func(t *testing.T) {
		terminal := terminals[0]
		terminal.Id = 0
		terminal.SequenceNumber = 2
		terminal.Mrid = uuid.New()
		_, err := db.NewInsert().Model(&terminal).Exec(ctx)
		require.NoError(t, err)
		_, err = GetTargetTerminalSequenceNumber(ctx, db, lines)
		require.Error(t, err)
		require.ErrorContains(t, err, terminal.ConductingEquipmentMrid.String())
	})

	t.Run("failed to fetch", func(t *testing.T) {
		ctxWithCancel, cancel := context.WithCancel(ctx)
		cancel()
		_, err = GetTargetTerminalSequenceNumber(ctxWithCancel, db, lines)
		require.Error(t, err)
		require.ErrorContains(t, err, "fetch existing")
	})

}
