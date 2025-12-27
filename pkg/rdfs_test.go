package pkg

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/graph/formats/rdf"
)

func TestLoadEquipmentRdfs(t *testing.T) {
	reader, err := resource.Open("resources/equipment.nq")
	assert.Nil(t, err)

	data, err := LoadObjects(reader)
	assert.Nil(t, err)
	assert.Equal(t, data.Nodes().Len(), 1930)
}

func TestErrorOnNoRdfs(t *testing.T) {
	buf := bytes.NewBufferString("not rdf content")
	_, err := LoadObjects(buf)
	assert.NotNil(t, err)

	assert.Contains(t, err.Error(), "statement no.")
}

func equipmentRdfsGraph() *RdfsGraph {
	reader, err := resource.Open("resources/equipment.nq")
	if err != nil {
		panic(err)
	}
	data, err := LoadObjects(reader)
	if err != nil {
		panic(err)
	}

	return &RdfsGraph{Graph: data}
}

func TestIsClass(t *testing.T) {
	eqGraph := equipmentRdfsGraph()
	stmt := FindFirstMatch(eqGraph.Graph, SubjectEndswith("#Terminal>"))
	assert.NotNil(t, stmt)
	assert.Equal(t, eqGraph.IsClass(stmt.Subject), true)
}

func TestIsClassFalse(t *testing.T) {
	eqGraph := equipmentRdfsGraph()
	stmt := FindFirstMatch(eqGraph.Graph, SubjectEndswith("#Terminal.ConductingEquipment>"))
	assert.NotNil(t, stmt)
	assert.Equal(t, eqGraph.IsClass(stmt.Subject), false)
}

func TestProperties(t *testing.T) {
	eqGraph := equipmentRdfsGraph()
	propertyList := eqGraph.Properties()
	assert.Greater(t, len(propertyList.properties), 0)
	assert.Greater(t, len(propertyList.superclass), 0)
	stmt := FindFirstMatch(eqGraph.Graph, SubjectEndswith("#Terminal>"))
	assert.NotNil(t, stmt)

	assert.Equal(t, eqGraph.IsClass(stmt.Subject), true)

	terminalProps := propertyList.GetProperties(stmt.Subject)
	assert.Equal(t, 9, len(terminalProps))

	// Smoke test some known
	exist := map[string]bool{
		"Terminal.ConductingEquipment": false,
		"IdentifiedObject.name":        false,
	}

	for _, item := range terminalProps {
		for k := range exist {
			if strings.Contains(item.Value, k) {
				exist[k] = true
			}
		}
	}

	// Check that all exist
	for _, v := range exist {
		assert.True(t, v)
	}
}

func TestGolangTypes(t *testing.T) {
	eqGraph := equipmentRdfsGraph()
	dtypes := eqGraph.GolangTypes()

	highVoltageLimit := Must(rdf.NewIRITerm(Cim16 + "VoltageLevel.highVoltageLimit")).Value
	dtype, ok := dtypes[highVoltageLimit]
	assert.True(t, ok)
	assert.Equal(t, "float64", dtype)

	name := Must(rdf.NewIRITerm(Cim16 + "IdentifiedObject.name")).Value
	dtype, ok = dtypes[name]
	assert.True(t, ok)
	assert.Equal(t, "string", dtype)

	normalOpen := Must(rdf.NewIRITerm(Cim16 + "Switch.normalOpen")).Value
	dtype, ok = dtypes[normalOpen]
	assert.True(t, ok)
	assert.Equal(t, "bool", dtype)
}

func TestUnusedAssociations(t *testing.T) {
	eqGraph := equipmentRdfsGraph()
	unused := eqGraph.UnusedAssociations()
	assert.Greater(t, len(unused), 0)
}

func TestGeographicalRegionInSuperclasses(t *testing.T) {
	eqGraph := equipmentRdfsGraph()
	properties := eqGraph.Properties()
	geo := "<" + Cim16 + "GeographicalRegion>"
	exists := false
	for k := range properties.superclass {
		if k.Value == geo {
			exists = true
			break
		}
	}
	require.True(t, exists)
}
