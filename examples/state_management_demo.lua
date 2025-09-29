-- MODERN DSL ONLY - State Management Demonstration
-- Advanced state management with TTL, atomic operations, and persistence

-- Task 1: Initialize application state
local init_state_task = task("initialize_state")
    :description("Initialize application state with modern features")
    :command(function()
        log.info("🔧 Initializing application state...")
        
        -- Initialize with structured data
        local app_state = {
            version = "2.0.0",
            started_at = os.time(),
            environment = "development",
            features = {
                logging = true,
                metrics = true,
                circuit_breaker = true
            },
            counters = {
                requests = 0,
                errors = 0,
                successes = 0
            }
        }
        
        -- Store with TTL (Time To Live)
        local success, err = state.set("app_config", app_state, {
            ttl = 3600, -- 1 hour
            atomic = true
        })
        
        if not success then
            log.error("❌ Failed to initialize state: " .. err)
            return false, "State initialization failed"
        end
        
        log.info("✅ Application state initialized")
        return true, "State initialized", {
            state_key = "app_config",
            ttl_seconds = 3600,
            timestamp = os.time()
        }
    end)
    :timeout("30s")
    :build()

-- Task 2: Simulate application operations
local simulate_operations_task = task("simulate_operations")
    :description("Simulate application operations with state updates")
    :depends_on({"initialize_state"})
    :command(function(params, deps)
        log.info("⚡ Simulating application operations...")
        
        -- Perform atomic counter updates
        local operations = {
            {type = "request", success = true},
            {type = "request", success = false},
            {type = "request", success = true},
            {type = "request", success = true},
            {type = "request", success = false}
        }
        
        local results = {
            operations_performed = 0,
            successes = 0,
            errors = 0
        }
        
        for i, op in ipairs(operations) do
            -- Atomic increment operations
            if op.success then
                state.increment("app_config.counters.successes", 1)
                results.successes = results.successes + 1
            else
                state.increment("app_config.counters.errors", 1)
                results.errors = results.errors + 1
            end
            
            state.increment("app_config.counters.requests", 1)
            results.operations_performed = results.operations_performed + 1
            
            log.info("🔄 Operation " .. i .. ": " .. (op.success and "SUCCESS" or "ERROR"))
        end
        
        log.info("📊 Operations completed:")
        log.info("  Total: " .. results.operations_performed)
        log.info("  Successes: " .. results.successes)
        log.info("  Errors: " .. results.errors)
        
        return true, "Operations simulated", results
    end)
    :build()

-- Task 3: Advanced state queries
local query_state_task = task("query_state")
    :description("Perform advanced state queries and analysis")
    :depends_on({"simulate_operations"})
    :command(function(params, deps)
        log.info("🔍 Querying application state...")
        
        -- Get current state
        local app_state, err = state.get("app_config")
        if not app_state then
            log.error("❌ Failed to query state: " .. err)
            return false, "State query failed"
        end
        
        -- Calculate metrics
        local total_requests = app_state.counters.requests or 0
        local successes = app_state.counters.successes or 0
        local errors = app_state.counters.errors or 0
        local success_rate = total_requests > 0 and (successes / total_requests * 100) or 0
        
        local metrics = {
            total_requests = total_requests,
            successes = successes,
            errors = errors,
            success_rate = string.format("%.2f", success_rate),
            uptime_seconds = os.time() - app_state.started_at,
            state_ttl = state.ttl("app_config") or 0
        }
        
        log.info("📈 Application Metrics:")
        log.info("  Total Requests: " .. metrics.total_requests)
        log.info("  Success Rate: " .. metrics.success_rate .. "%")
        log.info("  Uptime: " .. metrics.uptime_seconds .. " seconds")
        log.info("  State TTL: " .. metrics.state_ttl .. " seconds")
        
        return true, "State queried successfully", metrics
    end)
    :build()

