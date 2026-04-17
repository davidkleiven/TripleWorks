package pkg

import (
	"bytes"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
	"github.com/parquet-go/parquet-go"
	"github.com/stretchr/testify/require"
)

func TestMustCreatePtdf(t *testing.T) {
	acLine := repository.InMemLister[models.ACLineSegment]{
		Items: []models.ACLineSegment{{}},
	}

	substation := repository.InMemLister[models.Substation]{
		Items: []models.Substation{{}},
	}

	ptdf := MustCreateRandomPtdf(&acLine, &substation)
	require.Equal(t, 1, len(ptdf))
	require.GreaterOrEqual(t, ptdf[0].Ptdf, -1.0)
	require.LessOrEqual(t, ptdf[0].Ptdf, 1.0)
}

func TestGetReturnsEmptyMapOnNoData(t *testing.T) {
	var ptdf PtdfMatrix
	result := ptdf.Flow(map[string]float64{"000-000": 1.0})
	require.NotNil(t, result)
	require.Equal(t, 0, len(result))
}

func TestPtdfMatrix(t *testing.T) {
	ptdfs := []PtdfRecord{
		{Node: "A", Line: "L1", Ptdf: 1.0},
		{Node: "A", Line: "L2", Ptdf: 0.5},
	}

	ptdf := NewPtdfMatrix(ptdfs)
	r, c := ptdf.Data.Dims()
	require.Equal(t, 2, r)
	require.Equal(t, 1, c)
	require.Equal(t, 2, len(ptdf.Lines))
	require.Equal(t, 1, len(ptdf.Nodes))

	flow := ptdf.Flow(map[string]float64{"A": 1.0, "B": 0.5})
	f1, ok := flow["L1"]
	require.True(t, ok)
	require.InDelta(t, 1.0, f1, 1e-6)

	f2, ok := flow["L2"]
	require.True(t, ok)
	require.InDelta(t, 0.5, f2, 1e-6)
}

func TestRemoveMetadataFromMrid(t *testing.T) {
	require.Equal(t, "a", RemoveMetadataFromMrid("a"))
	mrid := uuid.New().String()
	require.Equal(t, mrid, RemoveMetadataFromMrid(mrid))

	mridWithMeta := mrid + "_bus23"
	require.Equal(t, mrid, RemoveMetadataFromMrid(mridWithMeta))
}

func TestNoPanicOnDescribeWithDefault(t *testing.T) {
	var ptdf PtdfMatrix
	var buf bytes.Buffer
	require.NotPanics(t, func() { ptdf.Describe(&buf) })
}

func TestWriteReadRoundTrip(t *testing.T) {
	records := []PtdfRecord{{}}
	var buf bytes.Buffer
	writer := parquet.NewGenericWriter[PtdfRecord](&buf)
	_, err := writer.Write(records)
	writer.Close()
	require.NoError(t, err)

	result, err := LoadParquetPtdf(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	require.Equal(t, 1, len(result))
}
