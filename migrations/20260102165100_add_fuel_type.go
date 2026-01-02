package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addFuelTypeEnum, revertAddFuelTypeEnum)
}

func addFuelTypeEnum(ctx context.Context, db *bun.DB) error {
	data := []models.FuelType{
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#FuelType.coal", Code: "coal",
			Comment: "Generic coal, not including lignite type."}},
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#FuelType.gas", Code: "gas",
			Comment: "Natural gas."}},
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#FuelType.hardCoal", Code: "hardCoal",
			Comment: "Hard coal"}},
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#FuelType.lignite", Code: "lignite",
			Comment: "The fuel is lignite coal.  Note that this is a special type of coal, so the other enum of coal is reserved for hard coal types or if the exact type of coal is not known."}},
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#FuelType.oil", Code: "oil",
			Comment: "Oil."}},
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#FuelType.oilShale", Code: "oilShale",
			Comment: "Oil Shale"}},
	}

	for i := range data {
		data[i].Id = i + 1
	}
	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddFuelTypeEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.FuelType{}).Where("1=1").Exec(ctx)
	return err
}
