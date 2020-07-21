package util_test

import (
	"os"
	"testing"

	"github.com/plumming/chilly/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestConfigDir(t *testing.T) {
	os.Setenv("HOME", "/random/home/dir")
	dir := util.ConfigDir()
	assert.Equal(t, "/random/home/dir/.chilly", dir)
}

func TestChillyConfigFile(t *testing.T) {
	os.Setenv("HOME", "/random/home/dir")
	dir := util.ChillyConfigFile()
	assert.Equal(t, "/random/home/dir/.chilly/config.yml", dir)
}

func TestChillyGhConfigFile(t *testing.T) {
	os.Setenv("HOME", "/random/home/dir")
	dir := util.GhConfigDir()
	assert.Equal(t, "/random/home/dir/.config/gh", dir)
}
