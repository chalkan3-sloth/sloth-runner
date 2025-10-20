package hooks

import (
	"testing"
	"time"
)

// Test NewEventQueue
func TestNewEventQueue(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := NewEventQueue(repo.db)

	if eq == nil {
		t.Error("Expected non-nil event queue")
	}

	if eq.db == nil {
		t.Error("Expected db to be set")
	}
}

// Test InitializeSchema
func TestInitializeSchema(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := NewEventQueue(repo.db)

	err = eq.InitializeSchema()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestInitializeSchema_CreatesEventsTable(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := NewEventQueue(repo.db)

	err = eq.InitializeSchema()
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Verify table exists by querying it
	_, err = eq.db.Exec("SELECT COUNT(*) FROM events")
	if err != nil {
		t.Errorf("Events table was not created: %v", err)
	}
}

func TestInitializeSchema_CreatesFileWatchersTable(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := NewEventQueue(repo.db)

	err = eq.InitializeSchema()
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Verify table exists
	_, err = eq.db.Exec("SELECT COUNT(*) FROM file_watchers")
	if err != nil {
		t.Errorf("File watchers table was not created: %v", err)
	}
}

func TestInitializeSchema_CreatesEventHookExecutionsTable(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := NewEventQueue(repo.db)

	err = eq.InitializeSchema()
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Verify table exists
	_, err = eq.db.Exec("SELECT COUNT(*) FROM event_hook_executions")
	if err != nil {
		t.Errorf("Event hook executions table was not created: %v", err)
	}
}

// Test EnqueueEvent
func TestEnqueueEvent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"task": "test",
		},
	}

	err = eq.EnqueueEvent(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify ID was assigned
	if event.ID == "" {
		t.Error("Expected event ID to be assigned")
	}
}

func TestEnqueueEvent_SetsDefaults(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = eq.EnqueueEvent(event)
	if err != nil {
		t.Fatalf("Failed to enqueue: %v", err)
	}

	if event.Status != EventStatusPending {
		t.Errorf("Expected status 'pending', got '%s'", event.Status)
	}

	if event.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestEnqueueEvent_PreservesID(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	customID := "custom-id-123"
	event := &Event{
		ID:        customID,
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = eq.EnqueueEvent(event)
	if err != nil {
		t.Fatalf("Failed to enqueue: %v", err)
	}

	if event.ID != customID {
		t.Errorf("Expected ID '%s', got '%s'", customID, event.ID)
	}
}

// Test GetPendingEvents
func TestGetPendingEvents(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Enqueue some events
	for i := 0; i < 5; i++ {
		event := &Event{
			Type:      EventTaskStarted,
			Timestamp: time.Now(),
			Data:      map[string]interface{}{},
		}
		eq.EnqueueEvent(event)
	}

	events, err := eq.GetPendingEvents(10)
	if err != nil {
		t.Fatalf("Failed to get pending events: %v", err)
	}

	if len(events) < 5 {
		t.Errorf("Expected at least 5 events, got %d", len(events))
	}
}

func TestGetPendingEvents_RespectsLimit(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Enqueue 10 events
	for i := 0; i < 10; i++ {
		event := &Event{
			Type:      EventTaskStarted,
			Timestamp: time.Now(),
			Data:      map[string]interface{}{},
		}
		eq.EnqueueEvent(event)
	}

	events, err := eq.GetPendingEvents(3)
	if err != nil {
		t.Fatalf("Failed to get pending events: %v", err)
	}

	if len(events) > 3 {
		t.Errorf("Expected max 3 events, got %d", len(events))
	}
}

// Test GetStuckProcessingEvents
func TestGetStuckProcessingEvents(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create an event and mark as processing
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Update to processing status with old timestamp
	oldTime := time.Now().Add(-60 * time.Second)
	_, err = eq.db.Exec("UPDATE events SET status = ?, processed_at = ? WHERE id = ?",
		EventStatusProcessing, oldTime.Unix(), event.ID)
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}

	// Get stuck events
	stuckEvents, err := eq.GetStuckProcessingEvents(30, 10)
	if err != nil {
		t.Fatalf("Failed to get stuck events: %v", err)
	}

	if len(stuckEvents) == 0 {
		t.Error("Expected at least one stuck event")
	}
}

func TestGetStuckProcessingEvents_RespectsTimeThreshold(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create a recently processed event
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Update to processing status with recent timestamp
	recentTime := time.Now().Add(-5 * time.Second)
	_, err = eq.db.Exec("UPDATE events SET status = ?, processed_at = ? WHERE id = ?",
		EventStatusProcessing, recentTime.Unix(), event.ID)
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}

	// Should not return recent processing events
	stuckEvents, err := eq.GetStuckProcessingEvents(30, 10)
	if err != nil {
		t.Fatalf("Failed to get stuck events: %v", err)
	}

	// Find our event in results
	found := false
	for _, e := range stuckEvents {
		if e.ID == event.ID {
			found = true
			break
		}
	}

	if found {
		t.Error("Recent processing event should not be considered stuck")
	}
}

