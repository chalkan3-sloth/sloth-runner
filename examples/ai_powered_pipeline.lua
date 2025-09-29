-- AI-Powered Task Intelligence Example
-- This demonstrates the new AI capabilities integrated with Sloth Runner

local ai = require("ai")
local exec = require("exec")
local log = require("log")

-- Configure AI settings
ai.configure({
    enabled = true,
    learning_mode = "adaptive",
    optimization_level = 7,
    failure_prediction = true,
    auto_optimize = true
})

-- Example 1: AI-Optimized Build Task
local build_task = task("build_with_ai")
    :description("Build application with AI optimization")
    :ai_optimization(true)
    :learning_mode("adaptive")
    :command(function(params, deps)
        local original_command = "go build -o app ./cmd/main.go"
        
        -- Get AI optimization suggestions
        local ai_suggestions = ai.optimize_command(original_command, {
            history = ai.get_task_history(original_command),
            system_resources = {
                cpu_usage = 45,
                memory_usage = 60,
                load_avg = 1.2
            }
        })
        
        if ai_suggestions and ai_suggestions.confidence_score > 0.7 then
            log.info("🤖 AI Optimization Applied!")
            log.info("📈 Expected Speedup: " .. string.format("%.1fx", ai_suggestions.expected_speedup))
            log.info("💡 Rationale: " .. ai_suggestions.rationale)
            
            -- Use optimized command
            local result = exec.run(ai_suggestions.optimized_command)
            
            -- Record execution for learning
            ai.record_execution({
                task_name = "build_with_ai",
                command = ai_suggestions.optimized_command,
                success = result.success,
                execution_time = result.duration or "0s",
                optimizations = {"ai_optimized"}
            })
            
            return result
        else
            log.info("ℹ️ Using original command (low AI confidence)")
            local result = exec.run(original_command)
            
            -- Still record execution for learning
            ai.record_execution({
                task_name = "build_with_ai",
                command = original_command,
                success = result.success,
                execution_time = result.duration or "0s"
            })
            
            return result
        end
    end)
    :timeout("5m")
    :retries(3, "exponential")
    :on_success(function(params, output)
        log.info("✅ Build completed successfully with AI assistance!")
    end)
    :build()

-- Example 2: Predictive Failure Detection
local deploy_task = task("smart_deploy")
    :description("Deploy with AI failure prediction")
    :command(function(params, deps)
        local deploy_command = "kubectl apply -f k8s/production/"
        
        -- Predict failure probability
        local prediction = ai.predict_failure("smart_deploy", deploy_command)
        
        log.info("🔮 AI Failure Prediction:")
        log.info("📊 Failure Probability: " .. string.format("%.1f%%", prediction.failure_probability * 100))
        log.info("🎯 Confidence: " .. string.format("%.1f%%", prediction.confidence * 100))
        
        -- Check if failure probability is too high
        if prediction.failure_probability > 0.3 then
            log.warn("⚠️ High failure risk detected!")
            
            -- Show risk factors
            for i, factor in ipairs(prediction.risk_factors) do
                log.warn("❌ " .. factor.description .. " (Impact: " .. string.format("%.1f", factor.impact) .. ")")
            end
            
            -- Show recommendations
            log.info("💡 AI Recommendations:")
            for i, rec in ipairs(prediction.recommendations) do
                log.info("  " .. i .. ". " .. rec)
            end
            
            -- Ask for confirmation
            print("Proceed with deployment? (y/N)")
            local response = io.read()
            if response:lower() ~= "y" and response:lower() ~= "yes" then
                return {success = false, message = "Deployment cancelled due to high failure risk"}
            end
        else
            log.info("✅ Low failure risk - proceeding with deployment")
        end
        
        -- Execute deployment
        local result = exec.run(deploy_command)
        
        -- Record execution
        ai.record_execution({
            task_name = "smart_deploy",
            command = deploy_command,
            success = result.success,
            execution_time = result.duration or "0s",
            error_message = result.error
        })
        
        return result
    end)
    :build()

