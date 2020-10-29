package kube

import (
	"fmt"

	"errors"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type factory struct {
}

func (f *factory) SetKubeContext(context string, config *api.Config) (*api.Config, error) {
	fmt.Printf("you selected context %s", context)
	ctx := config.Contexts[context]
	if ctx == nil {
		return nil, fmt.Errorf("could not find Kubernetes context %s", context)
	}

	newConfig := *config
	newConfig.CurrentContext = context
	err := clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), newConfig, false)
	if err != nil {
		return nil, err
	}
	return &newConfig, nil
}

func NewKuber() Kuber {
	return &factory{}
}

func (f *factory) LoadConfig() (*api.Config, error) {
	po := clientcmd.NewDefaultPathOptions()
	if po == nil {
		return nil, errors.New("unable to get kube config path options")
	}
	config, err := po.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}
