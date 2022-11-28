package domain

import (
	"fmt"

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
	if c.Path == "" {
		return fmt.Errorf("path must be set")
	}

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
	log.Logger().Infof("Importing config from %s", c.Path)

	kuber := c.Kuber()

	newConfig := c.Config

	// get the location of origin for the first context
	locationOfOrigin := c.Config.Contexts[c.Config.CurrentContext].LocationOfOrigin
	log.Logger().Debugf("locationOfOrigin %s", locationOfOrigin)

	for k, v := range c.ConfigToImport.Contexts {
		log.Logger().Debugf("context %s from %s", k, v.LocationOfOrigin)
		newConfig.Contexts[k] = v
		newConfig.Contexts[k].LocationOfOrigin = locationOfOrigin
	}

	for k, v := range c.ConfigToImport.AuthInfos {
		log.Logger().Debugf("authInfo %s from %s", k, v.LocationOfOrigin)
		newConfig.AuthInfos[k] = v
		newConfig.AuthInfos[k].LocationOfOrigin = locationOfOrigin
	}

	for k, v := range c.ConfigToImport.Clusters {
		log.Logger().Debugf("cluster %s from %s", k, v.LocationOfOrigin)
		newConfig.Clusters[k] = v
		newConfig.Clusters[k].LocationOfOrigin = locationOfOrigin
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
