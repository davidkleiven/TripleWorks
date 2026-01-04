package migrations

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(createCim16OpsTables, revertCreateCim16OpsTables)
}

var opsTables = []any{
	(*models.Accumulator)(nil),
	(*models.AccumulatorLimit)(nil),
	(*models.AccumulatorLimitSet)(nil),
	(*models.AccumulatorReset)(nil),
	(*models.AccumulatorValue)(nil),
	(*models.ActivePowerLimit)(nil),
	(*models.Analog)(nil),
	(*models.AnalogControl)(nil),
	(*models.AnalogLimit)(nil),
	(*models.AnalogLimitSet)(nil),
	(*models.AnalogValue)(nil),
	(*models.ApparentPowerLimit)(nil),
	(*models.Bay)(nil),
	(*models.Command)(nil),
	(*models.ConnectivityNode)(nil),
	(*models.Control)(nil),
	(*models.DayType)(nil),
	(*models.Discrete)(nil),
	(*models.DiscreteValue)(nil),
	(*models.EnergyArea)(nil),
	(*models.GrossToNetActivePowerCurve)(nil),
	(*models.Ground)(nil),
	(*models.GroundDisconnector)(nil),
	(*models.Limit)(nil),
	(*models.LimitSet)(nil),
	(*models.LoadArea)(nil),
	(*models.Measurement)(nil),
	(*models.MeasurementValue)(nil),
	(*models.MeasurementValueQuality)(nil),
	(*models.MeasurementValueSource)(nil),
	(*models.Quality61850)(nil),
	(*models.RaiseLowerCommand)(nil),
	(*models.RegularTimePoint)(nil),
	(*models.RegulationSchedule)(nil),
	(*models.Season)(nil),
	(*models.SetPoint)(nil),
	(*models.StationSupply)(nil),
	(*models.StringMeasurement)(nil),
	(*models.StringMeasurementValue)(nil),
	(*models.SubLoadArea)(nil),
	(*models.SwitchSchedule)(nil),
	(*models.TapSchedule)(nil),
	(*models.ValueAliasSet)(nil),
	(*models.ValueToAlias)(nil),
}

func createCim16OpsTables(ctx context.Context, db *bun.DB) error {
	for i, table := range opsTables {
		_, err := db.NewCreateTable().
			Model(table).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to create table no. %d: %w", i, err)
		}
	}
	return nil
}

func revertCreateCim16OpsTables(ctx context.Context, db *bun.DB) error {
	for i := len(opsTables); i > 0; i-- {
		_, err := db.NewDropTable().Model(opsTables[i-1]).IfExists().Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to drop table no. %d: %w", i, err)
		}
	}
	return nil
}
