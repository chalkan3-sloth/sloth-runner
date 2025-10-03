# State Management and Idempotency

## Overview

Sloth Runner includes a built-in state management system that enables tracking of configuration state and resource management. This provides a foundation for implementing idempotent operations.

## State Storage

State is stored in SQLite databases per agent:

```
$HOME/.sloth-runner/state/
  ├── local.db          # Local agent state
  ├── mariaguica.db     # mariaguica agent state
  └── production.db     # production agent state
```

## CLI Commands

### List States

View all tracked states:

```bash
# List all states for local agent
sloth-runner state list

# List states for specific agent
sloth-runner state list --agent mariaguica

# Filter by prefix
sloth-runner state list file

# Output as JSON
sloth-runner state list --output json
```

### Show State Details

View detailed information about a specific state:

```bash
# Show state details
sloth-runner state show deployment:version

# JSON output
sloth-runner state show deployment:version --output json
```

### Delete States

Remove state entries:

```bash
# Delete specific state (with confirmation)
sloth-runner state delete deployment:version

# Skip confirmation
sloth-runner state delete deployment:version --yes
```

### Clear All States

Remove all state entries:

```bash
# Clear all states (with confirmation)
sloth-runner state clear

# Skip confirmation
sloth-runner state clear --yes
```

### State Statistics

View state database statistics:

```bash
# Show statistics
sloth-runner state stats

# JSON output
sloth-runner state stats --output json
```

## Using State in Tasks

### Basic State Operations

```lua
task({
    name = "track-deployment",
    run = function()
        -- Store deployment version
        state.set("deployment:version", "v1.2.3")
        
        -- Retrieve stored value
        local version = state.get("deployment:version")
        print("Current version: " .. version)
        
        -- Check if key exists
        if state.exists("deployment:rollback") then
            print("Rollback state available")
        end
        
        -- Delete state
        state.delete("deployment:old-version")
    end
})
```

### State Locks

Prevent concurrent modifications:

```lua
task({
    name = "critical-update",
    run = function()
        -- Acquire lock with 5 minute timeout
        state.lock("deployment", "task-123", 300)
        
        -- Critical operations here
        state.set("deployment:status", "in-progress")
        
        -- Release lock
        state.unlock("deployment", "task-123")
    end
})
```

### Using With Lock Helper

```lua
task({
    name = "safe-update",
    run = function()
        -- Automatically manages lock lifecycle
        state.with_lock("deployment", "task-123", 300, function()
            -- Operations inside lock
            state.set("deployment:status", "updating")
            -- ... perform update ...
            state.set("deployment:status", "complete")
        end)
    end
})
```

### Implementing Idempotency

You can implement idempotent operations using state tracking:

```lua
task({
    name = "install-package",
    run = function()
        local package_name = "nginx"
        local state_key = "package:installed:" .. package_name
        
        -- Check if already installed
        if state.exists(state_key) then
            print("Package " .. package_name .. " already installed (skipping)")
            return
        end
        
        -- Install package
        cmd({
            command = "apt-get install -y " .. package_name,
            delegate_to = values.host
        })
        
        -- Track installation
        state.set(state_key, tostring(os.time()))
        print("Package " .. package_name .. " installed")
    end
})
```

### Configuration File Management

Track configuration changes:

```lua
task({
    name = "update-config",
    run = function()
        local config_file = "/etc/app/config.ini"
        local template_src = "./templates/config.ini.tmpl"
        
        -- Read current template
        local new_content = template.render({
            src = template_src,
            vars = values
        })
        
        -- Calculate checksum
        local new_hash = crypto.sha256(new_content)
        local state_key = "config:hash:" .. config_file
        
        -- Get stored hash
        local old_hash = state.get(state_key)
        
        if old_hash == new_hash then
            print("Configuration unchanged (skipping)")
            return
        end
        
        -- Deploy new configuration
        file.copy({
            src = template_src,
            dest = config_file,
            delegate_to = values.host
        })
        
        -- Update hash
        state.set(state_key, new_hash)
        print("Configuration updated")
    end
})
```

### Multi-Resource Tracking

Track multiple related resources:

