# Stack Management & Idempotency Guide

## Overview

Sloth Runner implements a sophisticated **stack-based state management system** similar to Pulumi and Terraform. This ensures **idempotent** infrastructure automation where resources are only created or modified when necessary.

## Key Concepts

### 1. Stacks

A **stack** is an isolated execution environment for your workflows. Each stack:
- Maintains its own state database
- Tracks all managed resources
- Provides idempotency guarantees
- Records execution history

### 2. Resources

A **resource** represents any managed entity (file, package, service, cloud resource, etc.). Each resource:
- Has a unique identifier within the stack
- Tracks its current state and properties
- Maintains a checksum for drift detection
- Is only applied when changes are detected

### 3. Idempotency

**Idempotency** means running the same workflow multiple times produces the same result. Resources are:
- **Created** if they don't exist
- **Updated** if they changed
- **Skipped** if they're already in the desired state

## CLI Commands

### Stack Management

```bash
# Create a new stack
sloth-runner stack new my-infrastructure

# List all stacks
sloth-runner stack list

# Show stack details
sloth-runner stack show my-infrastructure

# Delete a stack
sloth-runner stack delete my-infrastructure
```

### State Management

```bash
# Set a key-value pair
sloth-runner state set key value

# Get a value
sloth-runner state get key

# List all keys
sloth-runner state list

# Delete a key
sloth-runner state delete key

# View statistics
sloth-runner state stats
```

## Using Stacks in Workflows

### Automatic Stack Integration

Every workflow automatically gets a stack. The stack functions are available globally in Lua:

```lua
-- Get current stack information
local stack_name = stack.get_name()
local stack_id = stack.get_id()
local stack_status = stack.get_status()

-- Set/get outputs
stack.set_output("web_url", "https://example.com")
local url = stack.get_output("web_url")
```

### Resource Registration

Modules can register resources for tracking. The stack system automatically handles idempotency:

```lua
-- Register a resource
local status, resource = stack.register_resource({
    type = "package",
    name = "nginx",
    module = "pkg",
    properties = {
        version = "1.18.0",
        state = "installed"
    }
})

-- status can be: "created", "changed", or "unchanged"
if status == "unchanged" then
    print("Package already installed with correct version")
elseif status == "changed" then
    print("Package version was updated")
elseif status == "created" then
    print("Package was installed")
end
```

### Resource State Updates

After applying changes, update the resource state:

```lua
-- Mark resource as successfully applied
stack.update_resource("package", "nginx", {
    state = "applied"
})

-- Mark resource as failed
stack.update_resource("package", "nginx", {
    state = "failed",
    error = "Installation failed: permission denied"
})
```

## Complete Example: Idempotent Web Server Setup

```lua
workflow({
    name = "web-server-setup",
    description = "Idempotent web server configuration"
})

-- This task will only execute changes when needed
task({
    name = "install-nginx",
    run = function()
        -- Check and install nginx
        local status = pkg.install({
            name = "nginx",
            state = "present"
        })
        
        if not status.changed then
            print("✓ nginx already installed")
        else
            print("✓ nginx installed")
        end
    end
})

task({
    name = "configure-nginx",
    depends_on = {"install-nginx"},
    run = function()
        -- Copy configuration file
        local result = file_ops.copy({
            src = "/configs/nginx.conf",
            dest = "/etc/nginx/nginx.conf",
            mode = "0644"
        })
        
        if not result.changed then
            print("✓ nginx.conf already up to date")
        else
            print("✓ nginx.conf updated")
            
            -- Only restart if config changed
            systemd.restart({name = "nginx"})
        end
    end
})

task({
    name = "ensure-service-running",
    depends_on = {"configure-nginx"},
    run = function()
        local status = systemd.ensure({
            name = "nginx",
            state = "started",
            enabled = true
        })
        
        if not status.changed then
            print("✓ nginx already running and enabled")
        else
            print("✓ nginx started and enabled")
        end
        
        -- Export service status
        stack.set_output("nginx_status", "running")
        stack.set_output("nginx_port", "80")
    end
})
```

## How Idempotency Works Internally

### 1. Checksum-Based Change Detection

For file operations, checksums are computed and compared:

```lua
-- Internal implementation in file_ops.copy
local src_checksum = compute_checksum(src)
local dst_checksum = compute_checksum(dst)

if src_checksum == dst_checksum then
    return {changed = false}  -- Skip copy
else
    -- Perform copy
    return {changed = true}
end
```

### 2. State Comparison

For configuration resources, properties are hashed and compared:

```lua
-- Internal stack resource tracking
local existing_resource = stack.get_resource("package", "nginx")

if existing_resource then
    local new_checksum = sha256(json.encode(new_properties))
    if new_checksum == existing_resource.checksum then
        -- No changes needed
        return "unchanged"
    else
        -- Update needed
        return "changed"
    end
else
    -- New resource
    return "created"
end
```

