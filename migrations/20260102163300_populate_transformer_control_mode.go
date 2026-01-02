package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addTransformerControlModeEnum, revertAddTransformerControlModeEnum)
}

func addTransformerControlModeEnum(ctx context.Context, db *bun.DB) error {
	data := []models.TransformerControlMode{
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#TransformerControlMode.reactive", Code: "reactive",
			Id:      1,
			Comment: "Reactive power flow control"}},
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#TransformerControlMode.volt", Code: "volt",
			Id:      2,
			Comment: "Voltage control"}},
	}
	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddTransformerControlModeEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.TransformerControlMode{}).Where("1=1").Exec(ctx)
	return err
}
