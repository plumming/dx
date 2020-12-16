package kube

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Kuber

type Kuber interface {
	SetKubeContext(string, *api.Config) (*api.Config, error)
	SetKubeNamespace(string, *api.Config) (*api.Config, error)
	LoadAPIConfig() (*api.Config, error)
	LoadClientConfig() (*rest.Config, error)
	GetCurrentContext(*api.Config) *api.Context
	GetCurrentNamespace(*api.Config) string
}
