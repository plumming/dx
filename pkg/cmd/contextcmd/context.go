package contextcmd

import (
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/domain"
	"github.com/spf13/cobra"
)

type ContextCmd struct {
	Cmd  *cobra.Command
	Args []string
}

func NewContextCmd() *cobra.Command {
	c := &ContextCmd{}
	cmd := &cobra.Command{
		Use:     "context",
		Short:   "Change the current Kubernetes context",
		Long:    "",
		Example: "",
		Aliases: []string{"ctx", "c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
		Args: cobra.MaximumNArgs(1),
	}
	return cmd
}

func (c *ContextCmd) Run() error {
	d := domain.NewContext()

	if len(c.Args) == 1 {
		d.Context = c.Args[0]
	}

	err := d.Validate()
	if err != nil {
		return errors.Wrap(err, "validate failed")
	}

	err = d.Run()
	if err != nil {
		return errors.Wrap(err, "run failed")
	}
	return nil
}
