package deletecmd

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
	}

	cmd.AddCommand(NewDeleteReposCmd())

	return cmd
}

// Run get help.
func (c *DeleteCmd) Run() error {
	return c.Cmd.Help()
}
