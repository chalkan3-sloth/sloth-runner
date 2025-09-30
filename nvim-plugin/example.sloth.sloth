-- Example Sloth Runner DSL file for testing syntax highlighting
-- File: example.sloth.lua

local log = require("log")
local exec = require("exec")
local state = require("state")
local fs = require("fs")

-- Build task with enhanced features
local build_task = task("build_application")
    :description("Build the Go application with optimizations")
    :command(function(params, deps)
        log.info("🔨 Building application...")
        
        -- Check if source files exist
        if not fs.exists("./cmd/main.go") then
            return false, "main.go not found", { error = "missing_source" }
        end
        
        -- Build with version from environment
        local version = os.getenv("BUILD_VERSION") or "dev"
        local result = exec.run(string.format(
            "go build -ldflags '-X main.version=%s' -o app ./cmd/main.go",
            version
        ))
        
        if result.success then
            -- Store build artifact info in state
            state.set("last_build", {
                version = version,
                timestamp = os.time(),
                binary_path = "./app"
            })
            
            return true, result.stdout, { 
                artifact = "app",
                version = version,
                size = fs.size("./app")
            }
        else
            return false, result.stderr, { error = "build_failed" }
        end
    end)
    :timeout("10m")
    :retries(2, "exponential") 
    :tags({"build", "golang", "binary"})
    :artifacts({"app"})
    :on_success(function(params, output)
        log.info("✅ Build completed successfully")
        log.info("📦 Artifact: " .. output.artifact)
        log.info("📏 Size: " .. output.size .. " bytes")
    end)
    :on_failure(function(params, error)
        log.error("❌ Build failed: " .. error)
        -- TODO: Send notification to Slack
    end)
    :build()

-- Test task with dependencies
local test_task = task("run_tests")
    :description("Run comprehensive test suite")
    :depends_on({"build_application"})
    :command(function(params, deps)
        log.info("🧪 Running tests...")
        
        -- Get build info from previous task
        local build_info = deps.build_application
        log.info("Testing binary: " .. build_info.artifact)
        
        -- Run unit tests
        local unit_result = exec.run("go test -v ./...")
        if not unit_result.success then
            return false, unit_result.stderr, { phase = "unit_tests" }
        end
        
        -- Run integration tests  
        local integration_result = exec.run("go test -tags=integration ./tests/...")
        if not integration_result.success then
            return false, integration_result.stderr, { phase = "integration_tests" }
        end
        
        return true, "All tests passed", {
            unit_tests = "passed",
            integration_tests = "passed",
            coverage = "85%"
        }
    end)
    :timeout("15m")
    :condition(function(params)
        -- Only run tests if we're not in production
        return os.getenv("ENVIRONMENT") ~= "production"
    end)
    :build()

-- Deployment task with environment-specific logic
local deploy_task = task("deploy_application")
    :description("Deploy application to target environment")
    :depends_on({"build_application", "run_tests"})
    :command(function(params, deps)
        local env = params.environment or "staging"
        log.info("🚀 Deploying to " .. env .. "...")
        
        -- Environment-specific configuration
        local config = {
            staging = {
                namespace = "staging",
                replicas = 2,
                resources = { cpu = "100m", memory = "128Mi" }
            },
            production = {
                namespace = "production", 
                replicas = 5,
                resources = { cpu = "500m", memory = "512Mi" }
            }
        }
        
        local env_config = config[env]
        if not env_config then
            return false, "Unknown environment: " .. env
        end
        
        -- Deploy using kubectl with templating
        local template = string.format([[
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  namespace: %s
spec:
  replicas: %d
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:%s
        resources:
          requests:
            cpu: %s
            memory: %s
]], env_config.namespace, env_config.replicas, 
    deps.build_application.version,
    env_config.resources.cpu, 
    env_config.resources.memory)
        
        -- Write manifest to temp file
        local manifest_file = "/tmp/deployment.yaml"
        fs.write(manifest_file, template)
        
        -- Apply deployment
        local result = exec.run("kubectl apply -f " .. manifest_file)
        
        if result.success then
            -- Wait for rollout to complete
            local rollout_result = exec.run(string.format(
                "kubectl rollout status deployment/myapp -n %s --timeout=300s",
                env_config.namespace
            ))
            
            if rollout_result.success then
                return true, "Deployment successful", {
                    environment = env,
                    namespace = env_config.namespace,
                    replicas = env_config.replicas,
                    status = "running"
                }
            else
                return false, "Rollout failed: " .. rollout_result.stderr
            end
        else
            return false, "Deployment failed: " .. result.stderr
        end
    end)
    :timeout("20m")
    :agent("deploy-agent")
    :run_on("production_cluster")
    :circuit_breaker({
        failure_threshold = 3,
        recovery_timeout = "5m"
    })
    :build()

