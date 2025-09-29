package gitops

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/state"
)

// GitOpsManager manages GitOps workflows and operations
type GitOpsManager struct {
	stateManager   *state.StateManager
	repositories   map[string]*Repository
	workflows      map[string]*Workflow
	syncController *SyncController
	diffEngine     *DiffEngine
	rollbackEngine *RollbackEngine
}

// Repository represents a GitOps repository
type Repository struct {
	ID           string            `json:"id"`
	URL          string            `json:"url"`
	Branch       string            `json:"branch"`
	Path         string            `json:"path"`
	Credentials  *Credentials      `json:"credentials,omitempty"`
	PollInterval time.Duration     `json:"poll_interval"`
	LastSync     time.Time         `json:"last_sync"`
	LastCommit   string            `json:"last_commit"`
	Status       RepositoryStatus  `json:"status"`
	Metadata     map[string]string `json:"metadata"`
}

// Workflow represents a GitOps workflow configuration
type Workflow struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Repository        string                 `json:"repository"`
	TargetPath        string                 `json:"target_path"`
	AutoSync          bool                   `json:"auto_sync"`
	DiffPreview       bool                   `json:"diff_preview"`
	RollbackOnFailure bool                   `json:"rollback_on_failure"`
	SyncPolicy        SyncPolicy             `json:"sync_policy"`
	Hooks             WorkflowHooks          `json:"hooks"`
	Status            WorkflowStatus         `json:"status"`
	LastSyncResult    *SyncResult            `json:"last_sync_result,omitempty"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// SyncPolicy defines how synchronization should be performed
type SyncPolicy struct {
	AutoPrune     bool          `json:"auto_prune"`
	SyncOptions   []string      `json:"sync_options"`
	Retry         RetryPolicy   `json:"retry"`
	HealthCheck   HealthCheck   `json:"health_check"`
	PreSyncHooks  []string      `json:"pre_sync_hooks"`
	PostSyncHooks []string      `json:"post_sync_hooks"`
}

// RetryPolicy defines retry behavior for failed syncs
type RetryPolicy struct {
	Limit    int           `json:"limit"`
	Backoff  BackoffPolicy `json:"backoff"`
	Timeout  time.Duration `json:"timeout"`
}

// BackoffPolicy defines exponential backoff parameters
type BackoffPolicy struct {
	Duration    time.Duration `json:"duration"`
	Factor      float64       `json:"factor"`
	Jitter      float64       `json:"jitter"`
	MaxDuration time.Duration `json:"max_duration"`
}

// HealthCheck defines health check configuration
type HealthCheck struct {
	Enabled     bool          `json:"enabled"`
	Timeout     time.Duration `json:"timeout"`
	FailureMode string        `json:"failure_mode"` // ignore, fail, rollback
}

// WorkflowHooks defines callback functions for workflow events
type WorkflowHooks struct {
	PreSync    func(ctx context.Context, workflow *Workflow) error
	PostSync   func(ctx context.Context, workflow *Workflow, result *SyncResult) error
	OnConflict func(ctx context.Context, workflow *Workflow, conflicts []Conflict) error
	OnRollback func(ctx context.Context, workflow *Workflow, reason string) error
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	ID            string           `json:"id"`
	WorkflowID    string           `json:"workflow_id"`
	StartTime     time.Time        `json:"start_time"`
	EndTime       time.Time        `json:"end_time"`
	Status        SyncStatus       `json:"status"`
	CommitHash    string           `json:"commit_hash"`
	ChangesApplied []Change        `json:"changes_applied"`
	Conflicts     []Conflict       `json:"conflicts"`
	HealthChecks  []HealthResult   `json:"health_checks"`
	Message       string           `json:"message"`
	Metrics       SyncMetrics      `json:"metrics"`
}

// Change represents a single change applied during sync
type Change struct {
	Type        ChangeType `json:"type"`
	Resource    string     `json:"resource"`
	Action      string     `json:"action"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Error       string     `json:"error,omitempty"`
}

// Conflict represents a sync conflict
type Conflict struct {
	Resource    string                 `json:"resource"`
	Type        ConflictType           `json:"type"`
	LocalState  map[string]interface{} `json:"local_state"`
	DesiredState map[string]interface{} `json:"desired_state"`
	Resolution  string                 `json:"resolution"`
}

// HealthResult represents a health check result
type HealthResult struct {
	Resource string        `json:"resource"`
	Status   HealthStatus  `json:"status"`
	Message  string        `json:"message"`
	Duration time.Duration `json:"duration"`
}

// SyncMetrics contains metrics about the sync operation
type SyncMetrics struct {
	Duration           time.Duration `json:"duration"`
	ResourcesProcessed int           `json:"resources_processed"`
	ResourcesApplied   int           `json:"resources_applied"`
	ResourcesSkipped   int           `json:"resources_skipped"`
	ResourcesFailed    int           `json:"resources_failed"`
	ConflictsResolved  int           `json:"conflicts_resolved"`
}

