-- Advanced Modules Showcase
-- This example demonstrates the new crypto, time, data, and database modules

print("🚀 Sloth Runner - Advanced Modules Showcase")
print("=" .. string.rep("=", 50))

-- 🔐 Crypto Module Demo
print("\n🔐 CRYPTO MODULE DEMO")
print("-" .. string.rep("-", 30))

-- Generate secure password for database
local db_password = crypto.generate_password(16, true)
print("Generated secure password:", db_password)

-- Hash the password for storage
local password_hash = crypto.sha256(db_password)
print("Password hash (SHA256):", password_hash)

-- Generate UUID for session
local session_id = crypto.uuid()
print("Session ID:", session_id)

-- Encrypt sensitive data
local encryption_key = "my-super-secret-key-32-chars-!"
local sensitive_data = "credit_card_number_1234567890"
local encrypted = crypto.aes_encrypt(encryption_key, sensitive_data)
print("Encrypted data:", encrypted)

-- Decrypt for verification
local decrypted = crypto.aes_decrypt(encryption_key, encrypted)
print("Decrypted data:", decrypted)

-- 📅 Time Module Demo
print("\n📅 TIME MODULE DEMO")
print("-" .. string.rep("-", 30))

-- Current time operations
local now = time.now()
print("Current timestamp:", now)
print("Formatted time:", time.format(now, "2006-01-02 15:04:05"))

-- Schedule future task
local future_task = time.add(now, "2h30m")
print("Task scheduled for:", time.format(future_task, "2006-01-02 15:04:05"))

-- Calculate duration until task
local seconds_until = time.until(future_task)
print("Seconds until task:", math.floor(seconds_until))

-- Time zone handling
local utc_time = time.utc(now)
print("UTC time:", time.rfc3339(utc_time))

-- 📊 Data Module Demo
print("\n📊 DATA MODULE DEMO")
print("-" .. string.rep("-", 30))

-- Create complex data structure
local deployment_config = {
    app = {
        name = "sloth-runner",
        version = "3.1.0",
        environment = "production"
    },
    database = {
        host = "db.example.com",
        port = 5432,
        ssl_enabled = true
    },
    deployment = {
        timestamp = now,
        user = "deploy-bot",
        session_id = session_id
    }
}

-- Convert to different formats
local json_config = data.json_pretty(deployment_config)
print("Deployment config (JSON):")
print(json_config)

local yaml_config = data.yaml_encode(deployment_config)
print("\nDeployment config (YAML):")
print(yaml_config)

-- Extract specific values using path
local app_name = data.get_path(deployment_config, "app.name")
local db_port = data.get_path(deployment_config, "database.port")
print("\nExtracted values:")
print("App name:", app_name)
print("DB port:", db_port)

-- Flatten configuration for environment variables
local flattened = data.flatten(deployment_config, "_")
print("\nFlattened config (for env vars):")
print(data.json_pretty(flattened))

-- 🗄️ Database Module Demo
print("\n🗄️ DATABASE MODULE DEMO")
print("-" .. string.rep("-", 30))

-- Connect to SQLite database
local connected, conn_name = db.connect("sqlite3", ":memory:")
if connected then
    print("Connected to in-memory database:", conn_name)
    
    -- Create deployment logs table
    local create_result = db.exec([[
        CREATE TABLE deployment_logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            app_name TEXT NOT NULL,
            version TEXT NOT NULL,
            environment TEXT NOT NULL,
            session_id TEXT NOT NULL,
            status TEXT NOT NULL,
            created_at INTEGER NOT NULL,
            metadata TEXT
        )
    ]])
    
    if create_result then
        print("Deployment logs table created successfully")
        
        -- Insert deployment record using transaction
        local tx_success = db.transaction(function()
            -- Insert main deployment record
            local insert_result = db.exec([[
                INSERT INTO deployment_logs 
                (app_name, version, environment, session_id, status, created_at, metadata)
                VALUES (?, ?, ?, ?, ?, ?, ?)
            ]], app_name, "3.1.0", "production", session_id, "started", now, json_config)
            
            if not insert_result then
                return false
            end
            
            print("Deployment record inserted with ID:", insert_result.last_insert_id)
            
            -- Insert completion record
            local completion_time = time.add(now, "5m")
            local completion_result = db.exec([[
                INSERT INTO deployment_logs 
                (app_name, version, environment, session_id, status, created_at, metadata)
                VALUES (?, ?, ?, ?, ?, ?, ?)
            ]], app_name, "3.1.0", "production", session_id, "completed", completion_time, 
            data.json_encode({duration_seconds = 300, success = true}))
            
            return completion_result ~= nil
        end)
        
        if tx_success then
            print("Deployment transaction completed successfully")
            
            -- Query deployment history
            local deployments = db.query([[
                SELECT id, app_name, version, status, 
                       datetime(created_at, 'unixepoch') as created_at_formatted
                FROM deployment_logs 
                WHERE session_id = ?
                ORDER BY created_at
            ]], session_id)
            
            if deployments and #deployments > 0 then
                print("\nDeployment history:")
                for i = 1, #deployments do
                    local deployment = deployments[i]
                    print(string.format("  %s: %s v%s - %s at %s", 
                        deployment.id, deployment.app_name, deployment.version, 
                        deployment.status, deployment.created_at_formatted))
                end
            end
        else
            print("Deployment transaction failed")
        end
    else
        print("Failed to create deployment logs table")
    end
    
    -- Close database connection
    db.close()
    print("Database connection closed")
else
    print("Failed to connect to database:", conn_name)
end

-- 🎯 Integration Example
print("\n🎯 INTEGRATION EXAMPLE")
print("-" .. string.rep("-", 30))

-- Create audit log entry combining all modules
local audit_entry = {
    event_id = crypto.uuid(),
    timestamp = time.rfc3339(now),
    event_type = "module_showcase_completed",
    security = {
        session_id = session_id,
        data_hash = crypto.sha256(json_config)
    },
    metadata = {
        modules_demonstrated = {"crypto", "time", "data", "database"},
        total_duration_seconds = math.floor(time.since(now)),
        data_formats_used = {"json", "yaml", "sql"}
    }
}

print("Audit log entry:")
print(data.json_pretty(audit_entry))

-- Simulate storing audit log in secure format
local audit_json = data.json_encode(audit_entry)
local audit_encrypted = crypto.aes_encrypt(encryption_key, audit_json)
print("\nAudit log encrypted and ready for secure storage")
print("Encrypted size:", string.len(audit_encrypted), "characters")

print("\n✅ Advanced modules showcase completed successfully!")
print("All modules are working together seamlessly. 🎉")