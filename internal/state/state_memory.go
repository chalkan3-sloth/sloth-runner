//go:build !cgo
// +build !cgo

// Package state provides state management for sloth-runner without SQLite
package state

import (
	"fmt"
	"sync"
	"time"
)

// StateManager manages in-memory state when SQLite is not available
type StateManager struct {
	data map[string]interface{}
	mu   sync.RWMutex
	path string
}

// NewStateManager creates a new in-memory state manager
func NewStateManager(dbPath string) (*StateManager, error) {
	return &StateManager{
		data: make(map[string]interface{}),
		path: dbPath,
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

// List returns all keys
func (sm *StateManager) List() ([]string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	keys := make([]string, 0, len(sm.data))
	for key := range sm.data {
		keys = append(keys, key)
	}
	return keys, nil
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