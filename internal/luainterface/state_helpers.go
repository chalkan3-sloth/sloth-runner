package luainterface

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// Helper methods for StateModule

// setValue stores a value with optional TTL
func (s *StateModule) setValue(key string, luaValue lua.LValue, ttlSeconds int64) error {
	value, valueType := s.serializeLuaValue(luaValue)
	
	var expiresAt *string
	if ttlSeconds > 0 {
		expiry := time.Now().Add(time.Duration(ttlSeconds) * time.Second).Format(time.RFC3339)
		expiresAt = &expiry
	}

	query := `
	INSERT INTO state_data (key, value, type, expires_at, updated_at, version) 
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, 1)
	ON CONFLICT(key) DO UPDATE SET 
		value = excluded.value,
		type = excluded.type,
		expires_at = excluded.expires_at,
		updated_at = CURRENT_TIMESTAMP,
		version = version + 1
	`

	_, err := s.db.Exec(query, key, value, valueType, expiresAt)
	return err
}

// getValue retrieves and deserializes a value
func (s *StateModule) getValue(key string) (lua.LValue, error) {
	var value, valueType string
	var expiresAt *string

	query := `
	SELECT value, type, expires_at FROM state_data 
	WHERE key = ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`

	err := s.db.QueryRow(query, key).Scan(&value, &valueType, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return s.deserializeValue(value, valueType), nil
}

// deleteKey removes a key
func (s *StateModule) deleteKey(key string) error {
	_, err := s.db.Exec("DELETE FROM state_data WHERE key = ?", key)
	return err
}

// keyExists checks if a key exists
func (s *StateModule) keyExists(key string) (bool, error) {
	var count int
	query := `
	SELECT COUNT(*) FROM state_data 
	WHERE key = ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`
	
	err := s.db.QueryRow(query, key).Scan(&count)
	return count > 0, err
}

// getKeys returns keys matching a pattern
func (s *StateModule) getKeys(pattern string) ([]string, error) {
	var keys []string
	
	// Convert shell-style pattern to SQL LIKE pattern
	sqlPattern := strings.ReplaceAll(pattern, "*", "%")
	sqlPattern = strings.ReplaceAll(sqlPattern, "?", "_")
	
	query := `
	SELECT key FROM state_data 
	WHERE key LIKE ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	ORDER BY key
	`
	
	rows, err := s.db.Query(query, sqlPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, rows.Err()
}

// clearKeys removes keys matching a pattern
func (s *StateModule) clearKeys(pattern string) error {
	sqlPattern := strings.ReplaceAll(pattern, "*", "%")
	sqlPattern = strings.ReplaceAll(sqlPattern, "?", "_")
	
	_, err := s.db.Exec("DELETE FROM state_data WHERE key LIKE ?", sqlPattern)
	return err
}

// acquireLock attempts to acquire a distributed lock with timeout
func (s *StateModule) acquireLock(lockName string, timeout time.Duration) (bool, error) {
	owner := fmt.Sprintf("sloth-runner-%d", os.Getpid())
	expiresAt := time.Now().Add(timeout).Format(time.RFC3339)
	
	// Clean expired locks first
	_, err := s.db.Exec("DELETE FROM state_locks WHERE expires_at < CURRENT_TIMESTAMP")
	if err != nil {
		return false, err
	}

	// Try to acquire lock
	query := `
	INSERT INTO state_locks (lock_name, owner, expires_at) 
	VALUES (?, ?, ?)
	ON CONFLICT(lock_name) DO NOTHING
	`
	
	result, err := s.db.Exec(query, lockName, owner, expiresAt)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// tryLock attempts to acquire a lock without waiting
func (s *StateModule) tryLock(lockName string, ttl time.Duration) (bool, error) {
	return s.acquireLock(lockName, ttl)
}

// releaseLock releases a distributed lock
func (s *StateModule) releaseLock(lockName string) error {
	owner := fmt.Sprintf("sloth-runner-%d", os.Getpid())
	
	query := "DELETE FROM state_locks WHERE lock_name = ? AND owner = ?"
	result, err := s.db.Exec(query, lockName, owner)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("lock not held by this process")
	}

	return nil
}

// setTTL sets expiration for an existing key
func (s *StateModule) setTTL(key string, ttlSeconds int64) error {
	var expiresAt *string
	if ttlSeconds > 0 {
		expiry := time.Now().Add(time.Duration(ttlSeconds) * time.Second).Format(time.RFC3339)
		expiresAt = &expiry
	}

	query := "UPDATE state_data SET expires_at = ?, updated_at = CURRENT_TIMESTAMP WHERE key = ?"
	result, err := s.db.Exec(query, expiresAt, key)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("key not found: %s", key)
	}

	return nil
}

// getTTL gets remaining TTL for a key
func (s *StateModule) getTTL(key string) (int64, error) {
	var expiresAt *string
	
	query := "SELECT expires_at FROM state_data WHERE key = ?"
	err := s.db.QueryRow(query, key).Scan(&expiresAt)
	if err == sql.ErrNoRows {
		return -2, nil // Key doesn't exist
	}
	if err != nil {
		return -1, err
	}

	if expiresAt == nil {
		return -1, nil // No expiration set
	}

	expiry, err := time.Parse(time.RFC3339, *expiresAt)
	if err != nil {
		return -1, err
	}

	remaining := time.Until(expiry)
	if remaining <= 0 {
		return 0, nil // Already expired
	}

	return int64(remaining.Seconds()), nil
}

// increment atomically increments a numeric value
func (s *StateModule) increment(key string, delta float64) (float64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var value, valueType string
	query := `
	SELECT value, type FROM state_data 
	WHERE key = ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`
	
	err = tx.QueryRow(query, key).Scan(&value, &valueType)
	
	var currentValue float64
	if err == sql.ErrNoRows {
		currentValue = 0 // Initialize to 0 if key doesn't exist
	} else if err != nil {
		return 0, err
	} else {
		if valueType != "number" {
			return 0, fmt.Errorf("key %s is not a number", key)
		}
		currentValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number format: %v", err)
		}
	}

	newValue := currentValue + delta

	// Update or insert the new value
	updateQuery := `
	INSERT INTO state_data (key, value, type, updated_at, version) 
	VALUES (?, ?, 'number', CURRENT_TIMESTAMP, 1)
	ON CONFLICT(key) DO UPDATE SET 
		value = ?,
		type = 'number',
		updated_at = CURRENT_TIMESTAMP,
		version = version + 1
	`
	
	newValueStr := strconv.FormatFloat(newValue, 'f', -1, 64)
	_, err = tx.Exec(updateQuery, key, newValueStr, newValueStr)
	if err != nil {
		return 0, err
	}

	return newValue, tx.Commit()
}

