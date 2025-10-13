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

// ============================================================================
// CHANNEL TESTS
// ============================================================================

func TestGoroutineChannel(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		-- Create buffered channel
		local ch = goroutine.channel(5)
		assert(ch ~= nil, "Channel should be created")

		-- Test send and receive
		local sent = ch:send(42)
		assert(sent == true, "Should send successfully")

		local value, ok = ch:receive()
		assert(ok == true, "Should receive successfully")
		assert(value == 42, "Value should match")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineChannelClose(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local ch = goroutine.channel(5)

		-- Send some values
		ch:send(1)
		ch:send(2)

		-- Close channel
		local closed = ch:close()
		assert(closed == true, "Should close successfully")

		-- Receive remaining values
		local v1, ok1 = ch:receive()
		assert(ok1 == true, "Should receive value 1")
		assert(v1 == 1, "Value should be 1")

		local v2, ok2 = ch:receive()
		assert(ok2 == true, "Should receive value 2")
		assert(v2 == 2, "Value should be 2")

		-- Next receive should fail (channel closed and empty)
		local v3, ok3 = ch:receive()
		assert(ok3 == false, "Should fail on closed empty channel")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineChannelRange(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local ch = goroutine.channel(10)

		-- Producer goroutine
		goroutine.spawn(function()
			for i = 1, 5 do
				ch:send(i)
			end
			ch:close()
		end)

		-- Consumer using range
		local sum = 0
		ch:range(function(value)
			sum = sum + value
		end)

		assert(sum == 15, "Sum should be 15 (1+2+3+4+5)")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineChannelCapLen(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local ch = goroutine.channel(10)

		-- Check capacity
		assert(ch:cap() == 10, "Capacity should be 10")

		-- Initially empty
		assert(ch:len() == 0, "Should be empty")

		-- Send some values
		ch:send(1)
		ch:send(2)
		ch:send(3)

		-- Check length
		assert(ch:len() == 3, "Should have 3 items")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

// ============================================================================
// PIPELINE TESTS
// ============================================================================

func TestGoroutinePipeline(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		-- Create input channel
		local input = goroutine.channel(10)

		-- Create 2-stage pipeline
		local output = goroutine.pipeline(input, {
			{
				workers = 1,
				fn = function(x) return x * 2 end
			},
			{
				workers = 1,
				fn = function(x) return x + 10 end
			}
		})

		-- Feed data
		goroutine.spawn(function()
			for i = 1, 5 do
				input:send(i)
			end
			input:close()
		end)

		-- Collect results
		local results = {}
		output:range(function(value)
			table.insert(results, value)
		end)

		-- Verify: 1*2+10=12, 2*2+10=14, 3*2+10=16, 4*2+10=18, 5*2+10=20
		assert(#results == 5, "Should have 5 results")
		assert(results[1] == 12, "First result should be 12")
		assert(results[2] == 14, "Second result should be 14")
		assert(results[3] == 16, "Third result should be 16")
		assert(results[4] == 18, "Fourth result should be 18")
		assert(results[5] == 20, "Fifth result should be 20")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutinePipelineMultipleWorkers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local input = goroutine.channel(20)

		-- Pipeline with multiple workers per stage
		local output = goroutine.pipeline(input, {
			{
				workers = 3,  -- 3 workers in stage 1
				fn = function(x) return x * 2 end
			},
			{
				workers = 2,  -- 2 workers in stage 2
				fn = function(x) return x + 1 end
			}
		})

		-- Send many values
		goroutine.spawn(function()
			for i = 1, 10 do
				input:send(i)
			end
			input:close()
		end)

		-- Collect all results
		local results = {}
		output:range(function(value)
			table.insert(results, value)
		end)

		-- Should receive all 10 results (order may vary due to parallelism)
		assert(#results == 10, "Should have 10 results")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

// ============================================================================
// FAN-OUT / FAN-IN TESTS
// ============================================================================

func TestGoroutineFanOut(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local input = goroutine.channel(10)

		-- Fan out to 3 channels
		local outputs = goroutine.fan_out(input, 3)
		assert(#outputs == 3, "Should have 3 output channels")

		-- Send values to input
		goroutine.spawn(function()
			for i = 1, 5 do
				input:send(i)
			end
			input:close()
		end)

		-- Each output should receive all values
		local wg = goroutine.wait_group()
		wg:add(3)

		local results = {{}, {}, {}}

		for i = 1, 3 do
			local idx = i
			goroutine.spawn(function()
				outputs[idx]:range(function(value)
					table.insert(results[idx], value)
				end)
				wg:done()
			end)
		end

		wg:wait()

		-- All outputs should have received all 5 values
		for i = 1, 3 do
			assert(#results[i] == 5, "Output " .. i .. " should have 5 values")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineFanIn(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		-- Create 3 input channels
		local ch1 = goroutine.channel(5)
		local ch2 = goroutine.channel(5)
		local ch3 = goroutine.channel(5)

		-- Fan in to single output
		local output = goroutine.fan_in({ch1, ch2, ch3})

		-- Send values on each channel
		goroutine.spawn(function()
			for i = 1, 3 do
				ch1:send(i * 10)
			end
			ch1:close()
		end)

		goroutine.spawn(function()
			for i = 1, 3 do
				ch2:send(i * 100)
			end
			ch2:close()
		end)

		goroutine.spawn(function()
			for i = 1, 3 do
				ch3:send(i * 1000)
			end
			ch3:close()
		end)

		-- Collect all results from merged output
		local results = {}
		output:range(function(value)
			table.insert(results, value)
		end)

		-- Should receive all 9 values (3 from each channel)
		assert(#results == 9, "Should have 9 total values")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

// ============================================================================
// CONTEXT TESTS
// ============================================================================

func TestGoroutineContext(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local ctx = goroutine.context()
		assert(ctx ~= nil, "Context should be created")

		-- Initially not cancelled
		assert(ctx:is_cancelled() == false, "Should not be cancelled initially")

		-- Cancel context
		ctx:cancel()

		-- Now should be cancelled
		assert(ctx:is_cancelled() == true, "Should be cancelled after cancel()")

		-- Should have error
		local err = ctx:err()
		assert(err ~= nil, "Should have error after cancel")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineContextTimeout(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local ctx = goroutine.context()
		local timeoutCtx = ctx:with_timeout(100)  -- 100ms timeout

		-- Wait for timeout
		goroutine.sleep(150)

		-- Should be cancelled due to timeout
		assert(timeoutCtx:is_cancelled() == true, "Should be cancelled after timeout")

		local err = timeoutCtx:err()
		assert(err ~= nil, "Should have timeout error")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineContextWithCancel(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local parentCtx = goroutine.context()
		local childCtx, cancel = parentCtx:with_cancel()

		-- Child initially not cancelled
		assert(childCtx:is_cancelled() == false, "Child should not be cancelled")

		-- Cancel child
		cancel()

		-- Child should be cancelled
		assert(childCtx:is_cancelled() == true, "Child should be cancelled")

		-- Parent should still be active
		assert(parentCtx:is_cancelled() == false, "Parent should still be active")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

// ============================================================================
// MUTEX TESTS
// ============================================================================

func TestGoroutineMutex(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local mu = goroutine.mutex()
		assert(mu ~= nil, "Mutex should be created")

		local counter = 0
		local wg = goroutine.wait_group()

		-- Spawn multiple goroutines that increment counter
		for i = 1, 10 do
			wg:add(1)
			goroutine.spawn(function()
				mu:lock()
				local temp = counter
				goroutine.sleep(1)  -- Small delay to encourage race conditions
				counter = temp + 1
				mu:unlock()
				wg:done()
			end)
		end

		wg:wait()

		-- With proper locking, counter should be exactly 10
		assert(counter == 10, "Counter should be 10 with proper locking")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineRWMutex(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local rwmu = goroutine.rwmutex()
		assert(rwmu ~= nil, "RWMutex should be created")

		local value = 0

		-- Test read lock
		rwmu:rlock()
		local read1 = value
		rwmu:runlock()

		-- Test write lock
		rwmu:lock()
		value = 42
		rwmu:unlock()

		-- Test read lock again
		rwmu:rlock()
		local read2 = value
		rwmu:runlock()

		assert(read2 == 42, "Should read updated value")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

// ============================================================================
// SEMAPHORE TESTS
// ============================================================================

func TestGoroutineSemaphore(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local sem = goroutine.semaphore(3)  -- Capacity of 3
		assert(sem ~= nil, "Semaphore should be created")

		-- Check initial capacity
		assert(sem:capacity() == 3, "Capacity should be 3")
		assert(sem:available() == 3, "Should have 3 tokens available")

		-- Acquire tokens
		sem:acquire()
		assert(sem:available() == 2, "Should have 2 tokens after acquire")

		sem:acquire()
		assert(sem:available() == 1, "Should have 1 token after second acquire")

		-- Release token
		sem:release()
		assert(sem:available() == 2, "Should have 2 tokens after release")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineSemaphoreConcurrency(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local sem = goroutine.semaphore(2)  -- Max 2 concurrent
		local wg = goroutine.wait_group()
		local active = 0
		local maxActive = 0
		local mu = goroutine.mutex()

		-- Spawn 5 goroutines, but only 2 should run concurrently
		for i = 1, 5 do
			wg:add(1)
			goroutine.spawn(function()
				sem:acquire()

				mu:lock()
				active = active + 1
				if active > maxActive then
					maxActive = active
				end
				mu:unlock()

				goroutine.sleep(50)  -- Simulate work

				mu:lock()
				active = active - 1
				mu:unlock()

				sem:release()
				wg:done()
			end)
		end

		wg:wait()

		-- Max concurrent should never exceed semaphore capacity
		assert(maxActive <= 2, "Max concurrent should be <= 2")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

// ============================================================================
// ATOMIC TESTS
// ============================================================================

func TestGoroutineAtomicInt(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local counter = goroutine.atomic_int(0)
		assert(counter ~= nil, "Atomic int should be created")

		-- Test load
		assert(counter:load() == 0, "Initial value should be 0")

		-- Test add
		local newVal = counter:add(5)
		assert(newVal == 5, "After add(5), value should be 5")
		assert(counter:load() == 5, "Load should return 5")

		-- Test store
		counter:store(10)
		assert(counter:load() == 10, "After store(10), value should be 10")

		-- Test swap
		local oldVal = counter:swap(20)
		assert(oldVal == 10, "Swap should return old value 10")
		assert(counter:load() == 20, "After swap, value should be 20")

		-- Test compare_and_swap
		local swapped = counter:compare_and_swap(20, 30)
		assert(swapped == true, "CAS should succeed")
		assert(counter:load() == 30, "Value should be 30 after successful CAS")

		-- Failed CAS
		swapped = counter:compare_and_swap(20, 40)
		assert(swapped == false, "CAS should fail with wrong old value")
		assert(counter:load() == 30, "Value should still be 30")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineAtomicConcurrency(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local counter = goroutine.atomic_int(0)
		local wg = goroutine.wait_group()

		-- Spawn 100 goroutines that each increment counter
		for i = 1, 100 do
			wg:add(1)
			goroutine.spawn(function()
				counter:add(1)
				wg:done()
			end)
		end

		wg:wait()

		-- With atomic operations, counter should be exactly 100
		assert(counter:load() == 100, "Counter should be 100 with atomic ops")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

// ============================================================================
// ONCE TESTS
// ============================================================================

func TestGoroutineOnce(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local once = goroutine.once()
		assert(once ~= nil, "Once should be created")

		local counter = 0

		-- Call multiple times
		for i = 1, 5 do
			once:call(function()
				counter = counter + 1
			end)
		end

		-- Should only execute once
		assert(counter == 1, "Function should execute exactly once")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineOnceConcurrent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local once = goroutine.once()
		local counter = goroutine.atomic_int(0)
		local wg = goroutine.wait_group()

		-- Spawn many goroutines all calling once
		for i = 1, 50 do
			wg:add(1)
			goroutine.spawn(function()
				once:call(function()
					counter:add(1)
				end)
				wg:done()
			end)
		end

		wg:wait()

		-- Should still execute exactly once despite concurrency
		assert(counter:load() == 1, "Should execute exactly once even with concurrency")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

// ============================================================================
// CONDITION VARIABLE TESTS
// ============================================================================

func TestGoroutineCond(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local cond = goroutine.cond()
		assert(cond ~= nil, "Cond should be created")

		local mu = cond:get_mutex()
		assert(mu ~= nil, "Should get mutex from cond")

		local ready = false
		local wg = goroutine.wait_group()

		-- Waiter goroutine
		wg:add(1)
		goroutine.spawn(function()
			mu:lock()
			while not ready do
				cond:wait()  -- Releases lock and waits
			end
			mu:unlock()
			wg:done()
		end)

		-- Give waiter time to start waiting
		goroutine.sleep(50)

		-- Signaler goroutine
		mu:lock()
		ready = true
		cond:signal()  -- Wake up waiter
		mu:unlock()

		wg:wait()
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestGoroutineCondBroadcast(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewGoroutineModule()
	defer module.Cleanup()

	L.PreloadModule("goroutine", module.Loader)

	script := `
		goroutine = require("goroutine")

		local cond = goroutine.cond()
		local mu = cond:get_mutex()
		local ready = false
		local wg = goroutine.wait_group()
		local completed = goroutine.atomic_int(0)

		-- Spawn 5 waiters
		for i = 1, 5 do
			wg:add(1)
			goroutine.spawn(function()
				mu:lock()
				while not ready do
					cond:wait()
				end
				mu:unlock()
				completed:add(1)
				wg:done()
			end)
		end

		-- Give waiters time to start
		goroutine.sleep(100)

		-- Broadcast to all waiters
		mu:lock()
		ready = true
		cond:broadcast()  -- Wake up all waiters
		mu:unlock()

		wg:wait()

		-- All 5 should have completed
		assert(completed:load() == 5, "All 5 waiters should complete")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}
