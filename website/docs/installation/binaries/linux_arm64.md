---
sidebar_position: 3
---

# Linux (arm64)

Install **slotalk** on linux (arm64).

## Requirements

To install the tool using this method you'll require:

* cURL
* tar
* wget (optional)

Present on your host machine.

## Installation

Simply run, in your terminal:

```shell
curl -s -L https://github.com/tfadeyi/slotalk/releases/latest/download/slotalk-linux-arm64.tar.gz | tar xzv
# might require sudo
mv slotalk-linux-arm64/slotalk /usr/local/bin
```

This will install the latest slotalk binary under the path: `/usr/local/bin/slotalk`.

> You can install different versions by setting the tag to the target version: https://github.com/tfadeyi/slotalk/releases/v0.1.0-alpha.1/download/slotalk-linux-arm64.tar.gz

## Verify Installation

```shell
slotalk --help
```

The binary should return something similar to:

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

## Uninstall 😢

To uninstall the tool you can simply delete the binary from the following directory.

```shell
# might require sudo
rm /usr/local/bin/slotalk
```
