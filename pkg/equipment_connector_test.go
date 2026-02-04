package pkg

import (
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIsConnected(t *testing.T) {
	var (
		cn1 = uuid.New()
		cn2 = uuid.New()
		cn3 = uuid.New()

		equip1 = uuid.New()
		equip2 = uuid.New()
		equip3 = uuid.New()
	)

	terminals := make([]models.Terminal, 4)

	// Line
	terminals[0].ConnectivityNodeMrid = cn1
	terminals[0].ConductingEquipmentMrid = equip1

	// Switch
	terminals[1].ConnectivityNodeMrid = cn1
	terminals[1].ConductingEquipmentMrid = equip2
	terminals[2].ConnectivityNodeMrid = cn2
	terminals[2].ConductingEquipmentMrid = equip2

	// Generator
	terminals[3].ConnectivityNodeMrid = cn2
	terminals[3].ConductingEquipmentMrid = equip3

	t.Run("connected", func(t *testing.T) {
		connector := NewEmptyConnector()
		connector.AddTerminals(terminals...)
		require.True(t, connector.IsConnected(equip1, equip3, 5))
		require.False(t, connector.IsConnected(equip1, equip3, 4))
	})

	terminals[3].ConnectivityNodeMrid = cn3
	t.Run("not connected", func(t *testing.T) {
		connector := NewEmptyConnector()
		connector.AddTerminals(terminals...)
		require.False(t, connector.IsConnected(equip1, equip3, 5))
	})

	t.Run("not connected with unknown node", func(t *testing.T) {
		connector := NewEmptyConnector()
		connector.AddTerminals(terminals...)
		require.False(t, connector.IsConnected(cn1, uuid.UUID{}, 3))
	})
}

func TestGetTerminal(t *testing.T) {
	connector := NewEmptyConnector()
	var terminal models.Terminal
	terminal.ConductingEquipmentMrid = uuid.New()
	terminal.Mrid = uuid.New()
	connector.AddTerminals(terminal)
	require.Equal(t, terminal.Mrid, connector.GetTerminal(terminal.ConductingEquipmentMrid).Mrid)
}

func TestOnlySwitchTerminalsOnExistingTerminals(t *testing.T) {
	connector := NewEmptyConnector()
	terminals := make([]models.Terminal, 2)
	for i := range terminals {
		terminals[i].Mrid = uuid.New()
		terminals[i].ConnectivityNodeMrid = uuid.New()
		terminals[i].ConductingEquipmentMrid = uuid.New()
	}
	connector.AddTerminals(terminals...)

	params := ConnectParams{
		Mrid1: terminals[0].ConductingEquipmentMrid,
		Mrid2: terminals[1].ConductingEquipmentMrid,
	}
	result := connector.MustConnect(&params)
	require.Equal(t, 2, len(result.Terminals))
}
