package kube

import "k8s.io/client-go/tools/clientcmd/api"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Kuber

type Kuber interface {
	SetKubeContext(string, *api.Config) error
	LoadConfig() (*api.Config, error)
}
