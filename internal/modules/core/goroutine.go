package core

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yuin/gopher-lua"
)

// GoroutineModule provides concurrent execution capabilities for Lua tasks
type GoroutineModule struct {
	info         CoreModuleInfo
	pools        map[string]*goroutinePool
	mu           sync.RWMutex
	globalCtx    context.Context
	globalCancel context.CancelFunc
}

// goroutinePool manages a pool of goroutines
type goroutinePool struct {
	name      string
	workers   int
	tasks     chan *poolTask
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	active    int64
	completed int64
	failed    int64
	results   sync.Map
}

// poolTask represents a task to be executed
type poolTask struct {
	id       string
	fn       *lua.LFunction
	L        *lua.LState
	args     []lua.LValue
	resultCh chan taskResult
}

// taskResult holds the result of a task execution
type taskResult struct {
	success bool
	values  []lua.LValue
	err     error
}

// NewGoroutineModule creates a new goroutine module
func NewGoroutineModule() *GoroutineModule {
	ctx, cancel := context.WithCancel(context.Background())
	
	info := CoreModuleInfo{
		Name:        "goroutine",
		Version:     "1.0.0",
		Description: "Concurrent execution with goroutines, worker pools, and async operations",
		Author:      "Sloth Runner Team",
		Category:    "core",
		Dependencies: []string{},
	}
	
	return &GoroutineModule{
		info:         info,
		pools:        make(map[string]*goroutinePool),
		globalCtx:    ctx,
		globalCancel: cancel,
	}
}

// Info returns module information
func (g *GoroutineModule) Info() CoreModuleInfo {
	return g.info
}

// Loader returns the Lua loader function
func (g *GoroutineModule) Loader(L *lua.LState) int {
	goroutineTable := L.NewTable()
	
	// Register functions
	L.SetFuncs(goroutineTable, map[string]lua.LGFunction{
		"spawn":       g.luaSpawn,
		"spawn_many":  g.luaSpawnMany,
		"wait_group":  g.luaWaitGroup,
		"pool_create": g.luaPoolCreate,
		"pool_submit": g.luaPoolSubmit,
		"pool_wait":   g.luaPoolWait,
		"pool_close":  g.luaPoolClose,
		"pool_stats":  g.luaPoolStats,
		"async":       g.luaAsync,
		"await":       g.luaAwait,
		"await_all":   g.luaAwaitAll,
		"sleep":       g.luaSleep,
		"timeout":     g.luaTimeout,
	})
	
	L.Push(goroutineTable)
	return 1
}

// luaSpawn spawns a single goroutine
// Usage: goroutine.spawn(function() ... end)
func (g *GoroutineModule) luaSpawn(L *lua.LState) int {
	fn := L.CheckFunction(1)
	
	// Create a new Lua state for the goroutine
	newL := lua.NewState()
	
	go func() {
		defer newL.Close()
		defer func() {
			if r := recover(); r != nil {
				// Log panic but don't crash
				fmt.Printf("Goroutine panic: %v\n", r)
			}
		}()
		
		// Copy the function to the new state
		if err := newL.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}); err != nil {
			fmt.Printf("Goroutine error: %v\n", err)
		}
	}()
	
	return 0
}

// luaSpawnMany spawns multiple goroutines
// Usage: goroutine.spawn_many(count, function(id) ... end)
func (g *GoroutineModule) luaSpawnMany(L *lua.LState) int {
	count := L.CheckInt(1)
	fn := L.CheckFunction(2)
	
	var wg sync.WaitGroup
	for i := 1; i <= count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			newL := lua.NewState()
			defer newL.Close()
			
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Goroutine %d panic: %v\n", id, r)
				}
			}()
			
			if err := newL.CallByParam(lua.P{
				Fn:      fn,
				NRet:    0,
				Protect: true,
			}, lua.LNumber(id)); err != nil {
				fmt.Printf("Goroutine %d error: %v\n", id, err)
			}
		}(i)
	}
	
	// Wait in a separate goroutine to not block Lua
	go wg.Wait()
	
	return 0
}

