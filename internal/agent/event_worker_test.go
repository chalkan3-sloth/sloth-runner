package agent

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Mock gRPC client for testing
type mockAgentRegistryClient struct {
	pb.AgentRegistryClient
	sendEventBatchFunc func(ctx context.Context, in *pb.SendEventBatchRequest, opts ...grpc.CallOption) (*pb.SendEventBatchResponse, error)
	mu                 sync.Mutex
	receivedEvents     []*pb.EventData
}

func (m *mockAgentRegistryClient) SendEventBatch(ctx context.Context, in *pb.SendEventBatchRequest, opts ...grpc.CallOption) (*pb.SendEventBatchResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sendEventBatchFunc != nil {
		return m.sendEventBatchFunc(ctx, in, opts...)
	}

	// Store received events
	m.receivedEvents = append(m.receivedEvents, in.Events...)

	return &pb.SendEventBatchResponse{
		Success:          true,
		EventsProcessed:  int32(len(in.Events)),
		FailedEventIds:   []string{},
		Message:          "Success",
	}, nil
}

// TestNewEventWorker tests creating a new event worker
func TestNewEventWorker(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:     "test-agent",
		MasterAddr:    "localhost:50051",
		BatchSize:     100,
		FlushInterval: 5 * time.Second,
	}

	worker := NewEventWorker(config)

	if worker == nil {
		t.Fatal("Expected non-nil worker")
	}

	if worker.agentName != "test-agent" {
		t.Errorf("Expected agent name 'test-agent', got %s", worker.agentName)
	}

	if worker.batchSize != 100 {
		t.Errorf("Expected batch size 100, got %d", worker.batchSize)
	}

	if worker.flushInterval != 5*time.Second {
		t.Errorf("Expected flush interval 5s, got %v", worker.flushInterval)
	}
}

// TestNewEventWorker_Defaults tests default values
func TestNewEventWorker_Defaults(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		// BatchSize and FlushInterval not set
	}

	worker := NewEventWorker(config)

	if worker.batchSize != 50 {
		t.Errorf("Expected default batch size 50, got %d", worker.batchSize)
	}

	if worker.flushInterval != 10*time.Second {
		t.Errorf("Expected default flush interval 10s, got %v", worker.flushInterval)
	}
}

// TestNewEventWorker_EmptyConfig tests with minimal config
func TestNewEventWorker_EmptyConfig(t *testing.T) {
	config := EventWorkerConfig{}

	worker := NewEventWorker(config)

	if worker == nil {
		t.Fatal("Expected non-nil worker even with empty config")
	}

	if worker.agentName != "" {
		t.Error("Expected empty agent name")
	}

	if worker.events == nil {
		t.Error("Expected initialized events slice")
	}
}

// TestEventWorkerConfig_Various tests various configurations
func TestEventWorkerConfig_Various(t *testing.T) {
	testCases := []struct {
		name          string
		config        EventWorkerConfig
		expectedBatch int
		expectedFlush time.Duration
	}{
		{
			name: "Custom values",
			config: EventWorkerConfig{
				AgentName:     "agent1",
				MasterAddr:    "localhost:50051",
				BatchSize:     200,
				FlushInterval: 30 * time.Second,
			},
			expectedBatch: 200,
			expectedFlush: 30 * time.Second,
		},
		{
			name: "Only batch size",
			config: EventWorkerConfig{
				AgentName:  "agent2",
				MasterAddr: "localhost:50051",
				BatchSize:  75,
			},
			expectedBatch: 75,
			expectedFlush: 10 * time.Second,
		},
		{
			name: "Only flush interval",
			config: EventWorkerConfig{
				AgentName:     "agent3",
				MasterAddr:    "localhost:50051",
				FlushInterval: 20 * time.Second,
			},
			expectedBatch: 50,
			expectedFlush: 20 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker := NewEventWorker(tc.config)

			if worker.batchSize != tc.expectedBatch {
				t.Errorf("Expected batch size %d, got %d", tc.expectedBatch, worker.batchSize)
			}

			if worker.flushInterval != tc.expectedFlush {
				t.Errorf("Expected flush interval %v, got %v", tc.expectedFlush, worker.flushInterval)
			}
		})
	}
}

// TestEventWorker_SendEvent tests sending a single event
func TestEventWorker_SendEvent(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:     "test-agent",
		MasterAddr:    "localhost:50051",
		BatchSize:     10,
		FlushInterval: 1 * time.Minute,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	err := worker.SendEvent("test.event", "test-stack", "run-123", data)
	if err != nil {
		t.Fatalf("SendEvent failed: %v", err)
	}

	worker.mu.Lock()
	bufferSize := len(worker.events)
	worker.mu.Unlock()

	if bufferSize != 1 {
		t.Errorf("Expected 1 event in buffer, got %d", bufferSize)
	}
}

