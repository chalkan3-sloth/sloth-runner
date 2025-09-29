package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/yuin/gopher-lua"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DatabaseModule provides database operations
type DatabaseModule struct {
	info        CoreModuleInfo
	connections map[string]*sql.DB
}

// NewDatabaseModule creates a new database module
func NewDatabaseModule() *DatabaseModule {
	info := CoreModuleInfo{
		Name:         "database",
		Version:      "1.0.0",
		Description:  "Database operations with support for SQLite, MySQL, PostgreSQL",
		Author:       "Sloth Runner Team",
		Category:     "core",
		Dependencies: []string{},
	}

	return &DatabaseModule{
		info:        info,
		connections: make(map[string]*sql.DB),
	}
}

// Info returns module information
func (d *DatabaseModule) Info() CoreModuleInfo {
	return d.info
}

// Loader loads the database module into Lua
func (d *DatabaseModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"connect":    d.luaConnect,
		"disconnect": d.luaDisconnect,
		"query":      d.luaQuery,
		"exec":       d.luaExec,
		"transaction": d.luaTransaction,
		"prepare":    d.luaPrepare,
		"ping":       d.luaPing,
		"close_all":  d.luaCloseAll,
	})

	L.Push(mod)
	return 1
}

// luaConnect connects to a database
func (d *DatabaseModule) luaConnect(L *lua.LState) int {
	name := L.CheckString(1)
	driver := L.CheckString(2)
	dsn := L.CheckString(3)
	
	db, err := sql.Open(driver, dsn)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Close existing connection if any
	if existing, exists := d.connections[name]; exists {
		existing.Close()
	}
	
	d.connections[name] = db
	L.Push(lua.LBool(true))
	return 1
}

// luaDisconnect disconnects from a database
func (d *DatabaseModule) luaDisconnect(L *lua.LState) int {
	name := L.CheckString(1)
	
	if db, exists := d.connections[name]; exists {
		err := db.Close()
		delete(d.connections, name)
		L.Push(lua.LBool(err == nil))
		if err != nil {
			L.Push(lua.LString(err.Error()))
			return 2
		}
	} else {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("Connection not found"))
		return 2
	}
	
	return 1
}

// luaQuery executes a SELECT query
func (d *DatabaseModule) luaQuery(L *lua.LState) int {
	name := L.CheckString(1)
	query := L.CheckString(2)
	
	db, exists := d.connections[name]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("Connection not found"))
		return 2
	}
	
	// Get parameters if provided
	var args []interface{}
	if L.GetTop() > 2 {
		if paramsTable := L.CheckTable(3); paramsTable != nil {
			paramsTable.ForEach(func(_, value lua.LValue) {
				args = append(args, d.luaValueToInterface(value))
			})
		}
	}
	
	rows, err := db.Query(query, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer rows.Close()
	
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Prepare result table
	result := L.NewTable()
	rowIndex := 1
	
	for rows.Next() {
		// Create slice of interface{} to hold row values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		
		// Scan the row
		if err := rows.Scan(valuePtrs...); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		
		// Create row table
		row := L.NewTable()
		for i, col := range columns {
			row.RawSetString(col, d.interfaceToLuaValue(L, values[i]))
		}
		
		result.RawSetInt(rowIndex, row)
		rowIndex++
	}
	
	if err := rows.Err(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(result)
	return 1
}

// luaExec executes an INSERT, UPDATE, or DELETE query
func (d *DatabaseModule) luaExec(L *lua.LState) int {
	name := L.CheckString(1)
	query := L.CheckString(2)
	
	db, exists := d.connections[name]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("Connection not found"))
		return 2
	}
	
	// Get parameters if provided
	var args []interface{}
	if L.GetTop() > 2 {
		if paramsTable := L.CheckTable(3); paramsTable != nil {
			paramsTable.ForEach(func(_, value lua.LValue) {
				args = append(args, d.luaValueToInterface(value))
			})
		}
	}
	
	result, err := db.Exec(query, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	resultTable := L.NewTable()
	
	if rowsAffected, err := result.RowsAffected(); err == nil {
		resultTable.RawSetString("rows_affected", lua.LNumber(rowsAffected))
	}
	
	if lastInsertId, err := result.LastInsertId(); err == nil {
		resultTable.RawSetString("last_insert_id", lua.LNumber(lastInsertId))
	}
	
	L.Push(resultTable)
	return 1
}

