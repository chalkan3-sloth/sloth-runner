# Advanced Features

Sloth Runner provides powerful advanced features for enterprise deployments using the Modern DSL pattern.

---

## Table of Contents

- [Conditional Execution](#conditional-execution)
- [Lifecycle Hooks](#lifecycle-hooks)
- [Artifact Management](#artifact-management)
- [Error Handling & Retries](#error-handling--retries)
- [Parallel Execution](#parallel-execution)
- [Task Dependencies](#task-dependencies)
- [Asynchronous Tasks](#asynchronous-tasks)
- [Distributed Execution](#distributed-execution)

---

## Conditional Execution

Control when tasks run using `run_if` and `abort_if` conditions.

### Basic Conditional Execution

```lua
task("deploy_to_production")
    :description("Deploy only in production environment")
    :command(function(this, params)
        log.info("Deploying to production...")
        return true, "Deployment successful"
    end)
    :run_if(function(this, params)
        return os.getenv("ENV") == "production"
    end)
    :build()

task("skip_tests")
    :description("Skip tests if SKIP_TESTS is set")
    :command(function(this, params)
        exec.run("npm test")
        return true, "Tests passed"
    end)
    :run_if(function(this, params)
        return os.getenv("SKIP_TESTS") ~= "true"
    end)
    :build()
```

### Abort Workflow on Condition

```lua
task("security_check")
    :description("Abort workflow if security vulnerabilities found")
    :command(function(this, params)
        local result = exec.run("npm audit --json")
        local vulnerabilities = json.decode(result.stdout)
        return true, "Security check completed", { vulnerabilities = vulnerabilities }
    end)
    :abort_if(function(this, params)
        local result = exec.run("npm audit")
        return result.exit_code ~= 0
    end)
    :build()
```

---

## Lifecycle Hooks

Execute custom logic at different stages of task execution.

### Complete Lifecycle Example

```lua
task("deploy_with_hooks")
    :description("Deployment with full lifecycle management")
    :pre_hook(function(this, params)
        log.info("Pre-deployment: Taking snapshot...")
        exec.run("aws ec2 create-snapshot --volume-id vol-123")
        return true, "Snapshot created"
    end)
    :command(function(this, params)
        log.info("Deploying application...")
        exec.run("kubectl apply -f deployment.yaml")
        return true, "Application deployed"
    end)
    :post_hook(function(this, params)
        log.info("Post-deployment: Running smoke tests...")
        exec.run("./smoke-tests.sh")
        return true, "Smoke tests passed"
    end)
    :on_success(function(this, params)
        log.success("Deployment successful! Sending notification...")
        exec.run("slack-notify 'Deployment successful'")
        return true, "Notification sent"
    end)
    :on_failure(function(this, params, err)
        log.error("Deployment failed: " .. err)
        exec.run("slack-notify 'Deployment failed: " .. err .. "'")
        exec.run("kubectl rollback deployment/myapp")
        return false, "Rollback initiated"
    end)
    :on_timeout(function(this, params)
        log.warn("Deployment timed out, rolling back...")
        exec.run("kubectl rollback deployment/myapp")
        return false, "Timeout rollback"
    end)
    :timeout("10m")
    :build()
```

### Multi-Stage Deployment with Hooks

```lua
task("blue_green_deployment")
    :description("Blue-green deployment with validation")
    :pre_hook(function(this, params)
        log.info("Preparing green environment...")
        exec.run("terraform apply -target=module.green")
        return true, "Green environment ready"
    end)
    :command(function(this, params)
        log.info("Deploying to green environment...")
        exec.run("ansible-playbook deploy-green.yml")
        return true, "Deployed to green"
    end)
    :post_hook(function(this, params)
        log.info("Validating green environment...")
        local result = exec.run("./health-check.sh green")
        if result.exit_code == 0 then
            log.info("Switching traffic to green...")
            exec.run("./switch-traffic.sh green")
            return true, "Traffic switched"
        else
            return false, "Health check failed"
        end
    end)
    :on_success(function(this, params)
        log.info("Decommissioning blue environment...")
        exec.run("terraform destroy -target=module.blue")
        return true, "Blue environment removed"
    end)
    :on_failure(function(this, params, err)
        log.error("Deployment failed, cleaning up green...")
        exec.run("terraform destroy -target=module.green")
        return false, "Cleanup completed"
    end)
    :build()
```

---

## Artifact Management

Share files between tasks using the artifact system.

### Basic Artifact Workflow

```lua
local build_task = task("build")
    :description("Build application and produce artifacts")
    :command(function(this, params)
        log.info("Building application...")
        exec.run("go build -o myapp")
        exec.run("tar czf myapp.tar.gz myapp")
        return true, "Build successful"
    end)
    :artifacts({"myapp.tar.gz", "myapp"})
    :build()

local test_task = task("test")
    :description("Test using built artifacts")
    :depends_on({"build"})
    :consumes({"myapp"})
    :command(function(this, params)
        log.info("Running tests with artifact...")
        exec.run("./myapp --version")
        exec.run("go test ./...")
        return true, "Tests passed"
    end)
    :build()

local deploy_task = task("deploy")
    :description("Deploy using built artifacts")
    :depends_on({"test"})
    :consumes({"myapp.tar.gz"})
    :command(function(this, params)
        log.info("Deploying artifact...")
        exec.run("scp myapp.tar.gz server:/opt/app/")
        exec.run("ssh server 'cd /opt/app && tar xzf myapp.tar.gz'")
        return true, "Deployment successful"
    end)
    :build()

workflow.define("ci_cd_pipeline")
    :description("Complete CI/CD with artifact management")
    :version("1.0.0")
    :tasks({build_task, test_task, deploy_task})
```

### Multi-Stage Artifact Pipeline

```lua
local compile_task = task("compile")
    :description("Compile source code")
    :command(function(this, params)
        exec.run("gcc -o program main.c")
        return true, "Compilation successful"
    end)
    :artifacts({"program"})
    :build()

local package_task = task("package")
    :description("Package compiled program")
    :depends_on({"compile"})
    :consumes({"program"})
    :command(function(this, params)
        exec.run("docker build -t myapp:latest .")
        exec.run("docker save myapp:latest > myapp.tar")
        return true, "Package created"
    end)
    :artifacts({"myapp.tar"})
    :build()

local distribute_task = task("distribute")
    :description("Distribute packaged application")
    :depends_on({"package"})
    :consumes({"myapp.tar"})
    :command(function(this, params)
        exec.run("aws s3 cp myapp.tar s3://artifacts/myapp/")
        return true, "Distribution complete"
    end)
    :build()

workflow.define("artifact_pipeline")
    :description("Multi-stage artifact pipeline")
    :version("1.0.0")
    :tasks({compile_task, package_task, distribute_task})
```

---

## Error Handling & Retries

Implement robust error handling with automatic retries.

### Retry Strategies

```lua
task("flaky_api_call")
    :description("Call external API with exponential backoff")
    :command(function(this, params)
        local response = http.get("https://api.example.com/data")
        if response.status ~= 200 then
            return false, "API call failed with status " .. response.status
        end
        return true, "API call successful"
    end)
    :retries(5, "exponential")
    :timeout("30s")
    :build()

task("database_migration")
    :description("Run database migration with fixed retry interval")
    :command(function(this, params)
        local result = exec.run("./migrate.sh")
        return result.exit_code == 0, result.stdout
    end)
    :retries(3, "fixed")
    :timeout("5m")
    :build()

task("network_operation")
    :description("Network operation with linear backoff")
    :command(function(this, params)
        exec.run("rsync -avz /data remote:/backup/")
        return true, "Sync completed"
    end)
    :retries(3, "linear")
    :build()
```

### Error Recovery Workflow

```lua
local main_task = task("main_operation")
    :description("Main operation with error recovery")
    :command(function(this, params)
        local result = exec.run("./critical-operation.sh")
        if result.exit_code ~= 0 then
            return false, "Operation failed: " .. result.stderr
        end
        return true, "Operation successful"
    end)
    :retries(3, "exponential")
    :on_failure(function(this, params, err)
        log.error("All retries exhausted, initiating recovery...")
        exec.run("./recovery.sh")
        return false, "Recovery initiated"
    end)
    :build()

workflow.define("resilient_workflow")
    :description("Workflow with comprehensive error handling")
    :version("1.0.0")
    :tasks({main_task})
```

---

## Parallel Execution

Execute multiple tasks concurrently for improved performance.

### Parallel Task Groups

```lua
task("lint")
    :description("Run linting")
    :command(function(this, params)
        exec.run("eslint .")
        return true, "Linting complete"
    end)
    :build()

task("type_check")
    :description("Run type checking")
    :command(function(this, params)
        exec.run("tsc --noEmit")
        return true, "Type checking complete"
    end)
    :build()

task("unit_tests")
    :description("Run unit tests")
    :command(function(this, params)
        exec.run("jest")
        return true, "Unit tests complete"
    end)
    :build()

task("integration_tests")
    :description("Run integration tests")
    :command(function(this, params)
        exec.run("npm run test:integration")
        return true, "Integration tests complete"
    end)
    :build()

task("parallel_quality_checks")
    :description("Run all quality checks in parallel")
    :command(function(this, params)
        local results, err = parallel({
            { name = "lint", command = "eslint ." },
            { name = "typecheck", command = "tsc --noEmit" },
            { name = "test", command = "jest" },
            { name = "audit", command = "npm audit" }
        })

        if err then
            return false, "Some checks failed: " .. err
        end

        return true, "All checks passed"
    end)
    :build()

workflow.define("quality_pipeline")
    :description("Parallel quality checks workflow")
    :version("1.0.0")
    :tasks({parallel_quality_checks})
```

---

## Task Dependencies

Create complex dependency chains between tasks.

### Linear Dependencies

```lua
local install = task("install")
    :description("Install dependencies")
    :command(function(this, params)
        exec.run("npm install")
        return true, "Dependencies installed"
    end)
    :build()

local build = task("build")
    :description("Build application")
    :depends_on({"install"})
    :command(function(this, params)
        exec.run("npm run build")
        return true, "Build complete"
    end)
    :build()

local test = task("test")
    :description("Run tests")
    :depends_on({"build"})
    :command(function(this, params)
        exec.run("npm test")
        return true, "Tests passed"
    end)
    :build()

local deploy = task("deploy")
    :description("Deploy application")
    :depends_on({"test"})
    :command(function(this, params)
        exec.run("./deploy.sh")
        return true, "Deployment complete"
    end)
    :build()

workflow.define("linear_pipeline")
    :description("Linear dependency pipeline")
    :version("1.0.0")
    :tasks({install, build, test, deploy})
```

### Complex Dependencies

```lua
local setup_db = task("setup_db")
    :description("Setup database")
    :command(function(this, params)
        exec.run("docker-compose up -d postgres")
        return true, "Database ready"
    end)
    :build()

local setup_cache = task("setup_cache")
    :description("Setup cache")
    :command(function(this, params)
        exec.run("docker-compose up -d redis")
        return true, "Cache ready"
    end)
    :build()

local migrate = task("migrate")
    :description("Run migrations")
    :depends_on({"setup_db"})
    :command(function(this, params)
        exec.run("./migrate.sh")
        return true, "Migrations complete"
    end)
    :build()

local seed = task("seed")
    :description("Seed database")
    :depends_on({"migrate"})
    :command(function(this, params)
        exec.run("./seed.sh")
        return true, "Seeding complete"
    end)
    :build()

local start_app = task("start_app")
    :description("Start application")
    :depends_on({"seed", "setup_cache"})
    :command(function(this, params)
        exec.run("npm start")
        return true, "Application started"
    end)
    :build()

workflow.define("complex_startup")
    :description("Complex dependency workflow")
    :version("1.0.0")
    :tasks({setup_db, setup_cache, migrate, seed, start_app})
```

---

## Asynchronous Tasks

Run tasks in the background without blocking workflow execution.

### Background Tasks

```lua
task("monitoring")
    :description("Start monitoring in background")
    :command(function(this, params)
        exec.run("./monitor.sh &")
        return true, "Monitoring started"
    end)
    :async(true)
    :build()

task("main_process")
    :description("Main application process")
    :command(function(this, params)
        log.info("Starting main process...")
        exec.run("./app")
        return true, "Process complete"
    end)
    :build()

workflow.define("async_workflow")
    :description("Workflow with async background tasks")
    :version("1.0.0")
    :tasks({
        task("monitoring"):async(true):command(function(this, params)
            exec.run("./monitor.sh")
            return true, "Monitoring"
        end):build(),
        task("main_process"):command(function(this, params)
            exec.run("./app")
            return true, "Complete"
        end):build()
    })
```

---

## Distributed Execution

Execute tasks across multiple agents in a distributed environment.

### Agent-Targeted Tasks

```lua
task("deploy_web_servers")
    :description("Deploy to web server agents")
    :target_agent("web-*")
    :command(function(this, params)
        exec.run("systemctl restart nginx")
        return true, "Web server restarted"
    end)
    :build()

task("backup_databases")
    :description("Backup on database agents")
    :target_agent("db-*")
    :command(function(this, params)
        exec.run("pg_dump mydb > /backup/mydb.sql")
        return true, "Backup complete"
    end)
    :build()

task("update_monitoring")
    :description("Update monitoring agents")
    :target_agent("monitor-*")
    :command(function(this, params)
        exec.run("./update-prometheus.sh")
        return true, "Monitoring updated"
    end)
    :build()

workflow.define("distributed_operations")
    :description("Distributed multi-agent operations")
    :version("1.0.0")
    :tasks({deploy_web_servers, backup_databases, update_monitoring})
```

### Tag-Based Execution

```lua
task("security_scan")
    :description("Run security scan on production servers")
    :target_tags({"production", "web"})
    :command(function(this, params)
        exec.run("trivy filesystem /")
        return true, "Security scan complete"
    end)
    :build()

task("performance_test")
    :description("Run performance tests on staging")
    :target_tags({"staging", "app"})
    :command(function(this, params)
        exec.run("k6 run performance-test.js")
        return true, "Performance test complete"
    end)
    :build()

workflow.define("environment_specific")
    :description("Environment-specific distributed tasks")
    :version("1.0.0")
    :tasks({security_scan, performance_test})
```

---

## Related Documentation

- [AI Integration](../ai-integration.md) - Intelligent workflow optimization
- [Distributed Agents](../distributed-agents.md) - Multi-agent orchestration
- [Web Dashboard](../web-dashboard.md) - Visual management interface
- [Advanced Scheduler](../advanced-scheduler.md) - Cron-based scheduling
- [Multi-Cloud Excellence](../multi-cloud-excellence.md) - Cloud provider integration
- [Stack Management](../stack-management.md) - State and environment management

---

## Getting Started

See [Getting Started Guide](./getting-started.md) for installation instructions.

---

## More Examples

Check out [Advanced Examples](./advanced-examples.md) for complete end-to-end use cases.

---

[English](./advanced-features.md) | [Português](../pt/advanced-features.md) | [中文](../zh/advanced-features.md)
