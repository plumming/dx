package domain_test

import (
	"github.com/plumming/dx/pkg/domain"
	"github.com/plumming/dx/pkg/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCanDetermineBranchName(t *testing.T) {
	dir, err := ioutil.TempDir("", "domain_test__TestCanDetermineBranchName")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	c := util.Command{
		Name: "git",
		Args: []string{"init"},
		Dir: dir,
	}
	output, err := c.RunWithoutRetry()
	assert.NoError(t, err)
	t.Log(output)

	bn, err := domain.CurrentBranchName(dir)
	assert.NoError(t, err)
	t.Log(bn)
	assert.Equal(t, "master", bn)
}

func TestCanStash(t *testing.T) {
	dir, err := ioutil.TempDir("", "domain_test__TestCanStash")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	c := util.Command{
		Name: "git",
		Args: []string{"init"},
		Dir: dir,
	}
	output, err := c.RunWithoutRetry()
	assert.NoError(t, err)
	t.Log(output)

	d1 := []byte("# domain_test__TestCanStash\n")
	err = ioutil.WriteFile(path.Join(dir, "README.md"), d1, 0644)
	assert.NoError(t, err)

	output, err = domain.Add(dir, "README.md")
	assert.NoError(t, err)
	t.Log(output)

	output, err = domain.Commit(dir,"Initial Commit")
	assert.NoError(t, err)
	t.Log(output)

	localChanges, err := domain.LocalChanges(dir)
	assert.NoError(t, err)
	assert.False(t, localChanges)

	d1 = []byte("hello\ngo\n")
	err = ioutil.WriteFile(path.Join(dir, "README.md"), d1, 0644)
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
}



