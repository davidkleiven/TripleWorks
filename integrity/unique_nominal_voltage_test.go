package integrity

import (
	"bytes"
	"encoding/json"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUniqueNominalVoltage(t *testing.T) {
	bvs := make([]models.BaseVoltage, 4)
	for i := range bvs {
		bvs[i].Mrid = uuid.New()
		bvs[i].NominalVoltage = float64(10 * i)
	}

	uCheck := UniqueNominalVoltage{BaseVoltages: bvs}
	result := uCheck.Check()
	obj, ok := result.(*UniqueNominalVoltageResult)
	require.True(t, ok)
	require.Equal(t, 0, len(obj.BaseVoltages))
}

func TestUniqueNominalVoltageSame(t *testing.T) {
	bvs := make([]models.BaseVoltage, 4)
	for i := range bvs {
		bvs[i].Mrid = uuid.New()
	}

	uCheck := UniqueNominalVoltage{BaseVoltages: bvs}
	result := uCheck.Check()
	obj, ok := result.(*UniqueNominalVoltageResult)
	require.True(t, ok)
	require.Equal(t, 1, len(obj.BaseVoltages))
	require.Equal(t, 4, len(obj.BaseVoltages[0]))
}

func TestUniqueNominalVoltageReport(t *testing.T) {
	rep := UniqueNominalVoltageResult{}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := rep.Report(enc)
	require.NoError(t, err)
}
