package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	configDir    = ".config"
	ghConfigDir  = "gh"
	dxConfigFile = "config.yml"
	dxConfigDir  = ".dx"
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

func ReadFile(fn string) ([]byte, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
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
	configDir := filepath.Join(home, dxConfigDir)
	return configDir
}

func CreateConfigPath(configPath, configFile string) *string {
	home := HomeDir()
	path := filepath.Join(home, configPath, configFile)
	return &path
}

func GhConfigDir() string {
	home := HomeDir()
	configDir := filepath.Join(home, configDir, ghConfigDir)
	return configDir
}

func DxConfigFile() string {
	return filepath.Join(ConfigDir(), dxConfigFile)
}

// credit https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}

	return nil
}

// credit https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyDir(src string, dst string, force bool) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		if force {
			os.RemoveAll(dst)
		} else {
			return os.ErrExist
		}
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath, force)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
