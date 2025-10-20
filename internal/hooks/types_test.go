package hooks

import (
	"testing"
	"time"
)

func TestEventTypeConstants(t *testing.T) {
	eventTypes := []EventType{
		// Agent events
		EventAgentRegistered,
		EventAgentDisconnected,
		EventAgentHeartbeatFailed,
		EventAgentUpdated,
		EventAgentConnected,
		EventAgentVersionMismatch,
		EventAgentResourceHigh,

		// Task events
		EventTaskStarted,
		EventTaskCompleted,
		EventTaskFailed,
		EventTaskTimeout,
		EventTaskRetrying,
		EventTaskCancelled,

		// Workflow events
		EventWorkflowStarted,
		EventWorkflowCompleted,
		EventWorkflowFailed,
		EventWorkflowPaused,
		EventWorkflowResumed,
		EventWorkflowCancelled,

		// System events
		EventSystemStartup,
		EventSystemShutdown,
		EventSystemError,
		EventSystemWarning,
		EventSystemResourceHigh,
		EventSystemDiskFull,
		EventSystemMemoryLow,
		EventSystemCPUHigh,

		// Scheduler events
		EventScheduleTriggered,
		EventScheduleMissed,
		EventScheduleCreated,
		EventScheduleDeleted,
		EventScheduleUpdated,
		EventScheduleEnabled,
		EventScheduleDisabled,

		// State events
		EventStateCreated,
		EventStateUpdated,
		EventStateDeleted,
		EventStateCorrupted,
		EventStateLocked,
		EventStateUnlocked,

		// Secret events
		EventSecretCreated,
		EventSecretAccessed,
		EventSecretDeleted,
		EventSecretUpdated,
		EventSecretRotationNeeded,
		EventSecretExpired,

		// Stack events
		EventStackDeployed,
		EventStackDestroyed,
		EventStackUpdated,
		EventStackDriftDetected,
		EventStackFailed,
		EventStackSnapshotCreated,
		EventStackRolledBack,
		EventStackRollbackFailed,
		EventStackLocked,
		EventStackUnlocked,
		EventStackTagged,
		EventStackUntagged,

		// Resource events
		EventResourceCreated,
		EventResourceUpdated,
		EventResourceDeleted,
		EventResourceFailed,

		// Backup events
		EventBackupStarted,
		EventBackupCompleted,
		EventBackupFailed,
		EventRestoreStarted,
		EventRestoreCompleted,
		EventRestoreFailed,

		// Database events
		EventDBConnected,
		EventDBDisconnected,
		EventDBQuerySlow,
		EventDBError,
		EventDBMigration,

		// Network events
		EventNetworkDown,
		EventNetworkUp,
		EventNetworkSlow,
		EventNetworkLatency,

		// Security events
		EventSecurityBreach,
		EventSecurityUnauthorized,
		EventSecurityLoginFailed,
		EventSecurityLoginSuccess,
		EventSecurityPermissionDenied,

		// File system events
		EventFileCreated,
		EventFileModified,
		EventFileDeleted,
		EventFileRenamed,
		EventDirCreated,
		EventDirDeleted,

		// Deploy events
		EventDeployStarted,
		EventDeployCompleted,
		EventDeployFailed,
		EventDeployRollback,

		// Health check events
		EventHealthCheckPassed,
		EventHealthCheckFailed,
		EventHealthDegraded,
		EventHealthRecovered,

		// Custom events
		EventCustom,
	}

	for _, eventType := range eventTypes {
		if eventType == "" {
			t.Error("EventType constant is empty")
		}
	}

	// Test count (should have all defined constants)
	if len(eventTypes) != 97 {
		t.Errorf("Expected 97 event types, got %d", len(eventTypes))
	}
}

func TestEventStatusConstants(t *testing.T) {
	statuses := []EventStatus{
		EventStatusPending,
		EventStatusProcessing,
		EventStatusCompleted,
		EventStatusFailed,
	}

	for _, status := range statuses {
		if status == "" {
			t.Error("EventStatus constant is empty")
		}
	}
}

