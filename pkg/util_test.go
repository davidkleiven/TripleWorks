package pkg

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustPanicOnErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Wanted panic")
		}
	}()

	Must(0, errors.New("Something went wrong"))
}

func TestStopEarly(t *testing.T) {
	// Basic test that it is possible to break early
	mapping := map[int]int{1: 2, 3: 5, 0: 7}
	result := make([]int, 0, len(mapping))
	num := 0
	for k := range Keys(mapping) {
		result = append(result, k)
		num += 1
		if num == 2 {
			break
		}
	}
	assert.Equal(t, len(result), 2)
}

func TestStructName(t *testing.T) {
	model := models.ACLineSegment{}
	name := StructName(model)
	require.Equal(t, name, "ACLineSegment")
	require.Equal(t, name, StructName(&model))
}

func TestMustBeValid(t *testing.T) {
	require.Panics(t, func() { MustBeValid(reflect.Value{}) })
}

func TestReturnOnFirstError(t *testing.T) {
	num, err := ReturnOnFirstError(
		func() error { return nil },
		func() error { return fmt.Errorf("What??") },
		func() error { return nil },
	)
	require.Error(t, err)
	require.Equal(t, 1, num)

	num, err = ReturnOnFirstError()
	require.NoError(t, err)
	require.Equal(t, 0, num)
}

func TestPanicOnErr(t *testing.T) {
	require.Panics(t, func() { PanicOnErr(fmt.Errorf("What??")) })
}

func TestSetCommitId(t *testing.T) {
	t.Run("invalid struct", func(t *testing.T) {
		err := SetCommitId(struct{}{}, 0)
		require.Error(t, err)
	})

	t.Run("valid struct", func(t *testing.T) {
		bv := models.BaseVoltage{}
		err := SetCommitId(&bv, 4)
		require.NoError(t, err)
		require.Equal(t, 4, bv.CommitId)
	})
}

func TestUnsetFields(t *testing.T) {
	model := struct {
		Name        string
		LastName    string `json:"last_name"`
		BunRelation string `bun:"rel:belongs-to"`
	}{}

	data := map[string]any{
		"Name": true,
	}

	unset := UnsetFields(data, model)
	require.Equal(t, []string{"last_name"}, unset)

}

func TestSubtypes(t *testing.T) {
	subtypes := Subtypes(&models.RotatingMachine{})
	require.Equal(t, len(subtypes), 2)

	names := make(map[string]struct{})
	for _, t := range subtypes {
		names[StructName(t)] = struct{}{}
	}
	want := map[string]struct{}{
		"AsynchronousMachine": {},
		"SynchronousMachine":  {},
	}
	require.Equal(t, want, names)

	t.Run("Substation is EquipmentContainer", func(t *testing.T) {
		subtypes := Subtypes(&models.EquipmentContainer{})
		isSubtype := false
		for _, subtype := range subtypes {
			if StructName(subtype) == "Substation" {
				isSubtype = true
				break
			}
		}
		require.True(t, isSubtype, fmt.Sprintf("%v\n", subtypes))
	})
}

func TestMustGet(t *testing.T) {
	mapping := map[string]bool{
		"field": true,
	}
	require.True(t, MustGet(mapping, "field"))
	require.Panics(t, func() { MustGet(mapping, "field2") })
}

func TestAssertNotNil(t *testing.T) {
	require.NotPanics(t, func() { AssertNotNil("what") })
	require.Panics(t, func() { AssertNotNil(nil) })
}

func TestRequireStruct(t *testing.T) {
	require.NotPanics(t, func() { RequireStruct(reflect.TypeOf(struct{}{})) })
	require.Panics(t, func() { RequireStruct(reflect.TypeOf(2)) })
}
