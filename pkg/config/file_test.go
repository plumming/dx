package config_test

import (
	"strings"
	"testing"

	"github.com/plumming/dx/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type test struct {
		name         string
		config       string
		expectations func(t *testing.T, c config.Config)
	}

	tests := []test{
		{
			name: "simple config",
			config: `---
botAccounts:
- bot1
- bot2
hiddenLabels:
- hide-this
maxAgeOfPRs: 100
maxNumberOfPRs: 10
repos:
- plumming/repo1
- plumming/repo2`,

			expectations: func(t *testing.T, c config.Config) {
				assert.Equal(t, []string{"bot1", "bot2"}, c.GetBotAccounts())
				assert.Equal(t, []string{"hide-this"}, c.GetHiddenLabels())
				assert.Equal(t, 100, c.GetMaxAgeOfPRs())
				assert.Equal(t, 10, c.GetMaxNumberOfPRs())
				assert.Equal(t, []string{"github.com"}, c.GetConfiguredServers())
				assert.Equal(t, []string{"repo:plumming/repo1", "repo:plumming/repo2"}, c.GetReposToQuery("github.com"))
			},
		},
		{
			name: "default config",
			config: `---
`,

			expectations: func(t *testing.T, c config.Config) {
				assert.Equal(t, []string{"dependabot", "dependabot-preview"}, c.GetBotAccounts())
				assert.Equal(t, []string{"hide-this"}, c.GetHiddenLabels())
				assert.Equal(t, -1, c.GetMaxAgeOfPRs())
				assert.Equal(t, 100, c.GetMaxNumberOfPRs())
				assert.Equal(t, []string{"github.com"}, c.GetConfiguredServers())
				assert.Equal(t, []string{"repo:plumming/dx"}, c.GetReposToQuery("github.com"))
			},
		},
		{
			name: "partial config",
			config: `---
botAccounts:
- bot1
- bot2
maxAgeOfPRs: 100
maxNumberOfPRs: 10`,

			expectations: func(t *testing.T, c config.Config) {
				assert.Equal(t, []string{"bot1", "bot2"}, c.GetBotAccounts())
				assert.Equal(t, []string{"hide-this"}, c.GetHiddenLabels())
				assert.Equal(t, 100, c.GetMaxAgeOfPRs())
				assert.Equal(t, 10, c.GetMaxNumberOfPRs())
				assert.Equal(t, []string{"github.com"}, c.GetConfiguredServers())
				assert.Equal(t, []string{"repo:plumming/dx"}, c.GetReposToQuery("github.com"))
			},
		},
		{
			name: "migrated config",
			config: `---
botAccounts:
- bot1
- bot2
hiddenLabels:
- hide-this
maxAgeOfPRs: 100
maxNumberOfPRs: 10
repositories:
  github.com:
  - plumming/repo1
  - plumming/repo2`,

			expectations: func(t *testing.T, c config.Config) {
				assert.Equal(t, []string{"bot1", "bot2"}, c.GetBotAccounts())
				assert.Equal(t, []string{"hide-this"}, c.GetHiddenLabels())
				assert.Equal(t, 100, c.GetMaxAgeOfPRs())
				assert.Equal(t, 10, c.GetMaxNumberOfPRs())
				assert.Equal(t, []string{"github.com"}, c.GetConfiguredServers())
				assert.Equal(t, []string{"repo:plumming/repo1", "repo:plumming/repo2"}, c.GetReposToQuery("github.com"))
			},
		},
		{
			name: "multiple git server config",
			config: `---
botAccounts:
- bot1
- bot2
hiddenLabels:
- hide-this
maxAgeOfPRs: 100
maxNumberOfPRs: 10
repositories:
  github.com:
  - plumming/repo1
  - plumming/repo2
  other.com:
  - other/repo1
  - other/repo2`,

			expectations: func(t *testing.T, c config.Config) {
				assert.Equal(t, []string{"bot1", "bot2"}, c.GetBotAccounts())
				assert.Equal(t, []string{"hide-this"}, c.GetHiddenLabels())
				assert.Equal(t, 100, c.GetMaxAgeOfPRs())
				assert.Equal(t, 10, c.GetMaxNumberOfPRs())
				assert.Equal(t, []string{"github.com", "other.com"}, c.GetConfiguredServers())
				assert.Equal(t, []string{"repo:plumming/repo1", "repo:plumming/repo2"}, c.GetReposToQuery("github.com"))
				assert.Equal(t, []string{"repo:other/repo1", "repo:other/repo2"}, c.GetReposToQuery("other.com"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := config.Load(strings.NewReader(tc.config))
			assert.NoError(t, err)
			assert.NotNil(t, c)

			tc.expectations(t, c)
		})
	}
}
