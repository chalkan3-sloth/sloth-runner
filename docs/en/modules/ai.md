# 🤖 AI Module - Complete API Reference

The AI module provides artificial intelligence capabilities for task optimization, failure prediction, and performance analytics.

## 📋 Module Overview

```lua
local ai = require("ai")
```

The AI module is the core of Sloth Runner's intelligence features, providing:

- **🔮 Predictive Failure Detection** - Predict task failures before they happen
- **⚡ Intelligent Optimization** - Automatically optimize commands for better performance  
- **📊 Performance Analytics** - Analyze execution patterns and trends
- **🧠 Adaptive Learning** - Continuous improvement from execution history

## 🔧 Configuration

### `ai.configure(config)`

Configure AI behavior and capabilities.

```lua
ai.configure({
    enabled = true,                    -- Enable/disable AI features
    learning_mode = "adaptive",        -- adaptive | aggressive | conservative
    optimization_level = 8,            -- 1-10 (higher = more aggressive)
    failure_prediction = true,         -- Enable failure prediction
    auto_optimize = true,              -- Automatically apply optimizations
    confidence_threshold = 0.7         -- Minimum confidence for auto-apply
})
```

**Parameters:**
- `enabled` (boolean): Enable or disable all AI features
- `learning_mode` (string): Learning aggressiveness level
- `optimization_level` (number): Optimization aggressiveness (1-10)
- `failure_prediction` (boolean): Enable predictive failure detection
- `auto_optimize` (boolean): Automatically apply high-confidence optimizations
- `confidence_threshold` (number): Minimum confidence score for auto-application

### `ai.get_config()`

Get current AI configuration.

```lua
local config = ai.get_config()
-- Returns: {enabled: true, learning_mode: "adaptive", ...}
```

## ⚡ Optimization

### `ai.optimize_command(command, options)`

Get AI optimization suggestions for a command.

```lua
local result = ai.optimize_command("go build -o app ./cmd/main.go", {
    history = ai.get_task_history("go build"),
    system_resources = {
        cpu_usage = 45,
        memory_usage = 60,
        load_avg = 1.2
    },
    similar_tasks = ai.find_similar_tasks("go build", 10),
    environment = "production"
})
```

**Parameters:**
- `command` (string): Original command to optimize
- `options` (table): Optimization context
  - `history` (array): Historical executions of this command
  - `system_resources` (table): Current system resource usage
  - `similar_tasks` (array): Similar task executions
  - `environment` (string): Execution environment (dev/staging/prod)

**Returns:**
```lua
{
    original_command = "go build -o app ./cmd/main.go",
    optimized_command = "go build -p 4 -ldflags='-s -w' -o app ./cmd/main.go",
    confidence_score = 0.85,           -- 0.0-1.0
    expected_speedup = 2.3,            -- Expected performance multiplier
    optimizations = {                  -- Applied optimizations
        {
            type = "parallelization",
            description = "Added -p 4 for parallel compilation",
            impact = 1.8
        },
        {
            type = "size_optimization", 
            description = "Added -ldflags='-s -w' to reduce binary size",
            impact = 0.5
        }
    },
    resource_savings = {
        estimated_time_saved = "1.2s",
        memory_efficiency = "+15%"
    },
    rationale = "Command shows parallelization opportunities based on system CPU count"
}
```

## 🔮 Failure Prediction

### `ai.predict_failure(task_name, command, options)`

Predict the probability of task failure.

```lua
local prediction = ai.predict_failure("deploy_task", "kubectl apply -f deployment.yaml", {
    history = ai.get_task_history("kubectl apply"),
    environment = "production",
    system_state = {
        disk_usage = 85,
        network_latency = 120
    }
})
```

**Parameters:**
- `task_name` (string): Name of the task being analyzed
- `command` (string): Command to be executed
- `options` (table): Prediction context
  - `history` (array): Historical executions
  - `environment` (string): Execution environment
  - `system_state` (table): Current system state

