# 🚀 Sloth Runner Agent Bootstrap

Quick and easy agent installation and configuration script for Sloth Runner.

## Features

- ✅ **One-line installation** - Install agent with a single command
- ✅ **Automatic systemd setup** - Creates and enables systemd service (when available)
- ✅ **Direct start fallback** - Works in containers and environments without systemd
- ✅ **Auto-reconnection** - Agent reconnects automatically on failures
- ✅ **Production-ready** - Includes proper configuration and resource limits
- ✅ **Cross-platform** - Works on Linux and macOS
- ✅ **Flexible configuration** - Many options for customization
- ✅ **Smart detection** - Auto-detects platform, IP, and systemd availability

## Quick Start

### Basic Installation

Install agent with just a name:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \
  --name myagent
```

This will:
1. Install sloth-runner to `/usr/local/bin`
2. Create systemd service
3. Configure agent to connect to `localhost:50053`
4. Enable and start the service

### Production Installation

Install agent with full configuration:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \
  --name production-agent-1 \
  --master 192.168.1.10:50053 \
  --port 50051 \
  --bind-address 192.168.1.20 \
  --user slothrunner
```

### User Installation (No Sudo)

Install to user directory without systemd:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \
  --name myagent \
  --no-sudo \
  --no-systemd
```

## Options

### Required

| Option | Description | Example |
|--------|-------------|---------|
| `--name` | Agent name (must be unique) | `--name myagent` |

### Optional

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `--master` | Master server address | `localhost:50053` | `--master 192.168.1.10:50053` |
| `--port` | Agent listening port | `50051` | `--port 50051` |
| `--bind-address` | IP address to bind to | Auto-detected | `--bind-address 192.168.1.20` |
| `--user` | User to run agent as | Current user | `--user slothrunner` |
| `--install-dir` | Installation directory | `/usr/local/bin` | `--install-dir /opt/bin` |
| `--version` | Specific version to install | Latest | `--version v3.23.1` |
| `--no-systemd` | Skip systemd service creation | - | `--no-systemd` |
| `--no-sudo` | Install to ~/.local/bin | - | `--no-sudo` |

## What It Does

1. **Downloads and installs** sloth-runner binary
2. **Auto-detects** platform (Linux/macOS, amd64/arm64)
3. **Checks systemd availability** and functionality
4. **Creates systemd service** (if available) with:
   - Type=forking with --daemon mode
   - Automatic restart on failure
   - Proper logging to journald
   - Resource limits
   - PID file management
5. **OR starts agent directly** (if systemd not available):
   - Runs agent with --daemon flag
   - Suitable for containers and limited environments
6. **Enables and starts** the service or agent
7. **Verifies** agent is running

## Post-Installation

### With Systemd

#### Check Agent Status

```bash
# View service status
sudo systemctl status sloth-runner-agent

# View live logs
sudo journalctl -u sloth-runner-agent -f

# Restart agent
sudo systemctl restart sloth-runner-agent

# Stop agent
sudo systemctl stop sloth-runner-agent
```

### Without Systemd (Direct Mode)

#### Manage Agent

```bash
# Check if agent is running
ps aux | grep sloth-runner | grep -v grep

# Stop agent
sudo pkill -f 'sloth-runner agent'

# View logs
cat agent.log

# Restart agent manually
/usr/local/bin/sloth-runner agent start \
  --name myagent \
  --master localhost:50053 \
  --port 50051 \
  --bind-address 192.168.1.20 \
  --daemon
```

### Verify Agent on Master

```bash
# List all registered agents
sloth-runner agent list

# Test agent
sloth-runner agent run myagent "hostname"

# Delete agent
sloth-runner agent delete myagent
```

## Systemd Service Details

The bootstrap script creates a systemd service at `/etc/systemd/system/sloth-runner-agent.service`:

```ini
[Unit]
Description=Sloth Runner Agent - myagent
Documentation=https://chalkan3.github.io/sloth-runner/

[Service]
Type=forking
User=root
WorkingDirectory=/root
PIDFile=/var/run/sloth-runner-agent-myagent.pid
Restart=on-failure
RestartSec=5s
StartLimitInterval=60s
StartLimitBurst=3
Environment="HOME=/root"
ExecStartPre=/usr/bin/mkdir -p /var/run
ExecStartPost=/bin/bash -c 'if [ -f /tmp/sloth-runner-agent-myagent.pid ]; then mv /tmp/sloth-runner-agent-myagent.pid /var/run/sloth-runner-agent-myagent.pid; fi'

