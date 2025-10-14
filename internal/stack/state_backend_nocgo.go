//go:build !cgo
// +build !cgo

package stack

import (
	"fmt"
	"time"
)

// StateBackend stub for non-CGO builds
type StateBackend struct{}

// StateSnapshot stub for non-CGO builds
type StateSnapshot struct {
	Version     int                    `json:"version"`
	CreatedAt   time.Time              `json:"created_at"`
	CreatedBy   string                 `json:"created_by"`
	Description string                 `json:"description"`
}

// DriftInfo stub for non-CGO builds
type DriftInfo struct {
	ID               int64                  `json:"id"`
	StackID          string                 `json:"stack_id"`
	ResourceID       string                 `json:"resource_id"`
	DetectedAt       time.Time              `json:"detected_at"`
	ExpectedState    map[string]interface{} `json:"expected_state"`
	ActualState      map[string]interface{} `json:"actual_state"`
	DriftedFields    []string               `json:"drifted_fields"`
	IsDrifted        bool                   `json:"is_drifted"`
	ResolutionStatus string                 `json:"resolution_status"`
}

// StateLock stub for non-CGO builds
type StateLock struct {
	StackID   string    `json:"stack_id"`
	LockID    string    `json:"lock_id"`
	Operation string    `json:"operation"`
	Who       string    `json:"who"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Info      string    `json:"info"`
}

// NewStateBackend returns an error for non-CGO builds
func NewStateBackend(dbPath string) (*StateBackend, error) {
	return nil, fmt.Errorf("StateBackend requires CGO support (SQLite). Please use a CGO-enabled build or use filesystem-based state management")
}

// GetStackManager stub
func (sb *StateBackend) GetStackManager() *StackManager {
	return nil
}

// Close stub
func (sb *StateBackend) Close() error {
	return nil
}

// CreateSnapshot stub
func (sb *StateBackend) CreateSnapshot(stackID, createdBy, description string) (int, error) {
	return 0, fmt.Errorf("state backend not available in non-CGO builds")
}

// GetSnapshot stub
func (sb *StateBackend) GetSnapshot(stackID string, version int) (*StateSnapshot, error) {
	return nil, fmt.Errorf("state backend not available in non-CGO builds")
}

// ListSnapshots stub
func (sb *StateBackend) ListSnapshots(stackID string) ([]StateSnapshot, error) {
	return nil, fmt.Errorf("state backend not available in non-CGO builds")
}

// RollbackToSnapshot stub
func (sb *StateBackend) RollbackToSnapshot(stackID string, version int, performedBy string) error {
	return fmt.Errorf("state backend not available in non-CGO builds")
}

// DetectDrift stub
func (sb *StateBackend) DetectDrift(stackID, resourceID string, expectedState, actualState map[string]interface{}) error {
	return fmt.Errorf("state backend not available in non-CGO builds")
}

// GetDriftInfo stub
func (sb *StateBackend) GetDriftInfo(stackID string) ([]*DriftInfo, error) {
	return nil, fmt.Errorf("state backend not available in non-CGO builds")
}

// LockState stub
func (sb *StateBackend) LockState(stackID, lockID, operation, who string, duration time.Duration) error {
	return fmt.Errorf("state backend not available in non-CGO builds")
}

// UnlockState stub
func (sb *StateBackend) UnlockState(stackID, lockID string) error {
	return fmt.Errorf("state backend not available in non-CGO builds")
}

// GetLockInfo stub
func (sb *StateBackend) GetLockInfo(stackID string) (*StateLock, error) {
	return nil, fmt.Errorf("state backend not available in non-CGO builds")
}

// AddTag stub
func (sb *StateBackend) AddTag(stackID, tag string) error {
	return fmt.Errorf("state backend not available in non-CGO builds")
}

// GetTags stub
func (sb *StateBackend) GetTags(stackID string) ([]string, error) {
	return nil, fmt.Errorf("state backend not available in non-CGO builds")
}

// RemoveTag stub
func (sb *StateBackend) RemoveTag(stackID, tag string) error {
	return fmt.Errorf("state backend not available in non-CGO builds")
}

// AddResourceDependency stub
func (sb *StateBackend) AddResourceDependency(resourceID, dependsOnID, depType string) error {
	return fmt.Errorf("state backend not available in non-CGO builds")
}

// GetResourceDependencies stub
func (sb *StateBackend) GetResourceDependencies(resourceID string) ([]string, error) {
	return nil, fmt.Errorf("state backend not available in non-CGO builds")
}

// GetActivity stub
func (sb *StateBackend) GetActivity(stackID string, limit int) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("state backend not available in non-CGO builds")
}
