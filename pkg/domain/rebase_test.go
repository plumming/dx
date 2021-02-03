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
		name                  string
		remotes               string
		originDefaultBranch   string
		upstreamDefaultBranch string
		expectedCommands      []string
		expectedRequests      []string
	}

	tests := []test{
		{
			name: "simple rebase on master",
			remotes: `origin https://github.com/origin/clone (fetch)
origin https://github.com/origin/clone (push)
upstream https://github.com/upstream/repo (fetch)
upstream https://github.com/upstream/repo (push)`,
			originDefaultBranch:   "master",
			upstreamDefaultBranch: "master",
			expectedCommands: []string{
				"git remote -v",
				"git remote -v",
				"git status --porcelain",
				"git branch --show-current",
				"git fetch --tags upstream master",
				"git rebase upstream/master",
				"git push origin master",
			},
			expectedRequests: []string{
				"https://api.github.com/repos/origin/clone",
				"https://api.github.com/repos/upstream/repo",
			},
		},
		{
			name: "simple rebase on main",
			remotes: `origin https://github.com/origin/clone (fetch)
origin https://github.com/origin/clone (push)
upstream https://github.com/upstream/repo (fetch)
upstream https://github.com/upstream/repo (push)`,
			originDefaultBranch:   "main",
			upstreamDefaultBranch: "main",
			expectedCommands: []string{
				"git remote -v",
				"git remote -v",
				"git status --porcelain",
				"git branch --show-current",
				"git fetch --tags upstream main",
				"git rebase upstream/main",
				"git push origin main",
			},
			expectedRequests: []string{
				"https://api.github.com/repos/origin/clone",
				"https://api.github.com/repos/upstream/repo",
			},
		},
		{
			name: "simple rebase on main with no upstream",
			remotes: `origin https://github.com/origin/clone (fetch)
origin https://github.com/origin/clone (push)`,
			originDefaultBranch: "main",
			expectedCommands: []string{
				"git remote -v",
				"git remote -v",
				"git status --porcelain",
				"git branch --show-current",
				"git pull --tags origin main",
			},
			expectedRequests: []string{
				"https://api.github.com/repos/origin/clone",
			},
		},
		{
			name: "complex rebase on differing branches",
			remotes: `origin https://github.com/origin/clone (fetch)
origin https://github.com/origin/clone (push)
upstream https://github.com/upstream/repo (fetch)
upstream https://github.com/upstream/repo (push)`,
			originDefaultBranch:   "master",
			upstreamDefaultBranch: "main",
			expectedCommands: []string{
				"git remote -v",
				"git remote -v",
				"git status --porcelain",
				"git branch --show-current",
				"git fetch --tags upstream main",
				"git rebase upstream/main",
				"git push origin master",
			},
			expectedRequests: []string{
				"https://api.github.com/repos/origin/clone",
				"https://api.github.com/repos/upstream/repo",
			},
		},
		{
			name: "simple rebase on main with .git extension",
			remotes: `origin https://github.com/origin/clone.git (fetch)
origin https://github.com/origin/clone.git (push)`,
			originDefaultBranch: "main",
			expectedCommands: []string{
				"git remote -v",
				"git remote -v",
				"git status --porcelain",
				"git branch --show-current",
				"git pull --tags origin main",
			},
			expectedRequests: []string{
				"https://api.github.com/repos/origin/clone",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			http := &api.FakeHTTP{}
			client := api.NewClient(api.ReplaceTripper(http))

			http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf("{ \"default_branch\":\"%s\"}", tc.originDefaultBranch)))
			if tc.upstreamDefaultBranch != "" {
				http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf("{ \"default_branch\":\"%s\"}", tc.upstreamDefaultBranch)))
			}

			rb := domain.NewRebase()
			rb.OriginDefaultBranch = tc.originDefaultBranch
			rb.UpstreamDefaultBranch = tc.upstreamDefaultBranch
			rb.SetGithubClient(client)

			r := mocks.MockCommandRunner{}
			domain.Runner = &r
			mocks.GetRunWithoutRetryFunc = func(c *util.Command) (string, error) {
				if c.String() == "git branch --show-current" {
					return tc.originDefaultBranch, nil
				}

				if c.String() == "git remote -v" {
					return tc.remotes, nil
				}

				if c.String() == "git status --porcelain" {
					return "", nil
				}

				return "<dummy output>", nil
			}

			err := rb.Validate()
			assert.NoError(t, err)

			err = rb.Run()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCommands, r.Commands)

			assert.Equal(t, len(tc.expectedRequests), len(http.Requests))

			for i, request := range tc.expectedRequests {
				assert.Equal(t, request, http.Requests[i].URL.String())
			}
		})
	}
}
