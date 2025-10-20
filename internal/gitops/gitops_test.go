package gitops

import (
	"testing"
	"time"
)

// Test Repository structure
func TestRepository_Structure(t *testing.T) {
	repo := &Repository{
		ID:           "repo-1",
		URL:          "https://github.com/user/repo.git",
		Branch:       "main",
		Path:         "/workflows",
		PollInterval: 5 * time.Minute,
		LastSync:     time.Now(),
		LastCommit:   "abc123",
		Status:       RepositoryStatusActive,
		Metadata:     map[string]string{"env": "production"},
	}

	if repo.ID != "repo-1" {
		t.Error("Expected ID to be set")
	}

	if repo.URL == "" {
		t.Error("Expected URL to be set")
	}

	if repo.Branch != "main" {
		t.Error("Expected Branch to be 'main'")
	}
}

func TestRepository_DefaultBranch(t *testing.T) {
	repo := &Repository{
		Branch: "main",
	}

	if repo.Branch != "main" {
		t.Error("Expected main branch")
	}
}

func TestRepository_AlternativeBranch(t *testing.T) {
	repo := &Repository{
		Branch: "develop",
	}

	if repo.Branch != "develop" {
		t.Error("Expected develop branch")
	}
}

func TestRepository_PollInterval(t *testing.T) {
	intervals := []time.Duration{
		1 * time.Minute,
		5 * time.Minute,
		15 * time.Minute,
		30 * time.Minute,
	}

	for _, interval := range intervals {
		repo := &Repository{
			PollInterval: interval,
		}

		if repo.PollInterval != interval {
			t.Errorf("Expected interval %v, got %v", interval, repo.PollInterval)
		}
	}
}

// Test Workflow structure
func TestWorkflow_Structure(t *testing.T) {
	workflow := &Workflow{
		ID:                "workflow-1",
		Name:              "production-sync",
		Repository:        "repo-1",
		TargetPath:        "/deployments",
		AutoSync:          true,
		DiffPreview:       true,
		RollbackOnFailure: true,
		Status:            WorkflowStatusActive,
		Metadata:          map[string]interface{}{"team": "platform"},
	}

	if workflow.ID == "" {
		t.Error("Expected ID to be set")
	}

	if workflow.Name == "" {
		t.Error("Expected Name to be set")
	}

	if !workflow.AutoSync {
		t.Error("Expected AutoSync to be true")
	}
}

func TestWorkflow_AutoSync(t *testing.T) {
	workflow := &Workflow{
		AutoSync: true,
	}

	if !workflow.AutoSync {
		t.Error("Expected AutoSync to be enabled")
	}
}

func TestWorkflow_DiffPreview(t *testing.T) {
	workflow := &Workflow{
		DiffPreview: true,
	}

	if !workflow.DiffPreview {
		t.Error("Expected DiffPreview to be enabled")
	}
}

func TestWorkflow_RollbackOnFailure(t *testing.T) {
	workflow := &Workflow{
		RollbackOnFailure: true,
	}

	if !workflow.RollbackOnFailure {
		t.Error("Expected RollbackOnFailure to be enabled")
	}
}

// Test SyncPolicy structure
func TestSyncPolicy_Structure(t *testing.T) {
	policy := SyncPolicy{
		AutoPrune:     true,
		SyncOptions:   []string{"CreateNamespace=true"},
		PreSyncHooks:  []string{"validate"},
		PostSyncHooks: []string{"notify"},
	}

	if !policy.AutoPrune {
		t.Error("Expected AutoPrune to be true")
	}

	if len(policy.SyncOptions) == 0 {
		t.Error("Expected sync options")
	}
}

func TestSyncPolicy_AutoPrune(t *testing.T) {
	policy := SyncPolicy{
		AutoPrune: true,
	}

	if !policy.AutoPrune {
		t.Error("Expected AutoPrune enabled")
	}
}

func TestSyncPolicy_Hooks(t *testing.T) {
	policy := SyncPolicy{
		PreSyncHooks:  []string{"hook1", "hook2"},
		PostSyncHooks: []string{"hook3"},
	}

	if len(policy.PreSyncHooks) != 2 {
		t.Errorf("Expected 2 pre-sync hooks, got %d", len(policy.PreSyncHooks))
	}

	if len(policy.PostSyncHooks) != 1 {
		t.Errorf("Expected 1 post-sync hook, got %d", len(policy.PostSyncHooks))
	}
}

// Test RetryPolicy structure
func TestRetryPolicy_Structure(t *testing.T) {
	retry := RetryPolicy{
		Limit:   3,
		Timeout: 10 * time.Minute,
	}

	if retry.Limit != 3 {
		t.Error("Expected Limit to be 3")
	}

	if retry.Timeout != 10*time.Minute {
		t.Error("Expected Timeout to be 10 minutes")
	}
}

