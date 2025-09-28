# 🛡️ Reliability Module

The **Reliability** module provides enterprise-grade reliability patterns including circuit breakers, retry logic with exponential backoff, and failure handling strategies. These patterns help build resilient systems that can gracefully handle failures and recover automatically.

## 🚀 Key Features

- **Circuit Breaker Pattern**: Prevents cascading failures by stopping calls to failing services
- **Retry Logic**: Configurable retry strategies with backoff algorithms
- **Failure Tracking**: Persistent failure state across task executions
- **Multiple Strategies**: Fixed delay, exponential backoff, linear backoff, custom
- **Jitter Support**: Randomization to prevent thundering herd problems
- **State Integration**: Uses state module for persistent failure tracking
- **Callback Support**: Custom callbacks for retry and state change events

## 📋 Basic Usage

### Simple Retry

```lua
-- Retry a function up to 3 times with 1 second initial delay
local result = reliability.retry(3, 1, function()
    -- Your potentially failing code here
    if math.random() > 0.7 then
        return "Success!"
    else
        return nil, "Random failure"
    end
end)

if result then
    log.info("Operation succeeded: " .. result)
else 
    log.error("All retries failed")
end
```

### Advanced Retry Configuration

```lua
local config = {
    max_attempts = 5,
    initial_delay = 0.5,  -- 500ms
    max_delay = 10,       -- 10 seconds max
    strategy = reliability.strategy.EXPONENTIAL_BACKOFF,
    multiplier = 2.0,
    jitter = true,
    on_retry = function(attempt, delay, error)
        log.warn("Retry attempt " .. attempt .. " in " .. delay .. "s: " .. error)
    end
}

local result = reliability.retry_with_config(config, function()
    -- Your code here
    return call_external_service()
end)
```

### Circuit Breaker

```lua
local cb_config = {
    max_failures = 3,     -- Open after 3 failures
    timeout = 30,         -- Wait 30 seconds before trying half-open
    success_threshold = 2, -- Need 2 successes to close circuit
    on_state_change = function(from_state, to_state)
        log.info("Circuit breaker: " .. from_state .. " -> " .. to_state)
    end
}

local result = reliability.circuit_breaker("external_api", cb_config, function()
    -- Call that might fail
    return http.get("https://api.example.com/data")
end)
```

## 🔄 Retry Strategies

### Available Strategy Types

```lua
-- Fixed delay between retries
reliability.strategy.FIXED_DELAY

-- Exponential backoff (delay doubles each time)
reliability.strategy.EXPONENTIAL_BACKOFF  

-- Linear backoff (delay increases linearly)
reliability.strategy.LINEAR_BACKOFF

-- Custom delay function
reliability.strategy.CUSTOM_BACKOFF
```

### Custom Delay Function

```lua
local config = {
    max_attempts = 5,
    strategy = reliability.strategy.CUSTOM_BACKOFF,
    custom_delay = function(attempt)
        -- Custom fibonacci-like delays
        if attempt == 1 then return 1 end
        if attempt == 2 then return 1 end
        return (attempt - 1) + (attempt - 2)
    end
}
```

## ⚡ Circuit Breaker States

### State Transitions

- **Closed** → **Open**: After max_failures consecutive failures
- **Open** → **Half-Open**: After timeout period expires  
- **Half-Open** → **Closed**: After success_threshold successes
- **Half-Open** → **Open**: After any failure

### Monitoring Circuit State

```lua
-- Get current statistics
local stats = reliability.get_circuit_stats("my_service")
if stats then
    log.info("Circuit state: " .. stats.state)
    log.info("Total requests: " .. stats.requests)
    log.info("Success rate: " .. (stats.total_success / stats.requests * 100) .. "%")
end

-- List all circuit breakers
local circuits = reliability.list_circuits()
for _, name in ipairs(circuits) do
    log.info("Circuit: " .. name)
end

-- Reset circuit breaker
reliability.reset_circuit("my_service")
```

## 🔗 Integration with State Module

### Persistent Failure Tracking

```lua
-- Track failures across task executions
local service_name = "payment_service"
local failure_key = "failures:" .. service_name

local function make_payment_call()
    local success = make_api_call()
    
    if success then
        -- Reset failure count on success
        state.set(failure_key, "0")
        return true
    else
        -- Increment failure counter
        local failures = state.increment(failure_key, 1)
        
        -- Circuit break if too many failures
        if failures >= 5 then
            return nil, "Service circuit opened - too many failures"
        end
        
        return nil, "Temporary service failure"
    end
end

-- Use with retry
local result = reliability.retry(3, 2, make_payment_call)
```

