package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/chilly/pkg/util"
)

var (
	defaultRepos = []string{"plumming/chilly"}
)

// Config defines repos to watch.
type Config struct {
	Repos        []string `json:"repos"`
	HiddenLabels []string `json:"hiddenLabels"`
}

func LoadFromFile(path string) (Config, error) {
	if exists, err := util.FileExists(path); err == nil && exists {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return Config{}, err
		}

		config := Config{}
		err = yaml.Unmarshal(content, &config)
		if err != nil {
			log.Logger().Infof("no repos configured in %s", path)
			return Config{}, err
		}

		if config.Repos == nil || len(config.Repos) == 0 {
			config.Repos = defaultRepos
		}

		return config, nil
	}

	config := Config{}
	config.Repos = defaultRepos

	return config, nil
}

func LoadFromDefaultLocation() (Config, error) {
	return LoadFromFile(util.ChillyConfigFile())
}

func (c *Config) SaveToDefaultLocation() error {
	return c.SaveToFile(util.ChillyConfigFile())
}

func (c *Config) SaveToFile(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		return err
	}
	return nil
}
