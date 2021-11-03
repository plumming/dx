package auth

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
)

const (
	APIGithubCom = "api.github.com"
)

// Config interface describes interacting with some persistent configuration for gh.
type Config interface {
	GetToken(hostname string) string
	GetUser(hostname string) string
	HasHosts() bool
	GetHosts() []string
}

type fileConfig struct {
	Hosts map[string]*HostConfig `json:"hosts"`
}

// HostConfig a struct containing the host level information.
type HostConfig struct {
	User  string `yaml:"user"`
	Token string `yaml:"oauth_token"`
}

func (c *fileConfig) GetToken(hostname string) string {
	log.Logger().Debugf("Getting token for host %s", hostname)
	if c.Hosts[hostname] != nil {
		return c.Hosts[hostname].Token
	}
	return ""
}

func (c *fileConfig) GetUser(hostname string) string {
	if c.Hosts[hostname] != nil {
		return c.Hosts[hostname].User
	}
	return ""
}

func (c *fileConfig) HasHosts() bool {
	return len(c.Hosts) > 0
}

func (c *fileConfig) GetHosts() []string {
	var hosts []string
	for k := range c.Hosts {
		hosts = append(hosts, k)
	}
	return hosts
}