**Returns:**
```lua
{
    failure_probability = 0.23,        -- 0.0-1.0
    confidence = 0.78,                 -- Confidence in prediction
    risk_factors = {                   -- Identified risk factors
        {
            type = "resource_contention",
            description = "High disk usage detected (85%)",
            impact = 0.6,
            severity = "medium"
        },
        {
            type = "network_latency",
            description = "Elevated network latency (120ms)",
            impact = 0.3,
            severity = "low"
        }
    },
    recommendations = {                -- AI-generated recommendations
        "Consider waiting for disk usage to decrease below 80%",
        "Add timeout configuration to handle network latency",
        "Implement retry logic with exponential backoff"
    },
    similar_failures = {               -- Historical similar failures
        count = 3,
        common_causes = ["network_timeout", "resource_exhaustion"]
    }
}
```

## 📊 Performance Analytics

### `ai.analyze_performance(command, options)`

Analyze performance patterns for a command or task.

```lua
local analysis = ai.analyze_performance("go build", {
    time_range = "30d",                -- 1d, 7d, 30d, 90d
    environment = "all",               -- all, dev, staging, prod
    include_failures = true
})
```

**Parameters:**
- `command` (string): Command to analyze
- `options` (table): Analysis options
  - `time_range` (string): Time range for analysis
  - `environment` (string): Environment filter
  - `include_failures` (boolean): Include failed executions

**Returns:**
```lua
{
    total_executions = 156,
    success_rate = 0.94,               -- 94% success rate
    avg_execution_time = "2.3s",
    fastest_execution = "1.1s",
    slowest_execution = "5.7s",
    performance_trend = "improving",    -- improving | stable | degrading
    insights = {                       -- AI-generated insights
        "Performance improved 23% over the last 30 days",
        "Failures primarily occur during high system load",
        "Consider caching to improve cold-start performance"
    },
    recommendations = {
        "Enable build caching to reduce average execution time",
        "Implement resource monitoring for failure prevention"
    },
    patterns = {                       -- Detected patterns
        peak_hours = ["09:00-10:00", "14:00-15:00"],
        failure_correlation = ["high_cpu_usage", "memory_pressure"]
    }
}
```

### `ai.get_task_stats(task_name)`

Get aggregated statistics for a specific task.

```lua
local stats = ai.get_task_stats("build_application")
```

**Returns:**
```lua
{
    task_name = "build_application",
    total_runs = 89,
    success_count = 84,
    failure_count = 5,
    success_rate = 0.944,              -- 94.4%
    total_time = "3m 45s",
    avg_time = "2.5s",
    fastest_time = "1.2s",
    slowest_time = "8.1s",
    last_execution = "2024-01-15T10:30:00Z",
    trend = "stable"
}
```

## 🧠 Learning & History

### `ai.record_execution(execution_data)`

Record task execution for AI learning.

```lua
ai.record_execution({
    task_name = "build_application",
    command = "go build -o app ./cmd/main.go",
    success = true,
    execution_time = "2.5s",
    start_time = os.time(),
    end_time = os.time() + 2.5,
    parameters = {
        environment = "development",
        go_version = "1.21.0",
        parallel = true
    },
    system_resources = {
        cpu_usage = 45,
        memory_usage = 60,
        disk_usage = 30
    },
    error_message = nil,               -- If success = false
    optimization_applied = true,
    ai_confidence = 0.85
})
```

**Parameters:**
- `task_name` (string): Name of the executed task
- `command` (string): Command that was executed
- `success` (boolean): Whether execution was successful
- `execution_time` (string): Time taken to execute
- `parameters` (table): Execution parameters and context
- `system_resources` (table): System resource state during execution
- `error_message` (string): Error message if failed
- `optimization_applied` (boolean): Whether AI optimization was used
- `ai_confidence` (number): Confidence score if optimization was applied

### `ai.get_task_history(command, limit)`

Get execution history for a command.

```lua
local history = ai.get_task_history("go build", 20)
-- Returns array of execution records
```

### `ai.find_similar_tasks(command, limit)`

Find tasks similar to the given command.

```lua
local similar = ai.find_similar_tasks("go build -o app", 10)
-- Returns array of similar task executions
```

## 💡 Insights & Recommendations

### `ai.generate_insights(options)`

Generate AI-powered insights about task execution patterns.

```lua
local insights = ai.generate_insights({
    scope = "global",                  -- global | task | command
    task_name = "build_application",   -- if scope = "task"
    time_range = "7d"
})
```

