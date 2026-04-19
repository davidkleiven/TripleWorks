package migrations

import (
	"context"
	"fmt"
	"reflect"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

var indexTables = []any{
	&models.ACDCConverterDCTerminal{},
	&models.ACDCConverter{},
	&models.ACDCTerminal{},
	&models.ACLineSegment{},
	&models.AsynchronousMachine{},
	&models.BaseVoltage{},
	&models.BasicIntervalSchedule{},
	&models.Breaker{},
	&models.BusNameMarker{},
	&models.BusbarSection{},
	&models.ConductingEquipment{},
	&models.Conductor{},
	&models.ConformLoadGroup{},
	&models.ConformLoad{},
	&models.ConnectivityNodeContainer{},
	&models.ConnectivityNode{},
	&models.Connector{},
	&models.ControlAreaGeneratingUnit{},
	&models.ControlArea{},
	&models.CsConverter{},
	&models.CurrentLimit{},
	&models.Curve{},
	&models.DCBaseTerminal{},
	&models.DCBreaker{},
	&models.DCBusbar{},
	&models.DCChopper{},
	&models.DCConductingEquipment{},
	&models.DCConverterUnit{},
	&models.DCDisconnector{},
	&models.DCEquipmentContainer{},
	&models.DCGround{},
	&models.DCLineSegment{},
	&models.DCLine{},
	&models.DCNode{},
	&models.DCSeriesDevice{},
	&models.DCShunt{},
	&models.DCSwitch{},
	&models.DCTerminal{},
	&models.Disconnector{},
	&models.EnergyConsumer{},
	&models.EnergySchedulingType{},
	&models.EnergySource{},
	&models.EquipmentContainer{},
	&models.Equipment{},
	&models.EquivalentBranch{},
	&models.EquivalentEquipment{},
	&models.EquivalentInjection{},
	&models.EquivalentNetwork{},
	&models.EquivalentShunt{},
	&models.ExternalNetworkInjection{},
	&models.FossilFuel{},
	&models.GeneratingUnit{},
	&models.GeographicalRegion{},
	&models.HydroGeneratingUnit{},
	&models.HydroPowerPlant{},
	&models.HydroPump{},
	&models.IdentifiedObject{},
	&models.Junction{},
	&models.LinearShuntCompensator{},
	&models.Line{},
	&models.LoadBreakSwitch{},
	&models.LoadGroup{},
	&models.LoadResponseCharacteristic{},
	&models.Location{},
	&models.NonConformLoadGroup{},
	&models.NonConformLoad{},
	&models.NonlinearShuntCompensator{},
	&models.NuclearGeneratingUnit{},
	&models.OperationalLimitSet{},
	&models.OperationalLimitType{},
	&models.OperationalLimit{},
	&models.PhaseTapChangerAsymmetrical{},
	&models.PhaseTapChangerLinear{},
	&models.PhaseTapChangerNonLinear{},
	&models.PhaseTapChangerSymmetrical{},
	&models.PhaseTapChangerTable{},
	&models.PhaseTapChangerTabular{},
	&models.PhaseTapChanger{},
	&models.PowerSystemResource{},
	&models.PowerTransformerEnd{},
	&models.PowerTransformer{},
	&models.ProtectedSwitch{},
	&models.RatioTapChangerTable{},
	&models.RatioTapChanger{},
	&models.ReactiveCapabilityCurve{},
	&models.RegularIntervalSchedule{},
	&models.RegulatingCondEq{},
	&models.RegulatingControl{},
	&models.ReportingGroup{},
	&models.RotatingMachine{},
	&models.SeriesCompensator{},
	&models.ShuntCompensator{},
	&models.SolarGeneratingUnit{},
	&models.StaticVarCompensator{},
	&models.SubGeographicalRegion{},
	&models.Substation{},
	&models.Switch{},
	&models.SynchronousMachine{},
	&models.TapChangerControl{},
	&models.TapChanger{},
	&models.Terminal{},
	&models.ThermalGeneratingUnit{},
	&models.TransformerEnd{},
	&models.VoltageLevel{},
	&models.VoltageLimit{},
	&models.VsCapabilityCurve{},
	&models.VsConverter{},
	&models.WindGeneratingUnit{},
}

func init() {
	migrations.MustRegister(addMridCommitIdIndex, revertAddMridCommitIdIndex)
}

func addMridCommitIdIndex(ctx context.Context, db *bun.DB) error {
	for i, item := range viewTables {
		name := db.Table(reflect.TypeOf(item).Elem()).Name
		sql := fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_mrid_commit_id ON %s (mrid, commit_id)", name, name)
		_, err := db.ExecContext(ctx, sql)
		if err != nil {
			return fmt.Errorf("Failed for %d (%s): %w", i, name, err)
		}
	}
	return nil
}

func revertAddMridCommitIdIndex(ctx context.Context, db *bun.DB) error {
	for i, item := range viewTables {
		name := db.Table(reflect.TypeOf(item).Elem()).Name
		sql := fmt.Sprintf("DROP INDEX IF EXISTS idx_%s_mrid_commit_id", name)
		_, err := db.ExecContext(ctx, sql)
		if err != nil {
			return fmt.Errorf("Failed for %d (%s): %w", i, name, err)
		}
	}
	return nil
}
