# 🔄 Agent Auto-Reconnection

## Overview

The Sloth Runner agent now includes **automatic reconnection** capabilities, ensuring high availability and resilience in distributed environments. When an agent loses connection to the master server, it will automatically attempt to reconnect without manual intervention.

## Features

### 🛡️ Robust Connection Management

- **Automatic reconnection** when connection to master is lost
- **Exponential backoff** strategy to avoid overwhelming the master
- **Connection health monitoring** via heartbeats
- **Smart failure detection** with configurable thresholds
- **Zero downtime** - agent continues running while reconnecting

### 📊 Connection Parameters

| Parameter | Default Value | Description |
|-----------|--------------|-------------|
| Initial Reconnect Delay | 5 seconds | Time to wait before first reconnection attempt |
| Max Reconnect Delay | 60 seconds | Maximum delay between reconnection attempts |
| Heartbeat Interval | 5 seconds | Frequency of heartbeat messages to master |
| Max Consecutive Failures | 3 | Number of failed heartbeats before reconnecting |
| Connection Timeout | 10 seconds | Timeout for establishing connection to master |

### 🔍 How It Works

1. **Initial Connection**
   - Agent starts and connects to the master server
   - Registers with the master using its name and address
   - Begins sending periodic heartbeats

2. **Health Monitoring**
   - Every 5 seconds, the agent sends a heartbeat to the master
   - Tracks consecutive heartbeat failures
   - After 3 consecutive failures, triggers reconnection

3. **Reconnection Process**
   - Closes the failed connection
   - Waits with exponential backoff (5s → 10s → 20s → ... → 60s max)
   - Attempts to establish a new connection
   - Re-registers with the master
   - Resumes normal operation

4. **Recovery**
   - On successful heartbeat after failures, resets failure counter
   - Logs recovery status for monitoring

## Usage Examples

### Starting an Agent with Auto-Reconnection

```bash
# Start agent in daemon mode with auto-reconnection enabled
sloth-runner agent start \
  --name my-agent \
  --port 50051 \
  --master master.example.com:50050 \
  --daemon
```

### Starting Multiple Agents

```bash
# Agent 1
sloth-runner agent start --name ladyguica --port 50051 \
  --master 192.168.1.15:50050 --daemon

# Agent 2
sloth-runner agent start --name keiteguica --port 50052 \
  --master 192.168.1.15:50050 --daemon

# Agent 3 (on remote host)
ssh user@192.168.1.16 "sloth-runner agent start \
  --name mariguica --port 50051 \
  --master 192.168.1.15:50050 --daemon"
```

## Monitoring and Logs

### Log Messages

The agent logs provide detailed information about connection status:

```
INFO  Agent registered with master at 192.168.1.15:50050, reporting address 192.168.1.16:50051
✓ Agent registered with master at 192.168.1.15:50050 (reporting address: 192.168.1.16:50051)

WARN  Heartbeat failed (1/3): rpc error: code = Unavailable desc = connection error
WARN  Heartbeat failed (2/3): rpc error: code = Unavailable desc = connection error
WARN  Heartbeat failed (3/3): rpc error: code = Unavailable desc = connection error

ERROR Lost connection to master after 3 failed heartbeats. Reconnecting...
⚠ Connection to master lost. Attempting to reconnect...

INFO  🔄 Reconnecting to master in 5s...
✓ Agent registered with master at 192.168.1.15:50050 (reporting address: 192.168.1.16:50051)
INFO  Heartbeat recovered, connection stable
```

### Viewing Agent Logs

```bash
# View real-time logs
tail -f agent.log

# Search for reconnection events
grep "Reconnecting" agent.log

# Count reconnection attempts
grep -c "Lost connection" agent.log
```

## Network Scenarios

### Scenario 1: Master Server Restart

**What happens:**
1. Master server goes down for maintenance
2. Agent detects 3 consecutive heartbeat failures
3. Agent enters reconnection mode
4. Master comes back online
5. Agent automatically reconnects and re-registers
6. Normal operation resumes

**Downtime:** Minimal (15-20 seconds)

### Scenario 2: Network Interruption

**What happens:**
1. Network connection between agent and master is interrupted
2. Agent detects heartbeat failures
3. Agent attempts reconnection with exponential backoff
4. Network is restored
5. Agent reconnects immediately
6. Tasks can be executed again

