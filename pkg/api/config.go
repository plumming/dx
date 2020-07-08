package api

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/plumming/chilly/pkg/util"
)

func ConfigFile() string {
	return path.Join(util.GhConfigDir(), "config.yml")
}

func ParseDefaultConfig() (Config, error) {
	return ParseConfig(ConfigFile())
}

var ReadConfigFile = func(fn string) ([]byte, error) {
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

func parseConfigFile(fn string) ([]byte, Config, error) {
	data, err := ReadConfigFile(fn)
	if err != nil {
		return nil, nil, err
	}

	var root fileConfig
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return data, nil, err
	}

	return data, &root, nil
}

func ParseConfig(fn string) (Config, error) {
	_, root, err := parseConfigFile(fn)
	if err != nil {
		return nil, err
	}

	return root, nil
}
