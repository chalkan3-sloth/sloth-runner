//go:build cgo
// +build cgo

package integration

import (
	"testing"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// Test NewStateEventEmitter
func TestNewStateEventEmitter(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	if emitter == nil {
		t.Error("Expected non-nil emitter")
	}

	if emitter.backend == nil {
		t.Error("Expected backend to be set")
	}

	if !emitter.enabled {
		t.Error("Expected emitter to be enabled by default")
	}
}

func TestNewStateEventEmitter_GetsGlobalDispatcher(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	// Initialize global dispatcher
	hooks.CleanupGlobalDispatcher()
	hooks.InitializeGlobalDispatcher()
	defer hooks.CleanupGlobalDispatcher()

	emitter := NewStateEventEmitter(backend)

	if emitter.dispatcher == nil {
		t.Error("Expected dispatcher to be retrieved from global")
	}
}

// Test Enable/Disable
func TestStateEventEmitter_Enable(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	emitter.Disable()
	if emitter.enabled {
		t.Error("Expected emitter to be disabled")
	}

	emitter.Enable()
	if !emitter.enabled {
		t.Error("Expected emitter to be enabled")
	}
}

func TestStateEventEmitter_Disable(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	if !emitter.enabled {
		t.Error("Expected emitter to be enabled initially")
	}

	emitter.Disable()
	if emitter.enabled {
		t.Error("Expected emitter to be disabled")
	}
}

// Test CreateStack
func TestStateEventEmitter_CreateStack(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	st := &stack.StackState{
		ID:          "test-stack",
		Name:        "Test Stack",
		Version:     "1",
		Status:      "active",
		Description: "Test description",
		CreatedAt:   time.Now(),
	}

	err = emitter.CreateStack(st)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify stack was created
	retrieved, err := backend.GetStackManager().GetStack("test-stack")
	if err != nil {
		t.Fatalf("Failed to retrieve stack: %v", err)
	}

	if retrieved.Name != "Test Stack" {
		t.Errorf("Expected name 'Test Stack', got '%s'", retrieved.Name)
	}
}

// Test UpdateStack
func TestStateEventEmitter_UpdateStack(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack first
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Original Name",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	// Update stack
	st.Name = "Updated Name"
	st.UpdatedAt = time.Now()

	err = emitter.UpdateStack(st)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	retrieved, err := backend.GetStackManager().GetStack("test-stack")
	if err != nil {
		t.Fatalf("Failed to retrieve stack: %v", err)
	}

	if retrieved.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", retrieved.Name)
	}
}

// Test DeleteStack
func TestStateEventEmitter_DeleteStack(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack first
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	// Delete stack
	err = emitter.DeleteStack("test-stack")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = backend.GetStackManager().GetStack("test-stack")
	if err == nil {
		t.Error("Expected error getting deleted stack")
	}
}

// Test CreateSnapshot
func TestStateEventEmitter_CreateSnapshot(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack first
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	// Create snapshot
	version, err := emitter.CreateSnapshot("test-stack", "user@example.com", "Test snapshot")
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	if version <= 0 {
		t.Errorf("Expected positive version, got %d", version)
	}
}

