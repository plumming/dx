## dx get issues

Gets your open issues

```
dx get issues [flags]
```

### Examples

```
Get a list of open issues:

  dx get issues

Get a list of your issues:

  dx get issues --me

Get a list of issues requiring review:

  dx get issues --review

Get a list of issues with a custom query:

  dx get issues --raw is:private


```

### Options

```
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

