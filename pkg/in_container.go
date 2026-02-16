package pkg

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"golang.org/x/sync/errgroup"
)

type InVoltageLevel struct {
	ConNodes        []models.ConnectivityNode   `json:"connectivity_nodes"`
	ConformLoads    []models.ConformLoad        `json:"conform_loads"`
	Gens            []models.SynchronousMachine `json:"sync_machines"`
	Lines           []models.ACLineSegment      `json:"lines"`
	NonConformLoads []models.NonConformLoad     `json:"non_conform_loads"`
	Switches        []models.Switch             `json:"switches"`
	Terminals       []models.Terminal           `json:"terminals"`
	Transformer     []models.PowerTransformer   `json:"transformers"`
}

func (invl *InVoltageLevel) PickOnlyLatest() {
	invl.ConNodes = OnlyActiveLatest(invl.ConNodes)
	invl.ConformLoads = OnlyActiveLatest(invl.ConformLoads)
	invl.Gens = OnlyActiveLatest(invl.Gens)
	invl.Lines = OnlyActiveLatest(invl.Lines)
	invl.NonConformLoads = OnlyActiveLatest(invl.NonConformLoads)
	invl.Switches = OnlyActiveLatest(invl.Switches)
	invl.Terminals = OnlyActiveLatest(invl.Terminals)
	invl.Transformer = OnlyActiveLatest(invl.Transformer)
}

type InVoltageLevelDataSources struct {
	VoltageLevel    repository.ReadRepository[models.VoltageLevel]
	ConNodes        repository.ConnectivityNodeReadRepository
	Terminals       repository.TerminalReadRepository
	Generators      repository.ReadRepository[models.SynchronousMachine]
	Lines           repository.ReadRepository[models.ACLineSegment]
	Switches        repository.ReadRepository[models.Switch]
	ConformLoads    repository.ReadRepository[models.ConformLoad]
	NonConformLoads repository.ReadRepository[models.NonConformLoad]
	Transformers    repository.ReadRepository[models.PowerTransformer]
}

func FetchInVoltageLevelData(ctx context.Context, sources *InVoltageLevelDataSources, vlMrid string) (*InVoltageLevel, error) {
	var result InVoltageLevel
	vls, err := sources.VoltageLevel.GetByMrid(ctx, vlMrid)
	if err != nil {
		return &result, fmt.Errorf("Failed to fetch voltage levels: %w", err)
	}

	result.ConNodes, err = sources.ConNodes.InContainer(ctx, vls.Mrid.String())
	if err != nil {
		return &result, fmt.Errorf("Failed to fetch connectivity nodes: %w", err)
	}

	conNodeMridIter := func(yield func(v string) bool) {
		for _, cn := range result.ConNodes {
			if !yield(cn.Mrid.String()) {
				return
			}
		}
	}

	result.Terminals, err = sources.Terminals.WithConnectivityNode(ctx, conNodeMridIter)
	if err != nil {
		return &result, fmt.Errorf("Failed to fetch terminals: %w", err)
	}

	condEquipmentMrid := func(yield func(v string) bool) {
		for _, term := range result.Terminals {
			if !yield(term.ConductingEquipmentMrid.String()) {
				return
			}
		}
	}

	group, grContext := errgroup.WithContext(ctx)
	group.Go(
		func() error {
			var ierr error
			result.Gens, ierr = sources.Generators.ListByMrids(grContext, condEquipmentMrid)
			return ierr
		})
	group.Go(
		func() error {
			var ierr error
			result.Lines, ierr = sources.Lines.ListByMrids(grContext, condEquipmentMrid)
			return ierr
		})
	group.Go(func() error {
		var ierr error
		result.Switches, ierr = sources.Switches.ListByMrids(grContext, condEquipmentMrid)
		return ierr
	})
	group.Go(
		func() error {
			var ierr error
			result.ConformLoads, ierr = sources.ConformLoads.ListByMrids(grContext, condEquipmentMrid)
			return ierr
		})
	group.Go(func() error {
		var ierr error
		result.NonConformLoads, ierr = sources.NonConformLoads.ListByMrids(grContext, condEquipmentMrid)
		return ierr
	})
	group.Go(func() error {
		var ierr error
		result.Transformer, ierr = sources.Transformers.ListByMrids(grContext, condEquipmentMrid)
		return ierr
	})
	return &result, group.Wait()
}
