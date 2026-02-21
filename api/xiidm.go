package api

import (
	"context"
	"encoding/xml"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
)

type XiidmExport struct {
	BusBreakerRepo repository.BusBreakerRepo
	Timeout        time.Duration
}

func (x *XiidmExport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), x.Timeout)
	defer cancel()

	data, err := x.BusBreakerRepo.Fetch(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load connections", "error", err)
		http.Error(w, "Failed to load connections: "+err.Error(), http.StatusInternalServerError)
		return
	}
	result := pkg.XiidmBusBreakerModel(data)
	result.LogSummary(ctx)
	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(result.Network)
}
