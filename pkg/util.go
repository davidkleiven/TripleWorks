package pkg

import (
	"cmp"
	"fmt"
	"iter"
	"math"
	"reflect"
	"slices"
	"strings"

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

	fTypes := FormTypes()
	for _, v := range fTypes {
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

func AssertNotNil(v any) {
	if v == nil {
		panic("Value should not be nil")
	}
}

func AssertDifferent[K comparable](v1, v2 K) {
	if v1 == v2 {
		panic(fmt.Sprintf("%v is not different from %v", v1, v2))
	}
}

func AssertGreater[K cmp.Ordered](v1, v2 K) {
	if v1 < v2 {
		panic(fmt.Sprintf("%v is smaller than %v", v1, v2))
	}
}

func Chain[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(v T) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func CosineSimilarity(word1, word2 string) float64 {
	if word1 == "" && word2 == "" {
		return 1.0
	}
	set1 := make(map[string]struct{})
	set2 := make(map[string]struct{})
	union := make(map[string]struct{})
	for token := range Ngrams(word1, 3) {
		set1[token] = struct{}{}
	}

	for token := range Ngrams(word2, 3) {
		set2[token] = struct{}{}
	}

	for k := range set1 {
		union[k] = struct{}{}
	}
	for k := range set2 {
		union[k] = struct{}{}
	}

	var (
		dot   float64
		norm1 float64
		norm2 float64
	)
	for k := range union {
		elem1 := 0.0
		if _, ok := set1[k]; ok {
			elem1 = 1.0
		}
		elem2 := 0.0
		if _, ok := set2[k]; ok {
			elem2 = 1.0
		}
		v1 := elem1
		v2 := elem2
		dot += v1 * v2
		norm1 += v1 * v1
		norm2 += v2 * v2
	}
	if norm1 == 0 && norm2 == 0 {
		return 1.0
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}
	return dot / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

func Ngrams(word string, n int) iter.Seq[string] {
	return func(yield func(token string) bool) {
		for i := range len(word) {
			if i+n >= len(word) {
				return
			}
			if !yield(word[i : i+n]) {
				return
			}
		}
	}
}

func Normalizename(word string) string {
	return strings.ReplaceAll(strings.ToLower(word), "-", "")
}

func Tokenize(word string) []string {
	return strings.Fields(word)
}

func ExactTokenSimilarity(sourceTokens []string, tokens []string) float64 {
	tokenPool := make(map[string]struct{})
	for _, token := range tokens {
		tokenPool[token] = struct{}{}
	}

	var (
		sourceWeight float64
		matchWeight  float64
	)
	for _, token := range sourceTokens {
		w := math.Log(1 + float64(len(token)))
		sourceWeight += w
		if _, ok := tokenPool[token]; ok {
			matchWeight += w
		}
	}
	return matchWeight / sourceWeight
}

func IndexBy[T any, K comparable](items []T, keyFn func(T) K) map[K]T {
	m := make(map[K]T, len(items))
	for _, item := range items {
		m[keyFn(item)] = item
	}
	return m
}

func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T {
	m := make(map[K][]T, len(items))
	for _, item := range items {
		k := keyFn(item)
		m[k] = append(m[k], item)
	}
	return m
}

func IndirectDescendingSort[T cmp.Ordered](values []T) []int {
	indices := make([]int, len(values))
	for i := range indices {
		indices[i] = i
	}

	slices.SortFunc(indices, func(i, j int) int {
		return -cmp.Compare(values[i], values[j])
	})
	return indices
}
