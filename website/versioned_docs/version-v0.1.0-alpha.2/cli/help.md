---
sidebar_position: 2
---

# Help Command

The `help` command will return the list of available CLI commands.

```shell
slotalk help
```

Output:

```shell
Generate Sloth SLO/SLI definitions from code annotations.

Usage:
  slotalk [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Init generates the Sloth definition specification from source code comments.
  version     Returns the binary build information.

Flags:
  -h, --help               help for slotalk
      --log-level string   Only log messages with the given severity or above. One of: [none, debug, info, warn], errors will always be printed (default "info")

Use "slotalk [command] --help" for more information about a command.
```

To view information about a specific command the `--help` flag can be passed, this will return
the list child commands and flags.

```shell
slotalk version --help
```

Output:

```shell
Returns the binary build information.

Usage:
  slotalk version [flags]

Flags:
  -h, --help   help for version

Global Flags:
      --log-level string   Only log messages with the given severity or above. One of: [none, debug, info, warn], errors will always be printed (default "info")
```
