# 📚 Modern DSL Reference Guide

This guide provides comprehensive examples and patterns for using the **Modern DSL** syntax in Sloth Runner.

## 🎯 Modern DSL Overview

The Modern DSL is the primary syntax for defining workflows in Sloth Runner, providing a powerful fluent API for task orchestration.

### Modern DSL Structure

```lua
-- Define tasks with fluent API
local my_task = task("task_name")
    :description("Task description")
    :command("shell command or function")
    :depends_on({"other_task"})
    :timeout("5m")
    :retries(3, "exponential")
    :build()

-- Define workflow
workflow.define("workflow_name", {
    description = "Workflow description",
    version = "2.0.0",
    tasks = { my_task }
})
```

## 🚀 Essential Modern DSL Patterns

### Basic Task Definition

```lua
local build_task = task("build_app")
    :description("Build the application")
    :command("npm run build")
    :timeout("10m")
    :build()
```

### Task with Dependencies

```lua
local deploy_task = task("deploy_app")
    :command("kubectl apply -f deployment.yaml")
    :depends_on({"build_app"})  -- Always use array format
    :build()

local notify_task = task("notify_team")
    :command("slack-notify.sh")
    :depends_on({"test_app", "security_scan"})  -- Multiple dependencies
    :build()
```

### Function Commands with Error Handling

```lua
local process_task = task("process_data")
    :command(function(params, deps)  -- 'deps' parameter for dependency outputs
        local data = deps.fetch_data.result
        log.info("Processing: " .. data)
        
        if not data then
            return false, "No data to process"
        end
        
        return true, "Processed", {processed = data}
    end)
    :build()
```

### Advanced Task with Hooks

```lua
local deploy_task = task("deploy")
    :command("deploy.sh")
    :pre_hook(function(params, deps)  -- pre_hook for setup
        log.info("Preparing deployment...")
        return true, "Ready"
    end)
    :post_hook(function(params, output)  -- post_hook for cleanup
        log.info("Deployment completed")
        return true, "Done"
    end)
    :on_success(function(params, output)  -- Success-specific hook
        notifications.send("slack", "Deployment successful!")
    end)
    :on_failure(function(params, error)  -- Failure-specific hook
        alerts.send("pagerduty", "Deployment failed: " .. error)
    end)
    :build()
```

## 🔧 Enhanced Modern DSL Features

### Circuit Breaker Pattern

```lua
local api_task = task("call_external_api")
    :command(function()
        return circuit.protect("payment_api", function()
            return net.http_post("https://api.payment.com/charge", data)
        end)
    end)
    :retries(3, "exponential")
    :build()
```

### Conditional Execution

```lua
local enhanced_task = task("enhanced_deploy")
    :description("Deploy with modern features")
    :command("deploy.sh")
    :depends_on({"build", "test"})
    :timeout("15m")
    :retries(3, "exponential")  -- Enhanced retry with strategy
    :condition(when("env.ENVIRONMENT == 'production'"))  -- Conditional execution
    :artifacts({"deployment.yaml", "logs/"})  -- Artifact management
    :metadata({  -- Rich metadata
        owner = "platform-team",
        cost_center = "engineering"
    })
    :build()
```

### Parallel Task Execution

```lua
local parallel_build = task("parallel_build")
    :command(function()
        -- Modern async patterns
        local results = async.parallel({
            frontend = function()
                return exec.run("npm run build:frontend")
            end,
            backend = function() 
                return exec.run("go build ./cmd/server")
            end,
            docs = function()
                return exec.run("mkdocs build")
            end
        }, {
            max_workers = 3,
            timeout = "10m",
            fail_fast = false
        })
        
        return true, "All builds completed", results
    end)
    :build()
```

## 🌟 Complete Workflow Examples

### Simple Workflow

```lua
-- Simple build and test workflow
local build_task = task("build")
    :description("Build application")
    :command("go build -o app ./cmd/main.go")
    :timeout("5m")
    :artifacts({"app"})
    :build()

local test_task = task("test")
    :description("Run tests")
    :command("go test ./...")
    :depends_on({"build"})
    :timeout("10m")
    :build()

workflow.define("ci_pipeline", {
    description = "Simple CI Pipeline",
    version = "2.0.0",
    tasks = { build_task, test_task },
    
    config = {
        timeout = "30m",
        retry_policy = "exponential"
    }
})
```

### Enterprise Workflow with Monitoring

```lua
workflow.define("enterprise_deployment", {
    description = "Enterprise deployment pipeline",
    version = "3.0.0",
    
    metadata = {
        author = "Platform Team",
        team = "infrastructure",
        tags = {"deployment", "production", "enterprise"},
        cost_center = "engineering",
        criticality = "high"
    },
    
    tasks = { build_task, test_task, security_task, deploy_task },
    
    config = {
        timeout = "2h",
        max_parallel_tasks = 4,
        retry_policy = "exponential",
        cleanup_on_failure = true,
        
        monitoring = {
            metrics = true,
            alerts = true,
            dashboard = "grafana://deployment-pipeline"
        },
        
        security = {
            required_secrets = ["k8s_token", "registry_password"],
            rbac_role = "deployment-executor"
        }
    },
    
    pre_conditions = {
        cluster_available = function()
            local result = exec.run("kubectl cluster-info")
            return result.success, "Kubernetes cluster not available"
        end
    },
    
    on_start = function()
        log.info("🚀 Starting enterprise deployment...")
        metrics.increment("deployment_starts_total")
        return true
    end,
    
    on_complete = function(success, results)
        local duration = metrics.stop_timer("deployment_duration")
        
        if success then
            log.info("✅ Deployment completed successfully!")
            metrics.increment("deployment_success_total")
        else
            log.error("❌ Deployment failed!")
            metrics.increment("deployment_failure_total")
        end
        
        return true
    end
})
```

## 🎯 Best Practices

### Task Definition
1. **Always use descriptive names** for tasks and workflows
2. **Set appropriate timeouts** for all tasks
3. **Use exponential backoff** for retry strategies
4. **Add metadata** for tracking and documentation
5. **Implement proper error handling** in function commands

### Workflow Organization
1. **Group related tasks** logically
2. **Use meaningful version numbers** for workflows
3. **Add comprehensive metadata** for maintainability
4. **Set resource limits** for performance
5. **Enable monitoring** for production workflows

### Error Handling
1. **Use circuit breakers** for external dependencies
2. **Implement compensation logic** for critical operations
3. **Add proper logging** at all levels
4. **Set up alerts** for failures
5. **Plan rollback strategies** for deployments

---

**🎯 The Modern DSL provides powerful capabilities for building robust, maintainable workflows. Use these patterns as building blocks for your automation needs!**