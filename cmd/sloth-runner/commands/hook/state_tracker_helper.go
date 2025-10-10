package hook

import (
	"log/slog"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// trackHookOperation tracks hook operations
func trackHookOperation(opType stack.OperationType, hookName string, hookType string, success bool) {
	tracker, err := services.GetGlobalStateTracker()
	if err != nil {
		slog.Warn("Failed to get state tracker", "error", err)
		return
	}

	status := "completed"
	if !success {
		status = "failed"
	}

	err = tracker.TrackHookOperation(opType, hookName, hookType, status, "cli-user")
	if err != nil {
		slog.Warn("Failed to track hook operation", "error", err, "operation", opType, "hook", hookName)
	} else {
		slog.Debug("Hook operation tracked", "operation", opType, "hook", hookName, "status", status)
	}
}

// trackHookRegister tracks hook registration
func trackHookRegister(hookName string, hookType string, success bool) {
	trackHookOperation(stack.OpHookRegister, hookName, hookType, success)
}

// trackHookUpdate tracks hook updates
func trackHookUpdate(hookName string, hookType string, success bool) {
	trackHookOperation(stack.OpHookUpdate, hookName, hookType, success)
}

// trackHookDelete tracks hook deletion
func trackHookDelete(hookName string, success bool) {
	trackHookOperation(stack.OpHookDelete, hookName, "", success)
}
