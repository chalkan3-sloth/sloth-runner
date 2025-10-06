# Sloth Runner Documentation

Comprehensive command-line documentation for Sloth Runner - a distributed task automation and orchestration system.

## Documentation Structure

This documentation is organized as man-page style references for each command. Each document provides detailed usage information, options, examples, and best practices.

## Quick Start

```bash
# Basic workflow execution
sloth-runner run production --file deployments/deploy.sloth --yes

# Install a remote agent
sloth-runner agent install prod-web-01 --ssh-host 192.168.1.100 --ssh-user root

# Register an event hook
sloth-runner hook add task_failure_alert --file hooks/alert.lua --event task.failed

# View recent events
sloth-runner events list --limit 20
```

## Core Documentation

### [sloth-runner(1)](sloth-runner.md)
Main command overview and system architecture. Start here for a high-level understanding of Sloth Runner.

**Topics covered:**
- System architecture and key concepts
- Workflows, agents, hooks, events, and stacks
- Environment variables and configuration
- Quick examples and common patterns

## Command Reference

### Workflow Execution

#### [sloth-runner-run(1)](commands/run.md)
Execute workflows and tasks locally or on remote agents.

**Key features:**
- Stack-based state management
- Local and remote execution
- Parameter passing and values files
- Interactive mode
- Multiple output formats

**Common use cases:**
```bash
# Deploy application
sloth-runner run production --file deploy.sloth --delegate-to web-01 --yes

# Run with parameters
sloth-runner run staging --file deploy.sloth --param version=v2.0.0 --yes

# Interactive execution
sloth-runner run dev --file workflow.sloth --interactive
```

### Agent Management

#### [sloth-runner-agent(1)](commands/agent.md)
Manage distributed remote agents for task execution.

**Subcommands:**
- `install` - Bootstrap agents on remote hosts via SSH
- `list` - List all registered agents
- `get` - Get detailed agent information
- `start` - Start an agent locally
- `stop` - Stop a running agent
- `delete` - Remove an agent from registry
- `update` - Update agent to latest version
- `exec` - Execute commands on agents
- `modules` - Check available tools on agents
- `metrics` - View agent metrics and telemetry

**Common use cases:**
```bash
# Install agent on remote server
sloth-runner agent install prod-web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user root \
  --master 192.168.1.10:50051

# Check agent status
sloth-runner agent list

# Execute command on agent
sloth-runner agent exec prod-web-01 "systemctl status nginx"

# Check available modules
sloth-runner agent modules prod-web-01
```

### Event and Hook System

#### [sloth-runner-hook(1)](commands/hook.md)
Manage event hooks for reactive automation.

**Subcommands:**
- `add` - Register a new event hook
- `list` - List all hooks
- `get` - Get hook details in JSON
- `show` - Show detailed hook information with history
- `delete` - Remove a hook
- `enable` - Enable a disabled hook
- `disable` - Temporarily disable a hook
- `test` - Test a hook with mock data

**Common use cases:**
```bash
# Register hook for task failures
sloth-runner hook add task_failure_alert \
  --file hooks/alert.lua \
  --event task.failed \
  --description "Send alerts on task failure"

# List all hooks
sloth-runner hook list

# Test hook before deployment
sloth-runner hook test task_failure_alert \
  --data '{"task":{"task_name":"test","error":"connection timeout"}}'

# View hook execution history
sloth-runner hook show task_failure_alert --history
```

#### [sloth-runner-events(1)](commands/events.md)
View and manage the event queue.

**Subcommands:**
- `list` - List events with filtering
- `show` - Show detailed event information
- `get` - Get event details in JSON
- `delete` - Delete a specific event
- `cleanup` - Remove old events

**Common use cases:**
```bash
# List recent events
sloth-runner events list --limit 50

# Filter by event type
sloth-runner events list --type task.failed

# Show event details with hook executions
sloth-runner events show evt_abc123

# Clean up old events
sloth-runner events cleanup --hours 168
```

