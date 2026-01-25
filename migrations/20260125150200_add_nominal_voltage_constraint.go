package migrations

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

func init() {
	migrations.MustRegister(createNominalVoltageConstraint, revertCreateNominalVoltageConstraint)
}

func createNominalVoltageConstraint(ctx context.Context, db *bun.DB) error {
	if db.Dialect().Name() != dialect.PG {
		return nil
	}
	query := `
		ALTER TABLE base_voltages DROP CONSTRAINT IF EXISTS nominal_voltage_positive;
		UPDATE base_voltages SET nominal_voltage = nominal_voltage + 1 WHERE nominal_voltage = 0;
		ALTER TABLE base_voltages
		ADD CONSTRAINT nominal_voltage_positive
	   	CHECK (nominal_voltage > 0);
		`
	_, err := db.Exec(query)
	return err
}

func revertCreateNominalVoltageConstraint(ctx context.Context, db *bun.DB) error {
	if db.Dialect().Name() != dialect.PG {
		return nil
	}

	query := "ALTER TABLE base_voltages DROP CONSTRAINT IF EXISTS nominal_voltage_positive"
	_, err := db.Exec(query)
	return err
}
