-- Test workflow to demonstrate stack functionality
TaskDefinitions = {
    test_stack = {
        description = "Test workflow for stack functionality",
        tasks = {
            {
                name = "setup",
                description = "Setup test environment",
                command = "echo 'Setting up test environment...'"
            },
            {
                name = "process",
                description = "Process some data",
                depends_on = {"setup"},
                command = function()
                    print("Processing data...")
                    -- Simulate some processing
                    local result = {
                        processed_items = 42,
                        timestamp = os.time()
                    }
                    
                    -- Export outputs
                    if not outputs then
                        outputs = {}
                    end
                    outputs.processed_items = result.processed_items
                    outputs.timestamp = result.timestamp
                    outputs.status = "success"
                    
                    return true
                end
            },
            {
                name = "cleanup",
                description = "Cleanup resources",
                depends_on = {"process"},
                command = "echo 'Cleaning up resources...'"
            }
        }
    }
}