package metrics

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/webui/services"
)

// Collector periodically collects metrics from agents with optimizations
type Collector struct {
	metricsDB     *MetricsDB
	agentClient   *services.AgentClient
	interval      time.Duration
	retentionDays int
	mu            sync.RWMutex
	running       bool
	stopCh        chan struct{}
	wg            sync.WaitGroup

	// Optimizations
	batchSize     int
	timeout       time.Duration
	lastAgents    []AgentInfo
	lastCheck     time.Time
	agentCache    sync.RWMutex
}

// CollectorConfig holds configuration for the metrics collector
type CollectorConfig struct {
	MetricsDB     *MetricsDB
	AgentClient   *services.AgentClient
	Interval      time.Duration // How often to collect metrics
	RetentionDays int           // How long to keep metrics
	BatchSize     int           // How many agents to collect in parallel
	Timeout       time.Duration // Timeout per agent request
}

// NewCollector creates a new metrics collector with optimizations
func NewCollector(cfg CollectorConfig) *Collector {
	if cfg.Interval == 0 {
		cfg.Interval = 60 * time.Second // Optimized: collect every 60 seconds (was 30s)
	}
	if cfg.RetentionDays == 0 {
		cfg.RetentionDays = 7 // Default: keep 7 days of metrics
	}
	if cfg.BatchSize == 0 {
		cfg.BatchSize = 5 // Process 5 agents at a time
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 2 * time.Second // 2s timeout per agent
	}

	return &Collector{
		metricsDB:     cfg.MetricsDB,
		agentClient:   cfg.AgentClient,
		interval:      cfg.Interval,
		retentionDays: cfg.RetentionDays,
		batchSize:     cfg.BatchSize,
		timeout:       cfg.Timeout,
		stopCh:        make(chan struct{}),
	}
}

// Start starts the metrics collector
func (c *Collector) Start(ctx context.Context, getAgents func() []AgentInfo) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = true
	c.mu.Unlock()

	slog.Info("Starting metrics collector", "interval", c.interval, "retention_days", c.retentionDays)

	c.wg.Add(2)

	// Goroutine for collecting metrics
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		// Collect immediately on start
		c.collectAllMetrics(ctx, getAgents())

		for {
			select {
			case <-ticker.C:
				c.collectAllMetrics(ctx, getAgents())
			case <-c.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	// Goroutine for cleanup old metrics
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(1 * time.Hour) // Cleanup every hour
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.cleanupOldMetrics(ctx)
			case <-c.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// AgentInfo contains agent information for collection
type AgentInfo struct {
	Name    string
	Address string
}

// collectAllMetrics collects metrics from all agents (optimized with batching)
func (c *Collector) collectAllMetrics(ctx context.Context, agents []AgentInfo) {
	if len(agents) == 0 {
		slog.Debug("No agents found for metrics collection")
		return
	}

	// Cache agent list for 5 minutes to avoid redundant checks
	c.agentCache.RLock()
	if time.Since(c.lastCheck) < 5*time.Minute && len(c.lastAgents) == len(agents) {
		// Use cached agent list
		agents = c.lastAgents
	}
	c.agentCache.RUnlock()

	slog.Info("Collecting metrics from agents", "count", len(agents), "batch_size", c.batchSize)

	// Batch processing: collect from N agents at a time
	for i := 0; i < len(agents); i += c.batchSize {
		end := i + c.batchSize
		if end > len(agents) {
			end = len(agents)
		}

		batch := agents[i:end]
		var wg sync.WaitGroup

		for _, agent := range batch {
			wg.Add(1)
			go func(a AgentInfo) {
				defer wg.Done()

				// Create timeout context
				timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
				defer cancel()

				c.collectAgentMetrics(timeoutCtx, a)
			}(agent)
		}

		wg.Wait()

		// Small delay between batches to avoid overwhelming the system
		if end < len(agents) {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Update cache
	c.agentCache.Lock()
	c.lastAgents = agents
	c.lastCheck = time.Now()
	c.agentCache.Unlock()
}

// collectAgentMetrics collects metrics from a single agent (optimized)
func (c *Collector) collectAgentMetrics(ctx context.Context, agent AgentInfo) {
	// Get resource usage from agent
	// Note: Connection pooling happens at gRPC layer automatically
	resp, err := c.agentClient.GetResourceUsage(ctx, agent.Address)
	if err != nil {
		slog.Warn("Failed to collect metrics from agent", "agent", agent.Name, "address", agent.Address, "error", err)
		return
	}

	slog.Debug("Got metrics from agent", "agent", agent.Name, "cpu", resp.CpuPercent, "memory", resp.MemoryPercent)

	// Create metric point
	metric := MetricPoint{
		Timestamp:       time.Now().Unix(),
		CPUPercent:      resp.CpuPercent,
		MemoryPercent:   resp.MemoryPercent,
		MemoryUsedBytes: resp.MemoryUsedBytes,
		DiskPercent:     resp.DiskPercent,
		LoadAvg1Min:     resp.LoadAvg_1Min,
		LoadAvg5Min:     resp.LoadAvg_5Min,
		LoadAvg15Min:    resp.LoadAvg_15Min,
		ProcessCount:    int(resp.ProcessCount),
	}

	// Store in database (batched)
	if err := c.metricsDB.StoreMetric(ctx, agent.Name, metric); err != nil {
		slog.Error("Failed to store metric", "agent", agent.Name, "error", err)
		return
	}

	slog.Debug("✅ Stored metric in database", "agent", agent.Name,
		"cpu", resp.CpuPercent, "memory", resp.MemoryPercent)
}

// cleanupOldMetrics removes metrics older than retention period
func (c *Collector) cleanupOldMetrics(ctx context.Context) {
	retention := time.Duration(c.retentionDays) * 24 * time.Hour
	if err := c.metricsDB.CleanupOldMetrics(ctx, retention); err != nil {
		slog.Error("Failed to cleanup old metrics", "error", err)
		return
	}
	slog.Debug("Cleaned up old metrics", "retention", retention)
}

// Stop stops the metrics collector
func (c *Collector) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return
	}

	slog.Info("Stopping metrics collector")
	c.running = false
	close(c.stopCh)
	c.wg.Wait()
}

// IsRunning returns true if collector is running
func (c *Collector) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}
