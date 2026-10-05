package api

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"com.github/davidkleiven/tripleworks/components"
	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubstationConnectorBody struct {
	ModelId        int    `json:"modelId"`
	FromSubstation string `json:"fromSubstation"`
	ToSubstation   string `json:"toSubstation"`
}

type SubstationConnector struct {
	LineRepo         repository.ReadRepository[models.ACLineSegment]
	SubstationRepo   repository.ReadRepository[models.Substation]
	TerminalRepo     repository.ReadRepository[models.Terminal]
	VoltageLevelRepo repository.ReadRepository[models.VoltageLevel]
	Inserter         repository.Inserter
	Timeout          time.Duration
}

func (s *SubstationConnector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	acLineMrid := r.PathValue("mrid")

	ctx, cancel := context.WithTimeout(r.Context(), s.Timeout)
	defer cancel()

	var (
		modelId             int
		selectedSubstations []string
		line                models.ACLineSegment
		substations         []models.Substation
		terminals           []models.Terminal
		vls                 []models.VoltageLevel
	)

	failNo, err := pkg.ReturnOnFirstError(
		func() error {
			return r.ParseForm()
		},
		func() error {
			selectedSubstations = r.PostForm["substation-mrid"]
			modelIdStr := r.FormValue("modelId")
			var ierr error
			modelId, ierr = strconv.Atoi(modelIdStr)
			return ierr
		},
		func() error {
			var ierr error
			line, ierr = s.LineRepo.GetByMrid(ctx, acLineMrid)
			return ierr
		},
		func() error {
			var ierr error
			substations, ierr = s.SubstationRepo.ListByMrids(ctx, slices.Values(selectedSubstations))
			return ierr
		},
		func() error {
			var ierr error
			terminals, ierr = s.TerminalRepo.List(ctx)
			return ierr
		},
		func() error {
			var ierr error
			vls, ierr = s.VoltageLevelRepo.List(ctx)
			return ierr
		},
	)

	if err != nil {
		http.Error(w, "Could not fetch data", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Could not fetch data", "error", err, "failNo", failNo)
		return
	}

	substations = pkg.OnlyActiveLatest(substations)
	terminals = pkg.OnlyActiveLatest(terminals)
	vls = pkg.OnlyActiveLatest(vls)

	if n := len(substations); n != 2 {
		msg := fmt.Sprintf("Both substations must exist. Found only %d", n)
		http.Error(w, msg, http.StatusBadRequest)
		slog.ErrorContext(ctx, msg)
		return
	}

	hasTerminal := false
	for _, terminal := range terminals {
		if terminal.ConductingEquipmentMrid == line.Mrid {
			hasTerminal = true
			break
		}
	}

	if hasTerminal {
		http.Error(w, "Lines connected via the substation connector can not have terminals", http.StatusConflict)
		slog.Info("Line already have terminals")
		return
	}

	var newItems []iter.Seq[any]
	for _, substation := range substations {
		params := pkg.LineConnectionParams{
			Substation: substation,
			Line:       line,
			Terminals:  terminals,
		}
		for _, vl := range vls {
			if vl.SubstationMrid == substation.Mrid {
				params.VoltageLevels = append(params.VoltageLevels, vl)
			}
		}

		// Error cases should be covered by earlier checks, thus we panic on error here
		result := pkg.Must(pkg.ConnectLineToSubstation(params))
		newItems = append(newItems, result.All(modelId))
	}

	commit := models.Commit{
		Message: fmt.Sprintf("Connect %s to %s and %s", line.Name, substations[0].Name, substations[1].Name),
		Author:  UserFromCtx(r.Context()),
	}

	err = pkg.InsertAllInserter(ctx, s.Inserter, commit, pkg.Chain(newItems...), pkg.NoOpOnInsert)
	if err != nil {
		w.Write([]byte("Could not connect lines to substation: " + err.Error()))
	} else {
		w.Write([]byte("Successfully committed: " + commit.Message))
	}
}

type SubstationConnectorWorkbench struct {
	LineRepo repository.ReadRepository[models.ACLineSegment]
	Timeout  time.Duration
}

func (s *SubstationConnectorWorkbench) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mrid := r.PathValue("mrid")
	ctx, cancel := context.WithTimeout(r.Context(), s.Timeout)
	defer cancel()

	line, err := s.LineRepo.GetByMrid(ctx, mrid)
	if err != nil {
		http.Error(w, "Could not fetch line: "+mrid, http.StatusInternalServerError)
		slog.ErrorContext(r.Context(), "Could not fetch  line", "mrid", mrid, "error", err)
		return
	}
	params := components.SubstationSelectorParams{
		FromSelector: components.SearchablePickerParams{
			Endpoint: "/substation-list",
			Name:     "from-substation",
		},
		ToSelector: components.SearchablePickerParams{
			Endpoint: "/substation-list",
			Name:     "to-substation",
		},
		LineMrid: line.Mrid.String(),
		LineName: line.Name,
	}

	wbComponent := components.SubstationConnectionWorkbench(params)
	wbComponent.Render(ctx, w)
}

