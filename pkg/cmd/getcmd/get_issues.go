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

type GetIssuesCmd struct {
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

func NewGetIssuesCmd() *cobra.Command {
	c := &GetIssuesCmd{}
	cmd := &cobra.Command{
		Use:   "issues",
		Short: "Gets your open issues",
		Long:  "",
		Example: `Get a list of open issues:

  dx get issues

Get a list of your issues:

  dx get issues --me

Get a list of issues requiring review:

  dx get issues --review

Get a list of issues with a custom query:

  dx get issues --raw is:private

`,
		Aliases: []string{"issues", "is"},
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

	cmd.Flags().StringVarP(&c.Raw, "raw", "", "",
		"Additional raw search parameters to use when querying")

	return cmd
}

func (c *GetIssuesCmd) Run() error {
	d := domain.NewGetIssues()
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
		fmt.Println(c.Filter(d.Issues))
		return nil
	}

	table := table.NewTable(c.Cmd.OutOrStdout())

	pullURL := ""
	if !c.Quiet {
		table.AddRow(
			"Issue#",
			"Author",
			"Title",
			"Age",
			"Labels",
			"Comments",
		)
	}

	for _, i := range d.Issues {
		issue := i
		if pullURL != issue.IssueString() {
			table.AddRow(fmt.Sprintf("# %s", util.ColorAnswer(issue.IssueString())))
			pullURL = issue.IssueString()
		}
		table.AddRow(
			util.ColorInfo(fmt.Sprintf("#%d", issue.Number)),
			issue.Author.Login,
			issue.ColoredTitle(),
			util.SafeTime(&issue.CreatedAt),
			issue.LabelsString(),
			util.SafeIfAboveZero(issue.Comments.TotalCount),
		)
	}

	table.Render()

	if len(d.Issues) > 0 {
		fmt.Printf("\nDisplaying %d Issue(s)\n", len(d.Issues))
	}

	if (d.FilteredBotAccounts + d.FilteredLabels) > 0 {
		flags := []string{}
		if d.FilteredBotAccounts > 0 {
			flags = append(flags, "--show-bots")
		}
		if d.FilteredLabels > 0 {
			flags = append(flags, "--show-hidden")
		}
		fmt.Printf("\nFiltered %d Issues, use %s to view them\n", (d.FilteredBotAccounts + d.FilteredLabels), strings.Join(flags, ", "))
	}

	return nil
}
