# SLOTH-RUNNER-AGENT(1) - Remote Agent Management

## NAME

**sloth-runner agent** - Manage distributed remote agents for task execution

## SYNOPSIS

```
sloth-runner agent <command> [options]
```

## DESCRIPTION

The **agent** command provides complete lifecycle management for remote agents in Sloth Runner. Agents are worker nodes that execute delegated tasks, enabling distributed automation across multiple hosts.

Agents provide:

- **Remote task execution** - Run tasks on specific infrastructure nodes
- **System information collection** - Query OS, hardware, network details
- **Module availability checking** - Verify installed tools (docker, systemd, incus, etc.)
- **Metrics and telemetry** - Export Prometheus metrics for monitoring
- **Automatic registration** - Self-register with master on startup
- **Heartbeat monitoring** - Continuous health checks via gRPC

The agent system uses gRPC for efficient bidirectional communication between the master and agents, with automatic reconnection and failure handling.

## ARCHITECTURE

```
┌─────────────┐         gRPC/50051        ┌──────────────┐
│   Master    │ <────────────────────────>│  Agent 1     │
│  (Laptop)   │                           │ (Server A)   │
└─────────────┘                           └──────────────┘
       │
       │         gRPC/50051
       │
       v
┌──────────────┐
│  Agent 2     │
│ (Server B)   │
└──────────────┘
```

## AVAILABLE COMMANDS

- **install** - Bootstrap a new agent on a remote host via SSH
- **start** - Start an agent locally or as daemon
- **stop** - Stop a running agent
- **list** - List all registered agents
- **get** - Get detailed agent information
- **delete** - Remove an agent from the registry
- **update** - Update an agent to the latest version
- **exec** - Execute arbitrary commands on an agent
- **modules** - Check available modules/tools on an agent
- **metrics** - View agent metrics and telemetry

## AGENT INSTALL

Bootstrap a new agent on a remote host via SSH. This command:

1. Connects to the remote host via SSH
2. Downloads the latest sloth-runner binary
3. Installs to `/usr/local/bin/sloth-runner`
4. Creates and enables a systemd service
5. Starts the agent
6. Agent automatically registers with the master

### Synopsis

```
sloth-runner agent install <agent-name> [options]
```

### Options

```
--ssh-host <host>          SSH hostname or IP address (required)
--ssh-user <user>          SSH username (default: root)
--ssh-port <port>          SSH port (default: 22)
--ssh-key <path>           SSH private key path (default: ~/.ssh/id_rsa)
--master <addr>            Master server address (default: localhost:50051)
--port <port>              Agent port (default: 50051)
--bind-address <addr>      Agent bind address (default: 0.0.0.0)
--report-address <addr>    Address agent reports to master (optional)
```

### Examples

Install agent on a remote server:

```bash
sloth-runner agent install prod-web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user root \
  --master 192.168.1.10:50051
```

Install with custom SSH key and port:

```bash
sloth-runner agent install staging-db \
  --ssh-host db.staging.example.com \
  --ssh-user deploy \
  --ssh-port 2222 \
  --ssh-key ~/.ssh/deploy_key \
  --master master.example.com:50051
```

Install agent that binds to specific interface:

```bash
sloth-runner agent install edge-node-01 \
  --ssh-host 10.0.5.20 \
  --ssh-user root \
  --bind-address 10.0.5.20 \
  --port 50060 \
  --master 10.0.1.10:50051
```

Install agent with custom report address (for NAT/firewall):

```bash
sloth-runner agent install cloud-vm-01 \
  --ssh-host 192.168.100.50 \
  --ssh-user ubuntu \
  --bind-address 0.0.0.0 \
  --port 50051 \
  --report-address 203.0.113.45:50051 \
  --master 192.168.1.10:50051
```

Sample output:

```
Installing agent 'prod-web-01' on 192.168.1.100...

✓ Connected to remote host via SSH
✓ Downloaded sloth-runner v5.0.0
✓ Installed binary to /usr/local/bin/sloth-runner
✓ Created systemd service: sloth-runner-agent
✓ Enabled service for autostart
✓ Started agent service

Agent 'prod-web-01' successfully installed and running!

Connection details:
  Address:  192.168.1.100:50051
  Version:  v5.0.0
  Status:   Active (registered with master)
```

