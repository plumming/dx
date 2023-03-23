package getcmd

import (
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/domain"
	"github.com/plumming/dx/pkg/util"

	"github.com/plumming/dx/pkg/table"
	"github.com/spf13/cobra"
)

type GetSecurityConfigCmd struct {
	cmd.CommonCmd
	Quiet bool
	Cmd   *cobra.Command
	Args  []string
}

func NewGetSecurityConfigCmd() *cobra.Command {
	c := &GetSecurityConfigCmd{}
	cmd := &cobra.Command{
		Use:   "security-config",
		Short: "Displays the security config for configured repos",
		Long:  "",
		Example: `Gets the security config for each repository, if you have permission:

  dx get security-config
`,
		Aliases: []string{"securityconfig", "sc"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
		Args: cobra.NoArgs,
	}

	c.AddOptions(cmd)

	cmd.Flags().BoolVarP(&c.Quiet, "quiet", "", false,
		"Hide the column headings")

	return cmd
}

func (c *GetSecurityConfigCmd) Run() error {
	d := domain.NewGetSecurityConfig()

	err := d.Validate()
	if err != nil {
		return errors.Wrap(err, "validate failed")
	}

	err = d.Run()
	if err != nil {
		return errors.Wrap(err, "run failed")
	}

	table := table.NewTable(c.Cmd.OutOrStdout())

	if !c.Quiet {
		table.AddRow(
			"Repository",
			"Vulnerability Alerts",
		)
	}

	for _, c := range d.Config {
		table.AddRow(
			c.NameWithOwner,
			colourBool(c.HasVulnerabilityAlertsEnabled),
		)
	}

	table.Render()

	return nil
}

func colourBool(b bool) string {
	if b {
		return util.ColorInfo("Yes")
	}
	return util.ColorError("No")
}
