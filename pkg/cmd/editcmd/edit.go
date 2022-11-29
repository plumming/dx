package editcmd

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
	}

	cmd.AddCommand(NewEditConfigCmd())

	return cmd
}

// Run get help.
func (c *EditCmd) Run() error {
	return c.Cmd.Help()
}
