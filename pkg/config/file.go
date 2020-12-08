package config

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/util"
)

var (
	defaultRepos          = []string{"plumming/dx"}
	defaultHiddenLabels   = []string{"hide-this"}
	defaultMaxNumberOfPRs = 100
	defaultMaxAge         = -1
)

// Config defines repos to watch.
type Config struct {
	Repos          []string `json:"repos"`
	HiddenLabels   []string `json:"hiddenLabels"`
	MaxNumberOfPRs int      `json:"maxNumberOfPRs"`
	MaxAge         int      `json:"maxAgeOfPRs"`
}

func (c *Config) ReposToQuery() []string {
	var repoList []string
	for _, r := range c.Repos {
		repoList = append(repoList, fmt.Sprintf("repo:%s", r))
	}
	return repoList
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

		config.SetDefaults()

		return config, nil
	}

	config := Config{}
	config.SetDefaults()

	return config, nil
}

func (c *Config) SetDefaults() {
	if c.Repos == nil || len(c.Repos) == 0 {
		c.Repos = defaultRepos
	}

	if c.HiddenLabels == nil || len(c.HiddenLabels) == 0 {
		c.HiddenLabels = defaultHiddenLabels
	}

	if c.MaxNumberOfPRs == 0 {
		c.MaxNumberOfPRs = defaultMaxNumberOfPRs
	}

	if c.MaxAge == 0 {
		c.MaxAge = defaultMaxAge
	}
}

func LoadFromDefaultLocation() (Config, error) {
	return LoadFromFile(util.DxConfigFile())
}

func (c *Config) SaveToDefaultLocation() error {
	return c.SaveToFile(util.DxConfigFile())
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
