//go:build cgo
// +build cgo

package integration

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// StateEventEmitter wraps StateBackend to emit events for all state operations
type StateEventEmitter struct {
	backend    *stack.StateBackend
	dispatcher *hooks.Dispatcher
	enabled    bool
}

// NewStateEventEmitter creates a new state event emitter
func NewStateEventEmitter(backend *stack.StateBackend) *StateEventEmitter {
	return &StateEventEmitter{
		backend:    backend,
		dispatcher: hooks.GetGlobalDispatcher(),
		enabled:    true,
	}
}

// Enable enables event emission
func (see *StateEventEmitter) Enable() {
	see.enabled = true
}

// Disable disables event emission
func (see *StateEventEmitter) Disable() {
	see.enabled = false
}

// emitEvent emits an event if dispatcher is available and events are enabled
func (see *StateEventEmitter) emitEvent(eventType hooks.EventType, stackID string, data map[string]interface{}) {
	if !see.enabled || see.dispatcher == nil {
		return
	}

	event := &hooks.Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Stack:     stackID,
		Data:      data,
	}

	if err := see.dispatcher.Dispatch(event); err != nil {
		slog.Warn("failed to dispatch state event",
			"event_type", eventType,
			"stack_id", stackID,
			"error", err)
	}
}

// CreateStack creates a stack and emits a stack.deployed event
func (see *StateEventEmitter) CreateStack(st *stack.StackState) error {
	if err := see.backend.GetStackManager().CreateStack(st); err != nil {
		return err
	}

	see.emitEvent(hooks.EventStackDeployed, st.ID, map[string]interface{}{
		"stack": map[string]interface{}{
			"id":          st.ID,
			"name":        st.Name,
			"version":     st.Version,
			"status":      st.Status,
			"description": st.Description,
			"created_at":  st.CreatedAt,
		},
	})

	return nil
}

// UpdateStack updates a stack and emits a stack.updated event
func (see *StateEventEmitter) UpdateStack(st *stack.StackState) error {
	if err := see.backend.GetStackManager().UpdateStack(st); err != nil {
		return err
	}

	see.emitEvent(hooks.EventStackUpdated, st.ID, map[string]interface{}{
		"stack": map[string]interface{}{
			"id":          st.ID,
			"name":        st.Name,
			"version":     st.Version,
			"status":      st.Status,
			"updated_at":  st.UpdatedAt,
		},
	})

	return nil
}

// DeleteStack deletes a stack and emits a stack.destroyed event
func (see *StateEventEmitter) DeleteStack(stackID string) error {
	// Get stack info before deletion
	st, err := see.backend.GetStackManager().GetStack(stackID)
	if err != nil {
		return err
	}

	if err := see.backend.GetStackManager().DeleteStack(stackID); err != nil {
		return err
	}

	see.emitEvent(hooks.EventStackDestroyed, stackID, map[string]interface{}{
		"stack": map[string]interface{}{
			"id":      st.ID,
			"name":    st.Name,
			"version": st.Version,
		},
	})

	return nil
}

// CreateSnapshot creates a snapshot and emits a stack.snapshot_created event
func (see *StateEventEmitter) CreateSnapshot(stackID, createdBy, description string) (int, error) {
	version, err := see.backend.CreateSnapshot(stackID, createdBy, description)
	if err != nil {
		return 0, err
	}

	see.emitEvent(hooks.EventStackSnapshotCreated, stackID, map[string]interface{}{
		"snapshot": map[string]interface{}{
			"stack_id":    stackID,
			"version":     version,
			"created_by":  createdBy,
			"description": description,
			"created_at":  time.Now(),
		},
	})

	return version, nil
}

// RollbackToSnapshot rolls back to a snapshot and emits a stack.rollback event
func (see *StateEventEmitter) RollbackToSnapshot(stackID string, version int, performedBy string) error {
	// Get snapshot info before rollback
	snapshot, err := see.backend.GetSnapshot(stackID, version)
	if err != nil {
		return err
	}

	if err := see.backend.RollbackToSnapshot(stackID, version, performedBy); err != nil {
		// Emit rollback failed event
		see.emitEvent(hooks.EventStackRollbackFailed, stackID, map[string]interface{}{
			"rollback": map[string]interface{}{
				"stack_id":       stackID,
				"target_version": version,
				"performed_by":   performedBy,
				"error":          err.Error(),
				"timestamp":      time.Now(),
			},
		})
		return err
	}

	see.emitEvent(hooks.EventStackRolledBack, stackID, map[string]interface{}{
		"rollback": map[string]interface{}{
			"stack_id":     stackID,
			"from_version": snapshot.Version,
			"to_version":   version,
			"performed_by": performedBy,
			"timestamp":    time.Now(),
		},
	})

	return nil
}

