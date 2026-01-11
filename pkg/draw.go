package pkg

import (
	"context"
	"image/color"
	"log/slog"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

type SubstationDiagramData struct {
	VoltageLevels     []models.VoltageLevel
	ConnectivityNodes []models.ConnectivityNode
	Terminals         []models.Terminal
	ACLineSegments    []models.ACLineSegment
	SyncMachines      []models.SynchronousMachine
}

func Substation(data *SubstationDiagramData) *vgsvg.Canvas {
	graph := simple.NewUndirectedGraph()
	nodeByMrid := make(map[uuid.UUID]DiagramNode)

	nodeId := int64(0)
	for _, machine := range data.SyncMachines {
		node := DiagramNode{id: nodeId, Name: machine.ShortName, Type: StructName(machine)}
		nodeByMrid[machine.Mrid] = node
		graph.AddNode(&node)
		nodeId += 1
	}

	for _, con := range data.ConnectivityNodes {
		node := DiagramNode{id: nodeId, Name: con.ShortName, Type: StructName(con)}
		nodeByMrid[con.Mrid] = node
		graph.AddNode(&node)
		nodeId += 1
	}

	for _, line := range data.ACLineSegments {
		node := DiagramNode{id: nodeId, Name: line.ShortName, Type: StructName(line)}
		nodeByMrid[line.Mrid] = node
		graph.AddNode(&node)
		nodeId += 1
	}

	for _, terminal := range data.Terminals {
		node := DiagramNode{id: nodeId, Name: terminal.ShortName, Type: StructName(terminal)}
		nodeByMrid[terminal.Mrid] = node
		graph.AddNode(&node)
		nodeId += 1

		n2 := MustGet(nodeByMrid, terminal.ConductingEquipment.Mrid)
		graph.SetEdge(graph.NewEdge(&node, &n2))

		n3 := MustGet(nodeByMrid, terminal.ConnectivityNode.Mrid)
		graph.SetEdge(graph.NewEdge(&node, &n3))
	}

	eades := layout.EadesR2{Repulsion: 1, Rate: 0.05, Updates: 30, Theta: 0.2}
	positionOptimizer := layout.NewOptimizerR2(graph, eades.Update)
	var n int
	for positionOptimizer.Update() {
		n++
	}
	slog.Info("Finished optimizing point positions", "numIterations", n)

	p := plot.New()
	p.Add(render{positionOptimizer})
	p.HideAxes()
	img := vgsvg.New(4*vg.Inch, 4*vg.Inch)
	dc := draw.New(img)
	p.Draw(dc)
	return img
}

func CollectSubstationData(ctx context.Context, db *bun.DB, s *models.Substation) (SubstationDiagramData, error) {
	var result SubstationDiagramData
	_, overallErr := ReturnOnFirstError(
		func() error {
			return db.NewSelect().Model(&result.VoltageLevels).Where("substation_mrid = ?", s.Mrid).Scan(ctx)
		},
		func() error {
			mrids := make([]uuid.UUID, len(result.VoltageLevels))
			for i, vl := range result.VoltageLevels {
				mrids[i] = vl.Mrid
			}
			return db.NewSelect().
				Model(&result.ConnectivityNodes).
				Relation("ConnectivityNodeContainer").
				Where("connectivity_node_container_mrid IN (?)", bun.In(mrids)).
				Scan(ctx)
		},
		func() error {
			mrids := make([]uuid.UUID, len(result.ConnectivityNodes))
			for i, vl := range result.ConnectivityNodes {
				mrids[i] = vl.Mrid
			}
			err := db.NewSelect().
				Model(&result.Terminals).
				Relation("ConnectivityNode").
				Relation("ConductingEquipment").
				Where("connectivity_node_mrid IN (?)", bun.In(mrids)).
				Scan(ctx)

			for _, terminal := range result.Terminals {
				AssertNotNil(terminal.ConnectivityNode)
				AssertNotNil(terminal.ConductingEquipment)
			}
			return err
		},
		func() error {
			mrids := make([]uuid.UUID, len(result.Terminals))
			for i, terminal := range result.Terminals {
				mrids[i] = terminal.ConductingEquipmentMrid
			}
			err := db.NewSelect().
				Model(&result.ACLineSegments).
				Where("mrid IN (?)", bun.In(mrids)).Scan(ctx)
			return err
		},
		func() error {
			mrids := make([]uuid.UUID, len(result.Terminals))
			for i, terminal := range result.Terminals {
				mrids[i] = terminal.ConductingEquipmentMrid
			}
			err := db.NewSelect().
				Model(&result.SyncMachines).
				Where("mrid IN (?)", bun.In(mrids)).Scan(ctx)
			return err
		},
	)

	// Extract only latest active versions
	result.VoltageLevels = OnlyActiveLatest(result.VoltageLevels)
	result.ConnectivityNodes = OnlyActiveLatest(result.ConnectivityNodes)
	result.Terminals = OnlyActiveLatest(result.Terminals)
	result.ACLineSegments = OnlyActiveLatest(result.ACLineSegments)
	result.SyncMachines = OnlyActiveLatest(result.SyncMachines)

	slog.InfoContext(
		ctx,
		"Loaded data for diagram",
		"substation-mrid", s.Mrid,
		"voltageLevels", len(result.VoltageLevels),
		"ConnectivityNodes", len(result.ConnectivityNodes),
		"Terminals", len(result.Terminals),
		"ACLineSegments", len(result.ACLineSegments),
		"SyncMachines", len(result.SyncMachines),
	)
	return result, overallErr
}

type render struct {
	layout.GraphR2
}

func (r render) Plot(c draw.Canvas, plt *plot.Plot) {
	nodes := r.GraphR2.Nodes()
	if nodes.Len() == 0 {
		slog.Info("Passed graph had zero nodes")
		return
	}

	trX, trY := plt.Transforms(&c)
	lineStyle := draw.LineStyle{Color: color.Black, Width: vg.Points(2)}

	// Draw edges
	for nodes.Next() {
		u := nodes.Node()
		uId := u.ID()
		ur2 := r.GraphR2.LayoutNodeR2(uId)

		to := r.GraphR2.From(uId)
		pointU := vg.Point{
			X: trX(ur2.Coord2.X),
			Y: trY(ur2.Coord2.Y),
		}

		for to.Next() {
			v := to.Node()
			vid := v.ID()
			vr2 := r.GraphR2.LayoutNodeR2(vid)
			pointV := vg.Point{
				X: trX(vr2.Coord2.X),
				Y: trY(vr2.Coord2.Y),
			}
			c.StrokeLine2(lineStyle, pointU.X, pointU.Y, pointV.X, pointV.Y)
		}
	}

	nodes = r.GraphR2.Nodes()
	squareSize := vg.Length(15)
	for nodes.Next() {
		u := nodes.Node()
		uId := u.ID()
		ur2 := r.GraphR2.LayoutNodeR2(uId)
		dNode := mustGetDiagramNode(u)
		pointU := vg.Point{
			X: trX(ur2.Coord2.X),
			Y: trY(ur2.Coord2.Y),
		}
		dNode.Draw(&c, pointU, squareSize)
	}
}
func (r render) DataRange() (xmin, xmax, ymin, ymax float64) {
	nodes := r.GraphR2.Nodes()
	if nodes.Len() == 0 {
		return
	}

	var xys plotter.XYs
	xys = make(plotter.XYs, 0, nodes.Len())
	for nodes.Next() {
		u := nodes.Node()
		uid := u.ID()
		ur2 := r.GraphR2.LayoutNodeR2(uid)
		xys = append(xys, plotter.XY(ur2.Coord2))
	}

	xmin, xmax, ymin, ymax = plotter.XYRange(xys)
	dx := xmax - xmin
	dy := ymax - ymin
	return xmin - 0.05*dx, xmax + 0.05*dx, ymin - 0.05*dy, ymax + 0.05*dy
}

func mustGetDiagramNode(u graph.Node) *DiagramNode {
	n, ok := u.(*DiagramNode)
	if !ok {
		panic("All nodes must be of kind 'DiagramNode'")
	}
	return n
}
