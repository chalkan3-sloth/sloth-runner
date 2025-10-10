package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// GlobalStateTracker is the singleton instance of the state tracker
var (
	globalStateTracker *StateTrackerService
	trackerMutex       sync.Mutex
	trackerInitialized bool
)

// StateTrackerService wraps the state tracker for use by commands
type StateTrackerService struct {
	tracker *stack.StateTracker
}

// GetGlobalStateTracker returns the singleton state tracker instance
func GetGlobalStateTracker() (*StateTrackerService, error) {
	trackerMutex.Lock()
	defer trackerMutex.Unlock()

	if globalStateTracker == nil {
		tracker, err := stack.NewStateTracker("")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize state tracker: %w", err)
		}

		globalStateTracker = &StateTrackerService{
			tracker: tracker,
		}
		trackerInitialized = true
	}

	return globalStateTracker, nil
}

// CloseGlobalStateTracker closes the global state tracker
func CloseGlobalStateTracker() error {
	trackerMutex.Lock()
	defer trackerMutex.Unlock()

	if globalStateTracker != nil {
		err := globalStateTracker.tracker.Close()
		globalStateTracker = nil
		trackerInitialized = false
		return err
	}

	return nil
}

// IsTrackerInitialized returns whether the tracker has been initialized
func IsTrackerInitialized() bool {
	trackerMutex.Lock()
	defer trackerMutex.Unlock()
	return trackerInitialized
}

// TrackAgentOperation tracks an agent operation
func (sts *StateTrackerService) TrackAgentOperation(
	opType stack.OperationType,
	agentName string,
	status string,
	metadata map[string]interface{},
	performedBy string,
) error {
	op := &stack.Operation{
		Type:        opType,
		StackName:   "agent-operations",
		ResourceID:  agentName,
		Status:      status,
		StartedAt:   time.Now(),
		Metadata:    metadata,
		PerformedBy: performedBy,
	}

	if status == "completed" || status == "failed" {
		now := time.Now()
		op.CompletedAt = &now
		op.Duration = now.Sub(op.StartedAt)
	}

	return sts.tracker.TrackOperation(op)
}

// TrackSchedulerOperation tracks a scheduler operation
func (sts *StateTrackerService) TrackSchedulerOperation(
	opType stack.OperationType,
	workflowName string,
	schedule string,
	status string,
	duration time.Duration,
	errorMsg string,
	performedBy string,
) error {
	op := &stack.Operation{
		Type:       opType,
		StackName:  "scheduler-operations",
		ResourceID: workflowName,
		Status:     status,
		StartedAt:  time.Now().Add(-duration),
		Duration:   duration,
		Metadata: map[string]interface{}{
			"schedule": schedule,
		},
		Error:       errorMsg,
		PerformedBy: performedBy,
	}

	if status == "completed" || status == "failed" {
		now := time.Now()
		op.CompletedAt = &now
	}

	return sts.tracker.TrackOperation(op)
}

// TrackSecretOperation tracks a secret operation
func (sts *StateTrackerService) TrackSecretOperation(
	opType stack.OperationType,
	secretKey string,
	stackID string,
	status string,
	performedBy string,
) error {
	op := &stack.Operation{
		Type:       opType,
		StackName:  "secret-operations",
		ResourceID: secretKey,
		Status:     status,
		StartedAt:  time.Now(),
		Metadata: map[string]interface{}{
			"stack_id": stackID,
		},
		PerformedBy: performedBy,
	}

	now := time.Now()
	op.CompletedAt = &now
	op.Duration = time.Since(op.StartedAt)

	return sts.tracker.TrackOperation(op)
}

// TrackHookOperation tracks a hook operation
func (sts *StateTrackerService) TrackHookOperation(
	opType stack.OperationType,
	hookName string,
	hookType string,
	status string,
	performedBy string,
) error {
	op := &stack.Operation{
		Type:       opType,
		StackName:  "hook-operations",
		ResourceID: hookName,
		Status:     status,
		StartedAt:  time.Now(),
		Metadata: map[string]interface{}{
			"hook_type": hookType,
		},
		PerformedBy: performedBy,
	}

	now := time.Now()
	op.CompletedAt = &now
	op.Duration = time.Since(op.StartedAt)

	return sts.tracker.TrackOperation(op)
}

// TrackSlothOperation tracks a sloth file operation
func (sts *StateTrackerService) TrackSlothOperation(
	opType stack.OperationType,
	slothName string,
	filePath string,
	status string,
	performedBy string,
) error {
	op := &stack.Operation{
		Type:       opType,
		StackName:  "sloth-operations",
		ResourceID: slothName,
		Status:     status,
		StartedAt:  time.Now(),
		Metadata: map[string]interface{}{
			"file_path": filePath,
		},
		PerformedBy: performedBy,
	}

	now := time.Now()
	op.CompletedAt = &now
	op.Duration = time.Since(op.StartedAt)

	return sts.tracker.TrackOperation(op)
}

