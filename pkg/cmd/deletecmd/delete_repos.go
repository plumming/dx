package deletecmd

import (
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/domain"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/spf13/cobra"
)

type DeleteReposCmd struct {
	cmd.CommonCmd
	Org string
	Cmd            *cobra.Command
	Args           []string
}

func NewDeleteReposCmd() *cobra.Command {
	c := &DeleteReposCmd{}
	cmd := &cobra.Command{
		Use:     "repos",
		Short:   "Delete your repositories",
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

	cmd.Flags().StringVarP(&c.Org, "org", "", "",
		"Organization to query")

	c.AddOptions(cmd)

	return cmd
}

func (c *DeleteReposCmd) Run() error {
	if c.Org == "" {
		return errors.New("Need to select an --org to query")
	}

	d := domain.Repo{}
	err := d.DeleteRepositories(c.Org)
	if err != nil {
		return err
	}

	return nil
}
