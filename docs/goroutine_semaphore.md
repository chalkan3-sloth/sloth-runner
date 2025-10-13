# Goroutine Semaphore Documentation

## Overview

O módulo `goroutine` fornece suporte completo para **Semaphores** (contadores de recursos), permitindo limitar o acesso concorrente a recursos compartilhados. Semaphores são fundamentais para controle de taxa (rate limiting), pools de conexão, e gerenciamento de recursos limitados.

## Why Semaphores?

Quando você tem recursos limitados (como conexões de banco de dados, slots de memória, ou limites de API), precisa controlar quantas goroutines podem acessar esses recursos simultaneamente. **Semaphores resolvem isso** atuando como um contador de "tokens" disponíveis.

**Exemplo prático**: Se você tem um banco de dados que suporta no máximo 10 conexões simultâneas, um semaphore com capacidade 10 garante que nunca mais de 10 goroutines tentarão conectar ao mesmo tempo.

## Semaphore vs Mutex vs Channel

| Feature | Semaphore | Mutex | Channel |
|---------|-----------|-------|---------|
| **Purpose** | Limit concurrent access | Protect data | Communicate |
| **Capacity** | N concurrent users | 1 user only | Message queue |
| **Use case** | Rate limiting, resource pools | Exclusive access | Data passing |
| **Blocking** | ✅ Yes (acquire) | ✅ Yes (lock) | ✅ Yes (send/receive) |
| **Non-blocking** | ✅ Yes (try_acquire) | ✅ Yes (try_lock) | ✅ Yes (try_send/receive) |

**When to use Semaphore:**
- ✅ Limiting concurrent connections (database, API)
- ✅ Rate limiting requests
- ✅ Managing resource pools
- ✅ Throttling concurrent operations
- ✅ Worker pool size control

**When to use Mutex:**
- ✅ Protecting shared data (exclusive access)
- ✅ Critical sections

**When to use Channel:**
- ✅ Communication between goroutines
- ✅ Passing data ownership

## API Reference

### Creating a Semaphore

```lua
-- Create semaphore with capacity N
local sem = goroutine.semaphore(N)

-- Examples:
local db_pool = goroutine.semaphore(10)    -- Max 10 database connections
local rate_limiter = goroutine.semaphore(5) -- Max 5 requests at a time
local worker_pool = goroutine.semaphore(4)  -- Max 4 concurrent workers
```

**Parameters:**
- `capacity` (number): Maximum number of concurrent acquisitions (must be > 0)

**Returns:**
- Semaphore object with methods: `acquire()`, `release()`, `try_acquire()`, `available()`, `capacity()`

### Methods

#### sem:acquire()
Acquires a token from the semaphore. **Blocks** if no tokens are available.

```lua
local sem = goroutine.semaphore(3)

sem:acquire()
-- Critical section: using resource
sem:release()
```

**Behavior:**
- Blocks until a token is available
- Decrements available token count
- Thread-safe

**Use when:** You **must** acquire the resource and are willing to wait.

---

#### sem:release()
Returns a token to the semaphore.

```lua
sem:acquire()
-- Do work
sem:release()  -- Always release!
```

**Behavior:**
- Increments available token count
- Errors if releasing more tokens than capacity (release without acquire)
- Thread-safe

**Important:** Always pair `acquire()` with `release()` to avoid token leaks.

---

#### sem:try_acquire()
Attempts to acquire a token **without blocking**.

```lua
local sem = goroutine.semaphore(2)

local acquired = sem:try_acquire()
if acquired then
  log.info("Got token!")
  -- Do work
  sem:release()
else
  log.info("No token available, doing alternative work")
  -- Fallback behavior
end
```

**Returns:**
- `true` if token was acquired
- `false` if no tokens available

**Use when:** You want to try acquiring but have a fallback if unavailable.

---

#### sem:available()
Returns the number of currently available tokens.

```lua
local sem = goroutine.semaphore(10)
sem:acquire()
sem:acquire()

log.info("Available: " .. sem:available())  -- Output: Available: 8
```

