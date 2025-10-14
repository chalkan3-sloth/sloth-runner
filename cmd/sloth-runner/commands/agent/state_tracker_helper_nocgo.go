//go:build !cgo
// +build !cgo

package agent

import (
	"log/slog"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// trackAgentOperation is a no-op for non-CGO builds
func trackAgentOperation(opType stack.OperationType, agentName string, status string, metadata map[string]interface{}) {
	slog.Debug("State tracking not available in non-CGO builds", "operation", opType, "agent", agentName, "status", status)
}

// trackAgentRegistration is a no-op for non-CGO builds
func trackAgentRegistration(agentName string, host string, port int, success bool) {
	// No-op
}

// trackAgentUpdate is a no-op for non-CGO builds
func trackAgentUpdate(agentName string, version string, success bool) {
	// No-op
}

// trackAgentDelete is a no-op for non-CGO builds
func trackAgentDelete(agentName string, success bool) {
	// No-op
}

// trackAgentStop is a no-op for non-CGO builds
func trackAgentStop(agentName string, success bool) {
	// No-op
}
