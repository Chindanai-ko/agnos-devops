# Observability Patterns for Agnos DevOps

## Structured JSON Logging
Use a library like `slog` (Go 1.21+) or `logrus`.

```go
import "log/slog"

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)
    
    slog.Info("Service started", "port", 8080, "env", os.Getenv("APP_ENV"))
}
```

## Metrics and Alerts
- **Latency:** Monitor `http_request_duration_seconds`.
- **Error Rate:** Monitor `http_requests_total{code=~"5.."}`.
- **Worker Status:** Track timestamp of last successful run.

### Alerting Rules (Prometheus)
```yaml
groups:
- name: agnos-alerts
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{code=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
    for: 2m
    labels:
      severity: critical
  - alert: WorkerStalled
    expr: time() - worker_last_success_timestamp_seconds > 3600
    for: 5m
    labels:
      severity: warning
```
