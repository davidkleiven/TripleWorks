package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addSvcControlModeEnum, revertAddSvcControlModeEnum)
}

func addSvcControlModeEnum(ctx context.Context, db *bun.DB) error {
	data := []models.SVCControlMode{
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SVCControlMode.reactivePower", Code: "reactivePower", Id: 1}},
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#SVCControlMode.voltage", Code: "voltage", Id: 2}},
	}
	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddSvcControlModeEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.SVCControlMode{}).Where("1=1").Exec(ctx)
	return err
}
