-- Sloth Runner - Advanced Core Modules Examples
-- This file demonstrates all the new native core modules

print("🦥 Sloth Runner - Core Modules Demo")
print("===================================")

-- System Module Examples
print("\n📁 System Module Examples:")
print("---------------------------")

-- Get system information
local sys = require("system")
print("Platform:", sys.platform())
print("Architecture:", sys.arch())
print("CPU Count:", sys.cpu_count())
print("Hostname:", sys.hostname())
print("Current Directory:", sys.pwd())

-- File operations
local test_file = "/tmp/sloth_test.txt"
if sys.exists(test_file) then
    sys.rmdir(test_file, false)
end

-- Create a test file
local result = sys.exec("echo", {"Hello from Sloth Runner!"})
if result.success then
    print("Command executed successfully:", result.output)
end

-- Environment variables
sys.setenv("SLOTH_TEST", "hello world")
print("Environment variable SLOTH_TEST:", sys.env("SLOTH_TEST"))

-- Crypto Module Examples
print("\n🔐 Crypto Module Examples:")
print("----------------------------")

local crypto = require("crypto")

-- Hashing examples
local data = "Hello, Sloth Runner!"
print("Original:", data)
print("MD5:", crypto.md5(data))
print("SHA256:", crypto.sha256(data))
print("SHA512:", crypto.sha512(data))

-- Password hashing
local password = "mySecretPassword123"
local hash = crypto.bcrypt_hash(password, 12)
print("BCrypt hash:", hash)
print("Password verify:", crypto.bcrypt_check(password, hash))

-- Encoding/Decoding
local encoded = crypto.base64_encode(data)
print("Base64 encoded:", encoded)
print("Base64 decoded:", crypto.base64_decode(encoded))

-- Random generation
print("Random bytes (16):", crypto.random_bytes(16))
print("Random string (32):", crypto.random_string(32))

-- AES Encryption
local key = "myEncryptionKey123456789012345678"
local encrypted = crypto.aes_encrypt(data, key)
print("AES encrypted:", encrypted)
if encrypted then
    local decrypted = crypto.aes_decrypt(encrypted, key)
    print("AES decrypted:", decrypted)
end

-- Notification Module Examples  
print("\n📧 Notification Module Examples:")
print("----------------------------------")

local notify = require("notify")

-- Slack notification example (webhook URL would be real)
--[[
local slack_success = notify.slack("https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK", {
    text = "Hello from Sloth Runner! 🦥",
    username = "Sloth Bot",
    channel = "#general",
    icon_emoji = ":sloth:",
    attachments = {
        {
            color = "good",
            title = "Task Completed",
            text = "All systems are running smoothly!",
            footer = "Sloth Runner",
            fields = {
                {title = "Status", value = "Success", short = true},
                {title = "Duration", value = "2.5s", short = true}
            }
        }
    }
})
--]]

-- Email notification example
--[[
local email_success = notify.email({
    smtp_host = "smtp.gmail.com",
    smtp_port = "587",
    username = "your.email@gmail.com",
    password = "your_app_password",
    from = "your.email@gmail.com",
    to = "recipient@example.com",
    subject = "Sloth Runner Notification",
    body = "Hello from Sloth Runner! All tasks completed successfully."
})
--]]

-- Generic webhook example
--[[
local webhook_success = notify.webhook("https://httpbin.org/post", {
    message = "Hello from Sloth Runner!",
    status = "success",
    timestamp = os.time()
})
--]]

print("Notification examples prepared (commented out - need real endpoints)")

-- Database Module Examples
print("\n🗄️  Database Module Examples:")
print("-------------------------------")

local db = require("database")

-- SQLite example
local db_success = db.connect("test", "sqlite3", ":memory:")
if db_success then
    print("Connected to SQLite database")
    
    -- Create table
    local create_result = db.exec("test", [[
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    ]])
    
    if create_result then
        print("Table created successfully")
        
        -- Insert data
        local insert_result = db.exec("test", 
            "INSERT INTO users (name, email) VALUES (?, ?)", 
            {"John Doe", "john@example.com"}
        )
        
        if insert_result then
            print("Inserted user, Last ID:", insert_result.last_insert_id)
            
            -- Query data
            local users = db.query("test", "SELECT * FROM users WHERE name = ?", {"John Doe"})
            if users and #users > 0 then
                print("Found user:", users[1].name, users[1].email)
            end
        end
    end
    
    -- Transaction example
    local tx_success, tx_results = db.transaction("test", {
        {
            query = "INSERT INTO users (name, email) VALUES (?, ?)",
            params = {"Jane Doe", "jane@example.com"}
        },
        {
            query = "INSERT INTO users (name, email) VALUES (?, ?)", 
            params = {"Bob Smith", "bob@example.com"}
        }
    })
    
    if tx_success then
        print("Transaction completed successfully")
        
        -- Query all users
        local all_users = db.query("test", "SELECT COUNT(*) as count FROM users")
        if all_users then
            print("Total users:", all_users[1].count)
        end
    end
    
    db.disconnect("test")
else
    print("Failed to connect to database")
end

-- Monitoring Module Examples
print("\n📊 Monitoring Module Examples:")
print("--------------------------------")

local monitor = require("monitor")

-- Counter examples
monitor.counter_inc("http_requests_total", {method = "GET", endpoint = "/api"})
monitor.counter_inc("http_requests_total", {method = "GET", endpoint = "/api"})
monitor.counter_add("http_requests_total", 5, {method = "POST", endpoint = "/api"})

-- Gauge examples
monitor.gauge_set("memory_usage_bytes", 1024*1024*512) -- 512MB
monitor.gauge_inc("active_connections", 1)
monitor.gauge_inc("active_connections", 3)
monitor.gauge_dec("active_connections", 1)

-- Timer examples
local timer_key = monitor.timer_start("task_duration", {task = "data_processing"})
-- Simulate some work
local start_time = os.clock()
while os.clock() - start_time < 0.1 do end -- Wait 100ms
local duration = monitor.timer_end(timer_key)
print("Task duration:", duration, "seconds")

-- Histogram examples
monitor.histogram_observe("request_duration_seconds", 0.150, {method = "GET"})
monitor.histogram_observe("request_duration_seconds", 0.300, {method = "POST"})
monitor.histogram_observe("request_duration_seconds", 0.075, {method = "GET"})

-- Get metrics
local metric = monitor.get_metric("http_requests_total", {method = "GET", endpoint = "/api"})
if metric then
    print("HTTP GET requests:", metric.value)
end

-- List all metrics
local all_metrics = monitor.list_metrics()
print("Total metrics tracked:", #all_metrics)

-- System metrics
local sys_metrics = monitor.system_metrics()
print("System - Goroutines:", sys_metrics.goroutines)
print("System - Memory Alloc:", sys_metrics.memory_alloc, "bytes")

-- Export examples
local prometheus_export = monitor.export_prometheus()
print("\nPrometheus Export (first 200 chars):")
print(string.sub(prometheus_export, 1, 200) .. "...")

print("\n✅ All core modules demonstrated successfully!")
print("🦥 Sloth Runner is ready for advanced task automation!")