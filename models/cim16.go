package models

import (
	"github.com/google/uuid"
	"time"
)

type Entity struct {
	ModelEntity
	Mrid uuid.UUID `bun:"mrid,type:uuid,pk"`
}
type DCBaseTerminal struct {
	ACDCTerminal
	DCNodeMrid uuid.UUID `bun:"dcnode_mrid,type:uuid" json:"dcnode_mrid"`
	DCNode     *Entity   `bun:"rel:belongs-to,join:dcnode_mrid=mrid" json:"dcnode,omitempty"`
}
type PhaseTapChangerNonLinear struct {
	PhaseTapChanger
	XMax                 float64 `bun:"x_max" json:"x_max" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerNonLinear.xMax"`
	XMin                 float64 `bun:"x_min" json:"x_min" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerNonLinear.xMin"`
	VoltageStepIncrement float64 `bun:"voltage_step_increment" json:"voltage_step_increment" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerNonLinear.voltageStepIncrement"`
}
type Length struct {
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Length.value"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type ControlArea struct {
	PowerSystemResource
	TypeId int                  `bun:"type_id" json:"type_id"`
	Type   *ControlAreaTypeKind `bun:"rel:belongs-to,join:type_id=id" json:"type,omitempty"`
}
type BasicIntervalSchedule struct {
	IdentifiedObject
	Value1UnitId int         `bun:"value1_unit_id" json:"value1_unit_id"`
	Value1Unit   *UnitSymbol `bun:"rel:belongs-to,join:value1_unit_id=id" json:"value1_unit,omitempty"`
	StartTime    time.Time   `bun:"start_time" json:"start_time" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#BasicIntervalSchedule.startTime"`
	Value2UnitId int         `bun:"value2_unit_id" json:"value2_unit_id"`
	Value2Unit   *UnitSymbol `bun:"rel:belongs-to,join:value2_unit_id=id" json:"value2_unit,omitempty"`
}
type TapChangerTablePoint struct {
	R     float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.r"`
	Step  int     `bun:"step" json:"step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.step"`
	Ratio float64 `bun:"ratio" json:"ratio" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.ratio"`
	B     float64 `bun:"b" json:"b" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.b"`
	G     float64 `bun:"g" json:"g" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.g"`
	X     float64 `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.x"`
}
type VoltagePerReactivePower struct {
	Value                   float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltagePerReactivePower.value"`
	DenominatorUnitId       int             `bun:"denominator_unit_id" json:"denominator_unit_id"`
	DenominatorUnit         *UnitSymbol     `bun:"rel:belongs-to,join:denominator_unit_id=id" json:"denominator_unit,omitempty"`
	MultiplierId            int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier              *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId                  int             `bun:"unit_id" json:"unit_id"`
	Unit                    *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	DenominatorMultiplierId int             `bun:"denominator_multiplier_id" json:"denominator_multiplier_id"`
	DenominatorMultiplier   *UnitMultiplier `bun:"rel:belongs-to,join:denominator_multiplier_id=id" json:"denominator_multiplier,omitempty"`
}
type Frequency struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Frequency.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type Terminal struct {
	ACDCTerminal
	ConductingEquipmentMrid uuid.UUID  `bun:"conducting_equipment_mrid,type:uuid" json:"conducting_equipment_mrid"`
	ConductingEquipment     *Entity    `bun:"rel:belongs-to,join:conducting_equipment_mrid=mrid" json:"conducting_equipment,omitempty"`
	PhasesId                int        `bun:"phases_id" json:"phases_id"`
	Phases                  *PhaseCode `bun:"rel:belongs-to,join:phases_id=id" json:"phases,omitempty"`
}
type Seconds struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Seconds.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type Conductor struct {
	ConductingEquipment
	Length float64 `bun:"length" json:"length" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Conductor.length"`
}
type ResistancePerLength struct {
	DenominatorUnitId       int             `bun:"denominator_unit_id" json:"denominator_unit_id"`
	DenominatorUnit         *UnitSymbol     `bun:"rel:belongs-to,join:denominator_unit_id=id" json:"denominator_unit,omitempty"`
	Value                   float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ResistancePerLength.value"`
	MultiplierId            int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier              *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	DenominatorMultiplierId int             `bun:"denominator_multiplier_id" json:"denominator_multiplier_id"`
	DenominatorMultiplier   *UnitMultiplier `bun:"rel:belongs-to,join:denominator_multiplier_id=id" json:"denominator_multiplier,omitempty"`
	UnitId                  int             `bun:"unit_id" json:"unit_id"`
	Unit                    *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type BaseVoltage struct {
	IdentifiedObject
	NominalVoltage float64 `bun:"nominal_voltage" json:"nominal_voltage" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#BaseVoltage.nominalVoltage"`
}
type Substation struct {
	EquipmentContainer
	RegionMrid uuid.UUID `bun:"region_mrid,type:uuid" json:"region_mrid"`
	Region     *Entity   `bun:"rel:belongs-to,join:region_mrid=mrid" json:"region,omitempty"`
}
type DCLine struct {
	DCEquipmentContainer
	RegionMrid uuid.UUID `bun:"region_mrid,type:uuid" json:"region_mrid"`
	Region     *Entity   `bun:"rel:belongs-to,join:region_mrid=mrid" json:"region,omitempty"`
}
type PhaseTapChanger struct {
	TapChanger
	TransformerEndMrid uuid.UUID `bun:"transformer_end_mrid,type:uuid" json:"transformer_end_mrid"`
	TransformerEnd     *Entity   `bun:"rel:belongs-to,join:transformer_end_mrid=mrid" json:"transformer_end,omitempty"`
}
type VoltageLimit struct {
	OperationalLimit
	Value float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltageLimit.value"`
}
type Reactance struct {
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Reactance.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
}
type LoadResponseCharacteristic struct {
	IdentifiedObject
	PConstantImpedance float64 `bun:"pconstant_impedance" json:"pconstant_impedance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pConstantImpedance"`
	QVoltageExponent   float64 `bun:"qvoltage_exponent" json:"qvoltage_exponent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qVoltageExponent"`
	ExponentModel      bool    `bun:"exponent_model" json:"exponent_model" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.exponentModel"`
	QConstantCurrent   float64 `bun:"qconstant_current" json:"qconstant_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qConstantCurrent"`
	PConstantCurrent   float64 `bun:"pconstant_current" json:"pconstant_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pConstantCurrent"`
	QConstantPower     float64 `bun:"qconstant_power" json:"qconstant_power" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qConstantPower"`
	PVoltageExponent   float64 `bun:"pvoltage_exponent" json:"pvoltage_exponent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pVoltageExponent"`
	QConstantImpedance float64 `bun:"qconstant_impedance" json:"qconstant_impedance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qConstantImpedance"`
	PConstantPower     float64 `bun:"pconstant_power" json:"pconstant_power" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pConstantPower"`
	PFrequencyExponent float64 `bun:"pfrequency_exponent" json:"pfrequency_exponent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pFrequencyExponent"`
	QFrequencyExponent float64 `bun:"qfrequency_exponent" json:"qfrequency_exponent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qFrequencyExponent"`
}
type InductancePerLength struct {
	UnitId                  int             `bun:"unit_id" json:"unit_id"`
	Unit                    *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	MultiplierId            int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier              *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	Value                   float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#InductancePerLength.value"`
	DenominatorUnitId       int             `bun:"denominator_unit_id" json:"denominator_unit_id"`
	DenominatorUnit         *UnitSymbol     `bun:"rel:belongs-to,join:denominator_unit_id=id" json:"denominator_unit,omitempty"`
	DenominatorMultiplierId int             `bun:"denominator_multiplier_id" json:"denominator_multiplier_id"`
	DenominatorMultiplier   *UnitMultiplier `bun:"rel:belongs-to,join:denominator_multiplier_id=id" json:"denominator_multiplier,omitempty"`
}
type RatioTapChanger struct {
	TapChanger
	StepVoltageIncrement     float64                 `bun:"step_voltage_increment" json:"step_voltage_increment" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RatioTapChanger.stepVoltageIncrement"`
	TculControlModeId        int                     `bun:"tcul_control_mode_id" json:"tcul_control_mode_id"`
	TculControlMode          *TransformerControlMode `bun:"rel:belongs-to,join:tcul_control_mode_id=id" json:"tcul_control_mode,omitempty"`
	TransformerEndMrid       uuid.UUID               `bun:"transformer_end_mrid,type:uuid" json:"transformer_end_mrid"`
	TransformerEnd           *Entity                 `bun:"rel:belongs-to,join:transformer_end_mrid=mrid" json:"transformer_end,omitempty"`
	RatioTapChangerTableMrid uuid.UUID               `bun:"ratio_tap_changer_table_mrid,type:uuid" json:"ratio_tap_changer_table_mrid"`
	RatioTapChangerTable     *Entity                 `bun:"rel:belongs-to,join:ratio_tap_changer_table_mrid=mrid" json:"ratio_tap_changer_table,omitempty"`
}
type HydroPump struct {
	Equipment
	HydroPowerPlantMrid uuid.UUID `bun:"hydro_power_plant_mrid,type:uuid" json:"hydro_power_plant_mrid"`
	HydroPowerPlant     *Entity   `bun:"rel:belongs-to,join:hydro_power_plant_mrid=mrid" json:"hydro_power_plant,omitempty"`
	RotatingMachineMrid uuid.UUID `bun:"rotating_machine_mrid,type:uuid" json:"rotating_machine_mrid"`
	RotatingMachine     *Entity   `bun:"rel:belongs-to,join:rotating_machine_mrid=mrid" json:"rotating_machine,omitempty"`
}
type Capacitance struct {
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Capacitance.value"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type PowerTransformerEnd struct {
	TransformerEnd
	X                    float64            `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.x"`
	ConnectionKindId     int                `bun:"connection_kind_id" json:"connection_kind_id"`
	ConnectionKind       *WindingConnection `bun:"rel:belongs-to,join:connection_kind_id=id" json:"connection_kind,omitempty"`
	G                    float64            `bun:"g" json:"g" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.g"`
	RatedS               float64            `bun:"rated_s" json:"rated_s" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.ratedS"`
	PowerTransformerMrid uuid.UUID          `bun:"power_transformer_mrid,type:uuid" json:"power_transformer_mrid"`
	PowerTransformer     *Entity            `bun:"rel:belongs-to,join:power_transformer_mrid=mrid" json:"power_transformer,omitempty"`
	R                    float64            `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.r"`
	RatedU               float64            `bun:"rated_u" json:"rated_u" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.ratedU"`
	B                    float64            `bun:"b" json:"b" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.b"`
}
type Line struct {
	EquipmentContainer
	RegionMrid uuid.UUID `bun:"region_mrid,type:uuid" json:"region_mrid"`
	Region     *Entity   `bun:"rel:belongs-to,join:region_mrid=mrid" json:"region,omitempty"`
}
type EnergySource struct {
	ConductingEquipment
	X                        float64   `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.x"`
	EnergySchedulingTypeMrid uuid.UUID `bun:"energy_scheduling_type_mrid,type:uuid" json:"energy_scheduling_type_mrid"`
	EnergySchedulingType     *Entity   `bun:"rel:belongs-to,join:energy_scheduling_type_mrid=mrid" json:"energy_scheduling_type,omitempty"`
	R                        float64   `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.r"`
	Xn                       float64   `bun:"xn" json:"xn" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.xn"`
	NominalVoltage           float64   `bun:"nominal_voltage" json:"nominal_voltage" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.nominalVoltage"`
	Rn                       float64   `bun:"rn" json:"rn" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.rn"`
	VoltageAngle             float64   `bun:"voltage_angle" json:"voltage_angle" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.voltageAngle"`
	VoltageMagnitude         float64   `bun:"voltage_magnitude" json:"voltage_magnitude" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.voltageMagnitude"`
	R0                       float64   `bun:"r0" json:"r0" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.r0"`
	X0                       float64   `bun:"x0" json:"x0" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.x0"`
}
type NonlinearShuntCompensatorPoint struct {
	G                             float64   `bun:"g" json:"g" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#NonlinearShuntCompensatorPoint.g"`
	B                             float64   `bun:"b" json:"b" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#NonlinearShuntCompensatorPoint.b"`
	NonlinearShuntCompensatorMrid uuid.UUID `bun:"nonlinear_shunt_compensator_mrid,type:uuid" json:"nonlinear_shunt_compensator_mrid"`
	NonlinearShuntCompensator     *Entity   `bun:"rel:belongs-to,join:nonlinear_shunt_compensator_mrid=mrid" json:"nonlinear_shunt_compensator,omitempty"`
	SectionNumber                 int       `bun:"section_number" json:"section_number" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#NonlinearShuntCompensatorPoint.sectionNumber"`
}
type VoltageLevel struct {
	EquipmentContainer
	HighVoltageLimit float64   `bun:"high_voltage_limit" json:"high_voltage_limit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltageLevel.highVoltageLimit"`
	SubstationMrid   uuid.UUID `bun:"substation_mrid,type:uuid" json:"substation_mrid"`
	Substation       *Entity   `bun:"rel:belongs-to,join:substation_mrid=mrid" json:"substation,omitempty"`
	BaseVoltageMrid  uuid.UUID `bun:"base_voltage_mrid,type:uuid" json:"base_voltage_mrid"`
	BaseVoltage      *Entity   `bun:"rel:belongs-to,join:base_voltage_mrid=mrid" json:"base_voltage,omitempty"`
	LowVoltageLimit  float64   `bun:"low_voltage_limit" json:"low_voltage_limit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltageLevel.lowVoltageLimit"`
}
type FossilFuel struct {
	IdentifiedObject
	FossilFuelTypeId          int       `bun:"fossil_fuel_type_id" json:"fossil_fuel_type_id"`
	FossilFuelType            *FuelType `bun:"rel:belongs-to,join:fossil_fuel_type_id=id" json:"fossil_fuel_type,omitempty"`
	ThermalGeneratingUnitMrid uuid.UUID `bun:"thermal_generating_unit_mrid,type:uuid" json:"thermal_generating_unit_mrid"`
	ThermalGeneratingUnit     *Entity   `bun:"rel:belongs-to,join:thermal_generating_unit_mrid=mrid" json:"thermal_generating_unit,omitempty"`
}
type DCShunt struct {
	DCConductingEquipment
	Capacitance float64 `bun:"capacitance" json:"capacitance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCShunt.capacitance"`
	RatedUdc    float64 `bun:"rated_udc" json:"rated_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCShunt.ratedUdc"`
	Resistance  float64 `bun:"resistance" json:"resistance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCShunt.resistance"`
}
type RegularIntervalSchedule struct {
	BasicIntervalSchedule
	TimeStep float64   `bun:"time_step" json:"time_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegularIntervalSchedule.timeStep"`
	EndTime  time.Time `bun:"end_time" json:"end_time" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegularIntervalSchedule.endTime"`
}
type LinearShuntCompensator struct {
	ShuntCompensator
	GPerSection float64 `bun:"gper_section" json:"gper_section" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LinearShuntCompensator.gPerSection"`
	BPerSection float64 `bun:"bper_section" json:"bper_section" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LinearShuntCompensator.bPerSection"`
}
type ACDCConverter struct {
	ConductingEquipment
	ValveU0         float64   `bun:"valve_u0" json:"valve_u0" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.valveU0"`
	RatedUdc        float64   `bun:"rated_udc" json:"rated_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.ratedUdc"`
	MinUdc          float64   `bun:"min_udc" json:"min_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.minUdc"`
	SwitchingLoss   float64   `bun:"switching_loss" json:"switching_loss" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.switchingLoss"`
	BaseS           float64   `bun:"base_s" json:"base_s" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.baseS"`
	IdleLoss        float64   `bun:"idle_loss" json:"idle_loss" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.idleLoss"`
	MaxUdc          float64   `bun:"max_udc" json:"max_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.maxUdc"`
	ResistiveLoss   float64   `bun:"resistive_loss" json:"resistive_loss" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.resistiveLoss"`
	NumberOfValves  int       `bun:"number_of_valves" json:"number_of_valves" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.numberOfValves"`
	PccTerminalMrid uuid.UUID `bun:"pcc_terminal_mrid,type:uuid" json:"pcc_terminal_mrid"`
	PccTerminal     *Entity   `bun:"rel:belongs-to,join:pcc_terminal_mrid=mrid" json:"pcc_terminal,omitempty"`
}
type DCConverterUnit struct {
	DCEquipmentContainer
	OperationModeId int                           `bun:"operation_mode_id" json:"operation_mode_id"`
	OperationMode   *DCConverterOperatingModeKind `bun:"rel:belongs-to,join:operation_mode_id=id" json:"operation_mode,omitempty"`
	SubstationMrid  uuid.UUID                     `bun:"substation_mrid,type:uuid" json:"substation_mrid"`
	Substation      *Entity                       `bun:"rel:belongs-to,join:substation_mrid=mrid" json:"substation,omitempty"`
}
type PU struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PU.value"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
}
type IdentifiedObject struct {
	BaseEntity
	Mrid               string `bun:"mrid" json:"mrid" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.mRID"`
	ShortName          string `bun:"short_name" json:"short_name" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#IdentifiedObject.shortName"`
	Description        string `bun:"description" json:"description" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.description"`
	Name               string `bun:"name" json:"name" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.name"`
	EnergyIdentCodeEic string `bun:"energy_ident_code_eic" json:"energy_ident_code_eic" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#IdentifiedObject.energyIdentCodeEic"`
}
type Resistance struct {
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Resistance.value"`
}
type SeriesCompensator struct {
	ConductingEquipment
	VaristorVoltageThreshold float64 `bun:"varistor_voltage_threshold" json:"varistor_voltage_threshold" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.varistorVoltageThreshold"`
	X                        float64 `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.x"`
	R                        float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.r"`
	VaristorRatedCurrent     float64 `bun:"varistor_rated_current" json:"varistor_rated_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.varistorRatedCurrent"`
	VaristorPresent          bool    `bun:"varistor_present" json:"varistor_present" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.varistorPresent"`
}
type TapChanger struct {
	PowerSystemResource
	NeutralStep           int       `bun:"neutral_step" json:"neutral_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.neutralStep"`
	HighStep              int       `bun:"high_step" json:"high_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.highStep"`
	LtcFlag               bool      `bun:"ltc_flag" json:"ltc_flag" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.ltcFlag"`
	NeutralU              float64   `bun:"neutral_u" json:"neutral_u" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.neutralU"`
	LowStep               int       `bun:"low_step" json:"low_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.lowStep"`
	NormalStep            int       `bun:"normal_step" json:"normal_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.normalStep"`
	TapChangerControlMrid uuid.UUID `bun:"tap_changer_control_mrid,type:uuid" json:"tap_changer_control_mrid"`
	TapChangerControl     *Entity   `bun:"rel:belongs-to,join:tap_changer_control_mrid=mrid" json:"tap_changer_control,omitempty"`
}
type CsConverter struct {
	ACDCConverter
	RatedIdc float64 `bun:"rated_idc" json:"rated_idc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.ratedIdc"`
	MaxIdc   float64 `bun:"max_idc" json:"max_idc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.maxIdc"`
	MaxAlpha float64 `bun:"max_alpha" json:"max_alpha" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.maxAlpha"`
	MinAlpha float64 `bun:"min_alpha" json:"min_alpha" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.minAlpha"`
	MaxGamma float64 `bun:"max_gamma" json:"max_gamma" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.maxGamma"`
	MinIdc   float64 `bun:"min_idc" json:"min_idc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.minIdc"`
	MinGamma float64 `bun:"min_gamma" json:"min_gamma" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.minGamma"`
}
type Voltage struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Voltage.value"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
}
type PhaseTapChangerLinear struct {
	PhaseTapChanger
	StepPhaseShiftIncrement float64 `bun:"step_phase_shift_increment" json:"step_phase_shift_increment" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerLinear.stepPhaseShiftIncrement"`
	XMin                    float64 `bun:"x_min" json:"x_min" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerLinear.xMin"`
	XMax                    float64 `bun:"x_max" json:"x_max" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerLinear.xMax"`
}
type OperationalLimitType struct {
	IdentifiedObject
	AcceptableDuration float64                        `bun:"acceptable_duration" json:"acceptable_duration" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#OperationalLimitType.acceptableDuration"`
	LimitTypeId        int                            `bun:"limit_type_id" json:"limit_type_id"`
	LimitType          *LimitTypeKind                 `bun:"rel:belongs-to,join:limit_type_id=id" json:"limit_type,omitempty"`
	DirectionId        int                            `bun:"direction_id" json:"direction_id"`
	Direction          *OperationalLimitDirectionKind `bun:"rel:belongs-to,join:direction_id=id" json:"direction,omitempty"`
}
type TieFlow struct {
	ControlAreaMrid uuid.UUID `bun:"control_area_mrid,type:uuid" json:"control_area_mrid"`
	ControlArea     *Entity   `bun:"rel:belongs-to,join:control_area_mrid=mrid" json:"control_area,omitempty"`
	PositiveFlowIn  bool      `bun:"positive_flow_in" json:"positive_flow_in" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TieFlow.positiveFlowIn"`
	TerminalMrid    uuid.UUID `bun:"terminal_mrid,type:uuid" json:"terminal_mrid"`
	Terminal        *Entity   `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
}
type ShuntCompensator struct {
	RegulatingCondEq
	NomU               float64   `bun:"nom_u" json:"nom_u" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.nomU"`
	MaximumSections    int       `bun:"maximum_sections" json:"maximum_sections" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.maximumSections"`
	SwitchOnDate       time.Time `bun:"switch_on_date" json:"switch_on_date" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.switchOnDate"`
	SwitchOnCount      int       `bun:"switch_on_count" json:"switch_on_count" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.switchOnCount"`
	VoltageSensitivity float64   `bun:"voltage_sensitivity" json:"voltage_sensitivity" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.voltageSensitivity"`
	Grounded           bool      `bun:"grounded" json:"grounded" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.grounded"`
	AVRDelay           float64   `bun:"avrdelay" json:"avrdelay" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.aVRDelay"`
	NormalSections     int       `bun:"normal_sections" json:"normal_sections" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.normalSections"`
}
type TransformerEnd struct {
	IdentifiedObject
	EndNumber       int       `bun:"end_number" json:"end_number" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TransformerEnd.endNumber"`
	BaseVoltageMrid uuid.UUID `bun:"base_voltage_mrid,type:uuid" json:"base_voltage_mrid"`
	BaseVoltage     *Entity   `bun:"rel:belongs-to,join:base_voltage_mrid=mrid" json:"base_voltage,omitempty"`
	TerminalMrid    uuid.UUID `bun:"terminal_mrid,type:uuid" json:"terminal_mrid"`
	Terminal        *Entity   `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
}
type AsynchronousMachine struct {
	RotatingMachine
	NominalSpeed     float64 `bun:"nominal_speed" json:"nominal_speed" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AsynchronousMachine.nominalSpeed"`
	NominalFrequency float64 `bun:"nominal_frequency" json:"nominal_frequency" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AsynchronousMachine.nominalFrequency"`
}
type DCTerminal struct {
	DCBaseTerminal
	DCConductingEquipmentMrid uuid.UUID `bun:"dcconducting_equipment_mrid,type:uuid" json:"dcconducting_equipment_mrid"`
	DCConductingEquipment     *Entity   `bun:"rel:belongs-to,join:dcconducting_equipment_mrid=mrid" json:"dcconducting_equipment,omitempty"`
}
type Temperature struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Temperature.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type ApparentPower struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ApparentPower.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type SubGeographicalRegion struct {
	IdentifiedObject
	RegionMrid uuid.UUID `bun:"region_mrid,type:uuid" json:"region_mrid"`
	Region     *Entity   `bun:"rel:belongs-to,join:region_mrid=mrid" json:"region,omitempty"`
}
type CurrentFlow struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurrentFlow.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type Inductance struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Inductance.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type CurveData struct {
	Xvalue    float64   `bun:"xvalue" json:"xvalue" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurveData.xvalue"`
	CurveMrid uuid.UUID `bun:"curve_mrid,type:uuid" json:"curve_mrid"`
	Curve     *Entity   `bun:"rel:belongs-to,join:curve_mrid=mrid" json:"curve,omitempty"`
	Y2value   float64   `bun:"y2value" json:"y2value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurveData.y2value"`
	Y1value   float64   `bun:"y1value" json:"y1value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurveData.y1value"`
}
type PhaseTapChangerTabular struct {
	PhaseTapChanger
	PhaseTapChangerTableMrid uuid.UUID `bun:"phase_tap_changer_table_mrid,type:uuid" json:"phase_tap_changer_table_mrid"`
	PhaseTapChangerTable     *Entity   `bun:"rel:belongs-to,join:phase_tap_changer_table_mrid=mrid" json:"phase_tap_changer_table,omitempty"`
}
type RatioTapChangerTablePoint struct {
	TapChangerTablePoint
	RatioTapChangerTableMrid uuid.UUID `bun:"ratio_tap_changer_table_mrid,type:uuid" json:"ratio_tap_changer_table_mrid"`
	RatioTapChangerTable     *Entity   `bun:"rel:belongs-to,join:ratio_tap_changer_table_mrid=mrid" json:"ratio_tap_changer_table,omitempty"`
}
type BusNameMarker struct {
	IdentifiedObject
	ReportingGroupMrid uuid.UUID `bun:"reporting_group_mrid,type:uuid" json:"reporting_group_mrid"`
	ReportingGroup     *Entity   `bun:"rel:belongs-to,join:reporting_group_mrid=mrid" json:"reporting_group,omitempty"`
	Priority           int       `bun:"priority" json:"priority" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#BusNameMarker.priority"`
}
type AngleRadians struct {
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AngleRadians.value"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type OperationalLimit struct {
	IdentifiedObject
	OperationalLimitSetMrid  uuid.UUID `bun:"operational_limit_set_mrid,type:uuid" json:"operational_limit_set_mrid"`
	OperationalLimitSet      *Entity   `bun:"rel:belongs-to,join:operational_limit_set_mrid=mrid" json:"operational_limit_set,omitempty"`
	OperationalLimitTypeMrid uuid.UUID `bun:"operational_limit_type_mrid,type:uuid" json:"operational_limit_type_mrid"`
	OperationalLimitType     *Entity   `bun:"rel:belongs-to,join:operational_limit_type_mrid=mrid" json:"operational_limit_type,omitempty"`
}
type DCNode struct {
	IdentifiedObject
	DCEquipmentContainerMrid uuid.UUID `bun:"dcequipment_container_mrid,type:uuid" json:"dcequipment_container_mrid"`
	DCEquipmentContainer     *Entity   `bun:"rel:belongs-to,join:dcequipment_container_mrid=mrid" json:"dcequipment_container,omitempty"`
}
type HydroGeneratingUnit struct {
	GeneratingUnit
	HydroPowerPlantMrid          uuid.UUID                  `bun:"hydro_power_plant_mrid,type:uuid" json:"hydro_power_plant_mrid"`
	HydroPowerPlant              *Entity                    `bun:"rel:belongs-to,join:hydro_power_plant_mrid=mrid" json:"hydro_power_plant,omitempty"`
	EnergyConversionCapabilityId int                        `bun:"energy_conversion_capability_id" json:"energy_conversion_capability_id"`
	EnergyConversionCapability   *HydroEnergyConversionKind `bun:"rel:belongs-to,join:energy_conversion_capability_id=id" json:"energy_conversion_capability,omitempty"`
}
type ConformLoadSchedule struct {
	SeasonDayTypeSchedule
	ConformLoadGroupMrid uuid.UUID `bun:"conform_load_group_mrid,type:uuid" json:"conform_load_group_mrid"`
	ConformLoadGroup     *Entity   `bun:"rel:belongs-to,join:conform_load_group_mrid=mrid" json:"conform_load_group,omitempty"`
}
type DCLineSegment struct {
	DCConductingEquipment
	PerLengthParameterMrid uuid.UUID `bun:"per_length_parameter_mrid,type:uuid" json:"per_length_parameter_mrid"`
	PerLengthParameter     *Entity   `bun:"rel:belongs-to,join:per_length_parameter_mrid=mrid" json:"per_length_parameter,omitempty"`
	Capacitance            float64   `bun:"capacitance" json:"capacitance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCLineSegment.capacitance"`
	Resistance             float64   `bun:"resistance" json:"resistance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCLineSegment.resistance"`
	Inductance             float64   `bun:"inductance" json:"inductance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCLineSegment.inductance"`
	Length                 float64   `bun:"length" json:"length" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCLineSegment.length"`
}
type EquivalentShunt struct {
	EquivalentEquipment
	B float64 `bun:"b" json:"b" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentShunt.b"`
	G float64 `bun:"g" json:"g" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentShunt.g"`
}
type Conductance struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Conductance.value"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
}
type RegulatingControl struct {
	PowerSystemResource
	ModeId       int                        `bun:"mode_id" json:"mode_id"`
	Mode         *RegulatingControlModeKind `bun:"rel:belongs-to,join:mode_id=id" json:"mode,omitempty"`
	TerminalMrid uuid.UUID                  `bun:"terminal_mrid,type:uuid" json:"terminal_mrid"`
	Terminal     *Entity                    `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
}
type ActivePowerPerCurrentFlow struct {
	MultiplierId            int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier              *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId                  int             `bun:"unit_id" json:"unit_id"`
	Unit                    *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	DenominatorUnitId       int             `bun:"denominator_unit_id" json:"denominator_unit_id"`
	DenominatorUnit         *UnitSymbol     `bun:"rel:belongs-to,join:denominator_unit_id=id" json:"denominator_unit,omitempty"`
	Value                   float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerCurrentFlow.value"`
	DenominatorMultiplierId int             `bun:"denominator_multiplier_id" json:"denominator_multiplier_id"`
	DenominatorMultiplier   *UnitMultiplier `bun:"rel:belongs-to,join:denominator_multiplier_id=id" json:"denominator_multiplier,omitempty"`
}
type EquivalentEquipment struct {
	ConductingEquipment
	EquivalentNetworkMrid uuid.UUID `bun:"equivalent_network_mrid,type:uuid" json:"equivalent_network_mrid"`
	EquivalentNetwork     *Entity   `bun:"rel:belongs-to,join:equivalent_network_mrid=mrid" json:"equivalent_network,omitempty"`
}
type ActivePower struct {
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePower.value"`
}
type ExternalNetworkInjection struct {
	RegulatingCondEq
	MaxP        float64 `bun:"max_p" json:"max_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.maxP"`
	MaxQ        float64 `bun:"max_q" json:"max_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.maxQ"`
	GovernorSCD float64 `bun:"governor_scd" json:"governor_scd" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.governorSCD"`
	MinQ        float64 `bun:"min_q" json:"min_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.minQ"`
	MinP        float64 `bun:"min_p" json:"min_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.minP"`
}
type SynchronousMachine struct {
	RotatingMachine
	MaxQ                               float64                 `bun:"max_q" json:"max_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachine.maxQ"`
	TypeId                             int                     `bun:"type_id" json:"type_id"`
	Type                               *SynchronousMachineKind `bun:"rel:belongs-to,join:type_id=id" json:"type,omitempty"`
	QPercent                           float64                 `bun:"qpercent" json:"qpercent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachine.qPercent"`
	MinQ                               float64                 `bun:"min_q" json:"min_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachine.minQ"`
	InitialReactiveCapabilityCurveMrid uuid.UUID               `bun:"initial_reactive_capability_curve_mrid,type:uuid" json:"initial_reactive_capability_curve_mrid"`
	InitialReactiveCapabilityCurve     *Entity                 `bun:"rel:belongs-to,join:initial_reactive_capability_curve_mrid=mrid" json:"initial_reactive_capability_curve,omitempty"`
}
type ControlAreaGeneratingUnit struct {
	IdentifiedObject
	ControlAreaMrid    uuid.UUID `bun:"control_area_mrid,type:uuid" json:"control_area_mrid"`
	ControlArea        *Entity   `bun:"rel:belongs-to,join:control_area_mrid=mrid" json:"control_area,omitempty"`
	GeneratingUnitMrid uuid.UUID `bun:"generating_unit_mrid,type:uuid" json:"generating_unit_mrid"`
	GeneratingUnit     *Entity   `bun:"rel:belongs-to,join:generating_unit_mrid=mrid" json:"generating_unit,omitempty"`
}
type PhaseTapChangerTablePoint struct {
	TapChangerTablePoint
	PhaseTapChangerTableMrid uuid.UUID `bun:"phase_tap_changer_table_mrid,type:uuid" json:"phase_tap_changer_table_mrid"`
	PhaseTapChangerTable     *Entity   `bun:"rel:belongs-to,join:phase_tap_changer_table_mrid=mrid" json:"phase_tap_changer_table,omitempty"`
	Angle                    float64   `bun:"angle" json:"angle" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerTablePoint.angle"`
}
type EquivalentInjection struct {
	EquivalentEquipment
	MaxQ                        float64   `bun:"max_q" json:"max_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.maxQ"`
	RegulationCapability        bool      `bun:"regulation_capability" json:"regulation_capability" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.regulationCapability"`
	MaxP                        float64   `bun:"max_p" json:"max_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.maxP"`
	ReactiveCapabilityCurveMrid uuid.UUID `bun:"reactive_capability_curve_mrid,type:uuid" json:"reactive_capability_curve_mrid"`
	ReactiveCapabilityCurve     *Entity   `bun:"rel:belongs-to,join:reactive_capability_curve_mrid=mrid" json:"reactive_capability_curve,omitempty"`
	MinP                        float64   `bun:"min_p" json:"min_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.minP"`
	MinQ                        float64   `bun:"min_q" json:"min_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.minQ"`
}
type OperationalLimitSet struct {
	IdentifiedObject
	TerminalMrid  uuid.UUID `bun:"terminal_mrid,type:uuid" json:"terminal_mrid"`
	Terminal      *Entity   `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
	EquipmentMrid uuid.UUID `bun:"equipment_mrid,type:uuid" json:"equipment_mrid"`
	Equipment     *Entity   `bun:"rel:belongs-to,join:equipment_mrid=mrid" json:"equipment,omitempty"`
}
type Curve struct {
	IdentifiedObject
	XUnitId      int         `bun:"xunit_id" json:"xunit_id"`
	XUnit        *UnitSymbol `bun:"rel:belongs-to,join:xunit_id=id" json:"xunit,omitempty"`
	CurveStyleId int         `bun:"curve_style_id" json:"curve_style_id"`
	CurveStyle   *CurveStyle `bun:"rel:belongs-to,join:curve_style_id=id" json:"curve_style,omitempty"`
	Y1UnitId     int         `bun:"y1_unit_id" json:"y1_unit_id"`
	Y1Unit       *UnitSymbol `bun:"rel:belongs-to,join:y1_unit_id=id" json:"y1_unit,omitempty"`
	Y2UnitId     int         `bun:"y2_unit_id" json:"y2_unit_id"`
	Y2Unit       *UnitSymbol `bun:"rel:belongs-to,join:y2_unit_id=id" json:"y2_unit,omitempty"`
}
type DCSeriesDevice struct {
	DCConductingEquipment
	RatedUdc   float64 `bun:"rated_udc" json:"rated_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCSeriesDevice.ratedUdc"`
	Inductance float64 `bun:"inductance" json:"inductance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCSeriesDevice.inductance"`
	Resistance float64 `bun:"resistance" json:"resistance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCSeriesDevice.resistance"`
}
type NonConformLoadSchedule struct {
	SeasonDayTypeSchedule
	NonConformLoadGroupMrid uuid.UUID `bun:"non_conform_load_group_mrid,type:uuid" json:"non_conform_load_group_mrid"`
	NonConformLoadGroup     *Entity   `bun:"rel:belongs-to,join:non_conform_load_group_mrid=mrid" json:"non_conform_load_group,omitempty"`
}
type ReactivePower struct {
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ReactivePower.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
}
type HydroPowerPlant struct {
	PowerSystemResource
	HydroPlantStorageTypeId int                    `bun:"hydro_plant_storage_type_id" json:"hydro_plant_storage_type_id"`
	HydroPlantStorageType   *HydroPlantStorageKind `bun:"rel:belongs-to,join:hydro_plant_storage_type_id=id" json:"hydro_plant_storage_type,omitempty"`
}
type Susceptance struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Susceptance.value"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
}
type StaticVarCompensator struct {
	RegulatingCondEq
	SVCControlModeId int             `bun:"svccontrol_mode_id" json:"svccontrol_mode_id"`
	SVCControlMode   *SVCControlMode `bun:"rel:belongs-to,join:svccontrol_mode_id=id" json:"svccontrol_mode,omitempty"`
	Slope            float64         `bun:"slope" json:"slope" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.slope"`
	InductiveRating  float64         `bun:"inductive_rating" json:"inductive_rating" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.inductiveRating"`
	VoltageSetPoint  float64         `bun:"voltage_set_point" json:"voltage_set_point" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.voltageSetPoint"`
	CapacitiveRating float64         `bun:"capacitive_rating" json:"capacitive_rating" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.capacitiveRating"`
}
type EquipmentVersion struct {
	EntsoeURIcore         string    `bun:"entsoe_uricore" json:"entsoe_uricore" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.entsoeURIcore"`
	EntsoeURIshortCircuit string    `bun:"entsoe_urishort_circuit" json:"entsoe_urishort_circuit" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.entsoeURIshortCircuit"`
	EntsoeUML             string    `bun:"entsoe_uml" json:"entsoe_uml" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.entsoeUML"`
	ShortName             string    `bun:"short_name" json:"short_name" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.shortName"`
	BaseURIshortCircuit   string    `bun:"base_urishort_circuit" json:"base_urishort_circuit" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.baseURIshortCircuit"`
	BaseURIcore           string    `bun:"base_uricore" json:"base_uricore" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.baseURIcore"`
	ModelDescriptionURI   string    `bun:"model_description_uri" json:"model_description_uri" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.modelDescriptionURI"`
	NamespaceRDF          string    `bun:"namespace_rdf" json:"namespace_rdf" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.namespaceRDF"`
	Date                  time.Time `bun:"date" json:"date" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.date"`
	NamespaceUML          string    `bun:"namespace_uml" json:"namespace_uml" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.namespaceUML"`
	BaseUML               string    `bun:"base_uml" json:"base_uml" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.baseUML"`
	DifferenceModelURI    string    `bun:"difference_model_uri" json:"difference_model_uri" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.differenceModelURI"`
	BaseURIoperation      string    `bun:"base_urioperation" json:"base_urioperation" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.baseURIoperation"`
	EntsoeURIoperation    string    `bun:"entsoe_urioperation" json:"entsoe_urioperation" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.entsoeURIoperation"`
}
type GeneratingUnit struct {
	Equipment
	VariableCost                    float64                 `bun:"variable_cost" json:"variable_cost" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.variableCost"`
	StartupCost                     float64                 `bun:"startup_cost" json:"startup_cost" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.startupCost"`
	RatedNetMaxP                    float64                 `bun:"rated_net_max_p" json:"rated_net_max_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.ratedNetMaxP"`
	MaxOperatingP                   float64                 `bun:"max_operating_p" json:"max_operating_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.maxOperatingP"`
	LongPF                          float64                 `bun:"long_pf" json:"long_pf" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.longPF"`
	InitialP                        float64                 `bun:"initial_p" json:"initial_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.initialP"`
	RatedGrossMaxP                  float64                 `bun:"rated_gross_max_p" json:"rated_gross_max_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.ratedGrossMaxP"`
	TotalEfficiency                 float64                 `bun:"total_efficiency" json:"total_efficiency" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.totalEfficiency"`
	RatedGrossMinP                  float64                 `bun:"rated_gross_min_p" json:"rated_gross_min_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.ratedGrossMinP"`
	MaximumAllowableSpinningReserve float64                 `bun:"maximum_allowable_spinning_reserve" json:"maximum_allowable_spinning_reserve" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.maximumAllowableSpinningReserve"`
	MinOperatingP                   float64                 `bun:"min_operating_p" json:"min_operating_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.minOperatingP"`
	GenControlSourceId              int                     `bun:"gen_control_source_id" json:"gen_control_source_id"`
	GenControlSource                *GeneratorControlSource `bun:"rel:belongs-to,join:gen_control_source_id=id" json:"gen_control_source,omitempty"`
	GovernorSCD                     float64                 `bun:"governor_scd" json:"governor_scd" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.governorSCD"`
	NominalP                        float64                 `bun:"nominal_p" json:"nominal_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.nominalP"`
	ShortPF                         float64                 `bun:"short_pf" json:"short_pf" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.shortPF"`
}
type Simple_Float struct {
	Value float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Simple_Float.value"`
}
type ConformLoad struct {
	EnergyConsumer
	LoadGroupMrid uuid.UUID `bun:"load_group_mrid,type:uuid" json:"load_group_mrid"`
	LoadGroup     *Entity   `bun:"rel:belongs-to,join:load_group_mrid=mrid" json:"load_group,omitempty"`
}
type PhaseTapChangerAsymmetrical struct {
	PhaseTapChangerNonLinear
	WindingConnectionAngle float64 `bun:"winding_connection_angle" json:"winding_connection_angle" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerAsymmetrical.windingConnectionAngle"`
}
type AngleDegrees struct {
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AngleDegrees.value"`
}
type CapacitancePerLength struct {
	MultiplierId            int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier              *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId                  int             `bun:"unit_id" json:"unit_id"`
	Unit                    *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	DenominatorUnitId       int             `bun:"denominator_unit_id" json:"denominator_unit_id"`
	DenominatorUnit         *UnitSymbol     `bun:"rel:belongs-to,join:denominator_unit_id=id" json:"denominator_unit,omitempty"`
	DenominatorMultiplierId int             `bun:"denominator_multiplier_id" json:"denominator_multiplier_id"`
	DenominatorMultiplier   *UnitMultiplier `bun:"rel:belongs-to,join:denominator_multiplier_id=id" json:"denominator_multiplier,omitempty"`
	Value                   float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CapacitancePerLength.value"`
}
type ACDCConverterDCTerminal struct {
	DCBaseTerminal
	PolarityId                int             `bun:"polarity_id" json:"polarity_id"`
	Polarity                  *DCPolarityKind `bun:"rel:belongs-to,join:polarity_id=id" json:"polarity,omitempty"`
	DCConductingEquipmentMrid uuid.UUID       `bun:"dcconducting_equipment_mrid,type:uuid" json:"dcconducting_equipment_mrid"`
	DCConductingEquipment     *Entity         `bun:"rel:belongs-to,join:dcconducting_equipment_mrid=mrid" json:"dcconducting_equipment,omitempty"`
}
type PerLengthDCLineParameter struct {
	Resistance  float64 `bun:"resistance" json:"resistance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerLengthDCLineParameter.resistance"`
	Capacitance float64 `bun:"capacitance" json:"capacitance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerLengthDCLineParameter.capacitance"`
	Inductance  float64 `bun:"inductance" json:"inductance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerLengthDCLineParameter.inductance"`
}
type VsConverter struct {
	ACDCConverter
	CapabilityCurveMrid uuid.UUID `bun:"capability_curve_mrid,type:uuid" json:"capability_curve_mrid"`
	CapabilityCurve     *Entity   `bun:"rel:belongs-to,join:capability_curve_mrid=mrid" json:"capability_curve,omitempty"`
	MaxValveCurrent     float64   `bun:"max_valve_current" json:"max_valve_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VsConverter.maxValveCurrent"`
	MaxModulationIndex  float64   `bun:"max_modulation_index" json:"max_modulation_index" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VsConverter.maxModulationIndex"`
}
type DCGround struct {
	DCConductingEquipment
	Inductance float64 `bun:"inductance" json:"inductance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCGround.inductance"`
	R          float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCGround.r"`
}
type Equipment struct {
	PowerSystemResource
	Aggregate              bool      `bun:"aggregate" json:"aggregate" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Equipment.aggregate"`
	EquipmentContainerMrid uuid.UUID `bun:"equipment_container_mrid,type:uuid" json:"equipment_container_mrid"`
	EquipmentContainer     *Entity   `bun:"rel:belongs-to,join:equipment_container_mrid=mrid" json:"equipment_container,omitempty"`
}
type EnergyConsumer struct {
	ConductingEquipment
	LoadResponseMrid uuid.UUID `bun:"load_response_mrid,type:uuid" json:"load_response_mrid"`
	LoadResponse     *Entity   `bun:"rel:belongs-to,join:load_response_mrid=mrid" json:"load_response,omitempty"`
}
type Money struct {
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *Currency       `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Money.value"`
}
type RotationSpeed struct {
	MultiplierId            int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier              *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	DenominatorUnitId       int             `bun:"denominator_unit_id" json:"denominator_unit_id"`
	DenominatorUnit         *UnitSymbol     `bun:"rel:belongs-to,join:denominator_unit_id=id" json:"denominator_unit,omitempty"`
	Value                   float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotationSpeed.value"`
	DenominatorMultiplierId int             `bun:"denominator_multiplier_id" json:"denominator_multiplier_id"`
	DenominatorMultiplier   *UnitMultiplier `bun:"rel:belongs-to,join:denominator_multiplier_id=id" json:"denominator_multiplier,omitempty"`
	UnitId                  int             `bun:"unit_id" json:"unit_id"`
	Unit                    *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type WindGeneratingUnit struct {
	GeneratingUnit
	WindGenUnitTypeId int              `bun:"wind_gen_unit_type_id" json:"wind_gen_unit_type_id"`
	WindGenUnitType   *WindGenUnitKind `bun:"rel:belongs-to,join:wind_gen_unit_type_id=id" json:"wind_gen_unit_type,omitempty"`
}
type CurrentLimit struct {
	OperationalLimit
	Value float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurrentLimit.value"`
}
type RotatingMachine struct {
	RegulatingCondEq
	RatedU             float64   `bun:"rated_u" json:"rated_u" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotatingMachine.ratedU"`
	GeneratingUnitMrid uuid.UUID `bun:"generating_unit_mrid,type:uuid" json:"generating_unit_mrid"`
	GeneratingUnit     *Entity   `bun:"rel:belongs-to,join:generating_unit_mrid=mrid" json:"generating_unit,omitempty"`
	RatedS             float64   `bun:"rated_s" json:"rated_s" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotatingMachine.ratedS"`
	RatedPowerFactor   float64   `bun:"rated_power_factor" json:"rated_power_factor" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotatingMachine.ratedPowerFactor"`
}
type NonConformLoad struct {
	EnergyConsumer
	LoadGroupMrid uuid.UUID `bun:"load_group_mrid,type:uuid" json:"load_group_mrid"`
	LoadGroup     *Entity   `bun:"rel:belongs-to,join:load_group_mrid=mrid" json:"load_group,omitempty"`
}
type ACLineSegment struct {
	Conductor
	X   float64 `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACLineSegment.x"`
	Bch float64 `bun:"bch" json:"bch" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACLineSegment.bch"`
	Gch float64 `bun:"gch" json:"gch" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACLineSegment.gch"`
	R   float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACLineSegment.r"`
}
type Switch struct {
	ConductingEquipment
	NormalOpen   bool    `bun:"normal_open" json:"normal_open" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Switch.normalOpen"`
	RatedCurrent float64 `bun:"rated_current" json:"rated_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Switch.ratedCurrent"`
	Retained     bool    `bun:"retained" json:"retained" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Switch.retained"`
}
type EquivalentBranch struct {
	EquivalentEquipment
	R21 float64 `bun:"r21" json:"r21" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentBranch.r21"`
	X   float64 `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentBranch.x"`
	R   float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentBranch.r"`
	X21 float64 `bun:"x21" json:"x21" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentBranch.x21"`
}
type PerCent struct {
	Value        float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerCent.value"`
	MultiplierId int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier   *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
	UnitId       int             `bun:"unit_id" json:"unit_id"`
	Unit         *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
}
type ActivePowerPerFrequency struct {
	DenominatorUnitId       int             `bun:"denominator_unit_id" json:"denominator_unit_id"`
	DenominatorUnit         *UnitSymbol     `bun:"rel:belongs-to,join:denominator_unit_id=id" json:"denominator_unit,omitempty"`
	UnitId                  int             `bun:"unit_id" json:"unit_id"`
	Unit                    *UnitSymbol     `bun:"rel:belongs-to,join:unit_id=id" json:"unit,omitempty"`
	DenominatorMultiplierId int             `bun:"denominator_multiplier_id" json:"denominator_multiplier_id"`
	DenominatorMultiplier   *UnitMultiplier `bun:"rel:belongs-to,join:denominator_multiplier_id=id" json:"denominator_multiplier,omitempty"`
	Value                   float64         `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerFrequency.value"`
	MultiplierId            int             `bun:"multiplier_id" json:"multiplier_id"`
	Multiplier              *UnitMultiplier `bun:"rel:belongs-to,join:multiplier_id=id" json:"multiplier,omitempty"`
}
type ConductingEquipment struct {
	Equipment
	BaseVoltageMrid uuid.UUID `bun:"base_voltage_mrid,type:uuid" json:"base_voltage_mrid"`
	BaseVoltage     *Entity   `bun:"rel:belongs-to,join:base_voltage_mrid=mrid" json:"base_voltage,omitempty"`
}
type ACDCTerminal struct {
	IdentifiedObject
	SequenceNumber    int       `bun:"sequence_number" json:"sequence_number" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCTerminal.sequenceNumber"`
	BusNameMarkerMrid uuid.UUID `bun:"bus_name_marker_mrid,type:uuid" json:"bus_name_marker_mrid"`
	BusNameMarker     *Entity   `bun:"rel:belongs-to,join:bus_name_marker_mrid=mrid" json:"bus_name_marker,omitempty"`
}
type RegulatingCondEq struct {
	ConductingEquipment
	RegulatingControlMrid uuid.UUID `bun:"regulating_control_mrid,type:uuid" json:"regulating_control_mrid"`
	RegulatingControl     *Entity   `bun:"rel:belongs-to,join:regulating_control_mrid=mrid" json:"regulating_control,omitempty"`
}
type DCBreaker struct {
	DCSwitch
}
type DCSwitch struct {
	DCConductingEquipment
}
type LoadBreakSwitch struct {
	ProtectedSwitch
}
type ProtectedSwitch struct {
	Switch
}
type NuclearGeneratingUnit struct {
	GeneratingUnit
}
type DCConductingEquipment struct {
	Equipment
}
type PhaseTapChangerTable struct {
	IdentifiedObject
}
type PowerTransformer struct {
	ConductingEquipment
}
type PowerSystemResource struct {
	IdentifiedObject
}
type RatioTapChangerTable struct {
	IdentifiedObject
}
type SeasonDayTypeSchedule struct {
}
type EquivalentNetwork struct {
	ConnectivityNodeContainer
}
type ConnectivityNodeContainer struct {
	PowerSystemResource
}
type ThermalGeneratingUnit struct {
	GeneratingUnit
}
type TapChangerControl struct {
	RegulatingControl
}
type Junction struct {
	Connector
}
type Connector struct {
	ConductingEquipment
}
type EquipmentContainer struct {
	ConnectivityNodeContainer
}
type NonlinearShuntCompensator struct {
	ShuntCompensator
}
type ReactiveCapabilityCurve struct {
	Curve
}
type DCBusbar struct {
	DCConductingEquipment
}
type NonConformLoadGroup struct {
	LoadGroup
}
type LoadGroup struct {
	IdentifiedObject
}
type DCEquipmentContainer struct {
	EquipmentContainer
}
type DCChopper struct {
	DCConductingEquipment
}
type SolarGeneratingUnit struct {
	GeneratingUnit
}
type GeographicalRegion struct {
	IdentifiedObject
}
type Breaker struct {
	ProtectedSwitch
}
type Disconnector struct {
	Switch
}
type EnergySchedulingType struct {
	IdentifiedObject
}
type VsCapabilityCurve struct {
	Curve
}
type ReportingGroup struct {
	IdentifiedObject
}
type BusbarSection struct {
	Connector
}
type DCDisconnector struct {
	DCSwitch
}
type PhaseTapChangerSymmetrical struct {
	PhaseTapChangerNonLinear
}
type ConformLoadGroup struct {
	LoadGroup
}
