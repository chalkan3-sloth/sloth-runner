package luainterface

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
	_ "github.com/mattn/go-sqlite3"
)

// StateModule provides state management and persistence functionality
type StateModule struct {
	db       *sql.DB
	locks    map[string]*sync.RWMutex
	locksMux *sync.RWMutex
	dbPath   string
}

var (
	globalStateModule *StateModule
	stateModuleOnce   sync.Once
)

// GetGlobalStateModule returns the singleton state module instance
func GetGlobalStateModule() *StateModule {
	stateModuleOnce.Do(func() {
		globalStateModule = NewStateModule("")
	})
	return globalStateModule
}

// NewStateModule creates a new StateModule with SQLite backend
func NewStateModule(dbPath string) *StateModule {
	if dbPath == "" {
		homeDir, _ := os.UserHomeDir()
		dbPath = filepath.Join(homeDir, ".sloth-runner", "state.db")
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("Failed to create state directory", "error", err, "path", dir)
		return nil
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_timeout=5000")
	if err != nil {
		slog.Error("Failed to open state database", "error", err, "path", dbPath)
		return nil
	}

	module := &StateModule{
		db:       db,
		locks:    make(map[string]*sync.RWMutex),
		locksMux: &sync.RWMutex{},
		dbPath:   dbPath,
	}

	if err := module.initDB(); err != nil {
		slog.Error("Failed to initialize state database", "error", err)
		return nil
	}

	slog.Info("State module initialized", "database", dbPath)
	return module
}

// initDB creates the necessary tables
func (s *StateModule) initDB() error {
	schema := `
	CREATE TABLE IF NOT EXISTS state_data (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		type TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NULL,
		version INTEGER DEFAULT 1
	);

	CREATE INDEX IF NOT EXISTS idx_state_expires_at ON state_data(expires_at);
	CREATE INDEX IF NOT EXISTS idx_state_updated_at ON state_data(updated_at);

	CREATE TABLE IF NOT EXISTS state_locks (
		lock_name TEXT PRIMARY KEY,
		owner TEXT NOT NULL,
		acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_locks_expires_at ON state_locks(expires_at);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Clean expired data and locks on startup
	go s.cleanupExpired()
	
	return nil
}

// cleanupExpired removes expired entries and locks
func (s *StateModule) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().Format(time.RFC3339)
			
			// Clean expired state data
			if _, err := s.db.Exec("DELETE FROM state_data WHERE expires_at IS NOT NULL AND expires_at < ?", now); err != nil {
				slog.Error("Failed to clean expired state data", "error", err)
			}
			
			// Clean expired locks
			if _, err := s.db.Exec("DELETE FROM state_locks WHERE expires_at < ?", now); err != nil {
				slog.Error("Failed to clean expired locks", "error", err)
			}
		}
	}
}

// Loader is the module loader function for Lua
func (s *StateModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set":           s.luaStateSet,
		"get":           s.luaStateGet,
		"delete":        s.luaStateDelete,
		"exists":        s.luaStateExists,
		"keys":          s.luaStateKeys,
		"clear":         s.luaStateClear,
		"lock":          s.luaStateLock,
		"try_lock":      s.luaStateTryLock,
		"unlock":        s.luaStateUnlock,
		"with_lock":     s.luaStateWithLock,
		"set_ttl":       s.luaStateSetTTL,
		"get_ttl":       s.luaStateGetTTL,
		"increment":     s.luaStateIncrement,
		"decrement":     s.luaStateDecrement,
		"append":        s.luaStateAppend,
		"list_push":     s.luaStateListPush,
		"list_pop":      s.luaStateListPop,
		"list_length":   s.luaStateListLength,
		"compare_swap":  s.luaStateCompareSwap,
		"stats":         s.luaStateStats,
	})
	L.Push(mod)
	return 1
}

// StateLoader is the global loader function
func StateLoader(L *lua.LState) int {
	return GetGlobalStateModule().Loader(L)
}

// luaStateSet sets a value in the state store
func (s *StateModule) luaStateSet(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckAny(2)
	
	var ttlSeconds int64 = 0
	if L.GetTop() >= 3 {
		ttlSeconds = int64(L.CheckNumber(3))
	}

	err := s.setValue(key, value, ttlSeconds)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

// luaStateGet gets a value from the state store
func (s *StateModule) luaStateGet(L *lua.LState) int {
	key := L.CheckString(1)
	defaultValue := lua.LNil
	
	if L.GetTop() >= 2 {
		defaultValue = L.CheckAny(2)
	}

	value, err := s.getValue(key)
	if err != nil {
		L.Push(defaultValue)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if value == nil {
		L.Push(defaultValue)
		return 1
	}

	L.Push(value)
	return 1
}

// luaStateDelete deletes a key from the state store
func (s *StateModule) luaStateDelete(L *lua.LState) int {
	key := L.CheckString(1)
	
	err := s.deleteKey(key)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

// luaStateExists checks if a key exists
func (s *StateModule) luaStateExists(L *lua.LState) int {
	key := L.CheckString(1)
	
	exists, err := s.keyExists(key)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(exists))
	return 1
}

// luaStateKeys returns all keys matching a pattern
func (s *StateModule) luaStateKeys(L *lua.LState) int {
	pattern := "*"
	if L.GetTop() >= 1 {
		pattern = L.CheckString(1)
	}

	keys, err := s.getKeys(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	table := L.NewTable()
	for i, key := range keys {
		table.RawSetInt(i+1, lua.LString(key))
	}
	
	L.Push(table)
	return 1
}

// luaStateClear clears all state data
func (s *StateModule) luaStateClear(L *lua.LState) int {
	pattern := "*"
	if L.GetTop() >= 1 {
		pattern = L.CheckString(1)
	}

	err := s.clearKeys(pattern)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

// luaStateLock acquires a distributed lock
func (s *StateModule) luaStateLock(L *lua.LState) int {
	lockName := L.CheckString(1)
	timeoutSeconds := 30.0
	
	if L.GetTop() >= 2 {
		timeoutSeconds = float64(L.CheckNumber(2))
	}

	success, err := s.acquireLock(lockName, time.Duration(timeoutSeconds)*time.Second)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(success))
	return 1
}

// luaStateTryLock attempts to acquire a lock without blocking
func (s *StateModule) luaStateTryLock(L *lua.LState) int {
	lockName := L.CheckString(1)
	ttlSeconds := 30.0
	
	if L.GetTop() >= 2 {
		ttlSeconds = float64(L.CheckNumber(2))
	}

	success, err := s.tryLock(lockName, time.Duration(ttlSeconds)*time.Second)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(success))
	return 1
}

// luaStateUnlock releases a distributed lock
func (s *StateModule) luaStateUnlock(L *lua.LState) int {
	lockName := L.CheckString(1)
	
	err := s.releaseLock(lockName)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

// luaStateWithLock executes a function with a lock held
func (s *StateModule) luaStateWithLock(L *lua.LState) int {
	lockName := L.CheckString(1)
	fn := L.CheckFunction(2)
	timeoutSeconds := 30.0
	
	if L.GetTop() >= 3 {
		timeoutSeconds = float64(L.CheckNumber(3))
	}

	// Acquire lock
	success, err := s.acquireLock(lockName, time.Duration(timeoutSeconds)*time.Second)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to acquire lock: %v", err)))
		return 2
	}

	if !success {
		L.Push(lua.LNil)
		L.Push(lua.LString("Failed to acquire lock: timeout"))
		return 2
	}

	// Execute function with lock held
	defer func() {
		if unlockErr := s.releaseLock(lockName); unlockErr != nil {
			slog.Error("Failed to release lock", "lock", lockName, "error", unlockErr)
		}
	}()

	// Call the Lua function
	L.Push(fn)
	err = L.PCall(0, lua.MultRet, nil)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Function execution failed: %v", err)))
		return 2
	}

	// Return whatever the function returned
	return L.GetTop()
}

// luaStateSetTTL sets TTL for an existing key
func (s *StateModule) luaStateSetTTL(L *lua.LState) int {
	key := L.CheckString(1)
	ttlSeconds := int64(L.CheckNumber(2))

	err := s.setTTL(key, ttlSeconds)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

// luaStateGetTTL gets TTL for a key
func (s *StateModule) luaStateGetTTL(L *lua.LState) int {
	key := L.CheckString(1)

	ttl, err := s.getTTL(key)
	if err != nil {
		L.Push(lua.LNumber(-1))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(ttl))
	return 1
}

// luaStateIncrement atomically increments a numeric value
func (s *StateModule) luaStateIncrement(L *lua.LState) int {
	key := L.CheckString(1)
	delta := 1.0
	
	if L.GetTop() >= 2 {
		delta = float64(L.CheckNumber(2))
	}

	result, err := s.increment(key, delta)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(result))
	return 1
}

// luaStateDecrement atomically decrements a numeric value
func (s *StateModule) luaStateDecrement(L *lua.LState) int {
	key := L.CheckString(1)
	delta := 1.0
	
	if L.GetTop() >= 2 {
		delta = float64(L.CheckNumber(2))
	}

	result, err := s.increment(key, -delta)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(result))
	return 1
}

// luaStateAppend appends to a string value
func (s *StateModule) luaStateAppend(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)

	newLength, err := s.appendString(key, value)
	if err != nil {
		L.Push(lua.LNumber(-1))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(newLength))
	return 1
}

// luaStateListPush pushes an item to a list
func (s *StateModule) luaStateListPush(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckAny(2)

	length, err := s.listPush(key, value)
	if err != nil {
		L.Push(lua.LNumber(-1))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(length))
	return 1
}

// luaStateListPop pops an item from a list
func (s *StateModule) luaStateListPop(L *lua.LState) int {
	key := L.CheckString(1)

	value, err := s.listPop(key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if value == nil {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(value)
	return 1
}

// luaStateListLength gets the length of a list
func (s *StateModule) luaStateListLength(L *lua.LState) int {
	key := L.CheckString(1)

	length, err := s.listLength(key)
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(length))
	return 1
}

// luaStateCompareSwap performs atomic compare-and-swap
func (s *StateModule) luaStateCompareSwap(L *lua.LState) int {
	key := L.CheckString(1)
	oldValue := L.CheckAny(2)
	newValue := L.CheckAny(3)

	success, err := s.compareAndSwap(key, oldValue, newValue)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(success))
	return 1
}

// luaStateStats returns state storage statistics
func (s *StateModule) luaStateStats(L *lua.LState) int {
	stats, err := s.getStats()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	table := L.NewTable()
	for key, value := range stats {
		switch v := value.(type) {
		case string:
			table.RawSetString(key, lua.LString(v))
		case int64:
			table.RawSetString(key, lua.LNumber(v))
		case float64:
			table.RawSetString(key, lua.LNumber(v))
		}
	}

	L.Push(table)
	return 1
}

// Helper methods implementation will be in the next part due to length constraints...