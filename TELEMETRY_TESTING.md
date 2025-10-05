# Testing Prometheus Telemetry Module on lady-arch

## Overview

The Prometheus telemetry module has been implemented and is ready for testing on the `lady-arch` agent.

## Architecture

### Components

1. **internal/telemetry/metrics.go** - Prometheus metrics definitions
   - Counters: tasks_total, grpc_requests_total, errors_total
   - Gauges: tasks_running, agent_uptime, memory, goroutines
   - Histograms: task_duration, grpc_duration

2. **internal/telemetry/server.go** - HTTP server for metrics
   - Port: 9090 (default, configurable)
   - Endpoint: `/metrics` (Prometheus format)
   - Endpoint: `/health` (Health check)
   - Endpoint: `/info` (Service info)

3. **internal/telemetry/collector.go** - Global singleton
   - Convenience functions for recording metrics
   - Thread-safe

### CLI Commands

```bash
# Get metrics endpoint URL
sloth-runner agent metrics prom <agent_name>

# View current metrics snapshot
sloth-runner agent metrics prom <agent_name> --snapshot
```

## Testing Steps

### 1. Deploy Binary to lady-arch

```bash
# From macOS (task-runner directory)
GOOS=linux GOARCH=arm64 go build -o sloth-runner-linux-arm64-telemetry cmd/sloth-runner/*.go

# Copy to lady-arch via incus
ssh igor@192.168.1.16 "incus file push sloth-runner-linux-arm64-telemetry lady-arch/home/igor/sloth-runner-new"
```

### 2. Stop Current Agent

```bash
# SSH into host
ssh igor@192.168.1.16

# Enter container
incus exec lady-arch -- bash

# Stop running agent
pkill sloth-runner

# Or find and kill by PID
ps aux | grep sloth-runner
kill <PID>
```

### 3. Start Agent with Telemetry

```bash
# Inside lady-arch container
cd /home/igor

# Make executable
chmod +x sloth-runner-new

# Start with telemetry enabled (default)
./sloth-runner-new agent start --name lady-arch --master 192.168.1.2:50053 --telemetry --metrics-port 9090 &

# Or explicitly disable telemetry
./sloth-runner-new agent start --name lady-arch --master 192.168.1.2:50053 --telemetry=false
```

### 4. Verify Telemetry Server

```bash
# Inside lady-arch container
curl http://localhost:9090/metrics

# Expected output: Prometheus format metrics
# HELP sloth_tasks_total Total number of tasks executed by status
# TYPE sloth_tasks_total counter
# sloth_tasks_total{group="default",status="success"} 0
...
```

### 5. Test CLI Commands (from macOS)

```bash
# Get metrics endpoint
./sloth-runner-telemetry agent metrics prom lady-arch

# Expected output:
# ✅ Metrics Endpoint:
#   URL: http://192.168.1.16:9090/metrics
#
# 📝 Usage:
#   # View metrics in browser:
#   open http://192.168.1.16:9090/metrics
#   ...

# View metrics snapshot
./sloth-runner-telemetry agent metrics prom lady-arch --snapshot

# Expected output: Full metrics dump

# 🆕 View detailed Grafana-style dashboard
./sloth-runner-telemetry agent metrics grafana lady-arch

# Expected output: Rich terminal dashboard with:
# - Agent information (version, OS, arch, uptime)
# - System resources (goroutines, memory) with progress bars
# - Task metrics summary (success, failed, skipped)
# - Task performance (P50/P99 latency)
# - gRPC metrics (request count, latency)
# - Error summary
# - Overall summary box

# Watch mode: auto-refresh every 5 seconds
./sloth-runner-telemetry agent metrics grafana lady-arch --watch

# Custom refresh interval (10 seconds)
./sloth-runner-telemetry agent metrics grafana lady-arch --watch --interval 10
```

### 6. Execute Tasks to Generate Metrics

```bash
# Run a simple task
./sloth-runner-telemetry agent run lady-arch "echo 'test task'"

# Check metrics again
curl http://192.168.1.16:9090/metrics | grep sloth_tasks
```

### 7. Prometheus Integration (Optional)

Create a Prometheus config to scrape the agent:

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'sloth-runner-agents'
    static_configs:
      - targets: ['192.168.1.16:9090']
        labels:
          agent: 'lady-arch'
```

Run Prometheus:

```bash
docker run -p 9091:9090 -v $PWD/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

Access Prometheus UI: http://localhost:9091

## Metrics Reference