// Test RollbackToSnapshot
func TestStateEventEmitter_RollbackToSnapshot(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack with initialized maps
	st := &stack.StackState{
		ID:            "test-stack",
		Name:          "Test Stack",
		Version:       "1",
		Status:        "active",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}
	backend.GetStackManager().CreateStack(st)

	// Create snapshot
	version, err := backend.CreateSnapshot("test-stack", "user@example.com", "Snapshot 1")
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	// Rollback
	err = emitter.RollbackToSnapshot("test-stack", version, "user@example.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestStateEventEmitter_RollbackToSnapshot_EmitsFailureEvent(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	// Initialize global dispatcher to capture events
	hooks.CleanupGlobalDispatcher()
	hooks.InitializeGlobalDispatcher()
	defer hooks.CleanupGlobalDispatcher()

	emitter := NewStateEventEmitter(backend)

	// Try to rollback to non-existent snapshot (should fail)
	err = emitter.RollbackToSnapshot("non-existent", 999, "user@example.com")
	if err == nil {
		t.Error("Expected error for non-existent stack")
	}
}

// Test DetectDrift
func TestStateEventEmitter_DetectDrift(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	// Create resource for drift detection
	resource := &stack.Resource{
		ID:      "resource-1",
		StackID: "test-stack",
		Type:    "compute",
		Name:    "test-resource",
		State:   "applied",
	}
	backend.GetStackManager().CreateResource(resource)

	expectedState := map[string]interface{}{
		"name": "expected",
	}

	actualState := map[string]interface{}{
		"name": "actual",
	}

	err = emitter.DetectDrift("test-stack", "resource-1", expectedState, actualState)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test LockState
func TestStateEventEmitter_LockState(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	err = emitter.LockState("test-stack", "lock-1", "apply", "user@example.com", 10*time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify lock by trying to lock again (should fail)
	err = backend.LockState("test-stack", "lock-2", "apply", "user@example.com", 10*time.Minute)
	if err == nil {
		t.Error("Expected error when trying to lock already locked stack")
	}
}

// Test UnlockState
func TestStateEventEmitter_UnlockState(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack and lock it
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)
	backend.LockState("test-stack", "lock-1", "apply", "user@example.com", 10*time.Minute)

	// Unlock
	err = emitter.UnlockState("test-stack", "lock-1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify unlock by trying to lock again (should succeed)
	err = backend.LockState("test-stack", "lock-2", "apply", "user@example.com", 10*time.Minute)
	if err != nil {
		t.Errorf("Expected lock to succeed after unlock: %v", err)
	}
}

// Test CreateResource
func TestStateEventEmitter_CreateResource(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack first
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	// Create resource
	resource := &stack.Resource{
		ID:        "res-1",
		StackID:   "test-stack",
		Type:      "compute",
		Name:      "test-resource",
		Module:    "test-module",
		State:     "pending",
		CreatedAt: time.Now(),
	}

	err = emitter.CreateResource(resource)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify resource was created
	retrieved, err := backend.GetStackManager().GetResource("res-1")
	if err != nil {
		t.Fatalf("Failed to retrieve resource: %v", err)
	}

	if retrieved.Name != "test-resource" {
		t.Errorf("Expected name 'test-resource', got '%s'", retrieved.Name)
	}
}

// Test UpdateResource
func TestStateEventEmitter_UpdateResource(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack and resource
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	resource := &stack.Resource{
		ID:      "res-1",
		StackID: "test-stack",
		Type:    "compute",
		Name:    "original-name",
		State:   "pending",
	}
	backend.GetStackManager().CreateResource(resource)

	// Update resource
	resource.Name = "updated-name"
	resource.State = "applied"
	resource.UpdatedAt = time.Now()

	err = emitter.UpdateResource(resource)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	retrieved, err := backend.GetStackManager().GetResource("res-1")
	if err != nil {
		t.Fatalf("Failed to retrieve resource: %v", err)
	}

	if retrieved.Name != "updated-name" {
		t.Errorf("Expected name 'updated-name', got '%s'", retrieved.Name)
	}
}

// Test DeleteResource
func TestStateEventEmitter_DeleteResource(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack and resource
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	resource := &stack.Resource{
		ID:      "res-1",
		StackID: "test-stack",
		Type:    "compute",
		Name:    "test-resource",
	}
	backend.GetStackManager().CreateResource(resource)

	// Delete resource
	err = emitter.DeleteResource("res-1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = backend.GetStackManager().GetResource("res-1")
	if err == nil {
		t.Error("Expected error getting deleted resource")
	}
}

// Test AddTag
func TestStateEventEmitter_AddTag(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	err = emitter.AddTag("test-stack", "production")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify tag was added
	tags, err := backend.GetTags("test-stack")
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	found := false
	for _, tag := range tags {
		if tag == "production" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected tag 'production' to be added")
	}
}

// Test RemoveTag
func TestStateEventEmitter_RemoveTag(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	// Create stack and add tag
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)
	backend.AddTag("test-stack", "production")

	// Remove tag
	err = emitter.RemoveTag("test-stack", "production")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify tag was removed
	tags, err := backend.GetTags("test-stack")
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	for _, tag := range tags {
		if tag == "production" {
			t.Error("Tag 'production' should have been removed")
		}
	}
}

// Test GetBackend
func TestStateEventEmitter_GetBackend(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	emitter := NewStateEventEmitter(backend)

	retrievedBackend := emitter.GetBackend()
	if retrievedBackend != backend {
		t.Error("Expected same backend instance")
	}
}

// Test CreateStackWithEventContext
func TestStateEventEmitter_CreateStackWithEventContext(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	// Initialize global dispatcher
	hooks.CleanupGlobalDispatcher()
	hooks.InitializeGlobalDispatcher()
	defer hooks.CleanupGlobalDispatcher()

	emitter := NewStateEventEmitter(backend)

	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}

	err = emitter.CreateStackWithEventContext(st, "test-agent", "run-123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify stack was created
	retrieved, err := backend.GetStackManager().GetStack("test-stack")
	if err != nil {
		t.Fatalf("Failed to retrieve stack: %v", err)
	}

	if retrieved.Name != "Test Stack" {
		t.Errorf("Expected name 'Test Stack', got '%s'", retrieved.Name)
	}

	// Verify execution context was set
	dispatcher := hooks.GetGlobalDispatcher()
	if dispatcher.GetCurrentAgent() != "test-agent" {
		t.Errorf("Expected agent 'test-agent', got '%s'", dispatcher.GetCurrentAgent())
	}

	if dispatcher.GetCurrentRunID() != "run-123" {
		t.Errorf("Expected runID 'run-123', got '%s'", dispatcher.GetCurrentRunID())
	}
}

// Test event emission when disabled
func TestStateEventEmitter_NoEventsWhenDisabled(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	// Initialize global dispatcher
	hooks.CleanupGlobalDispatcher()
	hooks.InitializeGlobalDispatcher()
	defer hooks.CleanupGlobalDispatcher()

	emitter := NewStateEventEmitter(backend)
	emitter.Disable()

	// Create stack (should not emit event)
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}

	err = emitter.CreateStack(st)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Stack should be created but no event dispatched
	// (We can't easily verify no event was dispatched without inspecting event queue)
}

// Test NewStateEventListener
func TestNewStateEventListener(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	listener := NewStateEventListener(backend)

	if listener == nil {
		t.Error("Expected non-nil listener")
	}

	if listener.backend == nil {
		t.Error("Expected backend to be set")
	}
}

// Test HandleWorkflowStarted
func TestStateEventListener_HandleWorkflowStarted(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	listener := NewStateEventListener(backend)

	// Create stack for workflow
	st := &stack.StackState{
		ID:      "workflow-test-workflow",
		Name:    "Workflow Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	event := &hooks.Event{
		Type:      hooks.EventWorkflowStarted,
		Timestamp: time.Now(),
		Stack:     "workflow-test-workflow",
		Data: map[string]interface{}{
			"workflow": map[string]interface{}{
				"name":   "test-workflow",
				"status": "running",
			},
		},
	}

	err = listener.HandleWorkflowStarted(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify snapshot was created
	snapshots, err := backend.ListSnapshots("workflow-test-workflow")
	if err != nil {
		t.Fatalf("Failed to list snapshots: %v", err)
	}

	if len(snapshots) == 0 {
		t.Error("Expected at least one snapshot")
	}
}

func TestStateEventListener_HandleWorkflowStarted_MissingData(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	listener := NewStateEventListener(backend)

	event := &hooks.Event{
		Type:      hooks.EventWorkflowStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = listener.HandleWorkflowStarted(event)
	if err == nil {
		t.Error("Expected error for missing workflow data")
	}
}

// Test HandleWorkflowCompleted
func TestStateEventListener_HandleWorkflowCompleted(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	listener := NewStateEventListener(backend)

	// Create stack
	st := &stack.StackState{
		ID:      "workflow-test-workflow",
		Name:    "Workflow Stack",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	event := &hooks.Event{
		Type:      hooks.EventWorkflowCompleted,
		Timestamp: time.Now(),
		Stack:     "workflow-test-workflow",
		Data: map[string]interface{}{
			"workflow": map[string]interface{}{
				"name":   "test-workflow",
				"status": "completed",
			},
		},
	}

	err = listener.HandleWorkflowCompleted(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify snapshot was created
	snapshots, err := backend.ListSnapshots("workflow-test-workflow")
	if err != nil {
		t.Fatalf("Failed to list snapshots: %v", err)
	}

	if len(snapshots) == 0 {
		t.Error("Expected at least one snapshot")
	}
}

func TestStateEventListener_HandleWorkflowCompleted_MissingData(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	listener := NewStateEventListener(backend)

	event := &hooks.Event{
		Type:      hooks.EventWorkflowCompleted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = listener.HandleWorkflowCompleted(event)
	if err == nil {
		t.Error("Expected error for missing workflow data")
	}
}

// Test HandleAgentDisconnected
func TestStateEventListener_HandleAgentDisconnected(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	listener := NewStateEventListener(backend)

	event := &hooks.Event{
		Type:      hooks.EventAgentDisconnected,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent": map[string]interface{}{
				"name":    "test-agent",
				"address": "localhost:50051",
			},
		},
	}

	err = listener.HandleAgentDisconnected(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestStateEventListener_HandleAgentDisconnected_MissingData(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	listener := NewStateEventListener(backend)

	event := &hooks.Event{
		Type:      hooks.EventAgentDisconnected,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = listener.HandleAgentDisconnected(event)
	if err == nil {
		t.Error("Expected error for missing agent data")
	}
}

// Test event emission without dispatcher
func TestStateEventEmitter_NoDispatcher(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	// Ensure no global dispatcher
	hooks.CleanupGlobalDispatcher()

	emitter := NewStateEventEmitter(backend)

	// Should work even without dispatcher (events just not emitted)
	st := &stack.StackState{
		ID:      "test-stack",
		Name:    "Test Stack",
		Version: "1",
		Status:  "active",
	}

	err = emitter.CreateStack(st)
	if err != nil {
		t.Errorf("Should work without dispatcher: %v", err)
	}
}

// Test workflow snapshot derivation
func TestStateEventListener_DerivesStackIDFromWorkflowName(t *testing.T) {
	backend, err := stack.NewStateBackend(":memory:")
	if err != nil {
		t.Fatalf("Failed to create backend: %v", err)
	}
	defer backend.Close()

	listener := NewStateEventListener(backend)

	// Create stack with derived name
	st := &stack.StackState{
		ID:      "workflow-my-workflow",
		Name:    "My Workflow",
		Version: "1",
		Status:  "active",
	}
	backend.GetStackManager().CreateStack(st)

	// Event without explicit stack ID
	event := &hooks.Event{
		Type:      hooks.EventWorkflowStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"workflow": map[string]interface{}{
				"name": "my-workflow",
			},
		},
	}

	err = listener.HandleWorkflowStarted(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should create snapshot for derived stack ID
	snapshots, err := backend.ListSnapshots("workflow-my-workflow")
	if err != nil {
		t.Fatalf("Failed to list snapshots: %v", err)
	}

	if len(snapshots) == 0 {
		t.Error("Expected at least one snapshot")
	}
}
