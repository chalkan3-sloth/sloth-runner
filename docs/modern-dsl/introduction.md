# 🎨 Modern DSL Introduction

## Overview

The **Modern DSL** is Sloth Runner's Lua-based workflow definition language using a **Builder Pattern** with method chaining. Workflows are defined in `.sloth` files with clean, expressive syntax.

## Why Modern DSL?

- **🎯 Fluent & Expressive**: Builder pattern with method chaining
- **📦 Global Modules**: All modules available without imports
- **🔄 Dynamic**: Use Lua's full power - loops, conditionals, functions
- **⚡ Fast**: Direct Lua execution
- **🧩 Composable**: Build complex workflows from reusable tasks

## Basic Structure

Every `.sloth` file defines tasks using the `task()` builder, then composes them into a workflow with `workflow.define()`:

```lua
-- Define a task using the builder pattern
local my_task = task("task-name")
    :description("What this task does")
    :command(function(this, params)
        -- Your task code here
        return true, "Task completed successfully"
    end)
    :build()

-- Compose tasks into a workflow
workflow
    .define("workflow_name")
    :description("What this workflow does")
    :version("1.0.0")
    :tasks({my_task})
    :on_complete(function(success, results)
        if success then
            log.info("✅ Workflow completed!")
        end
    end)
```

## Complete Example

```lua
-- Install and configure nginx on a remote server
local install_nginx = task("install-nginx")
    :description("Install nginx package")
    :delegate_to("web-server")
    :command(function(this, params)
        local success, msg = pkg.install({
            packages = {"nginx", "certbot"}
        })

        if not success then
            return false, "Failed to install: " .. tostring(msg)
        end

        return true, "Nginx installed successfully"
    end)
    :build()

local start_nginx = task("start-nginx")
    :description("Start and enable nginx service")
    :delegate_to("web-server")
    :command(function(this, params)

        local success, msg = systemd.start("nginx")
        if not success then
            return false, "Failed to start nginx"
        end

        systemd.enable("nginx")
        return true, "Nginx started and enabled"
    end)
    :build()

workflow
    .define("nginx_deployment")
    :description("Deploy nginx web server")
    :version("1.0.0")
    :tasks({install_nginx, start_nginx})
    :config({
        timeout = "10m",
        max_parallel_tasks = 1
    })
    :on_complete(function(success, results)
        if success then
            log.info("🎉 Nginx deployment completed!")
        else
            log.error("❌ Deployment failed")
        end
    end)
```

## Task Builder API

The `task()` builder provides these methods:

```lua
local my_task = task("unique-task-name")
    :description("Human-readable description")
    :delegate_to("agent-name")         -- Execute on remote agent
    :user("username")                   -- Run as specific user
    :workdir("/path/to/directory")      -- Set working directory
    :timeout("5m")                      -- Task timeout
    :retries(3)                         -- Retry count on failure
    :command(function(this, params)     -- Task function
        -- Task implementation
        return true, "Success message"
        -- or
        -- return false, "Error message"
    end)
    :build()  -- MUST call .build() to finalize the task
```

### Task Return Values

Tasks return two values: `(success, message)`

```lua
:command(function(this, params)
    -- Success
    return true, "Operation completed successfully"

    -- Failure
    return false, "Error: operation failed"

    -- Optional third return value for data
    return true, "Data fetched", {count = 42, items = {...}}
end)
```

## Workflow Builder API

The `workflow.define()` builder provides these methods:

```lua
workflow
    .define("workflow_name")
    :description("What this workflow does")
    :version("1.0.0")
    :tasks({task1, task2, task3})  -- Array of built tasks
    :config({
        timeout = "30m",
        max_parallel_tasks = 2
    })
    :on_complete(function(success, results)
        -- Called after workflow completes
        if success then
            log.info("All tasks completed")
        else
            log.error("Workflow failed")
        end
    end)
```

## Global Modules (No require!)

Most modules are **globally available** - just use them:

```lua
:command(function(this, params)
    -- Package management
    pkg.install({packages = {"nginx", "postgresql"}})
    pkg.update()
    pkg.remove({package = "oldpackage"})

    -- User management
    user.create({
        username = "webuser",
        password = "changeme123",
        home = "/home/webuser",
        shell = "/bin/bash",
        groups = {"wheel", "docker"},
        create_home = true
    })

    -- Git operations
    git.clone({
        url = "https://github.com/user/repo",
        local_path = "/opt/repo",
        clean = false
    })

    -- File operations
    file_ops.copy({src = "/source", dest = "/dest"})
    file_ops.mkdir({path = "/opt/app", mode = "0755"})

    -- Stow (dotfiles management)
    stow.link({
        package = "zsh",
        source_dir = "/home/user/dotfiles",
        target_dir = "/home/user",
        create_target = true
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

    return true, "All operations completed"
end)
```

