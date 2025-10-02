# 🚀 Goroutine Examples

This directory contains comprehensive examples demonstrating how to use goroutines in Sloth Runner for parallel task execution.

## What are Goroutines?

Goroutines are lightweight threads managed by the Go runtime. Sloth Runner exposes Go's powerful concurrency model to Lua scripts, allowing you to execute multiple operations in parallel within a single task.

## 📁 Available Examples

### 1. **Parallel Deployment** (`parallel_deployment.sloth`)
Deploy your application to multiple servers simultaneously.

**Performance:** 
- Sequential: 6 servers × 30s = 3 minutes
- Parallel: 30 seconds total (6x faster!)

```bash
sloth-runner run -f examples/goroutines/parallel_deployment.sloth
```

**What it demonstrates:**
- Creating multiple goroutines for parallel deployment
- Collecting and processing results from all goroutines
- Comprehensive error handling and reporting
- Real-world deployment scenario

---

### 2. **Parallel Health Check** (`parallel_health_check.sloth`)
Check the health of multiple services concurrently.

**Performance:**
- Sequential: 5 services × 5s = 25 seconds
- Parallel: 5 seconds total (5x faster!)

```bash
sloth-runner run -f examples/goroutines/parallel_health_check.sloth
```

**What it demonstrates:**
- Parallel service monitoring
- System health validation
- Fast failure detection
- Service status aggregation

---

### 3. **Parallel Data Processing** (`parallel_data_processing.sloth`)
Process large datasets by splitting them into parallel chunks.

**Performance:**
- Sequential: 1000 items × 10ms = 10 seconds
- Parallel: ~1 second (10x faster!)

```bash
sloth-runner run -f examples/goroutines/parallel_data_processing.sloth
```

**What it demonstrates:**
- Data chunking for parallel processing
- Result aggregation from multiple goroutines
- Large dataset handling
- Performance optimization techniques

---

## 🎯 Quick Start

### Basic Goroutine Pattern

```lua
local task_example = task("goroutine_example")
    :command(function(this, params)
        local go = require("goroutine")
        
        -- Create goroutines
        local g1 = go.create(function()
            -- Do work in parallel
            return "result 1"
        end)
        
        local g2 = go.create(function()
            -- Do work in parallel
            return "result 2"
        end)
        
        -- Wait for all to complete
        local results = go.wait_all({g1, g2}, 30)
        
        -- Process results
        for _, result in ipairs(results) do
            if result.success then
                log.info("Got result: " .. result.value)
            end
        end
        
        return true, "Completed"
    end)
    :build()
```

## 📚 API Reference

### `go.create(function)`
Creates a new goroutine that executes the given function in parallel.

```lua
local g = go.create(function()
    -- Your parallel work here
    return result
end)
```

**Returns:** Goroutine handle

---

### `go.wait_all(goroutines, timeout)`
Waits for all goroutines to complete or until timeout.

```lua
local results = go.wait_all({g1, g2, g3}, 60)
```

**Parameters:**
- `goroutines`: Table of goroutine handles
- `timeout`: Maximum wait time in seconds

**Returns:** Table of results, each containing:
```lua
{
    success = true/false,  -- Whether goroutine completed successfully
    value = ...,           -- Return value from goroutine
    error = "...",         -- Error message if failed
    duration = 1.5         -- Execution time in seconds
}
```

---

### `go.wait_any(goroutines, timeout)`
Waits for the first goroutine to complete.

```lua
local result = go.wait_any({g1, g2, g3}, 60)
```

**Returns:** Single result (first to complete)

---

## 💡 Best Practices

### ✅ When to Use Goroutines

- **Network operations:** API calls, health checks, remote deployments
- **I/O operations:** File processing, database queries
- **Multiple independent tasks:** Deployments, data processing, monitoring
- **Time-sensitive operations:** When you need results quickly

### ❌ When NOT to Use Goroutines

- **Simple sequential tasks:** No benefit for single operations
- **CPU-bound work:** Go runtime already optimizes CPU usage
- **Very short operations:** Overhead may outweigh benefits
- **Shared state without synchronization:** Can lead to race conditions

---

## 🎓 Learning Path

1. **Start with:** `parallel_deployment.sloth` - Simple, practical example
2. **Then try:** `parallel_health_check.sloth` - System integration
3. **Finally explore:** `parallel_data_processing.sloth` - Advanced patterns

---

## 🔧 Troubleshooting

### Goroutine Timeout
```lua
-- Increase timeout for slow operations
local results = go.wait_all(goroutines, 120)  -- 2 minutes
```

### Handling Errors
```lua
for _, result in ipairs(results) do
    if not result.success then
        log.error("Goroutine failed: " .. result.error)
        -- Handle failure
    end
end
```

### Too Many Goroutines
```lua
-- Process in batches
local batch_size = 10
for i = 1, #items, batch_size do
    local batch_goroutines = {}
    for j = i, math.min(i + batch_size - 1, #items) do
        -- Create goroutine
    end
    local results = go.wait_all(batch_goroutines, 60)
end
```

---

## 📊 Performance Comparison

| Scenario | Sequential | Parallel (Goroutines) | Speedup |
|----------|------------|----------------------|---------|
| 10 API calls @ 2s each | 20s | 2s | 10x |
| 5 deployments @ 30s each | 2.5m | 30s | 5x |
| 1000 items processing | 10s | 1s | 10x |
| Health check 20 services | 1m | 5s | 12x |

---

## 🌟 Advanced Examples

Looking for more complex scenarios? Check out:

- **Multi-region deployment:** Deploy to multiple cloud regions in parallel
- **Distributed testing:** Run test suites across multiple agents
- **Batch processing:** Process thousands of files concurrently
- **Real-time monitoring:** Monitor multiple metrics simultaneously

---

## 🤝 Contributing

Have a great goroutine example? We'd love to see it! Please:

1. Create your example file
2. Add comprehensive comments
3. Include performance metrics
4. Submit a pull request

---

## 📖 Additional Resources

- [Complete Goroutine Documentation](../../docs/modules/goroutine.md)
- [Modern DSL Reference](../../docs/LUA_API.md)
- [Performance Best Practices](../../docs/performance.md)

---

**Questions?** Open an issue or discussion on GitHub!

**Tip:** Start with `parallel_deployment.sloth` - it's the most practical example for getting started with goroutines in Sloth Runner.
