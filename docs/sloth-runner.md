# SLOTH-RUNNER(1) - Automation and Orchestration Tool

## NAME

**sloth-runner** - Distributed task automation and orchestration system with Lua scripting

## SYNOPSIS

```
sloth-runner [global-flags] <command> [command-flags] [arguments]
```

## DESCRIPTION

Sloth Runner is a powerful distributed automation framework that enables you to:

- **Define and execute workflows** using Lua-based task definitions
- **Orchestrate tasks** across multiple remote agents
- **React to system events** using a sophisticated hooks system
- **Manage infrastructure** with built-in modules for SSH, Incus/LXD, package management
- **Track execution history** with comprehensive state management
- **Monitor and observe** tasks and agents in real-time

The tool is designed for infrastructure automation, deployment orchestration, and distributed task execution with a focus on simplicity and extensibility.

## KEY CONCEPTS

### Workflows (.sloth files)
Lua scripts that define tasks, their dependencies, and execution logic. Tasks can run locally or be delegated to remote agents.

### Agents
Remote workers that execute delegated tasks. Agents communicate with the master via gRPC and can be managed through the CLI.

### Hooks
Lua scripts that execute in response to system events (task lifecycle, agent status, workflow events). Hooks enable reactive automation and monitoring.

### Events
System-generated notifications about task execution, agent status, and workflow state. Events trigger registered hooks.

### Stacks
Isolated environments for managing workflows, state, and configuration. Stacks enable multi-environment deployments.

## GLOBAL FLAGS

```
-V, --version          Show version information
-h, --help            Show help for any command
```

## MAIN COMMANDS

### Core Workflow Execution

- **run** - Execute a workflow or task directly
- **workflow** - Manage workflows (list, preview, run)
- **sloth** - Manage workflow files repository

### Agent Management

- **agent** - Manage remote agents (install, list, update, metrics, etc.)
- **master** - Start the master server for agent coordination

### Event and Hook System

- **hook** - Manage event hooks (add, list, enable, disable)
- **events** - View and manage event queue

### State and Configuration

- **stack** - Manage isolated execution environments
- **state** - Query and manage execution state
- **secrets** - Manage sensitive configuration values
- **scheduler** - Manage scheduled task execution

### Utilities

- **list** - List available tasks in workflows
- **ui** - Start web UI for monitoring and management
- **completion** - Generate shell completion scripts

## EXAMPLES

### Basic Workflow Execution

Execute a task from a workflow file:

```bash
sloth-runner run deploy_app --file deployments/app.sloth --yes
```

Execute with parameters:

```bash
sloth-runner run deploy_app \
  --file deployments/app.sloth \
  --param environment=production \
  --param version=v2.1.0 \
  --yes
```

### Delegated Execution

Run a task on a remote agent:

```bash
sloth-runner run install_packages \
  --file maintenance/packages.sloth \
  --delegate-to web-server-01 \
  --yes
```

Run across multiple agents:

```bash
sloth-runner run health_check \
  --file monitoring/checks.sloth \
  --delegate-to web-server-01 \
  --delegate-to web-server-02 \
  --delegate-to web-server-03 \
  --yes
```

### Hook Management

Register a hook for task failures:

```bash
sloth-runner hook add failure_alert \
  --event task.failed \
  --script hooks/send-alert.lua \
  --description "Send alert on task failure"
```

Test a hook with mock data:

```bash
sloth-runner hook test failure_alert \
  --event-type task.failed \
  --mock-data '{"task_name":"deploy_app","error":"connection timeout"}'
```

### Agent Management

Install a new agent on a remote host:

```bash
sloth-runner agent install prod-server-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user root \
  --ssh-port 22 \
  --master 192.168.1.10:50053 \
  --port 50060
```

Check agent status and modules:

```bash
sloth-runner agent list
sloth-runner agent modules prod-server-01 --check pkg
```

### Event Monitoring

View recent events:

```bash
sloth-runner events list --limit 20
```

View detailed event with hook execution history:

```bash
sloth-runner events show <event-id>
```

Clean up old events:

```bash
sloth-runner events cleanup --older-than 7d
```

### Stack Management

Create and activate a stack for staging environment:

```bash
sloth-runner stack create staging --description "Staging environment"
sloth-runner stack activate staging
```

List all stacks:

```bash
sloth-runner stack list
```

## WORKFLOW FILE FORMAT

Sloth Runner workflows are Lua scripts with a specific structure:

```lua
-- examples/simple-deployment.sloth

-- Define a task
local build_task = task("build_application")
    :description("Build the application")
    :command(function(this, params)
        log.info("Building application version: " .. params.version)

        local result = exec.run("go build -o app ./cmd/app")
        if result.exit_code ~= 0 then
            error("Build failed: " .. result.stderr)
        end

        log.info("Build completed successfully")
        return true, "Build completed"
    end)
    :build()

-- Define another task with dependency
local deploy_task = task("deploy_application")
    :description("Deploy the built application")
    :depends_on("build_application")
    :command(function(this, params)
        log.info("Deploying to " .. params.environment)

        -- Copy binary to remote server
        local result = ssh.exec({
            host = params.target_host,
            user = "deploy",
            command = "systemctl restart app"
        })

        if result.exit_code ~= 0 then
            error("Deployment failed")
        end

        return true, "Deployed successfully"
    end)
    :build()

-- Return tasks for execution
return {build_task, deploy_task}
```

## HOOK SCRIPT FORMAT

Hooks are Lua scripts that respond to events:

```lua
-- hooks/task-failure-alert.lua

function on_event()
    -- Access event data
    local task = event.task or event.data.task

    -- Log the failure
    local alert = string.format(
        "ALERT: Task '%s' failed on agent '%s'\nError: %s",
        task.task_name,
        task.agent_name,
        task.error
    )

    -- Use globally available modules
    file_ops.lineinfile({
        path = "/var/log/sloth-runner/failures.log",
        line = alert,
        create = true
    })

    -- Could also send notification via HTTP
    -- http.post({
    --     url = "https://alerts.example.com/webhook",
    --     body = alert
    -- })

    print("Alert sent for failed task: " .. task.task_name)
    return true
end
```

## ENVIRONMENT VARIABLES

```
SLOTH_RUNNER_MASTER_ADDR    Master server address (default: localhost:50053)
SLOTH_RUNNER_CONFIG_DIR     Configuration directory (default: ~/.sloth-runner)
SLOTH_RUNNER_CACHE_DIR      Cache directory (default: ./.sloth-cache)
```

## FILES

```
~/.sloth-runner/config.yaml         Main configuration file
./.sloth-cache/agents.db           Agent registry database
./.sloth-cache/hooks.db            Hooks and events database
./.sloth-cache/state.db            Task execution state database
```

## EXIT CODES

```
0    Success
1    General error
2    Invalid arguments or flags
3    Task execution failed
4    Agent communication error
5    Hook execution failed
```

## SEE ALSO

- **sloth-runner-hook(1)** - Manage event hooks
- **sloth-runner-agent(1)** - Manage remote agents
- **sloth-runner-events(1)** - View and manage events
- **sloth-runner-workflow(1)** - Manage workflows
- **sloth-runner-run(1)** - Execute workflows and tasks

Full documentation: https://github.com/yourusername/task-runner

## AUTHOR

Written by the Sloth Runner development team.

## COPYRIGHT

Copyright © 2025. Released under MIT License.