## Modules That Need require()

Only a few modules need `require()`:

```lua
:command(function(this, params)
    -- Systemd module
    systemd.start("nginx")
    systemd.enable("nginx")

    -- Parallel execution
    local goroutine = require("goroutine")
    local handle = goroutine.async(function()
        -- runs in parallel
    end)
    local results = goroutine.await_all({handle})

    return true, "Done"
end)
```

## Task Dependencies

Tasks are executed in the order they appear in the `tasks` array. For complex dependencies, use multiple workflows or order tasks explicitly:

```lua
-- Tasks execute in order: build → test → deploy
workflow
    .define("deployment")
    :tasks({
        build_task,    -- Runs first
        test_task,     -- Runs after build
        deploy_task    -- Runs after test
    })
```

## Remote Execution

Use `:delegate_to()` to execute tasks on remote agents:

```lua
local setup_server = task("setup-server")
    :description("Setup remote server")
    :delegate_to("production-server")  -- Execute on agent
    :user("deployer")                   -- Run as deployer user
    :workdir("/opt/app")                -- Set working directory
    :command(function(this, params)
        pkg.install({packages = {"nginx"}})
        return true, "Server configured"
    end)
    :build()
```

## Dynamic Workflows

Use Lua to generate tasks programmatically:

```lua
-- Generate tasks for multiple servers
local servers = {"web-01", "web-02", "web-03"}
local tasks = {}

for _, server in ipairs(servers) do
    local t = task("setup-" .. server)
        :description("Setup " .. server)
        :delegate_to(server)
        :command(function(this, params)
            pkg.install({packages = {"nginx"}})
            log.info("✓ " .. server .. " configured")
            return true, server .. " ready"
        end)
        :build()

    table.insert(tasks, t)
end

workflow
    .define("multi_server_setup")
    :description("Setup multiple servers")
    :tasks(tasks)
```

## Conditional Logic

Use Lua conditionals inside task commands:

```lua
local os_specific_setup = task("os-setup")
    :description("Install OS-specific packages")
    :command(function(this, params)
        local os_info = facts.os()

        if os_info.family == "debian" then
            pkg.install({packages = {"apt-transport-https"}})
        elseif os_info.family == "redhat" then
            pkg.install({packages = {"yum-utils"}})
        else
            log.warn("Unknown OS: " .. os_info.family)
        end

        return true, "OS-specific setup completed"
    end)
    :build()
```

## Error Handling

Use Lua's `pcall` for safe error handling:

```lua
local safe_operation = task("safe-op")
    :description("Operation with error handling")
    :command(function(this, params)
        local success, err = pcall(function()
            file_ops.copy({
                src = "/important/file",
                dest = "/backup/file"
            })
        end)

        if success then
            return true, "File backed up successfully"
        else
            log.error("Backup failed: " .. tostring(err))
            return false, "Backup failed: " .. tostring(err)
        end
    end)
    :build()
```

## Parallel Execution

Execute operations in parallel using the goroutine module:

```lua
local parallel_deploy = task("parallel-deploy")
    :description("Deploy to multiple servers in parallel")
    :command(function(this, params)
        local goroutine = require("goroutine")

        local servers = {"web-01", "web-02", "web-03"}
        local handles = {}

        -- Start parallel deployments
        for _, server in ipairs(servers) do
            local handle = goroutine.async(function()
                log.info("Deploying to " .. server)
                -- Deployment logic
                goroutine.sleep(1000)
                return server, "success"
            end)
            table.insert(handles, handle)
        end

        -- Wait for all to complete
        local results = goroutine.await_all(handles)

        -- Process results
        for _, result in ipairs(results) do
            if result.success then
                local server_name = result.values[1]
                log.info("✅ " .. server_name .. " deployed")
            else
                log.error("❌ Failed: " .. result.error)
            end
        end

        return true, "All deployments completed"
    end)
    :build()
```

## Complete Real-World Example

