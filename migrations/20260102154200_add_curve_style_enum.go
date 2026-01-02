package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addCurveStyleEnum, revertAddCurveStyleEnum)
}

func addCurveStyleEnum(ctx context.Context, db *bun.DB) error {
	data := []models.CurveStyle{
		{RdfsEnum: models.RdfsEnum{
			Id:      1,
			Code:    "constantYValue",
			Comment: "The Y-axis values are assumed constant until the next curve point and prior to the first curve point.",
			Iri:     "http://iec.ch/TC57/2013/CIM-schema-cim16#CurveStyle.constantYValue",
		}},
		{RdfsEnum: models.RdfsEnum{
			Id:      2,
			Code:    "straightLineYValues",
			Comment: "The Y-axis values are assumed to be a straight line between values.  Also known as linear interpolation.",
			Iri:     "http://iec.ch/TC57/2013/CIM-schema-cim16#CurveStyle.straightLineYValues",
		}},
	}

	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddCurveStyleEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.CurveStyle{}).Where("1=1").Exec(ctx)
	return err
}
