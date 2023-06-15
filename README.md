<div align="center">

# Slotalk

<p align="center">
<a href="#Motivation">Motivation</a> | <a href="#Prerequisites">Prerequisites</a> | <a href="#Try-it!">Try it!</a> | <a href="#Installation">Installation</a> | <a href="#Get-Started">Get Started</a> | <a href="#Examples">Examples</a>
</p>

[![Nix Devshell](https://img.shields.io/badge/nix-devshell-blue?logo=NixOS&style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments#Nix)
[![Continuous Integration](https://img.shields.io/github/actions/workflow/status/tfadeyi/sloth-simple-comments/ci.yml?branch=main&style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-yellowgreen.svg?style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments/blob/main/LICENSE)
[![Language](https://img.shields.io/github/go-mod/go-version/tfadeyi/sloth-simple-comments?style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments)
[![GitHub release](https://img.shields.io/github/v/release/tfadeyi/sloth-simple-comments?color=green&style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments/releases)
[![Code size](https://img.shields.io/github/languages/code-size/tfadeyi/sloth-simple-comments?color=orange&style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments)
[![Go Report Card](https://goreportcard.com/badge/github.com/tfadeyi/sloth-simple-comments?style=for-the-badge)](https://goreportcard.com/report/github.com/tfadeyi/sloth-simple-comments)
</div>

> ⚠ The tool is still not ready for real production use yet.

Slotalk is a CLI tool that allows developers to embed [Sloth](https://github.com/slok/sloth) SLO/SLI [specifications](https://github.com/slok/sloth/blob/main/pkg/prometheus/api/v1/README.md) as in-code annotations rather than a YAML file.

<p align="center">
<img src="./static/images/simple-usage-example.gif">
</p>

Similar to how [Swaggo](https://github.com/swaggo/swag) does for Swagger docs,  Slotalk moves the SLO/SLI specification closer to where its relevant Prometheus metric was defined.

Slotalk can be used in tandem with the [Sloth CLI](https://github.com/slok/sloth#getting-started) to generate Prometheus alerts groups from the in-code annotations, which can be used in any Prometheus/Grafana monitoring system to keep track of the service's SLOs. See examples below.



## Motivation

* **Experimentation**, this was the main motivation behind development, testing libraries like: [go/ast](https://pkg.go.dev/go/ast), [wazero](https://github.com/tetratelabs/wazero), [participle](https://github.com/alecthomas/participle).
* **Developer experience**, finding ways to improve developer experience when it comes to more platform engineering concepts like SLIs and SLOs. I want to see if moving these concepts closer to devs,
  would make them less of an afterthought.
* **More Experimentation**, many of the cloud native tools I've seen, and I've worked on have been very targeted towards DevOps/SecOps and Platform Engineering personas,
  so I wanted try my hand on building something for developers.
* Trying ways to avoid writing YAML...

## Prerequisites

* [Sloth CLI](https://github.com/slok/sloth) (optional)
* [Go](https://go.dev/doc/install)
* [Nix](https://zero-to-nix.com/start/install) (optional)

## Try it!

### Nix
Generate Prometheus SLO alert rules from an example [metrics.go](https://gist.githubusercontent.com/tfadeyi/df60aebd858d1c76428c045d4df7b114/raw/dfb96773dfb64086280845b9a0776012cbd7d26b/metrics.go).
   ```shell
   # creates a nix demo shell with slotalk and sloth. just follow the shell instructions
   nix develop github:tfadeyi/slotalk#demo
   ```

### Source
Generate Prometheus SLO alert rules from an example [metrics.go](https://gist.githubusercontent.com/tfadeyi/df60aebd858d1c76428c045d4df7b114/raw/dfb96773dfb64086280845b9a0776012cbd7d26b/metrics.go).

1. Install Slotalk
   ```shell
   # install the latest version of slotalk
   go install github.com/tfadeyi/slotalk@latest
   # install the latest version of sloth
   go install github.com/slok/sloth/cmd/sloth@latest
   ```
2. Run `slotalk` and `sloth` to generate Prometheus alert rules from code annotations.
    ```shell
    curl https://gist.githubusercontent.com/tfadeyi/df60aebd858d1c76428c045d4df7b114/raw/dfb96773dfb64086280845b9a0776012cbd7d26b/metrics.go > metrics.go
    cat metrics.go | slotalk init -f - > ./sloth_defs.yaml
    sloth generate -i ./sloth_defs.yaml -o ./rules.yml
    ```

You now should have Prometheus alerting rules that can be added to your Prometheus configuration.
<details>
  <summary>Prometheus configuration</summary>

```yaml
# my global config
global:
  scrape_interval: 5s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 5s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
 - "rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "exporter"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9301"]
```

</details>

## Installation

<details>
     <summary><strong>Go Install</strong></summary>

   ```shell
   # install the latest version of slotalk
   go install github.com/tfadeyi/slotalk@latest
   ```

</details>

<details>
     <summary><strong>Nix</strong></summary>

```shell
nix run github:tfadeyi/slotalk
```

</details>

<details>
     <summary><strong>Docker</strong></summary>

```shell
docker pull ghcr.io/tfadeyi/slotalk:latest
```

</details>

<details>
     <summary><strong>Install Script</strong></summary>

```shell
curl -sfL https://raw.githubusercontent.com/tfadeyi/slotalk/main/install.sh | sh -
```

</details>

<strong>Pre-released binaries</strong>

<details>
     <summary><strong>Linux (x84_64)</strong></summary>

```shell
  curl -s -L https://github.com/tfadeyi/slotalk/releases/latest/download/slotalk-linux-amd64.tar.gz | tar xzv
```

</details>

<details>
     <summary><strong>Linux (arm)</strong></summary>

```shell
  curl -s -L https://github.com/tfadeyi/slotalk/releases/latest/download/slotalk-linux-arm.tar.gz | tar xzv
```

</details>

<details>
     <summary><strong>Linux (arm64)</strong></summary>

```shell
  curl -s -L https://github.com/tfadeyi/slotalk/releases/latest/download/slotalk-linux-arm64.tar.gz | tar xzv
```

</details>

<details>
     <summary><strong>MacOS (amd64)</strong></summary>

```shell
  curl -s -L https://github.com/tfadeyi/slotalk/releases/latest/download/slotalk-darwin-amd64.tar.gz | tar xzv
```

</details>

<details>
     <summary><strong>MacOS (Apple Silicon)</strong></summary>

```shell
  curl -s -L https://github.com/tfadeyi/slotalk/releases/latest/download/slotalk-darwin-arm64.tar.gz | tar xzv
```

</details>

## Get Started

1. Add comments to your source code. See [Declarative Comments](#Declarative-Comments).
2. Run `slotalk` init in the project's root. This will parse your source code annotations and print the sloth definitions to standard out.
    ```shell
    ./slotalk init
    ```

   You can also specify the specific file to parse by using the `-f` flag.

    ```shell
    ./slotalk init -f metrics.go
    ```

   Another way would be to pass the input file through pipe.

    ```shell
    cat metrics.go | ./slotalk init -f -
    ```

## CLI usage

```text
Usage:
  slotalk init [flags]

Flags:
      --dirs strings     Comma separated list of directories to be parses by the tool (default [/home/jetstack-oluwole/go/src/github.com/tfadeyi/slotalk])
  -f, --file string      Source file to parse.
      --format strings   Output format (yaml,json). (default [yaml])
  -h, --help             help for init
      --lang string      Language of the source files. (go) (default "go")

Global Flags:
      --log-level string   Only log messages with the given severity or above. One of: [none, debug, info, warn], errors will always be printed (default "info")
```

## Declarative Comments (Sloth)

The Sloth definitions are added through declarative comments, as shown below.

```go
// @sloth service chatgpt
// @sloth.slo name chat-gpt-availability
// @sloth.slo objective 95.0
// @sloth.sli error_query sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[{{.window}}])) OR on() vector(0)
// @sloth.sli total_query sum(rate(tenant_login_operations_total{client="chat-gpt"}[{{.window}}]))
// @sloth.slo description 95% of logins to the chat-gpt app should be successful.
// @sloth.alerting name ChatGPTAvailability
```

### Service definitions
| Annotation | Description                                                     | Example                                      |
|------------|-----------------------------------------------------------------|----------------------------------------------|
| service    | **Required**. The name of the service the definitions refer to. | @sloth service chat-gpt                      |
| version    | The version of the Sloth specification.                         | @sloth version prometheus/v1                 |
| labels     | The labels associated to the Sloth service.                     | @sloth labels foo bar @sloth labels test slo |

### SLO definitions

| Annotation  | Description                                                                                                                                                         | Example                                                                                                        |
|-------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------|
| name        | **Required**. The name of the SLO.                                                                                                                                  | @sloth.slo name availability                                                                                   |
| objective   | **Required**. The SLO Objective is target of the SLO the percentage (0, 100] (e.g 99.9).                                                                            | @sloth.slo objective 95.0                                                                                      |
| description | Description is the description of the SLO.                                                                                                                          | @sloth.slo description 95% of logins to the chat-gpt app should be successful annotations. (can be multilined) |
| labels      | Labels are the Prometheus labels that will have all the recording and alerting rules for this specific SLO. These labels are merged with the previous level labels. | @sloth.slo labels foo bar @sloth labels test slo                                                               |

### Alerting definitions

| Annotation  | Description                                                                                     | Example                                                                                                                                               |
|-------------|-------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| name        | **Required**. Name is the name used by the alerts generated for this SLO.                       | @sloth.alerting name ChatGPTAvailability                                                                                                              |
| labels      | Labels are the Prometheus labels that will have all the alerts generated by this SLO.           | @sloth.alerting labels severity critical                                                                                                              |
| annotations | Annotations are the Prometheus annotations that will have all the alerts generated by this SLO. | @sloth.alerting annotations runbook: "https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubeapilatencyhigh" |

### Page Alerting definitions

| Annotation  | Description                                                                                                                           | Example                                           |
|-------------|---------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------|
| labels      | Labels are the Prometheus labels for the specific alert. For example can be useful to route the Page alert to specific Slack channel. | @sloth.alerting.page labels severity critical     |
| annotations | Annotations are the Prometheus annotations for the specific alert.                                                                    | @sloth.alerting.page annotations tier application |

### Ticket Alerting definitions

| Annotation  | Description                                                                                                                           | Example                                             |
|-------------|---------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|
| labels      | Labels are the Prometheus labels for the specific alert. For example can be useful to route the Page alert to specific Slack channel. | @sloth.alerting.ticket labels severity critical     |
| annotations | Annotations are the Prometheus annotations for the specific alert.                                                                    | @sloth.alerting.ticket annotations tier application |

## Examples

### Basic usage - Generate Sloth definitions using go:generate
The following example shows how to use `go:generate` to generate Sloth definitions from in code annotations.

**metrics.go**
```go
    // @sloth service chatgpt
    var (
        // @sloth.slo name chat-gpt-availability
        // @sloth.slo objective 95.0
        // @sloth.sli error_query sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[{{.window}}])) OR on() vector(0)
        // @sloth.sli total_query sum(rate(tenant_login_operations_total{client="chat-gpt"}[{{.window}}]))
        // @sloth.slo description 95% of logins to the chat-gpt app should be successful.
        // @sloth.alerting name ChatGPTAvailability
        metricGaugeCertInventoryProcessingMessages = prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: "chatgpt",
                Subsystem: "auth0",
                Name:      "tenant_login_operations_total",
            })
        tenantFailedLogins = prometheus.NewCounter(
            prometheus.CounterOpts{
            Namespace: "chatgpt",
            Subsystem: "auth0",
            Name:      "tenant_failed_login_operations_total",
        })
    )
```

**main.go**
```go
//go:generate slotalk init

package main

import (
)

// @sloth service chatgpt
func main() {
}
```

Running go generate, will allow the `slotalk` to walk through the different packages parsing the in code annotations and
generate Sloth definitions.

```shell
go generate ./...
```

<details>
  <summary>Result Sloth Definitions.</summary>

```yaml
# Code generated by slotalk: https://github.com/tfadeyi/slotalk.
# DO NOT EDIT.
version: prometheus/v1
service: chatgpt
slos:
    - name: chat-gpt-availability
      description: 95% of logins to the chat-gpt app should be successful.
      objective: 95
      sli:
        events:
            error_query: sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[{{.window}}])) OR on() vector(0)
            total_query: sum(rate(tenant_login_operations_total{client="chat-gpt"}[{{.window}}]))
      alerting:
        name: ChatGPTAvailability
```

</details>

### Basic usage - Generate Prometheus alert groups from code annotations

This example shows how sloth's comments can be added next to the prometheus metrics defined in a `metrics.go` file. 
```go
    // @sloth service chatgpt
    
    var (
        // @sloth.slo name chat-gpt-availability
        // @sloth.slo objective 95.0
        // @sloth.sli error_query sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[{{.window}}])) OR on() vector(0)
        // @sloth.sli total_query sum(rate(tenant_login_operations_total{client="chat-gpt"}[{{.window}}]))
        // @sloth.slo description 95% of logins to the chat-gpt app should be successful.
        // @sloth.alerting name ChatGPTAvailability
        metricGaugeCertInventoryProcessingMessages = prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: "chatgpt",
                Subsystem: "auth0",
                Name:      "tenant_login_operations_total",
            })
        tenantFailedLogins = prometheus.NewCounter(
            prometheus.CounterOpts{
            Namespace: "chatgpt",
            Subsystem: "auth0",
            Name:      "tenant_failed_login_operations_total",
        })
    )
```

Now running the following command from the root of the project.

```shell
./slotalk init
```

This will generate the following sloth definitions being outputted to standard out.

```yaml
version: prometheus/v1
service: "chatgpt"
slos:
    - name: chat-gpt-availability
      description: 95% of logins to the chat-gpt app should be successful.
      objective: 95
      sli:
        events:
            error_query: sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[{{.window}}])) OR on() vector(0)
            total_query: sum(rate(tenant_login_operations_total{client="chat-gpt"}[{{.window}}]))
      alerting:
        name: "ChatGPTAvailability"
```

This specification can then be passed to the Sloth CLI to generate Prometheus alerting groups.

```shell
./slotalk init > sloth_defs.yaml && sloth generate -i sloth_defs.yaml
```

<details>
  <summary>Resulting alert groups.</summary>

```yaml
# Code generated by Sloth (v0.11.0): https://github.com/slok/sloth.
# DO NOT EDIT.

groups:
- name: sloth-slo-sli-recordings-foo-chat-gpt-availability
  rules:
  - record: slo:sli_error:ratio_rate5m
    expr: |
      (sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[5m])) OR on() vector(0))
      /
      (sum(rate(tenant_login_operations_total{client="chat-gpt"}[5m])))
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_window: 5m
  - record: slo:sli_error:ratio_rate30m
    expr: |
      (sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[30m])) OR on() vector(0))
      /
      (sum(rate(tenant_login_operations_total{client="chat-gpt"}[30m])))
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_window: 30m
  - record: slo:sli_error:ratio_rate1h
    expr: |
      (sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[1h])) OR on() vector(0))
      /
      (sum(rate(tenant_login_operations_total{client="chat-gpt"}[1h])))
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_window: 1h
  - record: slo:sli_error:ratio_rate2h
    expr: |
      (sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[2h])) OR on() vector(0))
      /
      (sum(rate(tenant_login_operations_total{client="chat-gpt"}[2h])))
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_window: 2h
  - record: slo:sli_error:ratio_rate6h
    expr: |
      (sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[6h])) OR on() vector(0))
      /
      (sum(rate(tenant_login_operations_total{client="chat-gpt"}[6h])))
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_window: 6h
  - record: slo:sli_error:ratio_rate1d
    expr: |
      (sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[1d])) OR on() vector(0))
      /
      (sum(rate(tenant_login_operations_total{client="chat-gpt"}[1d])))
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_window: 1d
  - record: slo:sli_error:ratio_rate3d
    expr: |
      (sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[3d])) OR on() vector(0))
      /
      (sum(rate(tenant_login_operations_total{client="chat-gpt"}[3d])))
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_window: 3d
  - record: slo:sli_error:ratio_rate30d
    expr: |
      sum_over_time(slo:sli_error:ratio_rate5m{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"}[30d])
      / ignoring (sloth_window)
      count_over_time(slo:sli_error:ratio_rate5m{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"}[30d])
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_window: 30d
- name: sloth-slo-meta-recordings-foo-chat-gpt-availability
  rules:
  - record: slo:objective:ratio
    expr: vector(0.95)
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
  - record: slo:error_budget:ratio
    expr: vector(1-0.95)
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
  - record: slo:time_period:days
    expr: vector(30)
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
  - record: slo:current_burn_rate:ratio
    expr: |
      slo:sli_error:ratio_rate5m{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"}
      / on(sloth_id, sloth_slo, sloth_service) group_left
      slo:error_budget:ratio{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"}
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
  - record: slo:period_burn_rate:ratio
    expr: |
      slo:sli_error:ratio_rate30d{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"}
      / on(sloth_id, sloth_slo, sloth_service) group_left
      slo:error_budget:ratio{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"}
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
  - record: slo:period_error_budget_remaining:ratio
    expr: 1 - slo:period_burn_rate:ratio{sloth_id="foo-chat-gpt-availability", sloth_service="foo",
      sloth_slo="chat-gpt-availability"}
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_service: foo
      sloth_slo: chat-gpt-availability
  - record: sloth_slo_info
    expr: vector(1)
    labels:
      foo: bar
      sloth_id: foo-chat-gpt-availability
      sloth_mode: cli-gen-prom
      sloth_objective: "95"
      sloth_service: foo
      sloth_slo: chat-gpt-availability
      sloth_spec: prometheus/v1
      sloth_version: v0.11.0
- name: sloth-slo-alerts-foo-chat-gpt-availability
  rules:
  - alert: K8sApiserverAvailabilityAlert
    expr: |
      (
          max(slo:sli_error:ratio_rate5m{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"} > (14.4 * 0.05)) without (sloth_window)
          and
          max(slo:sli_error:ratio_rate1h{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"} > (14.4 * 0.05)) without (sloth_window)
      )
      or
      (
          max(slo:sli_error:ratio_rate30m{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"} > (6 * 0.05)) without (sloth_window)
          and
          max(slo:sli_error:ratio_rate6h{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"} > (6 * 0.05)) without (sloth_window)
      )
    labels:
      sloth_severity: page
    annotations:
      summary: '{{$labels.sloth_service}} {{$labels.sloth_slo}} SLO error budget burn
        rate is over expected.'
      title: (page) {{$labels.sloth_service}} {{$labels.sloth_slo}} SLO error budget
        burn rate is too fast.
  - alert: K8sApiserverAvailabilityAlert
    expr: |
      (
          max(slo:sli_error:ratio_rate2h{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"} > (3 * 0.05)) without (sloth_window)
          and
          max(slo:sli_error:ratio_rate1d{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"} > (3 * 0.05)) without (sloth_window)
      )
      or
      (
          max(slo:sli_error:ratio_rate6h{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"} > (1 * 0.05)) without (sloth_window)
          and
          max(slo:sli_error:ratio_rate3d{sloth_id="foo-chat-gpt-availability", sloth_service="foo", sloth_slo="chat-gpt-availability"} > (1 * 0.05)) without (sloth_window)
      )
    labels:
      sloth_severity: ticket
    annotations:
      summary: '{{$labels.sloth_service}} {{$labels.sloth_slo}} SLO error budget burn
        rate is over expected.'
      title: (ticket) {{$labels.sloth_service}} {{$labels.sloth_slo}} SLO error budget
        burn rate is too fast.
```

</details>


## License
MIT, see [LICENSE.md](./LICENSE).