// Credentials for accessing private repositories
type Credentials struct {
	Type       CredentialType `json:"type"`
	Username   string         `json:"username,omitempty"`
	Password   string         `json:"password,omitempty"`
	Token      string         `json:"token,omitempty"`
	SSHKey     string         `json:"ssh_key,omitempty"`
	SSHKeyPath string         `json:"ssh_key_path,omitempty"`
}

// Enums
type RepositoryStatus string
type WorkflowStatus string
type SyncStatus string
type ChangeType string
type ConflictType string
type HealthStatus string
type CredentialType string

const (
	// Repository Status
	RepositoryStatusActive      RepositoryStatus = "active"
	RepositoryStatusSyncing     RepositoryStatus = "syncing"
	RepositoryStatusError       RepositoryStatus = "error"
	RepositoryStatusUnreachable RepositoryStatus = "unreachable"

	// Workflow Status
	WorkflowStatusActive    WorkflowStatus = "active"
	WorkflowStatusSyncing   WorkflowStatus = "syncing"
	WorkflowStatusSynced    WorkflowStatus = "synced"
	WorkflowStatusDegraded  WorkflowStatus = "degraded"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusSuspended WorkflowStatus = "suspended"

	// Sync Status
	SyncStatusRunning   SyncStatus = "running"
	SyncStatusSucceeded SyncStatus = "succeeded"
	SyncStatusFailed    SyncStatus = "failed"
	SyncStatusCancelled SyncStatus = "cancelled"

	// Change Type
	ChangeTypeCreate ChangeType = "create"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
	ChangeTypePatch  ChangeType = "patch"

	// Conflict Type
	ConflictTypeResourceExists ConflictType = "resource_exists"
	ConflictTypeOutOfSync      ConflictType = "out_of_sync"
	ConflictTypeValidation     ConflictType = "validation"
	ConflictTypePermission     ConflictType = "permission"

	// Health Status
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusProgressing HealthStatus = "progressing"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusSuspended HealthStatus = "suspended"
	HealthStatusUnknown   HealthStatus = "unknown"

	// Credential Type
	CredentialTypeUserPass CredentialType = "userpass"
	CredentialTypeToken    CredentialType = "token"
	CredentialTypeSSH      CredentialType = "ssh"
)

// NewGitOpsManager creates a new GitOps manager
func NewGitOpsManager(stateManager *state.StateManager) *GitOpsManager {
	return &GitOpsManager{
		stateManager:   stateManager,
		repositories:   make(map[string]*Repository),
		workflows:      make(map[string]*Workflow),
		syncController: NewSyncController(),
		diffEngine:     NewDiffEngine(),
		rollbackEngine: NewRollbackEngine(),
	}
}

// RegisterRepository registers a new repository for GitOps
func (gm *GitOpsManager) RegisterRepository(ctx context.Context, repo *Repository) error {
	slog.Info("Registering GitOps repository",
		"id", repo.ID,
		"url", repo.URL,
		"branch", repo.Branch)

	// Validate repository access
	if err := gm.validateRepositoryAccess(ctx, repo); err != nil {
		return fmt.Errorf("repository validation failed: %w", err)
	}

	// Store repository
	gm.repositories[repo.ID] = repo
	repo.Status = RepositoryStatusActive
	repo.LastSync = time.Now()

	// Persist to state
	if err := gm.persistRepository(ctx, repo); err != nil {
		return fmt.Errorf("failed to persist repository: %w", err)
	}

	slog.Info("GitOps repository registered successfully", "id", repo.ID)
	return nil
}

// CreateWorkflow creates a new GitOps workflow
func (gm *GitOpsManager) CreateWorkflow(ctx context.Context, workflow *Workflow) error {
	slog.Info("Creating GitOps workflow",
		"id", workflow.ID,
		"name", workflow.Name,
		"repository", workflow.Repository)

	// Validate repository exists
	if _, exists := gm.repositories[workflow.Repository]; !exists {
		return fmt.Errorf("repository %s not found", workflow.Repository)
	}

	// Set default values
	if workflow.SyncPolicy.Retry.Limit == 0 {
		workflow.SyncPolicy.Retry.Limit = 3
	}
	if workflow.SyncPolicy.Retry.Backoff.Duration == 0 {
		workflow.SyncPolicy.Retry.Backoff.Duration = 10 * time.Second
	}
	if workflow.SyncPolicy.Retry.Backoff.Factor == 0 {
		workflow.SyncPolicy.Retry.Backoff.Factor = 2.0
	}

	workflow.Status = WorkflowStatusActive

	// Store workflow
	gm.workflows[workflow.ID] = workflow

	// Persist to state
	if err := gm.persistWorkflow(ctx, workflow); err != nil {
		return fmt.Errorf("failed to persist workflow: %w", err)
	}

	// Start initial sync if auto_sync is enabled
	if workflow.AutoSync {
		go func() {
			if err := gm.SyncWorkflow(context.Background(), workflow.ID); err != nil {
				slog.Error("Initial sync failed", "workflow", workflow.ID, "error", err)
			}
		}()
	}

	slog.Info("GitOps workflow created successfully", "id", workflow.ID)
	return nil
}

