## dx get prs

Gets your open prs

```
dx get prs [flags]
```

### Examples

```
Get a list of open PRs:

  dx get prs

Get a list of your PRs:

  dx get prs --me

Get a list of PRs requiring review:

  dx get prs --review

Get a list of PRs with a custom query:

  dx get prs --raw is:private


```

### Options

```
  -c, --copy           Output is copy and pasteable
  -m, --me             Show all PRs that are created by the author
  -q, --query string   JMESPath query filter
      --quiet          Hide the column headings
      --raw string     Additional raw search parameters to use when querying
      --review         Show PRs that are ready for review
      --show-bots      Show bot account PRs (default: false)
      --show-hidden    Show PRs that are filtered by hidden labels (default: false)
```

### Options inherited from parent commands

```
      --help   Show help for command
```

### SEE ALSO

* [dx get](dx_get.md)	 - 

