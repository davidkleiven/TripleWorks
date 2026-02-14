package integrity

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"testing/quick"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFixInconsistentTerminals(t *testing.T) {
	terminals := make([]models.Terminal, 9)
	for i := range 3 {
		terminals[i].Mrid = uuid.New()
		terminals[i].ConductingEquipmentMrid = uuid.New()
		terminals[i].CommitId = i
		terminals[i].SequenceNumber = 1
	}

	// Should be sequence number 2
	for i := range 3 {
		terminals[i+3].Mrid = uuid.New()
		terminals[i+3].ConductingEquipmentMrid = terminals[i].ConductingEquipmentMrid
		terminals[i+3].CommitId = i + 3
		terminals[i+3].SequenceNumber = 1
		if i == 0 {
			// This should make the first terminal valid
			terminals[i+3].SequenceNumber = 2
		}
	}

	// Copies of the first three with later commit
	for i := range 3 {
		terminals[i+6].Mrid = terminals[i].Mrid
		terminals[i+6].ConductingEquipmentMrid = terminals[i].ConductingEquipmentMrid
		terminals[i+6].CommitId = i + 6
		terminals[i+6].SequenceNumber = terminals[i].SequenceNumber
	}

	check := UniqueSequenceNumberPerConductingEquipment{terminals: terminals}
	result := check.Check()
	invalid := result.Fix()
	num := 0
	for range invalid {
		num++
	}
	require.Equal(t, 2, num)

	t.Run("report", func(t *testing.T) {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		err := result.Report(enc)
		require.NoError(t, err)
	})
}

func TestDoubleFixIsAlwaysEmpty(t *testing.T) {
	prop := func() bool {
		var terminals []models.Terminal
		err := faker.FakeData(&terminals, options.WithRandomMapAndSliceMinSize(1))
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(terminals), 1)

		for i := range terminals {
			terminals[i].CommitId = i
		}

		check := UniqueSequenceNumberPerConductingEquipment{terminals: terminals}
		result := check.Check()

		invalidTermials := []models.Terminal{}
		for term := range result.Fix() {
			asTerm, ok := term.(models.Terminal)
			if !ok {
				return false
			}
			asTerm.CommitId = len(terminals) + 1
			invalidTermials = append(invalidTermials, asTerm)
		}

		// Each terminal should only occure once among invalid
		grouped := pkg.GroupBy(invalidTermials, func(trm models.Terminal) uuid.UUID { return trm.Mrid })
		for mrid, group := range grouped {
			if len(group) != 1 {
				t.Logf("%s occure %d times among terminals", mrid, len(group))
				return false
			}
		}

		// Append the new terminals to pretend they are commited and rerun fix
		check2 := UniqueSequenceNumberPerConductingEquipment{terminals: append(terminals, invalidTermials...)}
		result = check2.Check()
		stillInvalid := []models.Terminal{}
		for term := range result.Fix() {
			asTerm, ok := term.(models.Terminal)
			require.True(t, ok)
			stillInvalid = append(stillInvalid, asTerm)
		}
		t.Logf("Still %d invalid terminals", len(stillInvalid))
		return len(stillInvalid) == 0
	}
	err := quick.Check(prop, nil)
	t.Log(err)
	require.NoError(t, err)
}

func TestFetchTerminals(t *testing.T) {
	db := pkg.NewTestConfig(pkg.WithDbName(t.Name())).DatabaseConnection()
	ctx := context.Background()
	_, err := migrations.RunUp(ctx, db)
	require.NoError(t, err)

	var terminals []models.Terminal
	err = faker.FakeData(&terminals)
	require.NoError(t, err)

	_, err = db.NewInsert().Model(&terminals).Exec(ctx)
	require.NoError(t, err)

	check := UniqueSequenceNumberPerConductingEquipment{}
	check.Fetch(ctx, db)
	require.Equal(t, len(terminals), len(check.terminals))
}
