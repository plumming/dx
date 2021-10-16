package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/plumming/dx/pkg/util"
)

var (
	defaultGitHubRepos    = []string{"plumming/dx"}
	defaultHiddenLabels   = []string{"hide-this"}
	defaultMaxNumberOfPRs = 100
	defaultMaxAge         = -1
	defaultBotAccounts    = []string{"dependabot", "dependabot-preview"}
)

// Config defines repos to watch.
type Config struct {
	// Deprecated
	Repos []string `yaml:"repos"`
	//Repositories   map[string][]string `json:"repositories"`
	HiddenLabels   []string `yaml:"hiddenLabels"`
	MaxNumberOfPRs int      `yaml:"maxNumberOfPRs"`
	MaxAge         int      `yaml:"maxAgeOfPRs"`
	BotAccounts    []string `yaml:"botAccounts"`
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

		return Load(bytes.NewReader(content))
	}

	config := Config{}
	config.SetDefaults()

	return config, nil
}

func (c *Config) SetDefaults() {
	if c.Repos == nil || len(c.Repos) == 0 {
		c.Repos = defaultGitHubRepos
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

	if c.BotAccounts == nil || len(c.BotAccounts) == 0 {
		c.BotAccounts = defaultBotAccounts
	}
}

func Load(reader io.Reader) (Config, error) {
	config := Config{}
	err := yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	config.SetDefaults()
	return config, nil
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
