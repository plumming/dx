package domain

import (
	"fmt"

	"github.com/plumming/dx/pkg/api"
)

func GetDefaultBranch(client *api.Client, org string, repo string) (string, error) {
	repository := repository{}
	err := client.REST("GET", fmt.Sprintf("repos/%s/%s", org, repo), nil, &repository)
	if err != nil {
		return "", err
	}

	return repository.DefaultBranch, nil
}

func GetCurrentUser(client *api.Client) (string, error) {
	currentUser := currentUser{}
	err := client.REST("GET", "user", nil, &currentUser)
	if err != nil {
		return "", err
	}

	return currentUser.Login, nil
}

type currentUser struct {
	Login string `json:"login"`
}

type repository struct {
	DefaultBranch string `json:"default_branch"`
}
