package api

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
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
				entity.CommitId = int(commit.Id)
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

func (e *EntityStore) SimpleUpload(w http.ResponseWriter, r *http.Request) {
	hundredMb := int64(100 << 20)
	kind := r.PathValue("kind")
	doCommit := r.URL.Query().Get("commit")

	r.Body = http.MaxBytesReader(w, r.Body, hundredMb)
	defer r.Body.Close()

	var (
		substations = "substations"
		generators  = "generators"
		loads       = "loads"
		lines       = "lines"
	)

	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	modelId := 0
	existing, err := pkg.ExistingMrids(ctx, e.db, modelId)
	if err != nil {
		slog.ErrorContext(ctx, "Could not get existing mrids", "error", err)
		http.Error(w, "Could not get existing mrids: "+err.Error(), http.StatusInternalServerError)
		return
	}

	existingSet := make(map[uuid.UUID]struct{})
	for _, mrid := range existing {
		existingSet[mrid] = struct{}{}
	}

	scanner := bufio.NewScanner(r.Body)
	num := 0
	itemIterators := []iter.Seq[any]{}
	var rawBytes bytes.Buffer
	for scanner.Scan() {
		num++
		line := scanner.Bytes()

		var (
			err          error
			itemIterator iter.Seq[any]
		)
		_, err = rawBytes.Write(line)
		switch kind {
		case substations:
			var substation pkg.SubstationLight
			err = json.Unmarshal(line, &substation)
			itemIterator = pkg.OnlyNewItems(existingSet, substation.CimItems(modelId))
		case generators:
			var generator pkg.GeneratorLight
			err = json.Unmarshal(line, &generator)
			itemIterator = pkg.OnlyNewItems(existingSet, generator.CimItems(modelId))
		case loads:
			var load pkg.LoadLight
			err = json.Unmarshal(line, &load)
			itemIterator = pkg.OnlyNewItems(existingSet, load.CimItems(modelId))
		case lines:
			var acline pkg.LineLight
			err = json.Unmarshal(line, &acline)
			itemIterator = pkg.OnlyNewItems(existingSet, acline.CimItems(modelId))
		default:
			err = fmt.Errorf("Unknown type %s", kind)
		}

		if err != nil {
			slog.ErrorContext(ctx, "Could not unmarshal line", "kind", kind, "lineNo", num, "error", err)
			http.Error(w, "Could not unmarshal line: "+err.Error(), http.StatusBadRequest)
			return
		}
		itemIterators = append(itemIterators, itemIterator)
	}

	rawData := models.SimpleUpload{Data: rawBytes.Bytes()}
	rawDataIter := func(yield func(v any) bool) {
		yield(&rawData)
	}
	itemIterators = append(itemIterators, rawDataIter)

	w.Header().Set("Content-Type", "application/n-triples")
	itemIterator := pkg.Chain(itemIterators...)
	if doCommit == "true" {
		msg := fmt.Sprintf("Add %d %s", num, kind)
		err := pkg.InsertAll(ctx, e.db, msg, itemIterator, writeNTriplesCallback(w))
		if err != nil {
			slog.ErrorContext(ctx, "Could not insert new items", "error", err)
			http.Error(w, "Could not insert new items: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		writer := writeNTriplesCallback(w)
		for item := range itemIterator {
			writer(item)
		}
	}
}

func (e *EntityStore) Commits(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	var commits []models.Commit
	err := e.db.NewSelect().Model(&commits).Scan(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Could not fetch commits", "error", err)
		http.Error(w, "Could not fetch commits: "+err.Error(), http.StatusInternalServerError)
		return
	}
	pkg.PanicOnErr(json.NewEncoder(w).Encode(&commits))
	w.Header().Set("Content-Type", "application/json")
}

func (e *EntityStore) DeleteCommit(w http.ResponseWriter, r *http.Request) {
	commitIdStr := r.PathValue("id")
	commitId, err := strconv.Atoi(commitIdStr)
	if err != nil {
		slog.ErrorContext(r.Context(), "Commit id is not an integer", "commitId", commitIdStr)
		http.Error(w, "'%s' is not an integer", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	var (
		affectedRows  int
		skippedTables []string
	)
	txErr := e.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		for _, itemPtr := range pkg.FormTypes() {
			_, isVersionedObject := itemPtr.(models.VersionedIdentifiedObject)
			if !isVersionedObject {
				skippedTables = append(skippedTables, pkg.StructName(itemPtr))
				continue
			}

			res, err := tx.NewDelete().Model(itemPtr).Where("commit_id = ?", commitId).Exec(ctx)
			if err != nil {
				return fmt.Errorf("Failed to delete from table %s: %w", pkg.StructName(itemPtr), err)
			}
			rows, _ := res.RowsAffected()
			affectedRows += int(rows)
		}
		_, err := tx.NewDelete().Model((*models.Commit)(nil)).Where("id = ?", commitId).Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to delete from commit table: %w", err)
		}

		_, err = tx.NewDelete().Model((*models.Entity)(nil)).Where("commit_id = ?", commitId).Exec(ctx)
		return err
	})

	slog.InfoContext(ctx, "Skipped delete for tables", "affectedRows", affectedRows, "skipped", skippedTables)
	if txErr != nil {
		slog.ErrorContext(ctx, "Failed to delete commit", "error", txErr)
		http.Error(w, "Failed to delete commit: "+txErr.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Successfully deleted commit %d", commitId)
}

func (e *EntityStore) Map(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	modelId := 0
	var (
		substations []models.Substation
		acLines     []models.ACLineSegment
		points      []models.PositionPoint
		bvs         []models.BaseVoltage
		vls         []models.VoltageLevel
		terminals   []models.Terminal
		cns         []models.ConnectivityNode
	)

	failNo, err := pkg.ReturnOnFirstError(
		func() error {
			return e.db.NewSelect().Model(&substations).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&acLines).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&points).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&bvs).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&vls).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&terminals).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&cns).Scan(ctx)
		},
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract substations", "error", err, "failNo", failNo)
		http.Error(w, "Failed to extract substations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	substations = pkg.OnlyActiveLatest(substations)
	bvs = pkg.OnlyActiveLatest(bvs)
	acLines = pkg.OnlyActiveLatest(acLines)
	vls = pkg.OnlyActiveLatest(vls)
	terminals = pkg.OnlyActiveLatest(terminals)
	cns = pkg.OnlyActiveLatest(cns)

	ptMap := pkg.IndexBy(points, func(p models.PositionPoint) uuid.UUID { return p.LocationMrid })
	bvMap := pkg.IndexBy(bvs, func(b models.BaseVoltage) uuid.UUID { return b.Mrid })
	tMap := pkg.GroupBy(terminals, func(t models.Terminal) uuid.UUID { return t.ConductingEquipmentMrid })
	cnsMap := pkg.IndexBy(cns, func(c models.ConnectivityNode) uuid.UUID { return c.Mrid })
	vlMap := pkg.IndexBy(vls, func(v models.VoltageLevel) uuid.UUID { return v.Mrid })
	subMap := pkg.IndexBy(substations, func(s models.Substation) uuid.UUID { return s.Mrid })

	acLineFromToMap := make(map[uuid.UUID]FromToLoc)
	ignoredLines := []string{}
	for _, line := range acLines {
		connectedTerminals, ok := tMap[line.Mrid]
		if !ok || len(connectedTerminals) != 2 {
			ignoredLines = append(ignoredLines, line.Name)
			continue
		}
		cn1 := cnsMap[connectedTerminals[0].ConnectivityNodeMrid]
		cn2 := cnsMap[connectedTerminals[1].ConnectivityNodeMrid]
		acLineFromToMap[line.Mrid] = FromToLoc{
			Pt1: ptMap[subMap[vlMap[cn1.ConnectivityNodeContainerMrid].SubstationMrid].LocationMrid],
			Pt2: ptMap[subMap[vlMap[cn2.ConnectivityNodeContainerMrid].SubstationMrid].LocationMrid],
		}
	}

	skipped := 0
	substationMapDataIter := func(yield func(v pkg.SubstationMapData) bool) {
		for _, sub := range substations {
			pt, ok := ptMap[sub.LocationMrid]
			if !ok {
				skipped++
				continue
			}

			data := pkg.SubstationMapData{
				Lat:  pt.YPosition,
				Lng:  pt.XPosition,
				Mrid: sub.Mrid.String(),
				Name: sub.Name,
			}

			if !yield(data) {
				return
			}
		}
	}
	slog.InfoContext(ctx, "Extracted substations", "num", len(substations), "modelId", modelId, "ignoredLines", ignoredLines, "numSkippedSubstations", skipped)

	lineIter := func(yield func(v pkg.LineMapData) bool) {
		for _, line := range acLines {
			fromTo := acLineFromToMap[line.Mrid]
			vl := bvMap[line.BaseVoltageMrid]
			data := pkg.LineMapData{
				LatFrom: fromTo.Pt1.YPosition,
				LatTo:   fromTo.Pt2.YPosition,
				LngFrom: fromTo.Pt1.XPosition,
				LngTo:   fromTo.Pt2.XPosition,
				Mrid:    line.Mrid.String(),
				Name:    line.Name,
				Voltage: int(vl.NominalVoltage),
			}

			if !yield(data) {
				return
			}
		}
	}
	pkg.RenderMap(w, substationMapDataIter, lineIter)
}

