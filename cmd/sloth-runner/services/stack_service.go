package services

import (
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/google/uuid"
)

// StackService handles stack operations
// This implements the Service Layer pattern
type StackService struct {
	manager *stack.StackManager
}

// NewStackService creates a new stack service
func NewStackService() (*StackService, error) {
	manager, err := stack.NewStackManager("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize stack manager: %w", err)
	}
	return &StackService{manager: manager}, nil
}

// Close closes the stack service
func (s *StackService) Close() error {
	return s.manager.Close()
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

// GetManager returns the underlying stack manager
func (s *StackService) GetManager() *stack.StackManager {
	return s.manager
}
