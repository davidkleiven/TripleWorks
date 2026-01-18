package pkg

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	var baseVoltage models.BaseVoltage
	baseVoltage.Mrid = uuid.UUID{}
	baseVoltage.Name = "bv 1"
	baseVoltage.ShortName = "b"
	baseVoltage.Description = "Base voltage"
	baseVoltage.NominalVoltage = 22.0

	var buf bytes.Buffer
	Export(&buf, func(yield func(item models.MridGetter) bool) { yield(baseVoltage) })

	triples := []string{
		"<urn:uuid:00000000-0000-0000-0000-000000000000> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://iec.ch/TC57/2013/CIM-schema-cim16#BaseVoltage> .",
		"<urn:uuid:00000000-0000-0000-0000-000000000000> <http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.mRID> \"00000000-0000-0000-0000-000000000000\" .",
		"<urn:uuid:00000000-0000-0000-0000-000000000000> <http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.name> \"bv 1\" .",
		"<urn:uuid:00000000-0000-0000-0000-000000000000> <http://entsoe.eu/CIM/SchemaExtension/3/1#IdentifiedObject.shortName> \"b\" .",
		"<urn:uuid:00000000-0000-0000-0000-000000000000> <http://iec.ch/TC57/2013/CIM-schema-cim16#IdentifiedObject.description> \"Base voltage\" .",
		"<urn:uuid:00000000-0000-0000-0000-000000000000> <http://iec.ch/TC57/2013/CIM-schema-cim16#BaseVoltage.nominalVoltage> \"22\"^^<http://www.w3.org/2001/XMLSchema#float> .",
	}

	content := buf.String()
	for i, token := range triples {
		require.Contains(t, content, token, fmt.Sprintf("Does not contain token: %d", i))
	}
}

func TestListAll(t *testing.T) {
	var bv models.BaseVoltage
	var sync models.SynchronousMachine

	db := NewTestConfig(WithDbName(t.Name())).DatabaseConnection()
	ctx := context.Background()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&bv).Exec(ctx)
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&sync).Exec(ctx)
	require.NoError(t, err)

	t.Run("receives all", func(t *testing.T) {
		allItems, err := LatestOfAllItems(ctx, db, 0)
		require.NoError(t, err)
		require.Equal(t, 2, len(allItems))
	})

	t.Run("struct name in error", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()
		_, err := LatestOfAllItems(cancelledCtx, db, 0)
		require.Error(t, err)
		require.ErrorContains(t, err, "Could not get")
	})
}

func TestExportAll(t *testing.T) {
	var bv models.BaseVoltage
	var sync models.SynchronousMachine
	iter := func(yield func(v models.MridGetter) bool) {
		if !yield(bv) {
			return
		}
		if !yield(sync) {
			return
		}
	}

	var buf bytes.Buffer
	Export(&buf, iter)
	content := buf.String()
	require.Contains(t, content, "BaseVoltage")
	require.Contains(t, content, "SynchronousMachine")
}

func TestTypeSpecifier(t *testing.T) {
	intSpecifier := typeSpecifier(2)
	require.Contains(t, intSpecifier, "XMLSchema#int")
}

func TestExportItemPointer(t *testing.T) {
	var bv models.BaseVoltage
	var buf bytes.Buffer
	ExportItem(&buf, &bv)
	content := buf.String()
	require.Contains(t, content, "BaseVoltage")
}
