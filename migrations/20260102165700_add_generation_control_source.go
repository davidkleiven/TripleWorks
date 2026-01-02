package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addGeneratorControlSourceEnum, revertAddGeneratorControlSourceEnum)
}

func addGeneratorControlSourceEnum(ctx context.Context, db *bun.DB) error {
	data := []models.GeneratorControlSource{
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratorControlSource.offAGC", Code: "offAGC",
			Comment: "Off of automatic generation control (AGC)."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratorControlSource.onAGC", Code: "onAGC",
			Comment: "On automatic generation control (AGC)."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratorControlSource.plantControl", Code: "plantControl",
			Comment: "Plant is controlling."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#GeneratorControlSource.unavailable", Code: "unavailable",
			Comment: "Not available."}},
	}

	for i := range data {
		data[i].Id = i + 1
	}
	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddGeneratorControlSourceEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.GeneratorControlSource{}).Where("1=1").Exec(ctx)
	return err
}
