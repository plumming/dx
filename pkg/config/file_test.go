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
				assert.Equal(t, []string{"bot1", "bot2"}, c.BotAccounts)
				assert.Equal(t, []string{"hide-this"}, c.HiddenLabels)
				assert.Equal(t, 100, c.MaxAge)
				assert.Equal(t, 10, c.MaxNumberOfPRs)
				assert.Equal(t, []string{"plumming/repo1", "plumming/repo2"}, c.Repos)
			},
		},
		{
			name: "default config",
			config: `---
`,

			expectations: func(t *testing.T, c config.Config) {
				assert.Equal(t, []string{"dependabot", "dependabot-preview"}, c.BotAccounts)
				assert.Equal(t, []string{"hide-this"}, c.HiddenLabels)
				assert.Equal(t, -1, c.MaxAge)
				assert.Equal(t, 100, c.MaxNumberOfPRs)
				assert.Equal(t, []string{"plumming/dx"}, c.Repos)
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
				assert.Equal(t, []string{"bot1", "bot2"}, c.BotAccounts)
				assert.Equal(t, []string{"hide-this"}, c.HiddenLabels)
				assert.Equal(t, 100, c.MaxAge)
				assert.Equal(t, 10, c.MaxNumberOfPRs)
				assert.Equal(t, []string{"plumming/dx"}, c.Repos)
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
