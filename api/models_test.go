package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/stretchr/testify/require"
)

func TestModelRepo(t *testing.T) {
	repo := repository.InMemLister[models.Model]{
		Items: []models.Model{{Id: 1, Name: "Model 1"}, {Id: 2, Name: "Model 2"}},
	}

	endpoint := ModelsEndpoint{Repo: &repo, Timeout: time.Second}
	req := httptest.NewRequest("GET", "/models", nil)
	rec := httptest.NewRecorder()
	endpoint.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	body := rec.Body.String()
	require.Contains(t, body, "<option value=\"1\">")

	repo.Err = errors.New("Something went wrong")
	rec2 := httptest.NewRecorder()
	endpoint.ServeHTTP(rec2, req)
	require.Equal(t, http.StatusInternalServerError, rec2.Code)
}
