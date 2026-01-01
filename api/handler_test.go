package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"com.github/davidkleiven/tripleworks/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	RootHandler(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestSetup(t *testing.T) {
	mux := http.NewServeMux()
	config := pkg.NewTestConfig()
	Setup(mux, config)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestCimTypes(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/cim-types", nil)
	CimTypes(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestEntityFormSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entity/form?type=BaseVoltage", nil)
	EntityForm(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestEntityFormUnknownType(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entity/form?type=MyBaseVoltage", nil)
	EntityForm(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPanicOnUnsuccessfulMigration(t *testing.T) {
	config := pkg.NewTestConfig()
	db := config.DatabaseConnection()
	require.Panics(t, func() { mustPerformMigrations(db, 0) })
}
