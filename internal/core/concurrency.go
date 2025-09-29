package core

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPool manages a pool of workers for concurrent task execution
type WorkerPool struct {
	workers     int
	tasks       chan func()
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	active      int64
	completed   int64
	failed      int64
	queued      int64
}

// NewWorkerPool creates a new worker pool with specified number of workers
func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	wp := &WorkerPool{
		workers: workers,
		tasks:   make(chan func(), workers*2), // Buffer to prevent blocking
		ctx:     ctx,
		cancel:  cancel,
	}
	
	// Start workers
	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
	
	return wp
}

// worker is the main worker goroutine
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	
	for {
		select {
		case task := <-wp.tasks:
			atomic.AddInt64(&wp.active, 1)
			atomic.AddInt64(&wp.queued, -1)
			
			func() {
				defer func() {
					atomic.AddInt64(&wp.active, -1)
					
					if recover() != nil {
						atomic.AddInt64(&wp.failed, 1)
					} else {
						atomic.AddInt64(&wp.completed, 1)
					}
				}()
				
				task()
			}()
			
		case <-wp.ctx.Done():
			return
		}
	}
}

// Submit submits a task to the worker pool
func (wp *WorkerPool) Submit(task func()) bool {
	select {
	case wp.tasks <- task:
		atomic.AddInt64(&wp.queued, 1)
		return true
	case <-wp.ctx.Done():
		return false
	default:
		// Channel is full, task rejected
		return false
	}
}

// SubmitWithTimeout submits a task with timeout
func (wp *WorkerPool) SubmitWithTimeout(task func(), timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(wp.ctx, timeout)
	defer cancel()
	
	select {
	case wp.tasks <- task:
		atomic.AddInt64(&wp.queued, 1)
		return true
	case <-ctx.Done():
		return false
	}
}

// Stats returns current worker pool statistics
func (wp *WorkerPool) Stats() WorkerPoolStats {
	return WorkerPoolStats{
		Workers:   wp.workers,
		Active:    atomic.LoadInt64(&wp.active),
		Queued:    atomic.LoadInt64(&wp.queued),
		Completed: atomic.LoadInt64(&wp.completed),
		Failed:    atomic.LoadInt64(&wp.failed),
	}
}

// Close gracefully shuts down the worker pool
func (wp *WorkerPool) Close() {
	wp.cancel()
	close(wp.tasks)
	wp.wg.Wait()
}

type WorkerPoolStats struct {
	Workers   int
	Active    int64
	Queued    int64
	Completed int64
	Failed    int64
}

// SafeMap provides a thread-safe map implementation
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewSafeMap creates a new thread-safe map
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]interface{}),
	}
}

// Set stores a key-value pair
func (sm *SafeMap) Set(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// Get retrieves a value by key
func (sm *SafeMap) Get(key string) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, exists := sm.data[key]
	return value, exists
}

// Delete removes a key-value pair
func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

// Keys returns all keys in the map
func (sm *SafeMap) Keys() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	keys := make([]string, 0, len(sm.data))
	for key := range sm.data {
		keys = append(keys, key)
	}
	return keys
}

// Len returns the number of elements in the map
func (sm *SafeMap) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.data)
}

// Clear removes all elements from the map
func (sm *SafeMap) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data = make(map[string]interface{})
}

// ForEach iterates over all key-value pairs
func (sm *SafeMap) ForEach(fn func(string, interface{})) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	for key, value := range sm.data {
		fn(key, value)
	}
}

// Semaphore provides a counting semaphore for resource limiting
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a new semaphore with specified capacity
func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

// Acquire acquires a permit from the semaphore
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// TryAcquire tries to acquire a permit without blocking
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireWithTimeout tries to acquire a permit with timeout
func (s *Semaphore) AcquireWithTimeout(timeout time.Duration) bool {
	select {
	case s.ch <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Release releases a permit back to the semaphore
func (s *Semaphore) Release() {
	<-s.ch
}

// Available returns the number of available permits
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}

// RWCounter provides a thread-safe counter with separate read/write operations
type RWCounter struct {
	mu    sync.RWMutex
	value int64
}

// NewRWCounter creates a new thread-safe counter
func NewRWCounter() *RWCounter {
	return &RWCounter{}
}

// Increment atomically increments the counter
func (c *RWCounter) Increment() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

// Decrement atomically decrements the counter
func (c *RWCounter) Decrement() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
	return c.value
}

// Add atomically adds a value to the counter
func (c *RWCounter) Add(delta int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
	return c.value
}

// Value returns the current counter value
func (c *RWCounter) Value() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Set sets the counter to a specific value
func (c *RWCounter) Set(value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = value
}

// Reset resets the counter to zero
func (c *RWCounter) Reset() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	old := c.value
	c.value = 0
	return old
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	rate     int64         // tokens per second
	capacity int64         // maximum tokens
	tokens   int64         // current tokens
	lastRefill time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, capacity int64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

// Allow checks if an operation is allowed under the rate limit
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	rl.refill()
	
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	
	return false
}

// AllowN checks if N operations are allowed
func (rl *RateLimiter) AllowN(n int64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	rl.refill()
	
	if rl.tokens >= n {
		rl.tokens -= n
		return true
	}
	
	return false
}

// refill adds tokens based on elapsed time
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	
	tokensToAdd := int64(elapsed.Seconds()) * rl.rate
	rl.tokens += tokensToAdd
	
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	
	rl.lastRefill = now
}

// Tokens returns the current number of available tokens
func (rl *RateLimiter) Tokens() int64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	rl.refill()
	return rl.tokens
}

// OncePool provides a pool of sync.Once instances for one-time initialization
type OncePool struct {
	pool sync.Pool
}

// NewOncePool creates a new pool of sync.Once instances
func NewOncePool() *OncePool {
	return &OncePool{
		pool: sync.Pool{
			New: func() interface{} {
				return &sync.Once{}
			},
		},
	}
}

// Get gets a sync.Once from the pool
func (op *OncePool) Get() *sync.Once {
	return op.pool.Get().(*sync.Once)
}

// Put returns a sync.Once to the pool (after resetting it)
func (op *OncePool) Put(once *sync.Once) {
	// sync.Once cannot be reused once Do() has been called
	// So we create a new one instead
	op.pool.Put(&sync.Once{})
}

// Barrier implements a reusable barrier for goroutine synchronization
type Barrier struct {
	mu       sync.Mutex
	cond     *sync.Cond
	total    int
	waiting  int
	generation int
}

// NewBarrier creates a new barrier for n goroutines
func NewBarrier(n int) *Barrier {
	b := &Barrier{
		total: n,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// Wait waits for all goroutines to reach the barrier
func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	generation := b.generation
	b.waiting++
	
	if b.waiting == b.total {
		// Last goroutine to reach the barrier
		b.waiting = 0
		b.generation++
		b.cond.Broadcast()
	} else {
		// Wait for all other goroutines
		for generation == b.generation {
			b.cond.Wait()
		}
	}
}