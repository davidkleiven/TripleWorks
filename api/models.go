package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/components"
	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
)

type ModelsEndpoint struct {
	Repo    repository.Lister[models.Model]
	Timeout time.Duration
}

func (m *ModelsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), m.Timeout)
	defer cancel()

	models, err := m.Repo.List(ctx)
	if err != nil {
		http.Error(w, "Could not fetch models: "+err.Error(), http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Coult not fetch models", "error", err)
		return
	}

	modelSelector := components.ModelSelector(models)
	modelSelector.Render(ctx, w)
}
