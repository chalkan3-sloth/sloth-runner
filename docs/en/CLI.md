# 🚀 Sloth Runner CLI Reference

Complete command-line interface reference for Sloth Runner - the AI-powered GitOps task orchestration platform.

---

## Overview

Sloth Runner provides a comprehensive CLI for task execution, agent management, scheduling, and workflow orchestration.

**Main Commands:**

```bash
sloth-runner [command] [flags]
```

| Command | Description |
|---------|-------------|
| `run` | Execute tasks from workflow files |
| `agent` | Manage distributed agents |
| `master` | Start the master coordination server |
| `scheduler` | Manage scheduled tasks |
| `stack` | Manage workflow stacks and state |
| `ui` | Start the web dashboard |
| `workflow` | Create and manage workflows |
| `list` | List tasks from workflow files |
| `version` | Show version information |

---

## `sloth-runner run`

Execute tasks defined in `.sloth` workflow files with multiple output formats and state persistence.

### Usage

```bash
sloth-runner run [file.sloth|stack-name] [flags]
```

### Flags

| Flag | Type | Description |
|------|------|-------------|
| `-f, --file` | string | Path to the Lua/Sloth task file |
| `-o, --output` | string | Output style: `basic`, `enhanced`, `rich`, `modern`, `json` (default: `basic`) |
| `-v, --values` | string | Path to values file (YAML/JSON) for parameterization |
| `--interactive` | bool | Run in interactive mode with prompts |
| `--yes` | bool | Skip confirmation prompts |

### Output Styles

- **basic**: Simple text output
- **enhanced**: Colored output with icons
- **rich**: Detailed output with progress bars
- **modern**: Modern UI with animations
- **json**: Machine-readable JSON output

### Examples

```bash
# Run with modern output style
sloth-runner run -f deploy.sloth -o modern

# Run with values file
sloth-runner run -f infra.sloth -v prod-values.yaml

# Run from stack
sloth-runner run prod-stack --yes

# Interactive mode
sloth-runner run -f tasks.sloth --interactive

# JSON output for CI/CD
sloth-runner run -f ci.sloth -o json
```

---

## `sloth-runner agent`

Manage distributed agents for remote task execution.

### Subcommands

#### `agent start`

Start an agent in agent mode to accept tasks from master server.

```bash
sloth-runner agent start [flags]
```

**Flags:**
- `--master string`: Master server address (default: `localhost:50053`)
- `--name string`: Agent name identifier
- `--tags string`: Comma-separated tags for agent capabilities
- `--daemon`: Run as background daemon

**Example:**
```bash
# Start agent with tags
sloth-runner agent start --master master.example.com:50053 \
  --name prod-agent-1 \
  --tags linux,docker,aws

# Start as daemon
sloth-runner agent start --daemon --name bg-agent
```

#### `agent list`

List all registered agents with their status.

```bash
sloth-runner agent list [flags]
```

**Flags:**
- `--master string`: Master server address

**Example:**
```bash
sloth-runner agent list --master master.example.com:50053
```

#### `agent exec`

Execute a command on a remote agent.

```bash
sloth-runner agent exec <agent_name> <command> [flags]
```

**Flags:**
- `--master string`: Master server address (or use SLOTH_RUNNER_MASTER_ADDR env var)
- `-o, --output string`: Output format: text or json (default: text)

**Example:**
```bash
# Using --master flag
sloth-runner agent exec prod-agent-1 "docker ps" --master master.example.com:50053

# Using environment variable
SLOTH_RUNNER_MASTER_ADDR=master.example.com:50053 sloth-runner agent exec prod-agent-1 "docker ps"
```

#### `agent stop`

Stop a remote agent gracefully.

```bash
sloth-runner agent stop [flags]
```

**Flags:**
- `--agent string`: Agent name to stop
- `--master string`: Master server address

#### `agent delete`

Delete an agent from the registry.

```bash
sloth-runner agent delete [flags]
```

---

## `sloth-runner master`

Start the master coordination server for managing distributed agents.

### Usage

```bash
sloth-runner master [flags]
```

### Flags

| Flag | Type | Description |
|------|------|-------------|
| `-p, --port` | int | Port to listen on (default: `50053`) |
| `--daemon` | bool | Run as background daemon |
| `--debug` | bool | Enable debug logging |

### Examples

```bash
# Start master server
sloth-runner master --port 50053

# Start as daemon with debug
sloth-runner master --daemon --debug

# Custom port
sloth-runner master --port 9000
```

---

## `sloth-runner scheduler`

Manage scheduled tasks for automated execution.

### Subcommands

#### `scheduler enable`

Enable the scheduler service.

```bash
sloth-runner scheduler enable
```

#### `scheduler disable`

Disable the scheduler service.

```bash
sloth-runner scheduler disable
```

