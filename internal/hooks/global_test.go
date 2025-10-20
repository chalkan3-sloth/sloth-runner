package hooks

import (
	"sync"
	"testing"
)

// Test InitializeGlobalDispatcher
func TestInitializeGlobalDispatcher(t *testing.T) {
	// Clean up any previous state
	CleanupGlobalDispatcher()

	err := InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	defer CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Error("Expected non-nil global dispatcher")
	}
}

func TestInitializeGlobalDispatcher_Idempotent(t *testing.T) {
	// Clean up any previous state
	CleanupGlobalDispatcher()

	// First initialization
	err := InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("First init failed: %v", err)
	}

	defer CleanupGlobalDispatcher()

	first := GetGlobalDispatcher()

	// Second initialization should be safe
	err = InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("Second init failed: %v", err)
	}

	second := GetGlobalDispatcher()

	// Should return same instance
	if first != second {
		t.Error("Expected same dispatcher instance")
	}
}

func TestInitializeGlobalDispatcher_CreatesRepository(t *testing.T) {
	// Clean up any previous state
	CleanupGlobalDispatcher()

	err := InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	defer CleanupGlobalDispatcher()

	globalMu.RLock()
	repo := globalRepo
	globalMu.RUnlock()

	if repo == nil {
		t.Error("Expected global repository to be created")
	}
}

// Test GetGlobalDispatcher
func TestGetGlobalDispatcher(t *testing.T) {
	// Clean up and initialize
	CleanupGlobalDispatcher()
	err := InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	defer CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Error("Expected non-nil dispatcher")
	}
}

func TestGetGlobalDispatcher_BeforeInit(t *testing.T) {
	// Clean up to ensure not initialized
	CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher != nil {
		t.Error("Expected nil dispatcher before initialization")
	}
}

func TestGetGlobalDispatcher_AfterCleanup(t *testing.T) {
	// Initialize and cleanup
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher != nil {
		t.Error("Expected nil dispatcher after cleanup")
	}
}

// Test CleanupGlobalDispatcher
func TestCleanupGlobalDispatcher(t *testing.T) {
	// Initialize first
	CleanupGlobalDispatcher()
	err := InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Cleanup
	CleanupGlobalDispatcher()

	globalMu.RLock()
	dispatcher := globalDispatcher
	repo := globalRepo
	globalMu.RUnlock()

	if dispatcher != nil {
		t.Error("Expected nil dispatcher after cleanup")
	}

	if repo != nil {
		t.Error("Expected nil repository after cleanup")
	}
}

func TestCleanupGlobalDispatcher_WhenNotInitialized(t *testing.T) {
	// Ensure not initialized
	CleanupGlobalDispatcher()

	// Cleanup should be safe even when not initialized
	CleanupGlobalDispatcher()
}

func TestCleanupGlobalDispatcher_MultipleCalls(t *testing.T) {
	// Initialize
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()

	// Multiple cleanups should be safe
	CleanupGlobalDispatcher()
	CleanupGlobalDispatcher()
	CleanupGlobalDispatcher()
}

// Test concurrent access
func TestGlobalDispatcher_ConcurrentInit(t *testing.T) {
	CleanupGlobalDispatcher()

	done := make(chan bool, 10)

	// Multiple goroutines trying to initialize
	for i := 0; i < 10; i++ {
		go func() {
			err := InitializeGlobalDispatcher()
			if err != nil {
				t.Logf("Init error: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	defer CleanupGlobalDispatcher()

	// Should have one valid dispatcher
	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Error("Expected non-nil dispatcher after concurrent init")
	}
}

func TestGlobalDispatcher_ConcurrentGet(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	done := make(chan bool, 20)

	// Multiple goroutines trying to get dispatcher
	for i := 0; i < 20; i++ {
		go func() {
			dispatcher := GetGlobalDispatcher()
			if dispatcher == nil {
				t.Error("Expected non-nil dispatcher")
			}
			done <- true
		}()
	}

	// Wait for all to complete
	for i := 0; i < 20; i++ {
		<-done
	}
}

func TestGlobalDispatcher_ConcurrentCleanup(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()

	done := make(chan bool, 5)

	// Multiple goroutines trying to cleanup
	for i := 0; i < 5; i++ {
		go func() {
			CleanupGlobalDispatcher()
			done <- true
		}()
	}

	// Wait for all to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Should be cleaned up
	dispatcher := GetGlobalDispatcher()
	if dispatcher != nil {
		t.Error("Expected nil dispatcher after concurrent cleanup")
	}
}

// Test lifecycle
func TestGlobalDispatcher_Lifecycle(t *testing.T) {
	// Start clean
	CleanupGlobalDispatcher()

	// Should be nil initially
	if GetGlobalDispatcher() != nil {
		t.Error("Expected nil dispatcher initially")
	}

	// Initialize
	err := InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Should be available
	if GetGlobalDispatcher() == nil {
		t.Error("Expected non-nil dispatcher after init")
	}

	// Cleanup
	CleanupGlobalDispatcher()

	// Should be nil again
	if GetGlobalDispatcher() != nil {
		t.Error("Expected nil dispatcher after cleanup")
	}
}

func TestGlobalDispatcher_MultipleLifecycles(t *testing.T) {
	// Multiple init/cleanup cycles
	for i := 0; i < 3; i++ {
		CleanupGlobalDispatcher()

		err := InitializeGlobalDispatcher()
		if err != nil {
			t.Fatalf("Init failed on cycle %d: %v", i, err)
		}

		dispatcher := GetGlobalDispatcher()
		if dispatcher == nil {
			t.Errorf("Expected non-nil dispatcher on cycle %d", i)
		}

		CleanupGlobalDispatcher()

		dispatcher = GetGlobalDispatcher()
		if dispatcher != nil {
			t.Errorf("Expected nil dispatcher after cleanup on cycle %d", i)
		}
	}
}

// Test global state isolation
func TestGlobalState_IsShared(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	// Get dispatcher in two different ways
	first := GetGlobalDispatcher()
	second := GetGlobalDispatcher()

	// Should be the same instance
	if first != second {
		t.Error("Expected same dispatcher instance")
	}
}

func TestGlobalState_ThreadSafe(t *testing.T) {
	CleanupGlobalDispatcher()

	var wg sync.WaitGroup

	// Start multiple operations concurrently
	wg.Add(3)

	// Initialize
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			InitializeGlobalDispatcher()
		}
	}()

	// Get
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			GetGlobalDispatcher()
		}
	}()

	// Cleanup
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			CleanupGlobalDispatcher()
		}
	}()

	wg.Wait()

	// Should not panic
}

