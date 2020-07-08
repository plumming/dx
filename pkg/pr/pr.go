package pr

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/plumming/chilly/pkg/util"
)

type PullRequest struct {
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	URL        string     `json:"url"`
	Mergeable  string     `json:"mergeable"`
	CreatedAt  time.Time  `json:"createdAt"`
	Author     Author     `json:"author"`
	Labels     Labels     `json:"labels"`
	Commits    Commits    `json:"commits"`
	Closed     bool       `json:"closed"`
	Repository Repository `json:"repository"`
}

const (
	success     = "SUCCESS"
	pending     = "PENDING"
	failure     = "FAILURE"
	conflicting = "CONFLICTING"
	unknown     = "UNKNOWN"
)

func (p *PullRequest) Display(showDependabot bool, showOnHold bool) bool {
	display := true
	// exit early
	if p.Closed {
		return false
	}

	if p.Author.Login == "dependabot-preview" {
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
	for _, c := range p.Commits.Nodes[0].Commit.Status.Contexts {
		if c.Context != "tide" && c.Context != "keeper" && c.Context != "Merge Status" {
			contexts = append(contexts, c.State)
		}
	}
	return contexts
}

func (p *PullRequest) ContextsString() string {
	c := unique(p.contexts())
	if stringInSlice(c, failure) {
		return "FAILURE"
	} else if stringInSlice(c, pending) {
		return pending
	} else {
		return success
	}
}

func (p *PullRequest) FailedContexts() []Context {
	var failedContexts []Context
	for _, c := range p.Commits.Nodes[0].Commit.Status.Contexts {
		if c.Context != "tide" && c.Context != "keeper" && c.Context != "Merge Status" && c.State == failure {
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
	for _, label := range p.Commits.Nodes[0].Commit.Status.Contexts {
		if name == label.Context {
			return true
		}
	}
	return false
}

func (p *PullRequest) RequiresReview() bool {
	if p.MergeableString() == "" &&
		!p.HasLabel("updatebot") &&
		!p.HasLabel("needs-ok-to-test") &&
		!p.OnHold() &&
		!p.HasLabel("approved") &&
		!p.HasLabel("do-not-merge/work-in-progress") &&
		p.ContextsString() == success {
		return true
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
	Status CommitStatus `json:"status"`
}

type CommitStatus struct {
	Contexts []Context `json:"contexts"`
}

type Context struct {
	State       string `json:"state"`
	TargetURL   string `json:"targetUrl"`
	Description string `json:"description"`
	Context     string `json:"context"`
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
