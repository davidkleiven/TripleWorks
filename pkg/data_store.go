package pkg

import (
	"cmp"
	"context"
	"fmt"
	"iter"
	"log/slog"
	"slices"
	"strings"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MridName struct {
	Mrid string
	Name string
}

type DataStore struct {
	Db *bun.DB
}

func FindAll[T any](db *bun.DB, ctx context.Context, modelId int) ([]T, error) {
	_ = modelId // TODO: add multi model support later

	var entities []T
	err := db.NewSelect().Model(&entities).Scan(ctx)
	return entities, err
}

func FindNameAndMrid[T models.VersionedObject](db *bun.DB, ctx context.Context, modelId int) ([]models.VersionedObject, error) {
	result, err := FindAll[T](db, ctx, modelId)
	if err != nil {
		return []models.VersionedObject{}, fmt.Errorf("Failed to fetch all items: %w", err)
	}

	resultInterfaces := make([]models.VersionedObject, len(result))
	for i, item := range result {
		resultInterfaces[i] = item
	}
	return OnlyActiveLatest(resultInterfaces), nil
}

func FindEnum[T models.Enum](ctx context.Context, db *bun.DB) ([]models.Enum, error) {
	var entities []T
	err := db.NewSelect().Model(&entities).Scan(ctx)

	result := make([]models.Enum, len(entities))
	for i, v := range entities {
		result[i] = v
	}
	return result, err
}

func OnlyLatestVersion[T models.VersionedIdentifiedObject](items []T) []T {
	grouped := GroupBy(items, func(item T) uuid.UUID { return item.GetMrid() })
	result := make([]T, 0, len(grouped))
	for _, group := range grouped {
		latest := slices.MaxFunc(group, func(a, b T) int { return cmp.Compare(a.GetCommitId(), b.GetCommitId()) })
		result = append(result, latest)
	}
	return result
}

func OnlyActiveLatest[T models.VersionedObject](items []T) []T {
	latest := OnlyLatestVersion(items)
	toDelete := DeletedIndices(latest)
	return RemoveIndices(latest, toDelete)
}

func DeletedIndices[T models.DeletedGetter](items []T) []int {
	result := []int{}
	for i, item := range items {
		if item.GetDeleted() {
			result = append(result, i)
		}
	}
	return result
}

func RemoveIndices[T any](items []T, indices []int) []T {
	out := items[:0]
	idx := 0
	for i, v := range items {
		if (idx < len(indices)) && (i == indices[idx]) {
			idx++
			continue
		}
		out = append(out, v)
	}
	return out
}

type MridAndName struct {
	Mrid uuid.UUID
	Name string
}

func CreateFilteredAllFinder(nameFilter string, typeFilter string) Finder {
	candidates := Finders
	if typeFilter != "" {
		filteredCandidates := make(map[string]Finder)
		for k, finder := range Finders {
			if strings.Contains(strings.ToLower(k), strings.ToLower(typeFilter)) {
				filteredCandidates[k] = finder
			}
		}
		candidates = filteredCandidates
	}
	nameFilterFunc := NoOpNameFilter
	if nameFilter != "" {
		nameFilterFunc = CreateContainsNameFilter(nameFilter)
	}

	return func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return AllFinder(ctx, db, modelId, candidates, nameFilterFunc)
	}

}

type NameFilter func(name string) bool

func NoOpNameFilter(name string) bool {
	return true
}

func CreateContainsNameFilter(substr string) NameFilter {
	return func(name string) bool {
		return strings.Contains(strings.ToLower(name), strings.ToLower(substr))
	}
}

func AllFinder(ctx context.Context, db *bun.DB, modelId int, candidates map[string]Finder, include NameFilter) ([]models.VersionedObject, error) {
	limit := 100
	result := make([]models.VersionedObject, 0, limit)
	for name, finder := range candidates {
		items, err := finder(ctx, db, modelId)
		if err != nil {
			return result, fmt.Errorf("Failed to find data for %s: %w", name, err)
		}

		for _, item := range items {
			if !include(item.GetName()) && !include(item.GetMrid().String()) {
				continue
			}
			obj := models.IdentifiedObject{
				Mrid: item.GetMrid(),
				Name: item.GetName(),
				BaseEntity: models.BaseEntity{
					CommitId: item.GetCommitId(),
					Deleted:  item.GetDeleted(),
				},
			}

			result = append(result, obj)
			if len(result) >= limit {
				return result, nil
			}
		}
	}
	return result, nil
}

