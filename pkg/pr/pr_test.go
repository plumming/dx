package pr

import (
	"testing"

	"github.com/plumming/dx/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestPullRequest_ContextsString(t *testing.T) {
	var tests = []struct {
		name     string
		contexts []Context
		exp      string
	}{
		{
			name: "build_success",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
			},
			exp: "SUCCESS",
		},
		{
			name: "build_pending",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "PENDING", Context: "Other-Build"},
			},
			exp: "PENDING",
		},
		{
			name: "build_failing",
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
			name: "build_pending",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "PENDING", Context: "Other-Build"},
				{State: "PENDING", Context: "Merge Status"},
			},
			exp: "PENDING",
		},
		{
			name: "build_failing",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "FAILURE", Context: "Other-Build"},
				{State: "PENDING", Context: "Merge Status"},
			},
			exp: "FAILURE",
		},
		{
			name: "build_error",
			contexts: []Context{
				{State: "SUCCESS", Context: "Build"},
				{State: "ERROR", Context: "Other-Build"},
				{State: "PENDING", Context: "Merge Status"},
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
	assert.Equal(t, pr.Display(true, true), false)
	assert.Equal(t, pr.Display(true, false), false)
	assert.Equal(t, pr.Display(false, true), false)
	assert.Equal(t, pr.Display(false, false), false)

	pr = PullRequest{Author: Author{Login: "jenkins-x-bot"}}

	assert.Equal(t, pr.Display(true, true), true)
	assert.Equal(t, pr.Display(true, false), true)
	assert.Equal(t, pr.Display(false, true), true)
	assert.Equal(t, pr.Display(false, false), true)

	pr = PullRequest{Author: Author{Login: "dependabot-preview"}}

	assert.Equal(t, pr.Display(true, true), true)
	assert.Equal(t, pr.Display(true, false), true)
	assert.Equal(t, pr.Display(false, true), false)
	assert.Equal(t, pr.Display(false, false), false)

	pr = PullRequest{Labels: Labels{
		Nodes: []Label{{Name: "do-not-merge/hold"}},
	}}

	assert.Equal(t, pr.Display(true, true), true)
	assert.Equal(t, pr.Display(true, false), false)
	assert.Equal(t, pr.Display(false, true), true)
	assert.Equal(t, pr.Display(false, false), false)
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
	pr := PullRequest{
		Title:   "this is a short PR title",
		Commits: Commits{Nodes: []CommitEntry{{Commit: Commit{StatusCheckRollup: StatusCheckRollup{Contexts: StatusContext{Nodes: []Context{{State: "SUCCESS"}}}}}}}},
	}

	assert.Equal(t, pr.ColoredTitle(), util.ColorInfo("this is a short PR title"))

	pr = PullRequest{
		Title:   "this is a short PR title",
		Commits: Commits{Nodes: []CommitEntry{{Commit: Commit{StatusCheckRollup: StatusCheckRollup{Contexts: StatusContext{Nodes: []Context{{State: "PENDING"}}}}}}}},
	}

	assert.Equal(t, pr.ColoredTitle(), util.ColorInfo("this is a short PR title"))

	pr = PullRequest{
		Title:   "this is a short PR title",
		Commits: Commits{Nodes: []CommitEntry{{Commit: Commit{StatusCheckRollup: StatusCheckRollup{Contexts: StatusContext{Nodes: []Context{{State: "FAILING"}}}}}}}},
	}

	assert.Equal(t, pr.ColoredTitle(), util.ColorInfo("this is a short PR title"))

	pr = PullRequest{
		Title:   "this is a really really really really really really really really really really really really really really long PR title",
		Commits: Commits{Nodes: []CommitEntry{{Commit: Commit{StatusCheckRollup: StatusCheckRollup{Contexts: StatusContext{Nodes: []Context{{State: "SUCCESS"}}}}}}}},
	}

	assert.Equal(t, pr.ColoredTitle(), util.ColorInfo("this is a really really really really really really really really really re..."))
}
