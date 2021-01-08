package getcmd

import (
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/domain"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/table"
	"github.com/spf13/cobra"
)

type GetReposCmd struct {
	cmd.CommonCmd
	Org  string
	User string
	Cmd  *cobra.Command
	Args []string
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

	cmd.Flags().StringVarP(&c.Org, "org", "o", "",
		"Organization to query")
	cmd.Flags().StringVarP(&c.User, "user", "u", "",
		"User to query")

	c.AddOptions(cmd)

	return cmd
}

func (c *GetReposCmd) Run() error {
	if c.Org == "" && c.User == "" {
		return errors.New("need to supply an --org or a --user to query")
	}

	d := domain.Repo{}

	table := table.NewTable(c.Cmd.OutOrStdout())
	table.AddRow("Repository")

	var err error
	var repos []domain.RepoInfo

	if c.Org != "" {
		repos, err = d.ListRepositoriesForOrg(c.Org)
		if err != nil {
			return err
		}
	} else {
		repos, err = d.ListRepositoriesForUser(c.User)
		if err != nil {
			return err
		}
	}

	for _, r := range repos {
		table.AddRow(r.Name)
	}

	table.Render()
	return nil
}