// TestEventWorker_SendEvent_BatchFlush tests auto-flush when batch is full
func TestEventWorker_SendEvent_BatchFlush(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:     "test-agent",
		MasterAddr:    "localhost:50051",
		BatchSize:     3, // Small batch for testing
		FlushInterval: 1 * time.Minute,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	// Send 3 events (should trigger flush)
	for i := 0; i < 3; i++ {
		data := map[string]interface{}{"index": i}
		err := worker.SendEvent("test.event", "stack", "run-1", data)
		if err != nil {
			t.Fatalf("SendEvent %d failed: %v", i, err)
		}
	}

	// Check mock client received events
	mockClient.mu.Lock()
	received := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != 3 {
		t.Errorf("Expected 3 events sent to master, got %d", received)
	}
}

// TestEventWorker_SendEvent_InvalidData tests with data that can't be marshaled
func TestEventWorker_SendEvent_InvalidData(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	// Channel can't be marshaled to JSON
	data := map[string]interface{}{
		"channel": make(chan int),
	}

	err := worker.SendEvent("test.event", "stack", "run-1", data)
	if err == nil {
		t.Error("Expected error when marshaling invalid data")
	}
}

// TestEventWorker_SendEvent_EmptyData tests with empty data map
func TestEventWorker_SendEvent_EmptyData(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	err := worker.SendEvent("test.event", "stack", "run-1", map[string]interface{}{})
	if err != nil {
		t.Fatalf("SendEvent with empty data failed: %v", err)
	}

	worker.mu.Lock()
	bufferSize := len(worker.events)
	worker.mu.Unlock()

	if bufferSize != 1 {
		t.Error("Expected 1 event in buffer")
	}
}

// TestEventWorker_SendEvent_NilData tests with nil data
func TestEventWorker_SendEvent_NilData(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	err := worker.SendEvent("test.event", "stack", "run-1", nil)
	if err != nil {
		t.Fatalf("SendEvent with nil data failed: %v", err)
	}

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	if event.DataJson != "null" {
		t.Errorf("Expected 'null' for nil data, got %s", event.DataJson)
	}
}

// TestEventWorker_SendEventWithSeverity tests sending events with custom severity
func TestEventWorker_SendEventWithSeverity(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	severities := []string{"info", "warning", "error", "critical"}

	for _, severity := range severities {
		data := map[string]interface{}{"severity_test": severity}
		err := worker.SendEventWithSeverity("test.event", "stack", "run-1", data, severity)
		if err != nil {
			t.Fatalf("SendEventWithSeverity failed for %s: %v", severity, err)
		}
	}

	worker.mu.Lock()
	bufferSize := len(worker.events)
	worker.mu.Unlock()

	if bufferSize != 4 {
		t.Errorf("Expected 4 events in buffer, got %d", bufferSize)
	}
}

// TestEventWorker_SendEventWithSeverity_AutoFlush tests auto-flush with severity
func TestEventWorker_SendEventWithSeverity_AutoFlush(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  2,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	data := map[string]interface{}{"test": "data"}

	// Send 2 events (should auto-flush)
	worker.SendEventWithSeverity("event1", "stack", "run-1", data, "info")
	worker.SendEventWithSeverity("event2", "stack", "run-1", data, "warning")

	mockClient.mu.Lock()
	received := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != 2 {
		t.Errorf("Expected 2 events flushed, got %d", received)
	}
}

// TestEventWorker_Flush tests manual flush
func TestEventWorker_Flush(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	// Send 5 events without triggering auto-flush
	for i := 0; i < 5; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.event", "stack", "run-1", data)
	}

	// Manual flush
	err := worker.flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	mockClient.mu.Lock()
	received := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != 5 {
		t.Errorf("Expected 5 events flushed, got %d", received)
	}

	// Buffer should be empty after flush
	worker.mu.Lock()
	bufferSize := len(worker.events)
	worker.mu.Unlock()

	if bufferSize != 0 {
		t.Errorf("Expected empty buffer after flush, got %d events", bufferSize)
	}
}

// TestEventWorker_Flush_EmptyBuffer tests flushing empty buffer
func TestEventWorker_Flush_EmptyBuffer(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	err := worker.flush()
	if err != nil {
		t.Errorf("Flush on empty buffer should not error, got: %v", err)
	}

	mockClient.mu.Lock()
	received := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != 0 {
		t.Errorf("Expected 0 events sent, got %d", received)
	}
}

