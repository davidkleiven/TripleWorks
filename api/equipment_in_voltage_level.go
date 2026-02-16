package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
)

type EquipmentInVoltageLevelOpt func(e *EquipmentInVoltageLevelEndpoint)

func WithTerminalRepo(getter repository.TerminalReadRepository) EquipmentInVoltageLevelOpt {
	return func(e *EquipmentInVoltageLevelEndpoint) {
		e.Terminals = getter
	}
}

func WithVoltageLevelGetter(getter repository.ReadRepository[models.VoltageLevel]) EquipmentInVoltageLevelOpt {
	return func(e *EquipmentInVoltageLevelEndpoint) {
		e.VoltageLevel = getter
	}
}

func NewInMemEquipmentInVoltageLevel(opts ...EquipmentInVoltageLevelOpt) *EquipmentInVoltageLevelEndpoint {
	d := EquipmentInVoltageLevelEndpoint{
		VoltageLevel:    &repository.InMemReadRepository[models.VoltageLevel]{},
		ConNodes:        &repository.InMemConnectivityNodeReadRepository{},
		Terminals:       &repository.InMemTerminalReadRepository{},
		Generators:      &repository.InMemReadRepository[models.SynchronousMachine]{},
		Lines:           &repository.InMemReadRepository[models.ACLineSegment]{},
		Switches:        &repository.InMemReadRepository[models.Switch]{},
		ConformLoads:    &repository.InMemReadRepository[models.ConformLoad]{},
		NonConformLoads: &repository.InMemReadRepository[models.NonConformLoad]{},
		Transformers:    &repository.InMemReadRepository[models.PowerTransformer]{},
	}
	for _, opt := range opts {
		opt(&d)
	}
	return &d
}

type EquipmentInVoltageLevelEndpoint struct {
	VoltageLevel    repository.ReadRepository[models.VoltageLevel]
	ConNodes        repository.ConnectivityNodeReadRepository
	Terminals       repository.TerminalReadRepository
	Generators      repository.ReadRepository[models.SynchronousMachine]
	Lines           repository.ReadRepository[models.ACLineSegment]
	Switches        repository.ReadRepository[models.Switch]
	ConformLoads    repository.ReadRepository[models.ConformLoad]
	NonConformLoads repository.ReadRepository[models.NonConformLoad]
	Transformers    repository.ReadRepository[models.PowerTransformer]
	timeout         time.Duration
}

func (e *EquipmentInVoltageLevelEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mrid := r.PathValue("mrid")
	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	sources := pkg.InVoltageLevelDataSources{
		VoltageLevel:    e.VoltageLevel,
		ConNodes:        e.ConNodes,
		Terminals:       e.Terminals,
		Generators:      e.Generators,
		Lines:           e.Lines,
		Switches:        e.Switches,
		ConformLoads:    e.ConformLoads,
		NonConformLoads: e.NonConformLoads,
		Transformers:    e.Transformers,
	}
	data, err := pkg.FetchInVoltageLevelData(ctx, &sources, mrid)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to find resources in voltage level", "error", err)
		http.Error(w, "Failed to find resources in voltage level: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data.PickOnlyLatest()
	respData := InVoltageLevelResp{
		Mrid:      mrid,
		Resources: data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&respData)
}

type InVoltageLevelResp struct {
	Mrid      string              `json:"mrid"`
	Resources *pkg.InVoltageLevel `json:"resources"`
}
