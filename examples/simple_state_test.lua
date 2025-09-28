-- Simple state test
TaskDefinitions = {
    simple_state_test = {
        description = "Simple test of state functionality",
        tasks = {
            test_basic_state = {
                name = "test_basic_state",
                description = "Test basic state operations",
                command = function()
                    log.info("Testing state module...")
                    
                    -- Test basic set/get
                    state.set("test_key", "test_value")
                    local value = state.get("test_key")
                    
                    if value == "test_value" then
                        log.info("✓ Basic set/get works: " .. value)
                    else
                        log.error("✗ Basic set/get failed")
                        return false, "Basic state test failed"
                    end
                    
                    -- Test numeric increment
                    local counter = state.increment("test_counter", 5)
                    log.info("✓ Counter incremented to: " .. counter)
                    
                    -- Test stats
                    local stats = state.stats()
                    log.info("✓ State stats: total_keys=" .. stats.total_keys)
                    
                    return true, "State module working correctly!"
                end,
            }
        }
    }
}