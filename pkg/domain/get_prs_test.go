package domain

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/plumming/dx/pkg/config"

	"github.com/plumming/dx/pkg/pr"

	"github.com/plumming/dx/pkg/api"

	"github.com/stretchr/testify/assert"
)

var expectedResponse = `{"data":
	{"search":
		{"nodes":[
			{"number":425,"title":"chore(deps): bump  versions","url":"https://github.com/plumming/test_repo/pull/425","createdAt":"2020-05-15T12:09:46Z","closed":false,"author":{"login":"dependabot-preview"},"repository":{"nameWithOwner":"plumming/test_repo"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"https://dashboard/PR-425/1","description":"Pipeline running stage(s): create","context":"test"},{"state":"SUCCESS","targetUrl":"","description":"In merge pool.","context":"Merge Status"}]}}}]}},
			{"number":273,"title":"chore(deps): bump https://github.com/plumming/test_repo from v0.0.427 to 0.0.428","url":"https://github.com/plumming/test_repo_two/pull/273","createdAt":"2020-05-15T11:47:46Z","closed":true,"author":{"login":"me-bot"},"repository":{"nameWithOwner":"plumming/test_repo_two"},"mergeable":"UNKNOWN","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"SUCCESS","targetUrl":"","description":"In merge pool.","context":"Merge Status"},{"state":"SUCCESS","targetUrl":"https://dashboard/PR-273/1","description":"Pipeline successful","context":"pr-build"}]}}}]}},
			{"number":793,"title":"chore(deps): bump https://github.com/plumming/jxui-frontend from 0.0.1243 to 0.0.1250","url":"https://github.com/test_repo_two/pull/793","createdAt":"2020-05-15T10:18:04Z","closed":false,"author":{"login":"cjxd-bot"},"repository":{"nameWithOwner":"my_owner"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"needs-ok-to-test"},{"name":"size/XS"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"}]}}}]}},
			{"number":157,"title":"chore(deps): bump https://github.com/plumming/jxui-frontend from 0.0.1249 to 0.0.1250","url":"https://github.com/test_repo_two/pull/157","createdAt":"2020-05-15T10:17:51Z","closed":false,"author":{"login":"cjxd-bot"},"repository":{"nameWithOwner":"my_owner"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"needs-ok-to-test"},{"name":"size/XS"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"}]}}}]}},
			{"number":792,"title":"chore: test to 0.0.923","url":"https://github.com/test_repo_two/pull/792","createdAt":"2020-05-15T10:07:30Z","closed":false,"author":{"login":"-bot"},"repository":{"nameWithOwner":"my_owner"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"},{"state":"SUCCESS","targetUrl":null,"description":"All Tasks have completed executing","context":"promotion-build"}]}}}]}},
			{"number":44,"title":"chore: prompt for version if not supplied","url":"https://github.com/plumming/dx/pull/44","createdAt":"2020-05-15T08:07:26Z","closed":false,"author":{"login":"me"},"repository":{"nameWithOwner":"plumming/dx"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XL"},{"name":"do-not-merge/hold"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"","description":"Not mergeable. Needs approved label.","context":"Merge Status"},{"state":"SUCCESS","targetUrl":"https://dashboard/PR-44/4","description":"Pipeline successful","context":"pr-build"}]}}}]}},
     		{"number":791,"title":"chore: my-service to 0.0.717","url":"https://github.com/test_repo_two/pull/791","createdAt":"2020-05-14T23:12:56Z","closed":false,"author":{"login":"bot"},"repository":{"nameWithOwner":"my_owner"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"},{"state":"SUCCESS","targetUrl":null,"description":"All Tasks have completed executing","context":"promotion-build"}]}}}]}},
			{"number":790,"title":"chore: test to 0.0.922","url":"https://github.com/test_repo_two/pull/790","createdAt":"2020-05-14T22:53:14Z","closed":false,"author":{"login":"bot"},"repository":{"nameWithOwner":"testOwners"},"mergeable":"MERGEABLE","labels":{"nodes":[{"name":"size/XS"},{"name":"updatebot"}]},"commits":{"nodes":[{"commit":{"status":{"contexts":[{"state":"PENDING","targetUrl":"http://dec","description":"Not mergeable. Needs approved label.","context":"tide"},{"state":"SUCCESS","targetUrl":null,"description":"All Tasks have completed executing","context":"promotion-build"}]}}}]}}
		]}
	}
}`

func TestGetPrs_Validate(t *testing.T) {
	d := NewGetPrs()

	err := d.Validate()
	assert.NoError(t, err)
}

func TestGetPrs_Run(t *testing.T) {
	d := NewGetPrs()

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	d.SetGithubClient(client)
	d.SetConfig(&config.Config{MaxAge: -1})

	http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf(expectedResponse)))

	err := d.Run()
	assert.NoError(t, err)

	assert.Equal(t, len(http.Requests), 1)
	assert.Equal(t, http.Requests[0].URL.Path, "/graphql")
	assert.Equal(t, http.Requests[0].URL.Host, "api.github.com")
	assert.Equal(t, http.Requests[0].URL.RawQuery, "")

	assert.Equal(t, 5, len(d.PullRequests))
}

