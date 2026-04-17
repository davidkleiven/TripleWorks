package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"strings"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/parquet-go/parquet-go"
	"gonum.org/v1/gonum/mat"
)

type PtdfRecord struct {
	Node string  `parquet:"node"`
	Line string  `parquet:"line"`
	Ptdf float64 `parquet:"ptdf"`
}

type PtdfProvider interface {
	Get(ctx context.Context, node string) map[string]float64
}

type PtdfMatrix struct {
	Data  *mat.Dense
	Lines map[string]int
	Nodes map[string]int
}

func (p *PtdfMatrix) Describe(w io.Writer) {
	var r, c int
	if p.Data != nil {
		r, c = p.Data.Dims()
	}

	maxNum := 3
	var num int
	nodeSamples := make([]string, 0, maxNum)
	for k := range p.Nodes {
		if num >= maxNum {
			break
		}
		nodeSamples = append(nodeSamples, k)
		num++
	}

	num = 0
	lineSamples := make([]string, 0, maxNum)
	for k := range p.Lines {
		if num >= maxNum {
			break
		}
		lineSamples = append(lineSamples, k)
		num++
	}

	fmt.Fprintf(w, "Dimensions: (%d, %d). Nodes: %s. Lines: %s\n", r, c, strings.Join(nodeSamples, ", "), strings.Join(lineSamples, ", "))
}

func (p *PtdfMatrix) InvLineIndex() []string {
	result := make([]string, len(p.Lines))
	for mrid, idx := range p.Lines {
		result[idx] = mrid
	}
	return result
}

func (p *PtdfMatrix) Flow(nodes map[string]float64) map[string]float64 {
	if p.Data == nil {
		return make(map[string]float64)
	}
	colIndices := make([]int, 0, len(nodes))
	colValues := make([]float64, 0, len(nodes))
	result := make(map[string]float64)
	var unknown []string
	for mrid, production := range nodes {
		idx, ok := p.Nodes[mrid]
		if ok {
			colIndices = append(colIndices, idx)
			colValues = append(colValues, production)
		} else {
			unknown = append(unknown, mrid)
		}
	}
	slog.Info("Calculating flow inputs", "num", len(colIndices), "unknown", unknown)

	flows := make([]float64, len(p.Lines))

	for j, production := range colValues {
		matIdx := colIndices[j]
		for i := range flows {
			flows[i] += p.Data.At(i, matIdx) * production
		}
	}

	invLineMap := p.InvLineIndex()
	for i, flow := range flows {
		result[invLineMap[i]] = flow
	}
	return result
}

func NewPtdfMatrix(records []PtdfRecord) *PtdfMatrix {
	lines := make(map[string]int)
	buses := make(map[string]int)
	nextLine := 0
	nextBus := 0
	for _, record := range records {
		lineMrid := RemoveMetadataFromMrid(record.Line)
		if _, ok := lines[lineMrid]; !ok {
			lines[lineMrid] = nextLine
			nextLine++
		}

		nodeMrid := RemoveMetadataFromMrid(record.Node)
		if _, ok := buses[nodeMrid]; !ok {
			buses[nodeMrid] = nextBus
			nextBus++
		}
	}

	numRows := len(lines)
	numCols := len(buses)
	var matrix *mat.Dense
	if numRows > 0 && numCols > 0 {
		matrix = mat.NewDense(numRows, numCols, nil)
		for _, record := range records {
			lineMrid := RemoveMetadataFromMrid(record.Line)
			nodeMrid := RemoveMetadataFromMrid(record.Node)
			row := MustGet(lines, lineMrid)
			col := MustGet(buses, nodeMrid)
			matrix.Set(row, col, record.Ptdf)
		}
	}
	ptdf := PtdfMatrix{
		Lines: lines,
		Nodes: buses,
		Data:  matrix,
	}
	var buf bytes.Buffer
	ptdf.Describe(&buf)
	slog.Info(buf.String())
	return &ptdf
}

// MustCreateRandomPtdf is intended for testing purposes only.
func MustCreateRandomPtdf(lineRepo repository.Lister[models.ACLineSegment], subRepo repository.Lister[models.Substation]) []PtdfRecord {
	ctx := context.Background()
	lines := Must(lineRepo.List(ctx))
	substations := Must(subRepo.List(ctx))
	var result []PtdfRecord
	for _, line := range lines {
		for _, sub := range substations {
			result = append(result, PtdfRecord{
				Node: sub.Mrid.String(),
				Line: line.Mrid.String(),
				Ptdf: 2.0*rand.Float64() - 1.0,
			})
		}
	}
	return result
}

func LoadParquetPtdf(r io.ReaderAt) ([]PtdfRecord, error) {
	parquetReader := parquet.NewGenericReader[PtdfRecord](r)
	var ptdfs []PtdfRecord
	_, err := parquetReader.Read(ptdfs)
	return ptdfs, err
}

func LoadParquetFromFactory(factory LatestReadCloserFactory, bucket string) []PtdfRecord {
	reader, err := factory.MakeReadCloser(context.Background(), bucket)
	if err != nil {
		slog.Error("Could not make ptdf read closer", "error", err)
		return nil
	}
	defer reader.Close()

	ptdf, err := LoadParquetPtdf(reader)
	if err != nil {
		slog.Error("Could not load ptdf: %w", "error", err)
	}
	return ptdf
}

// FilterMrid removes extra information that might have been added to the id
// The id itself should be the 36 first characters
func RemoveMetadataFromMrid(mrid string) string {
	if len(mrid) < 36 {
		return mrid
	}
	return mrid[:36]

}
