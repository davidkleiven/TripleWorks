package migrations

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

var tablesWithLatestView2 = []any{
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
	migrations.MustRegister(addDeletedToView, revertDeletedToView)
}

func addDeletedToView(ctx context.Context, db *bun.DB) error {
	invalidLines := MustGetQuery("invalid_lines.sql")
	_, err := db.ExecContext(ctx, "DROP VIEW IF EXISTS v_invalid_lines")
	if err != nil {
		return fmt.Errorf("Could not drop view v_invalid_lines: %w", err)
	}

	for i, item := range tablesWithLatestView2 {
		name := db.Table(reflect.TypeOf(item).Elem()).Name
		vName := fmt.Sprintf("v_%s_latest", name)
		_, err := db.ExecContext(ctx, fmt.Sprintf("DROP VIEW IF EXISTS %s", vName))
		if err != nil {
			slog.Error("Failed to remove view", "name", vName, "error", err)
		}

		sql := MustGetViewSql(name)
		_, err = db.ExecContext(ctx, sql)
		if err != nil {
			return fmt.Errorf("Failed for %d (%s): %w", i, name, err)
		}
	}

	_, err = db.ExecContext(ctx, invalidLines)
	if err != nil {
		return fmt.Errorf("Could not re-create view invalid_lines")
	}
	return nil
}

func revertDeletedToView(ctx context.Context, db *bun.DB) error {
	return nil
}