```lua
task({
    name = "setup-webserver",
    run = function()
        local resources = {
            {type = "user", name = "www-data"},
            {type = "dir", name = "/var/www"},
            {type = "service", name = "nginx"}
        }
        
        for _, res in ipairs(resources) do
            local state_key = res.type .. ":" .. res.name
            
            if not state.exists(state_key) then
                -- Create resource based on type
                if res.type == "user" then
                    cmd({command = "useradd " .. res.name})
                elseif res.type == "dir" then
                    cmd({command = "mkdir -p " .. res.name})
                elseif res.type == "service" then
                    cmd({command = "systemctl enable " .. res.name})
                end
                
                -- Mark as created
                state.set(state_key, "created")
            end
        end
    end
})
```

### Cross-Agent State Sharing

Share state between agents:

```lua
task({
    name = "leader-election",
    delegate_to = "agent1",
    run = function()
        -- Try to become leader
        if not state.exists("cluster:leader") then
            state.set("cluster:leader", "agent1")
            print("Became cluster leader")
        end
    end
})

task({
    name = "check-leader",
    delegate_to = "agent2",
    run = function()
        local leader = state.get("cluster:leader")
        print("Current leader: " .. (leader or "none"))
    end
})
```

## Advanced Patterns

### State-Based Conditionals

```lua
task({
    name = "conditional-deployment",
    run = function()
        local env = state.get("environment:type") or "development"
        
        if env == "production" then
            -- Production-specific logic
            state.set("deployment:replicas", "5")
        else
            -- Development logic
            state.set("deployment:replicas", "1")
        end
    end
})
```

### Versioned State

```lua
task({
    name = "versioned-config",
    run = function()
        -- Increment version
        local version = state.increment("config:version", 1)
        
        -- Store versioned config
        state.set("config:v" .. version, config_content)
        
        -- Keep reference to current
        state.set("config:current", tostring(version))
    end
})
```

### State Cleanup

```lua
task({
    name = "cleanup-old-state",
    run = function()
        -- List all states with prefix
        local all_states = state.list("temporary:")
        
        -- Clean up temporary states
        for key, _ in pairs(all_states) do
            state.delete(key)
        end
    end
})
```

## Best Practices

### 1. Use Consistent Key Naming

```lua
-- Good pattern: {resource_type}:{operation}:{identifier}
state.set("package:installed:nginx", "true")
state.set("config:hash:/etc/nginx/nginx.conf", checksum)
state.set("deployment:version:app", "v1.2.3")
```

### 2. Check Before Modify

```lua
-- Always check existence before operations
if not state.exists("resource:initialized") then
    -- Initialize resource
    state.set("resource:initialized", "true")
end
```

### 3. Use Locks for Critical Sections

```lua
-- Protect critical operations
state.with_lock("resource", "task-id", 300, function()
    -- Critical code here
end)
```

### 4. Clean Up State

```lua
-- Remove state when resource is deleted
state.delete("package:installed:old-package")
```

### 5. Monitor State Size

```bash
# Regular checks
sloth-runner state stats
```

## State API Reference

Available in Lua tasks:

### state.set(key, value)
Store a key-value pair
```lua
state.set("app:version", "1.0.0")
```

### state.get(key)
Retrieve a value
```lua
local version = state.get("app:version")
```

### state.exists(key)
Check if key exists
```lua
if state.exists("deployment:lock") then
    print("Deployment locked")
end
```

### state.delete(key)
Remove a key
```lua
state.delete("temp:session")
```

### state.list(prefix)
List keys with prefix
```lua
local configs = state.list("config:")
```

### state.increment(key, delta)
Increment numeric value
```lua
local count = state.increment("deploy:count", 1)
```

### state.lock(name, holder, timeout_seconds)
Acquire a lock
```lua
state.lock("deployment", "task-123", 300)
```

### state.unlock(name, holder)
Release a lock
```lua
state.unlock("deployment", "task-123")
```

### state.with_lock(name, holder, timeout, function)
Execute function with lock
```lua
state.with_lock("deploy", "task-1", 300, function()
    -- Protected code
end)
```

### state.is_locked(name)
Check if locked
```lua
local locked, holder = state.is_locked("deployment")
if locked then
    print("Locked by: " .. holder)
end
```

## Troubleshooting

### State Out of Sync

Reset specific state:
```bash
sloth-runner state delete package:installed:nginx
```

### Clear All State

Start fresh:
```bash
sloth-runner state clear --yes
```

### View State Contents

Inspect stored values:
```bash
sloth-runner state list
sloth-runner state show app:version
```

## Future Enhancements

Planned improvements:
- Remote state backends (S3, etcd, Consul)
- State encryption at rest
- State versioning and history
- State import/export
- Web UI for state visualization
- Automatic resource checksumming
- Built-in idempotency for all modules

