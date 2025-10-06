package core

import (
	"log/slog"
	"runtime"
	"sync"
	"time"
)

// GlobalCore holds all core components for the Sloth Runner
type GlobalCore struct {
	// Performance monitoring
	Performance *PerformanceMetrics
	
	// Resource management
	ResourcePool *ResourcePool
	MemoryManager *MemoryManager
	FileCache    *FileCache
	
	// Error handling
	ErrorRecovery *ErrorRecovery
	ErrorCollector *ErrorCollector
	TimeoutManager *TimeoutManager
	
	// Concurrency primitives
	WorkerPool *WorkerPool
	GlobalMap  *SafeMap
	
	// Circuit breakers for external dependencies
	CircuitBreakers map[string]*CircuitBreaker
	
	// Configuration
	Config *CoreConfig
	
	// Logger
	Logger *slog.Logger
	
	// Lifecycle
	mu       sync.RWMutex
	started  bool
	shutdown chan struct{}
}

// CoreConfig defines configuration for the core system
type CoreConfig struct {
	// Performance settings
	MaxWorkers       int
	MaxMemoryMB      int64
	CacheEnabled     bool
	CacheSizeMB      int64
	CacheCompression bool
	
	// Error handling settings
	MaxErrors        int
	PanicRecovery    bool
	TimeoutDefault   time.Duration
	
	// Circuit breaker settings
	CircuitBreakerMaxFailures int64
	CircuitBreakerResetTime   time.Duration
	
	// Resource limits
	MaxGoroutines    int
	GCInterval       time.Duration
	MetricsInterval  time.Duration
}

// DefaultCoreConfig returns default configuration
func DefaultCoreConfig() *CoreConfig {
	return &CoreConfig{
		MaxWorkers:                runtime.NumCPU() * 2,
		MaxMemoryMB:               512,
		CacheEnabled:              true,
		CacheSizeMB:               64,
		CacheCompression:          true,
		MaxErrors:                 100,
		PanicRecovery:             true,
		TimeoutDefault:            10 * time.Minute,
		CircuitBreakerMaxFailures: 5,
		CircuitBreakerResetTime:   60 * time.Second,
		MaxGoroutines:             1000,
		GCInterval:                5 * time.Minute,
		MetricsInterval:           30 * time.Second,
	}
}

// NewGlobalCore creates and initializes the global core system
func NewGlobalCore(config *CoreConfig, logger *slog.Logger) (*GlobalCore, error) {
	if config == nil {
		config = DefaultCoreConfig()
	}
	
	gc := &GlobalCore{
		Config:          config,
		Logger:          logger,
		CircuitBreakers: make(map[string]*CircuitBreaker),
		shutdown:        make(chan struct{}),
	}
	
	// Initialize components
	if err := gc.initializeComponents(); err != nil {
		return nil, err
	}
	
	return gc, nil
}

// initializeComponents initializes all core components
func (gc *GlobalCore) initializeComponents() error {
	// Performance metrics
	gc.Performance = NewPerformanceMetrics()
	
	// Resource pool
	gc.ResourcePool = NewResourcePool(gc.Config.MaxWorkers * 2)
	
	// Memory manager
	gc.MemoryManager = NewMemoryManager(gc.Config.MaxMemoryMB)
	
	// File cache
	if gc.Config.CacheEnabled {
		var err error
		gc.FileCache, err = NewFileCache(
			".sloth-cache",
			gc.Config.CacheSizeMB,
			gc.Config.CacheCompression,
		)
		if err != nil {
			gc.Logger.Warn("Failed to initialize file cache", "error", err)
		}
	}
	
	// Error handling
	gc.ErrorRecovery = NewErrorRecovery(gc.Logger)
	gc.ErrorCollector = NewErrorCollector(gc.Config.MaxErrors)
	gc.TimeoutManager = NewTimeoutManager(gc.Logger)
	
	// Concurrency
	gc.WorkerPool = NewWorkerPool(gc.Config.MaxWorkers)
	gc.GlobalMap = NewSafeMap()
	
	// Create common circuit breakers
	gc.createCircuitBreakers()
	
	return nil
}

// createCircuitBreakers creates circuit breakers for common external dependencies
func (gc *GlobalCore) createCircuitBreakers() {
	commonServices := []string{
		"http_external",
		"database",
		"filesystem", 
		"docker_daemon",
		"kubernetes_api",
		"cloud_provider",
	}
	
	for _, service := range commonServices {
		gc.CircuitBreakers[service] = NewCircuitBreaker(
			service,
			gc.Config.CircuitBreakerMaxFailures,
			gc.Config.CircuitBreakerResetTime,
		)
	}
}