func TestHookStruct(t *testing.T) {
	now := time.Now()
	lastRun := now.Add(-1 * time.Hour)

	hook := &Hook{
		ID:          "test-hook-1",
		Name:        "Test Hook",
		Description: "Test description",
		EventType:   EventTaskCompleted,
		FilePath:    "/path/to/hook.sh",
		Stack:       "test-stack",
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastRun:     &lastRun,
		RunCount:    5,
	}

	if hook.ID != "test-hook-1" {
		t.Errorf("Hook.ID = %v, want test-hook-1", hook.ID)
	}

	if hook.EventType != EventTaskCompleted {
		t.Errorf("Hook.EventType = %v, want %v", hook.EventType, EventTaskCompleted)
	}

	if !hook.Enabled {
		t.Error("Hook.Enabled = false, want true")
	}

	if hook.RunCount != 5 {
		t.Errorf("Hook.RunCount = %v, want 5", hook.RunCount)
	}

	if hook.LastRun == nil {
		t.Error("Hook.LastRun is nil")
	}
}

func TestEventStruct(t *testing.T) {
	now := time.Now()
	processedAt := now.Add(1 * time.Second)

	event := &Event{
		ID:        "event-1",
		Type:      EventAgentRegistered,
		Timestamp: now,
		Data: map[string]interface{}{
			"key": "value",
		},
		Status:      EventStatusCompleted,
		CreatedAt:   now,
		ProcessedAt: &processedAt,
		Stack:       "test-stack",
		Agent:       "test-agent",
		RunID:       "run-123",
	}

	if event.ID != "event-1" {
		t.Errorf("Event.ID = %v, want event-1", event.ID)
	}

	if event.Type != EventAgentRegistered {
		t.Errorf("Event.Type = %v, want %v", event.Type, EventAgentRegistered)
	}

	if event.Status != EventStatusCompleted {
		t.Errorf("Event.Status = %v, want %v", event.Status, EventStatusCompleted)
	}

	if event.Stack != "test-stack" {
		t.Errorf("Event.Stack = %v, want test-stack", event.Stack)
	}

	if event.Agent != "test-agent" {
		t.Errorf("Event.Agent = %v, want test-agent", event.Agent)
	}

	if len(event.Data) != 1 {
		t.Errorf("Event.Data length = %v, want 1", len(event.Data))
	}
}

func TestAgentEventStruct(t *testing.T) {
	agentEvent := &AgentEvent{
		Name:    "agent-1",
		Address: "192.168.1.100:5000",
		Tags:    []string{"production", "us-east"},
		Version: "1.0.0",
		SystemInfo: map[string]interface{}{
			"os":   "linux",
			"arch": "amd64",
		},
	}

	if agentEvent.Name != "agent-1" {
		t.Errorf("AgentEvent.Name = %v, want agent-1", agentEvent.Name)
	}

	if len(agentEvent.Tags) != 2 {
		t.Errorf("AgentEvent.Tags length = %v, want 2", len(agentEvent.Tags))
	}

	if agentEvent.Version != "1.0.0" {
		t.Errorf("AgentEvent.Version = %v, want 1.0.0", agentEvent.Version)
	}
}

func TestTaskEventStruct(t *testing.T) {
	taskEvent := &TaskEvent{
		TaskName:  "deploy-app",
		AgentName: "agent-1",
		Status:    "completed",
		ExitCode:  0,
		Duration:  "5.2s",
		Stack:     "production",
		RunID:     "run-456",
	}

	if taskEvent.TaskName != "deploy-app" {
		t.Errorf("TaskEvent.TaskName = %v, want deploy-app", taskEvent.TaskName)
	}

	if taskEvent.ExitCode != 0 {
		t.Errorf("TaskEvent.ExitCode = %v, want 0", taskEvent.ExitCode)
	}

	if taskEvent.Stack != "production" {
		t.Errorf("TaskEvent.Stack = %v, want production", taskEvent.Stack)
	}
}

func TestHookResultStruct(t *testing.T) {
	now := time.Now()
	duration := 2 * time.Second

	result := &HookResult{
		HookID:     "hook-1",
		Success:    true,
		Output:     "Hook executed successfully",
		Duration:   duration,
		ExecutedAt: now,
	}

	if !result.Success {
		t.Error("HookResult.Success = false, want true")
	}

	if result.Duration != duration {
		t.Errorf("HookResult.Duration = %v, want %v", result.Duration, duration)
	}

	if result.Output == "" {
		t.Error("HookResult.Output is empty")
	}
}

