package pkg

import (
	"context"
	"fmt"
	"slices"

	"com.github/davidkleiven/tripleworks/models"
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

func FindNameAndMrid[T models.MridNameGetter](db *bun.DB, ctx context.Context, modelId int) ([]models.MridNameGetter, error) {
	result, err := FindAll[T](db, ctx, modelId)
	if err != nil {
		return []models.MridNameGetter{}, fmt.Errorf("Failed to fetch all items: %w", err)
	}

	resultInterfaces := make([]models.MridNameGetter, len(result))
	for i, item := range result {
		resultInterfaces[i] = item
	}
	return resultInterfaces, nil
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
	// Pick max commit per mrid
	toKeep := make(map[uuid.UUID]int)
	for _, item := range items {
		mrid := item.GetMrid()
		commitId := item.GetCommitId()
		maxCommit, ok := toKeep[mrid]
		if !ok || commitId > maxCommit {
			toKeep[mrid] = commitId
		}
	}

	return slices.DeleteFunc(items, func(item T) bool {
		mrid := item.GetMrid()
		commitId := item.GetCommitId()
		maxCommitId := toKeep[mrid]
		return commitId != maxCommitId
	})
}

type MridAndName struct {
	Mrid uuid.UUID
	Name string
}

type Finder func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error)

