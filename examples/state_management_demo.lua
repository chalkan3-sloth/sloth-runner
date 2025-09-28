-- State Management Example for sloth-runner
-- This demonstrates the powerful state persistence capabilities

TaskDefinitions = {
    state_management_demo = {
        description = "Comprehensive demonstration of state management features",
        tasks = {
            -- Task 1: Basic state operations
            basic_state_operations = {
                name = "basic_state_operations",
                description = "Demonstrates basic set, get, delete operations",
                command = function()
                    log.info("=== Basic State Operations ===")
                    
                    -- Set various data types
                    state.set("app_version", "v1.2.3")
                    state.set("deployment_count", 42)
                    state.set("feature_flags", {
                        new_ui = true,
                        beta_features = false,
                        debug_mode = true
                    })
                    
                    -- Get values
                    local version = state.get("app_version")
                    local count = state.get("deployment_count")
                    local flags = state.get("feature_flags")
                    
                    log.info("App Version: " .. version)
                    log.info("Deployment Count: " .. count)
                    log.info("Feature Flags: " .. data.to_json(flags))
                    
                    -- Check existence
                    if state.exists("app_version") then
                        log.info("app_version key exists")
                    end
                    
                    -- Get with default
                    local missing = state.get("missing_key", "default_value")
                    log.info("Missing key with default: " .. missing)
                    
                    return true, "Basic operations completed successfully"
                end,
            },
            
            -- Task 2: TTL and expiration
            ttl_operations = {
                name = "ttl_operations",
                description = "Demonstrates TTL and key expiration",
                depends_on = "basic_state_operations",
                command = function()
                    log.info("=== TTL Operations ===")
                    
                    -- Set with TTL (5 seconds)
                    state.set("temp_token", "abc123xyz", 5)
                    log.info("Set temp token with 5 second TTL")
                    
                    -- Check TTL
                    local ttl = state.get_ttl("temp_token")
                    log.info("Token TTL: " .. ttl .. " seconds")
                    
                    -- Set TTL for existing key
                    state.set_ttl("deployment_count", 10)
                    log.info("Set TTL for deployment_count to 10 seconds")
                    
                    return true, "TTL operations completed"
                end,
            },
            
            -- Task 3: Atomic operations
            atomic_operations = {
                name = "atomic_operations",
                description = "Demonstrates atomic increment, decrement, append",
                depends_on = "ttl_operations",
                command = function()
                    log.info("=== Atomic Operations ===")
                    
                    -- Atomic increment
                    local counter = state.increment("request_counter", 1)
                    log.info("Request counter: " .. counter)
                    
                    counter = state.increment("request_counter", 5)
                    log.info("Request counter after +5: " .. counter)
                    
                    -- Atomic decrement
                    counter = state.decrement("request_counter", 2)
                    log.info("Request counter after -2: " .. counter)
                    
                    -- String append
                    state.set("log_messages", "Starting deployment")
                    local length = state.append("log_messages", " -> Downloading images")
                    log.info("Log message length: " .. length)
                    
                    length = state.append("log_messages", " -> Configuring services")
                    log.info("Log message length: " .. length)
                    
                    local messages = state.get("log_messages")
                    log.info("Full log: " .. messages)
                    
                    return true, "Atomic operations completed"
                end,
            },
            
            -- Task 4: List operations
            list_operations = {
                name = "list_operations",
                description = "Demonstrates list push, pop, length operations",
                depends_on = "atomic_operations",
                command = function()
                    log.info("=== List Operations ===")
                    
                    -- Push items to list
                    state.list_push("deployment_history", {
                        version = "v1.0.0",
                        timestamp = os.time(),
                        success = true
                    })
                    
                    state.list_push("deployment_history", {
                        version = "v1.1.0",
                        timestamp = os.time(),
                        success = true
                    })
                    
                    state.list_push("deployment_history", {
                        version = "v1.2.0",
                        timestamp = os.time(),
                        success = false,
                        error = "Database connection failed"
                    })
                    
                    -- Check list length
                    local length = state.list_length("deployment_history")
                    log.info("Deployment history length: " .. length)
                    
                    -- Pop last item
                    local last_deployment = state.list_pop("deployment_history")
                    if last_deployment then
                        log.info("Last deployment: " .. data.to_json(last_deployment))
                    end
                    
                    local new_length = state.list_length("deployment_history")
                    log.info("New deployment history length: " .. new_length)
                    
                    return true, "List operations completed"
                end,
            },
            
            -- Task 5: Distributed locking
            locking_demo = {
                name = "locking_demo",
                description = "Demonstrates distributed locking capabilities",
                depends_on = "list_operations",
                command = function()
                    log.info("=== Distributed Locking Demo ===")
                    
                    -- Try to acquire a lock
                    local lock_acquired = state.try_lock("deployment_lock", 30)
                    if lock_acquired then
                        log.info("Successfully acquired deployment lock")
                        
                        -- Simulate some work
                        exec.run("sleep 2")
                        
                        -- Release the lock
                        local released = state.unlock("deployment_lock")
                        if released then
                            log.info("Lock released successfully")
                        else
                            log.error("Failed to release lock")
                        end
                    else
                        log.error("Failed to acquire deployment lock")
                    end
                    
                    return true, "Locking demo completed"
                end,
            },
            
            -- Task 6: Critical section with automatic lock management
            critical_section_demo = {
                name = "critical_section_demo",
                description = "Demonstrates automatic lock management with critical sections",
                depends_on = "locking_demo",
                command = function()
                    log.info("=== Critical Section Demo ===")
                    
                    -- Execute critical section with automatic locking
                    state.with_lock("critical_resource", function()
                        log.info("Inside critical section - lock is held")
                        
                        -- Simulate critical work
                        local config = state.get("system_config", {})
                        config.last_maintenance = os.time()
                        config.maintenance_count = (config.maintenance_count or 0) + 1
                        
                        state.set("system_config", config)
                        
                        log.info("Critical work completed - maintenance count: " .. config.maintenance_count)
                        
                        -- Lock will be automatically released when function returns
                    end, 10) -- 10 second timeout
                    
                    return true, "Critical section demo completed"
                end,
            },
            
            -- Task 7: Compare and swap
            compare_swap_demo = {
                name = "compare_swap_demo",
                description = "Demonstrates atomic compare-and-swap operations",
                depends_on = "critical_section_demo",
                command = function()
                    log.info("=== Compare and Swap Demo ===")
                    
                    -- Set initial value
                    state.set("cas_counter", 10)
                    
                    -- Successful compare-and-swap
                    local success = state.compare_swap("cas_counter", 10, 20)
                    if success then
                        log.info("CAS successful: 10 -> 20")
                    else
                        log.error("CAS failed")
                    end
                    
                    -- Failed compare-and-swap (wrong expected value)
                    success = state.compare_swap("cas_counter", 10, 30)
                    if success then
                        log.info("CAS successful: 10 -> 30")
                    else
                        log.info("CAS failed as expected (value was not 10)")
                    end
                    
                    local current = state.get("cas_counter")
                    log.info("Current CAS counter value: " .. current)
                    
                    return true, "Compare and swap demo completed"
                end,
            },
            
            -- Task 8: State inspection and management
            state_inspection = {
                name = "state_inspection",
                description = "Demonstrates state inspection and management features",
                depends_on = "compare_swap_demo",
                command = function()
                    log.info("=== State Inspection ===")
                    
                    -- List all keys
                    local all_keys = state.keys()
                    log.info("All keys count: " .. #all_keys)
                    for i, key in ipairs(all_keys) do
                        log.info("Key " .. i .. ": " .. key)
                    end
                    
                    -- List keys with pattern
                    local app_keys = state.keys("app_*")
                    log.info("App keys: " .. data.to_json(app_keys))
                    
                    -- Get statistics
                    local stats = state.stats()
                    log.info("State statistics: " .. data.to_json(stats))
                    
                    return true, "State inspection completed"
                end,
            },
            
            -- Task 9: Cleanup demonstration
            cleanup_demo = {
                name = "cleanup_demo",
                description = "Demonstrates selective cleanup of state data",
                depends_on = "state_inspection",
                command = function()
                    log.info("=== Cleanup Demo ===")
                    
                    -- Create some temporary keys
                    state.set("temp_key1", "value1")
                    state.set("temp_key2", "value2")
                    state.set("temp_key3", "value3")
                    
                    local before_keys = state.keys("temp_*")
                    log.info("Temp keys before cleanup: " .. #before_keys)
                    
                    -- Clear temp keys
                    state.clear("temp_*")
                    
                    local after_keys = state.keys("temp_*")
                    log.info("Temp keys after cleanup: " .. #after_keys)
                    
                    -- Final statistics
                    local final_stats = state.stats()
                    log.info("Final statistics: " .. data.to_json(final_stats))
                    
                    return true, "Cleanup demo completed"
                end,
            }
        }
    }
}