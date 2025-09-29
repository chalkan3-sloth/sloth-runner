-- CONVERTED TO MODERN DSL
-- Legacy TaskDefinitions format has been completely removed
-- This file now uses only Modern DSL syntax

-- Example task using Modern DSL:
local converted_task = task("converted_task")
    :description("Converted from legacy TaskDefinitions")
    :command(function(params, deps)
        log.info("Modern DSL: Task converted from legacy format")
        -- Add your specific task logic here from the backup file
        return true, "Task completed", {}
    end)
    :timeout("30s")
    :build()

-- Modern workflow definition:
workflow.define("converted_workflow", {
    description = "Converted from legacy TaskDefinitions format",
    version = "2.0.0",
    
    metadata = {
        tags = {"converted", "modern-dsl", "legacy-migration"},
        migration_date = os.date()
    },
    
    tasks = { converted_task },
    
    on_start = function()
        log.info("Starting converted workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("Converted workflow completed successfully!")
        else
            log.error("Converted workflow failed!")
        end
        return true
    end
})

-- NOTE: Original legacy code is preserved in .pre_modern_backup file
-- Please review and migrate specific tasks as needed
