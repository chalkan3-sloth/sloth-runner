# 📈 Prometheus Metrics Reference

## Overview

Sloth Runner exposes comprehensive metrics in Prometheus text format. All metrics are prefixed with `sloth_` and follow Prometheus naming conventions.

## Metrics Endpoint

### HTTP Endpoints

| Endpoint | Description |
|----------|-------------|
| `/metrics` | Prometheus metrics in text format |
| `/health` | Health check endpoint (returns `OK`) |
| `/info` | Service information JSON |

### Access Methods

=== "CLI"
    ```bash
    # Get endpoint URL
    ./sloth-runner agent metrics prom <agent_name>

    # View raw metrics
    ./sloth-runner agent metrics prom <agent_name> --snapshot
    ```

=== "cURL"
    ```bash
    curl http://agent-host:9090/metrics
    ```

=== "Browser"
    ```
    http://agent-host:9090/metrics
    ```

## Metric Types

Sloth Runner uses three types of Prometheus metrics:

- **Counter**: Monotonically increasing value (e.g., total tasks executed)
- **Gauge**: Value that can go up or down (e.g., running tasks, memory usage)
- **Histogram**: Distribution of values with quantiles (e.g., task duration)

## Task Metrics

### sloth_tasks_total

**Type**: Counter
**Description**: Total number of tasks executed
**Labels**:
- `status`: Task result status (`success`, `failed`, `skipped`)
- `group`: Task group name (from `.sloth` file)

**Example**:
```prometheus
# HELP sloth_tasks_total Total number of tasks executed by status
# TYPE sloth_tasks_total counter
sloth_tasks_total{group="web-deployment",status="success"} 145
sloth_tasks_total{group="web-deployment",status="failed"} 3
sloth_tasks_total{group="database-setup",status="success"} 67
```

**Use Cases**:
- Track total tasks executed per group
- Calculate success rate: `success / (success + failed)`
- Identify failing task groups

**PromQL Examples**:
```promql
# Total successful tasks
sum(sloth_tasks_total{status="success"})

# Task failure rate
sum(rate(sloth_tasks_total{status="failed"}[5m]))

# Success rate by group
sum(sloth_tasks_total{status="success"}) by (group) /
sum(sloth_tasks_total) by (group)
```

---

### sloth_tasks_running

**Type**: Gauge
**Description**: Number of currently executing tasks
**Labels**: None

**Example**:
```prometheus
# HELP sloth_tasks_running Number of tasks currently executing
# TYPE sloth_tasks_running gauge
sloth_tasks_running 3
```

**Use Cases**:
- Monitor concurrent task execution
- Detect task queue buildup
- Track agent capacity utilization

**PromQL Examples**:
```promql
# Current running tasks
sloth_tasks_running

# Average running tasks over time
avg_over_time(sloth_tasks_running[1h])

# Max concurrent tasks in last hour
max_over_time(sloth_tasks_running[1h])
```

---

### sloth_task_duration_seconds

**Type**: Histogram
**Description**: Task execution duration in seconds
**Labels**:
- `group`: Task group name
- `task`: Task name

**Buckets**: `.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10` (seconds)

**Example**:
```prometheus
# HELP sloth_task_duration_seconds Task execution duration in seconds
# TYPE sloth_task_duration_seconds histogram
sloth_task_duration_seconds_bucket{group="web",task="deploy",le="0.005"} 0
sloth_task_duration_seconds_bucket{group="web",task="deploy",le="0.01"} 0
sloth_task_duration_seconds_bucket{group="web",task="deploy",le="1"} 45
sloth_task_duration_seconds_bucket{group="web",task="deploy",le="+Inf"} 145
sloth_task_duration_seconds_sum{group="web",task="deploy"} 234.56
sloth_task_duration_seconds_count{group="web",task="deploy"} 145
```

**Quantiles Available**:
- `quantile="0.5"` (P50, median)
- `quantile="0.9"` (P90)
- `quantile="0.99"` (P99)

**Use Cases**:
- Identify slow tasks
- Track performance degradation
- Set SLO/SLA thresholds

