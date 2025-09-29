# 📚 Modern DSL Migration Guide

This guide helps you migrate from the legacy `Modern DSLs` format to the new Modern DSL syntax.

## 🎯 Migration Strategy

### Phase 1: Understanding the Differences

#### Legacy Format Structure
```lua
Modern DSLs = {
    workflow_name = {
        description = "Workflow description",
        tasks = {
            {
                name = "task_name",
                description = "Task description",
                command = "shell command or function",
                depends_on = "other_task",
                timeout = "5m",
                retries = 3
            }
        }
    }
}
```

#### Modern DSL Structure
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

### Phase 2: Automated Migration

Use the built-in migration tool:

```bash
# Migrate single file
./sloth-runner migrate -f legacy-workflow.lua -o modern-workflow.lua

# Migrate all files in directory
./sloth-runner migrate -d examples/ -o modern-examples/ --format modern-dsl

# Dry run to preview changes
./sloth-runner migrate -f workflow.lua --dry-run
```

### Phase 3: Manual Migration Steps

## 🔄 Step-by-Step Migration

### Step 1: Basic Task Conversion

**Before (Legacy):**
```lua
{
    name = "build_app",
    description = "Build the application",
    command = "npm run build",
    timeout = "10m"
}
```

**After (Modern DSL):**
```lua
local build_task = task("build_app")
    :description("Build the application")
    :command("npm run build")
    :timeout("10m")
    :build()
```

### Step 2: Dependencies

**Before (Legacy):**
```lua
{
    name = "deploy_app",
    command = "kubectl apply -f deployment.yaml",
    depends_on = "build_app"  -- Single dependency
}

{
    name = "notify_team",
    command = "slack-notify.sh",
    depends_on = {"test_app", "security_scan"}  -- Multiple dependencies
}
```

**After (Modern DSL):**
```lua
local deploy_task = task("deploy_app")
    :command("kubectl apply -f deployment.yaml")
    :depends_on({"build_app"})  -- Always use array
    :build()

local notify_task = task("notify_team")
    :command("slack-notify.sh")
    :depends_on({"test_app", "security_scan"})
    :build()
```

### Step 3: Function Commands

**Before (Legacy):**
```lua
{
    name = "process_data",
    command = function(params, input_from_dependency)
        local data = input_from_dependency.fetch_data.result
        log.info("Processing: " .. data)
        return true, "Processed", {processed = data}
    end
}
```

**After (Modern DSL):**
```lua
local process_task = task("process_data")
    :command(function(params, deps)  -- 'deps' instead of 'input_from_dependency'
        local data = deps.fetch_data.result
        log.info("Processing: " .. data)
        return true, "Processed", {processed = data}
    end)
    :build()
```

### Step 4: Hooks Migration

**Before (Legacy):**
```lua
{
    name = "deploy",
    command = "deploy.sh",
    pre_exec = function(params, deps)
        log.info("Preparing deployment...")
        return true, "Ready"
    end,
    post_exec = function(params, output)
        log.info("Deployment completed")
        return true, "Done"
    end
}
```

**After (Modern DSL):**
```lua
local deploy_task = task("deploy")
    :command("deploy.sh")
    :pre_hook(function(params, deps)  -- pre_exec becomes pre_hook
        log.info("Preparing deployment...")
        return true, "Ready"
    end)
    :post_hook(function(params, output)  -- post_exec becomes post_hook
        log.info("Deployment completed")
        return true, "Done"
    end)
    :build()
```

### Step 5: Enhanced Features

Add modern DSL exclusive features:

```lua
local enhanced_task = task("enhanced_deploy")
    :description("Deploy with modern features")
    :command("deploy.sh")
    :depends_on({"build", "test"})
    :timeout("15m")
    :retries(3, "exponential")  -- Enhanced retry with strategy
    :condition(when("env.ENVIRONMENT == 'production'"))  -- Conditional execution
    :artifacts({"deployment.yaml", "logs/"})  -- Artifact management
    :on_success(function(params, output)  -- Success-specific hook
        notifications.send("slack", "Deployment successful!")
    end)
    :on_failure(function(params, error)  -- Failure-specific hook
        alerts.send("pagerduty", "Deployment failed: " .. error)
    end)
    :metadata({  -- Rich metadata
        owner = "platform-team",
        cost_center = "engineering"
    })
    :build()
```

