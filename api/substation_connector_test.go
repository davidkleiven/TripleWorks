package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSetSelectedSubstation(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/selected?mrid=000&name=componentName", nil)
	SetSelectedSubstation(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "name=\"substation-mrid\"")
	require.Contains(t, body, "value=\"000\"")
}

func TestSubstationListQueryHandler(t *testing.T) {
	sub1 := uuid.New()
	substations := make([]models.Substation, 30)

	substations[0].Mrid = sub1
	substations[0].Name = "Sub A"
	substations[0].CommitId = 1

	substations[1].Mrid = sub1
	substations[1].Name = "Sub B"
	substations[1].CommitId = 2

	substations[2].Mrid = uuid.New()
	substations[2].Name = "Other station"
	substations[2].CommitId = 2

	for i := 3; i < len(substations); i++ {
		substations[i].Name = "Oslo"
		substations[i].Mrid = uuid.New()
	}

	lister := repository.InMemLister[models.Substation]{Items: substations}

	listEndpoint := SubstationListQueryHandler{SubstationRepo: &lister, Timeout: time.Second}

	t.Run("return only latest", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/substation-list?q=su", nil)
		rec := httptest.NewRecorder()
		listEndpoint.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		body := rec.Body.String()
		require.Contains(t, body, "Sub B")
		require.NotContains(t, body, "Sub A")
	})

	t.Run("return at most 20", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/substation-list?q=osl", nil)
		rec := httptest.NewRecorder()
		listEndpoint.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		body := rec.Body.String()
		require.Equal(t, strings.Count(body, "<span"), 20, body)

	})

	t.Run("404 on error", func(t *testing.T) {
		defer func() {
			lister.Err = nil
		}()
		lister.Err = errors.New("error")
		req := httptest.NewRequest("GET", "/substation-list?q=osl", nil)
		rec := httptest.NewRecorder()
		listEndpoint.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)

	})
}

func TestSubstationConnectorWorkbench(t *testing.T) {
	items := make([]models.ACLineSegment, 1)
	items[0].Mrid = uuid.New()
	items[0].Name = "Brottem - Klabu"

	lineRepo := repository.InMemReadRepository[models.ACLineSegment]{Items: items}
	wb := SubstationConnectorWorkbench{
		LineRepo: &lineRepo,
		Timeout:  time.Second,
	}

	mux := http.NewServeMux()
	mux.Handle("/wb/{mrid}", &wb)

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/wb/"+items[0].Mrid.String(), nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "id=\"connect-substations-btn\"")
		require.Contains(t, rec.Body.String(), "id=\"move-substations-btn\"")
		require.Contains(t, rec.Body.String(), "/move/"+items[0].Mrid.String())
	})

	t.Run("failure on unknown component", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/wb/0000-0000", nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestSubstationConnector(t *testing.T) {
	lines := make([]models.ACLineSegment, 2)
	lines[0].Mrid = uuid.New()
	lines[1].Mrid = uuid.New()

	terminals := make([]models.Terminal, 1)
	terminals[0].ConductingEquipmentMrid = lines[1].Mrid // Line 1 can not be connected

	substations := make([]models.Substation, 2)
	substations[0].Mrid = uuid.New()
	substations[1].Mrid = uuid.New()

	vls := make([]models.VoltageLevel, 1)
	vls[0].Mrid = uuid.New()
	vls[0].SubstationMrid = substations[0].Mrid

	lineRepo := repository.InMemReadRepository[models.ACLineSegment]{Items: lines}
	substationRepo := repository.InMemReadRepository[models.Substation]{Items: substations}
	terminalRepo := repository.InMemReadRepository[models.Terminal]{Items: terminals}
	vlRepo := repository.InMemReadRepository[models.VoltageLevel]{Items: vls}
	inserter := repository.InMemInserter{}

	connector := SubstationConnector{
		LineRepo:         &lineRepo,
		SubstationRepo:   &substationRepo,
		TerminalRepo:     &terminalRepo,
		VoltageLevelRepo: &vlRepo,
		Inserter:         &inserter,
		Timeout:          time.Second,
	}

	form := url.Values{}
	form.Set("modelId", "1")
	form.Add("substation-mrid", substations[0].Mrid.String())
	form.Add("substation-mrid", substations[1].Mrid.String())
	validValues := form.Encode()

	mrids, ok := form["substation-mrid"]
	require.True(t, ok)
	require.Equal(t, len(mrids), 2)
	mrids[0] = "0000-0000"
	form["substation-mrid"] = mrids
	unknownSubstation := form.Encode()

	form.Set("modelId", "not an int")
	wrongModelId := form.Encode()

	mux := http.NewServeMux()
	mux.Handle("/connect/{mrid}", &connector)

	t.Run("successful connection", func(t *testing.T) {
		defer func() {
			inserter.Items = inserter.Items[:0]
		}()

		req := httptest.NewRequest("POST", "/connect/"+lines[0].Mrid.String(), bytes.NewBufferString(validValues))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Greater(t, len(inserter.Items), 0)
	})

	t.Run("bad request on unknown substation", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/connect/"+lines[0].Mrid.String(), bytes.NewBufferString(unknownSubstation))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, len(inserter.Items))
	})

	t.Run("internal server error on parsing error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/connect/"+lines[0].Mrid.String(), bytes.NewBufferString(wrongModelId))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Equal(t, 0, len(inserter.Items))
	})

	t.Run("conflict if line already has terminals", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/connect/"+lines[1].Mrid.String(), bytes.NewBufferString(validValues))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusConflict, rec.Code)
		require.Equal(t, 0, len(inserter.Items))
	})

	t.Run("insert failure", func(t *testing.T) {
		inserter.InsertError = errors.New("errors")
		defer func() {
			inserter.InsertError = nil
			inserter.Items = inserter.Items[:0]
		}()

		req := httptest.NewRequest("POST", "/connect/"+lines[0].Mrid.String(), bytes.NewBufferString(validValues))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, len(inserter.Items), 1) // Commit is inserted
		require.Contains(t, rec.Body.String(), "Could not connect")
	})

}