#### `scheduler list`

List all scheduled tasks with their configuration.

```bash
sloth-runner scheduler list [flags]
```

**Output:**
- Task name
- Schedule (cron expression)
- Next run time
- Status (enabled/disabled)

**Example:**
```bash
sloth-runner scheduler list
```

#### `scheduler delete`

Delete a scheduled task.

```bash
sloth-runner scheduler delete [task-name]
```

---

## `sloth-runner stack`

Manage workflow stacks for state persistence and environment isolation.

### Subcommands

#### `stack new`

Create a new workflow stack.

```bash
sloth-runner stack new [stack-name] [flags]
```

**Flags:**
- `-f, --file string`: Workflow file to associate
- `--description string`: Stack description

**Example:**
```bash
sloth-runner stack new prod-infra \
  -f infrastructure.sloth \
  --description "Production infrastructure stack"
```

#### `stack list`

List all workflow stacks.

```bash
sloth-runner stack list
```

**Output:**
- Stack name
- Workflow file
- State status
- Last updated

#### `stack show`

Show detailed information about a stack.

```bash
sloth-runner stack show [stack-name]
```

**Output:**
- Stack configuration
- State variables
- Execution history
- Associated resources

#### `stack delete`

Delete a workflow stack and its state.

```bash
sloth-runner stack delete [stack-name] [flags]
```

**Flags:**
- `--force`: Force deletion without confirmation

---

## `sloth-runner ui`

Start the web-based dashboard for visual management.

### Usage

```bash
sloth-runner ui [flags]
```

### Flags

| Flag | Type | Description |
|------|------|-------------|
| `-p, --port` | int | Port for UI server (default: `8080`) |
| `--daemon` | bool | Run as background daemon |
| `--debug` | bool | Enable debug logging |

### Features

- 📊 Real-time task monitoring
- 🤖 Agent health dashboard
- 📅 Scheduler management
- 📦 Stack browser
- 📈 Metrics and analytics

### Examples

```bash
# Start UI on default port
sloth-runner ui

# Custom port
sloth-runner ui --port 3000

# Run as daemon
sloth-runner ui --daemon --port 8080
```

Access at: `http://localhost:8080`

---

## `sloth-runner workflow`

Create and manage workflow projects with scaffolding.

### Subcommands

#### `workflow init`

Initialize a new workflow project with templates.

```bash
sloth-runner workflow init [project-name] [flags]
```

**Flags:**
- `--template string`: Template to use (default: `basic`)
- `--path string`: Target directory

**Available Templates:**
- `basic`: Simple task workflow
- `cicd`: CI/CD pipeline
- `infra`: Infrastructure automation
- `gitops`: GitOps deployment

**Example:**
```bash
# Create CI/CD project
sloth-runner workflow init my-pipeline --template cicd

# Custom path
sloth-runner workflow init my-project --template infra --path ./projects/
```

#### `workflow list-templates`

List all available workflow templates.

```bash
sloth-runner workflow list-templates
```

---

## `sloth-runner list`

List tasks and task groups from a workflow file without execution.

### Usage

```bash
sloth-runner list [flags]
sloth-runner list [flags]
```

**Flags:**

*   `-f, --file string`: **(Required)** Path to the Lua task configuration file.
*   `-v, --values string`: Path to a YAML values file, in case your task definitions depend on it.

---

## `sloth-runner new`

Generates a new boilerplate Lua task definition file from a template.

**Usage:**
```bash
sloth-runner new <group-name> [flags]
```

**Arguments:**

*   `<group-name>`: The name of the main task group to be created in the file.

**Flags:**

*   `-t, --template string`: The template to use. Default is `simple`. Run `sloth-runner template list` to see all available options.
*   `-o, --output string`: The path to the output file. If not provided, the generated content will be printed to stdout.

```bash
sloth-runner list [flags]
```

### Flags

| Flag | Type | Description |
|------|------|-------------|
| `-f, --file` | string | Path to workflow file |

### Output

- Task groups
- Task names
- Descriptions
- Dependencies
- Conditions

### Example

```bash
sloth-runner list -f deploy.sloth
```

---

## `sloth-runner version`

Display version and build information.

### Usage

```bash
sloth-runner version
```

### Output

- Version number
- Git commit hash
- Build date
- Go version

---

## Global Flags

Available for all commands:

| Flag | Description |
|------|-------------|
| `-h, --help` | Show command help |
| `--debug` | Enable debug output |
| `--config string` | Config file path (default: `~/.sloth-runner/config.yaml`) |

---

## Configuration File

Sloth Runner supports configuration via `~/.sloth-runner/config.yaml`:

