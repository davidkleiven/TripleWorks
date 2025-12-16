package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/formats/rdf"
)

func TestNoPropertiesOnUnknownNode(t *testing.T) {
	properties := NewPropertyList()
	bnode := Must(rdf.NewBlankTerm("bnode"))
	assert.Equal(t, len(properties.GetProperties(bnode)), 0)
}

func TestClasses(t *testing.T) {
	properties := NewPropertyList()
	bnode := Must(rdf.NewBlankTerm("bnode"))
	bnodeTarget := Must(rdf.NewBlankTerm("bode1"))
	properties.AddProperty(bnode, bnodeTarget)
	num := 0
	for range properties.Classes() {
		num++
	}
	assert.Equal(t, num, 1)
}
