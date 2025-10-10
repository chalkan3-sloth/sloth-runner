//go:build cgo
// +build cgo

package state

import (
	"archive/tar"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WorkflowStateExport represents an exported workflow state
type WorkflowStateExport struct {
	State    *WorkflowState  `json:"state"`
	Versions []StateVersion  `json:"versions"`
	Tags     []string        `json:"tags"`
	ExportedAt time.Time     `json:"exported_at"`
	ExportedBy string        `json:"exported_by"`
}

// StateDiff represents differences between two state versions
type StateDiff struct {
	WorkflowID    string                 `json:"workflow_id"`
	FromVersion   int                    `json:"from_version"`
	ToVersion     int                    `json:"to_version"`
	StatusChange  string                 `json:"status_change,omitempty"`
	AddedResources    []Resource         `json:"added_resources"`
	RemovedResources  []Resource         `json:"removed_resources"`
	ModifiedResources []ResourceDiff     `json:"modified_resources"`
	MetadataChanges   map[string]string  `json:"metadata_changes"`
	OutputChanges     map[string]string  `json:"output_changes"`
}

// ResourceDiff represents changes in a resource
type ResourceDiff struct {
	ResourceID string                 `json:"resource_id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Changes    map[string]interface{} `json:"changes"`
}

// StateAnalytics represents analytics data
type StateAnalytics struct {
	TotalWorkflows     int                       `json:"total_workflows"`
	TotalResources     int                       `json:"total_resources"`
	StatusDistribution map[string]int            `json:"status_distribution"`
	ResourceTypes      map[string]int            `json:"resource_types"`
	AverageDuration    float64                   `json:"average_duration"`
	SuccessRate        float64                   `json:"success_rate"`
	TopWorkflows       []WorkflowExecutionStats  `json:"top_workflows"`
	RecentActivity     []ActivityEntry           `json:"recent_activity"`
}

// WorkflowExecutionStats represents execution statistics
type WorkflowExecutionStats struct {
	WorkflowName    string  `json:"workflow_name"`
	ExecutionCount  int     `json:"execution_count"`
	SuccessCount    int     `json:"success_count"`
	FailureCount    int     `json:"failure_count"`
	AverageDuration float64 `json:"average_duration"`
}

// ActivityEntry represents a recent activity
type ActivityEntry struct {
	WorkflowID   string    `json:"workflow_id"`
	WorkflowName string    `json:"workflow_name"`
	Action       string    `json:"action"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
}

// StateSearchQuery represents search parameters
type StateSearchQuery struct {
	Name           string
	Status         []WorkflowStateStatus
	Tags           []string
	ResourceType   string
	DateFrom       *time.Time
	DateTo         *time.Time
	HasErrors      bool
	MinDuration    int64
	MaxDuration    int64
	Limit          int
}

// ExtendWorkflowSchema adds tables for advanced features
func (sm *StateManager) ExtendWorkflowSchema() error {
	schema := `
	-- Tags Table
	CREATE TABLE IF NOT EXISTS workflow_tags (
		workflow_id TEXT NOT NULL,
		tag TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (workflow_id, tag),
		FOREIGN KEY (workflow_id) REFERENCES workflow_states(id) ON DELETE CASCADE
	);

	-- Backups Table
	CREATE TABLE IF NOT EXISTS workflow_backups (
		id TEXT PRIMARY KEY,
		workflow_id TEXT NOT NULL,
		backup_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		description TEXT,
		size_bytes INTEGER,
		FOREIGN KEY (workflow_id) REFERENCES workflow_states(id) ON DELETE CASCADE
	);

	-- Activity Log Table
	CREATE TABLE IF NOT EXISTS workflow_activity (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workflow_id TEXT NOT NULL,
		action TEXT NOT NULL,
		details TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		FOREIGN KEY (workflow_id) REFERENCES workflow_states(id) ON DELETE CASCADE
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_workflow_tags_tag ON workflow_tags(tag);
	CREATE INDEX IF NOT EXISTS idx_workflow_backups_workflow ON workflow_backups(workflow_id);
	CREATE INDEX IF NOT EXISTS idx_workflow_activity_workflow ON workflow_activity(workflow_id);
	CREATE INDEX IF NOT EXISTS idx_workflow_activity_created ON workflow_activity(created_at);
	`

	_, err := sm.db.Exec(schema)
	return err
}

// AddTag adds a tag to a workflow
func (sm *StateManager) AddTag(workflowID string, tag string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, err := sm.db.Exec(`
		INSERT OR IGNORE INTO workflow_tags (workflow_id, tag)
		VALUES (?, ?)
	`, workflowID, tag)

	return err
}

// RemoveTag removes a tag from a workflow
func (sm *StateManager) RemoveTag(workflowID string, tag string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, err := sm.db.Exec(`
		DELETE FROM workflow_tags
		WHERE workflow_id = ? AND tag = ?
	`, workflowID, tag)

	return err
}

// GetTags retrieves all tags for a workflow
func (sm *StateManager) GetTags(workflowID string) ([]string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	rows, err := sm.db.Query(`
		SELECT tag FROM workflow_tags
		WHERE workflow_id = ?
		ORDER BY tag
	`, workflowID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// ExportWorkflowState exports a workflow state to JSON
func (sm *StateManager) ExportWorkflowState(workflowID, exportedBy string) (*WorkflowStateExport, error) {
	state, err := sm.GetWorkflowState(workflowID)
	if err != nil {
		return nil, err
	}

	versions, err := sm.GetVersions(workflowID)
	if err != nil {
		return nil, err
	}

	tags, err := sm.GetTags(workflowID)
	if err != nil {
		tags = []string{} // Continue even if tags fail
	}

	export := &WorkflowStateExport{
		State:      state,
		Versions:   versions,
		Tags:       tags,
		ExportedAt: time.Now(),
		ExportedBy: exportedBy,
	}

	return export, nil
}

// ImportWorkflowState imports a workflow state from JSON
func (sm *StateManager) ImportWorkflowState(export *WorkflowStateExport, overwrite bool) error {
	// Check if workflow already exists
	existing, err := sm.GetWorkflowState(export.State.ID)
	if err == nil && existing != nil && !overwrite {
		return fmt.Errorf("workflow %s already exists (use overwrite flag)", export.State.ID)
	}

	// Delete existing if overwrite
	if existing != nil && overwrite {
		if err := sm.DeleteWorkflowState(export.State.ID); err != nil {
			return fmt.Errorf("failed to delete existing workflow: %w", err)
		}
	}

	// Create workflow state
	if err := sm.CreateWorkflowState(export.State); err != nil {
		return fmt.Errorf("failed to create workflow state: %w", err)
	}

	// Add resources
	for _, resource := range export.State.Resources {
		if err := sm.AddResource(export.State.ID, &resource); err != nil {
			return fmt.Errorf("failed to add resource: %w", err)
		}
	}

	// Add outputs
	for key, value := range export.State.Outputs {
		if err := sm.SetOutput(export.State.ID, key, value); err != nil {
			return fmt.Errorf("failed to set output: %w", err)
		}
	}

	// Add tags
	for _, tag := range export.Tags {
		sm.AddTag(export.State.ID, tag)
	}

	// Log activity
	sm.LogActivity(export.State.ID, "import", "Imported workflow state", export.ExportedBy)

	return nil
}

// BackupWorkflowState creates a compressed backup of workflow state
func (sm *StateManager) BackupWorkflowState(workflowID, backupDir, createdBy string) (string, error) {
	// Export state
	export, err := sm.ExportWorkflowState(workflowID, createdBy)
	if err != nil {
		return "", err
	}

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.tar.gz", export.State.Name, timestamp)
	backupPath := filepath.Join(backupDir, filename)

	// Create tar.gz file
	file, err := os.Create(backupPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Marshal export to JSON
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return "", err
	}

	// Write to tar
	header := &tar.Header{
		Name:    "state.json",
		Mode:    0644,
		Size:    int64(len(data)),
		ModTime: time.Now(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return "", err
	}

	if _, err := tarWriter.Write(data); err != nil {
		return "", err
	}

	// Get file size
	fileInfo, _ := os.Stat(backupPath)
	var size int64
	if fileInfo != nil {
		size = fileInfo.Size()
	}

	// Record backup in database
	backupID := fmt.Sprintf("%s-%s", workflowID, timestamp)
	sm.mu.Lock()
	_, err = sm.db.Exec(`
		INSERT INTO workflow_backups (id, workflow_id, backup_path, created_by, description, size_bytes)
		VALUES (?, ?, ?, ?, ?, ?)
	`, backupID, workflowID, backupPath, createdBy, fmt.Sprintf("Backup created at %s", timestamp), size)
	sm.mu.Unlock()

	if err != nil {
		return "", fmt.Errorf("failed to record backup: %w", err)
	}

	// Log activity
	sm.LogActivity(workflowID, "backup", fmt.Sprintf("Created backup: %s", filename), createdBy)

	return backupPath, nil
}

// RestoreWorkflowState restores a workflow state from backup
func (sm *StateManager) RestoreWorkflowState(backupPath string, overwrite bool) error {
	// Open backup file
	file, err := os.Open(backupPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	// Read tar contents
	var export WorkflowStateExport
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Name == "state.json" {
			data, err := io.ReadAll(tarReader)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(data, &export); err != nil {
				return err
			}
			break
		}
	}

	// Import the state
	return sm.ImportWorkflowState(&export, overwrite)
}

// DiffVersions compares two versions of a workflow state
func (sm *StateManager) DiffVersions(workflowID string, fromVersion, toVersion int) (*StateDiff, error) {
	// Get both versions
	versions, err := sm.GetVersions(workflowID)
	if err != nil {
		return nil, err
	}

	var fromState, toState *WorkflowState
	for _, v := range versions {
		if v.Version == fromVersion {
			if err := json.Unmarshal([]byte(v.State), &fromState); err != nil {
				return nil, err
			}
		}
		if v.Version == toVersion {
			if err := json.Unmarshal([]byte(v.State), &toState); err != nil {
				return nil, err
			}
		}
	}

	if fromState == nil || toState == nil {
		return nil, fmt.Errorf("one or both versions not found")
	}

	diff := &StateDiff{
		WorkflowID:  workflowID,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		MetadataChanges: make(map[string]string),
		OutputChanges:   make(map[string]string),
	}

	// Status change
	if fromState.Status != toState.Status {
		diff.StatusChange = fmt.Sprintf("%s -> %s", fromState.Status, toState.Status)
	}

	// Resource diff
	fromResources := make(map[string]Resource)
	toResources := make(map[string]Resource)

	for _, r := range fromState.Resources {
		fromResources[r.ID] = r
	}
	for _, r := range toState.Resources {
		toResources[r.ID] = r
	}

	// Find added and removed resources
	for id, r := range toResources {
		if _, exists := fromResources[id]; !exists {
			diff.AddedResources = append(diff.AddedResources, r)
		}
	}

	for id, r := range fromResources {
		if _, exists := toResources[id]; !exists {
			diff.RemovedResources = append(diff.RemovedResources, r)
		}
	}

	// Find modified resources
	for id, toRes := range toResources {
		if fromRes, exists := fromResources[id]; exists {
			if !compareJSON(fromRes.Attributes, toRes.Attributes) || fromRes.Status != toRes.Status {
				resourceDiff := ResourceDiff{
					ResourceID: id,
					Name:       toRes.Name,
					Type:       toRes.Type,
					Changes:    make(map[string]interface{}),
				}

				if fromRes.Status != toRes.Status {
					resourceDiff.Changes["status"] = map[string]string{
						"from": fromRes.Status,
						"to":   toRes.Status,
					}
				}

				diff.ModifiedResources = append(diff.ModifiedResources, resourceDiff)
			}
		}
	}

	// Metadata changes
	for key, toVal := range toState.Metadata {
		if fromVal, exists := fromState.Metadata[key]; !exists || fromVal != toVal {
			diff.MetadataChanges[key] = fmt.Sprintf("%s -> %s", fromVal, toVal)
		}
	}

	// Output changes
	for key, toVal := range toState.Outputs {
		if fromVal, exists := fromState.Outputs[key]; !exists || fromVal != toVal {
			diff.OutputChanges[key] = fmt.Sprintf("%s -> %s", fromVal, toVal)
		}
	}

	return diff, nil
}

// SearchWorkflows performs advanced search on workflows
func (sm *StateManager) SearchWorkflows(query StateSearchQuery) ([]*WorkflowState, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sql := "SELECT DISTINCT ws.id FROM workflow_states ws"
	args := []interface{}{}
	conditions := []string{}

	// Join with tags if needed
	if len(query.Tags) > 0 {
		sql += " LEFT JOIN workflow_tags wt ON ws.id = wt.workflow_id"
		placeholders := []string{}
		for _, tag := range query.Tags {
			placeholders = append(placeholders, "?")
			args = append(args, tag)
		}
		conditions = append(conditions, fmt.Sprintf("wt.tag IN (%s)", strings.Join(placeholders, ",")))
	}

	// Join with resources if needed
	if query.ResourceType != "" {
		sql += " LEFT JOIN workflow_resources wr ON ws.id = wr.workflow_id"
		conditions = append(conditions, "wr.type = ?")
		args = append(args, query.ResourceType)
	}

	// Name filter
	if query.Name != "" {
		conditions = append(conditions, "ws.name LIKE ?")
		args = append(args, "%"+query.Name+"%")
	}

	// Status filter
	if len(query.Status) > 0 {
		placeholders := []string{}
		for _, status := range query.Status {
			placeholders = append(placeholders, "?")
			args = append(args, string(status))
		}
		conditions = append(conditions, fmt.Sprintf("ws.status IN (%s)", strings.Join(placeholders, ",")))
	}

	// Date range
	if query.DateFrom != nil {
		conditions = append(conditions, "ws.started_at >= ?")
		args = append(args, query.DateFrom)
	}
	if query.DateTo != nil {
		conditions = append(conditions, "ws.started_at <= ?")
		args = append(args, query.DateTo)
	}

	// Duration range
	if query.MinDuration > 0 {
		conditions = append(conditions, "ws.duration_seconds >= ?")
		args = append(args, query.MinDuration)
	}
	if query.MaxDuration > 0 {
		conditions = append(conditions, "ws.duration_seconds <= ?")
		args = append(args, query.MaxDuration)
	}

	// Errors filter
	if query.HasErrors {
		conditions = append(conditions, "ws.error_msg IS NOT NULL AND ws.error_msg != ''")
	}

	// Build WHERE clause
	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	sql += " ORDER BY ws.started_at DESC"

	// Limit
	if query.Limit > 0 {
		sql += " LIMIT ?"
		args = append(args, query.Limit)
	}

	rows, err := sm.db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search workflows: %w", err)
	}
	defer rows.Close()

	var states []*WorkflowState
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		sm.mu.RUnlock()
		state, err := sm.GetWorkflowState(id)
		sm.mu.RLock()
		if err != nil {
			return nil, err
		}

		states = append(states, state)
	}

	return states, rows.Err()
}

// PruneOldStates removes workflow states older than specified duration
func (sm *StateManager) PruneOldStates(olderThan time.Duration, keepSuccessful bool) (int, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	cutoffTime := time.Now().Add(-olderThan)

	var query string
	if keepSuccessful {
		query = `DELETE FROM workflow_states WHERE started_at < ? AND status != ?`
	} else {
		query = `DELETE FROM workflow_states WHERE started_at < ?`
	}

	var result sql.Result
	var err error

	if keepSuccessful {
		result, err = sm.db.Exec(query, cutoffTime, string(WorkflowStatusSuccess))
	} else {
		result, err = sm.db.Exec(query, cutoffTime)
	}

	if err != nil {
		return 0, err
	}

	affected, _ := result.RowsAffected()
	return int(affected), nil
}

// GetAnalytics returns comprehensive analytics
func (sm *StateManager) GetAnalytics() (*StateAnalytics, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	analytics := &StateAnalytics{
		StatusDistribution: make(map[string]int),
		ResourceTypes:      make(map[string]int),
	}

	// Total workflows
	sm.db.QueryRow("SELECT COUNT(*) FROM workflow_states").Scan(&analytics.TotalWorkflows)

	// Total resources
	sm.db.QueryRow("SELECT COUNT(*) FROM workflow_resources").Scan(&analytics.TotalResources)

	// Status distribution
	rows, err := sm.db.Query("SELECT status, COUNT(*) FROM workflow_states GROUP BY status")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var status string
			var count int
			rows.Scan(&status, &count)
			analytics.StatusDistribution[status] = count
		}
	}

	// Resource types
	rows, err = sm.db.Query("SELECT type, COUNT(*) FROM workflow_resources GROUP BY type")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rtype string
			var count int
			rows.Scan(&rtype, &count)
			analytics.ResourceTypes[rtype] = count
		}
	}

	// Average duration
	sm.db.QueryRow("SELECT AVG(duration_seconds) FROM workflow_states WHERE duration_seconds > 0").Scan(&analytics.AverageDuration)

	// Success rate
	var successCount int
	sm.db.QueryRow("SELECT COUNT(*) FROM workflow_states WHERE status = ?", string(WorkflowStatusSuccess)).Scan(&successCount)
	if analytics.TotalWorkflows > 0 {
		analytics.SuccessRate = float64(successCount) / float64(analytics.TotalWorkflows) * 100
	}

	// Top workflows
	rows, err = sm.db.Query(`
		SELECT name,
			   COUNT(*) as exec_count,
			   SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_count,
			   SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failure_count,
			   AVG(duration_seconds) as avg_duration
		FROM workflow_states
		GROUP BY name
		ORDER BY exec_count DESC
		LIMIT 10
	`)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stats WorkflowExecutionStats
			rows.Scan(&stats.WorkflowName, &stats.ExecutionCount, &stats.SuccessCount, &stats.FailureCount, &stats.AverageDuration)
			analytics.TopWorkflows = append(analytics.TopWorkflows, stats)
		}
	}

	// Recent activity
	rows, err = sm.db.Query(`
		SELECT id, name, status, started_at
		FROM workflow_states
		ORDER BY started_at DESC
		LIMIT 20
	`)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var activity ActivityEntry
			rows.Scan(&activity.WorkflowID, &activity.WorkflowName, &activity.Status, &activity.Timestamp)
			activity.Action = "execution"
			analytics.RecentActivity = append(analytics.RecentActivity, activity)
		}
	}

	return analytics, nil
}

// LogActivity logs an activity for a workflow
func (sm *StateManager) LogActivity(workflowID, action, details, createdBy string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, err := sm.db.Exec(`
		INSERT INTO workflow_activity (workflow_id, action, details, created_by)
		VALUES (?, ?, ?, ?)
	`, workflowID, action, details, createdBy)

	return err
}

// GetActivity retrieves activity log for a workflow
func (sm *StateManager) GetActivity(workflowID string, limit int) ([]ActivityEntry, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := `
		SELECT workflow_id, action, details, created_at, created_by
		FROM workflow_activity
		WHERE workflow_id = ?
		ORDER BY created_at DESC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := sm.db.Query(query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []ActivityEntry
	for rows.Next() {
		var activity ActivityEntry
		var createdBy sql.NullString
		rows.Scan(&activity.WorkflowID, &activity.Action, &activity.Status, &activity.Timestamp, &createdBy)
		if createdBy.Valid {
			activity.Status = createdBy.String
		}
		activities = append(activities, activity)
	}

	return activities, rows.Err()
}
