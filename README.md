# sloth-simple-comments
Embed Sloth SLO and SLI specifications within your application's code.

Example:
```go
    // @sloth service chatgpt
    
    var (
        // @sloth name chat-gpt-availability
        // @sloth objective 95.0
        // @sloth.sli error_query sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[{{.window}}])) OR on() vector(0)
        // @sloth.sli total_query sum(rate(tenant_login_operations_total{client="chat-gpt"}[{{.window}}]))
        // @sloth description 95% of logins to the chat-gpt app should be successful.
        metricGaugeCertInventoryProcessingMessages = prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: "chatgpt",
                Subsystem: "auth0",
                Name:      "tenant_login_operations_total",
            })
    )
```

Result:
```yaml
version: prometheus/v1
service: "chatgpt"
slos:
    - name: chat-gpt-availability
      description: 95% of logins to the chat-gpt app should be successful.
      objective: 95
      sli:
        raw:
            error_ratio_query: ""
        events:
            error_query: sum(rate(tenant_failed_login_operations_total{client="chat-gpt"}[{{.window}}])) OR on() vector(0)
            total_query: sum(rate(tenant_login_operations_total{client="chat-gpt"}[{{.window}}]))
      alerting:
        name: ""
```
