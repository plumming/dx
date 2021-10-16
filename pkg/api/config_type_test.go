package api

import (
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

var (
	fileContent = `hosts:
  github.com:
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
)

func TestCanLoadConfig(t *testing.T) {
	f := fileConfig{}
	err := yaml.Unmarshal([]byte(fileContent), &f)
	assert.NoError(t, err)

	user := f.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := f.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}
