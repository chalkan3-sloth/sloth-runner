package hooks

import (
	"time"
)

// EventType represents the type of event that triggers a hook
type EventType string

const (
	// Agent events
	EventAgentRegistered      EventType = "agent.registered"
	EventAgentDisconnected    EventType = "agent.disconnected"
	EventAgentHeartbeatFailed EventType = "agent.heartbeat_failed"
	EventAgentUpdated         EventType = "agent.updated"
	EventAgentConnected       EventType = "agent.connected"
	EventAgentVersionMismatch EventType = "agent.version_mismatch"
	EventAgentResourceHigh    EventType = "agent.resource_high" // CPU/Memory alta

	// Task events
	EventTaskStarted   EventType = "task.started"
	EventTaskCompleted EventType = "task.completed"
	EventTaskFailed    EventType = "task.failed"
	EventTaskTimeout   EventType = "task.timeout"
	EventTaskRetrying  EventType = "task.retrying"
	EventTaskCancelled EventType = "task.cancelled"

	// Workflow events
	EventWorkflowStarted   EventType = "workflow.started"
	EventWorkflowCompleted EventType = "workflow.completed"
	EventWorkflowFailed    EventType = "workflow.failed"
	EventWorkflowPaused    EventType = "workflow.paused"
	EventWorkflowResumed   EventType = "workflow.resumed"
	EventWorkflowCancelled EventType = "workflow.cancelled"

	// System events
	EventSystemStartup       EventType = "system.startup"
	EventSystemShutdown      EventType = "system.shutdown"
	EventSystemError         EventType = "system.error"
	EventSystemWarning       EventType = "system.warning"
	EventSystemResourceHigh  EventType = "system.resource_high"  // CPU, memória, disco
	EventSystemDiskFull      EventType = "system.disk_full"
	EventSystemMemoryLow     EventType = "system.memory_low"
	EventSystemCPUHigh       EventType = "system.cpu_high"

	// Scheduler events
	EventScheduleTriggered EventType = "schedule.triggered"
	EventScheduleMissed    EventType = "schedule.missed"
	EventScheduleCreated   EventType = "schedule.created"
	EventScheduleDeleted   EventType = "schedule.deleted"
	EventScheduleUpdated   EventType = "schedule.updated"
	EventScheduleEnabled   EventType = "schedule.enabled"
	EventScheduleDisabled  EventType = "schedule.disabled"

	// State events
	EventStateCreated   EventType = "state.created"
	EventStateUpdated   EventType = "state.updated"
	EventStateDeleted   EventType = "state.deleted"
	EventStateCorrupted EventType = "state.corrupted"
	EventStateLocked    EventType = "state.locked"
	EventStateUnlocked  EventType = "state.unlocked"

	// Secret events
	EventSecretCreated        EventType = "secret.created"
	EventSecretAccessed       EventType = "secret.accessed"
	EventSecretDeleted        EventType = "secret.deleted"
	EventSecretUpdated        EventType = "secret.updated"
	EventSecretRotationNeeded EventType = "secret.rotation_needed"
	EventSecretExpired        EventType = "secret.expired"

	// Stack events
	EventStackDeployed         EventType = "stack.deployed"
	EventStackDestroyed        EventType = "stack.destroyed"
	EventStackUpdated          EventType = "stack.updated"
	EventStackDriftDetected    EventType = "stack.drift_detected"
	EventStackFailed           EventType = "stack.failed"
	EventStackSnapshotCreated  EventType = "stack.snapshot_created"
	EventStackRolledBack       EventType = "stack.rolled_back"
	EventStackRollbackFailed   EventType = "stack.rollback_failed"
	EventStackLocked           EventType = "stack.locked"
	EventStackUnlocked         EventType = "stack.unlocked"
	EventStackTagged           EventType = "stack.tagged"
	EventStackUntagged         EventType = "stack.untagged"

	// Resource events
	EventResourceCreated EventType = "resource.created"
	EventResourceUpdated EventType = "resource.updated"
	EventResourceDeleted EventType = "resource.deleted"
	EventResourceFailed  EventType = "resource.failed"

	// Backup events
	EventBackupStarted   EventType = "backup.started"
	EventBackupCompleted EventType = "backup.completed"
	EventBackupFailed    EventType = "backup.failed"
	EventRestoreStarted  EventType = "restore.started"
	EventRestoreCompleted EventType = "restore.completed"
	EventRestoreFailed   EventType = "restore.failed"

	// Database events
	EventDBConnected    EventType = "db.connected"
	EventDBDisconnected EventType = "db.disconnected"
	EventDBQuerySlow    EventType = "db.query_slow"
	EventDBError        EventType = "db.error"
	EventDBMigration    EventType = "db.migration"

	// Network events
	EventNetworkDown     EventType = "network.down"
	EventNetworkUp       EventType = "network.up"
	EventNetworkSlow     EventType = "network.slow"
	EventNetworkLatency  EventType = "network.latency_high"

	// Security events
	EventSecurityBreach        EventType = "security.breach"
	EventSecurityUnauthorized  EventType = "security.unauthorized"
	EventSecurityLoginFailed   EventType = "security.login_failed"
	EventSecurityLoginSuccess  EventType = "security.login_success"
	EventSecurityPermissionDenied EventType = "security.permission_denied"

	// File system events
	EventFileCreated  EventType = "file.created"
	EventFileModified EventType = "file.modified"
	EventFileDeleted  EventType = "file.deleted"
	EventFileRenamed  EventType = "file.renamed"
	EventDirCreated   EventType = "dir.created"
	EventDirDeleted   EventType = "dir.deleted"

	// Deploy events
	EventDeployStarted   EventType = "deploy.started"
	EventDeployCompleted EventType = "deploy.completed"
	EventDeployFailed    EventType = "deploy.failed"
	EventDeployRollback  EventType = "deploy.rollback"

	// Health check events
	EventHealthCheckPassed EventType = "health.check_passed"
	EventHealthCheckFailed EventType = "health.check_failed"
	EventHealthDegraded    EventType = "health.degraded"
	EventHealthRecovered   EventType = "health.recovered"

	// Custom events
	EventCustom EventType = "custom"
)

