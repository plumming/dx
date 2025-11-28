package auth

import (
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
	emptyConfigFile = ""

	dummyHostsFile = ` github.com:
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
	errorHostsFile = ` github.com
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
	emptyHostsFile = ""
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
	configFile, err := os.CreateTemp(os.TempDir(), "TestCanLoadConfigFromFile")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(configFile.Name()) }()

	err = os.WriteFile(configFile.Name(), []byte(dummyConfigFile), 0600)
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
	hostsFile, err := os.CreateTemp(os.TempDir(), "TestCanLoadHostsFromFile")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(hostsFile.Name()) }()

	err = os.WriteFile(hostsFile.Name(), []byte(dummyHostsFile), 0600)
	assert.NoError(t, err)

	config, err := parseHostsFile(hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseConfigFile(t *testing.T) {
	configFile, err := os.CreateTemp(os.TempDir(), "TestParseConfigFile")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(configFile.Name()) }()

	err = os.WriteFile(configFile.Name(), []byte(dummyConfigFile), 0600)
	assert.NoError(t, err)

	config, err := parseConfigFile(configFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseConfigFileWithYamlError(t *testing.T) {
	configFile, err := os.CreateTemp(os.TempDir(), "TestParseConfigFileWithYamlError")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(configFile.Name()) }()

	err = os.WriteFile(configFile.Name(), []byte(errorConfigFile), 0600)
	assert.NoError(t, err)

	config, err := parseConfigFile(configFile.Name())
	assert.Error(t, err)
	assert.Equal(t, nil, config)
}

func TestParseHostsFile(t *testing.T) {
	hostsFile, err := os.CreateTemp(os.TempDir(), "TestParseHostsFile")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(hostsFile.Name()) }()

	err = os.WriteFile(hostsFile.Name(), []byte(dummyHostsFile), 0600)
	assert.NoError(t, err)

	config, err := parseHostsFile(hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseHostsFileWithYamlError(t *testing.T) {
	hostsFile, err := os.CreateTemp(os.TempDir(), "TestParseHostsFileWithYamlError")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(hostsFile.Name()) }()

	err = os.WriteFile(hostsFile.Name(), []byte(errorHostsFile), 0600)
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

func TestParseDefaultConfigWithHostsFile(t *testing.T) {
	hostsFile, err := os.CreateTemp(os.TempDir(), "TestParseDefaultConfigWithHostsFile1")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(hostsFile.Name()) }()
	configFile, err := os.CreateTemp(os.TempDir(), "TestParseDefaultConfigWithHostsFile2")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(configFile.Name()) }()

	err = os.WriteFile(hostsFile.Name(), []byte(dummyHostsFile), 0600)
	assert.NoError(t, err)

	err = os.WriteFile(configFile.Name(), []byte(emptyConfigFile), 0600)
	assert.NoError(t, err)

	config, err := ParseDefaultConfig(configFile.Name(), hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseDefaultConfigWithConfigFile(t *testing.T) {
	hostsFile, err := os.CreateTemp(os.TempDir(), "TestParseDefaultConfigWithConfigFile1")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(hostsFile.Name()) }()
	configFile, err := os.CreateTemp(os.TempDir(), "TestParseDefaultConfigWithConfigFile2")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(configFile.Name()) }()

	err = os.WriteFile(hostsFile.Name(), []byte(emptyHostsFile), 0600)
	assert.NoError(t, err)

	err = os.WriteFile(configFile.Name(), []byte(dummyConfigFile), 0600)
	assert.NoError(t, err)

	config, err := ParseDefaultConfig(configFile.Name(), hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseDefaultConfigWithNoHostsFile(t *testing.T) {
	configFile, err := os.CreateTemp(os.TempDir(), "TestParseDefaultConfigWithNoHostsFile2")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(configFile.Name()) }()

	err = os.WriteFile(configFile.Name(), []byte(dummyConfigFile), 0600)
	assert.NoError(t, err)

	config, err := ParseDefaultConfig(configFile.Name(), "hosts-file-does-not-exist")
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}
