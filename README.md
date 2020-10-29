# chilly

[![codecov](https://codecov.io/gh/plumming/chilly/branch/master/graph/badge.svg)](https://codecov.io/gh/plumming/chilly)
[![Go Report Card](https://goreportcard.com/badge/github.com/plumming/chilly)](https://goreportcard.com/report/github.com/plumming/chilly)
![golangci-lint](https://github.com/plumming/chilly/workflows/golangci-lint/badge.svg)
![Check PR can be merged](https://github.com/plumming/chilly/workflows/Check%20PR%20can%20be%20merged/badge.svg)
![Set Label](https://github.com/plumming/chilly/workflows/Set%20Label/badge.svg)
![Build and test Go](https://github.com/plumming/chilly/workflows/Build%20and%20test%20Go/badge.svg)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)

# Installation via Homebrew

To install `chilly` using homebrew, run the following:

```
brew tap plumming/chilly
brew install chilly
```

## Basic Usage

### View a list of Open PRs

```
chilly get prs
```

### Configure the list of repos to watch

```
chilly edit config
```

Add more entries under the `repos:` block e.g.

```
repos:
- org/repo1
- org/repo2
```

### Exclude PRs based on labels

```
chilly edit config
```

```
hiddenLabels:
- wip
- do not merge
```