// Hook represents a registered hook
type Hook struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	EventType   EventType  `json:"event_type"`
	FilePath    string     `json:"file_path"`
	Stack       string     `json:"stack,omitempty"`        // Stack name for hook isolation
	Enabled     bool       `json:"enabled"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastRun     *time.Time `json:"last_run,omitempty"`
	RunCount    int64      `json:"run_count"`
}

// EventStatus represents the processing status of an event
type EventStatus string

const (
	EventStatusPending    EventStatus = "pending"
	EventStatusProcessing EventStatus = "processing"
	EventStatusCompleted  EventStatus = "completed"
	EventStatusFailed     EventStatus = "failed"
)

// Event represents an event that can trigger hooks
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Status    EventStatus            `json:"status"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	ProcessedAt *time.Time           `json:"processed_at,omitempty"`

	// Execution context
	Stack  string `json:"stack,omitempty"`   // Stack name being executed
	Agent  string `json:"agent,omitempty"`   // Agent executing the workflow
	RunID  string `json:"run_id,omitempty"`  // Unique run identifier
}

// AgentEvent contains agent-specific event data
type AgentEvent struct {
	Name       string            `json:"name"`
	Address    string            `json:"address"`
	Tags       []string          `json:"tags"`
	Version    string            `json:"version"`
	SystemInfo map[string]interface{} `json:"system_info,omitempty"`
}

// TaskEvent contains task-specific event data
type TaskEvent struct {
	TaskName   string `json:"task_name"`
	AgentName  string `json:"agent_name"`
	Status     string `json:"status"`
	ExitCode   int32  `json:"exit_code,omitempty"`
	Error      string `json:"error,omitempty"`
	Duration   string `json:"duration"`

	// Execution context
	Stack  string `json:"stack,omitempty"`   // Stack name being executed
	RunID  string `json:"run_id,omitempty"`  // Unique run identifier
}

// HookResult represents the result of hook execution
type HookResult struct {
	HookID     string        `json:"hook_id"`
	Success    bool          `json:"success"`
	Output     string        `json:"output"`
	Error      string        `json:"error,omitempty"`
	Duration   time.Duration `json:"duration"`
	ExecutedAt time.Time     `json:"executed_at"`
}

// EventHookExecution represents a hook execution linked to an event
type EventHookExecution struct {
	ID         int64         `json:"id"`
	EventID    string        `json:"event_id"`
	HookID     string        `json:"hook_id"`
	HookName   string        `json:"hook_name"`
	Success    bool          `json:"success"`
	Output     string        `json:"output"`
	Error      string        `json:"error,omitempty"`
	Duration   time.Duration `json:"duration"`
	ExecutedAt time.Time     `json:"executed_at"`

	// Additional event info (populated when querying by agent)
	EventType  string `json:"event_type,omitempty"`
	EventAgent string `json:"event_agent,omitempty"`
	EventStack string `json:"event_stack,omitempty"`
	EventRunID string `json:"event_run_id,omitempty"`
}

// FileEvent contains file system event data
type FileEvent struct {
	Path      string `json:"path"`
	Operation string `json:"operation"` // create, modify, delete
	Size      int64  `json:"size,omitempty"`
	Mode      string `json:"mode,omitempty"`
}

// CustomEvent contains custom event data
type CustomEvent struct {
	Name    string                 `json:"name"`
	Payload map[string]interface{} `json:"payload"`
	Source  string                 `json:"source"` // workflow name or source identifier
}