// SyncWorkflow synchronizes a workflow with its Git repository
func (gm *GitOpsManager) SyncWorkflow(ctx context.Context, workflowID string) error {
	workflow, exists := gm.workflows[workflowID]
	if !exists {
		return fmt.Errorf("workflow %s not found", workflowID)
	}

	repo, exists := gm.repositories[workflow.Repository]
	if !exists {
		return fmt.Errorf("repository %s not found", workflow.Repository)
	}

	slog.Info("Starting GitOps sync",
		"workflow", workflowID,
		"repository", repo.URL)

	// Create sync result
	syncResult := &SyncResult{
		ID:         fmt.Sprintf("sync-%d", time.Now().Unix()),
		WorkflowID: workflowID,
		StartTime:  time.Now(),
		Status:     SyncStatusRunning,
	}

	// Update workflow status
	workflow.Status = WorkflowStatusSyncing
	workflow.LastSyncResult = syncResult

	// Execute sync
	if err := gm.syncController.ExecuteSync(ctx, workflow, repo, syncResult); err != nil {
		syncResult.Status = SyncStatusFailed
		syncResult.Message = err.Error()
		workflow.Status = WorkflowStatusFailed

		// Handle rollback if enabled
		if workflow.RollbackOnFailure {
			slog.Warn("Sync failed, initiating rollback", "workflow", workflowID)
			if rollbackErr := gm.rollbackEngine.Rollback(ctx, workflow, err.Error()); rollbackErr != nil {
				slog.Error("Rollback failed", "workflow", workflowID, "error", rollbackErr)
			}
		}

		return fmt.Errorf("sync failed: %w", err)
	}

	syncResult.Status = SyncStatusSucceeded
	syncResult.EndTime = time.Now()
	workflow.Status = WorkflowStatusSynced

	// Update repository last sync
	repo.LastSync = time.Now()
	repo.LastCommit = syncResult.CommitHash

	slog.Info("GitOps sync completed successfully",
		"workflow", workflowID,
		"duration", syncResult.EndTime.Sub(syncResult.StartTime))

	return nil
}

// GetWorkflowStatus returns the current status of a workflow
func (gm *GitOpsManager) GetWorkflowStatus(workflowID string) (*Workflow, error) {
	workflow, exists := gm.workflows[workflowID]
	if !exists {
		return nil, fmt.Errorf("workflow %s not found", workflowID)
	}
	return workflow, nil
}

// ListWorkflows returns all registered workflows
func (gm *GitOpsManager) ListWorkflows() map[string]*Workflow {
	return gm.workflows
}

// GenerateDiff generates a diff preview for pending changes
func (gm *GitOpsManager) GenerateDiff(ctx context.Context, workflowID string) (*DiffResult, error) {
	workflow, exists := gm.workflows[workflowID]
	if !exists {
		return nil, fmt.Errorf("workflow %s not found", workflowID)
	}

	repo, exists := gm.repositories[workflow.Repository]
	if !exists {
		return nil, fmt.Errorf("repository %s not found", workflow.Repository)
	}

	return gm.diffEngine.GenerateDiff(ctx, workflow, repo)
}

// RollbackWorkflow rolls back a workflow to the previous state
func (gm *GitOpsManager) RollbackWorkflow(ctx context.Context, workflowID string, reason string) error {
	workflow, exists := gm.workflows[workflowID]
	if !exists {
		return fmt.Errorf("workflow %s not found", workflowID)
	}

	return gm.rollbackEngine.Rollback(ctx, workflow, reason)
}

// StartAutoSync starts automatic synchronization for all auto-sync enabled workflows
func (gm *GitOpsManager) StartAutoSync(ctx context.Context) {
	slog.Info("Starting GitOps auto-sync controller")

	for _, workflow := range gm.workflows {
		if workflow.AutoSync {
			go gm.autoSyncLoop(ctx, workflow)
		}
	}
}

// autoSyncLoop runs the auto-sync loop for a workflow
func (gm *GitOpsManager) autoSyncLoop(ctx context.Context, workflow *Workflow) {
	repo := gm.repositories[workflow.Repository]
	ticker := time.NewTicker(repo.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := gm.SyncWorkflow(ctx, workflow.ID); err != nil {
				slog.Error("Auto-sync failed",
					"workflow", workflow.ID,
					"error", err)
			}
		}
	}
}

// Helper methods
func (gm *GitOpsManager) validateRepositoryAccess(ctx context.Context, repo *Repository) error {
	// Implementation would validate Git repository access
	// For now, return success
	return nil
}

func (gm *GitOpsManager) persistRepository(ctx context.Context, repo *Repository) error {
	// Implementation would persist repository to state manager
	return nil
}

func (gm *GitOpsManager) persistWorkflow(ctx context.Context, workflow *Workflow) error {
	// Implementation would persist workflow to state manager
	return nil
}