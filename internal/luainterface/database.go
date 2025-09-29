package luainterface

import (
	"database/sql"
	"fmt"
	"reflect"

	lua "github.com/yuin/gopher-lua"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DatabaseModule provides database connectivity for Lua scripts
type DatabaseModule struct {
	connections map[string]*sql.DB
}

// NewDatabaseModule creates a new database module
func NewDatabaseModule() *DatabaseModule {
	return &DatabaseModule{
		connections: make(map[string]*sql.DB),
	}
}

// RegisterDatabaseModule registers the database module with the Lua state
func RegisterDatabaseModule(L *lua.LState) {
	module := NewDatabaseModule()
	
	// Create the db table
	dbTable := L.NewTable()
	
	// Connection management
	L.SetField(dbTable, "connect", L.NewFunction(module.luaConnect))
	L.SetField(dbTable, "close", L.NewFunction(module.luaClose))
	L.SetField(dbTable, "ping", L.NewFunction(module.luaPing))
	
	// Query operations
	L.SetField(dbTable, "query", L.NewFunction(module.luaQuery))
	L.SetField(dbTable, "exec", L.NewFunction(module.luaExec))
	L.SetField(dbTable, "prepare", L.NewFunction(module.luaPrepare))
	
	// Transaction operations
	L.SetField(dbTable, "begin", L.NewFunction(module.luaBegin))
	L.SetField(dbTable, "commit", L.NewFunction(module.luaCommit))
	L.SetField(dbTable, "rollback", L.NewFunction(module.luaRollback))
	L.SetField(dbTable, "transaction", L.NewFunction(module.luaTransaction))
	
	// Utility functions
	L.SetField(dbTable, "escape", L.NewFunction(module.luaEscape))
	L.SetField(dbTable, "last_insert_id", L.NewFunction(module.luaLastInsertID))
	L.SetField(dbTable, "rows_affected", L.NewFunction(module.luaRowsAffected))
	
	// Store module reference for cleanup
	ud := L.NewUserData()
	ud.Value = module
	L.SetGlobal("__db_module", ud)
	
	// Register the db table globally
	L.SetGlobal("db", dbTable)
}

// Connection management
func (db *DatabaseModule) luaConnect(L *lua.LState) int {
	driverName := L.CheckString(1)
	dataSourceName := L.CheckString(2)
	connectionName := L.OptString(3, "default")
	
	// Close existing connection with same name
	if existingDB, exists := db.connections[connectionName]; exists {
		existingDB.Close()
	}
	
	sqlDB, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	db.connections[connectionName] = sqlDB
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(connectionName))
	return 2
}

func (db *DatabaseModule) luaClose(L *lua.LState) int {
	connectionName := L.OptString(1, "default")
	
	if sqlDB, exists := db.connections[connectionName]; exists {
		err := sqlDB.Close()
		delete(db.connections, connectionName)
		
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(err.Error()))
			return 2
		}
	}
	
	L.Push(lua.LBool(true))
	return 1
}

func (db *DatabaseModule) luaPing(L *lua.LState) int {
	connectionName := L.OptString(1, "default")
	
	sqlDB, exists := db.connections[connectionName]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("connection not found: " + connectionName))
		return 2
	}
	
	err := sqlDB.Ping()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// Query operations
func (db *DatabaseModule) luaQuery(L *lua.LState) int {
	query := L.CheckString(1)
	connectionName := L.OptString(2, "default")
	
	// Get parameters
	var args []interface{}
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, luaValueToInterface(L.Get(i)))
	}
	
	sqlDB, exists := db.connections[connectionName]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("connection not found: " + connectionName))
		return 2
	}
	
	rows, err := sqlDB.Query(query, args...)
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
	
	// Create result table
	result := L.NewTable()
	rowIndex := 1
	
	for rows.Next() {
		// Create slice to hold row values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		
		// Scan row
		if err := rows.Scan(valuePtrs...); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		
		// Create row table
		row := L.NewTable()
		for i, col := range columns {
			value := values[i]
			row.RawSetString(col, interfaceToLuaValue(L, value))
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

func (db *DatabaseModule) luaExec(L *lua.LState) int {
	query := L.CheckString(1)
	connectionName := L.OptString(2, "default")
	
	// Get parameters
	var args []interface{}
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, luaValueToInterface(L.Get(i)))
	}
	
	sqlDB, exists := db.connections[connectionName]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("connection not found: " + connectionName))
		return 2
	}
	
	result, err := sqlDB.Exec(query, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Get last insert ID and rows affected
	lastInsertID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	
	resultTable := L.NewTable()
	L.SetField(resultTable, "last_insert_id", lua.LNumber(lastInsertID))
	L.SetField(resultTable, "rows_affected", lua.LNumber(rowsAffected))
	
	L.Push(resultTable)
	return 1
}

