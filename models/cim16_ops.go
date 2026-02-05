package models

import (
	"github.com/google/uuid"
	"time"
)

type TerminalOperations struct {
	ConnectivityNodeMrid uuid.UUID `bun:"connectivity_node_mrid,type:uuid" json:"connectivity_node_mrid" iri:"cim:TerminalOperations.ConnectivityNode"`
	ConnectivityNode     *Entity   `bun:"rel:belongs-to,join:connectivity_node_mrid=mrid" json:"connectivity_node,omitempty"`
}

type LoadGroupOperations struct {
	SubLoadAreaMrid uuid.UUID `bun:"sub_load_area_mrid,type:uuid" json:"sub_load_area_mrid" iri:"cim:LoadGroupOperations.SubLoadArea"`
	SubLoadArea     *Entity   `bun:"rel:belongs-to,join:sub_load_area_mrid=mrid" json:"sub_load_area,omitempty"`
}

type AnalogLimit struct {
	Limit
	LimitSetMrid uuid.UUID `bun:"limit_set_mrid,type:uuid" json:"limit_set_mrid" iri:"cim:AnalogLimit.LimitSet"`
	LimitSet     *Entity   `bun:"rel:belongs-to,join:limit_set_mrid=mrid" json:"limit_set,omitempty"`
	Value        float64   `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AnalogLimit.value"`
}
type DiscreteValue struct {
	MeasurementValue
	Value        int       `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DiscreteValue.value"`
	DiscreteMrid uuid.UUID `bun:"discrete_mrid,type:uuid" json:"discrete_mrid" iri:"cim:DiscreteValue.Discrete"`
	Discrete     *Entity   `bun:"rel:belongs-to,join:discrete_mrid=mrid" json:"discrete,omitempty"`
}
type RaiseLowerCommand struct {
	AnalogControl
	ValueAliasSetMrid uuid.UUID `bun:"value_alias_set_mrid,type:uuid" json:"value_alias_set_mrid" iri:"cim:RaiseLowerCommand.ValueAliasSet"`
	ValueAliasSet     *Entity   `bun:"rel:belongs-to,join:value_alias_set_mrid=mrid" json:"value_alias_set,omitempty"`
}
type Discrete struct {
	Measurement
	ValueAliasSetMrid uuid.UUID `bun:"value_alias_set_mrid,type:uuid" json:"value_alias_set_mrid" iri:"cim:Discrete.ValueAliasSet"`
	ValueAliasSet     *Entity   `bun:"rel:belongs-to,join:value_alias_set_mrid=mrid" json:"value_alias_set,omitempty"`
}
type AccumulatorLimitSet struct {
	LimitSet
	MeasurementsMrid uuid.UUID `bun:"measurements_mrid,type:uuid" json:"measurements_mrid" iri:"cim:AccumulatorLimitSet.Measurements"`
	Measurements     *Entity   `bun:"rel:belongs-to,join:measurements_mrid=mrid" json:"measurements,omitempty"`
}
type SetPoint struct {
	AnalogControl
	Value       float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SetPoint.value"`
	NormalValue float64 `bun:"normal_value" json:"normal_value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SetPoint.normalValue"`
}
type GrossToNetActivePowerCurve struct {
	Curve
	GeneratingUnitMrid uuid.UUID `bun:"generating_unit_mrid,type:uuid" json:"generating_unit_mrid" iri:"cim:GrossToNetActivePowerCurve.GeneratingUnit"`
	GeneratingUnit     *Entity   `bun:"rel:belongs-to,join:generating_unit_mrid=mrid" json:"generating_unit,omitempty"`
}
type SubLoadArea struct {
	EnergyArea
	LoadAreaMrid uuid.UUID `bun:"load_area_mrid,type:uuid" json:"load_area_mrid" iri:"cim:SubLoadArea.LoadArea"`
	LoadArea     *Entity   `bun:"rel:belongs-to,join:load_area_mrid=mrid" json:"load_area,omitempty"`
}
type SwitchSchedule struct {
	SeasonDayTypeSchedule
	SwitchMrid uuid.UUID `bun:"switch_mrid,type:uuid" json:"switch_mrid" iri:"cim:SwitchSchedule.Switch"`
	Switch     *Entity   `bun:"rel:belongs-to,join:switch_mrid=mrid" json:"switch,omitempty"`
}
type AccumulatorLimit struct {
	Limit
	LimitSetMrid uuid.UUID `bun:"limit_set_mrid,type:uuid" json:"limit_set_mrid" iri:"cim:AccumulatorLimit.LimitSet"`
	LimitSet     *Entity   `bun:"rel:belongs-to,join:limit_set_mrid=mrid" json:"limit_set,omitempty"`
	Value        int       `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AccumulatorLimit.value"`
}
type ActivePowerLimit struct {
	OperationalLimit
	Value float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerLimit.value"`
}
type StringMeasurementValue struct {
	MeasurementValue
	StringMeasurement string `bun:"string_measurement" json:"string_measurement" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StringMeasurementValue.StringMeasurement"`
	Value             string `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StringMeasurementValue.value"`
}
type RegulationSchedule struct {
	SeasonDayTypeSchedule
	RegulatingControlMrid uuid.UUID `bun:"regulating_control_mrid,type:uuid" json:"regulating_control_mrid" iri:"cim:RegulationSchedule.RegulatingControl"`
	RegulatingControl     *Entity   `bun:"rel:belongs-to,join:regulating_control_mrid=mrid" json:"regulating_control,omitempty"`
}
type AccumulatorReset struct {
	Control
	AccumulatorValue int `bun:"accumulator_value" json:"accumulator_value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AccumulatorReset.AccumulatorValue"`
}
type RegularTimePoint struct {
	Value1           float64 `bun:"value1" json:"value1" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegularTimePoint.value1"`
	Value2           float64 `bun:"value2" json:"value2" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegularTimePoint.value2"`
	IntervalSchedule int     `bun:"interval_schedule" json:"interval_schedule" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegularTimePoint.IntervalSchedule"`
	SequenceNumber   int     `bun:"sequence_number" json:"sequence_number" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegularTimePoint.sequenceNumber"`
}
type AnalogControl struct {
	Control
	MinValue    float64 `bun:"min_value" json:"min_value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AnalogControl.minValue"`
	AnalogValue float64 `bun:"analog_value" json:"analog_value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AnalogControl.AnalogValue"`
	MaxValue    float64 `bun:"max_value" json:"max_value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AnalogControl.maxValue"`
}
type AnalogLimitSet struct {
	LimitSet
	MeasurementsMrid uuid.UUID `bun:"measurements_mrid,type:uuid" json:"measurements_mrid" iri:"cim:AnalogLimitSet.Measurements"`
	Measurements     *Entity   `bun:"rel:belongs-to,join:measurements_mrid=mrid" json:"measurements,omitempty"`
}
type AccumulatorValue struct {
	MeasurementValue
	AccumulatorMrid uuid.UUID `bun:"accumulator_mrid,type:uuid" json:"accumulator_mrid" iri:"cim:AccumulatorValue.Accumulator"`
	Accumulator     *Entity   `bun:"rel:belongs-to,join:accumulator_mrid=mrid" json:"accumulator,omitempty"`
	Value           int       `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AccumulatorValue.value"`
}
type Analog struct {
	Measurement
	PositiveFlowIn bool `bun:"positive_flow_in" json:"positive_flow_in" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Analog.positiveFlowIn"`
}
type LimitSet struct {
	IdentifiedObject
	IsPercentageLimits bool `bun:"is_percentage_limits" json:"is_percentage_limits" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LimitSet.isPercentageLimits"`
}