-- Task 4: State persistence and backup
local backup_state_task = task("backup_state")
    :description("Create persistent backup of application state")
    :depends_on({"query_state"})
    :command(function(params, deps)
        log.info("💾 Creating state backup...")
        
        -- Get all application state
        local app_state, err = state.get("app_config")
        if not app_state then
            log.error("❌ Failed to get state for backup: " .. err)
            return false, "Backup failed"
        end
        
        -- Add backup metadata
        app_state.backup_metadata = {
            backed_up_at = os.time(),
            backup_version = "1.0",
            backed_up_by = "sloth-runner",
            original_ttl = state.ttl("app_config")
        }
        
        -- Export to JSON file
        local backup_filename = "state_backup_" .. os.time() .. ".json"
        local backup_content = data.to_json(app_state)
        
        local success, write_err = fs.write(backup_filename, backup_content)
        if not success then
            log.error("❌ Failed to write backup file: " .. write_err)
            return false, "Backup file creation failed"
        end
        
        -- Also store backup in state with different key
        state.set("app_config_backup", app_state, {
            ttl = 86400, -- 24 hours
            atomic = true
        })
        
        log.info("✅ State backup created: " .. backup_filename)
        return true, "Backup completed", {
            backup_file = backup_filename,
            backup_size_bytes = string.len(backup_content),
            backup_timestamp = os.time()
        }
    end)
    :artifacts({backup_filename})
    :build()

-- Task 5: State validation and cleanup
local validate_cleanup_task = task("validate_and_cleanup")
    :description("Validate state integrity and perform cleanup")
    :depends_on({"backup_state"})
    :command(function(params, deps)
        log.info("🧹 Validating state and performing cleanup...")
        
        -- Validate state structure
        local app_state, err = state.get("app_config")
        if not app_state then
            log.error("❌ State validation failed: " .. err)
            return false, "Validation failed"
        end
        
        local validation_results = {
            has_version = app_state.version ~= nil,
            has_counters = app_state.counters ~= nil,
            has_features = app_state.features ~= nil,
            counters_valid = false,
            backup_exists = false
        }
        
        -- Validate counters
        if app_state.counters then
            validation_results.counters_valid = 
                app_state.counters.requests and 
                app_state.counters.successes and 
                app_state.counters.errors and
                app_state.counters.requests >= 0
        end
        
        -- Check backup exists
        local backup_state, _ = state.get("app_config_backup")
        validation_results.backup_exists = backup_state ~= nil
        
        -- Count validation results
        local valid_count = 0
        local total_checks = 0
        for check, result in pairs(validation_results) do
            total_checks = total_checks + 1
            if result then valid_count = valid_count + 1 end
        end
        
        local validation_score = (valid_count / total_checks) * 100
        
        log.info("✅ Validation completed:")
        log.info("  Validation Score: " .. string.format("%.1f", validation_score) .. "%")
        log.info("  Checks Passed: " .. valid_count .. "/" .. total_checks)
        
        -- Cleanup if requested
        if params.cleanup == "true" then
            log.info("🧹 Performing cleanup...")
            state.delete("app_config")
            state.delete("app_config_backup")
            log.info("✅ State cleanup completed")
        end
        
        return true, "Validation and cleanup completed", {
            validation_score = validation_score,
            checks_passed = valid_count,
            total_checks = total_checks,
            cleanup_performed = params.cleanup == "true"
        }
    end)
    :build()

-- Modern State Management Workflow
workflow.define("state_management_demo", {
    description = "Advanced state management demonstration - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"state", "persistence", "atomic", "ttl", "modern-dsl"},
        complexity = "advanced",
        estimated_duration = "3m"
    },
    
    tasks = {
        init_state_task,
        simulate_operations_task,
        query_state_task,
        backup_state_task,
        validate_cleanup_task
    },
    
    config = {
        timeout = "15m",
        retry_policy = "exponential",
        max_parallel_tasks = 2,
        state_persistence = {
            enabled = true,
            backend = "sqlite",
            cleanup_on_completion = false
        }
    },
    
    on_start = function()
        log.info("🚀 Starting state management demonstration...")
        log.info("💾 Demonstrating: TTL, atomic operations, persistence, validation")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 State management demo completed successfully!")
            log.info("📊 All state operations completed")
            
            -- Display final summary
            if results.validate_and_cleanup then
                log.info("🏆 Final validation score: " .. 
                        results.validate_and_cleanup.validation_score .. "%")
            end
        else
            log.error("❌ State management demo failed!")
            log.warn("🔧 Check state backend and permissions")
        end
        return true
    end
})
