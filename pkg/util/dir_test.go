package util_test

import (
	"testing"

	"github.com/plumming/dx/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestDxBinaryLocation(t *testing.T) {
	_, err := util.DxBinaryLocation()
	assert.NoError(t, err)
}
