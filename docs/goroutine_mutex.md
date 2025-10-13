# Goroutine Mutex Documentation

## Overview

O módulo `goroutine` fornece suporte completo para **Mutex** (Mutual Exclusion) e **RWMutex** (Read-Write Mutex), permitindo proteger dados compartilhados entre goroutines de forma segura. Mutexes são fundamentais para evitar race conditions em programação concorrente.

## Why Mutexes?

Quando múltiplas goroutines acessam e modificam os mesmos dados simultaneamente, podem ocorrer **race conditions** (condições de corrida), levando a:
- Dados inconsistentes
- Valores corrompidos
- Comportamento imprevisível

**Mutexes resolvem isso** garantindo que apenas uma goroutine por vez possa acessar a seção crítica do código.

## Mutex Types

### 1. **Mutex** - Mutual Exclusion Lock
Um mutex padrão que permite apenas **um** acesso por vez (exclusivo).

### 2. **RWMutex** - Read-Write Mutex
Um mutex especializado que permite:
- **Múltiplos leitores** simultaneamente (read lock)
- **Um escritor** exclusivo por vez (write lock)

## API Reference

### Mutex

#### Creating a Mutex

```lua
local mu = goroutine.mutex()
```

#### Methods

| Method | Description | Blocking | Returns |
|--------|-------------|----------|---------|
| `mu:lock()` | Acquire exclusive lock | ✅ Yes | - |
| `mu:unlock()` | Release lock | ❌ No | - |
| `mu:try_lock()` | Try to acquire lock (non-blocking) | ❌ No | `boolean` |

#### Example

```lua
local mu = goroutine.mutex()

mu:lock()
-- Critical section: only one goroutine at a time
counter = counter + 1
mu:unlock()
```

### RWMutex

#### Creating a RWMutex

```lua
local rwmu = goroutine.rwmutex()
```

#### Methods

| Method | Description | Blocking | Returns |
|--------|-------------|----------|---------|
| `rwmu:lock()` | Acquire exclusive write lock | ✅ Yes | - |
| `rwmu:unlock()` | Release write lock | ❌ No | - |
| `rwmu:rlock()` | Acquire shared read lock | ✅ Yes | - |
| `rwmu:runlock()` | Release read lock | ❌ No | - |
| `rwmu:try_lock()` | Try to acquire write lock (non-blocking) | ❌ No | `boolean` |
| `rwmu:try_rlock()` | Try to acquire read lock (non-blocking) | ❌ No | `boolean` |

#### Example

```lua
local rwmu = goroutine.rwmutex()

-- Multiple readers can read concurrently
rwmu:rlock()
local value = shared_data
rwmu:runlock()

-- Writer has exclusive access
rwmu:lock()
shared_data = new_value
rwmu:unlock()
```

## Usage Patterns

### 1. Protecting Shared Counter

**Problem**: Multiple goroutines incrementing a counter causes race conditions.

**Solution**: Use mutex to protect the counter.

```lua
workflow("counter-demo")
    :task("safe-counter", function()
        local counter = 0
        local mu = goroutine.mutex()
        local wg = goroutine.wait_group()

        -- Spawn 10 goroutines
        wg:add(10)
        for i = 1, 10 do
            goroutine.spawn(function()
                for j = 1, 100 do
                    mu:lock()
                    counter = counter + 1  -- Critical section
                    mu:unlock()
                end
                wg:done()
            end)
        end

        wg:wait()
        log.info("Final counter: " .. counter)  -- Should be 1000
        return true
    end)
```

**Output:**
```
Final counter: 1000  ✓ Correct (with mutex)
```

Without mutex, you'd get unpredictable results like 987, 943, etc.

### 2. Protecting Shared Data Structure

```lua
local shared_map = {}
local mu = goroutine.mutex()

-- Writer goroutine
goroutine.spawn(function()
    mu:lock()
    shared_map["key1"] = "value1"
    shared_map["key2"] = "value2"
    mu:unlock()
end)

-- Reader goroutine
goroutine.spawn(function()
    mu:lock()
    local value = shared_map["key1"]
    mu:unlock()

    if value then
        log.info("Got: " .. value)
    end
end)
```

### 3. Try Lock Pattern (Non-blocking)

Use `try_lock()` when you don't want to block:

```lua
local mu = goroutine.mutex()

local acquired = mu:try_lock()
if acquired then
    log.info("Got the lock!")
    -- Do work
    mu:unlock()
else
    log.info("Lock is busy, doing other work instead")
    -- Do alternative work
end
```

**Use Cases:**
- Avoid blocking when there's alternative work to do
- Implement timeout patterns
- Prevent deadlocks (see below)

### 4. RWMutex for Read-Heavy Workloads

When you have many readers and few writers, use RWMutex for better performance:

```lua
local config = { timeout = 30, retries = 3 }
local rwmu = goroutine.rwmutex()

-- Multiple readers (can run concurrently)
for i = 1, 10 do
    goroutine.spawn(function()
        rwmu:rlock()
        log.info("Timeout: " .. config.timeout)
        rwmu:runlock()
    end)
end

-- Occasional writer (exclusive access)
goroutine.spawn(function()
    rwmu:lock()
    config.timeout = 60  -- Update
    rwmu:unlock()
end)
```

**Performance:**
```
┌─────────────────────────────────────────┐
│  Multiple readers with read locks       │
│  (All run CONCURRENTLY!)                │
│                                         │
│  Reader-1 ──┐                           │
│  Reader-2 ──┼─→ All reading at same time│
│  Reader-3 ──┘                           │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  Writer with write lock                 │
│  (BLOCKS all other access)              │
│                                         │
│  Writer ─────→ Exclusive access         │
│  (readers and other writers wait)       │
└─────────────────────────────────────────┘
```

### 5. Combining Mutex with Channels

Mutexes and channels solve different problems:

- **Channels**: For communication between goroutines
- **Mutexes**: For protecting shared data

You can combine them:

```lua
local shared_map = {}
local mu = goroutine.mutex()
local results = goroutine.channel(10)

-- Worker updates shared map and sends notification
local worker = function(id)
    local key = "worker_" .. id
    local value = id * 100

    -- Protect shared map with mutex
    mu:lock()
    shared_map[key] = value
    mu:unlock()

    -- Notify via channel
    results:send("Worker " .. id .. " done")
end

-- Spawn workers
for i = 1, 5 do
    goroutine.spawn(function() worker(i) end)
end

-- Collect results
for i = 1, 5 do
    local msg, ok = results:receive()
    if ok then
        log.info(msg)
    end
end
```

### 6. Deadlock Prevention with Try Lock

**Deadlock** occurs when goroutines wait for each other in a cycle:

```
Goroutine A: has lock1, wants lock2
Goroutine B: has lock2, wants lock1
→ Deadlock! Both wait forever.
```

**Solution**: Use `try_lock()` with backoff:

```lua
local safe_double_lock = function(name, first, second)
    for attempt = 1, 5 do
        first:lock()

        local got_second = second:try_lock()
        if got_second then
            -- Success! Have both locks
            log.info(name .. ": Got both locks")

            -- Critical section
            -- ...

            second:unlock()
            first:unlock()
            return true
        else
            -- Couldn't get second lock, release first and retry
            first:unlock()
            goroutine.sleep(math.random(10, 50))  -- Random backoff
        end
    end

    log.info(name .. ": Failed to acquire both locks")
    return false
end
```

## Examples

O arquivo `examples/goroutine_mutex_example.sloth` contém 8 exemplos práticos:

1. **Basic Mutex** - Protecting shared counter
2. **Critical Section** - Protecting data structure modifications
3. **Try Lock** - Non-blocking lock attempts
4. **RWMutex Readers** - Multiple concurrent readers
5. **RWMutex Writers** - Exclusive writer access
6. **RWMutex Try Lock** - Non-blocking RWMutex operations
7. **Mutex with Channels** - Combining mutexes and channels
8. **Deadlock Prevention** - Using try_lock to avoid deadlocks

### Running Examples

```bash
# Basic mutex example
sloth-runner run basic_mutex --file examples/goroutine_mutex_example.sloth --yes

# RWMutex example
sloth-runner run rwmutex_readers --file examples/goroutine_mutex_example.sloth --yes

# Deadlock prevention
sloth-runner run deadlock_prevention --file examples/goroutine_mutex_example.sloth --yes
```

## Best Practices

### 1. Always Unlock

**BAD**:
```lua
mu:lock()
-- Do work
-- Forgot to unlock!
```

**GOOD**:
```lua
mu:lock()
-- Do work
mu:unlock()
```

**TIP**: Keep critical sections short!

### 2. Don't Lock Twice

**BAD** (causes deadlock):
```lua
mu:lock()
mu:lock()  -- DEADLOCK! Already locked
mu:unlock()
mu:unlock()
```

**GOOD**:
```lua
mu:lock()
-- Do work
mu:unlock()
```

### 3. Unlock in Reverse Order

When using multiple locks:

**GOOD**:
```lua
mu1:lock()
mu2:lock()
-- Critical section
mu2:unlock()  -- Unlock in reverse order
mu1:unlock()
```

### 4. Use RWMutex for Read-Heavy Workloads

**BAD** (slower):
```lua
local mu = goroutine.mutex()

-- Many readers compete for exclusive lock
for i = 1, 100 do
    goroutine.spawn(function()
        mu:lock()
        local _ = shared_data  -- Just reading
        mu:unlock()
    end)
end
```

**GOOD** (faster):
```lua
local rwmu = goroutine.rwmutex()

-- Many readers can read concurrently
for i = 1, 100 do
    goroutine.spawn(function()
        rwmu:rlock()  -- Shared read lock
        local _ = shared_data
        rwmu:runlock()
    end)
end
```

### 5. Use Try Lock to Avoid Blocking