## AGENT START

Start an agent process locally. Useful for:
- Running an agent on the same machine as the master
- Testing agent functionality
- Manual agent management (without systemd)

### Synopsis

```
sloth-runner agent start [options]
```

### Options

```
--name <name>              Agent name (default: default-agent)
--master <addr>            Master server address (default: localhost:50051)
--port <port>              Port to listen on (default: 50052)
--bind-address <addr>      Address to bind to (default: 0.0.0.0)
--report-address <addr>    Address to report to master
--daemon                   Run as background daemon
--telemetry                Enable telemetry and metrics server
--metrics-port <port>      Port for metrics server (default: 9090)
```

### Examples

Start a local agent:

```bash
sloth-runner agent start --name local-agent
```

Start agent with custom port:

```bash
sloth-runner agent start \
  --name build-agent \
  --port 50060 \
  --master 192.168.1.10:50051
```

Start agent as daemon with telemetry:

```bash
sloth-runner agent start \
  --name monitoring-agent \
  --daemon \
  --telemetry \
  --metrics-port 9090
```

Start agent on specific interface:

```bash
sloth-runner agent start \
  --name edge-agent \
  --bind-address 10.0.5.20 \
  --port 50051 \
  --master 10.0.1.10:50051
```

Sample output:

```
Starting agent 'local-agent'...

✓ Agent listening on 0.0.0.0:50052
✓ Connected to master at localhost:50051
✓ Registered with master
✓ Heartbeat started (interval: 30s)

Agent 'local-agent' is running
Press Ctrl+C to stop
```

## AGENT STOP

Stop a running agent by sending a shutdown request via the master.

### Synopsis

```
sloth-runner agent stop <agent-name> [options]
```

### Options

```
--master <addr>    Master server address (default: localhost:50051)
```

### Examples

Stop an agent:

```bash
sloth-runner agent stop prod-web-01
```

Stop with custom master address:

```bash
sloth-runner agent stop staging-db \
  --master master.example.com:50051
```

## AGENT LIST

List all registered agents. Reads from local database by default, or queries the master server if specified.

### Synopsis

```
sloth-runner agent list [options]
```

### Options

```
--master <addr>    Query master server instead of local database
--local            Force reading from local database only
--debug            Enable debug logging
```

### Examples

List agents from local database:

```bash
sloth-runner agent list
```

List agents from master server:

```bash
sloth-runner agent list --master 192.168.1.10:50051
```

Sample output:

```
Registered Agents:
┌─────────────────┬──────────────────────┬──────────┬───────────┬─────────────────────┐
│ Name            │ Address              │ Version  │ Status    │ Last Heartbeat      │
├─────────────────┼──────────────────────┼──────────┼───────────┼─────────────────────┤
│ prod-web-01     │ 192.168.1.100:50051  │ v5.0.0   │ Online    │ 2025-10-06 17:42:15 │
│ prod-web-02     │ 192.168.1.101:50051  │ v5.0.0   │ Online    │ 2025-10-06 17:42:18 │
│ staging-db      │ 192.168.2.50:50051   │ v5.0.0   │ Online    │ 2025-10-06 17:42:10 │
│ build-server    │ 192.168.1.200:50060  │ v4.9.2   │ Offline   │ 2025-10-06 15:20:33 │
└─────────────────┴──────────────────────┴──────────┴───────────┴─────────────────────┘

Total: 4 agents (3 online, 1 offline)
```

## AGENT GET

Retrieve detailed system information from a specific agent.

### Synopsis

```
sloth-runner agent get <agent-name> [options]
```

### Options

```
--master <addr>        Master server address (default: localhost:50051)
-o, --output <format>  Output format: text or json (default: text)
```

### Examples

Get agent information:

```bash
sloth-runner agent get prod-web-01
```

Get as JSON for scripting:

```bash
sloth-runner agent get prod-web-01 --output json
```

Get from specific master:

```bash
sloth-runner agent get prod-web-01 \
  --master 192.168.1.10:50051 \
  --output text
```

Sample output (text format):

```
Agent: prod-web-01
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Connection:
  Address:         192.168.1.100:50051
  Version:         v5.0.0
  Status:          Online
  Last Heartbeat:  2025-10-06 17:42:45

System Information:
  Hostname:        web-server-01
  OS:              Ubuntu 22.04.3 LTS
  Kernel:          5.15.0-89-generic
  Architecture:    x86_64
  Uptime:          15 days, 3:22:15

Hardware:
  CPU:             Intel Xeon E5-2680 v4 @ 2.40GHz
  Cores:           8 (16 threads)
  Memory:          32 GB
  Disk:            500 GB SSD

Network:
  Interfaces:      eth0 (192.168.1.100), lo (127.0.0.1)
  Hostname:        web-server-01.local

Available Modules:
  ✓ docker          Docker version 24.0.6
  ✓ systemd         systemd 249
  ✓ pkg             apt
  ✓ ssh             OpenSSH_8.9
  ✓ git             git version 2.34.1
  ✗ incus           Not installed
```

Sample output (JSON format):

```json
{
  "name": "prod-web-01",
  "address": "192.168.1.100:50051",
  "version": "v5.0.0",
  "status": "online",
  "last_heartbeat": "2025-10-06T17:42:45Z",
  "system_info": {
    "hostname": "web-server-01",
    "os": "Ubuntu 22.04.3 LTS",
    "kernel": "5.15.0-89-generic",
    "architecture": "x86_64",
    "uptime_seconds": 1321335,
    "cpu": "Intel Xeon E5-2680 v4 @ 2.40GHz",
    "cpu_cores": 8,
    "memory_total_gb": 32,
    "disk_total_gb": 500
  },
  "modules": {
    "docker": {"available": true, "version": "24.0.6"},
    "systemd": {"available": true, "version": "249"},
    "pkg": {"available": true, "type": "apt"},
    "ssh": {"available": true, "version": "OpenSSH_8.9"},
    "git": {"available": true, "version": "2.34.1"},
    "incus": {"available": false}
  }
}
```

## AGENT DELETE

Remove an agent from the registry. Note: This does not stop the agent if it's running.

### Synopsis

```
sloth-runner agent delete <agent-name> [options]
```

### Options

```
--master <addr>    Master server address (default: localhost:50051)
--force            Skip confirmation prompt
```

### Examples

Delete an agent (with confirmation):

```bash
sloth-runner agent delete old-build-server
```

Delete without confirmation:

```bash
sloth-runner agent delete old-build-server --force
```

Sample output:

```
Warning: This will remove agent 'old-build-server' from the registry.
The agent will continue running if active, but will not be reachable.

Are you sure? (y/N): y

✓ Agent 'old-build-server' removed from registry
```

## AGENT UPDATE

Update an agent to a specific version or the latest version from GitHub releases.

### Synopsis

```
sloth-runner agent update <agent-name> [options]
```

### Options

```
--master <addr>        Master server address (default: localhost:50051)
--version <version>    Version to update to (default: latest)
--restart              Restart agent service after update (default: true)
```

### Examples

Update to latest version:

```bash
sloth-runner agent update prod-web-01
```

Update to specific version:

```bash
sloth-runner agent update prod-web-01 --version v5.1.0
```

Update without restarting:

```bash
sloth-runner agent update staging-db --restart=false
```

Sample output:

```
Updating agent 'prod-web-01'...

✓ Current version: v5.0.0
✓ Latest version:  v5.1.0
✓ Downloaded update binary
✓ Stopped agent service
✓ Replaced binary at /usr/local/bin/sloth-runner
✓ Restarted agent service
✓ Agent reconnected successfully

Agent 'prod-web-01' updated to v5.1.0
```

## AGENT EXEC

Execute arbitrary shell commands on a remote agent. Useful for:
- Ad-hoc system administration
- Quick diagnostics
- Testing agent connectivity

### Synopsis

```
sloth-runner agent exec <agent-name> <command> [options]
```

### Options

