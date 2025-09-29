# 🚀 Sloth Runner - Quick Tutorial

> **Welcome to the world's most intelligent task orchestration platform!**  
> This tutorial will get you up and running with AI-powered automation and GitOps workflows in just 5 minutes.

## 🎯 What You'll Learn

By the end of this tutorial, you'll have:
- ✅ Installed Sloth Runner
- ✅ Created your first AI-optimized task
- ✅ Set up a GitOps workflow
- ✅ Built an intelligent automation pipeline

## 📦 Step 1: Installation

### Quick Install (Recommended)

```bash
# Install the latest version
curl -sSL https://get.sloth-runner.dev | bash

# Verify installation
sloth-runner --version
```

### Manual Install

```bash
# Download from GitHub releases
wget https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner-linux-amd64.tar.gz
tar -xzf sloth-runner-linux-amd64.tar.gz
sudo mv sloth-runner /usr/local/bin/
```

### Verify Installation

```bash
sloth-runner --help
```

## 🤖 Step 2: Your First AI-Powered Task

Create a file called `hello-ai.lua`:

```lua
-- hello-ai.lua
local ai = require("ai")
local log = require("log")

-- Configure AI for optimal performance
ai.configure({
    enabled = true,
    learning_mode = "adaptive",
    optimization_level = 8,
    failure_prediction = true
})

log.info("🤖 Welcome to AI-Powered Automation!")

-- Create an intelligent task
local hello_task = task("ai_hello_world")
    :description("My first AI-optimized task")
    :command(function(params, deps)
        local original_command = "echo 'Hello, AI-Powered World!'"
        
        -- Let AI optimize the command
        local ai_result = ai.optimize_command(original_command)
        
        if ai_result.confidence_score > 0.5 then
            log.info("🤖 AI optimized the command!")
            log.info("📈 Confidence: " .. string.format("%.1f%%", ai_result.confidence_score * 100))
            log.info("⚡ Expected speedup: " .. string.format("%.1fx", ai_result.expected_speedup))
            
            -- Use AI-optimized command
            return exec.run(ai_result.optimized_command)
        else
            log.info("ℹ️ Using original command (low AI confidence)")
            return exec.run(original_command)
        end
    end)
    :on_success(function(params, output)
        log.info("✅ Task completed successfully!")
        
        -- Record execution for AI learning
        ai.record_execution({
            task_name = "ai_hello_world",
            command = output.command,
            success = true,
            execution_time = output.duration or "0s"
        })
    end)
    :build()

-- Create a simple workflow
workflow.define("ai_tutorial", {
    description = "AI-Powered Hello World Tutorial",
    tasks = { hello_task },
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 AI Tutorial completed successfully!")
            log.info("🧠 AI is now learning from this execution")
        end
    end
})
```

### Run Your AI-Powered Task

```bash
sloth-runner run -f hello-ai.lua
```

**Expected Output:**
```
🤖 Welcome to AI-Powered Automation!
🤖 AI optimized the command!
📈 Confidence: 85.0%
⚡ Expected speedup: 1.2x
Hello, AI-Powered World!
✅ Task completed successfully!
🎉 AI Tutorial completed successfully!
🧠 AI is now learning from this execution
```

## 🔄 Step 3: Your First GitOps Workflow

Create a file called `hello-gitops.lua`:

```lua
-- hello-gitops.lua
local gitops = require("gitops")
local log = require("log")

log.info("🔄 Welcome to GitOps Native Automation!")

-- Create a GitOps workflow (you can use any public repo for testing)
local git_workflow = gitops.workflow({
    repo = "https://github.com/kubernetes/examples",  -- Example public repo
    branch = "master",
    auto_sync = false,  -- Manual sync for this demo
    diff_preview = true,
    rollback_on_failure = true
})

log.info("✅ GitOps workflow created!")
log.info("📋 Workflow ID: " .. git_workflow.workflow_id)
log.info("📦 Repository ID: " .. git_workflow.repository_id)

-- Create a deployment task
local deploy_task = task("gitops_deploy")
    :description("GitOps-powered deployment")
    :command(function(params, deps)
        -- Preview changes before deployment
        log.info("🔍 Generating diff preview...")
        local diff = gitops.generate_diff(git_workflow.workflow_id)
        
        if diff then
            log.info("📊 GitOps Change Summary:")
            log.info("  📝 Total changes: " .. diff.summary.total_changes)
            log.info("  ✨ Created resources: " .. diff.summary.created_resources)
            log.info("  🔄 Updated resources: " .. diff.summary.updated_resources)
            log.info("  🗑️ Deleted resources: " .. diff.summary.deleted_resources)
            
            if diff.summary.conflict_count > 0 then
                log.warn("⚠️ " .. diff.summary.conflict_count .. " conflict(s) detected")
                return {success = false, message = "Conflicts need resolution"}
            end
        else
            log.info("ℹ️ No changes detected")
        end
        
        -- Simulate deployment (since we're using a read-only example repo)
        log.info("🚀 Simulating GitOps deployment...")
        log.info("✅ GitOps deployment completed successfully!")
        
        return {success = true, message = "GitOps deployment successful"}
    end)
    :build()

-- Monitor workflow status
local status_task = task("check_status")
    :description("Check GitOps workflow status")
    :command(function(params, deps)
        local status = gitops.get_workflow_status(git_workflow.workflow_id)
        
        if status then
            log.info("📊 GitOps Workflow Status:")
            log.info("  🏷️ Name: " .. status.name)
            log.info("  📈 Status: " .. status.status)
            log.info("  🔄 Auto-sync: " .. tostring(status.auto_sync))
            log.info("  📦 Repository: " .. status.repository)
        end
        
        return {success = true, message = "Status check completed"}
    end)
    :build()

-- Create GitOps workflow
workflow.define("gitops_tutorial", {
    description = "GitOps Native Tutorial",
    tasks = { deploy_task, status_task },
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 GitOps Tutorial completed successfully!")
            log.info("🔄 You now understand GitOps native workflows!")
        end
    end
})
```

### Run Your GitOps Workflow

```bash
sloth-runner run -f hello-gitops.lua
```

**Expected Output:**
```
🔄 Welcome to GitOps Native Automation!
✅ GitOps workflow created!
📋 Workflow ID: workflow-1234567890
📦 Repository ID: repo-1234567890
🔍 Generating diff preview...
📊 GitOps Change Summary:
  📝 Total changes: 3
  ✨ Created resources: 2
  🔄 Updated resources: 1
  🗑️ Deleted resources: 0
🚀 Simulating GitOps deployment...
✅ GitOps deployment completed successfully!
📊 GitOps Workflow Status:
  🏷️ Name: lua-generated-workflow
  📈 Status: syncing
  🔄 Auto-sync: false
  📦 Repository: repo-1234567890
🎉 GitOps Tutorial completed successfully!
🔄 You now understand GitOps native workflows!
```

## ⚡ Step 4: Intelligent Pipeline (AI + GitOps)

Create a file called `intelligent-pipeline.lua`:

