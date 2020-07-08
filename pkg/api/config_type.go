package api

const defaultHostname = "github.com"

// This interface describes interacting with some persistent configuration for gh.
type Config interface {
	GetToken(hostname string) string
	GetUser(hostname string) string
}

type fileConfig struct {
	Hosts map[string]*HostConfig `json:"hosts"`
}

type NotFoundError struct {
	error
}

type HostConfig struct {
	User  string `json:"user"`
	Token string `json:"oauth_token"`
}

func (c *fileConfig) GetToken(hostname string) string {
	return c.Hosts[hostname].Token
}

func (c *fileConfig) GetUser(hostname string) string {
	return c.Hosts[hostname].User
}
