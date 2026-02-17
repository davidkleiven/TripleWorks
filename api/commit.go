package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
)

type CommitEndpoint struct {
	Db      repository.Inserter
	timeout time.Duration
}

func (c *CommitEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	ctx := r.Context()

	var (
		content       []byte
		modelMetaData ModelMetaData
		model         any
	)

	failedCallNum, err := pkg.ReturnOnFirstError(
		func() error {
			var ierr error
			content, ierr = io.ReadAll(r.Body)
			return ierr
		},
		func() error {
			return json.Unmarshal(content, &modelMetaData)
		},
		func() error {
			var ierr error
			model, ierr = pkg.FormInputFieldsForType(modelMetaData.CimType)
			return ierr
		},
		func() error {
			var rawJson map[string]any
			jsonErr := json.Unmarshal(content, &rawJson)
			unsetFields := pkg.UnsetFields(rawJson, model)
			errUnsetCheck := CheckUnsetFieldsCommit(unsetFields)
			return errors.Join(jsonErr, errUnsetCheck)
		},
		func() error {
			return json.Unmarshal(content, model)
		},
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse response", "error", err, "failedCall", failedCallNum)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newHash := pkg.MustGetHash(model)
	if newHash == modelMetaData.Checksum {
		slog.InfoContext(ctx, "No changes detected")
		w.Write([]byte("No changes detected. No commit performed"))
		return
	}

	entity := models.Entity{
		ModelEntity: models.ModelEntity{ModelId: 0},
		Mrid:        modelMetaData.Mrid,
		EntityType:  pkg.StructName(model),
	}

	gridModel := models.Model{
		Id:   modelMetaData.ModelId,
		Name: modelMetaData.ModelName,
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	commit := models.Commit{
		Branch:    "main",
		Message:   modelMetaData.CommitMessage,
		Author:    "TripleWorks",
		CreatedAt: time.Now(),
	}

	itemIter := func(yield func(v any) bool) {
		pkg.YieldMany(yield, &entity, &gridModel, &model)
	}
	err = pkg.InsertAllInserter(ctx, c.Db, commit, itemIter, pkg.NoOpOnInsert)

	if err != nil {
		slog.ErrorContext(ctx, "Could not insert data", "error", err)
		http.Error(w, "Could not insert data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, "Successfully upgraded data", "commitId", commit.Id, "commitMessage", commit.Message, "type", pkg.StructName(model))

	if isDeleted(model) {
		fmt.Fprintf(w, "Item %s was deleted", modelMetaData.Mrid)
		return
	}
	fmt.Fprintf(w, "Successfully updated object %s", modelMetaData.Mrid)
}

func CheckUnsetFieldsCommit(unset []string) error {
	allowedUnset := map[string]struct{}{
		"Id":        {},
		"CommitId":  {},
		"id":        {},
		"commit_id": {},
	}
	for _, k := range unset {
		_, allowed := allowedUnset[k]
		if !allowed {
			return fmt.Errorf("Field '%s' must be set by the provided payload", k)
		}
	}
	return nil
}
