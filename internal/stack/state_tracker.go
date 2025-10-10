//go:build cgo
// +build cgo

package stack

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// StateTracker is the central tracking system for all operations
// This is the CORE of the application - all operations flow through here
type StateTracker struct {
	backend  *StateBackend
	eventBus *EventBus
}

// NewStateTracker creates a new state tracker (singleton pattern)
func NewStateTracker(dbPath string) (*StateTracker, error) {
	backend, err := NewStateBackend(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize state backend: %w", err)
	}

	st := &StateTracker{
		backend:  backend,
		eventBus: NewEventBus(),
	}

	// Register default event handlers
	st.registerDefaultHandlers()

	return st, nil
}

// GetEventBus returns the event bus for external subscriptions
func (st *StateTracker) GetEventBus() *EventBus {
	return st.eventBus
}

// GetBackend returns the underlying state backend
func (st *StateTracker) GetBackend() *StateBackend {
	return st.backend
}

// Close closes the state tracker
func (st *StateTracker) Close() error {
	return st.backend.Close()
}

// OperationType represents the type of operation being tracked
type OperationType string

const (
	OpWorkflowExecution  OperationType = "workflow_execution"
	OpAgentRegistration  OperationType = "agent_registration"
	OpAgentUpdate        OperationType = "agent_update"
	OpAgentDelete        OperationType = "agent_delete"
	OpAgentStop          OperationType = "agent_stop"
	OpSchedulerEnable    OperationType = "scheduler_enable"
	OpSchedulerDisable   OperationType = "scheduler_disable"
	OpScheduledExecution OperationType = "scheduled_execution"
	OpSecretCreate       OperationType = "secret_create"
	OpSecretUpdate       OperationType = "secret_update"
	OpSecretDelete       OperationType = "secret_delete"
	OpHookRegister       OperationType = "hook_register"
	OpHookUpdate         OperationType = "hook_update"
	OpHookDelete         OperationType = "hook_delete"
	OpSlothAdd           OperationType = "sloth_add"
	OpSlothUpdate        OperationType = "sloth_update"
	OpSlothDelete        OperationType = "sloth_delete"
	OpBackup             OperationType = "backup"
	OpRestore            OperationType = "restore"
	OpDeployment         OperationType = "deployment"
	OpMaintenance        OperationType = "maintenance"
)

// Operation represents a tracked operation
type Operation struct {
	ID         string                 `json:"id"`
	Type       OperationType          `json:"type"`
	StackName  string                 `json:"stack_name"`
	ResourceID string                 `json:"resource_id"`
	Status     string                 `json:"status"` // pending, running, completed, failed
	StartedAt  time.Time              `json:"started_at"`
	CompletedAt *time.Time            `json:"completed_at,omitempty"`
	Duration   time.Duration          `json:"duration"`
	Metadata   map[string]interface{} `json:"metadata"`
	Error      string                 `json:"error,omitempty"`
	PerformedBy string                `json:"performed_by"`
}

// TrackOperation tracks any operation in the system
func (st *StateTracker) TrackOperation(op *Operation) error {
	// Create or get the appropriate stack based on operation type
	stackName := st.getStackNameForOperation(op.Type)

	stackID, err := st.getOrCreateOperationStack(stackName, op.Type)
	if err != nil {
		return fmt.Errorf("failed to get/create operation stack: %w", err)
	}

	// Generate operation ID if not provided
	if op.ID == "" {
		op.ID = uuid.New().String()
	}

	// Create resource for this operation
	resource := &Resource{
		ID:       fmt.Sprintf("%s-%s", op.Type, op.ID),
		StackID:  stackID,
		Type:     string(op.Type),
		Name:     op.ResourceID,
		Module:   "state_tracker",
		Properties: map[string]interface{}{
			"operation_id":  op.ID,
			"status":        op.Status,
			"started_at":    op.StartedAt,
			"completed_at":  op.CompletedAt,
			"duration":      op.Duration.String(),
			"metadata":      op.Metadata,
			"error":         op.Error,
			"performed_by":  op.PerformedBy,
		},
		Dependencies: []string{},
		State:        op.Status,
		Checksum:     op.ID,
		Metadata: map[string]interface{}{
			"operation_type": op.Type,
			"tracked_at":     time.Now(),
		},
	}

	// Store the resource
	if err := st.backend.GetStackManager().CreateResource(resource); err != nil {
		return fmt.Errorf("failed to create operation resource: %w", err)
	}

	// Create snapshot for important operations
	if st.shouldSnapshot(op.Type, op.Status) {
		description := fmt.Sprintf("%s: %s (%s)", op.Type, op.ResourceID, op.Status)
		_, err := st.backend.CreateSnapshot(stackID, op.PerformedBy, description)
		if err != nil {
			// Log warning but don't fail
			fmt.Printf("Warning: failed to create snapshot: %v\n", err)
		}
	}

	return nil
}

