//go:build cgo
// +build cgo

package stack

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// registerDefaultHandlers registers default event handlers
func (st *StateTracker) registerDefaultHandlers() {
	// Persist all events to database
	st.eventBus.SubscribeAll(func(ctx context.Context, event *StateEvent) error {
		// Persist event to database
		if err := st.persistEvent(event); err != nil {
			slog.Warn("Failed to persist event", "error", err, "event_id", event.ID)
		}
		return nil
	})

	// Log all events
	st.eventBus.SubscribeAll(func(ctx context.Context, event *StateEvent) error {
		slog.Info("State event",
			"type", event.Type,
			"source", event.Source,
			"stack", event.StackName,
			"severity", event.Severity,
		)
		return nil
	})

	// Handle critical events
	st.eventBus.Subscribe(EventValidationFailed, func(ctx context.Context, event *StateEvent) error {
		slog.Error("State validation failed",
			"stack", event.StackName,
			"data", event.Data,
		)
		// Auto-repair could be triggered here
		return nil
	})

	// Handle drift detection
	st.eventBus.Subscribe(EventDriftDetected, func(ctx context.Context, event *StateEvent) error {
		slog.Warn("State drift detected",
			"stack", event.StackName,
			"data", event.Data,
		)
		return nil
	})
}

// emitEvent emits an event to the event bus
func (st *StateTracker) emitEvent(eventType EventType, stackID, stackName, source string, data map[string]interface{}, severity string) {
	event := &StateEvent{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    source,
		StackID:   stackID,
		StackName: stackName,
		Data:      data,
		Severity:  severity,
	}

	// Emit asynchronously to avoid blocking operations
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := st.eventBus.Publish(ctx, event); err != nil {
			slog.Warn("Failed to publish event", "error", err, "event_type", eventType)
		}
	}()
}

// TrackOperationWithEvents tracks an operation and emits events
func (st *StateTracker) TrackOperationWithEvents(op *Operation) error {
	// Emit operation started event
	st.emitEvent(
		EventOperationStarted,
		"",
		op.StackName,
		"state_tracker",
		map[string]interface{}{
			"operation_type": op.Type,
			"operation_id":   op.ID,
			"resource_id":    op.ResourceID,
		},
		"info",
	)

	// Track the operation
	err := st.TrackOperation(op)

	// Emit completion or failure event
	if err != nil {
		st.emitEvent(
			EventOperationFailed,
			"",
			op.StackName,
			"state_tracker",
			map[string]interface{}{
				"operation_type": op.Type,
				"operation_id":   op.ID,
				"resource_id":    op.ResourceID,
				"error":          err.Error(),
			},
			"error",
		)
		return err
	}

	if op.Status == "completed" {
		st.emitEvent(
			EventOperationCompleted,
			"",
			op.StackName,
			"state_tracker",
			map[string]interface{}{
				"operation_type": op.Type,
				"operation_id":   op.ID,
				"resource_id":    op.ResourceID,
				"duration":       op.Duration.String(),
			},
			"info",
		)
	} else if op.Status == "failed" {
		st.emitEvent(
			EventOperationFailed,
			"",
			op.StackName,
			"state_tracker",
			map[string]interface{}{
				"operation_type": op.Type,
				"operation_id":   op.ID,
				"resource_id":    op.ResourceID,
				"error":          op.Error,
			},
			"error",
		)
	}

	return nil
}

// CreateSnapshotWithEvent creates a snapshot and emits an event
func (st *StateTracker) CreateSnapshotWithEvent(stackID, creator, description string) (int, error) {
	version, err := st.backend.CreateSnapshot(stackID, creator, description)

	if err == nil {
		stack, _ := st.backend.GetStackManager().GetStack(stackID)
		stackName := ""
		if stack != nil {
			stackName = stack.Name
		}

		st.emitEvent(
			EventSnapshotCreated,
			stackID,
			stackName,
			creator,
			map[string]interface{}{
				"version":     version,
				"description": description,
			},
			"info",
		)
	}

	return version, err
}

