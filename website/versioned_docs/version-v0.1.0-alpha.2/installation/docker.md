---
sidebar_position: 3
---

# Docker 🐋

You can run **slotalk** into your machine using `docker`.

This method won't install the binary in the host machine per se, but it will run the binary in a container.

## Requirements

To install the tool using this method you'll require:

* Docker

Present on your host machine.

## Running slotalk with Nix

Simply run, in your terminal:

```shell
docker run docker ghcr.io/tfadeyi/slotalk:latest
```

This gives you an ephemeral way to run slotalk, so you can try out the tool without installing it.

## Try it!

```shell
docker run docker ghcr.io/tfadeyi/slotalk:latest --help
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
