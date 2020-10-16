package domain

import (
	"testing"

	"github.com/plumming/chilly/pkg/kube/kubefakes"

	"github.com/plumming/chilly/pkg/prompter/prompterfakes"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestConnectWorkspace_Validate_AllDataSupplied(t *testing.T) {
	d := NewContext()

	c := testKubeConfig()

	var kuber = &kubefakes.FakeKuber{}
	d.SetKuber(kuber)
	kuber.LoadConfigReturns(&c, nil)

	var prompter = &prompterfakes.FakePrompter{}
	d.SetPrompter(prompter)
	prompter.SelectFromOptionsReturns("xxx", nil)

	err := d.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "xxx", d.Context)
}

func testKubeConfig() api.Config {
	clusterA := api.Cluster{
		Server: "",
	}
	contextA := api.Context{
		Cluster: "ClusterA",
	}

	conf := api.Config{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "",
		Contexts: map[string]*api.Context{
			"contextA": &contextA,
		},
		Clusters: map[string]*api.Cluster{
			"clusterA": &clusterA,
		},
		AuthInfos:   map[string]*api.AuthInfo{},
		Preferences: api.Preferences{},
	}
	return conf
}
