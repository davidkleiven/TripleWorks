package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func uuidFromString(item string) uuid.UUID {
	space := uuid.UUID{}
	return uuid.NewMD5(space, []byte(item))
}

func TestBusBreakerConnection(t *testing.T) {
	sqldb, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	ctx := context.Background()
	_, err = migrations.RunUp(ctx, db)
	require.NoError(t, err)

	for commitId := range 3 {
		var bv models.BaseVoltage
		bv.Mrid = uuidFromString("bv")
		bv.CommitId = commitId
		_, err := db.NewInsert().Model(&bv).Exec(ctx)
		require.NoError(t, err)

		substations := make([]models.Substation, 2)
		for i := range substations {
			substations[i].Mrid = uuidFromString(fmt.Sprintf("sub%d", i))
			substations[i].CommitId = commitId
		}

		_, err = db.NewInsert().Model(&substations).Exec(ctx)
		require.NoError(t, err)

		vls := make([]models.VoltageLevel, 2)
		for i := range vls {
			vls[i].Mrid = uuidFromString(fmt.Sprintf("vl%d", i))
			vls[i].CommitId = commitId
			vls[i].SubstationMrid = substations[i].Mrid
			vls[i].BaseVoltageMrid = bv.Mrid
		}
		_, err = db.NewInsert().Model(&vls).Exec(ctx)
		require.NoError(t, err)

		cns := make([]models.ConnectivityNode, 2)
		for i := range cns {
			cns[i].Mrid = uuidFromString(fmt.Sprintf("cn%d", i))
			cns[i].ConnectivityNodeContainerMrid = vls[i].Mrid
			cns[i].CommitId = commitId
		}
		_, err = db.NewInsert().Model(&cns).Exec(ctx)
		require.NoError(t, err)

		var line models.ACLineSegment
		line.Mrid = uuidFromString("acline")
		line.BaseVoltageMrid = bv.Mrid
		line.CommitId = commitId
		_, err = db.NewInsert().Model(&line).Exec(ctx)
		require.NoError(t, err)

		terminals := make([]models.Terminal, 2)
		for i := range terminals {
			terminals[i].Mrid = uuidFromString(fmt.Sprintf("t%d", i))
			terminals[i].ConductingEquipmentMrid = line.Mrid
			terminals[i].ConnectivityNodeMrid = cns[i].Mrid
			terminals[i].CommitId = commitId
		}

		_, err = db.NewInsert().Model(&terminals).Exec(ctx)
		require.NoError(t, err)
	}

	store := BunBusBreakerRepo{Db: db}
	result, err := store.Fetch(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(result))

	acLines := make(map[uuid.UUID]struct{})
	for _, res := range result {
		acLines[res.Mrid] = struct{}{}
	}
	require.Equal(t, len(acLines), 1)
}

func TestCachedBusBreakerRepo(t *testing.T) {
	var cb CachedBusbReakerrepo
	res, err := cb.Fetch(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
}
