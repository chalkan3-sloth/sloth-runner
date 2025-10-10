package sloth

import (
	"log/slog"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// trackSlothOperation tracks sloth file operations
func trackSlothOperation(opType stack.OperationType, slothName string, filePath string, success bool) {
	tracker, err := services.GetGlobalStateTracker()
	if err != nil {
		slog.Warn("Failed to get state tracker", "error", err)
		return
	}

	status := "completed"
	if !success {
		status = "failed"
	}

	err = tracker.TrackSlothOperation(opType, slothName, filePath, status, "cli-user")
	if err != nil {
		slog.Warn("Failed to track sloth operation", "error", err, "operation", opType, "sloth", slothName)
	} else {
		slog.Debug("Sloth operation tracked", "operation", opType, "sloth", slothName, "status", status)
	}
}

// trackSlothAdd tracks sloth file addition
func trackSlothAdd(slothName string, filePath string, success bool) {
	trackSlothOperation(stack.OpSlothAdd, slothName, filePath, success)
}

// trackSlothUpdate tracks sloth file updates
func trackSlothUpdate(slothName string, filePath string, success bool) {
	trackSlothOperation(stack.OpSlothUpdate, slothName, filePath, success)
}

// trackSlothDelete tracks sloth file deletion
func trackSlothDelete(slothName string, success bool) {
	trackSlothOperation(stack.OpSlothDelete, slothName, "", success)
}
