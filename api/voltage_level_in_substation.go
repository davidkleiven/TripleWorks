package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
)

type InVoltageLevelEndpoint struct {
	voltageLevelRepo repository.VoltageLevelReadRepository
	timeout          time.Duration
}

func (i *InVoltageLevelEndpoint) VoltageLevelsInSubstation(w http.ResponseWriter, r *http.Request) {
	mrid := r.PathValue("mrid")
	ctx, cancel := context.WithTimeout(r.Context(), i.timeout)
	defer cancel()

	vls, err := i.voltageLevelRepo.InSubstation(ctx, mrid)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch voltage levels", "error", err)
		http.Error(w, "Failed to fetch voltage levels: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := struct {
		Mrid          string                `json:"mrid"`
		VoltageLevels []models.VoltageLevel `json:"voltage_levels"`
	}{
		Mrid:          mrid,
		VoltageLevels: vls,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
