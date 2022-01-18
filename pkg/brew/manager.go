package brew

import (
	"os/exec"
	"strings"

	"github.com/plumming/dx/pkg/util"
)

var Runner util.CommandRunner = util.DefaultCommandRunner{}

func IsBrewInstalled() (bool, error) {
	path, err := exec.LookPath("brew")
	if err != nil {
		return false, err
	}
	return path != "", nil
}

func IsInstalledViaBrew() (bool, error) {
	ok, err := IsBrewInstalled()
	if err != nil {
		return ok, err
	}
	if ok {
		c := util.Command{
			Name: "brew",
			Args: []string{"list"},
		}
		output, err := Runner.RunWithoutRetry(&c)
		if err != nil {
			return false, err
		}
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if line == "dx" {
				return true, nil
			}
		}
	}
	return false, nil
}