// RollbackToSnapshotWithEvent rolls back to a snapshot and emits an event
func (st *StateTracker) RollbackToSnapshotWithEvent(stackID string, version int, performer string) error {
	stack, _ := st.backend.GetStackManager().GetStack(stackID)
	stackName := ""
	if stack != nil {
		stackName = stack.Name
	}

	err := st.backend.RollbackToSnapshot(stackID, version, performer)

	if err == nil {
		st.emitEvent(
			EventSnapshotRestored,
			stackID,
			stackName,
			performer,
			map[string]interface{}{
				"version": version,
			},
			"warning",
		)
	}

	return err
}

// DetectDriftWithEvent detects drift and emits an event
func (st *StateTracker) DetectDriftWithEvent(stackID string) (bool, []string, error) {
	stack, err := st.backend.GetStackManager().GetStack(stackID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get stack: %w", err)
	}

	// Get existing drift information from the database
	driftInfos, err := st.backend.GetDriftInfo(stackID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get drift info: %w", err)
	}

	// Build list of drift messages
	drifts := make([]string, 0)
	hasDrift := false

	for _, driftInfo := range driftInfos {
		if driftInfo.IsDrifted && driftInfo.ResolutionStatus == "pending" {
			hasDrift = true
			driftMsg := fmt.Sprintf("resource %s: %d field(s) drifted",
				driftInfo.ResourceID, len(driftInfo.DriftedFields))
			drifts = append(drifts, driftMsg)
		}
	}

	// Also check resources marked as "drift" state
	resources, err := st.backend.GetStackManager().ListResources(stackID)
	if err == nil {
		for _, resource := range resources {
			if resource.State == "drift" {
				hasDrift = true
				driftMsg := fmt.Sprintf("resource %s: state marked as drift", resource.Name)
				drifts = append(drifts, driftMsg)
			}
		}
	}

	if hasDrift {
		st.emitEvent(
			EventDriftDetected,
			stackID,
			stack.Name,
			"drift_detector",
			map[string]interface{}{
				"drift_count": len(drifts),
				"drifts":      drifts,
			},
			"warning",
		)
	}

	return hasDrift, drifts, nil
}

// LockStateWithEvent locks state and emits an event
func (st *StateTracker) LockStateWithEvent(stackID, lockedBy, reason string) error {
	// Generate a unique lock ID
	lockID := uuid.New().String()

	// Use 1 hour as default duration (0 means no expiration in some implementations)
	duration := 1 * time.Hour

	// Lock the state with proper parameters
	err := st.backend.LockState(stackID, lockID, reason, lockedBy, duration)

	if err == nil {
		stack, _ := st.backend.GetStackManager().GetStack(stackID)
		stackName := ""
		if stack != nil {
			stackName = stack.Name
		}

		st.emitEvent(
			EventStateLocked,
			stackID,
			stackName,
			lockedBy,
			map[string]interface{}{
				"lock_id":  lockID,
				"reason":   reason,
				"duration": duration.String(),
			},
			"info",
		)
	}

	return err
}

// UnlockStateWithEvent unlocks state and emits an event
func (st *StateTracker) UnlockStateWithEvent(stackID, unlockedBy string) error {
	// Get the current lock info to retrieve lock_id
	lockInfo, err := st.backend.GetLockInfo(stackID)
	if err != nil {
		return fmt.Errorf("failed to get lock info: %w", err)
	}

	if lockInfo == nil {
		return fmt.Errorf("stack is not locked")
	}

	// Unlock using the lock ID
	err = st.backend.UnlockState(stackID, lockInfo.LockID)

	if err == nil {
		stack, _ := st.backend.GetStackManager().GetStack(stackID)
		stackName := ""
		if stack != nil {
			stackName = stack.Name
		}

		st.emitEvent(
			EventStateUnlocked,
			stackID,
			stackName,
			unlockedBy,
			map[string]interface{}{
				"lock_id":   lockInfo.LockID,
				"locked_by": lockInfo.Who,
			},
			"info",
		)
	}

	return err
}

