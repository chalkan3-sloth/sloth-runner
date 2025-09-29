//go:build !cgo
// +build !cgo

// Package state provides state management for sloth-runner without SQLite
package state

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// StateManager manages in-memory state when SQLite is not available
type StateManager struct {
	data  map[string]interface{}
	locks map[string]*lockInfo
	mu    sync.RWMutex
	path  string
}

// lockInfo represents a lock in memory
type lockInfo struct {
	holder    string
	expiresAt time.Time
}

// NewStateManager creates a new in-memory state manager
func NewStateManager(dbPath string) (*StateManager, error) {
	return &StateManager{
		data:  make(map[string]interface{}),
		locks: make(map[string]*lockInfo),
		path:  dbPath,
	}, nil
}

// Close closes the state manager
func (sm *StateManager) Close() error {
	return nil
}

// Set sets a key-value pair
func (sm *StateManager) Set(key string, value interface{}) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
	return nil
}

// Get retrieves a value by key
func (sm *StateManager) Get(key string) (interface{}, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	value, exists := sm.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

// Delete removes a key
func (sm *StateManager) Delete(key string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
	return nil
}

// Exists checks if a key exists
func (sm *StateManager) Exists(key string) (bool, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, exists := sm.data[key]
	return exists, nil
}

// List returns all keys (signature compatible with SQLite version)
func (sm *StateManager) List(prefix string) (map[string]string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	result := make(map[string]string)
	for key, value := range sm.data {
		if prefix == "" || strings.HasPrefix(key, prefix) {
			// Convert value to string
			valueStr := fmt.Sprintf("%v", value)
			result[key] = valueStr
		}
	}
	return result, nil
}

// Clear removes all keys
func (sm *StateManager) Clear() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data = make(map[string]interface{})
	return nil
}

// Stats returns state statistics
func (sm *StateManager) Stats() (StateStats, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	return StateStats{
		TotalKeys:    len(sm.data),
		TotalSize:    0, // Not easily calculable for in-memory
		LastModified: time.Now().Unix(),
		Backend:      "memory",
	}, nil
}

// Increment increments a numeric value
func (sm *StateManager) Increment(key string, delta int64) (int64, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	current := int64(0)
	if val, exists := sm.data[key]; exists {
		if num, ok := val.(int64); ok {
			current = num
		}
	}
	
	newValue := current + delta
	sm.data[key] = newValue
	return newValue, nil
}

// SetWithTTL sets a key with TTL (no-op in memory implementation)
func (sm *StateManager) SetWithTTL(key string, value interface{}, ttlSeconds int) error {
	// For in-memory implementation, we'll just set the value without TTL
	// TTL functionality would require a background goroutine to clean up expired keys
	return sm.Set(key, value)
}

// StateStats represents state statistics
type StateStats struct {
	TotalKeys    int    `json:"total_keys"`
	TotalSize    int64  `json:"total_size"`
	LastModified int64  `json:"last_modified"`
	Backend      string `json:"backend"`
}

// StateMetadata contains metadata about a state entry
type StateMetadata struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Lock holds a distributed lock (in-memory implementation)
func (sm *StateManager) Lock(name string, holder string, timeout time.Duration) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	// Clean up expired locks
	sm.cleanupExpiredLocks()
	
	// Check if lock already exists
	if lock, exists := sm.locks[name]; exists {
		return fmt.Errorf("lock '%s' is already held by '%s'", name, lock.holder)
	}
	
	// Create new lock
	sm.locks[name] = &lockInfo{
		holder:    holder,
		expiresAt: time.Now().Add(timeout),
	}
	
	return nil
}

// Unlock releases a distributed lock
func (sm *StateManager) Unlock(name string, holder string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	lock, exists := sm.locks[name]
	if !exists {
		return fmt.Errorf("lock '%s' not found", name)
	}
	
	if lock.holder != holder {
		return fmt.Errorf("lock '%s' not held by '%s'", name, holder)
	}
	
	delete(sm.locks, name)
	return nil
}

// IsLocked checks if a lock is currently held
func (sm *StateManager) IsLocked(name string) (bool, string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	// Clean up expired locks
	sm.cleanupExpiredLocks()
	
	lock, exists := sm.locks[name]
	if !exists {
		return false, "", nil
	}
	
	return true, lock.holder, nil
}

// cleanupExpiredLocks removes expired locks (must be called with lock held)
func (sm *StateManager) cleanupExpiredLocks() {
	now := time.Now()
	for name, lock := range sm.locks {
		if now.After(lock.expiresAt) {
			delete(sm.locks, name)
		}
	}
}

// GetMetadata returns metadata about a key (simplified for memory implementation)
func (sm *StateManager) GetMetadata(key string) (*StateMetadata, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	value, exists := sm.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	
	// Convert value to string for metadata
	valueStr := fmt.Sprintf("%v", value)
	now := time.Now()
	
	return &StateMetadata{
		Key:       key,
		Value:     valueStr,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
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