package pkg

import (
	"context"
	"fmt"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

func prepSubstationDb(t *testing.T, db *bun.DB) *models.Substation {
	ctx := context.Background()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	model := models.Model{Name: "model"}
	_, err = db.NewInsert().Model(&model).Exec(ctx)
	require.NoError(t, err)

	commit := models.Commit{}
	_, err = db.NewInsert().Model(&commit).Exec(ctx)
	require.NoError(t, err)

	entities := make([]models.Entity, 7)
	for i := range entities {
		entities[i].ModelId = model.Id
		entities[i].Mrid = uuid.New()
	}

	var (
		vl    models.VoltageLevel
		cn    models.ConnectivityNode
		term  models.Terminal
		term2 models.Terminal
		ac    models.ACLineSegment
		sub   models.Substation
		sync  models.SynchronousMachine
	)

	vl.Mrid = entities[0].Mrid
	vl.SubstationMrid = entities[4].Mrid

	cn.Mrid = entities[1].Mrid
	cn.ConnectivityNodeContainerMrid = vl.Mrid

	ac.Mrid = entities[2].Mrid
	sync.Mrid = entities[5].Mrid

	term.Mrid = entities[3].Mrid
	term.ConductingEquipmentMrid = ac.Mrid
	term.ConnectivityNodeMrid = cn.Mrid

	term2.Mrid = entities[6].Mrid
	term2.ConductingEquipmentMrid = sync.Mrid
	term2.ConnectivityNodeMrid = cn.Mrid

	vl.CommitId = int(commit.Id)
	cn.CommitId = int(commit.Id)
	term.CommitId = int(commit.Id)
	term2.CommitId = int(commit.Id)
	ac.CommitId = int(commit.Id)
	sync.CommitId = int(commit.Id)

	_, err = db.NewInsert().Model(&entities).Exec(ctx)
	require.NoError(t, err)

	data := []any{&vl, &cn, &ac, &term, &term2, &sync}
	for i, item := range data {
		_, err = db.NewInsert().Model(item).Exec(ctx)
		require.NoError(t, err, fmt.Sprintf("Insert: %d", i))
	}

	sub.Mrid = vl.SubstationMrid
	return &sub
}

func TestCollectSubstationData(t *testing.T) {
	db := NewTestConfig().DatabaseConnection()
	substation := prepSubstationDb(t, db)
	result, err := CollectSubstationData(context.Background(), db, substation)
	require.NoError(t, err)
	require.Equal(t, 1, len(result.VoltageLevels), "VoltageLevels")
	require.Equal(t, 1, len(result.ConnectivityNodes), "ConnectivityNodes")
	require.Equal(t, 2, len(result.Terminals), "Terminals")
	require.Equal(t, 1, len(result.ACLineSegments), "ACLineSegments")
	require.Equal(t, 1, len(result.SyncMachines), "SyncMachines")

	require.NotPanics(t, func() { Substation(&result) })
}

func TestPlotEmptyGraph(t *testing.T) {
	canvas := vgsvg.New(vg.Length(4), vg.Length(4))
	dc := draw.New(canvas)
	graph := simple.NewUndirectedGraph()
	eades := layout.EadesR2{}
	optimizer := layout.NewOptimizerR2(graph, eades.Update)
	renderer := render{optimizer}
	require.NotPanics(t, func() { renderer.Plot(dc, plot.New()) })
}

func TestZeroDataRangeOnEmptyGraph(t *testing.T) {
	graph := simple.NewUndirectedGraph()
	eades := layout.EadesR2{}
	optimizer := layout.NewOptimizerR2(graph, eades.Update)
	xmin, xmax, ymin, ymax := render{optimizer}.DataRange()
	require.Equal(t, xmin, 0.0)
	require.Equal(t, xmax, 0.0)
	require.Equal(t, ymin, 0.0)
	require.Equal(t, ymax, 0.0)
}

func TestPanicsOnWrongNodeType(t *testing.T) {
	graph := simple.NewDirectedGraph().NewNode()
	require.Panics(t, func() { mustGetDiagramNode(graph) })
}
