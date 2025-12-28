package pkg

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	dir, err := os.MkdirTemp("", "data-models*")
	require.NoError(t, err)

	mainFile := filepath.Join(dir, "main.go")
	f, err := os.Create(mainFile)
	require.NoError(t, err)
	defer f.Close()
	defer os.RemoveAll(dir)

	cmd := exec.Command("go", "mod", "init", "tempmod")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err)

	types := eqGraph.GolangTypes()
	params := WriteBunModelParams{
		Types:    *NewTypes(types),
		UuidType: "Entity",
		Package:  "models",
	}
	properties.WriteAllBunModels(f, params)

	_, testFile, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(testFile))
	enumFilePath := filepath.Join(rootDir, "models", "rdfs_enum.go")
	enumFile, err := os.Open(enumFilePath)
	require.NoError(t, err)
	defer enumFile.Close()

	enumDest, err := os.Create(filepath.Join(dir, "enums.go"))
	require.NoError(t, err)
	io.Copy(enumDest, enumFile)
	enumDest.Close()

	cmd = exec.Command("go", "get", "github.com/google/uuid")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err)

	cmd = exec.Command("go", "build")
	cmd.Dir = dir
	err = cmd.Run()
	assert.Nil(t, err)
}

func TestCapitalizeFirstEmptyString(t *testing.T) {
	assert.Equal(t, capitalizeFirst(""), "")
}

func TestStringDefaultInTypes(t *testing.T) {
	types := NewTypes(make(map[string]string))
	require.Equal(t, "string", types.Get("what"))
}
