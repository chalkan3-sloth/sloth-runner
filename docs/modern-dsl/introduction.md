# 🎨 Modern DSL Introduction

## Overview

The **Modern DSL** is Sloth Runner's Lua-based workflow definition language. Workflows are defined in `.sloth` files using simple, readable Lua syntax with access to powerful built-in modules.

## Why Modern DSL?

- **🚀 Simple & Readable**: Clean Lua syntax, easy to understand
- **📦 Global Modules**: All modules available without imports
- **🔄 Dynamic**: Use Lua's full power - loops, conditionals, functions
- **⚡ Fast**: Direct Lua execution, no YAML parsing
- **🧩 Reusable**: Create functions for common patterns

## Basic Workflow Structure

Every `.sloth` file defines a workflow using the `workflow()` function:

```lua
workflow({
    name = "my_workflow",
    description = "What this workflow does",
    tasks = {
        {
            name = "task_name",
            description = "What this task does",
            run = function()
                -- Your task code here
                return {changed = true, message = "Task completed"}
            end
        }
    }
})
```

## Complete Example

```lua
-- Simple web server setup workflow
workflow({
    name = "web_server_setup",
    description = "Install and configure nginx",
    tasks = {
        {
            name = "update_packages",
            description = "Update package list",
            run = function()
                log.info("📦 Updating packages...")
                local success, output = pkg.update()

                if success then
                    return {changed = true, message = "Packages updated"}
                else
                    return {failed = true, message = "Update failed"}
                end
            end
        },
        {
            name = "install_nginx",
            description = "Install nginx web server",
            depends_on = {"update_packages"},
            run = function()
                log.info("📦 Installing nginx...")
                local success, output = pkg.install({"nginx", "certbot"})

                if success then
                    return {changed = true, message = "Nginx installed"}
                else
                    return {changed = false, message = "Already installed"}
                end
            end
        },
        {
            name = "start_nginx",
            description = "Start nginx service",
            depends_on = {"install_nginx"},
            run = function()
                log.info("▶️  Starting nginx...")

                -- Use systemd module
                local systemd = require("systemd")
                local success, output = systemd.start("nginx")

                if success then
                    systemd.enable("nginx")
                    return {changed = true, message = "Nginx started and enabled"}
                else
                    return {failed = true, message = "Failed to start nginx"}
                end
            end
        }
    }
})
```

## Task Structure

Each task in the `tasks` array has these fields:

```lua
{
    name = "task_name",              -- Required: unique task identifier
    description = "What it does",    -- Optional: task description
    timeout = "5m",                   -- Optional: task timeout (default: no limit)
    retries = 3,                      -- Optional: retry count on failure
    depends_on = {"other_task"},      -- Optional: task dependencies
    run = function()                  -- Required: task function
        -- Task code here
        return {changed = true, message = "Done"}
    end
}
```

## Task Return Values

Tasks must return a table with one of these formats:

```lua
-- Success with changes
return {changed = true, message = "Task completed successfully"}

-- Success without changes (idempotent)
return {changed = false, message = "Already in desired state"}

-- Failure
return {failed = true, message = "Error: something went wrong"}
```

## Task Dependencies

Use `depends_on` to control execution order:

```lua
workflow({
    name = "deployment",
    description = "Deploy application",
    tasks = {
        {
            name = "build",
            run = function()
                -- Build code
                return {changed = true, message = "Built"}
            end
        },
        {
            name = "test",
            depends_on = {"build"},  -- Runs after build
            run = function()
                -- Run tests
                return {changed = false, message = "Tests passed"}
            end
        },
        {
            name = "deploy",
            depends_on = {"build", "test"},  -- Runs after both
            run = function()
                -- Deploy application
                return {changed = true, message = "Deployed"}
            end
        }
    }
})
```

## Global Modules (No require!)

Most modules are **globally available** - just use them:

```lua
run = function()
    -- Package management
    pkg.install({"nginx", "postgresql"})
    pkg.update()
    pkg.remove("oldpackage")

    -- File operations
    file_ops.copy({src = "/source", dest = "/dest"})
    file_ops.delete({path = "/tmp/file"})
    file_ops.mkdir({path = "/opt/app", mode = "0755"})

    -- User management
    user.create({
        name = "appuser",
        shell = "/bin/bash",
        home = "/home/appuser",
        create_home = true
    })

    -- Git operations
    git.clone({
        repo = "https://github.com/user/repo",
        dest = "/opt/repo"
    })

    -- System facts
    local os_info = facts.os()
    local cpu_count = facts.cpu_count()

    -- Logging
    log.info("Information message")
    log.warn("Warning message")
    log.error("Error message")

    -- Shell commands
    local result = exec.run("hostname")
    print("Hostname: " .. result)

    return {changed = true, message = "Done"}
end
```

## Modules That Need require()

Only a few modules need `require()`:

```lua
run = function()
    -- Systemd module
    local systemd = require("systemd")
    systemd.start("nginx")
    systemd.enable("nginx")

    -- Parallel execution
    local goroutine = require("goroutine")
    local handle = goroutine.async(function()
        -- runs in parallel
    end)
    local results = goroutine.await_all({handle})

    return {changed = true, message = "Done"}
end
```

## Timeouts and Retries

Add resilience to your tasks:

```lua
{
    name = "flaky_network_call",
    timeout = "30s",      -- Task will timeout after 30 seconds
    retries = 3,          -- Retry up to 3 times on failure
    run = function()
        -- Make network call
        local response = http.get("https://api.example.com/data")
        return {changed = false, message = "Data fetched"}
    end
}
```

## Dynamic Workflows with Lua

Leverage Lua's programming capabilities:

```lua
workflow({
    name = "multi_server_setup",
    description = "Setup multiple servers",
    tasks = (function()
        local servers = {"web-01", "web-02", "web-03"}
        local tasks = {}

        -- Generate a task for each server
        for _, server in ipairs(servers) do
            table.insert(tasks, {
                name = "setup_" .. server,
                description = "Setup " .. server,
                run = function()
                    log.info("Setting up " .. server)
                    -- Setup code
                    return {changed = true, message = server .. " configured"}
                end
            })
        end

        return tasks
    end)()
})
```

## Conditional Logic

```lua
workflow({
    name = "os_specific_setup",
    description = "Setup based on OS",
    tasks = {
        {
            name = "install_packages",
            run = function()
                local os_info = facts.os()

                if os_info.family == "debian" then
                    pkg.install({"apt-transport-https", "ca-certificates"})
                    log.info("Installed Debian packages")
                elseif os_info.family == "redhat" then
                    pkg.install({"yum-utils", "device-mapper-persistent-data"})
                    log.info("Installed RedHat packages")
                else
                    log.warn("Unknown OS family: " .. os_info.family)
                end

                return {changed = true, message = "OS-specific packages installed"}
            end
        }
    }
})
```

## Error Handling

Use Lua's `pcall` for error handling:

```lua
{
    name = "safe_operation",
    run = function()
        local success, err = pcall(function()
            -- Risky operation
            file_ops.copy({
                src = "/important/file",
                dest = "/backup/file"
            })
        end)

        if success then
            return {changed = true, message = "File backed up"}
        else
            log.error("Backup failed: " .. tostring(err))
            return {failed = true, message = "Backup failed"}
        end
    end
}
```

## Parallel Execution

Execute tasks in parallel using the goroutine module:

```lua
workflow({
    name = "parallel_deployment",
    description = "Deploy to multiple servers in parallel",
    tasks = {
        {
            name = "deploy_all",
            run = function()
                local goroutine = require("goroutine")

                local servers = {"web-01", "web-02", "web-03"}
                local handles = {}

                -- Start parallel deployments
                for _, server in ipairs(servers) do
                    local handle = goroutine.async(function()
                        log.info("Deploying to " .. server)
                        -- Deployment logic here
                        goroutine.sleep(1000) -- simulate work
                        return server, "success"
                    end)
                    table.insert(handles, handle)
                end

                -- Wait for all to complete
                local results = goroutine.await_all(handles)

                -- Check results
                for _, result in ipairs(results) do
                    if result.success then
                        local server = result.values[1]
                        log.info("✅ " .. server .. " deployed")
                    else
                        log.error("❌ Deployment failed: " .. result.error)
                    end
                end

                return {changed = true, message = "All deployments completed"}
            end
        }
    }
})
```

## Complete Real-World Example

