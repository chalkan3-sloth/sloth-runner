//go:build cgo
// +build cgo

package services

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

// TestStackServiceIntegration tests the complete integration of stack services
func TestStackServiceIntegration(t *testing.T) {
	// Create temporary database for testing
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_stacks.db")

	// Set environment to use test database
	os.Setenv("SLOTH_RUNNER_DB_PATH", dbPath)
	defer os.Unsetenv("SLOTH_RUNNER_DB_PATH")

	// Create stack service
	stackService, err := NewStackService()
	if err != nil {
		t.Fatalf("Failed to create stack service: %v", err)
	}
	defer stackService.Close()

	// Test 1: Create stack
	t.Run("CreateStack", func(t *testing.T) {
		stackID, err := stackService.GetOrCreateStack("test-stack", "test-workflow", "/tmp/test.sloth")
		if err != nil {
			t.Fatalf("Failed to create stack: %v", err)
		}

		if stackID == "" {
			t.Error("Stack ID should not be empty")
		}
	})

	// Test 2: Create snapshot
	t.Run("CreateSnapshot", func(t *testing.T) {
		// Get the stack we just created
		testStack, err := stackService.GetManager().GetStackByName("test-stack")
		if err != nil {
			t.Fatalf("Failed to get test stack: %v", err)
		}

		version, err := stackService.CreateSnapshot(testStack.ID, "test-user", "Test snapshot")
		if err != nil {
			t.Fatalf("Failed to create snapshot: %v", err)
		}

		if version != 1 {
			t.Errorf("Expected version 1, got %d", version)
		}
	})

	// Test 3: List snapshots
	t.Run("ListSnapshots", func(t *testing.T) {
		testStack, err := stackService.GetManager().GetStackByName("test-stack")
		if err != nil {
			t.Fatalf("Failed to get test stack: %v", err)
		}

		snapshots, err := stackService.ListSnapshots(testStack.ID)
		if err != nil {
			t.Fatalf("Failed to list snapshots: %v", err)
		}

		if len(snapshots) != 1 {
			t.Errorf("Expected 1 snapshot, got %d", len(snapshots))
		}
	})

	// Test 4: Drift detection
	t.Run("DetectDrift", func(t *testing.T) {
		testStack, err := stackService.GetManager().GetStackByName("test-stack")
		if err != nil {
			t.Fatalf("Failed to get test stack: %v", err)
		}

		// Create a resource
		resource := &stack.Resource{
			ID:      "test-resource-1",
			StackID: testStack.ID,
			Type:    "test",
			Name:    "test-resource",
			Module:  "test",
			Properties: map[string]interface{}{
				"key1": "value1",
			},
			State:    "applied",
			Checksum: "test-checksum",
		}

		err = stackService.GetManager().CreateResource(resource)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		// Detect drift
		expectedState := map[string]interface{}{
			"key1": "value1",
		}
		actualState := map[string]interface{}{
			"key1": "value2", // Drifted value
		}

		err = stackService.DetectDrift(testStack.ID, resource.ID, expectedState, actualState)
		if err != nil {
			t.Fatalf("Failed to detect drift: %v", err)
		}

		// Get drift info
		drifts, err := stackService.GetDriftInfo(testStack.ID)
		if err != nil {
			t.Fatalf("Failed to get drift info: %v", err)
		}

		if len(drifts) != 1 {
			t.Errorf("Expected 1 drift, got %d", len(drifts))
		}

		if !drifts[0].IsDrifted {
			t.Error("Expected drift to be detected")
		}
	})

	// Test 5: Tags
	t.Run("Tags", func(t *testing.T) {
		testStack, err := stackService.GetManager().GetStackByName("test-stack")
		if err != nil {
			t.Fatalf("Failed to get test stack: %v", err)
		}

		// Add tags
		err = stackService.AddTag(testStack.ID, "production")
		if err != nil {
			t.Fatalf("Failed to add tag: %v", err)
		}

		err = stackService.AddTag(testStack.ID, "critical")
		if err != nil {
			t.Fatalf("Failed to add tag: %v", err)
		}

		// Get tags
		tags, err := stackService.GetTags(testStack.ID)
		if err != nil {
			t.Fatalf("Failed to get tags: %v", err)
		}

		if len(tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(tags))
		}
	})

	// Test 6: Activity log
	t.Run("ActivityLog", func(t *testing.T) {
		testStack, err := stackService.GetManager().GetStackByName("test-stack")
		if err != nil {
			t.Fatalf("Failed to get test stack: %v", err)
		}

		activities, err := stackService.GetActivity(testStack.ID, 10)
		if err != nil {
			t.Fatalf("Failed to get activity: %v", err)
		}

		// Should have at least snapshot activity
		if len(activities) == 0 {
			t.Error("Expected some activity, got none")
		}
	})
}

