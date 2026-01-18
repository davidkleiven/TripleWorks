package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupStore(t *testing.T) *EntityStore {
	dbName := fmt.Sprintf("%s_%d", t.Name(), time.Now().UnixNano())
	config := pkg.NewTestConfig(pkg.WithDbName(dbName))
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
	t.Run("no error on timeout", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/entities?kind=BaseVoltage", nil)

		store.GetEntityForKind(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

}
func TestCommit(t *testing.T) {
	store := setupStore(t)
	data := `{"mrid": "530dfa65-3158-4bdc-845f-3483a24374b9", "name": "Base voltage 420kV", "cim_type": "BaseVoltage",
	"checksum": "00", "modelId": 0, "modelName": "national", "energy_ident_code_eic": "EIC", "description": "Desc",
	"short_name": "name", "nominal_voltage": 22.0, "deleted": false}`

	t.Run("successful commit", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json")
		store.Commit(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "Successfully")
	})

	t.Run("successful deletion", func(t *testing.T) {
		deleteBody := strings.ReplaceAll(data, "false}", "true}")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(deleteBody))

		req.Header.Set("Content-Type", "application/json")
		store.Commit(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "deleted")
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

func TestACLineSegmentReturnedWhenRequestingConductingEquipment(t *testing.T) {
	store := setupStore(t)
	acLine := models.ACLineSegment{
		Conductor: models.Conductor{
			ConductingEquipment: models.ConductingEquipment{
				Equipment: models.Equipment{
					PowerSystemResource: models.PowerSystemResource{
						IdentifiedObject: models.IdentifiedObject{Name: "Demo line"},
					},
				},
			},
		},
	}
	ctx := context.Background()
	_, err := store.db.NewInsert().Model(&acLine).Exec(ctx)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entities?kind=ConductingEquipment", nil)
	store.GetEntityForKind(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	content := rec.Body.String()
	require.Contains(t, content, "Demo line")
}

func TestEntityListBadRequestOnUnknownType(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entities?type=NonExistentType", nil)
	store := setupStore(t)
	store.EntityList(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestInternalServerErrorOnTimeoutInEntityList(t *testing.T) {
	store := setupStore(t)
	store.timeout = 0

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entities?type=BaseVoltage", nil)
	store.EntityList(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestBaseVoltageContainsName(t *testing.T) {
	store := setupStore(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entities?type=BaseVoltage", nil)
	store.EntityList(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "No items")
}

func TestDrawDiagram(t *testing.T) {
	store := setupStore(t)

	var substation models.Substation
	substation.Mrid = uuid.New()
	_, err := store.db.NewInsert().Model(&substation).Exec(context.Background())
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /substation/{mrid}/diagram", store.SubstationDiagram)

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/substation/%s/diagram", substation.Mrid), nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "image/svg+xml", rec.Header().Get("Content-Type"))
	})

	t.Run("no substation", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/substation/0000-0000/diagram", nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("fail on write", func(t *testing.T) {
		rec := httptest.NewRecorder()
		respWriter := testutils.FailingResponseWriter{ResponseWriter: rec, WriteErr: errors.New("what is this??")}
		req := httptest.NewRequest("GET", fmt.Sprintf("/substation/%s/diagram", substation.Mrid), nil)
		mux.ServeHTTP(&respWriter, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, string(respWriter.WantToWrite), "write image")

	})
}

func TestEditComponentForm(t *testing.T) {
	store := setupStore(t)

	entity := models.Entity{Mrid: uuid.New(), EntityType: pkg.StructName(models.Substation{})}
	var substation models.Substation
	substation.Mrid = entity.Mrid

	ctx := context.Background()
	_, err := store.db.NewInsert().Model(&entity).Exec(ctx)
	require.NoError(t, err)

	_, err = store.db.NewInsert().Model(&substation).Exec(ctx)
	require.NoError(t, err)

	invalidEntity := models.Entity{Mrid: uuid.New(), EntityType: "UnknownDataType"}
	_, err = store.db.NewInsert().Model(&invalidEntity).Exec(ctx)
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc("/entity-form/{mrid}", store.EditComponentForm)
	mux.HandleFunc("/resource/{mrid}", store.Resource)

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/entity-form/%s", entity.Mrid), nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "hx-get")
	})

	t.Run("success corresonding resource", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/resource/%s", entity.Mrid), nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var content ResourceItem
		err := json.NewDecoder(rec.Body).Decode(&content)
		require.NoError(t, err)
		require.Equal(t, "Substation", content.Type)
	})

	t.Run("non existing mrid", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/entity-form/0000-0000", nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "Failed to find entity")
	})

	t.Run("non existing mrid get resource", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/resource/0000-0000", nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "Failed to fetch resource")
	})

	t.Run("non existing type", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/entity-form/%s", invalidEntity.Mrid), nil)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "for type")
	})
}

