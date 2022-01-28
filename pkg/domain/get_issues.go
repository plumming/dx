package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/plumming/dx/pkg/config"

	"github.com/plumming/dx/pkg/pr"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/cmd"
)

var (
	getIssuesQuery = `{
	search(query: "is:issue is:open %s", type: ISSUE, first: %d) {
      nodes {
      ... on Issue {
        number
        title
        url
        createdAt
        closed
        author {
          login
        }
        repository {
          nameWithOwner
        }
        comments {
          totalCount
        }
        labels(first: 10) {
          nodes {
            name
          }
        }
      }
    }
  }
}`
)

// GetIssues defines get pull issues response.
type GetIssues struct {
	cmd.CommonOptions
	ShowBots            bool
	ShowHidden          bool
	Me                  bool
	Review              bool
	Raw                 string
	Issues              []pr.Issue
	FilteredLabels      int
	FilteredBotAccounts int
}

// IssueData.
type IssueData struct {
	Search IssueSearch `json:"search"`
}

// IssueSearch.
type IssueSearch struct {
	Issues []pr.Issue `json:"nodes"`
}

// NewGetIssues.
func NewGetIssues() *GetIssues {
	g := &GetIssues{}
	return g
}

// Validate input.
func (g *GetIssues) Validate() error {
	return nil
}

// Run the cmd.
func (g *GetIssues) Run() error {
	cfg, err := g.DxConfig()
	if err != nil {
		return err
	}

	var issues []pr.Issue

	for _, host := range cfg.GetConfiguredServers() {
		i, err := g.GetIssuesForHost(host, cfg, getIssuesQuery)
		if err != nil {
			return err
		}
		issues = append(issues, i...)
	}

	sort.Sort(pr.ByIssuesString(issues))

	var issuesToReturn []pr.Issue

	filteredOnLabels := 0
	filteredOnAccounts := 0

	for _, issue := range issues {
		if issue.Display() {
			if g.filterOnLabels(issue, cfg.GetHiddenLabels()) {
				filteredOnLabels++
			} else if g.filterOnAccounts(issue, cfg.GetBotAccounts()) {
				filteredOnAccounts++
			} else {
				issuesToReturn = append(issuesToReturn, issue)
			}
		}
	}

	log.Logger().Debugf("Filtered %d/%d Issue(s)", filteredOnLabels, filteredOnAccounts)

	g.Issues = issuesToReturn
	g.FilteredLabels = filteredOnLabels
	g.FilteredBotAccounts = filteredOnAccounts

	return nil
}

func (g *GetIssues) GetIssuesForHost(host string, cfg config.Config, query string) ([]pr.Issue, error) {
	client, err := g.GithubClient()
	if err != nil {
		return nil, err
	}
	currentUser, err := GetCurrentUser(client, host)
	if err != nil {
		return nil, err
	}

	var queryString string
	if g.Raw != "" {
		queryString = g.Raw
	} else if g.Me {
		orgs, err := GetOrgsForUser(client, host)
		if err != nil {
			return nil, err
		}
		log.Logger().Debugf("User is a member of %s organisations", orgs)

		userQuery := "user:" + strings.Join(orgs, " user:")
		log.Logger().Debugf("User '%s'", userQuery)

		queryString = fmt.Sprintf("author:%s %s", currentUser, userQuery)
	} else if g.Review {
		queryString = fmt.Sprintf("review-requested:%s", currentUser)
	} else {
		queryString = strings.Join(cfg.GetReposToQuery(host), " ")
	}

	if cfg.GetMaxAgeOfPRs() != -1 {
		dateString := time.Now().AddDate(0, 0, -1*cfg.GetMaxAgeOfPRs()).Format("2006-01-02")
		queryString = queryString + " created:>" + dateString
	}

	queryToRun := fmt.Sprintf(query, queryString, cfg.GetMaxNumberOfPRs())
	log.Logger().Debugf("running query\n%s", queryToRun)

	data := IssueData{}
	err = client.GraphQL(host, queryToRun, nil, &data)
	if err != nil {
		return nil, err
	}

	return data.Search.Issues, nil
}

func (g *GetIssues) filterOnAccounts(pr pr.Issue, botAccounts []string) bool {
	for _, botAccount := range botAccounts {
		if pr.Author.Login == botAccount {
			return !g.ShowBots
		}
	}
	return false
}

func (g *GetIssues) filterOnLabels(pr pr.Issue, hiddenLabels []string) bool {
	for _, label := range hiddenLabels {
		if pr.HasLabel(label) {
			return !g.ShowHidden
		}
	}
	return false
}
