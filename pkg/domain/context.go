package domain

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/plumming/dx/pkg/cmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Context defines Kubernetes context.
type Context struct {
	cmd.CommonOptions
	Context string
	Config  *api.Config
}

// NewContext.
func NewContext() *Context {
	c := &Context{}
	return c
}

// Validate input.
func (c *Context) Validate() error {
	kuber := c.Kuber()
	var err error
	c.Config, err = kuber.LoadConfig()
	if err != nil {
		return err
	}
	c.Context, err = c.selectContext()
	if err != nil {
		return errors.Wrap(err, "failed to select context")
	}
	return nil
}

// Run the cmd.
func (c *Context) Run() error {
	fmt.Printf("you selected context %s", c.Context)
	k := c.Kuber()
	var err error
	c.Config, err = k.SetKubeContext(c.Context, c.Config)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) selectContext() (string, error) {
	contexts := c.loadContexts()
	prompter := c.Prompter()
	ctx, err := prompter.SelectFromOptions("Select a context:", contexts)
	if err != nil {
		return "", errors.Wrap(err, "failed selecting context from prompter")
	}
	return ctx, nil
}

func (c *Context) loadContexts() []string {
	var contexts []string
	for k := range c.Config.Contexts {
		contexts = append(contexts, k)
	}
	return contexts
}
