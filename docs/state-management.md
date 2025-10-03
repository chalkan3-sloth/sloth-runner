# State Management and Idempotency

## Overview

Sloth Runner includes a built-in state management system that enables idempotent operations. This ensures that running the same task multiple times produces the same result without unintended side effects.

## State Storage

State is stored in SQLite databases, one per agent:

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

# Filter by resource type
sloth-runner state list file

# Output as JSON
sloth-runner state list --output json
```

### Show State Details

View detailed information about a specific state:

```bash
# Show state details
sloth-runner state show user:john

# JSON output
sloth-runner state show user:john --output json
```

### Delete States

Remove state entries:

```bash
# Delete specific state (with confirmation)
sloth-runner state delete user:john

# Skip confirmation
sloth-runner state delete user:john --yes
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

## How Idempotency Works

### File Module Example

When you use file operations, the module automatically tracks state:

```lua
task({
    name = "create-config",
    run = function()
        -- This operation is idempotent
        file.copy({
            src = "/templates/nginx.conf",
            dest = "/etc/nginx/nginx.conf",
            mode = "0644"
        })
        
        -- Only updates if content changes
        file.lineinfile({
            path = "/etc/hosts",
            line = "192.168.1.100 myserver",
            state = "present"
        })
    end
})
```

### User Module Example

```lua
task({
    name = "create-users",
    run = function()
        -- Creates user only if doesn't exist
        -- Updates only if attributes changed
        user.create({
            name = "deploy",
            shell = "/bin/bash",
            home = "/home/deploy",
            groups = {"sudo", "docker"}
        })
    end
})
```

### SSH Module Example

```lua
task({
    name = "configure-ssh",
    run = function()
        -- Adds key only if not present
        ssh.authorized_key({
            user = "deploy",
            key = "ssh-rsa AAAAB3...",
            state = "present"
        })
        
        -- Updates SSH config only if changed
        ssh.config({
            host = "production",
            options = {
                Hostname = "192.168.1.100",
                User = "deploy",
                Port = 22
            }
        })
    end
})
```

## State Tracking Internals

### Automatic State Detection

Modules automatically:

1. **Calculate checksums** of resource configurations
2. **Compare with stored state** before executing
3. **Skip execution** if state matches (no changes)
4. **Update state** after successful execution
5. **Report changes** (changed/unchanged)

### State Keys

State keys follow the pattern:

```
{module}:{resource_type}:{resource_id}
```

Examples:
- `file:copy:/etc/nginx/nginx.conf`
- `user:create:deploy`
- `ssh:authorized_key:deploy:ssh-rsa-AAAA...`
- `systemd:service:nginx`

### Checksum Calculation

Checksums are SHA-256 hashes of:
- Resource configuration
- Target state
- Relevant attributes

Example for file copy:
```go
checksum = sha256({
    src_path,
    dest_path,
    file_mode,
    owner,
    group,
    content_hash
})
```

## Benefits

### 1. Safe Re-runs

Run tasks multiple times safely:

```bash
# First run: creates everything
sloth-runner run -f setup.sloth

# Second run: skips unchanged resources
sloth-runner run -f setup.sloth
```

### 2. Incremental Updates

Only applies changes:

```lua
-- First run: creates 3 users
for _, user in ipairs({"alice", "bob", "charlie"}) do
    user.create({name = user})
end

-- Second run with added user: only creates "dave"
for _, user in ipairs({"alice", "bob", "charlie", "dave"}) do
    user.create({name = user})
end
```

### 3. Drift Detection

Detect manual changes:

```lua
task({
    name = "check-config",
    run = function()
        -- Will report changes if file was modified
        local result = file.copy({
            src = "/templates/config.ini",
            dest = "/app/config.ini"
        })
        
        if result.changed then
            print("WARNING: Configuration was out of sync")
        end
    end
})
```

## Advanced Usage

### Manual State Control

Force resource recreation:

```lua
-- Delete state before running
state.delete("user:create:deploy")

-- Then recreate
user.create({name = "deploy"})
```

### State Locks

Prevent concurrent modifications:

```lua
task({
    name = "critical-update",
    run = function()
        state.lock("deployment", "task-123", 300) -- 5 minute timeout
        
        -- Critical operations here
        
        state.unlock("deployment", "task-123")
    end
})
```

### Cross-Agent State

Share state between agents:

```lua
-- On agent1: store result
state.set("deployment:version", "v1.2.3")

-- On agent2: read result
local version = state.get("deployment:version")
```

## Best Practices

### 1. Use Descriptive Resource Names

```lua
-- Good
user.create({name = "app-deploy"})

-- Better - includes purpose
user.create({
    name = "app-deploy",
    comment = "Application deployment user"
})
```

### 2. Group Related Resources

```lua
task({
    name = "setup-webserver",
    run = function()
        -- Group all webserver setup
        user.create({name = "www-data"})
        file.copy({src = "nginx.conf", dest = "/etc/nginx/nginx.conf"})
        systemd.enable({name = "nginx"})
    end
})
```

### 3. Handle State Cleanup

```lua
task({
    name = "cleanup",
    run = function()
        -- Remove resource and its state
        user.delete({name = "old-user"})
        state.delete("user:create:old-user")
    end
})
```

### 4. Monitor State Size

```bash
# Check state database size
sloth-runner state stats

# Clean old states periodically
sloth-runner state clear --yes
```

## Troubleshooting

### State Out of Sync

If manual changes cause state mismatch:

```bash
# Option 1: Delete specific state
sloth-runner state delete file:copy:/etc/nginx/nginx.conf

# Option 2: Clear all and re-run
sloth-runner state clear --yes
sloth-runner run -f setup.sloth
```

### Locked Resources

If task fails with lock error:

```bash
# Check locked resources
sloth-runner state list lock

# Manually unlock (use with caution)
state.unlock("resource-name", "holder-id")
```

## Modules with Idempotency

The following modules support idempotent operations:

### File Operations
- `file.copy()` - Copy files
- `file.template()` - Render templates
- `file.lineinfile()` - Manage file lines
- `file.blockinfile()` - Manage file blocks
- `file.replace()` - Replace content

### User Management
- `user.create()` - Create/update users
- `user.delete()` - Remove users
- `user.modify()` - Modify user attributes

### SSH Configuration
- `ssh.authorized_key()` - Manage SSH keys
- `ssh.config()` - Manage SSH config
- `ssh.known_hosts()` - Manage known hosts

### System Services
- `systemd.enable()` - Enable services
- `systemd.start()` - Start services
- `systemd.reload()` - Reload services

### Package Management
- `pkg.install()` - Install packages
- `pkg.remove()` - Remove packages
- `pkg.update()` - Update packages

## Future Enhancements

- Remote state backends (S3, etcd, Consul)
- State encryption
- State versioning and history
- State import/export
- Web UI for state visualization
