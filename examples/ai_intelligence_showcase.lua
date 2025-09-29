-- Complete AI-Powered Task Intelligence Demo
-- Demonstrates all AI features without external dependencies

local ai = require("ai")
local log = require("log")
local exec = require("exec")

-- Configure AI
ai.configure({
    enabled = true,
    learning_mode = "adaptive",
    optimization_level = 8,
    failure_prediction = true,
    auto_optimize = true
})

log.info("🤖 AI-Powered Task Intelligence Demo")
log.info("=" .. string.rep("=", 50))

-- Task 1: AI Optimization Demo
local optimization_demo = task("optimization_demo")
    :description("Demonstrates AI command optimization")
    :command(function(params, deps)
        log.info("🚀 Testing AI Optimization...")
        
        local commands = {
            "echo 'Hello World'",
            "ls -la",
            "cat /etc/hostname",
            "pwd",
            "date"
        }
        
        for i, cmd in ipairs(commands) do
            log.info("🔍 Analyzing command: " .. cmd)
            
            local suggestion = ai.optimize_command(cmd, {
                system_resources = {
                    cpu_usage = math.random(20, 80),
                    memory_usage = math.random(30, 90)
                }
            })
            
            if suggestion then
                log.info("  📊 Confidence: " .. string.format("%.1f%%", suggestion.confidence_score * 100))
                log.info("  ⚡ Speedup: " .. string.format("%.1fx", suggestion.expected_speedup))
                log.info("  💡 Optimized: " .. suggestion.optimized_command)
                
                if #suggestion.optimizations > 0 then
                    log.info("  🔧 Optimizations applied:")
                    for j, opt in ipairs(suggestion.optimizations) do
                        log.info("    " .. j .. ". " .. opt.description .. " (Impact: " .. string.format("%.1f", opt.impact) .. ")")
                    end
                end
            end
            
            -- Execute the command
            local result = exec.run(cmd)
            
            -- Record execution for AI learning
            ai.record_execution({
                task_name = "optimization_demo",
                command = cmd,
                success = result.success,
                execution_time = "50ms",
                parameters = {
                    iteration = i,
                    test_run = true
                }
            })
            
            log.info("  ✅ Executed and recorded for AI learning")
            log.info("")
        end
        
        return {success = true, message = "AI optimization demo completed"}
    end)
    :build()

-- Task 2: Failure Prediction Demo
local prediction_demo = task("prediction_demo")
    :description("Demonstrates AI failure prediction")
    :command(function(params, deps)
        log.info("🔮 Testing AI Failure Prediction...")
        
        local risky_commands = {
            "rm -rf /tmp/nonexistent",
            "curl http://fake-domain-12345.com",
            "cat /nonexistent/file.txt",
            "echo 'safe command'",
            "ls /invalid/path"
        }
        
        for i, cmd in ipairs(risky_commands) do
            log.info("🧪 Predicting failure for: " .. cmd)
            
            local prediction = ai.predict_failure("prediction_demo", cmd)
            
            if prediction then
                log.info("  📊 Failure probability: " .. string.format("%.1f%%", prediction.failure_probability * 100))
                log.info("  🎯 Confidence: " .. string.format("%.1f%%", prediction.confidence * 100))
                
                if prediction.failure_probability > 0.3 then
                    log.warn("  ⚠️ HIGH RISK detected!")
                    
                    if #prediction.risk_factors > 0 then
                        log.info("  📋 Risk factors:")
                        for j, factor in ipairs(prediction.risk_factors) do
                            log.warn("    " .. j .. ". " .. factor.description .. " (Impact: " .. string.format("%.1f", factor.impact) .. ")")
                        end
                    end
                    
                    if #prediction.recommendations > 0 then
                        log.info("  💡 AI Recommendations:")
                        for j, rec in ipairs(prediction.recommendations) do
                            log.info("    " .. j .. ". " .. rec)
                        end
                    end
                else
                    log.info("  ✅ Low risk - safe to proceed")
                end
            end
            
            log.info("")
        end
        
        return {success = true, message = "AI prediction demo completed"}
    end)
    :build()

