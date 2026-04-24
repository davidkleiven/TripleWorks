package api

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
)

type CrossRegionLine struct {
	LineMrid    string `bun:"line_mrid"`
	LineName    string `bun:"line_name"`
	FromBidzone string `bun:"from_bidzone"`
	ToBidzone   string `bun:"to_bidzone"`
}

type SubstationBidzone struct {
	Mrid    string `bun:"mrid"`
	Name    string `bun:"name"`
	Bidzone string `bun:"bidzone"`
}

type CrossBorderPtdf struct {
	Mrid           string  `json:"mrid"`
	Name           string  `json:"name"`
	SubstationMrid string  `json:"substation_mrid"`
	SubstationName string  `json:"substation_name"`
	FromBidzone    string  `json:"from_bidzone"`
	ToBidzone      string  `json:"to_bidzone"`
	Ptdf           float64 `json:"ptdf"`
}

type CrossBorderPtdfResp struct {
	Items []CrossBorderPtdf `json:"items"`
}

type FlowResponse struct {
	Flow map[string]float64 `json:"flow"`
}

type FlowEndpoint struct {
	PtdfMutex               sync.RWMutex
	Ptdf                    *pkg.PtdfMatrix
	MaxNumFlows             int
	Timeout                 time.Duration
	CrossRegionLineLister   repository.Lister[CrossRegionLine]
	SubstationBidzoneLister repository.Lister[SubstationBidzone]
}

func (f *FlowEndpoint) UpdatePtdf(newPtdfs chan []pkg.PtdfRecord) {
	slog.Info("Starting update ptdf task")
	for records := range newPtdfs {
		slog.Info("Received new ptdfs")
		ptdf := pkg.NewPtdfMatrix(records)
		f.PtdfMutex.Lock()
		f.Ptdf = ptdf
		f.PtdfMutex.Unlock()
	}
	slog.Info("Stopping update ptdf task")
}

func (f *FlowEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var production map[string]float64
	failNo, err := pkg.ReturnOnFirstError(
		func() error {
			return r.ParseForm()
		},
		func() error {
			production = make(map[string]float64)
			for k, v := range r.Form {
				floatV, ierr := strconv.ParseFloat(NthOrEmpty(v, 1), 64)
				if ierr != nil {
					return ierr
				}
				production[k] = floatV
			}
			return nil
		},
	)

	if err != nil {
		http.Error(w, "Could not parse data", http.StatusBadRequest)
		slog.Error("Could not parse form", "error", err, "failNo", failNo)
		return
	}

	f.PtdfMutex.RLock()
	flow := f.Ptdf.Flow(production)
	f.PtdfMutex.RUnlock()

	flow = NLargest(flow, f.MaxNumFlows)
	w.Header().Set("Content-Type", "application/json")

	resp := FlowResponse{Flow: flow}
	json.NewEncoder(w).Encode(resp)
}

func (f *FlowEndpoint) CrossRegionPtdf(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), f.Timeout)
	defer cancel()

	connections, errCon := f.CrossRegionLineLister.List(ctx)
	substations, errSub := f.SubstationBidzoneLister.List(ctx)

	if err := errors.Join(errCon, errSub); err != nil {
		http.Error(w, "Could not fetch data: "+err.Error(), http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Could fetch connections or substations", "error", err)
		return
	}

	conMrids := make([]string, len(connections))
	for i, con := range connections {
		conMrids[i] = con.LineMrid
	}

	substationToBidzoneMap := make(map[string]*SubstationBidzone)
	for _, s := range substations {
		substationToBidzoneMap[s.Mrid] = &s
	}

	consMap := make(map[string]*CrossRegionLine)
	for _, c := range connections {
		consMap[c.LineMrid] = &c
	}

	ptdfRecords := f.Ptdf.FilterLines(conMrids)
	var result []CrossBorderPtdf
	for record := range ptdfRecords {
		s := pkg.MustGet(substationToBidzoneMap, record.Node)
		c := pkg.MustGet(consMap, record.Line)
		result = append(result, CrossBorderPtdf{
			Mrid:           record.Line,
			Name:           c.LineName,
			SubstationMrid: record.Node,
			SubstationName: s.Name,
			FromBidzone:    c.FromBidzone,
			ToBidzone:      c.ToBidzone,
			Ptdf:           record.Ptdf,
		})
	}

	respBody := CrossBorderPtdfResp{Items: result}
	json.NewEncoder(w).Encode(respBody)
}

func NLargest(flow map[string]float64, n int) map[string]float64 {
	if len(flow) < n {
		return flow
	}

	keys := make([]string, 0, len(flow))
	for k := range flow {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b string) int {
		va := math.Abs(pkg.MustGet(flow, a))
		vb := math.Abs(pkg.MustGet(flow, b))
		return -cmp.Compare(va, vb)
	})
	result := make(map[string]float64)
	for _, k := range keys[:n] {
		result[k] = pkg.MustGet(flow, k)
	}
	return result

}

func NthOrEmpty(v []string, n int) string {
	if len(v) <= n {
		return ""
	}
	return v[n]
}