// Test ListEvents
func TestListEvents(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Enqueue events
	for i := 0; i < 3; i++ {
		event := &Event{
			Type:      EventTaskStarted,
			Timestamp: time.Now(),
			Data:      map[string]interface{}{},
		}
		eq.EnqueueEvent(event)
	}

	events, err := eq.ListEvents("", "", 10)
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	if len(events) < 3 {
		t.Errorf("Expected at least 3 events, got %d", len(events))
	}
}

func TestListEvents_FilterByType(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Enqueue different event types
	eq.EnqueueEvent(&Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	})

	eq.EnqueueEvent(&Event{
		Type:      EventTaskCompleted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	})

	// Filter by type
	events, err := eq.ListEvents(EventTaskStarted, "", 10)
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	for _, event := range events {
		if event.Type != EventTaskStarted {
			t.Errorf("Expected only TaskStarted events, got %s", event.Type)
		}
	}
}

func TestListEvents_FilterByStatus(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Enqueue events
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Filter by status
	events, err := eq.ListEvents("", EventStatusPending, 10)
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	for _, e := range events {
		if e.Status != EventStatusPending {
			t.Errorf("Expected only pending events, got status %s", e.Status)
		}
	}
}

// Test ListEventsByAgent
func TestListEventsByAgent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Enqueue events for different agents
	eq.EnqueueEvent(&Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
		Agent:     "agent-1",
	})

	eq.EnqueueEvent(&Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
		Agent:     "agent-2",
	})

	// Filter by agent
	events, err := eq.ListEventsByAgent("agent-1", "", "", 10)
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	for _, event := range events {
		if event.Agent != "agent-1" {
			t.Errorf("Expected only agent-1 events, got agent %s", event.Agent)
		}
	}
}

// Test GetEvent
func TestGetEvent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create event
	original := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	eq.EnqueueEvent(original)

	// Retrieve event
	retrieved, err := eq.GetEvent(original.ID)
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}

	if retrieved.ID != original.ID {
		t.Errorf("Expected ID %s, got %s", original.ID, retrieved.ID)
	}

	if retrieved.Type != original.Type {
		t.Errorf("Expected type %s, got %s", original.Type, retrieved.Type)
	}
}

func TestGetEvent_NonExistent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	_, err = eq.GetEvent("non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent event")
	}
}

// Test UpdateEventStatus
func TestUpdateEventStatus(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create event
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Update status
	err = eq.UpdateEventStatus(event.ID, EventStatusCompleted, "")
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	// Verify update
	retrieved, err := eq.GetEvent(event.ID)
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}

	if retrieved.Status != EventStatusCompleted {
		t.Errorf("Expected status 'completed', got '%s'", retrieved.Status)
	}
}

func TestUpdateEventStatus_WithError(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create event
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Update with error
	errorMsg := "test error"
	err = eq.UpdateEventStatus(event.ID, EventStatusFailed, errorMsg)
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	// Verify error was stored
	retrieved, err := eq.GetEvent(event.ID)
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}

	if retrieved.Error != errorMsg {
		t.Errorf("Expected error '%s', got '%s'", errorMsg, retrieved.Error)
	}
}

