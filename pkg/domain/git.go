package domain

import (
	"errors"
	"io"
	url2 "net/url"
	"strings"

	"github.com/plumming/dx/pkg/util"
)

func GetOrgAndRepoFromCurrentDir() (string, string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"remote", "-v"},
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return "", "", err
	}

	return ExtractOrgAndRepoFromGitRemotes(strings.NewReader(output))
}

func CurrentBranchName(dir string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"branch", "--show-current"},
		Dir:  dir,
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return "", err
	}
	return output, nil
}

func SwitchBranch(dir string, name string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"checkout", name},
		Dir:  dir,
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return "", err
	}
	return output, nil
}

func Stash(dir string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"stash"},
		Dir:  dir,
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return "", err
	}
	return output, nil
}

func StashPop(dir string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"stash", "pop"},
		Dir:  dir,
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return "", err
	}
	return output, nil
}

func Add(dir string, name string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"add", name},
		Dir:  dir,
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return "", err
	}
	return output, nil
}

func Commit(dir string, message string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"commit", "-m", message},
		Dir:  dir,
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return "", err
	}
	return output, nil
}

func Status(dir string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"status"},
		Dir:  dir,
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return "", err
	}
	return output, nil
}

func LocalChanges(dir string) (bool, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"status", "--porcelain"},
		Dir:  dir,
	}
	output, err := c.RunWithoutRetry()
	if err != nil {
		return false, err
	}
	return output != "", nil
}

func ExtractOrgAndRepoFromGitRemotes(reader io.Reader) (string, string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return "", "", err
	}
	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "origin") && strings.HasSuffix(line, "(push)") {
			urlString := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, "origin\t", ""), "(push)", ""))
			url, err := url2.Parse(urlString)
			if err != nil {
				return "", "", err
			}
			fragments := strings.Split(url.Path, "/")
			if len(fragments) != 3 {
				return "", "", errors.New("invalid url path '" + url.Path + "'")
			}
			return fragments[1], fragments[2], nil
		}
	}
	return "", "", errors.New("unable to find remote named 'origin'")
}
