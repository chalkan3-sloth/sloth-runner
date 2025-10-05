# 📊 Telemetry & Observability

## Overview

Sloth Runner provides comprehensive telemetry and observability features through native **Prometheus integration** and a rich **terminal-based Grafana-style dashboard**. Monitor your agent fleet, track task execution metrics, analyze performance, and gain deep insights into your infrastructure automation.

!!! success "Enterprise-Grade Observability"
    Built-in Prometheus metrics server with auto-discovery, real-time dashboards, and zero-configuration setup.

## Key Features

### 🎯 Prometheus Integration

- **Native Metrics Exporter**: Built-in HTTP server exposing Prometheus-compatible metrics
- **Auto-Discovery**: Metrics endpoint automatically configured on agent startup
- **Standard Format**: Compatible with Prometheus, Grafana, and all observability tools
- **Zero Configuration**: Telemetry enabled by default with sensible defaults

### 📊 Terminal Dashboard

- **Rich Visualization**: Beautiful terminal-based dashboard with tables, charts, and progress bars
- **Real-time Updates**: Watch mode with configurable refresh intervals
- **Comprehensive Metrics**: System resources, task performance, gRPC stats, and error tracking
- **Color-Coded Insights**: Visual indicators for performance and health status

### 📈 Metrics Categories

#### Task Metrics
- Total tasks executed (by status: success, failed, skipped)
- Currently running tasks
- Task duration histograms (P50, P99 latencies)
- Per-task and per-group performance tracking

#### System Metrics
- Agent uptime
- Memory allocation
- Goroutines count
- Agent version and build information

#### gRPC Metrics
- Request counts per method
- Request duration histograms
- Success/error rates

#### Error Tracking
- Error counts by type
- Failed task tracking
- System error monitoring

## Quick Start

### Enable Telemetry on Agent

Telemetry is enabled by default. Start your agent:

```bash
./sloth-runner agent start --name my-agent --master localhost:50053
```

To explicitly configure telemetry:

```bash
# Enable telemetry with custom port
./sloth-runner agent start \
  --name my-agent \
  --master localhost:50053 \
  --telemetry \
  --metrics-port 9090
```

To disable telemetry:

```bash
./sloth-runner agent start \
  --name my-agent \
  --master localhost:50053 \
  --telemetry=false
```

### Access Metrics

#### Get Prometheus Endpoint

```bash
./sloth-runner agent metrics prom my-agent
```

Output:
```
✅ Metrics Endpoint:
  URL: http://192.168.1.100:9090/metrics

📝 Usage:
  # View metrics in browser:
  open http://192.168.1.100:9090/metrics

  # Fetch metrics via curl:
  curl http://192.168.1.100:9090/metrics

  # Configure Prometheus scraper:
  - job_name: 'sloth-runner-agents'
    static_configs:
      - targets: ['192.168.1.100:9090']
```

#### View Snapshot

```bash
./sloth-runner agent metrics prom my-agent --snapshot
```

### View Dashboard

#### Single View

```bash
./sloth-runner agent metrics grafana my-agent
```

#### Watch Mode (Auto-Refresh)

```bash
# Refresh every 5 seconds (default)
./sloth-runner agent metrics grafana my-agent --watch

# Custom refresh interval (2 seconds)
./sloth-runner agent metrics grafana my-agent --watch --interval 2
```

## Architecture

```mermaid
graph LR
    A[Sloth Runner Agent] --> B[Telemetry Server :9090]
    B --> C[/metrics endpoint]
    B --> D[/health endpoint]
    B --> E[/info endpoint]

    C --> F[Prometheus Scraper]
    C --> G[CLI: agent metrics prom]
    C --> H[CLI: agent metrics grafana]

    F --> I[Prometheus Server]
    I --> J[Grafana Dashboards]

    style B fill:#4CAF50
    style C fill:#2196F3
    style H fill:#FF9800
```

### Components

