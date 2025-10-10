# 📦 Complete Modules Reference

Sloth Runner provides **14 built-in modules** for comprehensive task automation. All modules are **globally available** - no `require()` needed!

## 🚀 Quick Start

```lua
-- Modules are global! Just use them directly
http.get({url = "https://api.example.com"})
system.exec("ls", {"-la"})
crypto.sha256("my data")
```

## 📚 Module Categories

### ⚡ Core Execution (5 modules)
- **[system](#system-module)** - Process execution, file operations, system info
- **[http](#http-module)** - HTTP client with retries and timeouts
- **[goroutine](#goroutine-module)** - Parallel execution and worker pools
- **[event](#event-module)** - Event dispatching for hooks
- **[help](#help-module)** - Interactive help system

### 🔐 Security & Data (4 modules)
- **[crypto](#crypto-module)** - Hashing, encryption, random generation
- **[validate](#validate-module)** - Data validation and sanitization
- **[database](#database-module)** - SQLite, MySQL, PostgreSQL support
- **[monitor](#monitor-module)** - Metrics collection (counters, gauges, histograms)

### 📢 Integration (1 module)
- **[notify](#notify-module)** - Slack, Discord, Email, Teams, Telegram

### 🎯 Automation (3 modules)
- **[facts](#facts-module)** - Agent system information (requires master)
- **[git](#git-module)** - Git repository operations (idempotent)
- **[sloth](#sloth-module)** - Self-management automation

### 🏗️ Infrastructure (1 module)
- **[incus](#incus-module)** - Container/VM management (fluent API)

---

## System Module

**Category**: Core Execution
**Description**: System operations including process management, file operations, and system information

### Functions

#### Process Management
```lua
-- Execute command synchronously
local result = system.exec("ls", {"-la", "/tmp"})
if result.success then
    print(result.output)
    print("Exit code:", result.exit_code)
end

-- Execute asynchronously (returns PID)
local pid = system.exec_async("long-running-command")
print("Started process:", pid)

-- Kill process
system.kill(pid, "TERM")  -- Send SIGTERM
system.kill(pid, "KILL")  -- Send SIGKILL
```

#### File Operations
```lua
-- Check existence
if system.exists("/tmp/file.txt") then
    print("File exists")
end

-- Directory operations
system.mkdir("/tmp/mydir", true)  -- recursive
system.rmdir("/tmp/mydir", true)  -- recursive

-- Copy and move
system.copy("/tmp/src.txt", "/tmp/dst.txt")
system.move("/tmp/old.txt", "/tmp/new.txt")

-- Permissions
system.chmod("/tmp/script.sh", "755")
system.chmod("/tmp/file.txt", "644")

-- File info
local info = system.stat("/tmp/file.txt")
print("Size:", info.size)
print("Mode:", info.mode)
print("Modified:", info.mtime)
```

#### Environment Variables
```lua
-- Get environment variable
local home = system.env("HOME")
print("Home dir:", home)

-- Set environment variable
system.setenv("MY_VAR", "value")

-- Unset environment variable
system.unsetenv("MY_VAR")
```

#### System Information
```lua
-- Basic info
print("Hostname:", system.hostname())
print("Platform:", system.platform())  -- "linux", "darwin", "windows"
print("Architecture:", system.arch())  -- "amd64", "arm64"

-- Resource info
print("CPUs:", system.cpu_count())
local mem = system.memory()
print("Total memory:", mem.total)
print("Available:", mem.available)
print("Used percent:", mem.used_percent)

-- Uptime
local uptime_secs = system.uptime()
print("System uptime:", uptime_secs, "seconds")

-- Process list
local procs = system.processes()
for _, proc in ipairs(procs) do
    print(proc.pid, proc.name, proc.cpu_percent)
end
```

#### Directory Navigation
```lua
-- Get current directory
local cwd = system.pwd()
print("Current dir:", cwd)

-- Change directory
system.cd("/tmp")

-- Find executable
local path = system.which("git")
print("Git location:", path)

-- Special directories
print("Temp dir:", system.temp_dir())
print("Home dir:", system.home_dir())
```

### Complete Example

```lua
task({
    name = "system-operations",
    run = function()
        -- Check system requirements
        local mem = system.memory()
        if mem.used_percent > 90 then
            error("Insufficient memory!")
        end

        -- Create working directory
        local work_dir = "/tmp/my_deployment"
        if not system.exists(work_dir) then
            system.mkdir(work_dir, true)
        end

        -- Execute deployment script
        system.cd(work_dir)
        local result = system.exec("./deploy.sh", {"--env", "production"})

        if result.success then
            -- Set permissions on output
            system.chmod("output.log", "644")
            return true, "Deployment successful"
        else
            return false, "Deployment failed: " .. result.stderr
        end
    end
})
```

---

## HTTP Module

**Category**: Core Execution
**Description**: HTTP client with advanced features (retries, timeouts, response validation)

### Functions

```lua
-- GET request
local response = http.get({
    url = "https://api.example.com/users",
    headers = {
        ["Authorization"] = "Bearer token123",
        ["Accept"] = "application/json"
    },
    timeout = 5,           -- seconds
    max_retries = 2,       -- retry on failure
    follow_redirects = true
})

if response.status_code == 200 then
    print("Success!")
    print(response.body)

    -- Parse JSON response
    local data = response.json
    for _, user in ipairs(data.users) do
        print(user.name)
    end
end

-- POST with JSON
http.post({
    url = "https://api.example.com/users",
    json = {
        name = "John Doe",
        email = "john@example.com"
    },
    headers = {
        ["Authorization"] = "Bearer token123"
    }
})

-- POST with form data
http.post({
    url = "https://api.example.com/upload",
    body = "field1=value1&field2=value2",
    headers = {
        ["Content-Type"] = "application/x-www-form-urlencoded"
    }
})

-- PUT request
http.put({
    url = "https://api.example.com/users/123",
    json = {
        name = "Jane Doe"
    }
})

-- DELETE request
http.delete({
    url = "https://api.example.com/users/123"
})

-- PATCH request
http.patch({
    url = "https://api.example.com/users/123",
    json = {
        status = "active"
    }
})

-- Generic request with custom method
http.request({
    method = "OPTIONS",
    url = "https://api.example.com"
})

-- Download file
http.download({
    url = "https://example.com/file.tar.gz",
    destination = "/tmp/file.tar.gz",
    timeout = 30
})

-- URL encoding
local encoded = http.url_encode("hello world & stuff")
print(encoded)  -- "hello+world+%26+stuff"

local decoded = http.url_decode(encoded)
print(decoded)  -- "hello world & stuff"

-- Build URL with parameters
local url = http.build_url({
    base = "https://api.example.com/search",
    path = "/users",
    params = {
        q = "john",
        limit = "10",
        offset = "0"
    }
})
-- Result: https://api.example.com/search/users?q=john&limit=10&offset=0

-- JSON helpers
local json_str = http.to_json({name = "test", value = 42})
local data = http.parse_json(json_str)

-- Status helpers
if http.is_success(response.status_code) then
    print("Success (2xx)")
end

if http.is_error(response.status_code) then
    print("Error (4xx or 5xx)")
end
```

### Complete Example

```lua
task({
    name = "api-integration",
    run = function()
        -- Check API health
        local health = http.get({
            url = "https://api.example.com/health",
            timeout = 5,
            max_retries = 3
        })

        if not http.is_success(health.status_code) then
            return false, "API is not healthy"
        end

        -- Create resource
        local create_resp = http.post({
            url = "https://api.example.com/resources",
            json = {
                name = "my-resource",
                config = {
                    enabled = true,
                    replicas = 3
                }
            },
            headers = {
                ["Authorization"] = "Bearer " .. os.getenv("API_TOKEN")
            }
        })

        if create_resp.status_code == 201 then
            local resource_id = create_resp.json.id
            print("Created resource:", resource_id)
            return true, "Resource created: " .. resource_id
        else
            return false, "Failed to create resource: " .. create_resp.body
        end
    end
})
```

---

## Goroutine Module

**Category**: Core Execution
**Description**: Concurrent execution with goroutines, worker pools, and async operations

### Functions

```lua
-- Spawn single goroutine
goroutine.spawn(function()
    print("Running in goroutine")
    system.exec("sleep", {"1"})
    print("Done")
end)

-- Spawn multiple goroutines
goroutine.spawn_many(10, function(id)
    print("Goroutine", id)
end)

-- Worker pool pattern
goroutine.pool_create("workers", {workers = 10, queue_size = 100})

for i = 1, 100 do
    goroutine.pool_submit("workers", function(task_id)
        print("Processing task", task_id)
        system.exec("sleep", {"1"})
        return "result-" .. task_id
    end, i)
end

-- Wait for all tasks to complete
goroutine.pool_wait("workers")

-- Get pool statistics
local stats = goroutine.pool_stats("workers")
print("Completed:", stats.completed)
print("Pending:", stats.pending)
print("Failed:", stats.failed)

-- Close pool
goroutine.pool_close("workers")

-- Async/await pattern
local handle = goroutine.async(function()
    goroutine.sleep(1000)  -- milliseconds
    return "async result"
end)

-- Do other work...

-- Wait for result
local success, result = goroutine.await(handle)
if success then
    print("Got result:", result)
end

-- Wait for multiple async operations
local h1 = goroutine.async(function() return "result1" end)
local h2 = goroutine.async(function() return "result2" end)
local h3 = goroutine.async(function() return "result3" end)

local results = goroutine.await_all({h1, h2, h3})
for i, res in ipairs(results) do
    print("Result", i, ":", res.success, res.value)
end

-- Sleep (non-blocking)
goroutine.sleep(1000)  -- milliseconds

-- Execute with timeout
local success, result = goroutine.timeout(5000, function()
    -- Long running operation
    system.exec("sleep", {"10"})
    return "done"
end)

if not success then
    print("Operation timed out!")
end
```

### Complete Example

```lua
task({
    name = "parallel-deployment",
    run = function()
        local servers = {"web-01", "web-02", "web-03", "web-04"}

        -- Create worker pool for deployments
        goroutine.pool_create("deploy", {workers = 4})

        -- Submit deployment tasks
        for _, server in ipairs(servers) do
            goroutine.pool_submit("deploy", function(srv)
                print("Deploying to", srv)

                -- Copy files
                system.exec("scp", {"app.tar.gz", srv .. ":/tmp/"})

                -- Extract and restart
                system.exec("ssh", {srv, "tar -xzf /tmp/app.tar.gz -C /opt/app"})
                system.exec("ssh", {srv, "systemctl restart app"})

                -- Verify
                goroutine.sleep(2000)
                local result = system.exec("ssh", {srv, "systemctl is-active app"})

                if result.output:match("active") then
                    print(srv, "deployment successful")
                    return true
                else
                    error(srv .. " deployment failed")
                end
            end, server)
        end

        -- Wait for all deployments
        goroutine.pool_wait("deploy")

        local stats = goroutine.pool_stats("deploy")
        goroutine.pool_close("deploy")

        if stats.failed > 0 then
            return false, stats.failed .. " deployments failed"
        end

        return true, "All " .. #servers .. " servers deployed successfully"
    end
})
```

---

## Event Module

**Category**: Core Execution
**Description**: Event dispatching system for triggering hooks from workflows

### Functions

```lua
-- Dispatch generic event
event.dispatch("deployment.complete", {
    environment = "production",
    version = "v1.2.3",
    timestamp = os.time(),
    deployed_by = "deploy-bot",
    commit_sha = "abc123"
})

-- Dispatch custom event with simple message
event.dispatch_custom("backup_complete", "Database backup finished successfully")

-- Dispatch file event (created/modified/deleted)
event.dispatch_file("created", "/var/log/app.log", "log_watcher")
event.dispatch_file("modified", "/etc/nginx/nginx.conf", "config_watcher")
event.dispatch_file("deleted", "/tmp/old_file.txt", "cleanup_monitor")
```

### Complete Example

```lua
task({
    name = "monitored-deployment",
    run = function()
        -- Dispatch start event
        event.dispatch("deployment.started", {
            environment = values.env or "staging",
            initiated_by = os.getenv("USER"),
            timestamp = os.time()
        })

        -- Perform deployment
        local success = true
        local error_msg = nil

        local result = system.exec("./deploy.sh")
        if not result.success then
            success = false
            error_msg = result.stderr
        end

        -- Dispatch completion event
        if success then
            event.dispatch("deployment.complete", {
                environment = values.env or "staging",
                status = "success",
                duration = os.time() - start_time
            })
        else
            event.dispatch("deployment.failed", {
                environment = values.env or "staging",
                error = error_msg,
                duration = os.time() - start_time
            })
        end

        return success, success and "Deployment complete" or error_msg
    end
})
```

---

## Help Module

**Category**: Core Execution
**Description**: Interactive help system for exploring available modules and functions

### Functions

```lua
-- Show general help
help()

-- List all modules
help.modules()

-- Search for functions
help.search("http")       -- Find all http-related functions
help.search("execute")    -- Find all execute-related functions

-- Show examples
help.examples()           -- Show general examples
help.examples("http")     -- Show http module examples
```

### Usage

```lua
task({
    name = "explore-modules",
    run = function()
        -- List available modules
        local modules = help.modules()
        print("Available modules:", #modules)

        -- Search for file-related functions
        local results = help.search("file")
        for _, result in ipairs(results) do
            print(result.module .. "." .. result.function)
        end

        return true
    end
})
```

---

## Crypto Module

**Category**: Security & Data
**Description**: Cryptographic operations including hashing, encryption, and random generation

### Functions

```lua
-- Hashing
local md5 = crypto.md5("my data")
local sha1 = crypto.sha1("my data")
local sha256 = crypto.sha256("my data")
local sha512 = crypto.sha512("my data")

-- Password hashing (bcrypt)
local hashed = crypto.bcrypt_hash("password123", 10)  -- cost factor 10
print("Hashed:", hashed)

local is_valid = crypto.bcrypt_check("password123", hashed)
print("Password valid:", is_valid)

-- Base64 encoding
local encoded = crypto.base64_encode("hello world")
local decoded = crypto.base64_decode(encoded)

-- Hex encoding
local hex = crypto.hex_encode("binary data")
local bin = crypto.hex_decode(hex)

-- AES encryption (AES-256-GCM)
local key = "my-encryption-key-32-bytes-long"  -- Must be 32 bytes for AES-256
local plaintext = "secret message"

local encrypted = crypto.aes_encrypt(plaintext, key)
print("Encrypted:", encrypted)

local decrypted = crypto.aes_decrypt(encrypted, key)
print("Decrypted:", decrypted)

-- Random generation
local random_hex = crypto.random_bytes(16)  -- 16 bytes, hex-encoded
print("Random:", random_hex)

local random_str = crypto.random_string(32, "alphanumeric")  -- charset options: alphanumeric, alpha, numeric, hex
print("Random string:", random_str)

-- Generate secure password
local password = crypto.generate_password(16)  -- length
print("Generated password:", password)
```

### Complete Example

```lua
task({
    name = "secure-credentials",
    run = function()
        -- Generate API key
        local api_key = crypto.random_string(32, "hex")
        print("Generated API key:", api_key)

        -- Hash API key for storage
        local api_key_hash = crypto.sha256(api_key)

        -- Store hashed key (example)
        system.exec("echo", {api_key_hash, ">", "/var/secrets/api_key.hash"})

        -- Generate user password
        local user_password = crypto.generate_password(16)
        local password_hash = crypto.bcrypt_hash(user_password, 12)

        print("User password:", user_password)
        print("Password hash:", password_hash)

        -- Encrypt sensitive config
        local config = http.to_json({
            database_url = "postgresql://user:pass@localhost/db",
            api_secret = "secret123"
        })

        local encryption_key = os.getenv("ENCRYPTION_KEY") or error("ENCRYPTION_KEY not set")
        local encrypted_config = crypto.aes_encrypt(config, encryption_key)

        -- Save encrypted config
        local f = io.open("/etc/app/config.enc", "w")
        f:write(encrypted_config)
        f:close()

        return true, "Credentials secured successfully"
    end
})
```

---

## Validate Module

**Category**: Security & Data
**Description**: Data validation and sanitization utilities

### Functions

```lua
-- Email validation
local result = validate.email("test@example.com")
if result.valid then
    print("Valid email:", result.email)
    print("Domain:", result.domain)
end

-- URL validation
local url_result = validate.url("https://example.com:8080/path?query=value")
if url_result.valid then
    print("Scheme:", url_result.scheme)    -- "https"
    print("Host:", url_result.host)        -- "example.com"
    print("Port:", url_result.port)        -- "8080"
    print("Path:", url_result.path)        -- "/path"
    print("Query:", url_result.query)      -- "query=value"
end

-- IP address validation
local ip = validate.ip("192.168.1.1")
if ip.valid then
    print("IP version:", ip.version)  -- "v4" or "v6"
end

-- Regex validation with capture groups
local regex_result = validate.regex("test123", "([a-z]+)(\\d+)")
if regex_result.valid then
    print("Matched:", regex_result.matched)
    print("Groups:", regex_result.groups[1], regex_result.groups[2])  -- "test", "123"
end

-- Length validation
local len = validate.length("hello", {min = 3, max = 10})
print("Valid length:", len.valid)

-- Can also check exact length
local exact = validate.length("hello", {exact = 5})

-- Range validation (numeric)
local range = validate.range(42, {min = 0, max = 100})
print("In range:", range.valid)

-- Required fields validation
local data = {
    name = "John",
    email = "john@example.com"
}
local required = validate.required(data, {"name", "email", "age"})
if not required.valid then
    print("Missing fields:", table.concat(required.missing, ", "))
end

-- Schema validation
local schema = {
    name = {required = "true", type = "string", min_length = 3},
    age = {required = "true", type = "number", min = 18, max = 120},
    email = {required = "true", type = "string"}
}

local user_data = {
    name = "Alice",
    age = 25,
    email = "alice@example.com"
}

local schema_result = validate.schema(user_data, schema)
if not schema_result.valid then
    for _, error in ipairs(schema_result.errors) do
        print("Error:", error.field, error.message)
    end
end

-- Sanitization
local html_clean = validate.sanitize("<script>alert('xss')</script>", "html")
print("Clean HTML:", html_clean)

local sql_clean = validate.sanitize("'; DROP TABLE users--", "sql")
print("Clean SQL:", sql_clean)

-- Other sanitization options
local trimmed = validate.sanitize("  hello  ", "trim")
local lower = validate.sanitize("HELLO", "lower")
local upper = validate.sanitize("hello", "upper")
```

### Complete Example

```lua
task({
    name = "validate-user-input",
    run = function()
        -- Get user input (from values or environment)
        local user_email = values.email
        local user_age = tonumber(values.age)
        local user_website = values.website

        -- Validate email
        local email_check = validate.email(user_email)
        if not email_check.valid then
            return false, "Invalid email address"
        end

        -- Validate age range
        local age_check = validate.range(user_age, {min = 18, max = 120})
        if not age_check.valid then
            return false, "Age must be between 18 and 120"
        end

        -- Validate website URL
        if user_website then
            local url_check = validate.url(user_website)
            if not url_check.valid then
                return false, "Invalid website URL"
            end
            if url_check.scheme ~= "https" then
                return false, "Website must use HTTPS"
            end
        end

        -- All validations passed
        print("✓ All validations passed")
        print("Email:", email_check.email)
        print("Age:", user_age)
        print("Website:", user_website or "N/A")

        return true, "User input validated successfully"
    end
})
```

---

## Database Module

**Category**: Security & Data
**Description**: Database operations with support for SQLite, MySQL, PostgreSQL

### Functions

```lua
-- Connect to database
db.connect("mydb", "sqlite3", "/path/to/database.sqlite")

-- Or connect to MySQL
db.connect("mysql_db", "mysql", "user:password@tcp(localhost:3306)/dbname")

-- Or connect to PostgreSQL
db.connect("pg_db", "postgres", "host=localhost port=5432 user=user password=pass dbname=mydb sslmode=disable")

-- Query (SELECT)
local rows = db.query("mydb", "SELECT * FROM users WHERE active = ?", {1})
for _, row in ipairs(rows) do
    print(row.id, row.name, row.email)
end

-- Execute (INSERT/UPDATE/DELETE)
local result = db.exec("mydb", "INSERT INTO users (name, email) VALUES (?, ?)",
    {"John Doe", "john@example.com"})
print("Last insert ID:", result.last_insert_id)
print("Rows affected:", result.rows_affected)

-- Update
db.exec("mydb", "UPDATE users SET active = ? WHERE id = ?", {1, 123})

-- Delete
db.exec("mydb", "DELETE FROM users WHERE id = ?", {123})

-- Transaction (execute multiple queries atomically)
local tx_result = db.transaction("mydb", {
    {
        query = "INSERT INTO users (name, email) VALUES (?, ?)",
        params = {"Alice", "alice@example.com"}
    },
    {
        query = "UPDATE counters SET count = count + 1 WHERE key = ?",
        params = {"user_count"}
    },
    {
        query = "INSERT INTO audit_log (action, timestamp) VALUES (?, ?)",
        params = {"user_created", os.time()}
    }
})

if not tx_result.success then
    print("Transaction failed:", tx_result.error)
else
    print("Transaction successful")
end

-- Prepared statements
local stmt_id = db.prepare("mydb", "SELECT * FROM users WHERE email = ?")
local results = db.query("mydb", stmt_id, {"john@example.com"})

-- Test connection
local ok = db.ping("mydb")
print("Connection OK:", ok)

-- Disconnect
db.disconnect("mydb")

-- Close all connections
db.close_all()
```

### Complete Example

```lua
task({
    name = "database-migration",
    run = function()
        -- Connect to SQLite database
        db.connect("app_db", "sqlite3", "/var/lib/app/app.db")

        -- Check connection
        if not db.ping("app_db") then
            return false, "Cannot connect to database"
        end

        -- Run migration in transaction
        local migration = db.transaction("app_db", {
            {
                query = [[
                    CREATE TABLE IF NOT EXISTS users (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        name TEXT NOT NULL,
                        email TEXT UNIQUE NOT NULL,
                        created_at INTEGER NOT NULL,
                        active INTEGER DEFAULT 1
                    )
                ]]
            },
            {
                query = [[
                    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)
                ]]
            },
            {
                query = [[
                    CREATE TABLE IF NOT EXISTS audit_log (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        action TEXT NOT NULL,
                        user_id INTEGER,
                        timestamp INTEGER NOT NULL
                    )
                ]]
            }
        })

        if not migration.success then
            db.disconnect("app_db")
            return false, "Migration failed: " .. migration.error
        end

        -- Insert initial data
        db.exec("app_db", "INSERT OR IGNORE INTO users (name, email, created_at) VALUES (?, ?, ?)",
            {"Admin User", "admin@example.com", os.time()})

        -- Verify
        local count = db.query("app_db", "SELECT COUNT(*) as count FROM users", {})
        print("Total users:", count[1].count)

        db.disconnect("app_db")

        return true, "Database migration completed successfully"
    end
})
```

---

## Monitor Module

**Category**: Security & Data
**Description**: Metrics collection and reporting with support for counters, gauges, histograms

### Functions

```lua
-- Counter operations
monitor.counter_inc("http_requests", {endpoint = "/api/users", method = "GET"})
monitor.counter_add("bytes_sent", 1024, {service = "api"})

-- Gauge operations
monitor.gauge_set("active_connections", 42)
monitor.gauge_inc("queue_size", 1)
monitor.gauge_dec("queue_size", 1)

-- Timer operations
local timer_key = monitor.timer_start("request_duration", {endpoint = "/api/users"})
-- ... do work ...
local duration_ms = monitor.timer_end(timer_key)
print("Request took:", duration_ms, "ms")

-- Histogram operations
monitor.histogram_observe("response_time", 0.123, {service = "api", endpoint = "/users"})
monitor.histogram_observe("request_size_bytes", 4096, {service = "api"})

-- Get metric value
local metric = monitor.get_metric("http_requests")
print("Value:", metric.value)
print("Type:", metric.type)  -- counter, gauge, histogram
print("Labels:", metric.labels)

-- List all metrics
local metrics = monitor.list_metrics()
for _, metric in ipairs(metrics) do
    print(metric.name, metric.value, metric.type)
end

-- Reset metric
monitor.reset_metric("http_requests")

-- Clear all metrics
monitor.clear_all()

-- Export metrics
local prom = monitor.export_prometheus()
print(prom)  -- Prometheus format

local json = monitor.export_json()
print(json)  -- JSON format

-- Collect system metrics automatically
monitor.system_metrics()  -- Collects CPU, memory, disk metrics

-- Get detailed memory stats
local mem = monitor.memory_stats()
print("Heap alloc:", mem.heap_alloc)
print("Heap sys:", mem.heap_sys)
print("GC runs:", mem.num_gc)
```

### Complete Example

```lua
task({
    name = "monitored-api-call",
    run = function()
        -- Increment request counter
        monitor.counter_inc("api_requests", {
            endpoint = "/users",
            method = "GET"
        })

        -- Start timer
        local timer = monitor.timer_start("api_request_duration", {
            endpoint = "/users"
        })

        -- Set gauge for active requests
        monitor.gauge_inc("api_active_requests", 1)

        -- Make API call
        local response = http.get({
            url = "https://api.example.com/users",
            timeout = 5
        })

        -- End timer
        local duration = monitor.timer_end(timer)

        -- Decrement active requests
        monitor.gauge_dec("api_active_requests", 1)

        -- Record response size
        monitor.histogram_observe("api_response_size", #response.body, {
            endpoint = "/users"
        })

        -- Record status code
        monitor.counter_inc("api_responses", {
            endpoint = "/users",
            status_code = tostring(response.status_code)
        })

        if response.status_code == 200 then
            monitor.counter_inc("api_success", {endpoint = "/users"})
            return true, "API call successful"
        else
            monitor.counter_inc("api_errors", {endpoint = "/users"})
            return false, "API call failed"
        end
    end
})
```

---

## Notify Module

**Category**: Integration
**Description**: Multi-channel notification services (Slack, Discord, Email, Teams, Telegram)

### Functions

```lua
-- Slack notification
notify.slack("https://hooks.slack.com/services/YOUR/WEBHOOK/URL", {
    text = "Deployment completed successfully!",
    username = "Deploy Bot",
    icon_emoji = ":rocket:",
    channel = "#deployments",
    attachments = {{
        color = "good",  -- or "warning", "danger", or hex color "#FF0000"
        title = "Production Deployment",
        text = "Version 2.0.0 deployed to production",
        fields = {
            {title = "Environment", value = "production", short = true},
            {title = "Version", value = "v2.0.0", short = true},
            {title = "Duration", value = "5m 23s", short = true}
        },
        footer = "Sloth Runner",
        ts = os.time()
    }}
})

-- Discord notification
notify.discord("https://discord.com/api/webhooks/YOUR/WEBHOOK", {
    content = "Deployment completed!",
    username = "Deploy Bot",
    avatar_url = "https://example.com/avatar.png",
    embeds = {{
        title = "Production Deployment",
        description = "Version 2.0.0 has been deployed",
        color = 65280,  -- Decimal color (green = 65280)
        fields = {
            {name = "Environment", value = "production", inline = true},
            {name = "Version", value = "v2.0.0", inline = true},
            {name = "Status", value = "✅ Success", inline = false}
        },
        timestamp = os.date("!%Y-%m-%dT%H:%M:%SZ")
    }}
})

-- Email notification
notify.email({
    smtp_host = "smtp.gmail.com",
    smtp_port = "587",
    username = "your-email@gmail.com",
    password = os.getenv("EMAIL_PASSWORD"),
    from = "deployments@example.com",
    to = "team@example.com",
    cc = "manager@example.com",
    subject = "Deployment Notification",
    body = "Deployment to production completed successfully.",
    html = true  -- Send as HTML
})

-- Webhook notification (generic)
notify.webhook("https://your-webhook-endpoint.com/notify", {
    method = "POST",
    headers = {
        ["Content-Type"] = "application/json",
        ["Authorization"] = "Bearer token123"
    },
    data = {
        event = "deployment.complete",
        environment = "production",
        version = "v2.0.0",
        timestamp = os.time()
    }
})

-- Microsoft Teams notification
notify.teams("https://outlook.office.com/webhook/YOUR/WEBHOOK", {
    title = "Deployment Alert",
    text = "Production deployment completed",
    theme_color = "00FF00",  -- Hex color
    sections = {{
        activityTitle = "Deployment v2.0.0",
        activitySubtitle = "Production Environment",
        facts = {
            {name = "Status", value = "Success"},
            {name = "Duration", value = "5m 23s"},
            {name = "Deployed By", value = "deploy-bot"}
        }
    }}
})

-- Telegram notification
notify.telegram({
    bot_token = os.getenv("TELEGRAM_BOT_TOKEN"),
    chat_id = "-1001234567890",  -- Group chat ID
    text = "🚀 *Deployment Complete*\n\nVersion: v2.0.0\nEnvironment: production\nStatus: ✅ Success",
    parse_mode = "Markdown"  -- or "HTML"
})
```

### Complete Example

```lua
task({
    name = "deploy-with-notifications",
    run = function()
        local environment = values.env or "staging"
        local version = values.version or "latest"

        -- Send start notification
        notify.slack(os.getenv("SLACK_WEBHOOK"), {
            text = "Starting deployment...",
            attachments = {{
                color = "warning",
                title = "Deployment Started",
                fields = {
                    {title = "Environment", value = environment, short = true},
                    {title = "Version", value = version, short = true}
                }
            }}
        })

        -- Perform deployment
        local start_time = os.time()
        local result = system.exec("./deploy.sh", {environment, version})
        local duration = os.time() - start_time

        -- Send completion notification
        if result.success then
            -- Success - notify multiple channels
            notify.slack(os.getenv("SLACK_WEBHOOK"), {
                attachments = {{
                    color = "good",
                    title = "✅ Deployment Successful",
                    fields = {
                        {title = "Environment", value = environment, short = true},
                        {title = "Version", value = version, short = true},
                        {title = "Duration", value = duration .. "s", short = true}
                    }
                }}
            })

            notify.discord(os.getenv("DISCORD_WEBHOOK"), {
                embeds = {{
                    title = "Deployment Complete",
                    description = "Successfully deployed " .. version .. " to " .. environment,
                    color = 65280,
                    fields = {
                        {name = "Duration", value = duration .. "s", inline = true}
                    }
                }}
            })

            return true, "Deployment successful"
        else
            -- Failure - send alert
            notify.slack(os.getenv("SLACK_WEBHOOK"), {
                text = "<!channel> Deployment failed!",
                attachments = {{
                    color = "danger",
                    title = "❌ Deployment Failed",
                    text = result.stderr,
                    fields = {
                        {title = "Environment", value = environment, short = true},
                        {title = "Version", value = version, short = true}
                    }
                }}
            })

            -- Also send email to on-call
            notify.email({
                smtp_host = "smtp.gmail.com",
                smtp_port = "587",
                username = os.getenv("EMAIL_USER"),
                password = os.getenv("EMAIL_PASSWORD"),
                from = "alerts@example.com",
                to = "oncall@example.com",
                subject = "[ALERT] Deployment Failed: " .. environment,
                body = "Deployment to " .. environment .. " failed.\n\nError: " .. result.stderr
            })

            return false, "Deployment failed"
        end
    end
})
```

---

## Facts Module

**Category**: Automation
**Description**: Agent system information (requires master connection)
**Full Documentation**: [facts.md](./facts.md)

### Key Functions

```lua
-- Get all facts from agent
local info, err = facts.get_all({agent = "web-01"})
print("Platform:", info.platform.os)
print("Memory:", info.memory.total)

-- Get specific information
local platform = facts.get_platform({agent = "web-01"})
local memory = facts.get_memory({agent = "web-01"})
local disk = facts.get_disk({agent = "web-01"})
local network = facts.get_network({agent = "web-01"})

-- Check package installation
local pkg = facts.get_package({agent = "web-01", name = "nginx"})
if pkg.installed then
    print("nginx version:", pkg.version)
end

-- Check service status
local svc = facts.get_service({agent = "web-01", name = "nginx"})
print("nginx status:", svc.status)
```

### Quick Example

```lua
task({
    name = "validate-requirements",
    run = function()
        local mem = facts.get_memory({agent = "app-server"})
        if mem.used_percent > 90 then
            return false, "Insufficient memory"
        end

        local pkg = facts.get_package({agent = "app-server", name = "docker"})
        if not pkg.installed then
            return false, "Docker not installed"
        end

        return true, "Requirements met"
    end
})
```

---

## Git Module

**Category**: Automation
**Description**: Git repository operations with idempotency support

### Functions

```lua
-- Clone repository (idempotent - won't re-clone if exists)
local repo, err = git.clone({
    url = "https://github.com/user/repo.git",
    local_path = "/tmp/repo",
    branch = "main",
    depth = 1,        -- shallow clone
    clean = false,    -- if true, removes and re-clones if exists
    single_branch = true
})

if err then
    print("Clone failed:", err)
else
    print("Repository at:", repo.path)
end

-- Pull latest changes
git.pull({
    path = "/tmp/repo",
    branch = "main",
    rebase = false
})

-- Get repository status
local status = git.status({path = "/tmp/repo"})
print("Branch:", status.branch)
print("Clean:", status.is_clean)
print("Modified files:", #status.modified)
print("Untracked files:", #status.untracked)

-- Checkout branch or commit
git.checkout({
    path = "/tmp/repo",
    branch = "develop"
})

git.checkout({
    path = "/tmp/repo",
    commit = "abc123"
})

-- Create commit
git.commit({
    path = "/tmp/repo",
    message = "Update configuration",
    add_all = true,  -- git add -A
    author = "Bot <bot@example.com>"
})

-- Push to remote
git.push({
    path = "/tmp/repo",
    remote = "origin",
    branch = "main",
    force = false
})

-- Check if directory is a git repository
local is_repo = git.is_repo("/tmp/repo")
print("Is git repo:", is_repo)

-- Clean directory (remove if exists)
git.ensure_clean("/tmp/repo")  -- removes directory if it exists
```

### Complete Example

```lua
task({
    name = "deploy-from-git",
    run = function()
        local repo_path = "/tmp/my-app"
        local repo_url = "https://github.com/myorg/my-app.git"

        -- Clone or update repository
        if git.is_repo(repo_path) then
            print("Pulling latest changes...")
            git.pull({path = repo_path})
        else
            print("Cloning repository...")
            git.clone({
                url = repo_url,
                local_path = repo_path,
                branch = "main",
                depth = 1
            })
        end

        -- Checkout specific version
        git.checkout({
            path = repo_path,
            branch = values.version or "main"
        })

        -- Get commit info
        local status = git.status({path = repo_path})
        print("Deploying branch:", status.branch)
        print("Latest commit:", status.commit_sha)

        -- Build and deploy
        system.cd(repo_path)
        local build = system.exec("make", {"build"})
        if not build.success then
            return false, "Build failed"
        end

        local deploy = system.exec("make", {"deploy"})
        if not deploy.success then
            return false, "Deploy failed"
        end

        return true, "Deployed from commit " .. status.commit_sha
    end
})
```

---

## Sloth Module

**Category**: Automation
**Description**: Self-management automation for sloth-runner itself

### Agent Management

```lua
-- Install agent on remote host (idempotent)
sloth.agent.install({
    name = "web-01",
    ssh_host = "192.168.1.10",
    ssh_user = "root",
    ssh_port = 22,
    master = "192.168.1.29:50053",
    bind_address = "0.0.0.0",
    port = 50060,
    report_address = "192.168.1.10:50060"
})

-- Update agent to latest version
sloth.agent.update({
    name = "web-01"
})

-- List all agents
local agents = sloth.agent.list()
for _, agent in ipairs(agents) do
    print(agent.name, agent.status, agent.version)
end

-- Get agent details
local agent = sloth.agent.get({name = "web-01"})
print("Agent:", agent.name)
print("Status:", agent.status)
print("Address:", agent.address)

-- Delete agent
sloth.agent.delete({name = "web-01"})

-- Start local agent
sloth.agent.start({
    bind_address = "0.0.0.0",
    port = 50060,
    master = "192.168.1.29:50053"
})

-- Stop local agent
sloth.agent.stop()
```

### Workflow Management

```lua
-- Add workflow to database
sloth.workflow.add({
    name = "deploy",
    file = "/path/to/deploy.sloth",
    description = "Production deployment workflow",
    active = true,
    tags = {"deployment", "production"}
})

-- List all workflows
local workflows = sloth.workflow.list()
for _, wf in ipairs(workflows) do
    print(wf.name, wf.active and "active" or "inactive")
end

-- Get workflow details
local wf = sloth.workflow.get({name = "deploy"})
print("Workflow:", wf.name)
print("File:", wf.file)
print("Active:", wf.active)

-- Remove workflow
sloth.workflow.remove({name = "deploy"})

-- Activate/deactivate workflow
sloth.workflow.activate({name = "deploy"})
sloth.workflow.deactivate({name = "deploy"})
```

### SSH Profile Management

```lua
-- Add SSH profile
sloth.ssh.add({
    name = "prod-server",
    host = "192.168.1.10",
    user = "deploy",
    port = 22,
    key_file = "/home/user/.ssh/id_rsa"
})

-- List SSH profiles
local profiles = sloth.ssh.list()
for _, profile in ipairs(profiles) do
    print(profile.name, profile.host)
end

-- Remove SSH profile
sloth.ssh.remove({name = "prod-server"})
```

### Stack Management

```lua
-- List workflow stacks
local stacks = sloth.stack.list()
for _, stack in ipairs(stacks) do
    print(stack.name, stack.status)
end

-- Get stack details
local stack = sloth.stack.get({name = "my-stack"})
print("Stack:", stack.name)
print("Status:", stack.status)

-- Delete stack
sloth.stack.delete({name = "my-stack"})
```

### Run Workflow

```lua
-- Execute workflow
local result = sloth.run({
    sloth = "deploy",                    -- workflow name
    delegate_to = "web-01",              -- target agent
    values = "env=production,version=v2.0.0",  -- key=value pairs
    yes = true                           -- auto-confirm
})
```

### Complete Example

```lua
task({
    name = "bootstrap-infrastructure",
    run = function()
        -- Install agents on all servers
        local servers = {
            {name = "web-01", host = "192.168.1.10"},
            {name = "web-02", host = "192.168.1.11"},
            {name = "db-01", host = "192.168.1.20"}
        }

        for _, server in ipairs(servers) do
            print("Installing agent on", server.name)
            sloth.agent.install({
                name = server.name,
                ssh_host = server.host,
                ssh_user = "root",
                master = "192.168.1.29:50053",
                port = 50060
            })
        end

        -- Register deployment workflows
        sloth.workflow.add({
            name = "deploy-web",
            file = "/etc/sloth-runner/workflows/deploy-web.sloth",
            description = "Web server deployment",
            active = true
        })

        sloth.workflow.add({
            name = "deploy-db",
            file = "/etc/sloth-runner/workflows/deploy-db.sloth",
            description = "Database deployment",
            active = true
        })

        -- Run initial deployment
        sloth.run({
            sloth = "deploy-db",
            delegate_to = "db-01",
            values = "env=production",
            yes = true
        })

        sloth.run({
            sloth = "deploy-web",
            delegate_to = "web-01,web-02",
            values = "env=production",
            yes = true
        })

        return true, "Infrastructure bootstrapped successfully"
    end
})
```

---

## Incus Module

**Category**: Infrastructure
**Description**: Incus/LXD container and VM management with fluent API
**Full Documentation**: [incus.md](./incus.md)

### Instance Management

```lua
-- Create and launch instance (fluent API)
incus.instance("web-01")
    :image("ubuntu/22.04")
    :profile("default")
    :config({
        ["limits.cpu"] = "2",
        ["limits.memory"] = "2GB"
    })
    :device("web-data", {
        type = "disk",
        source = "/srv/www",
        path = "/var/www"
    })
    :proxy("http", "tcp:0.0.0.0:80", "tcp:127.0.0.1:8080")
    :delegate_to("lxd-host")
    :launch()  -- create and start

-- Manage existing instance
incus.instance("web-01")
    :delegate_to("lxd-host")
    :start()

incus.instance("web-01")
    :delegate_to("lxd-host")
    :stop()

-- Execute command in instance
local output = incus.exec({
    instance = "web-01",
    command = "systemctl status nginx",
    user = "root",
    cwd = "/var/www",
    delegate_to = "lxd-host"
})
print(output)
```

### Network Management

```lua
-- Create network
incus.network("br0")
    :type("bridge")
    :config({
        ["ipv4.address"] = "10.0.0.1/24",
        ["ipv4.nat"] = "true",
        ["ipv6.address"] = "none"
    })
    :delegate_to("lxd-host")
    :create()

-- Attach network to instance
incus.network("br0")
    :delegate_to("lxd-host")
    :attach("web-01")
```

### Quick Example

```lua
task({
    name = "deploy-containers",
    run = function()
        -- Create network
        incus.network("web-net")
            :type("bridge")
            :config({["ipv4.address"] = "10.10.0.1/24"})
            :delegate_to("lxd-host")
            :create()

        -- Deploy web servers
        local servers = {"web-01", "web-02", "web-03"}
        for _, name in ipairs(servers) do
            incus.instance(name)
                :image("ubuntu:22.04")
                :profile("default")
                :config({["limits.cpu"] = "2"})
                :delegate_to("lxd-host")
                :launch()

            -- Install nginx
            incus.exec({
                instance = name,
                command = "apt update && apt install -y nginx",
                delegate_to = "lxd-host"
            })
        end

        return true, "Deployed " .. #servers .. " web servers"
    end
})
```

---

## Best Practices

### 1. Error Handling

Always check for errors and handle them gracefully:

```lua
local result, err = some_function()
if err then
    return false, "Operation failed: " .. err
end
```

### 2. Use Environment Variables for Secrets

Never hardcode credentials:

```lua
local api_key = os.getenv("API_KEY")
if not api_key then
    error("API_KEY environment variable not set")
end
```

### 3. Idempotency

Design tasks to be idempotent (can be run multiple times safely):

```lua
-- Check before creating
if not system.exists("/tmp/mydir") then
    system.mkdir("/tmp/mydir", true)
end

-- Git clone with clean=false is idempotent
git.clone({
    url = "https://github.com/user/repo.git",
    local_path = "/tmp/repo",
    clean = false  -- won't re-clone if exists
})
```

### 4. Parallel Execution

Use goroutines for independent parallel tasks:

```lua
goroutine.pool_create("workers", {workers = 10})

for _, server in ipairs(servers) do
    goroutine.pool_submit("workers", function(srv)
        -- Deploy to server
    end, server)
end

goroutine.pool_wait("workers")
goroutine.pool_close("workers")
```

### 5. Monitoring and Metrics

Add metrics to track operations:

```lua
monitor.counter_inc("deployments", {env = "production"})

local timer = monitor.timer_start("deployment_duration")
-- ... perform deployment ...
monitor.timer_end(timer)
```

### 6. Notifications

Notify on success and failure:

```lua
if success then
    notify.slack(webhook, {text = "✅ Deployment successful"})
else
    notify.slack(webhook, {text = "❌ Deployment failed"})
end
```

---

## Complete Integration Example

This example demonstrates using multiple modules together:

```lua
task({
    name = "full-stack-deployment",
    run = function()
        -- Start metrics
        monitor.counter_inc("deployments_started", {env = values.env})
        local timer = monitor.timer_start("deployment_duration")

        -- Validate target system
        local mem = facts.get_memory({agent = values.target})
        if mem.used_percent > 90 then
            monitor.counter_inc("deployments_failed", {reason = "insufficient_memory"})
            return false, "Insufficient memory"
        end

        -- Clone/update repository
        local repo_path = "/tmp/my-app"
        if git.is_repo(repo_path) then
            git.pull({path = repo_path})
        else
            git.clone({
                url = "https://github.com/myorg/app.git",
                local_path = repo_path,
                branch = "main"
            })
        end

        -- Build application
        system.cd(repo_path)
        local build = system.exec("make", {"build"})
        if not build.success then
            monitor.counter_inc("deployments_failed", {reason = "build_error"})
            notify.slack(os.getenv("SLACK_WEBHOOK"), {
                text = "❌ Build failed",
                attachments = {{
                    color = "danger",
                    text = build.stderr
                }}
            })
            return false, "Build failed"
        end

        -- Deploy to containers in parallel
        goroutine.pool_create("deploy", {workers = 3})
        local servers = {"web-01", "web-02", "web-03"}

        for _, server in ipairs(servers) do
            goroutine.pool_submit("deploy", function(srv)
                -- Copy binary
                system.exec("scp", {"./app", srv .. ":/opt/app/app"})

                -- Restart service
                incus.exec({
                    instance = srv,
                    command = "systemctl restart app",
                    delegate_to = "lxd-host"
                })

                monitor.counter_inc("deployments_per_server", {server = srv})
            end, server)
        end

        goroutine.pool_wait("deploy")
        goroutine.pool_close("deploy")

        -- Record metrics
        local duration = monitor.timer_end(timer)
        monitor.counter_inc("deployments_successful")

        -- Send success notification
        notify.slack(os.getenv("SLACK_WEBHOOK"), {
            text = "✅ Deployment successful",
            attachments = {{
                color = "good",
                fields = {
                    {title = "Servers", value = #servers, short = true},
                    {title = "Duration", value = duration .. "ms", short = true}
                }
            }}
        })

        return true, "Deployed to " .. #servers .. " servers in " .. duration .. "ms"
    end
})
```

---

## See Also

- [Modern DSL Guide](../modern-dsl/introduction.md)
- [Core Concepts](../en/core-concepts.md)
- [Advanced Examples](../en/advanced-examples.md)
- [Individual Module Documentation](./index.md)

---

**Need help?** Use the `help` module:
```lua
help()              -- Show general help
help.modules()      -- List all modules
help.search("http") -- Search for functions
```
