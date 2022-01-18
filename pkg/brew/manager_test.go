package brew_test

import (
	"testing"

	"github.com/plumming/dx/pkg/brew"
	"github.com/stretchr/testify/assert"
)

func TestManager_IsBrewInstalled(t *testing.T) {
	_, err := brew.IsInstalledViaBrew()
	assert.NoError(t, err)
}
