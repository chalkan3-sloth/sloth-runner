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
	channels     map[string]*luaChannel
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

// luaChannel wraps a Go channel for use in Lua
type luaChannel struct {
	ch        chan lua.LValue
	capacity  int
	closed    bool
	closeMu   sync.Mutex
	direction string // "bidirectional", "send", "receive"
}

// selectCase represents a case in a select statement
type selectCase struct {
	caseType string // "send", "receive", "default"
	channel  *luaChannel
	value    lua.LValue // for send operations
	fn       *lua.LFunction
}

// luaSemaphore wraps a semaphore for use in Lua
type luaSemaphore struct {
	ch chan struct{}
	capacity int
}

// luaAtomicInt64 wraps an atomic int64 for use in Lua
type luaAtomicInt64 struct {
	value int64
}

// luaOnce wraps sync.Once for use in Lua
type luaOnce struct {
	once sync.Once
}

// luaCond wraps sync.Cond for use in Lua
type luaCond struct {
	cond *sync.Cond
	mu   *sync.Mutex
}

// luaContext wraps context.Context for use in Lua
type luaContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// pipelineStage represents a stage in a processing pipeline
type pipelineStage struct {
	fn       *lua.LFunction
	workers  int
	inputCh  *luaChannel
	outputCh *luaChannel
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
		channels:     make(map[string]*luaChannel),
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
		// Channel operations
		"channel":        g.luaChannelMake,
		"channel_send":   g.luaChannelSend,
		"channel_receive": g.luaChannelReceive,
		"channel_close":  g.luaChannelClose,
		"channel_len":    g.luaChannelLen,
		"channel_cap":    g.luaChannelCap,
		"select":         g.luaSelect,
		"select_timeout": g.luaSelectTimeout,
		// Mutex operations
		"mutex":    g.luaMutex,
		"rwmutex":  g.luaRWMutex,
		// Semaphore operations
		"semaphore": g.luaSemaphore,
		// Atomic operations
		"atomic_int": g.luaAtomicInt,
		// Synchronization primitives
		"once": g.luaOnce,
		"cond":  g.luaCond,
		// Context operations
		"context": g.luaContext,
		// Pipeline operations
		"pipeline": g.luaPipeline,
		"fan_out":  g.luaFanOut,
		"fan_in":   g.luaFanIn,
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

// Cleanup closes all pools and channels
func (g *GoroutineModule) Cleanup() {
	g.globalCancel()

	g.mu.Lock()
	defer g.mu.Unlock()

	for name, pool := range g.pools {
		pool.close()
		delete(g.pools, name)
	}

	for name, ch := range g.channels {
		ch.closeMu.Lock()
		if !ch.closed {
			close(ch.ch)
			ch.closed = true
		}
		ch.closeMu.Unlock()
		delete(g.channels, name)
	}
}

// ============================================================================
// CHANNEL OPERATIONS
// ============================================================================

// createChannelMetatable creates a metatable with all channel methods
func (g *GoroutineModule) createChannelMetatable(L *lua.LState) *lua.LTable {
	mt := L.NewTable()

	// send method: ch:send(value)
	L.SetField(mt, "send", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)
		if ch.direction == "receive" {
			L.RaiseError("cannot send on receive-only channel")
			return 0
		}

		value := L.Get(2)

		ch.closeMu.Lock()
		closed := ch.closed
		ch.closeMu.Unlock()

		if closed {
			L.Push(lua.LFalse)
			L.Push(lua.LString("send on closed channel"))
			return 2
		}

		ch.ch <- value
		L.Push(lua.LTrue)
		return 1
	}))

	// receive method: local value, ok = ch:receive()
	L.SetField(mt, "receive", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)
		if ch.direction == "send" {
			L.RaiseError("cannot receive on send-only channel")
			return 0
		}

		value, ok := <-ch.ch
		if !ok {
			L.Push(lua.LNil)
			L.Push(lua.LFalse)
			return 2
		}

		L.Push(value)
		L.Push(lua.LTrue)
		return 2
	}))

	// try_send method: local ok = ch:try_send(value) -- non-blocking
	L.SetField(mt, "try_send", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)
		if ch.direction == "receive" {
			L.RaiseError("cannot send on receive-only channel")
			return 0
		}

		value := L.Get(2)

		ch.closeMu.Lock()
		closed := ch.closed
		ch.closeMu.Unlock()

		if closed {
			L.Push(lua.LFalse)
			L.Push(lua.LString("send on closed channel"))
			return 2
		}

		select {
		case ch.ch <- value:
			L.Push(lua.LTrue)
			return 1
		default:
			L.Push(lua.LFalse)
			return 1
		}
	}))

	// try_receive method: local value, ok = ch:try_receive() -- non-blocking
	L.SetField(mt, "try_receive", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)
		if ch.direction == "send" {
			L.RaiseError("cannot receive on send-only channel")
			return 0
		}

		select {
		case value, ok := <-ch.ch:
			if !ok {
				L.Push(lua.LNil)
				L.Push(lua.LFalse)
				return 2
			}
			L.Push(value)
			L.Push(lua.LTrue)
			return 2
		default:
			L.Push(lua.LNil)
			L.Push(lua.LFalse)
			return 2
		}
	}))

	// close method: ch:close()
	L.SetField(mt, "close", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)

		ch.closeMu.Lock()
		defer ch.closeMu.Unlock()

		if ch.closed {
			L.Push(lua.LFalse)
			L.Push(lua.LString("channel already closed"))
			return 2
		}

		close(ch.ch)
		ch.closed = true
		L.Push(lua.LTrue)
		return 1
	}))

	// len method: local len = ch:len()
	L.SetField(mt, "len", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)
		L.Push(lua.LNumber(len(ch.ch)))
		return 1
	}))

	// cap method: local cap = ch:cap()
	L.SetField(mt, "cap", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)
		L.Push(lua.LNumber(ch.capacity))
		return 1
	}))

	// is_closed method: local closed = ch:is_closed()
	L.SetField(mt, "is_closed", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)
		ch.closeMu.Lock()
		closed := ch.closed
		ch.closeMu.Unlock()
		L.Push(lua.LBool(closed))
		return 1
	}))

	// range method: ch:range(function(value) ... end) -- iterate until channel closes
	L.SetField(mt, "range", L.NewFunction(func(L *lua.LState) int {
		ch := L.CheckUserData(1).Value.(*luaChannel)
		if ch.direction == "send" {
			L.RaiseError("cannot range over send-only channel")
			return 0
		}

		fn := L.CheckFunction(2)

		// Iterate over channel until closed
		for {
			value, ok := <-ch.ch
			if !ok {
				// Channel closed, stop iteration
				break
			}

			// Call handler function with value
			if err := L.CallByParam(lua.P{
				Fn:      fn,
				NRet:    0,
				Protect: true,
			}, value); err != nil {
				L.RaiseError("error in range handler: %v", err)
				return 0
			}
		}

		return 0
	}))

	L.SetField(mt, "__index", mt)
	return mt
}

