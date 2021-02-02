package domain_test

import (
	"testing"

	"github.com/plumming/dx/pkg/domain"
	"github.com/plumming/dx/pkg/util"
	"github.com/plumming/dx/pkg/util/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCanRebase(t *testing.T) {
	type test struct {
		name          string
		defaultBranch string
		expected      []string
	}

	tests := []test{
		{
			name:          "simple rebase on master",
			defaultBranch: "master",
			expected: []string{
				"git status --porcelain",
				"git branch --show-current",
				"git fetch --tags upstream master",
				"git rebase upstream/master",
				"git push origin master",
			},
		},
		{
			name:          "simple rebase on main",
			defaultBranch: "main",
			expected: []string{
				"git status --porcelain",
				"git branch --show-current",
				"git fetch --tags upstream main",
				"git rebase upstream/main",
				"git push origin main",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rb := domain.NewRebase()
			rb.DefaultBranch = tc.defaultBranch

			r := mocks.MockCommandRunner{}
			domain.Runner = &r
			mocks.GetRunWithoutRetryFunc = func(c *util.Command) (string, error) {
				if c.String() == "git branch --show-current" {
					return tc.defaultBranch, nil
				}
				return "", nil
			}

			err := rb.Run()

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, r.Commands)
		})
	}
}
