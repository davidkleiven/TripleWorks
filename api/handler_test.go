package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	RootHandler(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestSetup(t *testing.T) {
	mux := http.NewServeMux()
	Setup(mux)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