## Event Types Reference

### Task Events
- `task.started` - Task execution began
- `task.completed` - Task finished successfully
- `task.failed` - Task execution failed
- `task.timeout` - Task exceeded time limit
- `task.retrying` - Task retry attempt
- `task.cancelled` - Task was cancelled

### Agent Events
- `agent.registered` - New agent joined
- `agent.connected` - Agent established connection
- `agent.disconnected` - Agent lost connection
- `agent.heartbeat_failed` - Agent missed heartbeat
- `agent.updated` - Agent software updated
- `agent.version_mismatch` - Agent version incompatible
- `agent.resource_high` - Agent resources critically high

### Workflow Events
- `workflow.started` - Workflow execution began
- `workflow.completed` - Workflow finished
- `workflow.failed` - Workflow execution failed
- `workflow.paused` - Workflow paused
- `workflow.resumed` - Workflow resumed
- `workflow.cancelled` - Workflow cancelled

### System Events
- `system.startup` - System initialized
- `system.shutdown` - System shutting down
- `system.error` - System-level error
- `system.warning` - System warning

### Custom Events
- `custom` - User-defined events from workflows

## Common Workflows

### Complete Deployment Pipeline

```bash
# 1. Install agents on target servers
sloth-runner agent install prod-web-01 --ssh-host 192.168.1.100 --ssh-user root
sloth-runner agent install prod-web-02 --ssh-host 192.168.1.101 --ssh-user root
sloth-runner agent install prod-db --ssh-host 192.168.1.150 --ssh-user root

# 2. Register deployment hooks
sloth-runner hook add deployment_success \
  --file hooks/notify-success.lua \
  --event task.completed

sloth-runner hook add deployment_failure \
  --file hooks/alert-failure.lua \
  --event task.failed

# 3. Create production stack
sloth-runner stack new production --description "Production environment"

# 4. Execute deployment
sloth-runner run production \
  --file workflows/deploy-app.sloth \
  --values values/production.yaml \
  --delegate-to prod-web-01 \
  --delegate-to prod-web-02 \
  --yes

# 5. Monitor deployment
sloth-runner events list --type task.completed --limit 10
sloth-runner stack history production
```

### Monitoring Setup

```bash
# 1. Create monitoring hooks
sloth-runner hook add task_monitor \
  --file hooks/task-logger.lua \
  --event task.started

sloth-runner hook add failure_alert \
  --file hooks/send-alert.lua \
  --event task.failed

sloth-runner hook add agent_monitor \
  --file hooks/agent-health.lua \
  --event agent.disconnected

# 2. List active hooks
sloth-runner hook list

# 3. Monitor events in real-time
watch -n 2 'sloth-runner events list --limit 10'

# 4. View hook execution history
sloth-runner hook show failure_alert --history --limit 50
```

### Multi-Environment Testing

```bash
# Development
sloth-runner stack new dev --description "Development environment"
sloth-runner run dev --file workflows/tests.sloth --param suite=unit --yes

# Staging
sloth-runner stack new staging --description "Staging environment"
sloth-runner run staging \
  --file workflows/tests.sloth \
  --param suite=integration \
  --delegate-to staging-server \
  --yes

# Production smoke tests
sloth-runner stack new production --description "Production environment"
sloth-runner run production \
  --file workflows/tests.sloth \
  --param suite=smoke \
  --delegate-to prod-web-01 \
  --yes
```

## Workflow File Examples

### Basic Task

```lua
local deploy = task("deploy_application")
    :description("Deploy application to server")
    :command(function(this, params)
        log.info("Deploying version: " .. params.version)

        -- Copy application binary
        ssh.copy({
            src = "build/app",
            dest = params.target_host .. ":/opt/app/app",
            user = "deploy"
        })

        -- Restart service
        ssh.exec({
            host = params.target_host,
            user = "deploy",
            command = "systemctl restart app"
        })

        log.info("Deployment completed")
        return true
    end)
    :build()

return {deploy}
```