func (e *EntityStore) ConnectDanglingLines(w http.ResponseWriter, r *http.Request) {
	doCommit := r.URL.Query().Get("commit")
	var (
		substations []models.Substation
		lines       []models.ACLineSegment
		terminals   []models.Terminal
		vls         []models.VoltageLevel
	)

	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	failNo, err := pkg.ReturnOnFirstError(
		func() error {
			return e.db.NewSelect().Model(&substations).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&lines).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&terminals).Scan(ctx)
		},
		func() error {
			return e.db.NewSelect().Model(&vls).Scan(ctx)
		},
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to get data", "failNo", failNo, "error", err)
		http.Error(w, "Failed to extract data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	lines = pkg.OnlyActiveLatest(lines)
	terminals = pkg.OnlyActiveLatest(terminals)
	substations = pkg.OnlyActiveLatest(substations)
	vls = pkg.OnlyActiveLatest(vls)
	vlsPerSubstation := pkg.GroupBy(vls, func(vl models.VoltageLevel) uuid.UUID { return vl.SubstationMrid })
	unconnecteLines := pkg.DanglingLines(lines, terminals)

	substationaNames := make([]string, len(substations))
	lineNames := make([]string, 0, len(lines))
	lines = lines[:0] // Clear old
	parenthesisExpr := regexp.MustCompile(`\([^)]+\)`)
	voltageExpr := regexp.MustCompile(`(?i)[0-9\s]+kv`)
	for i, sub := range substations {
		name := parenthesisExpr.ReplaceAllString(sub.Name, "")
		name = voltageExpr.ReplaceAllString(name, "")
		substationaNames[i] = name
	}
	for line := range unconnecteLines {
		name := parenthesisExpr.ReplaceAllLiteralString(line.Name, "")
		name = voltageExpr.ReplaceAllString(name, "")
		name = strings.ReplaceAll(name, "-", " ")
		lineNames = append(lineNames, name)
		lines = append(lines, line)
	}

	selector := pkg.TopSelector{Num: 2}
	assignments := selector.Select(lineNames, substationaNames, pkg.NameSimilarity)

	results := make([]iter.Seq[any], 0, len(lines)*2)
	lines = pkg.MustSlice(lines)
	substations = pkg.MustSlice(substations)
	lines = pkg.RequireSameLength(assignments, lines)
	for lineIdx := range assignments {
		line := lines[lineIdx]
		for _, subIdx := range assignments[lineIdx] {
			sub := substations[subIdx]
			vls, ok := vlsPerSubstation[sub.Mrid]
			if !ok {
				vls = []models.VoltageLevel{}
			}

			params := pkg.LineConnectionParams{
				Substation:    sub,
				Line:          line,
				VoltageLevels: vls,
				Terminals:     terminals,
			}
			result := pkg.Must(pkg.ConnectLineToSubstation(params))

			if result.VoltageLevel != nil {
				vlsPerSubstation[sub.Mrid] = append(vlsPerSubstation[sub.Mrid], *result.VoltageLevel)
			}
			results = append(results, result.All(0))
		}
	}

	w.Header().Set("Content-Type", "application/n-triples")
	allItems := pkg.Chain(results...)
	if doCommit == "true" {
		lineName := ""
		for _, line := range lines {
			lineName += line.Name + " "
		}
		msg := fmt.Sprintf("Connect %d lines to two substations each", len(lines))
		err := pkg.InsertAll(ctx, e.db, msg, allItems, writeNTriplesCallback(w))
		if err != nil {
			slog.ErrorContext(ctx, "Could not insert items", "error", err)
			http.Error(w, "Could not insert items", http.StatusInternalServerError)
			return
		}
	} else {
		writer := writeNTriplesCallback(w)
		for item := range allItems {
			writer(item)
		}
	}

}