### Distributed Lock with Retry

```lua
-- Combine distributed locking with retry logic
local retry_config = {
    max_attempts = 5,
    initial_delay = 0.5,
    strategy = reliability.strategy.LINEAR_BACKOFF
}

local result = reliability.retry_with_config(retry_config, function()
    -- Try to acquire distributed lock
    if not state.try_lock("critical_resource", 10) then
        return nil, "Could not acquire lock"
    end
    
    -- Do critical work
    local work_result = perform_critical_operation()
    
    -- Release lock
    state.unlock("critical_resource")
    
    return work_result
end)
```

## 📊 Advanced Patterns

### Combine Multiple Patterns

```lua
-- Deployment with circuit breaker, retry, and state tracking
local deployment_steps = {"validate", "backup", "deploy", "verify"}

for _, step in ipairs(deployment_steps) do
    local step_result = reliability.retry_with_config({
        max_attempts = 3,
        initial_delay = 1,
        strategy = reliability.strategy.EXPONENTIAL_BACKOFF,
        on_retry = function(attempt, delay, error)
            state.append("deployment_log", 
                step .. " retry " .. attempt .. ": " .. error, "\n")
        end
    }, function()
        return reliability.circuit_breaker("deployment_service", {
            max_failures = 2,
            timeout = 30,
            on_state_change = function(from, to)
                state.set("deployment_cb_state", to)
            end
        }, function()
            return execute_deployment_step(step)
        end)
    end)
    
    if not step_result then
        state.set("deployment_status", "failed_at_" .. step)
        return false, "Deployment failed at: " .. step
    end
    
    -- Update progress
    local progress = math.floor((step_index / #deployment_steps) * 100)
    state.set("deployment_progress", progress)
end

state.set("deployment_status", "completed")
```

### Health Check with Backoff

```lua
-- Health check with exponential backoff
local health_config = {
    max_attempts = 10,
    initial_delay = 1,
    max_delay = 60,
    strategy = reliability.strategy.EXPONENTIAL_BACKOFF,
    multiplier = 1.5,
    jitter = true
}

local health_status = reliability.retry_with_config(health_config, function()
    local response = http.get("http://localhost:8080/health")
    
    if response.status == 200 then
        return response.body
    else
        return nil, "Health check failed: " .. response.status
    end
end)
```

## 🎛️ Configuration Reference

### Retry Configuration

```lua
{
    max_attempts = 3,           -- Maximum retry attempts
    initial_delay = 1,          -- Initial delay in seconds
    max_delay = 30,             -- Maximum delay in seconds  
    strategy = "exponential",   -- Retry strategy
    multiplier = 2.0,           -- Backoff multiplier
    jitter = true,              -- Add random jitter
    on_retry = function(attempt, delay, error)
        -- Retry callback
    end
}
```

### Circuit Breaker Configuration

```lua
{
    max_failures = 5,           -- Failures before opening
    timeout = 60,               -- Seconds before half-open
    success_threshold = 1,      -- Successes needed to close
    on_state_change = function(from, to)
        -- State change callback  
    end
}
```

## 🚨 Error Handling

### Custom Error Predicates

```lua
-- Retry only on specific errors
local config = {
    max_attempts = 3,
    should_retry = function(error)
        -- Only retry on network errors
        return string.find(error, "network") or string.find(error, "timeout")
    end
}
```

### Error Types

- **RetryableError**: Explicitly marked as retryable
- **NonRetryableError**: Should not be retried
- **CircuitBreakerError**: Circuit is open, don't retry immediately

## 📈 Monitoring and Observability

### Metrics Collection

```lua
-- Circuit breaker metrics
local cb_stats = reliability.get_circuit_stats("service_name")
-- Returns: requests, total_success, total_failures, consecutive_success, 
--          consecutive_failures, state, last_success_time, last_failure_time

-- State-based metrics
local failure_count = tonumber(state.get("service_failures", "0"))
local success_rate = calculate_success_rate()

-- Log metrics
log.info("Service metrics", {
    circuit_state = cb_stats.state,
    failure_count = failure_count,
    success_rate = success_rate
})
```

The reliability module provides the foundation for building resilient, fault-tolerant automation workflows that can handle failures gracefully and recover automatically.