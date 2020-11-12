package update

import (
	"fmt"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-version"
	"github.com/plumming/dx/pkg/api"

	"io/ioutil"
	"time"
)

// ReleaseInfo stores information about a release.
type ReleaseInfo struct {
	Version string `json:"tag_name"`
	URL     string `json:"html_url"`
}

type StateEntry struct {
	CheckedForUpdateAt time.Time   `yaml:"checked_for_update_at"`
	LatestRelease      ReleaseInfo `yaml:"latest_release"`
}

// CheckForUpdate checks whether this software has had a newer release on GitHub.
func CheckForUpdate(client *api.Client, stateFilePath, repo, currentVersion string) (*ReleaseInfo, error) {
	log.Logger().Debugf("checking for update - current version: %s", currentVersion)
	latestRelease, err := GetLatestReleaseInfo(client, stateFilePath, repo, false)
	if err != nil {
		return nil, err
	}

	log.Logger().Debugf("latest release is %s", latestRelease)
	if versionGreaterThan(latestRelease.Version, currentVersion) {
		return latestRelease, nil
	}

	return nil, nil
}

func GetLatestReleaseInfo(client *api.Client, stateFilePath, repo string, force bool) (*ReleaseInfo, error) {
	if !force {
		stateEntry, err := getStateEntry(stateFilePath)
		if err == nil && time.Since(stateEntry.CheckedForUpdateAt).Hours() < 24 {
			return &stateEntry.LatestRelease, nil
		}
	}

	var latestRelease ReleaseInfo
	err := client.REST("GET", fmt.Sprintf("repos/%s/releases/latest", repo), nil, &latestRelease)
	if err != nil {
		return nil, err
	}

	err = setStateEntry(stateFilePath, time.Now(), latestRelease)
	if err != nil {
		return nil, err
	}

	return &latestRelease, nil
}

func getStateEntry(stateFilePath string) (*StateEntry, error) {
	content, err := ioutil.ReadFile(stateFilePath)
	if err != nil {
		return nil, err
	}

	var stateEntry StateEntry
	err = yaml.Unmarshal(content, &stateEntry)
	if err != nil {
		return nil, err
	}

	return &stateEntry, nil
}

func setStateEntry(stateFilePath string, t time.Time, r ReleaseInfo) error {
	data := StateEntry{CheckedForUpdateAt: t, LatestRelease: r}
	content, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	_ = ioutil.WriteFile(stateFilePath, content, 0600)

	return nil
}

func versionGreaterThan(v, w string) bool {
	log.Logger().Debugf("checking if %s is greater than %s", v, w)
	vv, ve := version.NewVersion(v)
	vw, we := version.NewVersion(w)

	log.Logger().Debugf("checking if %s is greater than %s - parsed", vv, vw)

	return ve == nil && we == nil && vv.GreaterThan(vw)
}