func TestMoveLine(t *testing.T) {
	reconnect := func(store *EntityStore) *http.ServeMux {
		mux := http.NewServeMux()
		mux.Handle("/move/{mrid}", &LineReconnector{db: store.db, Timeout: time.Second})
		return mux
	}
	move := func(t *testing.T, mux *http.ServeMux, line uuid.UUID, substations ...string) *httptest.ResponseRecorder {
		form := url.Values{}
		for _, substation := range substations {
			form.Add("substation-mrid", substation)
		}
		req := httptest.NewRequest("POST", "/move/"+line.String(), bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
		return rec
	}

	t.Run("moves both connectivity nodes", func(t *testing.T) {
		store := setupStore(t)
		fx := seedMoveFixture(t, store)

		move(t, reconnect(store), fx.Line, fx.Substations["c"].String(), fx.Substations["d"].String())

		require.ElementsMatch(t,
			[]uuid.UUID{fx.Substations["c"], fx.Substations["d"]},
			lineSubstations(t, store, fx.Line),
		)
		require.Equal(t, seededConNodes+2, countRows(t, store, "connectivity_nodes"))
		require.Equal(t, 2, countRows(t, store, "commits"))
	})

	t.Run("moves only one connectivity node", func(t *testing.T) {
		store := setupStore(t)
		fx := seedMoveFixture(t, store)

		move(t, reconnect(store), fx.Line, fx.Substations["c"].String())

		spans := lineSubstations(t, store, fx.Line)
		require.Len(t, spans, 2)
		require.Contains(t, spans, fx.Substations["c"])
		stayed := slices.DeleteFunc(slices.Clone(spans), func(m uuid.UUID) bool {
			return m == fx.Substations["c"]
		})
		require.Subset(t, []uuid.UUID{fx.Substations["a"], fx.Substations["b"]}, stayed)
		require.Equal(t, seededConNodes+1, countRows(t, store, "connectivity_nodes"))
	})

	t.Run("keeps the node already in a selected substation", func(t *testing.T) {
		store := setupStore(t)
		fx := seedMoveFixture(t, store)

		// Line is at a-b. Selecting a and c must move the b node to c, not
		// the a node to c (which would leave the line at c-b)
		move(t, reconnect(store), fx.Line, fx.Substations["a"].String(), fx.Substations["c"].String())

		require.ElementsMatch(t,
			[]uuid.UUID{fx.Substations["a"], fx.Substations["c"]},
			lineSubstations(t, store, fx.Line),
		)
		require.Equal(t, seededConNodes+1, countRows(t, store, "connectivity_nodes"))
	})

	t.Run("no inserts when moving to the current substations", func(t *testing.T) {
		store := setupStore(t)
		fx := seedMoveFixture(t, store)

		rec := move(t, reconnect(store), fx.Line, fx.Substations["a"].String(), fx.Substations["b"].String())
		require.Contains(t, rec.Body.String(), "Nothing to move")

		require.ElementsMatch(t,
			[]uuid.UUID{fx.Substations["a"], fx.Substations["b"]},
			lineSubstations(t, store, fx.Line),
		)
		require.Equal(t, seededConNodes, countRows(t, store, "connectivity_nodes"))
		require.Equal(t, 1, countRows(t, store, "commits"))
	})

	t.Run("line never connects the same substation", func(t *testing.T) {
		store := setupStore(t)
		fx := seedMoveFixture(t, store)

		substation := fx.Substations["c"].String()
		move(t, reconnect(store), fx.Line, substation, substation)

		spans := lineSubstations(t, store, fx.Line)
		require.Len(t, spans, 2)
		require.Contains(t, spans, fx.Substations["c"])
		require.NotEqual(t, spans[0], spans[1], "Both ends of the line are at the same substation")
		require.Equal(t, seededConNodes+1, countRows(t, store, "connectivity_nodes"))
	})
}

// moveFixture is a line connected between substation "a" and "b" with
// substation "c" and "d" available as move targets
type moveFixture struct {
	Line        uuid.UUID
	Substations map[string]uuid.UUID
}

const seededConNodes = 2

func seedMoveFixture(t *testing.T, store *EntityStore) moveFixture {
	t.Helper()
	ctx := context.Background()

	var (
		bv   models.BaseVoltage
		line models.ACLineSegment
		fx   = moveFixture{Substations: map[string]uuid.UUID{}}
	)

	bv.Mrid = uuid.New()
	bv.NominalVoltage = 132
	line.Mrid = uuid.New()
	line.BaseVoltageMrid = bv.Mrid
	fx.Line = line.Mrid

	items := []any{&bv}
	voltageLevels := map[string]uuid.UUID{}
	for _, name := range []string{"a", "b", "c", "d"} {
		sub := models.Substation{Mrid: uuid.New(), Name: "Substation " + strings.ToUpper(name)}
		vl := models.VoltageLevel{Mrid: uuid.New(), BaseVoltageMrid: bv.Mrid, SubstationMrid: sub.Mrid}

		fx.Substations[name] = sub.Mrid
		voltageLevels[name] = vl.Mrid
		items = append(items, &sub, &vl)
	}

	conNodeA := models.ConnectivityNode{Mrid: uuid.New(), ConnectivityNodeContainerMrid: voltageLevels["a"]}
	conNodeB := models.ConnectivityNode{Mrid: uuid.New(), ConnectivityNodeContainerMrid: voltageLevels["b"]}

	t1 := models.Terminal{
		Mrid:                    uuid.New(),
		SequenceNumber:          1,
		ConductingEquipmentMrid: line.Mrid,
		ConnectivityNodeMrid:    conNodeA.Mrid,
	}
	t2 := models.Terminal{
		Mrid:                    uuid.New(),
		SequenceNumber:          2,
		ConductingEquipmentMrid: line.Mrid,
		ConnectivityNodeMrid:    conNodeB.Mrid,
	}

	items = append(items, &line, &conNodeA, &conNodeB, &t1, &t2)
	require.NoError(t, pkg.InsertAll(ctx, store.db, models.Commit{Message: "Seed move fixture"}, slices.Values(items), pkg.NoOpOnInsert))
	return fx
}

func lineSubstations(t *testing.T, store *EntityStore, line uuid.UUID) []uuid.UUID {
	t.Helper()
	var mrids []uuid.UUID
	err := store.db.NewSelect().
		TableExpr("v_terminals_latest t").
		ColumnExpr("s.mrid as substation_mrid").
		Join("INNER JOIN v_connectivity_nodes_latest c ON t.connectivity_node_mrid = c.mrid").
		Join("INNER JOIN v_voltage_levels_latest v ON c.connectivity_node_container_mrid = v.mrid").
		Join("INNER JOIN v_substations_latest s ON v.substation_mrid = s.mrid").
		Where("t.conducting_equipment_mrid = ?", line).
		Scan(context.Background(), &mrids)
	require.NoError(t, err)
	return mrids
}

func countRows(t *testing.T, store *EntityStore, table string) int {
	t.Helper()
	count, err := store.db.NewSelect().TableExpr(table).Count(context.Background())
	require.NoError(t, err)
	return count
}
