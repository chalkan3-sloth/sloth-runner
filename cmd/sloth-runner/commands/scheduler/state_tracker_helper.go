package scheduler

import (
	"log/slog"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// trackSchedulerOperation tracks scheduler operations
func trackSchedulerOperation(opType stack.OperationType, workflowName string, schedule string, status string, duration time.Duration, errorMsg string) {
	tracker, err := services.GetGlobalStateTracker()
	if err != nil {
		slog.Warn("Failed to get state tracker", "error", err)
		return
	}

	err = tracker.TrackSchedulerOperation(opType, workflowName, schedule, status, duration, errorMsg, "cli-user")
	if err != nil {
		slog.Warn("Failed to track scheduler operation", "error", err, "operation", opType, "workflow", workflowName)
	} else {
		slog.Debug("Scheduler operation tracked", "operation", opType, "workflow", workflowName, "status", status)
	}
}

// trackSchedulerEnable tracks scheduler enable operations
func trackSchedulerEnable(workflowName string, schedule string, success bool) {
	status := "completed"
	if !success {
		status = "failed"
	}

	trackSchedulerOperation(stack.OpSchedulerEnable, workflowName, schedule, status, 0, "")
}

// trackSchedulerDisable tracks scheduler disable operations
func trackSchedulerDisable(workflowName string, success bool) {
	status := "completed"
	if !success {
		status = "failed"
	}

	trackSchedulerOperation(stack.OpSchedulerDisable, workflowName, "", status, 0, "")
}

// trackScheduledExecution tracks scheduled workflow executions
func trackScheduledExecution(workflowName string, schedule string, success bool, duration time.Duration, errorMsg string) {
	status := "completed"
	if !success {
		status = "failed"
	}

	trackSchedulerOperation(stack.OpScheduledExecution, workflowName, schedule, status, duration, errorMsg)
}
