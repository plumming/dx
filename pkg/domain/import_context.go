package domain

import (
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/cmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// ImportContext defines Kubernetes context.
type ImportContext struct {
	cmd.CommonOptions
	Path           string
	Config         *api.Config
	ConfigToImport *api.Config
}

// NewImportContext creates a new ImportContext.
func NewImportContext(path string) *ImportContext {
	c := &ImportContext{
		Path: path,
	}
	return c
}

// Validate input.
func (c *ImportContext) Validate() error {
	k := c.Kuber()
	var err error

	c.Config, err = k.LoadAPIConfig()
	if err != nil {
		return err
	}

	// load api config from path
	c.ConfigToImport, err = k.LoadAPIConfigFromPath(c.Path)
	if err != nil {
		return err
	}

	return nil
}

// Run the cmd.
func (c *ImportContext) Run() error {
	log.Logger().Infof("Importing from %s", c.Path)

	kuber := c.Kuber()

	newConfig := c.Config

	for k, v := range c.ConfigToImport.Contexts {
		log.Logger().Debugf("context %s", k)
		newConfig.Contexts[k] = v
	}

	for k, v := range c.ConfigToImport.AuthInfos {
		log.Logger().Debugf("authInfo %s", k)
		newConfig.AuthInfos[k] = v
	}

	for k, v := range c.ConfigToImport.Clusters {
		log.Logger().Debugf("cluster %s", k)
		newConfig.Clusters[k] = v
	}

	for k, v := range c.ConfigToImport.Extensions {
		log.Logger().Debugf("extensions %s", k)
		newConfig.Extensions[k] = v
	}

	_, err := kuber.SetKubeConfig(newConfig)
	if err != nil {
		return err
	}

	return nil
}