### Tasks with Dependencies

```lua
local build = task("build")
    :description("Build application")
    :command(function()
        exec.run("go build -o app ./cmd/app")
        return true
    end)
    :build()

local test = task("test")
    :description("Run tests")
    :depends_on("build")
    :command(function()
        exec.run("go test ./...")
        return true
    end)
    :build()

local deploy = task("deploy")
    :description("Deploy application")
    :depends_on("test")
    :command(function(this, params)
        ssh.copy({
            src = "app",
            dest = params.host .. ":/opt/app/app"
        })
        return true
    end)
    :build()

return {build, test, deploy}
```

### Parallel Task Execution

```lua
-- Health check that runs on multiple agents in parallel
local health = task("health_check")
    :description("Check service health")
    :command(function()
        local result = http.get({url = "http://localhost:8080/health"})
        if result.code ~= 200 then
            error("Health check failed")
        end
        return true
    end)
    :build()

return {health}
```

Execute on multiple agents:
```bash
sloth-runner run prod \
  --file health.sloth \
  --delegate-to web-01 \
  --delegate-to web-02 \
  --delegate-to web-03 \
  --yes
```

## Hook Script Examples

### Task Failure Alert

```lua
-- hooks/task-failure-alert.lua
function on_event()
    local task = event.task

    local message = string.format(
        "ALERT: Task '%s' failed on %s\\nError: %s",
        task.task_name,
        task.agent_name,
        task.error or "Unknown error"
    )

    -- Log to file
    file_ops.lineinfile({
        path = "/var/log/sloth/failures.log",
        line = "[" .. os.date("%Y-%m-%d %H:%M:%S") .. "] " .. message,
        create = true
    })

    -- Send to monitoring system
    http.post({
        url = "https://monitoring.example.com/webhook",
        body = message
    })

    print("Alert sent for task: " .. task.task_name)
    return true
end
```

### Agent Health Monitor

```lua
-- hooks/agent-health-monitor.lua
function on_event()
    local agent = event.agent

    log.warn("Agent " .. agent.name .. " disconnected!")

    -- Attempt automatic recovery
    exec.run("sloth-runner agent start " .. agent.name .. " --ssh-reconnect")

    -- Alert operations team
    http.post({
        url = "https://alerts.example.com/agent-down",
        body = string.format("Agent %s disconnected at %s",
            agent.name, os.date("%Y-%m-%d %H:%M:%S"))
    })

    return true
end
```

## Available Modules in Workflows and Hooks

### Execution
- **exec** - Command execution
- **ssh** - SSH operations

### File Operations
- **file_ops** - File manipulation (lineinfile, blockinfile, copy, template)

### Network
- **http** - HTTP client
- **ssh** - SSH client and file transfer

### Infrastructure
- **docker** - Docker container management
- **incus** - Incus/LXD container and VM management
- **systemd** - System service management
- **pkg** - Package management (apt, yum, apk)

### Utilities
- **log** - Logging functions
- **json** - JSON encoding/decoding
- **os** - OS utilities
- **string** - String manipulation

## Configuration Files

### Agent Configuration
```yaml
# ~/.sloth-runner/config.yaml
master:
  address: "192.168.1.10:50051"

agent:
  name: "prod-web-01"
  port: 50051
  bind_address: "0.0.0.0"
  telemetry:
    enabled: true
    port: 9090
```

### Values File
```yaml
# values/production.yaml
version: "v2.1.0"
environment: "production"
replicas: 3

database:
  host: "prod-db.internal"
  port: 5432
  name: "myapp_prod"

deployment:
  strategy: "rolling"
  max_unavailable: 1
```

## Troubleshooting

### View Logs
```bash
# Event logs
sloth-runner events list --type task.failed --limit 20

# Hook execution logs
sloth-runner hook show hook_name --history

# Stack execution history
sloth-runner stack history production
```