**PromQL Examples**:
```promql
# P99 latency by task
histogram_quantile(0.99,
  rate(sloth_task_duration_seconds_bucket[5m]))

# Average task duration
rate(sloth_task_duration_seconds_sum[5m]) /
rate(sloth_task_duration_seconds_count[5m])

# Tasks slower than 5 seconds
count(sloth_task_duration_seconds_bucket{le="5"} == 0)
```

---

## gRPC Metrics

### sloth_grpc_requests_total

**Type**: Counter
**Description**: Total number of gRPC requests
**Labels**:
- `method`: gRPC method name (e.g., `ExecuteTask`, `ExecuteCommand`)
- `status`: Request status (`ok`, `error`)

**Example**:
```prometheus
# HELP sloth_grpc_requests_total Total number of gRPC requests
# TYPE sloth_grpc_requests_total counter
sloth_grpc_requests_total{method="ExecuteTask",status="ok"} 234
sloth_grpc_requests_total{method="ExecuteTask",status="error"} 5
sloth_grpc_requests_total{method="ExecuteCommand",status="ok"} 89
```

**Use Cases**:
- Monitor gRPC call volume
- Track error rates
- Identify failing methods

**PromQL Examples**:
```promql
# gRPC request rate
sum(rate(sloth_grpc_requests_total[5m])) by (method)

# gRPC error rate
sum(rate(sloth_grpc_requests_total{status="error"}[5m])) by (method)

# gRPC success percentage
sum(sloth_grpc_requests_total{status="ok"}) /
sum(sloth_grpc_requests_total) * 100
```

---

### sloth_grpc_request_duration_seconds

**Type**: Histogram
**Description**: gRPC request duration in seconds
**Labels**:
- `method`: gRPC method name

**Buckets**: `.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10` (seconds)

**Example**:
```prometheus
# HELP sloth_grpc_request_duration_seconds gRPC request duration in seconds
# TYPE sloth_grpc_request_duration_seconds histogram
sloth_grpc_request_duration_seconds_bucket{method="ExecuteTask",le="0.1"} 200
sloth_grpc_request_duration_seconds_bucket{method="ExecuteTask",le="+Inf"} 234
sloth_grpc_request_duration_seconds_sum{method="ExecuteTask"} 45.67
sloth_grpc_request_duration_seconds_count{method="ExecuteTask"} 234
```

**Use Cases**:
- Monitor gRPC latency
- Detect network issues
- Track API performance

**PromQL Examples**:
```promql
# P99 gRPC latency
histogram_quantile(0.99,
  rate(sloth_grpc_request_duration_seconds_bucket[5m]))

# Slow gRPC calls (>1s)
sum(rate(sloth_grpc_request_duration_seconds_bucket{le="1"}[5m])) == 0
```

---

## System Metrics

### sloth_agent_uptime_seconds

**Type**: Gauge
**Description**: Agent uptime in seconds
**Labels**: None

**Example**:
```prometheus
# HELP sloth_agent_uptime_seconds Agent uptime in seconds
# TYPE sloth_agent_uptime_seconds gauge
sloth_agent_uptime_seconds 3600
```

**Use Cases**:
- Monitor agent availability
- Track restart frequency
- Calculate uptime percentage

**PromQL Examples**:
```promql
# Uptime in hours
sloth_agent_uptime_seconds / 3600

# Agents that restarted recently (< 5 minutes)
sloth_agent_uptime_seconds < 300
```

---

### sloth_agent_info

**Type**: Gauge (always 1)
**Description**: Agent build information
**Labels**:
- `version`: Agent version
- `os`: Operating system (linux, darwin, windows)
- `arch`: Architecture (amd64, arm64)

**Example**:
```prometheus
# HELP sloth_agent_info Agent build information
# TYPE sloth_agent_info gauge
sloth_agent_info{version="v1.2.3",os="linux",arch="arm64"} 1
```

**Use Cases**:
- Track agent versions across fleet
- Identify agents needing updates
- Monitor platform distribution

**PromQL Examples**:
```promql
# Count agents by version
count(sloth_agent_info) by (version)

# Count agents by OS
count(sloth_agent_info) by (os)

# Find outdated agents
sloth_agent_info{version!="v1.2.3"}
```

---

### sloth_goroutines

**Type**: Gauge
**Description**: Number of goroutines
**Labels**: None

**Example**:
```prometheus
# HELP sloth_goroutines Number of goroutines
# TYPE sloth_goroutines gauge
sloth_goroutines 342
```

