# SLOTH-RUNNER-HOOK(1) - Event Hook Management

## NAME

**sloth-runner hook** - Manage event hooks for reactive automation

## SYNOPSIS

```
sloth-runner hook <command> [options]
```

## DESCRIPTION

The **hook** command provides a complete interface for managing event hooks in Sloth Runner. Hooks are Lua scripts that execute automatically in response to system events such as task completions, agent status changes, workflow events, and custom events.

Hooks enable you to:

- React to task failures with alerts and remediation
- Monitor agent health and connectivity
- Track workflow execution statistics
- Implement custom business logic triggered by events
- Automate incident response and logging
- Integrate with external systems (monitoring, ticketing, notifications)

The hook system runs 100 concurrent workers processing events from a buffered queue (1000 event capacity), ensuring high throughput and minimal latency.

## AVAILABLE COMMANDS

- **add** - Register a new event hook
- **list** - List all registered hooks
- **get** - Get hook details in JSON format
- **show** - Show detailed hook information with execution history
- **delete** - Remove a hook from the system
- **enable** - Enable a disabled hook
- **disable** - Temporarily disable a hook
- **test** - Test a hook with mock event data

## HOOK ADD

Register a new event hook that will execute when specific events occur.

### Synopsis

```
sloth-runner hook add <hook-name> --file <path> --event <type> [options]
```

### Options

```
-f, --file <path>          Path to the Lua hook script (required)
-e, --event <type>         Event type to trigger the hook (required)
-d, --description <text>   Human-readable description
-s, --stack <name>         Stack name for hook isolation
    --enabled              Enable immediately (default: true)
```

### Event Types

The system supports 100+ event types across these categories:

**Task Events:**
- `task.started` - Task execution began
- `task.completed` - Task finished successfully
- `task.failed` - Task execution failed
- `task.timeout` - Task exceeded time limit
- `task.retrying` - Task retry attempt
- `task.cancelled` - Task was cancelled

**Agent Events:**
- `agent.registered` - New agent joined
- `agent.connected` - Agent established connection
- `agent.disconnected` - Agent lost connection
- `agent.heartbeat_failed` - Agent missed heartbeat
- `agent.updated` - Agent software updated
- `agent.version_mismatch` - Agent version incompatible
- `agent.resource_high` - Agent resources critically high

**Workflow Events:**
- `workflow.started` - Workflow execution began
- `workflow.completed` - Workflow finished
- `workflow.failed` - Workflow execution failed
- `workflow.paused` - Workflow paused
- `workflow.resumed` - Workflow resumed
- `workflow.cancelled` - Workflow cancelled

**System Events:**
- `system.startup` - System initialized
- `system.shutdown` - System shutting down
- `system.error` - System-level error occurred
- `system.warning` - System warning issued

**Custom Events:**
- `custom` - User-defined events dispatched from workflows

### Hook Script Format

Hook scripts must define an `on_event()` function:

```lua
function on_event()
    -- Access event data through global 'event' variable
    local task = event.task or event.data.task

    -- Available event fields depend on event type:
    -- task.task_name, task.agent_name, task.status,
    -- task.exit_code, task.error, task.duration

    -- Use globally available modules (no require needed):
    -- - file_ops: File operations
    -- - exec: Command execution
    -- - log: Logging
    -- - ssh: SSH operations
    -- - http: HTTP requests

    -- Example: Log to file
    file_ops.lineinfile({
        path = "/var/log/sloth/events.log",
        line = "Task " .. task.task_name .. " completed",
        create = true
    })

    -- Return true for success, false for failure
    return true
end
```

### Examples

Register a hook for task failures:

```bash
sloth-runner hook add task_failure_alert \
  --file hooks/alert-on-failure.lua \
  --event task.failed \
  --description "Send alerts when tasks fail"
```

Register a multi-event hook for task lifecycle tracking:

```bash
sloth-runner hook add task_tracker \
  --file hooks/task-lifecycle.lua \
  --event task.started \
  --description "Track task execution lifecycle"

sloth-runner hook add task_tracker_completion \
  --file hooks/task-lifecycle.lua \
  --event task.completed \
  --description "Track task completions"
```

Register a hook isolated to a specific stack:

```bash
sloth-runner hook add prod_monitor \
  --file hooks/production-monitor.lua \
  --event task.failed \
  --stack production \
  --description "Production environment monitoring"
```

Register a hook for agent monitoring:

```bash
sloth-runner hook add agent_health \
  --file hooks/agent-health-check.lua \
  --event agent.heartbeat_failed \
  --description "Alert on agent connectivity issues"
```

Register initially disabled (for testing):

```bash
sloth-runner hook add experimental_hook \
  --file hooks/experimental.lua \
  --event custom \
  --enabled=false \
  --description "Experimental custom event handler"
```

## HOOK LIST

Display all registered hooks in the system.

### Synopsis

```
sloth-runner hook list [options]
```

### Options

```
-s, --stack <name>    Filter hooks by stack name
```

### Output Format

