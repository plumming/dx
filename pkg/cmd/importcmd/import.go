package importcmd

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
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
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				log.Logger().WithError(err).Fatal("unable to run command")
			}
		},
	}

	cmd.AddCommand(NewImportContextCmd())

	return cmd
}

// Run import help.
func (c *ImportCmd) Run() error {
	return c.Cmd.Help()
}
