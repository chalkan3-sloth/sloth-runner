-- Sloth Runner - Core Modules Test Task

-- Test system module
local test_system = task("test_system_module")
    :description("Test system module functionality")
    :command(function(params)
        print("🦥 Testing System Module")
        local sys = require("system")
        
        print("Platform:", sys.platform())
        print("Architecture:", sys.arch()) 
        print("CPU Count:", sys.cpu_count())
        print("Hostname:", sys.hostname())
        
        -- Test command execution
        local result = sys.exec("echo", {"Hello from system module!"})
        if result.success then
            print("Command output:", result.output)
        end
        
        -- Test environment
        sys.setenv("TEST_VAR", "hello")
        print("Environment test:", sys.env("TEST_VAR"))
        
        return true, "System module test completed", {success = true}
    end)
    :timeout("30s")
    :build()

-- Test crypto module
local test_crypto = task("test_crypto_module")
    :description("Test crypto module functionality")
    :depends_on("test_system_module")
    :command(function(params)
        print("🔐 Testing Crypto Module")
        local crypto = require("crypto")
        
        local data = "Hello Sloth Runner!"
        print("Original:", data)
        print("SHA256:", crypto.sha256(data))
        
        -- Test base64
        local encoded = crypto.base64_encode(data)
        print("Base64 encoded:", encoded)
        print("Base64 decoded:", crypto.base64_decode(encoded))
        
        -- Test random
        print("Random string:", crypto.random_string(16))
        
        return true, "Crypto module test completed", {success = true}
    end)
    :timeout("30s")
    :build()

-- Test monitoring module
local test_monitoring = task("test_monitoring_module")
    :description("Test monitoring module functionality")
    :depends_on("test_crypto_module")
    :command(function(params)
        print("📊 Testing Monitoring Module")
        local monitor = require("monitor")
        
        -- Test counter
        monitor.counter_inc("test_counter")
        monitor.counter_inc("test_counter")
        local counter_metric = monitor.get_metric("test_counter")
        if counter_metric then
            print("Counter value:", counter_metric.value)
        end
        
        -- Test gauge
        monitor.gauge_set("test_gauge", 42.5)
        local gauge_metric = monitor.get_metric("test_gauge")
        if gauge_metric then
            print("Gauge value:", gauge_metric.value)
        end
        
        -- Test timer
        local timer_key = monitor.timer_start("test_timer")
        -- Simulate work
        local start = os.clock()
        while os.clock() - start < 0.01 do end
        local duration = monitor.timer_end(timer_key)
        print("Timer duration:", duration, "seconds")
        
        -- List all metrics
        local metrics = monitor.list_metrics()
        print("Total metrics:", #metrics)
        
        return true, "Monitoring module test completed", {success = true}
    end)
    :timeout("30s")
    :build()

-- Test database module
local test_database = task("test_database_module")
    :description("Test database module functionality")
    :depends_on("test_monitoring_module")
    :command(function(params)
        print("🗄️ Testing Database Module")
        local db = require("database")
        
        -- Connect to in-memory SQLite
        local connected = db.connect("test", "sqlite3", ":memory:")
        if connected then
            print("Connected to SQLite database")
            
            -- Create table
            local create_result = db.exec("test", "CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)")
            if create_result then
                print("Table created successfully")
                
                -- Insert data
                local insert_result = db.exec("test", "INSERT INTO test_table (name) VALUES (?)", {"Test Entry"})
                if insert_result then
                    print("Data inserted, ID:", insert_result.last_insert_id)
                    
                    -- Query data
                    local results = db.query("test", "SELECT * FROM test_table")
                    if results and #results > 0 then
                        print("Found entry:", results[1].name)
                    end
                end
            end
            
            db.disconnect("test")
            print("Database test completed")
        else
            print("Failed to connect to database")
        end
        
        return true, "Database module test completed", {success = true}
    end)
    :timeout("30s")
    :build()

-- Create the group
local test_group = group("test_modules")
    :description("Test new core modules functionality")
    :add_task(test_system)
    :add_task(test_crypto)
    :add_task(test_monitoring)
    :add_task(test_database)
    :build()