### Step 6: Workflow Definition

**Before (Legacy):**
```lua
Modern DSLs = {
    ci_pipeline = {
        description = "CI/CD Pipeline",
        tasks = { task1, task2, task3 }
    }
}
```

**After (Modern DSL):**
```lua
workflow.define("ci_pipeline", {
    description = "CI/CD Pipeline - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "DevOps Team",
        tags = {"ci", "deployment"},
        created_at = os.date()
    },
    
    tasks = { task1, task2, task3 },
    
    config = {
        timeout = "2h",
        retry_policy = "exponential",
        max_parallel_tasks = 4
    },
    
    on_start = function()
        log.info("Starting CI/CD pipeline...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("Pipeline completed successfully!")
        end
        return true
    end
})
```

## 🔧 Migration Patterns

### Pattern 1: Simple Task Migration

```lua
-- Legacy to Modern DSL converter function
function migrate_simple_task(legacy_task)
    local modern_task = task(legacy_task.name)
        :description(legacy_task.description or "Migrated task")
        :command(legacy_task.command)
    
    if legacy_task.timeout then
        modern_task = modern_task:timeout(legacy_task.timeout)
    end
    
    if legacy_task.retries then
        modern_task = modern_task:retries(legacy_task.retries, "exponential")
    end
    
    if legacy_task.depends_on then
        local deps = type(legacy_task.depends_on) == "table" 
            and legacy_task.depends_on 
            or {legacy_task.depends_on}
        modern_task = modern_task:depends_on(deps)
    end
    
    return modern_task:build()
end
```

### Pattern 2: Workflow Migration

```lua
-- Migrate complete workflow
function migrate_workflow(workflow_name, legacy_def)
    local modern_tasks = {}
    
    for _, legacy_task in ipairs(legacy_def.tasks) do
        table.insert(modern_tasks, migrate_simple_task(legacy_task))
    end
    
    workflow.define(workflow_name, {
        description = legacy_def.description .. " - Migrated to Modern DSL",
        version = "2.0.0",
        tasks = modern_tasks,
        
        -- Add modern features
        config = {
            timeout = "2h",
            retry_policy = "exponential"
        }
    })
end
```

## 🚀 Advanced Migration Examples

### Complex Workflow Migration

**Before (Legacy):**
```lua
Modern DSLs = {
    microservice_deployment = {
        description = "Deploy microservice to Kubernetes",
        tasks = {
            {
                name = "build_image",
                description = "Build Docker image",
                command = function(params)
                    local result = exec.run("docker build -t myapp:latest .")
                    if not result.success then
                        return false, "Build failed"
                    end
                    return true, "Build completed", {
                        image_tag = "myapp:latest",
                        image_id = result.stdout:match("sha256:([a-f0-9]+)")
                    }
                end,
                timeout = "10m",
                retries = 2
            },
            {
                name = "run_tests",
                description = "Run unit tests",
                command = "npm test",
                depends_on = "build_image",
                timeout = "5m"
            },
            {
                name = "security_scan",
                description = "Security vulnerability scan",
                command = "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy myapp:latest",
                depends_on = "build_image",
                timeout = "3m"
            },
            {
                name = "deploy_to_k8s",
                description = "Deploy to Kubernetes",
                command = function(params, deps)
                    local image_info = deps.build_image
                    log.info("Deploying image: " .. image_info.image_tag)
                    
                    local deploy_result = exec.run("kubectl apply -f k8s/deployment.yaml")
                    if not deploy_result.success then
                        return false, "Deployment failed"
                    end
                    
                    return true, "Deployed successfully", {
                        deployment_time = os.time(),
                        namespace = "production"
                    }
                end,
                depends_on = {"run_tests", "security_scan"},
                timeout = "15m",
                retries = 1
            }
        }
    }
}
```

