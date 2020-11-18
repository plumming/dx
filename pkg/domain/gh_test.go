package domain_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/plumming/dx/pkg/api"
	"github.com/plumming/dx/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultBranch_Main(t *testing.T) {
	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))

	http.StubResponse(200, bytes.NewBufferString("{ \"default_branch\":\"main\" }"))

	defaultBranch, err := domain.GetDefaultBranch(client, "orgname", "reponame")
	assert.NoError(t, err)

	assert.Equal(t, defaultBranch, "main")
}

func TestGetDefaultBranch_Master(t *testing.T) {
	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))

	http.StubResponse(200, bytes.NewBufferString("{ \"default_branch\":\"master\"}"))

	defaultBranch, err := domain.GetDefaultBranch(client, "orgname", "reponame")
	assert.NoError(t, err)

	assert.Equal(t, defaultBranch, "master")
}

func TestGetOrgAndRepo(t *testing.T) {
	output := `origin	https://github.com/clone/chilly (fetch)
origin	https://github.com/clone/chilly (push)
upstream	https://github.com/plumming/dx (fetch)
upstream	https://github.com/plumming/dx (push)`

	org, repo, err := domain.ExtractOrgAndRepoFromGitRemotes(strings.NewReader(output))
	assert.NoError(t, err)
	assert.Equal(t, org, "clone")
	assert.Equal(t, repo, "chilly")
}
