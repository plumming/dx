package rebasecmd

import (
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/domain"
	"github.com/spf13/cobra"
)

type RebaseCmd struct {
	Cmd            *cobra.Command
	Args           []string
	ForceWithLease bool
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
		Args: cobra.NoArgs,
	}

	cmd.Flags().BoolVarP(&c.ForceWithLease, "force-with-lease", "f", false,
		"Use `--force-with-lease` when pushing back up to the forked repo")

	return cmd
}

func (c *RebaseCmd) Run() error {
	d := domain.NewRebase(c.ForceWithLease)

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
