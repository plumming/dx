package getcmd

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/spf13/cobra"
)

// GetCmd defines parent get.
type GetCmd struct {
	Cmd  *cobra.Command
	Args []string
}

// NewGetCmd creates get cmd.
func NewGetCmd() *cobra.Command {
	c := &GetCmd{}
	cmd := &cobra.Command{
		Use:     "get",
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

	cmd.AddCommand(NewGetPrsCmd())

	return cmd
}

// Run get help.
func (c *GetCmd) Run() error {
	return c.Cmd.Help()
}
