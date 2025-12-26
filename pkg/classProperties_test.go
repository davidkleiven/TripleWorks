package pkg

import (
	"os"
	"os/exec"
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

func TestWriteBunModel(t *testing.T) {
	eqGraph := equipmentRdfsGraph()
	properties := eqGraph.Properties()

	f, err := os.CreateTemp("", "models*.go")
	if err != nil {
		t.Fatal("Could not create temp file")
	}
	defer f.Close()
	//defer os.Remove(f.Name())

	types := eqGraph.GolangTypes()
	properties.WriteAllBunModels(f, *NewTypes(types), "")

	cmd := exec.Command("go", "build", f.Name())
	err = cmd.Run()
	assert.Nil(t, err)
}

func TestCapitalizeFirstEmptyString(t *testing.T) {
	assert.Equal(t, capitalizeFirst(""), "")
}
