-- MODERN DSL ONLY - Workflow description - Modern DSL
-- Converted from legacy Modern DSL format
-- Category: intermediate, Complexity: moderate

-- Main task using Modern DSL
local main_task = task("parallel-processing_task")
    :description("Workflow description - Modern DSL")
    :command(function(params, deps)
        log.info("🚀 Executing parallel-processing with Modern DSL...")
        
        -- TODO: Replace with actual implementation from original file
        -- Original logic should be migrated here
        
        return true, "Task completed successfully", {
            task_name = "parallel-processing",
            execution_time = os.time(),
            status = "success"
        }
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :on_success(function(params, output)
        log.info("✅ parallel-processing task completed successfully")
    end)
    :on_failure(function(params, error)
        log.error("❌ parallel-processing task failed: " .. error)
    end)
    :build()

-- Additional tasks can be added here following the same pattern
-- local secondary_task = task("parallel-processing_secondary")
--     :description("Secondary task for parallel-processing")
--     :depends_on({"parallel-processing_task"})
--     :command(function(params, deps)
--         -- Secondary logic here
--         return true, "Secondary task completed", {}
--     end)
--     :build()

-- Modern Workflow Definition
workflow.define("parallel-processing_workflow", {
    description = "Workflow description - Modern DSL - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        category = "intermediate",
        complexity = "moderate",
        tags = {"parallel-processing", "modern-dsl", "intermediate"},
        created_at = os.date(),
        migrated_from = "Modern DSL format"
    },
    
    tasks = {
        main_task
        -- Add additional tasks here
    },
    
    config = {
        timeout = "15m",
        retry_policy = "exponential",
        max_parallel_tasks = 2,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting parallel-processing workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ parallel-processing workflow completed successfully!")
        else
            log.error("❌ parallel-processing workflow failed!")
        end
        return true
    end
})

-- Migration Note:
-- This file has been converted from legacy Modern DSL format to Modern DSL
-- TODO: Review and implement the original logic in the Modern DSL structure above
-- Original backup saved as parallel-processing.lua.backup