// luaChannelMake creates a new channel
// Usage:
//   local ch = goroutine.channel()           -- unbuffered
//   local ch = goroutine.channel(10)         -- buffered with capacity 10
//   local ch = goroutine.channel(10, "send") -- send-only channel
func (g *GoroutineModule) luaChannelMake(L *lua.LState) int {
	capacity := L.OptInt(1, 0)
	direction := L.OptString(2, "bidirectional") // "bidirectional", "send", "receive"

	if capacity < 0 {
		L.ArgError(1, "capacity must be non-negative")
		return 0
	}

	if direction != "bidirectional" && direction != "send" && direction != "receive" {
		L.ArgError(2, "direction must be 'bidirectional', 'send', or 'receive'")
		return 0
	}

	ch := &luaChannel{
		ch:        make(chan lua.LValue, capacity),
		capacity:  capacity,
		closed:    false,
		direction: direction,
	}

	// Create userdata
	ud := L.NewUserData()
	ud.Value = ch

	// Set metatable with all channel methods
	mt := g.createChannelMetatable(L)
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

// luaChannelSend sends a value to a channel (standalone function)
// Usage: goroutine.channel_send(ch, value)
func (g *GoroutineModule) luaChannelSend(L *lua.LState) int {
	ud := L.CheckUserData(1)
	ch := ud.Value.(*luaChannel)
	value := L.Get(2)

	if ch.direction == "receive" {
		L.RaiseError("cannot send on receive-only channel")
		return 0
	}

	ch.closeMu.Lock()
	closed := ch.closed
	ch.closeMu.Unlock()

	if closed {
		L.Push(lua.LFalse)
		L.Push(lua.LString("send on closed channel"))
		return 2
	}

	ch.ch <- value
	L.Push(lua.LTrue)
	return 1
}

// luaChannelReceive receives a value from a channel (standalone function)
// Usage: local value, ok = goroutine.channel_receive(ch)
func (g *GoroutineModule) luaChannelReceive(L *lua.LState) int {
	ud := L.CheckUserData(1)
	ch := ud.Value.(*luaChannel)

	if ch.direction == "send" {
		L.RaiseError("cannot receive on send-only channel")
		return 0
	}

	value, ok := <-ch.ch
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LFalse)
		return 2
	}

	L.Push(value)
	L.Push(lua.LTrue)
	return 2
}

// luaChannelClose closes a channel (standalone function)
// Usage: goroutine.channel_close(ch)
func (g *GoroutineModule) luaChannelClose(L *lua.LState) int {
	ud := L.CheckUserData(1)
	ch := ud.Value.(*luaChannel)

	ch.closeMu.Lock()
	defer ch.closeMu.Unlock()

	if ch.closed {
		L.Push(lua.LFalse)
		L.Push(lua.LString("channel already closed"))
		return 2
	}

	close(ch.ch)
	ch.closed = true
	L.Push(lua.LTrue)
	return 1
}

// luaChannelLen returns the number of elements in the channel buffer
// Usage: local len = goroutine.channel_len(ch)
func (g *GoroutineModule) luaChannelLen(L *lua.LState) int {
	ud := L.CheckUserData(1)
	ch := ud.Value.(*luaChannel)
	L.Push(lua.LNumber(len(ch.ch)))
	return 1
}

