package pkg

import (
	"bytes"
	"strings"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEntityOptions(t *testing.T) {
	var buf bytes.Buffer
	EntityOptions(&buf, "")
	require.Contains(t, buf.String(), "<option")
}

func TestEntityOptionsTargetFirst(t *testing.T) {
	var buf bytes.Buffer
	EntityOptions(&buf, "Substation")
	require.True(t, strings.HasPrefix(buf.String(), "<option>Substation"))
}

func TestFormInputFieldsForType(t *testing.T) {
	_, err := FormInputFieldsForType("unknown type")
	require.Error(t, err)

	model, err := FormInputFieldsForType("BaseVoltage")
	require.NoError(t, err)

	_, isBaseVoltage := model.(*models.BaseVoltage)
	require.True(t, isBaseVoltage)
}

func TestInputFormFields(t *testing.T) {
	var buf bytes.Buffer

	var item models.Terminal
	FormInputFields(&buf, item)
	require.Contains(t, buf.String(), "SequenceNumber")
}

func TestWriteInputItemBool(t *testing.T) {
	var (
		buf  bytes.Buffer
		isOn bool
	)

	params := writeInputConfig{name: "isOn", value: isOn}
	writeInputItem(&buf, &params)
	require.Contains(t, buf.String(), "isOn", "checkbox")
}

func TestFlattenWorksWithPassedPtr(t *testing.T) {
	var bv models.BaseVoltage
	result := FlattenStruct(&bv)
	require.Greater(t, len(result), 0)
}

func TestMustGetQueryKindPanicsOnMissingValue(t *testing.T) {
	fieldmap := make(map[string]formField)
	require.Panics(t, func() { mustGetQueryKind("MyId", fieldmap) })
}

func TestRandomOrCurrentUUid(t *testing.T) {
	uuid, err := uuid.NewUUID()
	require.NoError(t, err)
	require.Equal(t, uuid, randomOrCurrentUuid(uuid))
}
