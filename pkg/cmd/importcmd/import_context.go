package importcmd

import (
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/domain"
	"github.com/spf13/cobra"
)

type ImportContextCmd struct {
	cmd.CommonCmd
	Path string
	Cmd  *cobra.Command
	Args []string
}

func NewImportContextCmd() *cobra.Command {
	c := &ImportContextCmd{}
	cmd := &cobra.Command{
		Use:     "context",
		Short:   "Import a kubernetes context",
		Long:    "",
		Example: "",
		Aliases: []string{"ctx", "c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
		Args: cobra.NoArgs,
	}

	cmd.Flags().StringVarP(&c.Path, "path", "f", "",
		"Path to the context file to import")
	err := cmd.MarkFlagRequired("path")
	if err != nil {
		panic(err)
	}

	return cmd
}

func (c *ImportContextCmd) Run() error {
	d := domain.NewImportContext(c.Path)

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