// Test DeleteEvent
func TestDeleteEvent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create event
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Delete event
	err = eq.DeleteEvent(event.ID)
	if err != nil {
		t.Fatalf("Failed to delete event: %v", err)
	}

	// Verify deletion
	_, err = eq.GetEvent(event.ID)
	if err == nil {
		t.Error("Expected error getting deleted event")
	}
}

// Test CleanupOldEvents
func TestCleanupOldEvents(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create old completed event
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Mark as completed with old timestamp
	oldTime := time.Now().Add(-48 * time.Hour)
	_, err = eq.db.Exec("UPDATE events SET status = ?, processed_at = ? WHERE id = ?",
		EventStatusCompleted, oldTime.Unix(), event.ID)
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}

	// Cleanup old events (older than 24 hours)
	deleted, err := eq.CleanupOldEvents(24 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to cleanup: %v", err)
	}

	if deleted == 0 {
		t.Error("Expected at least one event to be deleted")
	}
}

func TestCleanupOldEvents_PreservesPending(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create old pending event
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now().Add(-48 * time.Hour),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Cleanup should not delete pending events
	eq.CleanupOldEvents(24 * time.Hour)

	// Verify event still exists
	_, err = eq.GetEvent(event.ID)
	if err != nil {
		t.Error("Pending event should not be deleted")
	}
}

// Test File Watcher operations
func TestAddFileWatcher(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	watcher := &FileWatcher{
		Name:      "test-watcher",
		Path:      "/tmp/test",
		Pattern:   "*.txt",
		Events:    []string{"create", "modify"},
		Recursive: true,
		Enabled:   true,
	}

	err = eq.AddFileWatcher(watcher)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if watcher.ID == "" {
		t.Error("Expected ID to be assigned")
	}
}

func TestListFileWatchers(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Add watcher
	watcher := &FileWatcher{
		Name:    "test-watcher",
		Path:    "/tmp/test",
		Events:  []string{"create"},
		Enabled: true,
	}
	eq.AddFileWatcher(watcher)

	// List watchers
	watchers, err := eq.ListFileWatchers()
	if err != nil {
		t.Fatalf("Failed to list watchers: %v", err)
	}

	if len(watchers) == 0 {
		t.Error("Expected at least one watcher")
	}
}

func TestGetFileWatcher(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Add watcher
	original := &FileWatcher{
		Name:    "test-watcher",
		Path:    "/tmp/test",
		Events:  []string{"create"},
		Enabled: true,
	}
	eq.AddFileWatcher(original)

	// Get watcher
	retrieved, err := eq.GetFileWatcher("test-watcher")
	if err != nil {
		t.Fatalf("Failed to get watcher: %v", err)
	}

	if retrieved.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, retrieved.Name)
	}

	if retrieved.Path != original.Path {
		t.Errorf("Expected path %s, got %s", original.Path, retrieved.Path)
	}
}

func TestDeleteFileWatcher(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Add watcher
	watcher := &FileWatcher{
		Name:    "test-watcher",
		Path:    "/tmp/test",
		Events:  []string{"create"},
		Enabled: true,
	}
	eq.AddFileWatcher(watcher)

	// Delete watcher
	err = eq.DeleteFileWatcher("test-watcher")
	if err != nil {
		t.Fatalf("Failed to delete watcher: %v", err)
	}

	// Verify deletion
	_, err = eq.GetFileWatcher("test-watcher")
	if err == nil {
		t.Error("Expected error getting deleted watcher")
	}
}

