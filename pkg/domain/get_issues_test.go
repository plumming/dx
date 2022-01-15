package domain

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/plumming/dx/pkg/auth"

	"github.com/plumming/dx/pkg/config/configfakes"

	"github.com/plumming/dx/pkg/api"

	"github.com/stretchr/testify/assert"
)

var userIssuesResponse = `{}`

func TestGetIssues_Validate(t *testing.T) {
	d := NewGetIssues()

	err := d.Validate()
	assert.NoError(t, err)
}

func TestGetIssues_Run(t *testing.T) {
	d := NewGetIssues()

	var authConfig auth.Config = &auth.FakeConfig{
		Hosts: map[string]*auth.HostConfig{
			"{{.Host}}": {User: "user", Token: "token"},
			"other.com": {User: "otheruser", Token: "token2"},
		},
	}

	dxConfig := configfakes.FakeConfig{}
	dxConfig.GetMaxAgeOfPRsReturns(-1)
	dxConfig.GetBotAccountsReturns([]string{"dependabot-preview"})
	dxConfig.GetHiddenLabelsReturns([]string{"do-not-merge/hold"})
	dxConfig.GetConfiguredServersReturns([]string{"github.com", "other.com"})
	dxConfig.GetReposToQueryReturnsOnCall(0, []string{"github/repo1", "github/repo2"})
	dxConfig.GetReposToQueryReturnsOnCall(1, []string{"other/repo1", "other/repo2"})

	d.SetDxConfig(&dxConfig)

	http := &api.FakeHTTP{}
	client := api.NewClient(authConfig, api.ReplaceTripper(http))
	d.SetGithubClient(client)

	// github.com
	http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf(userIssuesResponse)))
	http.StubResponse(200, bytes.NewBufferString(expectedIssuesResponse("github.com")))
	// other.com
	http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf(userIssuesResponse)))
	http.StubResponse(200, bytes.NewBufferString(expectedIssuesResponse("other.com")))

	err := d.Run()
	assert.NoError(t, err)

	assert.Equal(t, 4, len(http.Requests))

	assert.Equal(t, "/user", http.Requests[0].URL.Path)
	assert.Equal(t, "api.github.com", http.Requests[0].URL.Host)
	assert.Equal(t, "", http.Requests[0].URL.RawQuery)

	assert.Equal(t, "/graphql", http.Requests[1].URL.Path)
	assert.Equal(t, "api.github.com", http.Requests[1].URL.Host)
	assert.Equal(t, "", http.Requests[1].URL.RawQuery)

	assert.Equal(t, "/api/v3/user", http.Requests[2].URL.Path)
	assert.Equal(t, "other.com", http.Requests[2].URL.Host)
	assert.Equal(t, "", http.Requests[2].URL.RawQuery)

	assert.Equal(t, "/api/graphql", http.Requests[3].URL.Path)
	assert.Equal(t, "other.com", http.Requests[3].URL.Host)
	assert.Equal(t, "", http.Requests[3].URL.RawQuery)

	assert.Equal(t, 10, len(d.Issues))
}

func TestGetIssues_Run_ShowOnHold(t *testing.T) {
	d := NewGetIssues()
	d.ShowHidden = true

	var authConfig auth.Config = &auth.FakeConfig{
		Hosts: map[string]*auth.HostConfig{
			"github.com": {User: "user", Token: "token"},
		},
	}

	http := &api.FakeHTTP{}
	client := api.NewClient(authConfig, api.ReplaceTripper(http))
	d.SetGithubClient(client)

	dxConfig := configfakes.FakeConfig{}
	dxConfig.GetMaxAgeOfPRsReturns(-1)
	dxConfig.GetBotAccountsReturns([]string{"dependabot-preview"})
	dxConfig.GetHiddenLabelsReturns([]string{"do-not-merge/hold"})
	dxConfig.GetConfiguredServersReturns([]string{"github.com"})
	dxConfig.GetReposToQueryReturns([]string{"github/repo1", "github/repo2"})

	d.SetDxConfig(&dxConfig)

	http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf(userIssuesResponse)))
	http.StubResponse(200, bytes.NewBufferString(expectedIssuesResponse("github.com")))

	err := d.Run()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(http.Requests))

	assert.Equal(t, "/user", http.Requests[0].URL.Path)
	assert.Equal(t, "api.github.com", http.Requests[0].URL.Host)
	assert.Equal(t, "", http.Requests[0].URL.RawQuery)

	assert.Equal(t, "/graphql", http.Requests[1].URL.Path)
	assert.Equal(t, "api.github.com", http.Requests[1].URL.Host)
	assert.Equal(t, "", http.Requests[1].URL.RawQuery)

	assert.Equal(t, len(d.Issues), 6)
}

