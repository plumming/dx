package domain_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/plumming/dx/pkg/util/mocks"

	"github.com/plumming/dx/pkg/domain"
	"github.com/plumming/dx/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestCanDetermineBranchName(t *testing.T) {
	cr := util.DefaultCommandRunner{}

	dir, err := ioutil.TempDir("", "domain_test__TestCanDetermineBranchName")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	c := util.Command{
		Name: "git",
		Args: []string{"init", "-b", "master"},
		Dir:  dir,
	}

	output, err := cr.RunWithoutRetry(&c)
	assert.NoError(t, err)
	t.Log(output)

	bn, err := domain.CurrentBranchName(dir)
	assert.NoError(t, err)
	t.Log(bn)
	assert.Equal(t, "master", bn)
}

func TestCanStash(t *testing.T) {
	cr := util.DefaultCommandRunner{}

	dir, err := ioutil.TempDir("", "domain_test__TestCanStash")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	c := util.Command{
		Name: "git",
		Args: []string{"init", "-b", "master"},
		Dir:  dir,
	}
	output, err := cr.RunWithoutRetry(&c)
	assert.NoError(t, err)
	t.Log(output)

	err = domain.ConfigCommitterInformation(dir, "test@test.com", "test user")
	assert.NoError(t, err)

	d1 := []byte("# domain_test__TestCanStash\n")
	err = ioutil.WriteFile(path.Join(dir, "README.md"), d1, 0600)
	assert.NoError(t, err)

	output, err = domain.Add(dir, "README.md")
	assert.NoError(t, err)
	t.Log(output)

	output, err = domain.Commit(dir, "Initial Commit")
	assert.NoError(t, err)
	t.Log(output)

	localChanges, err := domain.LocalChanges(dir)
	assert.NoError(t, err)
	assert.False(t, localChanges)

	d1 = []byte("hello\ngo\n")
	err = ioutil.WriteFile(path.Join(dir, "README.md"), d1, 0600)
	assert.NoError(t, err)

	localChanges, err = domain.LocalChanges(dir)
	assert.NoError(t, err)
	assert.True(t, localChanges)

	output, err = domain.Status(dir)
	assert.NoError(t, err)
	t.Log(output)

	output, err = domain.Stash(dir)
	assert.NoError(t, err)
	t.Log(output)

	localChanges, err = domain.LocalChanges(dir)
	assert.NoError(t, err)
	assert.False(t, localChanges)

	output, err = domain.StashPop(dir)
	assert.NoError(t, err)
	t.Log(output)

	localChanges, err = domain.LocalChanges(dir)
	assert.NoError(t, err)
	assert.True(t, localChanges)

	output, err = domain.Add(dir, "README.md")
	assert.NoError(t, err)
	t.Log(output)

	output, err = domain.Commit(dir, "Updated Commit")
	assert.NoError(t, err)
	t.Log(output)

	localChanges, err = domain.LocalChanges(dir)
	assert.NoError(t, err)
	assert.False(t, localChanges)

	d1 = []byte("hello\ngo\n")
	err = ioutil.WriteFile(path.Join(dir, "OTHER.md"), d1, 0600)
	assert.NoError(t, err)

	localChanges, err = domain.LocalChanges(dir)
	assert.NoError(t, err)
	assert.False(t, localChanges)
}

func TestLocalChanges(t *testing.T) {
	type test struct {
		name     string
		raw      string
		expected bool
	}

	tests := []test{
		{
			name:     "no changes",
			raw:      ``,
			expected: false,
		},
		{
			name:     "changes to existing files",
			raw:      ` M go.sum`,
			expected: true,
		},
		{
			name:     "new file",
			raw:      `?? ll`,
			expected: false,
		},
		{
			name: "changes to both existing and new files",
			raw: ` M go.sum
?? ll`,
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := mocks.MockCommandRunner{}
			domain.Runner = &r
			mocks.GetRunWithoutRetryFunc = func(c *util.Command) (string, error) {
				return tc.raw, nil
			}

			b, err := domain.LocalChanges("")
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, b)

			t.Logf("commands> %s", r.Commands)
		})
	}
}

func TestCanDetermineRemoteNames(t *testing.T) {
	type test struct {
		raw         string
		remote      string
		expectedURL string
	}

	tests := []test{
		{
			raw: `origin  https://github.com/garethjevans/chilly (fetch)
origin  https://github.com/garethjevans/chilly (push)
upstream        https://github.com/plumming/dx (fetch)
upstream        https://github.com/plumming/dx (push)`,
			remote:      "origin",
			expectedURL: "https://github.com/garethjevans/chilly",
		},
		{
			raw: `origin  https://github.com/garethjevans/chilly (fetch)
origin  https://github.com/garethjevans/chilly (push)
upstream        https://github.com/plumming/dx (fetch)
upstream        https://github.com/plumming/dx (push)`,
			remote:      "upstream",
			expectedURL: "https://github.com/plumming/dx",
		},
		{
			raw: `origin  https://github.com/garethjevans/chilly (fetch)
origin  https://github.com/garethjevans/chilly (push)`,
			remote:      "upstream",
			expectedURL: "",
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestCanDetermineRemoteNames-%s", tc.remote), func(t *testing.T) {
			url, err := domain.ExtractURLFromRemote(strings.NewReader(tc.raw), tc.remote)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedURL, url)
		})
	}
}
