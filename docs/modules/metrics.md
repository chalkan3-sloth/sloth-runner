# 📊 Metrics & Monitoring Module

The **Metrics & Monitoring** module provides comprehensive system monitoring, custom metrics collection, and health checking capabilities. It enables real-time observability of both system resources and application performance.

## 🚀 Key Features

- **System Metrics**: Automatic collection of CPU, memory, disk, and network metrics
- **Runtime Metrics**: Go runtime information (goroutines, heap, GC)
- **Custom Metrics**: Gauges, counters, histograms, and timers
- **Health Checks**: Automatic system health monitoring
- **HTTP Endpoints**: Prometheus-compatible metrics export
- **Alerting**: Threshold-based alerts
- **JSON API**: Complete metrics data for integrations

## 📊 System Metrics

### CPU, Memory, and Disk Monitoring

```lua
-- Get current CPU usage
local cpu_usage = metrics.system_cpu()
log.info("CPU Usage: " .. string.format("%.1f%%", cpu_usage))

-- Get memory information
local memory_info = metrics.system_memory()
log.info("Memory: " .. string.format("%.1f%% (%.0f/%.0f MB)", 
    memory_info.percent, memory_info.used_mb, memory_info.total_mb))

-- Get disk usage
local disk_info = metrics.system_disk("/")
log.info("Disk: " .. string.format("%.1f%% (%.1f/%.1f GB)", 
    disk_info.percent, disk_info.used_gb, disk_info.total_gb))

-- Check specific disk path
local var_disk = metrics.system_disk("/var")
log.info("Var disk usage: " .. string.format("%.1f%%", var_disk.percent))
```

### Runtime Information

```lua
-- Get Go runtime metrics
local runtime = metrics.runtime_info()
log.info("Runtime Information:")
log.info("  Goroutines: " .. runtime.goroutines)
log.info("  CPU cores: " .. runtime.num_cpu)
log.info("  Heap allocated: " .. string.format("%.1f MB", runtime.heap_alloc_mb))
log.info("  Heap system: " .. string.format("%.1f MB", runtime.heap_sys_mb))
log.info("  GC cycles: " .. runtime.num_gc)
log.info("  Go version: " .. runtime.go_version)
```

## 📈 Custom Metrics

### Gauge Metrics (Current Values)

```lua
-- Set simple gauge values
metrics.gauge("cpu_temperature", 65.4)
metrics.gauge("active_connections", 142)
metrics.gauge("queue_size", 23)

-- Set gauge with tags
metrics.gauge("memory_usage", memory_percent, {
    server = "web-01",
    environment = "production",
    region = "us-east-1"
})

-- Update deployment status
metrics.gauge("deployment_progress", 75.5, {
    app = "frontend",
    version = "v2.1.0"
})
```

### Counter Metrics (Incremental Values)

```lua
-- Increment counters
local total_requests = metrics.counter("http_requests_total", 1)
local error_count = metrics.counter("http_errors_total", 1, {
    status_code = "500",
    endpoint = "/api/users"
})

-- Bulk increment
local processed = metrics.counter("messages_processed", 50, {
    queue = "user_notifications",
    priority = "high"
})

log.info("Total requests processed: " .. total_requests)
```

### Histogram Metrics (Value Distributions)

```lua
-- Record response times
metrics.histogram("response_time_ms", 245.6, {
    endpoint = "/api/users",
    method = "GET"
})

-- Record payload sizes
metrics.histogram("payload_size_bytes", 1024, {
    content_type = "application/json"
})

-- Record batch sizes
metrics.histogram("batch_size", 150, {
    operation = "bulk_insert",
    table = "user_events"
})
```

### Timer Metrics (Function Execution Time)

```lua
-- Time function execution automatically
local duration = metrics.timer("database_query", function()
    -- Simulate database query
    local result = exec.run("sleep 0.1")
    return result
end, {
    query_type = "select",
    table = "users"
})

log.info("Database query took: " .. string.format("%.2f ms", duration))

-- Time complex operations
local processing_time = metrics.timer("data_processing", function()
    -- Process large dataset
    local data = {}
    for i = 1, 100000 do
        data[i] = math.sqrt(i) * 2.5
    end
    return #data
end, {
    operation = "mathematical_computation",
    size = "large"
})

log.info("Data processing completed in: " .. string.format("%.2f ms", processing_time))
```