func TestGetPrs_Run_ShowOnHold(t *testing.T) {
	d := NewGetPrs()
	d.ShowOnHold = true

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	d.SetGithubClient(client)
	d.SetConfig(&config.Config{MaxAge: -1})

	http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf(expectedResponse)))

	err := d.Run()
	assert.NoError(t, err)

	assert.Equal(t, len(http.Requests), 1)
	assert.Equal(t, http.Requests[0].URL.Path, "/graphql")
	assert.Equal(t, http.Requests[0].URL.Host, "api.github.com")
	assert.Equal(t, http.Requests[0].URL.RawQuery, "")

	assert.Equal(t, 6, len(d.PullRequests))
}

func TestGetPrs_Run_ShowDependabot(t *testing.T) {
	d := NewGetPrs()
	d.ShowDependabot = true

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	d.SetGithubClient(client)
	d.SetConfig(&config.Config{MaxAge: -1})

	http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf(expectedResponse)))

	err := d.Run()
	assert.NoError(t, err)

	assert.Equal(t, len(http.Requests), 1)
	assert.Equal(t, http.Requests[0].URL.Path, "/graphql")
	assert.Equal(t, http.Requests[0].URL.Host, "api.github.com")
	assert.Equal(t, http.Requests[0].URL.RawQuery, "")

	assert.Equal(t, 6, len(d.PullRequests))
}

func TestGetPrs_Retrigger(t *testing.T) {
	d := NewGetPrs()

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	d.SetGithubClient(client)

	prs := []pr.PullRequest{
		{
			URL:        "https://github.com/org/not-mergeable/pulls/1",
			Repository: pr.Repository{NameWithOwner: "org/not-mergeable"},
			Number:     1,
			Mergeable:  "CONFLICTING",
			Commits: pr.Commits{
				Nodes: []pr.CommitEntry{
					{
						Commit: pr.Commit{
							StatusCheckRollup: pr.StatusCheckRollup{
								Contexts: pr.StatusContext{[]pr.Context{}},
							},
						},
					},
				},
			},
		},
		{
			URL:        "https://github.com/org/mergeable/pulls/2",
			Repository: pr.Repository{NameWithOwner: "org/mergeable"},
			Number:     2,
			Mergeable:  "MERGEABLE",
			Commits: pr.Commits{
				Nodes: []pr.CommitEntry{
					{
						Commit: pr.Commit{
							StatusCheckRollup: pr.StatusCheckRollup{
								Contexts: pr.StatusContext{[]pr.Context{{Context: "build", State: "FAILURE"}}},
							},
						},
					},
				},
			},
			Labels: pr.Labels{Nodes: []pr.Label{{Name: "updatebot"}}},
		},
		{
			URL:        "https://github.com/org/mergeable-not-updatebot/pulls/3",
			Repository: pr.Repository{NameWithOwner: "org/mergeable-not-updatebot"},
			Number:     3,
			Mergeable:  "MERGEABLE",
			Commits: pr.Commits{
				Nodes: []pr.CommitEntry{
					{
						Commit: pr.Commit{
							StatusCheckRollup: pr.StatusCheckRollup{
								Contexts: pr.StatusContext{[]pr.Context{{Context: "build", State: "FAILURE"}}},
							},
						},
					},
				},
			},
			Labels: pr.Labels{Nodes: []pr.Label{{Name: "size/M"}}},
		},
		{
			URL:        "https://github.com/org/mergeable-but-pending/pulls/4",
			Repository: pr.Repository{NameWithOwner: "org/mergeable-but-pending"},
			Number:     4,
			Mergeable:  "MERGEABLE",
			Commits: pr.Commits{
				Nodes: []pr.CommitEntry{
					{
						Commit: pr.Commit{
							StatusCheckRollup: pr.StatusCheckRollup{
								Contexts: pr.StatusContext{[]pr.Context{{Context: "build", State: "PENDING"}}},
							},
						},
					},
				},
			},
			Labels: pr.Labels{Nodes: []pr.Label{{Name: "updatebot"}}},
		},
		{
			URL:        "https://github.com/org/plumming/pulls/5",
			Repository: pr.Repository{NameWithOwner: "org/plumming"},
			Number:     5,
			Mergeable:  "MERGEABLE",
			Commits: pr.Commits{
				Nodes: []pr.CommitEntry{
					{
						Commit: pr.Commit{
							StatusCheckRollup: pr.StatusCheckRollup{
								Contexts: pr.StatusContext{[]pr.Context{{Context: "plumming", State: "FAILURE"},
									{Context: "build", State: "FAILURE"}}},
							},
						},
					},
				},
			},
			Labels: pr.Labels{Nodes: []pr.Label{{Name: "updatebot"}}},
		},
	}

	d.PullRequests = prs

	http.StubResponse(200, bytes.NewBufferString("{}"))
	http.StubResponse(200, bytes.NewBufferString("{}"))
	http.StubResponse(200, bytes.NewBufferString("{}"))

	err := d.Retrigger()
	assert.NoError(t, err)

	assert.Equal(t, len(http.Requests), 3)
	assert.Equal(t, http.Requests[0].URL.Path, "/repos/org/mergeable/issues/2/comments")
	assert.Equal(t, http.Requests[0].URL.Host, "api.github.com")
	assert.Equal(t, http.Requests[0].URL.RawQuery, "")

	assert.Equal(t, http.Requests[1].URL.Path, "/repos/org/plumming/issues/5/comments")
	assert.Equal(t, http.Requests[1].URL.Host, "api.github.com")
	assert.Equal(t, http.Requests[1].URL.RawQuery, "")
}
