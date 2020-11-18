package domain

import (
	"errors"
	"fmt"
	"io"
	url2 "net/url"
	"strings"

	"github.com/plumming/dx/pkg/api"
	"github.com/plumming/dx/pkg/util"
)

func GetDefaultBranch(client *api.Client, org string, repo string) (string, error) {
	repository := repository{}
	err := client.REST("GET", fmt.Sprintf("repos/%s/%s", org, repo), nil, &repository)
	if err != nil {
		return "", err
	}

	return repository.DefaultBranch, nil
}

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

type repository struct {
	DefaultBranch string `json:"default_branch"`
}
