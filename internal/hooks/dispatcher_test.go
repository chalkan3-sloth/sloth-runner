package hooks

import (
	"testing"
	"time"
)

// Test NewDispatcher
func TestNewDispatcher(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)

	if dispatcher == nil {
		t.Error("Expected non-nil dispatcher")
	}

	if dispatcher.repo == nil {
		t.Error("Expected repository to be set")
	}

	if dispatcher.executor == nil {
		t.Error("Expected executor to be initialized")
	}

	if !dispatcher.enabled {
		t.Error("Expected dispatcher to be enabled by default")
	}
}

func TestNewDispatcher_HasStopChannel(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)

	if dispatcher.stopChan == nil {
		t.Error("Expected stopChan to be initialized")
	}
}

func TestNewDispatcher_HasEventChannel(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)

	if dispatcher.eventChannel == nil {
		t.Error("Expected eventChannel to be initialized")
	}

	// Channel should have buffer of 1000
	if cap(dispatcher.eventChannel) != 1000 {
		t.Errorf("Expected eventChannel buffer of 1000, got %d", cap(dispatcher.eventChannel))
	}
}

func TestNewDispatcher_MaxWorkers(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)

	if dispatcher.maxWorkers != 100 {
		t.Errorf("Expected maxWorkers to be 100, got %d", dispatcher.maxWorkers)
	}
}

// Test Dispatch
func TestDispatch_Success(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"task": map[string]interface{}{
				"task_name": "test-task",
			},
		},
	}

	err = dispatcher.Dispatch(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDispatch_WhenDisabled(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	dispatcher.Disable()

	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = dispatcher.Dispatch(event)
	if err != nil {
		t.Errorf("Expected no error when disabled, got %v", err)
	}
}

// Test Enable/Disable
func TestDispatcher_Enable(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	dispatcher.Disable()
	if dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be disabled")
	}

	dispatcher.Enable()
	if !dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be enabled")
	}
}

func TestDispatcher_Disable(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	if !dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be enabled initially")
	}

	dispatcher.Disable()
	if dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be disabled")
	}
}

func TestDispatcher_IsEnabled(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	if !dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be enabled by default")
	}
}

// Test DispatchAgentRegistered
func TestDispatchAgentRegistered(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	agent := &AgentEvent{
		Name:    "test-agent",
		Address: "localhost:50051",
		Tags:    []string{"tag1", "tag2"},
		Version: "1.0.0",
		SystemInfo: map[string]interface{}{
			"os": "linux",
		},
	}

	err = dispatcher.DispatchAgentRegistered(agent)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDispatchAgentRegistered_CreatesEvent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	agent := &AgentEvent{
		Name:    "test-agent",
		Address: "localhost:50051",
	}

	err = dispatcher.DispatchAgentRegistered(agent)
	if err != nil {
		t.Fatalf("Failed to dispatch: %v", err)
	}

	// Give it time to process
	time.Sleep(100 * time.Millisecond)

	// Verify event was queued
	events, err := repo.EventQueue.ListEvents(EventAgentRegistered, "", 10)
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected at least one event")
	}
}