type Finder func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error)

func GetFinder(name, nameFilter, typeFilter string) (Finder, error) {
	if name == "all" {
		return CreateFilteredAllFinder(nameFilter, typeFilter), nil
	}
	finder, ok := Finders[name]
	if !ok {
		return nil, fmt.Errorf("Could not find a finder for %s", name)
	}
	return finder, nil
}

func NoOpOnInsert(v any) error {
	return nil
}

func InsertAll(ctx context.Context, db *bun.DB, msg string, items iter.Seq[any], onInsert func(v any) error) error {
	commit := models.Commit{
		Author:    "TripleWorks",
		Message:   msg,
		CreatedAt: time.Now(),
	}
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&commit).Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to insert commit: %w", err)
		}

		num := 0
		for item := range items {
			setCommitIfPossible(item, int(commit.Id))
			_, err := tx.NewInsert().Model(item).On("CONFLICT DO NOTHING").Exec(ctx)
			if err != nil {
				return fmt.Errorf("Failed after %d: %w", num, err)
			}

			if err := onInsert(item); err != nil {
				return fmt.Errorf("On insert failed after %d: %w", num, err)
			}
			num++
		}
		slog.InfoContext(ctx, "Inserted records into the database", "num", num)
		return nil
	})
}

func InsertAllInserter(ctx context.Context, inserter repository.Inserter, commit models.Commit, items iter.Seq[any], onInsert func(v any) error) error {
	insertFn := func(ctx context.Context, inserter repository.Inserter) error {
		err := inserter.Insert(ctx, &commit)
		if err != nil {
			return fmt.Errorf("Failed to insert commit: %w", err)
		}

		num := 0
		for item := range items {
			setCommitIfPossible(item, int(commit.Id))
			err = inserter.Insert(ctx, item)
			if err != nil {
				return fmt.Errorf("Failed after %d: %w", num, err)
			}

			if err := onInsert(item); err != nil {
				return fmt.Errorf("On insert failed after %d: %w", num, err)
			}
			num++
		}
		slog.InfoContext(ctx, "Inserted records into the database", "num", num)
		return nil
	}
	insertFnTx := repository.WithTx(insertFn)
	return insertFnTx(ctx, inserter)
}

func ExistingMrids(ctx context.Context, db *bun.DB, modelId int) ([]uuid.UUID, error) {
	var mrids []uuid.UUID
	err := db.NewSelect().TableExpr("entities").Column("mrid").Where("model_id = ?", modelId).Scan(ctx, &mrids)
	return mrids, err
}

func OnlyNewItems(existing map[uuid.UUID]struct{}, items iter.Seq[any]) iter.Seq[any] {
	return func(yield func(v any) bool) {
		for item := range items {
			mrid := mridIfPossible(item)
			_, exists := existing[mrid]
			if exists {
				continue
			}

			if !yield(item) {
				return
			}
		}
	}
}

func setCommitIfPossible(v any, commitId int) {
	commitSetter, ok := v.(models.CommitIdSetter)
	if ok {
		commitSetter.SetCommitId(commitId)
	}
}

func mridIfPossible(v any) uuid.UUID {
	mridGetter, ok := v.(models.MridGetter)
	if ok {
		return mridGetter.GetMrid()
	}
	return uuid.UUID{}
}

type LonLat struct {
	Lat float64
	Lon float64
}