1. **Telemetry Server** (`internal/telemetry/server.go`)
   - HTTP server running on configurable port (default: 9090)
   - Serves Prometheus metrics in text format
   - Provides health check and service info endpoints

2. **Metrics Collector** (`internal/telemetry/metrics.go`)
   - Defines all Prometheus metrics (counters, gauges, histograms)
   - Thread-safe global singleton
   - Automatic runtime metrics collection

3. **Visualizer** (`internal/telemetry/visualizer.go`)
   - Fetches and parses Prometheus metrics
   - Rich terminal dashboard rendering
   - Historical trends support

4. **CLI Commands**
   - `agent metrics prom`: Get endpoint URL or raw metrics
   - `agent metrics grafana`: Display rich dashboard

## Use Cases

### Development

Monitor your tasks during development:

```bash
# Terminal 1: Watch dashboard
./sloth-runner agent metrics grafana dev-agent --watch --interval 1

# Terminal 2: Execute tasks
./sloth-runner run -f deploy.sloth --values dev.yaml
```

### Production Monitoring

Integrate with Prometheus and Grafana:

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'sloth-runner-production'
    static_configs:
      - targets:
          - 'agent1:9090'
          - 'agent2:9090'
          - 'agent3:9090'
        labels:
          environment: production
```

### Performance Analysis

Identify slow tasks and bottlenecks:

```bash
# View detailed performance metrics
./sloth-runner agent metrics grafana prod-agent

# Check P99 latencies in Task Performance section
# Tasks with 🔴 Slow indicator need optimization
```

### Debugging

Track errors and failures:

```bash
# View error counts
./sloth-runner agent metrics grafana my-agent

# Check Errors section for error types
# Cross-reference with Task Metrics for failed tasks
```

## Next Steps

- [Prometheus Metrics Reference](prometheus-metrics.md) - Complete metrics documentation
- [Grafana Dashboard Guide](grafana-dashboard.md) - Dashboard features and usage
- [Deployment Guide](deployment.md) - Production deployment and integration

## Supported Platforms

- ✅ Linux (amd64, arm64)
- ✅ macOS (Intel, Apple Silicon)
- ✅ Windows (via WSL2)
- ✅ Containers (Docker, Incus/LXC)
- ✅ Kubernetes (via DaemonSet)

## Performance Impact

Telemetry has **minimal performance overhead**:

- Memory: ~10-20MB additional
- CPU: <1% under normal load
- Network: Metrics served only on-demand (pull model)
- Storage: Metrics stored in-memory, no persistence

## Security Considerations

!!! warning "Network Exposure"
    The metrics endpoint is exposed on all network interfaces by default. In production:

    - Use firewall rules to restrict access
    - Consider binding to localhost only and using reverse proxy
    - Enable authentication via reverse proxy (Prometheus doesn't support auth natively)

!!! tip "Best Practice"
    Run agents in private networks and expose metrics only to monitoring infrastructure.

## Troubleshooting

### Telemetry Not Starting

Check agent logs for errors:

```bash
tail -f agent.log | grep -i telemetry
```

Verify port availability:

```bash
netstat -tuln | grep 9090
```

Try different port:

```bash
./sloth-runner agent start --name my-agent --metrics-port 9091
```

### Cannot Access Metrics

Test from agent host:

```bash
curl http://localhost:9090/metrics
```

Test from remote:

```bash
curl http://agent-ip:9090/metrics
```

Check firewall:

```bash
# Allow port 9090
sudo ufw allow 9090/tcp

# Or use firewalld
sudo firewall-cmd --permanent --add-port=9090/tcp
sudo firewall-cmd --reload
```

### Dashboard Shows No Data

Verify agent is running with telemetry:

```bash
./sloth-runner agent list
```

Check metrics endpoint directly:

```bash
./sloth-runner agent metrics prom my-agent --snapshot
```

Ensure tasks have been executed (initial metrics are zero):

```bash
./sloth-runner agent run my-agent "echo test"
```

## Further Reading

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [pterm Library](https://github.com/pterm/pterm) (used for terminal visualization)
