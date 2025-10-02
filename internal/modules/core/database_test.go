package core

import (
	"database/sql"
	"os"
	"testing"

	"github.com/yuin/gopher-lua"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (string, func()) {
	dbFile := "/tmp/test_sloth_runner_" + t.Name() + ".db"
	
	// Remove if exists
	os.Remove(dbFile)
	
	// Create test database
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	
	// Create test table
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT,
			age INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	
	// Insert test data
	_, err = db.Exec(`INSERT INTO users (name, email, age) VALUES ('Alice', 'alice@example.com', 30)`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
	
	_, err = db.Exec(`INSERT INTO users (name, email, age) VALUES ('Bob', 'bob@example.com', 25)`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
	
	db.Close()
	
	cleanup := func() {
		os.Remove(dbFile)
	}
	
	return dbFile, cleanup
}

func TestDatabaseModule_Info(t *testing.T) {
	module := NewDatabaseModule()
	info := module.Info()

	if info.Name != "database" {
		t.Errorf("Expected module name 'database', got '%s'", info.Name)
	}

	if info.Category != "core" {
		t.Errorf("Expected category 'core', got '%s'", info.Category)
	}
}

func TestDatabaseModule_Connect(t *testing.T) {
	dbFile, cleanup := setupTestDB(t)
	defer cleanup()

	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		local success, err = db.connect("test", "sqlite3", "` + dbFile + `")
		
		if not success then
			error("Failed to connect: " .. tostring(err))
		end
		
		-- Test ping
		local pingOk = db.ping("test")
		if not pingOk then
			error("Ping failed")
		end
		
		-- Disconnect
		db.disconnect("test")
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestDatabaseModule_Query(t *testing.T) {
	dbFile, cleanup := setupTestDB(t)
	defer cleanup()

	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		db.connect("test", "sqlite3", "` + dbFile + `")
		
		local results = db.query("test", "SELECT * FROM users ORDER BY id")
		
		if not results then
			error("Query failed")
		end
		
		-- Check we got 2 rows
		local count = 0
		for k, v in pairs(results) do
			if type(k) == "number" then
				count = count + 1
			end
		end
		
		if count ~= 2 then
			error("Expected 2 rows, got " .. count)
		end
		
		-- Check first row
		if results[1].name ~= "Alice" then
			error("Expected name 'Alice', got " .. tostring(results[1].name))
		end
		
		if results[1].age ~= 30 then
			error("Expected age 30, got " .. tostring(results[1].age))
		end
		
		db.disconnect("test")
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestDatabaseModule_QueryWithParams(t *testing.T) {
	dbFile, cleanup := setupTestDB(t)
	defer cleanup()

	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		db.connect("test", "sqlite3", "` + dbFile + `")
		
		local results = db.query("test", "SELECT * FROM users WHERE name = ?", {"Bob"})
		
		if not results or not results[1] then
			error("Query with params failed")
		end
		
		if results[1].name ~= "Bob" then
			error("Expected name 'Bob', got " .. tostring(results[1].name))
		end
		
		if results[1].age ~= 25 then
			error("Expected age 25, got " .. tostring(results[1].age))
		end
		
		db.disconnect("test")
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestDatabaseModule_Exec(t *testing.T) {
	dbFile, cleanup := setupTestDB(t)
	defer cleanup()

	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		db.connect("test", "sqlite3", "` + dbFile + `")
		
		-- Insert
		local result = db.exec("test", "INSERT INTO users (name, email, age) VALUES (?, ?, ?)", 
			{"Charlie", "charlie@example.com", 35})
		
		if not result then
			error("Insert failed")
		end
		
		if result.rows_affected ~= 1 then
			error("Expected 1 row affected, got " .. tostring(result.rows_affected))
		end
		
		if not result.last_insert_id or result.last_insert_id <= 0 then
			error("Expected valid last_insert_id")
		end
		
		-- Update
		local updateResult = db.exec("test", "UPDATE users SET age = ? WHERE name = ?", 
			{40, "Charlie"})
		
		if updateResult.rows_affected ~= 1 then
			error("Expected 1 row affected by update")
		end
		
		-- Verify
		local results = db.query("test", "SELECT * FROM users WHERE name = ?", {"Charlie"})
		if results[1].age ~= 40 then
			error("Update didn't work, age is " .. tostring(results[1].age))
		end
		
		db.disconnect("test")
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestDatabaseModule_Transaction(t *testing.T) {
	dbFile, cleanup := setupTestDB(t)
	defer cleanup()

	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		db.connect("test", "sqlite3", "` + dbFile + `")
		
		local queries = {
			{
				query = "INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
				params = {"David", "david@example.com", 28}
			},
			{
				query = "INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
				params = {"Eve", "eve@example.com", 32}
			},
			{
				query = "UPDATE users SET age = age + 1 WHERE name = ?",
				params = {"Alice"}
			}
		}
		
		local success, results = db.transaction("test", queries)
		
		if not success then
			error("Transaction failed: " .. tostring(results))
		end
		
		-- Verify all changes
		local allUsers = db.query("test", "SELECT COUNT(*) as count FROM users")
		if allUsers[1].count ~= 4 then  -- 2 original + 2 new
			error("Expected 4 users total, got " .. tostring(allUsers[1].count))
		end
		
		-- Verify Alice's age was incremented
		local alice = db.query("test", "SELECT age FROM users WHERE name = ?", {"Alice"})
		if alice[1].age ~= 31 then  -- Was 30, now 31
			error("Expected Alice's age to be 31, got " .. tostring(alice[1].age))
		end
		
		db.disconnect("test")
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestDatabaseModule_Disconnect(t *testing.T) {
	dbFile, cleanup := setupTestDB(t)
	defer cleanup()

	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		db.connect("test", "sqlite3", "` + dbFile + `")
		
		local success = db.disconnect("test")
		if not success then
			error("Disconnect failed")
		end
		
		-- Try to disconnect again (should fail)
		local success2, err = db.disconnect("test")
		if success2 then
			error("Expected disconnect to fail for non-existent connection")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestDatabaseModule_CloseAll(t *testing.T) {
	dbFile1, cleanup1 := setupTestDB(t)
	defer cleanup1()
	
	dbFile2 := "/tmp/test_sloth_runner_" + t.Name() + "_2.db"
	os.Remove(dbFile2)
	defer os.Remove(dbFile2)
	
	// Create second database
	db2, _ := sql.Open("sqlite3", dbFile2)
	db2.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY)`)
	db2.Close()

	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		db.connect("db1", "sqlite3", "` + dbFile1 + `")
		db.connect("db2", "sqlite3", "` + dbFile2 + `")
		
		-- Both should be pingable
		if not db.ping("db1") then
			error("db1 not connected")
		end
		if not db.ping("db2") then
			error("db2 not connected")
		end
		
		-- Close all
		db.close_all()
		
		-- Both should now fail to ping
		local ping1 = db.ping("db1")
		local ping2 = db.ping("db2")
		
		if ping1 or ping2 then
			error("Expected both connections to be closed")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestDatabaseModule_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		
		-- Try to query without connecting
		local results, err = db.query("nonexistent", "SELECT * FROM users")
		
		if results then
			error("Expected query to fail with no connection")
		end
		
		if not err then
			error("Expected error message")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestDatabaseModule_InvalidDSN(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewDatabaseModule()
	L.PreloadModule("database", module.Loader)

	code := `
		local db = require("database")
		
		local success, err = db.connect("test", "sqlite3", "/invalid/path/that/does/not/exist.db")
		
		if success then
			error("Expected connection to fail with invalid DSN")
		end
		
		if not err then
			error("Expected error message")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}
