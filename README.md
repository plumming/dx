# DX

[![codecov](https://codecov.io/gh/plumming/dx/branch/master/graph/badge.svg)](https://codecov.io/gh/plumming/dx)
[![Go Report Card](https://goreportcard.com/badge/github.com/plumming/dx)](https://goreportcard.com/report/github.com/plumming/dx)
![golangci-lint](https://github.com/plumming/dx/workflows/golangci-lint/badge.svg)
![Check PR can be merged](https://github.com/plumming/dx/workflows/Check%20PR%20can%20be%20merged/badge.svg)
![Set Label](https://github.com/plumming/dx/workflows/Set%20Label/badge.svg)
![Build and test Go](https://github.com/plumming/dx/workflows/Build%20and%20test%20Go/badge.svg)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Downloads](https://img.shields.io/github/downloads/plumming/dx/total.svg)]()

# Installation via Homebrew

To install `dx` using homebrew, run the following:

```
brew tap plumming/homebrew-tap
brew install dx
```

## Basic Usage

### View a list of Open PRs

```
dx get prs
```

### Configure the list of repos to watch

```
dx edit config
```

Add more entries under the `repos:` block e.g.

```
repos:
- org/repo1
- org/repo2
```

### Exclude PRs based on labels

```
dx edit config
```

```
hiddenLabels:
- wip
- do not merge
```

### Exclude PRs older than X days

```
dx edit config
```

```
maxAgeOfPRs: 180
```

or for all PRs 

```
maxAgeOfPRs: -1
```

### Rebase the local repository

This will rebase the local repository against the remote called 'upstream'.

```
dx rebase
```


### Commands

- [dx](./docs/dx.md)
- [dx context](./docs/dx_context.md)
- [dx delete](./docs/dx_delete.md)
- [dx delete repos](./docs/dx_delete_repos.md)
- [dx edit](./docs/dx_edit.md)
- [dx edit config](./docs/dx_edit_config.md)
- [dx get](./docs/dx_get.md)
- [dx get prs](./docs/dx_get_prs.md)
- [dx get repos](./docs/dx_get_repos.md)
- [dx namespace](./docs/dx_namespace.md)
- [dx rebase](./docs/dx_rebase.md)
- [dx upgrade](./docs/dx_upgrade.md)
- [dx upgrade cli](./docs/dx_upgrade_cli.md)
