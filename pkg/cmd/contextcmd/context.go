package contextcmd

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/pkg/errors"
	"github.com/plumming/chilly/pkg/domain"
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
		Short:   "View or change the current Kubernetes context",
		Long:    "",
		Example: "",
		Aliases: []string{"ctx"},
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

func (c *ContextCmd) Run() error {
	d := domain.NewContext()

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
