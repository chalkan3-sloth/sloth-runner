package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EventWorker monitors local events and sends them to master
type EventWorker struct {
	agentName     string
	masterAddr    string
	batchSize     int
	flushInterval time.Duration

	mu            sync.Mutex
	events        []*pb.EventData
	client        pb.AgentRegistryClient
	conn          *grpc.ClientConn

	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// EventWorkerConfig holds configuration for the event worker
type EventWorkerConfig struct {
	AgentName     string
	MasterAddr    string
	BatchSize     int           // Max events to buffer before sending
	FlushInterval time.Duration // Max time to wait before sending buffered events
}

// NewEventWorker creates a new event worker
func NewEventWorker(config EventWorkerConfig) *EventWorker {
	ctx, cancel := context.WithCancel(context.Background())

	// Set defaults
	if config.BatchSize == 0 {
		config.BatchSize = 50 // Default batch size
	}
	if config.FlushInterval == 0 {
		config.FlushInterval = 10 * time.Second // Default flush interval
	}

	return &EventWorker{
		agentName:     config.AgentName,
		masterAddr:    config.MasterAddr,
		batchSize:     config.BatchSize,
		flushInterval: config.FlushInterval,
		events:        make([]*pb.EventData, 0, config.BatchSize),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Start begins the event worker (connects to master and starts monitoring)
func (w *EventWorker) Start() error {
	// Connect to master
	conn, err := grpc.Dial(w.masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to master: %w", err)
	}

	w.conn = conn
	w.client = pb.NewAgentRegistryClient(conn)

	slog.Info("Event worker connected to master",
		"master_addr", w.masterAddr,
		"agent_name", w.agentName,
		"batch_size", w.batchSize,
		"flush_interval", w.flushInterval)

	// Start flush ticker
	w.wg.Add(1)
	go w.flushLoop()

	return nil
}

// Stop stops the event worker and flushes remaining events
func (w *EventWorker) Stop() error {
	w.cancel()
	w.wg.Wait()

	// Flush any remaining events
	if err := w.flush(); err != nil {
		slog.Error("Failed to flush remaining events", "error", err)
	}

	if w.conn != nil {
		w.conn.Close()
	}

	slog.Info("Event worker stopped", "agent_name", w.agentName)
	return nil
}

// SendEvent sends a single event to the master (adds to buffer)
func (w *EventWorker) SendEvent(eventType, stack, runID string, data map[string]interface{}) error {
	// Convert data to JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Create event data
	event := &pb.EventData{
		EventId:    uuid.New().String(),
		EventType:  eventType,
		AgentName:  w.agentName,
		Timestamp:  time.Now().Unix(),
		Stack:      stack,
		RunId:      runID,
		DataJson:   string(dataJSON),
		Severity:   "info",
	}

	// Add to buffer
	w.mu.Lock()
	w.events = append(w.events, event)
	shouldFlush := len(w.events) >= w.batchSize
	w.mu.Unlock()

	// Flush if batch is full
	if shouldFlush {
		return w.flush()
	}

	return nil
}

// SendEventWithSeverity sends a single event with custom severity
func (w *EventWorker) SendEventWithSeverity(eventType, stack, runID string, data map[string]interface{}, severity string) error {
	// Convert data to JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Create event data
	event := &pb.EventData{
		EventId:    uuid.New().String(),
		EventType:  eventType,
		AgentName:  w.agentName,
		Timestamp:  time.Now().Unix(),
		Stack:      stack,
		RunId:      runID,
		DataJson:   string(dataJSON),
		Severity:   severity,
	}

	// Add to buffer
	w.mu.Lock()
	w.events = append(w.events, event)
	shouldFlush := len(w.events) >= w.batchSize
	w.mu.Unlock()

	// Flush if batch is full
	if shouldFlush {
		return w.flush()
	}

	return nil
}

// flush sends all buffered events to master
func (w *EventWorker) flush() error {
	w.mu.Lock()
	if len(w.events) == 0 {
		w.mu.Unlock()
		return nil
	}

	// Take a copy and clear buffer
	events := make([]*pb.EventData, len(w.events))
	copy(events, w.events)
	w.events = w.events[:0]
	w.mu.Unlock()

	// Send batch to master
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := w.client.SendEventBatch(ctx, &pb.SendEventBatchRequest{
		Events:    events,
		BatchSize: int32(len(events)),
	})

	if err != nil {
		// Put events back in buffer on failure
		w.mu.Lock()
		w.events = append(events, w.events...)
		w.mu.Unlock()

		slog.Error("Failed to send event batch to master",
			"error", err,
			"batch_size", len(events))
		return fmt.Errorf("failed to send event batch: %w", err)
	}

	if !resp.Success {
		slog.Warn("Event batch partially processed",
			"processed", resp.EventsProcessed,
			"failed", len(resp.FailedEventIds),
			"message", resp.Message)
	} else {
		slog.Debug("Event batch sent successfully",
			"batch_size", len(events),
			"processed", resp.EventsProcessed)
	}

	return nil
}

// flushLoop periodically flushes buffered events
func (w *EventWorker) flushLoop() {
	defer w.wg.Done()

	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			if err := w.flush(); err != nil {
				slog.Debug("Periodic flush failed", "error", err)
			}
		}
	}
}

// MonitorSystemEvents monitors system-level events (CPU, memory, disk, etc.)
func (w *EventWorker) MonitorSystemEvents(interval time.Duration) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-w.ctx.Done():
				return
			case <-ticker.C:
				// Collect system info
				sysInfo, err := CollectSystemInfo()
				if err != nil {
					slog.Debug("Failed to collect system info", "error", err)
					continue
				}

				// Calculate average load for CPU usage estimate
				var avgLoad float64
				if len(sysInfo.LoadAverage) > 0 {
					avgLoad = sysInfo.LoadAverage[0]
				}
				cpuUsagePercent := (avgLoad / float64(sysInfo.CPUs)) * 100

				// Check for threshold violations and send events
				if cpuUsagePercent > 80.0 {
					w.SendEventWithSeverity("system.cpu_high", "", "", map[string]interface{}{
						"cpu_load_1min": avgLoad,
						"cpu_cores":     sysInfo.CPUs,
						"cpu_percent":   cpuUsagePercent,
					}, "warning")
				}

				if sysInfo.Memory != nil && sysInfo.Memory.UsedPercent > 85.0 {
					w.SendEventWithSeverity("system.memory_high", "", "", map[string]interface{}{
						"memory_used_percent": sysInfo.Memory.UsedPercent,
						"memory_used_bytes":   sysInfo.Memory.Used,
						"memory_total_bytes":  sysInfo.Memory.Total,
					}, "warning")
				}

				// Check first disk partition for space
				if len(sysInfo.Disk) > 0 && sysInfo.Disk[0].UsedPercent > 90.0 {
					w.SendEventWithSeverity("system.disk_full", "", "", map[string]interface{}{
						"disk_used_percent": sysInfo.Disk[0].UsedPercent,
						"disk_used_bytes":   sysInfo.Disk[0].Used,
						"disk_total_bytes":  sysInfo.Disk[0].Total,
						"mountpoint":        sysInfo.Disk[0].Mountpoint,
					}, "critical")
				}

				// Send periodic health check
				memUsedMB := uint64(0)
				if sysInfo.Memory != nil {
					memUsedMB = sysInfo.Memory.Used / (1024 * 1024)
				}

				diskUsedGB := uint64(0)
				if len(sysInfo.Disk) > 0 {
					diskUsedGB = sysInfo.Disk[0].Used / (1024 * 1024 * 1024)
				}

				w.SendEvent("agent.health_check", "", "", map[string]interface{}{
					"cpu_load":         avgLoad,
					"cpu_cores":        sysInfo.CPUs,
					"memory_used_mb":   memUsedMB,
					"disk_used_gb":     diskUsedGB,
					"uptime_seconds":   sysInfo.Uptime,
					"hostname":         sysInfo.Hostname,
				})
			}
		}
	}()

	slog.Info("System event monitoring started",
		"interval", interval,
		"agent", w.agentName)
}