// ForceUnlockStateWithEvent forcefully unlocks state (for admin use)
func (st *StateTracker) ForceUnlockStateWithEvent(stackID, unlockedBy string) error {
	// For force unlock, we can delete all locks for the stack directly
	// This bypasses the normal lock_id check
	lockInfo, _ := st.backend.GetLockInfo(stackID)

	var lockID string
	if lockInfo != nil {
		lockID = lockInfo.LockID
	} else {
		lockID = "force-unlock"
	}

	// Attempt to unlock with any lock ID
	err := st.backend.UnlockState(stackID, lockID)

	if err == nil {
		stack, _ := st.backend.GetStackManager().GetStack(stackID)
		stackName := ""
		if stack != nil {
			stackName = stack.Name
		}

		st.emitEvent(
			EventStateUnlocked,
			stackID,
			stackName,
			unlockedBy,
			map[string]interface{}{
				"forced":  true,
				"lock_id": lockID,
			},
			"warning",
		)
	}

	return err
}

// ValidateState validates a stack state
func (st *StateTracker) ValidateState(stackID string) (bool, []string, error) {
	stack, err := st.backend.GetStackManager().GetStack(stackID)
	if err != nil {
		return false, []string{fmt.Sprintf("stack not found: %v", err)}, err
	}

	issues := make([]string, 0)

	// Validate resources exist
	resources, err := st.backend.GetStackManager().ListResources(stackID)
	if err != nil {
		issues = append(issues, fmt.Sprintf("failed to list resources: %v", err))
	}

	// Check for orphaned dependencies
	for _, resource := range resources {
		for _, depID := range resource.Dependencies {
			found := false
			for _, r := range resources {
				if r.ID == depID {
					found = true
					break
				}
			}
			if !found {
				issues = append(issues, fmt.Sprintf("resource %s has orphaned dependency: %s", resource.Name, depID))
			}
		}
	}

	// Validate stack metadata
	if stack.Name == "" {
		issues = append(issues, "stack has no name")
	}
	if stack.Version == "" {
		issues = append(issues, "stack has no version")
	}

	valid := len(issues) == 0

	// Emit event
	if valid {
		st.emitEvent(
			EventValidationPassed,
			stackID,
			stack.Name,
			"validator",
			map[string]interface{}{
				"resources_count": len(resources),
			},
			"info",
		)
	} else {
		st.emitEvent(
			EventValidationFailed,
			stackID,
			stack.Name,
			"validator",
			map[string]interface{}{
				"issues_count": len(issues),
				"issues":       issues,
			},
			"error",
		)
	}

	return valid, issues, nil
}

// GetEventHistory returns recent events from database
func (st *StateTracker) GetEventHistory(limit int) []*StateEvent {
	events := make([]*StateEvent, 0)

	rows, err := st.backend.GetStackManager().db.Query(`
		SELECT id, event_type, timestamp, source, stack_id, stack_name, data, severity
		FROM state_events
		ORDER BY timestamp DESC
		LIMIT ?
	`, limit)
	if err != nil {
		slog.Warn("Failed to query event history", "error", err)
		return events
	}
	defer rows.Close()

	for rows.Next() {
		var event StateEvent
		var dataJSON string
		err := rows.Scan(&event.ID, &event.Type, &event.Timestamp, &event.Source, &event.StackID, &event.StackName, &dataJSON, &event.Severity)
		if err != nil {
			slog.Warn("Failed to scan event", "error", err)
			continue
		}

		// Unmarshal data
		if dataJSON != "" {
			if err := json.Unmarshal([]byte(dataJSON), &event.Data); err != nil {
				slog.Warn("Failed to unmarshal event data", "error", err)
			}
		}

		events = append(events, &event)
	}

	return events
}

