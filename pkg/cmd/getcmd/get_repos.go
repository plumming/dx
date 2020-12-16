package getcmd

import (
	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/domain"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/table"
	"github.com/spf13/cobra"
)

type GetReposCmd struct {
	cmd.CommonCmd
	Filter string
	Cmd    *cobra.Command
	Args   []string
}

func NewGetReposCmd() *cobra.Command {
	c := &GetReposCmd{}
	cmd := &cobra.Command{
		Use:     "repos",
		Short:   "Lists your repositories",
		Long:    "",
		Example: "",
		Aliases: []string{"repositories", "repository"},
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

	return cmd
}

func (c *GetReposCmd) Run() error {
	d := domain.Repo{}

	table := table.NewTable(c.Cmd.OutOrStdout())
	table.AddRow("Repository")

	repos, err := d.ListRepositories("garethjevans-test")
	if err != nil {
		return err
	}

	for _, r := range repos {
		table.AddRow(r.Name)
	}

	table.Render()
	return nil
}
