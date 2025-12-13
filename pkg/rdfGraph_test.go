package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/formats/rdf"
)

func TestNilOnNoMatch(t *testing.T) {
	graph := rdf.NewGraph()
	stmt := FindFirstMatch(graph, SubjectEndswith("whatever"))
	assert.Nil(t, stmt)
}
