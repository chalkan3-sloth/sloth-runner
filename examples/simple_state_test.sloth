-- MODERN DSL ONLY - Simple State Test
-- Basic state management operations for beginners

-- Task 1: Set basic state
local set_state_task = task("set_state")
    :description("Set basic application state")
    :command(function()
        log.info("📝 Setting basic state...")
        
        -- Simple state operations
        local user_data = {
            name = "Test User",
            age = 30,
            role = "developer",
            created_at = os.time()
        }
        
        local success, err = state.set("user_profile", user_data)
        if not success then
            log.error("❌ Failed to set state: " .. err)
            return false, "State set failed"
        end
        
        log.info("✅ User profile state set successfully")
        return true, "State set completed", {
            state_key = "user_profile",
            user_name = user_data.name
        }
    end)
    :timeout("10s")
    :build()

-- Task 2: Read and modify state
local modify_state_task = task("modify_state")
    :description("Read and modify existing state")
    :depends_on({"set_state"})
    :command(function(params, deps)
        log.info("🔄 Reading and modifying state...")
        
        -- Get existing state
        local user_data, err = state.get("user_profile")
        if not user_data then
            log.error("❌ Failed to get state: " .. err)
            return false, "State read failed"
        end
        
        log.info("👤 Current user: " .. user_data.name)
        
        -- Modify state
        user_data.last_login = os.time()
        user_data.login_count = (user_data.login_count or 0) + 1
        user_data.modified_at = os.time()
        
        local success, set_err = state.set("user_profile", user_data)
        if not success then
            log.error("❌ Failed to update state: " .. set_err)
            return false, "State update failed"
        end
        
        log.info("✅ User state updated - Login count: " .. user_data.login_count)
        return true, "State modified", {
            login_count = user_data.login_count,
            last_login = user_data.last_login
        }
    end)
    :build()

-- Task 3: Test state persistence
local test_persistence_task = task("test_persistence")
    :description("Test state persistence across operations")
    :depends_on({"modify_state"})
    :command(function(params, deps)
        log.info("🔍 Testing state persistence...")
        
        -- Read state again to verify persistence
        local user_data, err = state.get("user_profile")
        if not user_data then
            log.error("❌ State not found - persistence failed: " .. err)
            return false, "Persistence test failed"
        end
        
        -- Verify data integrity
        local checks = {
            has_name = user_data.name ~= nil,
            has_age = user_data.age ~= nil,
            has_login_count = user_data.login_count ~= nil,
            has_timestamps = user_data.created_at ~= nil and user_data.modified_at ~= nil
        }
        
        local passed_checks = 0
        for check, result in pairs(checks) do
            if result then 
                passed_checks = passed_checks + 1
                log.info("✅ " .. check .. ": passed")
            else
                log.error("❌ " .. check .. ": failed")
            end
        end
        
        local persistence_score = (passed_checks / 4) * 100
        
        log.info("📊 Persistence test completed:")
        log.info("  Score: " .. persistence_score .. "%")
        log.info("  User: " .. user_data.name)
        log.info("  Logins: " .. (user_data.login_count or 0))
        
        return true, "Persistence test completed", {
            persistence_score = persistence_score,
            checks_passed = passed_checks,
            user_data = user_data
        }
    end)
    :build()

-- Task 4: Cleanup state
local cleanup_state_task = task("cleanup_state")
    :description("Clean up test state")
    :depends_on({"test_persistence"})
    :command(function(params, deps)
        log.info("🧹 Cleaning up test state...")
        
        -- Optional cleanup based on parameter
        if params.cleanup ~= "false" then
            local success, err = state.delete("user_profile")
            if success then
                log.info("✅ Test state cleaned up successfully")
                return true, "Cleanup completed", {
                    cleanup_performed = true,
                    cleaned_key = "user_profile"
                }
            else
                log.warn("⚠️  Cleanup failed: " .. err)
                return true, "Cleanup attempted", {
                    cleanup_performed = false,
                    error = err
                }
            end
        else
            log.info("ℹ️  Cleanup skipped (cleanup=false)")
            return true, "Cleanup skipped", {
                cleanup_performed = false,
                reason = "skipped_by_parameter"
            }
        end
    end)
    :build()

-- Simple State Test Workflow
workflow.define("simple_state_test", {
    description = "Simple state management test - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"state", "basic", "test", "modern-dsl"},
        complexity = "beginner",
        estimated_duration = "1m"
    },
    
    tasks = {
        set_state_task,
        modify_state_task,
        test_persistence_task,
        cleanup_state_task
    },
    
    config = {
        timeout = "5m",
        retry_policy = "linear",
        max_parallel_tasks = 1, -- Sequential execution for state consistency
        state_backend = "memory" -- Use memory backend for testing
    },
    
    on_start = function()
        log.info("🚀 Starting simple state test...")
        log.info("📚 This is a basic state management demonstration")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Simple state test completed successfully!")
            
            if results.test_persistence then
                log.info("📈 Final persistence score: " .. 
                        results.test_persistence.persistence_score .. "%")
            end
            
            if results.cleanup_state and results.cleanup_state.cleanup_performed then
                log.info("🧹 Cleanup completed")
            end
        else
            log.error("❌ Simple state test failed!")
            log.warn("🔧 Check state backend availability")
        end
        return true
    end
})
