package migrations

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(createCim16Tables, revertCreateCim16Tables)
}

var tables = []any{
	// Enum tables
	(*models.ControlAreaTypeKind)(nil),
	(*models.Currency)(nil),
	(*models.CurveStyle)(nil),
	(*models.DCConverterOperatingModeKind)(nil),
	(*models.DCPolarityKind)(nil),
	(*models.FuelType)(nil),
	(*models.GeneratorControlSource)(nil),
	(*models.HydroEnergyConversionKind)(nil),
	(*models.HydroPlantStorageKind)(nil),
	(*models.LimitTypeKind)(nil),
	(*models.OperationalLimitDirectionKind)(nil),
	(*models.PetersenCoilModeKind)(nil),
	(*models.PhaseCode)(nil),
	(*models.RegulatingControlModeKind)(nil),
	(*models.SVCControlMode)(nil),
	(*models.ShortCircuitRotorKind)(nil),
	(*models.Source)(nil),
	(*models.SynchronousMachineKind)(nil),
	(*models.TransformerControlMode)(nil),
	(*models.UnitMultiplier)(nil),
	(*models.UnitSymbol)(nil),
	(*models.Validity)(nil),
	(*models.WindGenUnitKind)(nil),
	(*models.WindingConnection)(nil),

	// Data tables
	(*models.Entity)(nil),
	(*models.ACDCConverterDCTerminal)(nil),
	(*models.ACDCConverter)(nil),
	(*models.ACDCTerminal)(nil),
	(*models.ACLineSegment)(nil),
	(*models.AsynchronousMachine)(nil),
	(*models.BaseVoltage)(nil),
	(*models.BasicIntervalSchedule)(nil),
	(*models.Breaker)(nil),
	(*models.BusNameMarker)(nil),
	(*models.BusbarSection)(nil),
	(*models.ConductingEquipment)(nil),
	(*models.Conductor)(nil),
	(*models.ConformLoadGroup)(nil),
	(*models.ConformLoadSchedule)(nil),
	(*models.ConformLoad)(nil),
	(*models.ConnectivityNodeContainer)(nil),
	(*models.Connector)(nil),
	(*models.ControlAreaGeneratingUnit)(nil),
	(*models.ControlArea)(nil),
	(*models.CsConverter)(nil),
	(*models.CurveData)(nil),
	(*models.Curve)(nil),
	(*models.DCBaseTerminal)(nil),
	(*models.DCBreaker)(nil),
	(*models.DCBusbar)(nil),
	(*models.DCChopper)(nil),
	(*models.DCConductingEquipment)(nil),
	(*models.DCConverterUnit)(nil),
	(*models.DCDisconnector)(nil),
	(*models.DCEquipmentContainer)(nil),
	(*models.DCGround)(nil),
	(*models.DCLineSegment)(nil),
	(*models.DCLine)(nil),
	(*models.DCNode)(nil),
	(*models.DCSeriesDevice)(nil),
	(*models.DCShunt)(nil),
	(*models.DCSwitch)(nil),
	(*models.DCTerminal)(nil),
	(*models.Disconnector)(nil),
	(*models.EnergyConsumer)(nil),
	(*models.EnergySchedulingType)(nil),
	(*models.EnergySource)(nil),
	(*models.EquipmentContainer)(nil),
	(*models.EquipmentVersion)(nil),
	(*models.Equipment)(nil),
	(*models.EquivalentBranch)(nil),
	(*models.EquivalentEquipment)(nil),
	(*models.EquivalentInjection)(nil),
	(*models.EquivalentNetwork)(nil),
	(*models.EquivalentShunt)(nil),
	(*models.ExternalNetworkInjection)(nil),
	(*models.FossilFuel)(nil),
	(*models.GeneratingUnit)(nil),
	(*models.GeographicalRegion)(nil),
	(*models.HydroGeneratingUnit)(nil),
	(*models.HydroPowerPlant)(nil),
	(*models.HydroPump)(nil),
	(*models.IdentifiedObject)(nil),
	(*models.InductancePerLength)(nil),
	(*models.Inductance)(nil),
	(*models.Junction)(nil),
	(*models.Length)(nil),
	(*models.LinearShuntCompensator)(nil),
	(*models.Line)(nil),
	(*models.LoadBreakSwitch)(nil),
	(*models.LoadGroup)(nil),
	(*models.LoadResponseCharacteristic)(nil),
	(*models.NonConformLoadGroup)(nil),
	(*models.NonConformLoadSchedule)(nil),
	(*models.NonConformLoad)(nil),
	(*models.NonlinearShuntCompensatorPoint)(nil),
	(*models.NonlinearShuntCompensator)(nil),
	(*models.NuclearGeneratingUnit)(nil),
	(*models.OperationalLimitSet)(nil),
	(*models.OperationalLimitType)(nil),
	(*models.OperationalLimit)(nil),
	(*models.PhaseTapChangerAsymmetrical)(nil),
	(*models.PhaseTapChangerLinear)(nil),
	(*models.PhaseTapChangerNonLinear)(nil),
	(*models.PhaseTapChangerSymmetrical)(nil),
	(*models.PhaseTapChangerTablePoint)(nil),
	(*models.PhaseTapChangerTable)(nil),
	(*models.PhaseTapChangerTabular)(nil),
	(*models.PhaseTapChanger)(nil),
	(*models.PowerSystemResource)(nil),
	(*models.PowerTransformerEnd)(nil),
	(*models.PowerTransformer)(nil),
	(*models.ProtectedSwitch)(nil),
	(*models.RatioTapChangerTablePoint)(nil),
	(*models.RatioTapChangerTable)(nil),
	(*models.RatioTapChanger)(nil),
	(*models.ReactiveCapabilityCurve)(nil),
	(*models.ReactivePower)(nil),
	(*models.RegularIntervalSchedule)(nil),
	(*models.RegulatingCondEq)(nil),
	(*models.RegulatingControl)(nil),
	(*models.ReportingGroup)(nil),
	(*models.RotatingMachine)(nil),
	(*models.SeriesCompensator)(nil),
	(*models.ShuntCompensator)(nil),
	(*models.SolarGeneratingUnit)(nil),
	(*models.StaticVarCompensator)(nil),
	(*models.SubGeographicalRegion)(nil),
	(*models.Substation)(nil),
	(*models.Switch)(nil),
	(*models.SynchronousMachine)(nil),
	(*models.TapChangerControl)(nil),
	(*models.TapChangerTablePoint)(nil),
	(*models.TapChanger)(nil),
	(*models.Terminal)(nil),
	(*models.ThermalGeneratingUnit)(nil),
	(*models.TieFlow)(nil),
	(*models.TransformerEnd)(nil),
	(*models.VoltageLevel)(nil),
	(*models.VoltageLimit)(nil),
	(*models.VsCapabilityCurve)(nil),
	(*models.VsConverter)(nil),
	(*models.WindGeneratingUnit)(nil),
}

func createCim16Tables(ctx context.Context, db *bun.DB) error {
	for i, table := range tables {
		_, err := db.NewCreateTable().
			Model(table).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to create table no. %d: %w", i, err)
		}
	}
	return nil
}

func revertCreateCim16Tables(ctx context.Context, db *bun.DB) error {
	for i := len(tables); i > 0; i-- {
		_, err := db.NewDropTable().Model(tables[i-1]).IfExists().Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to drop table no. %d: %w", i, err)
		}
	}
	return nil
}
