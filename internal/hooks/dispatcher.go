package hooks

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Dispatcher manages event dispatching to hooks
type Dispatcher struct {
	repo         *Repository
	executor     *Executor
	mu           sync.RWMutex
	enabled      bool
	stopChan     chan struct{}
	workerWg     sync.WaitGroup
	processing   bool
	eventChannel chan *Event
	maxWorkers   int

	// Execution context for events
	currentStack  string
	currentAgent  string
	currentRunID  string
}

// NewDispatcher creates a new event dispatcher
func NewDispatcher(repo *Repository) *Dispatcher {
	d := &Dispatcher{
		repo:         repo,
		executor:     NewExecutor(repo),
		enabled:      true,
		stopChan:     make(chan struct{}),
		eventChannel: make(chan *Event, 1000), // Buffer para 1000 eventos
		maxWorkers:   100,                      // Máximo de 100 goroutines simultâneas
	}

	// Start event processor workers in background
	d.StartEventProcessor()

	return d
}

// Dispatch dispatches an event by sending it to the event channel
func (d *Dispatcher) Dispatch(event *Event) error {
	d.mu.RLock()
	if !d.enabled {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	// Enqueue event in database for persistence
	if err := d.repo.EventQueue.EnqueueEvent(event); err != nil {
		return fmt.Errorf("failed to enqueue event: %w", err)
	}

	slog.Debug("event enqueued",
		"event_id", event.ID,
		"event_type", event.Type)

	// Send event to channel for immediate processing
	select {
	case d.eventChannel <- event:
		slog.Debug("event sent to processing channel", "event_id", event.ID)
	default:
		// Channel buffer is full, log warning
		slog.Warn("event channel buffer full, event will be processed from database",
			"event_id", event.ID,
			"event_type", event.Type)
	}

	return nil
}

// executeHook executes a single hook
func (d *Dispatcher) executeHook(hook *Hook, event *Event) {
	slog.Debug("executing hook", "hook_name", hook.Name, "event_type", event.Type)

	result, err := d.executor.Execute(hook, event)
	if err != nil {
		slog.Error("hook execution failed",
			"hook_name", hook.Name,
			"error", err)
	} else {
		if result.Success {
			slog.Info("hook executed successfully",
				"hook_name", hook.Name,
				"duration", result.Duration)
		} else {
			slog.Warn("hook execution returned failure",
				"hook_name", hook.Name,
				"error", result.Error)
		}
	}

	// Record execution
	if result != nil {
		if err := d.repo.RecordExecution(result); err != nil {
			slog.Error("failed to record hook execution",
				"hook_name", hook.Name,
				"error", err)
		}
	}
}

// DispatchAgentRegistered dispatches an agent.registered event
func (d *Dispatcher) DispatchAgentRegistered(agent *AgentEvent) error {
	event := &Event{
		Type:      EventAgentRegistered,
		Timestamp: getCurrentTime(),
		Data: map[string]interface{}{
			"agent": map[string]interface{}{
				"name":        agent.Name,
				"address":     agent.Address,
				"tags":        agent.Tags,
				"version":     agent.Version,
				"system_info": agent.SystemInfo,
			},
		},
	}

	return d.Dispatch(event)
}

// DispatchAgentDisconnected dispatches an agent.disconnected event
func (d *Dispatcher) DispatchAgentDisconnected(agent *AgentEvent) error {
	event := &Event{
		Type:      EventAgentDisconnected,
		Timestamp: getCurrentTime(),
		Data: map[string]interface{}{
			"agent": map[string]interface{}{
				"name":    agent.Name,
				"address": agent.Address,
				"tags":    agent.Tags,
				"version": agent.Version,
			},
		},
	}

	return d.Dispatch(event)
}

// DispatchAgentHeartbeatFailed dispatches an agent.heartbeat_failed event
func (d *Dispatcher) DispatchAgentHeartbeatFailed(agent *AgentEvent) error {
	event := &Event{
		Type:      EventAgentHeartbeatFailed,
		Timestamp: getCurrentTime(),
		Data: map[string]interface{}{
			"agent": map[string]interface{}{
				"name":    agent.Name,
				"address": agent.Address,
			},
		},
	}

	return d.Dispatch(event)
}

// DispatchAgentUpdated dispatches an agent.updated event
func (d *Dispatcher) DispatchAgentUpdated(agent *AgentEvent) error {
	event := &Event{
		Type:      EventAgentUpdated,
		Timestamp: getCurrentTime(),
		Data: map[string]interface{}{
			"agent": map[string]interface{}{
				"name":    agent.Name,
				"address": agent.Address,
				"version": agent.Version,
			},
		},
	}

	return d.Dispatch(event)
}

// DispatchTaskStarted dispatches a task.started event
func (d *Dispatcher) DispatchTaskStarted(task *TaskEvent) error {
	event := &Event{
		Type:      EventTaskStarted,
		Timestamp: getCurrentTime(),
		Data: map[string]interface{}{
			"task": map[string]interface{}{
				"task_name":  task.TaskName,
				"agent_name": task.AgentName,
				"status":     task.Status,
			},
		},
		Stack:  task.Stack,
		Agent:  task.AgentName,
		RunID:  task.RunID,
	}

	return d.Dispatch(event)
}

// DispatchTaskCompleted dispatches a task.completed event
func (d *Dispatcher) DispatchTaskCompleted(task *TaskEvent) error {
	event := &Event{
		Type:      EventTaskCompleted,
		Timestamp: getCurrentTime(),
		Data: map[string]interface{}{
			"task": map[string]interface{}{
				"task_name":  task.TaskName,
				"agent_name": task.AgentName,
				"status":     task.Status,
				"exit_code":  task.ExitCode,
				"duration":   task.Duration,
			},
		},
		Stack:  task.Stack,
		Agent:  task.AgentName,
		RunID:  task.RunID,
	}

	return d.Dispatch(event)
}

// DispatchTaskFailed dispatches a task.failed event
func (d *Dispatcher) DispatchTaskFailed(task *TaskEvent) error {
	event := &Event{
		Type:      EventTaskFailed,
		Timestamp: getCurrentTime(),
		Data: map[string]interface{}{
			"task": map[string]interface{}{
				"task_name":  task.TaskName,
				"agent_name": task.AgentName,
				"status":     task.Status,
				"exit_code":  task.ExitCode,
				"error":      task.Error,
				"duration":   task.Duration,
			},
		},
		Stack:  task.Stack,
		Agent:  task.AgentName,
		RunID:  task.RunID,
	}

	return d.Dispatch(event)
}

// Enable enables the dispatcher
func (d *Dispatcher) Enable() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.enabled = true
	slog.Info("hook dispatcher enabled")
}

