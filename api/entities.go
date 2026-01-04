package api

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EntityStore struct {
	db           *bun.DB
	timeout      time.Duration
	allowedUnset map[string]struct{}
}

func (e *EntityStore) GetEnumOptions(w http.ResponseWriter, r *http.Request) {
	var (
		kind       = r.URL.Query().Get("kind")
		choice     = r.URL.Query().Get("choice")
		errCode    int
		choiceId   int
		enumFinder pkg.EnumFinder
		enumValues []models.Enum
	)

	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	failed, err := pkg.ReturnOnFirstError(
		func() error {
			if choice == "" {
				choiceId = -1
				return nil
			}
			var ierr error
			choiceId, ierr = strconv.Atoi(choice)
			errCode = http.StatusBadRequest
			return ierr
		},
		func() error {
			var ok bool
			enumFinder, ok = pkg.EnumFinders[kind]
			errCode = http.StatusBadRequest
			if !ok {
				return fmt.Errorf("Could not find enum for '%s'", kind)
			}
			return nil
		},
		func() error {
			var ierr error
			enumValues, ierr = enumFinder(ctx, e.db)
			errCode = http.StatusInternalServerError
			return ierr
		},
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to process enum equest", "error", err, "call no.", failed)
		http.Error(w, err.Error(), errCode)
		return
	}

	// Sort by choice
	slices.SortFunc(enumValues, func(a, b models.Enum) int {
		if a.GetId() == choiceId {
			return -1
		}

		if b.GetId() == choiceId {
			return 1
		}
		return cmp.Compare(a.GetCode(), b.GetCode())
	})

	for _, item := range enumValues {
		fmt.Fprintf(w, "<option value=\"%d\">%s</option>\n", item.GetId(), item.GetCode())
	}
}

func (e *EntityStore) GetEntityForKind(w http.ResponseWriter, r *http.Request) {
	kind := r.URL.Query().Get("kind")
	choice := r.URL.Query().Get("choice")

	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	result, err := getFinderForAllSubtypes(kind)
	if err != nil {

		slog.ErrorContext(ctx, "Failed to locate a finder", "kind", kind)
		http.Error(w, "Failed to locate a finder for "+kind, http.StatusBadRequest)
		return
	}
	result.LogNotfound(ctx)

	var items []models.MridNameGetter
	for _, finder := range result.finders {
		newItems, err := finder(ctx, e.db, 0)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to find all items of type: %v", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, newItems...)
	}

	// Sort result such that the current choice is first, and the remaining are in alphabetic order
	slices.SortFunc(items, func(a, b models.MridNameGetter) int {
		mridA := a.GetMrid().String()
		mridB := b.GetMrid().String()
		if mridA == choice {
			return -1
		}

		if mridB == choice {
			return 1
		}
		return cmp.Compare(a.GetName(), b.GetName())
	})

	if !choiceExists(items, choice) {
		fmt.Fprintf(w, "<option mrid=\"no-mrid\"></option>")
	}

	for _, item := range items {
		fmt.Fprintf(w, "<option mrid=\"%s\">%s</option>\n", item.GetMrid(), item.GetName())
	}
}

func (e *EntityStore) Commit(w http.ResponseWriter, r *http.Request) {
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
			errUnsetCheck := e.CheckUnsetFields(unsetFields)
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

	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	commit := models.Commit{
		Branch:    "main",
		Message:   modelMetaData.CommitMessage,
		Author:    "TripleWorks",
		CreatedAt: time.Now(),
	}

	err = e.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(&commit).
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = tx.NewInsert().
			Model(&entity).
			On("CONFLICT DO NOTHING").
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().
			Model(&gridModel).
			On("CONFLICT DO NOTHING").
			Exec(ctx)

		if err != nil {
			return err
		}

		if err := pkg.SetCommitId(model, int(commit.Id)); err != nil {
			return err
		}

		_, err = tx.NewInsert().
			Model(model).
			Exec(ctx)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "Could not insert data", "error", err)
		http.Error(w, "Could not insert data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, "Successfully upgraded data", "commitId", commit.Id, "commitMessage", commit.Message, "type", pkg.StructName(model))
	w.Write([]byte("Successfully upgraded data"))
}

func (e *EntityStore) CheckUnsetFields(unset []string) error {
	for _, k := range unset {
		_, allowed := e.allowedUnset[k]
		if !allowed {
			return fmt.Errorf("Field '%s' must be set by the provided payload", k)
		}
	}
	return nil
}

func (e *EntityStore) EntityList(w http.ResponseWriter, r *http.Request) {
	entityType := r.URL.Query().Get("type")
	finder, ok := pkg.Finders[entityType]
	if !ok {
		slog.ErrorContext(r.Context(), "Could not locate a finder", "type", entityType)
		http.Error(w, "Could not locate a finder for the provided type", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	items, err := finder(ctx, e.db, 0)
	if err != nil {
		slog.ErrorContext(ctx, "Could not retrieve items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pkg.CreateList(w, items)
}

func NewEntityStore(db *bun.DB, timeout time.Duration) *EntityStore {
	return &EntityStore{
		db:      db,
		timeout: timeout,
		allowedUnset: map[string]struct{}{
			"Id":        {},
			"CommitId":  {},
			"id":        {},
			"commit_id": {},
			"deleted":   {},
		},
	}
}

type ModelMetaData struct {
	CimType       string    `json:"cim_type"`
	Checksum      string    `json:"checksum"`
	Mrid          uuid.UUID `json:"mrid"`
	ModelId       int       `json:"modelId"`
	ModelName     string    `json:"modelName"`
	CommitMessage string    `json:"commitMessage"`
}

type finderForSubtypesResult struct {
	finders  []pkg.Finder
	notFound []string
}

func (f *finderForSubtypesResult) LogNotfound(ctx context.Context) {
	if len(f.notFound) > 0 {
		slog.InfoContext(ctx, "Could not find a finder", "types", f.notFound)
	}
}

func newFinderForSubtype() *finderForSubtypesResult {
	return &finderForSubtypesResult{
		finders:  []pkg.Finder{},
		notFound: []string{},
	}
}

func getFinderForAllSubtypes(kind string) (*finderForSubtypesResult, error) {
	result := newFinderForSubtype()
	current, err := pkg.FormInputFieldsForType(kind)
	if err != nil {
		return result, fmt.Errorf("Finder for subtypes failed: %w", err)
	}
	subtypes := pkg.Subtypes(current)
	subtypes = append(subtypes, current)
	for _, subtype := range subtypes {
		finder, ok := pkg.Finders[pkg.StructName(subtype)]
		if !ok {
			result.notFound = append(result.notFound, pkg.StructName(subtype))
			continue
		}
		result.finders = append(result.finders, finder)
	}
	return result, nil
}

func choiceExists(items []models.MridNameGetter, choice string) bool {
	for _, item := range items {
		if item.GetMrid().String() == choice {
			return true
		}
	}
	return false
}
