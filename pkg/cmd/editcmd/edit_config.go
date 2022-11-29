package editcmd

import (
	"os"
	"os/exec"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/config"
	"github.com/spf13/cobra"
)

type EditConfigCmd struct {
	ShowDependabot bool
	ShowOnHold     bool
	Retrigger      bool
	Review         bool
	Cmd            *cobra.Command
	Args           []string
}

func NewEditConfigCmd() *cobra.Command {
	c := &EditConfigCmd{}
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Edit the configuration for dx",
		Long:    "",
		Example: "",
		Aliases: []string{"configuration"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
		Args: cobra.NoArgs,
	}

	return cmd
}

func (c *EditConfigCmd) Run() error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	f, err := os.CreateTemp("", "dx-*.yaml")
	if err != nil {
		log.Logger().Fatal(err)
	}
	log.Logger().Infof("Tmp File %s", f.Name())
	err = f.Close()
	if err != nil {
		log.Logger().Fatal(err)
	}

	configuration, err := config.LoadFromDefaultLocation()
	if err != nil {
		log.Logger().Fatal(err)
	}

	err = configuration.SaveToFile(f.Name())
	if err != nil {
		log.Logger().Fatal(err)
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Logger().Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Logger().Errorf("Error while editing. Error: %v", err)
	} else {
		log.Logger().Printf("Successfully edited.")
	}

	configuration, err = config.LoadFromFile(f.Name())
	if err != nil {
		log.Logger().Fatal(err)
	}

	err = configuration.SaveToDefaultLocation()
	if err != nil {
		log.Logger().Fatal(err)
	}

	return nil
}
