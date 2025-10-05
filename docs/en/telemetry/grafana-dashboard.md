# 📊 Grafana-Style Terminal Dashboard

## Overview

The `agent metrics grafana` command provides a comprehensive, **Grafana-inspired dashboard directly in your terminal**. View real-time metrics, performance indicators, and system health without leaving the command line.

!!! tip "No External Dependencies"
    The terminal dashboard is completely self-contained. No Grafana installation required!

## Features

### 🎨 Rich Visualization

- **Tables**: Formatted data with headers and borders
- **Progress Bars**: Visual representation of resource usage
- **Color Coding**: Green/yellow/red indicators for status
- **Sections**: Organized layout with clear separators
- **Summary Box**: At-a-glance statistics

### 🔄 Real-Time Updates

- **Watch Mode**: Auto-refresh at configurable intervals
- **Live Metrics**: See changes as tasks execute
- **Screen Clearing**: Clean updates without scrolling

### 📈 Comprehensive Metrics

- Agent information and build details
- System resources (memory, goroutines)
- Task execution statistics
- Performance metrics (P50, P99 latencies)
- gRPC request statistics
- Error tracking

## Quick Start

### Basic Usage

View dashboard once:

```bash
./sloth-runner agent metrics grafana <agent_name>
```

### Watch Mode

Auto-refresh every 5 seconds (default):

```bash
./sloth-runner agent metrics grafana <agent_name> --watch
```

Custom refresh interval (e.g., every 2 seconds):

```bash
./sloth-runner agent metrics grafana <agent_name> --watch --interval 2
```

Stop watching: Press `Ctrl+C`

## Dashboard Sections

### 1. 🔧 Agent Information

Displays agent metadata and configuration.

**Fields**:

| Field | Description | Example |
|-------|-------------|---------|
| Version | Agent build version | `v1.2.3`, `dev` |
| OS | Operating system | `linux`, `darwin` |
| Architecture | CPU architecture | `arm64`, `amd64` |
| Uptime | Time since agent start | `2h 34m`, `5d 12h 45m` |
| Last Updated | Metrics snapshot timestamp | `2025-10-05 15:42:30` |

**Example Output**:
```
🔧 Agent Information
┌──────────────┬─────────────────────────────┐
│ Version      │ v1.2.3                      │
│ OS           │ linux                       │
│ Architecture │ arm64                       │
│ Uptime       │ 2h 34m                      │
│ Last Updated │ 2025-10-05 15:42:30         │
└──────────────┴─────────────────────────────┘
```

---

### 2. 💻 System Resources

Visual progress bars showing resource utilization.

**Metrics**:

#### Goroutines
- Current count vs. threshold (default: 1000)
- Color coding:
    - 🟢 Green: < 60% (healthy)
    - 🟡 Yellow: 60-80% (moderate)
    - 🔴 Red: > 80% (high)

#### Memory (MB)
- Allocated memory in megabytes
- Threshold: 512MB default
- Same color coding as goroutines

**Example Output**:
```
💻 System Resources

Goroutines: [████████████░░░░░░░░░░░░░░░░░░░░] 342/1000 (34.2%)
Memory (MB): [██████░░░░░░░░░░░░░░░░░░░░░░░░░░] 78/512 (15.2%)
```

**Interpretation**:
- Green bars: System is healthy
- Yellow bars: Monitor closely, may need attention
- Red bars: Resource pressure, investigate

---

### 3. 📋 Task Metrics

Summary of task execution results.

**Status Table**:

| Status | Icon | Color | Description |
|--------|------|-------|-------------|
| Success | ✓ | Green | Tasks completed successfully |
| Failed | ✗ | Red | Tasks that failed |
| Skipped | ⊘ | Yellow | Tasks skipped (conditions not met) |

**Running Tasks Bar**:
- Current concurrent tasks
- Threshold: 10 (configurable)

**Example Output**:
```
📋 Task Metrics

┌───────────┬────────┐
│  Status   │ Count  │
├───────────┼────────┤
│ ✓ Success │ 145    │
│ ✗ Failed  │ 3      │
│ ⊘ Skipped │ 12     │
└───────────┴────────┘

Running Tasks: [██░░░░░░░░░░░░░░░░░░░░░] 2/10 (20.0%)
```

**Use Cases**:
- Quick health check: Low failure rate is good
- Capacity monitoring: High running tasks may indicate bottleneck
- Audit trail: Total tasks executed

---

### 4. ⏱️ Task Performance

Detailed latency metrics for executed tasks.

**Columns**:

| Column | Description |
|--------|-------------|
| Task | Task name from `.sloth` file |
| P50 (ms) | Median execution time (50th percentile) |
| P99 (ms) | 99th percentile latency |
| Status | Performance indicator |

**Performance Indicators**:

| Indicator | Criteria | Meaning |
|-----------|----------|---------|
| 🟢 Fast | P99 < 1000ms | Excellent performance |
| 🟡 Normal | P99 < 5000ms | Acceptable performance |
| 🔴 Slow | P99 >= 5000ms | Needs optimization |

**Example Output**:
```
⏱️  Task Performance

┌──────────────────┬──────────┬──────────┬──────────┐
│      Task        │ P50 (ms) │ P99 (ms) │  Status  │
├──────────────────┼──────────┼──────────┼──────────┤
│ install_packages │ 234.56   │ 567.89   │ 🟡 Normal│
│ check_service    │ 12.34    │ 45.67    │ 🟢 Fast  │
│ deploy_app       │ 1234.56  │ 5678.90  │ 🔴 Slow  │
└──────────────────┴──────────┴──────────┴──────────┘
```

**Action Items**:
- 🟢 Fast tasks: No action needed
- 🟡 Normal tasks: Monitor trends
- 🔴 Slow tasks: Investigate and optimize

---

### 5. 🌐 gRPC Metrics

Statistics for master-agent communication.

**Columns**:

| Column | Description |
|--------|-------------|
| Method | gRPC method name |
| Requests | Total requests for this method |
| Avg Latency (ms) | P50 latency in milliseconds |

**Common Methods**:

| Method | Description |
|--------|-------------|
| `ExecuteTask` | Task execution requests |
| `ExecuteCommand` | Direct command execution |
| `GetAgentInfo` | Agent info queries |
| `RegisterAgent` | Agent registration |

**Example Output**:
```
🌐 gRPC Metrics

┌─────────────────┬──────────┬──────────────────┐
│     Method      │ Requests │ Avg Latency (ms) │
├─────────────────┼──────────┼──────────────────┤
│ ExecuteTask     │ 156      │ 234.56           │
│ ExecuteCommand  │ 45       │ 12.34            │
└─────────────────┴──────────┴──────────────────┘
```

**Interpretation**:
- High request count: Agent is actively used
- High latency: Network or master performance issues
- Low latency (<50ms): Excellent connectivity

---

### 6. ⚠️ Errors

Error tracking by type (only shown if errors exist).

**Columns**:

| Column | Description |
|--------|-------------|
| Error Type | Category of error |
| Count | Number of occurrences (red) |

**Common Error Types**:

| Type | Description |
|------|-------------|
| `task_execution` | Errors during task execution |
| `grpc_timeout` | gRPC request timeouts |
| `module_error` | Module-specific errors |
| `network_error` | Network connectivity issues |

**Example Output**:
```
⚠️  Errors

┌────────────────┬───────┐
│  Error Type    │ Count │
├────────────────┼───────┤
│ task_execution │ 12    │
│ grpc_timeout   │ 3     │
│ module_error   │ 5     │
└────────────────┴───────┘
```

**Action Items**:
- Investigate errors with highest counts
- Check logs for error details
- Review failing tasks in Task Metrics section

---

### 7. 📈 Summary

Consolidated overview in a highlighted box.

**Metrics**:

| Metric | Description | Color |
|--------|-------------|-------|
| Total Tasks | All tasks executed | Cyan |
| Running | Currently executing | Yellow |
| Memory | Current allocation (MB) | Green |
| Goroutines | Active goroutines | Magenta |

**Example Output**:
```
╔═════════════════════════════════════════════════════╗
║                    📈 Summary                       ║
╠═════════════════════════════════════════════════════╣
║ Total Tasks: 160 | Running: 2 | Memory: 78 MB |     ║
║ Goroutines: 342                                     ║
╚═════════════════════════════════════════════════════╝
```

---

## Use Cases

### Development Workflow

Monitor tasks during development:

=== "Scenario"
    You're developing a deployment script and want to see metrics in real-time.

=== "Command"
    ```bash
    # Terminal 1: Watch dashboard
    ./sloth-runner agent metrics grafana dev-agent --watch --interval 1

    # Terminal 2: Execute tasks repeatedly
    for i in {1..10}; do
      ./sloth-runner run -f deploy.sloth
      sleep 2
    done
    ```

=== "Benefit"
    - See task counts increment
    - Monitor latency changes
    - Catch errors immediately

### Performance Tuning

Identify bottlenecks:

=== "Scenario"
    Production tasks are slower than expected.