// appendString appends to a string value
func (s *StateModule) appendString(key string, appendValue string) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var value, valueType string
	query := `
	SELECT value, type FROM state_data 
	WHERE key = ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`
	
	err = tx.QueryRow(query, key).Scan(&value, &valueType)
	
	var currentValue string
	if err == sql.ErrNoRows {
		currentValue = "" // Initialize to empty string if key doesn't exist
	} else if err != nil {
		return 0, err
	} else {
		if valueType != "string" {
			return 0, fmt.Errorf("key %s is not a string", key)
		}
		currentValue = value
	}

	newValue := currentValue + appendValue

	// Update or insert the new value
	updateQuery := `
	INSERT INTO state_data (key, value, type, updated_at, version) 
	VALUES (?, ?, 'string', CURRENT_TIMESTAMP, 1)
	ON CONFLICT(key) DO UPDATE SET 
		value = ?,
		type = 'string',
		updated_at = CURRENT_TIMESTAMP,
		version = version + 1
	`
	
	_, err = tx.Exec(updateQuery, key, newValue, newValue)
	if err != nil {
		return 0, err
	}

	return len(newValue), tx.Commit()
}

// listPush adds an item to a list
func (s *StateModule) listPush(key string, luaValue lua.LValue) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var value, valueType string
	query := `
	SELECT value, type FROM state_data 
	WHERE key = ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`
	
	err = tx.QueryRow(query, key).Scan(&value, &valueType)
	
	var currentList []interface{}
	if err == sql.ErrNoRows {
		currentList = []interface{}{} // Initialize empty list
	} else if err != nil {
		return 0, err
	} else {
		if valueType != "list" {
			return 0, fmt.Errorf("key %s is not a list", key)
		}
		err = json.Unmarshal([]byte(value), &currentList)
		if err != nil {
			return 0, fmt.Errorf("invalid list format: %v", err)
		}
	}

	// Add new item
	newItem := s.luaValueToInterface(luaValue)
	currentList = append(currentList, newItem)

	// Serialize back
	newValue, err := json.Marshal(currentList)
	if err != nil {
		return 0, err
	}

	// Update or insert
	updateQuery := `
	INSERT INTO state_data (key, value, type, updated_at, version) 
	VALUES (?, ?, 'list', CURRENT_TIMESTAMP, 1)
	ON CONFLICT(key) DO UPDATE SET 
		value = ?,
		type = 'list',
		updated_at = CURRENT_TIMESTAMP,
		version = version + 1
	`
	
	_, err = tx.Exec(updateQuery, key, string(newValue), string(newValue))
	if err != nil {
		return 0, err
	}

	return len(currentList), tx.Commit()
}

