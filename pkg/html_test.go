package pkg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	var buf bytes.Buffer
	Index(&buf)

	assert.Contains(t, buf.String(), "TripleWorks")
}