=== "Command"
    ```bash
    ./sloth-runner agent metrics grafana prod-agent
    ```

=== "Action"
    1. Check **Task Performance** section
    2. Identify 🔴 Slow tasks
    3. Review P99 latencies
    4. Optimize those specific tasks

### Production Monitoring

Quick health checks:

=== "Scenario"
    Check agent health during on-call.

=== "Command"
    ```bash
    ./sloth-runner agent metrics grafana prod-agent
    ```

=== "Checks"
    - ✅ Low error count in Errors section
    - ✅ Green resource bars (System Resources)
    - ✅ Low failure rate in Task Metrics
    - ✅ Low gRPC latency

### Capacity Planning

Determine if you need more agents:

=== "Scenario"
    Deciding if current agent fleet is sufficient.

=== "Command"
    ```bash
    # Check multiple agents
    for agent in agent1 agent2 agent3; do
      echo "=== $agent ==="
      ./sloth-runner agent metrics grafana $agent | grep -A2 "System Resources"
    done
    ```

=== "Decision Criteria"
    - 🔴 Red resource bars: Add more agents
    - 🟡 Yellow consistently: Monitor closely
    - 🟢 Green: Current capacity is good

---

## Advanced Features

### Watch Mode

Continuous monitoring with auto-refresh:

```bash
# Refresh every 5 seconds (default)
./sloth-runner agent metrics grafana my-agent --watch

# Fast refresh (1 second) for development
./sloth-runner agent metrics grafana my-agent --watch --interval 1

# Slow refresh (30 seconds) for overview
./sloth-runner agent metrics grafana my-agent --watch --interval 30
```

**Features**:
- Clears screen between updates for clean display
- Press `Ctrl+C` to stop
- Ideal for monitoring during task execution

### Comparison

Compare metrics across multiple agents:

```bash
#!/bin/bash
# compare-agents.sh

agents=("agent1" "agent2" "agent3")

for agent in "${agents[@]}"; do
  echo "========================================="
  echo "Agent: $agent"
  echo "========================================="
  ./sloth-runner agent metrics grafana $agent
  echo ""
  read -p "Press Enter for next agent..."
done
```

### Scripting

Extract specific metrics for automation:

```bash
# Get current running tasks
./sloth-runner agent metrics grafana my-agent | grep "Running Tasks"

# Check for errors
./sloth-runner agent metrics grafana my-agent | grep -A10 "⚠️  Errors"

# Extract memory usage
./sloth-runner agent metrics grafana my-agent | grep "Memory (MB)"
```

---

## Color Reference

### Status Colors

| Color | Hex | Usage |
|-------|-----|-------|
| 🟢 Green | `#4CAF50` | Success, healthy, fast |
| 🟡 Yellow | `#FFC107` | Warning, moderate, skipped |
| 🔴 Red | `#F44336` | Error, high, slow |
| 🔵 Cyan | `#00BCD4` | Information, totals |
| 🟣 Magenta | `#9C27B0` | Secondary metrics |

### Visual Indicators

| Symbol | Meaning |
|--------|---------|
| ✓ | Success |
| ✗ | Failure |
| ⊘ | Skipped |
| 🟢 | Fast/Healthy |
| 🟡 | Normal/Warning |
| 🔴 | Slow/Critical |

---

## Troubleshooting

### Dashboard Shows "No Data"

**Symptoms**: All metrics are zero or empty

**Causes**:
1. Agent just started (no tasks executed yet)
2. Telemetry disabled
3. Metrics endpoint unreachable

**Solutions**:

```bash
# Check if agent has telemetry enabled
./sloth-runner agent list

# Verify metrics endpoint
./sloth-runner agent metrics prom my-agent --snapshot

# Execute a test task to generate metrics
./sloth-runner agent run my-agent "echo test"

# Try dashboard again
./sloth-runner agent metrics grafana my-agent
```

### Connection Refused

**Symptoms**: "Failed to fetch metrics: connection refused"

**Causes**:
1. Agent is down
2. Metrics port is blocked
3. Wrong agent name

**Solutions**:

```bash
# Verify agent is running
./sloth-runner agent list

# Check metrics endpoint
curl http://agent-ip:9090/health

# Check firewall
telnet agent-ip 9090
```

### Incomplete Dashboard

**Symptoms**: Some sections missing

**Causes**:
1. No data for that metric category (e.g., no errors = no Errors section)
2. Old agent version without all metrics

**Solutions**:
- This is normal! Sections only appear when data exists.
- For Errors section: Only shown when errors > 0
- For Task Performance: Only shown when tasks have been executed

### Watch Mode Not Updating

