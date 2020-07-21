package editcmd

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/spf13/cobra"
)

// EditCmd defines parent get.
type EditCmd struct {
	Cmd  *cobra.Command
	Args []string
}

// NewGetCmd creates get cmd.
func NewEditCmd() *cobra.Command {
	c := &EditCmd{}
	cmd := &cobra.Command{
		Use:     "edit",
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

	cmd.AddCommand(NewEditConfigCmd())

	return cmd
}

// Run get help.
func (c *EditCmd) Run() error {
	return c.Cmd.Help()
}
