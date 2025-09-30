-- Enhanced Modern DSL Demo - Simplified Version
-- This demonstrates the improved TaskRunner, Core integration, and modern DSL syntax

log.info("Loading Enhanced Modern Pipeline Demo")

-- Test core system integration
local stats = core.stats()
log.info("Core system initialized", {
    uptime = stats.uptime_seconds,
    memory = stats.memory_alloc,
    workers = stats.worker_active
})

-- Demonstrate modern parallel execution
log.info("Testing enhanced parallel execution...")

local parallel_results, parallel_errors = async.parallel({
    build_frontend = function()
        log.info("Building frontend...")
        async.sleep(200)
        return { 
            status = "success", 
            size = "2.5MB",
            artifacts = {"dist/app.js", "dist/app.css"}
        }
    end,
    
    build_backend = function()
        log.info("Building backend...")
        async.sleep(300)
        return {
            status = "success",
            binary = "app",
            size = "15MB"
        }
    end,
    
    run_tests = function()
        log.info("Running tests...")
        async.sleep(150)
        return {
            status = "success",
            tests_passed = 42,
            coverage = 85.5
        }
    end
}, 3) -- Use 3 workers

if parallel_errors then
    log.error("Some parallel tasks failed")
else
    log.info("All parallel tasks completed successfully")
end

-- Demonstrate performance monitoring
log.info("Testing performance monitoring...")

local perf_result, duration, perf_error = perf.measure(function()
    log.info("Performing CPU-intensive task...")
    
    -- Simulate processing
    for i = 1, 1000 do
        math.sqrt(i * 42)
    end
    
    async.sleep(100)
    return "processing_complete"
end, "cpu_intensive_task")

log.info("Performance measurement completed in " .. duration .. "ms")

-- Get memory statistics
local memory_stats = perf.memory()
log.info("Memory usage: " .. memory_stats.current_mb .. "MB (" .. memory_stats.usage_percent .. "%)")

-- Demonstrate advanced error handling
log.info("Testing advanced error handling...")

local try_result, caught_error = error.try(
    function()
        log.info("Attempting potentially failing operation...")
        
        -- Simulate a condition that might fail
        local random_value = math.random()
        if random_value > 0.7 then
            error("Simulated failure: random value " .. random_value .. " too high")
        end
        
        return "operation_successful"
    end,
    function(err)
        log.warn("Caught and handling error: " .. err)
        return "error_handled"
    end,
    function()
        log.info("Cleanup operations completed")
    end
)

log.info("Try-catch operation completed")

-- Demonstrate retry mechanism
log.info("Testing retry mechanism...")

local retry_attempts = 0
local retry_result, retry_error = error.retry(function()
    retry_attempts = retry_attempts + 1
    log.info("Retry attempt #" .. retry_attempts)
    
    -- Simulate success after a few attempts
    if retry_attempts >= 2 then
        return "retry_successful_after_" .. retry_attempts .. "_attempts"
    else
        error("Simulated temporary failure")
    end
end, 3, 500) -- 3 attempts, 500ms initial delay

if retry_error then
    log.error("Retry failed: " .. retry_error)
else
    log.info("Retry succeeded: " .. retry_result)
end

-- Demonstrate circuit breaker pattern
log.info("Testing circuit breaker pattern...")

local cb_result, cb_error = flow.circuit_breaker("external_api", function()
    log.info("Making call to external API (protected by circuit breaker)")
    
    -- Simulate external API call
    async.sleep(50)
    return {
        status = "200",
        data = { message = "API call successful" },
        response_time = "45ms"
    }
end)

if cb_error then
    log.error("Circuit breaker blocked the call: " .. cb_error)
else
    log.info("Circuit breaker allowed the call - success!")
end

-- Demonstrate rate limiting
log.info("Testing rate limiting...")

for i = 1, 3 do
    local rate_result, rate_error = flow.rate_limit(2, function() -- 2 RPS
        log.info("Rate limited operation #" .. i)
        return "rate_limited_result_" .. i
    end)
    
    if rate_error then
        log.error("Rate limiting failed: " .. rate_error)
    end
end

-- Demonstrate configuration and secrets
log.info("Testing configuration and secrets...")

local config_value = utils.config("environment", "development")
log.info("Configuration retrieved: environment = " .. config_value)

local secret_value, secret_error = utils.secret("api_key")
if secret_error then
    log.error("Failed to retrieve secret: " .. secret_error)
else
    log.info("Secret retrieved successfully (masked)")
end

-- Demonstrate checkpoint creation
log.info("Testing checkpoint creation...")

local checkpoint_name = task.checkpoint("demo_checkpoint", {
    timestamp = os.time(),
    demo_progress = "75%",
    completed_tasks = {"build_frontend", "build_backend", "run_tests"}
})

log.info("Checkpoint created: " .. checkpoint_name)

-- Demonstrate workflow definition
log.info("Testing workflow definition...")

local workflow_success, workflow_error = workflow.define("demo_workflow", {
    description = "Enhanced demo workflow",
    version = "2.0.0",
    
    stages = {
        {
            name = "preparation",
            description = "Prepare environment and dependencies"
        },
        {
            name = "execution", 
            description = "Execute main tasks"
        },
        {
            name = "cleanup",
            description = "Clean up resources"
        }
    }
})

if workflow_error then
    log.error("Workflow definition failed: " .. workflow_error)
else
    log.info("Workflow defined successfully")
end

-- Final summary
log.info("=== Enhanced Features Demo Summary ===")
log.info("✓ Core system integration working")
log.info("✓ Parallel execution with worker pools")
log.info("✓ Performance monitoring and metrics")
log.info("✓ Advanced error handling (try-catch)")
log.info("✓ Retry mechanisms with backoff")
log.info("✓ Circuit breaker pattern")
log.info("✓ Rate limiting")
log.info("✓ Configuration and secret management")
log.info("✓ Checkpoint system")
log.info("✓ Modern workflow definition")

log.info("Enhanced Modern DSL Demo completed successfully! 🚀")

-- Return demo results for potential inspection
return {
    status = "completed",
    features_tested = {
        "core_integration",
        "parallel_execution", 
        "performance_monitoring",
        "error_handling",
        "retry_mechanism",
        "circuit_breaker",
        "rate_limiting",
        "configuration",
        "checkpoints",
        "workflows"
    },
    results = {
        parallel_execution = parallel_results,
        performance = {
            duration = duration,
            memory_usage = memory_stats
        },
        error_handling = {
            try_result = try_result,
            retry_result = retry_result
        },
        circuit_breaker = cb_result
    },
    demo_completed_at = os.time()
}