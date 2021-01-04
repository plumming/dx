package domain

import (
	"context"
	"fmt"
	"sort"

	"k8s.io/client-go/rest"

	"github.com/pkg/errors"

	"github.com/plumming/dx/pkg/cmd"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Namespace defines Kubernetes namespace.
type Namespace struct {
	cmd.CommonOptions
	Namespace  string
	APIConfig  *api.Config
	RESTConfig *rest.Config
}

// NewNamespace.
func NewNamespace() *Namespace {
	n := &Namespace{}
	return n
}

// Validate input.
func (n *Namespace) Validate() error {
	kuber := n.Kuber()
	var err error
	n.APIConfig, err = kuber.LoadAPIConfig()
	if err != nil {
		return errors.Wrap(err, "failed to create load api config")
	}
	n.RESTConfig, err = kuber.LoadClientConfig()
	if err != nil {
		return errors.Wrap(err, "failed to create load client config")
	}
	n.Namespace, err = n.selectNamespace()
	if err != nil {
		return errors.Wrap(err, "failed to select namespace")
	}
	return nil
}

// Run the cmd.
func (n *Namespace) Run() error {
	fmt.Printf("you selected namespace %s", n.Namespace)
	kuber := n.Kuber()
	var err error
	n.APIConfig, err = kuber.SetKubeNamespace(n.Namespace, n.APIConfig)
	if err != nil {
		return err
	}
	return nil
}

func (n *Namespace) selectNamespace() (string, error) {
	currentNamespace := n.Kuber().GetCurrentNamespace(n.APIConfig)
	namespaces, err := n.loadNamespaces()
	if err != nil {
		return "", errors.Wrap(err, "while loading namespaces")
	}
	prompter := n.Prompter()
	namespace, err := prompter.SelectFromOptionsWithDefault("Select a namespace:", currentNamespace, namespaces)
	if err != nil {
		return "", errors.Wrap(err, "failed selecting namespace from prompter")
	}
	return namespace, nil
}

func (n *Namespace) loadNamespaces() ([]string, error) {
	var namespaces []string
	client, err := n.KubeClient(n.RESTConfig)
	if err != nil {
		return namespaces, errors.Wrap(err, "failed to create kube client")
	}
	ctx := context.TODO()
	if ctx != nil {
		list, err := client.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
		if err != nil {
			return namespaces, fmt.Errorf("loading namespaces %s", err)
		}
		for _, n := range list.Items {
			namespaces = append(namespaces, n.Name)
		}
		sort.Strings(namespaces)
	}
	return namespaces, nil
}
