package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/plumming/dx/pkg/util"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Config

var (
	defaultGitHubRepos    = []string{"plumming/dx"}
	defaultHiddenLabels   = []string{"hide-this"}
	defaultMaxNumberOfPRs = 100
	defaultMaxAge         = -1
	defaultBotAccounts    = []string{"dependabot", "dependabot-preview"}
)

// Config defines repos to watch.
type Config interface {
	GetReposToQuery(server string) []string
	GetConfiguredServers() []string
	GetBotAccounts() []string
	GetHiddenLabels() []string
	GetMaxAgeOfPRs() int
	GetMaxNumberOfPRs() int
	SaveToDefaultLocation() error
	SaveToFile(path string) error
}

// fileBasedConfig defines repos to watch.
type fileBasedConfig struct {
	// Deprecated
	Repos          []string            `yaml:"repos"`
	Repositories   map[string][]string `yaml:"repositories"`
	HiddenLabels   []string            `yaml:"hiddenLabels"`
	MaxNumberOfPRs int                 `yaml:"maxNumberOfPRs"`
	MaxAge         int                 `yaml:"maxAgeOfPRs"`
	BotAccounts    []string            `yaml:"botAccounts"`
}

func (c *fileBasedConfig) GetReposToQuery(server string) []string {
	var repoList []string
	for _, r := range c.Repositories[server] {
		repoList = append(repoList, fmt.Sprintf("repo:%s", r))
	}
	return repoList
}

func (c *fileBasedConfig) GetConfiguredServers() []string {
	var servers []string
	for k := range c.Repositories {
		servers = append(servers, k)
	}
	return servers
}

func LoadFromFile(path string) (Config, error) {
	if exists, err := util.FileExists(path); err == nil && exists {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return &fileBasedConfig{}, err
		}

		return Load(bytes.NewReader(content))
	}

	config := &fileBasedConfig{}
	config.setDefaults()

	return config, nil
}

func (c *fileBasedConfig) setDefaults() {
	// migrate from the old structure to the new structure
	if c.Repositories == nil || len(c.Repositories) == 0 {
		c.Repositories = make(map[string][]string)
		if c.Repos == nil || len(c.Repos) == 0 {
			c.Repositories["github.com"] = defaultGitHubRepos
		} else {
			c.Repositories["github.com"] = c.Repos
			c.Repos = nil
		}
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
	config := &fileBasedConfig{}
	err := yaml.NewDecoder(reader).Decode(config)
	if err != nil {
		return &fileBasedConfig{}, err
	}

	config.setDefaults()
	return config, nil
}

func LoadFromDefaultLocation() (Config, error) {
	return LoadFromFile(util.DxConfigFile())
}

func (c *fileBasedConfig) GetBotAccounts() []string {
	return c.BotAccounts
}

func (c *fileBasedConfig) GetHiddenLabels() []string {
	return c.HiddenLabels
}

func (c *fileBasedConfig) GetMaxAgeOfPRs() int {
	return c.MaxAge
}

func (c *fileBasedConfig) GetMaxNumberOfPRs() int {
	return c.MaxNumberOfPRs
}

func (c *fileBasedConfig) SaveToDefaultLocation() error {
	return c.SaveToFile(util.DxConfigFile())
}

func (c *fileBasedConfig) SaveToFile(path string) error {
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
