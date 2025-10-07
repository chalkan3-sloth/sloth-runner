# 📚 Complete CLI Commands Reference

## Overview

Sloth Runner offers a complete and powerful command-line interface (CLI) for managing workflows, agents, modules, hooks, events, and much more. This documentation covers **all** available commands with practical examples.

---

## 🎯 Main Commands

### `run` - Execute Workflows

Executes a Sloth workflow from a file.

```bash
# Basic syntax
sloth-runner run <workflow-name> --file <file.sloth> [options]

# Examples
sloth-runner run deploy --file deploy.sloth
sloth-runner run deploy --file deploy.sloth --yes                    # Non-interactive mode
sloth-runner run deploy --file deploy.sloth --group production       # Execute specific group
sloth-runner run deploy --file deploy.sloth --delegate-to agent1     # Delegate to agent
sloth-runner run deploy --file deploy.sloth --delegate-to agent1 --delegate-to agent2  # Multiple agents
sloth-runner run deploy --file deploy.sloth --values vars.yaml       # Pass variables
sloth-runner run deploy --file deploy.sloth --var "env=production"   # Inline variable
```

**Options:**
- `--file, -f` - Path to Sloth file
- `--yes, -y` - Non-interactive mode (no confirmation)
- `--group, -g` - Execute only a specific group
- `--delegate-to` - Delegate execution to remote agent(s)
- `--values` - YAML file with variables
- `--var` - Define inline variable (can use multiple times)
- `--verbose, -v` - Verbose mode

---

## 🤖 Agent Management

### `agent list` - List Agents

Lists all agents registered with the master server.

```bash
# Syntax
sloth-runner agent list [options]

# Examples
sloth-runner agent list                    # List all agents
sloth-runner agent list --format json      # JSON output
sloth-runner agent list --format yaml      # YAML output
sloth-runner agent list --status active    # Only active agents
```

**Options:**
- `--format` - Output format: table (default), json, yaml
- `--status` - Filter by status: active, inactive, all

---

### `agent get` - Agent Details

Gets detailed information about a specific agent.

```bash
# Syntax
sloth-runner agent get <agent-name> [options]

# Examples
sloth-runner agent get web-server-01
sloth-runner agent get web-server-01 --format json
sloth-runner agent get web-server-01 --show-metrics       # Include metrics
```

**Options:**
- `--format` - Output format: table, json, yaml
- `--show-metrics` - Show agent metrics

---

### `agent install` - Install Remote Agent

Installs the Sloth Runner agent on a remote server via SSH.

```bash
# Syntax
sloth-runner agent install <agent-name> --ssh-host <host> --ssh-user <user> [options]

# Examples
sloth-runner agent install web-01 --ssh-host 192.168.1.100 --ssh-user root
sloth-runner agent install web-01 --ssh-host 192.168.1.100 --ssh-user root --ssh-port 2222
sloth-runner agent install web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user root \
  --master 192.168.1.1:50053 \
  --bind-address 0.0.0.0 \
  --port 50060 \
  --report-address 192.168.1.100:50060
```

**Options:**
- `--ssh-host` - SSH host of remote server (required)
- `--ssh-user` - SSH user (required)
- `--ssh-port` - SSH port (default: 22)
- `--ssh-key` - Path to SSH private key
- `--master` - Master server address (default: localhost:50053)
- `--bind-address` - Agent bind address (default: 0.0.0.0)
- `--port` - Agent port (default: 50060)
- `--report-address` - Address the agent reports to master

---

### `agent update` - Update Agent

Updates the agent binary to the latest version.

```bash
# Syntax
sloth-runner agent update <agent-name> [options]

# Examples
sloth-runner agent update web-01
sloth-runner agent update web-01 --version v1.2.3
sloth-runner agent update web-01 --restart           # Restart after update
```

**Options:**
- `--version` - Specific version (default: latest)
- `--restart` - Restart agent after update
- `--force` - Force update even if version is the same

---

### `agent modules` - Agent Modules

Lists or checks available modules on an agent.

