package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"com.github/davidkleiven/tripleworks/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	endpoint := NewInMemValidationEndpoint()

	terminals := make([]models.Terminal, 2)
	terminals[0].Mrid = uuid.New()
	terminals[1].Mrid = uuid.New()
	store := &repository.InMemReadRepository[models.Terminal]{Items: terminals}

	endpoint.TerminalRepo = store

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/validate", nil)
	endpoint.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestErrorOnFailingRead(t *testing.T) {
	endpoint := NewInMemValidationEndpoint()
	endpoint.TerminalRepo = &repository.FailingReadRepo[models.Terminal]{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/validate", nil)
	endpoint.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestNoCrashOnFailingWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := testutils.FailingResponseWriter{
		ResponseWriter: rec,
		WriteErr:       errors.New("something went wrong"),
	}

	endpoint := NewInMemValidationEndpoint()

	req := httptest.NewRequest("GET", "/validate", nil)
	endpoint.ServeHTTP(&writer, req)

	// Expect OK since writing has started
	require.Equal(t, http.StatusOK, rec.Code)
}
