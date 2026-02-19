package pkg

import (
	"context"
	"fmt"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestApplyPatch(t *testing.T) {
	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	ctx := context.Background()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	var bv models.BaseVoltage
	bv.Mrid = uuid.New()

	_, err = db.NewInsert().Model(&bv).Exec(ctx)
	require.NoError(t, err)

	entity := models.Entity{
		Mrid:       bv.Mrid,
		EntityType: StructName(bv),
	}
	_, err = db.NewInsert().Model(&entity).Exec(ctx)
	require.NoError(t, err)

	patch := []JsonPatch{{
		Op:    "replace",
		Path:  fmt.Sprintf("/%s/nominal_voltage", bv.Mrid),
		Value: []byte{0x32, 0x32},
	}}

	t.Run("success", func(t *testing.T) {
		err = ApplyPatch(ctx, db, patch)
		require.NoError(t, err)

		var bvs []models.BaseVoltage
		err = db.NewSelect().Model(&bvs).Scan(ctx)
		require.NoError(t, err)

		require.Equal(t, 2, len(bvs))

		active := OnlyActiveLatest(bvs)
		require.Equal(t, 1, len(active))
		require.Equal(t, 22.0, active[0].NominalVoltage)
	})

	t.Run("unknown mrid", func(t *testing.T) {
		patch := []JsonPatch{{
			Op:   "replace",
			Path: "/0000-0000/nominal_voltage",
		}}
		err = ApplyPatch(ctx, db, patch)
		require.ErrorContains(t, err, "in result set")
	})

	t.Run("unknown wrong path format", func(t *testing.T) {
		patch := []JsonPatch{{
			Op:   "replace",
			Path: "/0000-0000/nominal_voltage/what",
		}}
		err = ApplyPatch(ctx, db, patch)
		require.ErrorContains(t, err, "Parse patch")
	})

	t.Run("Unsupported operation", func(t *testing.T) {
		patch := []JsonPatch{{
			Op:    "some-random-op",
			Path:  fmt.Sprintf("/%s/nominal_voltage", bv.Mrid),
			Value: []byte{0x32, 0x32},
		}}
		err = ApplyPatch(ctx, db, patch)
		require.ErrorContains(t, err, "Unsupported operation")
	})
}