func LinesConnectedToSubstationByName(ctx context.Context, db *bun.DB, substation *models.Substation) ([]models.ACLineSegment, error) {
	var (
		lines       []models.ACLineSegment
		substations []models.Substation
	)

	// Collect line and substations that matches the name of the target substation
	failNo, err := ReturnOnFirstError(
		func() error {
			return db.NewSelect().Model(&lines).Where("? LIKE ?", bun.Ident("name"), fmt.Sprintf("%%%s%%", substation.Name)).Scan(ctx)
		},
		func() error {
			return db.NewSelect().
				Model(&substations).
				Where("? LIKE ?", bun.Ident("name"), fmt.Sprintf("%%%s%%", substation.Name)).
				Scan(ctx)
		},
	)

	if err != nil {
		return lines, fmt.Errorf("Failed to find ac lines by name call no %d: %w", failNo, err)
	}
	lines = OnlyActiveLatest(lines)
	substations = OnlyActiveLatest(substations)

	// Lines should only be included if it has the highest similarity score with the current station
	lines = slices.DeleteFunc(lines, func(line models.ACLineSegment) bool {
		var (
			highestScoreMrid uuid.UUID
			highestScore     float64 = -1.0
		)
		lineName := Normalizename(line.Name)
		lineTokens := Tokenize(lineName)

		for _, sub := range substations {
			substationName := Normalizename(sub.Name)
			substationTokens := Tokenize(substationName)
			cosScore := CosineSimilarity(lineName, substationName)
			exactScore := ExactTokenSimilarity(substationTokens, lineTokens)

			// The score is weighted sum of exact matches and cosine similarity
			score := 0.8*exactScore + 0.2*cosScore
			if score > highestScore {
				highestScore = score
				highestScoreMrid = sub.Mrid
			}
		}
		return highestScoreMrid != substation.Mrid
	})
	return lines, nil
}