```bash
# Syntax
sloth-runner agent modules <agent-name> [options]

# Examples
sloth-runner agent modules web-01                      # List all modules
sloth-runner agent modules web-01 --check pkg          # Check if 'pkg' module is available
sloth-runner agent modules web-01 --check docker      # Check if Docker is installed
sloth-runner agent modules web-01 --format json       # JSON output
```

**Options:**
- `--check` - Check specific module
- `--format` - Output format: table, json, yaml

---

### `agent start` - Start Agent

Starts the agent service locally.

```bash
# Syntax
sloth-runner agent start [options]

# Examples
sloth-runner agent start                                    # Start with default settings
sloth-runner agent start --master 192.168.1.1:50053         # Connect to specific master
sloth-runner agent start --port 50060                       # Use specific port
sloth-runner agent start --name my-agent                    # Set agent name
sloth-runner agent start --bind 0.0.0.0                     # Bind to all interfaces
sloth-runner agent start --foreground                       # Run in foreground
```

**Options:**
- `--master` - Master server address (default: localhost:50053)
- `--port` - Agent port (default: 50060)
- `--name` - Agent name (default: hostname)
- `--bind` - Bind address (default: 0.0.0.0)
- `--report-address` - Address the agent reports
- `--foreground` - Run in foreground (not daemon)

---

### `agent stop` - Stop Agent

Stops the agent service.

```bash
# Syntax
sloth-runner agent stop [options]

# Examples
sloth-runner agent stop                # Stop local agent
sloth-runner agent stop --name web-01  # Stop specific agent
```

---

### `agent restart` - Restart Agent

Restarts the agent service.

```bash
# Syntax
sloth-runner agent restart [agent-name]

# Examples
sloth-runner agent restart               # Restart local agent
sloth-runner agent restart web-01        # Restart remote agent
```

---

### `agent metrics` - Agent Metrics

View agent performance and resource metrics.

```bash
# Syntax
sloth-runner agent metrics <agent-name> [options]

# Examples
sloth-runner agent metrics web-01
sloth-runner agent metrics web-01 --format json
sloth-runner agent metrics web-01 --watch              # Continuous updates
sloth-runner agent metrics web-01 --interval 5         # 5 second interval
```

**Options:**
- `--format` - Format: table, json, yaml, prometheus
- `--watch` - Continuous updates
- `--interval` - Update interval in seconds (default: 2)

---

### `agent metrics grafana` - Grafana Dashboard

Generates and displays Grafana dashboard for an agent.

```bash
# Syntax
sloth-runner agent metrics grafana <agent-name> [options]

# Examples
sloth-runner agent metrics grafana web-01
sloth-runner agent metrics grafana web-01 --export dashboard.json
```

**Options:**
- `--export` - Export dashboard to JSON file

---

## 📦 Sloth Management (Saved Workflows)

### `sloth list` - List Sloths

Lists all workflows saved in the local repository.

```bash
# Syntax
sloth-runner sloth list [options]

# Examples
sloth-runner sloth list                   # List all
sloth-runner sloth list --active          # Only active sloths
sloth-runner sloth list --inactive        # Only inactive sloths
sloth-runner sloth list --format json     # JSON output
```

**Options:**
- `--active` - Only active sloths
- `--inactive` - Only inactive sloths
- `--format` - Format: table, json, yaml

---

### `sloth add` - Add Sloth

Adds a new workflow to the repository.

```bash
# Syntax
sloth-runner sloth add <name> --file <path> [options]

# Examples
sloth-runner sloth add deploy --file deploy.sloth
sloth-runner sloth add deploy --file deploy.sloth --description "Production deploy"
sloth-runner sloth add deploy --file deploy.sloth --tags "prod,deploy"
```

**Options:**
- `--file` - Path to Sloth file (required)
- `--description` - Sloth description
- `--tags` - Comma-separated tags

---

### `sloth get` - Get Sloth

Displays details of a specific sloth.

```bash
# Syntax
sloth-runner sloth get <name> [options]

# Examples
sloth-runner sloth get deploy
sloth-runner sloth get deploy --format json
sloth-runner sloth get deploy --show-content    # Show workflow content
```

**Options:**
- `--format` - Format: table, json, yaml
- `--show-content` - Show complete workflow content

---

### `sloth update` - Update Sloth