// luaWaitGroup creates a WaitGroup for synchronization
// Usage: local wg = goroutine.wait_group(); wg.add(1); wg.done(); wg.wait()
func (g *GoroutineModule) luaWaitGroup(L *lua.LState) int {
	wg := &sync.WaitGroup{}
	
	ud := L.NewUserData()
	ud.Value = wg
	
	// Create metatable with methods
	mt := L.NewTable()
	L.SetField(mt, "add", L.NewFunction(func(L *lua.LState) int {
		wg := L.CheckUserData(1).Value.(*sync.WaitGroup)
		delta := L.OptInt(2, 1)
		wg.Add(delta)
		return 0
	}))
	
	L.SetField(mt, "done", L.NewFunction(func(L *lua.LState) int {
		wg := L.CheckUserData(1).Value.(*sync.WaitGroup)
		wg.Done()
		return 0
	}))
	
	L.SetField(mt, "wait", L.NewFunction(func(L *lua.LState) int {
		wg := L.CheckUserData(1).Value.(*sync.WaitGroup)
		wg.Wait()
		return 0
	}))
	
	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)
	
	L.Push(ud)
	return 1
}

// luaPoolCreate creates a worker pool
// Usage: goroutine.pool_create("mypool", { workers = 10 })
func (g *GoroutineModule) luaPoolCreate(L *lua.LState) int {
	name := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	workers := 4 // default
	if w := options.RawGetString("workers"); w != lua.LNil {
		workers = int(w.(lua.LNumber))
	}
	
	g.mu.Lock()
	defer g.mu.Unlock()
	
	// Close existing pool if it exists
	if existing, exists := g.pools[name]; exists {
		existing.close()
	}
	
	// Create new pool
	pool := newGoroutinePool(name, workers, g.globalCtx)
	g.pools[name] = pool
	
	L.Push(lua.LTrue)
	return 1
}

// luaPoolSubmit submits a task to a worker pool
// Usage: local id = goroutine.pool_submit("mypool", function() ... end)
func (g *GoroutineModule) luaPoolSubmit(L *lua.LState) int {
	poolName := L.CheckString(1)
	fn := L.CheckFunction(2)
	
	g.mu.RLock()
	pool, exists := g.pools[poolName]
	g.mu.RUnlock()
	
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("pool not found: %s", poolName)))
		return 2
	}
	
	// Generate task ID
	taskID := fmt.Sprintf("%s-%d", poolName, time.Now().UnixNano())
	
	// Collect arguments
	args := make([]lua.LValue, 0)
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.Get(i))
	}
	
	// Create result channel
	resultCh := make(chan taskResult, 1)
	
	task := &poolTask{
		id:       taskID,
		fn:       fn,
		L:        L,
		args:     args,
		resultCh: resultCh,
	}
	
	// Submit to pool
	select {
	case pool.tasks <- task:
		L.Push(lua.LString(taskID))
		return 1
	case <-pool.ctx.Done():
		L.Push(lua.LNil)
		L.Push(lua.LString("pool is closed"))
		return 2
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString("pool queue is full"))
		return 2
	}
}

// luaPoolWait waits for a pool to complete all tasks
// Usage: goroutine.pool_wait("mypool")
func (g *GoroutineModule) luaPoolWait(L *lua.LState) int {
	poolName := L.CheckString(1)
	
	g.mu.RLock()
	pool, exists := g.pools[poolName]
	g.mu.RUnlock()
	
	if !exists {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("pool not found: %s", poolName)))
		return 2
	}
	
	// Close tasks channel and wait
	close(pool.tasks)
	pool.wg.Wait()
	
	L.Push(lua.LTrue)
	return 1
}

// luaPoolClose closes a worker pool
// Usage: goroutine.pool_close("mypool")
func (g *GoroutineModule) luaPoolClose(L *lua.LState) int {
	poolName := L.CheckString(1)
	
	g.mu.Lock()
	pool, exists := g.pools[poolName]
	if exists {
		delete(g.pools, poolName)
	}
	g.mu.Unlock()
	
	if !exists {
		L.Push(lua.LFalse)
		return 1
	}
	
	pool.close()
	L.Push(lua.LTrue)
	return 1
}

