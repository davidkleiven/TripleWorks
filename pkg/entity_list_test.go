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
	mridGetter := make([]models.MridNameGetter, len(bvs))
	for i, bv := range bvs {
		mridGetter[i] = &bv
	}

	var buf bytes.Buffer
	CreateList(&buf, mridGetter)

	content := buf.String()
	require.Contains(t, content, "NominalVoltage")
	require.Contains(t, content, "Name")
	require.Equal(t, len(bvs)+1, strings.Count(content, "<tr>"))
}

func TestCreateListEmpty(t *testing.T) {
	var mridGetter []models.MridNameGetter
	var buf bytes.Buffer

	CreateList(&buf, mridGetter)
	require.Contains(t, buf.String(), "No items")
}
