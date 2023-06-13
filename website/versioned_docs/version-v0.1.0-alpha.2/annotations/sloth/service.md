---
sidebar_position: 1
---

# Service Annotations

Define details about the service that owns the SLOs being defined.

The following is an example of what the in-code annotations defining a simple service:

```go
// @sloth service chat-gpt
// @sloth version prometheus/v1
// @sloth labels tier ml-application  
// @sloth labels environment staging
```

The above annotations will generate the following sloth specification:

```yaml
version: prometheus/v1
service: chatgpt
labels:
  tier: ml-application
  environment: staging
```

## Table of Annotations

| Annotation | Description                                                     | Example                                      |
|------------|-----------------------------------------------------------------|----------------------------------------------|
| service    | **Required**. The name of the service the definitions refer to. | @sloth service chat-gpt                      |
| version    | The version of the Sloth specification.                         | @sloth version prometheus/v1                 |
| labels     | The labels associated to the Sloth service.                     | @sloth labels foo bar @sloth labels test slo |
