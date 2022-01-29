package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimTemplate(t *testing.T) {
	assert.Equal(t, TrimTemplate("my-instance-template-1.2.3", "my"), "1.2.3")
}

func TestSafeIfAboveZero(t *testing.T) {
	assert.Equal(t, SafeIfAboveZero(0), "")
	assert.Equal(t, SafeIfAboveZero(1), "1")
}
