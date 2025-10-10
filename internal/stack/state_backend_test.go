//go:build cgo
// +build cgo

package stack

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStateBackend_CreateSnapshot(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_state.db")

	backend, err := NewStateBackend(dbPath)
	if err != nil {
		t.Fatalf("Failed to create state backend: %v", err)
	}
	defer backend.Close()

	// Create a test stack
	stackID := uuid.New().String()
	stack := &StackState{
		ID:            stackID,
		Name:          "test-stack",
		Description:   "Test stack for snapshot",
		Version:       "1.0.0",
		Status:        "created",
		WorkflowFile:  "test.sloth",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := backend.GetStackManager().CreateStack(stack); err != nil {
		t.Fatalf("Failed to create stack: %v", err)
	}

	// Create a snapshot
	version, err := backend.CreateSnapshot(stackID, "test-user", "Initial snapshot")
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	if version != 1 {
		t.Errorf("Expected version 1, got %d", version)
	}

	// Retrieve the snapshot
	snapshot, err := backend.GetSnapshot(stackID, version)
	if err != nil {
		t.Fatalf("Failed to get snapshot: %v", err)
	}

	if snapshot.Version != version {
		t.Errorf("Expected version %d, got %d", version, snapshot.Version)
	}

	if snapshot.StackState.Name != "test-stack" {
		t.Errorf("Expected stack name 'test-stack', got '%s'", snapshot.StackState.Name)
	}
}

func TestStateBackend_DriftDetection(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_drift.db")

	backend, err := NewStateBackend(dbPath)
	if err != nil {
		t.Fatalf("Failed to create state backend: %v", err)
	}
	defer backend.Close()

	// Create stack and resource
	stackID := uuid.New().String()
	stack := &StackState{
		ID:            stackID,
		Name:          "drift-test-stack",
		Description:   "Test drift detection",
		Version:       "1.0.0",
		Status:        "created",
		WorkflowFile:  "test.sloth",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := backend.GetStackManager().CreateStack(stack); err != nil {
		t.Fatalf("Failed to create stack: %v", err)
	}

	// Create a resource
	resourceID := uuid.New().String()
	resource := &Resource{
		ID:           resourceID,
		StackID:      stackID,
		Type:         "test-resource",
		Name:         "my-resource",
		Module:       "test",
		Properties:   map[string]interface{}{"key": "value"},
		Dependencies: []string{},
		State:        "applied",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Metadata:     make(map[string]interface{}),
	}

	if err := backend.GetStackManager().CreateResource(resource); err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	// Detect drift (no drift)
	expectedState := map[string]interface{}{"key": "value"}
	actualState := map[string]interface{}{"key": "value"}

	if err := backend.DetectDrift(stackID, resourceID, expectedState, actualState); err != nil {
		t.Fatalf("Failed to detect drift: %v", err)
	}

	// Detect drift (with drift)
	actualStateDrifted := map[string]interface{}{"key": "different-value"}

	if err := backend.DetectDrift(stackID, resourceID, expectedState, actualStateDrifted); err != nil {
		t.Fatalf("Failed to detect drift: %v", err)
	}

	// Get drift info
	drifts, err := backend.GetDriftInfo(stackID)
	if err != nil {
		t.Fatalf("Failed to get drift info: %v", err)
	}

	if len(drifts) != 2 {
		t.Errorf("Expected 2 drift detections, got %d", len(drifts))
	}

	// Check that at least one is marked as drifted
	hasDrift := false
	for _, d := range drifts {
		if d.IsDrifted {
			hasDrift = true
			break
		}
	}

	if !hasDrift {
		t.Error("Expected at least one drifted detection")
	}
}

func TestStateBackend_StateLock(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_lock.db")

	backend, err := NewStateBackend(dbPath)
	if err != nil {
		t.Fatalf("Failed to create state backend: %v", err)
	}
	defer backend.Close()

	// Create stack
	stackID := uuid.New().String()
	stack := &StackState{
		ID:            stackID,
		Name:          "lock-test-stack",
		Description:   "Test locking",
		Version:       "1.0.0",
		Status:        "created",
		WorkflowFile:  "test.sloth",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := backend.GetStackManager().CreateStack(stack); err != nil {
		t.Fatalf("Failed to create stack: %v", err)
	}

	// Acquire lock
	lockID := "test-lock-1"
	if err := backend.LockState(stackID, lockID, "test-operation", "test-user", 5*time.Minute); err != nil {
		t.Fatalf("Failed to lock state: %v", err)
	}

	// Try to acquire lock again (should fail)
	lockID2 := "test-lock-2"
	err = backend.LockState(stackID, lockID2, "test-operation", "another-user", 5*time.Minute)
	if err == nil {
		t.Error("Expected lock to fail when already locked")
	}

	// Unlock
	if err := backend.UnlockState(stackID, lockID); err != nil {
		t.Fatalf("Failed to unlock state: %v", err)
	}

	// Now lock should succeed
	if err := backend.LockState(stackID, lockID2, "test-operation", "another-user", 5*time.Minute); err != nil {
		t.Fatalf("Failed to lock state after unlock: %v", err)
	}
}

func TestStateBackend_Rollback(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_rollback.db")

	backend, err := NewStateBackend(dbPath)
	if err != nil {
		t.Fatalf("Failed to create state backend: %v", err)
	}
	defer backend.Close()

	// Create stack
	stackID := uuid.New().String()
	stack := &StackState{
		ID:            stackID,
		Name:          "rollback-test-stack",
		Description:   "Test rollback",
		Version:       "1.0.0",
		Status:        "created",
		WorkflowFile:  "test.sloth",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := backend.GetStackManager().CreateStack(stack); err != nil {
		t.Fatalf("Failed to create stack: %v", err)
	}

	// Create initial snapshot
	version1, err := backend.CreateSnapshot(stackID, "test-user", "Version 1")
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	// Modify stack
	stack.Status = "updated"
	if err := backend.GetStackManager().UpdateStack(stack); err != nil {
		t.Fatalf("Failed to update stack: %v", err)
	}

	// Create second snapshot
	version2, err := backend.CreateSnapshot(stackID, "test-user", "Version 2")
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	if version2 <= version1 {
		t.Errorf("Expected version 2 (%d) > version 1 (%d)", version2, version1)
	}

	// Rollback to version 1
	if err := backend.RollbackToSnapshot(stackID, version1, "test-user"); err != nil {
		t.Fatalf("Failed to rollback: %v", err)
	}

	// Verify rollback
	currentStack, err := backend.GetStackManager().GetStack(stackID)
	if err != nil {
		t.Fatalf("Failed to get stack after rollback: %v", err)
	}

	if currentStack.Status != "rolled_back" {
		t.Errorf("Expected status 'rolled_back', got '%s'", currentStack.Status)
	}
}

func TestStateBackend_Tags(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_tags.db")

	backend, err := NewStateBackend(dbPath)
	if err != nil {
		t.Fatalf("Failed to create state backend: %v", err)
	}
	defer backend.Close()

	// Create stack
	stackID := uuid.New().String()
	stack := &StackState{
		ID:            stackID,
		Name:          "tags-test-stack",
		Description:   "Test tags",
		Version:       "1.0.0",
		Status:        "created",
		WorkflowFile:  "test.sloth",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := backend.GetStackManager().CreateStack(stack); err != nil {
		t.Fatalf("Failed to create stack: %v", err)
	}

	// Add tags
	tags := []string{"production", "critical", "monitored"}
	for _, tag := range tags {
		if err := backend.AddTag(stackID, tag); err != nil {
			t.Fatalf("Failed to add tag '%s': %v", tag, err)
		}
	}

	// Get tags
	retrievedTags, err := backend.GetTags(stackID)
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	if len(retrievedTags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(retrievedTags))
	}

	// Verify tags are sorted
	expectedTags := []string{"critical", "monitored", "production"}
	for i, tag := range retrievedTags {
		if tag != expectedTags[i] {
			t.Errorf("Tag at index %d: expected '%s', got '%s'", i, expectedTags[i], tag)
		}
	}
}

func TestStateBackend_Activity(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_activity.db")

	backend, err := NewStateBackend(dbPath)
	if err != nil {
		t.Fatalf("Failed to create state backend: %v", err)
	}
	defer backend.Close()

	// Create stack
	stackID := uuid.New().String()
	stack := &StackState{
		ID:            stackID,
		Name:          "activity-test-stack",
		Description:   "Test activity",
		Version:       "1.0.0",
		Status:        "created",
		WorkflowFile:  "test.sloth",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := backend.GetStackManager().CreateStack(stack); err != nil {
		t.Fatalf("Failed to create stack: %v", err)
	}

	// Perform various operations that should log activity
	backend.CreateSnapshot(stackID, "test-user", "Test snapshot")
	backend.LockState(stackID, "lock-1", "test", "test-user", 5*time.Minute)
	backend.UnlockState(stackID, "lock-1")

	// Get activity
	activities, err := backend.GetActivity(stackID, 10)
	if err != nil {
		t.Fatalf("Failed to get activity: %v", err)
	}

	if len(activities) == 0 {
		t.Error("Expected some activities to be logged")
	}

	// Verify activity structure
	for _, act := range activities {
		if _, ok := act["type"]; !ok {
			t.Error("Activity missing 'type' field")
		}
		if _, ok := act["created_at"]; !ok {
			t.Error("Activity missing 'created_at' field")
		}
	}
}

func TestStateBackend_ResourceDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_deps.db")

	backend, err := NewStateBackend(dbPath)
	if err != nil {
		t.Fatalf("Failed to create state backend: %v", err)
	}
	defer backend.Close()

	// Create stack and resources
	stackID := uuid.New().String()
	stack := &StackState{
		ID:            stackID,
		Name:          "deps-test-stack",
		Description:   "Test dependencies",
		Version:       "1.0.0",
		Status:        "created",
		WorkflowFile:  "test.sloth",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := backend.GetStackManager().CreateStack(stack); err != nil {
		t.Fatalf("Failed to create stack: %v", err)
	}

	// Create resources
	resource1ID := uuid.New().String()
	resource1 := &Resource{
		ID:           resource1ID,
		StackID:      stackID,
		Type:         "database",
		Name:         "db",
		Module:       "test",
		Properties:   map[string]interface{}{},
		Dependencies: []string{},
		State:        "applied",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Metadata:     make(map[string]interface{}),
	}

	resource2ID := uuid.New().String()
	resource2 := &Resource{
		ID:           resource2ID,
		StackID:      stackID,
		Type:         "application",
		Name:         "app",
		Module:       "test",
		Properties:   map[string]interface{}{},
		Dependencies: []string{resource1ID},
		State:        "applied",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Metadata:     make(map[string]interface{}),
	}

	backend.GetStackManager().CreateResource(resource1)
	backend.GetStackManager().CreateResource(resource2)

	// Add dependency
	if err := backend.AddResourceDependency(resource2ID, resource1ID, "explicit"); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	// Get dependencies
	deps, err := backend.GetResourceDependencies(resource2ID)
	if err != nil {
		t.Fatalf("Failed to get dependencies: %v", err)
	}

	if len(deps) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(deps))
	}

	if len(deps) > 0 && deps[0] != resource1ID {
		t.Errorf("Expected dependency to be '%s', got '%s'", resource1ID, deps[0])
	}
}

// Benchmark tests
func BenchmarkStateBackend_CreateSnapshot(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench_snapshot.db")

	backend, err := NewStateBackend(dbPath)
	if err != nil {
		b.Fatalf("Failed to create state backend: %v", err)
	}
	defer backend.Close()

	// Create test stack
	stackID := uuid.New().String()
	stack := &StackState{
		ID:            stackID,
		Name:          "bench-stack",
		Description:   "Benchmark stack",
		Version:       "1.0.0",
		Status:        "created",
		WorkflowFile:  "test.sloth",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	backend.GetStackManager().CreateStack(stack)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		backend.CreateSnapshot(stackID, "bench-user", "Benchmark snapshot")
	}
}
