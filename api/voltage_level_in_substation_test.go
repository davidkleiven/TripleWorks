package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/repository"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/require"
)

func TestVoltageLevelsInSubstation(t *testing.T) {
	var store repository.InMemVoltageLevelReadRepository
	err := faker.FakeData(&store.Items, options.WithRandomMapAndSliceMinSize(1))
	require.NoError(t, err)
	targetMrid := store.Items[0].SubstationMrid

	endpoint := InVoltageLevelEndpoint{voltageLevelRepo: &store, timeout: time.Second}

	mux := http.NewServeMux()
	mux.HandleFunc("/substations/{mrid}/voltage-levels", endpoint.VoltageLevelsInSubstation)

	t.Run("existing", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/substations/%s/voltage-levels", targetMrid), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		var res map[string]any
		err := json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)
		receivedVls, ok := res["voltage_levels"]
		require.True(t, ok)
		asSlice, ok := receivedVls.([]any)
		require.True(t, ok)
		require.Equal(t, 1, len(asSlice))

		mrid, ok := res["mrid"]
		require.True(t, ok)
		require.Equal(t, targetMrid.String(), mrid.(string))
	})

	t.Run("non existing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/substations/0000-0000/voltage-levels", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	endpoint.voltageLevelRepo = &repository.InMemVoltageLevelReadRepository{InSubstationErr: errors.New("something went wrong")}
	t.Run("503 on error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/substations/0000-0000/voltage-levels", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