// DetectDrift detects drift and emits a stack.drift_detected event if drift is found
func (see *StateEventEmitter) DetectDrift(stackID, resourceID string, expectedState, actualState map[string]interface{}) error {
	if err := see.backend.DetectDrift(stackID, resourceID, expectedState, actualState); err != nil {
		return err
	}

	// Check if drift was actually detected
	drifts, err := see.backend.GetDriftInfo(stackID)
	if err != nil {
		return err
	}

	// Find the specific drift we just detected
	for _, drift := range drifts {
		if drift.ResourceID == resourceID && drift.IsDrifted {
			see.emitEvent(hooks.EventStackDriftDetected, stackID, map[string]interface{}{
				"drift": map[string]interface{}{
					"stack_id":       stackID,
					"resource_id":    resourceID,
					"expected_state": expectedState,
					"actual_state":   actualState,
					"drifted_fields": drift.DriftedFields,
					"detected_at":    time.Now(),
				},
			})
			break
		}
	}

	return nil
}

// LockState locks state and emits a stack.locked event
func (see *StateEventEmitter) LockState(stackID, lockID, operation, who string, duration time.Duration) error {
	if err := see.backend.LockState(stackID, lockID, operation, who, duration); err != nil {
		return err
	}

	see.emitEvent(hooks.EventStackLocked, stackID, map[string]interface{}{
		"lock": map[string]interface{}{
			"stack_id":  stackID,
			"lock_id":   lockID,
			"operation": operation,
			"who":       who,
			"duration":  duration.String(),
			"locked_at": time.Now(),
		},
	})

	return nil
}

// UnlockState unlocks state and emits a stack.unlocked event
func (see *StateEventEmitter) UnlockState(stackID, lockID string) error {
	if err := see.backend.UnlockState(stackID, lockID); err != nil {
		return err
	}

	see.emitEvent(hooks.EventStackUnlocked, stackID, map[string]interface{}{
		"lock": map[string]interface{}{
			"stack_id":    stackID,
			"lock_id":     lockID,
			"unlocked_at": time.Now(),
		},
	})

	return nil
}

// CreateResource creates a resource and emits a resource.created event
func (see *StateEventEmitter) CreateResource(resource *stack.Resource) error {
	if err := see.backend.GetStackManager().CreateResource(resource); err != nil {
		return err
	}

	see.emitEvent(hooks.EventResourceCreated, resource.StackID, map[string]interface{}{
		"resource": map[string]interface{}{
			"id":         resource.ID,
			"stack_id":   resource.StackID,
			"type":       resource.Type,
			"name":       resource.Name,
			"module":     resource.Module,
			"state":      resource.State,
			"created_at": resource.CreatedAt,
		},
	})

	return nil
}

// UpdateResource updates a resource and emits a resource.updated event
func (see *StateEventEmitter) UpdateResource(resource *stack.Resource) error {
	// Get previous state
	oldResource, _ := see.backend.GetStackManager().GetResource(resource.ID)

	if err := see.backend.GetStackManager().UpdateResource(resource); err != nil {
		return err
	}

	eventData := map[string]interface{}{
		"resource": map[string]interface{}{
			"id":         resource.ID,
			"stack_id":   resource.StackID,
			"type":       resource.Type,
			"name":       resource.Name,
			"state":      resource.State,
			"updated_at": resource.UpdatedAt,
		},
	}

	// Add previous state if available
	if oldResource != nil {
		eventData["previous_state"] = oldResource.State
	}

	see.emitEvent(hooks.EventResourceUpdated, resource.StackID, eventData)

	return nil
}