```lua
-- User environment setup with dotfiles
local install_packages = task("install-packages")
    :description("Install default packages")
    :delegate_to("lady-arch")
    :command(function(this, params)
        local success, msg = pkg.install({
            packages = {"kitty-terminfo", "stow", "git", "zsh", "lsd", "fzf"}
        })

        if not success then
            return false, "Failed to install packages: " .. tostring(msg)
        end

        return true, "Packages installed successfully"
    end)
    :build()

local create_user = task("create-user")
    :description("Create and configure user")
    :delegate_to("lady-arch")
    :command(function(this, params)
        local success, msg = user.create({
            username = "igor",
            password = "changeme123",
            home = "/home/igor",
            shell = "/bin/zsh",
            groups = {"wheel"},
            create_home = true
        })

        if not success then
            return false, "Failed to create user: " .. tostring(msg)
        end

        return true, "User created successfully"
    end)
    :build()

local clone_dotfiles = task("clone-dotfiles")
    :description("Clone dotfiles repository")
    :delegate_to("lady-arch")
    :user("igor")
    :workdir("/home/igor")
    :command(function(this, params)
        local repo, err = git.clone({
            url = "https://github.com/chalkan3/dotfiles.git",
            local_path = "/home/igor/dotfiles",
            clean = false
        })

        if err then
            return false, "Failed to clone dotfiles: " .. err
        end

        if repo.exists then
            log.info("✓ Dotfiles repository already exists")
        else
            log.info("✓ Dotfiles cloned successfully")
        end

        return true, "Dotfiles ready"
    end)
    :build()

local stow_config = task("stow-zsh-config")
    :description("Stow zsh configuration files")
    :delegate_to("lady-arch")
    :user("igor")
    :command(function(this, params)
        -- Ensure target directory exists
        local ok_dir, msg_dir = stow.ensure_target({
            path = "/home/igor/.zsh",
            owner = "igor",
            mode = "0755"
        })

        if not ok_dir then
            return false, "Failed to create directory: " .. msg_dir
        end

        -- Stow configuration
        local ok_stow, msg_stow = stow.link({
            package = ".",
            source_dir = "/home/igor/dotfiles/zsh",
            target_dir = "/home/igor/.zsh",
            create_target = true,
            verbose = true
        })

        if not ok_stow then
            return false, "Failed to stow config: " .. msg_stow
        end

        return true, "Configuration stowed successfully"
    end)
    :build()

workflow
    .define("user_environment_setup")
    :description("Complete user environment setup with dotfiles")
    :version("2.0.0")
    :tasks({
        install_packages,
        create_user,
        clone_dotfiles,
        stow_config
    })
    :config({
        timeout = "30m",
        max_parallel_tasks = 1
    })
    :on_complete(function(success, results)
        if success then
            log.info("🎉 User environment setup completed successfully!")
            log.info("📋 Summary:")
            log.info("  ✓ Packages installed")
            log.info("  ✓ User created")
            log.info("  ✓ Dotfiles cloned")
            log.info("  ✓ Configuration stowed")
        else
            log.error("❌ Setup failed")
        end
    end)
```

## Available Modules

Run `sloth-runner modules list` to see all available modules:

- `pkg` - Package management (apt, yum, dnf, pacman)
- `user` - User and group management
- `file_ops` - File operations
- `git` - Git operations
- `stow` - Dotfiles management with GNU Stow
- `systemd` - Service management (requires `require()`)
- `incus` - LXC/VM container management
- `facts` - System information
- `goroutine` - Parallel execution (requires `require()`)
- `exec` - Shell command execution
- `log` - Logging functions
- And many more...

## Best Practices

1. **Always call `:build()`** - Tasks must end with `:build()`
2. **Use descriptive names** - Make task and workflow names self-documenting
3. **Add descriptions** - Document what each task does
4. **Handle errors** - Use `pcall` for critical operations
5. **Log appropriately** - Use `log.info`, `log.warn`, `log.error`
6. **Delegate wisely** - Use `:delegate_to()` for remote execution
7. **Set timeouts** - Prevent hanging tasks with `:timeout()`
8. **Keep focused** - One responsibility per task
9. **Test incrementally** - Build and test tasks individually

## Quick Reference Template

```lua
-- Define tasks
local my_task = task("task-name")
    :description("What this task does")
    :delegate_to("agent-name")  -- Optional: remote execution
    :user("username")            -- Optional: run as user
    :workdir("/path")            -- Optional: working directory
    :timeout("5m")               -- Optional: timeout
    :retries(3)                  -- Optional: retry count
    :command(function(this, params)
        -- Your code here

        -- Return success
        return true, "Success message"

        -- Or return failure
        -- return false, "Error message"
    end)
    :build()  -- Required!

-- Compose workflow
workflow
    .define("workflow_name")
    :description("What this workflow does")
    :version("1.0.0")
    :tasks({my_task})  -- Array of tasks
    :config({
        timeout = "30m",
        max_parallel_tasks = 1
    })
    :on_complete(function(success, results)
        if success then
            log.info("✅ Workflow completed!")
        end
    end)
```

## Next Steps

- [📚 Module API Examples](module-api-examples.md) - Real-world module usage
- [🎯 Best Practices](best-practices.md) - Advanced patterns
- [📖 Reference Guide](reference-guide.md) - Complete API reference
- [🔧 Modules List](../../modules-list-command.md) - All available modules

Start building powerful, composable workflows with the Modern DSL Builder Pattern!