**Returns:** `number` - count of available tokens

**Use for:** Monitoring, debugging, metrics

---

#### sem:capacity()
Returns the total capacity of the semaphore.

```lua
local sem = goroutine.semaphore(10)
log.info("Capacity: " .. sem:capacity())  -- Output: Capacity: 10
```

**Returns:** `number` - total capacity

**Use for:** Calculating usage percentage, monitoring

## Usage Patterns

### 1. Basic Rate Limiting

**Problem**: Limit concurrent operations to prevent resource exhaustion.

**Solution**: Use semaphore to control concurrency.

```lua
workflow("rate-limit-demo")
    :task("limit-concurrent-requests", function()
        local max_concurrent = 5
        local sem = goroutine.semaphore(max_concurrent)
        local wg = goroutine.wait_group()

        -- Spawn 20 requests, but only 5 run concurrently
        wg:add(20)
        for i = 1, 20 do
            goroutine.spawn(function()
                sem:acquire()
                log.info("Request " .. i .. " processing...")
                goroutine.sleep(100)  -- Simulate work
                sem:release()
                wg:done()
            end)
        end

        wg:wait()
        return true
    end)
```

**Result**: Only 5 requests run at any time, preventing overload.

---

### 2. Database Connection Pool

**Problem**: Database supports max N connections, need to limit concurrent queries.

**Solution**: Semaphore acts as connection pool manager.

```lua
local max_connections = 10
local db_sem = goroutine.semaphore(max_connections)

local execute_query = function(query)
    db_sem:acquire()  -- Wait for connection slot
    log.info("Executing query (Available: " .. db_sem:available() .. ")")

    -- Execute query
    local result = database.query(query)

    db_sem:release()  -- Return connection slot
    return result
end

-- Many goroutines can call execute_query
-- Only 10 will execute concurrently
for i = 1, 100 do
    goroutine.spawn(function()
        execute_query("SELECT * FROM users WHERE id = " .. i)
    end)
end
```

---

### 3. Try Acquire with Fallback

**Problem**: Want to try acquiring resource but have alternative if unavailable.

**Solution**: Use `try_acquire()` for non-blocking attempt.

```lua
local sem = goroutine.semaphore(5)

local process_with_fallback = function(task_id)
    local acquired = sem:try_acquire()

    if acquired then
        log.info("Task " .. task_id .. ": Using premium resource")
        -- Use expensive resource
        premium_service.process(task_id)
        sem:release()
    else
        log.info("Task " .. task_id .. ": Premium busy, using standard")
        -- Fallback to standard service
        standard_service.process(task_id)
    end
end
```

**Use cases:**
- Graceful degradation
- Fast path / slow path
- Premium / standard tiers

---

### 4. Worker Pool

**Problem**: Limit number of concurrent workers processing jobs.

**Solution**: Semaphore controls worker count.

```lua
local max_workers = 4
local worker_sem = goroutine.semaphore(max_workers)
local jobs = goroutine.channel(100)
local wg = goroutine.wait_group()

-- Worker function
local worker = function(job_id)
    worker_sem:acquire()
    log.info("Worker processing job " .. job_id)

    -- Process job
    local result = process_job(job_id)

    worker_sem:release()
    return result
end

-- Spawn jobs
wg:add(50)
for i = 1, 50 do
    jobs:send(i)
end

-- Process jobs (max 4 concurrent)
for i = 1, 50 do
    local job_id, ok = jobs:receive()
    if ok then
        goroutine.spawn(function()
            worker(job_id)
            wg:done()
        end)
    end
end

wg:wait()
```

---

### 5. Timeout Pattern with Try Acquire

**Problem**: Need to acquire resource but give up after timeout.

**Solution**: Loop with `try_acquire()` and time tracking.

