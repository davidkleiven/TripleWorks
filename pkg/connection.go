package pkg

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ConnectionCtx struct {
	Terminals     []models.Terminal
	ConNodes      []models.ConnectivityNode
	VoltageLevels []models.VoltageLevel
	Substations   []models.Substation
}

func FetchConnectionData(ctx context.Context, db *bun.DB, mrid string) (ConnectionCtx, error) {
	var c ConnectionCtx
	err := Pipe(&c,
		findTerminalsStep(ctx, db, mrid),
		findConNodes(ctx, db),
		findVoltageLevels(ctx, db),
		findSubstations(ctx, db),
	)
	return c, err
}

func FindConnection(c *ConnectionCtx) Connection {
	var con Connection
	c.Terminals = OnlyActiveLatest(c.Terminals)
	c.ConNodes = OnlyActiveLatest(c.ConNodes)
	c.VoltageLevels = OnlyActiveLatest(c.VoltageLevels)
	c.Substations = OnlyActiveLatest(c.Substations)

	cns := IndexBy(c.ConNodes, func(cn models.ConnectivityNode) uuid.UUID { return cn.Mrid })
	vls := IndexBy(c.VoltageLevels, func(vl models.VoltageLevel) uuid.UUID { return vl.Mrid })
	subs := IndexBy(c.Substations, func(sub models.Substation) uuid.UUID { return sub.Mrid })
	for _, term := range c.Terminals {
		cn := cns[term.ConnectivityNodeMrid]
		vl := vls[cn.ConnectivityNodeContainerMrid]
		sub := subs[vl.SubstationMrid]
		path := ConnectionVertex{
			TerminalMrid:     term.Mrid,
			TerminalName:     term.Name,
			TerminalSeqNo:    term.SequenceNumber,
			ConNodeMrid:      cn.Mrid,
			ConNodeName:      cn.Name,
			VoltageLevelMrid: vl.Mrid,
			VoltageLevelName: vl.Name,
			SubstationMrid:   sub.Mrid,
			SubstationName:   sub.Name,
		}
		con.Vertices = append(con.Vertices, path)
		con.Mrid = term.ConductingEquipmentMrid
	}
	return con
}

type ConnectionVertex struct {
	TerminalMrid     uuid.UUID `json:"terminal_mrid"`
	TerminalName     string    `json:"terminal_name"`
	TerminalSeqNo    int       `json:"terminal_sequence_number"`
	ConNodeMrid      uuid.UUID `json:"connectivity_node_mrid,omitempty"`
	ConNodeName      string    `json:"connectivity_node_name,omitempty"`
	VoltageLevelMrid uuid.UUID `json:"voltage_level_mrid,omitempty"`
	VoltageLevelName string    `json:"voltage_level_name,omitempty"`
	SubstationMrid   uuid.UUID `json:"substation_mrid,omitempty"`
	SubstationName   string    `json:"substation_name,omitempty"`
}

type Connection struct {
	Vertices []ConnectionVertex `json:"vertices"`
	Mrid     uuid.UUID          `json:"mrid"`
}

func findTerminalsStep(ctx context.Context, db *bun.DB, mrid string) Step[ConnectionCtx] {
	return Step[ConnectionCtx]{
		Name: "Find terminals",
		Run: func(c *ConnectionCtx) error {
			return db.NewSelect().Model(&c.Terminals).Where("conducting_equipment_mrid = ?", mrid).Scan(ctx)
		},
	}
}

func findConNodes(ctx context.Context, db *bun.DB) Step[ConnectionCtx] {
	return Step[ConnectionCtx]{
		Name: "Find connectivity nodes",
		Run: func(c *ConnectionCtx) error {
			mrids := make([]uuid.UUID, len(c.Terminals))
			for i, term := range c.Terminals {
				mrids[i] = term.ConnectivityNodeMrid
			}
			return db.NewSelect().Model(&c.ConNodes).Where("mrid IN (?)", bun.In(mrids)).Scan(ctx)
		},
	}
}

func findVoltageLevels(ctx context.Context, db *bun.DB) Step[ConnectionCtx] {
	return Step[ConnectionCtx]{
		Name: "Find voltage levels",
		Run: func(c *ConnectionCtx) error {
			mrids := make([]uuid.UUID, len(c.ConNodes))
			for i, con := range c.ConNodes {
				mrids[i] = con.ConnectivityNodeContainerMrid
			}
			return db.NewSelect().Model(&c.VoltageLevels).Where("mrid IN (?)", bun.In(mrids)).Scan(ctx)
		},
	}
}

func findSubstations(ctx context.Context, db *bun.DB) Step[ConnectionCtx] {
	return Step[ConnectionCtx]{
		Name: "Find substations",
		Run: func(c *ConnectionCtx) error {
			mrids := make([]uuid.UUID, len(c.VoltageLevels))
			for i, vl := range c.VoltageLevels {
				mrids[i] = vl.SubstationMrid
			}
			return db.NewSelect().Model(&c.Substations).Where("mrid IN (?)", bun.In(mrids)).Scan(ctx)
		},
	}
}
