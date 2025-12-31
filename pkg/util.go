package pkg

import (
	"iter"
	"reflect"
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

func StructName(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
