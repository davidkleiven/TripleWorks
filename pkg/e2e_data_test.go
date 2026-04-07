package pkg

import (
	"testing"

	"com.github/davidkleiven/tripleworks/repository"
	"github.com/stretchr/testify/require"
)

func TestE2eCanBeInserted(t *testing.T) {
	data := MakeE2eData()
	inserter := repository.InMemInserter{}
	require.NotPanics(t, func() { InsertE2eData(data, &inserter) })
}
