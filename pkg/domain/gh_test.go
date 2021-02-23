package domain_test

import (
	"bytes"
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

func TestGetCurrentUser(t *testing.T) {
	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))

	http.StubResponse(200, bytes.NewBufferString("{ \"login\":\"octocat\"}"))

	currentUser, err := domain.GetCurrentUser(client)
	assert.NoError(t, err)

	assert.Equal(t, currentUser, "octocat")
}

func TestGetOrgAndRepo(t *testing.T) {
	org, repo, err := domain.ExtractOrgAndRepoURL("https://github.com/clone/chilly")
	assert.NoError(t, err)
	assert.Equal(t, org, "clone")
	assert.Equal(t, repo, "chilly")

	org, repo, err = domain.ExtractOrgAndRepoURL("https://github.com/plumming/dx")
	assert.NoError(t, err)
	assert.Equal(t, org, "plumming")
	assert.Equal(t, repo, "dx")
}
