package api

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dummyConfigFile = `hosts:
  github.com:
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
	errorConfigFile = `hosts
  github.com:
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
	dummyHostsFile = ` github.com:
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
	errorHostsFile = ` github.com
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
)

func TestConfigFile(t *testing.T) {
	configFile := ConfigFile()
	assert.Contains(t, configFile, "config.yml")
}

func TestHostsFile(t *testing.T) {
	hostsFile := HostsFile()
	assert.Contains(t, hostsFile, "hosts.yml")
}

func TestCanLoadConfigFromFile(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "TestCanLoadConfigFromFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(configFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte(dummyConfigFile), 0600)
	assert.NoError(t, err)

	config, err := parseConfigFile(configFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseConfigWithNonExistentConfigFile(t *testing.T) {
	config, err := parseConfigFile("config-file-does-not-exist")
	assert.Nil(t, config)
	assert.Error(t, err)
	assert.EqualError(t, err, "open config-file-does-not-exist: no such file or directory")
}

func TestParseConfigWithNonExistentHostsFile(t *testing.T) {
	config, err := parseHostsFile("hosts-file-does-not-exist")
	assert.Equal(t, &fileConfig{}, config)
	assert.Error(t, err)
	assert.EqualError(t, err, "open hosts-file-does-not-exist: no such file or directory")
}

func TestCanLoadHostsFromFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile(os.TempDir(), "TestCanLoadHostsFromFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(hostsFile.Name())

	err = ioutil.WriteFile(hostsFile.Name(), []byte(dummyHostsFile), 0600)
	assert.NoError(t, err)

	config, err := parseHostsFile(hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseConfigFile(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "TestParseConfigFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(configFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte(dummyConfigFile), 0600)
	assert.NoError(t, err)

	config, err := parseConfigFile(configFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseConfigFileWithYamlError(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "TestParseConfigFileWithYamlError")
	if err != nil {
		panic(err)
	}
	defer os.Remove(configFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte(errorConfigFile), 0600)
	assert.NoError(t, err)

	config, err := parseConfigFile(configFile.Name())
	assert.Error(t, err)
	assert.Equal(t, nil, config)
}

func TestParseHostsFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile(os.TempDir(), "TestParseHostsFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(hostsFile.Name())

	err = ioutil.WriteFile(hostsFile.Name(), []byte(dummyHostsFile), 0600)
	assert.NoError(t, err)

	config, err := parseHostsFile(hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseHostsFileWithYamlError(t *testing.T) {
	hostsFile, err := ioutil.TempFile(os.TempDir(), "TestParseHostsFileWithYamlError")
	if err != nil {
		panic(err)
	}
	defer os.Remove(hostsFile.Name())

	err = ioutil.WriteFile(hostsFile.Name(), []byte(errorHostsFile), 0600)
	assert.NoError(t, err)

	config, err := parseHostsFile(hostsFile.Name())
	assert.Error(t, err)
	assert.Equal(t, &fileConfig{}, config)
}

func TestParseHostsWithNonExistentHostsFile(t *testing.T) {
	config, err := parseHostsFile("hosts-file-does-not-exist")
	assert.Equal(t, &fileConfig{}, config)
	assert.Error(t, err)
	assert.EqualError(t, err, "open hosts-file-does-not-exist: no such file or directory")
}