**BAD** (blocks forever if lock is held):
```lua
mu:lock()  -- Blocks indefinitely
-- Critical section
mu:unlock()
```

**GOOD** (has fallback):
```lua
local acquired = mu:try_lock()
if acquired then
    -- Critical section
    mu:unlock()
else
    -- Do alternative work
    log.info("Lock busy, skipping")
end
```

### 6. Keep Critical Sections Short

**BAD**:
```lua
mu:lock()
-- 100 lines of code
-- Heavy computation
-- Network calls
mu:unlock()
```

**GOOD**:
```lua
-- Do computation outside lock
local result = heavy_computation()

-- Lock only for data access
mu:lock()
shared_data = result
mu:unlock()
```

## When to Use What?

### Use **Mutex** when:
- ✅ Protecting shared mutable data
- ✅ Need exclusive access
- ✅ Both reads and writes need protection

### Use **RWMutex** when:
- ✅ Many readers, few writers
- ✅ Reading is much more common than writing
- ✅ Want concurrent reads for better performance

### Use **Channels** when:
- ✅ Communication between goroutines
- ✅ Passing data ownership
- ✅ Coordinating workflows

### Use **WaitGroup** when:
- ✅ Waiting for multiple goroutines to finish
- ✅ No data sharing needed

## Comparison Table

| Feature | Mutex | RWMutex | Channel |
|---------|-------|---------|---------|
| **Purpose** | Protect data | Protect data | Communicate |
| **Concurrent reads** | ❌ No | ✅ Yes | N/A |
| **Exclusive writes** | ✅ Yes | ✅ Yes | N/A |
| **Passes data** | ❌ No | ❌ No | ✅ Yes |
| **Blocking** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Non-blocking** | ✅ try_lock | ✅ try_lock/try_rlock | ✅ try_send/try_receive |
| **Use case** | Shared mutable data | Read-heavy data | Message passing |

## Common Pitfalls

### 1. Forgetting to Unlock

```lua
mu:lock()
if error_condition then
    return false  -- ⚠️ FORGOT TO UNLOCK!
end
mu:unlock()
```

### 2. Locking Already Locked Mutex

```lua
mu:lock()
some_function()  -- This function also calls mu:lock()
mu:unlock()

-- ⚠️ DEADLOCK if some_function locks the same mutex!
```

### 3. Accessing Shared Data Without Lock

```lua
mu:lock()
shared_data.field1 = value1
mu:unlock()

-- ⚠️ RACE CONDITION! No lock here
local x = shared_data.field2
```

### 4. Holding Lock During Expensive Operations

```lua
mu:lock()
make_http_request()  -- ⚠️ Blocks all other goroutines!
mu:unlock()
```

**Fix**: Do expensive work outside the lock:

```lua
local response = make_http_request()

mu:lock()
shared_data = response
mu:unlock()
```

## Advanced Patterns

### Pattern 1: Mutex with Defer-like Behavior

Use a helper function to ensure unlock:

```lua
local with_lock = function(mu, fn)
    mu:lock()
    local success, result = pcall(fn)
    mu:unlock()

    if not success then
        error(result)
    end
    return result
end

-- Usage
with_lock(mu, function()
    -- Critical section
    shared_data = shared_data + 1
end)
```

### Pattern 2: Conditional Lock

```lua
local conditional_update = function(condition, update_fn)
    if not condition then
        return false
    end

    mu:lock()
    update_fn()
    mu:unlock()
    return true
end
```

### Pattern 3: Read-Modify-Write

```lua
local increment = function()
    rwmu:lock()
    local current = counter
    counter = current + 1
    rwmu:unlock()
    return current
end
```

## Troubleshooting

### Problem: Deadlock

**Symptoms**: Program hangs, no progress

**Causes**:
- Locking same mutex twice
- Circular lock dependencies
- Forgetting to unlock

**Solutions**:
- Use `try_lock()` with backoff
- Always lock in same order
- Keep critical sections short

### Problem: Race Conditions

**Symptoms**: Inconsistent data, weird values

**Causes**:
- Accessing shared data without lock
- Unlocking too early

**Solutions**:
- Always lock before accessing shared data
- Use mutex for all reads AND writes

### Problem: Poor Performance

**Symptoms**: Slow execution with many goroutines

**Causes**:
- Using Mutex instead of RWMutex for read-heavy workloads
- Holding locks too long
- Too much contention

**Solutions**:
- Use RWMutex for concurrent reads
- Keep critical sections minimal
- Reduce contention (shard data)

## Conclusion

Mutexes são essenciais para programação concorrente segura. Use-os para proteger dados compartilhados, mas lembre-se:

- **Keep it simple**: Lock, work, unlock
- **Keep it short**: Minimize critical sections
- **Keep it safe**: Always unlock

Para mais informações sobre mutexes em Go (que inspirou esta implementação):
- [Go sync.Mutex](https://pkg.go.dev/sync#Mutex)
- [Go sync.RWMutex](https://pkg.go.dev/sync#RWMutex)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
