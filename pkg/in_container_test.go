package pkg

import (
	"context"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

var targetVl = uuid.New()

func populateVoltageLevels(t *testing.T) *bun.DB {
	terminals := make([]models.Terminal, 10)
	for i := range terminals {
		err := faker.FakeData(&terminals[i])
		require.NoError(t, err)
	}

	conNodes := make([]models.ConnectivityNode, len(terminals))

	for i := range conNodes {
		err := faker.FakeData(&conNodes[i])
		require.NoError(t, err)
		terminals[i].ConnectivityNodeMrid = conNodes[i].Mrid
	}

	vls := make([]models.VoltageLevel, 2)
	for i := range vls {
		err := faker.FakeData(&vls[i])
		require.NoError(t, err)
	}
	vls[0].Mrid = targetVl

	for i := range conNodes {
		conNodes[i].ConnectivityNodeContainerMrid = vls[i%len(vls)].Mrid
	}

	lines := make([]models.ACLineSegment, 10)
	for i := range lines {
		err := faker.FakeData(&lines[i])
		require.NoError(t, err)
	}

	conformLoads := make([]models.ConformLoad, 10)
	for i := range conformLoads {
		err := faker.FakeData(&conformLoads[i])
		require.NoError(t, err)
	}

	gens := make([]models.SynchronousMachine, 10)
	for i := range gens {
		err := faker.FakeData(&gens[i])
		require.NoError(t, err)
	}

	nonConformLoads := make([]models.NonConformLoad, 10)
	for i := range nonConformLoads {
		err := faker.FakeData(&nonConformLoads[i])
		require.NoError(t, err)
	}

	switches := make([]models.Switch, 10)
	for i := range switches {
		err := faker.FakeData(&switches[i])
		require.NoError(t, err)
	}

	transformer := make([]models.PowerTransformer, 10)
	for i := range transformer {
		err := faker.FakeData(&transformer[i])
		require.NoError(t, err)
	}

	terminals[0].ConductingEquipmentMrid = lines[0].Mrid
	terminals[1].ConductingEquipmentMrid = conformLoads[0].Mrid
	terminals[2].ConductingEquipmentMrid = gens[0].Mrid
	terminals[3].ConductingEquipmentMrid = nonConformLoads[0].Mrid
	terminals[4].ConductingEquipmentMrid = switches[0].Mrid
	terminals[5].ConductingEquipmentMrid = transformer[0].Mrid

	ctx := context.Background()
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	// Insert in dependency order: VoltageLevels -> ConnectivityNodes -> Terminals -> equipment
	_, err = db.NewInsert().Model(&vls).Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&conNodes).Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&terminals).Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&lines).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&conformLoads).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&gens).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&nonConformLoads).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&switches).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&transformer).Exec(ctx)
	require.NoError(t, err)
	return db
}

func TestEquipmentInVoltageLevel(t *testing.T) {
	db := populateVoltageLevels(t)

	ctx := context.Background()
	result, err := FetchInVoltageLevelData(ctx, db, targetVl.String())
	require.NoError(t, err)

	require.Equal(t, 5, len(result.Terminals))
	require.Equal(t, 5, len(result.ConNodes))

	vls := make(map[uuid.UUID]struct{})
	for _, cn := range result.ConNodes {
		vls[cn.ConnectivityNodeContainerMrid] = struct{}{}
	}
	require.Equal(t, 1, len(vls))
	require.Equal(t, 0, len(result.ConformLoads))
	require.Equal(t, 0, len(result.NonConformLoads))
	require.Equal(t, 0, len(result.Transformer))
	require.Equal(t, 1, len(result.Gens))
	require.Equal(t, 1, len(result.Lines))
	require.Equal(t, 1, len(result.Switches))
}

func TestOnlyPickLatest(t *testing.T) {
	var data InVoltageLevel
	err := faker.FakeData(&data)
	require.NoError(t, err)

	origNumConNodes := len(data.ConNodes)
	origNumConformLoads := len(data.ConformLoads)
	origNumGens := len(data.Gens)
	origNumLines := len(data.Lines)
	origNumNonConformLoads := len(data.NonConformLoads)
	origNumSwitches := len(data.Switches)
	origNumTerminals := len(data.Terminals)
	origNumTransformer := len(data.Transformer)

	uniqueConNodeMrids := make(map[uuid.UUID]struct{})
	for _, cn := range data.ConNodes {
		uniqueConNodeMrids[cn.Mrid] = struct{}{}
	}
	require.Equal(t, origNumConNodes, len(uniqueConNodeMrids), "Ensure unique mrids")

	// Duplicate everything
	data.ConNodes = append(data.ConNodes, data.ConNodes...)
	data.ConformLoads = append(data.ConformLoads, data.ConformLoads...)
	data.Gens = append(data.Gens, data.Gens...)
	data.Lines = append(data.Lines, data.Lines...)
	data.NonConformLoads = append(data.NonConformLoads, data.NonConformLoads...)
	data.Switches = append(data.Switches, data.Switches...)
	data.Terminals = append(data.Terminals, data.Terminals...)
	data.Transformer = append(data.Transformer, data.Transformer...)

	for i := range data.ConNodes {
		data.ConNodes[i].CommitId = i
		data.ConNodes[i].Deleted = false
	}
	for i := range data.ConformLoads {
		data.ConformLoads[i].CommitId = i
		data.ConformLoads[i].Deleted = false
	}
	for i := range data.Gens {
		data.Gens[i].CommitId = i
		data.Gens[i].Deleted = false
	}
	for i := range data.Lines {
		data.Lines[i].CommitId = i
		data.Lines[i].Deleted = false
	}
	for i := range data.NonConformLoads {
		data.NonConformLoads[i].CommitId = i
		data.NonConformLoads[i].Deleted = false
	}
	for i := range data.Switches {
		data.Switches[i].CommitId = i
		data.Switches[i].Deleted = false
	}
	for i := range data.Terminals {
		data.Terminals[i].CommitId = i
		data.Terminals[i].Deleted = false
	}
	for i := range data.Transformer {
		data.Transformer[i].CommitId = i
		data.Transformer[i].Deleted = false
	}

	data.PickOnlyLatest()

	require.Equal(t, origNumConNodes, len(data.ConNodes), "ConNodes")
	require.Equal(t, origNumConformLoads, len(data.ConformLoads), "ConformLoads")
	require.Equal(t, origNumGens, len(data.Gens), "Gens")
	require.Equal(t, origNumLines, len(data.Lines), "Lines")
	require.Equal(t, origNumNonConformLoads, len(data.NonConformLoads), "NonConformLoads")
	require.Equal(t, origNumSwitches, len(data.Switches), "Switches")
	require.Equal(t, origNumTerminals, len(data.Terminals), "Terminals")
	require.Equal(t, origNumTransformer, len(data.Transformer), "Transformer")
}

func TestErrorOnUnknownType(t *testing.T) {
	ctx := context.Background()
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	step := makeEquipmentStep(ctx, db, NamedEquipment("some random equipment"))

	var data InVoltageLevel
	err := step.Run(&data)
	require.Error(t, err)
}
