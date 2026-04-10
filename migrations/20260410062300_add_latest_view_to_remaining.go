package migrations

import (
	"context"
	"fmt"
	"reflect"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

var tablesWithLatestView = []any{
	&models.ACDCConverter{},
	&models.ACDCConverterDCTerminal{},
	&models.ACDCTerminal{},
	&models.ACLineSegment{},
	&models.AsynchronousMachine{},
	&models.BaseVoltage{},
	&models.Breaker{},
	&models.BusNameMarker{},
	&models.BusbarSection{},
	&models.ConductingEquipment{},
	&models.Conductor{},
	&models.ConformLoad{},
	&models.ConformLoadGroup{},
	&models.ConnectivityNode{},
	&models.ConnectivityNodeContainer{},
	&models.Connector{},
	&models.ControlArea{},
	&models.ControlAreaGeneratingUnit{},
	&models.CsConverter{},
	&models.CurrentLimit{},
	&models.DCBaseTerminal{},
	&models.DCBreaker{},
	&models.DCBusbar{},
	&models.DCChopper{},
	&models.DCConductingEquipment{},
	&models.DCConverterUnit{},
	&models.DCDisconnector{},
	&models.DCEquipmentContainer{},
	&models.DCGround{},
	&models.DCLine{},
	&models.DCLineSegment{},
	&models.DCNode{},
	&models.DCSeriesDevice{},
	&models.DCShunt{},
	&models.DCSwitch{},
	&models.DCTerminal{},
	&models.Disconnector{},
	&models.EnergyConsumer{},
	&models.EnergySchedulingType{},
	&models.EnergySource{},
	&models.Equipment{},
	&models.EquipmentContainer{},
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
	&models.Line{},
	&models.LinearShuntCompensator{},
	&models.LoadBreakSwitch{},
	&models.LoadGroup{},
	&models.NonConformLoad{},
	&models.NonConformLoadGroup{},
	&models.NonlinearShuntCompensator{},
	&models.NuclearGeneratingUnit{},
	&models.OperationalLimit{},
	&models.OperationalLimitSet{},
	&models.OperationalLimitType{},
	&models.PhaseTapChanger{},
	&models.PhaseTapChangerAsymmetrical{},
	&models.PhaseTapChangerLinear{},
	&models.PhaseTapChangerNonLinear{},
	&models.PhaseTapChangerSymmetrical{},
	&models.PowerSystemResource{},
	&models.PowerTransformer{},
	&models.PowerTransformerEnd{},
	&models.ProtectedSwitch{},
	&models.RatioTapChanger{},
	&models.RegulatingCondEq{},
	&models.RegulatingControl{},
	&models.ReportingGroup{},
	&models.RotatingMachine{},
	&models.SeriesCompensator{},
	&models.ShuntCompensator{},
	&models.SolarGeneratingUnit{},
	&models.StaticVarCompensator{},
	&models.SubGeographicalRegion{},
	&models.Switch{},
	&models.SynchronousMachine{},
	&models.TapChanger{},
	&models.TapChangerControl{},
	&models.Terminal{},
	&models.ThermalGeneratingUnit{},
	&models.TransformerEnd{},
	&models.VoltageLevel{},
	&models.VsConverter{},
	&models.WindGeneratingUnit{},
}

func init() {
	migrations.MustRegister(addLatestView, revertAddLatestView)
}

func addLatestView(ctx context.Context, db *bun.DB) error {
	for i, item := range tablesWithLatestView {
		name := db.Table(reflect.TypeOf(item).Elem()).Name
		sql := MustGetViewSql(name)
		_, err := db.ExecContext(ctx, sql)
		if err != nil {
			return fmt.Errorf("Failed for %d (%s): %w", i, name, err)
		}
	}
	return nil
}

func revertAddLatestView(ctx context.Context, db *bun.DB) error {
	for i, item := range tablesWithLatestView {
		name := db.Table(reflect.TypeOf(item).Elem()).Name
		sql := fmt.Sprintf("DROP VIEW v_%s_latest", name)
		_, err := db.ExecContext(ctx, sql)
		if err != nil {
			return fmt.Errorf("Failed for %d (%s): %w", i, name, err)
		}
	}
	return nil
}
