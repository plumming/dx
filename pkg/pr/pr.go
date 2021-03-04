package pr

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/plumming/dx/pkg/util"
)

type PullRequest struct {
	Number         int        `json:"number"`
	Title          string     `json:"title"`
	URL            string     `json:"url"`
	Mergeable      string     `json:"mergeable"`
	CreatedAt      time.Time  `json:"createdAt"`
	Author         Author     `json:"author"`
	Labels         Labels     `json:"labels"`
	Commits        Commits    `json:"commits"`
	Closed         bool       `json:"closed"`
	Repository     Repository `json:"repository"`
	ReviewDecision string     `json:"reviewDecision"`
}

const (
	success     = "SUCCESS"
	pending     = "PENDING"
	failure     = "FAILURE"
	error       = "ERROR"
	conflicting = "CONFLICTING"
	unknown     = "UNKNOWN"
)

func (p *PullRequest) Display(showDependabot bool, showOnHold bool, hiddenLabels ...string) bool {
	display := true
	// exit early
	if p.Closed {
		return false
	}

	for _, label := range hiddenLabels {
		if p.HasLabel(label) {
			return false
		}
	}

	if p.Author.Login == "dependabot-preview" || p.Author.Login == "dependabot" {
		display = showDependabot
	}

	if p.OnHold() {
		display = showOnHold
	}

	return display
}

func (p *PullRequest) LabelsString() string {
	labels := []string{}
	for _, label := range p.Labels.Nodes {
		labels = append(labels, label.Name)
	}
	return strings.Join(labels, ", ")
}

func (p *PullRequest) contexts() []string {
	contexts := []string{}
	for _, c := range p.Commits.Nodes[0].Commit.StatusCheckRollup.Contexts.Nodes {
		if c.Context != "" && c.Context != "tide" && c.Context != "keeper" && c.Context != "Merge Status" {
			contexts = append(contexts, c.State)
		}
		if c.Name != "" {
			contexts = append(contexts, c.Conclusion)
		}
	}
	return contexts
}

func (p *PullRequest) ContextsString() string {
	c := unique(p.contexts())
	if stringInSlice(c, error) {
		return "ERROR"
	} else if stringInSlice(c, failure) {
		return "FAILURE"
	} else if stringInSlice(c, pending) {
		return pending
	} else {
		return success
	}
}

func (p *PullRequest) FailedContexts() []Context {
	var failedContexts []Context
	for _, c := range p.Commits.Nodes[0].Commit.StatusCheckRollup.Contexts.Nodes {
		if c.Context != "" && c.Context != "tide" && c.Context != "keeper" && c.Context != "Merge Status" && c.State == failure {
			failedContexts = append(failedContexts, c)
		}
		if c.Name != "" && (c.Conclusion == failure || c.Conclusion == error) {
			failedContexts = append(failedContexts, c)
		}
	}
	return failedContexts
}

func (p *PullRequest) ColoredTitle() string {
	if p.ContextsString() == success {
		return util.ColorInfo(p.TrimmedTitle())
	} else if p.ContextsString() == pending {
		return util.ColorWarning(p.TrimmedTitle())
	} else {
		return util.ColorError(p.TrimmedTitle())
	}
}

func (p *PullRequest) ColoredReviewDecision() string {
	if p.ReviewDecision == "APPROVED" {
		return util.ColorInfo("Approved")
	} else if p.ReviewDecision == "REVIEW_REQUIRED" {
		return util.ColorWarning("Required")
	} else if p.ReviewDecision == "CHANGES_REQUESTED" {
		return util.ColorError("Changes Requested")
	} else {
		return p.ReviewDecision
	}
}

func (p *PullRequest) TrimmedTitle() string {
	if len(p.Title) > 75 {
		return fmt.Sprintf("%s...", p.Title[:75])
	}
	return p.Title
}

func (p *PullRequest) MergeableString() string {
	if p.Mergeable == conflicting {
		return "* Conflict"
	}
	if p.Mergeable == unknown {
		return "* ?"
	}
	return ""
}

func (p *PullRequest) OnHold() bool {
	return p.HasLabel("do-not-merge/hold")
}

func (p *PullRequest) PullsString() string {
	r, _ := regexp.Compile("pull/[0-9]+")
	return r.ReplaceAllString(p.URL, "pulls")
}

func (p *PullRequest) HasLabel(name string) bool {
	for _, label := range p.Labels.Nodes {
		if name == label.Name {
			return true
		}
	}
	return false
}

func (p *PullRequest) HasContext(name string) bool {
	for _, label := range p.Commits.Nodes[0].Commit.StatusCheckRollup.Contexts.Nodes {
		if name == label.Context {
			return true
		}
		if name == label.Name {
			return true
		}
	}
	return false
}

type Repository struct {
	NameWithOwner string `json:"nameWithOwner"`
}

type Author struct {
	Login string `json:"login"`
}

type Labels struct {
	Nodes []Label `json:"nodes"`
}

type Label struct {
	Name string `json:"name"`
}

type Commits struct {
	Nodes []CommitEntry `json:"nodes"`
}

type CommitEntry struct {
	Commit Commit `json:"commit"`
}

type Commit struct {
	StatusCheckRollup StatusCheckRollup `json:"statusCheckRollup"`
}

type StatusCheckRollup struct {
	State    string        `json:"state"`
	Contexts StatusContext `json:"contexts"`
}

type StatusContext struct {
	Nodes []Context `json:"nodes"`
}

type Context struct {
	State       string `json:"state"`
	Description string `json:"description"`
	Context     string `json:"context"`
	Conclusion  string `json:"conclusion"`
	Name        string `json:"name"`
	Title       string `json:"title"`
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func stringInSlice(stringSlice []string, a string) bool {
	for _, b := range stringSlice {
		if b == a {
			return true
		}
	}
	return false
}

type ByPullsString []PullRequest

func (p ByPullsString) Len() int {
	return len(p)
}

func (p ByPullsString) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByPullsString) Less(i, j int) bool {
	if p[i].PullsString() == p[j].PullsString() {
		return p[i].Number < p[j].Number
	}
	return p[i].PullsString() < p[j].PullsString()
}