// Start starts the core system
func (gc *GlobalCore) Start() error {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	
	if gc.started {
		return nil
	}
	
	gc.Logger.Info("Starting Sloth Runner core system")
	
	// Start background monitoring
	go gc.monitoringLoop()
	go gc.housekeepingLoop()
	
	gc.started = true
	gc.Logger.Info("Sloth Runner core system started successfully")
	
	return nil
}

// Stop gracefully stops the core system
func (gc *GlobalCore) Stop() error {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	
	if !gc.started {
		return nil
	}
	
	gc.Logger.Info("Stopping Sloth Runner core system")
	
	// Signal shutdown
	close(gc.shutdown)
	
	// Stop components
	if gc.WorkerPool != nil {
		gc.WorkerPool.Close()
	}
	
	if gc.TimeoutManager != nil {
		gc.TimeoutManager.CancelAll()
	}
	
	gc.started = false
	gc.Logger.Info("Sloth Runner core system stopped")
	
	return nil
}

// monitoringLoop runs background monitoring
func (gc *GlobalCore) monitoringLoop() {
	ticker := time.NewTicker(gc.Config.MetricsInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			gc.collectMetrics()
		case <-gc.shutdown:
			return
		}
	}
}

// housekeepingLoop runs periodic maintenance tasks
func (gc *GlobalCore) housekeepingLoop() {
	ticker := time.NewTicker(gc.Config.GCInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			gc.performHousekeeping()
		case <-gc.shutdown:
			return
		}
	}
}

// collectMetrics collects and logs system metrics
func (gc *GlobalCore) collectMetrics() {
	// Get performance metrics
	perfMetrics := gc.Performance.GetSnapshot()
	
	// Get worker pool stats
	var poolStats WorkerPoolStats
	if gc.WorkerPool != nil {
		poolStats = gc.WorkerPool.Stats()
	}
	
	// Get memory stats
	memStats := gc.MemoryManager.GetStats()
	
	// Get cache stats
	var cacheStats CacheStats
	if gc.FileCache != nil {
		cacheStats = gc.FileCache.Stats()
	}
	
	// Log metrics
	gc.Logger.Debug("Core system metrics",
		"goroutines", perfMetrics.GoroutineCount,
		"memory_alloc", perfMetrics.MemoryUsage.Alloc,
		"memory_peak", perfMetrics.PeakMemoryUsage,
		"tasks_executed", perfMetrics.TaskExecutions,
		"worker_active", poolStats.Active,
		"worker_queued", poolStats.Queued,
		"worker_completed", poolStats.Completed,
		"worker_failed", poolStats.Failed,
		"mem_tracked", memStats.TrackedAllocations,
		"mem_usage_pct", memStats.UsagePercent,
		"cache_entries", cacheStats.Entries,
		"cache_usage_pct", cacheStats.UsageRatio*100,
	)
	
	// Check for resource warnings
	gc.checkResourceWarnings(perfMetrics, memStats, poolStats)
}

// checkResourceWarnings checks for resource usage warnings
func (gc *GlobalCore) checkResourceWarnings(perf PerformanceMetrics, mem MemoryStats, pool WorkerPoolStats) {
	// Check goroutine count
	if perf.GoroutineCount > gc.Config.MaxGoroutines {
		gc.Logger.Warn("High goroutine count detected",
			"current", perf.GoroutineCount,
			"max", gc.Config.MaxGoroutines)
	}
	
	// Check memory usage
	if mem.UsagePercent > 80.0 {
		gc.Logger.Warn("High memory usage detected",
			"usage_percent", mem.UsagePercent,
			"current_mb", mem.CurrentSize/1024/1024,
			"max_mb", mem.MaxSize/1024/1024)
	}
	
	// Check worker pool queue
	if pool.Queued > int64(gc.Config.MaxWorkers*2) {
		gc.Logger.Warn("High worker queue detected",
			"queued", pool.Queued,
			"workers", pool.Workers)
	}
	
	// Check failed tasks
	if pool.Failed > 0 && pool.Completed > 0 {
		failureRate := float64(pool.Failed) / float64(pool.Completed+pool.Failed) * 100
		if failureRate > 10.0 {
			gc.Logger.Warn("High task failure rate detected",
				"failure_rate_percent", failureRate,
				"failed", pool.Failed,
				"completed", pool.Completed)
		}
	}
}