func TestRetryPolicy_Limits(t *testing.T) {
	limits := []int{1, 3, 5, 10}

	for _, limit := range limits {
		retry := RetryPolicy{
			Limit: limit,
		}

		if retry.Limit != limit {
			t.Errorf("Expected limit %d, got %d", limit, retry.Limit)
		}
	}
}

// Test BackoffPolicy structure
func TestBackoffPolicy_Structure(t *testing.T) {
	backoff := BackoffPolicy{
		Duration:    1 * time.Second,
		Factor:      2.0,
		Jitter:      0.1,
		MaxDuration: 1 * time.Minute,
	}

	if backoff.Duration != 1*time.Second {
		t.Error("Expected Duration to be 1 second")
	}

	if backoff.Factor != 2.0 {
		t.Error("Expected Factor to be 2.0")
	}

	if backoff.Jitter != 0.1 {
		t.Error("Expected Jitter to be 0.1")
	}
}

func TestBackoffPolicy_ExponentialFactors(t *testing.T) {
	factors := []float64{1.5, 2.0, 2.5, 3.0}

	for _, factor := range factors {
		backoff := BackoffPolicy{
			Factor: factor,
		}

		if backoff.Factor != factor {
			t.Errorf("Expected factor %f, got %f", factor, backoff.Factor)
		}
	}
}

// Test HealthCheck structure
func TestHealthCheck_Structure(t *testing.T) {
	health := HealthCheck{
		Enabled:     true,
		Timeout:     30 * time.Second,
		FailureMode: "rollback",
	}

	if !health.Enabled {
		t.Error("Expected Enabled to be true")
	}

	if health.Timeout != 30*time.Second {
		t.Error("Expected Timeout to be 30 seconds")
	}

	if health.FailureMode != "rollback" {
		t.Error("Expected FailureMode to be 'rollback'")
	}
}

func TestHealthCheck_FailureModes(t *testing.T) {
	modes := []string{"ignore", "fail", "rollback"}

	for _, mode := range modes {
		health := HealthCheck{
			FailureMode: mode,
		}

		if health.FailureMode != mode {
			t.Errorf("Expected mode %s, got %s", mode, health.FailureMode)
		}
	}
}

// Test SyncResult structure
func TestSyncResult_Structure(t *testing.T) {
	result := &SyncResult{
		ID:             "sync-1",
		WorkflowID:     "workflow-1",
		StartTime:      time.Now(),
		EndTime:        time.Now().Add(5 * time.Minute),
		Status:         SyncStatusSucceeded,
		CommitHash:     "abc123",
		ChangesApplied: []Change{},
		Conflicts:      []Conflict{},
		Message:        "Sync completed successfully",
	}

	if result.ID == "" {
		t.Error("Expected ID to be set")
	}

	if result.Status != SyncStatusSucceeded {
		t.Error("Expected Status to be succeeded")
	}
}

func TestSyncResult_Success(t *testing.T) {
	result := &SyncResult{
		Status: SyncStatusSucceeded,
	}

	if result.Status != SyncStatusSucceeded {
		t.Error("Expected succeeded status")
	}
}

func TestSyncResult_Failure(t *testing.T) {
	result := &SyncResult{
		Status: SyncStatusFailed,
	}

	if result.Status != SyncStatusFailed {
		t.Error("Expected failed status")
	}
}

// Test Change structure
func TestChange_Structure(t *testing.T) {
	change := Change{
		Type:        ChangeTypeCreate,
		Resource:    "deployment/app",
		Action:      "create",
		Description: "Create new deployment",
		Status:      "applied",
	}

	if change.Type != ChangeTypeCreate {
		t.Error("Expected ChangeTypeCreate")
	}

	if change.Resource == "" {
		t.Error("Expected Resource to be set")
	}
}

func TestChange_Types(t *testing.T) {
	types := []ChangeType{
		ChangeTypeCreate,
		ChangeTypeUpdate,
		ChangeTypeDelete,
		ChangeTypePatch,
	}

	for _, changeType := range types {
		change := Change{
			Type: changeType,
		}

		if change.Type != changeType {
			t.Errorf("Expected type %v, got %v", changeType, change.Type)
		}
	}
}