## 🏥 Health Monitoring

### Automatic Health Status

```lua
-- Get comprehensive health status
local health = metrics.health_status()
log.info("Overall Health Status: " .. health.overall)

-- Check individual components
local components = {"cpu", "memory", "disk"}
for _, component in ipairs(components) do
    local comp_info = health[component]
    if comp_info then
        local status_icon = "✅"
        if comp_info.status == "warning" then
            status_icon = "⚠️"
        elseif comp_info.status == "critical" then
            status_icon = "❌"
        end
        
        log.info(string.format("  %s %s: %.1f%% (%s)", 
            status_icon, component:upper(), comp_info.usage, comp_info.status))
    end
end
```

### Custom Health Checks

```lua
-- Create health check function
function check_application_health()
    local health_score = 100
    local issues = {}
    
    -- Check database connectivity
    local db_result = exec.run("pg_isready -h localhost -p 5432")
    if db_result ~= "" then
        health_score = health_score - 20
        table.insert(issues, "Database connection failed")
    end
    
    -- Check disk space
    local disk = metrics.system_disk("/")
    if disk.percent > 90 then
        health_score = health_score - 30
        table.insert(issues, "Disk space critical: " .. string.format("%.1f%%", disk.percent))
    end
    
    -- Check memory usage
    local memory = metrics.system_memory()
    if memory.percent > 85 then
        health_score = health_score - 25
        table.insert(issues, "Memory usage high: " .. string.format("%.1f%%", memory.percent))
    end
    
    -- Record health score
    metrics.gauge("application_health_score", health_score)
    
    if health_score < 70 then
        metrics.alert("application_health", {
            level = "warning",
            message = "Application health degraded: " .. table.concat(issues, ", "),
            score = health_score
        })
    end
    
    return health_score >= 70
end

-- Use in tasks
local health_check = task("health_check")
    :description("Monitor application health status")
    :command(function(this, params)
        local healthy = check_application_health()
        return healthy, healthy and "System healthy" or "System health issues detected"
    end)
    :build()

local health_monitoring = workflow.define("health_monitoring")
    :description("Health monitoring workflow")
    :version("1.0.0")
    :tasks({health_check})
```

## 🚨 Alerting System

### Creating Alerts

```lua
-- Simple threshold alert
local cpu = metrics.system_cpu()
if cpu > 80 then
    metrics.alert("high_cpu_usage", {
        level = "warning",
        message = "CPU usage is high: " .. string.format("%.1f%%", cpu),
        threshold = 80,
        value = cpu,
        severity = "medium"
    })
end

-- Complex alert with multiple conditions
local memory = metrics.system_memory()
local disk = metrics.system_disk()

if memory.percent > 90 and disk.percent > 85 then
    metrics.alert("resource_exhaustion", {
        level = "critical",
        message = string.format("Critical resource usage - Memory: %.1f%%, Disk: %.1f%%", 
            memory.percent, disk.percent),
        memory_usage = memory.percent,
        disk_usage = disk.percent,
        recommended_action = "Scale up resources immediately"
    })
end

-- Application-specific alerts
local queue_size = state.get("task_queue_size", 0)
if queue_size > 1000 then
    metrics.alert("queue_backlog", {
        level = "warning", 
        message = "Task queue backlog detected: " .. queue_size .. " items",
        queue_size = queue_size,
        estimated_processing_time = queue_size * 2 .. " seconds"
    })
end
```

## 🔍 Metrics Management

### Retrieving Custom Metrics

```lua
-- Get specific custom metric
local cpu_metric = metrics.get_custom("cpu_temperature")
if cpu_metric then
    log.info("CPU Temperature metric: " .. data.to_json(cpu_metric))
end

-- List all custom metrics
local all_metrics = metrics.list_custom()
log.info("Total custom metrics: " .. #all_metrics)
for i, metric_name in ipairs(all_metrics) do
    log.info("  " .. i .. ". " .. metric_name)
end
```

### Performance Monitoring Example

