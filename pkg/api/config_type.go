package api

const defaultHostname = "github.com"

// Config interface describes interacting with some persistent configuration for gh.
type Config interface {
	GetToken(hostname string) string
	GetUser(hostname string) string
	HasHosts() bool
}

type fileConfig struct {
	Hosts map[string]*HostConfig `yaml:"hosts"`
}

// NotFoundError.
type NotFoundError struct {
	error
}

// HostConfig.
type HostConfig struct {
	User  string `yaml:"user"`
	Token string `yaml:"oauth_token"`
}

func (c *fileConfig) GetToken(hostname string) string {
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