func TestEventHookExecutionStruct(t *testing.T) {
	now := time.Now()
	duration := 1500 * time.Millisecond

	execution := &EventHookExecution{
		ID:         1,
		EventID:    "event-1",
		HookID:     "hook-1",
		HookName:   "Test Hook",
		Success:    true,
		Output:     "Execution output",
		Duration:   duration,
		ExecutedAt: now,
		EventType:  "task.completed",
		EventAgent: "agent-1",
		EventStack: "production",
		EventRunID: "run-789",
	}

	if execution.ID != 1 {
		t.Errorf("EventHookExecution.ID = %v, want 1", execution.ID)
	}

	if !execution.Success {
		t.Error("EventHookExecution.Success = false, want true")
	}

	if execution.Duration != duration {
		t.Errorf("EventHookExecution.Duration = %v, want %v", execution.Duration, duration)
	}
}

func TestFileEventStruct(t *testing.T) {
	fileEvent := &FileEvent{
		Path:      "/path/to/file.txt",
		Operation: "create",
		Size:      1024,
		Mode:      "0644",
	}

	if fileEvent.Path != "/path/to/file.txt" {
		t.Errorf("FileEvent.Path = %v, want /path/to/file.txt", fileEvent.Path)
	}

	if fileEvent.Operation != "create" {
		t.Errorf("FileEvent.Operation = %v, want create", fileEvent.Operation)
	}

	if fileEvent.Size != 1024 {
		t.Errorf("FileEvent.Size = %v, want 1024", fileEvent.Size)
	}
}

func TestCustomEventStruct(t *testing.T) {
	customEvent := &CustomEvent{
		Name: "custom.event",
		Payload: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
		Source: "workflow-123",
	}

	if customEvent.Name != "custom.event" {
		t.Errorf("CustomEvent.Name = %v, want custom.event", customEvent.Name)
	}

	if len(customEvent.Payload) != 2 {
		t.Errorf("CustomEvent.Payload length = %v, want 2", len(customEvent.Payload))
	}

	if customEvent.Source != "workflow-123" {
		t.Errorf("CustomEvent.Source = %v, want workflow-123", customEvent.Source)
	}
}

