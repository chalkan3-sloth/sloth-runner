//go:build cgo
// +build cgo

package stack

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// EventType represents the type of state event
type EventType string

const (
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

// StateEvent represents an event in the state system
type StateEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	StackID   string                 `json:"stack_id,omitempty"`
	StackName string                 `json:"stack_name,omitempty"`
	Operation *Operation             `json:"operation,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Severity  string                 `json:"severity"` // info, warning, error, critical
}

// EventHandler is a function that handles state events
type EventHandler func(context.Context, *StateEvent) error

// EventBus is the central event distribution system
type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
	eventLog []*StateEvent
	logMu    sync.RWMutex
	maxLog   int
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
		eventLog: make([]*StateEvent, 0),
		maxLog:   1000, // Keep last 1000 events in memory
	}
}

// Subscribe registers a handler for a specific event type
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.handlers[eventType] == nil {
		eb.handlers[eventType] = make([]EventHandler, 0)
	}

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	slog.Debug("Event handler subscribed", "event_type", eventType)
}

// SubscribeAll registers a handler for all event types
func (eb *EventBus) SubscribeAll(handler EventHandler) {
	allTypes := []EventType{
		EventOperationStarted,
		EventOperationCompleted,
		EventOperationFailed,
		EventSnapshotCreated,
		EventSnapshotRestored,
		EventDriftDetected,
		EventStateLocked,
		EventStateUnlocked,
		EventResourceCreated,
		EventResourceUpdated,
		EventResourceDeleted,
		EventValidationFailed,
		EventValidationPassed,
		EventBackupCreated,
		EventBackupRestored,
	}

	for _, eventType := range allTypes {
		eb.Subscribe(eventType, handler)
	}
}

// Publish publishes an event to all registered handlers
func (eb *EventBus) Publish(ctx context.Context, event *StateEvent) error {
	// Add to event log
	eb.logEvent(event)

	// Get handlers for this event type
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		slog.Debug("No handlers for event type", "event_type", event.Type)
		return nil
	}

	// Execute handlers concurrently
	var wg sync.WaitGroup
	errors := make([]error, 0)
	errorsMu := sync.Mutex{}

	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()

			if err := h(ctx, event); err != nil {
				errorsMu.Lock()
				errors = append(errors, err)
				errorsMu.Unlock()
				slog.Warn("Event handler error", "event_type", event.Type, "error", err)
			}
		}(handler)
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("event handling errors: %d handlers failed", len(errors))
	}

	return nil
}

// logEvent adds an event to the in-memory log
func (eb *EventBus) logEvent(event *StateEvent) {
	eb.logMu.Lock()
	defer eb.logMu.Unlock()

	eb.eventLog = append(eb.eventLog, event)

	// Trim log if it exceeds maxLog
	if len(eb.eventLog) > eb.maxLog {
		eb.eventLog = eb.eventLog[len(eb.eventLog)-eb.maxLog:]
	}
}

// GetRecentEvents returns recent events
func (eb *EventBus) GetRecentEvents(limit int) []*StateEvent {
	eb.logMu.RLock()
	defer eb.logMu.RUnlock()

	if limit <= 0 || limit > len(eb.eventLog) {
		limit = len(eb.eventLog)
	}

	// Return most recent events
	start := len(eb.eventLog) - limit
	if start < 0 {
		start = 0
	}

	events := make([]*StateEvent, limit)
	copy(events, eb.eventLog[start:])

	return events
}

// GetEventsByType returns events filtered by type
func (eb *EventBus) GetEventsByType(eventType EventType, limit int) []*StateEvent {
	eb.logMu.RLock()
	defer eb.logMu.RUnlock()

	filtered := make([]*StateEvent, 0)

	// Iterate in reverse to get most recent first
	for i := len(eb.eventLog) - 1; i >= 0 && len(filtered) < limit; i-- {
		if eb.eventLog[i].Type == eventType {
			filtered = append(filtered, eb.eventLog[i])
		}
	}

	return filtered
}

// GetEventsByStack returns events for a specific stack
func (eb *EventBus) GetEventsByStack(stackID string, limit int) []*StateEvent {
	eb.logMu.RLock()
	defer eb.logMu.RUnlock()

	filtered := make([]*StateEvent, 0)

	for i := len(eb.eventLog) - 1; i >= 0 && len(filtered) < limit; i-- {
		if eb.eventLog[i].StackID == stackID {
			filtered = append(filtered, eb.eventLog[i])
		}
	}

	return filtered
}

// ClearLog clears the event log
func (eb *EventBus) ClearLog() {
	eb.logMu.Lock()
	defer eb.logMu.Unlock()

	eb.eventLog = make([]*StateEvent, 0)
}

// GetStats returns statistics about events
func (eb *EventBus) GetStats() map[string]interface{} {
	eb.logMu.RLock()
	defer eb.logMu.RUnlock()

	stats := map[string]interface{}{
		"total_events": len(eb.eventLog),
		"by_type":      make(map[EventType]int),
		"by_severity":  make(map[string]int),
	}

	byType := stats["by_type"].(map[EventType]int)
	bySeverity := stats["by_severity"].(map[string]int)

	for _, event := range eb.eventLog {
		byType[event.Type]++
		bySeverity[event.Severity]++
	}

	return stats
}