### Debug Mode
```bash
# Enable debug logging
sloth-runner run dev --file workflow.sloth --debug --yes

# Test connectivity
sloth-runner agent exec prod-web-01 "echo test"

# Check agent modules
sloth-runner agent modules prod-web-01
```

### Common Issues

**Agent not connecting:**
```bash
# Check agent status
sloth-runner agent list

# Restart agent on remote host
ssh root@192.168.1.100 "systemctl restart sloth-runner-agent"

# Check agent logs
ssh root@192.168.1.100 "journalctl -u sloth-runner-agent -n 50"
```

**Hook not triggering:**
```bash
# Check hook is enabled
sloth-runner hook list

# Test hook manually
sloth-runner hook test hook_name --data '{"task":{"task_name":"test"}}'

# Check recent events
sloth-runner events list --limit 10
```

**Task delegation failing:**
```bash
# Verify agent is online
sloth-runner agent list

# Check agent modules
sloth-runner agent modules agent_name

# Test command execution
sloth-runner agent exec agent_name "echo test"
```

## Best Practices

### 1. Use Stacks for Environment Isolation
```bash
sloth-runner stack new dev
sloth-runner stack new staging
sloth-runner stack new production
```

### 2. Implement Comprehensive Hooks
```bash
# Monitor all task lifecycle events
sloth-runner hook add task_start --event task.started --file hooks/logger.lua
sloth-runner hook add task_complete --event task.completed --file hooks/logger.lua
sloth-runner hook add task_fail --event task.failed --file hooks/alert.lua
```

### 3. Regular Maintenance
```bash
# Daily event cleanup
sloth-runner events cleanup --hours 168

# Weekly agent updates
sloth-runner agent update prod-web-01
```

### 4. Parameter Files for Complex Deployments
```bash
# Use values files instead of inline parameters
sloth-runner run production \
  --file deploy.sloth \
  --values values/production.yaml \
  --yes
```

### 5. Test in Lower Environments First
```bash
# Test in dev
sloth-runner run dev --file workflow.sloth --yes

# Promote to staging
sloth-runner run staging --file workflow.sloth --yes

# Deploy to production
sloth-runner run production --file workflow.sloth --yes
```

## Getting Help

### Command Help
```bash
# General help
sloth-runner --help

# Command-specific help
sloth-runner run --help
sloth-runner agent --help
sloth-runner hook --help

# Subcommand help
sloth-runner agent install --help
sloth-runner hook add --help
```

### Version Information
```bash
sloth-runner --version
```

## Additional Resources

- **Main Documentation:** [sloth-runner(1)](sloth-runner.md)
- **GitHub Repository:** https://github.com/yourusername/task-runner
- **Issue Tracker:** https://github.com/yourusername/task-runner/issues

## Index of All Commands

### Main Commands
- [sloth-runner](sloth-runner.md) - Main command overview
- [run](commands/run.md) - Execute workflows

### Agent Management
- [agent](commands/agent.md) - Manage remote agents
- agent install - Install agent on remote host
- agent list - List all agents
- agent get - Get agent details
- agent start - Start an agent
- agent stop - Stop an agent
- agent delete - Remove an agent
- agent update - Update agent version
- agent exec - Execute commands on agent
- agent modules - Check available modules
- agent metrics - View agent metrics

### Event and Hook System
- [hook](commands/hook.md) - Manage event hooks
- hook add - Register new hook
- hook list - List all hooks
- hook get - Get hook in JSON
- hook show - Show hook with history
- hook delete - Remove hook
- hook enable - Enable hook
- hook disable - Disable hook
- hook test - Test hook

- [events](commands/events.md) - Manage event queue
- events list - List events
- events show - Show event details
- events get - Get event in JSON
- events delete - Delete event
- events cleanup - Remove old events

## Copyright

Copyright © 2025 Sloth Runner Development Team. Released under MIT License.
