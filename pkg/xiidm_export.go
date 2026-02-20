package pkg

import (
	"context"
	"log/slog"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"com.github/davidkleiven/tripleworks/xiidm"
	"github.com/google/uuid"
)

type ExportData struct {
	Lines       []models.ACLineSegment
	Terminals   []models.Terminal
	Substations []models.Substation
}

type XiidmResult struct {
	Network       xiidm.Network
	DanglingLines []uuid.UUID
}

func (x *XiidmResult) LogSummary(ctx context.Context) {
	if len(x.DanglingLines) == 0 {
		return
	}
	slog.InfoContext(ctx, "XiidmSummary", "numSkippedLines", len(x.DanglingLines), "skippedLines", x.DanglingLines)
}

func XiidmBusBreakerModel(data []repository.BusBreakerConnection) *XiidmResult {
	pu := PerUnit{Sbase: 100.0}

	nodeNums := make(map[uuid.UUID]int)
	subMap := make(map[uuid.UUID]*xiidm.Substation)
	nextNum := 0
	for _, row := range data {
		if _, ok := nodeNums[row.SubstationMrid]; !ok {
			nodeNums[row.SubstationMrid] = nextNum
			nextNum++
		}
	}

	network := xiidm.Network{
		CaseDateAttr:               time.Now().Format(time.RFC3339),
		Xmlns:                      xiidm.IidmNs,
		IdAttr:                     uuid.New().String(),
		MinimumValidationLevelAttr: "STRICT",
	}
	for subMrid := range nodeNums {
		bus := xiidm.Bus{
			Identifiable: xiidm.Identifiable{
				IdAttr: subMrid.String() + "_bus",
			},
		}

		bbTop := xiidm.BusBreakerTopology{
			Bus: []xiidm.Bus{bus},
		}

		vl := xiidm.VoltageLevel{
			Identifiable: xiidm.Identifiable{
				IdAttr: subMrid.String() + "_vl",
			},
			NominalVAttr:         1.0, // This is a p.u. model
			TopologyKindAttr:     "BUS_BREAKER",
			LowVoltageLimitAttr:  0.9,
			HighVoltageLimitAttr: 1.1,
			BusBreakerTopology:   &bbTop,
		}

		substation := xiidm.Substation{
			Identifiable: xiidm.Identifiable{IdAttr: subMrid.String()},
			VoltageLevel: []xiidm.VoltageLevel{vl},
		}
		network.Substation = append(network.Substation, substation)
		subMap[subMrid] = &substation
	}

	lines := GroupBy(data, func(v repository.BusBreakerConnection) uuid.UUID { return v.Mrid })
	var dangling []uuid.UUID
	for mrid, con := range lines {
		if len(con) != 2 {
			dangling = append(dangling, mrid)
			continue
		}
		v := con[0].NominalVoltage
		r := con[0].R
		x := con[0].X
		sub1Mrid := con[0].SubstationMrid
		sub2Mrid := con[1].SubstationMrid

		node1 := MustGet(nodeNums, sub1Mrid)
		node2 := MustGet(nodeNums, sub2Mrid)
		sub1 := MustGet(subMap, sub1Mrid)
		sub2 := MustGet(subMap, sub2Mrid)

		line := xiidm.Line{
			RAttr: pu.R(r, v),
			XAttr: pu.X(x, v),
			Branch: xiidm.Branch{
				Node1Attr:           node1,
				Node2Attr:           node2,
				VoltageLevelId1Attr: sub1.VoltageLevel[0].IdAttr,
				VoltageLevelId2Attr: sub2.VoltageLevel[0].IdAttr,
				Identifiable:        xiidm.Identifiable{IdAttr: mrid.String(), NameAttr: con[0].Name},
			},
		}
		network.Line = append(network.Line, line)
	}
	return &XiidmResult{
		Network:       network,
		DanglingLines: dangling,
	}
}