type SubstationListQueryHandler struct {
	SubstationRepo repository.Lister[models.Substation]
	Timeout        time.Duration
}

func (s *SubstationListQueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	display := r.URL.Query().Get("selection-display")

	ctx, cancel := context.WithTimeout(r.Context(), s.Timeout)
	defer cancel()

	substations, err := s.SubstationRepo.List(ctx)
	if err != nil {
		http.Error(w, "Could not fetch substatoons: "+err.Error(), http.StatusInternalServerError)
		slog.Error("Could not fetch substations", "query", query, "error", err)
		return
	}
	substations = pkg.OnlyActiveLatest(substations)

	var keep []models.Substation
	for _, substation := range substations {
		name := strings.ToLower(substation.Name)
		if strings.Contains(name, query) {
			keep = append(keep, substation)
			if len(keep) >= 20 {
				break
			}
		}
	}

	listComponent := components.SubstationPickResult(keep, "#"+display)
	listComponent.Render(ctx, w)
}

func SetSelectedSubstation(w http.ResponseWriter, r *http.Request) {
	mrid := r.URL.Query().Get("mrid")
	name := r.URL.Query().Get("name")

	fmt.Fprintf(w, `<span class="tag is-primary is-light">%s</span><span class="is-size-7 has-text-grey-light">%s</span><input name="substation-mrid" type="hidden" value="%s"/>`, name, mrid, mrid)
}

type ConnectivityNodeToMove struct {
	Mrid             uuid.UUID `bun:"mrid"`
	VoltageLevelMrid uuid.UUID `bun:"vl_mrid"`
	SubstationMrid   uuid.UUID `bun:"substation_mrid"`
	CurrentSubMrid   uuid.UUID `bun:"current_substation_mrid"`
}

type ConNodeMoveCtx struct {
	ToMove   []ConnectivityNodeToMove
	Original []models.ConnectivityNode
}

type LineReconnector struct {
	db      *bun.DB
	Timeout time.Duration
}

