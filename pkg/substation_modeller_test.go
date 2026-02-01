package pkg

import (
	"context"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNoVoltageLevels(t *testing.T) {
	data := SubstationData{VoltageLevels: []ConnectableVoltageLevel{}}
	result := CreateFullyConnectedSubstation(data, NewEmptyConnector())
	require.Equal(t, 0, len(result.Transformers))
	require.Equal(t, 0, len(result.TransformerEnds))
	require.Equal(t, 0, len(result.Terminals))
	require.Equal(t, 0, len(result.ConnectivityNodes))
	require.Equal(t, 0, len(result.Switches))
	require.Equal(t, 0, len(result.BusNameMarkers))
}

func TestSingleVoltageLevels(t *testing.T) {
	data := SubstationData{
		VoltageLevels: []ConnectableVoltageLevel{
			{
				ConnectivityNodes: []models.ConnectivityNode{{}, {}, {}, {}},
			}},
	}
	result := CreateFullyConnectedSubstation(data, NewEmptyConnector())
	require.Equal(t, 0, len(result.Transformers))
	require.Equal(t, 0, len(result.TransformerEnds))
	require.Equal(t, 0, len(result.Terminals))
	require.Equal(t, 0, len(result.ConnectivityNodes))
	require.Equal(t, 0, len(result.Switches))
	require.Equal(t, 0, len(result.BusNameMarkers))
}

func requireNoReconnection(t *testing.T, data SubstationData, connector *EquipmentConnector) {
	origNumTerminals := len(connector.terminals)
	result := CreateFullyConnectedSubstation(data, connector)
	require.Equal(t, 0, len(result.Transformers))
	require.Equal(t, 0, len(result.TransformerEnds))
	require.Equal(t, 0, len(result.Terminals))
	require.Equal(t, 0, len(result.ConnectivityNodes))
	require.Equal(t, 0, len(result.Switches))
	require.Equal(t, 0, len(result.BusNameMarkers))
	require.Equal(t, origNumTerminals, len(connector.terminals))
}

func TestOneCnInEachLevel(t *testing.T) {
	var (
		bv1 models.BaseVoltage
		bv2 models.BaseVoltage
		cn1 models.ConnectivityNode
		cn2 models.ConnectivityNode
		sub models.Substation
		vl1 models.VoltageLevel
		vl2 models.VoltageLevel
	)

	bv1.Mrid = uuid.New()
	bv2.Mrid = uuid.New()
	cn1.Mrid = uuid.New()
	cn2.Mrid = uuid.New()
	sub.Mrid = uuid.New()
	vl1.BaseVoltageMrid = bv1.Mrid
	vl2.BaseVoltageMrid = bv2.Mrid

	data := SubstationData{
		Substation: sub,
		VoltageLevels: []ConnectableVoltageLevel{
			{
				BaseVoltage:       bv1,
				VoltageLevel:      vl1,
				ConnectivityNodes: []models.ConnectivityNode{cn1},
			},
			{
				BaseVoltage:       bv2,
				VoltageLevel:      vl2,
				ConnectivityNodes: []models.ConnectivityNode{cn2},
			},
		},
	}

	connector := NewEmptyConnector()
	result := CreateFullyConnectedSubstation(data, connector)
	require.Equal(t, 1, len(result.Switches))
	require.Equal(t, 1, len(result.ConnectivityNodes))

	// Two windings + one switch
	require.Equal(t, 4, len(result.Terminals))
	require.Equal(t, 2, len(result.TransformerEnds))
	require.Equal(t, 1, len(result.Transformers))

	// Switch should be located at vl1
	require.Equal(t, result.Switches[0].BaseVoltageMrid, vl1.BaseVoltageMrid)
	require.Equal(t, vl1.BaseVoltageMrid, result.Transformers[0].BaseVoltageMrid)

	require.Equal(t, result.TransformerEnds[0].BaseVoltageMrid, vl1.BaseVoltageMrid)
	require.Equal(t, result.TransformerEnds[1].BaseVoltageMrid, vl2.BaseVoltageMrid)
	requireNoReconnection(t, data, connector)
}

func TestSubstationDataWrite(t *testing.T) {
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

	var sub models.Substation
	sub.Mrid = uuid.New()

	substationEntity := models.Entity{Mrid: sub.Mrid, ModelEntity: models.ModelEntity{ModelId: model.Id}, CommitId: int(commit.Id)}
	_, err = db.NewInsert().Model(&substationEntity).Exec(ctx)
	require.NoError(t, err)

	vls := make([]ConnectableVoltageLevel, 3)
	conNodeEntities := []models.Entity{}
	for i := range vls {
		vls[i].ConnectivityNodes = make([]models.ConnectivityNode, 3)
		for j := range vls[i].ConnectivityNodes {
			mrid := uuid.New()
			vls[i].ConnectivityNodes[j].Mrid = mrid
			conNodeEntity := models.Entity{Mrid: mrid, ModelEntity: models.ModelEntity{ModelId: model.Id}, CommitId: int(commit.Id)}
			conNodeEntities = append(conNodeEntities, conNodeEntity)
		}
	}

	_, err = db.NewInsert().Model(&conNodeEntities).Exec(ctx)
	require.NoError(t, err)

	data := SubstationData{
		Substation:    sub,
		VoltageLevels: vls,
	}

	connector := NewEmptyConnector()
	result := CreateFullyConnectedSubstation(data, connector)
	err = result.Write(ctx, db, model.Id, "Add substation")
	require.NoError(t, err)
	requireNoReconnection(t, data, connector)
}

func TestPanicOnInconsistentMrids(t *testing.T) {
	var (
		vl models.VoltageLevel
		bv models.BaseVoltage
	)

	vl.BaseVoltageMrid = uuid.New()
	bv.Mrid = uuid.New()

	connectableVl := ConnectableVoltageLevel{
		VoltageLevel: vl,
		BaseVoltage:  bv,
	}
	require.Panics(t, func() { connectableVl.RequireConsistentVoltageMrid() })
}

func TestConnectLinesToSubstation(t *testing.T) {
	var (
		substation models.Substation
		line       models.ACLineSegment
		terminal1  models.Terminal     // Connected to line
		terminal2  models.Terminal     // No connected to line
		vl         models.VoltageLevel // Same as line
	)
	substation.Mrid = uuid.New()
	line.Mrid = uuid.New()
	line.BaseVoltageMrid = uuid.New()

	vl.Mrid = uuid.New()
	vl.SubstationMrid = substation.Mrid
	vl.BaseVoltageMrid = line.BaseVoltageMrid

	terminal1.Mrid = uuid.New()
	terminal1.SequenceNumber = 2
	terminal1.ConductingEquipmentMrid = line.Mrid

	terminal2.Mrid = uuid.New()
	terminal2.SequenceNumber = 1
	terminal2.ConductingEquipmentMrid = uuid.New()

	params := LineConnectionParams{
		Substation:    substation,
		Line:          line,
		Terminals:     []models.Terminal{terminal1, terminal2},
		VoltageLevels: []models.VoltageLevel{vl},
	}

	t.Run("connect existing voltage level", func(t *testing.T) {
		result, err := ConnectLineToSubstation(params)
		require.NoError(t, err)
		require.Nil(t, result.VoltageLevel) // Voltage lever already exist
		require.Equal(t, 1, result.Terminal.SequenceNumber)
	})

	t.Run("connect vl does not exist", func(t *testing.T) {
		origBvMrid := params.Line.BaseVoltageMrid
		defer func() {
			params.Line.BaseVoltageMrid = origBvMrid
		}()
		params.Line.BaseVoltageMrid = uuid.New()

		result, err := ConnectLineToSubstation(params)
		require.NoError(t, err)
		require.NotNil(t, result.VoltageLevel) // Voltage does not exist
	})

	t.Run("no connection if line has two terminals", func(t *testing.T) {
		origConductingEquipmentMrid := params.Terminals[1].ConductingEquipmentMrid
		defer func() {
			params.Terminals[1].ConductingEquipmentMrid = origConductingEquipmentMrid
		}()

		params.Terminals[1].ConductingEquipmentMrid = line.Mrid
		_, err := ConnectLineToSubstation(params)
		require.ErrorContains(t, err, "already has two terminals")
	})

	t.Run("no connection if a voltage level belongs to another substation", func(t *testing.T) {
		defer func() {
			params.VoltageLevels[0].SubstationMrid = params.Substation.Mrid
		}()
		params.VoltageLevels[0].SubstationMrid = uuid.New()
		_, err := ConnectLineToSubstation(params)
		require.ErrorContains(t, err, "does not belong to substation")
	})
}

func TestAllLineConnectionResult(t *testing.T) {
	params := LineConnectionResult{}
	count := 0
	for range params.All(0) {
		count++
	}
	require.Equal(t, 8, count)

	params.VoltageLevel = &models.VoltageLevel{}
	count = 0
	for range params.All(0) {
		count++
	}
	require.Equal(t, 10, count)
}