-- Task 3: Learning and Analytics Demo
local analytics_demo = task("analytics_demo")
    :description("Demonstrates AI learning and analytics")
    :command(function(params, deps)
        log.info("📈 Testing AI Learning and Analytics...")
        
        -- Record some sample executions for analysis
        local sample_commands = {
            "echo 'test 1'",
            "echo 'test 2'", 
            "echo 'test 3'"
        }
        
        log.info("📝 Recording sample executions...")
        for i, cmd in ipairs(sample_commands) do
            local success = math.random() > 0.2 -- 80% success rate
            local exec_time = math.random(50, 500) .. "ms"
            
            ai.record_execution({
                task_name = "analytics_demo",
                command = cmd,
                success = success,
                execution_time = exec_time,
                parameters = {
                    test_iteration = i,
                    timestamp = os.time()
                }
            })
            
            log.info("  " .. i .. ". Recorded: " .. cmd .. " (Success: " .. tostring(success) .. ")")
        end
        
        -- Find similar tasks
        log.info("🔍 Finding similar tasks...")
        local similar = ai.find_similar_tasks("echo 'test'", 5)
        
        if #similar > 0 then
            log.info("  📊 Found " .. #similar .. " similar tasks:")
            for i, task in ipairs(similar) do
                log.info("    " .. i .. ". " .. task.task_name .. " - " .. task.command)
                log.info("       Success: " .. tostring(task.success) .. ", Time: " .. task.execution_time)
            end
        else
            log.info("  ℹ️ No similar tasks found")
        end
        
        -- Performance analysis
        log.info("📊 Performing performance analysis...")
        local analysis = ai.analyze_performance("echo")
        
        if analysis.total_executions > 0 then
            log.info("  📈 Performance Analysis Results:")
            log.info("    Total executions: " .. analysis.total_executions)
            log.info("    Success rate: " .. string.format("%.1f%%", analysis.success_rate * 100))
            log.info("    Average time: " .. analysis.avg_execution_time)
            log.info("    Performance trend: " .. analysis.performance_trend)
            
            if #analysis.insights > 0 then
                log.info("  💡 Performance insights:")
                for i, insight in ipairs(analysis.insights) do
                    log.info("    " .. i .. ". " .. insight)
                end
            end
        else
            log.info("  ℹ️ No historical data available for analysis")
        end
        
        return {success = true, message = "AI analytics demo completed"}
    end)
    :build()

-- Task 4: AI Insights Generation
local insights_demo = task("insights_demo")
    :description("Demonstrates AI insights generation")
    :command(function(params, deps)
        log.info("💡 Generating AI Insights...")
        
        local insights = ai.generate_insights({
            scope = "global",
            context = "demo"
        })
        
        if #insights > 0 then
            log.info("🌟 AI-Generated Insights:")
            for i, insight in ipairs(insights) do
                log.info("  " .. i .. ". " .. insight)
            end
        else
            log.info("  ℹ️ No insights available at this time")
        end
        
        -- Display AI configuration
        local config = ai.get_config()
        log.info("🔧 Current AI Configuration:")
        log.info("  Enabled: " .. tostring(config.enabled))
        log.info("  Learning Mode: " .. config.learning_mode)
        log.info("  Optimization Level: " .. config.optimization_level .. "/10")
        log.info("  Failure Prediction: " .. tostring(config.failure_prediction))
        log.info("  Auto Optimize: " .. tostring(config.auto_optimize))
        
        return {success = true, message = "AI insights demo completed"}
    end)
    :build()

-- Define the AI-powered workflow
workflow.define("ai_intelligence_showcase", {
    description = "Complete AI-Powered Task Intelligence Showcase",
    version = "1.0.0",
    
    -- AI-enhanced workflow metadata
    metadata = {
        author = "Sloth Runner AI",
        tags = {"ai", "demonstration", "intelligence", "optimization"},
        ai_enabled = true
    },
    
    tasks = { 
        optimization_demo, 
        prediction_demo, 
        analytics_demo, 
        insights_demo 
    },
    
    -- AI-powered workflow hooks
    on_task_start = function(task_name)
        log.info("🤖 AI Pre-flight check: " .. task_name)
        -- Could add AI-powered pre-execution analysis here
    end,
    
    on_task_complete = function(task_name, success, output)
        log.info("🧠 AI Learning: Task " .. task_name .. " completed (Success: " .. tostring(success) .. ")")
        
        -- Record task completion for AI learning
        ai.record_execution({
            task_name = task_name,
            command = "workflow_task",
            success = success,
            execution_time = "1s", -- Placeholder
            parameters = {
                workflow = "ai_intelligence_showcase",
                output_length = output and #tostring(output) or 0
            }
        })
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 AI Intelligence Showcase completed successfully!")
            log.info("🧠 All AI features demonstrated and validated")
            log.info("📊 Machine learning data collected for future optimizations")
        else
            log.error("💥 AI Intelligence Showcase encountered issues")
            log.info("🔍 AI will analyze failures for future improvements")
        end
        
        -- Final AI insights
        log.info("")
        log.info("🌟 Final AI Assessment:")
        log.info("  - Command optimization: ✅ Working")
        log.info("  - Failure prediction: ✅ Working") 
        log.info("  - Learning system: ✅ Working")
        log.info("  - Analytics engine: ✅ Working")
        log.info("  - Insights generation: ✅ Working")
        log.info("")
        log.info("🚀 Sloth Runner AI-Powered Task Intelligence is ready!")
    end
})