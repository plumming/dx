package util_test

import (
	"os"
	"testing"

	"github.com/plumming/dx/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestConfigDir(t *testing.T) {
	os.Setenv("HOME", "/random/home/dir")
	dir := util.ConfigDir()
	assert.Equal(t, "/random/home/dir/.dx", dir)
}

func TestDxConfigFile(t *testing.T) {
	os.Setenv("HOME", "/random/home/dir")
	dir := util.DxConfigFile()
	assert.Equal(t, "/random/home/dir/.dx/config.yml", dir)
}

func TestDxGhConfigFile(t *testing.T) {
	os.Setenv("HOME", "/random/home/dir")
	dir := util.GhConfigDir()
	assert.Equal(t, "/random/home/dir/.config/gh", dir)
}
