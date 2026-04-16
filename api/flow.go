package api

import (
	"cmp"
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"slices"
	"strconv"
	"time"

	"com.github/davidkleiven/tripleworks/pkg"
)

type FlowResponse struct {
	Flow map[string]float64 `json:"flow"`
}

type FlowEndpoint struct {
	Ptdf        *pkg.PtdfMatrix
	MaxNumFlows int
	Timeout     time.Duration
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

	flow := f.Ptdf.Flow(production)
	flow = NLargest(flow, f.MaxNumFlows)
	w.Header().Set("Content-Type", "application/json")

	resp := FlowResponse{Flow: flow}
	json.NewEncoder(w).Encode(resp)
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