**Use Cases**:
- Monitor goroutine leaks
- Track concurrency patterns
- Identify resource issues

**PromQL Examples**:
```promql
# Goroutine growth rate
rate(sloth_goroutines[5m])

# High goroutine count (> 1000)
sloth_goroutines > 1000

# Average goroutines
avg_over_time(sloth_goroutines[1h])
```

---

### sloth_memory_allocated_bytes

**Type**: Gauge
**Description**: Memory allocated by Go runtime in bytes
**Labels**: None

**Example**:
```prometheus
# HELP sloth_memory_allocated_bytes Memory allocated by Go runtime
# TYPE sloth_memory_allocated_bytes gauge
sloth_memory_allocated_bytes 81788928
```

**Use Cases**:
- Monitor memory usage
- Detect memory leaks
- Capacity planning

**PromQL Examples**:
```promql
# Memory in MB
sloth_memory_allocated_bytes / 1024 / 1024

# Memory growth rate
rate(sloth_memory_allocated_bytes[5m])

# High memory usage (> 500MB)
sloth_memory_allocated_bytes > 524288000
```

---

## Error Metrics

### sloth_errors_total

**Type**: Counter
**Description**: Total number of errors
**Labels**:
- `type`: Error type (e.g., `task_execution`, `grpc_timeout`, `module_error`)

**Example**:
```prometheus
# HELP sloth_errors_total Total number of errors by type
# TYPE sloth_errors_total counter
sloth_errors_total{type="task_execution"} 12
sloth_errors_total{type="grpc_timeout"} 3
sloth_errors_total{type="module_error"} 5
```

**Use Cases**:
- Track error frequency
- Identify error patterns
- Alert on error spikes

**PromQL Examples**:
```promql
# Error rate by type
sum(rate(sloth_errors_total[5m])) by (type)

# Total errors
sum(sloth_errors_total)

# Error increase in last hour
increase(sloth_errors_total[1h])
```

---

## Metric Update Frequency

| Metric Category | Update Trigger |
|----------------|----------------|
| Task Metrics | On task completion |
| gRPC Metrics | On each gRPC call |
| System Metrics | Every 15 seconds (automatic) |
| Error Metrics | When errors occur |

## Best Practices

### Querying

1. **Use rate() for counters**:
   ```promql
   rate(sloth_tasks_total[5m])
   ```

2. **Use histogram_quantile() for latencies**:
   ```promql
   histogram_quantile(0.99, rate(sloth_task_duration_seconds_bucket[5m]))
   ```

3. **Aggregate with by()**:
   ```promql
   sum(sloth_tasks_total) by (group, status)
   ```

### Alerting

Example Prometheus alert rules:

```yaml
groups:
  - name: sloth_runner_alerts
    interval: 30s
    rules:
      - alert: HighTaskFailureRate
        expr: |
          sum(rate(sloth_tasks_total{status="failed"}[5m]))
          /
          sum(rate(sloth_tasks_total[5m]))
          > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High task failure rate on {{ $labels.instance }}"

      - alert: SlowTaskExecution
        expr: |
          histogram_quantile(0.99,
            rate(sloth_task_duration_seconds_bucket[5m])
          ) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Slow task execution detected"

      - alert: AgentDown
        expr: up{job="sloth-runner"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Agent {{ $labels.instance }} is down"

      - alert: HighMemoryUsage
        expr: sloth_memory_allocated_bytes > 536870912
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage (> 512MB)"
```

## Recording Rules

Precompute expensive queries:

```yaml
groups:
  - name: sloth_runner_recording_rules
    interval: 30s
    rules:
      - record: job:sloth_task_success_rate:5m
        expr: |
          sum(rate(sloth_tasks_total{status="success"}[5m])) by (job)
          /
          sum(rate(sloth_tasks_total[5m])) by (job)

      - record: job:sloth_task_p99_latency:5m
        expr: |
          histogram_quantile(0.99,
            sum(rate(sloth_task_duration_seconds_bucket[5m])) by (job, le)
          )
```

## Next Steps

- [Grafana Dashboard Guide](grafana-dashboard.md) - Visualize these metrics
- [Deployment Guide](deployment.md) - Set up Prometheus scraping
- [Telemetry Overview](index.md) - Back to overview
