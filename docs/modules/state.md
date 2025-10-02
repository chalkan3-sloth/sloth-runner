# 💾 State Management Module

The **State Management** module provides powerful persistent state capabilities with atomic operations, distributed locks, and TTL (Time To Live) functionality. All data is stored locally using SQLite with WAL mode for maximum performance and reliability.

## 🚀 Key Features

- **SQLite Persistence**: Reliable storage with WAL mode
- **Atomic Operations**: Thread-safe increment, compare-and-swap, append
- **Distributed Locks**: Critical sections with automatic timeout
- **TTL (Time To Live)**: Automatic key expiration
- **Data Types**: String, number, boolean, table, list
- **Pattern Matching**: Wildcard key searches
- **Auto Cleanup**: Background cleanup of expired data
- **Statistics**: Usage and performance metrics

## 📋 Basic Usage

### Setting and Getting Values

```lua
-- Set values
state.set("app_version", "v1.2.3")
state.set("user_count", 1000)
state.set("config", {
    debug = true,
    max_connections = 100
})

-- Get values
local version = state.get("app_version")
local count = state.get("user_count")
local config = state.get("config")

-- Get with default value
local theme = state.get("ui_theme", "dark")

-- Check existence
if state.exists("app_version") then
    log.info("App version is configured")
end

-- Delete key
state.delete("old_key")
```

### TTL (Time To Live)

```lua
-- Set with TTL (60 seconds)
state.set("session_token", "abc123", 60)

-- Set TTL for existing key
state.set_ttl("user_session", 300) -- 5 minutes

-- Check remaining TTL
local ttl = state.get_ttl("session_token")
log.info("Token expires in " .. ttl .. " seconds")
```

### Atomic Operations

```lua
-- Atomic increment
local counter = state.increment("page_views", 1)
local bulk_counter = state.increment("downloads", 50)

-- Atomic decrement  
local remaining = state.decrement("inventory", 5)

-- String append
state.set("log_messages", "Starting application")
local new_length = state.append("log_messages", " -> Connecting to database")

-- Atomic compare-and-swap
local old_version = state.get("config_version")
local success = state.compare_swap("config_version", old_version, old_version + 1)
if success then
    log.info("Configuration updated safely")
end
```

### List Operations

```lua
-- Add items to list
state.list_push("deployment_queue", {
    app = "frontend",
    version = "v2.1.0",
    environment = "staging"
})

-- Check list size
local queue_size = state.list_length("deployment_queue")
log.info("Items in queue: " .. queue_size)

-- Process list (pop removes last item)
while state.list_length("deployment_queue") > 0 do
    local deployment = state.list_pop("deployment_queue")
    log.info("Processing deployment: " .. deployment.app)
    -- Process deployment...
end
```

### Distributed Locks and Critical Sections

```lua
-- Try to acquire lock (no waiting)
local lock_acquired = state.try_lock("deployment_lock", 30) -- 30 seconds TTL
if lock_acquired then
    -- Critical work
    state.unlock("deployment_lock")
end

-- Lock with wait and timeout
local acquired = state.lock("database_migration", 60) -- wait up to 60s
if acquired then
    -- Execute migration
    state.unlock("database_migration")
end

-- Critical section with automatic lock management
state.with_lock("critical_section", function()
    log.info("Executing critical operation...")
    
    -- Update global counter
    local counter = state.increment("global_counter", 1)
    
    -- Update timestamp
    state.set("last_operation", os.time())
    
    log.info("Critical operation completed - counter: " .. counter)
    
    -- Lock is automatically released when function returns
    return "operation_success"
end, 15) -- 15 second timeout
```

## 🔍 API Reference

### Basic Operations
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `state.set(key, value, ttl?)` | key: string, value: any, ttl?: number | success: boolean | Set a value with optional TTL |
| `state.get(key, default?)` | key: string, default?: any | value: any | Get a value or return default |
| `state.delete(key)` | key: string | success: boolean | Remove a key |
| `state.exists(key)` | key: string | exists: boolean | Check if key exists |
| `state.clear(pattern?)` | pattern?: string | success: boolean | Remove keys by pattern |

### TTL Operations
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `state.set_ttl(key, seconds)` | key: string, seconds: number | success: boolean | Set TTL for existing key |
| `state.get_ttl(key)` | key: string | ttl: number | Get remaining TTL (-1 = no TTL, -2 = not exists) |

### Atomic Operations
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `state.increment(key, delta?)` | key: string, delta?: number | new_value: number | Atomically increment value |
| `state.decrement(key, delta?)` | key: string, delta?: number | new_value: number | Atomically decrement value |
| `state.append(key, value)` | key: string, value: string | new_length: number | Atomically append string |
| `state.compare_swap(key, old, new)` | key: string, old: any, new: any | success: boolean | Atomic compare-and-swap |

### List Operations
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `state.list_push(key, item)` | key: string, item: any | length: number | Add item to end of list |
| `state.list_pop(key)` | key: string | item: any \| nil | Remove and return last item |
| `state.list_length(key)` | key: string | length: number | Get list length |

