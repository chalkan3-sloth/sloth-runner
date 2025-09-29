-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:30 -03

local parallel_demo_task = task("run_in_parallel")
local comparison_task = task("sequential_vs_parallel")

local parallel_demo_task = task("run_in_parallel")
local comparison_task = task("sequential_vs_parallel")
local parallel_demo_task = task("run_in_parallel")
    :description("Demonstrates parallel execution with modern DSL")
    :command(function()
        log.info("🚀 Modern DSL: Starting parallel execution demo")

        -- Define sub-tasks for parallel execution using modern async patterns
        local start_time = os.time()
        log.info("Executing 3 tasks in parallel (2s, 4s, 6s sleep times)...")

        -- Use modern async.parallel with enhanced features
        local results, err = async.parallel({
            short_task = function() 
                log.info("Short task starting...")
                local success, output = exec.run("echo 'Sub-task 1 starting...'; sleep 2; echo 'Sub-task 1 finished.'")
                return { name = "Short sleep", duration = 2, success = success, output = output }
            end,
            medium_task = function() 
                log.info("Medium task starting...")
                local success, output = exec.run("echo 'Sub-task 2 starting...'; sleep 4; echo 'Sub-task 2 finished.'")
                return { name = "Medium sleep", duration = 4, success = success, output = output }
            end,
            long_task = function() 
                log.info("Long task starting...")
                local success, output = exec.run("echo 'Sub-task 3 starting...'; sleep 6; echo 'Sub-task 3 finished.'")
                return { name = "Long sleep", duration = 6, success = success, output = output }
            end
        }, {
            max_workers = 3,
            timeout = "10s",
            fail_fast = false
        })

        local end_time = os.time()
        local duration = end_time - start_time
        
        log.info("✅ Parallel execution completed in " .. duration .. " seconds")
        log.info("📊 Expected time: ~6 seconds (duration of longest task)")

        if err then
            log.error("❌ Parallel execution failed: " .. err)
            return false, "Parallel execution failed"
        end

        -- Enhanced result processing
        log.info("📋 Results from parallel execution:")
        local success_count = 0
        for task_name, result in pairs(results or {}) do
            if result.success then
                success_count = success_count + 1
                log.info("  ✅ " .. result.name .. " - SUCCESS")
            else
                log.error("  ❌ " .. result.name .. " - FAILED")
            end
        end

        local efficiency = math.floor((6 / duration) * 100)
        
        return true, "Parallel execution demo completed successfully", {
            total_duration = duration,
            tasks_executed = 3,
            successful_tasks = success_count,
            parallel_efficiency = efficiency .. "%",
            performance_score = efficiency >= 90 and "Excellent" or efficiency >= 70 and "Good" or "Needs Improvement"
        }
    end)
    :timeout("15s")
    :retries(2, "linear")
    :on_success(function(params, output)
        log.info("🎉 Parallel demo completed!")
        log.info("📈 Efficiency: " .. output.parallel_efficiency)
        log.info("🏆 Performance: " .. output.performance_score)
    end)
    :on_failure(function(params, error)
        log.error("🚨 Parallel demo failed: " .. error)
    end)
    :build()
local comparison_task = task("sequential_vs_parallel")
    :description("Compare sequential vs parallel execution")
    :depends_on({"run_in_parallel"})
    :command(function(params, deps)
        log.info("📊 Comparison Results:")
        log.info("  ⏱️  Parallel execution time: " .. deps.run_in_parallel.total_duration .. "s")
        log.info("  📈 Parallel efficiency: " .. deps.run_in_parallel.parallel_efficiency)
        log.info("  🔄 Sequential would take: ~12s (2+4+6)")
        
        local time_saved = 12 - deps.run_in_parallel.total_duration
        log.info("  💰 Time saved: " .. time_saved .. "s")
        
        return true, "Comparison completed", {
            parallel_time = deps.run_in_parallel.total_duration,
            sequential_time = 12,
            time_saved = time_saved,
            efficiency_gain = math.floor((time_saved / 12) * 100) .. "%"
        }
    end)
    :build()

workflow.define("parallel_execution_demo", {
    description = "Parallel execution demonstration - Modern DSL Only",
    version = "2.0.0",
    
    metadata = {
        category = "demonstration",
        tags = {"parallel", "async", "performance", "modern-dsl"},
        author = "Sloth Runner Team",
        complexity = "intermediate"
    },
    
    tasks = {
        parallel_demo_task,
        comparison_task
    },
    
    config = {
        timeout = "20m",
        retry_policy = "exponential",
        max_parallel_tasks = 2,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🎬 Starting parallel execution demonstration...")
        log.info("💡 This demo shows the power of async parallel processing")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Parallel execution demo completed successfully!")
            log.info("📊 Results summary:")
            for task_name, result in pairs(results) do
                if result.efficiency_gain then
                    log.info("  💫 Efficiency gain: " .. result.efficiency_gain)
                end
            end
        else
            log.error("❌ Parallel execution demo failed!")
        end
        return true
    end
})
