package domain

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"sort"

	"github.com/pkg/errors"

	"github.com/plumming/dx/pkg/cmd"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Namespace defines Kubernetes namespace.
type Namespace struct {
	cmd.CommonOptions
	Namespace string
	Config  *api.Config
}

// NewNamespace
func NewNamespace() *Namespace {
	n := &Namespace{}
	return n
}

// Validate input.
func (n *Namespace) Validate() error {
	kuber := n.Kuber()
	var err error
	n.Config, err = kuber.LoadConfig()
	if err != nil {
		return errors.Wrap(err, "failed to create api config")
	}
	restConfig, err := kuber.CreateKubeClientConfig()
	if err != nil {
		return errors.Wrap(err, "failed to create rest config")
	}
	client, err := n.KubeClient(restConfig)
	if err != nil {
		return errors.Wrap(err, "failed to creat kube client")
	}
	n.Namespace, err = n.selectNamespace(client)
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
	n.Config, err = kuber.SetKubeNamespace(n.Namespace, n.Config)
	if err != nil {
		return err
	}
	return nil
}

func (n *Namespace) selectNamespace(client kubernetes.Interface) (string, error) {
	namespaces, err := n.loadNamespaces(client)
	if err != nil {
		errors.Wrap(err, "while loading namespaces")
	}
	prompter := n.Prompter()
	ns, err := prompter.SelectFromOptions("Select a namespace:", namespaces)
	if err != nil {
		return "", errors.Wrap(err, "failed selecting namespace from prompter")
	}

	return ns, nil
}

func (n *Namespace) loadNamespaces(client kubernetes.Interface) ([]string, error) {
	var namespaces []string
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

