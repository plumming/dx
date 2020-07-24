package util_test

import (
	"testing"

	"github.com/plumming/chilly/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestChillyBinaryLocation(t *testing.T) {
	_, err := util.ChillyBinaryLocation()
	assert.NoError(t, err)
}
