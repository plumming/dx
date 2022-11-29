package namespacecmd

import (
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/domain"
	"github.com/spf13/cobra"
)

type NamespaceCmd struct {
	Cmd  *cobra.Command
	Args []string
}

func NewNamespaceCmd() *cobra.Command {
	c := &NamespaceCmd{}
	cmd := &cobra.Command{
		Use:     "namespace",
		Short:   "View or change the current Kubernetes cluster namespace",
		Long:    "",
		Example: "",
		Aliases: []string{"ns"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
		Args: cobra.MaximumNArgs(1),
	}
	return cmd
}

func (c *NamespaceCmd) Run() error {
	d := domain.NewNamespace()

	if len(c.Args) == 1 {
		d.Namespace = c.Args[0]
	} else {
		err := d.Validate()
		if err != nil {
			return errors.Wrap(err, "validate failed")
		}
	}

	err := d.Run()
	if err != nil {
		return errors.Wrap(err, "run failed")
	}
	return nil
}
