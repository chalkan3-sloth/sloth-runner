//go:build cgo
// +build cgo

package secrets

import (
	"log/slog"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// trackSecretOperation tracks secret operations
func trackSecretOperation(opType stack.OperationType, secretKey string, stackID string, success bool) {
	tracker, err := services.GetGlobalStateTracker()
	if err != nil {
		slog.Warn("Failed to get state tracker", "error", err)
		return
	}

	status := "completed"
	if !success {
		status = "failed"
	}

	err = tracker.TrackSecretOperation(opType, secretKey, stackID, status, "cli-user")
	if err != nil {
		slog.Warn("Failed to track secret operation", "error", err, "operation", opType, "secret", secretKey)
	} else {
		slog.Debug("Secret operation tracked", "operation", opType, "secret", secretKey, "status", status)
	}
}

// trackSecretCreate tracks secret creation
func trackSecretCreate(secretKey string, stackID string, success bool) {
	trackSecretOperation(stack.OpSecretCreate, secretKey, stackID, success)
}

// trackSecretUpdate tracks secret updates
func trackSecretUpdate(secretKey string, stackID string, success bool) {
	trackSecretOperation(stack.OpSecretUpdate, secretKey, stackID, success)
}

// trackSecretDelete tracks secret deletion
func trackSecretDelete(secretKey string, stackID string, success bool) {
	trackSecretOperation(stack.OpSecretDelete, secretKey, stackID, success)
}