// TestSysadminStackIntegration tests sysadmin operations tracking
func TestSysadminStackIntegration(t *testing.T) {
	// Create temporary database for testing
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_sysadmin.db")

	os.Setenv("SLOTH_RUNNER_DB_PATH", dbPath)
	defer os.Unsetenv("SLOTH_RUNNER_DB_PATH")

	// Create services
	stackService, err := NewStackService()
	if err != nil {
		t.Fatalf("Failed to create stack service: %v", err)
	}
	defer stackService.Close()

	sysadminService := NewSysadminStackService(stackService)

	// Test backup operation tracking
	t.Run("TrackBackupOperation", func(t *testing.T) {
		err := sysadminService.TrackBackupOperation(
			"backup-001",
			"Test backup",
			[]string{"/etc", "/var"},
			1024*1024, // 1MB
		)
		if err != nil {
			t.Fatalf("Failed to track backup operation: %v", err)
		}

		// Verify the operation was tracked
		history, err := sysadminService.GetSysadminOperationHistory()
		if err != nil {
			t.Fatalf("Failed to get operation history: %v", err)
		}

		if len(history) != 1 {
			t.Errorf("Expected 1 operation, got %d", len(history))
		}
	})

	// Test deployment operation tracking
	t.Run("TrackDeploymentOperation", func(t *testing.T) {
		err := sysadminService.TrackDeploymentOperation(
			"v1.0.0",
			[]string{"agent-1", "agent-2"},
			"rolling",
			true,
			5*time.Minute,
		)
		if err != nil {
			t.Fatalf("Failed to track deployment operation: %v", err)
		}

		// Verify the operation was tracked
		history, err := sysadminService.GetSysadminOperationHistory()
		if err != nil {
			t.Fatalf("Failed to get operation history: %v", err)
		}

		if len(history) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(history))
		}
	})

	// Test maintenance operation tracking
	t.Run("TrackMaintenanceOperation", func(t *testing.T) {
		err := sysadminService.TrackMaintenanceOperation(
			"cleanup",
			"agent-1",
			"Cleanup old logs",
			true,
		)
		if err != nil {
			t.Fatalf("Failed to track maintenance operation: %v", err)
		}

		// Verify the operation was tracked
		history, err := sysadminService.GetSysadminOperationHistory()
		if err != nil {
			t.Fatalf("Failed to get operation history: %v", err)
		}

		if len(history) != 3 {
			t.Errorf("Expected 3 operations, got %d", len(history))
		}
	})
}

// TestSchedulerStackIntegration tests scheduler operations tracking
func TestSchedulerStackIntegration(t *testing.T) {
	// Create temporary database for testing
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_scheduler.db")

	os.Setenv("SLOTH_RUNNER_DB_PATH", dbPath)
	defer os.Unsetenv("SLOTH_RUNNER_DB_PATH")

	// Create services
	stackService, err := NewStackService()
	if err != nil {
		t.Fatalf("Failed to create stack service: %v", err)
	}
	defer stackService.Close()

	schedulerService := NewSchedulerStackService(stackService)

	// Test scheduled execution tracking
	t.Run("TrackScheduledExecution", func(t *testing.T) {
		err := schedulerService.TrackScheduledExecution(
			"daily-backup",
			"0 0 * * *",
			"exec-001",
			true,
			2*time.Minute,
			"",
		)
		if err != nil {
			t.Fatalf("Failed to track scheduled execution: %v", err)
		}

		// Get history
		history, err := schedulerService.GetSchedulerHistory(10)
		if err != nil {
			t.Fatalf("Failed to get scheduler history: %v", err)
		}

		if len(history) != 1 {
			t.Errorf("Expected 1 execution, got %d", len(history))
		}
	})

	// Test schedule change tracking
	t.Run("TrackScheduleChange", func(t *testing.T) {
		err := schedulerService.TrackScheduleChange(
			"daily-backup",
			"0 0 * * *",
			"enabled",
		)
		if err != nil {
			t.Fatalf("Failed to track schedule change: %v", err)
		}

		// Get history
		history, err := schedulerService.GetSchedulerHistory(10)
		if err != nil {
			t.Fatalf("Failed to get scheduler history: %v", err)
		}

		if len(history) != 2 {
			t.Errorf("Expected 2 records, got %d", len(history))
		}
	})

	// Test execution statistics
	t.Run("GetScheduleExecutionStats", func(t *testing.T) {
		// Add more executions
		schedulerService.TrackScheduledExecution(
			"daily-backup",
			"0 0 * * *",
			"exec-002",
			true,
			3*time.Minute,
			"",
		)

		schedulerService.TrackScheduledExecution(
			"weekly-report",
			"0 0 * * 0",
			"exec-003",
			false,
			1*time.Minute,
			"Connection timeout",
		)

		// Get stats
		stats, err := schedulerService.GetScheduleExecutionStats()
		if err != nil {
			t.Fatalf("Failed to get execution stats: %v", err)
		}

		if stats["total_executions"].(int) != 3 {
			t.Errorf("Expected 3 total executions, got %d", stats["total_executions"])
		}

		if stats["successful"].(int) != 2 {
			t.Errorf("Expected 2 successful executions, got %d", stats["successful"])
		}

		if stats["failed"].(int) != 1 {
			t.Errorf("Expected 1 failed execution, got %d", stats["failed"])
		}
	})
}