// TestEventWorker_Flush_MultipleTimes tests multiple flushes
func TestEventWorker_Flush_MultipleTimes(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	totalEvents := 0

	// Flush 3 times with different event counts
	for round := 1; round <= 3; round++ {
		eventCount := round * 2

		for i := 0; i < eventCount; i++ {
			data := map[string]interface{}{"round": round, "index": i}
			worker.SendEvent("test.event", "stack", "run-1", data)
		}

		worker.flush()
		totalEvents += eventCount
	}

	mockClient.mu.Lock()
	received := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != totalEvents {
		t.Errorf("Expected %d total events, got %d", totalEvents, received)
	}
}

// TestEventWorker_ConcurrentSendEvent tests concurrent event sending
func TestEventWorker_ConcurrentSendEvent(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  1000, // Large batch to avoid auto-flush
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	var wg sync.WaitGroup
	goroutines := 10
	eventsPerGoroutine := 10

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				data := map[string]interface{}{"goroutine": id, "event": j}
				worker.SendEvent("test.concurrent", "stack", "run-1", data)
			}
		}(i)
	}

	wg.Wait()

	worker.mu.Lock()
	bufferSize := len(worker.events)
	worker.mu.Unlock()

	expectedEvents := goroutines * eventsPerGoroutine
	if bufferSize != expectedEvents {
		t.Errorf("Expected %d events in buffer, got %d", expectedEvents, bufferSize)
	}
}

// TestEventWorker_Stop tests stopping the worker
func TestEventWorker_Stop(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:     "test-agent",
		MasterAddr:    "localhost:50051",
		BatchSize:     100,
		FlushInterval: 100 * time.Millisecond,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	// Add some events
	for i := 0; i < 5; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.event", "stack", "run-1", data)
	}

	// Stop should flush remaining events
	err := worker.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	mockClient.mu.Lock()
	received := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != 5 {
		t.Errorf("Expected 5 events flushed on stop, got %d", received)
	}
}

// TestEventWorker_EventDataFields tests that event data fields are populated correctly
func TestEventWorker_EventDataFields(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent-123",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	data := map[string]interface{}{
		"field1": "value1",
		"field2": 42,
	}

	worker.SendEvent("custom.event.type", "production-stack", "run-abc-123", data)

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	if event.EventType != "custom.event.type" {
		t.Errorf("Expected event type 'custom.event.type', got %s", event.EventType)
	}

	if event.AgentName != "test-agent-123" {
		t.Errorf("Expected agent name 'test-agent-123', got %s", event.AgentName)
	}

	if event.Stack != "production-stack" {
		t.Errorf("Expected stack 'production-stack', got %s", event.Stack)
	}

	if event.RunId != "run-abc-123" {
		t.Errorf("Expected run ID 'run-abc-123', got %s", event.RunId)
	}

	if event.Severity != "info" {
		t.Errorf("Expected default severity 'info', got %s", event.Severity)
	}

	if event.EventId == "" {
		t.Error("Expected non-empty event ID")
	}

	if event.Timestamp == 0 {
		t.Error("Expected non-zero timestamp")
	}

	// Verify data JSON
	var parsedData map[string]interface{}
	err := json.Unmarshal([]byte(event.DataJson), &parsedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal data JSON: %v", err)
	}

	if parsedData["field1"] != "value1" {
		t.Errorf("Expected field1='value1', got %v", parsedData["field1"])
	}

	if parsedData["field2"].(float64) != 42 {
		t.Errorf("Expected field2=42, got %v", parsedData["field2"])
	}
}

// TestEventWorker_UniqueEventIDs tests that event IDs are unique
func TestEventWorker_UniqueEventIDs(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)

	eventIDs := make(map[string]bool)

	for i := 0; i < 50; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.event", "stack", "run-1", data)
	}

	worker.mu.Lock()
	for _, event := range worker.events {
		if eventIDs[event.EventId] {
			t.Errorf("Duplicate event ID found: %s", event.EventId)
		}
		eventIDs[event.EventId] = true
	}
	worker.mu.Unlock()

	if len(eventIDs) != 50 {
		t.Errorf("Expected 50 unique event IDs, got %d", len(eventIDs))
	}
}

// TestEventWorker_EmptyStrings tests with empty string parameters
func TestEventWorker_EmptyStrings(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	data := map[string]interface{}{"test": "data"}

	err := worker.SendEvent("", "", "", data)
	if err != nil {
		t.Fatalf("SendEvent with empty strings failed: %v", err)
	}

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	if event.EventType != "" {
		t.Errorf("Expected empty event type, got %s", event.EventType)
	}

	if event.Stack != "" {
		t.Errorf("Expected empty stack, got %s", event.Stack)
	}

	if event.RunId != "" {
		t.Errorf("Expected empty run ID, got %s", event.RunId)
	}
}