```lua
local acquire_with_timeout = function(sem, timeout_ms)
    local elapsed = 0

    while elapsed < timeout_ms do
        local acquired = sem:try_acquire()
        if acquired then
            return true
        end

        goroutine.sleep(50)
        elapsed = elapsed + 50
    end

    log.info("Timeout: Could not acquire resource after " .. timeout_ms .. "ms")
    return false
end

-- Usage
local sem = goroutine.semaphore(5)

if acquire_with_timeout(sem, 1000) then
    -- Got resource within 1 second
    do_work()
    sem:release()
else
    -- Timeout: do fallback
    do_fallback()
end
```

---

### 6. API Rate Limiting

**Problem**: API allows N requests per time window.

**Solution**: Semaphore with periodic refill.

```lua
local rate_limit = 10  -- 10 requests per second
local limiter = goroutine.semaphore(rate_limit)
local refill_running = true

-- Refill tokens every second
goroutine.spawn(function()
    while refill_running do
        goroutine.sleep(1000)

        -- Consume all tokens and re-add them (reset to full capacity)
        for i = 1, rate_limit do
            local acquired = limiter:try_acquire()
            if not acquired then
                break
            end
        end
        log.info("Rate limiter refilled")
    end
end)

-- Make API request
local make_request = function(req_id)
    limiter:acquire()  -- Block if rate limit exceeded
    log.info("Request " .. req_id .. " executing")

    -- Make API call
    api.call(req_id)
end

-- Many requests will be throttled
for i = 1, 100 do
    goroutine.spawn(function()
        make_request(i)
    end)
end
```

---

### 7. Monitoring Semaphore Usage

**Problem**: Need to track resource usage and availability.

**Solution**: Monitor goroutine reading `available()` and `capacity()`.

```lua
local sem = goroutine.semaphore(10)
local monitoring = true

-- Monitor goroutine
goroutine.spawn(function()
    while monitoring do
        local available = sem:available()
        local capacity = sem:capacity()
        local in_use = capacity - available
        local percent = (in_use / capacity) * 100

        log.info(string.format("Resource usage: %d/%d (%.1f%%)", in_use, capacity, percent))

        goroutine.sleep(500)
    end
end)

-- Workers
for i = 1, 20 do
    goroutine.spawn(function()
        sem:acquire()
        goroutine.sleep(math.random(1000, 3000))
        sem:release()
    end)
end
```

**Use for:**
- Metrics collection
- Alerting on high usage
- Debugging resource bottlenecks

---

### 8. Combining Semaphore with Channels

**Problem**: Coordinate job queue with limited workers.

**Solution**: Channel for jobs, semaphore for worker limit.

```lua
local sem = goroutine.semaphore(3)
local jobs = goroutine.channel(10)
local results = goroutine.channel(10)

-- Producer
goroutine.spawn(function()
    for i = 1, 10 do
        jobs:send("Job-" .. i)
    end
    jobs:close()
end)

-- Worker pool
goroutine.spawn(function()
    while true do
        local job, ok = jobs:receive()
        if not ok then break end

        goroutine.spawn(function()
            sem:acquire()
            log.info("Processing " .. job)
            goroutine.sleep(200)
            results:send(job .. " done")
            sem:release()
        end)
    end
end)

-- Consumer
for i = 1, 10 do
    local result, ok = results:receive()
    if ok then
        log.info(result)
    end
end
```

## Examples

O arquivo `examples/goroutine_semaphore_example.sloth` contém 8 exemplos práticos:

1. **Basic Semaphore** - Rate limiting básico
2. **Connection Pool** - Pool de conexões de banco de dados
3. **Try Acquire** - Aquisição não-bloqueante
4. **Worker Pool** - Pool de workers limitado
5. **Rate Limiter** - Limitador de requisições com refill
6. **Semaphore with Channels** - Coordenação com channels
7. **Resource Pool Timeout** - Aquisição com timeout
8. **Monitor Semaphore** - Monitoramento de uso

### Running Examples

