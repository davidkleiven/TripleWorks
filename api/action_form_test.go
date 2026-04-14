package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestActionFormAdd(t *testing.T) {
	form := make(url.Values)
	form.Add("000", "producer")
	form.Add("000", "1.0")

	buf := bytes.NewBufferString(form.Encode())
	req := httptest.NewRequest("POST", "/action-form?mrid=111&name=station", buf)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	action := ActionFormEndpoint{Timeout: time.Second}
	action.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "producer")
	require.Contains(t, body, "000")
	require.Contains(t, body, "station")
	require.Contains(t, body, "111")
}

func TestDeleteItem(t *testing.T) {
	form := make(url.Values)
	form.Add("000", "producer")
	form.Add("000", "1.0")
	form.Add("001", "station")
	form.Add("001", "1.0")

	buf := bytes.NewBufferString(form.Encode())
	req := httptest.NewRequest("POST", "/action-form?mrid=000&action=delete", buf)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	action := ActionFormEndpoint{Timeout: time.Second}
	action.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.NotContains(t, body, "producer")
	require.NotContains(t, body, "000")
	require.Contains(t, body, "station")
	require.Contains(t, body, "001")
}

func TestErrorOnInvalidForm(t *testing.T) {
	buf := bytes.NewBufferString("not;a;form;")
	req := httptest.NewRequest("POST", "/action-form?mrid=000&action=delete", buf)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	action := ActionFormEndpoint{Timeout: time.Second}
	action.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
