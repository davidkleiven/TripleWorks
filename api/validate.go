package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/integrity"
	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/uptrace/bun"
)

func NewBunValidationEndpoint(db *bun.DB, timeout time.Duration) *ValidateEndpoint {
	return &ValidateEndpoint{
		TerminalRepo: &repository.BunReadRepository[models.Terminal]{Db: db},
		Timeout:      timeout,
	}
}

func NewInMemValidationEndpoint() *ValidateEndpoint {
	return &ValidateEndpoint{
		TerminalRepo: &repository.InMemReadRepository[models.Terminal]{},
	}
}

type ValidateEndpoint struct {
	TerminalRepo repository.ReadRepository[models.Terminal]
	Timeout      time.Duration
}

func (v *ValidateEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var terminals []models.Terminal

	ctx, cancel := context.WithTimeout(r.Context(), v.Timeout)
	defer cancel()

	failNo, err := pkg.ReturnOnFirstError(
		func() error {
			var ierr error
			terminals, ierr = v.TerminalRepo.List(ctx)
			return ierr
		},
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch data", "failNo", failNo, "error", err)
		http.Error(w, fmt.Sprintf("Failed to fetch data (%d): %s", failNo, err), http.StatusInternalServerError)
		return
	}

	checks := []integrity.QualityCheck{
		&integrity.UniqueSequenceNumberPerConductingEquipment{Terminals: terminals},
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	encoder := json.NewEncoder(w)

	for _, check := range checks {
		result := check.Check()
		err := result.Report(encoder)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to write report", "error", err)
			encoder.Encode(JsonlError{Error: err.Error()})
			return
		}
	}
}

type JsonlError struct {
	Error string `json:"error"`
}
