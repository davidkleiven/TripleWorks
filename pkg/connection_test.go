package pkg

import (
	"context"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestConnectNoData(t *testing.T) {
	var data ConnectionCtx
	result := FindConnection(&data)
	require.Equal(t, 0, len(result.Vertices))
}

func validConnectionContext() *ConnectionCtx {
	data := ConnectionCtx{
		Terminals:     make([]models.Terminal, 2),
		ConNodes:      make([]models.ConnectivityNode, 2),
		VoltageLevels: make([]models.VoltageLevel, 2),
		Substations:   make([]models.Substation, 2),
	}

	conductingEquipmentMrid := uuid.New()
	for i := range 2 {
		data.Substations[i].Mrid = uuid.New()
		data.VoltageLevels[i].Mrid = uuid.New()
		data.VoltageLevels[i].SubstationMrid = data.Substations[i].Mrid

		data.ConNodes[i].Mrid = uuid.New()
		data.ConNodes[i].ConnectivityNodeContainerMrid = data.VoltageLevels[i].Mrid

		data.Terminals[i].SequenceNumber = i
		data.Terminals[i].Mrid = uuid.New()
		data.Terminals[i].ConnectivityNodeMrid = data.ConNodes[i].Mrid
		data.Terminals[i].ConductingEquipmentMrid = conductingEquipmentMrid
	}
	return &data
}

func TestAcLineConnectsTwoSubstations(t *testing.T) {
	data := validConnectionContext()
	result := FindConnection(data)
	data.Terminals = MustSlice(data.Terminals)
	require.Equal(t, result.Mrid, data.Terminals[0].ConductingEquipmentMrid)
	require.Equal(t, 2, len(result.Vertices))
}

func TestFetchConnectionData(t *testing.T) {
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	ctx := context.Background()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	data := validConnectionContext()

	entities := make([]models.Entity, 0, 8)
	for _, term := range data.Terminals {
		entities = append(entities, MakeEntity(term, 0))
	}
	for _, con := range data.ConNodes {
		entities = append(entities, MakeEntity(con, 0))
	}
	for _, vl := range data.VoltageLevels {
		entities = append(entities, MakeEntity(vl, 0))
	}
	for _, sub := range data.Substations {
		entities = append(entities, MakeEntity(sub, 0))
	}

	_, err = db.NewInsert().Model(&entities).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&data.Terminals).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&data.ConNodes).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&data.VoltageLevels).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&data.Substations).Exec(ctx)
	require.NoError(t, err)

	mrid := MustSlice(data.Terminals)[0].ConductingEquipmentMrid

	receivedCtx, err := FetchConnectionData(ctx, db, mrid.String())
	require.NoError(t, err)
	require.Equal(t, 2, len(receivedCtx.Terminals))
	require.Equal(t, 2, len(receivedCtx.ConNodes))
	require.Equal(t, 2, len(receivedCtx.VoltageLevels))
	require.Equal(t, 2, len(receivedCtx.Substations))

}
