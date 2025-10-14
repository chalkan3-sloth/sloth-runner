//go:build !cgo
// +build !cgo

package stack

import (
	"context"
	"fmt"
	"time"
)

// EventType represents the type of state event
type EventType string

const (
	// EventType constants for non-CGO builds
	EventOperationStarted   EventType = "operation.started"
	EventOperationCompleted EventType = "operation.completed"
	EventOperationFailed    EventType = "operation.failed"
	EventSnapshotCreated    EventType = "snapshot.created"
	EventSnapshotRestored   EventType = "snapshot.restored"
	EventDriftDetected      EventType = "drift.detected"
	EventStateLocked        EventType = "state.locked"
	EventStateUnlocked      EventType = "state.unlocked"
	EventResourceCreated    EventType = "resource.created"
	EventResourceUpdated    EventType = "resource.updated"
	EventResourceDeleted    EventType = "resource.deleted"
	EventValidationFailed   EventType = "validation.failed"
	EventValidationPassed   EventType = "validation.passed"
	EventBackupCreated      EventType = "backup.created"
	EventBackupRestored     EventType = "backup.restored"
)

// StateEvent represents a state change event
type StateEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	StackID   string                 `json:"stack_id,omitempty"`
	StackName string                 `json:"stack_name,omitempty"`
	Operation *Operation             `json:"operation,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Severity  string                 `json:"severity"`
}

// EventHandler is a function that handles state events
type EventHandler func(context.Context, *StateEvent) error

// EventBus stub for non-CGO builds
type EventBus struct{}

// NewEventBus returns an error for non-CGO builds
func NewEventBus() *EventBus {
	return &EventBus{}
}

// Subscribe stub
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	// no-op for non-CGO builds
}

// Unsubscribe stub
func (eb *EventBus) Unsubscribe(eventType EventType, handler EventHandler) {
	// no-op for non-CGO builds
}

// Publish stub
func (eb *EventBus) Publish(ctx context.Context, event *StateEvent) error {
	return fmt.Errorf("event bus not available in non-CGO builds")
}

// GetEventLog stub
func (eb *EventBus) GetEventLog(limit int) []*StateEvent {
	return nil
}

// Clear stub
func (eb *EventBus) Clear() {
	// no-op for non-CGO builds
}