### Distributed Locks
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `state.try_lock(name, ttl)` | name: string, ttl: number | success: boolean | Try to acquire lock without waiting |
| `state.lock(name, timeout?)` | name: string, timeout?: number | success: boolean | Acquire lock with timeout |
| `state.unlock(name)` | name: string | success: boolean | Release lock |
| `state.with_lock(name, fn, timeout?)` | name: string, fn: function, timeout?: number | result: any | Execute function with automatic lock |

### Utilities
| Function | Parameters | Return | Description |
|----------|------------|---------|-------------|
| `state.keys(pattern?)` | pattern?: string | keys: table | List keys by pattern |
| `state.stats()` | - | stats: table | Get system statistics |

## 💡 Practical Use Cases

### 1. Deployment Version Control

```lua
Modern DSLs = {
    deployment_pipeline = {
        tasks = {
            prepare_deploy = {
                command = function()
                    -- Check last deployed version
                    local last_version = state.get("last_deployed_version", "v0.0.0")
                    local new_version = "v1.2.3"
                    
                    -- Check if already deployed
                    if last_version == new_version then
                        log.warn("Version " .. new_version .. " already deployed")
                        return false, "Version already deployed"
                    end
                    
                    -- Register deployment start
                    state.set("deploy_status", "in_progress")
                    state.set("deploy_start_time", os.time())
                    state.increment("total_deploys", 1)
                    
                    return true, "Deploy preparation completed"
                end
            },
            
            execute_deploy = {
                depends_on = "prepare_deploy",
                command = function()
                    -- Critical section for deployment
                    return state.with_lock("deployment_lock", function()
                        log.info("Executing deployment with lock...")
                        
                        -- Simulate deployment
                        exec.run("sleep 5")
                        
                        -- Update state
                        state.set("last_deployed_version", "v1.2.3")
                        state.set("deploy_status", "completed")
                        state.set("deploy_end_time", os.time())
                        
                        -- Record history
                        state.list_push("deploy_history", {
                            version = "v1.2.3",
                            timestamp = os.time(),
                            duration = state.get("deploy_end_time") - state.get("deploy_start_time")
                        })
                        
                        return true, "Deploy completed successfully"
                    end, 300) -- 5 minutes timeout
                end
            }
        }
    }
}
```

### 2. Intelligent Caching with TTL

```lua
-- Helper function for caching
function get_cached_data(cache_key, fetch_function, ttl)
    local cached = state.get(cache_key)
    if cached then
        log.info("Cache hit: " .. cache_key)
        return cached
    end
    
    log.info("Cache miss: " .. cache_key .. " - fetching...")
    local data = fetch_function()
    state.set(cache_key, data, ttl or 300) -- 5 minutes default
    return data
end

-- Usage in tasks
Modern DSLs = {
    data_processing = {
        tasks = {
            fetch_user_data = {
                command = function()
                    local user_data = get_cached_data("user:123:profile", function()
                        -- Simulate expensive fetch
                        return {
                            name = "Alice",
                            email = "alice@example.com",
                            preferences = {"dark_mode", "notifications"}
                        }
                    end, 600) -- Cache for 10 minutes
                    
                    log.info("User data: " .. data.to_json(user_data))
                    return true, "User data retrieved"
                end
            }
        }
    }
}
```

### 3. Rate Limiting

```lua
function check_rate_limit(identifier, max_requests, window_seconds)
    local key = "rate_limit:" .. identifier
    local current_count = state.get(key, 0)
    
    if current_count >= max_requests then
        return false, "Rate limit exceeded"
    end
    
    -- Increment counter
    if current_count == 0 then
        -- First request in window
        state.set(key, 1, window_seconds)
    else
        -- Increment existing counter
        state.increment(key, 1)
    end
    
    return true, "Request allowed"
end

-- Usage in tasks
Modern DSLs = {
    api_tasks = {
        tasks = {
            make_api_call = {
                command = function()
                    local allowed, msg = check_rate_limit("api_calls", 100, 3600) -- 100 calls/hour
                    
                    if not allowed then
                        log.error(msg)
                        return false, msg
                    end
                    
                    -- Make API call
                    log.info("Making API call...")
                    return true, "API call completed"
                end
            }
        }
    }
}
```

## ⚙️ Configuration and Storage

### Database Location

By default, the SQLite database is created at:
- **Linux/macOS**: `~/.sloth-runner/state.db`
- **Windows**: `%USERPROFILE%\.sloth-runner\state.db`

### Technical Characteristics

- **Engine**: SQLite 3 with WAL mode
- **Concurrent Access**: Support for multiple simultaneous connections
- **Auto-cleanup**: Automatic cleanup of expired data every 5 minutes
- **Lock Timeout**: Expired locks are cleaned automatically
- **Serialization**: JSON for complex objects, native format for simple types

### Limitations

- **Local Scope**: State is persisted only on local machine
- **Concurrency**: Locks are effective only within local process
- **Size**: Suitable for small to medium datasets (< 1GB)

## 🔄 Best Practices

1. **Use TTL for temporary data** to prevent storage bloat
2. **Use locks for critical sections** to avoid race conditions  
3. **Use patterns for bulk operations** to manage related keys
4. **Monitor storage size** using `state.stats()` 
5. **Use atomic operations** instead of read-modify-write patterns
6. **Clean up expired keys** regularly with `state.clear(pattern)`

The **State Management** module transforms sloth-runner into a stateful, reliable platform for complex task orchestration! 🚀