# Idempotency in Sloth Runner Modules

## Overview

Idempotency is a critical feature in configuration management and infrastructure automation. An operation is idempotent if running it multiple times produces the same result as running it once. In Sloth Runner, all configuration modules now support idempotency.

## What is Idempotency?

Idempotency means that you can run the same task multiple times safely without causing unwanted side effects. For example:

- Installing a package that's already installed should not reinstall it
- Creating a user that already exists should not fail or recreate
- Starting a service that's already running should not restart it

## Benefits

1. **Safety**: Run playbooks multiple times without fear
2. **Performance**: Skip unnecessary operations
3. **Predictability**: Always know what changed
4. **Debugging**: Clear feedback on what was modified

## Idempotent Modules

### Package Management (pkg)

The `pkg` module now checks if packages are already installed before attempting installation:

```lua
task({
    name = "install-nginx",
    run = function()
        -- First run: installs nginx
        -- Subsequent runs: returns changed=false
        local result = pkg.install({packages = "nginx"})
        
        if result.changed then
            print("Nginx was installed")
        else
            print("Nginx already installed")
        end
    end
})
```

**Behavior:**
- `pkg.install()`: Only installs packages that aren't already installed
- `pkg.remove()`: Only removes packages that are actually installed
- Returns `changed=true` only when actual changes are made

### User Management (user)

User operations check if the user exists before creating:

```lua
task({
    name = "create-app-user",
    run = function()
        local result = user.create({
            username = "appuser",
            home = "/home/appuser",
            shell = "/bin/bash",
            groups = {"docker", "sudo"}
        })
        
        if result.changed then
            print("User created")
        else
            print("User already exists")
        end
    end
})
```

**Behavior:**
- `user.create()`: Returns `changed=false` if user already exists
- Future: Will verify properties match desired state

### Systemd Service Management (systemd)

Service operations check current state before making changes:

```lua
task({
    name = "manage-nginx",
    run = function()
        -- Check if already running
        local start_result = systemd.start({name = "nginx"})
        
        -- Check if already enabled
        local enable_result = systemd.enable({name = "nginx"})
        
        print("Started: " .. tostring(start_result.changed))
        print("Enabled: " .. tostring(enable_result.changed))
    end
})
```

**Behavior:**
- `systemd.start()`: Returns `changed=false` if already active
- `systemd.stop()`: Returns `changed=false` if already inactive
- `systemd.enable()`: Returns `changed=false` if already enabled
- `systemd.disable()`: Returns `changed=false` if already disabled

### File Operations (file_ops)

File operations compare checksums to detect changes:

```lua
task({
    name = "copy-config",
    run = function()
        local result = file_ops.copy({
            src = "/source/nginx.conf",
            dest = "/etc/nginx/nginx.conf",
            mode = "0644"
        })
        
        if result.changed then
            print("Configuration updated")
            systemd.restart({name = "nginx"})
        else
            print("Configuration unchanged")
        end
    end
})
```

**Behavior:**
- `file_ops.copy()`: Compares checksums, only copies if different
- `file_ops.lineinfile()`: Only modifies file if line doesn't match
- `file_ops.blockinfile()`: Only updates if block content differs
- `file_ops.replace()`: Only writes if replacements are made

## Understanding the Response Format

All idempotent operations now return a table with:

```lua
{
    changed = true|false,  -- Did the operation make changes?
    message = "...",       -- Human-readable description
    -- Module-specific fields
}
```

### Example Response Patterns

**Package Installation (already installed):**
```lua
{
    changed = false,
    message = "All packages already installed"
}
```

**Package Installation (newly installed):**
```lua
{
    changed = true,
    installed = "nginx, vim",
    output = "..."
}
```

**Service Start (already running):**
```lua
{
    changed = false,
    message = "Service nginx is already active"
}
```

**Service Start (started now):**
```lua
{
    changed = true,
    message = "Service nginx started"
}
```

## Best Practices

### 1. Check Changed Flag Before Dependent Actions

```lua
task({
    name = "update-and-restart",
    run = function()
        local result = file_ops.copy({
            src = "app.conf",
            dest = "/etc/app/app.conf"
        })
        
        -- Only restart if configuration changed
        if result.changed then
            systemd.restart({name = "app"})
        end
    end
})
```

### 2. Use Idempotency for Convergent State

```lua
task({
    name = "ensure-state",
    run = function()
        -- Run multiple times, always converges to desired state
        pkg.install({packages = {"nginx", "vim", "git"}})
        
        user.create({
            username = "webuser",
            groups = {"www-data"}
        })
        
        systemd.enable({name = "nginx"})
        systemd.start({name = "nginx"})
    end
})
```

### 3. Conditional Logic Based on Changes

```lua
task({
    name = "deploy-app",
    run = function()
        local deps_changed = pkg.install({
            packages = {"python3", "python3-pip"}
        }).changed
        
        local code_changed = file_ops.copy({
            src = "app.py",
            dest = "/opt/app/app.py"
        }).changed
        
        -- Only restart if dependencies or code changed
        if deps_changed or code_changed then
            print("Changes detected, restarting service")
            systemd.restart({name = "myapp"})
        else
            print("No changes, service continues running")
        end
    end
})
```

## Testing Idempotency

To test if your tasks are idempotent:

1. **First Run**: Should make changes
   ```bash
   sloth-runner run deployment.sloth
   # Output: changed=true
   ```

2. **Second Run**: Should skip already-done work
   ```bash
   sloth-runner run deployment.sloth
   # Output: changed=false
   ```

3. **Verify No Side Effects**: Check that running twice doesn't cause issues
   ```bash
   # Run multiple times
   for i in {1..5}; do
       sloth-runner run deployment.sloth
   done
   ```

## Migration Guide

### Old Code (Non-Idempotent)

```lua
task({
    name = "setup",
    run = function()
        -- Always attempts to install
        pkg.install({packages = "nginx"})
        
        -- Always attempts to create
        user.create({username = "webuser"})
        
        -- Always starts (might fail if running)
        systemd.start({name = "nginx"})
    end
})
```

### New Code (Idempotent)

```lua
task({
    name = "setup",
    run = function()
        -- Checks first, installs only if needed
        local pkg_result = pkg.install({packages = "nginx"})
        
        -- Checks if user exists first
        local user_result = user.create({username = "webuser"})
        
        -- Checks if already running
        local start_result = systemd.start({name = "nginx"})
        
        -- Report what changed
        if pkg_result.changed or user_result.changed or start_result.changed then
            print("System state updated")
        else
            print("System already in desired state")
        end
    end
})
```

## Module Support Status

| Module | Idempotent | Notes |
|--------|-----------|-------|
| pkg | ✅ | Checks package installation status |
| user | ✅ | Checks user existence |
| systemd | ✅ | Checks service state |
| file_ops.copy | ✅ | Compares checksums |
| file_ops.lineinfile | ✅ | Checks line existence |
| file_ops.blockinfile | ✅ | Checks block content |
| file_ops.replace | ✅ | Compares before/after |
| ssh | 🔄 | Planned |
| incus | 🔄 | Planned |

## Future Enhancements

1. **Property Verification**: For `user.create()`, verify all properties match (not just existence)
2. **Atomic Operations**: Ensure all-or-nothing for complex operations
3. **Diff Mode**: Show what would change without making changes
4. **Check Mode**: Dry-run to preview changes
5. **Change Tracking**: Detailed logs of what changed

## Conclusion

Idempotency makes Sloth Runner safe, predictable, and efficient. You can now run your automation scripts confidently, knowing they won't cause unwanted side effects on repeated executions.
