package services

import (
	"fmt"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/google/uuid"
)

// StackService handles stack operations with state management
// This implements the Service Layer pattern with StateBackend for advanced features
type StackService struct {
	backend *stack.StateBackend
	manager *stack.StackManager // Direct access to manager for backward compatibility
}

// NewStackService creates a new stack service with state backend
func NewStackService() (*StackService, error) {
	// Check for test database path from environment
	dbPath := os.Getenv("SLOTH_RUNNER_DB_PATH")

	backend, err := stack.NewStateBackend(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize state backend: %w", err)
	}
	return &StackService{
		backend: backend,
		manager: backend.GetStackManager(),
	}, nil
}

// Close closes the stack service
func (s *StackService) Close() error {
	return s.backend.Close()
}

// GetOrCreateStack gets an existing stack or creates a new one
func (s *StackService) GetOrCreateStack(stackName, workflowName, filePath string) (string, error) {
	existingStack, err := s.manager.GetStackByName(stackName)
	if err == nil {
		return existingStack.ID, nil
	}

	// Create new stack
	stackID := uuid.New().String()
	newStack := &stack.StackState{
		ID:            stackID,
		Name:          stackName,
		Description:   fmt.Sprintf("Stack for workflow: %s", workflowName),
		Version:       "1.0.0",
		WorkflowFile:  filePath,
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := s.manager.CreateStack(newStack); err != nil {
		return "", fmt.Errorf("failed to create stack: %w", err)
	}

	return stackID, nil
}

// UpdateStackStatus updates the stack status
func (s *StackService) UpdateStackStatus(stackID, status string) error {
	currentStack, err := s.manager.GetStack(stackID)
	if err != nil {
		return err
	}
	currentStack.Status = status
	return s.manager.UpdateStack(currentStack)
}

// RecordExecution records a stack execution
func (s *StackService) RecordExecution(stackID string, execution *stack.StackExecution) error {
	return s.manager.RecordExecution(stackID, execution)
}

// UpdateStackAfterExecution updates stack state after execution
func (s *StackService) UpdateStackAfterExecution(
	stackID string,
	status string,
	duration time.Duration,
	errorMessage string,
	outputs map[string]interface{},
) error {
	stackState, err := s.manager.GetStack(stackID)
	if err != nil {
		return err
	}

	stackState.Status = status
	stackState.LastDuration = duration
	stackState.LastError = errorMessage
	stackState.ExecutionCount++
	stackState.Outputs = outputs

	if status == "completed" {
		now := time.Now()
		stackState.CompletedAt = &now
	}

	return s.manager.UpdateStack(stackState)
}

// GetStack returns the stack by ID
func (s *StackService) GetStack(stackID string) (*stack.StackState, error) {
	return s.manager.GetStack(stackID)
}

// GetStackByName returns the stack by name
func (s *StackService) GetStackByName(stackName string) (*stack.StackState, error) {
	return s.manager.GetStackByName(stackName)
}

// GetManager returns the underlying stack manager
func (s *StackService) GetManager() *stack.StackManager {
	return s.manager
}

// GetBackend returns the underlying state backend
func (s *StackService) GetBackend() *stack.StateBackend {
	return s.backend
}

// CreateSnapshot creates a new state snapshot (version)
func (s *StackService) CreateSnapshot(stackID, createdBy, description string) (int, error) {
	return s.backend.CreateSnapshot(stackID, createdBy, description)
}

// GetSnapshot retrieves a specific snapshot
func (s *StackService) GetSnapshot(stackID string, version int) (*stack.StateSnapshot, error) {
	return s.backend.GetSnapshot(stackID, version)
}

// ListSnapshots lists all snapshots for a stack
func (s *StackService) ListSnapshots(stackID string) ([]stack.StateSnapshot, error) {
	return s.backend.ListSnapshots(stackID)
}

// RollbackToSnapshot rolls back stack to a specific version
func (s *StackService) RollbackToSnapshot(stackID string, version int, performedBy string) error {
	return s.backend.RollbackToSnapshot(stackID, version, performedBy)
}

// DetectDrift checks for drift between expected and actual state
func (s *StackService) DetectDrift(stackID, resourceID string, expectedState, actualState map[string]interface{}) error {
	return s.backend.DetectDrift(stackID, resourceID, expectedState, actualState)
}

// GetDriftInfo retrieves drift information for a stack
func (s *StackService) GetDriftInfo(stackID string) ([]*stack.DriftInfo, error) {
	return s.backend.GetDriftInfo(stackID)
}

// LockState acquires a lock on the state
func (s *StackService) LockState(stackID, lockID, operation, who string, duration time.Duration) error {
	return s.backend.LockState(stackID, lockID, operation, who, duration)
}

// UnlockState releases a state lock
func (s *StackService) UnlockState(stackID, lockID string) error {
	return s.backend.UnlockState(stackID, lockID)
}

// AddTag adds a tag to a stack
func (s *StackService) AddTag(stackID, tag string) error {
	return s.backend.AddTag(stackID, tag)
}

// GetTags retrieves tags for a stack
func (s *StackService) GetTags(stackID string) ([]string, error) {
	return s.backend.GetTags(stackID)
}

// AddResourceDependency records a dependency between resources
func (s *StackService) AddResourceDependency(resourceID, dependsOnID, depType string) error {
	return s.backend.AddResourceDependency(resourceID, dependsOnID, depType)
}

// GetResourceDependencies retrieves dependencies for a resource
func (s *StackService) GetResourceDependencies(resourceID string) ([]string, error) {
	return s.backend.GetResourceDependencies(resourceID)
}

// GetActivity retrieves activity log for a stack
func (s *StackService) GetActivity(stackID string, limit int) ([]map[string]interface{}, error) {
	return s.backend.GetActivity(stackID, limit)
}

// ResourceDependency represents a dependency between resources
type ResourceDependency struct {
	ResourceID  string `json:"resource_id"`
	DependsOnID string `json:"depends_on_id"`
}

// GetStackResourceDependencies retrieves all dependencies for resources in a stack
func (s *StackService) GetStackResourceDependencies(stackID string) ([]ResourceDependency, error) {
	// Get all resources in the stack
	resources, err := s.manager.ListResources(stackID)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	deps := make([]ResourceDependency, 0)
	for _, resource := range resources {
		for _, depID := range resource.Dependencies {
			deps = append(deps, ResourceDependency{
				ResourceID:  resource.ID,
				DependsOnID: depID,
			})
		}
	}

	return deps, nil
}

// ListStackResources returns all resources for a stack
func (s *StackService) ListStackResources(stackID string) ([]*stack.Resource, error) {
	return s.manager.ListResources(stackID)
}

// ListStacks lists all stacks
func (s *StackService) ListStacks() ([]*stack.StackState, error) {
	return s.manager.ListStacks()
}

// IsLocked checks if a stack is currently locked
func (s *StackService) IsLocked(stackID string) (bool, error) {
	lockInfo, err := s.backend.GetLockInfo(stackID)
	if err != nil {
		return false, err
	}
	return lockInfo != nil, nil
}

// GetLockInfo retrieves lock information for a stack
func (s *StackService) GetLockInfo(stackID string) (*stack.StateLock, error) {
	return s.backend.GetLockInfo(stackID)
}
