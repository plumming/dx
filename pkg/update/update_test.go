package update

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/plumming/dx/pkg/auth"

	"github.com/plumming/dx/pkg/api"
)

func TestCheckForUpdate(t *testing.T) {
	scenarios := []struct {
		Name           string
		CurrentVersion string
		LatestVersion  string
		LatestURL      string
		ExpectsResult  bool
	}{
		{
			Name:           "latest is newer",
			CurrentVersion: "v0.0.1",
			LatestVersion:  "v1.0.0",
			LatestURL:      "https://www.spacejam.com/archive/spacejam/movie/jam.htm",
			ExpectsResult:  true,
		},
		{
			Name:           "current is prerelease",
			CurrentVersion: "v1.0.0-pre.1",
			LatestVersion:  "v1.0.0",
			LatestURL:      "https://www.spacejam.com/archive/spacejam/movie/jam.htm",
			ExpectsResult:  true,
		},
		{
			Name:           "latest is current",
			CurrentVersion: "v1.0.0",
			LatestVersion:  "v1.0.0",
			LatestURL:      "https://www.spacejam.com/archive/spacejam/movie/jam.htm",
			ExpectsResult:  false,
		},
		{
			Name:           "latest is older",
			CurrentVersion: "v0.10.0-pre.1",
			LatestVersion:  "v0.9.0",
			LatestURL:      "https://www.spacejam.com/archive/spacejam/movie/jam.htm",
			ExpectsResult:  false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			var authConfig auth.Config = &auth.FakeConfig{
				Hosts: map[string]*auth.HostConfig{
					"github.com": {User: "user", Token: "token"},
					"other.com":  {User: "otheruser", Token: "token2"},
				},
			}

			http := &api.FakeHTTP{}
			client := api.NewClient(authConfig, api.ReplaceTripper(http))
			http.StubResponse(200, bytes.NewBufferString(fmt.Sprintf(`{
				"tag_name": "%s",
				"html_url": "%s"
			}`, s.LatestVersion, s.LatestURL)))

			rel, err := CheckForUpdate(client, tempFilePath(), "OWNER/REPO", s.CurrentVersion)
			if err != nil {
				t.Fatal(err)
			}

			if len(http.Requests) != 1 {
				t.Fatalf("expected 1 HTTP request, got %d", len(http.Requests))
			}
			requestPath := http.Requests[0].URL.Path
			if requestPath != "/repos/OWNER/REPO/releases/latest" {
				t.Errorf("HTTP path: %q", requestPath)
			}

			if !s.ExpectsResult {
				if rel != nil {
					t.Fatal("expected no new release")
				}
				return
			}
			if rel == nil {
				t.Fatal("expected to report new release")
			}

			if rel.Version != s.LatestVersion {
				t.Errorf("Version: %q", rel.Version)
			}
			if rel.URL != s.LatestURL {
				t.Errorf("URL: %q", rel.URL)
			}
		})
	}
}

func tempFilePath() string {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(file.Name())
	return file.Name()
}
