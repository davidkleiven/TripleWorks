package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatchFormEndpoint(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/patch", nil)
	PatchForm(rec, req)
	require.Equal(t, rec.Code, http.StatusOK)
}
