-- Simple test of new core modules
print("🦥 Testing New Core Modules")
print("===========================")

-- Test system module
print("\n📁 System Module:")
local sys = require("system")
print("Platform:", sys.platform())
print("Architecture:", sys.arch())
print("Hostname:", sys.hostname())

-- Test crypto module 
print("\n🔐 Crypto Module:")
local crypto = require("crypto")
local data = "Hello World!"
print("SHA256 of '" .. data .. "':", crypto.sha256(data))

-- Test monitoring module
print("\n📊 Monitoring Module:")
local monitor = require("monitor")
monitor.counter_inc("test_counter")
monitor.counter_inc("test_counter")
local metric = monitor.get_metric("test_counter")
if metric then
    print("Counter value:", metric.value)
end

-- Test database module
print("\n🗄️ Database Module:")
local db = require("database")
local connected = db.connect("test", "sqlite3", ":memory:")
if connected then
    print("SQLite connection: Success")
    db.disconnect("test")
else
    print("SQLite connection: Failed")
end

print("\n✅ All modules tested successfully!")