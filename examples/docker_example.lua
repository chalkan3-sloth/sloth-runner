-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:30 -03

local config = {
local build_image_task = task("build_image")
local run_container_task = task("run_container")
local cleanup_task = task("cleanup")

local build_image_task = task("build_image")
local run_container_task = task("run_container")
local cleanup_task = task("cleanup")
local build_image_task = task("build_image")
    :description("Builds Docker image with modern DSL")
    :command(function()
        log.info("Modern DSL: Building Docker image with tag: " .. config.image_tag)
        
        -- Enhanced Docker build with better error handling
        local result = docker.build({
            tag = config.image_tag,
            path = config.dockerfile_path,
            no_cache = false,
            progress = "auto"
        })

        if not result.success then
            log.error("Modern DSL: Docker build failed: " .. result.stderr)
            return false, "Docker build failed: " .. result.stderr
        end

        log.info("Modern DSL: Docker image built successfully")
        return true, "Docker image built successfully", {
            image_tag = config.image_tag,
            image_id = result.image_id or "unknown",
            build_time = result.duration or 0
        }
    end)
    :timeout("10m")
    :on_success(function(params, output)
        log.info("Modern DSL: Image ready: " .. output.image_tag)
    end)
    :build()
local run_container_task = task("run_container")
    :description("Runs container to verify image with modern DSL")
    :depends_on({"build_image"})
    :command(function(params, deps)
        local image_info = deps.build_image
        log.info("Modern DSL: Running container from image: " .. image_info.image_tag)
        
        -- Enhanced container run with configuration
        local result = docker.run({
            image = config.image_tag,
            name = config.container_name,
            remove = true, -- Auto remove after execution
            interactive = false,
            tty = false,
            env = {
                "ENV=test",
                "VERSION=" .. (image_info.image_id or "unknown")
            }
        })

        if not result.success then
            log.error("Modern DSL: Container run failed: " .. result.stderr)
            return false, "Container run failed: " .. result.stderr
        end

        log.info("Modern DSL: Container output: " .. result.stdout)
        return true, "Container ran successfully", {
            container_name = config.container_name,
            output = result.stdout,
            exit_code = result.exit_code or 0
        }
    end)
    :timeout("5m")
    :build()
local cleanup_task = task("cleanup")
    :description("Cleans up Docker resources with modern DSL")
    :depends_on({"run_container"})
    :command(function()
        log.info("Modern DSL: Cleaning up Docker resources...")
        
        local cleanup_results = {}
        
        -- Remove container (if not auto-removed)
        local container_result = docker.container.remove({
            name = config.container_name,
            force = true
        })
        cleanup_results.container_removed = container_result.success
        
        if not container_result.success then
            log.warn("Container cleanup: " .. (container_result.stderr or "already removed"))
        end
        
        -- Remove image
        local image_result = docker.image.remove({
            tag = config.image_tag,
            force = true
        })
        cleanup_results.image_removed = image_result.success
        
        if image_result.success then
            log.info("Modern DSL: Docker image removed successfully")
        else
            log.warn("Image cleanup: " .. (image_result.stderr or "already removed"))
        end
        
        return true, "Cleanup completed", cleanup_results
    end)
    :always_run(true) -- Run even if previous tasks failed
    :build()

workflow.define("docker_build_pipeline_modern", {
    description = "Docker build and test pipeline - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        category = "infrastructure",
        tags = {"docker", "build", "container", "modern-dsl"}
    },
    
    tasks = {
        build_image_task,
        run_container_task,
        cleanup_task
    },
    
    config = {
        timeout = "20m",
        cleanup_on_failure = true,
        docker_config = config
    },
    
    on_start = function()
        log.info("Starting Docker build pipeline...")
        -- Verify Docker is available
        local docker_version = docker.version()
        if docker_version.success then
            log.info("Docker version: " .. docker_version.version)
        else
            log.error("Docker not available!")
            return false
        end
        return true
    end,
    
    on_failure = function(task_name, error)
        log.warn("Docker pipeline failure at " .. task_name .. ": " .. error)
        -- Force cleanup on failure
        docker.container.remove({name = config.container_name, force = true})
        docker.image.remove({tag = config.image_tag, force = true})
        return true
    end
})