// luaPoolStats returns statistics for a pool
// Usage: local stats = goroutine.pool_stats("mypool")
func (g *GoroutineModule) luaPoolStats(L *lua.LState) int {
	poolName := L.CheckString(1)
	
	g.mu.RLock()
	pool, exists := g.pools[poolName]
	g.mu.RUnlock()
	
	if !exists {
		L.Push(lua.LNil)
		return 1
	}
	
	stats := L.NewTable()
	stats.RawSetString("name", lua.LString(pool.name))
	stats.RawSetString("workers", lua.LNumber(pool.workers))
	stats.RawSetString("active", lua.LNumber(atomic.LoadInt64(&pool.active)))
	stats.RawSetString("completed", lua.LNumber(atomic.LoadInt64(&pool.completed)))
	stats.RawSetString("failed", lua.LNumber(atomic.LoadInt64(&pool.failed)))
	stats.RawSetString("queued", lua.LNumber(len(pool.tasks)))
	
	L.Push(stats)
	return 1
}

// luaAsync executes a function asynchronously and returns a promise-like handle
// Usage: local handle = goroutine.async(function() return "result" end)
func (g *GoroutineModule) luaAsync(L *lua.LState) int {
	fn := L.CheckFunction(1)
	
	// Create result channel
	resultCh := make(chan taskResult, 1)
	taskID := fmt.Sprintf("async-%d", time.Now().UnixNano())
	
	// Execute in goroutine
	go func() {
		newL := lua.NewState()
		defer newL.Close()
		
		defer func() {
			if r := recover(); r != nil {
				resultCh <- taskResult{
					success: false,
					err:     fmt.Errorf("panic: %v", r),
				}
			}
		}()
		
		// Call function
		err := newL.CallByParam(lua.P{
			Fn:      fn,
			NRet:    lua.MultRet,
			Protect: true,
		})
		
		if err != nil {
			resultCh <- taskResult{
				success: false,
				err:     err,
			}
			return
		}
		
		// Collect results
		results := make([]lua.LValue, newL.GetTop())
		for i := 1; i <= newL.GetTop(); i++ {
			results[i-1] = newL.Get(i)
		}
		
		resultCh <- taskResult{
			success: true,
			values:  results,
		}
	}()
	
	// Create handle userdata
	handle := &asyncHandle{
		id:       taskID,
		resultCh: resultCh,
		done:     false,
	}
	
	ud := L.NewUserData()
	ud.Value = handle
	
	// Create metatable
	mt := L.NewTable()
	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)
	
	L.Push(ud)
	return 1
}

// asyncHandle represents a handle to an async operation
type asyncHandle struct {
	id       string
	resultCh chan taskResult
	result   taskResult
	done     bool
	mu       sync.Mutex
}

// luaAwait waits for an async operation to complete
// Usage: local success, result = goroutine.await(handle)
func (g *GoroutineModule) luaAwait(L *lua.LState) int {
	ud := L.CheckUserData(1)
	handle := ud.Value.(*asyncHandle)
	
	handle.mu.Lock()
	defer handle.mu.Unlock()
	
	if !handle.done {
		handle.result = <-handle.resultCh
		handle.done = true
	}
	
	if handle.result.err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(handle.result.err.Error()))
		return 2
	}
	
	L.Push(lua.LTrue)
	for _, val := range handle.result.values {
		L.Push(val)
	}
	return 1 + len(handle.result.values)
}

// luaAwaitAll waits for all async operations to complete
// Usage: local results = goroutine.await_all({handle1, handle2, ...})
func (g *GoroutineModule) luaAwaitAll(L *lua.LState) int {
	handles := L.CheckTable(1)
	
	results := L.NewTable()
	index := 1
	
	handles.ForEach(func(k, v lua.LValue) {
		if ud, ok := v.(*lua.LUserData); ok {
			if handle, ok := ud.Value.(*asyncHandle); ok {
				handle.mu.Lock()
				if !handle.done {
					handle.result = <-handle.resultCh
					handle.done = true
				}
				handle.mu.Unlock()
				
				resultTable := L.NewTable()
				resultTable.RawSetString("success", lua.LBool(handle.result.err == nil))
				if handle.result.err != nil {
					resultTable.RawSetString("error", lua.LString(handle.result.err.Error()))
				} else {
					values := L.NewTable()
					for i, val := range handle.result.values {
						values.RawSetInt(i+1, val)
					}
					resultTable.RawSetString("values", values)
				}
				
				results.RawSetInt(index, resultTable)
				index++
			}
		}
	})
	
	L.Push(results)
	return 1
}

