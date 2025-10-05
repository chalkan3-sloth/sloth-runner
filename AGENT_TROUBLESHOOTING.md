# 🔧 Agent Troubleshooting Guide: lady-arch

## Problem Summary

**Agent:** `lady-arch` (Arch Linux container inside Incus)
**Host:** `192.168.1.16` (physical machine)
**Status:** Active (sends heartbeats) but unreachable for commands
**Error:** "Agent Unreachable - The master server cannot reach the agent"

## Root Cause Analysis

The agent `lady-arch` is running **inside an Incus container** but is registered with the **host IP** (`192.168.1.16:50052`) instead of being accessible from outside.

### Current Setup:
```
┌─────────────────────────────────────────┐
│  Master (Mac - 192.168.1.29)            │
│  Tries to connect to: 192.168.1.16:50052│
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│  Host (192.168.1.16)                     │
│  ┌───────────────────────────────────┐  │
│  │ Incus Container: lady-arch        │  │
│  │ Agent listening on: 0.0.0.0:50052 │  │
│  │ Container IP: 10.x.x.x (internal) │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

**Problem:** The master tries to connect to `192.168.1.16:50052`, but the container's port 50052 is **not forwarded** to the host.

## Solutions

### ✅ Solution 1: Incus Proxy Device (Recommended)

Add a proxy device to forward traffic from host port to container:

```bash
# On the host (192.168.1.16)
incus config device add lady-arch sloth-proxy proxy \
  listen=tcp:0.0.0.0:50052 \
  connect=tcp:127.0.0.1:50052

# Restart the container to apply
incus restart lady-arch
```

**Verification:**
```bash
# From master or any machine
nc -zv 192.168.1.16 50052
# Should show: Connection succeeded!
```

### ✅ Solution 2: Use Host Network Mode

Run the container with host networking:

```bash
# Stop current container
incus stop lady-arch

# Create new profile with host networking
incus profile create host-network
incus profile set host-network raw.lxc "lxc.net.0.type = none"

# Apply to container
incus profile add lady-arch host-network

# Start container
incus start lady-arch

# Restart agent inside container
incus exec lady-arch -- systemctl restart sloth-runner-agent
```

### ✅ Solution 3: Change Agent Report Address

Make the agent report the correct accessible address:

```bash
# Inside the container (lady-arch)
# Edit the agent service to use correct report address

# Stop the agent
systemctl stop sloth-runner-agent

# Start with correct report address
sloth-runner agent start \
  --name lady-arch \
  --master 192.168.1.29:50053 \
  --port 50052 \
  --bind-address 0.0.0.0 \
  --report-address 192.168.1.16:50052 \
  --daemon
```

**Note:** The `--report-address` flag tells the agent what address to register with the master.

### ✅ Solution 4: NAT/Firewall Rule on Host

Add iptables rule on the host to forward:

```bash
# On host (192.168.1.16)
sudo iptables -t nat -A PREROUTING -p tcp --dport 50052 \
  -j DNAT --to-destination <CONTAINER_IP>:50052

sudo iptables -A FORWARD -p tcp -d <CONTAINER_IP> --dport 50052 -j ACCEPT

# Make persistent
sudo iptables-save > /etc/iptables/rules.v4
```

## Quick Fix (Recommended)

**Use Incus Proxy Device:**

```bash
# SSH to the host
ssh user@192.168.1.16

# Add proxy device
incus config device add lady-arch sloth-agent-proxy proxy \
  listen=tcp:0.0.0.0:50052 \
  connect=tcp:127.0.0.1:50052

# Verify
incus config show lady-arch | grep -A 3 devices

# Test from master
nc -zv 192.168.1.16 50052
```

## Testing

After applying the fix:

```bash
# Test from master (Mac)
./sloth-runner agent run lady-arch 'ls -la'

# Should show file listing from lady-arch container
```

## Verification Checklist

- [ ] Agent shows as "Active" in `sloth-runner agent list`
- [ ] Port 50052 is accessible from master: `nc -zv 192.168.1.16 50052`
- [ ] Agent responds to commands: `sloth-runner agent run lady-arch 'whoami'`
- [ ] Heartbeats are being received (check master.log)

## Current Status

**Agent Registration:**
```
Name:          lady-arch
Address:       192.168.1.16:50052
Status:        Active
Last Heartbeat: Recent (< 1 min ago)
```

**Network Test Results:**
```
✅ Ping to host (192.168.1.16): SUCCESS
✅ TCP connection to port 50052: SUCCESS (from Mac)
❌ gRPC connection to agent: FAILED
```

**Diagnosis:** The port is open but the gRPC service inside the container is not reachable because Incus needs a proxy device to forward the traffic.

## Next Steps

1. **Immediate Fix**: Add Incus proxy device (Solution 1)
2. **Verify**: Test with `sloth-runner agent run lady-arch 'ls -la'`
3. **Long-term**: Consider using `--report-address` flag for all container agents

## Related Issues

- Agent inside Docker containers: Use `--network host` or port mapping
- Agent inside VMs: Ensure port forwarding in VM hypervisor
- Agent behind NAT: Use correct `--report-address` flag

## Additional Notes

The agent `lady-arch` successfully sends **heartbeats** to the master, which means:
- ✅ Container has network connectivity
- ✅ Agent can reach master at 192.168.1.29:50053
- ✅ Outbound traffic works

But the master **cannot reach the agent back** because:
- ❌ Inbound traffic to container port 50052 is not forwarded
- ❌ Container IP is not directly accessible from outside
