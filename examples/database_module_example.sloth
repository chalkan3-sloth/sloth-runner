-- Database Module Examples

-- SQLite connection
local connected, conn_id = db.connect("sqlite3", "./example.db")
if not connected then
    print("Failed to connect:", conn_id)
    return
end

print("Connected to database:", conn_id)

-- Create table
local result = db.exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, email TEXT, age INTEGER)")
if result then
    print("Table created successfully")
    print("Rows affected:", result.rows_affected)
else
    print("Failed to create table")
end

-- Insert data
local insert_result = db.exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", "John Doe", "john@example.com", 30)
if insert_result then
    print("User inserted with ID:", insert_result.last_insert_id)
end

-- Insert more data
db.exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", "Jane Smith", "jane@example.com", 25)
db.exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", "Bob Johnson", "bob@example.com", 35)

-- Query data
local users = db.query("SELECT * FROM users WHERE age > ?", 20)
if users then
    print("Found", #users, "users:")
    for i = 1, #users do
        local user = users[i]
        print(string.format("ID: %s, Name: %s, Email: %s, Age: %s", 
                          user.id, user.name, user.email, user.age))
    end
else
    print("Failed to query users")
end

-- Transaction example
local tx_success = db.transaction(function(tx)
    -- These operations will be rolled back if any fails
    local result1 = db.exec("UPDATE users SET age = age + 1 WHERE name = ?", "John Doe")
    local result2 = db.exec("UPDATE users SET age = age + 1 WHERE name = ?", "Jane Smith")
    
    if not result1 or not result2 then
        return false -- This will cause a rollback
    end
    
    return true -- Commit the transaction
end)

if tx_success then
    print("Transaction completed successfully")
else
    print("Transaction failed and was rolled back")
end

-- Test connection
local ping_ok = db.ping()
print("Database connection is healthy:", ping_ok)

-- Close connection
local closed = db.close()
print("Database connection closed:", closed)

-- PostgreSQL example (commented out - requires PostgreSQL)
--[[
local pg_connected, pg_conn = db.connect("postgres", "user=username password=password dbname=mydb sslmode=disable", "pg_conn")
if pg_connected then
    print("Connected to PostgreSQL:", pg_conn)
    
    -- PostgreSQL specific operations
    local pg_result = db.query("SELECT version()", "pg_conn")
    if pg_result and #pg_result > 0 then
        print("PostgreSQL version:", pg_result[1].version)
    end
    
    db.close("pg_conn")
end
--]]

-- MySQL example (commented out - requires MySQL)
--[[
local mysql_connected, mysql_conn = db.connect("mysql", "username:password@tcp(localhost:3306)/dbname", "mysql_conn")
if mysql_connected then
    print("Connected to MySQL:", mysql_conn)
    
    -- MySQL specific operations
    local mysql_result = db.query("SELECT VERSION() as version", "mysql_conn")
    if mysql_result and #mysql_result > 0 then
        print("MySQL version:", mysql_result[1].version)
    end
    
    db.close("mysql_conn")
end
--]]