// listPop removes and returns the last item from a list
func (s *StateModule) listPop(key string) (lua.LValue, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var value, valueType string
	query := `
	SELECT value, type FROM state_data 
	WHERE key = ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`
	
	err = tx.QueryRow(query, key).Scan(&value, &valueType)
	if err == sql.ErrNoRows {
		return nil, nil // Empty list
	}
	if err != nil {
		return nil, err
	}

	if valueType != "list" {
		return nil, fmt.Errorf("key %s is not a list", key)
	}

	var currentList []interface{}
	err = json.Unmarshal([]byte(value), &currentList)
	if err != nil {
		return nil, fmt.Errorf("invalid list format: %v", err)
	}

	if len(currentList) == 0 {
		return nil, nil // Empty list
	}

	// Pop last item
	lastItem := currentList[len(currentList)-1]
	currentList = currentList[:len(currentList)-1]

	// Update list
	newValue, err := json.Marshal(currentList)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE state_data SET value = ?, updated_at = CURRENT_TIMESTAMP, version = version + 1 WHERE key = ?", string(newValue), key)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Convert back to Lua value
	L := lua.NewState()
	defer L.Close()
	return GoValueToLua(L, lastItem), nil
}

// listLength returns the length of a list
func (s *StateModule) listLength(key string) (int, error) {
	var value, valueType string
	query := `
	SELECT value, type FROM state_data 
	WHERE key = ? AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`
	
	err := s.db.QueryRow(query, key).Scan(&value, &valueType)
	if err == sql.ErrNoRows {
		return 0, nil // Empty list
	}
	if err != nil {
		return 0, err
	}

	if valueType != "list" {
		return 0, fmt.Errorf("key %s is not a list", key)
	}

	var currentList []interface{}
	err = json.Unmarshal([]byte(value), &currentList)
	if err != nil {
		return 0, fmt.Errorf("invalid list format: %v", err)
	}

	return len(currentList), nil
}

// compareAndSwap performs atomic compare-and-swap
func (s *StateModule) compareAndSwap(key string, oldValue, newValue lua.LValue) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	currentValue, err := s.getValue(key)
	if err != nil {
		return false, err
	}

	// Compare values (simplified comparison)
	if !s.luaValuesEqual(currentValue, oldValue) {
		return false, nil // Values don't match
	}

	// Swap with new value
	err = s.setValue(key, newValue, 0)
	if err != nil {
		return false, err
	}

	return true, tx.Commit()
}

