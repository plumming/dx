package util

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	configDir        = ".config"
	ghConfigDir      = "gh"
	chillyConfigFile = "config.yml"
	chillyConfigDir  = ".chilly"
)

// FileExists checks if path exists and is a file.
func FileExists(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return !fileInfo.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.Wrapf(err, "failed to check if file exists %s", path)
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	h := os.Getenv("USERPROFILE") // windows
	if h == "" {
		h = "."
	}
	return h
}

func ConfigDir() string {
	home := HomeDir()
	configDir := filepath.Join(home, chillyConfigDir)
	return configDir
}

func GhConfigDir() string {
	home := HomeDir()
	configDir := filepath.Join(home, configDir, ghConfigDir)
	return configDir
}

func ChillyConfigFile() string {
	return filepath.Join(ConfigDir(), chillyConfigFile)
}