func TestGetIssues_Run_ShowBots(t *testing.T) {
	d := NewGetIssues()
	d.ShowBots = true

	var authConfig auth.Config = &auth.FakeConfig{
		Hosts: map[string]*auth.HostConfig{
			"github.com": {User: "user", Token: "token"},
		},
	}

	http := &api.FakeHTTP{}
	client := api.NewClient(authConfig, api.ReplaceTripper(http))
	d.SetGithubClient(client)

	dxConfig := configfakes.FakeConfig{}
	dxConfig.GetMaxAgeOfPRsReturns(-1)
	dxConfig.GetBotAccountsReturns([]string{"dependabot-preview"})
	dxConfig.GetHiddenLabelsReturns([]string{"do-not-merge/hold"})
	dxConfig.GetConfiguredServersReturns([]string{"github.com"})
	dxConfig.GetReposToQueryReturns([]string{"github/repo1", "github/repo2"})

	d.SetDxConfig(&dxConfig)

	http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf(userIssuesResponse)))
	http.StubResponse(200, bytes.NewBufferString(expectedIssuesResponse("github.com")))

	err := d.Run()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(http.Requests))

	assert.Equal(t, "/user", http.Requests[0].URL.Path)
	assert.Equal(t, "api.github.com", http.Requests[0].URL.Host)
	assert.Equal(t, "", http.Requests[0].URL.RawQuery)

	assert.Equal(t, "/graphql", http.Requests[1].URL.Path)
	assert.Equal(t, "api.github.com", http.Requests[1].URL.Host)
	assert.Equal(t, "", http.Requests[1].URL.RawQuery)

	assert.Equal(t, 6, len(d.Issues))
}

func expectedIssuesResponse(host string) string {
	variables := make(map[string]string)
	variables["Host"] = host

	tmpl, err := template.New("test").Parse(`{"data":
	{"search":
		{"nodes":[
			{"number":425,"title":"chore(deps): bump  versions","url":"https://{{.Host}}/plumming/test_repo/pull/425","createdAt":"2020-05-15T12:09:46Z","closed":false,"author":{"login":"dependabot-preview"},"repository":{"nameWithOwner":"plumming/test_repo"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"https://dashboard/PR-425/1","description":"Pipeline running stage(s): create","context":"test"},{"state":"SUCCESS","targetUrl":"","description":"In merge pool.","context":"Merge Status"}]}}}]}},
			{"number":273,"title":"chore(deps): bump https://{{.Host}}/plumming/test_repo from v0.0.427 to 0.0.428","url":"https://{{.Host}}/plumming/test_repo_two/pull/273","createdAt":"2020-05-15T11:47:46Z","closed":true,"author":{"login":"me-bot"},"repository":{"nameWithOwner":"plumming/test_repo_two"},"mergeable":"UNKNOWN","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"SUCCESS","targetUrl":"","description":"In merge pool.","context":"Merge Status"},{"state":"SUCCESS","targetUrl":"https://dashboard/PR-273/1","description":"Pipeline successful","context":"pr-build"}]}}}]}},
			{"number":793,"title":"chore(deps): bump https://{{.Host}}/plumming/jxui-frontend from 0.0.1243 to 0.0.1250","url":"https://{{.Host}}/test_repo_two/pull/793","createdAt":"2020-05-15T10:18:04Z","closed":false,"author":{"login":"cjxd-bot"},"repository":{"nameWithOwner":"my_owner"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"needs-ok-to-test"},{"name":"size/XS"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"}]}}}]}},
			{"number":157,"title":"chore(deps): bump https://{{.Host}}/plumming/jxui-frontend from 0.0.1249 to 0.0.1250","url":"https://{{.Host}}/test_repo_two/pull/157","createdAt":"2020-05-15T10:17:51Z","closed":false,"author":{"login":"cjxd-bot"},"repository":{"nameWithOwner":"my_owner"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"needs-ok-to-test"},{"name":"size/XS"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"}]}}}]}},
			{"number":792,"title":"chore: test to 0.0.923","url":"https://{{.Host}}/test_repo_two/pull/792","createdAt":"2020-05-15T10:07:30Z","closed":false,"author":{"login":"-bot"},"repository":{"nameWithOwner":"my_owner"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"},{"state":"SUCCESS","targetUrl":null,"description":"All Tasks have completed executing","context":"promotion-build"}]}}}]}},
			{"number":44,"title":"chore: prompt for version if not supplied","url":"https://{{.Host}}/plumming/dx/pull/44","createdAt":"2020-05-15T08:07:26Z","closed":false,"author":{"login":"me"},"repository":{"nameWithOwner":"plumming/dx"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XL"},{"name":"do-not-merge/hold"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"","description":"Not mergeable. Needs approved label.","context":"Merge Status"},{"state":"SUCCESS","targetUrl":"https://dashboard/PR-44/4","description":"Pipeline successful","context":"pr-build"}]}}}]}},
     		{"number":791,"title":"chore: my-service to 0.0.717","url":"https://{{.Host}}/test_repo_two/pull/791","createdAt":"2020-05-14T23:12:56Z","closed":false,"author":{"login":"bot"},"repository":{"nameWithOwner":"my_owner"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"},{"state":"SUCCESS","targetUrl":null,"description":"All Tasks have completed executing","context":"promotion-build"}]}}}]}},
			{"number":790,"title":"chore: test to 0.0.922","url":"https://{{.Host}}/test_repo_two/pull/790","createdAt":"2020-05-14T22:53:14Z","closed":false,"author":{"login":"bot"},"repository":{"nameWithOwner":"testOwners"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"},{"state":"SUCCESS","targetUrl":null,"description":"All Tasks have completed executing","context":"promotion-build"}]}}}]}}
		]}
	}
}`)
	if err != nil {
		panic(err)
	}
	var output bytes.Buffer
	err = tmpl.Execute(&output, variables)
	if err != nil {
		panic(err)
	}

	return output.String()
}
