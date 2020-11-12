package editcmd

import (
	"os"
	"os/exec"
	"path"

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

func (c *EditConfigCmd) Run() error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	fpath := path.Join(os.TempDir(), "thetemporaryfile.txt")
	log.Logger().Infof("Tmp File %s", fpath)
	f, err := os.Create(fpath)
	if err != nil {
		log.Logger().Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Logger().Fatal(err)
	}

	configuration, err := config.LoadFromDefaultLocation()
	if err != nil {
		log.Logger().Fatal(err)
	}

	err = configuration.SaveToFile(fpath)
	if err != nil {
		log.Logger().Fatal(err)
	}

	cmd := exec.Command(editor, fpath)
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

	configuration, err = config.LoadFromFile(fpath)
	if err != nil {
		log.Logger().Fatal(err)
	}

	err = configuration.SaveToDefaultLocation()
	if err != nil {
		log.Logger().Fatal(err)
	}

	return nil
}
