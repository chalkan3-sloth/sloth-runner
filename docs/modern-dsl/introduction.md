# 🎨 Modern DSL Introduction

## Overview

The **Modern DSL (Domain Specific Language)** is a powerful, Lua-based configuration language for Sloth Runner that makes task automation intuitive and expressive. It provides a clean, readable syntax for defining complex automation workflows while leveraging the full power of Lua and Sloth Runner's extensive module ecosystem.

## Why Modern DSL?

The Modern DSL offers several advantages over traditional YAML-based configuration:

- **🚀 Dynamic Configuration**: Use variables, loops, and conditionals
- **📦 Module System**: Access powerful built-in and custom modules
- **🔄 Reusability**: Create functions and templates for common patterns
- **🎯 Type Safety**: Leverage Lua's type system for safer configurations
- **⚡ Performance**: Direct execution without YAML parsing overhead
- **🧩 Composability**: Build complex workflows from simple components

## Basic Structure

A Modern DSL file is a Lua script that defines tasks using Sloth Runner's API:

```lua
-- Basic task definition
name = "my-automation"
version = "1.0.0"
description = "My first Modern DSL automation"

-- Define a task group
group "setup" {
    description = "Initial setup tasks",

    task "create_directory" {
        module = "fs",
        action = "mkdir",
        path = "/opt/myapp",
        mode = "0755"
    }
}
```

## Module System

The power of Modern DSL comes from its seamless integration with Sloth Runner's module system:

```lua
-- Using modules directly without require
group "system_setup" {
    task "install_packages" {
        module = "pkg",
        action = "install",
        packages = {"nginx", "postgresql", "redis"},
        state = "present"
    }

    task "configure_service" {
        module = "systemd",
        action = "service",
        name = "nginx",
        state = "started",
        enabled = true
    }
}
```

## Dynamic Capabilities

Unlike static YAML, Modern DSL supports full programming constructs:

```lua
-- Variables and conditionals
local environments = {"dev", "staging", "prod"}
local base_path = "/opt/applications"

for _, env in ipairs(environments) do
    group("setup_" .. env) {
        task("create_env_directory") {
            module = "fs",
            action = "mkdir",
            path = base_path .. "/" .. env,
            mode = "0755"
        }

        if env == "prod" then
            task("setup_monitoring") {
                module = "systemd",
                action = "service",
                name = "prometheus-node-exporter",
                state = "started"
            }
        end
    }
end
```

## Global Modules

Sloth Runner provides global modules that don't require explicit imports:

```lua
-- State management
state.set("app_version", "2.0.1")
local version = state.get("app_version")

-- Logging
log.info("Starting deployment for version: " .. version)

-- Facts (system information)
local os_info = facts.get_os()
if os_info.family == "debian" then
    -- Debian-specific tasks
end

-- HTTP operations
local response = http.get("https://api.github.com/repos/myorg/myrepo/releases/latest")
local latest_release = json.decode(response.body)
```

## Parallel Execution

Execute tasks in parallel for better performance:

```lua
group "parallel_deployment" {
    parallel = true,  -- Enable parallel execution for this group

    task "deploy_web" {
        module = "docker",
        action = "container",
        name = "web-app",
        image = "myapp:latest",
        state = "started"
    }

    task "deploy_api" {
        module = "docker",
        action = "container",
        name = "api-server",
        image = "myapi:latest",
        state = "started"
    }

    task "deploy_worker" {
        module = "docker",
        action = "container",
        name = "background-worker",
        image = "myworker:latest",
        state = "started"
    }
}
```

## Error Handling

Robust error handling with retry logic:

```lua
group "deployment" {
    on_error = "continue",  -- Continue on error
    max_retries = 3,
    retry_delay = 5,  -- seconds

    task "download_artifact" {
        module = "net",
        action = "download",
        url = "https://releases.example.com/app.tar.gz",
        dest = "/tmp/app.tar.gz",

        on_success = function(result)
            log.info("Download completed: " .. result.size .. " bytes")
        end,

        on_failure = function(error)
            log.error("Download failed: " .. error.message)
            -- Send notification
            notification.send("slack", {
                channel = "#alerts",
                message = "Deployment failed: " .. error.message
            })
        end
    }
}
```

## Templates and Functions

Create reusable components:

```lua
-- Define a reusable function for creating services
function create_service(name, port, replicas)
    return task("deploy_" .. name) {
        module = "docker",
        action = "service",
        name = name,
        image = name .. ":latest",
        replicas = replicas or 1,
        ports = {port .. ":" .. port},
        networks = {"app-network"},
        environment = {
            SERVICE_NAME = name,
            SERVICE_PORT = port
        }
    }
end

-- Use the function
group "microservices" {
    create_service("auth-service", 3000, 2),
    create_service("user-service", 3001, 3),
    create_service("payment-service", 3002, 2),
    create_service("notification-service", 3003, 1)
}
```

## Integration with External Data

Load configuration from external sources:

```lua
-- Load configuration from JSON
local config = json.decode(fs.read("/etc/myapp/config.json"))

-- Load secrets from environment
local db_password = os.getenv("DB_PASSWORD") or vault.get("database/password")

group "database_setup" {
    task "configure_database" {
        module = "postgresql",
        action = "database",
        name = config.database.name,
        owner = config.database.user,
        encoding = "UTF-8"
    }

    task "create_user" {
        module = "postgresql",
        action = "user",
        name = config.database.user,
        password = db_password,
        privileges = ["CREATEDB", "CREATEROLE"]
    }
}
```

## Best Practices

1. **Use descriptive names**: Make your tasks and groups self-documenting
2. **Leverage modules**: Don't reinvent the wheel - use existing modules
3. **Handle errors gracefully**: Always consider failure scenarios
4. **Keep it DRY**: Use functions and templates for repeated patterns
5. **Document complex logic**: Add comments for non-obvious implementations
6. **Test incrementally**: Use `--dry-run` flag to preview changes
7. **Version control**: Track your DSL files in Git

## Next Steps

- [📚 Module API Examples](module-api-examples.md) - Comprehensive module usage examples
- [🎯 Best Practices](best-practices.md) - Advanced patterns and techniques
- [📖 Reference Guide](reference-guide.md) - Complete API reference

## Quick Example: Complete Infrastructure Setup

```lua
name = "infrastructure-setup"
version = "2.0.0"
description = "Complete infrastructure automation with Modern DSL"

-- Global configuration
local app_name = "myapp"
local domain = "example.com"
local environments = {"dev", "staging", "prod"}

-- Helper function for environment-specific configs
function get_replicas(env)
    local replicas_map = {
        dev = 1,
        staging = 2,
        prod = 5
    }
    return replicas_map[env] or 1
end

-- Setup for each environment
for _, env in ipairs(environments) do
    group("setup_" .. env) {
        description = "Setup " .. env .. " environment",

        -- Create namespace
        task("create_namespace") {
            module = "k8s",
            action = "namespace",
            name = app_name .. "-" .. env,
            state = "present"
        },

        -- Deploy application
        task("deploy_app") {
            module = "k8s",
            action = "deployment",
            name = app_name,
            namespace = app_name .. "-" .. env,
            replicas = get_replicas(env),
            image = app_name .. ":" .. state.get("version", "latest"),

            -- Environment-specific configuration
            environment = {
                APP_ENV = env,
                APP_NAME = app_name,
                DB_HOST = env .. "-db." .. domain,
                CACHE_HOST = env .. "-redis." .. domain
            },

            -- Health checks
            readiness_probe = {
                http_get = {
                    path = "/health",
                    port = 8080
                },
                initial_delay_seconds = 10,
                period_seconds = 5
            }
        },

        -- Configure monitoring (production only)
        when = (env == "prod"),
        task("setup_monitoring") {
            module = "prometheus",
            action = "scrape_config",
            job_name = app_name .. "_metrics",
            targets = {app_name .. "-" .. env .. "." .. domain .. ":9090"}
        }
    }
end

-- Post-deployment validation
group "validation" {
    task "check_deployments" {
        module = "http",
        action = "check",
        urls = {
            "https://dev." .. domain .. "/health",
            "https://staging." .. domain .. "/health",
            "https://prod." .. domain .. "/health"
        },
        expected_status = 200,
        timeout = 30
    }
}
```

This introduction provides a foundation for understanding and using the Modern DSL. Continue to the [Module API Examples](module-api-examples.md) for detailed examples of working with specific modules.