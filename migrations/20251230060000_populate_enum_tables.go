package migrations

import (
	"context"
	"fmt"
	"strings"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(populateEnumTables, revertPopulateEnumTables)
}

func populateEnumTables(ctx context.Context, db *bun.DB) error {
	enumData := []models.RdfsEnum{
		// LimitTypeKind
		{Id: 1, Code: "highVoltage", Iri: "http://entsoe.eu/CIM/SchemaExtension/3/1#LimitTypeKind.highVoltage", Comment: "Referring to the rating of the equipments, a voltage too high can lead to accelerated ageing or the destruction of the equipment. This limit type may or may not have duration."},
		{Id: 2, Code: "lowVoltage", Iri: "http://entsoe.eu/CIM/SchemaExtension/3/1#LimitTypeKind.lowVoltage", Comment: "A too low voltage can disturb the normal operation of some protections and transformer equipped with on-load tap changers, electronic power devices or can affect the behaviour of the auxiliaries of generation units.This limit type may or may not have duration."},
		{Id: 3, Code: "patl", Iri: "http://entsoe.eu/CIM/SchemaExtension/3/1#LimitTypeKind.patl", Comment: "The Permanent Admissible Transmission Loading (PATL) is the loading in Amps, MVA or MW that can be accepted by a network branch for an unlimited duration without any risk for the material.The duration attribute is not used and shall be excluded for the PATL limit type. Hence only one limit value exists for the PATL type."},
		{Id: 4, Code: "patlt", Iri: "http://entsoe.eu/CIM/SchemaExtension/3/1#LimitTypeKind.patlt", Comment: "Permanent Admissible Transmission Loading Threshold  (PATLT) is a value in engineering units defined for PATL and calculated using percentage less than 100 of the PATL type intended to alert operators of an arising condition. The percentage should be given in the name of the OperationalLimitSet. The aceptableDuration is another way to express the severity of the limit."},
		{Id: 5, Code: "tatl", Iri: "http://entsoe.eu/CIM/SchemaExtension/3/1#LimitTypeKind.tatl", Comment: "Temporarily Admissible Transmission Loading (TATL) which is the loading in Amps, MVA or MW that can be accepted by a branch for a certain limited duration.The TATL can be defined in different ways:as a fixed percentage of the PATL for a given time (for example, 115% of the PATL that can be accepted during 15 minutes),pairs of TATL type and Duration calculated for each line taking into account its particular configuration and conditions of functioning (for example, it can define a TATL acceptable during 20 minutes and another one acceptable during 10 minutes).Such a definition of TATL can depend on the initial operating conditions of the network element (sag situation of a line).The duration attribute can be used define several TATL limit types. Hence multiple TATL limit values may exist having different durations."},
		{Id: 6, Code: "tc", Iri: "http://entsoe.eu/CIM/SchemaExtension/3/1#LimitTypeKind.tc", Comment: "Tripping Current (TC) is the ultimate intensity without any delay. It is defined as the threshold the line will trip without any possible remedial actions.The tripping of the network element is ordered by protections against short circuits or by overload protections, but in any case, the activation delay of these protections is not compatible with the reaction delay of an operator (less than one minute).The duration is always zero and the duration attribute may be left out. Hence only one limit value exists for the TC type."},
		{Id: 7, Code: "tct", Iri: "http://entsoe.eu/CIM/SchemaExtension/3/1#LimitTypeKind.tct", Comment: "Tripping Current Threshold  (TCT) is a value in engineering units defined for TC and calculated using percentage less than 100 of the TC type intended to alert operators of an arising condition. The percentage should be given in the name of the OperationalLimitSet. The aceptableDuration is another way to express the severity of the limit."},

		// ControlAreaTypeKind
		{Id: 1, Code: "AGC", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#ControlAreaTypeKind.AGC", Comment: "Used for automatic generation control."},
		{Id: 2, Code: "Forecast", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#ControlAreaTypeKind.Forecast", Comment: "Used for load forecast."},
		{Id: 3, Code: "Interchange", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#ControlAreaTypeKind.Interchange", Comment: "Used for interchange specification or control."},

		// DCConverterOperatingModeKind
		{Id: 1, Code: "bipolar", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#DCConverterOperatingModeKind.bipolar", Comment: "Bipolar operation."},
		{Id: 2, Code: "monopolarGroundReturn", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#DCConverterOperatingModeKind.monopolarGroundReturn", Comment: "Monopolar operation with ground return"},
		{Id: 3, Code: "monopolarMetallicReturn", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#DCConverterOperatingModeKind.monopolarMetallicReturn", Comment: "Monopolar operation with metallic return"},

		// DCPolarityKind
		{Id: 1, Code: "middle", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#DCPolarityKind.middle", Comment: "Middle pole, potentially grounded."},
		{Id: 2, Code: "negative", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#DCPolarityKind.negative", Comment: "Negative pole."},
		{Id: 3, Code: "positive", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#DCPolarityKind.positive", Comment: "Positive pole."},

		// HydroEnergyConversionKind
		{Id: 1, Code: "generator", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#HydroEnergyConversionKind.generator", Comment: "Able to generate power, but not able to pump water for energy storage."},
		{Id: 2, Code: "pumpAndGenerator", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#HydroEnergyConversionKind.pumpAndGenerator", Comment: "Able to both generate power and pump water for energy storage."},

		// HydroPlantStorageKind
		{Id: 1, Code: "pumpedStorage", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#HydroPlantStorageKind.pumpedStorage", Comment: "Pumped storage."},
		{Id: 2, Code: "runOfRiver", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#HydroPlantStorageKind.runOfRiver", Comment: "Run of river."},
		{Id: 3, Code: "storage", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#HydroPlantStorageKind.storage", Comment: "Storage."},

		// OperationalLimitDirectionKind
		{Id: 1, Code: "absoluteValue", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#OperationalLimitDirectionKind.absoluteValue", Comment: "An absoluteValue limit means that a monitored absolute value above the limit value is a violation."},
		{Id: 2, Code: "high", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#OperationalLimitDirectionKind.high", Comment: "High means that a monitored value above the limit value is a violation.   If applied to a terminal flow, the positive direction is into the terminal."},
		{Id: 3, Code: "low", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#OperationalLimitDirectionKind.low", Comment: "Low means a monitored value below the limit is a violation.  If applied to a terminal flow, the positive direction is into the terminal."},

		// PetersenCoilModeKind
		{Id: 1, Code: "automaticPositioning", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PetersenCoilModeKind.automaticPositioning", Comment: "Automatic positioning."},
		{Id: 2, Code: "fixed", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PetersenCoilModeKind.fixed", Comment: "Fixed position."},
		{Id: 3, Code: "manual", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PetersenCoilModeKind.manual", Comment: "Manual positioning."},

		// RegulatingControlModeKind
		{Id: 1, Code: "activePower", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControlModeKind.activePower", Comment: "Active power is specified."},
		{Id: 2, Code: "admittance", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControlModeKind.admittance", Comment: "Admittance is specified."},
		{Id: 3, Code: "currentFlow", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControlModeKind.currentFlow", Comment: "Current flow is specified."},
		{Id: 4, Code: "powerFactor", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControlModeKind.powerFactor", Comment: "Power factor is specified."},
		{Id: 5, Code: "reactivePower", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControlModeKind.reactivePower", Comment: "Reactive power is specified."},
		{Id: 6, Code: "temperature", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControlModeKind.temperature", Comment: "Control switches on/off based on the local temperature (i.e., a thermostat)."},
		{Id: 7, Code: "timeScheduled", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControlModeKind.timeScheduled", Comment: "Control switches on/off by time of day. The times may change on the weekend, or in different seasons."},
		{Id: 8, Code: "voltage", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#RegulatingControlModeKind.voltage", Comment: "Voltage is specified."},

		// ShortCircuitRotorKind
		{Id: 1, Code: "salientPole1", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#ShortCircuitRotorKind.salientPole1", Comment: "Salient pole 1 in the IEC 60909"},
		{Id: 2, Code: "salientPole2", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#ShortCircuitRotorKind.salientPole2", Comment: "Salient pole 2 in IEC 60909"},
		{Id: 3, Code: "turboSeries1", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#ShortCircuitRotorKind.turboSeries1", Comment: "Turbo Series 1 in the IEC 60909"},
		{Id: 4, Code: "turboSeries2", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#ShortCircuitRotorKind.turboSeries2", Comment: "Turbo series 2 in IEC 60909"},

		// SynchronousMachineKind
		{Id: 1, Code: "condenser", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachineKind.condenser", Comment: ""},
		{Id: 2, Code: "generator", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachineKind.generator", Comment: ""},
		{Id: 3, Code: "generatorOrCondenser", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachineKind.generatorOrCondenser", Comment: ""},
		{Id: 4, Code: "generatorOrCondenserOrMotor", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachineKind.generatorOrCondenserOrMotor", Comment: ""},
		{Id: 5, Code: "generatorOrMotor", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachineKind.generatorOrMotor", Comment: ""},
		{Id: 6, Code: "motor", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachineKind.motor", Comment: ""},
		{Id: 7, Code: "motorOrCondenser", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SynchronousMachineKind.motorOrCondenser", Comment: ""},

		// WindGenUnitKind
		{Id: 1, Code: "offshore", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindGenUnitKind.offshore", Comment: "The wind generating unit is located offshore."},
		{Id: 2, Code: "onshore", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindGenUnitKind.onshore", Comment: "The wind generating unit is located onshore."},
	}

	// Clean up comments
	for i := range enumData {
		comment := strings.TrimSuffix(enumData[i].Comment, "^^<http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral>")
		comment = strings.Trim(comment, `"`)
		enumData[i].Comment = comment
	}

	enums := []any{
		&models.LimitTypeKind{RdfsEnum: enumData[0]},
		&models.LimitTypeKind{RdfsEnum: enumData[1]},
		&models.LimitTypeKind{RdfsEnum: enumData[2]},
		&models.LimitTypeKind{RdfsEnum: enumData[3]},
		&models.LimitTypeKind{RdfsEnum: enumData[4]},
		&models.LimitTypeKind{RdfsEnum: enumData[5]},
		&models.LimitTypeKind{RdfsEnum: enumData[6]},

		&models.ControlAreaTypeKind{RdfsEnum: enumData[7]},
		&models.ControlAreaTypeKind{RdfsEnum: enumData[8]},
		&models.ControlAreaTypeKind{RdfsEnum: enumData[9]},

		&models.DCConverterOperatingModeKind{RdfsEnum: enumData[10]},
		&models.DCConverterOperatingModeKind{RdfsEnum: enumData[11]},
		&models.DCConverterOperatingModeKind{RdfsEnum: enumData[12]},

		&models.DCPolarityKind{RdfsEnum: enumData[13]},
		&models.DCPolarityKind{RdfsEnum: enumData[14]},
		&models.DCPolarityKind{RdfsEnum: enumData[15]},

		&models.HydroEnergyConversionKind{RdfsEnum: enumData[16]},
		&models.HydroEnergyConversionKind{RdfsEnum: enumData[17]},

		&models.HydroPlantStorageKind{RdfsEnum: enumData[18]},
		&models.HydroPlantStorageKind{RdfsEnum: enumData[19]},
		&models.HydroPlantStorageKind{RdfsEnum: enumData[20]},

		&models.OperationalLimitDirectionKind{RdfsEnum: enumData[21]},
		&models.OperationalLimitDirectionKind{RdfsEnum: enumData[22]},
		&models.OperationalLimitDirectionKind{RdfsEnum: enumData[23]},

		&models.PetersenCoilModeKind{RdfsEnum: enumData[24]},
		&models.PetersenCoilModeKind{RdfsEnum: enumData[25]},
		&models.PetersenCoilModeKind{RdfsEnum: enumData[26]},

		&models.RegulatingControlModeKind{RdfsEnum: enumData[27]},
		&models.RegulatingControlModeKind{RdfsEnum: enumData[28]},
		&models.RegulatingControlModeKind{RdfsEnum: enumData[29]},
		&models.RegulatingControlModeKind{RdfsEnum: enumData[30]},
		&models.RegulatingControlModeKind{RdfsEnum: enumData[31]},
		&models.RegulatingControlModeKind{RdfsEnum: enumData[32]},
		&models.RegulatingControlModeKind{RdfsEnum: enumData[33]},
		&models.RegulatingControlModeKind{RdfsEnum: enumData[34]},

		&models.ShortCircuitRotorKind{RdfsEnum: enumData[35]},
		&models.ShortCircuitRotorKind{RdfsEnum: enumData[36]},
		&models.ShortCircuitRotorKind{RdfsEnum: enumData[37]},
		&models.ShortCircuitRotorKind{RdfsEnum: enumData[38]},

		&models.SynchronousMachineKind{RdfsEnum: enumData[39]},
		&models.SynchronousMachineKind{RdfsEnum: enumData[40]},
		&models.SynchronousMachineKind{RdfsEnum: enumData[41]},
		&models.SynchronousMachineKind{RdfsEnum: enumData[42]},
		&models.SynchronousMachineKind{RdfsEnum: enumData[43]},
		&models.SynchronousMachineKind{RdfsEnum: enumData[44]},
		&models.SynchronousMachineKind{RdfsEnum: enumData[45]},

		&models.WindGenUnitKind{RdfsEnum: enumData[46]},
		&models.WindGenUnitKind{RdfsEnum: enumData[47]},
	}

	for _, enum := range enums {
		_, err := db.NewInsert().Model(enum).Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to insert enum: %w", err)
		}
	}

	return nil
}

func revertPopulateEnumTables(ctx context.Context, db *bun.DB) error {
	enumTables := []string{
		"limit_type_kinds", "control_area_type_kinds", "dc_converter_operating_mode_kinds",
		"dc_polarity_kinds", "hydro_energy_conversion_kinds", "hydro_plant_storage_kinds",
		"operational_limit_direction_kinds", "petersen_coil_mode_kinds",
		"regulating_control_mode_kinds", "short_circuit_rotor_kinds",
		"synchronous_machine_kinds", "wind_gen_unit_kinds",
	}

	for _, table := range enumTables {
		_, err := db.NewDelete().
			Table(table).
			Where("1=1").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to clear table %s: %w", table, err)
		}
	}

	return nil
}
