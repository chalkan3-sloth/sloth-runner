-- Test AI Module Loading
local ai = require("ai")
local log = require("log")

log.info("🤖 Testing AI module...")

-- Test AI configuration
local success, err = pcall(function()
    ai.configure({
        enabled = true,
        learning_mode = "adaptive",
        optimization_level = 5
    })
end)

if success then
    log.info("✅ AI configuration successful!")
    
    -- Get configuration
    local config = ai.get_config()
    log.info("📊 AI Config - Enabled: " .. tostring(config.enabled))
    log.info("📊 AI Config - Mode: " .. config.learning_mode)
    log.info("📊 AI Config - Level: " .. config.optimization_level)
else
    log.error("❌ AI configuration failed: " .. tostring(err))
end

-- Test optimization
local opt_success, opt_err = pcall(function()
    local suggestion = ai.optimize_command("go build", {
        system_resources = {
            cpu_usage = 30,
            memory_usage = 50
        }
    })
    
    if suggestion then
        log.info("🚀 AI Optimization suggestion received!")
        log.info("📈 Confidence: " .. string.format("%.1f%%", suggestion.confidence_score * 100))
        log.info("⚡ Expected speedup: " .. string.format("%.1fx", suggestion.expected_speedup))
        log.info("💡 Optimized command: " .. suggestion.optimized_command)
    else
        log.info("ℹ️ No optimization suggestions available")
    end
end)

if not opt_success then
    log.error("❌ AI optimization test failed: " .. tostring(opt_err))
end

log.info("🏁 AI module test completed!")

-- Define a simple workflow that uses the AI module
workflow.define("ai_test", {
    description = "Test AI module functionality",
    version = "1.0.0",
    
    tasks = {
        {
            name = "ai_test_task",
            description = "Test task to validate AI functionality",
            command = function(params, deps)
                log.info("🧪 Running AI test task...")
                
                -- Record a test execution
                ai.record_execution({
                    task_name = "ai_test_task",
                    command = "echo 'AI test'",
                    success = true,
                    execution_time = "100ms"
                })
                
                log.info("✅ AI test task completed!")
                return {success = true, output = "AI test successful"}
            end
        }
    },
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 AI test workflow completed successfully!")
        else
            log.error("💥 AI test workflow failed!")
        end
    end
})