// getStackNameForOperation determines which stack to use for an operation type
func (st *StateTracker) getStackNameForOperation(opType OperationType) string {
	switch opType {
	case OpWorkflowExecution:
		return "workflow-executions"
	case OpAgentRegistration, OpAgentUpdate, OpAgentDelete, OpAgentStop:
		return "agent-operations"
	case OpSchedulerEnable, OpSchedulerDisable, OpScheduledExecution:
		return "scheduler-operations"
	case OpSecretCreate, OpSecretUpdate, OpSecretDelete:
		return "secret-operations"
	case OpHookRegister, OpHookUpdate, OpHookDelete:
		return "hook-operations"
	case OpSlothAdd, OpSlothUpdate, OpSlothDelete:
		return "sloth-operations"
	case OpBackup, OpRestore, OpDeployment, OpMaintenance:
		return "sysadmin-operations"
	default:
		return "general-operations"
	}
}

// getOrCreateOperationStack creates or retrieves an operation stack
func (st *StateTracker) getOrCreateOperationStack(stackName string, opType OperationType) (string, error) {
	stackManager := st.backend.GetStackManager()

	// Try to get existing stack
	existingStack, err := stackManager.GetStackByName(stackName)
	if err == nil {
		return existingStack.ID, nil
	}

	// Create new stack
	stackID := uuid.New().String()
	newStack := &StackState{
		ID:            stackID,
		Name:          stackName,
		Description:   fmt.Sprintf("Stack for %s operations", opType),
		Version:       "1.0.0",
		WorkflowFile:  "<state-tracker>",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata: map[string]interface{}{
			"operation_category": opType,
			"managed_by":         "state_tracker",
		},
	}

	if err := stackManager.CreateStack(newStack); err != nil {
		return "", fmt.Errorf("failed to create operation stack: %w", err)
	}

	return stackID, nil
}

// shouldSnapshot determines if a snapshot should be created for this operation
func (st *StateTracker) shouldSnapshot(opType OperationType, status string) bool {
	// Snapshot on completion or failure of important operations
	importantOps := map[OperationType]bool{
		OpWorkflowExecution:  true,
		OpAgentRegistration:  true,
		OpAgentUpdate:        true,
		OpSchedulerEnable:    true,
		OpDeployment:         true,
		OpBackup:             true,
	}

	if !importantOps[opType] {
		return false
	}

	// Only snapshot on completed or failed
	return status == "completed" || status == "failed"
}

// GetOperationHistory retrieves the history of operations
func (st *StateTracker) GetOperationHistory(opType OperationType, limit int) ([]*Resource, error) {
	stackName := st.getStackNameForOperation(opType)

	stack, err := st.backend.GetStackManager().GetStackByName(stackName)
	if err != nil {
		return nil, fmt.Errorf("operation stack not found: %w", err)
	}

	resources, err := st.backend.GetStackManager().ListResources(stack.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list operation resources: %w", err)
	}

	// Limit results if needed
	if limit > 0 && len(resources) > limit {
		resources = resources[:limit]
	}

	return resources, nil
}

