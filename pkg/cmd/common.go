package cmd

import (
	"github.com/plumming/dx/pkg/api"
	"github.com/plumming/dx/pkg/config"
	"github.com/plumming/dx/pkg/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/plumming/dx/pkg/prompter"
	"k8s.io/client-go/kubernetes/fake"
)

type CommonOptions struct {
	prompter        prompter.Prompter
	nonGithubClient *api.Client
	githubClient    *api.Client
	kuber           kube.Kuber
	config          config.Config
	kubeClient      kubernetes.Interface
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

func (c *CommonOptions) SetNonGithubClient(h *api.Client) {
	c.nonGithubClient = h
}

func (c *CommonOptions) NonGithubClient() *api.Client {
	if c.nonGithubClient == nil {
		c.nonGithubClient = api.NewClient()
	}
	return c.nonGithubClient
}

func (c *CommonOptions) SetGithubClient(h *api.Client) {
	c.githubClient = h
}

func (c *CommonOptions) GithubClient() (*api.Client, error) {
	if c.githubClient == nil {
		var err error
		c.githubClient, err = api.BasicClient()
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

func (c *CommonOptions) SetConfig(config config.Config) {
	c.config = config
}

func (c *CommonOptions) Config() (config.Config, error) {
	if c.config == nil {
		con, err := config.LoadFromDefaultLocation()
		if err != nil {
			return nil, err
		}
		c.config = con
	}
	return c.config, nil
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