```
--master <addr>        Master server address
-o, --output <format>  Output format: text or json (default: text)
```

### Examples

Execute a simple command:

```bash
sloth-runner agent exec prod-web-01 "hostname"
```

Check disk usage:

```bash
sloth-runner agent exec prod-web-01 "df -h /"
```

Run complex commands:

```bash
sloth-runner agent exec prod-web-01 \
  "docker ps --format '{{.Names}}: {{.Status}}'"
```

Get output as JSON:

```bash
sloth-runner agent exec prod-web-01 "uptime" --output json
```

Sample output (text):

```
Executing on agent 'prod-web-01': hostname

Output:
web-server-01

Exit Code: 0
Duration: 124ms
```

Sample output (JSON):

```json
{
  "agent": "prod-web-01",
  "command": "hostname",
  "exit_code": 0,
  "stdout": "web-server-01\n",
  "stderr": "",
  "duration_ms": 124
}
```

Common use cases:

```bash
# Check system load
sloth-runner agent exec prod-web-01 "uptime"

# View running processes
sloth-runner agent exec prod-web-01 "ps aux | head -20"

# Check service status
sloth-runner agent exec prod-web-01 "systemctl status nginx"

# View logs
sloth-runner agent exec prod-web-01 "journalctl -u nginx -n 50"

# Test connectivity
sloth-runner agent exec prod-web-01 "ping -c 3 google.com"
```

## AGENT MODULES

Check which external tools and modules are available on an agent. This helps ensure tasks will work correctly when delegated.

### Synopsis

```
sloth-runner agent modules <agent-name> [options]
```

### Options

```
--master <addr>    Master server address (default: localhost:50051)
--check <module>   Check specific module availability
```

### Available Modules

- **docker** - Docker container runtime
- **systemd** - System and service manager
- **pkg** - Package manager (apt, yum, apk, etc.)
- **ssh** - SSH client
- **git** - Git version control
- **incus** - Incus container/VM manager
- **http** - HTTP client utilities (curl/wget)

### Examples

Check all modules:

```bash
sloth-runner agent modules prod-web-01
```

Check specific module:

```bash
sloth-runner agent modules prod-web-01 --check docker
```

Sample output:

```
Modules available on agent 'prod-web-01':

┌──────────┬───────────┬─────────────────────────┐
│ Module   │ Available │ Version/Details         │
├──────────┼───────────┼─────────────────────────┤
│ docker   │ ✓         │ Docker version 24.0.6   │
│ systemd  │ ✓         │ systemd 249             │
│ pkg      │ ✓         │ apt (Ubuntu)            │
│ ssh      │ ✓         │ OpenSSH_8.9             │
│ git      │ ✓         │ git version 2.34.1      │
│ incus    │ ✗         │ Not installed           │
│ http     │ ✓         │ curl 7.81.0             │
└──────────┴───────────┴─────────────────────────┘

7 modules checked: 6 available, 1 unavailable
```

Use in automation:

```bash
# Verify Docker is available before deployment
if sloth-runner agent modules prod-web-01 --check docker | grep -q "Available.*✓"; then
    echo "Docker available, proceeding with deployment"
    sloth-runner run deploy_containers --delegate-to prod-web-01 --yes
else
    echo "Docker not available on prod-web-01"
    exit 1
fi
```

## AGENT METRICS

View and manage agent metrics and telemetry. Agents can export Prometheus metrics for monitoring.

### Synopsis

```
sloth-runner agent metrics <command> [options]
```

### Available Commands

- **prom** - Get Prometheus metrics endpoint
- **grafana** - Display detailed metrics dashboard

### Examples

View Prometheus metrics:

```bash
sloth-runner agent metrics prom prod-web-01
```

Display metrics dashboard:

```bash
sloth-runner agent metrics grafana prod-web-01
```

Sample Prometheus metrics:

