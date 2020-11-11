package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/plumming/dx/pkg/cmd/contextcmd"

	"github.com/plumming/dx/pkg/cmd/editcmd"

	"github.com/plumming/dx/pkg/deprecation"

	"github.com/spf13/cobra/doc"

	"github.com/plumming/dx/pkg/cmd/upgradecmd"

	"github.com/plumming/dx/pkg/api"
	"github.com/plumming/dx/pkg/update"
	"github.com/plumming/dx/pkg/util"
	"github.com/plumming/dx/pkg/version"

	"github.com/plumming/dx/pkg/cmd/getcmd"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Version is dynamically set by the toolchain or overridden by the Makefile.
var Version = version.Version

// BuildDate is dynamically set at build time in the Makefile.
var BuildDate = version.BuildDate

var versionOutput = ""
var updaterEnabled = "plumming/dx"

func init() {
	log.Logger()

	if strings.Contains(Version, "dev") {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
	Version = strings.TrimPrefix(Version, "v")
	if BuildDate == "" {
		RootCmd.Version = Version
	} else {
		RootCmd.Version = fmt.Sprintf("%s (%s)", Version, BuildDate)
	}
	versionOutput = fmt.Sprintf("dx version %s\n%s\n", RootCmd.Version, changelogURL(Version))
	RootCmd.AddCommand(versionCmd)
	RootCmd.SetVersionTemplate(versionOutput)

	RootCmd.AddCommand(docsCmd)

	RootCmd.PersistentFlags().Bool("help", false, "Show help for command")
	RootCmd.Flags().Bool("version", false, "Show version")

	RootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if err == pflag.ErrHelp {
			return err
		}
		return &FlagError{Err: err}
	})

	RootCmd.AddCommand(getcmd.NewGetCmd())
	RootCmd.AddCommand(editcmd.NewEditCmd())
	RootCmd.AddCommand(upgradecmd.NewUpgradeCmd())
	RootCmd.AddCommand(contextcmd.NewContextCmd())

	c := completionCmd
	c.Flags().StringP("shell", "s", "bash", "Shell type: {bash|zsh|fish|powershell}")
	RootCmd.AddCommand(c)
}

// FlagError is the kind of error raised in flag processing.
type FlagError struct {
	Err error
}

// Error.
func (fe FlagError) Error() string {
	return fe.Err.Error()
}

// Unwrap FlagError.
func (fe FlagError) Unwrap() error {
	return fe.Err
}

// RootCmd is the entry point of command-line execution.
var RootCmd = &cobra.Command{
	Use:   "dx",
	Short: "Plumming dx",
	Long:  `Have you got the chillies.`,

	SilenceErrors: false,
	SilenceUsage:  false,
}

var versionCmd = &cobra.Command{
	Use:    "version",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(versionOutput)
	},
}

var linkHandler = func(in string) string {
	return in
}

var alphaOnlyTemplate = `## Alpha Only
Please use '%s' when using the beta environment
`

var betaOnlyTemplate = `## Beta Only
'%s' is only available from beta onwards
`

var deprecatedTemplate = `## Deprecated
Please use '%s' instead

This command will be removed on or around the %s.

%s`

var filePrepender = func(in string) string {
	log.Logger().Infof("Generating %s", in)

	name := filepath.Base(in)
	base := strings.TrimSuffix(name, path.Ext(name))
	command := strings.Replace(base, "_", " ", -1)

	content := ""

	alpha, ok := deprecation.AlphaOnly[command]
	if ok {
		log.Logger().Infof("Cmd is alpha only '%s'", command)
		content += fmt.Sprintf(alphaOnlyTemplate, alpha.Replacement)
	}

	_, ok = deprecation.BetaOnly[command]
	if ok {
		log.Logger().Infof("Cmd is beta only '%s'", command)
		content += fmt.Sprintf(betaOnlyTemplate, command)
	}

	dc, ok := deprecation.DeprecatedCommands[command]
	if ok {
		log.Logger().Infof("Cmd is deprecated '%s'", command)
		content += fmt.Sprintf(deprecatedTemplate, dc.Replacement, dc.Date, dc.Info)
	}

	return content
}

var docsCmd = &cobra.Command{
	Use:    "docs",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.DisableAutoGenTag = true

		err := doc.GenMarkdownTreeCustom(RootCmd, "./docs", filePrepender, linkHandler)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var completionCmd = &cobra.Command{
	Use:    "completion",
	Hidden: true,
	Short:  "Generate shell completion scripts",
	Long: `Generate shell completion scripts for GitHub CLI commands.

The output of this command will be computer code and is meant to be saved to a
file or immediately evaluated by an interactive shell.

For example, for bash you could add this to your '~/.bash_profile':

	eval "$(gh completion -s bash)"

When installing GitHub CLI through a package manager, however, it's possible that
no additional shell configuration is necessary to gain completion support. For
Homebrew, see <https://docs.brew.sh/Shell-Completion>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shellType, err := cmd.Flags().GetString("shell")
		if err != nil {
			return err
		}

		if shellType == "" {
			shellType = "bash"
		}

		switch shellType {
		case "bash":
			return RootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return RootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return RootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			return RootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
		default:
			return fmt.Errorf("unsupported shell type %q", shellType)
		}
	},
}

func changelogURL(version string) string {
	path := "https://github.com/plumming/dx"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}

	url := fmt.Sprintf("%s/releases/tag/v%s", path, strings.TrimPrefix(version, "v"))
	return url
}

func main() {
	currentVersion := version.Version
	updateMessageChan := make(chan *update.ReleaseInfo)
	go func() {
		rel, _ := checkForUpdate(currentVersion)
		updateMessageChan <- rel
	}()

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	newRelease := <-updateMessageChan
	if newRelease != nil {
		msg := fmt.Sprintf("%s %s → %s\n%s\n\n%s",
			util.ColorInfo("A new release of dx is available:"),
			util.ColorWarning(currentVersion),
			util.ColorWarning(newRelease.Version),
			util.ColorInfo(newRelease.URL),
			util.ColorInfo("dx upgrade cli"))

		stderr := os.Stderr
		fmt.Fprintf(stderr, "\n\n%s\n\n", msg)
	}

	os.Exit(0)
}

func shouldCheckForUpdate() bool {
	return updaterEnabled != "" && !isCompletionCommand() // && utils.IsTerminal(os.Stderr)
}

func isCompletionCommand() bool {
	return len(os.Args) > 1 && os.Args[1] == "completion"
}

func checkForUpdate(currentVersion string) (*update.ReleaseInfo, error) {
	if !shouldCheckForUpdate() {
		log.Logger().Debug("checking for updates is disabled")
		return nil, nil
	}

	client, err := api.BasicClient()
	if err != nil {
		return nil, err
	}

	repo := updaterEnabled

	if _, err := os.Stat(util.ConfigDir()); os.IsNotExist(err) {
		err = os.Mkdir(util.ConfigDir(), 0755)
		if err != nil {
			log.Logger().Warnf("unable to create %s directory", util.ConfigDir())
		}
	}
	stateFilePath := path.Join(util.ConfigDir(), "state.yml")
	return update.CheckForUpdate(client, stateFilePath, repo, currentVersion)
}
