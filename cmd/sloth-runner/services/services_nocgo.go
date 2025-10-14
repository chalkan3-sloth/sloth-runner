//go:build !cgo
// +build !cgo

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

var errNoCGO = fmt.Errorf("this operation requires CGO support (SQLite).\nPlease use a CGO-enabled build (available for Linux) or compile with CGO_ENABLED=1")

// StackService stub for non-CGO builds
type StackService struct{}

// StateTrackerService stub for non-CGO builds
type StateTrackerService struct{}

// SecretsService stub for non-CGO builds
type SecretsService struct{}

// Secret stub for non-CGO builds
type Secret struct {
	Name      string
	CreatedAt time.Time
}

// ResourceDependency stub for non-CGO builds
type ResourceDependency struct {
	ResourceID   string
	DependsOnID  string
	ResourceName string
	DependsOnName string
	Type         string
}

// NewStackService returns a stub error
func NewStackService(dbPath string) (*StackService, error) {
	return nil, errNoCGO
}

// NewStateTrackerService returns a stub error
func NewStateTrackerService(dbPath string) (*StateTrackerService, error) {
	return nil, errNoCGO
}

// NewSecretsService returns a stub error
func NewSecretsService(ctx context.Context, dbPath string, masterKey []byte) (*SecretsService, error) {
	return nil, errNoCGO
}

// GetOrCreateSalt returns a stub error
func GetOrCreateSalt(stackService *StackService, stackID string) ([]byte, error) {
	return nil, errNoCGO
}

// GetGlobalStateTracker returns nil for non-CGO builds
func GetGlobalStateTracker() (*StateTrackerService, error) {
	return nil, errNoCGO
}

// StackService stub methods
func (s *StackService) Close() error { return errNoCGO }
func (s *StackService) GetOrCreateStack(stackName, workflowName, filePath string) (string, error) { return "", errNoCGO }
func (s *StackService) UpdateStackStatus(stackID, status string) error { return errNoCGO }
func (s *StackService) RecordExecution(stackID string, execution *stack.StackExecution) error { return errNoCGO }
func (s *StackService) UpdateStackAfterExecution(stackID string, executionID string, status string, resources []*stack.Resource, outputs map[string]string, metadata map[string]interface{}) error { return errNoCGO }
func (s *StackService) GetStack(stackID string) (*stack.StackState, error) { return nil, errNoCGO }
func (s *StackService) GetStackByName(stackName string) (*stack.StackState, error) { return nil, errNoCGO }
func (s *StackService) GetManager() *stack.StackManager { return nil }
func (s *StackService) GetBackend() *stack.StateBackend { return nil }
func (s *StackService) CreateSnapshot(stackID, createdBy, description string) (int, error) { return 0, errNoCGO }
func (s *StackService) GetSnapshot(stackID string, version int) (*stack.StateSnapshot, error) { return nil, errNoCGO }
func (s *StackService) ListSnapshots(stackID string) ([]stack.StateSnapshot, error) { return nil, errNoCGO }
func (s *StackService) RollbackToSnapshot(stackID string, version int, performedBy string) error { return errNoCGO }
func (s *StackService) DetectDrift(stackID, resourceID string, expectedState, actualState map[string]interface{}) error { return errNoCGO }
func (s *StackService) GetDriftInfo(stackID string) ([]*stack.DriftInfo, error) { return nil, errNoCGO }
func (s *StackService) LockState(stackID, lockID, operation, who string, duration time.Duration) error { return errNoCGO }
func (s *StackService) UnlockState(stackID, lockID string) error { return errNoCGO }
func (s *StackService) AddTag(stackID, tag string) error { return errNoCGO }
func (s *StackService) GetTags(stackID string) ([]string, error) { return nil, errNoCGO }
func (s *StackService) AddResourceDependency(resourceID, dependsOnID, depType string) error { return errNoCGO }
func (s *StackService) GetResourceDependencies(resourceID string) ([]string, error) { return nil, errNoCGO }
func (s *StackService) GetActivity(stackID string, limit int) ([]map[string]interface{}, error) { return nil, errNoCGO }
func (s *StackService) GetStackResourceDependencies(stackID string) ([]ResourceDependency, error) { return nil, errNoCGO }
func (s *StackService) ListStackResources(stackID string) ([]*stack.Resource, error) { return nil, errNoCGO }
func (s *StackService) ListStacks() ([]*stack.StackState, error) { return nil, errNoCGO }
func (s *StackService) IsLocked(stackID string) (bool, error) { return false, errNoCGO }
func (s *StackService) GetLockInfo(stackID string) (*stack.StateLock, error) { return nil, errNoCGO }

// SecretsService stub methods
func (s *SecretsService) Close() error { return errNoCGO }
func (s *SecretsService) AddSecret(ctx context.Context, stackID, name, value, password string, salt []byte) error { return errNoCGO }
func (s *SecretsService) GetSecret(ctx context.Context, stackID, name, password string, salt []byte) (string, error) { return "", errNoCGO }
func (s *SecretsService) ListSecrets(ctx context.Context, stackID string) ([]Secret, error) { return nil, errNoCGO }
func (s *SecretsService) GetAllSecrets(ctx context.Context, stackID, password string, salt []byte) (map[string]string, error) { return nil, errNoCGO }
func (s *SecretsService) RemoveSecret(ctx context.Context, stackID, name string) error { return errNoCGO }
func (s *SecretsService) RemoveAllSecrets(ctx context.Context, stackID string) error { return errNoCGO }
func (s *SecretsService) HasSecrets(ctx context.Context, stackID string) (bool, error) { return false, errNoCGO }

// StateTrackerService stub methods
func (s *StateTrackerService) Close() error { return errNoCGO }
func (s *StateTrackerService) RecordOperation(op *stack.Operation) error { return errNoCGO }
func (s *StateTrackerService) GetOperationHistory(stackName string, limit int) ([]stack.Operation, error) { return nil, errNoCGO }