var Finders = map[string]Finder{
	"ACDCConverter": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ACDCConverter](db, ctx, modelId)
	},
	"ACDCConverterDCTerminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ACDCConverterDCTerminal](db, ctx, modelId)
	},
	"ACDCTerminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ACDCTerminal](db, ctx, modelId)
	},
	"ACLineSegment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ACLineSegment](db, ctx, modelId)
	},
	"AsynchronousMachine": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.AsynchronousMachine](db, ctx, modelId)
	},
	"BaseVoltage": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.BaseVoltage](db, ctx, modelId)
	},
	"BasicIntervalSchedule": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.BasicIntervalSchedule](db, ctx, modelId)
	},
	"Breaker": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Breaker](db, ctx, modelId)
	},
	"BusNameMarker": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.BusNameMarker](db, ctx, modelId)
	},
	"BusbarSection": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.BusbarSection](db, ctx, modelId)
	},
	"ConductingEquipment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ConductingEquipment](db, ctx, modelId)
	},
	"Conductor": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Conductor](db, ctx, modelId)
	},
	"ConformLoad": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ConformLoad](db, ctx, modelId)
	},
	"ConformLoadGroup": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ConformLoadGroup](db, ctx, modelId)
	},
	"ConnectivityNode": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ConnectivityNode](db, ctx, modelId)
	},
	"ConnectivityNodeContainer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ConnectivityNodeContainer](db, ctx, modelId)
	},
	"Connector": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Connector](db, ctx, modelId)
	},
	"ControlArea": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ControlArea](db, ctx, modelId)
	},
	"ControlAreaGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ControlAreaGeneratingUnit](db, ctx, modelId)
	},
	"CsConverter": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.CsConverter](db, ctx, modelId)
	},
	"CurrentLimit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.CurrentLimit](db, ctx, modelId)
	},
	"Curve": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Curve](db, ctx, modelId)
	},
	"DCBaseTerminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCBaseTerminal](db, ctx, modelId)
	},
	"DCBreaker": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCBreaker](db, ctx, modelId)
	},
	"DCBusbar": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCBusbar](db, ctx, modelId)
	},
	"DCChopper": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCChopper](db, ctx, modelId)
	},
	"DCConductingEquipment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCConductingEquipment](db, ctx, modelId)
	},
	"DCConverterUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCConverterUnit](db, ctx, modelId)
	},
	"DCDisconnector": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCDisconnector](db, ctx, modelId)
	},
	"DCEquipmentContainer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCEquipmentContainer](db, ctx, modelId)
	},
	"DCGround": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCGround](db, ctx, modelId)
	},
	"DCLine": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCLine](db, ctx, modelId)
	},
	"DCLineSegment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCLineSegment](db, ctx, modelId)
	},
	"DCNode": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCNode](db, ctx, modelId)
	},
	"DCSeriesDevice": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCSeriesDevice](db, ctx, modelId)
	},
	"DCShunt": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCShunt](db, ctx, modelId)
	},
	"DCSwitch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCSwitch](db, ctx, modelId)
	},
	"DCTerminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.DCTerminal](db, ctx, modelId)
	},
	"Disconnector": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Disconnector](db, ctx, modelId)
	},
	"EnergyConsumer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EnergyConsumer](db, ctx, modelId)
	},
	"EnergySchedulingType": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EnergySchedulingType](db, ctx, modelId)
	},
	"EnergySource": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EnergySource](db, ctx, modelId)
	},
	"Equipment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Equipment](db, ctx, modelId)
	},
	"EquipmentContainer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EquipmentContainer](db, ctx, modelId)
	},
	"EquivalentBranch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EquivalentBranch](db, ctx, modelId)
	},
	"EquivalentEquipment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EquivalentEquipment](db, ctx, modelId)
	},
	"EquivalentInjection": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EquivalentInjection](db, ctx, modelId)
	},
	"EquivalentNetwork": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EquivalentNetwork](db, ctx, modelId)
	},
	"EquivalentShunt": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.EquivalentShunt](db, ctx, modelId)
	},
	"ExternalNetworkInjection": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ExternalNetworkInjection](db, ctx, modelId)
	},
	"FossilFuel": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.FossilFuel](db, ctx, modelId)
	},
	"GeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.GeneratingUnit](db, ctx, modelId)
	},
	"GeographicalRegion": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.GeographicalRegion](db, ctx, modelId)
	},
	"HydroGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.HydroGeneratingUnit](db, ctx, modelId)
	},
	"HydroPowerPlant": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.HydroPowerPlant](db, ctx, modelId)
	},
	"HydroPump": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.HydroPump](db, ctx, modelId)
	},
	"IdentifiedObject": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.IdentifiedObject](db, ctx, modelId)
	},
	"InitialReactiveCapabilityCurve": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ReactiveCapabilityCurve](db, ctx, modelId)
	},
	"Junction": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Junction](db, ctx, modelId)
	},
	"Line": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Line](db, ctx, modelId)
	},
	"LinearShuntCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.LinearShuntCompensator](db, ctx, modelId)
	},
	"LoadBreakSwitch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.LoadBreakSwitch](db, ctx, modelId)
	},
	"LoadGroup": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.LoadGroup](db, ctx, modelId)
	},
	"LoadResponseCharacteristic": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.LoadResponseCharacteristic](db, ctx, modelId)
	},
	"NonConformLoad": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.NonConformLoad](db, ctx, modelId)
	},
	"NonConformLoadGroup": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.NonConformLoadGroup](db, ctx, modelId)
	},
	"NonlinearShuntCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.NonlinearShuntCompensator](db, ctx, modelId)
	},
	"NuclearGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.NuclearGeneratingUnit](db, ctx, modelId)
	},
	"OperationalLimit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.OperationalLimit](db, ctx, modelId)
	},
	"OperationalLimitSet": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.OperationalLimitSet](db, ctx, modelId)
	},
	"OperationalLimitType": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.OperationalLimitType](db, ctx, modelId)
	},
	"PhaseTapChanger": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PhaseTapChanger](db, ctx, modelId)
	},
	"PhaseTapChangerAsymmetrical": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PhaseTapChangerAsymmetrical](db, ctx, modelId)
	},
	"PhaseTapChangerLinear": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PhaseTapChangerLinear](db, ctx, modelId)
	},
	"PhaseTapChangerNonLinear": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PhaseTapChangerNonLinear](db, ctx, modelId)
	},
	"PhaseTapChangerSymmetrical": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PhaseTapChangerSymmetrical](db, ctx, modelId)
	},
	"PhaseTapChangerTable": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PhaseTapChangerTable](db, ctx, modelId)
	},
	"PhaseTapChangerTabular": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PhaseTapChangerTabular](db, ctx, modelId)
	},
	"PowerSystemResource": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PowerSystemResource](db, ctx, modelId)
	},
	"PowerTransformer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PowerTransformer](db, ctx, modelId)
	},
	"PowerTransformerEnd": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.PowerTransformerEnd](db, ctx, modelId)
	},
	"ProtectedSwitch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ProtectedSwitch](db, ctx, modelId)
	},
	"RatioTapChanger": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.RatioTapChanger](db, ctx, modelId)
	},
	"RatioTapChangerTable": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.RatioTapChangerTable](db, ctx, modelId)
	},
	"ReactiveCapabilityCurve": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ReactiveCapabilityCurve](db, ctx, modelId)
	},
	"Region": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.SubGeographicalRegion](db, ctx, modelId)
	},
	"RegularIntervalSchedule": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.RegularIntervalSchedule](db, ctx, modelId)
	},
	"RegulatingCondEq": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.RegulatingCondEq](db, ctx, modelId)
	},
	"RegulatingControl": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.RegulatingControl](db, ctx, modelId)
	},
	"ReportingGroup": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ReportingGroup](db, ctx, modelId)
	},
	"RotatingMachine": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.RotatingMachine](db, ctx, modelId)
	},
	"SeriesCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.SeriesCompensator](db, ctx, modelId)
	},
	"ShuntCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ShuntCompensator](db, ctx, modelId)
	},
	"SolarGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.SolarGeneratingUnit](db, ctx, modelId)
	},
	"StaticVarCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.StaticVarCompensator](db, ctx, modelId)
	},
	"SubGeographicalRegion": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.SubGeographicalRegion](db, ctx, modelId)
	},
	"Substation": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Substation](db, ctx, modelId)
	},
	"Switch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Switch](db, ctx, modelId)
	},
	"SynchronousMachine": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.SynchronousMachine](db, ctx, modelId)
	},
	"TapChanger": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.TapChanger](db, ctx, modelId)
	},
	"TapChangerControl": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.TapChangerControl](db, ctx, modelId)
	},
	"Terminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.Terminal](db, ctx, modelId)
	},
	"ThermalGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.ThermalGeneratingUnit](db, ctx, modelId)
	},
	"TransformerEnd": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.TransformerEnd](db, ctx, modelId)
	},
	"VoltageLevel": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.VoltageLevel](db, ctx, modelId)
	},
	"VoltageLimit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.VoltageLimit](db, ctx, modelId)
	},
	"VsCapabilityCurve": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.VsCapabilityCurve](db, ctx, modelId)
	},
	"VsConverter": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.VsConverter](db, ctx, modelId)
	},
	"WindGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
		return FindNameAndMrid[models.WindGeneratingUnit](db, ctx, modelId)
	},
}