**Recovery time:** Depends on network restoration + backoff delay

### Scenario 3: Agent Restart

**What happens:**
1. Agent is restarted (manual or automatic)
2. Agent immediately connects to master
3. Registers with its configured name and address
4. Ready to execute tasks

**Downtime:** Time to restart the agent process (~1-2 seconds)

## Best Practices

### 1. Use Daemon Mode in Production

Always run agents in daemon mode for production environments:

```bash
sloth-runner agent start --name prod-agent \
  --master master.prod.com:50050 --daemon
```

### 2. Monitor Agent Logs

Set up log monitoring and alerting:

```bash
# Example: Alert on multiple reconnection attempts
grep -c "Lost connection" agent.log | \
  awk '$1 > 10 { print "WARNING: Agent reconnected", $1, "times" }'
```

### 3. Configure Report Address for NAT/Firewalls

When agents are behind NAT or firewalls:

```bash
sloth-runner agent start --name agent-behind-nat \
  --bind-address 0.0.0.0 \
  --report-address public-ip.example.com \
  --master master.example.com:50050 \
  --daemon
```

### 4. Load Balancing

For high-availability setups, consider running multiple agents:

```bash
# Start 3 agents on different machines
for i in {1..3}; do
  ssh node-$i "sloth-runner agent start \
    --name agent-node-$i \
    --master master.example.com:50050 \
    --daemon"
done
```

## Troubleshooting

### Agent Cannot Connect Initially

**Problem:** Agent fails to connect on startup

**Solutions:**
- Verify master address is correct
- Check network connectivity: `ping master.example.com`
- Ensure master is running: `netstat -tulpn | grep 50050`
- Check firewall rules

### Frequent Reconnections

**Problem:** Agent reconnects too often

**Possible causes:**
- Network instability
- Master server overloaded
- Firewall dropping idle connections
- DNS resolution issues

**Solutions:**
- Check network quality
- Monitor master server resources
- Adjust firewall timeout settings
- Use IP address instead of hostname

### Agent Shows as Inactive

**Problem:** `sloth-runner agent list` shows agent as inactive

**Solutions:**
- Check if agent process is running
- Review agent logs for errors
- Verify network connectivity
- Restart the agent

## Technical Details

### Connection State Machine

```
[Start] → [Connecting] → [Registered] → [Active]
                ↑              ↓
                └──[Reconnecting]←─[Disconnected]
```

### Exponential Backoff Algorithm

```
delay = min(initial_delay * 2^(attempt-1), max_delay)

Example:
Attempt 1: 5s
Attempt 2: 10s
Attempt 3: 20s
Attempt 4: 40s
Attempt 5: 60s (capped)
Attempt 6: 60s (capped)
```

### Heartbeat Protocol

- Uses gRPC streaming for efficient communication
- Lightweight messages (~100 bytes)
- Timeout-based failure detection
- Graceful degradation on network issues

## Performance Impact

### Resource Usage

- **CPU:** Negligible (<0.1% on idle)
- **Memory:** ~2MB for connection management
- **Network:** ~200 bytes/5s (heartbeat traffic)

### Latency

- **Detection time:** 15 seconds (3 × 5s heartbeats)
- **Reconnection time:** 5-60 seconds (depends on backoff)
- **Total recovery time:** 20-75 seconds worst case

## Security Considerations

Currently, the agent uses **insecure gRPC connections**. For production environments:

1. Consider implementing TLS/SSL encryption
2. Use authentication tokens or certificates
3. Implement network segmentation
4. Use VPN or secure tunnels for agent-master communication

## Future Enhancements

- [ ] Configurable reconnection parameters via flags
- [ ] TLS/SSL support for encrypted connections
- [ ] Authentication and authorization
- [ ] Metrics export (Prometheus/Grafana)
- [ ] Health check endpoint
- [ ] Graceful shutdown handling
- [ ] Connection pooling for multiple masters
- [ ] Circuit breaker pattern implementation

## Related Documentation

- [Agent Setup Guide](agent-setup.md)
- [Master Configuration](master-setup.md)
- [Troubleshooting Guide](troubleshooting.md)
- [Security Best Practices](security.md)
