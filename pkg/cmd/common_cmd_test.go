package cmd_test

import (
	cmd2 "github.com/plumming/chilly/pkg/cmd"
	"github.com/plumming/chilly/pkg/pr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommonCmd_Filter_AllData(t *testing.T) {
	prs := []pr.PullRequest{
		{Number: 1, Author: pr.Author{Login: "author1"}},
		{Number: 2, Author: pr.Author{Login: "author2"}},
		{Number: 3, Author: pr.Author{Login: "author3"}},
	}

	cmd := cmd2.CommonCmd{Query: "[]"}
	result, err := cmd.Filter(prs)
	assert.NoError(t, err)
	expected := `[
  {
    "author": {
      "login": "author1"
    },
    "closed": false,
    "commits": {
      "nodes": null
    },
    "createdAt": "0001-01-01T00:00:00Z",
    "labels": {
      "nodes": null
    },
    "mergeable": "",
    "number": 1,
    "repository": {
      "nameWithOwner": ""
    },
    "title": "",
    "url": ""
  },
  {
    "author": {
      "login": "author2"
    },
    "closed": false,
    "commits": {
      "nodes": null
    },
    "createdAt": "0001-01-01T00:00:00Z",
    "labels": {
      "nodes": null
    },
    "mergeable": "",
    "number": 2,
    "repository": {
      "nameWithOwner": ""
    },
    "title": "",
    "url": ""
  },
  {
    "author": {
      "login": "author3"
    },
    "closed": false,
    "commits": {
      "nodes": null
    },
    "createdAt": "0001-01-01T00:00:00Z",
    "labels": {
      "nodes": null
    },
    "mergeable": "",
    "number": 3,
    "repository": {
      "nameWithOwner": ""
    },
    "title": "",
    "url": ""
  }
]`
	assert.Equal(t,expected, result)
}

func TestCommonCmd_Filter_FilterOnAuthor(t *testing.T) {
	prs := []pr.PullRequest{
		{Number: 1, Author: pr.Author{Login: "author1"}},
		{Number: 2, Author: pr.Author{Login: "author2"}},
		{Number: 3, Author: pr.Author{Login: "author3"}},
	}

	cmd := cmd2.CommonCmd{Query: "[?author.login=='author1']"}
	result, err := cmd.Filter(prs)
	assert.NoError(t, err)
	expected := `[
  {
    "author": {
      "login": "author1"
    },
    "closed": false,
    "commits": {
      "nodes": null
    },
    "createdAt": "0001-01-01T00:00:00Z",
    "labels": {
      "nodes": null
    },
    "mergeable": "",
    "number": 1,
    "repository": {
      "nameWithOwner": ""
    },
    "title": "",
    "url": ""
  }
]`
	assert.Equal(t,expected, result)
}

func TestCommonCmd_Filter_Rewrite(t *testing.T) {
	prs := []pr.PullRequest{
		{Number: 1, Author: pr.Author{Login: "author1"}},
		{Number: 2, Author: pr.Author{Login: "author2"}},
		{Number: 3, Author: pr.Author{Login: "author3"}},
	}

	cmd := cmd2.CommonCmd{Query: "[].{url:url, number:number, author:author.login}"}
	result, err := cmd.Filter(prs)
	assert.NoError(t, err)
	expected := `[
  {
    "author": "author1",
    "number": 1,
    "url": ""
  },
  {
    "author": "author2",
    "number": 2,
    "url": ""
  },
  {
    "author": "author3",
    "number": 3,
    "url": ""
  }
]`
	assert.Equal(t,expected, result)
}