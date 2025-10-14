//go:build cgo
// +build cgo

package agent

import (
	"log/slog"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// trackAgentOperation is a helper to track agent operations
func trackAgentOperation(opType stack.OperationType, agentName string, status string, metadata map[string]interface{}) {
	tracker, err := services.GetGlobalStateTracker()
	if err != nil {
		slog.Warn("Failed to get state tracker", "error", err)
		return
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	err = tracker.TrackAgentOperation(opType, agentName, status, metadata, "cli-user")
	if err != nil {
		slog.Warn("Failed to track agent operation", "error", err, "operation", opType, "agent", agentName)
	} else {
		slog.Debug("Agent operation tracked", "operation", opType, "agent", agentName, "status", status)
	}
}

// trackAgentRegistration tracks agent registration
func trackAgentRegistration(agentName string, host string, port int, success bool) {
	status := "completed"
	if !success {
		status = "failed"
	}

	metadata := map[string]interface{}{
		"host": host,
		"port": port,
	}

	trackAgentOperation(stack.OpAgentRegistration, agentName, status, metadata)
}

// trackAgentUpdate tracks agent updates
func trackAgentUpdate(agentName string, version string, success bool) {
	status := "completed"
	if !success {
		status = "failed"
	}

	metadata := map[string]interface{}{
		"version": version,
	}

	trackAgentOperation(stack.OpAgentUpdate, agentName, status, metadata)
}

// trackAgentDelete tracks agent deletion
func trackAgentDelete(agentName string, success bool) {
	status := "completed"
	if !success {
		status = "failed"
	}

	trackAgentOperation(stack.OpAgentDelete, agentName, status, nil)
}

// trackAgentStop tracks agent stop
func trackAgentStop(agentName string, success bool) {
	status := "completed"
	if !success {
		status = "failed"
	}

	trackAgentOperation(stack.OpAgentStop, agentName, status, nil)
}
