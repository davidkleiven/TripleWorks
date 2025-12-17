package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogRequestCallPassedHandler(t *testing.T) {
	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}

	wrappedHandler := LogRequest(http.HandlerFunc(handler))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.True(t, handlerCalled)

}