// TestEventWorker_LargeData tests with large data payload
func TestEventWorker_LargeData(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	// Create large data structure
	largeData := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		largeData[string(rune('a'+i%26))+string(rune('0'+i%10))] = i
	}

	err := worker.SendEvent("test.large", "stack", "run-1", largeData)
	if err != nil {
		t.Fatalf("SendEvent with large data failed: %v", err)
	}

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	if len(event.DataJson) == 0 {
		t.Error("Expected non-empty data JSON for large data")
	}
}

// TestEventWorker_NestedData tests with nested data structures
func TestEventWorker_NestedData(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	nestedData := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": "deep value",
			},
		},
		"array": []interface{}{1, 2, 3, "four"},
	}

	err := worker.SendEvent("test.nested", "stack", "run-1", nestedData)
	if err != nil {
		t.Fatalf("SendEvent with nested data failed: %v", err)
	}

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	// Verify nested data can be unmarshaled
	var parsedData map[string]interface{}
	err = json.Unmarshal([]byte(event.DataJson), &parsedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal nested data: %v", err)
	}

	level1, ok := parsedData["level1"].(map[string]interface{})
	if !ok {
		t.Error("Expected nested level1 to be map")
	}

	level2, ok := level1["level2"].(map[string]interface{})
	if !ok {
		t.Error("Expected nested level2 to be map")
	}

	if level2["level3"] != "deep value" {
		t.Errorf("Expected deep value, got %v", level2["level3"])
	}
}

// TestEventWorker_SpecialCharactersInData tests with special characters
func TestEventWorker_SpecialCharactersInData(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	data := map[string]interface{}{
		"special_chars": "Hello 世界 🎉 \n\t\r",
		"unicode":       "Ñoño",
		"emoji":         "😀😁😂",
	}

	err := worker.SendEvent("test.special", "stack", "run-1", data)
	if err != nil {
		t.Fatalf("SendEvent with special characters failed: %v", err)
	}

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	// Verify data can be unmarshaled
	var parsedData map[string]interface{}
	err = json.Unmarshal([]byte(event.DataJson), &parsedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal data with special characters: %v", err)
	}

	if parsedData["emoji"] != "😀😁😂" {
		t.Errorf("Emoji not preserved correctly, got: %v", parsedData["emoji"])
	}
}

// TestEventWorker_BufferCapacity tests buffer capacity management
func TestEventWorker_BufferCapacity(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	// Check initial capacity
	worker.mu.Lock()
	initialCap := cap(worker.events)
	worker.mu.Unlock()

	if initialCap != 10 {
		t.Errorf("Expected initial capacity 10, got %d", initialCap)
	}
}

// TestEventWorker_ContextCancellation tests context cancellation
func TestEventWorker_ContextCancellation(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:     "test-agent",
		MasterAddr:    "localhost:50051",
		BatchSize:     100,
		FlushInterval: 100 * time.Millisecond,
	}

	worker := NewEventWorker(config)

	if worker.ctx == nil {
		t.Error("Expected non-nil context")
	}

	if worker.cancel == nil {
		t.Error("Expected non-nil cancel function")
	}

	// Cancel context
	worker.cancel()

	// Context should be done
	select {
	case <-worker.ctx.Done():
		// Expected
	default:
		t.Error("Expected context to be done after cancel")
	}
}

// TestEventWorker_NilClient tests behavior with nil client
func TestEventWorker_NilClient(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10, // Large batch size to prevent auto-flush
	}

	worker := NewEventWorker(config)
	// Don't set client - it's nil

	data := map[string]interface{}{"test": "data"}

	// Send events without triggering flush
	err1 := worker.SendEvent("event1", "stack", "run-1", data)
	err2 := worker.SendEvent("event2", "stack", "run-1", data)

	// Events should be added to buffer successfully (no error)
	if err1 != nil {
		t.Errorf("Expected no error when adding event to buffer, got %v", err1)
	}
	if err2 != nil {
		t.Errorf("Expected no error when adding event to buffer, got %v", err2)
	}

	// Verify events are in buffer
	worker.mu.Lock()
	bufferSize := len(worker.events)
	worker.mu.Unlock()

	if bufferSize != 2 {
		t.Errorf("Expected 2 events in buffer, got %d", bufferSize)
	}
}

// TestEventWorker_BatchSizeBoundary tests batch size boundary conditions
func TestEventWorker_BatchSizeBoundary(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  5,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	// Send exactly batch size - 1 events (should not flush)
	for i := 0; i < 4; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.event", "stack", "run-1", data)
	}

	mockClient.mu.Lock()
	received := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != 0 {
		t.Errorf("Expected 0 events flushed (under threshold), got %d", received)
	}

	// Send one more to trigger flush
	worker.SendEvent("test.event", "stack", "run-1", map[string]interface{}{"index": 4})

	mockClient.mu.Lock()
	received = len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != 5 {
		t.Errorf("Expected 5 events flushed (at threshold), got %d", received)
	}
}