### Task Metrics

- `sloth_tasks_total{status, group}` - Counter of all tasks
  - Labels: status (success, failed, skipped), group (task group name)

- `sloth_tasks_running` - Gauge of currently executing tasks

- `sloth_task_duration_seconds{group, task}` - Histogram of task durations
  - Buckets: .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10

### gRPC Metrics

- `sloth_grpc_requests_total{method, status}` - Counter of gRPC calls
  - Labels: method (ExecuteTask, etc.), status (ok, error)

- `sloth_grpc_request_duration_seconds{method}` - Histogram of gRPC latency
  - Buckets: .001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10

### System Metrics

- `sloth_agent_uptime_seconds` - Agent uptime in seconds
- `sloth_agent_info{version, os, arch}` - Agent build information (always 1)
- `sloth_goroutines` - Number of goroutines
- `sloth_memory_allocated_bytes` - Memory allocated by Go runtime

### Error Metrics

- `sloth_errors_total{type}` - Counter of errors by type

## Troubleshooting

### Telemetry server not starting

```bash
# Check if port 9090 is already in use
netstat -tuln | grep 9090

# Try different port
./sloth-runner-new agent start --name lady-arch --metrics-port 9091
```

### Cannot access metrics from macOS

```bash
# Check firewall on lady-arch
sudo ufw status

# Allow port if needed
sudo ufw allow 9090/tcp

# Test connectivity
nc -zv 192.168.1.16 9090
```

### Metrics not updating

```bash
# Check agent logs
tail -f agent.log

# Verify telemetry is enabled
curl http://localhost:9090/info
# Should return: {"service":"sloth-runner","metrics_port":9090,"metrics_endpoint":"/metrics"}
```

## Expected Behavior

1. **Agent startup**: Should log "✓ Telemetry server started at http://localhost:9090/metrics"

2. **Metrics endpoint**: Should respond with Prometheus-format metrics

3. **Runtime updates**: Metrics like goroutines, memory should update every 15s

4. **Task execution**: Each task should increment counters and record durations

5. **CLI command**: Should show formatted endpoint info and usage examples

## Success Criteria

- [ ] Agent starts with telemetry enabled
- [ ] Metrics endpoint accessible at http://192.168.1.16:9090/metrics
- [ ] Metrics update in real-time
- [ ] CLI command `agent metrics prom lady-arch` works
- [ ] CLI command `agent metrics grafana lady-arch` displays rich dashboard
- [ ] Watch mode (`--watch` flag) updates dashboard continuously
- [ ] Task execution reflects in metrics
- [ ] Prometheus can scrape metrics (optional)

## Grafana Command Features

The new `agent metrics grafana` command provides:

### 📊 Dashboard Sections

1. **Agent Information**
   - Version, OS, Architecture
   - Uptime (formatted as days/hours/minutes)
   - Last update timestamp

2. **System Resources**
   - Goroutines count with visual progress bar
   - Memory allocation (MB) with visual progress bar
   - Color-coded thresholds (green < 60%, yellow < 80%, red >= 80%)

3. **Task Metrics**
   - Summary table by status (success, failed, skipped)
   - Color-coded status indicators (✓, ✗, ⊘)
   - Running tasks count with progress bar

4. **Task Performance**
   - P50 and P99 latency per task (in milliseconds)
   - Performance status indicators:
     - 🟢 Fast (P99 < 1s)
     - 🟡 Normal (P99 < 5s)
     - 🔴 Slow (P99 >= 5s)

5. **gRPC Metrics**
   - Request count per method
   - Average latency (P50) in milliseconds

6. **Error Summary**
   - Error count by type
   - Color-coded error counts

7. **Summary Box**
   - Total tasks executed
   - Currently running tasks
   - Current memory usage
   - Current goroutine count

### 🔄 Watch Mode

Enable continuous updates with `--watch`:

```bash
# Default: refresh every 5 seconds
./sloth-runner-telemetry agent metrics grafana lady-arch --watch

# Custom interval: refresh every 10 seconds
./sloth-runner-telemetry agent metrics grafana lady-arch --watch --interval 10
```

Features:
- Clears screen between updates
- Real-time metrics tracking
- Press Ctrl+C to stop
- Useful for monitoring during task execution

## Next Steps

After successful testing:

1. Update agent deployment automation
2. Add Grafana dashboards for visualization
3. Set up Prometheus Alert Manager for critical metrics
4. Document metric interpretation guide
5. Add custom business metrics as needed