func TestFileWatcherStruct(t *testing.T) {
	now := time.Now()

	watcher := &FileWatcher{
		ID:        "watcher-1",
		Name:      "Config Watcher",
		Path:      "/etc/config",
		Pattern:   "*.yml",
		Events:    []string{"create", "modify"},
		Recursive: true,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if watcher.ID != "watcher-1" {
		t.Errorf("FileWatcher.ID = %v, want watcher-1", watcher.ID)
	}

	if watcher.Pattern != "*.yml" {
		t.Errorf("FileWatcher.Pattern = %v, want *.yml", watcher.Pattern)
	}

	if !watcher.Recursive {
		t.Error("FileWatcher.Recursive = false, want true")
	}

	if len(watcher.Events) != 2 {
		t.Errorf("FileWatcher.Events length = %v, want 2", len(watcher.Events))
	}
}

func TestSystemEventStruct(t *testing.T) {
	systemEvent := &SystemEvent{
		Type:      "error",
		Message:   "System error occurred",
		Component: "database",
		Severity:  "critical",
		Metrics: map[string]interface{}{
			"cpu":    90.5,
			"memory": 85.2,
		},
	}

	if systemEvent.Type != "error" {
		t.Errorf("SystemEvent.Type = %v, want error", systemEvent.Type)
	}

	if systemEvent.Severity != "critical" {
		t.Errorf("SystemEvent.Severity = %v, want critical", systemEvent.Severity)
	}

	if len(systemEvent.Metrics) != 2 {
		t.Errorf("SystemEvent.Metrics length = %v, want 2", len(systemEvent.Metrics))
	}
}

func TestResourceEventStruct(t *testing.T) {
	resourceEvent := &ResourceEvent{
		Resource:  "cpu",
		Current:   95.5,
		Threshold: 90.0,
		Total:     100.0,
		Available: 4.5,
		AgentName: "agent-1",
	}

	if resourceEvent.Resource != "cpu" {
		t.Errorf("ResourceEvent.Resource = %v, want cpu", resourceEvent.Resource)
	}

	if resourceEvent.Current != 95.5 {
		t.Errorf("ResourceEvent.Current = %v, want 95.5", resourceEvent.Current)
	}

	if resourceEvent.Threshold != 90.0 {
		t.Errorf("ResourceEvent.Threshold = %v, want 90.0", resourceEvent.Threshold)
	}
}

func TestScheduleEventStruct(t *testing.T) {
	scheduleEvent := &ScheduleEvent{
		ScheduleID:   "sched-1",
		ScheduleName: "Daily Backup",
		CronExpr:     "0 2 * * *",
		WorkflowName: "backup-workflow",
		NextRun:      "2024-01-20 02:00:00",
		LastRun:      "2024-01-19 02:00:00",
	}

	if scheduleEvent.ScheduleID != "sched-1" {
		t.Errorf("ScheduleEvent.ScheduleID = %v, want sched-1", scheduleEvent.ScheduleID)
	}

	if scheduleEvent.CronExpr != "0 2 * * *" {
		t.Errorf("ScheduleEvent.CronExpr = %v, want '0 2 * * *'", scheduleEvent.CronExpr)
	}
}

func TestStateEventStruct(t *testing.T) {
	stateEvent := &StateEvent{
		StateKey:  "app.config",
		Namespace: "production",
		Value:     "new-value",
		PrevValue: "old-value",
		Metadata: map[string]interface{}{
			"user": "admin",
		},
	}

	if stateEvent.StateKey != "app.config" {
		t.Errorf("StateEvent.StateKey = %v, want app.config", stateEvent.StateKey)
	}

	if stateEvent.Namespace != "production" {
		t.Errorf("StateEvent.Namespace = %v, want production", stateEvent.Namespace)
	}
}

func TestSecretEventStruct(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)

	secretEvent := &SecretEvent{
		SecretName: "api-key",
		Namespace:  "production",
		AccessedBy: "service-1",
		ExpiresAt:  expiresAt,
		RotatedAt:  now,
	}

	if secretEvent.SecretName != "api-key" {
		t.Errorf("SecretEvent.SecretName = %v, want api-key", secretEvent.SecretName)
	}

	if secretEvent.AccessedBy != "service-1" {
		t.Errorf("SecretEvent.AccessedBy = %v, want service-1", secretEvent.AccessedBy)
	}
}

func TestStackEventStruct(t *testing.T) {
	stackEvent := &StackEvent{
		StackName: "prod-infrastructure",
		Provider:  "terraform",
		Status:    "deployed",
		Resources: []string{"vpc", "subnet", "instance"},
		DriftInfo: map[string]interface{}{
			"drifted": false,
		},
		Duration: "5m30s",
	}

	if stackEvent.StackName != "prod-infrastructure" {
		t.Errorf("StackEvent.StackName = %v, want prod-infrastructure", stackEvent.StackName)
	}

	if stackEvent.Provider != "terraform" {
		t.Errorf("StackEvent.Provider = %v, want terraform", stackEvent.Provider)
	}

	if len(stackEvent.Resources) != 3 {
		t.Errorf("StackEvent.Resources length = %v, want 3", len(stackEvent.Resources))
	}
}

func TestBackupEventStruct(t *testing.T) {
	backupEvent := &BackupEvent{
		BackupID:    "backup-123",
		BackupType:  "full",
		Source:      "/data",
		Destination: "s3://bucket/backup",
		Size:        1024 * 1024 * 500, // 500MB
		Duration:    "10m",
	}

	if backupEvent.BackupType != "full" {
		t.Errorf("BackupEvent.BackupType = %v, want full", backupEvent.BackupType)
	}

	if backupEvent.Size == 0 {
		t.Error("BackupEvent.Size should not be zero")
	}
}