// DeleteResource deletes a resource and emits a resource.deleted event
func (see *StateEventEmitter) DeleteResource(resourceID string) error {
	// Get resource info before deletion
	resource, err := see.backend.GetStackManager().GetResource(resourceID)
	if err != nil {
		return err
	}

	if err := see.backend.GetStackManager().DeleteResource(resourceID); err != nil {
		return err
	}

	see.emitEvent(hooks.EventResourceDeleted, resource.StackID, map[string]interface{}{
		"resource": map[string]interface{}{
			"id":       resource.ID,
			"stack_id": resource.StackID,
			"type":     resource.Type,
			"name":     resource.Name,
		},
	})

	return nil
}

// AddTag adds a tag and emits a stack.tagged event
func (see *StateEventEmitter) AddTag(stackID, tag string) error {
	if err := see.backend.AddTag(stackID, tag); err != nil {
		return err
	}

	see.emitEvent(hooks.EventStackTagged, stackID, map[string]interface{}{
		"tag": map[string]interface{}{
			"stack_id": stackID,
			"tag":      tag,
			"added_at": time.Now(),
		},
	})

	return nil
}

// RemoveTag removes a tag and emits a stack.untagged event
func (see *StateEventEmitter) RemoveTag(stackID, tag string) error {
	if err := see.backend.RemoveTag(stackID, tag); err != nil {
		return err
	}

	see.emitEvent(hooks.EventStackUntagged, stackID, map[string]interface{}{
		"tag": map[string]interface{}{
			"stack_id":   stackID,
			"tag":        tag,
			"removed_at": time.Now(),
		},
	})

	return nil
}

// GetBackend returns the underlying StateBackend
func (see *StateEventEmitter) GetBackend() *stack.StateBackend {
	return see.backend
}

// CreateStackWithEventContext creates a stack with execution context
func (see *StateEventEmitter) CreateStackWithEventContext(st *stack.StackState, agent, runID string) error {
	// Set execution context for events
	if see.dispatcher != nil {
		see.dispatcher.SetExecutionContext(st.ID, agent, runID)
	}

	return see.CreateStack(st)
}

// StateEventListener listens for events and performs actions
type StateEventListener struct {
	backend *stack.StateBackend
}

// NewStateEventListener creates a new state event listener
func NewStateEventListener(backend *stack.StateBackend) *StateEventListener {
	return &StateEventListener{
		backend: backend,
	}
}

// HandleWorkflowStarted handles workflow.started events by creating a pre-execution snapshot
func (sel *StateEventListener) HandleWorkflowStarted(event *hooks.Event) error {
	workflowData, ok := event.Data["workflow"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("workflow data not found in event")
	}

	workflowName, _ := workflowData["name"].(string)

	stackID := event.Stack
	if stackID == "" {
		// Try to derive stack ID from workflow name
		stackID = fmt.Sprintf("workflow-%s", workflowName)
	}

	// Create pre-execution snapshot
	_, err := sel.backend.CreateSnapshot(
		stackID,
		"system",
		fmt.Sprintf("Pre-execution snapshot for workflow: %s", workflowName),
	)

	if err != nil {
		slog.Warn("failed to create pre-execution snapshot",
			"workflow", workflowName,
			"error", err)
	}

	return err
}

// HandleWorkflowCompleted handles workflow.completed events by creating a post-execution snapshot
func (sel *StateEventListener) HandleWorkflowCompleted(event *hooks.Event) error {
	workflowData, ok := event.Data["workflow"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("workflow data not found in event")
	}

	workflowName, _ := workflowData["name"].(string)
	status, _ := workflowData["status"].(string)

	stackID := event.Stack
	if stackID == "" {
		stackID = fmt.Sprintf("workflow-%s", workflowName)
	}

	// Create post-execution snapshot
	_, err := sel.backend.CreateSnapshot(
		stackID,
		"system",
		fmt.Sprintf("Post-execution snapshot: %s (status: %s)", workflowName, status),
	)

	if err != nil {
		slog.Warn("failed to create post-execution snapshot",
			"workflow", workflowName,
			"error", err)
	}

	return err
}

// HandleAgentDisconnected handles agent.disconnected events
func (sel *StateEventListener) HandleAgentDisconnected(event *hooks.Event) error {
	agentData, ok := event.Data["agent"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("agent data not found in event")
	}

	agentName, _ := agentData["name"].(string)

	slog.Info("agent disconnected, resources may require drift detection",
		"agent", agentName)

	return nil
}
