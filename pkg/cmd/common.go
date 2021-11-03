package cmd

import (
	"github.com/plumming/dx/pkg/api"
	"github.com/plumming/dx/pkg/auth"
	"github.com/plumming/dx/pkg/config"
	"github.com/plumming/dx/pkg/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/plumming/dx/pkg/prompter"
	"k8s.io/client-go/kubernetes/fake"
)

type CommonOptions struct {
	prompter     prompter.Prompter
	githubClient *api.Client
	authConfig   auth.Config
	kuber        kube.Kuber
	dxConfig     config.Config
	kubeClient   kubernetes.Interface
}

func (c *CommonOptions) SetPrompter(p prompter.Prompter) {
	c.prompter = p
}

func (c *CommonOptions) Prompter() prompter.Prompter {
	if c.prompter == nil {
		c.prompter = prompter.NewPrompter()
	}
	return c.prompter
}

func (c *CommonOptions) SetGithubClient(h *api.Client) {
	c.githubClient = h
}

func (c *CommonOptions) GithubClient() (*api.Client, error) {
	if c.githubClient == nil {
		authConfig, err := c.AuthConfig()
		if err != nil {
			return nil, err
		}
		c.githubClient, err = api.BasicClient(authConfig)
		if err != nil {
			return nil, err
		}
	}
	return c.githubClient, nil
}

func (c *CommonOptions) SetKuber(k kube.Kuber) {
	c.kuber = k
}

func (c *CommonOptions) Kuber() kube.Kuber {
	if c.kuber == nil {
		c.kuber = kube.NewKuber()
	}
	return c.kuber
}

func (c *CommonOptions) SetDxConfig(dxConfig config.Config) {
	c.dxConfig = dxConfig
}

func (c *CommonOptions) DxConfig() (config.Config, error) {
	if c.dxConfig == nil {
		con, err := config.LoadFromDefaultLocation()
		if err != nil {
			return nil, err
		}
		c.dxConfig = con
	}
	return c.dxConfig, nil
}

func (c *CommonOptions) FakeKubeClient() kubernetes.Interface {
	if c.kubeClient == nil {
		c.kubeClient = fake.NewSimpleClientset()
	}
	return c.kubeClient
}

func (c *CommonOptions) KubeClient(config *rest.Config) (kubernetes.Interface, error) {
	if c.kubeClient == nil {
		var err error
		c.kubeClient, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
	}
	return c.kubeClient, nil
}

func (c *CommonOptions) AuthConfig() (auth.Config, error) {
	if c.authConfig == nil {
		apiConfig, err := auth.NewDefaultConfig()
		if err != nil {
			return nil, err
		}
		c.authConfig = apiConfig
	}
	return c.authConfig, nil
}