// TestEventWorker_TimeStampAccuracy tests timestamp generation
func TestEventWorker_TimeStampAccuracy(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	before := time.Now().Unix()
	worker.SendEvent("test.timestamp", "stack", "run-1", map[string]interface{}{"test": "timestamp"})
	after := time.Now().Unix()

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	if event.Timestamp < before || event.Timestamp > after {
		t.Errorf("Timestamp %d not within range [%d, %d]", event.Timestamp, before, after)
	}
}

// TestEventWorker_MultipleAgents tests events from different agents
func TestEventWorker_MultipleAgents(t *testing.T) {
	agents := []string{"agent1", "agent2", "agent3"}
	workers := make([]*EventWorker, len(agents))
	mockClients := make([]*mockAgentRegistryClient, len(agents))

	for i, agent := range agents {
		config := EventWorkerConfig{
			AgentName:  agent,
			MasterAddr: "localhost:50051",
			BatchSize:  10,
		}
		workers[i] = NewEventWorker(config)
		mockClients[i] = &mockAgentRegistryClient{}
		workers[i].client = mockClients[i]

		// Send event from each agent
		data := map[string]interface{}{"agent_id": i}
		workers[i].SendEvent("test.multi_agent", "stack", "run-1", data)
	}

	// Verify each agent has events
	for i, worker := range workers {
		worker.mu.Lock()
		bufferSize := len(worker.events)
		worker.mu.Unlock()

		if bufferSize != 1 {
			t.Errorf("Agent %d expected 1 event, got %d", i, bufferSize)
		}
	}
}

// TestEventWorker_RapidFire tests sending many events quickly
func TestEventWorker_RapidFire(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  1000,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	eventCount := 500

	start := time.Now()
	for i := 0; i < eventCount; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.rapid", "stack", "run-1", data)
	}
	elapsed := time.Since(start)

	worker.mu.Lock()
	bufferSize := len(worker.events)
	worker.mu.Unlock()

	if bufferSize != eventCount {
		t.Errorf("Expected %d events in buffer, got %d", eventCount, bufferSize)
	}

	// Verify reasonable performance (should take less than 1 second)
	if elapsed > time.Second {
		t.Errorf("Rapid fire took too long: %v", elapsed)
	}
}

// TestEventWorker_DifferentEventTypes tests various event types
func TestEventWorker_DifferentEventTypes(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)

	eventTypes := []string{
		"system.cpu",
		"system.memory",
		"application.error",
		"user.login",
		"deployment.complete",
	}

	for _, eventType := range eventTypes {
		data := map[string]interface{}{"type": eventType}
		worker.SendEvent(eventType, "stack", "run-1", data)
	}

	worker.mu.Lock()
	if len(worker.events) != len(eventTypes) {
		t.Errorf("Expected %d events, got %d", len(eventTypes), len(worker.events))
	}

	// Verify event types
	foundTypes := make(map[string]bool)
	for _, event := range worker.events {
		foundTypes[event.EventType] = true
	}
	worker.mu.Unlock()

	for _, expectedType := range eventTypes {
		if !foundTypes[expectedType] {
			t.Errorf("Event type %s not found", expectedType)
		}
	}
}

// TestEventWorker_DifferentStacks tests events from different stacks
func TestEventWorker_DifferentStacks(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	stacks := []string{"development", "staging", "production"}

	for _, stack := range stacks {
		data := map[string]interface{}{"stack": stack}
		worker.SendEvent("test.stack", stack, "run-1", data)
	}

	worker.mu.Lock()
	for i, event := range worker.events {
		if event.Stack != stacks[i] {
			t.Errorf("Event %d expected stack %s, got %s", i, stacks[i], event.Stack)
		}
	}
	worker.mu.Unlock()
}

// TestEventWorker_DifferentRunIDs tests events with different run IDs
func TestEventWorker_DifferentRunIDs(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	runIDs := []string{"run-001", "run-002", "run-003"}

	for _, runID := range runIDs {
		data := map[string]interface{}{"run_id": runID}
		worker.SendEvent("test.run", "stack", runID, data)
	}

	worker.mu.Lock()
	for i, event := range worker.events {
		if event.RunId != runIDs[i] {
			t.Errorf("Event %d expected run ID %s, got %s", i, runIDs[i], event.RunId)
		}
	}
	worker.mu.Unlock()
}

// TestEventWorker_EventDataTypes tests various data types in event data
func TestEventWorker_EventDataTypes(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	data := map[string]interface{}{
		"string_field":  "test string",
		"int_field":     42,
		"float_field":   3.14,
		"bool_field":    true,
		"null_field":    nil,
		"array_field":   []interface{}{1, 2, 3},
		"object_field":  map[string]interface{}{"nested": "value"},
	}

	err := worker.SendEvent("test.datatypes", "stack", "run-1", data)
	if err != nil {
		t.Fatalf("Failed to send event with various data types: %v", err)
	}

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	// Verify data can be unmarshaled
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(event.DataJson), &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal event data: %v", err)
	}

	if parsed["string_field"] != "test string" {
		t.Errorf("String field not preserved")
	}

	if parsed["bool_field"] != true {
		t.Errorf("Bool field not preserved")
	}
}

