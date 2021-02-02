package domain

import (
	"fmt"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/pkg/errors"
	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/util"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Rebase struct {
	cmd.CommonOptions
	OriginOrg     string
	OriginRepo    string
	UpstreamOrg   string
	UpstreamRepo  string
	DefaultBranch string
	Config        *api.Config
}

// NewRebase.
func NewRebase() *Rebase {
	c := &Rebase{}
	return c
}

// Validate input.
func (c *Rebase) Validate() error {
	gh, err := c.GithubClient()
	if err != nil {
		return err
	}

	origin, err := GetRemote("origin")
	if err != nil {
		return err
	}
	upstream, err := GetRemote("upstream")
	if err != nil {
		return err
	}

	if upstream == "" {
		log.Logger().Warnf("No remote named 'upstream' found")
	}

	if origin == upstream {
		return errors.New("origin & upstream appear to be the same: " + origin)
	}

	c.OriginOrg, c.OriginRepo, err = ExtractOrgAndRepoURL(origin)
	if err != nil {
		return err
	}
	log.Logger().Debugf("determined origin repo as %s/%s", c.OriginOrg, c.OriginRepo)

	if upstream != "" {
		c.UpstreamOrg, c.UpstreamRepo, err = ExtractOrgAndRepoURL(upstream)
		if err != nil {
			return err
		}
		log.Logger().Debugf("determined upstream repo as %s/%s", c.UpstreamOrg, c.UpstreamRepo)
	}

	c.DefaultBranch, err = GetDefaultBranch(gh, c.OriginOrg, c.OriginOrg)
	log.Logger().Debugf("determined default branch as %s", c.DefaultBranch)
	if err != nil {
		return err
	}

	return nil
}

// Run the cmd.
func (c *Rebase) Run() error {
	// should check if there are local changes
	localChanges, err := LocalChanges("")
	if err != nil {
		return err
	}

	if localChanges {
		log.Logger().Error("There appear to be local changes, please stash and try again")
		return nil
	}

	// should check if we are on the non default branch
	currentBranch, err := CurrentBranchName("")
	if err != nil {
		return err
	}
	if c.DefaultBranch != currentBranch {
		log.Logger().Errorf("You appear to not be on the default branch, please switch to %s", c.DefaultBranch)
		return nil
	}

	if c.UpstreamRepo == "" && c.UpstreamOrg == "" {
		// git fetch --tags upstream master
		cmd := util.Command{
			Name: "git",
			Args: []string{"pull", "--tags", "origin", c.DefaultBranch},
		}
		output, err := Runner.RunWithoutRetry(&cmd)
		if err != nil {
			return err
		}
		log.Logger().Info(output)
	} else {
		// git fetch --tags upstream master
		cmd := util.Command{
			Name: "git",
			Args: []string{"fetch", "--tags", "upstream", c.DefaultBranch},
		}
		output, err := Runner.RunWithoutRetry(&cmd)
		if err != nil {
			return err
		}
		log.Logger().Info(output)

		// git rebase upstream/master
		cmd = util.Command{
			Name: "git",
			Args: []string{"rebase", fmt.Sprintf("upstream/%s", c.DefaultBranch)},
		}
		output, err = Runner.RunWithoutRetry(&cmd)
		if err != nil {
			return err
		}
		log.Logger().Info(output)

		// git push origin master
		cmd = util.Command{
			Name: "git",
			Args: []string{"push", "origin", c.DefaultBranch},
		}
		output, err = Runner.RunWithoutRetry(&cmd)
		if err != nil {
			return err
		}
		log.Logger().Info(output)
	}

	return nil
}
