package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTop2Selection(t *testing.T) {
	selector := TopSelector{Num: 2}
	lines := []string{"Trondheim - Brottem", "Trondheim-Namsos"}
	substations := []string{"Trondheim", "Brottem", "Trond", "Nams", "Namsos", "Brott"}

	result := selector.Select(lines, substations, NameSimilarity)
	want := [][]int{{0, 1}, {0, 4}}
	require.Equal(t, want, result)
}
