package namespacecmd

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
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
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				log.Logger().Fatalf("unable to run command: %s", err)
			}
		},
		Args: cobra.NoArgs,
	}
	return cmd
}

func (c *NamespaceCmd) Run() error {
	d := domain.NewNamespace()

	err := d.Validate()
	if err != nil {
		return errors.Wrap(err, "validate failed")
	}
	err = d.Run()
	if err != nil {
		return errors.Wrap(err, "validate failed")
	}
	return nil
}
