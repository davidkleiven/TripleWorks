package pkg

import (
	"context"
	"io"
	"log/slog"
	"math/rand/v2"

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
		if _, ok := lines[record.Line]; !ok {
			lines[record.Line] = nextLine
			nextLine++
		}

		if _, ok := buses[record.Node]; !ok {
			buses[record.Node] = nextBus
			nextBus++
		}
	}

	numRows := len(lines)
	numCols := len(buses)
	var matrix *mat.Dense
	if numRows > 0 && numCols > 0 {
		matrix = mat.NewDense(numRows, numCols, nil)
		for _, record := range records {
			row := MustGet(lines, record.Line)
			col := MustGet(buses, record.Node)
			matrix.Set(row, col, record.Ptdf)
		}
	}
	return &PtdfMatrix{
		Lines: lines,
		Nodes: buses,
		Data:  matrix,
	}
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
