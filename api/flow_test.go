package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/stretchr/testify/require"
)

func formData() *bytes.Buffer {
	form := make(url.Values)
	form.Add("0000", "station A")
	form.Add("0000", "0.5")
	form.Add("1111", "station B")
	form.Add("1111", "0.8")
	return bytes.NewBufferString(form.Encode())
}

func TestFlow(t *testing.T) {

	ptdfRecords := []pkg.PtdfRecord{
		{Node: "0000", Line: "L1", Ptdf: 1.0},
		{Node: "1111", Line: "L1", Ptdf: 0.5},
	}

	flow := FlowEndpoint{
		Ptdf:        pkg.NewPtdfMatrix(ptdfRecords),
		MaxNumFlows: 10,
		Timeout:     time.Second,
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/flow", formData())
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		flow.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var result FlowResponse
		err := json.NewDecoder(rec.Body).Decode(&result)
		require.NoError(t, err)

		want := 0.5*1.0 + 0.8*0.5
		flow, ok := result.Flow["L1"]
		require.True(t, ok)
		require.InDelta(t, want, flow, 1e-6)
	})

	t.Run("bad request when flow is not a number", func(t *testing.T) {
		form := make(url.Values)
		form.Add("0000", "station")
		form.Add("0000", "another name")

		buf := bytes.NewBufferString(form.Encode())
		req := httptest.NewRequest("POST", "/flow", buf)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		flow.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestNthOrEmpty(t *testing.T) {
	data := []string{"A", "B"}
	require.Equal(t, NthOrEmpty(data, 2), "")
	require.Equal(t, NthOrEmpty(data, 1), "B")
	require.Equal(t, NthOrEmpty(data, 0), "A")
}

func TestNthLargest(t *testing.T) {
	data := map[string]float64{
		"a": 5.0,
		"b": 3.0,
		"c": -4.0,
	}

	result := NLargest(data, 2)
	_, aOk := result["a"]
	_, bOk := result["b"]
	_, cOk := result["c"]
	require.True(t, aOk)
	require.False(t, bOk)
	require.True(t, cOk)
}

func TestReceiveNewPtdfOnChannel(t *testing.T) {
	ptdfChannel := make(chan []pkg.PtdfRecord)
	flow := FlowEndpoint{}
	go flow.UpdatePtdf(ptdfChannel)
	defer func() {
		close(ptdfChannel)
	}()

	data := []pkg.PtdfRecord{{Node: "A", Line: "B", Ptdf: 1.0}}
	ptdfChannel <- data

	require.Eventually(t, func() bool {
		flow.PtdfMutex.RLock()
		defer flow.PtdfMutex.RUnlock()
		return flow.Ptdf != nil
	}, time.Second, 10*time.Millisecond)

	require.NotNil(t, flow.Ptdf.Data)
	_, ok := flow.Ptdf.Nodes["A"]
	require.True(t, ok)
	_, ok = flow.Ptdf.Lines["B"]
	require.True(t, ok)
}

func TestCrossBorder(t *testing.T) {
	records := []pkg.PtdfRecord{{Node: "A", Line: "L1"}}
	substationLister := repository.InMemLister[SubstationBidzone]{
		Items: []SubstationBidzone{{Mrid: "A"}},
	}

	crossRegionLineLister := repository.InMemLister[CrossRegionLine]{
		Items: []CrossRegionLine{{LineMrid: "L1"}},
	}

	flow := FlowEndpoint{
		Ptdf: pkg.NewPtdfMatrix(records),
		Timeout: time.Second,
		SubstationBidzoneLister: &substationLister,
		CrossRegionLineLister: &crossRegionLineLister,
	}

	t.Run("success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cross-region", nil)
		flow.CrossRegionPtdf(rec, req)
		require.Equal(t, rec.Code, http.StatusOK)
	})

	t.Run("read failure", func(t *testing.T) {
		defer func() {
			substationLister.Err = nil
		}()
		substationLister.Err = errors.New("something went wrong")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cross-region", nil)
		flow.CrossRegionPtdf(rec, req)
		require.Equal(t, rec.Code, http.StatusInternalServerError)	
	})
}
