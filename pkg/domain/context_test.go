package domain

import (
	"testing"

	"github.com/plumming/dx/pkg/kube/kubefakes"

	"github.com/plumming/dx/pkg/prompter/prompterfakes"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestConnectWorkspace_Validate_AllDataSupplied(t *testing.T) {
	d := NewContext()

	c := testKubeConfig()
	var kuber = &kubefakes.FakeKuber{}
	d.SetKuber(kuber)
	kuber.LoadConfigReturns(&c, nil)

	testContext := getContextKey(c.Contexts, 0)

	var prompter = &prompterfakes.FakePrompter{}
	d.SetPrompter(prompter)
	prompter.SelectFromOptionsReturns(testContext, nil)

	err := d.Validate()
	assert.NoError(t, err)
	assert.Equal(t, testContext, d.Context)
}

func TestConnectWorkspace_Run(t *testing.T) {
	d := NewContext()

	c := testKubeConfig()
	d.Config = &c
	d.Context = getContextKey(c.Contexts, 0)

	x := &c
	x.CurrentContext = getContextKey(c.Contexts, 0)

	var kuber = &kubefakes.FakeKuber{}
	d.SetKuber(kuber)
	kuber.SetKubeContextReturns(x, nil)

	err := d.Run()
	assert.NoError(t, err)
	assert.Equal(t, x.CurrentContext, d.Context)
}

//TODO: unmarshal config from a file.
func testKubeConfig() api.Config {
	clusterA := api.Cluster{
		Server: "",
	}
	contextA := api.Context{
		Cluster: "ClusterA",
	}
	contextB := api.Context{
		Cluster: "ClusterB",
	}

	conf := api.Config{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "contextB",
		Contexts: map[string]*api.Context{
			"contextA": &contextA,
			"contextB": &contextB,
		},
		Clusters: map[string]*api.Cluster{
			"clusterA": &clusterA,
		},
		AuthInfos:   map[string]*api.AuthInfo{},
		Preferences: api.Preferences{},
	}
	return conf
}

func getContextKey(m map[string]*api.Context, index int) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys[index]
}
