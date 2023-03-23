package domain

import (
	"fmt"
	"strings"

	"github.com/plumming/dx/pkg/securityconfig"

	"github.com/plumming/dx/pkg/config"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/cmd"
)

var (
	getSecurityConfigQuery = `{
	search(type:REPOSITORY, query: "%s", first: 100) {
      nodes {
      ... on Repository {
        nameWithOwner
        url
        hasVulnerabilityAlertsEnabled
      }
    }
  }
}`
)

// GetSecurityConfig defines get vulnerability alerts request response.
type GetSecurityConfig struct {
	cmd.CommonOptions
	Config []securityconfig.Repository
}

// SecurityConfigData.
type SecurityConfigData struct {
	Search SecurityConfigSearch `json:"search"`
}

// SecurityConfigSearch.
type SecurityConfigSearch struct {
	SecurityConfig []securityconfig.Repository `json:"nodes"`
}

// NewGetSecurityConfig.
func NewGetSecurityConfig() *GetSecurityConfig {
	g := &GetSecurityConfig{}
	return g
}

// Validate input.
func (g *GetSecurityConfig) Validate() error {
	return nil
}

// Run the cmd.
func (g *GetSecurityConfig) Run() error {
	cfg, err := g.DxConfig()
	if err != nil {
		return err
	}

	var config []securityconfig.Repository

	for _, host := range cfg.GetConfiguredServers() {
		c, err := g.GetSecurityConfigForHost(host, cfg, getSecurityConfigQuery)
		if err != nil {
			return err
		}
		config = append(config, c...)
	}

	g.Config = config

	return nil
}

func (g *GetSecurityConfig) GetSecurityConfigForHost(host string, cfg config.Config, query string) ([]securityconfig.Repository, error) {
	client, err := g.GithubClient()
	if err != nil {
		return nil, err
	}

	queryString := strings.Join(cfg.GetReposToQuery(host), " ")
	queryToRun := fmt.Sprintf(query, queryString)
	log.Logger().Debugf("running query\n%s", queryToRun)

	data := SecurityConfigData{}
	err = client.GraphQL(host, queryToRun, nil, &data)
	if err != nil {
		return nil, err
	}

	return data.Search.SecurityConfig, nil
}