func (db *DatabaseModule) luaPrepare(L *lua.LState) int {
	query := L.CheckString(1)
	connectionName := L.OptString(2, "default")
	
	sqlDB, exists := db.connections[connectionName]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("connection not found: " + connectionName))
		return 2
	}
	
	stmt, err := sqlDB.Prepare(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Return a prepared statement userdata
	ud := L.NewUserData()
	ud.Value = stmt
	L.Push(ud)
	return 1
}

// Transaction operations
func (db *DatabaseModule) luaBegin(L *lua.LState) int {
	connectionName := L.OptString(1, "default")
	
	sqlDB, exists := db.connections[connectionName]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("connection not found: " + connectionName))
		return 2
	}
	
	tx, err := sqlDB.Begin()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Return transaction userdata
	ud := L.NewUserData()
	ud.Value = tx
	L.Push(ud)
	return 1
}

func (db *DatabaseModule) luaCommit(L *lua.LState) int {
	txUserData := L.CheckUserData(1)
	
	tx, ok := txUserData.Value.(*sql.Tx)
	if !ok {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("invalid transaction object"))
		return 2
	}
	
	err := tx.Commit()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

func (db *DatabaseModule) luaRollback(L *lua.LState) int {
	txUserData := L.CheckUserData(1)
	
	tx, ok := txUserData.Value.(*sql.Tx)
	if !ok {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("invalid transaction object"))
		return 2
	}
	
	err := tx.Rollback()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

func (db *DatabaseModule) luaTransaction(L *lua.LState) int {
	fn := L.CheckFunction(1)
	connectionName := L.OptString(2, "default")
	
	sqlDB, exists := db.connections[connectionName]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("connection not found: " + connectionName))
		return 2
	}
	
	tx, err := sqlDB.Begin()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Execute the function
	L.Push(fn)
	txUserData := L.NewUserData()
	txUserData.Value = tx
	L.Push(txUserData)
	
	err = L.PCall(1, 1, nil)
	if err != nil {
		tx.Rollback()
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Check if function returned false (indicating rollback)
	result := L.Get(-1)
	L.Pop(1)
	
	if lua.LVAsBool(result) == false {
		tx.Rollback()
		L.Push(lua.LBool(false))
		L.Push(lua.LString("transaction rolled back by user function"))
		return 2
	}
	
	err = tx.Commit()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// Utility functions
func (db *DatabaseModule) luaEscape(L *lua.LState) int {
	value := L.CheckString(1)
	// Simple SQL escape - replace single quotes
	escaped := fmt.Sprintf("'%s'", fmt.Sprintf("%q", value))
	L.Push(lua.LString(escaped))
	return 1
}

func (db *DatabaseModule) luaLastInsertID(L *lua.LState) int {
	resultUserData := L.CheckUserData(1)
	
	result, ok := resultUserData.Value.(sql.Result)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString("invalid result object"))
		return 2
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LNumber(id))
	return 1
}

func (db *DatabaseModule) luaRowsAffected(L *lua.LState) int {
	resultUserData := L.CheckUserData(1)
	
	result, ok := resultUserData.Value.(sql.Result)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString("invalid result object"))
		return 2
	}
	
	count, err := result.RowsAffected()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LNumber(count))
	return 1
}

// Helper functions
func luaValueToInterface(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LNilType:
		return nil
	default:
		return lua.LVAsString(v)
	}
}

func interfaceToLuaValue(L *lua.LState, value interface{}) lua.LValue {
	if value == nil {
		return lua.LNil
	}
	
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return lua.LNil
	}
	
	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return lua.LNil
		}
		v = v.Elem()
	}
	
	switch v.Kind() {
	case reflect.Bool:
		return lua.LBool(v.Bool())
	case reflect.String:
		return lua.LString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(v.Uint())
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(v.Float())
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// Handle []byte
			return lua.LString(string(v.Bytes()))
		}
		fallthrough
	default:
		return lua.LString(fmt.Sprintf("%v", value))
	}
}