```bash
# Basic semaphore
sloth-runner run basic_semaphore --file examples/goroutine_semaphore_example.sloth --yes

# Connection pool
sloth-runner run connection_pool --file examples/goroutine_semaphore_example.sloth --yes

# Worker pool
sloth-runner run worker_pool --file examples/goroutine_semaphore_example.sloth --yes

# Monitoring
sloth-runner run monitor_semaphore --file examples/goroutine_semaphore_example.sloth --yes
```

## Best Practices

### 1. Always Release

**BAD**:
```lua
sem:acquire()
-- Do work
-- Forgot to release!
```

**GOOD**:
```lua
sem:acquire()
-- Do work
sem:release()  -- Always release!
```

**TIP**: Tokens are finite. Forgetting to release causes token leaks.

---

### 2. Release Exactly Once

**BAD** (causes error):
```lua
sem:acquire()
sem:release()
sem:release()  -- ERROR: release without acquire
```

**GOOD**:
```lua
sem:acquire()
-- Do work
sem:release()  -- Release exactly once
```

---

### 3. Use Try Acquire for Non-Critical Operations

**BAD** (blocks unnecessarily):
```lua
sem:acquire()  -- Blocks even if not critical
nice_to_have_operation()
sem:release()
```

**GOOD** (non-blocking):
```lua
local acquired = sem:try_acquire()
if acquired then
    nice_to_have_operation()
    sem:release()
else
    log.info("Skipping non-critical operation")
end
```

---

### 4. Choose Appropriate Capacity

**Too Low**:
```lua
local sem = goroutine.semaphore(1)  -- Too restrictive, poor throughput
```

**Too High**:
```lua
local sem = goroutine.semaphore(1000)  -- No effective limiting
```

**Just Right**:
```lua
-- Based on actual resource limits
local max_db_connections = 20
local sem = goroutine.semaphore(max_db_connections)
```

**Guideline**: Set capacity to **actual resource limit** (not arbitrary number).

---

### 5. Monitor in Production

```lua
-- Periodically log usage
goroutine.spawn(function()
    while true do
        local usage = sem:capacity() - sem:available()
        log.info("Semaphore usage: " .. usage .. "/" .. sem:capacity())
        goroutine.sleep(60000)  -- Every minute
    end
end)
```

**Use for:**
- Detecting bottlenecks
- Capacity planning
- Alert on high usage

---

### 6. Combine with Timeout

**BAD** (infinite wait):
```lua
sem:acquire()  -- Might wait forever
```

**GOOD** (with timeout):
```lua
if acquire_with_timeout(sem, 5000) then
    -- Got resource within 5 seconds
    do_work()
    sem:release()
else
    -- Timeout: do fallback
    log.error("Could not acquire resource, using fallback")
    do_fallback()
end
```

## Common Pitfalls

### 1. Token Leak

**Problem**: Forgetting to release tokens.

```lua
sem:acquire()
if error_condition then
    return false  -- ⚠️ FORGOT TO RELEASE!
end
sem:release()
```

**Solution**: Always release, even on error paths.

```lua
sem:acquire()
local success = pcall(function()
    -- Do work
    if error_condition then
        error("Something went wrong")
    end
end)
sem:release()  -- Always releases

if not success then
    log.error("Work failed but token released")
end
```

---

### 2. Over-releasing

**Problem**: Calling release more times than acquire.

```lua
sem:acquire()
sem:release()
sem:release()  -- ⚠️ ERROR: release without acquire
```

**Solution**: Track acquire/release pairs carefully.

---

### 3. Using Semaphore for Mutual Exclusion

**BAD** (use mutex instead):
```lua
local sem = goroutine.semaphore(1)  -- Don't do this!
sem:acquire()
shared_data = shared_data + 1
sem:release()
```

**GOOD** (use mutex):
```lua
local mu = goroutine.mutex()
mu:lock()
shared_data = shared_data + 1
mu:unlock()
```

**Why**: Mutex has better semantics for exclusive access. Semaphore is for counting resources.

---

### 4. Deadlock with Multiple Semaphores

