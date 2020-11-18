package rebasecmd

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/domain"
	"github.com/spf13/cobra"
)

type RebaseCmd struct {
	Cmd  *cobra.Command
	Args []string
}

func NewRebaseCmd() *cobra.Command {
	c := &RebaseCmd{}
	cmd := &cobra.Command{
		Use:   "rebase",
		Short: "Rebase the local clone",
		Long: "Performs a 'git fetch upstream master && git rebase upstream/master && git push origin master'.  " +
			"Uses the default_branch name determined from the GitHub API.",
		Example: "",
		Aliases: []string{"rb"},
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
	return cmd
}

func (c *RebaseCmd) Run() error {
	d := domain.NewRebase()

	err := d.Validate()
	if err != nil {
		return errors.Wrap(err, "validate failed")
	}
	err = d.Run()
	if err != nil {
		return errors.Wrap(err, "validate failed")
	}
	return nil
}
