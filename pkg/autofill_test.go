package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsServer(t *testing.T) {
	server := JsServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/js/autofill.js", nil)
	server.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "function")
}

func TestGetAutofill(t *testing.T) {
	formState := FormState{
		Kind:   "ACLineSegment",
		Length: 1.0,
		Name:   "Demo line no 1",
	}

	names := []string{}
	for k := range floatAutofillers {
		names = append(names, k)
	}

	for k := range stringAutofillers {
		names = append(names, k)
	}

	for _, name := range names {
		_, err := GetAutofillValue(name, &formState)
		require.NoError(t, err)
	}

	_, err := GetAutofillValue("NonExistentField", &formState)
	require.Error(t, err)
}
