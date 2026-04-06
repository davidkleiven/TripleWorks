package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSetSelectedSubstation(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/selected?mrid=000&name=componentName&fieldName=field", nil)
	SetSelectedSubstation(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "name=\"field\"")
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
	form.Set("fromSubstation", substations[0].Mrid.String())
	form.Set("toSubstation", substations[1].Mrid.String())
	validValues := form.Encode()

	form.Set("fromSubstation", "0000-0000")
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