```lua
local monitor_api_performance = task("monitor_api_performance")
    :description("Monitor API performance with detailed metrics")
    :command(function(this, params)
        -- Start monitoring session
        log.info("Starting API performance monitoring...")

        -- Simulate API calls and measure performance
        for i = 1, 10 do
            local api_time = metrics.timer("api_call_" .. i, function()
                -- Simulate API call
                exec.run("curl -s -o /dev/null -w '%{time_total}' https://api.example.com/health")
            end, {
                endpoint = "health",
                call_number = tostring(i)
            })

            -- Record response time
            metrics.histogram("api_response_time", api_time, {
                endpoint = "health"
            })

            -- Check if response time is acceptable
            if api_time > 1000 then -- 1 second
                metrics.counter("slow_api_calls", 1, {
                    endpoint = "health"
                })

                metrics.alert("slow_api_response", {
                    level = "warning",
                    message = string.format("Slow API response: %.2f ms", api_time),
                    response_time = api_time,
                    threshold = 1000
                })
            end

            -- Brief delay between calls
            exec.run("sleep 0.1")
        end

        -- Get summary statistics
        local system_health = metrics.health_status()
        log.info("System health after API tests: " .. system_health.overall)

        return true, "API performance monitoring completed"
    end)
    :build()

local performance_monitoring = workflow.define("performance_monitoring")
    :description("Performance monitoring workflow")
    :version("1.0.0")
    :tasks({monitor_api_performance})
```

## 🌐 HTTP Endpoints

The metrics module automatically exposes HTTP endpoints for external monitoring systems:

### Prometheus Format (`/metrics`)
```bash
# Access Prometheus-compatible metrics
curl http://agent:8080/metrics

# Example output:
# sloth_agent_cpu_usage_percent 15.4
# sloth_agent_memory_usage_mb 2048.5
# sloth_agent_disk_usage_percent 67.2
# sloth_agent_tasks_total 142
```

### JSON Format (`/metrics/json`)
```bash
# Get complete metrics in JSON format
curl http://agent:8080/metrics/json

# Example response:
{
  "agent_name": "myagent1",
  "timestamp": "2024-01-15T10:30:00Z",
  "system": {
    "cpu_usage_percent": 15.4,
    "memory_usage_mb": 2048.5,
    "disk_usage_percent": 67.2
  },
  "runtime": {
    "num_goroutines": 25,
    "heap_alloc_mb": 45.2
  },
  "custom": {
    "api_response_time": {...},
    "deployment_progress": 85.5
  }
}
```

### Health Check (`/health`)
```bash
# Check agent health status
curl http://agent:8080/health

# Example response:
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "checks": {
    "cpu": {"usage": 15.4, "status": "healthy"},
    "memory": {"usage": 45.8, "status": "healthy"},
    "disk": {"usage": 67.2, "status": "healthy"}
  }
}
```

## 📋 API Reference

### System Metrics
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `metrics.system_cpu()` | - | usage: number | Get current CPU usage percentage |
| `metrics.system_memory()` | - | info: table | Get memory usage information |
| `metrics.system_disk(path?)` | path?: string | info: table | Get disk usage for path (default: "/") |
| `metrics.runtime_info()` | - | info: table | Get Go runtime information |

### Custom Metrics
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `metrics.gauge(name, value, tags?)` | name: string, value: number, tags?: table | success: boolean | Set gauge metric |
| `metrics.counter(name, increment?, tags?)` | name: string, increment?: number, tags?: table | new_value: number | Increment counter |
| `metrics.histogram(name, value, tags?)` | name: string, value: number, tags?: table | success: boolean | Record histogram value |
| `metrics.timer(name, function, tags?)` | name: string, func: function, tags?: table | duration: number | Time function execution |

### Health and Monitoring
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `metrics.health_status()` | - | status: table | Get comprehensive health status |
| `metrics.alert(name, data)` | name: string, data: table | success: boolean | Create alert |

### Utilities
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `metrics.get_custom(name)` | name: string | metric: table \| nil | Get custom metric by name |
| `metrics.list_custom()` | - | names: table | List all custom metric names |

## 🎯 Best Practices

1. **Use appropriate metric types** - gauges for current values, counters for totals, histograms for distributions
2. **Add meaningful tags** to categorize and filter metrics
3. **Set reasonable alert thresholds** to avoid alert fatigue
4. **Monitor performance impact** of extensive metrics collection
5. **Use timers for performance-critical operations** to identify bottlenecks
6. **Implement health checks** for all critical system components
7. **Export metrics to external systems** like Prometheus for long-term storage

The **Metrics & Monitoring** module provides comprehensive observability for your distributed sloth-runner environment! 📊🚀