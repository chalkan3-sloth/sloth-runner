# 🖥️ Complete System Monitoring with Grafana Dashboard

## Overview

The enhanced Grafana dashboard now provides **comprehensive system monitoring** for the entire machine, not just Sloth Runner agent metrics. This gives you a complete view of system health and performance.

## Features

### 📊 System Metrics Collected

#### CPU Monitoring
- **Overall CPU usage** percentage
- **Per-core utilization** with visual bars
- **Load averages** (1, 5, 15 minutes)
- **CPU model and specifications**
- **Number of cores and threads**
- **CPU speed in GHz**

#### Memory Monitoring
- **RAM usage** (used, free, available)
- **Swap usage** and statistics
- **Memory pressure indicators**
- **Visual progress bars** with color coding

#### Disk Monitoring
- **All mount points** with usage percentages
- **Filesystem types** and device names
- **Total, used, and free space**
- **Inode utilization**
- **Color-coded usage warnings**

#### Network Monitoring
- **Network interfaces** statistics
- **Bytes sent/received**
- **Packets transmitted/received**
- **Network errors and drops**
- **Active connection count**

#### Process Monitoring
- **Top 10 processes** by CPU usage
- **Process memory consumption**
- **PID, username, and status**
- **Zombie process detection**
- **Total process count**

#### System Information
- **Hostname and OS details**
- **Kernel version and architecture**
- **System uptime**
- **Boot time**
- **Platform information**

## Usage

### Basic Command

```bash
# Monitor a running agent with complete system metrics
sloth-runner agent metrics grafana <agent-name>

# Example
sloth-runner agent metrics grafana lady-arch
```

### Continuous Monitoring

The dashboard automatically refreshes to show real-time metrics:

```bash
# Monitor with auto-refresh (press Ctrl+C to exit)
sloth-runner agent metrics grafana lady-arch --refresh 5
```

## Dashboard Layout

```
📊 Sloth Runner Complete System Monitor - Agent: lady-arch

🖥️ System Overview
┌──────────────────┬────────────────────────────────┐
│ Hostname         │ lady-arch                      │
│ OS               │ Ubuntu 22.04                   │
│ Kernel           │ 5.15.0-89 (x86_64)            │
│ Uptime           │ 15d 4h 32m                     │
│ Boot Time        │ 2025-09-20 08:15:23           │
│ Processes        │ 245 (Zombies: 0)              │
│ Network Conn     │ 142                           │
└──────────────────┴────────────────────────────────┘

🔥 CPU Metrics
┌──────────────────┬────────────────────────────────┐
│ Model            │ Intel Core i7-10700K           │
│ Cores            │ 8                              │
│ Threads          │ 16                             │
│ Speed            │ 3.80 GHz                       │
│ Load Average     │ 1.23, 1.45, 1.32              │
└──────────────────┴────────────────────────────────┘
CPU Usage: [████████████░░░░░░░░░░░░░░░░] 32.5%

┌─ CPU Cores Usage ──────────────────────────────────┐
│ Core  0: ████████░░░░░░░░░░░  42.1% | Core  1: ███████░░░░░░░░░░░░  38.5% │
│ Core  2: ██████░░░░░░░░░░░░░  31.2% | Core  3: █████░░░░░░░░░░░░░░  28.7% │
│ Core  4: ████░░░░░░░░░░░░░░░  22.3% | Core  5: ███░░░░░░░░░░░░░░░░  18.9% │
│ Core  6: ██████████░░░░░░░░░  52.4% | Core  7: ████████░░░░░░░░░░░  44.6% │
└────────────────────────────────────────────────────┘

💾 Memory Metrics
RAM Usage:  [████████████████████░░░░░░░░░░░░░░░░░] 55.2%
Swap Usage: [████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 12.3%

┌─────┬─────────┬─────────┬─────────┬───────────┐
│ Type│ Total   │ Used    │ Free    │ Available │
├─────┼─────────┼─────────┼─────────┼───────────┤
│ RAM │ 32.0 GB │ 17.6 GB │ 14.4 GB │ 13.2 GB   │
│ Swap│ 8.0 GB  │ 984 MB  │ 7.0 GB  │ -         │
└─────┴─────────┴─────────┴─────────┴───────────┘

💿 Disk Metrics
┌────────────┬─────────────┬──────┬─────────┬─────────┬─────────┬──────┐
│ Mount      │ Device      │ Type │ Size    │ Used    │ Free    │ Use% │
├────────────┼─────────────┼──────┼─────────┼─────────┼─────────┼──────┤
│ /          │ /dev/sda2   │ ext4 │ 500 GB  │ 245 GB  │ 255 GB  │ 49.0%│
│ /home      │ /dev/sda3   │ ext4 │ 1.5 TB  │ 890 GB  │ 610 GB  │ 59.3%│
│ /boot      │ /dev/sda1   │ vfat │ 512 MB  │ 128 MB  │ 384 MB  │ 25.0%│
└────────────┴─────────────┴──────┴─────────┴─────────┴─────────┴──────┘

🌐 Network Metrics
┌───────────┬──────────┬──────────┬─────────────┬────────┬───────┐
│ Interface │ Sent     │ Received │ Packets     │ Errors │ Drops │
├───────────┼──────────┼──────────┼─────────────┼────────┼───────┤
│ eth0      │ 15.2 GB  │ 48.7 GB  │ 25M/89M     │ 0/0    │ 0/0   │
│ wlan0     │ 892 MB   │ 2.3 GB   │ 1.2M/3.5M   │ 0/0    │ 0/0   │
└───────────┴──────────┴──────────┴─────────────┴────────┴───────┘

📋 Top Processes
┌──────┬──────────────────┬──────────┬───────┬─────────┬────────┐
│ PID  │ Name             │ User     │ CPU%  │ Memory  │ Status │
├──────┼──────────────────┼──────────┼───────┼─────────┼────────┤
│ 2451 │ chrome           │ chalkan3 │ 12.3% │ 892 MB  │ S      │
│ 1823 │ vscode           │ chalkan3 │ 8.7%  │ 1.2 GB  │ S      │
│ 9821 │ sloth-runner     │ chalkan3 │ 5.2%  │ 45 MB   │ S      │
│ 3421 │ postgres         │ postgres │ 3.8%  │ 256 MB  │ S      │
│ 5632 │ nginx            │ www-data │ 2.1%  │ 32 MB   │ S      │
└──────┴──────────────────┴──────────┴───────┴─────────┴────────┘

🦥 Sloth Runner Agent
┌──────────────┬────────────────────┐
│ Version      │ v4.17.2           │
│ Uptime       │ 2h 45m            │
│ Goroutines   │ 42                │
│ Memory       │ 45 MB             │
└──────────────┴────────────────────┘

┌─ 📈 System Summary ─────────────────────────────────────┐
│ 🖥️  CPU: 32.5% | RAM: 55.2% | Disk: 3 mounts         │
│ 📊 Processes: 245 | Network: 142 conn | Uptime: 15d   │
│ 🦥 Tasks: 1523 | Running: 2 | Agent Memory: 45 MB    │
└─────────────────────────────────────────────────────────┘
```

