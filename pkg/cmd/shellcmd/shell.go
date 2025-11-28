package shellcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/plumming/dx/pkg/kube"
	"github.com/plumming/dx/pkg/util"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"k8s.io/client-go/tools/clientcmd"
)

type ShellCmd struct {
	Filter string

	Cmd  *cobra.Command
	Args []string
}

const (
	defaultPromptCommand = "dx prompt"
	defaultRcFile        = `
if [ -f /etc/bashrc ]; then
    source /etc/bashrc
fi
if [ -f ~/.bashrc ]; then
    source ~/.bashrc
fi
`

	zshRcFile = `
if [ -f /etc/zshrc ]; then
    source /etc/zshrc
fi
if [ -f ~/.zshrc ]; then
    source ~/.zshrc
fi
`
)

func NewShellCmd() *cobra.Command {
	c := &ShellCmd{}
	cmd := &cobra.Command{
		Use:     "shell",
		Short:   "Launch a new shell with a copy of the current kubernetes config",
		Long:    "",
		Example: "dx shell",
		Aliases: []string{"sh", "s"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
		Args: cobra.NoArgs,
	}
	cmd.Flags().StringVarP(&c.Filter, "filter", "f", "", "Filter the list of contexts to switch between using the given text")

	return cmd
}

func (c *ShellCmd) Run() error {
	config, _, err := kube.LoadConfig()
	if err != nil {
		return err
	}

	if config == nil || config.Contexts == nil || len(config.Contexts) == 0 {
		return fmt.Errorf("no kubernetes contexts available! try to create or connect to cluster")
	}

	contextNames := c.filterContexts(config)

	ctxName := ""

	defaultCtxName := config.CurrentContext
	pick, err := c.PickContext(contextNames, defaultCtxName)
	if err != nil {
		return err
	}
	ctxName = pick

	newConfig := *config
	newConfig.CurrentContext = ctxName

	c.cleanup()
	tmpDirName, err := os.MkdirTemp("/tmp", ".dx-shell-")
	if err != nil {
		return err
	}
	tmpConfigFileName := tmpDirName + "/config"
	err = clientcmd.WriteToFile(newConfig, tmpConfigFileName)
	if err != nil {
		return err
	}

	shell := filepath.Base(os.Getenv("SHELL"))
	prompt := c.createNewBashPrompt(os.Getenv("PS1"))
	rcFile := defaultRcFile + "\nexport PS1=" + prompt + "\nexport KUBECONFIG=\"" + tmpConfigFileName + "\"\n"
	tmpRCFileName := tmpDirName + "/.bashrc"

	if shell == "zsh" {
		prompt = c.createNewZshPrompt(os.Getenv("PS1"))
		rcFile = zshRcFile + "\nexport PS1=" + prompt + "\nexport KUBECONFIG=\"" + tmpConfigFileName + "\"\n"
		tmpRCFileName = tmpDirName + "/.zshrc"
	}
	err = os.WriteFile(tmpRCFileName, []byte(rcFile), 0600)
	if err != nil {
		return err
	}

	info := util.ColorInfo
	fmt.Printf("Creating a new shell using the Kubernetes context %s\n", info(ctxName))
	fmt.Printf("Shell RC file is %s\n\n", tmpRCFileName)
	fmt.Printf("All changes to the Kubernetes context like changing environment, namespace or context will be local to this shell\n")
	fmt.Printf("To return to the global context use the command: exit\n\n")

	e := exec.Command(shell, "-rcfile", tmpRCFileName, "-i")
	if shell == "zsh" {
		env := os.Environ()
		env = append(env, fmt.Sprintf("ZDOTDIR=%s", tmpDirName))
		e = exec.Command(shell, "-i")
		e.Env = env
	}

	return e.Run()
}

func (c *ShellCmd) filterContexts(config *api.Config) []string {
	var contextNames []string
	for k, v := range config.Contexts {
		if k != "" && v != nil {
			if c.Filter == "" || strings.Contains(k, c.Filter) {
				contextNames = append(contextNames, k)
			}
		}
	}

	sort.Strings(contextNames)
	return contextNames
}

func (c *ShellCmd) cleanup() {
	//clean old folders
	files, err := filepath.Glob("/tmp/.dx-shell-*")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.RemoveAll(f); err != nil {
			panic(err)
		}
	}
}

func (c *ShellCmd) PickContext(names []string, defaultValue string) (string, error) {
	if len(names) == 0 {
		return "", nil
	}
	if len(names) == 1 {
		return names[0], nil
	}
	name := ""
	prompt := &survey.Select{
		Message: "Change Kubernetes context:",
		Options: names,
		Default: defaultValue,
	}
	err := survey.AskOne(prompt, &name, nil)
	return name, err
}

func (c *ShellCmd) createNewBashPrompt(prompt string) string {
	if prompt == "" {
		return fmt.Sprintf("'[\\u@\\h \\W \\$(%s) ]\\$ '", defaultPromptCommand)
	}
	if strings.Contains(prompt, defaultPromptCommand) {
		return prompt
	}
	if prompt[0] == '"' {
		return prompt[0:1] + fmt.Sprintf("\\$(%s) ", defaultPromptCommand) + prompt[1:]
	}
	if prompt[0] == '\'' {
		return prompt[0:1] + fmt.Sprintf("$(%s) ", defaultPromptCommand) + prompt[1:]
	}
	return fmt.Sprintf("'$(%s) ", defaultPromptCommand) + prompt + "'"
}

func (c *ShellCmd) createNewZshPrompt(prompt string) string {
	if prompt == "" {
		return fmt.Sprintf("'[$(%s) %%n@%%m %%c]\\$ '", defaultPromptCommand)
	}
	if strings.Contains(prompt, defaultPromptCommand) {
		return prompt
	}
	if prompt[0] == '"' {
		return prompt[0:1] + fmt.Sprintf("$(%s) ", defaultPromptCommand) + prompt[1:]
	}
	if prompt[0] == '\'' {
		return prompt[0:1] + fmt.Sprintf("$(%s) ", defaultPromptCommand) + prompt[1:]
	}
	return fmt.Sprintf("'$(%s) ", defaultPromptCommand) + prompt + "'"
}