**BAD**:
```lua
-- Goroutine 1
sem1:acquire()
sem2:acquire()  -- Might deadlock
sem2:release()
sem1:release()

-- Goroutine 2
sem2:acquire()
sem1:acquire()  -- Might deadlock
sem1:release()
sem2:release()
```

**GOOD**: Always acquire in same order.

```lua
-- Both goroutines acquire in same order
sem1:acquire()
sem2:acquire()
-- Critical section
sem2:release()
sem1:release()
```

## Implementation Details

### Thread Safety

Semaphores são implementados usando channels do Go, garantindo thread-safety nativa.

### Memory Management

Semaphores são automaticamente garbage collected quando não há mais referências.

### Performance

- **Acquire**: O(1) - channel receive
- **Release**: O(1) - channel send
- **Try Acquire**: O(1) - non-blocking select
- **Available**: O(1) - channel len
- **Capacity**: O(1) - field access

## Comparison with Go

O Sloth Runner semaphore segue a semântica do package `golang.org/x/sync/semaphore`:

| Go (golang.org/x/sync/semaphore) | Sloth Runner |
|----------------------------------|--------------|
| `semaphore.NewWeighted(n)` | `goroutine.semaphore(n)` |
| `sem.Acquire(ctx, 1)` | `sem:acquire()` |
| `sem.Release(1)` | `sem:release()` |
| `sem.TryAcquire(1)` | `sem:try_acquire()` |
| N/A | `sem:available()` |
| N/A | `sem:capacity()` |

**Diferenças**:
- Go usa `context` para timeout, Sloth Runner usa loop manual
- Sloth Runner adiciona `available()` e `capacity()` para monitoramento

## When NOT to Use Semaphore

### Don't Use Semaphore For:

1. **Protecting shared data** → Use **Mutex** instead
   ```lua
   -- BAD
   local sem = goroutine.semaphore(1)
   sem:acquire()
   counter = counter + 1
   sem:release()

   -- GOOD
   local mu = goroutine.mutex()
   mu:lock()
   counter = counter + 1
   mu:unlock()
   ```

2. **Communication between goroutines** → Use **Channel** instead
   ```lua
   -- BAD
   local sem = goroutine.semaphore(1)
   -- Try to pass data via global + semaphore

   -- GOOD
   local ch = goroutine.channel()
   ch:send(data)
   ```

3. **One-time initialization** → Use **Once** (future feature)

4. **Waiting for goroutines** → Use **WaitGroup** instead

## Troubleshooting

### Problem: Goroutines Never Complete

**Symptom**: Program hangs, goroutines waiting forever.

**Cause**: Token leak - acquired but never released.

**Solution**:
```lua
-- Wrap in pcall to ensure release
sem:acquire()
local success = pcall(function()
    -- Work that might error
end)
sem:release()  -- Always releases
```

---

### Problem: "release without acquire" Error

**Symptom**: Error when calling `release()`.

**Cause**: Calling release more times than acquire.

**Solution**: Ensure acquire/release pairs match 1:1.

---

### Problem: Low Throughput

**Symptom**: Operations are slow despite having capacity.

**Cause**: Semaphore capacity too low.

**Solution**:
```lua
-- Check usage
local usage_percent = (sem:capacity() - sem:available()) / sem:capacity() * 100
if usage_percent > 90 then
    log.warn("Semaphore at " .. usage_percent .. "% usage, consider increasing capacity")
end
```

## Conclusion

Semaphores são ferramentas poderosas para:
- ✅ Limitar concorrência
- ✅ Gerenciar recursos finitos
- ✅ Rate limiting
- ✅ Connection pooling
- ✅ Worker pools

**Lembre-se**:
- **Always release** tokens
- **Choose capacity** based on actual limits
- **Use try_acquire** when you have fallback
- **Monitor usage** in production
- **Use mutex** for mutual exclusion (not semaphore with capacity 1)

Para mais informações sobre semaphores:
- [Go sync.Semaphore](https://pkg.go.dev/golang.org/x/sync/semaphore)
- [Wikipedia - Semaphore](https://en.wikipedia.org/wiki/Semaphore_(programming))