// GetEventsByType returns events filtered by type from database
func (st *StateTracker) GetEventsByType(eventType EventType, limit int) []*StateEvent {
	events := make([]*StateEvent, 0)

	rows, err := st.backend.GetStackManager().db.Query(`
		SELECT id, event_type, timestamp, source, stack_id, stack_name, data, severity
		FROM state_events
		WHERE event_type = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`, string(eventType), limit)
	if err != nil {
		slog.Warn("Failed to query events by type", "error", err)
		return events
	}
	defer rows.Close()

	for rows.Next() {
		var event StateEvent
		var dataJSON string
		err := rows.Scan(&event.ID, &event.Type, &event.Timestamp, &event.Source, &event.StackID, &event.StackName, &dataJSON, &event.Severity)
		if err != nil {
			slog.Warn("Failed to scan event", "error", err)
			continue
		}

		// Unmarshal data
		if dataJSON != "" {
			if err := json.Unmarshal([]byte(dataJSON), &event.Data); err != nil {
				slog.Warn("Failed to unmarshal event data", "error", err)
			}
		}

		events = append(events, &event)
	}

	return events
}

// GetEventsByStack returns events for a specific stack from database
func (st *StateTracker) GetEventsByStack(stackID string, limit int) []*StateEvent {
	events := make([]*StateEvent, 0)

	rows, err := st.backend.GetStackManager().db.Query(`
		SELECT id, event_type, timestamp, source, stack_id, stack_name, data, severity
		FROM state_events
		WHERE stack_id = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`, stackID, limit)
	if err != nil {
		slog.Warn("Failed to query events by stack", "error", err)
		return events
	}
	defer rows.Close()

	for rows.Next() {
		var event StateEvent
		var dataJSON string
		err := rows.Scan(&event.ID, &event.Type, &event.Timestamp, &event.Source, &event.StackID, &event.StackName, &dataJSON, &event.Severity)
		if err != nil {
			slog.Warn("Failed to scan event", "error", err)
			continue
		}

		// Unmarshal data
		if dataJSON != "" {
			if err := json.Unmarshal([]byte(dataJSON), &event.Data); err != nil {
				slog.Warn("Failed to unmarshal event data", "error", err)
			}
		}

		events = append(events, &event)
	}

	return events
}

// GetEventStats returns event statistics from database
func (st *StateTracker) GetEventStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// Count total events
	var totalEvents int
	err := st.backend.GetStackManager().db.QueryRow("SELECT COUNT(*) FROM state_events").Scan(&totalEvents)
	if err != nil {
		slog.Warn("Failed to count total events", "error", err)
		totalEvents = 0
	}
	stats["total_events"] = totalEvents

	// Count by type
	byType := make(map[EventType]int)
	rows, err := st.backend.GetStackManager().db.Query(`
		SELECT event_type, COUNT(*) as count
		FROM state_events
		GROUP BY event_type
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var eventType string
			var count int
			if err := rows.Scan(&eventType, &count); err == nil {
				byType[EventType(eventType)] = count
			}
		}
	}
	stats["by_type"] = byType

	// Count by severity
	bySeverity := make(map[string]int)
	rows2, err := st.backend.GetStackManager().db.Query(`
		SELECT severity, COUNT(*) as count
		FROM state_events
		GROUP BY severity
	`)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var severity string
			var count int
			if err := rows2.Scan(&severity, &count); err == nil {
				bySeverity[severity] = count
			}
		}
	}
	stats["by_severity"] = bySeverity

	return stats
}

// persistEvent persists an event to the database
func (st *StateTracker) persistEvent(event *StateEvent) error {
	// Debug log
	slog.Info("Persisting event to database",
		"event_id", event.ID,
		"event_type", event.Type,
		"stack_name", event.StackName,
	)

	// Serialize event data to JSON
	dataJSON, err := json.Marshal(event.Data)
	if err != nil {
		slog.Error("Failed to marshal event data", "error", err)
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Insert into state_events table
	result, err := st.backend.GetStackManager().db.Exec(`
		INSERT OR IGNORE INTO state_events (id, event_type, timestamp, source, stack_id, stack_name, data, severity)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, event.ID, string(event.Type), event.Timestamp, event.Source, event.StackID, event.StackName, string(dataJSON), event.Severity)

	if err != nil {
		slog.Error("Failed to insert event", "error", err)
		return fmt.Errorf("failed to persist event: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	slog.Info("Event persisted successfully", "rows_affected", rowsAffected)

	return nil
}