func (l *LineReconnector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), l.Timeout)
	defer cancel()
	lineMrid := r.PathValue("mrid")

	conNodesToMove := pkg.Result[bool]{Err: r.ParseForm()}.
		Map(func(bool) []string {
			return r.PostForm["substation-mrid"]
		}).
		Apply(func(substations []string) ([]ConnectivityNodeToMove, error) {
			if len(substations) == 0 {
				return nil, nil
			}
			var conNodesToMove []ConnectivityNodeToMove
			err := l.db.NewSelect().
				TableExpr("v_ac_line_segments_latest as lines").
				ColumnExpr("c.mrid, v_new.mrid as vl_mrid, v_new.substation_mrid, v.substation_mrid as current_substation_mrid").
				Join("INNER JOIN v_terminals_latest t ON t.conducting_equipment_mrid = lines.mrid").
				Join("INNER JOIN v_connectivity_nodes_latest c ON t.connectivity_node_mrid = c.mrid").
				Join("INNER JOIN v_voltage_levels_latest v ON c.connectivity_node_container_mrid = v.mrid").
				Join("INNER JOIN v_substations_latest s ON v.substation_mrid = s.mrid").
				Join("INNER JOIN v_voltage_levels_latest v_new ON v_new.base_voltage_mrid = lines.base_voltage_mrid").
				Where("lines.mrid = ?", lineMrid).
				Where("v_new.substation_mrid IN (?)", bun.List(substations)).
				Scan(ctx, &conNodesToMove)
			return conNodesToMove, err
		}).
		Apply(func(toMove []ConnectivityNodeToMove) (ConNodeMoveCtx, error) {
			if len(toMove) == 0 {
				return ConNodeMoveCtx{}, nil
			}
			mrids := make([]string, 0, len(toMove))
			seen := make(map[string]struct{})
			for _, item := range toMove {
				if _, ok := seen[item.Mrid.String()]; !ok {
					mrids = append(mrids, item.Mrid.String())
					seen[item.Mrid.String()] = struct{}{}
				}
			}

			var conNodes []models.ConnectivityNode
			err := l.db.NewSelect().TableExpr("v_connectivity_nodes_latest").Where("mrid IN (?)", bun.List(mrids)).Scan(ctx, &conNodes)
			return ConNodeMoveCtx{
				ToMove:   toMove,
				Original: conNodes,
			}, err

		})

	if conNodesToMove.Err != nil {
		slog.Error("Could not find connectivity nodes to move", "error", conNodesToMove.Err)
		http.Error(w, "Could not move connectivity nodes: "+conNodesToMove.Err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove cases where candidate substation is among one of current substations
	currentSubstations := make(map[uuid.UUID]struct{})
	for _, toMove := range conNodesToMove.Value.ToMove {
		currentSubstations[toMove.CurrentSubMrid] = struct{}{}
	}

	byMrids := pkg.IndexBy(conNodesToMove.Value.Original, func(c models.ConnectivityNode) uuid.UUID { return c.Mrid })

	movedNodes := make(map[uuid.UUID]struct{})
	newConNodes := make([]models.ConnectivityNode, 0, len(byMrids))
	for _, move := range conNodesToMove.Value.ToMove {
		if _, ok := currentSubstations[move.SubstationMrid]; ok {
			// Substation is already one of the ends of the line. A line should
			// never connect the same substation twice
			continue
		}
		if _, ok := movedNodes[move.Mrid]; ok {
			// Connectivity node already has a new home in this request
			continue
		}

		node, ok := byMrids[move.Mrid]
		pkg.Assert(ok, "Connectivity node must be present")
		node.Id = 0
		node.ConnectivityNodeContainerMrid = move.VoltageLevelMrid
		currentSubstations[move.SubstationMrid] = struct{}{}
		movedNodes[move.Mrid] = struct{}{}
		newConNodes = append(newConNodes, node)
	}

	if len(newConNodes) == 0 {
		w.Write([]byte("Nothing to move"))
		return
	}

	commit := models.Commit{
		Message: fmt.Sprintf("Move %d connectivity nodes to new substations", len(newConNodes)),
		Author:  UserFromCtx(ctx),
	}
	items := func(yield func(v any) bool) {
		for i := range newConNodes {
			if !yield(&newConNodes[i]) {
				return
			}
		}
	}
	if err := pkg.InsertAll(ctx, l.db, commit, items, pkg.NoOpOnInsert); err != nil {
		slog.ErrorContext(ctx, "Could not move connectivity nodes", "error", err, "num", len(newConNodes))
		http.Error(w, "Could not move connectivity nodes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, "Updated voltage levels for connectivity nodes", "num", len(newConNodes))
	fmt.Fprintf(w, "Successfully moved %d connectivity nodes", len(newConNodes))
}
