package kube

import (
	"errors"
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"net/http"
	"os"
)

type factory struct {
}

func (f *factory) SetKubeContext(context string, config *api.Config) (*api.Config, error) {
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

func (f *factory) SetKubeNamespace(namespace string, config *api.Config) (*api.Config, error) {
	newConfig := *config
	ctx := f.GetCurrentContext(config)
	if ctx == nil {
		return nil, errors.New("could not find Kubernetes context")
	}
	if ctx.Namespace == namespace {
		return config, nil
	}

	ctx.Namespace = namespace
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

func (f *factory) GetCurrentContext(config *api.Config) *api.Context {
	if config != nil {
		name := config.CurrentContext
		if name != "" && config.Contexts != nil {
			return config.Contexts[name]
		}
	}
	return nil
}

func (f *factory) GetCurrentNamespace(config *api.Config) string {
	ctx := f.GetCurrentContext(config)
	if ctx != nil {
		n := ctx.Namespace
		if n != "" {
			return n
		}
	}
	return "default"
}

func (f *factory) CreateKubeClientConfig() (*rest.Config, error) {
	po := clientcmd.NewDefaultPathOptions()
	if po == nil {
		return nil, errors.New("unable to get kube config path options")
	}
	restConfig, err := clientcmd.BuildConfigFromFlags("", po.GlobalFile)
	if err != nil {
		return nil, err
	}
	// for testing purposes one can enable tracing of Kube REST API calls
	traceKubeAPI := os.Getenv("TRACE_KUBE_API")
	if traceKubeAPI == "1" || traceKubeAPI == "on" {
		restConfig.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
			return &Tracer{RoundTripper: rt}
		}
	}
	return restConfig, nil
}