**Returns:**
```lua
{
    "Tasks executed during business hours have 15% lower failure rate",
    "Commands with parallel flags show 40% better performance", 
    "Memory-intensive tasks perform better with explicit heap size settings",
    "Network-dependent tasks should include timeout and retry configurations"
}
```

## 🎯 Best Practices

### 1. **Always Record Executions**
```lua
-- Record every execution for AI learning
workflow.define("my_pipeline", {
    on_task_complete = function(task_name, success, output)
        ai.record_execution({
            task_name = task_name,
            command = output.command,
            success = success,
            execution_time = output.duration
        })
    end
})
```

### 2. **Use Confidence Thresholds**
```lua
-- Only apply high-confidence optimizations
local optimization = ai.optimize_command(command)
if optimization.confidence_score > 0.8 then
    command = optimization.optimized_command
    log.info("Applied AI optimization with " .. (optimization.confidence_score * 100) .. "% confidence")
end
```

### 3. **Monitor Predictions**
```lua
-- Always check predictions for critical tasks
local prediction = ai.predict_failure(task_name, command)
if prediction.failure_probability > 0.3 then
    log.warn("High failure risk detected: " .. (prediction.failure_probability * 100) .. "%")
    for _, rec in ipairs(prediction.recommendations) do
        log.info("Recommendation: " .. rec)
    end
end
```

### 4. **Regular Analysis**
```lua
-- Periodic performance analysis
local analysis = ai.analyze_performance("critical_task")
if analysis.performance_trend == "degrading" then
    log.warn("Performance degradation detected for critical_task")
    -- Take action
end
```

## 🔬 Advanced Features

### Learning Modes

- **Adaptive**: Balanced learning and optimization (recommended)
- **Aggressive**: Maximum optimization attempts, higher risk
- **Conservative**: Minimal changes, maximum safety

### Optimization Strategies

The AI system includes multiple built-in optimization strategies:
- **Parallelization**: Detect parallel execution opportunities
- **Memory Optimization**: Adjust memory settings for optimal performance
- **Compiler Optimization**: Suggest better compiler flags and options
- **Caching**: Implement intelligent caching strategies
- **Network Optimization**: Optimize network operations and timeouts
- **I/O Optimization**: Improve file and disk operations

### Custom Metrics

You can provide custom metrics to improve AI analysis:

```lua
ai.record_execution({
    task_name = "custom_task",
    success = true,
    execution_time = "1.5s",
    custom_metrics = {
        memory_peak = "512MB",
        cache_hit_rate = 0.85,
        network_requests = 15,
        database_queries = 8
    }
})
```

## 🚀 Integration Examples

### With Modern DSL

```lua
local build_task = task("ai_optimized_build")
    :description("Build with AI optimization")
    :command(function(params, deps)
        local cmd = "go build -o app ./cmd/main.go"
        local optimization = ai.optimize_command(cmd, {
            history = ai.get_task_history(cmd)
        })
        
        if optimization.confidence_score > 0.7 then
            return exec.run(optimization.optimized_command)
        else
            return exec.run(cmd)
        end
    end)
    :on_success(function(params, output)
        ai.record_execution({
            task_name = "ai_optimized_build",
            command = output.command,
            success = true,
            execution_time = output.duration
        })
    end)
    :build()
```

### With GitOps

```lua
local gitops_task = task("intelligent_deploy")
    :command(function(params, deps)
        local deploy_cmd = "kubectl apply -f manifests/"
        
        -- AI failure prediction
        local prediction = ai.predict_failure("intelligent_deploy", deploy_cmd)
        if prediction.failure_probability > 0.25 then
            log.warn("High deployment risk detected")
            return {success = false, message = "Deployment blocked by AI risk assessment"}
        end
        
        -- GitOps deployment
        return gitops.sync_workflow(params.workflow_id)
    end)
    :build()
```

## 📚 See Also

- [AI Features Overview](../ai-features.md)
- [Performance Optimization Guide](../ai/optimization.md)
- [Failure Prediction Guide](../ai/prediction.md)
- [AI Best Practices](../ai/best-practices.md)