package domain_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/plumming/dx/pkg/api"

	"github.com/plumming/dx/pkg/domain"
	"github.com/plumming/dx/pkg/util"
	"github.com/plumming/dx/pkg/util/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCanRebase(t *testing.T) {
	type test struct {
		name          string
		remotes       string
		defaultBranch string
		expected      []string
	}

	tests := []test{
		{
			name: "simple rebase on master",
			remotes: `origin https://github.com/origin/clone (fetch)
origin https://github.com/origin/clone (push)
upstream https://github.com/upstream/repo (fetch)
upstream https://github.com/upstream/repo (push)`,
			defaultBranch: "master",
			expected: []string{
				"git remote -v",
				"git remote -v",
				"git status --porcelain",
				"git branch --show-current",
				"git fetch --tags upstream master",
				"git rebase upstream/master",
				"git push origin master",
			},
		},
		{
			name: "simple rebase on main",
			remotes: `origin https://github.com/origin/clone (fetch)
origin https://github.com/origin/clone (push)
upstream https://github.com/upstream/repo (fetch)
upstream https://github.com/upstream/repo (push)`,
			defaultBranch: "main",
			expected: []string{
				"git remote -v",
				"git remote -v",
				"git status --porcelain",
				"git branch --show-current",
				"git fetch --tags upstream main",
				"git rebase upstream/main",
				"git push origin main",
			},
		},
		{
			name: "simple rebase on main with no upstream",
			remotes: `origin https://github.com/origin/clone (fetch)
origin https://github.com/origin/clone (push)`,
			defaultBranch: "main",
			expected: []string{
				"git remote -v",
				"git remote -v",
				"git status --porcelain",
				"git branch --show-current",
				"git pull --tags origin main",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			http := &api.FakeHTTP{}
			client := api.NewClient(api.ReplaceTripper(http))

			http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf("{ \"default_branch\":\"%s\"}", tc.defaultBranch)))

			rb := domain.NewRebase()
			rb.DefaultBranch = tc.defaultBranch
			rb.SetGithubClient(client)

			r := mocks.MockCommandRunner{}
			domain.Runner = &r
			mocks.GetRunWithoutRetryFunc = func(c *util.Command) (string, error) {
				if c.String() == "git branch --show-current" {
					return tc.defaultBranch, nil
				}

				if c.String() == "git remote -v" {
					return tc.remotes, nil
				}

				return "", nil
			}

			err := rb.Validate()
			assert.NoError(t, err)

			err = rb.Run()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, r.Commands)
		})
	}
}
