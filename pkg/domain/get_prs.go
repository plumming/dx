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
	getPrsQuery = `{
	search(query: "is:pr is:open %s", type: ISSUE, first: %d) {
      nodes {
      ... on PullRequest {
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
        mergeable
        labels(first: 10) {
          nodes {
            name
          }
        }
        reviewDecision
        commits(last: 1){
          nodes{
            commit{
              statusCheckRollup {
                state
                contexts(last:100) {
                  totalCount
                  nodes {
                    ...on StatusContext {
                      state
                      context
                      description
                    }
                    ...on CheckRun {
                      conclusion
                      name
                      title
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`
)

// GetPrs defines get pull request response.
type GetPrs struct {
	cmd.CommonOptions
	ShowBots            bool
	ShowHidden          bool
	Me                  bool
	Review              bool
	Raw                 string
	PullRequests        []pr.PullRequest
	FilteredLabels      int
	FilteredBotAccounts int
}

// PrData.
type PrData struct {
	Search PrSearch `json:"search"`
}

// PrSearch.
type PrSearch struct {
	PullRequests []pr.PullRequest `json:"nodes"`
}

// NewGetPrs.
func NewGetPrs() *GetPrs {
	g := &GetPrs{}
	return g
}

// Validate input.
func (g *GetPrs) Validate() error {
	return nil
}

// Run the cmd.
func (g *GetPrs) Run() error {
	cfg, err := g.DxConfig()
	if err != nil {
		return err
	}

	var pulls []pr.PullRequest

	for _, host := range cfg.GetConfiguredServers() {
		p, err := g.getPrsForHost(host, cfg, getPrsQuery)
		if err != nil {
			return err
		}
		pulls = append(pulls, p...)
	}

	sort.Sort(pr.ByPullsString(pulls))

	var pullsToReturn []pr.PullRequest

	filteredOnLabels := 0
	filteredOnAccounts := 0

	for _, pullRequest := range pulls {
		if pullRequest.Display() {
			if g.filterOnLabels(pullRequest, cfg.GetHiddenLabels()) {
				filteredOnLabels++
			} else if g.filterOnAccounts(pullRequest, cfg.GetBotAccounts()) {
				filteredOnAccounts++
			} else {
				pullsToReturn = append(pullsToReturn, pullRequest)
			}
		}
	}

	log.Logger().Debugf("Filtered %d/%d PR(s)", filteredOnLabels, filteredOnAccounts)

	g.PullRequests = pullsToReturn
	g.FilteredLabels = filteredOnLabels
	g.FilteredBotAccounts = filteredOnAccounts

	return nil
}

func (g *GetPrs) getPrsForHost(host string, cfg config.Config, query string) ([]pr.PullRequest, error) {
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

	data := PrData{}
	err = client.GraphQL(host, queryToRun, nil, &data)
	if err != nil {
		return nil, err
	}

	return data.Search.PullRequests, nil
}

func (g *GetPrs) filterOnAccounts(pr pr.PullRequest, botAccounts []string) bool {
	for _, botAccount := range botAccounts {
		if pr.Author.Login == botAccount {
			return !g.ShowBots
		}
	}
	return false
}

func (g *GetPrs) filterOnLabels(pr pr.PullRequest, hiddenLabels []string) bool {
	for _, label := range hiddenLabels {
		if pr.HasLabel(label) {
			return !g.ShowHidden
		}
	}
	return false
}

// Retrigger failed prs.
func (g *GetPrs) Retrigger() error {
	client, err := g.GithubClient()
	if err != nil {
		return err
	}

	log.Logger().Infof("Retriggering Failed & Non Conflicting PRs...")

	for _, pr := range g.PullRequests {
		if pr.ContextsString() == "FAILURE" && pr.Mergeable == "MERGEABLE" && pr.HasLabel("updatebot") {
			failedContexts := pr.FailedContexts()
			for _, f := range failedContexts {
				testCommand := fmt.Sprintf("/test %s", f.Context)
				if f.Context == "pr-build" {
					testCommand = "/test this"
				}
				log.Logger().Infof("%s with '%s'", pr.URL, testCommand)

				url := fmt.Sprintf("repos/%s/issues/%d/comments", pr.Repository.NameWithOwner, pr.Number)
				body := fmt.Sprintf("{ \"body\": \"%s\" }", testCommand)

				err := client.REST("github.com", "POST", url, strings.NewReader(body), nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
