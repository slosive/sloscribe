---
sidebar_position: 2
---

# SLO Annotations

Define details about the SLOs.

The following is an example of what the in-code annotations defining a simple SLO:

```go
// @sloth.slo name chat-gpt-availability
// @sloth.slo objective 95.0
// @sloth.slo description 95% of logins to the chat-gpt app
// should be successful annotations
// @sloth.labels severity critical
```

The above annotations will generate the following sloth specification:

```yaml
version: prometheus/v1
service: ...
labels: ...
slos:
  - name: chat-gpt-availability
    description: 95% of logins to the chat-gpt app should be successful.
    objective: 95
    labels:
      severity: critical
```

## Table of Annotations
| Annotation  | Description                                                                                                                                                         | Example                                                                                                        |
|-------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------|
| name        | **Required**. The name of the SLO.                                                                                                                                  | @sloth.slo name availability                                                                                   |
| objective   | **Required**. The SLO Objective is target of the SLO the percentage (0, 100] (e.g 99.9).                                                                            | @sloth.slo objective 95.0                                                                                      |
| description | Description is the description of the SLO.                                                                                                                          | @sloth.slo description 95% of logins to the chat-gpt app should be successful annotations. (can be multilined) |
| labels      | Labels are the Prometheus labels that will have all the recording and alerting rules for this specific SLO. These labels are merged with the previous level labels. | @sloth.slo labels foo bar @sloth labels test slo                                                               |