// getStats returns storage statistics
func (s *StateModule) getStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count total keys
	var totalKeys int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM state_data WHERE expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP").Scan(&totalKeys)
	if err != nil {
		return nil, err
	}
	stats["total_keys"] = totalKeys

	// Count expired keys
	var expiredKeys int64
	err = s.db.QueryRow("SELECT COUNT(*) FROM state_data WHERE expires_at IS NOT NULL AND expires_at <= CURRENT_TIMESTAMP").Scan(&expiredKeys)
	if err != nil {
		return nil, err
	}
	stats["expired_keys"] = expiredKeys

	// Count active locks
	var activeLocks int64
	err = s.db.QueryRow("SELECT COUNT(*) FROM state_locks WHERE expires_at > CURRENT_TIMESTAMP").Scan(&activeLocks)
	if err != nil {
		return nil, err
	}
	stats["active_locks"] = activeLocks

	// Database file size
	if info, err := os.Stat(s.dbPath); err == nil {
		stats["db_size_bytes"] = info.Size()
	}

	stats["db_path"] = s.dbPath

	return stats, nil
}

// Serialization helpers

func (s *StateModule) serializeLuaValue(value lua.LValue) (string, string) {
	switch value.Type() {
	case lua.LTString:
		return lua.LVAsString(value), "string"
	case lua.LTNumber:
		return fmt.Sprintf("%v", lua.LVAsNumber(value)), "number"
	case lua.LTBool:
		return fmt.Sprintf("%v", lua.LVAsBool(value)), "boolean"
	case lua.LTTable:
		// Serialize as JSON
		goValue := LuaToGoValue(nil, value)
		jsonBytes, _ := json.Marshal(goValue)
		
		// Detect if it's an array or object
		if table := value.(*lua.LTable); table.Len() > 0 {
			return string(jsonBytes), "list"
		}
		return string(jsonBytes), "table"
	default:
		return fmt.Sprintf("%v", value), "string"
	}
}

func (s *StateModule) deserializeValue(value, valueType string) lua.LValue {
	L := lua.NewState()
	defer L.Close()
	
	switch valueType {
	case "string":
		return lua.LString(value)
	case "number":
		if num, err := strconv.ParseFloat(value, 64); err == nil {
			return lua.LNumber(num)
		}
		return lua.LString(value)
	case "boolean":
		return lua.LBool(value == "true")
	case "list", "table":
		var goValue interface{}
		if err := json.Unmarshal([]byte(value), &goValue); err == nil {
			return GoValueToLua(L, goValue)
		}
		return lua.LString(value)
	default:
		return lua.LString(value)
	}
}

func (s *StateModule) luaValueToInterface(value lua.LValue) interface{} {
	switch value.Type() {
	case lua.LTString:
		return lua.LVAsString(value)
	case lua.LTNumber:
		return lua.LVAsNumber(value)
	case lua.LTBool:
		return lua.LVAsBool(value)
	case lua.LTTable:
		return LuaToGoValue(nil, value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (s *StateModule) luaValuesEqual(v1, v2 lua.LValue) bool {
	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}
	
	if v1.Type() != v2.Type() {
		return false
	}
	
	switch v1.Type() {
	case lua.LTString:
		return lua.LVAsString(v1) == lua.LVAsString(v2)
	case lua.LTNumber:
		return lua.LVAsNumber(v1) == lua.LVAsNumber(v2)
	case lua.LTBool:
		return lua.LVAsBool(v1) == lua.LVAsBool(v2)
	default:
		// Simplified comparison for complex types
		return fmt.Sprintf("%v", v1) == fmt.Sprintf("%v", v2)
	}
}

// OpenState initializes the state module in Lua
func OpenState(L *lua.LState) {
	L.PreloadModule("state", StateLoader)
	if err := L.DoString(`state = require("state")`); err != nil {
		panic(err)
	}
}