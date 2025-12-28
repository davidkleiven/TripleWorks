package models

import (
	"github.com/google/uuid"
	"time"
)

type Entity struct {
	ModelEntity
	Mrid uuid.UUID `bun:"mrid,type:uuid,pk"`
}
type AngleDegrees struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AngleDegrees.unit"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AngleDegrees.multiplier"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AngleDegrees.value"`
}
type Resistance struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Resistance.unit"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Resistance.value"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Resistance.multiplier"`
}
type ConformLoad struct {
	EnergyConsumer
	LoadGroupMrid uuid.UUID `bun:"load_group_mrid,type:uuid" json:"load_group_mrid"`
	LoadGroup     *Entity   `bun:"rel:belongs-to,join:load_group_mrid=mrid" json:"load_group,omitempty"`
}
type ACDCTerminal struct {
	IdentifiedObject
	BusNameMarkerMrid uuid.UUID `bun:"bus_name_marker_mrid,type:uuid" json:"bus_name_marker_mrid"`
	BusNameMarker     *Entity   `bun:"rel:belongs-to,join:bus_name_marker_mrid=mrid" json:"bus_name_marker,omitempty"`
	SequenceNumber    int       `bun:"sequence_number" json:"sequence_number" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCTerminal.sequenceNumber"`
}
type ControlAreaGeneratingUnit struct {
	IdentifiedObject
	ControlAreaMrid    uuid.UUID `bun:"control_area_mrid,type:uuid" json:"control_area_mrid"`
	ControlArea        *Entity   `bun:"rel:belongs-to,join:control_area_mrid=mrid" json:"control_area,omitempty"`
	GeneratingUnitMrid uuid.UUID `bun:"generating_unit_mrid,type:uuid" json:"generating_unit_mrid"`
	GeneratingUnit     *Entity   `bun:"rel:belongs-to,join:generating_unit_mrid=mrid" json:"generating_unit,omitempty"`
}
type DCNode struct {
	IdentifiedObject
	DCEquipmentContainerMrid uuid.UUID `bun:"dcequipment_container_mrid,type:uuid" json:"dcequipment_container_mrid"`
	DCEquipmentContainer     *Entity   `bun:"rel:belongs-to,join:dcequipment_container_mrid=mrid" json:"dcequipment_container,omitempty"`
}
type InductancePerLength struct {
	Value                 float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#InductancePerLength.value"`
	DenominatorUnit       string  `bun:"denominator_unit" json:"denominator_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#InductancePerLength.denominatorUnit"`
	Unit                  string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#InductancePerLength.unit"`
	DenominatorMultiplier string  `bun:"denominator_multiplier" json:"denominator_multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#InductancePerLength.denominatorMultiplier"`
	Multiplier            string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#InductancePerLength.multiplier"`
}
type CurrentFlow struct {
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurrentFlow.value"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurrentFlow.multiplier"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurrentFlow.unit"`
}
type Susceptance struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Susceptance.unit"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Susceptance.value"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Susceptance.multiplier"`
}
type SynchronousMachine struct {
	RotatingMachine
	QPercent                           float64   `bun:"qpercent" json:"qpercent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachine.qPercent"`
	InitialReactiveCapabilityCurveMrid uuid.UUID `bun:"initial_reactive_capability_curve_mrid,type:uuid" json:"initial_reactive_capability_curve_mrid"`
	InitialReactiveCapabilityCurve     *Entity   `bun:"rel:belongs-to,join:initial_reactive_capability_curve_mrid=mrid" json:"initial_reactive_capability_curve,omitempty"`
	MinQ                               float64   `bun:"min_q" json:"min_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachine.minQ"`
	MaxQ                               float64   `bun:"max_q" json:"max_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachine.maxQ"`
	Type                               string    `bun:"type" json:"type" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachine.type"`
}
type Frequency struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Frequency.unit"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Frequency.value"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Frequency.multiplier"`
}
type ReactivePower struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ReactivePower.unit"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ReactivePower.value"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ReactivePower.multiplier"`
}
type EquivalentShunt struct {
	EquivalentEquipment
	B float64 `bun:"b" json:"b" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentShunt.b"`
	G float64 `bun:"g" json:"g" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentShunt.g"`
}
type EquipmentVersion struct {
	Date                  time.Time `bun:"date" json:"date" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.date"`
	ShortName             string    `bun:"short_name" json:"short_name" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.shortName"`
	BaseUML               string    `bun:"base_uml" json:"base_uml" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.baseUML"`
	NamespaceRDF          string    `bun:"namespace_rdf" json:"namespace_rdf" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.namespaceRDF"`
	EntsoeURIoperation    string    `bun:"entsoe_urioperation" json:"entsoe_urioperation" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.entsoeURIoperation"`
	BaseURIcore           string    `bun:"base_uricore" json:"base_uricore" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.baseURIcore"`
	NamespaceUML          string    `bun:"namespace_uml" json:"namespace_uml" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.namespaceUML"`
	EntsoeURIshortCircuit string    `bun:"entsoe_urishort_circuit" json:"entsoe_urishort_circuit" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.entsoeURIshortCircuit"`
	BaseURIoperation      string    `bun:"base_urioperation" json:"base_urioperation" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.baseURIoperation"`
	EntsoeURIcore         string    `bun:"entsoe_uricore" json:"entsoe_uricore" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.entsoeURIcore"`
	ModelDescriptionURI   string    `bun:"model_description_uri" json:"model_description_uri" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.modelDescriptionURI"`
	BaseURIshortCircuit   string    `bun:"base_urishort_circuit" json:"base_urishort_circuit" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.baseURIshortCircuit"`
	EntsoeUML             string    `bun:"entsoe_uml" json:"entsoe_uml" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.entsoeUML"`
	DifferenceModelURI    string    `bun:"difference_model_uri" json:"difference_model_uri" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#EquipmentVersion.differenceModelURI"`
}
type NonConformLoadSchedule struct {
	SeasonDayTypeSchedule
	NonConformLoadGroupMrid uuid.UUID `bun:"non_conform_load_group_mrid,type:uuid" json:"non_conform_load_group_mrid"`
	NonConformLoadGroup     *Entity   `bun:"rel:belongs-to,join:non_conform_load_group_mrid=mrid" json:"non_conform_load_group,omitempty"`
}
type Capacitance struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Capacitance.unit"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Capacitance.value"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Capacitance.multiplier"`
}
type PhaseTapChanger struct {
	TapChanger
	TransformerEndMrid uuid.UUID `bun:"transformer_end_mrid,type:uuid" json:"transformer_end_mrid"`
	TransformerEnd     *Entity   `bun:"rel:belongs-to,join:transformer_end_mrid=mrid" json:"transformer_end,omitempty"`
}
type RegulatingControl struct {
	PowerSystemResource
	TerminalMrid uuid.UUID `bun:"terminal_mrid,type:uuid" json:"terminal_mrid"`
	Terminal     *Entity   `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
	Mode         string    `bun:"mode" json:"mode" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControl.mode"`
}
type WindGeneratingUnit struct {
	GeneratingUnit
	WindGenUnitType string `bun:"wind_gen_unit_type" json:"wind_gen_unit_type" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#WindGeneratingUnit.windGenUnitType"`
}
type PhaseTapChangerLinear struct {
	PhaseTapChanger
	StepPhaseShiftIncrement float64 `bun:"step_phase_shift_increment" json:"step_phase_shift_increment" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerLinear.stepPhaseShiftIncrement"`
	XMax                    float64 `bun:"x_max" json:"x_max" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerLinear.xMax"`
	XMin                    float64 `bun:"x_min" json:"x_min" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerLinear.xMin"`
}
type VoltageLimit struct {
	OperationalLimit
	Value float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltageLimit.value"`
}
type IdentifiedObject struct {
	BaseEntity
	ShortName          string `bun:"short_name" json:"short_name" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#IdentifiedObject.shortName"`
	Mrid               string `bun:"mrid" json:"mrid" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.mRID"`
	Description        string `bun:"description" json:"description" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.description"`
	Name               string `bun:"name" json:"name" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.name"`
	EnergyIdentCodeEic string `bun:"energy_ident_code_eic" json:"energy_ident_code_eic" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#IdentifiedObject.energyIdentCodeEic"`
}
type FossilFuel struct {
	IdentifiedObject
	ThermalGeneratingUnitMrid uuid.UUID `bun:"thermal_generating_unit_mrid,type:uuid" json:"thermal_generating_unit_mrid"`
	ThermalGeneratingUnit     *Entity   `bun:"rel:belongs-to,join:thermal_generating_unit_mrid=mrid" json:"thermal_generating_unit,omitempty"`
	FossilFuelType            string    `bun:"fossil_fuel_type" json:"fossil_fuel_type" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#FossilFuel.fossilFuelType"`
}
type BaseVoltage struct {
	IdentifiedObject
	NominalVoltage float64 `bun:"nominal_voltage" json:"nominal_voltage" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#BaseVoltage.nominalVoltage"`
}
type Voltage struct {
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Voltage.value"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Voltage.unit"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Voltage.multiplier"`
}
type NonConformLoad struct {
	EnergyConsumer
	LoadGroupMrid uuid.UUID `bun:"load_group_mrid,type:uuid" json:"load_group_mrid"`
	LoadGroup     *Entity   `bun:"rel:belongs-to,join:load_group_mrid=mrid" json:"load_group,omitempty"`
}
type AsynchronousMachine struct {
	RotatingMachine
	NominalSpeed     float64 `bun:"nominal_speed" json:"nominal_speed" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AsynchronousMachine.nominalSpeed"`
	NominalFrequency float64 `bun:"nominal_frequency" json:"nominal_frequency" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AsynchronousMachine.nominalFrequency"`
}
type DCShunt struct {
	DCConductingEquipment
	Resistance  float64 `bun:"resistance" json:"resistance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCShunt.resistance"`
	Capacitance float64 `bun:"capacitance" json:"capacitance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCShunt.capacitance"`
	RatedUdc    float64 `bun:"rated_udc" json:"rated_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCShunt.ratedUdc"`
}
type EquivalentBranch struct {
	EquivalentEquipment
	R   float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentBranch.r"`
	X   float64 `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentBranch.x"`
	X21 float64 `bun:"x21" json:"x21" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentBranch.x21"`
	R21 float64 `bun:"r21" json:"r21" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentBranch.r21"`
}
type ACDCConverter struct {
	ConductingEquipment
	NumberOfValves  int       `bun:"number_of_valves" json:"number_of_valves" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.numberOfValves"`
	PccTerminalMrid uuid.UUID `bun:"pcc_terminal_mrid,type:uuid" json:"pcc_terminal_mrid"`
	PccTerminal     *Entity   `bun:"rel:belongs-to,join:pcc_terminal_mrid=mrid" json:"pcc_terminal,omitempty"`
	SwitchingLoss   float64   `bun:"switching_loss" json:"switching_loss" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.switchingLoss"`
	BaseS           float64   `bun:"base_s" json:"base_s" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.baseS"`
	MinUdc          float64   `bun:"min_udc" json:"min_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.minUdc"`
	MaxUdc          float64   `bun:"max_udc" json:"max_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.maxUdc"`
	ResistiveLoss   float64   `bun:"resistive_loss" json:"resistive_loss" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.resistiveLoss"`
	IdleLoss        float64   `bun:"idle_loss" json:"idle_loss" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.idleLoss"`
	ValveU0         float64   `bun:"valve_u0" json:"valve_u0" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.valveU0"`
	RatedUdc        float64   `bun:"rated_udc" json:"rated_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverter.ratedUdc"`
}
type DCLine struct {
	DCEquipmentContainer
	RegionMrid uuid.UUID `bun:"region_mrid,type:uuid" json:"region_mrid"`
	Region     *Entity   `bun:"rel:belongs-to,join:region_mrid=mrid" json:"region,omitempty"`
}
type CurveData struct {
	Y1value   float64   `bun:"y1value" json:"y1value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurveData.y1value"`
	Xvalue    float64   `bun:"xvalue" json:"xvalue" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurveData.xvalue"`
	Y2value   float64   `bun:"y2value" json:"y2value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurveData.y2value"`
	CurveMrid uuid.UUID `bun:"curve_mrid,type:uuid" json:"curve_mrid"`
	Curve     *Entity   `bun:"rel:belongs-to,join:curve_mrid=mrid" json:"curve,omitempty"`
}
type ActivePower struct {
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePower.multiplier"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePower.unit"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePower.value"`
}
type Simple_Float struct {
	Value float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Simple_Float.value"`
}
type ConductingEquipment struct {
	Equipment
	BaseVoltageMrid uuid.UUID `bun:"base_voltage_mrid,type:uuid" json:"base_voltage_mrid"`
	BaseVoltage     *Entity   `bun:"rel:belongs-to,join:base_voltage_mrid=mrid" json:"base_voltage,omitempty"`
}
type PowerTransformerEnd struct {
	TransformerEnd
	PowerTransformerMrid uuid.UUID `bun:"power_transformer_mrid,type:uuid" json:"power_transformer_mrid"`
	PowerTransformer     *Entity   `bun:"rel:belongs-to,join:power_transformer_mrid=mrid" json:"power_transformer,omitempty"`
	RatedS               float64   `bun:"rated_s" json:"rated_s" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.ratedS"`
	G                    float64   `bun:"g" json:"g" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.g"`
	ConnectionKind       string    `bun:"connection_kind" json:"connection_kind" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.connectionKind"`
	X                    float64   `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.x"`
	RatedU               float64   `bun:"rated_u" json:"rated_u" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.ratedU"`
	B                    float64   `bun:"b" json:"b" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.b"`
	R                    float64   `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PowerTransformerEnd.r"`
}
type ApparentPower struct {
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ApparentPower.value"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ApparentPower.unit"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ApparentPower.multiplier"`
}
type Length struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Length.unit"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Length.value"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Length.multiplier"`
}
type TapChangerTablePoint struct {
	Ratio float64 `bun:"ratio" json:"ratio" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.ratio"`
	R     float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.r"`
	B     float64 `bun:"b" json:"b" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.b"`
	X     float64 `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.x"`
	Step  int     `bun:"step" json:"step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.step"`
	G     float64 `bun:"g" json:"g" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChangerTablePoint.g"`
}
type OperationalLimit struct {
	IdentifiedObject
	OperationalLimitSetMrid  uuid.UUID `bun:"operational_limit_set_mrid,type:uuid" json:"operational_limit_set_mrid"`
	OperationalLimitSet      *Entity   `bun:"rel:belongs-to,join:operational_limit_set_mrid=mrid" json:"operational_limit_set,omitempty"`
	OperationalLimitTypeMrid uuid.UUID `bun:"operational_limit_type_mrid,type:uuid" json:"operational_limit_type_mrid"`
	OperationalLimitType     *Entity   `bun:"rel:belongs-to,join:operational_limit_type_mrid=mrid" json:"operational_limit_type,omitempty"`
}
type EquivalentEquipment struct {
	ConductingEquipment
	EquivalentNetworkMrid uuid.UUID `bun:"equivalent_network_mrid,type:uuid" json:"equivalent_network_mrid"`
	EquivalentNetwork     *Entity   `bun:"rel:belongs-to,join:equivalent_network_mrid=mrid" json:"equivalent_network,omitempty"`
}
type VsConverter struct {
	ACDCConverter
	MaxModulationIndex  float64   `bun:"max_modulation_index" json:"max_modulation_index" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VsConverter.maxModulationIndex"`
	MaxValveCurrent     float64   `bun:"max_valve_current" json:"max_valve_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VsConverter.maxValveCurrent"`
	CapabilityCurveMrid uuid.UUID `bun:"capability_curve_mrid,type:uuid" json:"capability_curve_mrid"`
	CapabilityCurve     *Entity   `bun:"rel:belongs-to,join:capability_curve_mrid=mrid" json:"capability_curve,omitempty"`
}
type HydroPowerPlant struct {
	PowerSystemResource
	HydroPlantStorageType string `bun:"hydro_plant_storage_type" json:"hydro_plant_storage_type" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#HydroPowerPlant.hydroPlantStorageType"`
}
type OperationalLimitSet struct {
	IdentifiedObject
	TerminalMrid  uuid.UUID `bun:"terminal_mrid,type:uuid" json:"terminal_mrid"`
	Terminal      *Entity   `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
	EquipmentMrid uuid.UUID `bun:"equipment_mrid,type:uuid" json:"equipment_mrid"`
	Equipment     *Entity   `bun:"rel:belongs-to,join:equipment_mrid=mrid" json:"equipment,omitempty"`
}
type ExternalNetworkInjection struct {
	RegulatingCondEq
	MaxQ        float64 `bun:"max_q" json:"max_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.maxQ"`
	MaxP        float64 `bun:"max_p" json:"max_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.maxP"`
	MinP        float64 `bun:"min_p" json:"min_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.minP"`
	GovernorSCD float64 `bun:"governor_scd" json:"governor_scd" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.governorSCD"`
	MinQ        float64 `bun:"min_q" json:"min_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ExternalNetworkInjection.minQ"`
}
type Equipment struct {
	PowerSystemResource
	EquipmentContainerMrid uuid.UUID `bun:"equipment_container_mrid,type:uuid" json:"equipment_container_mrid"`
	EquipmentContainer     *Entity   `bun:"rel:belongs-to,join:equipment_container_mrid=mrid" json:"equipment_container,omitempty"`
	Aggregate              bool      `bun:"aggregate" json:"aggregate" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Equipment.aggregate"`
}
type PhaseTapChangerTabular struct {
	PhaseTapChanger
	PhaseTapChangerTableMrid uuid.UUID `bun:"phase_tap_changer_table_mrid,type:uuid" json:"phase_tap_changer_table_mrid"`
	PhaseTapChangerTable     *Entity   `bun:"rel:belongs-to,join:phase_tap_changer_table_mrid=mrid" json:"phase_tap_changer_table,omitempty"`
}
type AngleRadians struct {
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AngleRadians.value"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AngleRadians.unit"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#AngleRadians.multiplier"`
}
type TapChanger struct {
	PowerSystemResource
	TapChangerControlMrid uuid.UUID `bun:"tap_changer_control_mrid,type:uuid" json:"tap_changer_control_mrid"`
	TapChangerControl     *Entity   `bun:"rel:belongs-to,join:tap_changer_control_mrid=mrid" json:"tap_changer_control,omitempty"`
	NormalStep            int       `bun:"normal_step" json:"normal_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.normalStep"`
	LtcFlag               bool      `bun:"ltc_flag" json:"ltc_flag" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.ltcFlag"`
	LowStep               int       `bun:"low_step" json:"low_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.lowStep"`
	NeutralStep           int       `bun:"neutral_step" json:"neutral_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.neutralStep"`
	HighStep              int       `bun:"high_step" json:"high_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.highStep"`
	NeutralU              float64   `bun:"neutral_u" json:"neutral_u" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TapChanger.neutralU"`
}
type DCSeriesDevice struct {
	DCConductingEquipment
	RatedUdc   float64 `bun:"rated_udc" json:"rated_udc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCSeriesDevice.ratedUdc"`
	Inductance float64 `bun:"inductance" json:"inductance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCSeriesDevice.inductance"`
	Resistance float64 `bun:"resistance" json:"resistance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCSeriesDevice.resistance"`
}
type Temperature struct {
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Temperature.multiplier"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Temperature.value"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Temperature.unit"`
}
type Terminal struct {
	ACDCTerminal
	Phases                  string    `bun:"phases" json:"phases" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Terminal.phases"`
	ConductingEquipmentMrid uuid.UUID `bun:"conducting_equipment_mrid,type:uuid" json:"conducting_equipment_mrid"`
	ConductingEquipment     *Entity   `bun:"rel:belongs-to,join:conducting_equipment_mrid=mrid" json:"conducting_equipment,omitempty"`
}
type PhaseTapChangerNonLinear struct {
	PhaseTapChanger
	XMin                 float64 `bun:"x_min" json:"x_min" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerNonLinear.xMin"`
	VoltageStepIncrement float64 `bun:"voltage_step_increment" json:"voltage_step_increment" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerNonLinear.voltageStepIncrement"`
	XMax                 float64 `bun:"x_max" json:"x_max" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerNonLinear.xMax"`
}
type PhaseTapChangerAsymmetrical struct {
	PhaseTapChangerNonLinear
	WindingConnectionAngle float64 `bun:"winding_connection_angle" json:"winding_connection_angle" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerAsymmetrical.windingConnectionAngle"`
}
type Seconds struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Seconds.unit"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Seconds.multiplier"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Seconds.value"`
}
type TransformerEnd struct {
	IdentifiedObject
	BaseVoltageMrid uuid.UUID `bun:"base_voltage_mrid,type:uuid" json:"base_voltage_mrid"`
	BaseVoltage     *Entity   `bun:"rel:belongs-to,join:base_voltage_mrid=mrid" json:"base_voltage,omitempty"`
	TerminalMrid    uuid.UUID `bun:"terminal_mrid,type:uuid" json:"terminal_mrid"`
	Terminal        *Entity   `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
	EndNumber       int       `bun:"end_number" json:"end_number" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TransformerEnd.endNumber"`
}
type LoadResponseCharacteristic struct {
	IdentifiedObject
	QConstantCurrent   float64 `bun:"qconstant_current" json:"qconstant_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qConstantCurrent"`
	PVoltageExponent   float64 `bun:"pvoltage_exponent" json:"pvoltage_exponent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pVoltageExponent"`
	PConstantPower     float64 `bun:"pconstant_power" json:"pconstant_power" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pConstantPower"`
	PConstantImpedance float64 `bun:"pconstant_impedance" json:"pconstant_impedance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pConstantImpedance"`
	QConstantPower     float64 `bun:"qconstant_power" json:"qconstant_power" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qConstantPower"`
	PConstantCurrent   float64 `bun:"pconstant_current" json:"pconstant_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pConstantCurrent"`
	QConstantImpedance float64 `bun:"qconstant_impedance" json:"qconstant_impedance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qConstantImpedance"`
	PFrequencyExponent float64 `bun:"pfrequency_exponent" json:"pfrequency_exponent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.pFrequencyExponent"`
	ExponentModel      bool    `bun:"exponent_model" json:"exponent_model" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.exponentModel"`
	QFrequencyExponent float64 `bun:"qfrequency_exponent" json:"qfrequency_exponent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qFrequencyExponent"`
	QVoltageExponent   float64 `bun:"qvoltage_exponent" json:"qvoltage_exponent" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LoadResponseCharacteristic.qVoltageExponent"`
}
type ShuntCompensator struct {
	RegulatingCondEq
	NomU               float64   `bun:"nom_u" json:"nom_u" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.nomU"`
	SwitchOnDate       time.Time `bun:"switch_on_date" json:"switch_on_date" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.switchOnDate"`
	MaximumSections    int       `bun:"maximum_sections" json:"maximum_sections" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.maximumSections"`
	NormalSections     int       `bun:"normal_sections" json:"normal_sections" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.normalSections"`
	Grounded           bool      `bun:"grounded" json:"grounded" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.grounded"`
	AVRDelay           float64   `bun:"avrdelay" json:"avrdelay" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.aVRDelay"`
	SwitchOnCount      int       `bun:"switch_on_count" json:"switch_on_count" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.switchOnCount"`
	VoltageSensitivity float64   `bun:"voltage_sensitivity" json:"voltage_sensitivity" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ShuntCompensator.voltageSensitivity"`
}
type DCBaseTerminal struct {
	ACDCTerminal
	DCNodeMrid uuid.UUID `bun:"dcnode_mrid,type:uuid" json:"dcnode_mrid"`
	DCNode     *Entity   `bun:"rel:belongs-to,join:dcnode_mrid=mrid" json:"dcnode,omitempty"`
}
type DCTerminal struct {
	DCBaseTerminal
	DCConductingEquipmentMrid uuid.UUID `bun:"dcconducting_equipment_mrid,type:uuid" json:"dcconducting_equipment_mrid"`
	DCConductingEquipment     *Entity   `bun:"rel:belongs-to,join:dcconducting_equipment_mrid=mrid" json:"dcconducting_equipment,omitempty"`
}
type CapacitancePerLength struct {
	Value                 float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CapacitancePerLength.value"`
	Unit                  string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CapacitancePerLength.unit"`
	DenominatorMultiplier string  `bun:"denominator_multiplier" json:"denominator_multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CapacitancePerLength.denominatorMultiplier"`
	Multiplier            string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CapacitancePerLength.multiplier"`
	DenominatorUnit       string  `bun:"denominator_unit" json:"denominator_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CapacitancePerLength.denominatorUnit"`
}
type ACDCConverterDCTerminal struct {
	DCBaseTerminal
	DCConductingEquipmentMrid uuid.UUID `bun:"dcconducting_equipment_mrid,type:uuid" json:"dcconducting_equipment_mrid"`
	DCConductingEquipment     *Entity   `bun:"rel:belongs-to,join:dcconducting_equipment_mrid=mrid" json:"dcconducting_equipment,omitempty"`
	Polarity                  string    `bun:"polarity" json:"polarity" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACDCConverterDCTerminal.polarity"`
}
type PU struct {
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PU.value"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PU.multiplier"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PU.unit"`
}
type ActivePowerPerFrequency struct {
	Unit                  string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerFrequency.unit"`
	Multiplier            string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerFrequency.multiplier"`
	Value                 float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerFrequency.value"`
	DenominatorUnit       string  `bun:"denominator_unit" json:"denominator_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerFrequency.denominatorUnit"`
	DenominatorMultiplier string  `bun:"denominator_multiplier" json:"denominator_multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerFrequency.denominatorMultiplier"`
}
type RegularIntervalSchedule struct {
	BasicIntervalSchedule
	EndTime  time.Time `bun:"end_time" json:"end_time" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegularIntervalSchedule.endTime"`
	TimeStep float64   `bun:"time_step" json:"time_step" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RegularIntervalSchedule.timeStep"`
}
type NonlinearShuntCompensatorPoint struct {
	SectionNumber                 int       `bun:"section_number" json:"section_number" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#NonlinearShuntCompensatorPoint.sectionNumber"`
	B                             float64   `bun:"b" json:"b" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#NonlinearShuntCompensatorPoint.b"`
	NonlinearShuntCompensatorMrid uuid.UUID `bun:"nonlinear_shunt_compensator_mrid,type:uuid" json:"nonlinear_shunt_compensator_mrid"`
	NonlinearShuntCompensator     *Entity   `bun:"rel:belongs-to,join:nonlinear_shunt_compensator_mrid=mrid" json:"nonlinear_shunt_compensator,omitempty"`
	G                             float64   `bun:"g" json:"g" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#NonlinearShuntCompensatorPoint.g"`
}
type BasicIntervalSchedule struct {
	IdentifiedObject
	Value2Unit string    `bun:"value2_unit" json:"value2_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#BasicIntervalSchedule.value2Unit"`
	Value1Unit string    `bun:"value1_unit" json:"value1_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#BasicIntervalSchedule.value1Unit"`
	StartTime  time.Time `bun:"start_time" json:"start_time" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#BasicIntervalSchedule.startTime"`
}
type EnergyConsumer struct {
	ConductingEquipment
	LoadResponseMrid uuid.UUID `bun:"load_response_mrid,type:uuid" json:"load_response_mrid"`
	LoadResponse     *Entity   `bun:"rel:belongs-to,join:load_response_mrid=mrid" json:"load_response,omitempty"`
}
type HydroPump struct {
	Equipment
	HydroPowerPlantMrid uuid.UUID `bun:"hydro_power_plant_mrid,type:uuid" json:"hydro_power_plant_mrid"`
	HydroPowerPlant     *Entity   `bun:"rel:belongs-to,join:hydro_power_plant_mrid=mrid" json:"hydro_power_plant,omitempty"`
	RotatingMachineMrid uuid.UUID `bun:"rotating_machine_mrid,type:uuid" json:"rotating_machine_mrid"`
	RotatingMachine     *Entity   `bun:"rel:belongs-to,join:rotating_machine_mrid=mrid" json:"rotating_machine,omitempty"`
}
type ControlArea struct {
	PowerSystemResource
	Type string `bun:"type" json:"type" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ControlArea.type"`
}
type DCConverterUnit struct {
	DCEquipmentContainer
	SubstationMrid uuid.UUID `bun:"substation_mrid,type:uuid" json:"substation_mrid"`
	Substation     *Entity   `bun:"rel:belongs-to,join:substation_mrid=mrid" json:"substation,omitempty"`
	OperationMode  string    `bun:"operation_mode" json:"operation_mode" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCConverterUnit.operationMode"`
}
type RotationSpeed struct {
	Multiplier            string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotationSpeed.multiplier"`
	Unit                  string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotationSpeed.unit"`
	Value                 float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotationSpeed.value"`
	DenominatorUnit       string  `bun:"denominator_unit" json:"denominator_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotationSpeed.denominatorUnit"`
	DenominatorMultiplier string  `bun:"denominator_multiplier" json:"denominator_multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotationSpeed.denominatorMultiplier"`
}
type SubGeographicalRegion struct {
	IdentifiedObject
	RegionMrid uuid.UUID `bun:"region_mrid,type:uuid" json:"region_mrid"`
	Region     *Entity   `bun:"rel:belongs-to,join:region_mrid=mrid" json:"region,omitempty"`
}
type ConformLoadSchedule struct {
	SeasonDayTypeSchedule
	ConformLoadGroupMrid uuid.UUID `bun:"conform_load_group_mrid,type:uuid" json:"conform_load_group_mrid"`
	ConformLoadGroup     *Entity   `bun:"rel:belongs-to,join:conform_load_group_mrid=mrid" json:"conform_load_group,omitempty"`
}
type PhaseTapChangerTablePoint struct {
	TapChangerTablePoint
	Angle                    float64   `bun:"angle" json:"angle" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseTapChangerTablePoint.angle"`
	PhaseTapChangerTableMrid uuid.UUID `bun:"phase_tap_changer_table_mrid,type:uuid" json:"phase_tap_changer_table_mrid"`
	PhaseTapChangerTable     *Entity   `bun:"rel:belongs-to,join:phase_tap_changer_table_mrid=mrid" json:"phase_tap_changer_table,omitempty"`
}
type ActivePowerPerCurrentFlow struct {
	Value                 float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerCurrentFlow.value"`
	Multiplier            string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerCurrentFlow.multiplier"`
	Unit                  string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerCurrentFlow.unit"`
	DenominatorUnit       string  `bun:"denominator_unit" json:"denominator_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerCurrentFlow.denominatorUnit"`
	DenominatorMultiplier string  `bun:"denominator_multiplier" json:"denominator_multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ActivePowerPerCurrentFlow.denominatorMultiplier"`
}
type Money struct {
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Money.multiplier"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Money.value"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Money.unit"`
}
type TieFlow struct {
	PositiveFlowIn  bool      `bun:"positive_flow_in" json:"positive_flow_in" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#TieFlow.positiveFlowIn"`
	TerminalMrid    uuid.UUID `bun:"terminal_mrid,type:uuid" json:"terminal_mrid"`
	Terminal        *Entity   `bun:"rel:belongs-to,join:terminal_mrid=mrid" json:"terminal,omitempty"`
	ControlAreaMrid uuid.UUID `bun:"control_area_mrid,type:uuid" json:"control_area_mrid"`
	ControlArea     *Entity   `bun:"rel:belongs-to,join:control_area_mrid=mrid" json:"control_area,omitempty"`
}
type StaticVarCompensator struct {
	RegulatingCondEq
	InductiveRating  float64 `bun:"inductive_rating" json:"inductive_rating" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.inductiveRating"`
	VoltageSetPoint  float64 `bun:"voltage_set_point" json:"voltage_set_point" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.voltageSetPoint"`
	SVCControlMode   string  `bun:"svccontrol_mode" json:"svccontrol_mode" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.sVCControlMode"`
	CapacitiveRating float64 `bun:"capacitive_rating" json:"capacitive_rating" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.capacitiveRating"`
	Slope            float64 `bun:"slope" json:"slope" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#StaticVarCompensator.slope"`
}
type PerCent struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerCent.unit"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerCent.multiplier"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerCent.value"`
}
type HydroGeneratingUnit struct {
	GeneratingUnit
	HydroPowerPlantMrid        uuid.UUID `bun:"hydro_power_plant_mrid,type:uuid" json:"hydro_power_plant_mrid"`
	HydroPowerPlant            *Entity   `bun:"rel:belongs-to,join:hydro_power_plant_mrid=mrid" json:"hydro_power_plant,omitempty"`
	EnergyConversionCapability string    `bun:"energy_conversion_capability" json:"energy_conversion_capability" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#HydroGeneratingUnit.energyConversionCapability"`
}
type EnergySource struct {
	ConductingEquipment
	EnergySchedulingTypeMrid uuid.UUID `bun:"energy_scheduling_type_mrid,type:uuid" json:"energy_scheduling_type_mrid"`
	EnergySchedulingType     *Entity   `bun:"rel:belongs-to,join:energy_scheduling_type_mrid=mrid" json:"energy_scheduling_type,omitempty"`
	VoltageMagnitude         float64   `bun:"voltage_magnitude" json:"voltage_magnitude" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.voltageMagnitude"`
	NominalVoltage           float64   `bun:"nominal_voltage" json:"nominal_voltage" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.nominalVoltage"`
	X                        float64   `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.x"`
	Rn                       float64   `bun:"rn" json:"rn" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.rn"`
	X0                       float64   `bun:"x0" json:"x0" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.x0"`
	R0                       float64   `bun:"r0" json:"r0" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.r0"`
	Xn                       float64   `bun:"xn" json:"xn" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.xn"`
	R                        float64   `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.r"`
	VoltageAngle             float64   `bun:"voltage_angle" json:"voltage_angle" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EnergySource.voltageAngle"`
}
type RotatingMachine struct {
	RegulatingCondEq
	RatedPowerFactor   float64   `bun:"rated_power_factor" json:"rated_power_factor" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotatingMachine.ratedPowerFactor"`
	GeneratingUnitMrid uuid.UUID `bun:"generating_unit_mrid,type:uuid" json:"generating_unit_mrid"`
	GeneratingUnit     *Entity   `bun:"rel:belongs-to,join:generating_unit_mrid=mrid" json:"generating_unit,omitempty"`
	RatedS             float64   `bun:"rated_s" json:"rated_s" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotatingMachine.ratedS"`
	RatedU             float64   `bun:"rated_u" json:"rated_u" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RotatingMachine.ratedU"`
}
type EquivalentInjection struct {
	EquivalentEquipment
	MaxP                        float64   `bun:"max_p" json:"max_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.maxP"`
	RegulationCapability        bool      `bun:"regulation_capability" json:"regulation_capability" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.regulationCapability"`
	ReactiveCapabilityCurveMrid uuid.UUID `bun:"reactive_capability_curve_mrid,type:uuid" json:"reactive_capability_curve_mrid"`
	ReactiveCapabilityCurve     *Entity   `bun:"rel:belongs-to,join:reactive_capability_curve_mrid=mrid" json:"reactive_capability_curve,omitempty"`
	MaxQ                        float64   `bun:"max_q" json:"max_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.maxQ"`
	MinQ                        float64   `bun:"min_q" json:"min_q" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.minQ"`
	MinP                        float64   `bun:"min_p" json:"min_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#EquivalentInjection.minP"`
}
type Curve struct {
	IdentifiedObject
	CurveStyle string `bun:"curve_style" json:"curve_style" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Curve.curveStyle"`
	Y2Unit     string `bun:"y2_unit" json:"y2_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Curve.y2Unit"`
	Y1Unit     string `bun:"y1_unit" json:"y1_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Curve.y1Unit"`
	XUnit      string `bun:"xunit" json:"xunit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Curve.xUnit"`
}
type PerLengthDCLineParameter struct {
	Resistance  float64 `bun:"resistance" json:"resistance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerLengthDCLineParameter.resistance"`
	Capacitance float64 `bun:"capacitance" json:"capacitance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerLengthDCLineParameter.capacitance"`
	Inductance  float64 `bun:"inductance" json:"inductance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PerLengthDCLineParameter.inductance"`
}
type RegulatingCondEq struct {
	ConductingEquipment
	RegulatingControlMrid uuid.UUID `bun:"regulating_control_mrid,type:uuid" json:"regulating_control_mrid"`
	RegulatingControl     *Entity   `bun:"rel:belongs-to,join:regulating_control_mrid=mrid" json:"regulating_control,omitempty"`
}
type Switch struct {
	ConductingEquipment
	Retained     bool    `bun:"retained" json:"retained" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Switch.retained"`
	NormalOpen   bool    `bun:"normal_open" json:"normal_open" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Switch.normalOpen"`
	RatedCurrent float64 `bun:"rated_current" json:"rated_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Switch.ratedCurrent"`
}
type BusNameMarker struct {
	IdentifiedObject
	Priority           int       `bun:"priority" json:"priority" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#BusNameMarker.priority"`
	ReportingGroupMrid uuid.UUID `bun:"reporting_group_mrid,type:uuid" json:"reporting_group_mrid"`
	ReportingGroup     *Entity   `bun:"rel:belongs-to,join:reporting_group_mrid=mrid" json:"reporting_group,omitempty"`
}
type ACLineSegment struct {
	Conductor
	Bch float64 `bun:"bch" json:"bch" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACLineSegment.bch"`
	R   float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACLineSegment.r"`
	X   float64 `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACLineSegment.x"`
	Gch float64 `bun:"gch" json:"gch" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ACLineSegment.gch"`
}
type SeriesCompensator struct {
	ConductingEquipment
	VaristorRatedCurrent     float64 `bun:"varistor_rated_current" json:"varistor_rated_current" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.varistorRatedCurrent"`
	VaristorPresent          bool    `bun:"varistor_present" json:"varistor_present" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.varistorPresent"`
	VaristorVoltageThreshold float64 `bun:"varistor_voltage_threshold" json:"varistor_voltage_threshold" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.varistorVoltageThreshold"`
	X                        float64 `bun:"x" json:"x" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.x"`
	R                        float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#SeriesCompensator.r"`
}
type Conductance struct {
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Conductance.unit"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Conductance.multiplier"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Conductance.value"`
}
type LinearShuntCompensator struct {
	ShuntCompensator
	GPerSection float64 `bun:"gper_section" json:"gper_section" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LinearShuntCompensator.gPerSection"`
	BPerSection float64 `bun:"bper_section" json:"bper_section" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#LinearShuntCompensator.bPerSection"`
}
type Line struct {
	EquipmentContainer
	RegionMrid uuid.UUID `bun:"region_mrid,type:uuid" json:"region_mrid"`
	Region     *Entity   `bun:"rel:belongs-to,join:region_mrid=mrid" json:"region,omitempty"`
}
type CurrentLimit struct {
	OperationalLimit
	Value float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CurrentLimit.value"`
}
type Reactance struct {
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Reactance.multiplier"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Reactance.unit"`
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Reactance.value"`
}
type CsConverter struct {
	ACDCConverter
	RatedIdc float64 `bun:"rated_idc" json:"rated_idc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.ratedIdc"`
	MaxGamma float64 `bun:"max_gamma" json:"max_gamma" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.maxGamma"`
	MaxIdc   float64 `bun:"max_idc" json:"max_idc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.maxIdc"`
	MinAlpha float64 `bun:"min_alpha" json:"min_alpha" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.minAlpha"`
	MinIdc   float64 `bun:"min_idc" json:"min_idc" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.minIdc"`
	MaxAlpha float64 `bun:"max_alpha" json:"max_alpha" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.maxAlpha"`
	MinGamma float64 `bun:"min_gamma" json:"min_gamma" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CsConverter.minGamma"`
}
type GeneratingUnit struct {
	Equipment
	RatedNetMaxP                    float64 `bun:"rated_net_max_p" json:"rated_net_max_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.ratedNetMaxP"`
	StartupCost                     float64 `bun:"startup_cost" json:"startup_cost" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.startupCost"`
	InitialP                        float64 `bun:"initial_p" json:"initial_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.initialP"`
	ShortPF                         float64 `bun:"short_pf" json:"short_pf" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.shortPF"`
	MaxOperatingP                   float64 `bun:"max_operating_p" json:"max_operating_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.maxOperatingP"`
	MinOperatingP                   float64 `bun:"min_operating_p" json:"min_operating_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.minOperatingP"`
	NominalP                        float64 `bun:"nominal_p" json:"nominal_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.nominalP"`
	RatedGrossMaxP                  float64 `bun:"rated_gross_max_p" json:"rated_gross_max_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.ratedGrossMaxP"`
	MaximumAllowableSpinningReserve float64 `bun:"maximum_allowable_spinning_reserve" json:"maximum_allowable_spinning_reserve" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.maximumAllowableSpinningReserve"`
	VariableCost                    float64 `bun:"variable_cost" json:"variable_cost" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.variableCost"`
	TotalEfficiency                 float64 `bun:"total_efficiency" json:"total_efficiency" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.totalEfficiency"`
	GenControlSource                string  `bun:"gen_control_source" json:"gen_control_source" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.genControlSource"`
	LongPF                          float64 `bun:"long_pf" json:"long_pf" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.longPF"`
	RatedGrossMinP                  float64 `bun:"rated_gross_min_p" json:"rated_gross_min_p" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.ratedGrossMinP"`
	GovernorSCD                     float64 `bun:"governor_scd" json:"governor_scd" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratingUnit.governorSCD"`
}
type RatioTapChangerTablePoint struct {
	TapChangerTablePoint
	RatioTapChangerTableMrid uuid.UUID `bun:"ratio_tap_changer_table_mrid,type:uuid" json:"ratio_tap_changer_table_mrid"`
	RatioTapChangerTable     *Entity   `bun:"rel:belongs-to,join:ratio_tap_changer_table_mrid=mrid" json:"ratio_tap_changer_table,omitempty"`
}
type RatioTapChanger struct {
	TapChanger
	RatioTapChangerTableMrid uuid.UUID `bun:"ratio_tap_changer_table_mrid,type:uuid" json:"ratio_tap_changer_table_mrid"`
	RatioTapChangerTable     *Entity   `bun:"rel:belongs-to,join:ratio_tap_changer_table_mrid=mrid" json:"ratio_tap_changer_table,omitempty"`
	TransformerEndMrid       uuid.UUID `bun:"transformer_end_mrid,type:uuid" json:"transformer_end_mrid"`
	TransformerEnd           *Entity   `bun:"rel:belongs-to,join:transformer_end_mrid=mrid" json:"transformer_end,omitempty"`
	StepVoltageIncrement     float64   `bun:"step_voltage_increment" json:"step_voltage_increment" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RatioTapChanger.stepVoltageIncrement"`
	TculControlMode          string    `bun:"tcul_control_mode" json:"tcul_control_mode" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#RatioTapChanger.tculControlMode"`
}
type VoltageLevel struct {
	EquipmentContainer
	BaseVoltageMrid  uuid.UUID `bun:"base_voltage_mrid,type:uuid" json:"base_voltage_mrid"`
	BaseVoltage      *Entity   `bun:"rel:belongs-to,join:base_voltage_mrid=mrid" json:"base_voltage,omitempty"`
	LowVoltageLimit  float64   `bun:"low_voltage_limit" json:"low_voltage_limit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltageLevel.lowVoltageLimit"`
	HighVoltageLimit float64   `bun:"high_voltage_limit" json:"high_voltage_limit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltageLevel.highVoltageLimit"`
	SubstationMrid   uuid.UUID `bun:"substation_mrid,type:uuid" json:"substation_mrid"`
	Substation       *Entity   `bun:"rel:belongs-to,join:substation_mrid=mrid" json:"substation,omitempty"`
}
type ResistancePerLength struct {
	DenominatorUnit       string  `bun:"denominator_unit" json:"denominator_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ResistancePerLength.denominatorUnit"`
	Multiplier            string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ResistancePerLength.multiplier"`
	Unit                  string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ResistancePerLength.unit"`
	Value                 float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ResistancePerLength.value"`
	DenominatorMultiplier string  `bun:"denominator_multiplier" json:"denominator_multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#ResistancePerLength.denominatorMultiplier"`
}
type Conductor struct {
	ConductingEquipment
	Length float64 `bun:"length" json:"length" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Conductor.length"`
}
type VoltagePerReactivePower struct {
	DenominatorUnit       string  `bun:"denominator_unit" json:"denominator_unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltagePerReactivePower.denominatorUnit"`
	Multiplier            string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltagePerReactivePower.multiplier"`
	Unit                  string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltagePerReactivePower.unit"`
	Value                 float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltagePerReactivePower.value"`
	DenominatorMultiplier string  `bun:"denominator_multiplier" json:"denominator_multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#VoltagePerReactivePower.denominatorMultiplier"`
}
type DCLineSegment struct {
	DCConductingEquipment
	Length                 float64   `bun:"length" json:"length" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCLineSegment.length"`
	Inductance             float64   `bun:"inductance" json:"inductance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCLineSegment.inductance"`
	PerLengthParameterMrid uuid.UUID `bun:"per_length_parameter_mrid,type:uuid" json:"per_length_parameter_mrid"`
	PerLengthParameter     *Entity   `bun:"rel:belongs-to,join:per_length_parameter_mrid=mrid" json:"per_length_parameter,omitempty"`
	Capacitance            float64   `bun:"capacitance" json:"capacitance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCLineSegment.capacitance"`
	Resistance             float64   `bun:"resistance" json:"resistance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCLineSegment.resistance"`
}
type OperationalLimitType struct {
	IdentifiedObject
	LimitType          string  `bun:"limit_type" json:"limit_type" iri:"http://entsoe.eu/CIM/SchemaExtension/3/1#OperationalLimitType.limitType"`
	Direction          string  `bun:"direction" json:"direction" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#OperationalLimitType.direction"`
	AcceptableDuration float64 `bun:"acceptable_duration" json:"acceptable_duration" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#OperationalLimitType.acceptableDuration"`
}
type Inductance struct {
	Value      float64 `bun:"value" json:"value" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Inductance.value"`
	Unit       string  `bun:"unit" json:"unit" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Inductance.unit"`
	Multiplier string  `bun:"multiplier" json:"multiplier" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#Inductance.multiplier"`
}
type DCGround struct {
	DCConductingEquipment
	Inductance float64 `bun:"inductance" json:"inductance" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCGround.inductance"`
	R          float64 `bun:"r" json:"r" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#DCGround.r"`
}
type Substation struct {
	EquipmentContainer
	RegionMrid uuid.UUID `bun:"region_mrid,type:uuid" json:"region_mrid"`
	Region     *Entity   `bun:"rel:belongs-to,join:region_mrid=mrid" json:"region,omitempty"`
}
type PowerSystemResource struct {
	IdentifiedObject
}
type DCConductingEquipment struct {
	Equipment
}
type DCDisconnector struct {
	DCSwitch
}
type DCSwitch struct {
	DCConductingEquipment
}
type ConnectivityNodeContainer struct {
	PowerSystemResource
}
type DCBusbar struct {
	DCConductingEquipment
}
type PowerTransformer struct {
	ConductingEquipment
}
type EnergySchedulingType struct {
	IdentifiedObject
}
type ThermalGeneratingUnit struct {
	GeneratingUnit
}
type NonlinearShuntCompensator struct {
	ShuntCompensator
}
type NonConformLoadGroup struct {
	LoadGroup
}
type LoadGroup struct {
	IdentifiedObject
}
type Connector struct {
	ConductingEquipment
}
type DCBreaker struct {
	DCSwitch
}
type EquipmentContainer struct {
	ConnectivityNodeContainer
}
type EquivalentNetwork struct {
	ConnectivityNodeContainer
}
type BusbarSection struct {
	Connector
}
type PhaseTapChangerTable struct {
	IdentifiedObject
}
type NuclearGeneratingUnit struct {
	GeneratingUnit
}
type SeasonDayTypeSchedule struct {
}
type Junction struct {
	Connector
}
type GeographicalRegion struct {
	IdentifiedObject
}
type DCEquipmentContainer struct {
	EquipmentContainer
}
type ProtectedSwitch struct {
	Switch
}
type ReportingGroup struct {
	IdentifiedObject
}
type DCChopper struct {
	DCConductingEquipment
}
type Disconnector struct {
	Switch
}
type PhaseTapChangerSymmetrical struct {
	PhaseTapChangerNonLinear
}
type ConformLoadGroup struct {
	LoadGroup
}
type VsCapabilityCurve struct {
	Curve
}
type LoadBreakSwitch struct {
	ProtectedSwitch
}
type SolarGeneratingUnit struct {
	GeneratingUnit
}
type RatioTapChangerTable struct {
	IdentifiedObject
}
type TapChangerControl struct {
	RegulatingControl
}
type ReactiveCapabilityCurve struct {
	Curve
}
type Breaker struct {
	ProtectedSwitch
}