func TestEnableFileWatcher(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Add disabled watcher
	watcher := &FileWatcher{
		Name:    "test-watcher",
		Path:    "/tmp/test",
		Events:  []string{"create"},
		Enabled: false,
	}
	eq.AddFileWatcher(watcher)

	// Enable watcher
	err = eq.EnableFileWatcher("test-watcher")
	if err != nil {
		t.Fatalf("Failed to enable watcher: %v", err)
	}

	// Verify enabled
	retrieved, err := eq.GetFileWatcher("test-watcher")
	if err != nil {
		t.Fatalf("Failed to get watcher: %v", err)
	}

	if !retrieved.Enabled {
		t.Error("Expected watcher to be enabled")
	}
}

func TestDisableFileWatcher(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Add enabled watcher
	watcher := &FileWatcher{
		Name:    "test-watcher",
		Path:    "/tmp/test",
		Events:  []string{"create"},
		Enabled: true,
	}
	eq.AddFileWatcher(watcher)

	// Disable watcher
	err = eq.DisableFileWatcher("test-watcher")
	if err != nil {
		t.Fatalf("Failed to disable watcher: %v", err)
	}

	// Verify disabled
	retrieved, err := eq.GetFileWatcher("test-watcher")
	if err != nil {
		t.Fatalf("Failed to get watcher: %v", err)
	}

	if retrieved.Enabled {
		t.Error("Expected watcher to be disabled")
	}
}

// Test RecordEventHookExecution
func TestRecordEventHookExecution(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create event
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	// Record execution
	result := &HookResult{
		Success:    true,
		Output:     "test output",
		Duration:   100 * time.Millisecond,
		ExecutedAt: time.Now(),
	}

	err = eq.RecordEventHookExecution(event.ID, "hook-id", "test-hook", result)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetEventHookExecutions(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create event and record execution
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}
	eq.EnqueueEvent(event)

	result := &HookResult{
		Success:    true,
		Output:     "test output",
		Duration:   100 * time.Millisecond,
		ExecutedAt: time.Now(),
	}
	eq.RecordEventHookExecution(event.ID, "hook-id", "test-hook", result)

	// Get executions
	executions, err := eq.GetEventHookExecutions(event.ID)
	if err != nil {
		t.Fatalf("Failed to get executions: %v", err)
	}

	if len(executions) == 0 {
		t.Error("Expected at least one execution")
	}
}

func TestGetHookExecutionsByAgent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	// Create event for specific agent
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
		Agent:     "test-agent",
	}
	eq.EnqueueEvent(event)

	// Record execution
	result := &HookResult{
		Success:    true,
		Output:     "test output",
		Duration:   100 * time.Millisecond,
		ExecutedAt: time.Now(),
	}
	eq.RecordEventHookExecution(event.ID, "hook-id", "test-hook", result)

	// Get executions by agent
	executions, err := eq.GetHookExecutionsByAgent("test-agent", 10)
	if err != nil {
		t.Fatalf("Failed to get executions: %v", err)
	}

	if len(executions) == 0 {
		t.Error("Expected at least one execution")
	}
}

// Test data marshaling
func TestEnqueueEvent_MarshalComplexData(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	complexData := map[string]interface{}{
		"string": "value",
		"number": 123,
		"array":  []string{"a", "b", "c"},
		"nested": map[string]interface{}{
			"key": "value",
		},
	}

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      complexData,
	}

	err = eq.EnqueueEvent(event)
	if err != nil {
		t.Fatalf("Failed to enqueue complex event: %v", err)
	}

	// Retrieve and verify
	retrieved, err := eq.GetEvent(event.ID)
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}

	if retrieved.Data["string"] != "value" {
		t.Error("Failed to preserve string data")
	}
}

// Test edge cases
func TestEnqueueEvent_EmptyData(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := repo.EventQueue

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = eq.EnqueueEvent(event)
	if err != nil {
		t.Errorf("Should handle empty data: %v", err)
	}
}

func TestGetPendingEvents_EmptyQueue(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	eq := NewEventQueue(repo.db)
	eq.InitializeSchema()

	// Clean all events first
	_, _ = eq.db.Exec("DELETE FROM events")

	events, err := eq.GetPendingEvents(10)
	if err != nil {
		t.Fatalf("Failed to get pending events: %v", err)
	}

	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}