**After (Modern DSL):**
```lua
-- Enhanced build task with modern features
local build_task = task("build_image")
    :description("Build Docker image with enhanced error handling")
    :command(function(params)
        log.info("Starting Docker build process...")
        
        local result = exec.run("docker build -t myapp:latest .", {
            timeout = "8m",
            env = {
                DOCKER_BUILDKIT = "1"
            }
        })
        
        if not result.success then
            return false, "Build failed: " .. result.stderr
        end
        
        local image_id = result.stdout:match("sha256:([a-f0-9]+)")
        
        return true, "Build completed", {
            image_tag = "myapp:latest",
            image_id = image_id,
            build_duration = result.duration,
            build_size = fs.size("Dockerfile")
        }
    end)
    :timeout("10m")
    :retries(2, "exponential")
    :artifacts({"Dockerfile", "docker-build.log"})
    :on_success(function(params, output)
        log.info("Docker image built successfully: " .. output.image_tag)
        log.info("Build took: " .. output.build_duration .. " seconds")
    end)
    :on_failure(function(params, error)
        log.error("Docker build failed: " .. error)
        -- Send alert to team
        notifications.send("slack", {
            channel = "#build-alerts",
            message = "🚨 Docker build failed: " .. error
        })
    end)
    :build()

local test_task = task("run_tests")
    :description("Run comprehensive test suite")
    :command("npm test")
    :depends_on({"build_image"})
    :timeout("5m")
    :condition(when("params.skip_tests != true"))
    :on_success(function(params, output)
        log.info("✅ All tests passed!")
    end)
    :build()

local security_task = task("security_scan")
    :description("Security vulnerability scan with Trivy")
    :command("docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy myapp:latest")
    :depends_on({"build_image"})
    :timeout("3m")
    :condition(when("env.ENVIRONMENT == 'production'"))
    :on_failure(function(params, error)
        log.error("🔒 Security vulnerabilities found!")
        -- Block deployment on security issues
        return false
    end)
    :build()

local deploy_task = task("deploy_to_k8s")
    :description("Deploy to Kubernetes with rollback capability")
    :command(function(params, deps)
        local image_info = deps.build_image
        log.info("Deploying image: " .. image_info.image_tag)
        
        -- Use circuit breaker for Kubernetes API calls
        local deploy_result = circuit.protect("k8s_api", function()
            return exec.run("kubectl apply -f k8s/deployment.yaml")
        end)
        
        if not deploy_result.success then
            return false, "Deployment failed: " .. deploy_result.error
        end
        
        -- Wait for rollout to complete
        local rollout_result = exec.run("kubectl rollout status deployment/myapp --timeout=300s")
        if not rollout_result.success then
            log.warn("Rollout did not complete, initiating rollback...")
            exec.run("kubectl rollout undo deployment/myapp")
            return false, "Rollout failed, rolled back to previous version"
        end
        
        return true, "Deployed successfully", {
            deployment_time = os.time(),
            namespace = "production",
            image_deployed = image_info.image_tag,
            replicas_ready = 3
        }
    end)
    :depends_on({"run_tests", "security_scan"})
    :timeout("15m")
    :retries(1, "fixed")
    :artifacts({"k8s/deployment.yaml", "deployment-logs.txt"})
    :on_success(function(params, output)
        log.info("🚀 Deployment successful!")
        
        -- Send success notification
        notifications.send("slack", {
            channel = "#deployments",
            message = "🎉 Microservice deployed successfully!\n" ..
                     "Image: " .. output.image_deployed .. "\n" ..
                     "Time: " .. os.date("%Y-%m-%d %H:%M:%S", output.deployment_time)
        })
        
        -- Update deployment registry
        registry.record_deployment({
            service = "myapp",
            version = output.image_deployed,
            environment = "production",
            timestamp = output.deployment_time
        })
    end)
    :on_failure(function(params, error)
        log.error("💥 Deployment failed: " .. error)
        
        -- Send critical alert
        alerts.send("pagerduty", {
            severity = "critical",
            summary = "Microservice deployment failed",
            details = error
        })
    end)
    :build()

-- Define modern workflow with enhanced features
workflow.define("microservice_deployment", {
    description = "Deploy microservice to Kubernetes - Modern DSL",
    version = "3.0.0",
    
    metadata = {
        author = "Platform Team",
        team = "infrastructure",
        tags = {"microservice", "kubernetes", "docker", "deployment"},
        repository = "github.com/company/microservice",
        documentation = "https://docs.company.com/deployment-guide",
        cost_center = "engineering",
        criticality = "high"
    },
    
    tasks = { build_task, test_task, security_task, deploy_task },
    
    config = {
        timeout = "45m",
        max_parallel_tasks = 2,  -- tests and security can run in parallel
        retry_policy = "exponential",
        cleanup_on_failure = true,
        
        monitoring = {
            metrics = true,
            alerts = true,
            dashboard = "grafana://microservice-deployment"
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
        end,
        
        registry_accessible = function()
            local result = exec.run("docker info")
            return result.success, "Docker registry not accessible"
        end
    },
    
    on_start = function()
        log.info("🚀 Starting microservice deployment pipeline...")
        
        -- Send start notification
        notifications.send("slack", {
            channel = "#deployments",
            message = "🔄 Starting microservice deployment...",
            color = "warning"
        })
        
        -- Record deployment start
        metrics.increment("deployment_starts_total", {
            service = "myapp",
            environment = "production"
        })
        
        return true
    end,
    
    on_complete = function(success, results)
        local duration = metrics.stop_timer("deployment_duration")
        
        if success then
            log.info("✅ Microservice deployment completed successfully!")
            
            metrics.increment("deployment_success_total", {
                service = "myapp",
                duration = duration
            })
        else
            log.error("❌ Microservice deployment failed!")
            
            metrics.increment("deployment_failure_total", {
                service = "myapp",
                duration = duration
            })
        end
        
        return true
    end
})
```

