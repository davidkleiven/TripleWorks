package migrations

import (
	"context"
	"fmt"
	"reflect"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

var viewTables = []any{
	&models.BasicIntervalSchedule{},
	&models.Curve{},
	&models.LoadResponseCharacteristic{},
	&models.Location{},
	&models.PhaseTapChangerTable{},
	&models.PhaseTapChangerTabular{},
	&models.RatioTapChangerTable{},
	&models.ReactiveCapabilityCurve{},
	&models.RegularIntervalSchedule{},
	&models.VoltageLimit{},
	&models.VsCapabilityCurve{},
}

func init() {
	migrations.MustRegister(addMoreLatestViews, revertAddMoreLatestView)
}

func addMoreLatestViews(ctx context.Context, db *bun.DB) error {
	for i, item := range viewTables {
		name := db.Table(reflect.TypeOf(item).Elem()).Name
		sql := MustGetViewSql(name)
		_, err := db.ExecContext(ctx, sql)
		if err != nil {
			return fmt.Errorf("Failed for %d (%s): %w", i, name, err)
		}
	}
	return nil
}

func revertAddMoreLatestView(ctx context.Context, db *bun.DB) error {
	for i, item := range viewTables {
		name := db.Table(reflect.TypeOf(item).Elem()).Name
		sql := fmt.Sprintf("DROP VIEW v_%s_latest", name)
		_, err := db.ExecContext(ctx, sql)
		if err != nil {
			return fmt.Errorf("Failed for %d (%s): %w", i, name, err)
		}
	}
	return nil
}
