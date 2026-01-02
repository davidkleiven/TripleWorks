package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addUnitMultiplierEnum, revertAddUnitMultiplierEnum)
}

func addUnitMultiplierEnum(ctx context.Context, db *bun.DB) error {
	data := []models.UnitMultiplier{
		{RdfsEnum: models.RdfsEnum{
			Id:      1,
			Iri:     "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.G",
			Code:    "G",
			Comment: "Giga 10**9.",
		}},
		{RdfsEnum: models.RdfsEnum{
			Id:      2,
			Iri:     "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.M",
			Code:    "M",
			Comment: "Mega 10**6.",
		}},
		{RdfsEnum: models.RdfsEnum{
			Id:      3,
			Iri:     "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.T",
			Code:    "T",
			Comment: "Tera 10**12."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.c", Code: "c",
			Id:      4,
			Comment: "Centi 10**-2."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.d", Code: "d",
			Id:      5,
			Comment: "Deci 10**-1."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.k", Code: "k",
			Id:      6,
			Comment: "Kilo 10**3."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.m", Code: "m",
			Id:      7,
			Comment: "Milli 10**-3."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.micro", Code: "micro",
			Id:      8,
			Comment: "Micro 10**-6."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.n", Code: "n",
			Id:      9,
			Comment: "Nano 10**-9."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.none", Code: "none",
			Id:      10,
			Comment: "No multiplier or equivalently multiply by 1."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#UnitMultiplier.p", Code: "p",
			Id:      11,
			Comment: "Pico 10**-12."}},
	}
	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddUnitMultiplierEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.UnitMultiplier{}).Where("1=1").Exec(ctx)
	return err
}
