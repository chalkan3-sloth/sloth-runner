# 📖 Modern DSL Reference Guide

Complete API reference for Sloth Runner's Modern DSL, covering all modules, functions, and patterns.

## Table of Contents
- [Core Syntax](#core-syntax)
- [Module Reference](#module-reference)
- [Global Functions](#global-functions)
- [Task Definition](#task-definition)
- [Group Definition](#group-definition)
- [State Management](#state-management)
- [Error Handling](#error-handling)
- [Advanced Patterns](#advanced-patterns)

## Core Syntax

### Basic Structure

```lua
-- File metadata
name = "workflow-name"
version = "1.0.0"
description = "Workflow description"
author = "Your Name"
tags = {"automation", "deployment"}

-- Import modules (optional - most are global)
local utils = require("lib/utils")

-- Define groups
group "group_name" {
    -- Group properties and tasks
}

-- Define individual tasks
task "task_name" {
    -- Task properties
}
```

## Module Reference

### Package Management (`pkg`)

```lua
-- Update package repositories
task "update_repos" {
    module = "pkg",
    action = "update"
}

-- Install packages
task "install_packages" {
    module = "pkg",
    action = "install",
    packages = {"nginx", "postgresql", "redis"},
    state = "present"  -- present, absent, latest
}

-- Remove packages
task "remove_packages" {
    module = "pkg",
    action = "remove",
    packages = {"apache2"},
    state = "absent",
    purge = true  -- Remove config files
}

-- Upgrade all packages
task "upgrade_system" {
    module = "pkg",
    action = "upgrade",
    dist_upgrade = true  -- For Debian/Ubuntu
}
```

### Systemd Service Management

```lua
-- Manage service state
task "manage_nginx" {
    module = "systemd",
    action = "service",
    name = "nginx",
    state = "started",  -- started, stopped, restarted, reloaded
    enabled = true,     -- Enable on boot
    daemon_reload = true
}

-- Multiple services
task "manage_services" {
    module = "systemd",
    action = "batch",
    services = {
        {name = "nginx", state = "started", enabled = true},
        {name = "postgresql", state = "started", enabled = true},
        {name = "redis", state = "started", enabled = false}
    }
}

-- Custom systemd operations
task "reload_daemon" {
    module = "systemd",
    action = "daemon_reload"
}

task "mask_service" {
    module = "systemd",
    action = "mask",
    name = "bluetooth"
}
```

### Docker Operations

```lua
-- Pull images
task "pull_images" {
    module = "docker",
    action = "pull",
    images = {"nginx:latest", "postgres:14", "redis:7-alpine"}
}

-- Run container
task "run_database" {
    module = "docker",
    action = "container",
    name = "myapp-db",
    image = "postgres:14",
    state = "started",
    ports = {"5432:5432"},
    volumes = {
        "/data/postgres:/var/lib/postgresql/data",
        "myapp-config:/etc/postgresql"
    },
    environment = {
        POSTGRES_DB = "myapp",
        POSTGRES_USER = "appuser",
        POSTGRES_PASSWORD = "${DB_PASSWORD}"  -- From environment
    },
    networks = {"myapp-network"},
    restart = "unless-stopped",
    healthcheck = {
        test = ["CMD", "pg_isready", "-U", "appuser"],
        interval = "30s",
        timeout = "10s",
        retries = 3
    }
}

-- Docker Compose
task "deploy_stack" {
    module = "docker",
    action = "compose",
    project = "myapp",
    compose_files = {"docker-compose.yml", "docker-compose.prod.yml"},
    env_file = ".env.production",
    pull = true,
    build = true,
    remove_orphans = true
}

-- Build image
task "build_app" {
    module = "docker",
    action = "build",
    tag = "myapp:${VERSION}",
    context = ".",
    dockerfile = "Dockerfile.prod",
    build_args = {
        VERSION = "${VERSION}",
        BUILD_DATE = os.date("%Y-%m-%d")
    },
    labels = {
        maintainer = "team@example.com",
        version = "${VERSION}"
    }
}
```

### Terraform Infrastructure

```lua
-- Initialize Terraform
task "terraform_init" {
    module = "terraform",
    action = "init",
    working_dir = "./terraform/production",
    backend = true,
    backend_config = {
        bucket = "terraform-state",
        key = "production/terraform.tfstate",
        region = "us-east-1"
    }
}

-- Plan changes
task "terraform_plan" {
    module = "terraform",
    action = "plan",
    working_dir = "./terraform/production",
    var_file = "production.tfvars",
    variables = {
        environment = "production",
        region = "us-east-1"
    },
    out = "production.tfplan",
    refresh = true
}

-- Apply changes
task "terraform_apply" {
    module = "terraform",
    action = "apply",
    working_dir = "./terraform/production",
    plan_file = "production.tfplan",
    auto_approve = true,
    parallelism = 10
}

-- Destroy infrastructure
task "terraform_destroy" {
    module = "terraform",
    action = "destroy",
    working_dir = "./terraform/production",
    var_file = "production.tfvars",
    auto_approve = false,  -- Require confirmation
    target = ["aws_instance.web"]  -- Specific resources
}
```

### Git Operations

```lua
-- Clone repository
task "clone_repo" {
    module = "git",
    action = "clone",
    repository = "https://github.com/myorg/myapp.git",
    dest = "/opt/myapp",
    branch = "main",
    depth = 1,  -- Shallow clone
    single_branch = true
}

-- Pull latest changes
task "update_code" {
    module = "git",
    action = "pull",
    repo_path = "/opt/myapp",
    branch = "main",
    rebase = true,
    force = false
}

-- Checkout branch/tag
task "checkout_version" {
    module = "git",
    action = "checkout",
    repo_path = "/opt/myapp",
    ref = "v${VERSION}",  -- Can be branch, tag, or commit
    create = false  -- Don't create new branch
}

-- Create and push tag
task "tag_release" {
    module = "git",
    action = "tag",
    repo_path = "/opt/myapp",
    tag = "v${VERSION}",
    message = "Release version ${VERSION}",
    push = true,
    force = false
}
```

### File System Operations

```lua
-- Create directory
task "create_dirs" {
    module = "fs",
    action = "mkdir",
    paths = {"/opt/myapp", "/var/log/myapp", "/etc/myapp"},
    mode = "0755",
    parents = true  -- Create parent directories
}

-- Copy files
task "copy_config" {
    module = "fs",
    action = "copy",
    src = "./configs/app.conf",
    dest = "/etc/myapp/app.conf",
    mode = "0644",
    owner = "appuser",
    group = "appgroup",
    backup = true  -- Backup existing file
}

-- Template file
task "render_config" {
    module = "fs",
    action = "template",
    src = "./templates/nginx.conf.j2",
    dest = "/etc/nginx/sites-available/myapp",
    variables = {
        server_name = "example.com",
        app_port = 3000,
        ssl_cert = "/etc/ssl/certs/example.com.crt"
    },
    mode = "0644",
    validate = "nginx -t -c %s"  -- Validate before replacing
}

-- Manage symlinks
task "create_symlink" {
    module = "fs",
    action = "symlink",
    src = "/etc/nginx/sites-available/myapp",
    dest = "/etc/nginx/sites-enabled/myapp",
    force = true
}

-- Archive operations
task "create_backup" {
    module = "fs",
    action = "archive",
    src = "/opt/myapp",
    dest = "/backups/myapp-${TIMESTAMP}.tar.gz",
    format = "tar.gz",
    exclude = {".git", "node_modules", "*.log"}
}
```

### Incus/LXD Container Management

```lua
-- Create container
task "create_container" {
    module = "incus",
    action = "container",
    name = "web-01",
    image = "ubuntu:22.04",
    state = "started",
    profiles = {"default", "web"},
    config = {
        "limits.cpu" = "2",
        "limits.memory" = "2GB",
        "boot.autostart" = "true"
    },
    devices = {
        root = {
            type = "disk",
            pool = "default",
            size = "20GB"
        },
        eth0 = {
            type = "nic",
            network = "lxdbr0",
            ["ipv4.address"] = "10.0.0.100"
        }
    }
}

-- Execute commands in container
task "configure_container" {
    module = "incus",
    action = "exec",
    container = "web-01",
    commands = {
        "apt-get update",
        "apt-get install -y nginx",
        "systemctl start nginx"
    }
}

-- Container snapshots
task "snapshot_container" {
    module = "incus",
    action = "snapshot",
    container = "web-01",
    snapshot_name = "before-upgrade",
    stateful = false
}

-- Copy files to container
task "deploy_to_container" {
    module = "incus",
    action = "file_push",
    container = "web-01",
    src = "./app.tar.gz",
    dest = "/tmp/app.tar.gz"
}
```

### Infrastructure Testing

```lua
-- Test suite
task "validate_infrastructure" {
    module = "infra_test",
    action = "suite",
    tests = {
        -- File tests
        {type = "file_exists", path = "/etc/nginx/nginx.conf"},
        {type = "file_contains", path = "/etc/nginx/nginx.conf",
         pattern = "worker_processes"},
        {type = "file_mode", path = "/etc/ssl/private", mode = "0700"},
        {type = "file_owner", path = "/var/www", owner = "www-data"},

        -- Directory tests
        {type = "directory_exists", path = "/opt/myapp"},
        {type = "directory_empty", path = "/tmp/build"},

        -- Service tests
        {type = "service_running", name = "nginx"},
        {type = "service_enabled", name = "nginx"},

        -- Port tests
        {type = "port_listening", port = 80, protocol = "tcp"},
        {type = "port_reachable", host = "example.com", port = 443},

        -- Process tests
        {type = "process_running", pattern = "nginx: master"},
        {type = "process_count", pattern = "worker", min = 2, max = 8},

        -- Command tests
        {type = "command_succeeds", command = "nginx -t"},
        {type = "command_output", command = "nginx -v",
         contains = "nginx/1."},

        -- HTTP tests
        {type = "http_ok", url = "http://localhost/health"},
        {type = "http_status", url = "https://example.com", status = 200},
        {type = "http_contains", url = "http://localhost",
         text = "Welcome"},

        -- DNS tests
        {type = "dns_resolves", hostname = "example.com"},
        {type = "dns_record", hostname = "example.com",
         record_type = "A", expected = "93.184.216.34"}
    }
}
```

## Global Functions

### Logging

```lua
-- Log levels
log.debug("Debug message")
log.info("Information message")
log.warn("Warning message")
log.error("Error message")
log.fatal("Fatal error - will exit")

-- Structured logging
log.info("Processing", {
    file = filename,
    size = filesize,
    user = username
})
```

### State Management

```lua
-- Set state values
state.set("deployment.version", "1.2.3")
state.set("deployment.timestamp", os.time())
state.set("servers", {"web-01", "web-02", "web-03"})

-- Get state values
local version = state.get("deployment.version")
local servers = state.get("servers", {})  -- Default value

-- Check existence
if state.has("deployment.version") then
    -- Key exists
end

-- Delete state
state.delete("temporary.data")

-- Clear all state
state.clear()
```

### Facts (System Information)

```lua
-- Get OS information
local os_info = facts.get_os()
-- Returns: {family = "debian", distro = "ubuntu", version = "22.04"}

-- Get network information
local network = facts.get_network()
-- Returns: {interfaces = {...}, routes = {...}}

-- Get hardware information
local hardware = facts.get_hardware()
-- Returns: {cpu_count = 8, memory_mb = 16384, disk_gb = 500}

-- Get all facts
local all_facts = facts.get_all()

-- Custom fact collection
facts.set_custom("app_version", get_app_version())
local custom = facts.get_custom("app_version")
```

### HTTP Operations

```lua
-- GET request
local response = http.get("https://api.example.com/data", {
    headers = {
        ["Authorization"] = "Bearer " .. token,
        ["Accept"] = "application/json"
    },
    timeout = 30
})

-- POST request
local response = http.post("https://api.example.com/data", {
    body = json.encode(data),
    headers = {
        ["Content-Type"] = "application/json"
    }
})

-- PUT request
local response = http.put("https://api.example.com/data/123", {
    body = json.encode(updated_data)
})

-- DELETE request
local response = http.delete("https://api.example.com/data/123")

-- Response handling
if response.status == 200 then
    local data = json.decode(response.body)
    log.info("Success: " .. data.message)
else
    log.error("HTTP error: " .. response.status)
end
```

### JSON Operations

```lua
-- Encode to JSON
local json_string = json.encode({
    name = "myapp",
    version = "1.2.3",
    features = {"auth", "api", "ui"}
})

-- Decode from JSON
local data = json.decode(json_string)

-- Pretty print
local pretty = json.encode_pretty(data, "  ")  -- 2-space indent

-- Read JSON file
local config = json.decode(fs.read("/etc/myapp/config.json"))

-- Write JSON file
fs.write("/etc/myapp/config.json", json.encode_pretty(config))
```

### Environment Variables

```lua
-- Get environment variable
local home = os.getenv("HOME")
local token = os.getenv("API_TOKEN") or "default-token"

-- Set environment variable (for child processes)
os.setenv("APP_ENV", "production")

-- Check if variable exists
if os.getenv("DEBUG") then
    log.debug_enabled = true
end
```

### Execution Control

```lua
-- Execute command
local result = exec.run("ls -la /opt")
if result.exit_code == 0 then
    log.info("Output: " .. result.stdout)
else
    log.error("Error: " .. result.stderr)
end

-- Execute with timeout
local result = exec.run_timeout("long-running-command", 60)

-- Execute in background
local pid = exec.background("server --port 8080")

-- Check if process is running
if exec.is_running(pid) then
    log.info("Server is running")
end

-- Kill process
exec.kill(pid, 15)  -- SIGTERM
```

## Task Definition

### Complete Task API

```lua
task "complete_example" {
    -- Basic properties
    module = "custom",  -- Module to use
    action = "execute", -- Action within module
    description = "Complete task example",

    -- Dependencies
    depends_on = {"other_task", "another_task"},

    -- Execution control
    enabled = true,  -- Enable/disable task
    when = function() return os.getenv("ENV") == "prod" end,
    unless = function() return file_exists("/tmp/skip") end,

    -- Timeout and retries
    timeout = 300,  -- Seconds
    retries = 3,
    retry_delay = 5,  -- Seconds between retries

    -- Main execution
    execute = function(params)
        -- Task logic here
        return true, "Success message", {output = "data"}
    end,

    -- Hooks
    before = function()
        log.info("Before task execution")
    end,

    after = function()
        log.info("After task execution")
    end,

    on_success = function(result)
        log.info("Task succeeded: " .. result.message)
    end,

    on_failure = function(error)
        log.error("Task failed: " .. error)
        -- Cleanup or rollback
    end,

    on_retry = function(attempt, error)
        log.warn("Retry " .. attempt .. ": " .. error)
    end,

    -- Error handling
    ignore_errors = false,
    continue_on_error = false,

    -- Resource limits
    max_memory = "2GB",
    max_cpu = "2",

    -- Variables
    vars = {
        custom_var = "value"
    },

    -- Register output
    register = "task_result"  -- Store result in variable
}
```

## Group Definition

### Complete Group API

```lua
group "group_name" {
    -- Basic properties
    description = "Group description",
    tags = {"deployment", "production"},

    -- Execution control
    enabled = true,
    when = function() return should_run() end,

    -- Parallel execution
    parallel = true,  -- Run tasks in parallel
    max_parallel = 3,  -- Max parallel tasks

    -- Error handling
    on_error = "continue",  -- continue, stop, rollback
    ignore_failures = false,

    -- Hooks
    before = function()
        log.info("Starting group execution")
    end,

    after = function()
        log.info("Group execution completed")
    end,

    on_failure = function(failed_tasks)
        log.error("Group failed: " .. #failed_tasks .. " tasks failed")
    end,

    -- Tasks
    task "task1" {
        module = "shell",
        command = "echo 'Task 1'"
    },

    task "task2" {
        module = "shell",
        command = "echo 'Task 2'"
    },

    -- Can include other groups
    include_group("other_group"),

    -- Variables scoped to group
    vars = {
        environment = "production"
    }
}
```

## Advanced Patterns

### Parallel Execution with Goroutines

```lua
task "parallel_deployment" {
    module = "goroutine",
    execute = function()
        local servers = {"web-01", "web-02", "web-03"}
        local wg = goroutine.WaitGroup()

        for _, server in ipairs(servers) do
            wg:Add(1)
            goroutine.Go(function()
                deploy_to_server(server)
                wg:Done()
            end)
        end

        wg:Wait()
        return true, "Deployed to all servers"
    end
}
```

### Dynamic Task Generation

```lua
-- Generate tasks based on configuration
local environments = {"dev", "staging", "prod"}

for _, env in ipairs(environments) do
    task("deploy_" .. env) {
        module = "k8s",
        action = "deploy",
        namespace = env,
        manifest = "./k8s/" .. env .. "/deployment.yaml",
        when = function()
            return os.getenv("DEPLOY_ENV") == env
        end
    }
end
```

### Template Functions

```lua
-- Create reusable task templates
function create_deployment_task(name, config)
    return task(name) {
        module = "k8s",
        action = "deploy",
        namespace = config.namespace,
        image = config.image,
        replicas = config.replicas,
        resources = config.resources
    }
end

-- Use template
group "deployments" {
    create_deployment_task("deploy_frontend", {
        namespace = "production",
        image = "frontend:latest",
        replicas = 3,
        resources = {cpu = "100m", memory = "256Mi"}
    }),

    create_deployment_task("deploy_backend", {
        namespace = "production",
        image = "backend:latest",
        replicas = 5,
        resources = {cpu = "500m", memory = "1Gi"}
    })
}
```

### Error Recovery Pattern

```lua
task "with_recovery" {
    module = "custom",
    execute = function()
        local success, result = pcall(function()
            -- Main operation
            return risky_operation()
        end)

        if not success then
            log.warn("Operation failed, attempting recovery")

            -- Recovery logic
            local recovered = attempt_recovery()
            if recovered then
                return true, "Recovered from failure"
            else
                return false, "Recovery failed: " .. result
            end
        end

        return true, "Operation successful", result
    end
}
```

## Complete Example

```lua
-- production-deployment.lua
name = "production-deployment"
version = "2.0.0"
description = "Production deployment with full validation"

-- Helper functions
local function validate_environment()
    local required_vars = {"API_TOKEN", "DB_PASSWORD", "DEPLOY_KEY"}
    for _, var in ipairs(required_vars) do
        if not os.getenv(var) then
            error("Missing required environment variable: " .. var)
        end
    end
end

-- Pre-flight checks
group "preflight" {
    task "validate_env" {
        module = "custom",
        execute = validate_environment
    }

    task "check_services" {
        module = "infra_test",
        tests = {
            {type = "service_running", name = "docker"},
            {type = "port_reachable", host = "registry.example.com", port = 443}
        }
    }
}

-- Build and test
group "build" {
    depends_on = {"preflight"},

    task "build_image" {
        module = "docker",
        action = "build",
        tag = "myapp:${VERSION}",
        context = ".",
        dockerfile = "Dockerfile.prod"
    }

    task "run_tests" {
        module = "docker",
        action = "run",
        image = "myapp:${VERSION}",
        command = "npm test",
        remove = true
    }

    task "push_image" {
        module = "docker",
        action = "push",
        image = "myapp:${VERSION}",
        registry = "registry.example.com"
    }
}

-- Deploy
group "deploy" {
    depends_on = {"build"},
    parallel = true,

    task "deploy_backend" {
        module = "k8s",
        action = "deploy",
        manifest = "./k8s/backend.yaml",
        namespace = "production",
        wait = true,
        timeout = 600
    }

    task "deploy_frontend" {
        module = "k8s",
        action = "deploy",
        manifest = "./k8s/frontend.yaml",
        namespace = "production",
        wait = true,
        timeout = 600
    }
}

-- Verify deployment
group "verify" {
    depends_on = {"deploy"},

    task "health_checks" {
        module = "infra_test",
        tests = {
            {type = "http_ok", url = "https://api.example.com/health"},
            {type = "http_ok", url = "https://example.com"},
            {type = "http_contains", url = "https://example.com",
             text = "Welcome"}
        }
    }

    task "smoke_tests" {
        module = "custom",
        execute = function()
            return run_smoke_tests()
        end
    }

    on_success = function()
        notification.send("slack", {
            channel = "#deployments",
            message = "Production deployment successful! Version: ${VERSION}"
        })
    end,

    on_failure = function()
        notification.send("pagerduty", {
            severity = "critical",
            message = "Production deployment failed!"
        })
        -- Trigger rollback
        rollback()
    end
}
```

## Summary

This reference guide covers the complete Modern DSL API for Sloth Runner. Key points:

1. **Modules** provide specialized functionality for different tasks
2. **Global functions** are available without imports
3. **Tasks and groups** organize your automation workflow
4. **Error handling** ensures reliability
5. **Parallel execution** improves performance
6. **Templates** enable code reuse

For more examples and patterns, see:
- [Module API Examples](module-api-examples.md)
- [Best Practices](best-practices.md)
- [Introduction](introduction.md)