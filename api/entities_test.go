package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupStore(t *testing.T) *EntityStore {
	config := pkg.NewTestConfig()
	store := NewEntityStore(config.DatabaseConnection(), config.Timeout)

	_, err := migrations.RunUp(context.Background(), store.db)
	require.NoError(t, err)
	return store
}

func TestGetEntityForKind(t *testing.T) {
	store := setupStore(t)

	bvs := make([]models.BaseVoltage, 10)
	for i := range bvs {
		mrid, err := uuid.NewUUID()
		require.NoError(t, err)

		bvs[i].Mrid = mrid
		bvs[i].NominalVoltage = 22.0 + float64(10*i)
	}

	chosen := bvs[5].Mrid.String()

	ctx := context.Background()
	_, err := store.db.NewInsert().Model(&bvs).Exec(ctx)
	require.NoError(t, err)

	t.Run("valid request", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/entities?kind=BaseVoltage&choice=%s", chosen), nil)

		store.GetEntityForKind(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		content := rec.Body.String()
		require.Contains(t, content, "<option")
		splitted := strings.Split(content, "\n")
		require.Equal(t, len(bvs), len(splitted)-1)

		// Confirm that the chosen item is the first
		require.Contains(t, splitted[0], chosen)
	})

	t.Run("invalid type", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/entities?kind=MyBaseVoltage", nil)

		store.GetEntityForKind(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	store.timeout = 0
	t.Run("error on timeout", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/entities?kind=BaseVoltage", nil)

		store.GetEntityForKind(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}
func TestCommit(t *testing.T) {
	store := setupStore(t)
	data := `{"mrid": "530dfa65-3158-4bdc-845f-3483a24374b9", "name": "Base voltage 420kV", "cim_type": "BaseVoltage",
	"checksum": "00", "modelId": 0, "modelName": "national", "energy_ident_code_eic": "EIC", "description": "Desc",
	"short_name": "name", "nominal_voltage": 22.0}`

	t.Run("successful commit", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json")
		store.Commit(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "Successfully")
	})

	t.Run("invalid json", func(t *testing.T) {
		body := "not json"
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		store.Commit(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("no insert on repeated checksum", func(t *testing.T) {
		obj := models.IdentifiedObject{}
		checksum := pkg.MustGetHash(obj)
		serialized, err := json.Marshal(obj)
		require.NoError(t, err)

		var generic map[string]any
		err = json.Unmarshal(serialized, &generic)
		require.NoError(t, err)

		generic["checksum"] = checksum
		generic["cim_type"] = "IdentifiedObject"
		generic["modelId"] = 0
		generic["modelName"] = "national"

		body, err := json.Marshal(generic)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		store.Commit(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "No changes detected")
	})

	store.timeout = 0
	t.Run("db insert error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json")
		store.Commit(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "not insert data")
	})
}

func TestGetEnumItems(t *testing.T) {
	store := setupStore(t)
	ctx := context.Background()
	_, err := migrations.RunUp(ctx, store.db)
	require.NoError(t, err)

	t.Run("valid request", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/enum?kind=WindingConnection&choice=4", nil)
		store.GetEnumOptions(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		content := rec.Body.String()
		splitted := strings.Split(content, "\n")
		require.Greater(t, len(splitted), 1)
		require.Contains(t, splitted[0], "value=\"4\"")
	})

	t.Run("no choice", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/enum?kind=WindingConnection", nil)
		store.GetEnumOptions(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "<option")
	})

	t.Run("bad request on unknown type", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/enum?kind=MyWindingConnection", nil)
		store.GetEnumOptions(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

}

func TestNoPanicInGetFinderForSubtype(t *testing.T) {
	formTypes := pkg.FormTypes()
	for name := range formTypes {
		require.NotPanics(t, func() { getFinderForAllSubtypes(name) }, name)
	}
}

func TestLogMsgOnNotEmpty(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{})
	logger := slog.New(handler)
	origLogger := slog.Default()
	slog.SetDefault(logger)
	defer slog.SetDefault(origLogger)

	finderResult := newFinderForSubtype()
	finderResult.notFound = append(finderResult.notFound, "MyType")
	finderResult.LogNotfound(context.Background())

	logContent := buf.String()
	require.Contains(t, logContent, "Could not find")
	require.Contains(t, logContent, "finder")
	require.Contains(t, logContent, "MyType")
}

func TestCheckUnsetFields(t *testing.T) {
	store := setupStore(t)
	err := store.CheckUnsetFields([]string{"reqruiedField"})
	require.Error(t, err)
}

func TestSubstationReturnedWhenRequestingEquipmentContainer(t *testing.T) {
	store := setupStore(t)
	uuid, err := uuid.NewUUID()
	require.NoError(t, err)
	substation := models.Substation{
		EquipmentContainer: models.EquipmentContainer{
			ConnectivityNodeContainer: models.ConnectivityNodeContainer{
				PowerSystemResource: models.PowerSystemResource{
					IdentifiedObject: models.IdentifiedObject{Mrid: uuid, Name: "Demo Substation"},
				},
			},
		},
	}

	ctx := context.Background()
	_, err = store.db.NewInsert().Model(&substation).Exec(ctx)
	require.NoError(t, err)

	finder := pkg.MustGet(pkg.Finders, "Substation")
	result, err := finder(ctx, store.db, 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(result))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entities?kind=EquipmentContainer", nil)
	store.GetEntityForKind(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	content := rec.Body.String()
	require.Contains(t, content, "Demo Substation")
}