// Disable disables the dispatcher
func (d *Dispatcher) Disable() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.enabled = false
	slog.Info("hook dispatcher disabled")
}

// IsEnabled returns whether the dispatcher is enabled
func (d *Dispatcher) IsEnabled() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.enabled
}

// StartEventProcessor starts the background event processor workers
func (d *Dispatcher) StartEventProcessor() {
	d.mu.Lock()
	if d.processing {
		d.mu.Unlock()
		return
	}
	d.processing = true
	d.mu.Unlock()

	// Start worker pool
	for i := 0; i < d.maxWorkers; i++ {
		d.workerWg.Add(1)
		go d.eventWorker(i)
	}

	// Start fallback processor for events that couldn't be sent to channel
	d.workerWg.Add(1)
	go d.fallbackProcessor()

	slog.Info("event processor started",
		"workers", d.maxWorkers,
		"channel_buffer", cap(d.eventChannel))
}

// StopEventProcessor stops the background event processor
func (d *Dispatcher) StopEventProcessor() {
	d.mu.Lock()
	if !d.processing {
		d.mu.Unlock()
		return
	}
	d.mu.Unlock()

	// Close stop channel to signal all workers
	close(d.stopChan)

	// Close event channel to drain remaining events
	close(d.eventChannel)

	// Wait for all workers to finish
	d.workerWg.Wait()

	d.mu.Lock()
	d.processing = false
	d.mu.Unlock()

	slog.Info("event processor stopped")
}

// eventWorker is a worker goroutine that processes events from the channel
func (d *Dispatcher) eventWorker(workerID int) {
	defer d.workerWg.Done()

	slog.Debug("event worker started", "worker_id", workerID)

	for {
		select {
		case <-d.stopChan:
			// Drain remaining events from channel before stopping
			for event := range d.eventChannel {
				slog.Debug("worker draining event",
					"worker_id", workerID,
					"event_id", event.ID)
				d.processEvent(event)
			}
			slog.Debug("event worker stopped", "worker_id", workerID)
			return
		case event, ok := <-d.eventChannel:
			if !ok {
				// Channel closed
				slog.Debug("event worker channel closed", "worker_id", workerID)
				return
			}
			slog.Debug("worker processing event",
				"worker_id", workerID,
				"event_id", event.ID,
				"event_type", event.Type)

			// Process event in this goroutine (each event gets its own processing)
			d.processEvent(event)
		}
	}
}

// fallbackProcessor processes pending events from database (fallback for channel overflow)
func (d *Dispatcher) fallbackProcessor() {
	defer d.workerWg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	slog.Debug("fallback processor started")

	for {
		select {
		case <-d.stopChan:
			slog.Debug("fallback processor stopped")
			return
		case <-ticker.C:
			d.processNextEvents()
		}
	}
}