```yaml
# Master server settings
master:
  host: localhost
  port: 50053

# Agent settings
agent:
  name: my-agent
  tags:
    - linux
    - docker
  reconnect: true
  
# UI settings
ui:
  port: 8080
  theme: dark
  
# Scheduler settings
scheduler:
  enabled: true
  timezone: UTC
```

---

## Environment Variables

Override configuration with environment variables:

| Variable | Description |
|----------|-------------|
| `SLOTH_MASTER_HOST` | Master server host |
| `SLOTH_MASTER_PORT` | Master server port |
| `SLOTH_AGENT_NAME` | Agent identifier |
| `SLOTH_UI_PORT` | UI server port |
| `SLOTH_DEBUG` | Enable debug mode |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |
| `2` | Task execution failed |
| `3` | Configuration error |
| `4` | Network/connection error |
| `5` | State management error |

---

## Examples

### Complete CI/CD Pipeline

```bash
# Initialize project
sloth-runner workflow init my-app --template cicd

# Run tests
sloth-runner run -f my-app/.sloth/test.sloth -o rich

# Deploy to staging
sloth-runner run staging-stack --yes

# Check agent status
sloth-runner agent list --master ci-master:50053

# Schedule nightly builds
sloth-runner scheduler add nightly-build \
  --cron "0 0 * * *" \
  --workflow build.sloth
```

### Infrastructure Automation

```bash
# Create infrastructure stack
sloth-runner stack new prod-infra -f infrastructure.sloth

# Apply with modern output
sloth-runner run prod-infra -o modern

# Show stack state
sloth-runner stack show prod-infra

# Teardown
sloth-runner stack delete prod-infra --force
```

### Distributed Task Execution

```bash
# Start master server
sloth-runner master --port 50053 --daemon

# Start agents on different servers
sloth-runner agent start --master master:50053 --name web-1 --tags web,nginx
sloth-runner agent start --master master:50053 --name db-1 --tags database,postgres

# Execute on specific agent
sloth-runner agent exec web-1 "systemctl status nginx" --master master:50053

# Start UI for monitoring
sloth-runner ui --port 8080
```

---

## Best Practices

### 1. Use Stacks for State Management

```bash
# Don't: Run without state
sloth-runner run -f deploy.sloth

# Do: Use stacks for persistence
sloth-runner stack new prod
sloth-runner run prod
```

### 2. Specify Output Format for CI/CD

```bash
# JSON for parsing
sloth-runner run -f ci.sloth -o json > results.json

# Rich for interactive
sloth-runner run -f deploy.sloth -o rich
```

### 3. Use Values Files for Environments

```bash
# Development
sloth-runner run -f app.sloth -v dev-values.yaml

# Production
sloth-runner run -f app.sloth -v prod-values.yaml
```

### 4. Tag Agents Appropriately

```bash
# Specific capabilities
sloth-runner agent start --tags "linux,docker,aws,x86_64"

# Environment-based
sloth-runner agent start --tags "prod,us-east-1"
```

---

## Troubleshooting

### Connection Issues

```bash
# Test master connectivity
curl http://master:50053/health

# Check agent logs
sloth-runner agent start --debug
```

### Task Execution Failures

```bash
# Run with debug output
sloth-runner run -f task.sloth --debug

# Interactive mode for troubleshooting
sloth-runner run -f task.sloth --interactive
```

### State Issues

```bash
# View stack state
sloth-runner stack show my-stack

# Reset stack (careful!)
sloth-runner stack delete my-stack
sloth-runner stack new my-stack -f workflow.sloth
```

---

## Related Documentation

- [Getting Started](/en/getting-started/)
- [Core Concepts](/en/core-concepts/)
- [Agent Architecture](/en/master-agent-architecture/)
- [Scheduler Guide](/advanced-scheduler/)
- [Web Dashboard](/web-dashboard/)
- [Stack Management](/stack-management/)

---

## See Also

- [REPL Interactive Shell](/en/repl/)
- [Modern DSL Syntax](/modern-dsl/introduction/)
- [Module Reference](/modules/)
- [Examples Repository](https://github.com/chalkan3-sloth/sloth-runner/tree/main/examples)

---

**Need more help?** Run `sloth-runner [command] --help` for detailed information about any command.

---

### `sloth-runner version`

Displays the current version of `sloth-runner`.

```bash
sloth-runner version
```

### `sloth-runner scheduler`

Manages the `sloth-runner` task scheduler, allowing you to enable, disable, list, and delete scheduled tasks.

For detailed information on scheduler commands and configuration, refer to the [Task Scheduler documentation](scheduler.md).

**Subcommands:**

*   `sloth-runner scheduler enable`: Starts the scheduler as a background process.
*   `sloth-runner scheduler disable`: Stops the running scheduler process.
*   `sloth-runner scheduler list`: Lists all configured scheduled tasks.
*   `sloth-runner scheduler delete <task_name>`: Deletes a specific scheduled task.

