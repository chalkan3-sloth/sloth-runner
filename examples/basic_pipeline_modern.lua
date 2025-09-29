-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:30 -03

local config = utils.config("pipeline_config", "development")
local fetch_task = task("fetch_data")
local process_task = task("process_data")
local store_task = task("store_result")

local fetch_task = task("fetch_data")
local process_task = task("process_data")
local store_task = task("store_result")
local fetch_task = task("fetch_data")
    :description("Simulates fetching raw data")
    :command(function(params)
        log.info("Modern DSL: Executing fetch_data...")
        
        -- Use performance monitoring (new feature)
        local result, duration = perf.measure(function()
            -- Simulate API call with circuit breaker protection
            return circuit.protect("external_api", function()
                -- Simulate success
                return {
                    success = true,
                    data = { raw_data = "some_data_from_api", source = "external_api" }
                }
            end)
        end)
        
        if result.success then
            return true, "echo 'Fetched raw data'", result.data
        else
            return false, "Failed to fetch data"
        end
    end)
    :post_hook(function(params, output)
        log.info("Modern DSL Hook: fetch_data completed. Raw data: " .. (output.raw_data or "N/A"))
        return true, "fetch_data post_exec successful"
    end)
    :timeout("30s")
    :retries(3, "exponential")
    :build()
local process_task = task("process_data")
    :description("Processes the raw data with modern features")
    :depends_on({"fetch_data"})
    :command(function(params, deps)
        local raw_data = deps.fetch_data.raw_data
        log.info("Modern DSL: Executing process_data with input: " .. raw_data)
        
        -- Enhanced conditional logic
        if raw_data == "invalid_data" then
            return false, "Invalid data received for processing"
        end
        
        -- Use async processing for better performance
        local result = async.timeout("5m", function()
            return {
                processed_data = "processed_" .. raw_data,
                status = "success",
                timestamp = os.time()
            }
        end)
        
        return true, "echo 'Processed data'", result
    end)
    :pre_hook(function(params, deps)
        log.info("Modern DSL Hook: process_data preparing. Input source: " .. (deps.fetch_data.source or "unknown"))
        return true, "process_data pre_exec successful"
    end)
    :condition(when("deps.fetch_data.success"))
    :build()
local store_task = task("store_result")
    :description("Stores the final processed data with modern features")
    :depends_on({"process_data"})
    :command(function(params, deps)
        local final_data = deps.process_data.processed_data
        log.info("Modern DSL: Executing store_result with final data: " .. final_data)
        
        -- Use modern file operations with error handling
        local success, err = fs.write("result.json", data.to_json({
            final_result = final_data,
            timestamp = os.time(),
            pipeline_version = "2.0.0"
        }))
        
        if err then
            return false, "Failed to store result: " .. err
        end
        
        return true, "echo 'Result stored'", {
            final_result = final_data,
            timestamp = os.time(),
            file_path = "result.json"
        }
    end)
    :artifacts({"result.json"})
    :on_success(function(params, output)
        log.info("Pipeline completed successfully! Result stored at: " .. output.file_path)
    end)
    :build()

workflow.define("basic_pipeline_modern", {
    description = "A simple data processing pipeline - Modern DSL Version",
    version = "2.0.0",
    
    -- Enhanced metadata
    metadata = {
        author = "Sloth Runner Team",
        created_at = os.date(),
        tags = {"data-processing", "basic", "modern-dsl"}
    },
    
    -- Modern task orchestration
    tasks = {
        fetch_task,
        process_task,
        store_task
    },
    
    -- Workflow-level configuration
    config = {
        max_parallel_tasks = 2,
        timeout = "30m",
        retry_policy = "exponential",
        cleanup_on_failure = true
    },
    
    -- Workflow hooks
    on_start = function()
        log.info("Starting modern basic pipeline...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("Modern basic pipeline completed successfully!")
        else
            log.error("Modern basic pipeline failed!")
        end
        return true
    end
})