// luaChannelCap returns the capacity of the channel
// Usage: local cap = goroutine.channel_cap(ch)
func (g *GoroutineModule) luaChannelCap(L *lua.LState) int {
	ud := L.CheckUserData(1)
	ch := ud.Value.(*luaChannel)
	L.Push(lua.LNumber(ch.capacity))
	return 1
}

// luaSelect implements select statement for multiplexing channel operations
// Usage:
//   goroutine.select({
//     { channel = ch1, receive = true, handler = function(value) ... end },
//     { channel = ch2, send = value, handler = function() ... end },
//     { default = true, handler = function() ... end }
//   })
func (g *GoroutineModule) luaSelect(L *lua.LState) int {
	cases := L.CheckTable(1)

	var selectCases []selectCase
	var defaultCase *selectCase

	// Parse cases
	cases.ForEach(func(k, v lua.LValue) {
		if caseTable, ok := v.(*lua.LTable); ok {
			sc := selectCase{}

			// Check if it's a default case
			if defaultVal := caseTable.RawGetString("default"); defaultVal != lua.LNil {
				if defaultVal == lua.LTrue {
					sc.caseType = "default"
					sc.fn = caseTable.RawGetString("handler").(*lua.LFunction)
					defaultCase = &sc
					return
				}
			}

			// Get channel
			chUD := caseTable.RawGetString("channel").(*lua.LUserData)
			sc.channel = chUD.Value.(*luaChannel)

			// Determine operation type
			if receiveVal := caseTable.RawGetString("receive"); receiveVal == lua.LTrue {
				sc.caseType = "receive"
			} else if sendVal := caseTable.RawGetString("send"); sendVal != lua.LNil {
				sc.caseType = "send"
				sc.value = sendVal
			}

			sc.fn = caseTable.RawGetString("handler").(*lua.LFunction)
			selectCases = append(selectCases, sc)
		}
	})

	// Execute select
	if len(selectCases) == 0 && defaultCase != nil {
		// Only default case
		if err := L.CallByParam(lua.P{
			Fn:      defaultCase.fn,
			NRet:    0,
			Protect: true,
		}); err != nil {
			L.RaiseError("error in default handler: %v", err)
		}
		return 0
	}

	// Build reflection select cases
	selectCh := make(chan int, 1)

	for i, sc := range selectCases {
		idx := i
		case_ := sc

		go func() {
			if case_.caseType == "receive" {
				value, ok := <-case_.channel.ch
				if ok {
					// Execute handler with received value
					if err := L.CallByParam(lua.P{
						Fn:      case_.fn,
						NRet:    0,
						Protect: true,
					}, value); err != nil {
						fmt.Printf("Error in receive handler: %v\n", err)
					}
					select {
					case selectCh <- idx:
					default:
					}
				}
			} else if case_.caseType == "send" {
				case_.channel.ch <- case_.value
				// Execute handler
				if err := L.CallByParam(lua.P{
					Fn:      case_.fn,
					NRet:    0,
					Protect: true,
				}); err != nil {
					fmt.Printf("Error in send handler: %v\n", err)
				}
				select {
				case selectCh <- idx:
				default:
				}
			}
		}()
	}

	// Wait for first case to complete or use default
	select {
	case <-selectCh:
		// One of the cases completed
	default:
		if defaultCase != nil {
			if err := L.CallByParam(lua.P{
				Fn:      defaultCase.fn,
				NRet:    0,
				Protect: true,
			}); err != nil {
				L.RaiseError("error in default handler: %v", err)
			}
		}
	}

	return 0
}

