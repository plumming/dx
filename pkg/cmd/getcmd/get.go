package getcmd

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
	}

	cmd.AddCommand(NewGetPrsCmd())
	cmd.AddCommand(NewGetReposCmd())
	cmd.AddCommand(NewGetIssuesCmd())
	cmd.AddCommand(NewGetVulnerabilityAlertsCmd())
	cmd.AddCommand(NewGetSecurityConfigCmd())

	return cmd
}

// Run get help.
func (c *GetCmd) Run() error {
	return c.Cmd.Help()
}