// performHousekeeping performs periodic cleanup and maintenance
func (gc *GlobalCore) performHousekeeping() {
	gc.Logger.Debug("Performing housekeeping tasks")
	
	// Force garbage collection if memory usage is high
	memStats := gc.MemoryManager.GetStats()
	if memStats.UsagePercent > 70.0 {
		runtime.GC()
		gc.Logger.Debug("Forced garbage collection due to high memory usage")
	}
	
	// Clear old errors
	if gc.ErrorCollector.HasErrors() {
		errors := gc.ErrorCollector.GetErrors()
		var criticalCount, highCount int
		
		for _, err := range errors {
			switch err.Severity {
			case SeverityCritical:
				criticalCount++
			case SeverityHigh:
				highCount++
			}
		}
		
		if criticalCount > 0 || highCount > 10 {
			gc.Logger.Warn("High error count detected",
				"critical", criticalCount,
				"high", highCount,
				"total", len(errors))
		}
		
		// Keep only recent errors
		if len(errors) > gc.Config.MaxErrors/2 {
			gc.ErrorCollector.Clear()
		}
	}
	
	// Log circuit breaker status
	for name, cb := range gc.CircuitBreakers {
		stats := cb.GetStats()
		if stats.State != "closed" || stats.Failures > 0 {
			gc.Logger.Debug("Circuit breaker status",
				"name", name,
				"state", stats.State,
				"failures", stats.Failures,
				"successes", stats.Successes)
		}
	}
}

// GetCircuitBreaker returns a circuit breaker by name
func (gc *GlobalCore) GetCircuitBreaker(name string) *CircuitBreaker {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	
	if cb, exists := gc.CircuitBreakers[name]; exists {
		return cb
	}
	
	// Create on-demand circuit breaker
	cb := NewCircuitBreaker(name, gc.Config.CircuitBreakerMaxFailures, gc.Config.CircuitBreakerResetTime)
	gc.CircuitBreakers[name] = cb
	
	return cb
}

// ExecuteWithRecovery executes a function with panic recovery and error collection
func (gc *GlobalCore) ExecuteWithRecovery(fn func() error, context string) error {
	start := time.Now()
	
	recovered, err := gc.ErrorRecovery.SafeExecute(fn)
	
	duration := time.Since(start)
	gc.Performance.RecordTaskExecution(duration)
	
	if err != nil {
		// Add context to error
		if se, ok := err.(*SlothError); ok {
			se.WithContext(context)
		}
		
		gc.ErrorCollector.Collect(err)
		
		if recovered {
			gc.Logger.Error("Panic recovered during execution",
				"context", context,
				"duration", duration,
				"error", err)
		} else {
			gc.Logger.Error("Error during execution",
				"context", context,
				"duration", duration,
				"error", err)
		}
	}
	
	return err
}

// SubmitTask submits a task to the worker pool
func (gc *GlobalCore) SubmitTask(task func(), context string) bool {
	wrappedTask := func() {
		gc.ExecuteWithRecovery(func() error {
			task()
			return nil
		}, context)
	}
	
	return gc.WorkerPool.Submit(wrappedTask)
}

// GetStats returns comprehensive core system statistics
func (gc *GlobalCore) GetStats() CoreStats {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	
	stats := CoreStats{
		Started: gc.started,
		Uptime:  time.Since(time.Now()), // This would be calculated properly
	}
	
	if gc.Performance != nil {
		stats.Performance = gc.Performance.GetSnapshot()
	}
	
	if gc.WorkerPool != nil {
		stats.WorkerPool = gc.WorkerPool.Stats()
	}
	
	if gc.MemoryManager != nil {
		stats.Memory = gc.MemoryManager.GetStats()
	}
	
	if gc.FileCache != nil {
		stats.Cache = gc.FileCache.Stats()
	}
	
	if gc.ErrorCollector != nil {
		stats.ErrorCount = len(gc.ErrorCollector.GetErrors())
	}
	
	// Circuit breaker stats
	stats.CircuitBreakers = make(map[string]CircuitStats)
	for name, cb := range gc.CircuitBreakers {
		stats.CircuitBreakers[name] = cb.GetStats()
	}
	
	return stats
}

// CoreStats represents comprehensive core system statistics
type CoreStats struct {
	Started         bool
	Uptime          time.Duration
	Performance     PerformanceMetrics
	WorkerPool      WorkerPoolStats
	Memory          MemoryStats
	Cache           CacheStats
	ErrorCount      int
	CircuitBreakers map[string]CircuitStats
}

// Global instance (singleton pattern)
var globalCoreInstance *GlobalCore
var globalCoreOnce sync.Once

// GetGlobalCore returns the global core instance
func GetGlobalCore() *GlobalCore {
	return globalCoreInstance
}

// InitializeGlobalCore initializes the global core instance
func InitializeGlobalCore(config *CoreConfig, logger *slog.Logger) error {
	var initErr error
	
	globalCoreOnce.Do(func() {
		globalCoreInstance, initErr = NewGlobalCore(config, logger)
		if initErr == nil {
			initErr = globalCoreInstance.Start()
		}
	})
	
	return initErr
}

// ShutdownGlobalCore shuts down the global core instance
func ShutdownGlobalCore() error {
	if globalCoreInstance != nil {
		return globalCoreInstance.Stop()
	}
	return nil
}