### 3. Drift Detection

The stack system can detect when resources have drifted from their desired state:

```bash
# Check for drift in a stack
sloth-runner stack drift my-infrastructure

# Show resources that have drifted
sloth-runner stack resources my-infrastructure --state drift
```

## Module-Specific Idempotency

### Package Module (pkg)

```lua
-- Only installs if package is missing or version differs
pkg.install({
    name = "docker",
    version = "20.10.0"
})
```

### User Module (user)

```lua
-- Only creates user if they don't exist
user.create({
    name = "appuser",
    shell = "/bin/bash",
    home = "/home/appuser"
})

-- Only modifies if properties changed
user.modify({
    name = "appuser",
    shell = "/bin/zsh"  -- Only updates shell if different
})
```

### Systemd Module

```lua
-- Only starts service if not running
-- Only enables if not enabled
systemd.ensure({
    name = "docker",
    state = "started",
    enabled = true
})
```

### File Operations

```lua
-- Only copies if files differ
file_ops.copy({
    src = "/src/file",
    dest = "/dst/file"
})

-- Only applies changes if line missing/different
file_ops.lineinfile({
    path = "/etc/config",
    line = "setting=value",
    regexp = "^setting="
})
```

## Best Practices

### 1. Always Use Stack Functions

```lua
-- Good: Track outputs in stack
stack.set_output("db_connection", connection_string)

-- Avoid: Using global variables (lost between runs)
_G.db_connection = connection_string
```

### 2. Handle Both Changed and Unchanged States

```lua
local result = pkg.install({name = "nginx"})

if result.changed then
    print("nginx was installed")
    -- Perform post-installation tasks
else
    print("nginx already present")
    -- Skip unnecessary work
end
```

### 3. Use Dependencies to Ensure Ordering

```lua
task({
    name = "configure",
    depends_on = {"install"},  -- Runs after install
    run = function()
        -- Configuration logic
    end
})
```

### 4. Register Custom Resources

For custom logic, explicitly register resources:

```lua
task({
    name = "custom-setup",
    run = function()
        local status, res = stack.register_resource({
            type = "custom",
            name = "my-resource",
            module = "custom",
            properties = {
                setting1 = "value1",
                setting2 = "value2"
            }
        })
        
        if status == "unchanged" then
            print("Resource already configured")
            return
        end
        
        -- Perform actual changes
        do_custom_setup()
        
        -- Mark as applied
        stack.update_resource("custom", "my-resource", {
            state = "applied"
        })
    end
})
```

## Querying Stack State

### From CLI

```bash
# Export stack state to JSON
sloth-runner stack export my-infrastructure > state.json

# List resources in a stack
sloth-runner stack resources my-infrastructure
```

### From Lua

```lua
-- Check if resource exists
if stack.resource_exists("package", "nginx") then
    local resource = stack.get_resource("package", "nginx")
    print("Resource state:", resource.state)
    print("Last applied:", resource.last_applied)
end
```

## Stack Persistence

Stacks are persisted in SQLite databases:

- **Default Location**: `/etc/sloth-runner/stacks.db`
- **User Location**: `~/.sloth-runner/stacks.db`
- **Custom Location**: Use `--db` flag

The database schema tracks:
- Stack metadata (name, version, status, created_at, updated_at)
- Resources (type, name, properties, checksum, state)
- Execution history
- Outputs and configuration

## Advanced Features

### Parallel Execution with Idempotency

```lua
-- Each goroutine gets idempotency guarantees
goroutine.map({"server1", "server2", "server3"}, function(server)
    local status = pkg.install({
        name = "nginx",
        delegate_to = server
    })
    
    -- Each server only installs if needed
    print(server .. ": " .. (status.changed and "installed" or "already present"))
end)
```

### Conditional Resource Management

```lua
task({
    name = "setup-database",
    run = function()
        local db_exists = stack.resource_exists("database", "mydb")
        
        if not db_exists then
            -- Create new database
            create_database("mydb")
            
            stack.register_resource({
                type = "database",
                name = "mydb",
                module = "custom",
                properties = {version = "1.0"}
            })
        else
            print("Database already exists")
        end
    end
})
```

## Summary

Sloth Runner's stack management provides:

1. **Idempotency**: Resources only change when needed
2. **State Tracking**: Full history of what was created/modified
3. **Drift Detection**: Know when infrastructure has changed
4. **Parallel Safety**: Goroutines work with idempotent resources
5. **Audit Trail**: Complete execution history

This makes Sloth Runner ideal for:
- Infrastructure as Code (IaC)
- Configuration Management
- Deployment Automation
- Compliance and Auditing
- GitOps Workflows