// TestEventWorker_EventSeverityLevels tests all severity levels
func TestEventWorker_EventSeverityLevels(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)

	severities := []string{"debug", "info", "warning", "error", "critical", "fatal"}

	for _, severity := range severities {
		data := map[string]interface{}{"level": severity}
		worker.SendEventWithSeverity("test.severity", "stack", "run-1", data, severity)
	}

	worker.mu.Lock()
	for i, event := range worker.events {
		if event.Severity != severities[i] {
			t.Errorf("Event %d expected severity %s, got %s", i, severities[i], event.Severity)
		}
	}
	worker.mu.Unlock()
}

// TestEventWorker_FlushPreservesOrder tests that flush preserves event order
func TestEventWorker_FlushPreservesOrder(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	// Send events in order
	for i := 0; i < 10; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.order", "stack", "run-1", data)
	}

	worker.flush()

	mockClient.mu.Lock()
	events := mockClient.receivedEvents
	mockClient.mu.Unlock()

	// Verify order
	for i, event := range events {
		var data map[string]interface{}
		json.Unmarshal([]byte(event.DataJson), &data)
		if data["index"].(float64) != float64(i) {
			t.Errorf("Event %d has wrong index: %v", i, data["index"])
		}
	}
}

// TestEventWorker_BufferClearAfterFlush tests buffer is cleared after flush
func TestEventWorker_BufferClearAfterFlush(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	// Fill buffer
	for i := 0; i < 20; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.buffer", "stack", "run-1", data)
	}

	worker.mu.Lock()
	bufferSizeBefore := len(worker.events)
	worker.mu.Unlock()

	if bufferSizeBefore != 20 {
		t.Errorf("Expected 20 events before flush, got %d", bufferSizeBefore)
	}

	// Flush
	worker.flush()

	worker.mu.Lock()
	bufferSizeAfter := len(worker.events)
	worker.mu.Unlock()

	if bufferSizeAfter != 0 {
		t.Errorf("Expected 0 events after flush, got %d", bufferSizeAfter)
	}
}

// TestEventWorker_PartialFlushFailure tests behavior when flush partially fails
func TestEventWorker_PartialFlushFailure(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{
		sendEventBatchFunc: func(ctx context.Context, in *pb.SendEventBatchRequest, opts ...grpc.CallOption) (*pb.SendEventBatchResponse, error) {
			return &pb.SendEventBatchResponse{
				Success:          false,
				EventsProcessed:  int32(len(in.Events) / 2),
				FailedEventIds:   []string{"event-1", "event-2"},
				Message:          "Partial failure",
			}, nil
		},
	}
	worker.client = mockClient

	// Send events
	for i := 0; i < 10; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.partial", "stack", "run-1", data)
	}

	// Flush (should report partial success)
	err := worker.flush()
	if err != nil {
		t.Errorf("Expected no error on partial flush, got: %v", err)
	}
}

// TestEventWorker_AgentNamePropagation tests agent name is set correctly
func TestEventWorker_AgentNamePropagation(t *testing.T) {
	agentName := "unique-test-agent-123"
	config := EventWorkerConfig{
		AgentName:  agentName,
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	data := map[string]interface{}{"test": "agent_name"}
	worker.SendEvent("test.agent_name", "stack", "run-1", data)

	worker.mu.Lock()
	event := worker.events[0]
	worker.mu.Unlock()

	if event.AgentName != agentName {
		t.Errorf("Expected agent name %s, got %s", agentName, event.AgentName)
	}
}

// TestEventWorker_MasterAddrConfiguration tests master address is stored
func TestEventWorker_MasterAddrConfiguration(t *testing.T) {
	masterAddr := "master.example.com:50051"
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: masterAddr,
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	if worker.masterAddr != masterAddr {
		t.Errorf("Expected master addr %s, got %s", masterAddr, worker.masterAddr)
	}
}

// TestEventWorker_BatchSizeConfiguration tests batch size is stored
func TestEventWorker_BatchSizeConfiguration(t *testing.T) {
	batchSize := 123
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  batchSize,
	}

	worker := NewEventWorker(config)

	if worker.batchSize != batchSize {
		t.Errorf("Expected batch size %d, got %d", batchSize, worker.batchSize)
	}
}

