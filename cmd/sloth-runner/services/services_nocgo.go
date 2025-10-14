//go:build !cgo
// +build !cgo

package services

import (
	"context"
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// StackService stub for non-CGO builds
type StackService struct{}

// StateTrackerService stub for non-CGO builds
type StateTrackerService struct{}

// SecretsService stub for non-CGO builds
type SecretsService struct{}

// NewStackService returns a stub error
func NewStackService(dbPath string) (*StackService, error) {
	return nil, fmt.Errorf("StackService requires CGO support (SQLite). Please use a CGO-enabled build")
}

// NewStateTrackerService returns a stub error
func NewStateTrackerService(dbPath string) (*StateTrackerService, error) {
	return nil, fmt.Errorf("StateTrackerService requires CGO support (SQLite). Please use a CGO-enabled build")
}

// NewSecretsService returns a stub error
func NewSecretsService(ctx context.Context, dbPath string, masterKey []byte) (*SecretsService, error) {
	return nil, fmt.Errorf("SecretsService requires CGO support (SQLite). Please use a CGO-enabled build")
}

// GetOrCreateSalt returns a stub error
func GetOrCreateSalt(stackService *StackService, stackID string) ([]byte, error) {
	return nil, fmt.Errorf("GetOrCreateSalt requires CGO support (SQLite). Please use a CGO-enabled build")
}

// GetGlobalStateTracker returns nil for non-CGO builds
func GetGlobalStateTracker() *stack.StateTracker {
	return nil
}
