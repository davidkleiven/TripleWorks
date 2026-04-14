package pkg

import "context"

type PtdfRecord struct {
	Node string  `parquet:"node"`
	Line string  `parquet:"line"`
	Ptdf float64 `parquet:"float64"`
}

type PtdfProvider interface {
	Get(ctx context.Context, node string) map[string]float64
}

type InMemPtdfProvider struct {
	Items []PtdfRecord
}

func (i *InMemPtdfProvider) Get(ctx context.Context, node string) map[string]float64 {
	result := make(map[string]float64)
	for _, record := range i.Items {
		if record.Node == node {
			result[record.Line] = record.Ptdf
		}
	}
	return result
}
