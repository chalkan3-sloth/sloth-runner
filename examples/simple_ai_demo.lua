-- Simple AI-Powered Task Example
-- Demonstrates basic AI optimization as shown in the original request

local ai = require("ai")
local exec = require("exec")
local log = require("log")

-- Simple task with AI optimization
local build_with_ai = task("build_with_ai")
    :description("Build with AI optimization")
    :command(function(params, deps)
        -- AI suggests optimizations based on historical data
        local ai_suggestions = ai.optimize_command("go build", {
            history = ai.get_task_history("go build"),
            system_resources = {
                cpu_usage = 45,
                memory_usage = 60
            },
            similar_tasks = ai.find_similar_tasks("go build", 10)
        })
        
        if ai_suggestions and ai_suggestions.confidence_score > 0.6 then
            log.info("🤖 AI optimized command: " .. ai_suggestions.optimized_command)
            log.info("📈 Expected speedup: " .. string.format("%.1fx", ai_suggestions.expected_speedup))
            
            return exec.run(ai_suggestions.optimized_command)
        else
            return exec.run("go build")
        end
    end)
    :build()

-- Task with failure prediction
local deploy_with_prediction = task("deploy_with_prediction")
    :command(function(params, deps)
        local command = "kubectl apply -f deployment.yaml"
        
        -- Predict failure probability
        local prediction = ai.predict_failure("deploy_with_prediction", command)
        
        log.info("🔮 Failure probability: " .. string.format("%.1f%%", prediction.failure_probability * 100))
        
        if prediction.failure_probability > 0.3 then
            log.warn("⚠️ High failure risk detected!")
            for _, rec in ipairs(prediction.recommendations) do
                log.info("💡 " .. rec)
            end
        end
        
        return exec.run(command)
    end)
    :build()

-- Workflow with AI features
workflow.define("ai_demo", {
    description = "Simple AI demonstration",
    tasks = { build_with_ai, deploy_with_prediction },
    
    on_task_complete = function(task_name, success, output)
        -- Record execution for AI learning
        ai.record_execution({
            task_name = task_name,
            command = output.command or "unknown",
            success = success,
            execution_time = output.duration or "0s"
        })
    end
})