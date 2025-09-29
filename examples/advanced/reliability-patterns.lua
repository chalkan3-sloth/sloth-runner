-- MODERN DSL - Advanced Reliability Patterns
-- This example demonstrates enterprise-grade reliability patterns
-- using the Modern DSL features

-- Task 1: Circuit Breaker Pattern
local resilient_api_task = task("resilient_api_call")
    :description("API call with circuit breaker protection")
    :command(function(params)
        log.info("🔄 Making resilient API call...")
        
        -- Use circuit breaker for external API protection
        local result = circuit.protect("external_api", function()
            -- Simulate API call with potential failure
            if math.random() > 0.3 then
                return {
                    success = true,
                    data = { response = "API success", timestamp = os.time() }
                }
            else
                error("API temporarily unavailable")
            end
        end)
        
        if result.success then
            log.info("✅ API call succeeded")
            return true, "API call completed", result.data
        else
            log.warn("⚠️  API call failed, circuit breaker activated")
            return false, "Circuit breaker open"
        end
    end)
    :retries(3, "exponential")
    :timeout("30s")
    :on_failure(function(params, error)
        log.error("API task failed: " .. error)
    end)
    :build()

-- Task 2: Retry with Exponential Backoff
local flaky_service_task = task("flaky_service")
    :description("Service with intelligent retry strategy")
    :depends_on({"resilient_api_call"})
    :command(function(params, deps)
        log.info("🔄 Calling flaky service...")
        
        -- Simulate a service that fails sometimes
        local success_rate = 0.7
        if math.random() < success_rate then
            log.info("✅ Flaky service succeeded")
            return true, "Service call successful", {
                service_data = "processed_data",
                attempts = 1
            }
        else
            log.warn("⚠️  Flaky service failed, will retry...")
            error("Service temporarily unavailable")
        end
    end)
    :retries(5, "exponential")  -- 5 retries with exponential backoff
    :retry_delay("1s")          -- Initial delay
    :max_retry_delay("30s")     -- Maximum delay cap
    :on_retry(function(attempt, error)
        log.warn("🔄 Retry attempt " .. attempt .. " due to: " .. error)
    end)
    :build()

-- Task 3: Saga Pattern Implementation
local saga_coordinator = task("saga_coordinator")
    :description("Coordinates distributed transaction with compensation")
    :depends_on({"flaky_service"})
    :command(function(params, deps)
        log.info("🎯 Starting distributed saga...")
        
        local saga_steps = {}
        local compensations = {}
        
        -- Step 1: Reserve resources
        local step1_success = true -- Simulate step
        if step1_success then
            table.insert(saga_steps, "resource_reserved")
            table.insert(compensations, "release_resources")
            log.info("✅ Step 1: Resources reserved")
        end
        
        -- Step 2: Process payment (simulate failure)
        local step2_success = math.random() > 0.3
        if step2_success then
            table.insert(saga_steps, "payment_processed")
            table.insert(compensations, "refund_payment")
            log.info("✅ Step 2: Payment processed")
        else
            log.error("❌ Step 2: Payment failed, executing compensation...")
            
            -- Execute compensations in reverse order
            for i = #compensations, 1, -1 do
                log.info("🔄 Compensating: " .. compensations[i])
            end
            
            return false, "Saga failed, compensated successfully"
        end
        
        log.info("🎉 Saga completed successfully!")
        return true, "Distributed transaction completed", {
            saga_steps = saga_steps,
            total_steps = #saga_steps
        }
    end)
    :timeout("2m")
    :build()

-- Task 4: Health Check with Monitoring
local health_monitor = task("health_monitor")
    :description("System health monitoring with alerts")
    :command(function(params)
        log.info("🔍 Performing health checks...")
        
        local health_status = {
            database = "healthy",
            api = "healthy",
            cache = "healthy"
        }
        
        -- Simulate health check logic
        local checks = {"database", "api", "cache"}
        for _, service in ipairs(checks) do
            -- Random health check (90% success rate)
            if math.random() > 0.9 then
                health_status[service] = "unhealthy"
                log.warn("⚠️  " .. service .. " health check failed")
            else
                log.info("✅ " .. service .. " is healthy")
            end
        end
        
        -- Count unhealthy services
        local unhealthy_count = 0
        for _, status in pairs(health_status) do
            if status == "unhealthy" then
                unhealthy_count = unhealthy_count + 1
            end
        end
        
        if unhealthy_count > 0 then
            log.error("⚠️  " .. unhealthy_count .. " services are unhealthy")
            return false, "Health check failed", health_status
        else
            log.info("🎉 All services are healthy")
            return true, "All systems operational", health_status
        end
    end)
    :async(true)
    :timeout("45s")
    :build()

-- Advanced Workflow with Reliability Patterns
workflow.define("reliability_patterns_demo", {
    description = "Advanced reliability patterns demonstration - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"reliability", "resilience", "enterprise", "modern-dsl"},
        complexity = "advanced",
        patterns = {"circuit_breaker", "retry", "saga", "health_check"}
    },
    
    tasks = {
        resilient_api_task,
        flaky_service_task,
        saga_coordinator,
        health_monitor
    },
    
    config = {
        timeout = "10m",
        retry_policy = "exponential",
        max_parallel_tasks = 2,
        fail_fast = false,
        circuit_breaker = {
            failure_threshold = 5,
            recovery_timeout = "30s",
            half_open_requests = 3
        }
    },
    
    on_start = function()
        log.info("🚀 Starting reliability patterns demonstration...")
        log.info("🔧 This workflow demonstrates:")
        log.info("   • Circuit Breaker Pattern")
        log.info("   • Exponential Backoff Retries")
        log.info("   • Saga Pattern for Distributed Transactions")
        log.info("   • Health Check Monitoring")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Reliability patterns demo completed successfully!")
            log.info("📊 All enterprise patterns executed properly")
        else
            log.warn("⚠️  Demo completed with some failures (expected for demonstration)")
            log.info("🔍 Check individual task results for pattern behavior")
        end
        
        -- Log circuit breaker stats
        local cb_stats = circuit.stats("external_api")
        if cb_stats then
            log.info("🔧 Circuit Breaker Stats:")
            log.info("   • State: " .. (cb_stats.state or "unknown"))
            log.info("   • Failures: " .. (cb_stats.failures or 0))
            log.info("   • Successes: " .. (cb_stats.successes or 0))
        end
        
        return true
    end
})