```lua
-- Production web server deployment
workflow({
    name = "production_deployment",
    description = "Deploy web application to production",
    tasks = {
        {
            name = "update_system",
            description = "Update system packages",
            timeout = "5m",
            run = function()
                pkg.update()
                pkg.install({"curl", "git", "build-essential"})
                return {changed = true, message = "System updated"}
            end
        },
        {
            name = "create_app_user",
            description = "Create application user",
            depends_on = {"update_system"},
            run = function()
                user.create({
                    name = "webapp",
                    shell = "/bin/bash",
                    home = "/opt/webapp",
                    create_home = true,
                    comment = "Web Application User"
                })
                return {changed = true, message = "User created"}
            end
        },
        {
            name = "clone_repository",
            description = "Clone application repository",
            depends_on = {"create_app_user"},
            timeout = "10m",
            retries = 3,
            run = function()
                git.clone({
                    repo = "https://github.com/company/webapp.git",
                    dest = "/opt/webapp/app",
                    branch = "main"
                })
                return {changed = true, message = "Repository cloned"}
            end
        },
        {
            name = "install_dependencies",
            description = "Install application dependencies",
            depends_on = {"clone_repository"},
            timeout = "15m",
            run = function()
                exec.run("cd /opt/webapp/app && npm install")
                return {changed = true, message = "Dependencies installed"}
            end
        },
        {
            name = "configure_nginx",
            description = "Configure nginx reverse proxy",
            depends_on = {"install_dependencies"},
            run = function()
                pkg.install({"nginx"})

                file_ops.template({
                    src = "/templates/nginx.conf.j2",
                    dest = "/etc/nginx/sites-available/webapp",
                    mode = "0644"
                })

                file_ops.link({
                    src = "/etc/nginx/sites-available/webapp",
                    dest = "/etc/nginx/sites-enabled/webapp"
                })

                local systemd = require("systemd")
                systemd.restart("nginx")
                systemd.enable("nginx")

                return {changed = true, message = "Nginx configured"}
            end
        },
        {
            name = "start_application",
            description = "Start web application",
            depends_on = {"configure_nginx"},
            run = function()
                local systemd = require("systemd")

                systemd.start("webapp")
                systemd.enable("webapp")

                -- Verify it's running
                local active, state = systemd.is_active("webapp")
                if active then
                    log.info("✅ Application is running!")
                    return {changed = true, message = "Application started"}
                else
                    return {failed = true, message = "Application failed to start: " .. state}
                end
            end
        }
    }
})
```

## Available Modules

Run `sloth-runner modules list` to see all available modules:

- `pkg` - Package management (apt, yum, dnf, pacman)
- `file_ops` - File operations
- `user` - User and group management
- `git` - Git operations
- `systemd` - Service management (requires `require()`)
- `incus` - LXC/VM container management
- `stow` - Dotfiles management
- `facts` - System information
- `goroutine` - Parallel execution (requires `require()`)
- `exec` - Shell command execution
- `log` - Logging functions
- And many more...

## Best Practices

1. **Use descriptive names** for workflows and tasks
2. **Add descriptions** to document what each task does
3. **Handle errors** with `pcall` for critical operations
4. **Use dependencies** to control execution order
5. **Set timeouts** for network operations
6. **Add retries** for flaky operations
7. **Log appropriately** with `log.info`, `log.warn`, `log.error`
8. **Keep tasks focused** - one responsibility per task
9. **Test incrementally** - start small, add complexity

## Quick Reference Template

```lua
workflow({
    name = "workflow_name",
    description = "What this workflow does",
    tasks = {
        {
            name = "task_name",
            description = "What this task does",
            timeout = "5m",
            retries = 3,
            depends_on = {"previous_task"},
            run = function()
                -- Your code here

                -- Return success with changes
                return {changed = true, message = "Task completed"}

                -- Or return success without changes
                -- return {changed = false, message = "Already in desired state"}

                -- Or return failure
                -- return {failed = true, message = "Error occurred"}
            end
        }
    }
})
```

## Next Steps

- [📚 Module API Examples](module-api-examples.md) - Real-world module usage
- [🎯 Best Practices](best-practices.md) - Advanced patterns
- [📖 Reference Guide](reference-guide.md) - Complete API reference
- [🔧 Modules List](../../modules-list-command.md) - All available modules

Start building your workflows with this simple, powerful DSL!