// luaSleep sleeps for a specified duration
// Usage: goroutine.sleep(1000) -- milliseconds
func (g *GoroutineModule) luaSleep(L *lua.LState) int {
	ms := L.CheckInt(1)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return 0
}

// luaTimeout executes a function with a timeout
// Usage: local success, result = goroutine.timeout(1000, function() ... end)
func (g *GoroutineModule) luaTimeout(L *lua.LState) int {
	ms := L.CheckInt(1)
	fn := L.CheckFunction(2)
	
	ctx, cancel := context.WithTimeout(g.globalCtx, time.Duration(ms)*time.Millisecond)
	defer cancel()
	
	resultCh := make(chan taskResult, 1)
	
	go func() {
		newL := lua.NewState()
		defer newL.Close()
		
		defer func() {
			if r := recover(); r != nil {
				resultCh <- taskResult{
					success: false,
					err:     fmt.Errorf("panic: %v", r),
				}
			}
		}()
		
		err := newL.CallByParam(lua.P{
			Fn:      fn,
			NRet:    lua.MultRet,
			Protect: true,
		})
		
		if err != nil {
			resultCh <- taskResult{
				success: false,
				err:     err,
			}
			return
		}
		
		results := make([]lua.LValue, newL.GetTop())
		for i := 1; i <= newL.GetTop(); i++ {
			results[i-1] = newL.Get(i)
		}
		
		resultCh <- taskResult{
			success: true,
			values:  results,
		}
	}()
	
	select {
	case result := <-resultCh:
		if result.err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(result.err.Error()))
			return 2
		}
		L.Push(lua.LTrue)
		for _, val := range result.values {
			L.Push(val)
		}
		return 1 + len(result.values)
	case <-ctx.Done():
		L.Push(lua.LFalse)
		L.Push(lua.LString("timeout exceeded"))
		return 2
	}
}

// Helper functions for pool management

func newGoroutinePool(name string, workers int, parentCtx context.Context) *goroutinePool {
	ctx, cancel := context.WithCancel(parentCtx)
	
	pool := &goroutinePool{
		name:    name,
		workers: workers,
		tasks:   make(chan *poolTask, workers*2),
		ctx:     ctx,
		cancel:  cancel,
	}
	
	// Start workers
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}
	
	return pool
}

func (p *goroutinePool) worker(id int) {
	defer p.wg.Done()
	
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			
			atomic.AddInt64(&p.active, 1)
			p.executeTask(task)
			atomic.AddInt64(&p.active, -1)
			
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *goroutinePool) executeTask(task *poolTask) {
	defer func() {
		if r := recover(); r != nil {
			atomic.AddInt64(&p.failed, 1)
			task.resultCh <- taskResult{
				success: false,
				err:     fmt.Errorf("panic: %v", r),
			}
		}
	}()
	
	// Create new Lua state for task
	newL := lua.NewState()
	defer newL.Close()
	
	// Execute function
	err := newL.CallByParam(lua.P{
		Fn:      task.fn,
		NRet:    lua.MultRet,
		Protect: true,
	}, task.args...)
	
	if err != nil {
		atomic.AddInt64(&p.failed, 1)
		task.resultCh <- taskResult{
			success: false,
			err:     err,
		}
		return
	}
	
	// Collect results
	results := make([]lua.LValue, newL.GetTop())
	for i := 1; i <= newL.GetTop(); i++ {
		results[i-1] = newL.Get(i)
	}
	
	atomic.AddInt64(&p.completed, 1)
	task.resultCh <- taskResult{
		success: true,
		values:  results,
	}
}

func (p *goroutinePool) close() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
}

// Cleanup closes all pools
func (g *GoroutineModule) Cleanup() {
	g.globalCancel()
	
	g.mu.Lock()
	defer g.mu.Unlock()
	
	for name, pool := range g.pools {
		pool.close()
		delete(g.pools, name)
	}
}
