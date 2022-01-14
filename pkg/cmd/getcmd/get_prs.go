package getcmd

import (
	"fmt"
	"strings"

	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/domain"

	"github.com/pkg/errors"

	"github.com/plumming/dx/pkg/util"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/table"
	"github.com/spf13/cobra"
)

type GetPrsCmd struct {
	cmd.CommonCmd
	ShowBots   bool
	ShowHidden bool
	Review     bool
	Quiet      bool
	Me         bool
	Raw        string
	Cmd        *cobra.Command
	Args       []string
}

func NewGetPrsCmd() *cobra.Command {
	c := &GetPrsCmd{}
	cmd := &cobra.Command{
		Use:   "prs",
		Short: "Gets your open prs",
		Long:  "",
		Example: `Get a list of open PRs:

  dx get prs

Get a list of your PRs:

  dx get prs --me

Get a list of PRs requiring review:

  dx get prs --review

Get a list of PRs with a custom query:

  dx get prs --raw is:private

`,
		Aliases: []string{"pr", "pulls", "pull-requests"},
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				log.Logger().Fatalf("unable to run command: %s", err)
			}
		},
		Args: cobra.NoArgs,
	}

	c.AddOptions(cmd)

	cmd.Flags().BoolVarP(&c.ShowBots, "show-bots", "", false,
		"Show bot account PRs (default: false)")
	cmd.Flags().BoolVarP(&c.ShowHidden, "show-hidden", "", false,
		"Show PRs that are filtered by hidden labels (default: false)")
	cmd.Flags().BoolVarP(&c.Review, "review", "", false,
		"Show PRs that are ready for review")
	cmd.Flags().BoolVarP(&c.Quiet, "quiet", "", false,
		"Hide the column headings")
	cmd.Flags().BoolVarP(&c.Me, "me", "m", false,
		"Show all PRs that are created by the author")

	cmd.Flags().StringVarP(&c.Raw, "raw", "", "",
		"Additional raw search parameters to use when querying")

	return cmd
}

func (c *GetPrsCmd) Run() error {
	d := domain.NewGetPrs()
	d.ShowHidden = c.ShowHidden
	d.ShowBots = c.ShowBots
	d.Me = c.Me
	d.Review = c.Review
	d.Raw = c.Raw

	err := d.Validate()
	if err != nil {
		return errors.Wrap(err, "validate failed")
	}

	err = d.Run()
	if err != nil {
		return errors.Wrap(err, "run failed")
	}

	if c.Query != "" {
		fmt.Println(c.Filter(d.PullRequests))
		return nil
	}

	table := table.NewTable(c.Cmd.OutOrStdout())

	pullURL := ""
	if !c.Quiet {
		table.AddRow(
			"PR#",
			"Author",
			"Title",
			"Age",
			"Review",
			"Labels",
			"Mergeable",
		)
	}

	for _, pr := range d.PullRequests {
		if pullURL != pr.PullsString() {
			table.AddRow(fmt.Sprintf("# %s", util.ColorAnswer(pr.PullsString())))
			pullURL = pr.PullsString()
		}
		table.AddRow(
			util.ColorInfo(fmt.Sprintf("#%d", pr.Number)),
			pr.Author.Login,
			pr.ColoredTitle(),
			util.SafeTime(&pr.CreatedAt),
			pr.ColoredReviewDecision(),
			pr.LabelsString(),
			pr.MergeableString(),
		)
	}

	table.Render()

	if len(d.PullRequests) > 0 {
		fmt.Printf("\nDisplaying %d PRs\n", len(d.PullRequests))
	}
	
	if (d.FilteredBotAccounts + d.FilteredLabels) > 0 {
		flags := []string{}
		if d.FilteredBotAccounts > 0 {
			flags = append(flags, "--show-bots")
		}
		if d.FilteredLabels > 0 {
			flags = append(flags, "--show-hidden")
		}
		fmt.Printf("\nFiltered %d PRs, use %s to view them\n", (d.FilteredBotAccounts + d.FilteredLabels), strings.Join(flags, ", "))
	}

	return nil
}