// TestEventWorker_FlushIntervalConfiguration tests flush interval is stored
func TestEventWorker_FlushIntervalConfiguration(t *testing.T) {
	flushInterval := 30 * time.Second
	config := EventWorkerConfig{
		AgentName:     "test-agent",
		MasterAddr:    "localhost:50051",
		BatchSize:     10,
		FlushInterval: flushInterval,
	}

	worker := NewEventWorker(config)

	if worker.flushInterval != flushInterval {
		t.Errorf("Expected flush interval %v, got %v", flushInterval, worker.flushInterval)
	}
}

// TestEventWorker_ZeroBatchSize tests zero batch size gets default
func TestEventWorker_ZeroBatchSize(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  0, // Zero should trigger default
	}

	worker := NewEventWorker(config)

	if worker.batchSize != 50 {
		t.Errorf("Expected default batch size 50, got %d", worker.batchSize)
	}
}

// TestEventWorker_ZeroFlushInterval tests zero flush interval gets default
func TestEventWorker_ZeroFlushInterval(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:     "test-agent",
		MasterAddr:    "localhost:50051",
		BatchSize:     10,
		FlushInterval: 0, // Zero should trigger default
	}

	worker := NewEventWorker(config)

	if worker.flushInterval != 10*time.Second {
		t.Errorf("Expected default flush interval 10s, got %v", worker.flushInterval)
	}
}

// TestEventWorker_VeryLargeBatchSize tests with very large batch size
func TestEventWorker_VeryLargeBatchSize(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10000,
	}

	worker := NewEventWorker(config)

	if worker.batchSize != 10000 {
		t.Errorf("Expected batch size 10000, got %d", worker.batchSize)
	}

	// Send events without triggering flush
	for i := 0; i < 100; i++ {
		data := map[string]interface{}{"index": i}
		worker.SendEvent("test.large_batch", "stack", "run-1", data)
	}

	worker.mu.Lock()
	bufferSize := len(worker.events)
	worker.mu.Unlock()

	if bufferSize != 100 {
		t.Errorf("Expected 100 events in buffer, got %d", bufferSize)
	}
}

// TestEventWorker_VerySmallBatchSize tests with batch size of 1
func TestEventWorker_VerySmallBatchSize(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  1,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	// Single event should trigger immediate flush
	data := map[string]interface{}{"test": "small_batch"}
	worker.SendEvent("test.small", "stack", "run-1", data)

	mockClient.mu.Lock()
	received := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if received != 1 {
		t.Errorf("Expected 1 event flushed immediately, got %d", received)
	}
}

// TestEventWorker_ConsecutiveFlushes tests multiple consecutive flushes
func TestEventWorker_ConsecutiveFlushes(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  100,
	}

	worker := NewEventWorker(config)
	mockClient := &mockAgentRegistryClient{}
	worker.client = mockClient

	// Flush multiple times in a row
	for i := 0; i < 5; i++ {
		// Add events
		for j := 0; j < 3; j++ {
			data := map[string]interface{}{"round": i, "index": j}
			worker.SendEvent("test.consecutive", "stack", "run-1", data)
		}

		// Flush
		err := worker.flush()
		if err != nil {
			t.Fatalf("Flush %d failed: %v", i, err)
		}
	}

	mockClient.mu.Lock()
	totalReceived := len(mockClient.receivedEvents)
	mockClient.mu.Unlock()

	if totalReceived != 15 {
		t.Errorf("Expected 15 total events across 5 flushes, got %d", totalReceived)
	}
}

// TestEventWorker_EventIDFormat tests that event IDs are valid UUIDs
func TestEventWorker_EventIDFormat(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	data := map[string]interface{}{"test": "uuid"}
	worker.SendEvent("test.uuid", "stack", "run-1", data)

	worker.mu.Lock()
	eventID := worker.events[0].EventId
	worker.mu.Unlock()

	// UUID should be 36 characters (with hyphens)
	if len(eventID) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(eventID))
	}

	// Check format (8-4-4-4-12)
	parts := 0
	for _, char := range eventID {
		if char == '-' {
			parts++
		}
	}

	if parts != 4 {
		t.Errorf("Expected 4 hyphens in UUID, got %d", parts)
	}
}