// processNextEvents processes pending events from the queue (fallback)
func (d *Dispatcher) processNextEvents() {
	// Get pending events (up to 10 at a time)
	events, err := d.repo.EventQueue.GetPendingEvents(10)
	if err != nil {
		slog.Error("failed to get pending events", "error", err)
		return
	}

	if len(events) == 0 {
		// Also check for stuck events in "processing" state for more than 30 seconds
		stuckEvents, err := d.repo.EventQueue.GetStuckProcessingEvents(30, 10)
		if err != nil {
			slog.Error("failed to get stuck processing events", "error", err)
			return
		}

		if len(stuckEvents) > 0 {
			slog.Warn("found stuck processing events, reprocessing", "count", len(stuckEvents))
			events = stuckEvents
		} else {
			return
		}
	}

	slog.Debug("fallback processing events", "count", len(events))

	// Process each event
	for _, event := range events {
		d.processEvent(event)
	}
}

// processEvent processes a single event
func (d *Dispatcher) processEvent(event *Event) {
	// Update event status to processing
	if err := d.repo.EventQueue.UpdateEventStatus(event.ID, EventStatusProcessing, ""); err != nil {
		slog.Error("failed to update event status", "event_id", event.ID, "error", err)
		return
	}

	// Get all enabled hooks for this event type
	hooks, err := d.repo.ListByEventType(event.Type)
	if err != nil {
		errMsg := fmt.Sprintf("failed to list hooks: %v", err)
		d.repo.EventQueue.UpdateEventStatus(event.ID, EventStatusFailed, errMsg)
		slog.Error("failed to list hooks for event", "event_id", event.ID, "error", err)
		return
	}

	if len(hooks) == 0 {
		// No hooks for this event type, mark as completed
		d.repo.EventQueue.UpdateEventStatus(event.ID, EventStatusCompleted, "")
		slog.Debug("no hooks registered for event", "event_id", event.ID, "event_type", event.Type)
		return
	}

	slog.Info("processing event",
		"event_id", event.ID,
		"event_type", event.Type,
		"hook_count", len(hooks))

	// Execute hooks concurrently
	var wg sync.WaitGroup
	var errorMsgs []string
	var errorMu sync.Mutex

	for _, hook := range hooks {
		wg.Add(1)
		go func(h *Hook) {
			defer wg.Done()

			result, err := d.executor.Execute(h, event)
			if err != nil {
				errorMu.Lock()
				errorMsgs = append(errorMsgs, fmt.Sprintf("hook %s: %v", h.Name, err))
				errorMu.Unlock()
				slog.Error("hook execution failed", "hook_name", h.Name, "error", err)
			} else {
				// Record hook execution for this event
				if err := d.repo.EventQueue.RecordEventHookExecution(event.ID, h.ID, h.Name, result); err != nil {
					slog.Error("failed to record event hook execution", "hook_name", h.Name, "error", err)
				}

				// Also record in regular hook_executions table
				if err := d.repo.RecordExecution(result); err != nil {
					slog.Error("failed to record hook execution", "hook_name", h.Name, "error", err)
				}

				if !result.Success {
					errorMu.Lock()
					errorMsgs = append(errorMsgs, fmt.Sprintf("hook %s: %s", h.Name, result.Error))
					errorMu.Unlock()
					slog.Warn("hook execution returned failure", "hook_name", h.Name, "error", result.Error)
				} else {
					slog.Info("hook executed successfully", "hook_name", h.Name, "duration", result.Duration)
				}
			}
		}(hook)
	}

	wg.Wait()

	// Update event status
	if len(errorMsgs) > 0 {
		errMsg := fmt.Sprintf("hooks failed: %v", errorMsgs)
		d.repo.EventQueue.UpdateEventStatus(event.ID, EventStatusFailed, errMsg)
	} else {
		d.repo.EventQueue.UpdateEventStatus(event.ID, EventStatusCompleted, "")
	}
}

// Helper function to get current time (for testing)
var getCurrentTime = func() time.Time {
	return time.Now()
}

// CreateEventDispatcherFunc creates a function that can be used with the event module
// This returns a closure that captures the dispatcher instance
func (d *Dispatcher) CreateEventDispatcherFunc() func(eventType string, data map[string]interface{}) error {
	return func(eventType string, data map[string]interface{}) error {
		event := &Event{
			Type:      EventType(eventType),
			Timestamp: time.Now(),
			Data:      data,
			Stack:     d.GetCurrentStack(),
			Agent:     d.GetCurrentAgent(),
			RunID:     d.GetCurrentRunID(),
		}
		return d.Dispatch(event)
	}
}

// GetCurrentStack returns the current stack from execution context
func (d *Dispatcher) GetCurrentStack() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.currentStack
}

// GetCurrentAgent returns the current agent from execution context
func (d *Dispatcher) GetCurrentAgent() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.currentAgent
}

// GetCurrentRunID returns the current run ID from execution context
func (d *Dispatcher) GetCurrentRunID() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.currentRunID
}

// SetExecutionContext sets the current execution context
func (d *Dispatcher) SetExecutionContext(stack, agent, runID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.currentStack = stack
	d.currentAgent = agent
	d.currentRunID = runID
}