-- Cleanup task
local cleanup_task = task("cleanup_artifacts")
    :description("Clean up temporary files and artifacts")
    :command(function(params, deps)
        log.info("🧹 Cleaning up...")
        
        -- Remove temporary files
        local files_to_remove = {
            "./app",
            "/tmp/deployment.yaml",
            "./coverage.out"
        }
        
        for _, file in ipairs(files_to_remove) do
            if fs.exists(file) then
                fs.remove(file)
                log.info("Removed: " .. file)
            end
        end
        
        return true, "Cleanup completed"
    end)
    :condition(function(params)
        -- Always run cleanup unless explicitly disabled
        return params.skip_cleanup ~= true
    end)
    :build()

-- Define the complete CI/CD workflow
workflow.define("ci_cd_pipeline", {
    description = "Complete CI/CD pipeline with build, test, and deployment",
    version = "2.1.0",
    
    metadata = {
        author = "DevOps Team",
        category = "deployment",
        complexity = "advanced",
        tags = {"ci", "cd", "kubernetes", "golang"},
        environments = {"staging", "production"},
        estimated_duration = "25m"
    },
    
    -- Task execution order
    tasks = {
        build_task,
        test_task, 
        deploy_task,
        cleanup_task
    },
    
    -- Global workflow configuration
    defaults = {
        timeout = "30m",
        retry_policy = "exponential",
        max_parallel_tasks = 2
    },
    
    -- Success handler
    on_success = function(results)
        log.info("🎉 CI/CD Pipeline completed successfully!")
        
        -- Extract deployment info
        local deploy_result = results.deploy_application
        if deploy_result then
            log.info("🌍 Environment: " .. deploy_result.environment)
            log.info("📍 Namespace: " .. deploy_result.namespace) 
            log.info("🔢 Replicas: " .. deploy_result.replicas)
        end
        
        -- Store pipeline success in state
        state.set("last_successful_pipeline", {
            timestamp = os.time(),
            version = results.build_application.version,
            environment = deploy_result.environment
        })
        
        -- Send success notification
        notification.slack({
            channel = "#deployments",
            message = string.format(
                "✅ Deployment successful!\nVersion: %s\nEnvironment: %s", 
                results.build_application.version,
                deploy_result.environment
            ),
            color = "good"
        })
    end,
    
    -- Failure handler with detailed error reporting
    on_failure = function(error, context)
        log.error("💥 CI/CD Pipeline failed!")
        log.error("Failed task: " .. context.failed_task)
        log.error("Error: " .. error.message)
        
        -- Enhanced error context
        if context.phase then
            log.error("Phase: " .. context.phase)
        end
        
        -- Store failure info for analysis
        state.set("last_pipeline_failure", {
            timestamp = os.time(),
            failed_task = context.failed_task,
            error = error.message,
            context = context
        })
        
        -- Send failure notification with details
        notification.slack({
            channel = "#alerts",
            message = string.format(
                "❌ Deployment failed!\nTask: %s\nError: %s\nPhase: %s",
                context.failed_task,
                error.message,
                context.phase or "unknown"
            ),
            color = "danger",
            mention = ["@devops-team"]
        })
        
        -- Trigger rollback if deployment failed
        if context.failed_task == "deploy_application" then
            log.info("🔄 Initiating automatic rollback...")
            
            -- TODO: Implement rollback logic
            local rollback_result = exec.run("kubectl rollout undo deployment/myapp")
            if rollback_result.success then
                log.info("✅ Rollback completed")
            else
                log.error("❌ Rollback failed: " .. rollback_result.stderr)
            end
        end
    end,
    
    -- Cleanup handler (always runs)
    on_cleanup = function(results, error)
        log.info("🧹 Running workflow cleanup...")
        
        -- Always clean up resources
        exec.run("docker system prune -f")
        
        -- Update metrics
        metrics.counter("pipeline_executions_total", 1, {
            status = error and "failed" or "success",
            environment = os.getenv("ENVIRONMENT") or "unknown"
        })
    end
})