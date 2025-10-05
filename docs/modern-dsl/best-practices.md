# 🎯 Modern DSL Best Practices

This guide provides comprehensive best practices for writing efficient, maintainable, and robust automation workflows using Sloth Runner's Modern DSL.

## 📋 Table of Contents
- [Code Organization](#code-organization)
- [Module Usage Patterns](#module-usage-patterns)
- [Error Handling](#error-handling)
- [Performance Optimization](#performance-optimization)
- [Security Best Practices](#security-best-practices)
- [Testing and Validation](#testing-and-validation)
- [Documentation Standards](#documentation-standards)

## Code Organization

### Structure Your Files

```lua
-- ✅ Good: Organized structure
-- file: deployments/production.lua

-- Configuration at the top
local config = require("config/production")
local utils = require("lib/utils")

-- Constants
local APP_NAME = "myapp"
local ENVIRONMENT = "production"

-- Helper functions
local function validate_deployment(params)
    -- validation logic
end

local function get_replica_count(env)
    local replicas = {
        dev = 1,
        staging = 2,
        production = 5
    }
    return replicas[env] or 1
end

-- Main workflow definition
name = "production-deployment"
version = "1.0.0"
description = "Production deployment workflow"

-- Task groups
group "preparation" {
    -- tasks
}

group "deployment" {
    -- tasks
}

group "verification" {
    -- tasks
}
```

### Modularize Common Patterns

```lua
-- file: lib/deployment_helpers.lua
local M = {}

-- Reusable function for creating deployment tasks
function M.create_deployment_task(environment, config)
    return task("deploy_to_" .. environment) {
        module = "k8s",
        action = "deploy",
        namespace = environment,
        image = config.image,
        replicas = config.replicas[environment],
        resources = config.resources[environment],

        on_success = function(result)
            state.set("last_deployment." .. environment, {
                version = config.version,
                timestamp = os.time(),
                pods = result.ready_pods
            })
        end
    }
end

-- Reusable health check
function M.health_check(service, environment)
    return task("health_check_" .. service) {
        module = "infra_test",
        action = "http_check",
        url = "https://" .. environment .. ".example.com/health",
        expected_status = 200,
        timeout = 30
    }
end

return M

-- Usage in main file:
local deploy_helpers = require("lib/deployment_helpers")

group "deployment" {
    deploy_helpers.create_deployment_task("production", config),
    deploy_helpers.health_check("api", "production")
}
```

## Module Usage Patterns

### Package Management Best Practices

```lua
-- ✅ Good: Idempotent package installation with verification
group "system_setup" {
    task "install_dependencies" {
        module = "pkg",
        action = "install",
        packages = function()
            -- Dynamically determine packages based on OS
            local os_family = facts.get_os().family
            local base_packages = {"curl", "git", "vim"}

            if os_family == "debian" then
                table.insert(base_packages, "apt-transport-https")
            elseif os_family == "redhat" then
                table.insert(base_packages, "yum-utils")
            end

            return base_packages
        end,
        state = "present",

        -- Verify installation
        post_install = function()
            for _, pkg in ipairs(packages) do
                if not pkg.is_installed(pkg) then
                    return false, pkg .. " failed to install"
                end
            end
            return true
        end
    }
}
```

### Service Management Patterns

```lua
-- ✅ Good: Graceful service management with checks
group "service_management" {
    task "restart_service_safely" {
        module = "systemd",
        action = "custom",

        execute = function(params)
            local service_name = params.service

            -- Check if service exists
            if not systemd.service_exists({service = service_name}) then
                return false, "Service " .. service_name .. " does not exist"
            end

            -- Check current status
            local status = systemd.status({service = service_name})

            -- Graceful restart with health check
            if status.active then
                -- Reload if config changed
                if params.config_changed then
                    systemd.reload({service = service_name})
                    goroutine.sleep(2000)
                end

                -- Graceful restart
                systemd.restart({service = service_name})
            else
                -- Start if not running
                systemd.start({service = service_name})
            end

            -- Wait and verify
            goroutine.sleep(5000)

            local new_status = systemd.status({service = service_name})
            if not new_status.active then
                return false, "Service failed to start after restart"
            end

            return true, "Service " .. service_name .. " restarted successfully"
        end
    }
}
```

### Docker Module Patterns

```lua
-- ✅ Good: Container lifecycle management
group "container_deployment" {
    task "deploy_with_zero_downtime" {
        module = "docker",
        action = "custom",

        execute = function(params)
            local container_name = params.container
            local image = params.image

            -- Pull new image first
            docker.pull({image = image})

            -- Check if container exists
            local existing = docker.inspect({container = container_name})

            if existing then
                -- Create new container with temporary name
                local temp_name = container_name .. "_new"
                docker.run({
                    name = temp_name,
                    image = image,
                    ports = params.ports,
                    environment = params.environment,
                    detach = true
                })

                -- Health check new container
                goroutine.sleep(5000)
                local health = docker.exec({
                    container = temp_name,
                    command = "curl -f http://localhost/health"
                })

                if health.exit_code == 0 then
                    -- Stop old container
                    docker.stop({container = container_name})
                    docker.remove({container = container_name})

                    -- Rename new container
                    docker.rename({
                        from = temp_name,
                        to = container_name
                    })
                else
                    -- Rollback
                    docker.stop({container = temp_name})
                    docker.remove({container = temp_name})
                    return false, "Health check failed on new container"
                end
            else
                -- Fresh deployment
                docker.run({
                    name = container_name,
                    image = image,
                    ports = params.ports,
                    environment = params.environment,
                    restart = "always",
                    detach = true
                })
            end

            return true, "Container deployed successfully"
        end
    }
}
```

## Error Handling

### Comprehensive Error Handling

```lua
-- ✅ Good: Multiple levels of error handling
group "robust_deployment" {
    on_error = "continue",  -- Continue with other tasks on error

    task "deploy_application" {
        module = "custom",
        action = "deploy",

        execute = function(params)
            -- Wrap in pcall for unexpected errors
            local success, result = pcall(function()
                -- Validate inputs
                if not params.version then
                    error("Version is required")
                end

                -- Try deployment with timeout
                local deploy_result = with_timeout(30, function()
                    return deploy_app(params)
                end)

                if not deploy_result.success then
                    -- Specific error handling
                    if deploy_result.error:match("timeout") then
                        error("Deployment timed out after 30 seconds")
                    elseif deploy_result.error:match("auth") then
                        error("Authentication failed - check credentials")
                    else
                        error("Deployment failed: " .. deploy_result.error)
                    end
                end

                return deploy_result
            end)

            if not success then
                -- Log error with context
                log.error("Deployment failed", {
                    version = params.version,
                    environment = params.environment,
                    error = result
                })

                -- Send alert
                notification.send("slack", {
                    channel = "#ops-alerts",
                    message = "Deployment failed: " .. result
                })

                return false, result
            end

            return true, "Deployment successful", result
        end,

        -- Retry with exponential backoff
        retry_policy = {
            max_attempts = 3,
            backoff = "exponential",
            initial_delay = 5,
            max_delay = 60,

            should_retry = function(error)
                -- Only retry on specific errors
                return error:match("timeout") or error:match("connection")
            end
        }
    }
}
```

### Circuit Breaker Pattern

```lua
-- ✅ Good: Circuit breaker for external services
local circuit_breakers = {}

function create_circuit_breaker(name, threshold, timeout)
    circuit_breakers[name] = {
        failures = 0,
        threshold = threshold or 5,
        timeout = timeout or 60,
        open_until = 0,
        state = "closed"  -- closed, open, half-open
    }
end

function call_with_circuit_breaker(name, fn)
    local cb = circuit_breakers[name]
    if not cb then
        create_circuit_breaker(name)
        cb = circuit_breakers[name]
    end

    -- Check if circuit is open
    if cb.state == "open" then
        if os.time() < cb.open_until then
            return false, "Circuit breaker is open"
        else
            -- Try half-open
            cb.state = "half-open"
        end
    end

    -- Execute function
    local success, result = pcall(fn)

    if success then
        -- Reset on success
        if cb.state == "half-open" then
            cb.state = "closed"
            cb.failures = 0
        end
        return true, result
    else
        -- Record failure
        cb.failures = cb.failures + 1

        if cb.failures >= cb.threshold then
            cb.state = "open"
            cb.open_until = os.time() + cb.timeout
            log.warn("Circuit breaker opened for: " .. name)
        end

        return false, result
    end
end

-- Usage
task "call_external_api" {
    module = "custom",
    execute = function(params)
        return call_with_circuit_breaker("payment_api", function()
            return http.post("https://api.payment.com/charge", {
                amount = params.amount
            })
        end)
    end
}
```

## Performance Optimization

### Parallel Execution

```lua
-- ✅ Good: Efficient parallel execution with worker pool
group "parallel_processing" {
    task "process_large_dataset" {
        module = "goroutine",
        action = "custom",

        execute = function(params)
            local files = fs.list(params.input_dir)
            local results = {}
            local errors = {}

            -- Create worker pool
            local worker_pool = goroutine.WorkerPool({
                size = params.workers or 4,
                queue_size = 100
            })

            -- Process files in parallel
            for _, file in ipairs(files) do
                worker_pool:submit(function()
                    local success, result = pcall(function()
                        return process_file(file)
                    end)

                    if success then
                        results[file] = result
                    else
                        errors[file] = result
                    end
                end)
            end

            -- Wait for completion
            worker_pool:wait()
            worker_pool:shutdown()

            -- Check results
            if #errors > 0 then
                local error_rate = #errors / #files
                if error_rate > 0.1 then  -- More than 10% errors
                    return false, "Too many errors: " .. #errors .. "/" .. #files
                end
            end

            return true, "Processing complete", {
                total = #files,
                successful = #results,
                failed = #errors
            }
        end
    }
}
```

### Resource Management

```lua
-- ✅ Good: Efficient resource usage
group "resource_efficient" {
    task "process_stream" {
        module = "custom",

        execute = function(params)
            -- Use streaming for large files
            local input = io.open(params.input_file, "r")
            local output = io.open(params.output_file, "w")

            local line_count = 0
            local batch = {}
            local batch_size = 1000

            -- Process in batches
            for line in input:lines() do
                table.insert(batch, line)
                line_count = line_count + 1

                if #batch >= batch_size then
                    -- Process batch
                    local processed = process_batch(batch)
                    for _, item in ipairs(processed) do
                        output:write(item .. "\n")
                    end

                    -- Clear batch
                    batch = {}

                    -- Yield periodically
                    if line_count % 10000 == 0 then
                        coroutine.yield()
                        log.info("Processed " .. line_count .. " lines")
                    end
                end
            end

            -- Process remaining
            if #batch > 0 then
                local processed = process_batch(batch)
                for _, item in ipairs(processed) do
                    output:write(item .. "\n")
                end
            end

            input:close()
            output:close()

            return true, "Processed " .. line_count .. " lines"
        end
    }
}
```

## Security Best Practices

### Secret Management

```lua
-- ✅ Good: Secure secret handling
group "secure_deployment" {
    task "deploy_with_secrets" {
        module = "custom",

        execute = function(params)
            -- Never hardcode secrets
            -- ❌ Bad: local password = "supersecret"

            -- ✅ Good: Get from secure source
            local db_password = os.getenv("DB_PASSWORD") or
                               vault.get("secrets/database/password")

            if not db_password then
                return false, "Database password not available"
            end

            -- Use secret without logging
            local conn = database.connect({
                host = params.db_host,
                username = params.db_user,
                password = db_password,  -- Never log this
                database = params.db_name
            })

            -- Clear from memory after use
            db_password = nil
            collectgarbage()

            return true, "Connected to database"
        end
    }
}
```

### Input Validation

```lua
-- ✅ Good: Comprehensive input validation
function validate_input(params)
    local validators = {
        -- Required fields
        required = {"username", "email", "action"},

        -- Type validation
        types = {
            username = "string",
            email = "string",
            age = "number",
            active = "boolean"
        },

        -- Pattern validation
        patterns = {
            email = "^[%w._%+-]+@[%w.-]+%.[%w]+$",
            username = "^[a-zA-Z0-9_]+$"
        },

        -- Range validation
        ranges = {
            age = {min = 18, max = 120}
        }
    }

    -- Check required fields
    for _, field in ipairs(validators.required) do
        if not params[field] then
            return false, "Missing required field: " .. field
        end
    end

    -- Check types
    for field, expected_type in pairs(validators.types) do
        if params[field] and type(params[field]) ~= expected_type then
            return false, field .. " must be " .. expected_type
        end
    end

    -- Check patterns
    for field, pattern in pairs(validators.patterns) do
        if params[field] and not params[field]:match(pattern) then
            return false, field .. " has invalid format"
        end
    end

    -- Check ranges
    for field, range in pairs(validators.ranges) do
        if params[field] then
            if params[field] < range.min or params[field] > range.max then
                return false, field .. " out of range"
            end
        end
    end

    return true
end

-- Usage
task "process_user_request" {
    module = "custom",

    execute = function(params)
        local valid, error = validate_input(params)
        if not valid then
            return false, "Validation failed: " .. error
        end

        -- Process validated input
        return process_request(params)
    end
}
```

## Testing and Validation

### Infrastructure Testing

```lua
-- ✅ Good: Comprehensive infrastructure validation
group "deployment_validation" {
    task "validate_deployment" {
        module = "infra_test",
        action = "suite",

        tests = {
            -- Service checks
            {
                name = "nginx_running",
                type = "service_running",
                service = "nginx",
                critical = true
            },

            -- Port checks
            {
                name = "http_port_open",
                type = "port_listening",
                port = 80,
                protocol = "tcp",
                critical = true
            },

            -- File checks
            {
                name = "config_exists",
                type = "file_exists",
                path = "/etc/nginx/nginx.conf",
                critical = false
            },

            -- HTTP checks
            {
                name = "health_endpoint",
                type = "http_check",
                url = "http://localhost/health",
                expected_status = 200,
                timeout = 10,
                retries = 3
            },

            -- Custom validation
            {
                name = "custom_check",
                type = "custom",
                validator = function()
                    local result = exec.run("nginx -t")
                    return result.exit_code == 0
                end
            }
        },

        on_failure = function(test, error)
            log.error("Validation failed", {
                test = test.name,
                type = test.type,
                error = error
            })

            if test.critical then
                -- Rollback if critical test fails
                rollback()
            end
        end
    }
}
```

### Dry Run Support

```lua
-- ✅ Good: Support for dry-run mode
group "deployment" {
    task "deploy_with_dry_run" {
        module = "custom",

        execute = function(params)
            local dry_run = params.dry_run or false

            if dry_run then
                log.info("[DRY-RUN] Would deploy version: " .. params.version)
                log.info("[DRY-RUN] Would update servers: " .. table.concat(params.servers, ", "))
                log.info("[DRY-RUN] Would use configuration: " .. params.config_file)

                -- Simulate validation
                local validation_result = validate_deployment_params(params)
                if not validation_result.valid then
                    return false, "[DRY-RUN] Validation failed: " .. validation_result.error
                end

                return true, "[DRY-RUN] Deployment would succeed"
            end

            -- Actual deployment
            return deploy_application(params)
        end
    }
}
```

## Documentation Standards

### Task Documentation

```lua
-- ✅ Good: Well-documented task
task "complex_deployment" {
    -- Clear description
    description = "Deploy application with blue-green strategy",

    -- Document parameters
    parameters = {
        version = {
            type = "string",
            required = true,
            description = "Version to deploy (e.g., v1.2.3)"
        },
        environment = {
            type = "string",
            required = true,
            enum = {"dev", "staging", "production"},
            description = "Target environment"
        },
        strategy = {
            type = "string",
            default = "blue-green",
            enum = {"blue-green", "canary", "rolling"},
            description = "Deployment strategy"
        }
    },

    -- Document outputs
    outputs = {
        deployment_id = "Unique deployment identifier",
        active_color = "Current active color (blue/green)",
        endpoint = "Application endpoint URL"
    },

    -- Examples
    examples = {
        {
            description = "Deploy to production",
            params = {
                version = "v1.2.3",
                environment = "production",
                strategy = "blue-green"
            }
        }
    },

    -- Actual implementation
    module = "custom",
    execute = function(params)
        -- Implementation
    end
}
```

### Workflow Documentation

```lua
-- ✅ Good: Complete workflow documentation
--[[
Workflow: Production Deployment Pipeline
Author: Platform Team
Version: 2.0.0
Last Updated: 2024-01-15

Description:
  This workflow handles the complete production deployment process including:
  - Building and testing the application
  - Deploying to staging for validation
  - Blue-green deployment to production
  - Automated rollback on failure

Prerequisites:
  - Docker installed and running
  - Kubernetes cluster access configured
  - Vault credentials available
  - Slack webhook configured

Usage:
  sloth-runner run production-deployment.lua --version=v1.2.3

Environment Variables:
  - DOCKER_REGISTRY: Docker registry URL
  - VAULT_TOKEN: Vault authentication token
  - SLACK_WEBHOOK: Slack webhook for notifications
  - KUBE_CONFIG: Path to Kubernetes config

Outputs:
  - Deployment report in ./reports/deployment-{timestamp}.json
  - Metrics exported to Prometheus
  - Notifications sent to #deployments Slack channel
--]]

name = "production-deployment"
version = "2.0.0"
-- ... rest of workflow
```

## Summary

Following these best practices will help you create:
- **Maintainable** workflows that are easy to understand and modify
- **Reliable** automation that handles errors gracefully
- **Efficient** processes that make optimal use of resources
- **Secure** implementations that protect sensitive data
- **Well-documented** code that others can understand and use

Remember to:
1. Start simple and iterate
2. Test thoroughly in non-production environments
3. Document your workflows and decisions
4. Monitor and measure performance
5. Keep security in mind at all times
6. Use modules effectively - they're tested and optimized
7. Contribute improvements back to the community

For more information, see:
- [Module API Examples](module-api-examples.md)
- [Reference Guide](reference-guide.md)
- [Introduction](introduction.md)