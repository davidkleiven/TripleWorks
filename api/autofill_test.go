package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.github/davidkleiven/tripleworks/pkg"
	"github.com/stretchr/testify/require"
)

func TestBasRequestOnWrongBodyAutofill(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/autofill", bytes.NewBufferString("not json"))
	AutofillHandler(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAutofill(t *testing.T) {
	data := pkg.AutofillInput{
		Fields: []pkg.AutofillTargetField{
			{Id: "field-value", Label: "ShortName", Checksum: pkg.MustGetHash("shrt"), Value: "shrt"},
			{Id: "field-value", Label: "SomeUnknownStuff", Checksum: pkg.MustGetHash("shrt"), Value: "shrt"},
		},
	}

	t.Run("altered value checksum ok", func(t *testing.T) {
		body, err := json.Marshal(data)
		require.NoError(t, err)

		buffer := bytes.NewBuffer(body)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/autofill", buffer)

		AutofillHandler(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res AutofillResult
		err = json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)

		// Since checksum matches we expect the data to be autofilled
		require.Equal(t, 1, len(res.Data))
		require.Equal(t, "", res.Data[0].Value)
	})

	t.Run("user altered value", func(t *testing.T) {
		origValue := data.Fields[0].Value
		defer func() {
			data.Fields[0].Value = origValue
		}()

		data.Fields[0].Value = "some other"

		body, err := json.Marshal(data)
		require.NoError(t, err)

		buffer := bytes.NewBuffer(body)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/autofill", buffer)

		AutofillHandler(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var res AutofillResult
		err = json.NewDecoder(rec.Body).Decode(&res)
		require.NoError(t, err)

		require.Equal(t, 0, len(res.Data))
	})
}
