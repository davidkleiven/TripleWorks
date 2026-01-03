package pkg

import (
	"fmt"
	"iter"
	"reflect"

	"com.github/davidkleiven/tripleworks/models"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Keys returns a sequence of keys from a map.
func Keys[K comparable, V any](m map[K]V) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return // stop early if consumer wants
			}
		}
	}
}

func MustGet[K comparable, V any](m map[K]V, key K) V {
	v, ok := m[key]
	if !ok {
		panic(fmt.Sprintf("key %v does not exist in map", key))
	}
	return v
}

func StructName(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

func MustBeValid(f reflect.Value) {
	if !f.IsValid() {
		panic("extracted and invalid field")
	}
}

func ReturnOnFirstError(fns ...func() error) (int, error) {
	for i, fn := range fns {
		err := fn()
		if err != nil {
			return i, err
		}
	}
	return 0, nil
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func SetCommitId(model any, commitId int) error {
	commitSetter, ok := model.(models.CommitIdSetter)
	if !ok {
		return fmt.Errorf("Could not convert %v into 'CommitIdSetter'", model)
	}
	commitSetter.SetCommitId(commitId)
	return nil
}

func UnsetFields(data map[string]any, target any) []string {
	fields := FlattenStruct(target)
	unset := []string{}
	for k, formField := range fields {
		if formField.IsBunRelation {
			continue
		}
		tag := formField.JsonTag
		if tag == "" {
			tag = k
		}
		_, ok := data[tag]
		if !ok {
			unset = append(unset, tag)
		}

	}
	return unset
}

func baseType(v any) reflect.Type {
	t := reflect.TypeOf(v)
	if t == nil {
		return nil
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func Subtypes(model any) []any {
	targetType := baseType(model)
	var result []any

	for _, v := range FormTypes {
		vType := baseType(v)
		for i := range vType.NumField() {
			f := vType.Field(i)
			if !f.Anonymous {
				continue
			}

			if f.Type == targetType {
				newTypes := Subtypes(v)
				result = append(result, v)
				result = append(result, newTypes...)
			}
		}
	}
	return result
}

func RequireStruct(v reflect.Type) {
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("'%v' is not a struct", v))
	}
}
