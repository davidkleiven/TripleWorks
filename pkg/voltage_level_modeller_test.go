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

func TestSingleLineModel(t *testing.T) {
	data := NewVoltageLevelEquipment(WithLines([]models.ACLineSegment{{}}))
	result := CreateFullyConnectedVoltageLevel(data, NewEmptyConnector())
	require.Equal(t, 0, len(result.Switches), "Switches")
	require.Equal(t, 0, len(result.Terminals), "Terminals")
}

func TestGeneratorAndLine(t *testing.T) {
	var line models.ACLineSegment
	line.Mrid = uuid.New()

	var machine models.SynchronousMachine
	machine.Mrid = uuid.New()
	data := NewVoltageLevelEquipment(
		WithLines([]models.ACLineSegment{line}),
		WithGenerators([]models.SynchronousMachine{machine}),
	)

	connector := NewEmptyConnector()
	for mrid := range data.EquipmentMrids() {
		var terminal models.Terminal
		terminal.Mrid = uuid.New()
		terminal.ConductingEquipmentMrid = mrid
		terminal.ConnectivityNodeMrid = uuid.New()
		connector.AddTerminals(terminal)
	}
	result := CreateFullyConnectedVoltageLevel(data, connector)
	require.Equal(t, 1, len(result.Switches), "Switches")
	require.Equal(t, 2, len(result.Terminals), "Terminals")
}

func TestLoadAndLine(t *testing.T) {
	var line models.ACLineSegment
	line.Mrid = uuid.New()

	var machine models.ConformLoad
	machine.Mrid = uuid.New()
	data := NewVoltageLevelEquipment(
		WithLines([]models.ACLineSegment{line}),
		WithConformLoads([]models.ConformLoad{machine}),
	)
	connector := NewEmptyConnector()
	for mrid := range data.EquipmentMrids() {
		var terminal models.Terminal
		terminal.Mrid = uuid.New()
		terminal.ConductingEquipmentMrid = mrid
		terminal.ConnectivityNodeMrid = uuid.New()
		connector.AddTerminals(terminal)
	}
	result := CreateFullyConnectedVoltageLevel(data, connector)
	require.Equal(t, 1, len(result.Switches), "Switches")
	require.Equal(t, 2, len(result.Terminals), "Terminals")
}

func TestGenAndLoadNotConnectedBySwitch(t *testing.T) {
	data := NewVoltageLevelEquipment(
		WithConformLoads([]models.ConformLoad{{}}),
		WithGenerators([]models.SynchronousMachine{{}}),
	)
	result := CreateFullyConnectedVoltageLevel(data, NewEmptyConnector())
	require.Equal(t, 0, len(result.Switches), "Switches")
	require.Equal(t, 0, len(result.Terminals), "Terminals")
}

