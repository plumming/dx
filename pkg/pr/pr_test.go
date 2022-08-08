package pr

import (
	"testing"

	"github.com/plumming/dx/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestPullRequest_ContextsString(t *testing.T) {
	var tests = []struct {
		name        string
		contexts    []Context
		rollupState string
		exp         string
	}{
		{
			name:        "build_success",
			rollupState: "SUCCESS",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
			},
			exp: "SUCCESS",
		},
		{
			name:        "build_pending",
			rollupState: "SUCCESS",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "PENDING", Context: "Other-Build"},
			},
			exp: "PENDING",
		},
		{
			name:        "build_failing",
			rollupState: "SUCCESS",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "FAILURE", Context: "Other-Build"},
			},
			exp: "FAILURE",
		},
		{
			name: "build_success",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "PENDING", Context: "Merge Status"},
			},
			exp: "SUCCESS",
		},
		{
			name:        "build_pending",
			rollupState: "SUCCESS",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "PENDING", Context: "Other-Build"},
				{State: "PENDING", Context: "Merge Status"},
			},
			exp: "PENDING",
		},
		{
			name:        "build_failing",
			rollupState: "SUCCESS",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "FAILURE", Context: "Other-Build"},
				{State: "PENDING", Context: "Merge Status"},
			},
			exp: "FAILURE",
		},
		{
			name:        "build_error",
			rollupState: "SUCCESS",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "ERROR", Context: "Other-Build"},
				{State: "PENDING", Context: "Merge Status"},
			},
			exp: "ERROR",
		},
		{
			name:        "in_progress_pending",
			rollupState: "PENDING",
			contexts: []Context{
				{Context: "Build"},
			},
			exp: "PENDING",
		},
		{
			name:        "in_progress_success",
			rollupState: "SUCCESS",
			contexts: []Context{
				{Context: "Build"},
			},
			exp: "SUCCESS",
		},
		{
			name:        "in_progress_error",
			rollupState: "ERROR",
			contexts: []Context{
				{Context: "Build"},
			},
			exp: "ERROR",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pr := PullRequest{
				Commits: Commits{
					Nodes: []CommitEntry{
						{
							Commit: Commit{
								StatusCheckRollup: StatusCheckRollup{
									State: test.rollupState,
									Contexts: StatusContext{
										Nodes: test.contexts,
									},
								},
							},
						},
					},
				},
			}
			assert.Equal(t, pr.ContextsString(), test.exp)
		})
	}
}

func TestPullRequest_TrimmedTitle(t *testing.T) {
	pr := PullRequest{Title: "chore(deps): bump https://github.com/plumming/test_repo from 0.0.694 to 0.0.695"}
	assert.Equal(t, pr.TrimmedTitle(), "chore(deps): bump https://github.com/plumming/test_repo from 0.0.694 to 0.0...")

	pr = PullRequest{Title: "chore(deps): bump dependency versions"}
	assert.Equal(t, pr.TrimmedTitle(), "chore(deps): bump dependency versions")
}

func TestPullRequest_MergableString(t *testing.T) {
	pr := PullRequest{Mergeable: "MERGEABLE"}
	assert.Equal(t, pr.MergeableString(), "")

	pr = PullRequest{Mergeable: "CONFLICTING"}
	assert.Equal(t, pr.MergeableString(), "* Conflict")

	pr = PullRequest{Mergeable: "UNKNOWN"}
	assert.Equal(t, pr.MergeableString(), "* ?")
}

func TestPullRequest_PullsString(t *testing.T) {
	pr := PullRequest{URL: "https://github.com/plumming/dx/pull/257"}
	assert.Equal(t, pr.PullsString(), "https://github.com/plumming/dx/pulls")

	pr = PullRequest{URL: "https://github.com/plumming/dx/pull/1083"}
	assert.Equal(t, pr.PullsString(), "https://github.com/plumming/dx/pulls")
}

func TestPullRequest_HasLabel(t *testing.T) {
	pr := PullRequest{
		Labels: Labels{
			Nodes: []Label{{Name: "lgtm"}},
		},
	}
	assert.Equal(t, pr.HasLabel("updatebot"), false)

	pr = PullRequest{
		Labels: Labels{
			Nodes: []Label{{Name: "updatebot"}},
		},
	}
	assert.Equal(t, pr.HasLabel("updatebot"), true)

	pr = PullRequest{
		Labels: Labels{
			Nodes: []Label{{Name: "lgtm"}, {Name: "updatebot"}},
		},
	}
	assert.Equal(t, pr.HasLabel("updatebot"), true)
}

