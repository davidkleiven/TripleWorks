package pkg

import (
	"bytes"
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
	require.Contains(t, content, "NominalVoltage")
	require.Contains(t, content, "Name")
	require.Equal(t, len(bvs)+1, strings.Count(content, "<tr"))
}

func TestCreateListEmpty(t *testing.T) {
	var versionedObject []models.VersionedObject
	var buf bytes.Buffer

	CreateList(&buf, versionedObject)
	require.Contains(t, buf.String(), "No items")
}
