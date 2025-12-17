package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config := NewDefaultConfig()
	assert.Equal(t, config.Port, 36000)
}
