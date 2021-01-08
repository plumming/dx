package deletecmd

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/spf13/cobra"
)

// DeleteCmd defines parent get.
type DeleteCmd struct {
	Cmd  *cobra.Command
	Args []string
}

// NewDeleteCmd creates delete cmd.
func NewDeleteCmd() *cobra.Command {
	c := &DeleteCmd{}
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "",
		Long:    "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				log.Logger().WithError(err).Fatal("unable to run command")
			}
		},
	}

	cmd.AddCommand(NewDeleteReposCmd())

	return cmd
}

// Run get help.
func (c *DeleteCmd) Run() error {
	return c.Cmd.Help()
}
