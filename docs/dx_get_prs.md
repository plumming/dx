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

Get a list of PRs including drafts:

  dx get prs --show-drafts

Get a list of PRs excluding approved ones:

  dx get prs --hide-approved

Get a list of PRs excluding queued ones:

  dx get prs --hide-queued


```

### Options

```
  -c, --copy[=true|false]            Output is copy and pasteable
      --hide-approved[=true|false]   Hide approved PRs (default: false)
      --hide-queued[=true|false]     Hide queued PRs (default: false)
  -m, --me[=true|false]              Show all PRs that are created by the author
  -q, --query string                 JMESPath query filter
      --quiet[=true|false]           Hide the column headings
      --raw string                   Additional raw search parameters to use when querying
      --review[=true|false]          Show PRs that are ready for review
      --show-bots[=true|false]       Show bot account PRs (default: false)
      --show-drafts[=true|false]     Show draft PRs (default: false)
      --show-hidden[=true|false]     Show PRs that are filtered by hidden labels (default: false)
```

### Options inherited from parent commands

```
      --help[=true|false]   Show help for command
```

### SEE ALSO

* [dx get](dx_get.md)	 - 