// luaTransaction executes multiple queries in a transaction
func (d *DatabaseModule) luaTransaction(L *lua.LState) int {
	name := L.CheckString(1)
	queriesTable := L.CheckTable(2)
	
	db, exists := d.connections[name]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("Connection not found"))
		return 2
	}
	
	tx, err := db.Begin()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Execute queries
	results := L.NewTable()
	resultIndex := 1
	
	queriesTable.ForEach(func(_, value lua.LValue) {
		if queryTable, ok := value.(*lua.LTable); ok {
			query := queryTable.RawGetString("query").String()
			
			var args []interface{}
			if paramsValue := queryTable.RawGetString("params"); paramsValue != lua.LNil {
				if paramsTable, ok := paramsValue.(*lua.LTable); ok {
					paramsTable.ForEach(func(_, paramValue lua.LValue) {
						args = append(args, d.luaValueToInterface(paramValue))
					})
				}
			}
			
			result, err := tx.Exec(query, args...)
			if err != nil {
				tx.Rollback()
				L.Push(lua.LBool(false))
				L.Push(lua.LString(err.Error()))
				return
			}
			
			resultTable := L.NewTable()
			if rowsAffected, err := result.RowsAffected(); err == nil {
				resultTable.RawSetString("rows_affected", lua.LNumber(rowsAffected))
			}
			if lastInsertId, err := result.LastInsertId(); err == nil {
				resultTable.RawSetString("last_insert_id", lua.LNumber(lastInsertId))
			}
			
			results.RawSetInt(resultIndex, resultTable)
			resultIndex++
		}
	})
	
	if err := tx.Commit(); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(results)
	return 2
}

// luaPrepare prepares a statement
func (d *DatabaseModule) luaPrepare(L *lua.LState) int {
	name := L.CheckString(1)
	query := L.CheckString(2)
	
	db, exists := d.connections[name]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("Connection not found"))
		return 2
	}
	
	stmt, err := db.Prepare(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// For simplicity, we'll execute and close the statement immediately
	// In a more sophisticated implementation, we'd manage prepared statements
	defer stmt.Close()
	
	L.Push(lua.LString("prepared")) // Placeholder
	return 1
}

// luaPing tests the database connection
func (d *DatabaseModule) luaPing(L *lua.LState) int {
	name := L.CheckString(1)
	
	db, exists := d.connections[name]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("Connection not found"))
		return 2
	}
	
	err := db.Ping()
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaCloseAll closes all database connections
func (d *DatabaseModule) luaCloseAll(L *lua.LState) int {
	for name, db := range d.connections {
		db.Close()
		delete(d.connections, name)
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// luaValueToInterface converts Lua value to interface{}
func (d *DatabaseModule) luaValueToInterface(value lua.LValue) interface{} {
	switch v := value.(type) {
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LNilType:
		return nil
	default:
		return value.String()
	}
}

// interfaceToLuaValue converts interface{} to Lua value
func (d *DatabaseModule) interfaceToLuaValue(L *lua.LState, value interface{}) lua.LValue {
	if value == nil {
		return lua.LNil
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return lua.LString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(v.Uint())
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(v.Float())
	case reflect.Bool:
		return lua.LBool(v.Bool())
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte -> string
			return lua.LString(string(value.([]byte)))
		}
		fallthrough
	default:
		// For complex types, convert to JSON string
		if jsonBytes, err := json.Marshal(value); err == nil {
			return lua.LString(string(jsonBytes))
		}
		return lua.LString(fmt.Sprintf("%v", value))
	}
}