func TestDatabaseEventStruct(t *testing.T) {
	dbEvent := &DatabaseEvent{
		Database:     "postgres",
		Operation:    "select",
		Query:        "SELECT * FROM users",
		Duration:     250.5,
		RowsAffected: 100,
	}

	if dbEvent.Database != "postgres" {
		t.Errorf("DatabaseEvent.Database = %v, want postgres", dbEvent.Database)
	}

	if dbEvent.Duration != 250.5 {
		t.Errorf("DatabaseEvent.Duration = %v, want 250.5", dbEvent.Duration)
	}
}

func TestNetworkEventStruct(t *testing.T) {
	networkEvent := &NetworkEvent{
		Interface:  "eth0",
		Latency:    15.5,
		Bandwidth:  100.0,
		PacketLoss: 0.5,
		RemoteHost: "192.168.1.100",
	}

	if networkEvent.Interface != "eth0" {
		t.Errorf("NetworkEvent.Interface = %v, want eth0", networkEvent.Interface)
	}

	if networkEvent.Latency != 15.5 {
		t.Errorf("NetworkEvent.Latency = %v, want 15.5", networkEvent.Latency)
	}
}

func TestSecurityEventStruct(t *testing.T) {
	securityEvent := &SecurityEvent{
		User:      "admin",
		IPAddress: "192.168.1.50",
		Action:    "login",
		Resource:  "/admin/panel",
		Result:    "failure",
		Reason:    "invalid password",
		Severity:  "high",
	}

	if securityEvent.User != "admin" {
		t.Errorf("SecurityEvent.User = %v, want admin", securityEvent.User)
	}

	if securityEvent.Severity != "high" {
		t.Errorf("SecurityEvent.Severity = %v, want high", securityEvent.Severity)
	}
}

func TestDeployEventStruct(t *testing.T) {
	deployEvent := &DeployEvent{
		DeployID:    "deploy-456",
		Environment: "production",
		Version:     "v2.0.0",
		Service:     "api-service",
		Status:      "completed",
		PrevVersion: "v1.9.0",
		Artifacts:   []string{"api.tar.gz", "config.yml"},
		Duration:    "15m",
	}

	if deployEvent.Environment != "production" {
		t.Errorf("DeployEvent.Environment = %v, want production", deployEvent.Environment)
	}

	if len(deployEvent.Artifacts) != 2 {
		t.Errorf("DeployEvent.Artifacts length = %v, want 2", len(deployEvent.Artifacts))
	}
}

func TestHealthCheckEventStruct(t *testing.T) {
	healthEvent := &HealthCheckEvent{
		CheckName:    "api-health",
		Service:      "api-service",
		Status:       "healthy",
		ResponseTime: 150.5,
		Message:      "All systems operational",
		Details: map[string]interface{}{
			"database": "connected",
			"cache":    "available",
		},
	}

	if healthEvent.Status != "healthy" {
		t.Errorf("HealthCheckEvent.Status = %v, want healthy", healthEvent.Status)
	}

	if healthEvent.ResponseTime != 150.5 {
		t.Errorf("HealthCheckEvent.ResponseTime = %v, want 150.5", healthEvent.ResponseTime)
	}
}

func TestWorkflowEventStruct(t *testing.T) {
	now := time.Now()
	endTime := now.Add(10 * time.Minute)

	workflowEvent := &WorkflowEvent{
		WorkflowName: "deploy-pipeline",
		WorkflowID:   "wf-789",
		Status:       "completed",
		StartTime:    now,
		EndTime:      &endTime,
		Duration:     "10m",
		TaskCount:    5,
		FailedTasks:  []string{},
		Metadata: map[string]interface{}{
			"triggered_by": "schedule",
		},
	}

	if workflowEvent.WorkflowName != "deploy-pipeline" {
		t.Errorf("WorkflowEvent.WorkflowName = %v, want deploy-pipeline", workflowEvent.WorkflowName)
	}

	if workflowEvent.TaskCount != 5 {
		t.Errorf("WorkflowEvent.TaskCount = %v, want 5", workflowEvent.TaskCount)
	}

	if len(workflowEvent.FailedTasks) != 0 {
		t.Errorf("WorkflowEvent.FailedTasks length = %v, want 0", len(workflowEvent.FailedTasks))
	}
}
