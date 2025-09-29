-- MIGRATED TO MODERN DSL - Basic Data Processing Pipeline
-- This example demonstrates the new Modern DSL syntax

-- Task 1: Fetch data using Modern DSL builder pattern
local fetch_data = task("fetch_data")
    :description("Simulates fetching raw data")
    :command(function(params)
        log.info("Modern DSL: Executing fetch_data...")
        return true, "echo 'Fetched raw data'", { 
            raw_data = "some_data_from_api", 
            source = "external_api" 
        }
    end)
    :on_success(function(params, output)
        log.info("Modern DSL Hook: fetch_data completed. Raw data: " .. (output.raw_data or "N/A"))
        return true, "fetch_data post_exec successful"
    end)
    :timeout("30s")
    :build()

-- Task 2: Process data with enhanced dependency handling
local process_data = task("process_data")
    :description("Processes the raw data")
    :depends_on({"fetch_data"})
    :command(function(params, deps)
        local raw_data = deps.fetch_data.raw_data
        log.info("Modern DSL: Executing process_data with input: " .. raw_data)
        
        if raw_data == "invalid_data" then
            return false, "Invalid data received for processing"
        end
        
        return true, "echo 'Processed data'", { 
            processed_data = "processed_" .. raw_data, 
            status = "success" 
        }
    end)
    :pre_hook(function(params, deps)
        log.info("Modern DSL Hook: process_data preparing. Input source: " .. (deps.fetch_data.source or "unknown"))
        return true, "process_data pre_exec successful"
    end)
    :build()

-- Task 3: Store result with Modern DSL features
local store_result = task("store_result")
    :description("Stores the final processed data")
    :depends_on({"process_data"})
    :command(function(params, deps)
        local final_data = deps.process_data.processed_data
        log.info("Modern DSL: Executing store_result with final data: " .. final_data)
        return true, "echo 'Result stored'", { 
            final_result = final_data, 
            timestamp = os.time() 
        }
    end)
    :build()

-- Modern Workflow Definition
workflow.define("basic_pipeline", {
    description = "A simple data processing pipeline - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"data-processing", "basic", "modern-dsl"},
        created_at = os.date()
    },
    
    tasks = {
        fetch_data,
        process_data,
        store_result
    },
    
    config = {
        timeout = "15m",
        retry_policy = "exponential",
        max_parallel_tasks = 2
    },
    
    on_start = function()
        log.info("🚀 Starting basic data processing pipeline...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Basic pipeline completed successfully!")
        else
            log.error("❌ Basic pipeline failed!")
        end
        return true
    end
})
