package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
	"github.com/parquet-go/parquet-go"
	"github.com/stretchr/testify/require"
)

type FixedDataDoer struct {
	StatusCode int
	Err        error
	Data       []byte
}

func (f *FixedDataDoer) Do(req *http.Request) (*http.Response, error) {
	header := make(http.Header)
	header.Set("Content-Type", "application/x-parquet")
	return &http.Response{
		StatusCode: f.StatusCode,
		Body:       io.NopCloser(bytes.NewBuffer(f.Data)),
		Header:     header,
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

	records := []pkg.PtdfRecord{{Node: "a", Line: "b", Ptdf: 1.0}}
	var buf bytes.Buffer
	writer := parquet.NewGenericWriter[pkg.PtdfRecord](&buf)
	_, err := writer.Write(records)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	successFullResp := FixedDataDoer{
		StatusCode: http.StatusOK,
		Data:       buf.Bytes(),
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

func TestSendToChannel(t *testing.T) {
	ptdfChan := make(chan []pkg.PtdfRecord)
	endpoint := RecalcPtdf{}
	data := []pkg.PtdfRecord{{}}
	require.NotPanics(t, func() { endpoint.Send(data) })

	endpoint.PtdfChan = ptdfChan
	var (
		result []pkg.PtdfRecord
		mu     sync.RWMutex
	)
	go func() {
		mu.Lock()
		defer mu.Unlock()
		result = <-ptdfChan
	}()
	endpoint.Send(data)
	require.Eventually(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return result != nil
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, data, result)
}
