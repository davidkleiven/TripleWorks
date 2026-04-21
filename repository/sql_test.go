package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPanicOnUnknownQuery(t *testing.T) {
	require.Panics(t, func() { MustGetQuery("unknown_query.sql") })
}
