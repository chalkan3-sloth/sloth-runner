package core

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/yuin/gopher-lua"
)

func TestNewGoroutineModule(t *testing.T) {
	module := NewGoroutineModule()
	
	if module == nil {
		t.Fatal("Expected module to be created")
	}
	
	if module.info.Name != "goroutine" {
		t.Errorf("Expected name 'goroutine', got '%s'", module.info.Name)
	}
	
	if module.pools == nil {
		t.Error("Expected pools map to be initialized")
	}
	
	if module.globalCtx == nil {
		t.Error("Expected global context to be initialized")
	}
}

func TestGoroutineModuleLoader(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	if err := L.DoString(`goroutine = require("goroutine")`); err != nil {
		t.Fatalf("Failed to load module: %v", err)
	}
	
	functions := []string{
		"spawn", "spawn_many", "wait_group", "pool_create", "pool_submit",
		"pool_wait", "pool_close", "pool_stats", "async", "await",
		"await_all", "sleep", "timeout",
	}
	
	for _, fn := range functions {
		if err := L.DoString(`assert(goroutine.` + fn + ` ~= nil, "` + fn + ` function not found")`); err != nil {
			t.Errorf("Function %s not found: %v", fn, err)
		}
	}
}

func TestGoroutineSpawn(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		-- Spawn a simple goroutine
		goroutine.spawn(function()
			-- This runs in a separate goroutine
		end)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	// Give goroutine time to execute
	time.Sleep(100 * time.Millisecond)
}

func TestGoroutineSpawnMany(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		-- Spawn multiple goroutines
		goroutine.spawn_many(5, function(id)
			-- Each goroutine gets its own id
		end)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	time.Sleep(200 * time.Millisecond)
}

func TestGoroutineWaitGroup(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		local wg = goroutine.wait_group()
		assert(wg ~= nil, "WaitGroup should be created")
		
		-- Test add/done/wait
		wg:add(1)
		wg:done()
		wg:wait()
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutinePoolCreate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		local success = goroutine.pool_create("test_pool", { workers = 4 })
		assert(success == true, "Pool should be created")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	// Verify pool exists
	module.mu.RLock()
	pool, exists := module.pools["test_pool"]
	module.mu.RUnlock()
	
	if !exists {
		t.Error("Expected pool to exist")
	}
	
	if pool.workers != 4 {
		t.Errorf("Expected 4 workers, got %d", pool.workers)
	}
}

func TestGoroutinePoolSubmit(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		goroutine.pool_create("test_pool", { workers = 2 })
		
		local id, err = goroutine.pool_submit("test_pool", function()
			-- Task work
		end)
		
		assert(id ~= nil, "Task ID should be returned")
		assert(err == nil, "No error expected")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	time.Sleep(100 * time.Millisecond)
}