// GetOperationStats returns statistics about operations
func (st *StateTracker) GetOperationStats(opType OperationType) (map[string]interface{}, error) {
	resources, err := st.GetOperationHistory(opType, 0) // Get all
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total":      len(resources),
		"completed":  0,
		"failed":     0,
		"running":    0,
		"pending":    0,
		"by_date":    make(map[string]int),
	}

	for _, resource := range resources {
		switch resource.State {
		case "completed":
			stats["completed"] = stats["completed"].(int) + 1
		case "failed":
			stats["failed"] = stats["failed"].(int) + 1
		case "running":
			stats["running"] = stats["running"].(int) + 1
		case "pending":
			stats["pending"] = stats["pending"].(int) + 1
		}

		// Count by date
		date := resource.CreatedAt.Format("2006-01-02")
		byDate := stats["by_date"].(map[string]int)
		byDate[date]++
	}

	return stats, nil
}

// GetAllOperationStats returns statistics for all operation types
func (st *StateTracker) GetAllOperationStats() (map[OperationType]map[string]interface{}, error) {
	allStats := make(map[OperationType]map[string]interface{})

	opTypes := []OperationType{
		OpWorkflowExecution,
		OpAgentRegistration,
		OpScheduledExecution,
		OpSecretCreate,
		OpHookRegister,
		OpSlothAdd,
		OpBackup,
		OpDeployment,
	}

	for _, opType := range opTypes {
		stats, err := st.GetOperationStats(opType)
		if err != nil {
			// Skip if stack doesn't exist yet
			continue
		}
		allStats[opType] = stats
	}

	return allStats, nil
}

// SearchOperations searches for operations matching criteria
func (st *StateTracker) SearchOperations(criteria map[string]interface{}) ([]*Resource, error) {
	// Get all operation stacks
	stackManager := st.backend.GetStackManager()
	allStacks, err := stackManager.ListStacks()
	if err != nil {
		return nil, err
	}

	var matchingResources []*Resource

	// Search through all operation stacks
	for _, stack := range allStacks {
		// Only search operation stacks
		if !st.isOperationStack(stack.Name) {
			continue
		}

		resources, err := stackManager.ListResources(stack.ID)
		if err != nil {
			continue
		}

		for _, resource := range resources {
			if st.matchesCriteria(resource, criteria) {
				matchingResources = append(matchingResources, resource)
			}
		}
	}

	return matchingResources, nil
}

// isOperationStack checks if a stack is an operation stack
func (st *StateTracker) isOperationStack(stackName string) bool {
	operationStacks := map[string]bool{
		"workflow-executions":   true,
		"agent-operations":      true,
		"scheduler-operations":  true,
		"secret-operations":     true,
		"hook-operations":       true,
		"sloth-operations":      true,
		"sysadmin-operations":   true,
		"general-operations":    true,
	}

	return operationStacks[stackName]
}

// matchesCriteria checks if a resource matches search criteria
func (st *StateTracker) matchesCriteria(resource *Resource, criteria map[string]interface{}) bool {
	for key, value := range criteria {
		switch key {
		case "type":
			if resource.Type != value.(string) {
				return false
			}
		case "status":
			if resource.State != value.(string) {
				return false
			}
		case "date_from":
			if resource.CreatedAt.Before(value.(time.Time)) {
				return false
			}
		case "date_to":
			if resource.CreatedAt.After(value.(time.Time)) {
				return false
			}
		}
	}

	return true
}

// RollbackOperation attempts to rollback an operation
func (st *StateTracker) RollbackOperation(opType OperationType, version int) error {
	stackName := st.getStackNameForOperation(opType)

	stack, err := st.backend.GetStackManager().GetStackByName(stackName)
	if err != nil {
		return fmt.Errorf("operation stack not found: %w", err)
	}

	return st.backend.RollbackToSnapshot(stack.ID, version, "state_tracker")
}
