-- MODERN DSL ONLY - Docker Operations Example
-- Demonstrates Docker operations using Modern DSL

-- Docker System Info task
local docker_info = task("docker_system_info")
    :description("Get Docker system information")
    :command(function(params)
        log.info("🐳 Getting Docker system information...")
        
        local result = exec.run("docker system info --format '{{.ServerVersion}}'", {
            timeout = "30s",
            capture_output = true
        })
        
        if result.success then
            local version = string.gsub(result.output or "", "%s+", "")
            log.info("🔍 Docker version: " .. version)
            
            return true, result.output, {
                docker_version = version,
                system_check = "passed"
            }
        else
            return false, "Docker not available or not running"
        end
    end)
    :timeout("60s")
    :retries(2, "exponential")
    :build()

-- Docker Build task
local docker_build = task("docker_build_image")
    :description("Build Docker image from Dockerfile")
    :depends_on({"docker_system_info"})
    :command(function(params, deps)
        log.info("🔨 Building Docker image...")
        
        -- Create a simple Dockerfile for demo
        local dockerfile_content = [[
FROM alpine:latest
RUN echo "Hello from Sloth Runner Modern DSL!" > /hello.txt
CMD cat /hello.txt
]]
        
        -- Write Dockerfile
        local success, err = fs.write("Dockerfile.sloth", dockerfile_content)
        if err then
            return false, "Failed to create Dockerfile: " .. err
        end
        
        -- Build the image
        local result = exec.run("docker build -f Dockerfile.sloth -t sloth-runner:demo .", {
            timeout = "300s",
            capture_output = true
        })
        
        if result.success then
            return true, "Docker image built successfully", {
                image_name = "sloth-runner:demo",
                dockerfile = "Dockerfile.sloth",
                docker_version = deps.docker_system_info.docker_version
            }
        else
            return false, "Failed to build Docker image: " .. (result.error or "unknown error")
        end
    end)
    :timeout("600s")
    :artifacts({"Dockerfile.sloth"})
    :on_success(function(params, output)
        log.info("✅ Docker image " .. output.image_name .. " built successfully")
    end)
    :build()

-- Docker Run task
local docker_run = task("docker_run_container")
    :description("Run Docker container")
    :depends_on({"docker_build_image"})
    :command(function(params, deps)
        log.info("🚀 Running Docker container...")
        
        local image_name = deps.docker_build_image.image_name
        local result = exec.run("docker run --rm " .. image_name, {
            timeout = "60s",
            capture_output = true
        })
        
        if result.success then
            log.info("📄 Container output: " .. (result.output or ""))
            return true, result.output, {
                container_output = result.output,
                image_used = image_name,
                execution_mode = "run_and_remove"
            }
        else
            return false, "Failed to run container"
        end
    end)
    :timeout("120s")
    :on_success(function(params, output)
        log.info("🎉 Container executed successfully")
    end)
    :build()

-- Docker Cleanup task
local docker_cleanup = task("docker_cleanup")
    :description("Clean up Docker resources")
    :depends_on({"docker_run_container"})
    :command(function(params, deps)
        log.info("🧹 Cleaning up Docker resources...")
        
        -- Remove the built image
        local image_name = deps.docker_build_image.image_name
        local result = exec.run("docker rmi " .. image_name, {
            timeout = "60s",
            capture_output = true
        })
        
        -- Remove Dockerfile
        local cleanup_success = fs.remove("Dockerfile.sloth")
        
        if result.success and cleanup_success then
            return true, "Cleanup completed", {
                image_removed = image_name,
                dockerfile_removed = "Dockerfile.sloth"
            }
        else
            return false, "Cleanup partially failed"
        end
    end)
    :timeout("90s")
    :build()

-- Modern Workflow Definition
workflow.define("docker_operations", {
    description = "Docker Operations Workflow - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"docker", "containers", "build", "run", "cleanup", "modern-dsl"},
        created_at = os.date(),
        prerequisites = "Docker installed and running"
    },
    
    tasks = {
        docker_info,
        docker_build,
        docker_run,
        docker_cleanup
    },
    
    config = {
        timeout = "20m",
        retry_policy = "exponential",
        max_parallel_tasks = 1,  -- Sequential execution for Docker operations
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting Docker operations workflow...")
        log.info("🐳 Ensure Docker is installed and running")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Docker operations workflow completed successfully!")
            log.info("🐳 Docker image built, run, and cleaned up")
        else
            log.error("❌ Docker operations workflow failed!")
            log.warn("🔍 Check Docker installation and daemon status")
        end
        return true
    end
})
