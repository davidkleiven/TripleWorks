package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/stretchr/testify/require"
)

func TestInVoltageLevel(t *testing.T) {
	vl := repository.InMemReadRepository[models.VoltageLevel]{Items: []models.VoltageLevel{{}}}
	mrid := vl.Items[0].Mrid
	store := NewInMemEquipmentInVoltageLevel(WithVoltageLevelGetter(&vl))
	mux := http.NewServeMux()
	mux.Handle("/resources/voltage-level/{mrid}", store)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/resources/voltage-level/%s", mrid), nil)
	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

}

func Test503OnFetchError(t *testing.T) {
	store := NewInMemEquipmentInVoltageLevel(WithTerminalRepo(&repository.InMemTerminalReadRepository{WithConNodeErr: errors.New("something went wrong")}))
	mux := http.NewServeMux()
	mux.Handle("/resources/voltage-level/{mrid}", store)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/resources/voltage-level/0000-0000", nil)
	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