// Test Conflict structure
func TestConflict_Structure(t *testing.T) {
	conflict := Conflict{
		Resource:     "service/api",
		Type:         ConflictTypeResourceExists,
		LocalState:   map[string]interface{}{"replicas": 3},
		DesiredState: map[string]interface{}{"replicas": 5},
		Resolution:   "merge",
	}

	if conflict.Resource == "" {
		t.Error("Expected Resource to be set")
	}

	if conflict.Type != ConflictTypeResourceExists {
		t.Error("Expected ConflictTypeResourceExists")
	}
}

func TestConflict_Types(t *testing.T) {
	types := []ConflictType{
		ConflictTypeResourceExists,
		ConflictTypeOutOfSync,
		ConflictTypeValidation,
		ConflictTypePermission,
	}

	for _, conflictType := range types {
		conflict := Conflict{
			Type: conflictType,
		}

		if conflict.Type != conflictType {
			t.Errorf("Expected type %v, got %v", conflictType, conflict.Type)
		}
	}
}

// Test HealthResult structure
func TestHealthResult_Structure(t *testing.T) {
	result := HealthResult{
		Resource: "deployment/app",
		Status:   HealthStatusHealthy,
		Message:  "All pods running",
		Duration: 2 * time.Second,
	}

	if result.Resource == "" {
		t.Error("Expected Resource to be set")
	}

	if result.Status != HealthStatusHealthy {
		t.Error("Expected HealthStatusHealthy")
	}
}

func TestHealthResult_Statuses(t *testing.T) {
	statuses := []HealthStatus{
		HealthStatusHealthy,
		HealthStatusProgressing,
		HealthStatusDegraded,
	}

	for _, status := range statuses {
		result := HealthResult{
			Status: status,
		}

		if result.Status != status {
			t.Errorf("Expected status %v, got %v", status, result.Status)
		}
	}
}

// Test SyncMetrics structure
func TestSyncMetrics_Structure(t *testing.T) {
	metrics := SyncMetrics{
		Duration:           5 * time.Minute,
		ResourcesProcessed: 10,
		ResourcesApplied:   8,
		ResourcesSkipped:   1,
		ResourcesFailed:    1,
		ConflictsResolved:  2,
	}

	if metrics.Duration == 0 {
		t.Error("Expected Duration to be set")
	}

	if metrics.ResourcesProcessed != 10 {
		t.Error("Expected ResourcesProcessed to be 10")
	}
}

func TestSyncMetrics_SuccessRate(t *testing.T) {
	metrics := SyncMetrics{
		ResourcesProcessed: 10,
		ResourcesApplied:   9,
		ResourcesFailed:    1,
	}

	if metrics.ResourcesApplied+metrics.ResourcesFailed != metrics.ResourcesProcessed {
		t.Error("Metrics should add up")
	}
}

// Test Credentials structure
func TestCredentials_Structure(t *testing.T) {
	creds := &Credentials{
		Type:     "token",
		Username: "user",
		Token:    "ghp_xxx",
	}

	if creds.Type == "" {
		t.Error("Expected Type to be set")
	}
}

func TestCredentials_Token(t *testing.T) {
	creds := &Credentials{
		Type:  "token",
		Token: "ghp_xxx",
	}

	if creds.Token == "" {
		t.Error("Expected Token to be set")
	}
}

func TestCredentials_SSHKey(t *testing.T) {
	creds := &Credentials{
		Type:       "ssh",
		SSHKeyPath: "/home/user/.ssh/id_rsa",
	}

	if creds.SSHKeyPath == "" {
		t.Error("Expected SSHKeyPath to be set")
	}
}

func TestCredentials_UsernamePassword(t *testing.T) {
	creds := &Credentials{
		Type:     "basic",
		Username: "user",
		Password: "pass",
	}

	if creds.Username == "" {
		t.Error("Expected Username to be set")
	}

	if creds.Password == "" {
		t.Error("Expected Password to be set")
	}
}

// Test Repository Status constants
func TestRepositoryStatus_Constants(t *testing.T) {
	if RepositoryStatusActive != "active" {
		t.Error("Expected 'active'")
	}

	if RepositoryStatusSyncing != "syncing" {
		t.Error("Expected 'syncing'")
	}

	if RepositoryStatusError != "error" {
		t.Error("Expected 'error'")
	}

	if RepositoryStatusUnreachable != "unreachable" {
		t.Error("Expected 'unreachable'")
	}
}

// Test Workflow Status constants
func TestWorkflowStatus_Constants(t *testing.T) {
	statuses := []WorkflowStatus{
		WorkflowStatusActive,
		WorkflowStatusSyncing,
		WorkflowStatusSynced,
		WorkflowStatusDegraded,
		WorkflowStatusFailed,
		WorkflowStatusSuspended,
	}

	for _, status := range statuses {
		if status == "" {
			t.Error("Status should not be empty")
		}
	}
}

