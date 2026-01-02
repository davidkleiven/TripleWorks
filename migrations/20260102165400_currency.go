package migrations

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(addCurrencyEnum, revertAddCurrencyEnum)
}

func addCurrencyEnum(ctx context.Context, db *bun.DB) error {
	data := []models.Currency{
		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.AUD", Code: "AUD",
			Comment: "Australian dollar"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.CAD", Code: "CAD",
			Comment: "Canadian dollar"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.CHF", Code: "CHF",
			Comment: "Swiss francs"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.CNY", Code: "CNY",
			Comment: "Chinese yuan renminbi"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.DKK", Code: "DKK",
			Comment: "Danish crown"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.EUR", Code: "EUR",
			Comment: "European euro"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.GBP", Code: "GBP",
			Comment: "British pound"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.INR", Code: "INR",
			Comment: "India rupees"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.JPY", Code: "JPY",
			Comment: "Japanese yen"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.NOK", Code: "NOK",
			Comment: "Norwegian crown"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.RUR", Code: "RUR",
			Comment: "Russian ruble"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.SEK", Code: "SEK",
			Comment: "Swedish crown"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.USD", Code: "USD",
			Comment: "US dollar"}},

		{RdfsEnum: models.RdfsEnum{Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#Currency.other", Code: "other",
			Comment: "Another type of currency."}},
	}

	for i := range data {
		data[i].Id = i + 1
	}
	_, err := db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func revertAddCurrencyEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().Model(&models.Currency{}).Where("1=1").Exec(ctx)
	return err
}
