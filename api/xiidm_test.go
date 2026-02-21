package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.github/davidkleiven/tripleworks/repository"
	"github.com/stretchr/testify/require"
)

type FailingBusBreakerRepo struct{}

func (f *FailingBusBreakerRepo) Fetch(ctx context.Context) ([]repository.BusBreakerConnection, error) {
	return nil, errors.New("failed to fetch connections")
}

func TestReceiveXmlData(t *testing.T) {
	endpoint := XiidmExport{
		BusBreakerRepo: &repository.CachedBusbReakerrepo{},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/xiidm", nil)
	endpoint.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/xml", rec.Header().Get("Content-Type"))
}

func TestInternalServerErrorOnFetchFailure(t *testing.T) {
	endpoint := XiidmExport{
		BusBreakerRepo: &FailingBusBreakerRepo{},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/xiidm", nil)
	endpoint.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)

}
