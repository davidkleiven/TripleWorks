package pkg

import (
	"context"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var targetVl = uuid.New()

func populateVoltageLevels() *InVoltageLevelDataSources {
	var (
		voltageLevel    repository.InMemVoltageLevelReadRepository
		conNodes        repository.InMemConnectivityNodeReadRepository
		terminals       repository.InMemTerminalReadRepository
		generators      repository.InMemReadRepository[models.SynchronousMachine]
		lines           repository.InMemReadRepository[models.ACLineSegment]
		switches        repository.InMemReadRepository[models.Switch]
		conformLoads    repository.InMemReadRepository[models.ConformLoad]
		nonConformLoads repository.InMemReadRepository[models.NonConformLoad]
		transformers    repository.InMemReadRepository[models.PowerTransformer]
	)
	terminals.Items = make([]models.Terminal, 10)
	for i := range terminals.Items {
		terminals.Items[i].Mrid = uuid.New()
		terminals.Items[i].ConnectivityNodeMrid = uuid.New()
		terminals.Items[i].ConductingEquipmentMrid = uuid.New()
	}

	voltageLevel.Items = make([]models.VoltageLevel, 2)
	voltageLevel.Items[0].Mrid = targetVl
	voltageLevel.Items[1].Mrid = uuid.New()

	vls := voltageLevel.Items

	conNodes.Items = make([]models.ConnectivityNode, len(terminals.Items))

	for i := range conNodes.Items {
		conNodes.Items[i].Mrid = terminals.Items[i].ConnectivityNodeMrid
		conNodes.Items[i].ConnectivityNodeContainerMrid = vls[i%len(vls)].Mrid
	}

	lines.Items = make([]models.ACLineSegment, 10)
	conformLoads.Items = make([]models.ConformLoad, 10)
	generators.Items = make([]models.SynchronousMachine, 10)
	nonConformLoads.Items = make([]models.NonConformLoad, 10)
	switches.Items = make([]models.Switch, 10)
	transformers.Items = make([]models.PowerTransformer, 10)

	lines.Items[0].Mrid = terminals.Items[0].ConductingEquipmentMrid
	conformLoads.Items[0].Mrid = terminals.Items[1].ConductingEquipmentMrid
	generators.Items[0].Mrid = terminals.Items[2].ConductingEquipmentMrid
	nonConformLoads.Items[0].Mrid = terminals.Items[3].ConductingEquipmentMrid
	switches.Items[0].Mrid = terminals.Items[4].ConductingEquipmentMrid
	transformers.Items[0].Mrid = terminals.Items[5].ConductingEquipmentMrid

	return &InVoltageLevelDataSources{
		VoltageLevel:    &voltageLevel,
		ConNodes:        &conNodes,
		Terminals:       &terminals,
		Generators:      &generators,
		Lines:           &lines,
		Switches:        &switches,
		ConformLoads:    &conformLoads,
		NonConformLoads: &nonConformLoads,
		Transformers:    &transformers,
	}
}

func TestEquipmentInVoltageLevel(t *testing.T) {
	sources := populateVoltageLevels()

	ctx := context.Background()
	result, err := FetchInVoltageLevelData(ctx, sources, targetVl.String())
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