func TestPullRequest_HasContext(t *testing.T) {
	var tests = []struct {
		name     string
		contexts []Context
		passing  string
		failing  string
	}{
		{
			name: "build_success",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
			},
			passing: "Build",
			failing: "Other",
		},
		{
			name: "build_pending",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "PENDING", Context: "Other-Build"},
			},
			passing: "Other-Build",
			failing: "Other",
		},
		{
			name: "build_success",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "PENDING", Context: "Merge Status"},
			},
			passing: "Merge Status",
			failing: "this",
		},
		{
			name: "build_pending_check",
			contexts: []Context{
				{Conclusion: "SUCCESS", Name: "Build"},
				{Conclusion: "PENDING", Name: "Other-Build"},
			},
			passing: "Other-Build",
			failing: "Other",
		},
		{
			name: "build_success_check",
			contexts: []Context{
				{Conclusion: "SUCCESS", Name: "Build"},
				{Conclusion: "PENDING", Name: "Merge Status"},
			},
			passing: "Merge Status",
			failing: "this",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pr := PullRequest{
				Commits: Commits{
					Nodes: []CommitEntry{
						{
							Commit: Commit{
								StatusCheckRollup: StatusCheckRollup{
									Contexts: StatusContext{
										Nodes: test.contexts,
									},
								},
							},
						},
					},
				},
			}
			assert.Equal(t, pr.HasContext(test.passing), true)
			assert.Equal(t, pr.HasContext(test.failing), false)
		})
	}
}

func TestPullRequest_Display(t *testing.T) {
	pr := PullRequest{Closed: true}
	assert.Equal(t, pr.Display(), false)

	pr = PullRequest{Closed: false}
	assert.Equal(t, pr.Display(), true)
}

func TestPullRequest_LabelsString(t *testing.T) {
	pr := PullRequest{Labels: Labels{
		Nodes: []Label{{Name: "do-not-merge/hold"}},
	}}

	assert.Equal(t, pr.LabelsString(), "do-not-merge/hold")

	pr = PullRequest{Labels: Labels{
		Nodes: []Label{{Name: "lgtm"}, {Name: "do-not-merge/hold"}},
	}}

	assert.Equal(t, pr.LabelsString(), "lgtm, do-not-merge/hold")

	pr = PullRequest{Labels: Labels{
		Nodes: []Label{},
	}}

	assert.Equal(t, pr.LabelsString(), "")
}

func TestPullRequest_ColoredTitle(t *testing.T) {
	var tests = []struct {
		title                  string
		statusCheckRollup      StatusCheckRollup
		expected               string
		expectedContextsString string
	}{
		{
			title:                  "this is a short PR title",
			statusCheckRollup:      StatusCheckRollup{Contexts: StatusContext{Nodes: []Context{{Context: "Test", State: "SUCCESS"}}}},
			expected:               util.ColorInfo("this is a short PR title"),
			expectedContextsString: "SUCCESS",
		},
		{
			title:                  "this is a short PR title",
			statusCheckRollup:      StatusCheckRollup{Contexts: StatusContext{Nodes: []Context{{Context: "Test", State: "PENDING"}}}},
			expected:               util.ColorWarning("this is a short PR title"),
			expectedContextsString: "PENDING",
		},
		{
			title:                  "this is a short PR title",
			statusCheckRollup:      StatusCheckRollup{State: "PENDING", Contexts: StatusContext{Nodes: []Context{{Context: "Test", Name: "PR"}}}},
			expected:               util.ColorWarning("this is a short PR title"),
			expectedContextsString: "PENDING",
		},
		{
			title:                  "this is a short PR title",
			statusCheckRollup:      StatusCheckRollup{Contexts: StatusContext{Nodes: []Context{{Context: "Test", State: "FAILURE"}}}},
			expected:               util.ColorError("this is a short PR title"),
			expectedContextsString: "FAILURE",
		},
		{
			title:                  "this is a really really really really really really really really really really really really really really long PR title",
			statusCheckRollup:      StatusCheckRollup{Contexts: StatusContext{Nodes: []Context{{Context: "Test", State: "SUCCESS"}}}},
			expected:               util.ColorInfo("this is a really really really really really really really really really re..."),
			expectedContextsString: "SUCCESS",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			pr := PullRequest{
				Title: test.title,
				Commits: Commits{
					Nodes: []CommitEntry{
						{
							Commit: Commit{
								StatusCheckRollup: test.statusCheckRollup,
							},
						},
					},
				},
			}
			//t.Logf("ContextsString: %s", pr.ContextsString())
			assert.Equal(t, test.expectedContextsString, pr.ContextsString())
			assert.Equal(t, test.expected, pr.ColoredTitle())
		})
	}
}
