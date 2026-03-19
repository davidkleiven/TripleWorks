package pkg

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/stretchr/testify/require"
)

func TestCreateList(t *testing.T) {
	bvs := []models.BaseVoltage{{}, {}}
	versionedObject := make([]models.VersionedObject, len(bvs))
	for i, bv := range bvs {
		versionedObject[i] = &bv
	}

	var buf bytes.Buffer
	CreateList(&buf, versionedObject)

	content := buf.String()
	require.Contains(t, content, "nominal_voltage")
	require.Contains(t, content, "name")
	require.Equal(t, len(bvs)+1, strings.Count(content, "<tr"))
}

func TestCreateListEmpty(t *testing.T) {
	var versionedObject []models.VersionedObject
	var buf bytes.Buffer

	CreateList(&buf, versionedObject)
	require.Contains(t, buf.String(), "No items")
}

func TestNonJsonObjectsNotInList(t *testing.T) {
	a := struct{ A float64 }{A: 1.0}
	items := []any{"not json str", a}
	var buf bytes.Buffer
	CreateList(&buf, items)
	require.NotContains(t, buf.String(), "json")
}

type unmarshable struct{}

func (u *unmarshable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("can not marshal to json")
}

func TestUnmarshableObjectsNotInList(t *testing.T) {
	items := []any{&unmarshable{}}
	var buf bytes.Buffer
	CreateList(&buf, items)
	require.NotContains(t, buf.String(), "json")
}
