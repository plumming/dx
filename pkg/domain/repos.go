package domain

import (
	"fmt"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/cmd"
)

type repoGetter func(string) ([]RepoInfo, error)

type Repo struct {
	cmd.CommonOptions
}

func (r *Repo) ListRepositoriesForOrg(org string) ([]RepoInfo, error) {
	client, err := r.GithubClient()
	if err != nil {
		return nil, err
	}

	var repos = []RepoInfo{}
	err = client.REST("GET", fmt.Sprintf("orgs/%s/repos", org), nil, &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *Repo) ListRepositoriesForUser(user string) ([]RepoInfo, error) {
	client, err := r.GithubClient()
	if err != nil {
		return nil, err
	}

	var repos = []RepoInfo{}
	err = client.REST("GET", fmt.Sprintf("users/%s/repos", user), nil, &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *Repo) SelectRepositories(filter string, in repoGetter) ([]string, error) {
	repos, err := in(filter)
	if err != nil {
		return nil, err
	}

	prompter := r.Prompter()

	selected, err := prompter.SelectMultipleFromOptions("Select repositories to delete", reposAsStringArray(repos))
	if err != nil {
		return nil, err
	}

	return selected, nil
}

func (r *Repo) DeleteRepositoriesFromOrg(org string) error {
	client, err := r.GithubClient()
	if err != nil {
		return err
	}

	selected, err := r.SelectRepositories(org, r.ListRepositoriesForOrg)
	if err != nil {
		return err
	}

	for _, s := range selected {
		log.Logger().Infof("deleting %s/%s", org, s)
		err = client.REST("DELETE", fmt.Sprintf("repos/%s/%s", org, s), nil, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) DeleteRepositoriesFromUser(user string) error {
	client, err := r.GithubClient()
	if err != nil {
		return err
	}

	selected, err := r.SelectRepositories(user, r.ListRepositoriesForUser)
	if err != nil {
		return err
	}

	for _, s := range selected {
		log.Logger().Infof("deleting %s/%s", user, s)
		err = client.REST("DELETE", fmt.Sprintf("repos/%s/%s", user, s), nil, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func reposAsStringArray(repos []RepoInfo) []string {
	var reposAsString = []string{}
	for _, r := range repos {
		reposAsString = append(reposAsString, r.Name)
	}
	return reposAsString
}

type RepoInfo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}
