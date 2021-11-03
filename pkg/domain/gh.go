package domain

import (
	"fmt"

	"github.com/plumming/dx/pkg/api"
)

func GetDefaultBranch(client *api.Client, host string, org string, repo string) (string, error) {
	var repository repository
	err := client.REST(host, "GET", fmt.Sprintf("repos/%s/%s", org, repo), nil, &repository)
	if err != nil {
		return "", err
	}

	return repository.DefaultBranch, nil
}

func GetCurrentUser(client *api.Client, host string) (string, error) {
	var currentUser currentUser
	err := client.REST(host, "GET", "user", nil, &currentUser)
	if err != nil {
		return "", err
	}

	return currentUser.Login, nil
}

func GetOrgsForUser(client *api.Client, host string) ([]string, error) {
	var organisations []organisation

	err := client.REST(host, "GET", "user/orgs", nil, &organisations)
	if err != nil {
		return nil, err
	}

	var orgs []string
	for _, o := range organisations {
		orgs = append(orgs, o.Login)
	}

	return orgs, nil
}

type organisation struct {
	Login string `json:"login"`
}

type currentUser struct {
	Login string `json:"login"`
}

type repository struct {
	DefaultBranch string `json:"default_branch"`
}
