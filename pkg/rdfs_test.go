package pkg

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	stmt := FindFirstMatch(eqGraph.Graph, SubjectEndswith("#Terminal>"))
	assert.NotNil(t, stmt)

	assert.Equal(t, eqGraph.IsClass(stmt.Subject), true)

	terminalProps := eqGraph.Properties(stmt.Subject)
	assert.Equal(t, 11, len(terminalProps))

	// Smoke test some known
	exist := map[string]bool{
		"Terminal.ConductingEquipment": false,
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
