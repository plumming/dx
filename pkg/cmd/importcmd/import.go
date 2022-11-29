package importcmd

import (
	"github.com/spf13/cobra"
)

// ImportCmd defines parent import.
type ImportCmd struct {
	Cmd  *cobra.Command
	Args []string
}

// NewImportCmd creates import cmd.
func NewImportCmd() *cobra.Command {
	c := &ImportCmd{}
	cmd := &cobra.Command{
		Use:     "import",
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
	}

	cmd.AddCommand(NewImportContextCmd())

	return cmd
}

// Run import help.
func (c *ImportCmd) Run() error {
	return c.Cmd.Help()
}