```lua
-- intelligent-pipeline.lua  
local ai = require("ai")
local gitops = require("gitops")
local log = require("log")

log.info("🧠 Building an Intelligent Pipeline with AI + GitOps!")

-- Configure AI for production-ready optimization
ai.configure({
    enabled = true,
    learning_mode = "adaptive",
    optimization_level = 7,
    failure_prediction = true,
    auto_optimize = true
})

-- Create GitOps workflow for deployment
local deployment_workflow = gitops.workflow({
    repo = "https://github.com/kubernetes/examples",
    branch = "master", 
    auto_sync = false,
    diff_preview = true,
    rollback_on_failure = true
})

-- AI-optimized build task
local build_task = task("intelligent_build")
    :description("AI-optimized build process")
    :command(function(params, deps)
        local build_command = "echo 'Building application...'"
        
        -- AI optimization with context
        local optimization = ai.optimize_command(build_command, {
            history = ai.get_task_history(build_command),
            system_resources = {
                cpu_usage = 30,
                memory_usage = 50
            }
        })
        
        if optimization.confidence_score > 0.6 then
            log.info("🤖 AI Build Optimization Applied!")
            log.info("📈 Confidence: " .. string.format("%.1f%%", optimization.confidence_score * 100))
            log.info("⚡ Expected Speedup: " .. string.format("%.1fx", optimization.expected_speedup))
            
            return exec.run(optimization.optimized_command)
        else
            log.info("🔧 Using standard build process")
            return exec.run(build_command)
        end
    end)
    :build()

-- AI-enhanced deployment with failure prediction
local deploy_task = task("intelligent_deploy")
    :description("AI-predicted GitOps deployment")
    :depends_on({"intelligent_build"})
    :command(function(params, deps)
        local deploy_command = "kubectl apply -f manifests/"
        
        -- AI failure prediction
        local prediction = ai.predict_failure("intelligent_deploy", deploy_command)
        
        log.info("🔮 AI Deployment Analysis:")
        log.info("📊 Failure Probability: " .. string.format("%.1f%%", prediction.failure_probability * 100))
        log.info("🎯 Prediction Confidence: " .. string.format("%.1f%%", prediction.confidence * 100))
        
        if prediction.failure_probability > 0.3 then
            log.warn("⚠️ HIGH RISK DEPLOYMENT DETECTED!")
            log.warn("🚨 AI recommends caution")
            
            for _, recommendation in ipairs(prediction.recommendations) do
                log.info("💡 AI Recommendation: " .. recommendation)
            end
        else
            log.info("✅ AI assessment: Low risk deployment")
        end
        
        -- GitOps deployment with intelligent preview
        log.info("🔍 GitOps intelligent diff analysis...")
        local diff = gitops.generate_diff(deployment_workflow.workflow_id)
        
        if diff and diff.summary.total_changes > 0 then
            log.info("📋 Deployment will apply " .. diff.summary.total_changes .. " changes")
            
            -- Check for conflicts
            if diff.summary.conflict_count > 0 then
                log.warn("💥 " .. diff.summary.conflict_count .. " conflicts detected!")
                return {success = false, message = "Conflicts require manual resolution"}
            end
        end
        
        -- Execute intelligent deployment
        log.info("🚀 Executing AI + GitOps intelligent deployment...")
        local deploy_success = true  -- Simulated for tutorial
        
        if deploy_success then
            log.info("✅ Intelligent deployment completed successfully!")
            
            -- Record execution for AI learning
            ai.record_execution({
                task_name = "intelligent_deploy",
                command = deploy_command,
                success = true,
                execution_time = "45s",
                ai_prediction_used = true,
                predicted_failure_probability = prediction.failure_probability,
                gitops_changes_applied = diff and diff.summary.total_changes or 0
            })
            
            return {success = true, message = "AI + GitOps deployment successful"}
        else
            log.error("💥 Deployment failed - AI will learn from this failure")
            return {success = false, message = "Deployment failed"}
        end
    end)
    :build()

-- Intelligent monitoring task
local monitor_task = task("intelligent_monitor")
    :description("AI-powered post-deployment monitoring")
    :depends_on({"intelligent_deploy"})
    :command(function(params, deps)
        log.info("📊 Starting AI-powered monitoring...")
        
        -- Get performance analysis
        local analysis = ai.analyze_performance("intelligent_deploy")
        
        if analysis.total_executions > 0 then
            log.info("📈 AI Performance Analysis:")
            log.info("  📊 Total executions: " .. analysis.total_executions)
            log.info("  ✅ Success rate: " .. string.format("%.1f%%", analysis.success_rate * 100))
            log.info("  ⏱️ Average time: " .. analysis.avg_execution_time)
            log.info("  📈 Trend: " .. analysis.performance_trend)
            
            if #analysis.insights > 0 then
                log.info("💡 AI Insights:")
                for _, insight in ipairs(analysis.insights) do
                    log.info("  🔮 " .. insight)
                end
            end
        else
            log.info("ℹ️ Insufficient data for AI analysis - will improve with more executions")
        end
        
        -- Check GitOps workflow status
        local gitops_status = gitops.get_workflow_status(deployment_workflow.workflow_id)
        log.info("🔄 GitOps Status: " .. gitops_status.status)
        
        return {success = true, message = "Intelligent monitoring completed"}
    end)
    :build()

-- Create the intelligent pipeline
workflow.define("intelligent_pipeline", {
    description = "AI + GitOps Intelligent Automation Pipeline",
    version = "2.0.0",
    
    metadata = {
        author = "AI-Powered DevOps",
        tags = {"ai", "gitops", "intelligent", "tutorial"},
        ai_enabled = true,
        gitops_enabled = true
    },
    
    tasks = { build_task, deploy_task, monitor_task },
    
    on_task_start = function(task_name)
        log.info("🤖 AI Pre-task Analysis: " .. task_name)
    end,
    
    on_task_complete = function(task_name, success, output)
        if success then
            log.info("✅ Intelligent task completed: " .. task_name)
        else
            log.error("❌ Task failed - AI will analyze failure: " .. task_name)
        end
        
        -- Always record for AI learning
        ai.record_execution({
            task_name = task_name,
            success = success,
            execution_time = output.duration or "0s",
            pipeline = "intelligent_pipeline"
        })
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 INTELLIGENT PIPELINE COMPLETED SUCCESSFULLY!")
            log.info("🧠 AI has learned from this execution")
            log.info("🔄 GitOps state is synchronized")
            log.info("📊 Performance data collected for future optimizations")
            
            -- Generate final AI insights
            local insights = ai.generate_insights({scope = "pipeline"})
            if #insights > 0 then
                log.info("🌟 Final AI Insights:")
                for _, insight in ipairs(insights) do
                    log.info("  💡 " .. insight)
                end
            end
        else
            log.error("💥 Intelligent pipeline failed")
            log.info("🔍 AI will analyze failure patterns for future improvements")
        end
    end
})
```

