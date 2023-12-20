## dx prompt

Configure a new command prompt

```
dx prompt [flags]
```

### Examples

```

		# Generate the current prompt
		dx prompt

		# Enable the prompt for bash
		PS1="[\u@\h \W \$(dx prompt)]\$ "

		# Enable the prompt for zsh
		PROMPT='$(dx prompt)'$PROMPT
	
```

### Options

```
      --context-color stringArray     The color for the Kubernetes context (default [cyan])
  -d, --divider string                The divider between the team and environment for the prompt (default ":")
  -i, --icon                          Uses an icon for the label in the prompt
  -l, --label string                  The label for the prompt (default "k8s")
      --label-color stringArray       The color for the label (default [blue])
      --namespace-color stringArray   The color for the namespace (default [green])
      --no-label                      Disables the use of the label in the prompt
  -p, --prefix string                 The prefix text for the prompt
  -s, --separator string              The separator between the label and the rest of the prompt (default ":")
  -x, --suffix string                 The suffix text for the prompt (default ">")
```

### Options inherited from parent commands

```
      --help   Show help for command
```

### SEE ALSO

* [dx](dx.md)	 - Plumming dx