var Finders = map[string]Finder{
	"ACDCConverter": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ACDCConverter](db, ctx, modelId)
	},
	"ACDCConverterDCTerminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ACDCConverterDCTerminal](db, ctx, modelId)
	},
	"ACDCTerminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ACDCTerminal](db, ctx, modelId)
	},
	"ACLineSegment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ACLineSegment](db, ctx, modelId)
	},
	"AsynchronousMachine": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.AsynchronousMachine](db, ctx, modelId)
	},
	"BaseVoltage": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.BaseVoltage](db, ctx, modelId)
	},
	"BasicIntervalSchedule": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.BasicIntervalSchedule](db, ctx, modelId)
	},
	"Breaker": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Breaker](db, ctx, modelId)
	},
	"BusNameMarker": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.BusNameMarker](db, ctx, modelId)
	},
	"BusbarSection": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.BusbarSection](db, ctx, modelId)
	},
	"ConductingEquipment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ConductingEquipment](db, ctx, modelId)
	},
	"Conductor": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Conductor](db, ctx, modelId)
	},
	"ConformLoad": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ConformLoad](db, ctx, modelId)
	},
	"ConformLoadGroup": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ConformLoadGroup](db, ctx, modelId)
	},
	"ConnectivityNode": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ConnectivityNode](db, ctx, modelId)
	},
	"ConnectivityNodeContainer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ConnectivityNodeContainer](db, ctx, modelId)
	},
	"Connector": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Connector](db, ctx, modelId)
	},
	"ControlArea": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ControlArea](db, ctx, modelId)
	},
	"ControlAreaGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ControlAreaGeneratingUnit](db, ctx, modelId)
	},
	"CsConverter": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.CsConverter](db, ctx, modelId)
	},
	"CurrentLimit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.CurrentLimit](db, ctx, modelId)
	},
	"Curve": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Curve](db, ctx, modelId)
	},
	"DCBaseTerminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCBaseTerminal](db, ctx, modelId)
	},
	"DCBreaker": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCBreaker](db, ctx, modelId)
	},
	"DCBusbar": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCBusbar](db, ctx, modelId)
	},
	"DCChopper": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCChopper](db, ctx, modelId)
	},
	"DCConductingEquipment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCConductingEquipment](db, ctx, modelId)
	},
	"DCConverterUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCConverterUnit](db, ctx, modelId)
	},
	"DCDisconnector": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCDisconnector](db, ctx, modelId)
	},
	"DCEquipmentContainer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCEquipmentContainer](db, ctx, modelId)
	},
	"DCGround": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCGround](db, ctx, modelId)
	},
	"DCLine": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCLine](db, ctx, modelId)
	},
	"DCLineSegment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCLineSegment](db, ctx, modelId)
	},
	"DCNode": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCNode](db, ctx, modelId)
	},
	"DCSeriesDevice": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCSeriesDevice](db, ctx, modelId)
	},
	"DCShunt": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCShunt](db, ctx, modelId)
	},
	"DCSwitch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCSwitch](db, ctx, modelId)
	},
	"DCTerminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.DCTerminal](db, ctx, modelId)
	},
	"Disconnector": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Disconnector](db, ctx, modelId)
	},
	"EnergyConsumer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EnergyConsumer](db, ctx, modelId)
	},
	"EnergySchedulingType": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EnergySchedulingType](db, ctx, modelId)
	},
	"EnergySource": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EnergySource](db, ctx, modelId)
	},
	"Equipment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Equipment](db, ctx, modelId)
	},
	"EquipmentContainer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EquipmentContainer](db, ctx, modelId)
	},
	"EquivalentBranch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EquivalentBranch](db, ctx, modelId)
	},
	"EquivalentEquipment": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EquivalentEquipment](db, ctx, modelId)
	},
	"EquivalentInjection": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EquivalentInjection](db, ctx, modelId)
	},
	"EquivalentNetwork": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EquivalentNetwork](db, ctx, modelId)
	},
	"EquivalentShunt": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.EquivalentShunt](db, ctx, modelId)
	},
	"ExternalNetworkInjection": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ExternalNetworkInjection](db, ctx, modelId)
	},
	"FossilFuel": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.FossilFuel](db, ctx, modelId)
	},
	"GeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.GeneratingUnit](db, ctx, modelId)
	},
	"GeographicalRegion": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.GeographicalRegion](db, ctx, modelId)
	},
	"HydroGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.HydroGeneratingUnit](db, ctx, modelId)
	},
	"HydroPowerPlant": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.HydroPowerPlant](db, ctx, modelId)
	},
	"HydroPump": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.HydroPump](db, ctx, modelId)
	},
	"IdentifiedObject": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.IdentifiedObject](db, ctx, modelId)
	},
	"Junction": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Junction](db, ctx, modelId)
	},
	"Line": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Line](db, ctx, modelId)
	},
	"LinearShuntCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.LinearShuntCompensator](db, ctx, modelId)
	},
	"LoadBreakSwitch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.LoadBreakSwitch](db, ctx, modelId)
	},
	"LoadGroup": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.LoadGroup](db, ctx, modelId)
	},
	"LoadResponseCharacteristic": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.LoadResponseCharacteristic](db, ctx, modelId)
	},
	"NonConformLoad": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.NonConformLoad](db, ctx, modelId)
	},
	"NonConformLoadGroup": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.NonConformLoadGroup](db, ctx, modelId)
	},
	"NonlinearShuntCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.NonlinearShuntCompensator](db, ctx, modelId)
	},
	"NuclearGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.NuclearGeneratingUnit](db, ctx, modelId)
	},
	"OperationalLimit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.OperationalLimit](db, ctx, modelId)
	},
	"OperationalLimitSet": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.OperationalLimitSet](db, ctx, modelId)
	},
	"OperationalLimitType": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.OperationalLimitType](db, ctx, modelId)
	},
	"PhaseTapChanger": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PhaseTapChanger](db, ctx, modelId)
	},
	"PhaseTapChangerAsymmetrical": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PhaseTapChangerAsymmetrical](db, ctx, modelId)
	},
	"PhaseTapChangerLinear": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PhaseTapChangerLinear](db, ctx, modelId)
	},
	"PhaseTapChangerNonLinear": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PhaseTapChangerNonLinear](db, ctx, modelId)
	},
	"PhaseTapChangerSymmetrical": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PhaseTapChangerSymmetrical](db, ctx, modelId)
	},
	"PhaseTapChangerTable": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PhaseTapChangerTable](db, ctx, modelId)
	},
	"PhaseTapChangerTabular": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PhaseTapChangerTabular](db, ctx, modelId)
	},
	"PowerSystemResource": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PowerSystemResource](db, ctx, modelId)
	},
	"PowerTransformer": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PowerTransformer](db, ctx, modelId)
	},
	"PowerTransformerEnd": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.PowerTransformerEnd](db, ctx, modelId)
	},
	"ProtectedSwitch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ProtectedSwitch](db, ctx, modelId)
	},
	"RatioTapChanger": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.RatioTapChanger](db, ctx, modelId)
	},
	"RatioTapChangerTable": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.RatioTapChangerTable](db, ctx, modelId)
	},
	"ReactiveCapabilityCurve": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ReactiveCapabilityCurve](db, ctx, modelId)
	},
	"Region": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.SubGeographicalRegion](db, ctx, modelId)
	},
	"RegularIntervalSchedule": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.RegularIntervalSchedule](db, ctx, modelId)
	},
	"RegulatingCondEq": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.RegulatingCondEq](db, ctx, modelId)
	},
	"RegulatingControl": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.RegulatingControl](db, ctx, modelId)
	},
	"ReportingGroup": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ReportingGroup](db, ctx, modelId)
	},
	"RotatingMachine": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.RotatingMachine](db, ctx, modelId)
	},
	"SeriesCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.SeriesCompensator](db, ctx, modelId)
	},
	"ShuntCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ShuntCompensator](db, ctx, modelId)
	},
	"SolarGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.SolarGeneratingUnit](db, ctx, modelId)
	},
	"StaticVarCompensator": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.StaticVarCompensator](db, ctx, modelId)
	},
	"SubGeographicalRegion": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.SubGeographicalRegion](db, ctx, modelId)
	},
	"Substation": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Substation](db, ctx, modelId)
	},
	"Switch": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Switch](db, ctx, modelId)
	},
	"SynchronousMachine": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.SynchronousMachine](db, ctx, modelId)
	},
	"TapChanger": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.TapChanger](db, ctx, modelId)
	},
	"TapChangerControl": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.TapChangerControl](db, ctx, modelId)
	},
	"Terminal": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.Terminal](db, ctx, modelId)
	},
	"ThermalGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.ThermalGeneratingUnit](db, ctx, modelId)
	},
	"TransformerEnd": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.TransformerEnd](db, ctx, modelId)
	},
	"VoltageLevel": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.VoltageLevel](db, ctx, modelId)
	},
	"VoltageLimit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.VoltageLimit](db, ctx, modelId)
	},
	"VsCapabilityCurve": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.VsCapabilityCurve](db, ctx, modelId)
	},
	"VsConverter": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
		return FindNameAndMrid[models.VsConverter](db, ctx, modelId)
	},
	"WindGeneratingUnit": func(ctx context.Context, db *bun.DB, modelId int) ([]models.MridNameGetter, error) {
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