The list displays:
- Hook name
- Event type it responds to
- Enabled/disabled status
- Associated stack (if any)
- Description

### Examples

List all hooks:

```bash
sloth-runner hook list
```

List hooks for a specific stack:

```bash
sloth-runner hook list --stack production
```

Sample output:

```
Hooks:
┌─────────────────────┬────────────────┬─────────┬────────────┬──────────────────────────────┐
│ Name                │ Event Type     │ Enabled │ Stack      │ Description                  │
├─────────────────────┼────────────────┼─────────┼────────────┼──────────────────────────────┤
│ task_failure_alert  │ task.failed    │ Yes     │ default    │ Send alerts when tasks fail  │
│ task_tracker        │ task.started   │ Yes     │ default    │ Track task lifecycle         │
│ agent_health        │ agent.hb_fail  │ Yes     │ production │ Monitor agent connectivity   │
│ experimental_hook   │ custom         │ No      │ default    │ Experimental handler         │
└─────────────────────┴────────────────┴─────────┴────────────┴──────────────────────────────┘
```

## HOOK GET

Retrieve hook details in JSON format for programmatic access.

### Synopsis

```
sloth-runner hook get <hook-name>
```

### Examples

Get hook as JSON:

```bash
sloth-runner hook get task_failure_alert
```

Sample output:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "task_failure_alert",
  "event_type": "task.failed",
  "script_path": "hooks/alert-on-failure.lua",
  "description": "Send alerts when tasks fail",
  "enabled": true,
  "stack": "default",
  "created_at": "2025-10-06T14:23:45Z",
  "updated_at": "2025-10-06T14:23:45Z"
}
```

Use in scripts:

```bash
# Check if hook is enabled
enabled=$(sloth-runner hook get task_failure_alert | jq -r '.enabled')
if [ "$enabled" = "true" ]; then
    echo "Hook is active"
fi
```

## HOOK SHOW

Display detailed human-readable information about a hook, including execution history.

### Synopsis

```
sloth-runner hook show <hook-name> [options]
```

### Options

```
--history           Include execution history
--limit <n>         Number of history entries to show (default: 10)
```

### Examples

Show basic hook information:

```bash
sloth-runner hook show task_failure_alert
```

Show with execution history:

```bash
sloth-runner hook show task_failure_alert --history
```

Show with more history entries:

```bash
sloth-runner hook show task_failure_alert --history --limit 50
```

Sample output:

```
Hook: task_failure_alert
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

ID:          550e8400-e29b-41d4-a716-446655440000
Event Type:  task.failed
Status:      Enabled
Stack:       default
Script Path: hooks/alert-on-failure.lua
Description: Send alerts when tasks fail

Created:     2025-10-06 14:23:45
Updated:     2025-10-06 14:23:45

Execution Statistics:
  Total Executions: 42
  Successful:       40
  Failed:           2
  Success Rate:     95.2%

Recent Execution History (Last 10):
┌──────────────────────┬─────────┬──────────┬─────────────────────────┐
│ Timestamp            │ Status  │ Duration │ Event ID                │
├──────────────────────┼─────────┼──────────┼─────────────────────────┤
│ 2025-10-06 17:32:22  │ Success │ 12ms     │ evt_abc123              │
│ 2025-10-06 16:15:33  │ Success │ 8ms      │ evt_def456              │
│ 2025-10-06 15:42:11  │ Failed  │ 2ms      │ evt_ghi789              │
└──────────────────────┴─────────┴──────────┴─────────────────────────┘
```

## HOOK DELETE

Remove a hook from the system permanently.

### Synopsis

```
sloth-runner hook delete <hook-name>
```

### Examples

Delete a hook:

```bash
sloth-runner hook delete experimental_hook
```

Delete with confirmation:

```bash
sloth-runner hook delete task_failure_alert
# System will prompt: "Are you sure you want to delete hook 'task_failure_alert'? (y/N)"
```

## HOOK ENABLE

Enable a previously disabled hook to resume event processing.

### Synopsis

```
sloth-runner hook enable <hook-name>
```

### Examples

Enable a hook:

```bash
sloth-runner hook enable experimental_hook
```

## HOOK DISABLE

Temporarily disable a hook without deleting it. Useful for maintenance or debugging.

### Synopsis

```
sloth-runner hook disable <hook-name>
```

### Examples

Disable a hook during maintenance:

```bash
sloth-runner hook disable task_failure_alert
# Perform maintenance...
sloth-runner hook enable task_failure_alert
```

## HOOK TEST

Test a hook by executing it with mock event data. This allows you to verify hook logic without waiting for real events.

### Synopsis

```
sloth-runner hook test <hook-name> [options]
```

### Options

```
-d, --data <json>    Custom event data as JSON string
```

### Examples

Test with default mock data:

```bash
sloth-runner hook test task_failure_alert
```

Test with custom event data:

```bash
sloth-runner hook test task_failure_alert \
  --data '{
    "task": {
      "task_name": "deploy_app",
      "agent_name": "prod-server-01",
      "status": "failed",
      "exit_code": 1,
      "error": "Connection timeout",
      "duration": "5m30s"
    }
  }'