Updates an existing sloth.

```bash
# Syntax
sloth-runner sloth update <name> [options]

# Examples
sloth-runner sloth update deploy --file deploy-v2.sloth
sloth-runner sloth update deploy --description "New description"
sloth-runner sloth update deploy --tags "prod,deploy,updated"
```

**Options:**
- `--file` - New Sloth file
- `--description` - New description
- `--tags` - New tags

---

### `sloth remove` - Remove Sloth

Removes a sloth from the repository.

```bash
# Syntax
sloth-runner sloth remove <name>

# Examples
sloth-runner sloth remove deploy
sloth-runner sloth remove deploy --force    # Remove without confirmation
```

**Options:**
- `--force` - Remove without confirmation

---

### `sloth activate` - Activate Sloth

Activates a deactivated sloth.

```bash
# Syntax
sloth-runner sloth activate <name>

# Examples
sloth-runner sloth activate deploy
```

---

### `sloth deactivate` - Deactivate Sloth

Deactivates a sloth (doesn't remove, just marks as inactive).

```bash
# Syntax
sloth-runner sloth deactivate <name>

# Examples
sloth-runner sloth deactivate deploy
```

---

## 🎣 Hook Management

### `hook list` - List Hooks

Lists all registered hooks.

```bash
# Syntax
sloth-runner hook list [options]

# Examples
sloth-runner hook list
sloth-runner hook list --format json
sloth-runner hook list --event workflow.started    # Filter by event
```

**Options:**
- `--format` - Format: table, json, yaml
- `--event` - Filter by event type

---

### `hook add` - Add Hook

Adds a new hook.

```bash
# Syntax
sloth-runner hook add <name> --event <event> --script <path> [options]

# Examples
sloth-runner hook add notify-slack --event workflow.completed --script notify.sh
sloth-runner hook add backup --event task.completed --script backup.lua
sloth-runner hook add validate --event workflow.started --script validate.lua --priority 10
```

**Options:**
- `--event` - Event type (required)
- `--script` - Script path (required)
- `--priority` - Execution priority (default: 0)
- `--enabled` - Hook enabled (default: true)

**Available events:**
- `workflow.started`
- `workflow.completed`
- `workflow.failed`
- `task.started`
- `task.completed`
- `task.failed`
- `agent.connected`
- `agent.disconnected`

---

### `hook remove` - Remove Hook

Removes a hook.

```bash
# Syntax
sloth-runner hook remove <name>

# Examples
sloth-runner hook remove notify-slack
sloth-runner hook remove notify-slack --force
```

---

### `hook enable` - Enable Hook

Enables a disabled hook.

```bash
# Syntax
sloth-runner hook enable <name>

# Examples
sloth-runner hook enable notify-slack
```

---

### `hook disable` - Disable Hook

Disables a hook.

```bash
# Syntax
sloth-runner hook disable <name>

# Examples
sloth-runner hook disable notify-slack
```

---

### `hook test` - Test Hook

Tests hook execution.

```bash
# Syntax
sloth-runner hook test <name> [options]

# Examples
sloth-runner hook test notify-slack
sloth-runner hook test notify-slack --payload '{"message": "test"}'
```

**Options:**
- `--payload` - JSON with test data

---

## 📡 Event Management

### `events list` - List Events

Lists recent system events.

```bash
# Syntax
sloth-runner events list [options]

# Examples
sloth-runner events list
sloth-runner events list --limit 50               # Last 50 events
sloth-runner events list --type workflow.started  # Filter by type
sloth-runner events list --since 1h               # Events from last hour
sloth-runner events list --format json
```

**Options:**
- `--limit` - Maximum number of events (default: 100)
- `--type` - Filter by event type
- `--since` - Filter by time (e.g., 1h, 30m, 24h)
- `--format` - Format: table, json, yaml

---

### `events watch` - Monitor Events

Monitors events in real-time.

```bash
# Syntax
sloth-runner events watch [options]

# Examples
sloth-runner events watch
sloth-runner events watch --type workflow.completed    # Only completed workflow events
sloth-runner events watch --filter "status=success"    # With filter
```

**Options:**
- `--type` - Filter by event type
- `--filter` - Filter expression

---

## 🗄️ Database Management

### `db backup` - Database Backup

Creates SQLite database backup.

```bash
# Syntax
sloth-runner db backup [options]

# Examples
sloth-runner db backup
sloth-runner db backup --output /backup/sloth-backup.db
sloth-runner db backup --compress                     # Compress with gzip
```

**Options:**
- `--output` - Backup file path
- `--compress` - Compress backup

---

### `db restore` - Restore Database

Restores database from backup.

```bash
# Syntax
sloth-runner db restore <backup-file> [options]

# Examples
sloth-runner db restore /backup/sloth-backup.db
sloth-runner db restore /backup/sloth-backup.db.gz --decompress
```

**Options:**
- `--decompress` - Decompress gzip backup

---

### `db vacuum` - Optimize Database

Optimizes and compacts the SQLite database.

```bash
# Syntax
sloth-runner db vacuum

# Examples
sloth-runner db vacuum
```

---

### `db stats` - Database Statistics

Shows database statistics.

```bash
# Syntax
sloth-runner db stats [options]

# Examples
sloth-runner db stats
sloth-runner db stats --format json
```

**Options:**
- `--format` - Format: table, json, yaml

---

## 🌐 SSH Management

### `ssh list` - List SSH Connections

Lists saved SSH connections.

```bash
# Syntax
sloth-runner ssh list [options]

# Examples
sloth-runner ssh list
sloth-runner ssh list --format json
```

**Options:**
- `--format` - Format: table, json, yaml

---

### `ssh add` - Add SSH Connection

Adds a new SSH connection.

```bash
# Syntax
sloth-runner ssh add <name> --host <host> --user <user> [options]

# Examples
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu --port 2222
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu --key ~/.ssh/id_rsa
```

**Options:**
- `--host` - SSH host (required)
- `--user` - SSH user (required)
- `--port` - SSH port (default: 22)
- `--key` - SSH private key path

---

### `ssh remove` - Remove SSH Connection

Removes a saved SSH connection.

```bash
# Syntax
sloth-runner ssh remove <name>

# Examples
sloth-runner ssh remove web-server
```

---

### `ssh test` - Test SSH Connection

Tests an SSH connection.

```bash
# Syntax
sloth-runner ssh test <name>

# Examples
sloth-runner ssh test web-server
```

---

## 📋 Modules

### `modules list` - List Modules

Lists all available modules.

```bash
# Syntax
sloth-runner modules list [options]

# Examples
sloth-runner modules list
sloth-runner modules list --format json
sloth-runner modules list --category cloud         # Filter by category
```

**Options:**
- `--format` - Format: table, json, yaml
- `--category` - Filter by category

---

### `modules info` - Module Information

Shows detailed information about a module.

```bash
# Syntax
sloth-runner modules info <module-name>

# Examples
sloth-runner modules info pkg
sloth-runner modules info docker
sloth-runner modules info terraform
```

---

## 🖥️ Server and UI

### `server` - Start Master Server

Starts the master server (gRPC).

```bash
# Syntax
sloth-runner server [options]

# Examples
sloth-runner server                          # Start on default port (50053)
sloth-runner server --port 50053             # Specify port
sloth-runner server --bind 0.0.0.0           # Bind to all interfaces
sloth-runner server --tls-cert cert.pem --tls-key key.pem  # With TLS
```

**Options:**
- `--port` - Server port (default: 50053)
- `--bind` - Bind address (default: 0.0.0.0)
- `--tls-cert` - TLS certificate
- `--tls-key` - TLS private key

---

### `ui` - Start Web UI

Starts the web interface.

```bash
# Syntax
sloth-runner ui [options]

# Examples
sloth-runner ui                              # Start on default port (8080)
sloth-runner ui --port 8080                  # Specify port
sloth-runner ui --bind 0.0.0.0               # Bind to all interfaces
```

**Options:**
- `--port` - Web UI port (default: 8080)
- `--bind` - Bind address (default: 0.0.0.0)

---

### `terminal` - Interactive Terminal

Opens interactive terminal to a remote agent.

```bash
# Syntax
sloth-runner terminal <agent-name>

# Examples
sloth-runner terminal web-01
```

---

## 🔧 Utilities

### `version` - Version

Shows the Sloth Runner version.

```bash
# Syntax
sloth-runner version

# Examples
sloth-runner version
sloth-runner version --format json
```

---

### `completion` - Auto-completion

Generates auto-completion scripts for the shell.

```bash
# Syntax
sloth-runner completion <shell>

# Examples
sloth-runner completion bash > /etc/bash_completion.d/sloth-runner
sloth-runner completion zsh > ~/.zsh/completion/_sloth-runner
sloth-runner completion fish > ~/.config/fish/completions/sloth-runner.fish
```

**Supported shells:** bash, zsh, fish, powershell

---

### `doctor` - Diagnostics

Runs system and configuration diagnostics.

```bash
# Syntax
sloth-runner doctor [options]

# Examples
sloth-runner doctor
sloth-runner doctor --format json
sloth-runner doctor --verbose             # Detailed output
```

**Options:**
- `--format` - Format: text, json
- `--verbose` - Detailed output

---

## 🔐 Environment Variables

Sloth Runner uses the following environment variables:

```bash
# Master server address
export SLOTH_RUNNER_MASTER_ADDR="192.168.1.1:50053"

# Agent port
export SLOTH_RUNNER_AGENT_PORT="50060"

# Web UI port
export SLOTH_RUNNER_UI_PORT="8080"

# Database path
export SLOTH_RUNNER_DB_PATH="~/.sloth-runner/sloth.db"

# Log level
export SLOTH_RUNNER_LOG_LEVEL="info"  # debug, info, warn, error

# Enable debug mode
export SLOTH_RUNNER_DEBUG="true"
```

---

## 📊 Common Usage Examples

### 1. Production Deploy with Delegation

```bash
sloth-runner run production-deploy \
  --file deployments/prod.sloth \
  --delegate-to web-01 \
  --delegate-to web-02 \
  --values prod-vars.yaml \
  --yes
```

### 2. Monitor Metrics of All Agents

```bash
# In one terminal
sloth-runner agent metrics web-01 --watch

# In another terminal
sloth-runner agent metrics web-02 --watch
```

### 3. Automated Backup

```bash
# Create compressed backup with timestamp
sloth-runner db backup \
  --output /backup/sloth-$(date +%Y%m%d-%H%M%S).db \
  --compress
```

### 4. Workflow with Notification Hook

```bash
# Add notification hook
sloth-runner hook add slack-notify \
  --event workflow.completed \
  --script /scripts/notify-slack.lua

# Run workflow (hook will be triggered automatically)
sloth-runner run deploy --file deploy.sloth --yes
```

### 5. Agent Installation on Multiple Servers

```bash
# Loop to install on multiple hosts
for host in 192.168.1.{10..20}; do
  sloth-runner agent install "agent-$host" \
    --ssh-host "$host" \
    --ssh-user ubuntu \
    --master 192.168.1.1:50053
done
```

---

## 🎓 Next Steps

- [📖 Modules Guide](modules-complete.md) - Complete documentation of all modules
- [🎨 Web UI](web-ui-complete.md) - Complete web interface guide
- [🎯 Advanced Examples](../en/advanced-examples.md) - Practical workflow examples
- [🏗️ Architecture](../architecture/sloth-runner-architecture.md) - System architecture

---

## 💡 Tips and Tricks

### Useful Aliases

Add to your `.bashrc` or `.zshrc`:

```bash
alias sr='sloth-runner'
alias sra='sloth-runner agent'
alias srr='sloth-runner run'
alias srl='sloth-runner sloth list'
alias srui='sloth-runner ui --port 8080'
```

### Auto-completion

```bash
# Bash
sloth-runner completion bash > /etc/bash_completion.d/sloth-runner
source /etc/bash_completion.d/sloth-runner

# Zsh
sloth-runner completion zsh > ~/.zsh/completion/_sloth-runner
```

### Debug Mode

```bash
export SLOTH_RUNNER_DEBUG=true
export SLOTH_RUNNER_LOG_LEVEL=debug
sloth-runner run deploy --file deploy.sloth --verbose
```

---

**Last updated:** 2025-10-07
