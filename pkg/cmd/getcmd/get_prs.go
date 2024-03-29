package getcmd

import (
	"fmt"
	"strings"

	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/domain"

	"github.com/pkg/errors"

	"github.com/plumming/dx/pkg/util"

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
	Copy       bool
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
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

	cmd.Flags().BoolVarP(&c.Copy, "copy", "c", false,
		"Output is copy and pasteable")

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

	if c.Copy {
		if !c.Quiet {
			table.AddRow(
				"PR#",
				"Author",
				"Title",
			)
		}

		for _, pr := range d.PullRequests {
			table.AddRow(
				util.ColorInfo(pr.URL),
				pr.Author.Login,
				pr.ColoredTitle(),
			)
		}
	} else {
		if !c.Quiet {
			table.AddRow(
				"PR#",
				"Author",
				"Title",
				"Age",
				"Review",
				"Labels",
				"Mergeable",
				"Comments",
			)
		}

		for _, p := range d.PullRequests {
			pr := p
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
				util.SafeIfAboveZero(pr.Comments.TotalCount),
			)
		}
	}

	table.Render()

	if len(d.PullRequests) > 0 {
		fmt.Printf("\nDisplaying %d PRs\n", len(d.PullRequests))
	}

	if (d.FilteredBotAccounts + d.FilteredLabels) > 0 {
		var flags []string
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