// Test DispatchAgentDisconnected
func TestDispatchAgentDisconnected(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	agent := &AgentEvent{
		Name:    "test-agent",
		Address: "localhost:50051",
		Tags:    []string{"tag1"},
		Version: "1.0.0",
	}

	err = dispatcher.DispatchAgentDisconnected(agent)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test DispatchAgentHeartbeatFailed
func TestDispatchAgentHeartbeatFailed(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	agent := &AgentEvent{
		Name:    "test-agent",
		Address: "localhost:50051",
	}

	err = dispatcher.DispatchAgentHeartbeatFailed(agent)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test DispatchAgentUpdated
func TestDispatchAgentUpdated(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	agent := &AgentEvent{
		Name:    "test-agent",
		Address: "localhost:50051",
		Version: "2.0.0",
	}

	err = dispatcher.DispatchAgentUpdated(agent)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test DispatchTaskStarted
func TestDispatchTaskStarted(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	task := &TaskEvent{
		TaskName:  "test-task",
		AgentName: "test-agent",
		Status:    "running",
		Stack:     "test-stack",
		RunID:     "run-123",
	}

	err = dispatcher.DispatchTaskStarted(task)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test DispatchTaskCompleted
func TestDispatchTaskCompleted(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	task := &TaskEvent{
		TaskName:  "test-task",
		AgentName: "test-agent",
		Status:    "completed",
		ExitCode:  0,
		Duration:  "5s",
		Stack:     "test-stack",
		RunID:     "run-123",
	}

	err = dispatcher.DispatchTaskCompleted(task)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test DispatchTaskFailed
func TestDispatchTaskFailed(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	task := &TaskEvent{
		TaskName:  "test-task",
		AgentName: "test-agent",
		Status:    "failed",
		ExitCode:  1,
		Error:     "task failed",
		Duration:  "2s",
		Stack:     "test-stack",
		RunID:     "run-123",
	}

	err = dispatcher.DispatchTaskFailed(task)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test StartEventProcessor
func TestStartEventProcessor(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := &Dispatcher{
		repo:         repo,
		executor:     NewExecutor(repo),
		enabled:      true,
		stopChan:     make(chan struct{}),
		eventChannel: make(chan *Event, 1000),
		maxWorkers:   10,
	}

	dispatcher.StartEventProcessor()
	defer dispatcher.StopEventProcessor()

	// Verify processing started
	dispatcher.mu.RLock()
	processing := dispatcher.processing
	dispatcher.mu.RUnlock()

	if !processing {
		t.Error("Expected processing to be true")
	}
}

func TestStartEventProcessor_AlreadyRunning(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	// Should already be processing
	dispatcher.mu.RLock()
	processing := dispatcher.processing
	dispatcher.mu.RUnlock()

	if !processing {
		t.Error("Expected processing to be true")
	}

	// Starting again should be safe
	dispatcher.StartEventProcessor()
}

// Test StopEventProcessor
func TestStopEventProcessor(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)

	// Start processor
	dispatcher.StartEventProcessor()

	// Stop processor
	dispatcher.StopEventProcessor()

	dispatcher.mu.RLock()
	processing := dispatcher.processing
	dispatcher.mu.RUnlock()

	if processing {
		t.Error("Expected processing to be false after stop")
	}
}

func TestStopEventProcessor_WhenNotRunning(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := &Dispatcher{
		repo:         repo,
		executor:     NewExecutor(repo),
		enabled:      true,
		stopChan:     make(chan struct{}),
		eventChannel: make(chan *Event, 1000),
		maxWorkers:   10,
		processing:   false,
	}

	// Stopping when not running should be safe
	dispatcher.StopEventProcessor()
}

// Test execution context
func TestSetExecutionContext(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	dispatcher.SetExecutionContext("test-stack", "test-agent", "run-123")

	if dispatcher.GetCurrentStack() != "test-stack" {
		t.Errorf("Expected stack 'test-stack', got '%s'", dispatcher.GetCurrentStack())
	}

	if dispatcher.GetCurrentAgent() != "test-agent" {
		t.Errorf("Expected agent 'test-agent', got '%s'", dispatcher.GetCurrentAgent())
	}

	if dispatcher.GetCurrentRunID() != "run-123" {
		t.Errorf("Expected runID 'run-123', got '%s'", dispatcher.GetCurrentRunID())
	}
}

func TestGetCurrentStack(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	dispatcher.SetExecutionContext("my-stack", "", "")

	stack := dispatcher.GetCurrentStack()
	if stack != "my-stack" {
		t.Errorf("Expected 'my-stack', got '%s'", stack)
	}
}

func TestGetCurrentAgent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	dispatcher.SetExecutionContext("", "my-agent", "")

	agent := dispatcher.GetCurrentAgent()
	if agent != "my-agent" {
		t.Errorf("Expected 'my-agent', got '%s'", agent)
	}
}

func TestGetCurrentRunID(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	dispatcher.SetExecutionContext("", "", "run-456")

	runID := dispatcher.GetCurrentRunID()
	if runID != "run-456" {
		t.Errorf("Expected 'run-456', got '%s'", runID)
	}
}

// Test CreateEventDispatcherFunc
func TestCreateEventDispatcherFunc(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	fn := dispatcher.CreateEventDispatcherFunc()

	if fn == nil {
		t.Error("Expected non-nil function")
	}
}

func TestCreateEventDispatcherFunc_DispatchEvent(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	fn := dispatcher.CreateEventDispatcherFunc()

	data := map[string]interface{}{
		"test": "value",
	}

	err = fn("custom.event", data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateEventDispatcherFunc_UsesExecutionContext(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	dispatcher.SetExecutionContext("ctx-stack", "ctx-agent", "ctx-run")

	fn := dispatcher.CreateEventDispatcherFunc()

	err = fn("custom.event", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to dispatch: %v", err)
	}

	// Give it time to process
	time.Sleep(100 * time.Millisecond)

	// Verify event has execution context
	events, err := repo.EventQueue.ListEvents(EventType("custom.event"), "", 10)
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected at least one event")
		return
	}

	event := events[0]
	if event.Stack != "ctx-stack" {
		t.Errorf("Expected stack 'ctx-stack', got '%s'", event.Stack)
	}

	if event.Agent != "ctx-agent" {
		t.Errorf("Expected agent 'ctx-agent', got '%s'", event.Agent)
	}

	if event.RunID != "ctx-run" {
		t.Errorf("Expected runID 'ctx-run', got '%s'", event.RunID)
	}
}

// Test concurrent safety
func TestDispatcher_ConcurrentEnable(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			dispatcher.Enable()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if !dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be enabled")
	}
}

func TestDispatcher_ConcurrentDisable(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			dispatcher.Disable()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if dispatcher.IsEnabled() {
		t.Error("Expected dispatcher to be disabled")
	}
}

func TestDispatcher_ConcurrentContextAccess(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	dispatcher := NewDispatcher(repo)
	defer dispatcher.StopEventProcessor()

	done := make(chan bool, 20)

	// Concurrent setters
	for i := 0; i < 10; i++ {
		go func(n int) {
			dispatcher.SetExecutionContext("stack", "agent", "run")
			done <- true
		}(i)
	}

	// Concurrent getters
	for i := 0; i < 10; i++ {
		go func() {
			_ = dispatcher.GetCurrentStack()
			_ = dispatcher.GetCurrentAgent()
			_ = dispatcher.GetCurrentRunID()
			done <- true
		}()
	}

	for i := 0; i < 20; i++ {
		<-done
	}
}

// Test getCurrentTime helper
func TestGetCurrentTime(t *testing.T) {
	before := time.Now()
	current := getCurrentTime()
	after := time.Now()

	if current.Before(before) || current.After(after.Add(time.Second)) {
		t.Error("getCurrentTime returned unexpected time")
	}
}
