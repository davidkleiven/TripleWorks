package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addUnitSymbolEnum, revertAddUnitSymbolEnum)
}

func addUnitSymbolEnum(ctx context.Context, db *bun.DB) error {
	data := []models.UnitSymbol{
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.A", Code: "A",
			Comment: "Current in ampere."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.F", Code: "F",
			Comment: "Capacitance in farad."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.H", Code: "H",
			Comment: "Inductance in henry."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.Hz", Code: "Hz",
			Comment: "Frequency in hertz."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.J", Code: "J",
			Comment: "Energy in joule."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.N", Code: "N",
			Comment: "Force in newton."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.Pa", Code: "Pa",
			Comment: "Pressure in pascal (n/m2)."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.S", Code: "S",
			Comment: "Conductance in siemens."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.V", Code: "V",
			Comment: "Voltage in volt."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.VA", Code: "VA",
			Comment: "Apparent power in volt ampere."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.VAh", Code: "VAh",
			Comment: "Apparent energy in volt ampere hours."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.VAr", Code: "VAr",
			Comment: "Reactive power in volt ampere reactive."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.VArh", Code: "VArh",
			Comment: "Reactive energy in volt ampere reactive hours."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.W", Code: "W",
			Comment: "Active power in watt."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.Wh", Code: "Wh",
			Comment: "Real energy in what hours."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.deg", Code: "deg",
			Comment: "Plane angle in degrees."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.degC", Code: "degC",
			Comment: "Relative temperature in degrees Celsius. In the SI unit system the symbol is \u00BAC. Electric charge is measured in coulomb that has the unit symbol C. To distinguish degree Celsius form coulomb the symbol used in the UML is degC. Reason for not using \u00BAC is the special character \u00BA is difficult to manage in software."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.g", Code: "g",
			Comment: "Mass in gram."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.h", Code: "h",
			Comment: "Time in hours."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.m2", Code: "m2",
			Comment: "Area in square meters."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.m3", Code: "m3",
			Comment: "Volume in cubic meters."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.m", Code: "m",
			Comment: "Length in meter."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.min", Code: "min",
			Comment: "Time in minutes."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.none", Code: "none",
			Comment: "Dimension less quantity, e.g. count, per unit, etc."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.ohm", Code: "ohm",
			Comment: "Resistance in ohm."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.rad", Code: "rad",
			Comment: "Plane angle in radians."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitSymbol.s", Code: "s",
			Comment: "Time in seconds."}},
	}

	for i := range data {
		data[i].Id = i + 1
	}
	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddUnitSymbolEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.UnitSymbol{}).Where("1=1").Exec(ctx)
	return err
}
