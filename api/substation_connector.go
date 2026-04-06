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
		body        SubstationConnectorBody
		line        models.ACLineSegment
		substations []models.Substation
		terminals   []models.Terminal
		vls         []models.VoltageLevel
	)

	failNo, err := pkg.ReturnOnFirstError(
		func() error {
			return r.ParseForm()
		},
		func() error {
			body.FromSubstation = r.FormValue("fromSubstation")
			body.ToSubstation = r.FormValue("toSubstation")
			modelIdStr := r.FormValue("modelId")
			var ierr error
			body.ModelId, ierr = strconv.Atoi(modelIdStr)
			return ierr
		},
		func() error {
			var ierr error
			line, ierr = s.LineRepo.GetByMrid(ctx, acLineMrid)
			return ierr
		},
		func() error {
			var ierr error
			substations, ierr = s.SubstationRepo.ListByMrids(ctx, slices.Values([]string{body.FromSubstation, body.ToSubstation}))
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
		newItems = append(newItems, result.All(body.ModelId))
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
			Endpoint:       "/substation-list",
			Name:           "from-substation",
			SelectedId:     "selected-from-substation",
			ResultTargetId: "selected-from-substation",
			FieldName:      "fromSubstation",
		},
		ToSelector: components.SearchablePickerParams{
			Endpoint:       "/substation-list",
			Name:           "to-substation",
			SelectedId:     "selected-to-substation",
			ResultTargetId: "selected-to-substation",
			FieldName:      "toSubstation",
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

	listComponent := components.SubstationPickResult(keep, r.URL.Query().Get("target"))
	listComponent.Render(ctx, w)
}

func SetSelectedSubstation(w http.ResponseWriter, r *http.Request) {
	mrid := r.URL.Query().Get("mrid")
	name := r.URL.Query().Get("name")
	fieldName := r.URL.Query().Get("fieldName")

	fmt.Fprintf(w, `
		<input type="hidden" name="%s" value="%s" />
		<span class="tag is-primary is-light">%s</span>
		<span class="is-size-7 has-text-grey-light">%s</span>
	`, fieldName, mrid, name, mrid)
}
