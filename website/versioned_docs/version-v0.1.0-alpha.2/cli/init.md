---
sidebar_position: 1
---

# Init Command

The `init` command parses the target file or directories for in code annotations describing a service's SLOs,
this will then output, to stdout or file, the generated SLOs definitions specification in either yaml or json.

```shell
slotalk init
```

## Usage

```shell
The init command parses files in the target directory for comments using the @sloth tags

Usage:
  slotalk init [flags]

Flags:
      --dirs strings     Comma separated list of directories to be parses by the tool (default [/home/oluwole/go/src/github.com/tfadeyi/slotalk])
  -f, --file string      Source file to parse.
      --format strings   Output format (yaml,json). (default [yaml])
  -h, --help             help for init
      --lang string      Language of the source files. (go) (default "go")
      --service string   Outputs only the selected service specification from the resulting parsed service specification.
      --to-file          Print the generated specifications to file, under ./slo_definitions.

Global Flags:
      --log-level string   Only log messages with the given severity or above. One of: [none, debug, info, warn], errors will always be printed (default "info")
```

By default, the command outputs the generated specification to stdout 

## Flags

### `--dirs`
Is the comma separated list of directories the tool will parse to find SLO annotations. If no arguments are passed
to the flag, it will default to parse the current working directory.

**Example:**

```shell
   
slotalk init --dirs
    
```

:::danger Take care

The `--dirs` flag cannot be used together with the `--file` flag.

:::


### `--file`
Is the name of the target file that the tool will parse for SLO annotations.

**Examples:**

```shell
slotalk init -f metrics.go
cat metrics.go | ./slotalk init -f -
```

### `--format`
Is the format of the output of the SLO specification, Yaml or JSON are the currently supported formats.

### `--lang`
Is the language of the target source code expected by the parser, ie: GoLang.

### `--service`
In the case multiple services SLO specifications are found when parsing the target directory, the `--service`
will allow to specify which service to output.

**Example:**
```shell
slotalk init --service chatgpt
```

### `--to-file`
Set the tool to parse and output SLO specification to file rather than stdout.