// TrackBackupOperation tracks a backup operation
func (sts *StateTrackerService) TrackBackupOperation(
	backupID string,
	backupPath string,
	size int64,
	status string,
	performedBy string,
) error {
	op := &stack.Operation{
		Type:       stack.OpBackup,
		StackName:  "sysadmin-operations",
		ResourceID: backupID,
		Status:     status,
		StartedAt:  time.Now(),
		Metadata: map[string]interface{}{
			"backup_path": backupPath,
			"size_bytes":  size,
		},
		PerformedBy: performedBy,
	}

	now := time.Now()
	op.CompletedAt = &now
	op.Duration = time.Since(op.StartedAt)

	return sts.tracker.TrackOperation(op)
}

// TrackDeploymentOperation tracks a deployment operation
func (sts *StateTrackerService) TrackDeploymentOperation(
	version string,
	agents []string,
	strategy string,
	success bool,
	duration time.Duration,
	performedBy string,
) error {
	status := "completed"
	if !success {
		status = "failed"
	}

	op := &stack.Operation{
		Type:       stack.OpDeployment,
		StackName:  "sysadmin-operations",
		ResourceID: version,
		Status:     status,
		StartedAt:  time.Now().Add(-duration),
		Duration:   duration,
		Metadata: map[string]interface{}{
			"agents":   agents,
			"strategy": strategy,
		},
		PerformedBy: performedBy,
	}

	now := time.Now()
	op.CompletedAt = &now

	return sts.tracker.TrackOperation(op)
}

// GetOperationHistory retrieves operation history
func (sts *StateTrackerService) GetOperationHistory(opType stack.OperationType, limit int) ([]*stack.Resource, error) {
	return sts.tracker.GetOperationHistory(opType, limit)
}

// GetOperationStats retrieves operation statistics
func (sts *StateTrackerService) GetOperationStats(opType stack.OperationType) (map[string]interface{}, error) {
	return sts.tracker.GetOperationStats(opType)
}

// GetAllOperationStats retrieves statistics for all operations
func (sts *StateTrackerService) GetAllOperationStats() (map[stack.OperationType]map[string]interface{}, error) {
	return sts.tracker.GetAllOperationStats()
}

// SearchOperations searches for operations matching criteria
func (sts *StateTrackerService) SearchOperations(criteria map[string]interface{}) ([]*stack.Resource, error) {
	return sts.tracker.SearchOperations(criteria)
}

// TrackOperationWithEvents tracks an operation and emits events
func (sts *StateTrackerService) TrackOperationWithEvents(op *stack.Operation) error {
	return sts.tracker.TrackOperationWithEvents(op)
}

// CreateSnapshotWithEvent creates a snapshot and emits an event
func (sts *StateTrackerService) CreateSnapshotWithEvent(stackID, creator, description string) (int, error) {
	return sts.tracker.CreateSnapshotWithEvent(stackID, creator, description)
}

// RollbackToSnapshotWithEvent rolls back to a snapshot and emits an event
func (sts *StateTrackerService) RollbackToSnapshotWithEvent(stackID string, version int, performer string) error {
	return sts.tracker.RollbackToSnapshotWithEvent(stackID, version, performer)
}

// DetectDriftWithEvent detects drift and emits an event
func (sts *StateTrackerService) DetectDriftWithEvent(stackID string) (bool, []string, error) {
	return sts.tracker.DetectDriftWithEvent(stackID)
}

// LockStateWithEvent locks state and emits an event
func (sts *StateTrackerService) LockStateWithEvent(stackID, lockedBy, reason string) error {
	return sts.tracker.LockStateWithEvent(stackID, lockedBy, reason)
}

// UnlockStateWithEvent unlocks state and emits an event
func (sts *StateTrackerService) UnlockStateWithEvent(stackID, unlockedBy string) error {
	return sts.tracker.UnlockStateWithEvent(stackID, unlockedBy)
}

// ForceUnlockStateWithEvent forcefully unlocks state (for admin use)
func (sts *StateTrackerService) ForceUnlockStateWithEvent(stackID, unlockedBy string) error {
	return sts.tracker.ForceUnlockStateWithEvent(stackID, unlockedBy)
}

// ValidateState validates a stack state
func (sts *StateTrackerService) ValidateState(stackID string) (bool, []string, error) {
	return sts.tracker.ValidateState(stackID)
}

// GetEventHistory returns recent events
func (sts *StateTrackerService) GetEventHistory(limit int) []*stack.StateEvent {
	return sts.tracker.GetEventHistory(limit)
}

// GetEventsByType returns events filtered by type
func (sts *StateTrackerService) GetEventsByType(eventType stack.EventType, limit int) []*stack.StateEvent {
	return sts.tracker.GetEventsByType(eventType, limit)
}

// GetEventsByStack returns events for a specific stack
func (sts *StateTrackerService) GetEventsByStack(stackID string, limit int) []*stack.StateEvent {
	return sts.tracker.GetEventsByStack(stackID, limit)
}

// GetEventStats returns event statistics
func (sts *StateTrackerService) GetEventStats() map[string]interface{} {
	return sts.tracker.GetEventStats()
}

// GetEventBus returns the event bus for subscriptions
func (sts *StateTrackerService) GetEventBus() *stack.EventBus {
	return sts.tracker.GetEventBus()
}

// GetBackend returns the underlying state backend
func (sts *StateTrackerService) GetBackend() *stack.StateBackend {
	return sts.tracker.GetBackend()
}