// Test dispatcher functionality after initialization
func TestGlobalDispatcher_IsEnabled(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Fatal("Expected non-nil dispatcher")
	}

	if !dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be enabled by default")
	}
}

func TestGlobalDispatcher_CanDispatch(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Fatal("Expected non-nil dispatcher")
	}

	// Try to dispatch an event
	agent := &AgentEvent{
		Name:    "test-agent",
		Address: "localhost:50051",
	}

	err := dispatcher.DispatchAgentRegistered(agent)
	if err != nil {
		t.Errorf("Failed to dispatch event: %v", err)
	}
}

func TestGlobalDispatcher_EventProcessorRunning(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Fatal("Expected non-nil dispatcher")
	}

	dispatcher.mu.RLock()
	processing := dispatcher.processing
	dispatcher.mu.RUnlock()

	if !processing {
		t.Error("Expected event processor to be running")
	}
}

// Test repository access
func TestGlobalRepository_IsCreated(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	globalMu.RLock()
	repo := globalRepo
	globalMu.RUnlock()

	if repo == nil {
		t.Error("Expected global repository to be created")
	}

	if repo.db == nil {
		t.Error("Expected repository to have database")
	}
}

func TestGlobalRepository_HasEventQueue(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	globalMu.RLock()
	repo := globalRepo
	globalMu.RUnlock()

	if repo == nil {
		t.Fatal("Expected non-nil repository")
	}

	if repo.EventQueue == nil {
		t.Error("Expected repository to have event queue")
	}
}

// Test cleanup stops event processor
func TestCleanup_StopsEventProcessor(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Fatal("Expected non-nil dispatcher")
	}

	// Verify processor is running
	dispatcher.mu.RLock()
	processing := dispatcher.processing
	dispatcher.mu.RUnlock()

	if !processing {
		t.Error("Expected processor to be running before cleanup")
	}

	// Cleanup
	CleanupGlobalDispatcher()

	// After cleanup, getting dispatcher should return nil
	dispatcher = GetGlobalDispatcher()
	if dispatcher != nil {
		t.Error("Expected nil dispatcher after cleanup")
	}
}

// Test that cleanup closes repository
func TestCleanup_ClosesRepository(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()

	globalMu.RLock()
	repo := globalRepo
	globalMu.RUnlock()

	if repo == nil {
		t.Fatal("Expected non-nil repository")
	}

	// Cleanup
	CleanupGlobalDispatcher()

	// Repository should be nil
	globalMu.RLock()
	repo = globalRepo
	globalMu.RUnlock()

	if repo != nil {
		t.Error("Expected nil repository after cleanup")
	}
}

// Test error handling
func TestInitializeGlobalDispatcher_HandlesErrors(t *testing.T) {
	// This test verifies error handling by attempting multiple initializations
	CleanupGlobalDispatcher()

	// First init should succeed
	err := InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("First init should succeed: %v", err)
	}

	// Second init should be safe (returns early if already initialized)
	err = InitializeGlobalDispatcher()
	if err != nil {
		t.Fatalf("Second init should be safe: %v", err)
	}

	CleanupGlobalDispatcher()
}

// Test dispatcher operations through global
func TestGlobalDispatcher_EnableDisable(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Fatal("Expected non-nil dispatcher")
	}

	// Disable
	dispatcher.Disable()
	if dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be disabled")
	}

	// Enable
	dispatcher.Enable()
	if !dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be enabled")
	}
}

func TestGlobalDispatcher_ExecutionContext(t *testing.T) {
	CleanupGlobalDispatcher()
	InitializeGlobalDispatcher()
	defer CleanupGlobalDispatcher()

	dispatcher := GetGlobalDispatcher()
	if dispatcher == nil {
		t.Fatal("Expected non-nil dispatcher")
	}

	// Set context
	dispatcher.SetExecutionContext("stack", "agent", "run")

	// Verify context
	if dispatcher.GetCurrentStack() != "stack" {
		t.Error("Expected stack to be set")
	}

	if dispatcher.GetCurrentAgent() != "agent" {
		t.Error("Expected agent to be set")
	}

	if dispatcher.GetCurrentRunID() != "run" {
		t.Error("Expected runID to be set")
	}
}