-- Example 3: Learning from Similar Tasks
local test_task = task("intelligent_test")
    :description("Run tests with insights from similar tasks")
    :command(function(params, deps)
        local test_command = "go test ./..."
        
        -- Find similar tasks
        local similar_tasks = ai.find_similar_tasks(test_command, 5)
        
        if #similar_tasks > 0 then
            log.info("🧠 Learning from " .. #similar_tasks .. " similar tasks")
            
            local success_count = 0
            local avg_time = 0
            
            for i, task in ipairs(similar_tasks) do
                if task.success then
                    success_count = success_count + 1
                end
                -- Parse execution time (simplified)
                local time_val = tonumber(task.execution_time:match("(%d+)"))
                if time_val then
                    avg_time = avg_time + time_val
                end
            end
            
            local success_rate = success_count / #similar_tasks
            avg_time = avg_time / #similar_tasks
            
            log.info("📈 Historical Success Rate: " .. string.format("%.1f%%", success_rate * 100))
            log.info("⏱️ Average Execution Time: ~" .. avg_time .. "s")
            
            if success_rate < 0.8 then
                log.warn("⚠️ Historical data shows lower success rate - adding extra validation")
                test_command = test_command .. " -v"
            end
        end
        
        -- Execute tests
        local result = exec.run(test_command)
        
        -- Record execution
        ai.record_execution({
            task_name = "intelligent_test",
            command = test_command,
            success = result.success,
            execution_time = result.duration or "0s",
            error_message = result.error
        })
        
        return result
    end)
    :build()

-- Example 4: Performance Analysis and Insights
local analysis_task = task("performance_analysis")
    :description("Analyze task performance using AI")
    :command(function(params, deps)
        local commands_to_analyze = {
            "go build -o app ./cmd/main.go",
            "go test ./...",
            "docker build -t myapp ."
        }
        
        for i, command in ipairs(commands_to_analyze) do
            log.info("🔍 Analyzing: " .. command)
            
            -- Get performance analysis
            local analysis = ai.analyze_performance(command)
            
            if analysis.total_executions > 0 then
                log.info("📊 Performance Analysis:")
                log.info("  Total Executions: " .. analysis.total_executions)
                log.info("  Success Rate: " .. string.format("%.1f%%", analysis.success_rate * 100))
                log.info("  Avg Time: " .. analysis.avg_execution_time)
                log.info("  Trend: " .. analysis.performance_trend)
                
                if #analysis.insights > 0 then
                    log.info("💡 Insights:")
                    for j, insight in ipairs(analysis.insights) do
                        log.info("  " .. j .. ". " .. insight)
                    end
                end
            else
                log.info("ℹ️ No historical data available")
            end
            
            log.info("---")
        end
        
        -- Generate general insights
        local general_insights = ai.generate_insights({
            scope = "global"
        })
        
        log.info("🌟 General AI Insights:")
        for i, insight in ipairs(general_insights) do
            log.info("  " .. i .. ". " .. insight)
        end
        
        return {success = true, message = "Performance analysis completed"}
    end)
    :build()

-- Example 5: AI-Enhanced Workflow
workflow.define("ai_enhanced_pipeline", {
    description = "CI/CD pipeline with AI optimization and prediction",
    version = "2.0.0",
    
    -- Enable AI features for the entire workflow
    ai_features = {
        optimization = true,
        failure_prediction = true,
        performance_monitoring = true,
        adaptive_learning = true
    },
    
    tasks = { build_task, test_task, deploy_task, analysis_task },
    
    -- AI-powered workflow hooks
    on_task_start = function(task_name)
        -- Predict failure before starting each task
        log.info("🤖 AI Pre-flight check for: " .. task_name)
        
        -- Get task statistics
        local stats = ai.get_task_stats(task_name)
        if stats and stats.total_runs > 0 then
            log.info("📊 Task History: " .. stats.total_runs .. " runs, " .. 
                    string.format("%.1f%%", stats.success_rate * 100) .. " success rate")
        end
    end,
    
    on_task_complete = function(task_name, success, output)
        -- Log completion for AI learning
        log.info("🧠 Recording task completion for AI learning")
        
        -- AI can analyze the output and learn from it
        if not success then
            log.warn("❌ Task failed - AI will learn from this failure")
        else
            log.info("✅ Task succeeded - AI will reinforce successful patterns")
        end
    end,
    
    on_workflow_complete = function(success, results)
        if success then
            log.info("🎉 AI-Enhanced Pipeline completed successfully!")
            
            -- Generate post-workflow insights
            local insights = ai.generate_insights({
                scope = "workflow",
                workflow_name = "ai_enhanced_pipeline"
            })
            
            log.info("🔮 Post-Workflow AI Insights:")
            for i, insight in ipairs(insights) do
                log.info("  " .. i .. ". " .. insight)
            end
        else
            log.error("💥 Pipeline failed - AI will analyze failure patterns")
        end
    end
})

-- Display AI configuration
local ai_config = ai.get_config()
log.info("🤖 AI Configuration:")
log.info("  Enabled: " .. tostring(ai_config.enabled))
log.info("  Learning Mode: " .. ai_config.learning_mode)
log.info("  Optimization Level: " .. ai_config.optimization_level .. "/10")
log.info("  Failure Prediction: " .. tostring(ai_config.failure_prediction))
log.info("  Auto Optimize: " .. tostring(ai_config.auto_optimize))