// FileWatcher represents a file system watcher configuration
type FileWatcher struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Pattern     string    `json:"pattern"`      // glob pattern
	Events      []string  `json:"events"`       // create, modify, delete
	Recursive   bool      `json:"recursive"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SystemEvent contains system-level event data
type SystemEvent struct {
	Type        string                 `json:"type"`        // startup, shutdown, error, warning
	Message     string                 `json:"message"`
	Component   string                 `json:"component"`   // which component triggered the event
	Severity    string                 `json:"severity"`    // info, warning, error, critical
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
}

// ResourceEvent contains resource monitoring event data
type ResourceEvent struct {
	Resource    string  `json:"resource"`    // cpu, memory, disk, network
	Current     float64 `json:"current"`     // current usage percentage
	Threshold   float64 `json:"threshold"`   // threshold that was exceeded
	Total       float64 `json:"total,omitempty"`
	Available   float64 `json:"available,omitempty"`
	AgentName   string  `json:"agent_name,omitempty"`
}

// ScheduleEvent contains scheduler event data
type ScheduleEvent struct {
	ScheduleID   string `json:"schedule_id"`
	ScheduleName string `json:"schedule_name"`
	CronExpr     string `json:"cron_expr,omitempty"`
	WorkflowName string `json:"workflow_name,omitempty"`
	NextRun      string `json:"next_run,omitempty"`
	LastRun      string `json:"last_run,omitempty"`
}

// StateEvent contains state management event data
type StateEvent struct {
	StateKey    string                 `json:"state_key"`
	Namespace   string                 `json:"namespace,omitempty"`
	Value       interface{}            `json:"value,omitempty"`
	PrevValue   interface{}            `json:"prev_value,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SecretEvent contains secret management event data
type SecretEvent struct {
	SecretName  string    `json:"secret_name"`
	Namespace   string    `json:"namespace,omitempty"`
	AccessedBy  string    `json:"accessed_by,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	RotatedAt   time.Time `json:"rotated_at,omitempty"`
}

// StackEvent contains infrastructure stack event data
type StackEvent struct {
	StackName   string                 `json:"stack_name"`
	Provider    string                 `json:"provider"` // terraform, pulumi, cloudformation
	Status      string                 `json:"status"`
	Resources   []string               `json:"resources,omitempty"`
	DriftInfo   map[string]interface{} `json:"drift_info,omitempty"`
	Duration    string                 `json:"duration,omitempty"`
}

// BackupEvent contains backup/restore event data
type BackupEvent struct {
	BackupID    string `json:"backup_id"`
	BackupType  string `json:"backup_type"`  // full, incremental, differential
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Size        int64  `json:"size,omitempty"`
	Duration    string `json:"duration,omitempty"`
}

// DatabaseEvent contains database event data
type DatabaseEvent struct {
	Database    string  `json:"database"`
	Operation   string  `json:"operation"`
	Query       string  `json:"query,omitempty"`
	Duration    float64 `json:"duration,omitempty"` // milliseconds
	RowsAffected int64  `json:"rows_affected,omitempty"`
	Error       string  `json:"error,omitempty"`
}

// NetworkEvent contains network event data
type NetworkEvent struct {
	Interface   string  `json:"interface,omitempty"`
	Latency     float64 `json:"latency,omitempty"`      // milliseconds
	Bandwidth   float64 `json:"bandwidth,omitempty"`    // Mbps
	PacketLoss  float64 `json:"packet_loss,omitempty"`  // percentage
	RemoteHost  string  `json:"remote_host,omitempty"`
}

// SecurityEvent contains security event data
type SecurityEvent struct {
	User        string `json:"user,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
	Action      string `json:"action"`
	Resource    string `json:"resource,omitempty"`
	Result      string `json:"result"`        // success, failure
	Reason      string `json:"reason,omitempty"`
	Severity    string `json:"severity"`      // low, medium, high, critical
}

// DeployEvent contains deployment event data
type DeployEvent struct {
	DeployID    string   `json:"deploy_id"`
	Environment string   `json:"environment"`   // dev, staging, production
	Version     string   `json:"version"`
	Service     string   `json:"service"`
	Status      string   `json:"status"`
	PrevVersion string   `json:"prev_version,omitempty"`
	Artifacts   []string `json:"artifacts,omitempty"`
	Duration    string   `json:"duration,omitempty"`
}

// HealthCheckEvent contains health check event data
type HealthCheckEvent struct {
	CheckName   string                 `json:"check_name"`
	Service     string                 `json:"service"`
	Status      string                 `json:"status"`        // healthy, degraded, unhealthy
	ResponseTime float64               `json:"response_time"` // milliseconds
	Message     string                 `json:"message,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// WorkflowEvent contains workflow event data
type WorkflowEvent struct {
	WorkflowName string                 `json:"workflow_name"`
	WorkflowID   string                 `json:"workflow_id,omitempty"`
	Status       string                 `json:"status"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      *time.Time             `json:"end_time,omitempty"`
	Duration     string                 `json:"duration,omitempty"`
	TaskCount    int                    `json:"task_count,omitempty"`
	FailedTasks  []string               `json:"failed_tasks,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}