func TestLineIsConnectedToLine(t *testing.T) {
	data := NewVoltageLevelEquipment(WithLines([]models.ACLineSegment{{}, {}}))
	for i := range data.Lines {
		data.Lines[i].Mrid = uuid.New()
	}

	connector := NewEmptyConnector()
	for mrid := range data.EquipmentMrids() {
		var terminal models.Terminal
		terminal.Mrid = uuid.New()
		terminal.ConductingEquipmentMrid = mrid
		terminal.ConnectivityNodeMrid = uuid.New()
		connector.AddTerminals(terminal)
	}

	result := CreateFullyConnectedVoltageLevel(data, connector)
	require.Equal(t, 1, len(result.Switches), "Switches")
	require.Equal(t, 2, len(result.Terminals), "Terminals")
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

	var commit models.Commit
	_, err = db.NewInsert().Model(&commit).Exec(ctx)
	require.NoError(t, err)

	zeroMrid := models.Entity{ModelEntity: models.ModelEntity{ModelId: model.Id}, CommitId: int(commit.Id)}
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

	for i := range data.Lines {
		data.Lines[i].Mrid = uuid.New()
	}
	for i := range data.Generators {
		data.Generators[i].Mrid = uuid.New()
	}
	for i := range data.ConformLoads {
		data.ConformLoads[i].Mrid = uuid.New()
	}

	connector := NewEmptyConnector()
	for mrid := range data.EquipmentMrids() {
		var terminal models.Terminal
		terminal.Mrid = uuid.New()
		terminal.ConductingEquipmentMrid = mrid
		terminal.ConnectivityNodeMrid = uuid.New()
		connector.AddTerminals(terminal)
	}

	conNodeEntities := make([]models.Entity, len(connector.terminals))
	for i, terminal := range connector.terminals {
		conNodeEntities[i].Mrid = terminal.ConnectivityNodeMrid
		conNodeEntities[i].EntityType = "ConnectivityNode"
		conNodeEntities[i].ModelId = model.Id
		conNodeEntities[i].CommitId = int(commit.Id)
	}
	_, err = db.NewInsert().Model(&conNodeEntities).Exec(ctx)
	require.NoError(t, err)

	result := CreateFullyConnectedVoltageLevel(data, connector)
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

func TestRepeatedConnectionIsSafe(t *testing.T) {
	data := NewVoltageLevelEquipment(
		WithConformLoads([]models.ConformLoad{{}, {}, {}}),
		WithGenerators([]models.SynchronousMachine{{}, {}, {}}),
		WithLines([]models.ACLineSegment{{}, {}, {}}),
	)

	for i := range data.Lines {
		data.Lines[i].Mrid = uuid.New()
	}
	for i := range data.Generators {
		data.Generators[i].Mrid = uuid.New()
	}
	for i := range data.ConformLoads {
		data.ConformLoads[i].Mrid = uuid.New()
	}

	connector := NewEmptyConnector()
	for mrid := range data.EquipmentMrids() {
		var terminal models.Terminal
		terminal.Mrid = uuid.New()
		terminal.ConductingEquipmentMrid = mrid
		terminal.ConnectivityNodeMrid = uuid.New()
		connector.AddTerminals(terminal)
	}
	result := CreateFullyConnectedVoltageLevel(data, connector)
	require.Greater(t, len(result.Terminals), 0)
	require.Greater(t, len(result.Switches), 0)
	require.Greater(t, len(result.BusNameMarkers), 0)

	result = CreateFullyConnectedVoltageLevel(data, connector)
	require.Equal(t, 0, len(result.Terminals))
	require.Equal(t, 0, len(result.Switches))
	require.Equal(t, 0, len(result.BusNameMarkers))
}

func TestVoltageLEvelFromToDb(t *testing.T) {
	ctx := context.Background()
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	var (
		substation        models.Substation
		bv200             models.BaseVoltage
		bv400             models.BaseVoltage
		vl200             models.VoltageLevel
		gen               models.SynchronousMachine
		load              models.ConformLoad
		line200           models.ACLineSegment
		line400           models.ACLineSegment
		line400NotRelated models.ACLineSegment
		terminal          models.Terminal
	)

	substation.Mrid = uuid.New()
	substation.Name = "Sub7"

	bv200.Mrid = uuid.New()
	bv200.Name = "200 kV"
	bv200.NominalVoltage = 200

	bv400.Mrid = uuid.New()
	bv400.Name = "400 kV"
	bv400.NominalVoltage = 400

	vl200.Mrid = uuid.New()
	vl200.SubstationMrid = substation.Mrid
	vl200.BaseVoltageMrid = bv200.Mrid

	gen.Mrid = uuid.New()
	gen.Name = "Sync machine"
	gen.EquipmentContainerMrid = vl200.Mrid

	load.Mrid = uuid.New()
	load.Name = "Load"
	load.EquipmentContainerMrid = vl200.Mrid

	line200.Mrid = uuid.New()
	line200.Name = "200 kV Sub6 - Sub7"
	line200.BaseVoltageMrid = bv200.Mrid

	line400.Mrid = uuid.New()
	line400.Name = "400 kV Sub7 - Sub10"
	line400.BaseVoltageMrid = bv400.Mrid

	line400NotRelated.Mrid = uuid.New()
	line400NotRelated.Name = "400 kV Sub10 - Sub50"
	line400NotRelated.BaseVoltageMrid = bv400.Mrid

	terminal.Mrid = uuid.New()
	terminal.ConductingEquipmentMrid = line200.Mrid
	terminal.SequenceNumber = 2

	// Add entities
	model := models.ModelEntity{ModelId: 1}
	entities := []models.Entity{
		{ModelEntity: model, Mrid: substation.Mrid, EntityType: StructName(substation)},
		{ModelEntity: model, Mrid: bv200.Mrid, EntityType: StructName(bv200)},
		{ModelEntity: model, Mrid: bv400.Mrid, EntityType: StructName(bv400)},
		{ModelEntity: model, Mrid: vl200.Mrid, EntityType: StructName(vl200)},
		{ModelEntity: model, Mrid: gen.Mrid, EntityType: StructName(gen)},
		{ModelEntity: model, Mrid: load.Mrid, EntityType: StructName(load)},
		{ModelEntity: model, Mrid: line200.Mrid, EntityType: StructName(line200)},
		{ModelEntity: model, Mrid: line400.Mrid, EntityType: StructName(line400)},
		{ModelEntity: model, Mrid: line400NotRelated.Mrid, EntityType: StructName(line400NotRelated)},
		{ModelEntity: model, Mrid: terminal.Mrid, EntityType: StructName(terminal)},
	}

	_, err = db.NewInsert().Model(&entities).Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&substation).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&bv200).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&bv400).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&vl200).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&gen).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&load).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&line200).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&line400).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&line400NotRelated).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&terminal).Exec(ctx)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		results, err := NewVoltageLevelEquipmentFromDb(ctx, db, &substation)
		require.NoError(t, err)
		require.Equal(t, 2, len(results))
		for _, data := range results {
			for _, line := range data.Lines {
				require.NotEqual(t, line.Mrid, line400NotRelated.Mrid)
			}
		}

		// One of the lines has a terminal with sequence number 2. Therefore, this line
		// should have the new terminal equal to 1
		require.Equal(t, gen.Mrid, results[0].Generators[0].Mrid)
		require.Equal(t, load.Mrid, results[0].ConformLoads[0].Mrid)
		require.Equal(t, 0, len(results[1].Generators))
		require.Equal(t, 0, len(results[1].ConformLoads))
	})

	t.Run("error on cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()
		_, err := NewVoltageLevelEquipmentFromDb(cancelledCtx, db, &substation)
		require.Error(t, err)
		require.ErrorContains(t, err, "Failed to fetch")
	})
}

func TestEarlyAbortInEquipmentMrid(t *testing.T) {
	var line models.ACLineSegment
	line.Mrid = uuid.New()

	var gen models.SynchronousMachine
	gen.Mrid = uuid.New()

	var load models.ConformLoad
	load.Mrid = uuid.New()

	data := NewVoltageLevelEquipment(WithLines([]models.ACLineSegment{line}), WithGenerators([]models.SynchronousMachine{gen}), WithConformLoads([]models.ConformLoad{load}))

	for i, stopMrid := range []uuid.UUID{line.Mrid, gen.Mrid, load.Mrid} {
		t.Run(fmt.Sprintf("Stop at %s", stopMrid), func(t *testing.T) {
			count := 0
			for mrid := range data.EquipmentMrids() {
				if mrid == stopMrid {
					break
				}
				count++
			}
			require.Equal(t, i, count)
		})
	}
}
