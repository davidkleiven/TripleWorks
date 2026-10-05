package pkg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResult(t *testing.T) {
	withErr := Result[bool]{Err: errors.New("error")}
	mapResult := withErr.Map(func(v bool) bool { return !v })
	require.Error(t, mapResult.Err)

	applyResult := withErr.Apply(func(v bool) (bool, error) { return !v, nil })
	require.Error(t, applyResult.Err)
}