**Symptoms**: Dashboard frozen in watch mode

**Causes**:
1. Terminal doesn't support ANSI escape codes
2. Very long refresh interval

**Solutions**:

```bash
# Use shorter interval
./sloth-runner agent metrics grafana my-agent --watch --interval 2

# Try different terminal
# (e.g., iTerm2, modern Terminal.app, Windows Terminal)

# Fallback: Run without watch mode
./sloth-runner agent metrics grafana my-agent
```

---

## Best Practices

### Refresh Intervals

| Use Case | Recommended Interval |
|----------|---------------------|
| Active development | 1-2 seconds |
| Task execution monitoring | 5 seconds (default) |
| Background monitoring | 10-30 seconds |
| Overview checks | Single run (no watch) |

### When to Use

✅ **Use Dashboard For**:
- Quick health checks
- Real-time task monitoring
- Performance troubleshooting
- Development feedback

❌ **Don't Use Dashboard For**:
- Historical analysis (use Grafana web UI)
- Alerting (use Prometheus alerts)
- Long-term trends (use time-series visualization)
- Multi-agent comparison (manually run for each)

### Complementary Tools

| Tool | When to Use |
|------|-------------|
| Terminal Dashboard | Quick checks, development |
| Prometheus | Historical queries, alerting |
| Grafana Web UI | Long-term trends, dashboards |
| `agent metrics prom` | Get endpoint URL, raw metrics |

---

## Examples

### Example 1: Healthy Agent

```
📊 Sloth Runner Metrics Dashboard - Agent: production-1

🔧 Agent Information
┌──────────────┬─────────────────────┐
│ Version      │ v1.2.3              │
│ OS           │ linux               │
│ Architecture │ amd64               │
│ Uptime       │ 7d 14h 23m          │
│ Last Updated │ 2025-10-05 10:30:15 │
└──────────────┴─────────────────────┘

💻 System Resources
Goroutines: [█████░░░░░░░░░] 125/1000 (12.5%)
Memory (MB): [███░░░░░░░░░░] 45/512 (8.8%)

📋 Task Metrics
┌───────────┬───────┐
│  Status   │ Count │
├───────────┼───────┤
│ ✓ Success │ 1,234 │
│ ✗ Failed  │ 5     │
└───────────┴───────┘

Running Tasks: [░░░░░░░░░░░░] 0/10 (0.0%)

⏱️  Task Performance
┌─────────────┬──────────┬──────────┬─────────┐
│    Task     │ P50 (ms) │ P99 (ms) │ Status  │
├─────────────┼──────────┼──────────┼─────────┤
│ health_check│ 5.23     │ 12.45    │ 🟢 Fast │
│ deploy      │ 456.78   │ 892.34   │ 🟡 Normal│
└─────────────┴──────────┴──────────┴─────────┘

╔════════════════════════════════════════════╗
║              📈 Summary                    ║
╠════════════════════════════════════════════╣
║ Total Tasks: 1,239 | Running: 0 |          ║
║ Memory: 45 MB | Goroutines: 125           ║
╚════════════════════════════════════════════╝
```

### Example 2: Agent Under Load

```
📊 Sloth Runner Metrics Dashboard - Agent: worker-3

💻 System Resources
Goroutines: [███████████████████████░░░] 857/1000 (85.7%)
Memory (MB): [████████████████████░░░░░░] 412/512 (80.5%)

📋 Task Metrics
Running Tasks: [████████░░] 8/10 (80.0%)

⏱️  Task Performance
┌────────────┬──────────┬──────────┬─────────┐
│    Task    │ P50 (ms) │ P99 (ms) │ Status  │
├────────────┼──────────┼──────────┼─────────┤
│ big_deploy │ 3456.78  │ 8932.12  │ 🔴 Slow │
└────────────┴──────────┴──────────┴─────────┘

⚠️  Errors
┌────────────────┬───────┐
│  Error Type    │ Count │
├────────────────┼───────┤
│ task_timeout   │ 23    │
└────────────────┴───────┘
```

**Interpretation**: This agent is under heavy load. Consider:
- Reducing concurrent tasks
- Optimizing slow tasks
- Adding more agents to distribute load
- Investigating task timeouts

---

## Next Steps

- [Prometheus Metrics Reference](prometheus-metrics.md) - Detailed metric documentation
- [Deployment Guide](deployment.md) - Set up production monitoring
- [Telemetry Overview](index.md) - Back to overview

## Further Reading

- [pterm Library Documentation](https://github.com/pterm/pterm) - Terminal visualization library used
- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/) - Metric naming and usage
