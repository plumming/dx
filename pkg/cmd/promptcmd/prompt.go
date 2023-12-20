package promptcmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/plumming/dx/pkg/kube"
	"github.com/plumming/dx/pkg/util"
	"github.com/spf13/cobra"
)

type PromptCmd struct {
	NoLabel        bool
	ShowIcon       bool
	Prefix         string
	Label          string
	Separator      string
	Divider        string
	Suffix         string
	LabelColor     []string
	NamespaceColor []string
	ContextColor   []string

	Cmd  *cobra.Command
	Args []string
}

func NewPromptCmd() *cobra.Command {
	c := &PromptCmd{}
	cmd := &cobra.Command{
		Use:   "prompt",
		Short: "Configure a new command prompt",
		Long:  "",
		Example: `
		# Generate the current prompt
		dx prompt

		# Enable the prompt for bash
		PS1="[\u@\h \W \$(dx prompt)]\$ "

		# Enable the prompt for zsh
		PROMPT='$(dx prompt)'$PROMPT
	`,
		Aliases: []string{"p"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Cmd = cmd
			c.Args = args
			return c.Run()
		},
		Args: cobra.NoArgs,
	}

	cmd.Flags().StringVarP(&c.Prefix, "prefix", "p", "", "The prefix text for the prompt")
	cmd.Flags().StringVarP(&c.Label, "label", "l", "k8s", "The label for the prompt")
	cmd.Flags().StringVarP(&c.Separator, "separator", "s", ":", "The separator between the label and the rest of the prompt")
	cmd.Flags().StringVarP(&c.Divider, "divider", "d", ":", "The divider between the team and environment for the prompt")
	cmd.Flags().StringVarP(&c.Suffix, "suffix", "x", ">", "The suffix text for the prompt")

	cmd.Flags().StringArrayVarP(&c.LabelColor, optionLabelColor, "", []string{"blue"}, "The color for the label")
	cmd.Flags().StringArrayVarP(&c.NamespaceColor, optionNamespaceColor, "", []string{"green"}, "The color for the namespace")
	cmd.Flags().StringArrayVarP(&c.ContextColor, optionContextColor, "", []string{"cyan"}, "The color for the Kubernetes context")

	cmd.Flags().BoolVarP(&c.NoLabel, "no-label", "", false, "Disables the use of the label in the prompt")
	cmd.Flags().BoolVarP(&c.ShowIcon, "icon", "i", false, "Uses an icon for the label in the prompt")

	return cmd
}

const (
	optionLabelColor     = "label-color"
	optionNamespaceColor = "namespace-color"
	optionContextColor   = "context-color"
)

// Run implements this command.
func (o *PromptCmd) Run() error {
	config, _, err := kube.LoadConfig()
	if err != nil {
		return err
	}

	context := config.CurrentContext
	namespace := kube.CurrentNamespace(config)

	// enable color
	color.NoColor = os.Getenv("TERM") == "dumb"

	label := o.Label
	separator := o.Separator
	divider := o.Divider
	prefix := o.Prefix
	suffix := o.Suffix

	labelColor, err := util.GetColor(optionLabelColor, o.LabelColor)
	if err != nil {
		return err
	}
	nsColor, err := util.GetColor(optionLabelColor, o.NamespaceColor)
	if err != nil {
		return err
	}
	ctxColor, err := util.GetColor(optionLabelColor, o.ContextColor)
	if err != nil {
		return err
	}
	if o.NoLabel {
		label = ""
		separator = ""
	} else {
		if o.ShowIcon {
			label = "☸️  "
			label = labelColor.Sprint(label)
		} else {
			label = labelColor.Sprint(label)
		}
	}
	if namespace == "" {
		divider = ""
	} else {
		namespace = nsColor.Sprint(namespace)
	}
	context = ctxColor.Sprint(context)
	fmt.Printf("%s\n", strings.Join([]string{prefix, label, separator, namespace, divider, context, suffix}, ""))
	return nil
}
