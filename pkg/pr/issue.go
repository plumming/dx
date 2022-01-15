package pr

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/plumming/dx/pkg/util"
)

type Issue struct {
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	URL        string     `json:"url"`
	CreatedAt  time.Time  `json:"createdAt"`
	Author     Author     `json:"author"`
	Labels     Labels     `json:"labels"`
	Closed     bool       `json:"closed"`
	Repository Repository `json:"repository"`
}

func (p *Issue) Display() bool {
	display := true
	// exit early
	if p.Closed {
		return false
	}

	return display
}

func (p *Issue) LabelsString() string {
	labels := []string{}
	for _, label := range p.Labels.Nodes {
		labels = append(labels, label.Name)
	}
	return strings.Join(labels, ", ")
}

func (p *Issue) ColoredTitle() string {
	return util.ColorInfo(p.TrimmedTitle())
}

func (p *Issue) TrimmedTitle() string {
	if len(p.Title) > 75 {
		return fmt.Sprintf("%s...", p.Title[:75])
	}
	return p.Title
}

func (p *Issue) IssueString() string {
	r, _ := regexp.Compile("issues/[0-9]+")
	return r.ReplaceAllString(p.URL, "issues")
}

func (p *Issue) HasLabel(name string) bool {
	for _, label := range p.Labels.Nodes {
		if name == label.Name {
			return true
		}
	}
	return false
}

type ByIssuesString []Issue

func (p ByIssuesString) Len() int {
	return len(p)
}

func (p ByIssuesString) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByIssuesString) Less(i, j int) bool {
	if p[i].IssueString() == p[j].IssueString() {
		return p[i].Number < p[j].Number
	}
	return p[i].IssueString() < p[j].IssueString()
}
