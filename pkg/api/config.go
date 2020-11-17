package api

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/plumming/dx/pkg/util"
)

func ConfigFile() string {
	return path.Join(util.GhConfigDir(), "config.yml")
}

func HostsFile() string {
	return path.Join(util.GhConfigDir(), "hosts.yml")
}

func ParseDefaultConfig(cf, hf string) (Config, error) {
	// since tokens are now stored as default within
	// ~/.config/gh/hosts.yml lets try and load from
	// there initially
	config, err := parseHostsFile(hf)
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		return nil, err
	}
	if config.HasHosts() {
		return config, nil
	}

	return parseConfigFile(cf)
}

var readConfigFile = func(fn string) ([]byte, error) {
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

func parseConfigFile(fn string) (Config, error) {
	data, err := readConfigFile(fn)
	if err != nil {
		return nil, err
	}

	var root fileConfig
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	return &root, nil
}

func parseHostsFile(hf string) (Config, error) {
	c := &fileConfig{}
	var h map[string]*HostConfig
	data, err := readConfigFile(hf)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(data, &h)
	if err != nil {
		return c, err
	}
	// check there are values in the map
	if len(h) > 0 {
		c.Hosts = h
	}

	return c, nil
}