func TestExport(t *testing.T) {

	store := setupStore(t)

	var bv models.BaseVoltage
	_, err := store.db.NewInsert().Model(&bv).Exec(context.Background())
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/export", nil)
		store.Export(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/n-triples", rec.Header().Get("Content-Type"))
	})

	t.Run("cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/export", nil)
		store.Export(rec, req.WithContext(cancelledCtx))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "fetch all items")
	})

}

func jsonlEncode[T any](t *testing.T, records ...T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, record := range records {
		err := enc.Encode(record)
		require.NoError(t, err)
	}
	return &buf
}

func TestSimpleUpload(t *testing.T) {
	store := setupStore(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/upload/{kind}", store.SimpleUpload)

	t.Run("substations no commit", func(t *testing.T) {
		rec := httptest.NewRecorder()

		body := jsonlEncode(t, pkg.SubstationLight{Name: "Sub A"}, pkg.SubstationLight{Name: "Sub B"})
		req := httptest.NewRequest("POST", "/upload/substations", body)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/n-triples", rec.Header().Get("Content-Type"))
		content := rec.Body.String()
		require.Contains(t, content, "Substation")
	})

	t.Run("generators no commit", func(t *testing.T) {
		rec := httptest.NewRecorder()

		body := jsonlEncode(t, pkg.GeneratorLight{Substation: "Sub A"}, pkg.GeneratorLight{Substation: "Sub B"})
		req := httptest.NewRequest("POST", "/upload/generators", body)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/n-triples", rec.Header().Get("Content-Type"))
		content := rec.Body.String()
		require.Contains(t, content, "SynchronousMachine")
	})

	t.Run("aclines no commit", func(t *testing.T) {
		rec := httptest.NewRecorder()

		body := jsonlEncode(t, pkg.LineLight{FromSubstation: "Sub A"}, pkg.LineLight{FromSubstation: "Sub B"})
		req := httptest.NewRequest("POST", "/upload/lines", body)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/n-triples", rec.Header().Get("Content-Type"))
		content := rec.Body.String()
		require.Contains(t, content, "ACLineSegment")
	})

	t.Run("loads no commit", func(t *testing.T) {
		rec := httptest.NewRecorder()

		body := jsonlEncode(t, pkg.LoadLight{Substation: "Sub A"}, pkg.LoadLight{Substation: "Sub B"})
		req := httptest.NewRequest("POST", "/upload/loads", body)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/n-triples", rec.Header().Get("Content-Type"))
		content := rec.Body.String()
		require.Contains(t, content, "ConformLoad")
	})

	t.Run("substations do commit", func(t *testing.T) {
		origMrids, err := pkg.ExistingMrids(context.Background(), store.db, 0)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		body := jsonlEncode(t, pkg.SubstationLight{Name: "Sub A"}, pkg.SubstationLight{Name: "Sub B"})
		req := httptest.NewRequest("POST", "/upload/substations?commit=true", body)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/n-triples", rec.Header().Get("Content-Type"))
		content := rec.Body.String()
		require.Contains(t, content, "Substation")

		finalMrids, err := pkg.ExistingMrids(context.Background(), store.db, 0)
		require.NoError(t, err)
		require.Greater(t, len(finalMrids), len(origMrids))
	})

	t.Run("internal server error on db failure", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/upload/substation", nil)
		mux.ServeHTTP(rec, req.WithContext(cancelledCtx))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "existing mrids")
	})

	t.Run("bad request on bad json", func(t *testing.T) {
		buf := bytes.NewBufferString("not jsonl")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/upload/substations", buf)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("bad request on unknown type", func(t *testing.T) {
		buf := bytes.NewBufferString("not jsonl")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/upload/transformers", buf)
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
