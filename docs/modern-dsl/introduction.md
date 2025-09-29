# 🎯 Modern DSL Introduction

Welcome to the **Modern DSL** (Domain Specific Language) for Sloth Runner - a revolutionary approach to defining workflows that combines the power of Lua with an intuitive, fluent API.

## 🚀 What is Modern DSL?

The Modern DSL is a new syntax layer built on top of Sloth Runner that provides:

- **🎯 Fluent API**: Chainable, intuitive method calls
- **📋 Declarative Workflows**: Configuration-driven workflow definitions  
- **🔄 Enhanced Features**: Built-in retry strategies, circuit breakers, and resilience patterns
- **🛡️ Type Safety**: Better validation and error detection
- **📊 Rich Metadata**: Comprehensive workflow and task information
- **⚡ Modern Patterns**: Async operations, performance monitoring, and observability

## 🎨 Syntax Comparison

### Legacy Format (Still Supported)
```lua
Modern DSLs = {
    my_pipeline = {
        description = "Traditional pipeline",
        tasks = {
            {
                name = "build_app",
                command = "go build -o app ./cmd/main.go",
                timeout = "5m",
                retries = 3,
                depends_on = "setup"
            },
            {
                name = "run_tests", 
                command = "go test ./...",
                depends_on = "build_app",
                timeout = "10m"
            }
        }
    }
}
```

### Modern DSL (New Approach)
```lua
-- Define tasks with fluent API
local build_task = task("build_app")
    :description("Build application with modern features")
    :command(function(params, deps)
        log.info("Building application...")
        local result = exec.run("go build -o app ./cmd/main.go")
        
        if not result.success then
            return false, "Build failed: " .. result.stderr
        end
        
        return true, "Build completed", {
            artifact = "app",
            size = fs.size("app"),
            build_time = result.duration
        }
    end)
    :timeout("5m")
    :retries(3, "exponential")
    :depends_on({"setup"})
    :artifacts({"app"})
    :on_success(function(params, output)
        log.info("Build completed! Artifact size: " .. output.size .. " bytes")
    end)
    :build()

local test_task = task("run_tests")
    :description("Run comprehensive test suite")
    :command("go test ./...")
    :depends_on({"build_app"})
    :timeout("10m")
    :condition(when("params.skip_tests != true"))
    :build()

-- Define workflow with rich configuration
workflow.define("my_pipeline", {
    description = "Modern CI/CD Pipeline",
    version = "2.0.0",
    
    metadata = {
        author = "DevOps Team",
        tags = {"ci", "golang", "build"},
        created_at = os.date(),
        repository = "github.com/company/project"
    },
    
    tasks = { build_task, test_task },
    
    config = {
        timeout = "30m",
        retry_policy = "exponential",
        max_parallel_tasks = 4,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting modern CI/CD pipeline...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Pipeline completed successfully!")
            -- Send notification, update status, etc.
        else
            log.error("❌ Pipeline failed!")
        end
        return true
    end
})
```

## 🎯 Key Benefits

### 1. **Enhanced Readability**
The fluent API makes workflows self-documenting and easier to understand:

```lua
-- Clear, expressive syntax
local deploy_task = task("deploy_to_production")
    :description("Deploy application to production environment")
    :command(function(params, deps)
        -- Business logic is clear and well-structured
        return deploy_application(deps.build_app.artifact)
    end)
    :condition(when("env.ENVIRONMENT == 'production'"))
    :retries(2, "exponential")
    :timeout("15m")
    :on_failure(function(params, error)
        alert.send("deployment_failed", {
            environment = "production",
            error = error
        })
    end)
    :build()
```

### 2. **Built-in Resilience Patterns**
Modern DSL includes enterprise-grade resilience patterns out of the box:

```lua
-- Circuit breaker for external dependencies
local api_task = task("call_external_api")
    :command(function()
        return circuit.protect("payment_api", function()
            return net.http_post("https://api.payment.com/charge", data)
        end)
    end)
    :retries(3, "exponential")
    :build()

-- Saga pattern for distributed transactions
local payment_saga = saga.define("payment_process")
    :step("validate_payment", validate_task)
    :step("charge_card", charge_task)
    :step("update_inventory", inventory_task)
    :compensate("charge_card", refund_task)
    :compensate("update_inventory", restore_inventory_task)
    :build()
```

### 3. **Advanced Async Operations**
Modern patterns for parallel and asynchronous execution:

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

### 4. **Rich Metadata and Observability**
Comprehensive tracking and monitoring capabilities:

```lua
workflow.define("data_pipeline", {
    description = "ETL Data Processing Pipeline",
    version = "3.1.0",
    
    metadata = {
        author = "Data Team",
        tags = {"etl", "data", "analytics"},
        sla = "4h",
        cost_center = "analytics",
        compliance = ["GDPR", "SOX"]
    },
    
    config = {
        monitoring = {
            metrics = true,
            alerts = true,
            dashboard = "grafana://data-pipeline"
        },
        performance = {
            expected_duration = "2h",
            memory_limit = "4GB",
            cpu_limit = "2 cores"
        }
    }
})
```

## 🔄 Migration Path

The Modern DSL provides a smooth migration path:

### 1. **100% Backward Compatibility**
Your existing scripts continue to work without any changes:

```lua
-- Your existing Modern DSLs still work perfectly
Modern DSLs = {
    legacy_pipeline = {
        description = "This still works!",
        tasks = { /* your existing tasks */ }
    }
}
```

### 2. **Gradual Adoption**
You can mix legacy and modern syntax in the same file:

```lua
-- New tasks using Modern DSL
local modern_task = task("modern_deploy")
    :description("Deploy with modern patterns")
    :command(function() return deploy_with_retry() end)
    :build()

-- Legacy tasks still work
Modern DSLs = {
    mixed_workflow = {
        description = "Mix of old and new",
        tasks = {
            {
                name = "legacy_task",
                command = "echo 'Still works!'"
            }
        }
    }
}
```

### 3. **Migration Tools**
Use built-in tools to convert existing workflows:

```bash
# Convert legacy to modern DSL
./sloth-runner migrate -f legacy-workflow.lua -o modern-workflow.lua

# Validate modern DSL syntax
./sloth-runner validate -f modern-workflow.lua --dsl-version 2.0
```

## 🎓 Learning Path

### Beginner
1. Start with simple task definitions
2. Learn the fluent API basics
3. Explore basic workflow configuration

### Intermediate  
4. Add error handling and retries
5. Use conditional execution
6. Implement parallel tasks

### Advanced
7. Master circuit breaker patterns
8. Implement saga patterns
9. Build enterprise-grade pipelines

## 🚀 Getting Started

Ready to start with Modern DSL? Check out these resources:

- [Task Definition API](./task-api.md) - Complete task builder reference
- [Workflow Definition](./workflow-api.md) - Workflow configuration guide
- [Migration Guide](./migration-guide.md) - Convert existing workflows
- [Best Practices](./best-practices.md) - Modern DSL patterns and guidelines
- [Examples](../../examples/) - Browse modernized examples

## 🤝 Community

The Modern DSL is designed with community feedback in mind:

- **🐛 Issues**: Report bugs and request features
- **💡 Ideas**: Propose new DSL features
- **📚 Documentation**: Help improve guides and examples
- **🔧 Tools**: Build migration and validation tools

---

**🎯 The Modern DSL represents the future of workflow automation - more powerful, more intuitive, and more maintainable than ever before!**