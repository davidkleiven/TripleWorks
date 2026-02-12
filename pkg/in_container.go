package pkg

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
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

func conNodesStep(ctx context.Context, db *bun.DB, vlMrid string) Step[InVoltageLevel] {
	return Step[InVoltageLevel]{
		Name: "Find connectivity nodes",
		Run: func(c *InVoltageLevel) error {
			return db.NewSelect().Model(&c.ConNodes).Where("connectivity_node_container_mrid = ?", vlMrid).Scan(ctx)
		},
	}
}

func terminalStep(ctx context.Context, db *bun.DB) Step[InVoltageLevel] {
	return Step[InVoltageLevel]{
		Name: "Find terminals",
		Run: func(c *InVoltageLevel) error {
			mrids := make([]uuid.UUID, len(c.ConNodes))
			for i, cn := range c.ConNodes {
				mrids[i] = cn.Mrid
			}
			return db.NewSelect().Model(&c.Terminals).Where("connectivity_node_mrid IN (?)", bun.In(mrids)).Scan(ctx)
		},
	}
}

type NamedEquipment string

const (
	ConformLoadName    NamedEquipment = "ConformLoad"
	GenName            NamedEquipment = "Gen"
	LineName           NamedEquipment = "Line"
	NonConformLoadName NamedEquipment = "NonConformLoad"
	SwitchName         NamedEquipment = "Switch"
	TransformerName    NamedEquipment = "Transformer"
)

func makeEquipmentStep(ctx context.Context, db *bun.DB, namedEquipment NamedEquipment) Step[InVoltageLevel] {
	return Step[InVoltageLevel]{
		Name: fmt.Sprintf("Find %s", namedEquipment),
		Run: func(c *InVoltageLevel) error {
			query := db.NewSelect()
			switch namedEquipment {
			case GenName:
				query = query.Model(&c.Gens)
			case LineName:
				query = query.Model(&c.Lines)
			case TransformerName:
				query = query.Model(&c.Transformer)
			case ConformLoadName:
				query = query.Model(&c.ConformLoads)
			case NonConformLoadName:
				query = query.Model(&c.NonConformLoads)
			case SwitchName:
				query = query.Model(&c.Switches)
			default:
				return fmt.Errorf("Unknown equipment kind: %s", namedEquipment)
			}

			mrids := make([]uuid.UUID, len(c.Terminals))
			for i, term := range c.Terminals {
				mrids[i] = term.ConductingEquipmentMrid
			}

			return query.Where("mrid IN (?)", bun.In(mrids)).Scan(ctx)
		},
	}
}

func FetchInVoltageLevelData(ctx context.Context, db *bun.DB, vlMrid string) (*InVoltageLevel, error) {
	var inVl InVoltageLevel
	err := Pipe(&inVl,
		conNodesStep(ctx, db, vlMrid),
		terminalStep(ctx, db),
		makeEquipmentStep(ctx, db, GenName),
		makeEquipmentStep(ctx, db, LineName),
		makeEquipmentStep(ctx, db, SwitchName),
		makeEquipmentStep(ctx, db, ConformLoadName),
		makeEquipmentStep(ctx, db, NonConformLoadName),
		makeEquipmentStep(ctx, db, TransformerName),
	)
	return &inVl, err
}
