package upgradecmd

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
	}

	cmd.AddCommand(NewUpgradeCliCmd())

	return cmd
}

// Run upgrade help.
func (c *UpgradeCmd) Run() error {
	return c.Cmd.Help()
}
