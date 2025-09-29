-- MODERN DSL ONLY: State Module Testing
-- Legacy TaskDefinitions removed - Modern DSL syntax only

-- Task 1: Basic state operations test
local state_test_task = task("test_basic_state")
    :description("Test basic state operations with modern DSL")
    :command(function()
        log.info("🧪 Testing state module...")
        
        -- Enhanced state testing with error handling
        local success, err = pcall(function()
            -- Test basic set/get
            state.set("test_key", "test_value")
            local value = state.get("test_key")
            
            if value == "test_value" then
                log.info("✅ Basic set/get works: " .. value)
            else
                error("Basic set/get failed")
            end
            
            -- Test numeric increment
            local counter = state.increment("test_counter", 5)
            log.info("✅ Counter incremented to: " .. counter)
            
            -- Test advanced features
            state.set_with_ttl("temp_key", "temp_value", 60) -- 60 seconds TTL
            log.info("✅ TTL key set successfully")
            
            -- Test stats
            local stats = state.stats()
            log.info("✅ State stats: total_keys=" .. (stats.total_keys or 0))
            
            return true
        end)
        
        if not success then
            log.error("❌ State test failed: " .. err)
            return false, "State module test failed: " .. err
        end
        
        return true, "State module working correctly!", {
            test_passed = true,
            features_tested = {"set", "get", "increment", "ttl", "stats"},
            test_timestamp = os.time()
        }
    end)
    :timeout("30s")
    :on_success(function(params, output)
        log.info("🎉 All state tests passed successfully!")
        log.info("📊 Features tested: " .. table.concat(output.features_tested, ", "))
    end)
    :build()

-- Task 2: Advanced state operations test
local advanced_state_task = task("test_advanced_state")
    :description("Test advanced state operations")
    :depends_on({"test_basic_state"})
    :command(function(params, deps)
        log.info("🔬 Testing advanced state operations...")
        
        -- Test complex data structures
        local complex_data = {
            name = "test_object",
            values = {1, 2, 3, 4, 5},
            metadata = {
                created_at = os.time(),
                version = "1.0.0"
            }
        }
        
        state.set("complex_object", complex_data)
        local retrieved = state.get("complex_object")
        
        if retrieved and retrieved.name == "test_object" then
            log.info("✅ Complex object storage works")
        else
            return false, "Complex object test failed"
        end
        
        -- Test concurrent operations
        for i = 1, 10 do
            state.increment("concurrent_counter", 1)
        end
        
        local final_count = state.get("concurrent_counter")
        log.info("✅ Concurrent counter final value: " .. (final_count or 0))
        
        return true, "Advanced state tests completed", {
            complex_object_test = "passed",
            concurrent_operations = 10,
            final_counter = final_count
        }
    end)
    :build()

-- Task 3: State cleanup test
local cleanup_task = task("cleanup_state")
    :description("Clean up test state data")
    :depends_on({"test_advanced_state"})
    :command(function(params, deps)
        log.info("🧹 Cleaning up test state data...")
        
        -- Remove test keys
        local cleanup_keys = {
            "test_key",
            "test_counter", 
            "temp_key",
            "complex_object",
            "concurrent_counter"
        }
        
        local cleaned_count = 0
        for _, key in ipairs(cleanup_keys) do
            if state.delete(key) then
                cleaned_count = cleaned_count + 1
            end
        end
        
        log.info("✅ Cleaned up " .. cleaned_count .. " test keys")
        
        return true, "State cleanup completed", {
            keys_cleaned = cleaned_count,
            cleanup_list = cleanup_keys
        }
    end)
    :build()

-- Define state testing workflow
workflow.define("state_module_test", {
    description = "Comprehensive state module testing - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        category = "testing",
        tags = {"state", "persistence", "modern-dsl"},
        author = "Sloth Runner Team"
    },
    
    tasks = {
        state_test_task,
        advanced_state_task,
        cleanup_task
    },
    
    config = {
        timeout = "5m",
        retry_policy = "linear",
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting comprehensive state module tests...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 All state tests completed successfully!")
            log.info("📊 Test summary:")
            for task_name, result in pairs(results) do
                if result.keys_cleaned then
                    log.info("  Cleanup: " .. result.keys_cleaned .. " keys")
                elseif result.features_tested then
                    log.info("  Features: " .. #result.features_tested)
                end
            end
        else
            log.error("❌ Some state tests failed!")
        end
        return true
    end
})