## 🎯 Migration Checklist

### ✅ Pre-Migration
- [ ] Backup original files
- [ ] Review current workflow functionality
- [ ] Identify dependencies and integrations
- [ ] Plan migration phases

### ✅ During Migration
- [ ] Convert tasks to fluent API
- [ ] Add enhanced error handling
- [ ] Implement modern retry strategies
- [ ] Add metadata and monitoring
- [ ] Test each migrated component

### ✅ Post-Migration
- [ ] Validate workflow functionality
- [ ] Performance testing
- [ ] Update documentation
- [ ] Train team on new syntax
- [ ] Remove legacy code (optional)

## 🚨 Common Migration Issues

### Issue 1: Function Parameter Changes
**Problem:** `input_from_dependency` parameter name changed
**Solution:** Update to `deps` parameter

```lua
-- Before
command = function(params, input_from_dependency)
    local data = input_from_dependency.task_name.result
end

-- After  
command = function(params, deps)
    local data = deps.task_name.result
end
```

### Issue 2: Hook Name Changes
**Problem:** Hook names changed in Modern DSL
**Solution:** Update hook names

```lua
-- Before
pre_exec = function() ... end
post_exec = function() ... end

-- After
pre_hook = function() ... end  -- or :on_start()
post_hook = function() ... end -- or :on_complete()
```

### Issue 3: Dependency Format
**Problem:** Single dependencies need to be arrays
**Solution:** Always use array format

```lua
-- Before
depends_on = "task_name"

-- After
depends_on = {"task_name"}
```

## 🎉 Migration Success Tips

1. **Start Small**: Migrate simple workflows first
2. **Test Frequently**: Validate each step of migration
3. **Use Tools**: Leverage automated migration tools
4. **Gradual Rollout**: Deploy migrated workflows incrementally
5. **Monitor**: Watch for performance and functionality changes
6. **Document**: Keep migration notes for team reference

---

**The migration to Modern DSL opens up powerful new capabilities while maintaining all existing functionality. Take your time and leverage the tools provided for a smooth transition!**