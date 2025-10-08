package pooling

import (
	"sync"

	"github.com/yuin/gopher-lua"
)

// LuaStatePool provides a pool of reusable Lua states
// Creating Lua states is expensive, so pooling reduces CPU and memory overhead
type LuaStatePool struct {
	pool sync.Pool
	size int
}

// NewLuaStatePool creates a new Lua state pool
func NewLuaStatePool(maxStates int) *LuaStatePool {
	return &LuaStatePool{
		size: maxStates,
		pool: sync.Pool{
			New: func() interface{} {
				return lua.NewState()
			},
		},
	}
}

// Get retrieves a Lua state from the pool
func (p *LuaStatePool) Get() *lua.LState {
	L := p.pool.Get().(*lua.LState)
	// Clear any leftover state
	L.SetTop(0)
	return L
}

// Put returns a Lua state to the pool
func (p *LuaStatePool) Put(L *lua.LState) {
	if L == nil {
		return
	}
	// Clear the state before returning to pool
	L.SetTop(0)
	p.pool.Put(L)
}

// Close closes all Lua states in the pool (call on shutdown)
func (p *LuaStatePool) Close() {
	// Note: sync.Pool doesn't provide a way to iterate over items
	// So we can't close all states. They'll be GC'd when pool is released.
	// This is fine since we're shutting down anyway.
}
