package pkg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	var buf bytes.Buffer
	Index(&buf)

	assert.Contains(t, buf.String(), "TripleWorks")
}

func TestPatachForm(t *testing.T) {
	var buf bytes.Buffer
	PatchForm(&buf)
	require.Contains(t, buf.String(), "Patch")
}
