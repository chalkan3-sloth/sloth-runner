//go:build cgo
// +build cgo

// Package state provides state management and persistence for sloth-runner
package state

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// StateManager manages persistent state
type StateManager struct {
	db   *sql.DB
	mu   sync.RWMutex
	path string
}

// NewStateManager creates a new state manager
func NewStateManager(dbPath string) (*StateManager, error) {
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dbPath = filepath.Join(homeDir, ".sloth-runner", "state.db")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sm := &StateManager{
		db:   db,
		path: dbPath,
	}

	if err := sm.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return sm, nil
}

// initSchema initializes the database schema
func (sm *StateManager) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS state (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS locks (
		name TEXT PRIMARY KEY,
		holder TEXT NOT NULL,
		acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL
	);

	CREATE TRIGGER IF NOT EXISTS update_state_timestamp 
	AFTER UPDATE ON state
	BEGIN
		UPDATE state SET updated_at = CURRENT_TIMESTAMP WHERE key = NEW.key;
	END;
	`

	_, err := sm.db.Exec(schema)
	return err
}

// Set stores a key-value pair
func (sm *StateManager) Set(key, value string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, err := sm.db.Exec(`
		INSERT INTO state (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP
	`, key, value)
	
	return err
}

// Get retrieves a value by key
func (sm *StateManager) Get(key string) (interface{}, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var value string
	err := sm.db.QueryRow("SELECT value FROM state WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return value, err
}

// Delete removes a key-value pair
func (sm *StateManager) Delete(key string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, err := sm.db.Exec("DELETE FROM state WHERE key = ?", key)
	return err
}

// List returns all keys with an optional prefix filter
func (sm *StateManager) List(prefix string) (map[string]string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var query string
	var args []interface{}

	if prefix != "" {
		query = "SELECT key, value FROM state WHERE key LIKE ? ORDER BY key"
		args = []interface{}{prefix + "%"}
	} else {
		query = "SELECT key, value FROM state ORDER BY key"
	}

	rows, err := sm.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}

	return result, rows.Err()
}

// Lock acquires a named lock with timeout
func (sm *StateManager) Lock(name string, holder string, timeout time.Duration) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Clean up expired locks first
	_, err := sm.db.Exec("DELETE FROM locks WHERE expires_at < datetime('now')")
	if err != nil {
		return fmt.Errorf("failed to cleanup expired locks: %w", err)
	}

	expiresAt := time.Now().Add(timeout)

	// Try to acquire lock
	_, err = sm.db.Exec(`
		INSERT INTO locks (name, holder, expires_at) VALUES (?, ?, ?)
	`, name, holder, expiresAt.Format(time.RFC3339))

	if err != nil {
		// Lock already exists or other error
		var existingHolder string
		err2 := sm.db.QueryRow("SELECT holder FROM locks WHERE name = ?", name).Scan(&existingHolder)
		if err2 == nil {
			return fmt.Errorf("lock '%s' is already held by '%s'", name, existingHolder)
		}
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	return nil
}

// Unlock releases a named lock
func (sm *StateManager) Unlock(name string, holder string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	result, err := sm.db.Exec("DELETE FROM locks WHERE name = ? AND holder = ?", name, holder)
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("lock '%s' not found or not held by '%s'", name, holder)
	}

	return nil
}

// IsLocked checks if a lock is currently held
func (sm *StateManager) IsLocked(name string) (bool, string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Clean up expired locks first
	_, err := sm.db.Exec("DELETE FROM locks WHERE expires_at < datetime('now')")
	if err != nil {
		return false, "", fmt.Errorf("failed to cleanup expired locks: %w", err)
	}

	var holder string
	err = sm.db.QueryRow("SELECT holder FROM locks WHERE name = ?", name).Scan(&holder)
	if err == sql.ErrNoRows {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return true, holder, nil
}

// Close closes the database connection
func (sm *StateManager) Close() error {
	if sm.db != nil {
		return sm.db.Close()
	}
	return nil
}

// WithLock executes a function while holding a lock
func (sm *StateManager) WithLock(name string, holder string, timeout time.Duration, fn func() error) error {
	if err := sm.Lock(name, holder, timeout); err != nil {
		return err
	}
	
	defer func() {
		if unlockErr := sm.Unlock(name, holder); unlockErr != nil {
			// Log the error but don't override the original error
			fmt.Printf("Warning: failed to unlock '%s': %v\n", name, unlockErr)
		}
	}()

	return fn()
}

// GetMetadata returns metadata about a key
func (sm *StateManager) GetMetadata(key string) (*StateMetadata, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var metadata StateMetadata
	err := sm.db.QueryRow(`
		SELECT key, value, created_at, updated_at 
		FROM state WHERE key = ?
	`, key).Scan(&metadata.Key, &metadata.Value, &metadata.CreatedAt, &metadata.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	
	return &metadata, err
}

// StateMetadata contains metadata about a state entry
type StateMetadata struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StateStats represents state statistics
type StateStats struct {
	TotalKeys    int    `json:"total_keys"`
	TotalSize    int64  `json:"total_size"`
	LastModified int64  `json:"last_modified"`
	Backend      string `json:"backend"`
}

// Stats returns state statistics
func (sm *StateManager) Stats() (StateStats, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var totalKeys int
	err := sm.db.QueryRow("SELECT COUNT(*) FROM state").Scan(&totalKeys)
	if err != nil {
		return StateStats{}, fmt.Errorf("failed to count keys: %w", err)
	}

	// Get approximate database size
	var totalSize int64
	err = sm.db.QueryRow("PRAGMA page_count").Scan(&totalSize)
	if err == nil {
		var pageSize int64
		sm.db.QueryRow("PRAGMA page_size").Scan(&pageSize)
		totalSize *= pageSize
	}

	return StateStats{
		TotalKeys:    totalKeys,
		TotalSize:    totalSize,
		LastModified: time.Now().Unix(),
		Backend:      "sqlite",
	}, nil
}

// Exists checks if a key exists
func (sm *StateManager) Exists(key string) (bool, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var count int
	err := sm.db.QueryRow("SELECT COUNT(*) FROM state WHERE key = ?", key).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Clear removes all state entries
func (sm *StateManager) Clear() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, err := sm.db.Exec("DELETE FROM state")
	return err
}

// Increment increments a numeric value
func (sm *StateManager) Increment(key string, delta int64) (int64, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var currentValue string
	err := sm.db.QueryRow("SELECT value FROM state WHERE key = ?", key).Scan(&currentValue)
	
	current := int64(0)
	if err == nil {
		// Parse existing value
		if parsed, parseErr := strconv.ParseInt(currentValue, 10, 64); parseErr == nil {
			current = parsed
		}
	}

	newValue := current + delta
	newValueStr := strconv.FormatInt(newValue, 10)

	// Use UPSERT
	_, err = sm.db.Exec(`
		INSERT INTO state (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET 
			value = excluded.value,
			updated_at = CURRENT_TIMESTAMP
	`, key, newValueStr)
	
	if err != nil {
		return 0, err
	}
	
	return newValue, nil
}

// SetWithTTL sets a key with TTL (TTL not implemented in this version)
func (sm *StateManager) SetWithTTL(key string, value interface{}, ttlSeconds int) error {
	// For simplicity, just ignore TTL and set the value
	valueStr := fmt.Sprintf("%v", value)
	return sm.Set(key, valueStr)
}