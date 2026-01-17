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
		slog.ErrorContext(ctx, "Failed to locate a finder", "kind", kind, "error", err)
		http.Error(w, "Failed to locate a finder for "+kind, http.StatusBadRequest)
		return
	}
	result.LogNotfound(ctx)

	var items []models.VersionedObject
	for _, finder := range result.finders {
		newItems, err := finder(ctx, e.db, 0)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to find all items of type: %v", "error", err)
		}
		items = append(items, newItems...)
	}
	items = pkg.OnlyLatestVersion(items)

	// Sort result such that the current choice is first, and the remaining are in alphabetic order
	slices.SortFunc(items, func(a, b models.VersionedObject) int {
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
		_, dberr := pkg.ReturnOnFirstError(
			func() error {
				_, ierr := tx.NewInsert().
					Model(&commit).
					Exec(ctx)
				return ierr
			},
			func() error {
				_, ierr := tx.NewInsert().
					Model(&entity).
					On("CONFLICT DO NOTHING").
					Exec(ctx)
				return ierr
			},
			func() error {
				_, ierr := tx.NewInsert().
					Model(&gridModel).
					On("CONFLICT DO NOTHING").
					Exec(ctx)
				return ierr
			},
			func() error {
				return pkg.SetCommitId(model, int(commit.Id))
			},
			func() error {
				_, ierr := tx.NewInsert().
					Model(model).
					Exec(ctx)
				return ierr

			},
		)
		return dberr
	})

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
	nameFilter := r.URL.Query().Get("name-filter")
	typeFilter := r.URL.Query().Get("type-filter")
	finder, err := pkg.GetFinder(entityType, nameFilter, typeFilter)
	if err != nil {
		slog.ErrorContext(r.Context(), "Could not locate a finder", "error", err, "type", entityType)
		http.Error(w, "Could not locate a finder for the provided type "+err.Error(), http.StatusBadRequest)
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

func (e *EntityStore) SubstationDiagram(w http.ResponseWriter, r *http.Request) {
	substationMrid := r.PathValue("mrid")

	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	var (
		substation models.Substation
		data       pkg.SubstationDiagramData
	)
	failedNo, err := pkg.ReturnOnFirstError(
		func() error {
			return e.db.NewSelect().Model(&substation).Where("mrid = ?", substationMrid).OrderBy("CommitId", bun.OrderDesc).Limit(1).Scan(ctx)
		},
		func() error {
			var ierr error
			data, ierr = pkg.CollectSubstationData(ctx, e.db, &substation)
			return ierr
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get data from database", "error", err, "failedNo", failedNo)
		http.Error(w, "Failed to get substation from database "+err.Error(), http.StatusInternalServerError)
		return
	}

	image := pkg.Substation(&data)
	w.Header().Set("Content-Type", "image/svg+xml")
	_, err = image.WriteTo(w)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write image", "error", err)
		http.Error(w, "Failed to write image "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (e *EntityStore) GetResource(ctx context.Context, mrid string) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	var entity models.Entity
	err := e.db.NewSelect().Model(&entity).Where("mrid = ?", mrid).Limit(1).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to find entity: %w", err)
	}

	resource, ok := pkg.FormTypes()[entity.EntityType]
	if !ok {
		return resource, fmt.Errorf("Could not find a form type for type %s", entity.EntityType)
	}

	err = e.db.NewSelect().Model(resource).Where("mrid = ?", mrid).OrderBy("commit_id", bun.OrderDesc).Limit(1).Scan(ctx)
	if err != nil {
		return resource, fmt.Errorf("Failed to collect data for editing resource: %w", err)
	}
	return resource, nil
}

func (e *EntityStore) EditComponentForm(w http.ResponseWriter, r *http.Request) {
	mrid := r.PathValue("mrid")
	resource, err := e.GetResource(r.Context(), mrid)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to fetch resource", "error", err)
		http.Error(w, "Failed to fetch resource "+err.Error(), http.StatusBadRequest)
		return
	}

	hxTriggerPayload := map[string]TriggerEditComponentForm{
		"editComponentFormChanged": {ResourceType: pkg.StructName(resource)},
	}
	hxTriggerPayloadBytes := pkg.Must(json.Marshal(hxTriggerPayload))
	w.Header().Set("HX-Trigger", string(hxTriggerPayloadBytes))

	pkg.FormInputFields(w, resource)
}

func (e *EntityStore) Resource(w http.ResponseWriter, r *http.Request) {
	mrid := r.PathValue("mrid")
	resource, err := e.GetResource(r.Context(), mrid)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to fetch resource", "error", err)
		http.Error(w, "Failed to fetch resource "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := ResourceItem{
		Data: resource,
		Type: pkg.StructName(resource),
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to write json", "error", err)
	}
}

func (e *EntityStore) Export(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	items, err := pkg.LatestOfAllItems(ctx, e.db, 0)
	if err != nil {
		slog.ErrorContext(ctx, "Could not fetch all items", "error", err)
		http.Error(w, "Could not fetch items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	itemIterator := func(yield func(v models.MridGetter) bool) {
		for _, item := range items {
			if !yield(item) {
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/n-triples")
	pkg.Export(w, itemIterator)
}

type ResourceItem struct {
	Data any    `json:"data"`
	Type string `json:"type"`
}

type TriggerEditComponentForm struct {
	ResourceType string `json:"resourceType"`
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

func choiceExists(items []models.VersionedObject, choice string) bool {
	for _, item := range items {
		if item.GetMrid().String() == choice {
			return true
		}
	}
	return false
}

func isDeleted(model any) bool {
	asDeleteGetter, ok := model.(models.DeletedGetter)
	return ok && asDeleteGetter.GetDeleted()
}
