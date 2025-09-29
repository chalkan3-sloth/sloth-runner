-- MODERN DSL ONLY - Parallel Execution Showcase
-- Demonstrates advanced parallel processing capabilities

-- Task 1: Fast CPU task
local cpu_task = task("cpu_intensive")
    :description("CPU-intensive task simulation")
    :command(function()
        log.info("💻 Starting CPU-intensive computation...")
        
        -- Simulate CPU work with performance monitoring
        local result = perf.measure(function()
            local sum = 0
            for i = 1, 1000000 do
                sum = sum + math.sqrt(i)
            end
            return sum
        end)
        
        log.info("⚡ CPU task completed in " .. result.duration .. "ms")
        return true, "CPU computation completed", {
            computation_result = result.value,
            duration_ms = result.duration,
            operations = 1000000
        }
    end)
    :timeout("30s")
    :build()

-- Task 2: IO-intensive task
local io_task = task("io_intensive")
    :description("IO-intensive file operations")
    :command(function()
        log.info("📁 Starting IO-intensive operations...")
        
        local start_time = os.time()
        
        -- Simulate file operations
        local files_created = {}
        for i = 1, 5 do
            local filename = "/tmp/sloth_test_" .. i .. ".txt"
            local content = "Test data " .. i .. " created at " .. os.date()
            
            local success, err = fs.write(filename, content)
            if success then
                table.insert(files_created, filename)
                log.info("📄 Created: " .. filename)
            else
                log.error("❌ Failed to create " .. filename .. ": " .. err)
            end
        end
        
        local duration = os.time() - start_time
        log.info("💾 IO operations completed in " .. duration .. " seconds")
        
        return true, "IO operations completed", {
            files_created = files_created,
            duration_seconds = duration,
            file_count = #files_created
        }
    end)
    :timeout("45s")
    :artifacts({"*.txt"})
    :build()

-- Task 3: Network task
local network_task = task("network_operations")
    :description("Network operations with circuit breaker")
    :command(function()
        log.info("🌐 Starting network operations...")
        
        -- Use circuit breaker for external calls
        local results = {}
        
        local api_result = circuit.protect("http_api", function()
            -- Simulate HTTP call
            log.info("📡 Making HTTP request...")
            return {
                status = 200,
                data = {message = "API call successful", timestamp = os.time()},
                response_time = math.random(100, 500)
            }
        end)
        
        table.insert(results, api_result)
        log.info("✅ Network operations completed")
        
        return true, "Network operations successful", {
            api_calls = #results,
            total_response_time = api_result.response_time,
            status = "success"
        }
    end)
    :timeout("60s")
    :retries(3, "exponential")
    :build()

-- Task 4: Parallel aggregator
local aggregator_task = task("aggregate_results")
    :description("Aggregates results from parallel tasks")
    :depends_on({"cpu_intensive", "io_intensive", "network_operations"})
    :command(function(params, deps)
        log.info("📊 Aggregating results from parallel tasks...")
        
        -- Collect all results
        local aggregate = {
            cpu_computation = deps.cpu_intensive.computation_result,
            cpu_duration = deps.cpu_intensive.duration_ms,
            files_created = deps.io_intensive.file_count,
            io_duration = deps.io_intensive.duration_seconds,
            network_calls = deps.network_operations.api_calls,
            network_time = deps.network_operations.total_response_time,
            aggregation_timestamp = os.time()
        }
        
        -- Calculate total processing time
        local total_time = (deps.cpu_intensive.duration_ms or 0) + 
                          ((deps.io_intensive.duration_seconds or 0) * 1000) + 
                          (deps.network_operations.total_response_time or 0)
        
        aggregate.total_processing_time_ms = total_time
        
        log.info("📈 Aggregation completed:")
        log.info("  Total processing time: " .. total_time .. "ms")
        log.info("  Files created: " .. aggregate.files_created)
        log.info("  Network calls: " .. aggregate.network_calls)
        
        return true, "Aggregation successful", aggregate
    end)
    :build()

-- Task 5: Cleanup task
local cleanup_task = task("cleanup")
    :description("Cleanup temporary files")
    :depends_on({"aggregate_results"})
    :command(function(params, deps)
        log.info("🧹 Starting cleanup operations...")
        
        -- Clean up test files
        local cleaned_files = 0
        local success, output = exec.run("rm -f /tmp/sloth_test_*.txt")
        if success then
            cleaned_files = deps.io_intensive.file_count or 0
            log.info("✅ Cleaned " .. cleaned_files .. " temporary files")
        else
            log.warn("⚠️  Cleanup had issues: " .. output)
        end
        
        return true, "Cleanup completed", {
            files_cleaned = cleaned_files,
            cleanup_time = os.time()
        }
    end)
    :build()

-- Modern Workflow with Parallel Execution
workflow.define("parallel_execution_demo", {
    description = "Advanced parallel execution demonstration - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"parallel", "performance", "async", "modern-dsl"},
        complexity = "intermediate",
        estimated_duration = "2m"
    },
    
    tasks = {
        cpu_task,
        io_task,
        network_task,
        aggregator_task,
        cleanup_task
    },
    
    config = {
        timeout = "10m",
        retry_policy = "exponential",
        max_parallel_tasks = 3, -- Allow 3 tasks to run in parallel
        fail_fast = false, -- Continue even if one parallel task fails
        circuit_breaker = {
            failure_threshold = 5,
            recovery_timeout = "30s"
        }
    },
    
    on_start = function()
        log.info("🚀 Starting parallel execution demonstration...")
        log.info("⚡ Tasks will run in parallel where possible")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Parallel execution demo completed successfully!")
            log.info("📊 Performance metrics collected")
            
            -- Display performance summary
            if results.aggregate_results then
                local total_time = results.aggregate_results.total_processing_time_ms
                log.info("⏱️  Total processing time: " .. total_time .. "ms")
            end
        else
            log.error("❌ Parallel execution demo failed!")
            log.warn("🔍 Check individual task results for details")
        end
        return true
    end
})