// Test Sync Status constants
func TestSyncStatus_Constants(t *testing.T) {
	statuses := []SyncStatus{
		SyncStatusRunning,
		SyncStatusSucceeded,
		SyncStatusFailed,
		SyncStatusCancelled,
	}

	for _, status := range statuses {
		if status == "" {
			t.Error("Status should not be empty")
		}
	}
}

// Test Change Type constants
func TestChangeType_Constants(t *testing.T) {
	types := []ChangeType{
		ChangeTypeCreate,
		ChangeTypeUpdate,
		ChangeTypeDelete,
		ChangeTypePatch,
	}

	for _, changeType := range types {
		if changeType == "" {
			t.Error("ChangeType should not be empty")
		}
	}
}

// Test edge cases
func TestRepository_EmptyMetadata(t *testing.T) {
	repo := &Repository{
		Metadata: map[string]string{},
	}

	if len(repo.Metadata) != 0 {
		t.Error("Expected empty metadata")
	}
}

func TestWorkflow_EmptyMetadata(t *testing.T) {
	workflow := &Workflow{
		Metadata: map[string]interface{}{},
	}

	if len(workflow.Metadata) != 0 {
		t.Error("Expected empty metadata")
	}
}

func TestSyncResult_NoChanges(t *testing.T) {
	result := &SyncResult{
		ChangesApplied: []Change{},
	}

	if len(result.ChangesApplied) != 0 {
		t.Error("Expected no changes")
	}
}

func TestSyncResult_NoConflicts(t *testing.T) {
	result := &SyncResult{
		Conflicts: []Conflict{},
	}

	if len(result.Conflicts) != 0 {
		t.Error("Expected no conflicts")
	}
}

func TestSyncMetrics_ZeroValues(t *testing.T) {
	metrics := SyncMetrics{}

	if metrics.ResourcesProcessed != 0 {
		t.Error("Expected zero resources processed")
	}

	if metrics.Duration != 0 {
		t.Error("Expected zero duration")
	}
}

func TestCredentials_EmptyToken(t *testing.T) {
	creds := &Credentials{
		Type:  "token",
		Token: "",
	}

	if creds.Token != "" {
		t.Error("Expected empty token")
	}
}

func TestRetryPolicy_ZeroLimit(t *testing.T) {
	retry := RetryPolicy{
		Limit: 0,
	}

	if retry.Limit != 0 {
		t.Error("Expected zero limit")
	}
}

func TestBackoffPolicy_DefaultJitter(t *testing.T) {
	backoff := BackoffPolicy{
		Jitter: 0.0,
	}

	if backoff.Jitter != 0.0 {
		t.Error("Expected zero jitter")
	}
}

func TestHealthCheck_Disabled(t *testing.T) {
	health := HealthCheck{
		Enabled: false,
	}

	if health.Enabled {
		t.Error("Expected health check to be disabled")
	}
}

func TestChange_WithError(t *testing.T) {
	change := Change{
		Status: "failed",
		Error:  "resource not found",
	}

	if change.Error == "" {
		t.Error("Expected error message")
	}
}

func TestConflict_EmptyResolution(t *testing.T) {
	conflict := Conflict{
		Resolution: "",
	}

	if conflict.Resolution != "" {
		t.Error("Expected empty resolution")
	}
}

func TestRepository_WithCredentials(t *testing.T) {
	repo := &Repository{
		Credentials: &Credentials{
			Type:  "token",
			Token: "xxx",
		},
	}

	if repo.Credentials == nil {
		t.Error("Expected credentials to be set")
	}
}

func TestWorkflow_WithLastSyncResult(t *testing.T) {
	workflow := &Workflow{
		LastSyncResult: &SyncResult{
			Status: SyncStatusSucceeded,
		},
	}

	if workflow.LastSyncResult == nil {
		t.Error("Expected last sync result")
	}
}

func TestSyncPolicy_NoHooks(t *testing.T) {
	policy := SyncPolicy{
		PreSyncHooks:  []string{},
		PostSyncHooks: []string{},
	}

	if len(policy.PreSyncHooks) != 0 {
		t.Error("Expected no pre-sync hooks")
	}

	if len(policy.PostSyncHooks) != 0 {
		t.Error("Expected no post-sync hooks")
	}
}

func TestSyncMetrics_HighResourceCount(t *testing.T) {
	metrics := SyncMetrics{
		ResourcesProcessed: 1000,
		ResourcesApplied:   950,
		ResourcesFailed:    50,
	}

	if metrics.ResourcesProcessed != 1000 {
		t.Error("Expected 1000 resources processed")
	}
}
