package integrity

import (
	"bytes"
	"encoding/json"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFixInconsistentTerminals(t *testing.T) {
	terminals := make([]models.Terminal, 9)
	for i := range 3 {
		terminals[i].Mrid = uuid.New()
		terminals[i].ConductingEquipmentMrid = uuid.New()
		terminals[i].CommitId = i
		terminals[i].SequenceNumber = 1
	}

	// Should be sequence number 2
	for i := range 3 {
		terminals[i+3].Mrid = uuid.New()
		terminals[i+3].ConductingEquipmentMrid = terminals[i].ConductingEquipmentMrid
		terminals[i+3].CommitId = i + 3
		terminals[i+3].SequenceNumber = 1
		if i == 0 {
			// This should make the first terminal valid
			terminals[i+3].SequenceNumber = 2
		}
	}

	// Copies of the first three with later commit
	for i := range 3 {
		terminals[i+6].Mrid = terminals[i].Mrid
		terminals[i+6].ConductingEquipmentMrid = terminals[i].ConductingEquipmentMrid
		terminals[i+6].CommitId = i + 6
		terminals[i+6].SequenceNumber = terminals[i].SequenceNumber
	}

	check := UniqueSequenceNumberPerConductingEquipment{Terminals: terminals}
	result := check.Check()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := result.Report(enc)
	require.NoError(t, err)

	var data struct {
		Terminals map[string][]models.Terminal `json:"terminals"`
	}
	err = json.Unmarshal(buf.Bytes(), &data)
	require.NoError(t, err)
	require.Equal(t, 2, len(data.Terminals))
}
