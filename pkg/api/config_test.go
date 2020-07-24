package api_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/plumming/chilly/pkg/api"

	"github.com/stretchr/testify/assert"
)

var (
	dummyConfigFile = `hosts:
  github.com:
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
)

func TestCanLoadConfigFromFile(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "TestCanLoadConfigFromFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	err = ioutil.WriteFile(file.Name(), []byte(dummyConfigFile), 0600)
	assert.NoError(t, err)

	config, err := api.ParseConfig(file.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}