## Color Coding

The dashboard uses intuitive color coding for quick status assessment:

- 🟢 **Green**: Normal (0-60%)
- 🟡 **Yellow**: Warning (60-80%)
- 🔴 **Red**: Critical (80-100%)

## Performance Considerations

### Resource Usage

The monitoring system is designed to be lightweight:
- CPU overhead: < 1%
- Memory usage: ~10-20 MB
- Network traffic: Minimal (local metrics only)

### Sampling Intervals

- CPU metrics: 1 second sample
- Memory metrics: Real-time
- Disk metrics: On-demand
- Network metrics: Cumulative counters
- Process metrics: Real-time snapshot

## Remote Monitoring

Monitor agents on remote machines:

```bash
# Monitor remote agent
sloth-runner agent metrics grafana lady-arch

# The agent's telemetry endpoint must be accessible
# Default port: 9090
```

## Troubleshooting

### Common Issues

1. **Metrics not available**
   - Ensure agent is running with telemetry enabled
   - Check firewall rules for port 9090

2. **High CPU usage in monitoring**
   - Increase refresh interval
   - Reduce number of processes tracked

3. **Permission errors**
   - Some metrics require elevated permissions
   - Run with appropriate privileges if needed

## Configuration

### Environment Variables

```bash
# Set custom metrics port
export SLOTH_METRICS_PORT=9091

# Set refresh interval (seconds)
export SLOTH_METRICS_REFRESH=10

# Disable color output
export NO_COLOR=1
```

### Agent Configuration

Ensure your agent is started with telemetry:

```bash
# Start agent with telemetry
sloth-runner agent start my-agent --metrics-port 9090
```

## API Integration

Access raw metrics programmatically:

```bash
# Get Prometheus metrics
curl http://agent-host:9090/metrics

# Get system info
curl http://agent-host:9090/info

# Health check
curl http://agent-host:9090/health
```

## Best Practices

1. **Regular Monitoring**: Set up continuous monitoring for production systems
2. **Alert Thresholds**: Configure alerts for critical metrics
3. **Historical Data**: Store metrics for trend analysis
4. **Capacity Planning**: Use metrics for resource planning
5. **Performance Baseline**: Establish normal operating parameters

## Examples

### Monitor Local System

```bash
# Start local agent with monitoring
sloth-runner agent start local --metrics-port 9090

# Open monitoring dashboard
sloth-runner agent metrics grafana local
```

### Monitor Container Environment

```bash
# Monitor agent inside container
sloth-runner agent metrics grafana container-agent

# Works with Docker, Incus, LXD
```

### Scripted Monitoring

```bash
#!/bin/bash
# Monitor and alert on high CPU

while true; do
    metrics=$(sloth-runner agent metrics prom my-agent)
    cpu_usage=$(echo "$metrics" | grep cpu_usage | awk '{print $2}')

    if (( $(echo "$cpu_usage > 80" | bc -l) )); then
        echo "ALERT: High CPU usage: $cpu_usage%"
        # Send notification
    fi

    sleep 60
done
```

## Summary

The enhanced Grafana dashboard transforms Sloth Runner into a comprehensive system monitoring solution, providing:

- **Complete visibility** into system resources
- **Real-time metrics** with visual indicators
- **Process tracking** for resource optimization
- **Network monitoring** for connectivity insights
- **Disk usage** tracking across all mount points
- **Integration** with Sloth Runner task metrics

This makes it an ideal tool for DevOps, system administrators, and developers who need comprehensive monitoring without complex setup.