ExecStart=/usr/local/bin/sloth-runner agent start \
  --name myagent \
  --master localhost:50053 \
  --port 50051 \
  --bind-address 192.168.1.20 \
  --daemon

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=sloth-runner-agent

# Performance
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

**Note**: The service uses `Type=forking` with `--daemon` flag, which makes the agent run in background mode and allows systemd to properly track the process.

## Examples

### Multiple Agents on Same Host

Install multiple agents with different ports:

```bash
# Agent 1
bash <(curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh) \
  --name agent-01 \
  --port 50051

# Agent 2
bash <(curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh) \
  --name agent-02 \
  --port 50052
```

Note: You'll need to modify the service name for additional agents.

### Remote Installation via SSH

Install agent on remote host:

```bash
ssh user@remote-host "bash <(curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh) \
  --name remote-agent \
  --master 192.168.1.10:50053 \
  --bind-address 192.168.1.20"
```

### Vagrant Installation

Install agent in Vagrant VM:

```bash
# From host machine
cd /path/to/vagrant
vagrant ssh -c "curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh | sudo bash -s -- \
  --name mariaguica \
  --master 192.168.1.29:50053 \
  --port 50051 \
  --bind-address 172.17.0.2 \
  --no-systemd"
```

**Note**: Use `--no-systemd` for Docker-based Vagrant providers where systemd may not work properly.

### Docker Container Installation

```bash
docker exec -it mycontainer bash -c "curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh | bash -s -- \
  --name container-agent \
  --master 192.168.1.10:50053 \
  --no-systemd"
```

## Troubleshooting

### Agent Not Registering

Check if agent can reach master:

```bash
# Test connectivity
telnet MASTER_IP 50053

# Check agent logs
sudo journalctl -u sloth-runner-agent -n 50

# Check agent is listening
sudo netstat -tulpn | grep sloth-runner
```

### Service Failed to Start

```bash
# Check service status
sudo systemctl status sloth-runner-agent

# View full logs
sudo journalctl -u sloth-runner-agent -n 100 --no-pager

# Test command manually
/usr/local/bin/sloth-runner agent start \
  --name myagent \
  --master localhost:50053
```

### Permission Issues

If running as specific user:

```bash
# Create user
sudo useradd -r -s /bin/false slothrunner

# Grant necessary permissions
sudo usermod -aG docker slothrunner  # If using Docker
```

## Uninstall

### Remove Service

```bash
# Stop and disable service
sudo systemctl stop sloth-runner-agent
sudo systemctl disable sloth-runner-agent

# Remove service file
sudo rm /etc/systemd/system/sloth-runner-agent.service

# Reload systemd
sudo systemctl daemon-reload
```

### Remove Binary

```bash
# Remove binary
sudo rm /usr/local/bin/sloth-runner

# Or if installed to user directory
rm ~/.local/bin/sloth-runner
```

### Delete from Master

```bash
# Delete agent from master registry
sloth-runner agent delete myagent --yes
```

## Advanced Configuration

### Custom Environment Variables

Edit the service file to add environment variables:

```bash
sudo systemctl edit sloth-runner-agent
```

Add:

```ini
[Service]
Environment="HTTPS_PROXY=http://proxy:8080"
Environment="LOG_LEVEL=debug"
```

### Custom Working Directory

```ini
[Service]
WorkingDirectory=/opt/sloth-runner
```

### Resource Limits

```ini
[Service]
CPUQuota=50%
MemoryLimit=512M
TasksMax=100
```

## See Also

- [Agent Documentation](https://chalkan3.github.io/sloth-runner/en/agents/)
- [Installation Guide](https://chalkan3.github.io/sloth-runner/en/installation/)
- [Quick Start](https://chalkan3.github.io/sloth-runner/en/quick-start/)
- [Troubleshooting](https://chalkan3.github.io/sloth-runner/en/troubleshooting/)

## License

Same as Sloth Runner - See [LICENSE](LICENSE)
