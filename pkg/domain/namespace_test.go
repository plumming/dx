package domain

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/plumming/dx/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/plumming/dx/pkg/kube/kubefakes"

	"github.com/plumming/dx/pkg/prompter/prompterfakes"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_Namespace_Validate_AllDataSupplied(t *testing.T) {
	want := NewNamespace()
	c := testKubeConfigForNamespace()
	want.APIConfig = &c

	var kuber = &kubefakes.FakeKuber{}
	want.SetKuber(kuber)
	kuber.LoadAPIConfigReturns(want.APIConfig, nil)

	testNamespace := getNamespaceKey(want.APIConfig, 0)

	var prompter = &prompterfakes.FakePrompter{}
	want.SetPrompter(prompter)
	prompter.SelectFromOptionsWithDefaultReturns(testNamespace, nil)

	want.FakeKubeClient()

	err := want.Validate()
	assert.NoError(t, err)
	assert.Equal(t, testNamespace, want.Namespace)
}

func TestSelectNamespace_Run(t *testing.T) {
	want := NewNamespace()
	apiConfig := testKubeConfigForNamespace()
	want.APIConfig = &apiConfig
	testNamespace := getNamespaceKey(want.APIConfig, 0)
	want.Namespace = testNamespace

	config := &apiConfig
	config.Contexts[config.CurrentContext].Namespace = testNamespace

	var kuber = &kubefakes.FakeKuber{}
	want.SetKuber(kuber)
	kuber.LoadAPIConfigReturns(config, nil)

	err := want.Run()
	assert.NoError(t, err)
	assert.Equal(t, config.Contexts[config.CurrentContext].Namespace, want.Namespace)
}

func TestSelectNamespaceUsingConfigFile_Run(t *testing.T) {
	want := NewNamespace()
	kuber := want.Kuber()

	// setup files to read and update
	baseDir, err := ioutil.TempDir("", "test_select_namespace")
	assert.NoError(t, err)
	testData := path.Join("test_data", "namespace_test")
	_, err = os.Stat(testData)
	assert.NoError(t, err)

	outDir := path.Join(baseDir, "output")
	err = util.CopyDir(testData, outDir, true)
	assert.NoError(t, err)

	defaultContext := filepath.Join(outDir, "default_context.yaml")
	os.Setenv("KUBECONFIG", defaultContext)
	want.APIConfig, err = kuber.LoadAPIConfig()
	assert.NoError(t, err)

	testNamespace := getNamespaceKey(want.APIConfig, 0)
	want.Namespace = testNamespace

	err = want.Run()
	assert.NoError(t, err)
	gotAPIConfig, err := want.Kuber().LoadAPIConfig()
	assert.NoError(t, err)
	assert.Equal(t, want.APIConfig.Contexts[want.APIConfig.CurrentContext].Namespace, gotAPIConfig.Contexts[gotAPIConfig.CurrentContext].Namespace)
}

func TestLoadConfig(t *testing.T) {
	n := NewNamespace()
	kuber := n.Kuber()

	wantAPIConfig := testKubeConfigForNamespace()
	testData := path.Join("test_data", "namespace_test")
	_, err := os.Stat(testData)
	assert.NoError(t, err)
	defaultContext := filepath.Join(testData, "default_context.yaml")
	_, err = os.Stat(defaultContext)
	assert.NoError(t, err)
	os.Setenv("KUBECONFIG", defaultContext)

	gotAPIConfig, err := kuber.LoadAPIConfig()
	assert.NoError(t, err)
	assert.Equal(t, wantAPIConfig.CurrentContext, gotAPIConfig.CurrentContext)
	assert.Equal(t, wantAPIConfig.Clusters["clusterA"], gotAPIConfig.Clusters["clusterA"])
}

func testKubeConfigForNamespace() api.Config {
	clusters := make(map[string]*api.Cluster)
	clusterA := api.Cluster{
		Server:           "https://127.0.0.1:59801",
		LocationOfOrigin: "test_data/namespace_test/default_context.yaml",
		Extensions:       map[string]runtime.Object{},
	}
	clusters["clusterA"] = &clusterA
	contextA := api.Context{
		Cluster:   "clusterA",
		Namespace: "namespaceA",
	}
	contextB := api.Context{
		Cluster:   "clusterB",
		Namespace: "namespaceB",
	}

	conf := api.Config{
		APIVersion:     "v1",
		Kind:           "DxConfig",
		CurrentContext: "contextB",
		Contexts: map[string]*api.Context{
			"contextA": &contextA,
			"contextB": &contextB,
		},
		Clusters:    clusters,
		AuthInfos:   map[string]*api.AuthInfo{},
		Preferences: api.Preferences{},
	}
	return conf
}

func getNamespaceKey(c *api.Config, index int) string {
	keys := make([]string, 0, len(c.Contexts))
	for _, v := range c.Contexts {
		keys = append(keys, v.Namespace)
	}
	return keys[index]
}
