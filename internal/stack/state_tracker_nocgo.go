//go:build !cgo
// +build !cgo

package stack

import (
	"fmt"
	"time"
)

// StateTracker stub for non-CGO builds
type StateTracker struct{}

// OperationType stub for non-CGO builds
type OperationType string

const (
	// OperationType constants for non-CGO builds
	OperationTypeCreate  OperationType = "create"
	OperationTypeUpdate  OperationType = "update"
	OperationTypeDelete  OperationType = "delete"
	OperationTypeRead    OperationType = "read"
	OperationTypeApply   OperationType = "apply"
	OperationTypeDestroy    OperationType = "destroy"
	OpBackup                OperationType = "backup"
	OpDeployment            OperationType = "deployment"
	OpWorkflowExecution     OperationType = "workflow_execution"
	OpSchedulerEnable       OperationType = "scheduler_enable"
	OpSchedulerDisable      OperationType = "scheduler_disable"
	OpScheduledExecution    OperationType = "scheduled_execution"
	OpSlothAdd              OperationType = "sloth_add"
	OpSlothUpdate           OperationType = "sloth_update"
	OpSlothDelete           OperationType = "sloth_delete"
	OpHookRegister          OperationType = "hook_register"
	OpHookUpdate            OperationType = "hook_update"
	OpHookDelete            OperationType = "hook_delete"
	OpAgentRegistration     OperationType = "agent_registration"
	OpAgentUpdate           OperationType = "agent_update"
	OpAgentDelete           OperationType = "agent_delete"
	OpAgentStop             OperationType = "agent_stop"
)

// Operation stub for non-CGO builds
type Operation struct {
	ID          string                 `json:"id"`
	Type        OperationType          `json:"type"`
	StackName   string                 `json:"stack_name"`
	ResourceID  string                 `json:"resource_id"`
	Status      string                 `json:"status"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Metadata    map[string]interface{} `json:"metadata"`
	Error       string                 `json:"error,omitempty"`
	PerformedBy string                 `json:"performed_by"`
}

// NewStateTracker returns an error for non-CGO builds
func NewStateTracker(dbPath string) (*StateTracker, error) {
	return nil, fmt.Errorf("StateTracker requires CGO support (SQLite). Please use a CGO-enabled build")
}

// TrackOperation stub
func (st *StateTracker) TrackOperation(op *Operation) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// GetOperationHistory stub
func (st *StateTracker) GetOperationHistory(stackID string, limit int) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("state tracker not available in non-CGO builds")
}

// GetResourceHistory stub
func (st *StateTracker) GetResourceHistory(resourceID string, limit int) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("state tracker not available in non-CGO builds")
}

// Close stub
func (st *StateTracker) Close() error {
	return nil
}

// RecordStateTransition stub
func (st *StateTracker) RecordStateTransition(stackID, resourceID, fromState, toState string, metadata map[string]interface{}) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// GetStateHistory stub
func (st *StateTracker) GetStateHistory(resourceID string, limit int) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("state tracker not available in non-CGO builds")
}

// TrackResourceMetrics stub
func (st *StateTracker) TrackResourceMetrics(resourceID string, metrics map[string]interface{}) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// GetResourceMetrics stub
func (st *StateTracker) GetResourceMetrics(resourceID string, since time.Time) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("state tracker not available in non-CGO builds")
}

// PurgeOldData stub
func (st *StateTracker) PurgeOldData(olderThan time.Duration) (int, error) {
	return 0, fmt.Errorf("state tracker not available in non-CGO builds")
}

// GetOperationStats stub
func (st *StateTracker) GetOperationStats(opType OperationType) (map[string]interface{}, error) {
	return nil, fmt.Errorf("state tracker not available in non-CGO builds")
}

// GetAllOperationStats stub
func (st *StateTracker) GetAllOperationStats() (map[string]interface{}, error) {
	return nil, fmt.Errorf("state tracker not available in non-CGO builds")
}

// SearchOperations stub
func (st *StateTracker) SearchOperations(criteria map[string]interface{}) ([]*Operation, error) {
	return nil, fmt.Errorf("state tracker not available in non-CGO builds")
}

// TrackOperationWithEvents stub
func (st *StateTracker) TrackOperationWithEvents(op *Operation) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// CreateSnapshotWithEvent stub
func (st *StateTracker) CreateSnapshotWithEvent(stackID, createdBy, description string) (int, error) {
	return 0, fmt.Errorf("state tracker not available in non-CGO builds")
}

// RollbackToSnapshotWithEvent stub
func (st *StateTracker) RollbackToSnapshotWithEvent(stackID string, version int, performedBy string) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// DetectDriftWithEvent stub
func (st *StateTracker) DetectDriftWithEvent(stackID, resourceID string, expectedState, actualState map[string]interface{}) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// LockStateWithEvent stub
func (st *StateTracker) LockStateWithEvent(stackID, lockID, operation, who string, duration time.Duration) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// UnlockStateWithEvent stub
func (st *StateTracker) UnlockStateWithEvent(stackID, lockID string) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// ForceUnlockStateWithEvent stub
func (st *StateTracker) ForceUnlockStateWithEvent(stackID string, force bool) error {
	return fmt.Errorf("state tracker not available in non-CGO builds")
}

// ValidateState stub
func (st *StateTracker) ValidateState(stackID string) (bool, []string, error) {
	return false, nil, fmt.Errorf("state tracker not available in non-CGO builds")
}

// GetEventHistory stub
func (st *StateTracker) GetEventHistory(stackID string, eventType EventType, limit int) ([]*StateEvent, error) {
	return nil, fmt.Errorf("state tracker not available in non-CGO builds")
}
