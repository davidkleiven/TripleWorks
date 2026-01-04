package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"com.github/davidkleiven/tripleworks/pkg"
)

func AutofillHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)

	var (
		body            []byte
		err             error
		autofillContent pkg.AutofillInput
	)

	failed, err := pkg.ReturnOnFirstError(
		func() error {
			var ierr error
			body, ierr = io.ReadAll(r.Body)
			return ierr
		},
		func() error {
			return json.Unmarshal(body, &autofillContent)
		},
	)

	if err != nil {
		slog.ErrorContext(r.Context(), "Could not read body", "error", err, "call no", failed)
		http.Error(w, "Could not read body "+err.Error(), http.StatusBadRequest)
		return
	}

	result := []pkg.AutofillTargetField{}
	var missing []string
	for _, field := range autofillContent.Fields {
		newChecksum := pkg.MustGetHash(field.Value)
		if newChecksum != field.Checksum {
			// User changed the field, we don't autofill
			continue
		}

		val, err := pkg.GetAutofillValue(field.Label, &autofillContent.State)
		if err == nil {
			field.Checksum = pkg.MustGetHash(val)
			field.Value = val
			result = append(result, field)
		} else {
			missing = append(missing, field.Label)
		}
	}

	if len(missing) > 0 {
		slog.InfoContext(r.Context(), "No autofill", "missing fields", missing)
	}

	wrappedResult := AutofillResult{Data: result}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(wrappedResult)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to encode json", "error", err)
	}
}

type AutofillResult struct {
	Data []pkg.AutofillTargetField `json:"data"`
}