// luaSelectTimeout implements select with timeout for multiplexing channel operations
// Usage:
//   local timedout, result = goroutine.select_timeout(timeout_ms, {
//     { channel = ch1, receive = true, handler = function(value) ... end },
//     { channel = ch2, send = value, handler = function() ... end },
//   })
// Returns: timedout (boolean), result (value from handler if any)
func (g *GoroutineModule) luaSelectTimeout(L *lua.LState) int {
	timeoutMs := L.CheckInt(1)
	cases := L.CheckTable(2)

	var selectCases []selectCase

	// Parse cases
	cases.ForEach(func(k, v lua.LValue) {
		if caseTable, ok := v.(*lua.LTable); ok {
			sc := selectCase{}

			// Get channel
			chUD := caseTable.RawGetString("channel").(*lua.LUserData)
			sc.channel = chUD.Value.(*luaChannel)

			// Determine operation type
			if receiveVal := caseTable.RawGetString("receive"); receiveVal == lua.LTrue {
				sc.caseType = "receive"
			} else if sendVal := caseTable.RawGetString("send"); sendVal != lua.LNil {
				sc.caseType = "send"
				sc.value = sendVal
			}

			sc.fn = caseTable.RawGetString("handler").(*lua.LFunction)
			selectCases = append(selectCases, sc)
		}
	})

	if len(selectCases) == 0 {
		L.Push(lua.LTrue) // Timed out immediately (no cases)
		L.Push(lua.LNil)
		return 2
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	// Result channel
	type selectResult struct {
		idx     int
		success bool
		err     error
	}
	resultCh := make(chan selectResult, 1)

	// Launch goroutines for each case
	for i, sc := range selectCases {
		idx := i
		case_ := sc

		go func() {
			if case_.caseType == "receive" {
				select {
				case value, ok := <-case_.channel.ch:
					if ok {
						// Execute handler with received value
						if err := L.CallByParam(lua.P{
							Fn:      case_.fn,
							NRet:    0,
							Protect: true,
						}, value); err != nil {
							select {
							case resultCh <- selectResult{idx: idx, success: false, err: err}:
							case <-ctx.Done():
							}
							return
						}
						select {
						case resultCh <- selectResult{idx: idx, success: true}:
						case <-ctx.Done():
						}
					}
				case <-ctx.Done():
					return
				}
			} else if case_.caseType == "send" {
				select {
				case case_.channel.ch <- case_.value:
					// Execute handler
					if err := L.CallByParam(lua.P{
						Fn:      case_.fn,
						NRet:    0,
						Protect: true,
					}); err != nil {
						select {
						case resultCh <- selectResult{idx: idx, success: false, err: err}:
						case <-ctx.Done():
						}
						return
					}
					select {
					case resultCh <- selectResult{idx: idx, success: true}:
					case <-ctx.Done():
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Wait for first result or timeout
	select {
	case result := <-resultCh:
		if result.err != nil {
			L.Push(lua.LFalse) // Not timed out
			L.Push(lua.LString(result.err.Error()))
			return 2
		}
		L.Push(lua.LFalse) // Not timed out
		L.Push(lua.LNumber(result.idx))
		return 2
	case <-ctx.Done():
		L.Push(lua.LTrue) // Timed out
		L.Push(lua.LNil)
		return 2
	}
}

// ============================================================================
// MUTEX OPERATIONS
// ============================================================================

// luaMutex creates a new mutex
// Usage:
//   local mu = goroutine.mutex()
//   mu:lock()
//   -- critical section
//   mu:unlock()
func (g *GoroutineModule) luaMutex(L *lua.LState) int {
	mu := &sync.Mutex{}

	ud := L.NewUserData()
	ud.Value = mu

	// Create metatable with methods
	mt := L.NewTable()

	// lock method: mu:lock()
	L.SetField(mt, "lock", L.NewFunction(func(L *lua.LState) int {
		mu := L.CheckUserData(1).Value.(*sync.Mutex)
		mu.Lock()
		return 0
	}))

	// unlock method: mu:unlock()
	L.SetField(mt, "unlock", L.NewFunction(func(L *lua.LState) int {
		mu := L.CheckUserData(1).Value.(*sync.Mutex)
		mu.Unlock()
		return 0
	}))

	// try_lock method: local ok = mu:try_lock() -- non-blocking
	L.SetField(mt, "try_lock", L.NewFunction(func(L *lua.LState) int {
		mu := L.CheckUserData(1).Value.(*sync.Mutex)
		locked := mu.TryLock()
		L.Push(lua.LBool(locked))
		return 1
	}))

	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

// luaRWMutex creates a new read-write mutex
// Usage:
//   local rwmu = goroutine.rwmutex()
//   rwmu:rlock()     -- read lock
//   rwmu:runlock()   -- read unlock
//   rwmu:lock()      -- write lock
//   rwmu:unlock()    -- write unlock
func (g *GoroutineModule) luaRWMutex(L *lua.LState) int {
	rwmu := &sync.RWMutex{}

	ud := L.NewUserData()
	ud.Value = rwmu

	// Create metatable with methods
	mt := L.NewTable()

	// lock method: rwmu:lock() -- write lock
	L.SetField(mt, "lock", L.NewFunction(func(L *lua.LState) int {
		rwmu := L.CheckUserData(1).Value.(*sync.RWMutex)
		rwmu.Lock()
		return 0
	}))

	// unlock method: rwmu:unlock() -- write unlock
	L.SetField(mt, "unlock", L.NewFunction(func(L *lua.LState) int {
		rwmu := L.CheckUserData(1).Value.(*sync.RWMutex)
		rwmu.Unlock()
		return 0
	}))

	// rlock method: rwmu:rlock() -- read lock
	L.SetField(mt, "rlock", L.NewFunction(func(L *lua.LState) int {
		rwmu := L.CheckUserData(1).Value.(*sync.RWMutex)
		rwmu.RLock()
		return 0
	}))

	// runlock method: rwmu:runlock() -- read unlock
	L.SetField(mt, "runlock", L.NewFunction(func(L *lua.LState) int {
		rwmu := L.CheckUserData(1).Value.(*sync.RWMutex)
		rwmu.RUnlock()
		return 0
	}))

	// try_lock method: local ok = rwmu:try_lock() -- non-blocking write lock
	L.SetField(mt, "try_lock", L.NewFunction(func(L *lua.LState) int {
		rwmu := L.CheckUserData(1).Value.(*sync.RWMutex)
		locked := rwmu.TryLock()
		L.Push(lua.LBool(locked))
		return 1
	}))

	// try_rlock method: local ok = rwmu:try_rlock() -- non-blocking read lock
	L.SetField(mt, "try_rlock", L.NewFunction(func(L *lua.LState) int {
		rwmu := L.CheckUserData(1).Value.(*sync.RWMutex)
		locked := rwmu.TryRLock()
		L.Push(lua.LBool(locked))
		return 1
	}))

	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

// ============================================================================
// SEMAPHORE OPERATIONS
// ============================================================================

// luaSemaphore creates a new semaphore
// Usage:
//   local sem = goroutine.semaphore(5)  -- Capacity of 5
//   sem:acquire()   -- Acquire a token (blocks if at capacity)
//   -- Use resource
//   sem:release()   -- Release token
func (g *GoroutineModule) luaSemaphore(L *lua.LState) int {
	capacity := L.CheckInt(1)

	if capacity <= 0 {
		L.ArgError(1, "capacity must be positive")
		return 0
	}

	sem := &luaSemaphore{
		ch:       make(chan struct{}, capacity),
		capacity: capacity,
	}

	// Fill semaphore with initial tokens
	for i := 0; i < capacity; i++ {
		sem.ch <- struct{}{}
	}

	// Create userdata
	ud := L.NewUserData()
	ud.Value = sem

	// Create metatable with methods
	mt := L.NewTable()

	// acquire method: sem:acquire() -- blocks if no tokens available
	L.SetField(mt, "acquire", L.NewFunction(func(L *lua.LState) int {
		sem := L.CheckUserData(1).Value.(*luaSemaphore)
		<-sem.ch // Wait for token
		return 0
	}))

	// release method: sem:release() -- returns a token
	L.SetField(mt, "release", L.NewFunction(func(L *lua.LState) int {
		sem := L.CheckUserData(1).Value.(*luaSemaphore)
		select {
		case sem.ch <- struct{}{}:
			// Token released
		default:
			L.RaiseError("semaphore: release without acquire")
		}
		return 0
	}))

	// try_acquire method: local ok = sem:try_acquire() -- non-blocking
	L.SetField(mt, "try_acquire", L.NewFunction(func(L *lua.LState) int {
		sem := L.CheckUserData(1).Value.(*luaSemaphore)
		select {
		case <-sem.ch:
			L.Push(lua.LTrue)
			return 1
		default:
			L.Push(lua.LFalse)
			return 1
		}
	}))

	// available method: local count = sem:available() -- how many tokens available
	L.SetField(mt, "available", L.NewFunction(func(L *lua.LState) int {
		sem := L.CheckUserData(1).Value.(*luaSemaphore)
		L.Push(lua.LNumber(len(sem.ch)))
		return 1
	}))

	// capacity method: local cap = sem:capacity() -- total capacity
	L.SetField(mt, "capacity", L.NewFunction(func(L *lua.LState) int {
		sem := L.CheckUserData(1).Value.(*luaSemaphore)
		L.Push(lua.LNumber(sem.capacity))
		return 1
	}))

	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

// ============================================================================
// ATOMIC OPERATIONS
// ============================================================================

// luaAtomicInt creates a new atomic integer
// Usage:
//   local counter = goroutine.atomic_int(0)  -- Initialize with 0
//   counter:add(1)       -- Atomic increment
//   counter:load()       -- Read value
//   counter:store(10)    -- Set value
//   counter:swap(20)     -- Swap and return old value
//   counter:compare_and_swap(10, 20)  -- CAS operation
func (g *GoroutineModule) luaAtomicInt(L *lua.LState) int {
	initialValue := L.OptInt64(1, 0)

	atomicInt := &luaAtomicInt64{
		value: initialValue,
	}

	// Create userdata
	ud := L.NewUserData()
	ud.Value = atomicInt

	// Create metatable with methods
	mt := L.NewTable()

	// add method: counter:add(delta) -- atomic add, returns new value
	L.SetField(mt, "add", L.NewFunction(func(L *lua.LState) int {
		ai := L.CheckUserData(1).Value.(*luaAtomicInt64)
		delta := L.CheckInt64(2)
		newValue := atomic.AddInt64(&ai.value, delta)
		L.Push(lua.LNumber(newValue))
		return 1
	}))

	// load method: local val = counter:load() -- atomic load
	L.SetField(mt, "load", L.NewFunction(func(L *lua.LState) int {
		ai := L.CheckUserData(1).Value.(*luaAtomicInt64)
		value := atomic.LoadInt64(&ai.value)
		L.Push(lua.LNumber(value))
		return 1
	}))

	// store method: counter:store(value) -- atomic store
	L.SetField(mt, "store", L.NewFunction(func(L *lua.LState) int {
		ai := L.CheckUserData(1).Value.(*luaAtomicInt64)
		value := L.CheckInt64(2)
		atomic.StoreInt64(&ai.value, value)
		return 0
	}))

	// swap method: local old = counter:swap(new) -- atomic swap, returns old value
	L.SetField(mt, "swap", L.NewFunction(func(L *lua.LState) int {
		ai := L.CheckUserData(1).Value.(*luaAtomicInt64)
		newValue := L.CheckInt64(2)
		oldValue := atomic.SwapInt64(&ai.value, newValue)
		L.Push(lua.LNumber(oldValue))
		return 1
	}))

	// compare_and_swap method: local swapped = counter:compare_and_swap(old, new)
	L.SetField(mt, "compare_and_swap", L.NewFunction(func(L *lua.LState) int {
		ai := L.CheckUserData(1).Value.(*luaAtomicInt64)
		oldValue := L.CheckInt64(2)
		newValue := L.CheckInt64(3)
		swapped := atomic.CompareAndSwapInt64(&ai.value, oldValue, newValue)
		L.Push(lua.LBool(swapped))
		return 1
	}))

	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

// ============================================================================
// SYNCHRONIZATION PRIMITIVES
// ============================================================================

// luaOnce creates a sync.Once for one-time initialization
// Usage:
//   local once = goroutine.once()
//   once:call(function()
//     log.info("This runs only once")
//   end)
func (g *GoroutineModule) luaOnce(L *lua.LState) int {
	onceObj := &luaOnce{
		once: sync.Once{},
	}

	// Create userdata
	ud := L.NewUserData()
	ud.Value = onceObj

	// Create metatable with methods
	mt := L.NewTable()

	// call method: once:call(function() ... end) -- execute function only once
	L.SetField(mt, "call", L.NewFunction(func(L *lua.LState) int {
		onceObj := L.CheckUserData(1).Value.(*luaOnce)
		fn := L.CheckFunction(2)

		onceObj.once.Do(func() {
			if err := L.CallByParam(lua.P{
				Fn:      fn,
				NRet:    0,
				Protect: true,
			}); err != nil {
				fmt.Printf("Error in once function: %v\n", err)
			}
		})

		return 0
	}))

	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

// ============================================================================
// CONDITION VARIABLES
// ============================================================================

// luaCond creates a condition variable for complex synchronization
// Usage:
//   local cond = goroutine.cond()
//   local mu = cond:get_mutex()  -- Get associated mutex
//   
//   -- Waiter goroutine
//   mu:lock()
//   while not condition do
//     cond:wait()  -- Releases lock, waits for signal, re-acquires lock
//   end
//   mu:unlock()
//   
//   -- Signaler goroutine
//   mu:lock()
//   condition = true
//   cond:signal()     -- Wake one waiter
//   -- or cond:broadcast()  -- Wake all waiters
//   mu:unlock()
func (g *GoroutineModule) luaCond(L *lua.LState) int {
	mu := &sync.Mutex{}
	condObj := &luaCond{
		cond: sync.NewCond(mu),
		mu:   mu,
	}

	// Create userdata
	ud := L.NewUserData()
	ud.Value = condObj

	// Create metatable with methods
	mt := L.NewTable()

	// wait method: cond:wait() -- releases lock, waits for signal, re-acquires lock
	L.SetField(mt, "wait", L.NewFunction(func(L *lua.LState) int {
		cond := L.CheckUserData(1).Value.(*luaCond)
		cond.cond.Wait()
		return 0
	}))

	// signal method: cond:signal() -- wakes one waiting goroutine
	L.SetField(mt, "signal", L.NewFunction(func(L *lua.LState) int {
		cond := L.CheckUserData(1).Value.(*luaCond)
		cond.cond.Signal()
		return 0
	}))

	// broadcast method: cond:broadcast() -- wakes all waiting goroutines
	L.SetField(mt, "broadcast", L.NewFunction(func(L *lua.LState) int {
		cond := L.CheckUserData(1).Value.(*luaCond)
		cond.cond.Broadcast()
		return 0
	}))

	// get_mutex method: local mu = cond:get_mutex() -- returns associated mutex
	L.SetField(mt, "get_mutex", L.NewFunction(func(L *lua.LState) int {
		cond := L.CheckUserData(1).Value.(*luaCond)
		
		// Create mutex userdata
		muUd := L.NewUserData()
		muUd.Value = cond.mu

		// Create metatable for mutex
		muMt := L.NewTable()

		// lock method
		L.SetField(muMt, "lock", L.NewFunction(func(L *lua.LState) int {
			mu := L.CheckUserData(1).Value.(*sync.Mutex)
			mu.Lock()
			return 0
		}))

		// unlock method
		L.SetField(muMt, "unlock", L.NewFunction(func(L *lua.LState) int {
			mu := L.CheckUserData(1).Value.(*sync.Mutex)
			mu.Unlock()
			return 0
		}))

		// try_lock method
		L.SetField(muMt, "try_lock", L.NewFunction(func(L *lua.LState) int {
			mu := L.CheckUserData(1).Value.(*sync.Mutex)
			locked := mu.TryLock()
			L.Push(lua.LBool(locked))
			return 1
		}))

		L.SetField(muMt, "__index", muMt)
		L.SetMetatable(muUd, muMt)

		L.Push(muUd)
		return 1
	}))

	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

// ============================================================================
// CONTEXT OPERATIONS
// ============================================================================

// luaContext creates a new context for cancellation and timeout management
// Usage:
//   -- Background context
//   local ctx = goroutine.context()
//
//   -- Context with cancellation
//   local ctx, cancel = goroutine.context()
//   cancel()  -- Cancel the context
//
//   -- Context with timeout (milliseconds)
//   local ctx = ctx:with_timeout(5000)
//
//   -- Context with deadline (milliseconds since epoch)
//   local ctx = ctx:with_deadline(os.time() * 1000 + 5000)
//
//   -- Check if context is cancelled
//   if ctx:is_cancelled() then ... end
//
//   -- Get error (returns "context canceled" or "context deadline exceeded")
//   local err = ctx:err()
func (g *GoroutineModule) luaContext(L *lua.LState) int {
	// Create background context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	ctxObj := &luaContext{
		ctx:    ctx,
		cancel: cancel,
	}

	// Create userdata
	ud := L.NewUserData()
	ud.Value = ctxObj

	// Create metatable with methods
	mt := L.NewTable()

	// with_cancel method: local childCtx, cancel = ctx:with_cancel()
	L.SetField(mt, "with_cancel", L.NewFunction(func(L *lua.LState) int {
		parentCtx := L.CheckUserData(1).Value.(*luaContext)
		childCtx, childCancel := context.WithCancel(parentCtx.ctx)

		child := &luaContext{
			ctx:    childCtx,
			cancel: childCancel,
		}

		childUd := L.NewUserData()
		childUd.Value = child
		L.SetMetatable(childUd, mt)

		// Create cancel function
		cancelFn := L.NewFunction(func(L *lua.LState) int {
			childCancel()
			return 0
		})

		L.Push(childUd)
		L.Push(cancelFn)
		return 2
	}))

	// with_timeout method: local ctx = ctx:with_timeout(milliseconds)
	L.SetField(mt, "with_timeout", L.NewFunction(func(L *lua.LState) int {
		parentCtx := L.CheckUserData(1).Value.(*luaContext)
		ms := L.CheckInt(2)

		childCtx, childCancel := context.WithTimeout(parentCtx.ctx, time.Duration(ms)*time.Millisecond)

		child := &luaContext{
			ctx:    childCtx,
			cancel: childCancel,
		}

		childUd := L.NewUserData()
		childUd.Value = child
		L.SetMetatable(childUd, mt)

		L.Push(childUd)
		return 1
	}))

	// with_deadline method: local ctx = ctx:with_deadline(deadline_ms)
	L.SetField(mt, "with_deadline", L.NewFunction(func(L *lua.LState) int {
		parentCtx := L.CheckUserData(1).Value.(*luaContext)
		deadlineMs := L.CheckInt64(2)

		// Convert milliseconds since epoch to time.Time
		deadline := time.Unix(0, deadlineMs*int64(time.Millisecond))

		childCtx, childCancel := context.WithDeadline(parentCtx.ctx, deadline)

		child := &luaContext{
			ctx:    childCtx,
			cancel: childCancel,
		}

		childUd := L.NewUserData()
		childUd.Value = child
		L.SetMetatable(childUd, mt)

		L.Push(childUd)
		return 1
	}))

	// is_cancelled method: local cancelled = ctx:is_cancelled()
	L.SetField(mt, "is_cancelled", L.NewFunction(func(L *lua.LState) int {
		ctxObj := L.CheckUserData(1).Value.(*luaContext)
		select {
		case <-ctxObj.ctx.Done():
			L.Push(lua.LTrue)
		default:
			L.Push(lua.LFalse)
		}
		return 1
	}))

	// err method: local err = ctx:err() -- returns error message if cancelled
	L.SetField(mt, "err", L.NewFunction(func(L *lua.LState) int {
		ctxObj := L.CheckUserData(1).Value.(*luaContext)
		err := ctxObj.ctx.Err()
		if err != nil {
			L.Push(lua.LString(err.Error()))
		} else {
			L.Push(lua.LNil)
		}
		return 1
	}))

	// cancel method: ctx:cancel()
	L.SetField(mt, "cancel", L.NewFunction(func(L *lua.LState) int {
		ctxObj := L.CheckUserData(1).Value.(*luaContext)
		if ctxObj.cancel != nil {
			ctxObj.cancel()
		}
		return 0
	}))

	// deadline method: local deadline_ms, ok = ctx:deadline()
	L.SetField(mt, "deadline", L.NewFunction(func(L *lua.LState) int {
		ctxObj := L.CheckUserData(1).Value.(*luaContext)
		deadline, ok := ctxObj.ctx.Deadline()
		if ok {
			// Convert to milliseconds since epoch
			ms := deadline.UnixNano() / int64(time.Millisecond)
			L.Push(lua.LNumber(ms))
			L.Push(lua.LTrue)
			return 2
		}
		L.Push(lua.LNil)
		L.Push(lua.LFalse)
		return 2
	}))

	L.SetField(mt, "__index", mt)
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

// ============================================================================
// PIPELINE OPERATIONS
// ============================================================================

// luaPipeline creates a processing pipeline from input to output through stages
// Usage:
//   local output = goroutine.pipeline(input_ch, {
//     { workers = 2, fn = function(value) return value * 2 end },
//     { workers = 3, fn = function(value) return value + 1 end }
//   })
func (g *GoroutineModule) luaPipeline(L *lua.LState) int {
	// Get input channel
	inputUD := L.CheckUserData(1)
	inputCh := inputUD.Value.(*luaChannel)

	// Get stages table
	stagesTable := L.CheckTable(2)

	currentCh := inputCh
	var wg sync.WaitGroup

	// Process each stage
	stagesTable.ForEach(func(k, v lua.LValue) {
		if stageTable, ok := v.(*lua.LTable); ok {
			// Get workers count (default 1)
			workers := 1
			if w := stageTable.RawGetString("workers"); w != lua.LNil {
				workers = int(w.(lua.LNumber))
			}

			// Get processing function
			fn := stageTable.RawGetString("fn").(*lua.LFunction)

			// Create output channel for this stage
			outputCh := &luaChannel{
				ch:        make(chan lua.LValue, 10),
				capacity:  10,
				closed:    false,
				direction: "bidirectional",
			}

			// Create a separate WaitGroup for this stage
			stageWg := &sync.WaitGroup{}

			// Launch workers for this stage
			for i := 0; i < workers; i++ {
				wg.Add(1)
				stageWg.Add(1)
				go func(input, output *luaChannel, processFn *lua.LFunction) {
					defer wg.Done()
					defer stageWg.Done()

					newL := lua.NewState()
					defer newL.Close()

					for {
						value, ok := <-input.ch
						if !ok {
							return
						}

						// Process value
						err := newL.CallByParam(lua.P{
							Fn:      processFn,
							NRet:    1,
							Protect: true,
						}, value)

						if err != nil {
							fmt.Printf("Pipeline stage error: %v\n", err)
							continue
						}

						result := newL.Get(-1)
						newL.Pop(1)

						// Send to next stage
						output.ch <- result
					}
				}(currentCh, outputCh, fn)
			}

			// Close this stage's output channel when all its workers are done
			go func(output *luaChannel, stageWg *sync.WaitGroup) {
				stageWg.Wait()
				output.closeMu.Lock()
				if !output.closed {
					close(output.ch)
					output.closed = true
				}
				output.closeMu.Unlock()
			}(outputCh, stageWg)

			// This stage's output becomes next stage's input
			currentCh = outputCh
		}
	})

	// Note: Each stage closes its own output channel when all its workers complete,
	// so the final output channel (currentCh) will be automatically closed by the
	// last stage's cleanup goroutine.

	// Return final output channel
	finalUD := L.NewUserData()
	finalUD.Value = currentCh

	// Set metatable with all channel methods
	finalMT := g.createChannelMetatable(L)
	L.SetMetatable(finalUD, finalMT)

	L.Push(finalUD)
	return 1
}

// luaFanOut distributes work from one channel to multiple output channels
// Usage:
//   local outputs = goroutine.fan_out(input_ch, 3) -- Creates 3 output channels
func (g *GoroutineModule) luaFanOut(L *lua.LState) int {
	inputUD := L.CheckUserData(1)
	inputCh := inputUD.Value.(*luaChannel)
	numOutputs := L.CheckInt(2)

	if numOutputs <= 0 {
		L.ArgError(2, "number of outputs must be positive")
		return 0
	}

	// Create output channels
	outputs := make([]*luaChannel, numOutputs)
	for i := 0; i < numOutputs; i++ {
		outputs[i] = &luaChannel{
			ch:        make(chan lua.LValue, 10),
			capacity:  10,
			closed:    false,
			direction: "bidirectional",
		}
	}

	// Fan out goroutine
	go func() {
		for value := range inputCh.ch {
			// Send to all outputs
			for _, outCh := range outputs {
				outCh.ch <- value
			}
		}

		// Close all outputs when input closes
		for _, outCh := range outputs {
			outCh.closeMu.Lock()
			if !outCh.closed {
				close(outCh.ch)
				outCh.closed = true
			}
			outCh.closeMu.Unlock()
		}
	}()

	// Return table of output channels
	outputTable := L.NewTable()
	for i, outCh := range outputs {
		outUD := L.NewUserData()
		outUD.Value = outCh

		// Set metatable with all channel methods
		outMT := g.createChannelMetatable(L)
		L.SetMetatable(outUD, outMT)

		outputTable.RawSetInt(i+1, outUD)
	}

	L.Push(outputTable)
	return 1
}

// luaFanIn merges multiple input channels into one output channel
// Usage:
//   local merged = goroutine.fan_in({ch1, ch2, ch3})
func (g *GoroutineModule) luaFanIn(L *lua.LState) int {
	inputsTable := L.CheckTable(1)

	// Collect input channels
	var inputs []*luaChannel
	inputsTable.ForEach(func(k, v lua.LValue) {
		if ud, ok := v.(*lua.LUserData); ok {
			if ch, ok := ud.Value.(*luaChannel); ok {
				inputs = append(inputs, ch)
			}
		}
	})

	if len(inputs) == 0 {
		L.ArgError(1, "at least one input channel required")
		return 0
	}

	// Create output channel
	output := &luaChannel{
		ch:        make(chan lua.LValue, 10),
		capacity:  10,
		closed:    false,
		direction: "bidirectional",
	}

	var wg sync.WaitGroup

	// Fan in goroutines (one per input)
	for _, inputCh := range inputs {
		wg.Add(1)
		go func(in *luaChannel) {
			defer wg.Done()
			for value := range in.ch {
				output.ch <- value
			}
		}(inputCh)
	}

	// Close output when all inputs are done
	go func() {
		wg.Wait()
		output.closeMu.Lock()
		if !output.closed {
			close(output.ch)
			output.closed = true
		}
		output.closeMu.Unlock()
	}()

	// Return output channel
	fanInUD := L.NewUserData()
	fanInUD.Value = output

	// Set metatable with all channel methods
	fanInMT := g.createChannelMetatable(L)
	L.SetMetatable(fanInUD, fanInMT)

	L.Push(fanInUD)
	return 1
}