// TestEventWorker_JSONEncodingPreservation tests JSON encoding is preserved
func TestEventWorker_JSONEncodingPreservation(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	originalData := map[string]interface{}{
		"number": 42,
		"string": "test",
		"bool":   true,
		"null":   nil,
	}

	worker.SendEvent("test.json", "stack", "run-1", originalData)

	worker.mu.Lock()
	dataJSON := worker.events[0].DataJson
	worker.mu.Unlock()

	// Parse back
	var parsedData map[string]interface{}
	err := json.Unmarshal([]byte(dataJSON), &parsedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if parsedData["number"].(float64) != 42 {
		t.Error("Number not preserved")
	}

	if parsedData["string"] != "test" {
		t.Error("String not preserved")
	}

	if parsedData["bool"] != true {
		t.Error("Bool not preserved")
	}

	if parsedData["null"] != nil {
		t.Error("Null not preserved")
	}
}

// TestEventWorker_DefaultSeverity tests default severity is "info"
func TestEventWorker_DefaultSeverity(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	data := map[string]interface{}{"test": "severity"}
	worker.SendEvent("test.default_severity", "stack", "run-1", data)

	worker.mu.Lock()
	severity := worker.events[0].Severity
	worker.mu.Unlock()

	if severity != "info" {
		t.Errorf("Expected default severity 'info', got %s", severity)
	}
}

// TestEventWorker_CustomSeverityPreserved tests custom severity is preserved
func TestEventWorker_CustomSeverityPreserved(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	data := map[string]interface{}{"test": "critical_error"}
	worker.SendEventWithSeverity("test.custom_severity", "stack", "run-1", data, "critical")

	worker.mu.Lock()
	severity := worker.events[0].Severity
	worker.mu.Unlock()

	if severity != "critical" {
		t.Errorf("Expected severity 'critical', got %s", severity)
	}
}

// TestEventWorker_InitializedFieldsNonNil tests all fields are initialized
func TestEventWorker_InitializedFieldsNonNil(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	if worker.events == nil {
		t.Error("Expected events slice to be initialized")
	}

	if worker.ctx == nil {
		t.Error("Expected context to be initialized")
	}

	if worker.cancel == nil {
		t.Error("Expected cancel function to be initialized")
	}
}

// TestEventWorker_EventsSliceInitialCapacity tests initial capacity
func TestEventWorker_EventsSliceInitialCapacity(t *testing.T) {
	batchSize := 25
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  batchSize,
	}

	worker := NewEventWorker(config)

	worker.mu.Lock()
	capacity := cap(worker.events)
	worker.mu.Unlock()

	if capacity != batchSize {
		t.Errorf("Expected initial capacity %d, got %d", batchSize, capacity)
	}
}

// TestEventWorker_SendEventReturnsError tests error handling
func TestEventWorker_SendEventReturnsError(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	// Send with invalid data that can't be marshaled
	data := map[string]interface{}{
		"func": func() {}, // Functions can't be marshaled
	}

	err := worker.SendEvent("test.error", "stack", "run-1", data)
	if err == nil {
		t.Error("Expected error when marshaling function, got nil")
	}
}

// TestEventWorker_SendEventWithSeverityReturnsError tests error with severity
func TestEventWorker_SendEventWithSeverityReturnsError(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	// Send with invalid data
	data := map[string]interface{}{
		"channel": make(chan int),
	}

	err := worker.SendEventWithSeverity("test.error", "stack", "run-1", data, "error")
	if err == nil {
		t.Error("Expected error when marshaling channel, got nil")
	}
}

// TestEventWorker_LongEventType tests with very long event type
func TestEventWorker_LongEventType(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	longType := ""
	for i := 0; i < 100; i++ {
		longType += "event."
	}

	data := map[string]interface{}{"test": "long_type"}
	err := worker.SendEvent(longType, "stack", "run-1", data)
	if err != nil {
		t.Fatalf("Failed to send event with long type: %v", err)
	}

	worker.mu.Lock()
	eventType := worker.events[0].EventType
	worker.mu.Unlock()

	if eventType != longType {
		t.Error("Long event type not preserved")
	}
}

// TestEventWorker_LongStack tests with very long stack name
func TestEventWorker_LongStack(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	longStack := ""
	for i := 0; i < 200; i++ {
		longStack += "s"
	}

	data := map[string]interface{}{"test": "long_stack"}
	err := worker.SendEvent("test.stack", longStack, "run-1", data)
	if err != nil {
		t.Fatalf("Failed to send event with long stack: %v", err)
	}

	worker.mu.Lock()
	stack := worker.events[0].Stack
	worker.mu.Unlock()

	if stack != longStack {
		t.Error("Long stack not preserved")
	}
}

// TestEventWorker_LongRunID tests with very long run ID
func TestEventWorker_LongRunID(t *testing.T) {
	config := EventWorkerConfig{
		AgentName:  "test-agent",
		MasterAddr: "localhost:50051",
		BatchSize:  10,
	}

	worker := NewEventWorker(config)

	longRunID := ""
	for i := 0; i < 200; i++ {
		longRunID += "r"
	}

	data := map[string]interface{}{"test": "long_run_id"}
	err := worker.SendEvent("test.run", "stack", longRunID, data)
	if err != nil {
		t.Fatalf("Failed to send event with long run ID: %v", err)
	}

	worker.mu.Lock()
	runID := worker.events[0].RunId
	worker.mu.Unlock()

	if runID != longRunID {
		t.Error("Long run ID not preserved")
	}
}