```
# HELP sloth_agent_tasks_total Total number of tasks executed
# TYPE sloth_agent_tasks_total counter
sloth_agent_tasks_total{agent="prod-web-01",status="success"} 1247
sloth_agent_tasks_total{agent="prod-web-01",status="failed"} 23

# HELP sloth_agent_task_duration_seconds Task execution duration
# TYPE sloth_agent_task_duration_seconds histogram
sloth_agent_task_duration_seconds_bucket{le="1"} 450
sloth_agent_task_duration_seconds_bucket{le="5"} 980
sloth_agent_task_duration_seconds_bucket{le="30"} 1240
sloth_agent_task_duration_seconds_sum 15234.5
sloth_agent_task_duration_seconds_count 1270

# HELP sloth_agent_cpu_usage_percent CPU usage percentage
# TYPE sloth_agent_cpu_usage_percent gauge
sloth_agent_cpu_usage_percent{agent="prod-web-01"} 23.5

# HELP sloth_agent_memory_usage_bytes Memory usage in bytes
# TYPE sloth_agent_memory_usage_bytes gauge
sloth_agent_memory_usage_bytes{agent="prod-web-01"} 8589934592
```

## COMPLETE DEPLOYMENT EXAMPLE

### 1. Install Agents on Multiple Servers

```bash
# Production web servers
sloth-runner agent install prod-web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user root \
  --master 192.168.1.10:50051

sloth-runner agent install prod-web-02 \
  --ssh-host 192.168.1.101 \
  --ssh-user root \
  --master 192.168.1.10:50051

# Database server
sloth-runner agent install prod-db \
  --ssh-host 192.168.1.150 \
  --ssh-user root \
  --master 192.168.1.10:50051

# Build server
sloth-runner agent install build-server \
  --ssh-host 192.168.1.200 \
  --ssh-user deploy \
  --port 50060 \
  --master 192.168.1.10:50051
```

### 2. Verify Agent Status

```bash
# List all agents
sloth-runner agent list --master 192.168.1.10:50051

# Check specific agent
sloth-runner agent get prod-web-01 --master 192.168.1.10:50051

# Verify Docker is available
sloth-runner agent modules prod-web-01 --check docker
```

### 3. Execute Tasks on Agents

```bash
# Deploy application to web servers
sloth-runner run deploy_app \
  --file workflows/deploy.sloth \
  --delegate-to prod-web-01 \
  --delegate-to prod-web-02 \
  --yes

# Run database migration
sloth-runner run migrate_db \
  --file workflows/database.sloth \
  --delegate-to prod-db \
  --yes
```

### 4. Monitor and Maintain

```bash
# Check agent health
watch -n 5 'sloth-runner agent list'

# View metrics
sloth-runner agent metrics grafana prod-web-01

# Update agents
sloth-runner agent update prod-web-01
sloth-runner agent update prod-web-02
sloth-runner agent update prod-db
```

## TROUBLESHOOTING

### Agent Not Connecting

```bash
# Check agent status
sloth-runner agent list

# Test SSH connectivity
ssh -p 22 root@192.168.1.100

# Check agent logs on remote host
ssh root@192.168.1.100 "journalctl -u sloth-runner-agent -n 50"

# Restart agent
ssh root@192.168.1.100 "systemctl restart sloth-runner-agent"
```

### Agent Version Mismatch

```bash
# Check current version
sloth-runner agent get prod-web-01 | grep Version

# Update to latest
sloth-runner agent update prod-web-01
```

### Firewall Issues

```bash
# Verify port is open
sloth-runner agent exec prod-web-01 "nc -zv 192.168.1.10 50051"

# Check firewall rules
sloth-runner agent exec prod-web-01 "iptables -L -n | grep 50051"
```

## FILES

```
.sloth-cache/agents.db              Local agent registry database
/usr/local/bin/sloth-runner         Agent binary (on remote hosts)
/etc/systemd/system/sloth-runner-agent.service    Agent systemd service
```

## ENVIRONMENT VARIABLES

```
SLOTH_RUNNER_MASTER_ADDR    Default master server address
```

## SEE ALSO

- **sloth-runner(1)** - Main sloth-runner command
- **sloth-runner-run(1)** - Execute workflows with delegation
- **sloth-runner-master(1)** - Start master server

## AUTHOR

Written by the Sloth Runner development team.

## COPYRIGHT

Copyright © 2025. Released under MIT License.
