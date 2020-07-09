package upgradecmd

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/spf13/cobra"
)

// UpgradeCmd defines the cmd.
type UpgradeCmd struct {
	Cmd  *cobra.Command
	Args []string
}

// NewUpgradeCmd defines a new cmd.
func NewUpgradeCmd() *cobra.Command {
	c := &UpgradeCmd{}
	cmd := &cobra.Command{
		Use:     "upgrade",
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

	cmd.AddCommand(NewUpgradeCliCmd())

	return cmd
}

// Run upgrade help.
func (c *UpgradeCmd) Run() error {
	return c.Cmd.Help()
}
