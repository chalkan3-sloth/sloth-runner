package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/state"
)

// StateAdapter provides context-aware methods for AI to use StateManager
type StateAdapter struct {
	sm *state.StateManager
}

// NewStateAdapter creates a new state adapter
func NewStateAdapter(sm *state.StateManager) *StateAdapter {
	return &StateAdapter{sm: sm}
}

// Set stores a key-value pair with TTL (TTL is ignored in this implementation)
func (sa *StateAdapter) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return sa.sm.Set(key, value)
}

// Get retrieves a value by key
func (sa *StateAdapter) Get(ctx context.Context, key string) (string, error) {
	value, err := sa.sm.Get(key)
	if err != nil {
		return "", err
	}
	if str, ok := value.(string); ok {
		return str, nil
	}
	return fmt.Sprintf("%v", value), nil
}

// Keys returns all keys matching a pattern
func (sa *StateAdapter) Keys(ctx context.Context, pattern string) ([]string, error) {
	// Extract prefix from pattern (simple implementation)
	prefix := strings.TrimSuffix(pattern, "*")
	
	result, err := sa.sm.List(prefix)
	if err != nil {
		return nil, err
	}
	
	keys := make([]string, 0, len(result))
	for key := range result {
		// Simple pattern matching
		if pattern == "*" || strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	
	return keys, nil
}

// Delete removes a key
func (sa *StateAdapter) Delete(ctx context.Context, key string) error {
	return sa.sm.Delete(key)
}