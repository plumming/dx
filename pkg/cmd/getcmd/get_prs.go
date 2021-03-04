package getcmd

import (
	"fmt"

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
	ShowDependabot bool
	ShowOnHold     bool
	Retrigger      bool
	Review         bool
	Quiet          bool
	Me             bool
	Cmd            *cobra.Command
	Args           []string
}

func NewGetPrsCmd() *cobra.Command {
	c := &GetPrsCmd{}
	cmd := &cobra.Command{
		Use:     "prs",
		Short:   "Gets your open prs",
		Long:    "",
		Example: "",
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

	cmd.Flags().BoolVarP(&c.ShowDependabot, "show-dependabot", "", false,
		"Show dependabot PRs (default: false)")
	cmd.Flags().BoolVarP(&c.ShowOnHold, "show-on-hold", "", false,
		"Show On Hold PRs (default: false)")
	cmd.Flags().BoolVarP(&c.Retrigger, "retrigger", "", false,
		"Retrigger failed PRs")
	cmd.Flags().BoolVarP(&c.Review, "review", "", false,
		"Show PRs that are ready for review")
	cmd.Flags().BoolVarP(&c.Quiet, "quiet", "", false,
		"Hide the column headings")
	cmd.Flags().BoolVarP(&c.Me, "me", "m", false,
		"Show all PRs that are created by the author")

	return cmd
}

func (c *GetPrsCmd) Run() error {
	d := domain.NewGetPrs()
	d.ShowOnHold = c.ShowOnHold
	d.ShowDependabot = c.ShowDependabot
	d.Me = c.Me
	d.Review = c.Review

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

	if !c.Retrigger {
		return nil
	}

	err = d.Retrigger()
	if err != nil {
		return errors.Wrap(err, "retrigger failed")
	}

	return nil
}
