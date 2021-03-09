package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/plumming/dx/pkg/pr"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/cmd"
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

// Data.
type Data struct {
	Search Search `json:"search"`
}

// Search.
type Search struct {
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
	client, err := g.GithubClient()
	if err != nil {
		return err
	}

	query := `{
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

	data := Data{}
	cfg, err := g.Config()
	if err != nil {
		return err
	}

	client, err = g.GithubClient()
	if err != nil {
		return err
	}
	currentUser, err := GetCurrentUser(client)
	if err != nil {
		return err
	}

	var queryString string
	if g.Raw != "" {
		queryString = g.Raw
	} else if g.Me {
		queryString = fmt.Sprintf("author:%s", currentUser)
	} else if g.Review {
		queryString = fmt.Sprintf("review-requested:%s", currentUser)
	} else {
		queryString = strings.Join(cfg.ReposToQuery(), " ")
	}

	if cfg.MaxAge != -1 {
		dateString := time.Now().AddDate(0, 0, -1*cfg.MaxAge).Format("2006-01-02")
		queryString = queryString + " created:>" + dateString
	}

	queryToRun := fmt.Sprintf(query, queryString, cfg.MaxNumberOfPRs)
	log.Logger().Debugf("running query\n%s", queryToRun)

	err = client.GraphQL(queryToRun, nil, &data)
	if err != nil {
		return err
	}

	pulls := data.Search.PullRequests
	sort.Sort(pr.ByPullsString(pulls))

	pullsToReturn := []pr.PullRequest{}

	filteredOnLabels := 0
	filteredOnAccounts := 0

	for _, pr := range pulls {
		if pr.Display() {
			if g.filterOnLabels(pr, cfg.HiddenLabels) {
				filteredOnLabels++
			} else if g.filterOnAccounts(pr, cfg.BotAccounts) {
				filteredOnAccounts++
			} else {
				pullsToReturn = append(pullsToReturn, pr)
			}
		}
	}

	log.Logger().Debugf("Filtered %d/%d PR(s)", filteredOnLabels, filteredOnAccounts)

	g.PullRequests = pullsToReturn
	g.FilteredLabels = filteredOnLabels
	g.FilteredBotAccounts = filteredOnAccounts

	return nil
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

				err := client.REST("POST", url, strings.NewReader(body), nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
