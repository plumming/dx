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

type repository struct {
	DefaultBranch string `json:"default_branch"`
}
