package api

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/ghodss/yaml"
	"github.com/plumming/dx/pkg/util"
)

func ConfigFile() string {
	return path.Join(util.GhConfigDir(), "config.yml")
}

func HostsFile() string {
	return path.Join(util.GhConfigDir(), "hosts.yml")
}

func ParseDefaultConfig() (Config, error) {
	return ParseConfig(ConfigFile(), HostsFile())
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

func parseConfigFile(fn, hf string) ([]byte, Config, error) {
	data, err := readConfigFile(fn)
	if err != nil {
		return nil, nil, err
	}

	var root fileConfig
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return data, nil, err
	}

	// First step to support new hosts configuration in gh config
	if root.Hosts == nil {
		hosts, err := parseHostsFile(hf)
		if err != nil {
			return data, nil, err
		}
		root.Hosts = hosts
	}

	return data, &root, nil
}

func parseHostsFile(hf string) (map[string]*HostConfig, error) {
	var h map[string]*HostConfig
	hostsFile, err := readConfigFile(hf)
	if err != nil {
		return h, err
	}
	err = yaml.Unmarshal(hostsFile, &h)
	if err != nil {
		return h, errors.Wrap(err, "while unmarshaling hosts file")
	}
	return h, nil
}

func ParseConfig(fn, hf string) (Config, error) {
	_, root, err := parseConfigFile(fn, hf)
	if err != nil {
		return nil, err
	}

	return root, nil
}
