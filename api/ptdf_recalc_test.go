package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type FixedDataDoer struct {
	StatusCode int
	Err        error
	Data       []byte
}

func (f *FixedDataDoer) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: f.StatusCode,
		Body:       io.NopCloser(bytes.NewBuffer(f.Data)),
		Header:     make(http.Header),
	}, f.Err
}

func TestRecalcPtdfEndpoint(t *testing.T) {
	line1 := uuid.New()
	sub1 := uuid.New()
	sub2 := uuid.New()
	busBranch := repository.CachedBusbReakerrepo{
		Items: []repository.BusBreakerConnection{
			{Mrid: line1, Name: "A-B", SubstationMrid: sub1},
			{Mrid: line1, Name: "A-B", SubstationMrid: sub2},
		},
	}

	writerFactory := pkg.InMemWriterFactory{}

	successFullResp := FixedDataDoer{
		StatusCode: http.StatusOK,
		Data:       []byte("some parquet data"),
	}

	recalcPtdf := RecalcPtdf{
		Model:             &busBranch,
		PtdfEndpoint:      "loadflowservice/ptdf",
		Bucket:            "/ptdf",
		Doer:              &successFullResp,
		Timeout:           time.Second,
		PtdfWriterFactory: &writerFactory,
	}

	clearCreatedWriters := func() {
		writerFactory.CreatedWriters = writerFactory.CreatedWriters[:0]
	}

	t.Run("success", func(t *testing.T) {
		defer clearCreatedWriters()
		req := httptest.NewRequest("POST", "/recalculate/ptdf", nil)
		rec := httptest.NewRecorder()
		recalcPtdf.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, 1, len(writerFactory.CreatedWriters))
		require.Greater(t, len(writerFactory.CreatedWriters[0].Data), 0)
	})

	t.Run("can write to multiple writers", func(t *testing.T) {
		wf2 := pkg.InMemWriterFactory{}
		wf := pkg.MultiWriterFactory{Factories: []pkg.WriterCloserFactory{&writerFactory, &wf2}}
		defer func() {
			clearCreatedWriters()
			recalcPtdf.PtdfWriterFactory = &writerFactory
		}()
		recalcPtdf.PtdfWriterFactory = &wf
		req := httptest.NewRequest("POST", "/recalculate/ptdf", nil)
		rec := httptest.NewRecorder()
		recalcPtdf.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, 1, len(writerFactory.CreatedWriters))
		require.Equal(t, 1, len(wf2.CreatedWriters))
		require.Greater(t, len(writerFactory.CreatedWriters[0].Data), 0)
		require.Equal(t, writerFactory.CreatedWriters[0].Data, wf2.CreatedWriters[0].Data)
	})

	t.Run("internal server error on failing write", func(t *testing.T) {
		defer func() {
			writerFactory.Err = nil
		}()
		writerFactory.Err = errors.New("something went wrong")
		req := httptest.NewRequest("POST", "/recalculate/ptdf", nil)
		rec := httptest.NewRecorder()
		recalcPtdf.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdatePtdfMessage(t *testing.T) {
	err := errors.New("something went wrong")

	var buf bytes.Buffer
	uploadPtdfRespMessage(&buf, err, "what?")
	require.Contains(t, buf.String(), "went wrong")

	buf.Reset()
	uploadPtdfRespMessage(&buf, nil, "what?")
	require.Contains(t, buf.String(), "what?")
}

func TestIsSuccessful(t *testing.T) {
	success, code := isSuccessful(nil)
	require.Equal(t, -1, code)
	require.False(t, success)

	resp := http.Response{StatusCode: http.StatusInternalServerError}
	success, code = isSuccessful(&resp)
	require.Equal(t, http.StatusInternalServerError, code)
	require.False(t, success)

	resp.StatusCode = http.StatusOK
	success, code = isSuccessful(&resp)
	require.Equal(t, http.StatusOK, code)
	require.True(t, success)
}
