package deletecmd

import (
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/domain"

	"github.com/spf13/cobra"
)

type DeleteReposCmd struct {
	cmd.CommonCmd
	Org  string
	User string
	Cmd  *cobra.Command
	Args []string
}

func NewDeleteReposCmd() *cobra.Command {
	c := &DeleteReposCmd{}
	cmd := &cobra.Command{
		Use:     "repos",
		Short:   "Delete your repositories",
		Long:    "",
		Example: "",
		Aliases: []string{"repositories", "repository"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
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

func (c *DeleteReposCmd) Run() error {
	if c.Org == "" && c.User == "" {
		return errors.New("need to supply an --org or a --user to query")
	}

	d := domain.Repo{}

	if c.Org != "" {
		err := d.DeleteRepositoriesFromOrg("github.com", c.Org)
		if err != nil {
			return err
		}
	} else {
		err := d.DeleteRepositoriesFromUser("github.com", c.User)
		if err != nil {
			return err
		}
	}

	return nil
}