type Quality61850 struct {
	Test              bool      `bun:"test" json:"test" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.test"`
	SourceId          int       `bun:"source_id" json:"source_id"`
	Source            *Source   `bun:"rel:belongs-to,join:source_id=id" json:"source,omitempty"`
	EstimatorReplaced bool      `bun:"estimator_replaced" json:"estimator_replaced" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.estimatorReplaced"`
	Failure           bool      `bun:"failure" json:"failure" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.failure"`
	OutOfRange        bool      `bun:"out_of_range" json:"out_of_range" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.outOfRange"`
	Oscillatory       bool      `bun:"oscillatory" json:"oscillatory" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.oscillatory"`
	ValidityId        int       `bun:"validity_id" json:"validity_id"`
	Validity          *Validity `bun:"rel:belongs-to,join:validity_id=id" json:"validity,omitempty"`
	OperatorBlocked   bool      `bun:"operator_blocked" json:"operator_blocked" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.operatorBlocked"`
	Suspect           bool      `bun:"suspect" json:"suspect" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.suspect"`
	OldData           bool      `bun:"old_data" json:"old_data" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.oldData"`
	BadReference      bool      `bun:"bad_reference" json:"bad_reference" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.badReference"`
	OverFlow          bool      `bun:"over_flow" json:"over_flow" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Quality61850.overFlow"`
}
type Season struct {
	IdentifiedObject
	EndDateMrid   uuid.UUID `bun:"end_date_mrid,type:uuid" json:"end_date_mrid" iri:"cim:Season.EndDate"`
	EndDate       *Entity   `bun:"rel:belongs-to,join:end_date_mrid=mrid" json:"end_date,omitempty"`
	StartDateMrid uuid.UUID `bun:"start_date_mrid,type:uuid" json:"start_date_mrid" iri:"cim:Season.StartDate"`
	StartDate     *Entity   `bun:"rel:belongs-to,join:start_date_mrid=mrid" json:"start_date,omitempty"`
}
type Bay struct {
	EquipmentContainer
	VoltageLevelMrid uuid.UUID `bun:"voltage_level_mrid,type:uuid" json:"voltage_level_mrid" iri:"cim:Bay.VoltageLevel"`
	VoltageLevel     *Entity   `bun:"rel:belongs-to,join:voltage_level_mrid=mrid" json:"voltage_level,omitempty"`
}
type AnalogValue struct {
	MeasurementValue
	AnalogMrid uuid.UUID `bun:"analog_mrid,type:uuid" json:"analog_mrid" iri:"cim:AnalogValue.Analog"`
	Analog     *Entity   `bun:"rel:belongs-to,join:analog_mrid=mrid" json:"analog,omitempty"`
	Value      float64   `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AnalogValue.value"`
}
type Command struct {
	Control
	DiscreteValue     int       `bun:"discrete_value" json:"discrete_value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Command.DiscreteValue"`
	ValueAliasSetMrid uuid.UUID `bun:"value_alias_set_mrid,type:uuid" json:"value_alias_set_mrid" iri:"cim:Command.ValueAliasSet"`
	ValueAliasSet     *Entity   `bun:"rel:belongs-to,join:value_alias_set_mrid=mrid" json:"value_alias_set,omitempty"`
	NormalValue       int       `bun:"normal_value" json:"normal_value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Command.normalValue"`
	Value             int       `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Command.value"`
}
type ValueToAlias struct {
	IdentifiedObject
	Value             int       `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ValueToAlias.value"`
	ValueAliasSetMrid uuid.UUID `bun:"value_alias_set_mrid,type:uuid" json:"value_alias_set_mrid" iri:"cim:ValueToAlias.ValueAliasSet"`
	ValueAliasSet     *Entity   `bun:"rel:belongs-to,join:value_alias_set_mrid=mrid" json:"value_alias_set,omitempty"`
}
type Measurement struct {
	IdentifiedObject
	PowerSystemResourceMrid uuid.UUID       `bun:"power_system_resource_mrid,type:uuid" json:"power_system_resource_mrid" iri:"cim:Measurement.PowerSystemResource"`
	PowerSystemResource     *Entity         `bun:"rel:belongs-to,join:power_system_resource_mrid=mrid" json:"power_system_resource,omitempty"`
	UnitMultiplierId        int             `bun:"unit_multiplier_id" json:"unit_multiplier_id"`
	UnitMultiplier          *UnitMultiplier `bun:"rel:belongs-to,join:unit_multiplier_id=id" json:"unit_multiplier,omitempty"`
	TerminalMrid            uuid.UUID       `bun:"terminal_mrid,type:uuid" json:"terminal_mrid" iri:"cim:Measurement.Terminal"`
	Terminal                *Entity         `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
	UnitSymbolId            int             `bun:"unit_symbol_id" json:"unit_symbol_id"`
	UnitSymbol              *UnitSymbol     `bun:"rel:belongs-to,join:unit_symbol_id=id" json:"unit_symbol,omitempty"`
	PhasesId                int             `bun:"phases_id" json:"phases_id"`
	Phases                  *PhaseCode      `bun:"rel:belongs-to,join:phases_id=id" json:"phases,omitempty"`
	MeasurementType         string          `bun:"measurement_type" json:"measurement_type" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Measurement.measurementType"`
}
type TapSchedule struct {
	SeasonDayTypeSchedule
	TapChangerMrid uuid.UUID `bun:"tap_changer_mrid,type:uuid" json:"tap_changer_mrid" iri:"cim:TapSchedule.TapChanger"`
	TapChanger     *Entity   `bun:"rel:belongs-to,join:tap_changer_mrid=mrid" json:"tap_changer,omitempty"`
}
type ConnectivityNode struct {
	IdentifiedObject
	ConnectivityNodeContainerMrid uuid.UUID `bun:"connectivity_node_container_mrid,type:uuid" json:"connectivity_node_container_mrid" iri:"cim:ConnectivityNode.ConnectivityNodeContainer"`
	ConnectivityNodeContainer     *Entity   `bun:"rel:belongs-to,join:connectivity_node_container_mrid=mrid" json:"connectivity_node_container,omitempty"`
}
type MeasurementValue struct {
	IdentifiedObject
	MeasurementValueSourceMrid uuid.UUID `bun:"measurement_value_source_mrid,type:uuid" json:"measurement_value_source_mrid" iri:"cim:MeasurementValue.MeasurementValueSource"`
	MeasurementValueSource     *Entity   `bun:"rel:belongs-to,join:measurement_value_source_mrid=mrid" json:"measurement_value_source,omitempty"`
	SensorAccuracy             float64   `bun:"sensor_accuracy" json:"sensor_accuracy" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#MeasurementValue.sensorAccuracy"`
	TimeStamp                  time.Time `bun:"time_stamp" json:"time_stamp" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#MeasurementValue.timeStamp"`
}
type MeasurementValueQuality struct {
	Quality61850
	MeasurementValueMrid uuid.UUID `bun:"measurement_value_mrid,type:uuid" json:"measurement_value_mrid" iri:"cim:MeasurementValueQuality.MeasurementValue"`
	MeasurementValue     *Entity   `bun:"rel:belongs-to,join:measurement_value_mrid=mrid" json:"measurement_value,omitempty"`
}
type Control struct {
	IdentifiedObject
	UnitSymbolId            int             `bun:"unit_symbol_id" json:"unit_symbol_id"`
	UnitSymbol              *UnitSymbol     `bun:"rel:belongs-to,join:unit_symbol_id=id" json:"unit_symbol,omitempty"`
	TimeStamp               time.Time       `bun:"time_stamp" json:"time_stamp" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Control.timeStamp"`
	UnitMultiplierId        int             `bun:"unit_multiplier_id" json:"unit_multiplier_id"`
	UnitMultiplier          *UnitMultiplier `bun:"rel:belongs-to,join:unit_multiplier_id=id" json:"unit_multiplier,omitempty"`
	PowerSystemResourceMrid uuid.UUID       `bun:"power_system_resource_mrid,type:uuid" json:"power_system_resource_mrid" iri:"cim:Control.PowerSystemResource"`
	PowerSystemResource     *Entity         `bun:"rel:belongs-to,join:power_system_resource_mrid=mrid" json:"power_system_resource,omitempty"`
	OperationInProgress     bool            `bun:"operation_in_progress" json:"operation_in_progress" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Control.operationInProgress"`
	ControlType             string          `bun:"control_type" json:"control_type" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Control.controlType"`
}
type ApparentPowerLimit struct {
	OperationalLimit
	Value float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ApparentPowerLimit.value"`
}
type StringMeasurement struct {
	Measurement
}
type ValueAliasSet struct {
	IdentifiedObject
}
type GroundDisconnector struct {
	Switch
}
type StationSupply struct {
	EnergyConsumer
}
type LoadArea struct {
	EnergyArea
}
type EnergyArea struct {
	IdentifiedObject
}
type Ground struct {
	ConductingEquipment
}
type DayType struct {
	IdentifiedObject
}
type Limit struct {
	IdentifiedObject
}
type Accumulator struct {
	Measurement
}
type MeasurementValueSource struct {
	IdentifiedObject
}
