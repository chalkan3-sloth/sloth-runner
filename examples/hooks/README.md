# Event Hooks Examples

This directory contains example hook scripts that can be triggered by various events in sloth-runner.

## Available Hooks

### 1. **notify-agent-join.lua**
- **Event**: `agent.registered`
- **Description**: Sends notifications when a new agent joins the cluster
- **Use case**: Automated agent onboarding, monitoring

### 2. **alert-agent-down.lua**
- **Event**: `agent.disconnected`
- **Description**: Sends alerts when an agent disconnects
- **Use case**: Infrastructure monitoring, alerting

## Hook DSL Reference

### Event Data Structure

Hooks receive an `event` global variable with the following structure:

```lua
event = {
    type = "agent.registered",  -- Event type
    timestamp = 1234567890,      -- Unix timestamp
    data = {...},                -- Event-specific data

    -- Convenience fields (extracted from data)
    agent = {
        name = "agent-name",
        address = "192.168.1.100:50051",
        tags = {"production", "linux"},
        version = "1.0.0",
        system_info = {
            os = "Linux",
            cpus = 4,
            memory = 8192
        }
    }
}
```

### Available Functions

#### Logging Functions
```lua
log.info("message")    -- Info level log
log.warn("message")    -- Warning level log
log.error("message")   -- Error level log
log.debug("message")   -- Debug level log
```

#### HTTP Functions
```lua
http.post(url)  -- Send HTTP POST request
```

#### Utility Functions
```lua
contains(table, value)  -- Check if table contains value
```

### Hook Function

Your hook must define an `on_event()` function:

```lua
function on_event()
    -- Your hook logic here

    -- Return true for success, false for failure
    return true
end
```

## Event Types

### Agent Events
- `agent.registered` - When a new agent joins
- `agent.disconnected` - When an agent disconnects
- `agent.heartbeat_failed` - When agent heartbeat fails
- `agent.updated` - When agent is updated

### Task Events
- `task.started` - When a task starts
- `task.completed` - When a task completes
- `task.failed` - When a task fails

### Workflow Events
- `workflow.started` - When a workflow starts
- `workflow.completed` - When a workflow completes
- `workflow.failed` - When a workflow fails

## Using Hooks

### 1. Add a Hook

```bash
sloth-runner hook add notify-agent-join \
  --file examples/hooks/notify-agent-join.lua \
  --event agent.registered \
  --description "Notify when agent joins"
```

### 2. List Hooks

```bash
sloth-runner hook list
```

### 3. Show Hook Details

```bash
sloth-runner hook show notify-agent-join
```

### 4. Test a Hook

```bash
sloth-runner hook test notify-agent-join
```

### 5. Disable/Enable Hooks

```bash
sloth-runner hook disable notify-agent-join
sloth-runner hook enable notify-agent-join
```

### 6. Delete a Hook

```bash
sloth-runner hook delete notify-agent-join
```

## Advanced Example

Here's a more complete example with conditional logic and multiple actions:

```lua
function on_event()
    local agent = event.agent

    -- Log the event
    log.info("Agent registered: " .. agent.name)

    -- Check if production agent
    if agent.tags and contains(agent.tags, "production") then
        log.warn("Production agent detected!")

        -- Send high-priority notification
        http.post("https://alerts.example.com/agent-join")

        -- Additional checks for production
        if agent.system_info then
            if agent.system_info.cpus < 4 then
                log.warn("Warning: Production agent has less than 4 CPUs")
            end

            if agent.system_info.memory < 4096 then
                log.warn("Warning: Production agent has less than 4GB RAM")
            end
        end
    else
        -- Regular notification for non-production
        log.info("Non-production agent joined")
        http.post("https://notifications.example.com/agent-join")
    end

    return true
end
```

## Best Practices

1. **Keep hooks simple**: Complex logic should be in separate tools/scripts
2. **Handle errors**: Always check if data exists before accessing it
3. **Use logging**: Log important steps for debugging
4. **Test hooks**: Use `sloth-runner hook test` before enabling
5. **Monitor execution**: Use `sloth-runner hook show --history` to check execution history

## Troubleshooting

If a hook fails:

1. Check the hook file exists and is readable
2. Test the hook: `sloth-runner hook test <hook-name>`
3. Check execution history: `sloth-runner hook show <hook-name> --history`
4. Verify event type matches the hook configuration
5. Check logs for Lua errors

## Next Steps

- Create custom hooks for your workflows
- Integrate with external monitoring systems
- Set up automated responses to events
- Build complex event-driven automation