### Run Your Intelligent Pipeline

```bash
sloth-runner run -f intelligent-pipeline.lua
```

**Expected Output:**
```
🧠 Building an Intelligent Pipeline with AI + GitOps!
🤖 AI Pre-task Analysis: intelligent_build
🤖 AI Build Optimization Applied!
📈 Confidence: 78.0%
⚡ Expected Speedup: 1.5x
Building application...
✅ Intelligent task completed: intelligent_build
🤖 AI Pre-task Analysis: intelligent_deploy
🔮 AI Deployment Analysis:
📊 Failure Probability: 15.2%
🎯 Prediction Confidence: 82.0%
✅ AI assessment: Low risk deployment
🔍 GitOps intelligent diff analysis...
📋 Deployment will apply 3 changes
🚀 Executing AI + GitOps intelligent deployment...
✅ Intelligent deployment completed successfully!
✅ Intelligent task completed: intelligent_deploy
🤖 AI Pre-task Analysis: intelligent_monitor
📊 Starting AI-powered monitoring...
📈 AI Performance Analysis:
  📊 Total executions: 1
  ✅ Success rate: 100.0%
  ⏱️ Average time: 45s
  📈 Trend: stable
🔄 GitOps Status: syncing
✅ Intelligent task completed: intelligent_monitor
🎉 INTELLIGENT PIPELINE COMPLETED SUCCESSFULLY!
🧠 AI has learned from this execution
🔄 GitOps state is synchronized
📊 Performance data collected for future optimizations
🌟 Final AI Insights:
  💡 Consider implementing parallel execution for better performance
  💡 Task execution patterns suggest optimal scheduling during off-peak hours
```

## 🎓 What You've Learned

Congratulations! You've just:

✅ **Mastered AI-Powered Automation**: Created tasks that use artificial intelligence for optimization and failure prediction

✅ **Implemented GitOps Native Workflows**: Set up Git-driven deployments with intelligent diff preview

✅ **Built Intelligent Pipelines**: Combined AI and GitOps for the ultimate automation experience

✅ **Used Modern DSL**: Experienced the clean, powerful syntax of Sloth Runner's domain-specific language

## 🚀 Next Steps

Now that you understand the basics, explore these advanced features:

### 🤖 Deep Dive into AI
- [AI Features Complete Guide](en/ai-features.md)
- [Performance Optimization Strategies](en/ai/optimization.md)
- [Failure Prediction Best Practices](en/ai/prediction.md)

### 🔄 Master GitOps
- [GitOps Features Complete Guide](en/gitops-features.md)
- [Multi-Environment Deployments](en/gitops/multi-env.md)
- [Kubernetes Integration](en/gitops/kubernetes.md)

### 🏢 Enterprise Features
- [Distributed Architecture](en/master-agent-architecture.md)
- [Security & RBAC](en/security.md)
- [Monitoring & Observability](en/monitoring.md)

### 📚 Real-World Examples
- [Production CI/CD Pipelines](en/examples/cicd.md)
- [Infrastructure as Code](en/examples/iac.md)
- [Multi-Cloud Deployments](en/examples/multi-cloud.md)

## 💬 Community & Support

- **📖 [Complete Documentation](https://sloth-runner.dev/docs)** - Comprehensive guides and references
- **💬 [Discord Community](https://discord.gg/sloth-runner)** - Get help and share experiences  
- **🐛 [GitHub Issues](https://github.com/chalkan3-sloth/sloth-runner/issues)** - Report bugs and request features
- **🏢 [Enterprise Support](mailto:enterprise@sloth-runner.dev)** - Commercial support and services

---

**🎉 Welcome to the future of intelligent automation!**

You're now ready to build production-grade workflows with the world's most advanced task orchestration platform. Start building amazing things with AI-powered automation and GitOps native workflows!