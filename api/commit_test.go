package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/stretchr/testify/require"
)

func TestCommit(t *testing.T) {
	var inserter repository.InMemInserter
	store := CommitEndpoint{Db: &inserter}
	data := `{"mrid": "530dfa65-3158-4bdc-845f-3483a24374b9", "name": "Base voltage 420kV", "cim_type": "BaseVoltage",
	"checksum": "00", "modelId": 0, "modelName": "national", "energy_ident_code_eic": "EIC", "description": "Desc",
	"short_name": "name", "nominal_voltage": 22.0, "deleted": false}`

	t.Run("successful commit", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json")
		store.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "Successfully")
	})

	t.Run("successful deletion", func(t *testing.T) {
		deleteBody := strings.ReplaceAll(data, "false}", "true}")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(deleteBody))

		req.Header.Set("Content-Type", "application/json")
		store.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "deleted")
	})

	t.Run("invalid json", func(t *testing.T) {
		body := "not json"
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		store.ServeHTTP(rec, req)
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
		store.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "No changes detected")
	})

	inserter.InsertError = errors.New("Something went wrong")
	t.Run("db insert error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/commit", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/json")
		store.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "not insert data")
	})
}
