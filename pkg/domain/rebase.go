package domain

import (
	"github.com/plumming/dx/pkg/cmd"
	"github.com/plumming/dx/pkg/util"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Rebase struct {
	cmd.CommonOptions
	Org string
	Repo string
	DefaultBranch string
	Config  *api.Config
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

	c.Org, c.Repo, err = GetOrgAndRepoFromCurrentDir()
	if err != nil {
		return err
	}
	logrus.Infof("determined repo as %s/%s", c.Org, c.Repo)

	c.DefaultBranch, err = GetDefaultBranch(gh, c.Org, c.Repo)
	logrus.Infof("determined default branch as %s", c.DefaultBranch)
	if err != nil {
		return err
	}

	return nil
}

// Run the cmd.
func (c *Rebase) Run() error {
	// should check if there are local changes

	// should check if we are on the non default branch

	// git fetch --tags upstream master
	cmd := util.Command{
		Name: "git",
		Args: []string{"fetch", "--tags", "upstream", c.DefaultBranch},
	}
	output, err := cmd.RunWithoutRetry()
	if err != nil {
		return err
	}
	logrus.Info(output)

	// git rebase upstream/master
	cmd = util.Command{
		Name: "git",
		Args: []string{"rebase", "upstream/" + c.DefaultBranch},
	}
	output, err = cmd.RunWithoutRetry()
	if err != nil {
		return err
	}
	logrus.Info(output)

	// git push origin master
	cmd = util.Command{
		Name: "git",
		Args: []string{"push", "origin", c.DefaultBranch},
	}
	output, err = cmd.RunWithoutRetry()
	if err != nil {
		return err
	}
	logrus.Info(output)

	return nil
}