type EnumFinder func(ctx context.Context, db *bun.DB) ([]models.Enum, error)

var EnumFinders = map[string]EnumFinder{
	"ControlAreaTypeKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.ControlAreaTypeKind](ctx, db)
	},
	"Currency": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.Currency](ctx, db)
	},
	"CurveStyle": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.CurveStyle](ctx, db)
	},
	"DCConverterOperatingModeKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.DCConverterOperatingModeKind](ctx, db)
	},
	"DCPolarityKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.DCPolarityKind](ctx, db)
	},
	"FuelType": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.FuelType](ctx, db)
	},
	"GeneratorControlSource": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.GeneratorControlSource](ctx, db)
	},
	"HydroEnergyConversionKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.HydroEnergyConversionKind](ctx, db)
	},
	"HydroPlantStorageKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.HydroPlantStorageKind](ctx, db)
	},
	"LimitTypeKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.LimitTypeKind](ctx, db)
	},
	"OperationalLimitDirectionKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.OperationalLimitDirectionKind](ctx, db)
	},
	"PetersenCoilModeKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.PetersenCoilModeKind](ctx, db)
	},
	"PhaseCode": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.PhaseCode](ctx, db)
	},
	"RegulatingControlModeKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.RegulatingControlModeKind](ctx, db)
	},
	"SVCControlMode": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.SVCControlMode](ctx, db)
	},
	"ShortCircuitRotorKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.ShortCircuitRotorKind](ctx, db)
	},
	"SynchronousMachineKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.SynchronousMachineKind](ctx, db)
	},
	"TransformerControlMode": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.TransformerControlMode](ctx, db)
	},
	"UnitMultiplier": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.UnitMultiplier](ctx, db)
	},
	"UnitSymbol": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.UnitSymbol](ctx, db)
	},
	"WindGenUnitKind": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.WindGenUnitKind](ctx, db)
	},
	"WindingConnection": func(ctx context.Context, db *bun.DB) ([]models.Enum, error) {
		return FindEnum[models.WindingConnection](ctx, db)
	},
}
