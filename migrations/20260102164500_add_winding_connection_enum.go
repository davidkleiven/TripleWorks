package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addWindingConnectionEnum, revertAddWindingConnectionEnum)
}

func addWindingConnectionEnum(ctx context.Context, db *bun.DB) error {
	data := []models.WindingConnection{
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindingConnection.A", Code: "A",
			Comment: "Autotransformer common winding"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindingConnection.D", Code: "D",
			Comment: "Delta"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindingConnection.I", Code: "I",
			Comment: "Independent winding, for single-phase connections"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindingConnection.Y", Code: "Y",
			Comment: "Wye"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindingConnection.Yn", Code: "Yn",
			Comment: "Wye, with neutral brought out for grounding."}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindingConnection.Z", Code: "Z",
			Comment: "ZigZag"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#WindingConnection.Zn", Code: "Zn",
			Comment: "ZigZag, with neutral brought out for grounding."}},
	}

	for i := range data {
		data[i].Id = i + 1
	}
	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddWindingConnectionEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.WindingConnection{}).Where("1=1").Exec(ctx)
	return err
}