```

Test agent event hook:

```bash
sloth-runner hook test agent_health \
  --data '{
    "agent": {
      "name": "web-server-03",
      "last_heartbeat": "2025-10-06T17:25:00Z",
      "missed_heartbeats": 3
    }
  }'
```

Sample output:

```
Testing hook: task_failure_alert
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Hook executed successfully

Output:
  🚨 ALERT: Task deploy_app failed
  Alert written to /var/log/sloth/failures.log

Execution Time: 15ms
Status: Success
```

## COMPLETE WORKFLOW EXAMPLE

Here's a complete example of setting up monitoring hooks:

### 1. Create Hook Scripts

**hooks/task-failure-alert.lua:**
```lua
function on_event()
    local task = event.task

    local message = string.format(
        "ALERT: Task '%s' failed on %s\\nError: %s\\nDuration: %s",
        task.task_name,
        task.agent_name,
        task.error or "Unknown error",
        task.duration or "Unknown"
    )

    -- Log to file
    file_ops.lineinfile({
        path = "/var/log/sloth/failures.log",
        line = string.format("[%s] %s", os.date("%Y-%m-%d %H:%M:%S"), message),
        create = true
    })

    -- Send to monitoring system
    http.post({
        url = "https://monitoring.example.com/webhook",
        headers = {["Content-Type"] = "application/json"},
        body = json.encode({
            alert_type = "task_failure",
            task = task.task_name,
            agent = task.agent_name,
            error = task.error
        })
    })

    print("Alert sent for task: " .. task.task_name)
    return true
end
```

**hooks/task-success-tracker.lua:**
```lua
function on_event()
    local task = event.task

    file_ops.lineinfile({
        path = "/var/log/sloth/completions.log",
        line = string.format(
            "[%s] ✓ %s completed on %s (exit: %d, duration: %s)",
            os.date("%Y-%m-%d %H:%M:%S"),
            task.task_name,
            task.agent_name,
            task.exit_code,
            task.duration
        ),
        create = true
    })

    return true
end
```

**hooks/agent-monitor.lua:**
```lua
function on_event()
    local agent = event.agent

    log.warn(string.format(
        "Agent %s disconnected! Last seen: %s",
        agent.name,
        agent.last_heartbeat
    ))

    -- Attempt automatic remediation
    exec.run(string.format(
        "sloth-runner agent start %s --ssh-reconnect",
        agent.name
    ))

    return true
end
```

### 2. Register Hooks

```bash
# Task failure monitoring
sloth-runner hook add task_failure_alert \
  --file hooks/task-failure-alert.lua \
  --event task.failed \
  --description "Alert on task failures"

# Task success tracking
sloth-runner hook add task_success_tracker \
  --file hooks/task-success-tracker.lua \
  --event task.completed \
  --description "Track successful task completions"

# Agent health monitoring
sloth-runner hook add agent_monitor \
  --file hooks/agent-monitor.lua \
  --event agent.disconnected \
  --description "Monitor and remediate agent disconnections"
```

### 3. Test Hooks

```bash
# Test failure alert
sloth-runner hook test task_failure_alert \
  --data '{"task":{"task_name":"test","agent_name":"local","error":"test error"}}'

# Test success tracker
sloth-runner hook test task_success_tracker \
  --data '{"task":{"task_name":"test","agent_name":"local","exit_code":0}}'
```

### 4. Monitor Hook Execution

```bash
# View hook execution history
sloth-runner hook show task_failure_alert --history --limit 20

# List all active hooks
sloth-runner hook list

# Check recent events
sloth-runner events list --limit 10
```

## HOOK SCRIPT BEST PRACTICES

1. **Always return a boolean** - Return `true` for success, `false` for failure
2. **Handle errors gracefully** - Use `pcall()` for operations that might fail
3. **Keep hooks fast** - Hooks run synchronously; avoid long-running operations
4. **Use appropriate logging** - Use `log.info()`, `log.warn()`, `log.error()`
5. **Test before deploying** - Use `hook test` to verify logic
6. **Document your hooks** - Add comments explaining the hook's purpose

## AVAILABLE MODULES IN HOOKS

Hooks have access to these modules globally (no `require()` needed):

- **file_ops** - File operations (lineinfile, blockinfile, copy, template)
- **exec** - Command execution
- **log** - Logging functions
- **ssh** - SSH operations
- **http** - HTTP client
- **json** - JSON encoding/decoding
- **os** - OS utilities (date, time, etc.)
- **string** - String manipulation

## FILES

```
.sloth-cache/hooks.db    SQLite database storing hooks and events
hooks/                   Recommended directory for hook scripts
```

## SEE ALSO

- **sloth-runner-events(1)** - View and manage events
- **sloth-runner(1)** - Main sloth-runner command
- **sloth-runner-workflow(1)** - Workflow management

## AUTHOR

Written by the Sloth Runner development team.

## COPYRIGHT

Copyright © 2025. Released under MIT License.
