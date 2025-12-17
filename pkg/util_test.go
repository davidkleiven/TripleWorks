package pkg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
