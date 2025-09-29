package core

import (
	"testing"
	"log/slog"
	"time"
	"os"
)

func TestGlobalCore(t *testing.T) {
	// Create test logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	
	// Test core initialization
	config := DefaultCoreConfig()
	config.MaxWorkers = 2
	config.MaxMemoryMB = 10
	config.CacheSizeMB = 5
	
	core, err := NewGlobalCore(config, logger)
	if err != nil {
		t.Fatalf("Failed to create core: %v", err)
	}
	
	// Test start
	if err := core.Start(); err != nil {
		t.Fatalf("Failed to start core: %v", err)
	}
	
	// Test task submission
	taskComplete := make(chan bool)
	success := core.SubmitTask(func() {
		time.Sleep(100 * time.Millisecond)
		taskComplete <- true
	}, "test_task")
	
	if !success {
		t.Error("Failed to submit task")
	}
	
	// Wait for task completion
	select {
	case <-taskComplete:
		// Task completed successfully
	case <-time.After(5 * time.Second):
		t.Error("Task did not complete within timeout")
	}
	
	// Test stats
	stats := core.GetStats()
	if !stats.Started {
		t.Error("Core should report as started")
	}
	
	if stats.WorkerPool.Workers != 2 {
		t.Errorf("Expected 2 workers, got %d", stats.WorkerPool.Workers)
	}
	
	// Test circuit breaker
	cb := core.GetCircuitBreaker("test_service")
	if cb == nil {
		t.Error("Expected circuit breaker to be created")
	}
	
	// Test error recovery
	err = core.ExecuteWithRecovery(func() error {
		panic("test panic")
	}, "test_panic")
	
	if err == nil {
		t.Error("Expected error from panic recovery")
	}
	
	// Test stop
	if err := core.Stop(); err != nil {
		t.Fatalf("Failed to stop core: %v", err)
	}
}

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker("test", 2, 100*time.Millisecond)
	
	// Test initial state (closed)
	stats := cb.GetStats()
	if stats.State != "closed" {
		t.Errorf("Expected initial state 'closed', got '%s'", stats.State)
	}
	
	// Test successful execution
	err := cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Test failures to open circuit
	for i := 0; i < 3; i++ {
		cb.Execute(func() error {
			return NewSlothError("TEST_ERROR", "test error", SeverityMedium)
		})
	}
	
	stats = cb.GetStats()
	if stats.State != "open" {
		t.Errorf("Expected state 'open' after failures, got '%s'", stats.State)
	}
	
	// Test execution blocked when open
	err = cb.Execute(func() error {
		return nil
	})
	if err == nil {
		t.Error("Expected error when circuit is open")
	}
}

func TestResourcePool(t *testing.T) {
	pool := NewResourcePool(10)
	
	// Test buffer pool
	buf1 := pool.GetBuffer()
	if buf1 == nil {
		t.Error("Expected buffer from pool")
	}
	
	if cap(buf1) != 4096 {
		t.Errorf("Expected buffer capacity 4096, got %d", cap(buf1))
	}
	
	// Use buffer and return
	buf1 = append(buf1, []byte("test data")...)
	pool.PutBuffer(buf1)
	
	// Get another buffer (should be reused)
	buf2 := pool.GetBuffer()
	if len(buf2) != 0 {
		t.Errorf("Expected empty buffer, got length %d", len(buf2))
	}
}

func TestWorkerPool(t *testing.T) {
	wp := NewWorkerPool(2)
	defer wp.Close()
	
	// Test task submission and execution
	completed := make(chan int, 5)
	
	for i := 0; i < 5; i++ {
		taskNum := i
		success := wp.Submit(func() {
			time.Sleep(10 * time.Millisecond)
			completed <- taskNum
		})
		
		if !success {
			t.Errorf("Failed to submit task %d", i)
		}
	}
	
	// Wait for all tasks to complete
	results := make(map[int]bool)
	for i := 0; i < 5; i++ {
		select {
		case taskNum := <-completed:
			results[taskNum] = true
		case <-time.After(5 * time.Second):
			t.Fatal("Tasks did not complete within timeout")
		}
	}
	
	// Verify all tasks completed
	if len(results) != 5 {
		t.Errorf("Expected 5 completed tasks, got %d", len(results))
	}
	
	// Check stats
	stats := wp.Stats()
	if stats.Workers != 2 {
		t.Errorf("Expected 2 workers, got %d", stats.Workers)
	}
	
	if stats.Completed < 5 {
		t.Errorf("Expected at least 5 completed tasks, got %d", stats.Completed)
	}
}

func TestErrorHandling(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	
	// Test error recovery
	recovery := NewErrorRecovery(logger)
	
	// Test normal execution
	recovered, err := recovery.SafeExecute(func() error {
		return NewSlothError("TEST_ERROR", "test error", SeverityLow)
	})
	
	if recovered {
		t.Error("Expected no panic recovery")
	}
	
	if err == nil {
		t.Error("Expected error to be returned")
	}
	
	// Test panic recovery
	recovered, err = recovery.SafeExecute(func() error {
		panic("test panic")
	})
	
	if !recovered {
		t.Error("Expected panic to be recovered")
	}
	
	if err == nil {
		t.Error("Expected error from panic recovery")
	}
	
	// Test error collector
	collector := NewErrorCollector(10)
	
	testErr := NewSlothError("TEST", "test", SeverityMedium)
	collector.Collect(testErr)
	
	if !collector.HasErrors() {
		t.Error("Expected collector to have errors")
	}
	
	errors := collector.GetErrors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
	
	if errors[0].Code != "TEST" {
		t.Errorf("Expected error code 'TEST', got '%s'", errors[0].Code)
	}
}

func TestSafeMap(t *testing.T) {
	sm := NewSafeMap()
	
	// Test set and get
	sm.Set("key1", "value1")
	sm.Set("key2", 42)
	
	val1, exists1 := sm.Get("key1")
	if !exists1 || val1 != "value1" {
		t.Errorf("Expected 'value1', got %v (exists: %v)", val1, exists1)
	}
	
	val2, exists2 := sm.Get("key2")
	if !exists2 || val2 != 42 {
		t.Errorf("Expected 42, got %v (exists: %v)", val2, exists2)
	}
	
	// Test non-existent key
	_, exists3 := sm.Get("key3")
	if exists3 {
		t.Error("Expected key3 to not exist")
	}
	
	// Test keys
	keys := sm.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
	
	// Test delete
	sm.Delete("key1")
	_, exists4 := sm.Get("key1")
	if exists4 {
		t.Error("Expected key1 to be deleted")
	}
	
	// Test clear
	sm.Clear()
	if sm.Len() != 0 {
		t.Errorf("Expected empty map after clear, got length %d", sm.Len())
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	wp := NewWorkerPool(4)
	defer wp.Close()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		wp.Submit(func() {
			// Simulate light work
			time.Sleep(time.Microsecond)
		})
	}
	
	// Wait for all tasks to complete
	for {
		stats := wp.Stats()
		if stats.Queued == 0 && stats.Active == 0 {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

func BenchmarkSafeMap(b *testing.B) {
	sm := NewSafeMap()
	
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "key" + string(rune(i%100))
			sm.Set(key, i)
			sm.Get(key)
			i++
		}
	})
}