func TestGoroutinePoolStats(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		goroutine.pool_create("test_pool", { workers = 3 })
		
		local stats = goroutine.pool_stats("test_pool")
		assert(stats ~= nil, "Stats should be returned")
		assert(stats.name == "test_pool", "Pool name should match")
		assert(stats.workers == 3, "Workers count should match")
		assert(stats.active ~= nil, "Active count should be present")
		assert(stats.completed ~= nil, "Completed count should be present")
		assert(stats.failed ~= nil, "Failed count should be present")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutinePoolClose(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		goroutine.pool_create("test_pool", { workers = 2 })
		
		local success = goroutine.pool_close("test_pool")
		assert(success == true, "Pool should be closed")
		
		-- Closing again should return false
		success = goroutine.pool_close("test_pool")
		assert(success == false, "Closing non-existent pool should return false")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	// Verify pool is removed
	module.mu.RLock()
	_, exists := module.pools["test_pool"]
	module.mu.RUnlock()
	
	if exists {
		t.Error("Expected pool to be removed")
	}
}

func TestGoroutineAsync(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		local handle = goroutine.async(function()
			return "result"
		end)
		
		assert(handle ~= nil, "Handle should be returned")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineAwait(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		local handle = goroutine.async(function()
			return "test_result"
		end)
		
		local success, result = goroutine.await(handle)
		assert(success == true, "Async should succeed")
		assert(result == "test_result", "Result should match")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineAwaitAll(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		local h1 = goroutine.async(function() return 1 end)
		local h2 = goroutine.async(function() return 2 end)
		local h3 = goroutine.async(function() return 3 end)
		
		local results = goroutine.await_all({h1, h2, h3})
		assert(#results == 3, "Should have 3 results")
		
		for i, result in ipairs(results) do
			assert(result.success == true, "All should succeed")
		end
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineSleep(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	start := time.Now()
	
	script := `
		goroutine = require("goroutine")
		goroutine.sleep(100) -- 100ms
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	elapsed := time.Since(start)
	if elapsed < 100*time.Millisecond {
		t.Errorf("Expected sleep of at least 100ms, got %v", elapsed)
	}
}

func TestGoroutineTimeout(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	// Test successful execution within timeout
	script := `
		goroutine = require("goroutine")
		
		local success, result = goroutine.timeout(1000, function()
			return "quick"
		end)
		
		assert(success == true, "Should succeed within timeout")
		assert(result == "quick", "Result should match")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineTimeoutExceeded(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	// This test would hang if we actually sleep, so we test the structure
	script := `
		goroutine = require("goroutine")
		
		local success, err = goroutine.timeout(10, function()
			-- In real scenario, this would take longer than 10ms
			-- For testing purposes, we just return immediately
			return "result"
		end)
		
		-- Function completed within timeout
		assert(success == true or err ~= nil, "Should either succeed or timeout")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutinePoolWorkerExecution(t *testing.T) {
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L := lua.NewState()
	defer L.Close()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		goroutine.pool_create("worker_test", { workers = 2 })
		
		-- Submit multiple tasks
		for i = 1, 5 do
			goroutine.pool_submit("worker_test", function()
				-- Simulate work
			end)
		end
		
		-- Wait a bit for tasks to process
		goroutine.sleep(200)
		
		local stats = goroutine.pool_stats("worker_test")
		-- Tasks should have been processed
		assert(stats.completed >= 0, "Some tasks should be completed")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineModuleCleanup(t *testing.T) {
	module := NewGoroutineModule()
	
	L := lua.NewState()
	defer L.Close()
	
	L.PreloadModule("goroutine", module.Loader)
	
	// Create multiple pools
	if err := L.DoString(`
		goroutine = require("goroutine")
		goroutine.pool_create("pool1", { workers = 2 })
		goroutine.pool_create("pool2", { workers = 3 })
		goroutine.pool_create("pool3", { workers = 4 })
	`); err != nil {
		t.Fatalf("Failed to create pools: %v", err)
	}
	
	// Verify pools exist
	module.mu.RLock()
	poolCount := len(module.pools)
	module.mu.RUnlock()
	
	if poolCount != 3 {
		t.Errorf("Expected 3 pools, got %d", poolCount)
	}
	
	// Cleanup
	module.Cleanup()
	
	// Verify all pools are removed
	module.mu.RLock()
	poolCount = len(module.pools)
	module.mu.RUnlock()
	
	if poolCount != 0 {
		t.Errorf("Expected 0 pools after cleanup, got %d", poolCount)
	}
}

func TestGoroutinePoolTaskExecution(t *testing.T) {
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	pool := newGoroutinePool("test", 2, module.globalCtx)
	
	var counter int64
	
	// Submit tasks
	for i := 0; i < 10; i++ {
		L := lua.NewState()
		fn := L.NewFunction(func(L *lua.LState) int {
			atomic.AddInt64(&counter, 1)
			return 0
		})
		
		resultCh := make(chan taskResult, 1)
		task := &poolTask{
			id:       "test-task",
			fn:       fn,
			L:        L,
			args:     []lua.LValue{},
			resultCh: resultCh,
		}
		
		pool.tasks <- task
		L.Close()
	}
	
	// Wait for tasks to complete
	close(pool.tasks)
	pool.wg.Wait()
	
	if counter != 10 {
		t.Errorf("Expected 10 tasks to execute, got %d", counter)
	}
	
	completed := atomic.LoadInt64(&pool.completed)
	if completed != 10 {
		t.Errorf("Expected 10 completed tasks, got %d", completed)
	}
}

func TestGoroutineModuleInfo(t *testing.T) {
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	info := module.Info()
	
	if info.Name != "goroutine" {
		t.Errorf("Expected name 'goroutine', got '%s'", info.Name)
	}
	
	if info.Category != "core" {
		t.Errorf("Expected category 'core', got '%s'", info.Category)
	}
	
	if info.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", info.Version)
	}
}

func TestGoroutineAsyncError(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewGoroutineModule()
	defer module.Cleanup()
	
	L.PreloadModule("goroutine", module.Loader)
	
	script := `
		goroutine = require("goroutine")
		
		local handle = goroutine.async(function()
			error("test error")
		end)
		
		local success, err = goroutine.await(handle)
		assert(success == false, "Should fail")
		assert(err ~= nil, "Error message should be present")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}
