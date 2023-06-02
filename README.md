<div align="center">

# Slotalk

[![Continuous Integration](https://img.shields.io/github/actions/workflow/status/tfadeyi/sloth-simple-comments/ci.yml?branch=main&style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-yellowgreen.svg?style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments/blob/main/LICENSE)
[![Language](https://img.shields.io/github/go-mod/go-version/tfadeyi/sloth-simple-comments?style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments)
[![GitHub release](https://img.shields.io/github/v/release/tfadeyi/sloth-simple-comments?color=green&style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments/releases)
[![Code size](https://img.shields.io/github/languages/code-size/tfadeyi/sloth-simple-comments?color=orange&style=for-the-badge)](https://github.com/tfadeyi/sloth-simple-comments)
[![Go Report Card](https://goreportcard.com/badge/github.com/tfadeyi/sloth-simple-comments?style=for-the-badge)](https://goreportcard.com/report/github.com/tfadeyi/sloth-simple-comments)
</div>


Slotalk is a CLI tool to help developers embed [Sloth](https://sloth.dev/) SLO/SLI [definitions](https://github.com/slok/sloth/tree/main/pkg/prometheus/api/v1) into their code base, without defining a separate
YAML file close to where the metrics used in the actual SLIs are defined.

The tool takes inspiration from [Swaggo](https://github.com/swaggo/swag), a tool to generate Swagger docs from Go code,
as such it uses a similar pattern when it comes to embedding the Sloth definitions within the code.

## Motivation

* **Experimentation**, this was the main motivation behind development, testing libraries like: [go/ast](https://pkg.go.dev/go/ast), [wazero](https://github.com/tetratelabs/wazero), [participle](https://github.com/alecthomas/participle).
* **Developer experience**, but also finding ways to improve developer experience when it comes to more platform engineering concepts like SLIs and SLOs. I want to see if moving these concepts closer to devs,
would make them less of an afterthought.
* **More Experimentation** Many of the cloud native tools I've seen have been very targeted towards DevOps/SecOps and Platform Engineering personas,
so I wanted try my hand on building something for developers.
* Trying ways to avoid writing YAML...

## Prerequisites

* [Sloth CLI](https://github.com/slok/sloth)
* [Go](https://go.dev/doc/install) (optional)
* [Nix](https://zero-to-nix.com/start/install) (optional)

## TL;DR

1. Add comments to your source code. See [Declarative Comments](#Declarative-Comments).

2. Slotalk Installation

   **Go install**

   If go is present on the host machine, you can just download the required binaries.
   
   ```shell
   # install the latest version of slotalk
   go install github.com/tfadeyi/slotalk@latest
   # (OPTIONAL) install the latest version of sloth
   go install github.com/slok/sloth/cmd/sloth@latest
   ```

   **Nix**

   If nix is present on the host machine, you can run the tool in the development shell
   
   ```shell
   # creates a nix development shell with slotalk and sloth
   nix develop github:tfadeyi/slotalk
   ```

   **Pre-released binaries**

   Download a pre-compiled binary from the release page.
   ```shell
     curl -LJO https://github.com/tfadeyi/slotalk/releases/download/v0.0.2/slotalk-linux-amd64.tar.gz && \
     tar -xzvf slotalk-linux-amd64.tar.gz && \
     cd slotalk-linux-amd64
   ```

3. Run `slotalk` init in the project's root. This will parse your comments and print to standard out.
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

4. Run `sloth` generate command using the `sloth_defs.yaml` as input.
    ```shell
    sloth generate -i sloth_defs.yaml
    ```

This will return the Prometheus alerting rules for the given SLOs.

## Declarative Comments

The definitions are added through declarative comments, as shown below.

```go
// @sloth.slo service chatgpt
// @sloth.slo name chat-gpt-availability
// @sloth.slo objective 95.0
// @sloth.sli error_query sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[{{.window}}])) OR on() vector(0)
// @sloth.sli total_query sum(rate(tenant_login_operations_total{client="chat-gpt"}[{{.window}}]))
// @sloth.slo description 95% of logins to the chat-gpt app should be successful.
// @sloth.alerting name ChatGPTAvailability
```

### Service definitions
| annotation | description                                                     | example                                         |
|------------|-----------------------------------------------------------------|-------------------------------------------------|
| service    | **Required**. The name of the service the definitions refer to. | @sloth service chat-gpt                         |
| version    | **Required**. The version of the Sloth specification.           | @sloth version prometheus/v1                    |
| labels     | The labels associated to the Sloth service.                     | @sloth labels foo bar \n @sloth labels test slo |

### SLO definitions

| annotation  | description                                                 | example                                                                                                        |
|-------------|-------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------|
| name        | **Required**. The name of the SLO.                          | @sloth.slo name availability                                                                                   |
| objective   | **Required**. The SLO objective in floating point notation. | @sloth.slo objective 95.0                                                                                      |
| description |                                                             | @sloth.slo description 95% of logins to the chat-gpt app should be successful annotations. (can be multilined) |
| labels      |                                                             |                                                                                                                |

### Alerting definitions

| annotation  | description | example                                                                                                                                               |
|-------------|-------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| name        |             | @sloth.alerting name ChatGPTAvailability                                                                                                              |
| labels      |             | @sloth.alerting labels severity critical (new labels should be in new line)                                                                           |
| annotations |             | @sloth.alerting annotations runbook: "https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubeapilatencyhigh" |

### Page Alerting definitions

| annotation  | description | example                                           |
|-------------|-------------|---------------------------------------------------|
| name        |             | @sloth.alerting.page name pageAlerting            |
| labels      |             | @sloth.alerting.page labels severity critical     |
| annotations |             | @sloth.alerting.page annotations tier application |

### Ticket Alerting definitions

| annotation  | description | example                                           |
|-------------|-------------|---------------------------------------------------|
| name        |             | @sloth.alerting.page name ticketAlerting          |
| labels      |             | @sloth.alerting.page labels severity critical     |
| annotations |             | @sloth.alerting.page annotations tier application |

### CLI usage

```text
Usage:
  sli-app init [flags]

Flags:
      --dirs strings     Comma separated list of directories to be parses by the tool (default [/home/jetstack-oluwole/go/src/github.com/tfadeyi/slotalk])
  -f, --file string      Source file to parse.
      --format strings   Output format (yaml,json). (default [yaml])
  -h, --help             help for init
      --lang string      Language of the source files. (go, wasm) (default "go")
```

### Examples

#### Basic usage

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