func (e *EntityStore) ApplyJsonPatch(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var jsonPatch pkg.JsonPatch
	if err := json.NewDecoder(r.Body).Decode(&jsonPatch); err != nil {
		slog.ErrorContext(r.Context(), "Failed to interpret json patch", "error", err)
		http.Error(w, "Failed to interpret json patch: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()
	if err := pkg.ApplyPatch(ctx, e.db, jsonPatch); err != nil {
		slog.ErrorContext(ctx, "Failed to apply patch", "error", err)
		http.Error(w, "Failed to apply patch: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Successfully updated database with patch")
}

func (e *EntityStore) Connection(w http.ResponseWriter, r *http.Request) {
	mrid := r.PathValue("mrid")
	ctx, cancel := context.WithTimeout(r.Context(), e.timeout)
	defer cancel()

	data, err := pkg.FetchConnectionData(ctx, e.db, mrid)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to find connection", "error", err)
		http.Error(w, "Failed to find connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	con := pkg.FindConnection(&data)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(con)
}

type ResourceItem struct {
	Data any    `json:"data"`
	Type string `json:"type"`
}

type TriggerEditComponentForm struct {
	ResourceType string `json:"resourceType"`
}

type FromToLoc struct {
	Pt1 models.PositionPoint
	Pt2 models.PositionPoint
}

func NewEntityStore(db *bun.DB, timeout time.Duration) *EntityStore {
	store := EntityStore{
		db:      db,
		timeout: timeout,
		allowedUnset: map[string]struct{}{
			"Id":        {},
			"CommitId":  {},
			"id":        {},
			"commit_id": {},
		},
	}
	return &store
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

func writeNTriplesCallback(w io.Writer) func(item any) error {
	return func(item any) error {
		mridGetter, ok := item.(models.MridGetter)
		if !ok {
			return nil
		}
		pkg.ExportItem(w, mridGetter)
		return nil
	}
}
