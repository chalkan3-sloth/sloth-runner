//go:build cgo
// +build cgo

package state

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// WorkflowStateStatus represents the status of a workflow execution
type WorkflowStateStatus string

const (
	WorkflowStatusPending  WorkflowStateStatus = "pending"
	WorkflowStatusRunning  WorkflowStateStatus = "running"
	WorkflowStatusSuccess  WorkflowStateStatus = "success"
	WorkflowStatusFailed   WorkflowStateStatus = "failed"
	WorkflowStatusRolledBack WorkflowStateStatus = "rolled_back"
)

// ResourceAction represents the action performed on a resource
type ResourceAction string

const (
	ResourceActionCreate ResourceAction = "create"
	ResourceActionUpdate ResourceAction = "update"
	ResourceActionDelete ResourceAction = "delete"
	ResourceActionRead   ResourceAction = "read"
	ResourceActionNoop   ResourceAction = "noop"
)

// WorkflowState represents the state of a workflow execution
type WorkflowState struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Version     int                 `json:"version"`
	Status      WorkflowStateStatus `json:"status"`
	StartedAt   time.Time           `json:"started_at"`
	CompletedAt *time.Time          `json:"completed_at,omitempty"`
	Duration    int64               `json:"duration_seconds"`
	Metadata    map[string]string   `json:"metadata"`
	Resources   []Resource          `json:"resources"`
	Outputs     map[string]string   `json:"outputs"`
	ErrorMsg    string              `json:"error_msg,omitempty"`
	LockedBy    string              `json:"locked_by,omitempty"`
}

// Resource represents a resource managed by a workflow
type Resource struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Action     ResourceAction         `json:"action"`
	Status     string                 `json:"status"`
	Attributes map[string]interface{} `json:"attributes"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// StateVersion represents a version of the workflow state
type StateVersion struct {
	ID          string    `json:"id"`
	WorkflowID  string    `json:"workflow_id"`
	Version     int       `json:"version"`
	State       string    `json:"state"` // JSON of WorkflowState
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	Description string    `json:"description"`
}

// DriftDetection represents drift between desired and actual state
type DriftDetection struct {
	WorkflowID   string                 `json:"workflow_id"`
	ResourceID   string                 `json:"resource_id"`
	ResourceType string                 `json:"resource_type"`
	DetectedAt   time.Time              `json:"detected_at"`
	Expected     map[string]interface{} `json:"expected"`
	Actual       map[string]interface{} `json:"actual"`
	Drifted      bool                   `json:"drifted"`
}

// InitWorkflowSchema initializes the workflow state schema
func (sm *StateManager) InitWorkflowSchema() error {
	schema := `
	-- Workflow States Table
	CREATE TABLE IF NOT EXISTS workflow_states (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		version INTEGER DEFAULT 1,
		status TEXT NOT NULL,
		started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		duration_seconds INTEGER DEFAULT 0,
		metadata TEXT, -- JSON
		error_msg TEXT,
		locked_by TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Workflow Resources Table
	CREATE TABLE IF NOT EXISTS workflow_resources (
		id TEXT PRIMARY KEY,
		workflow_id TEXT NOT NULL,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		action TEXT NOT NULL,
		status TEXT NOT NULL,
		attributes TEXT, -- JSON
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (workflow_id) REFERENCES workflow_states(id) ON DELETE CASCADE
	);

	-- Workflow Outputs Table
	CREATE TABLE IF NOT EXISTS workflow_outputs (
		workflow_id TEXT NOT NULL,
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (workflow_id, key),
		FOREIGN KEY (workflow_id) REFERENCES workflow_states(id) ON DELETE CASCADE
	);

	-- State Versions Table
	CREATE TABLE IF NOT EXISTS state_versions (
		id TEXT PRIMARY KEY,
		workflow_id TEXT NOT NULL,
		version INTEGER NOT NULL,
		state TEXT NOT NULL, -- JSON snapshot
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		description TEXT,
		FOREIGN KEY (workflow_id) REFERENCES workflow_states(id) ON DELETE CASCADE,
		UNIQUE(workflow_id, version)
	);

	-- Drift Detection Table
	CREATE TABLE IF NOT EXISTS drift_detections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workflow_id TEXT NOT NULL,
		resource_id TEXT NOT NULL,
		resource_type TEXT NOT NULL,
		detected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expected TEXT, -- JSON
		actual TEXT, -- JSON
		drifted BOOLEAN DEFAULT 0,
		FOREIGN KEY (workflow_id) REFERENCES workflow_states(id) ON DELETE CASCADE
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_workflow_states_name ON workflow_states(name);
	CREATE INDEX IF NOT EXISTS idx_workflow_states_status ON workflow_states(status);
	CREATE INDEX IF NOT EXISTS idx_workflow_resources_workflow ON workflow_resources(workflow_id);
	CREATE INDEX IF NOT EXISTS idx_workflow_resources_type ON workflow_resources(type);
	CREATE INDEX IF NOT EXISTS idx_state_versions_workflow ON state_versions(workflow_id);
	CREATE INDEX IF NOT EXISTS idx_drift_workflow ON drift_detections(workflow_id);

	-- Triggers for updated_at
	CREATE TRIGGER IF NOT EXISTS update_workflow_states_timestamp
	AFTER UPDATE ON workflow_states
	BEGIN
		UPDATE workflow_states SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

	CREATE TRIGGER IF NOT EXISTS update_workflow_resources_timestamp
	AFTER UPDATE ON workflow_resources
	BEGIN
		UPDATE workflow_resources SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;
	`

	_, err := sm.db.Exec(schema)
	return err
}

// CreateWorkflowState creates a new workflow state
func (sm *StateManager) CreateWorkflowState(state *WorkflowState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	metadataJSON, err := json.Marshal(state.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = sm.db.Exec(`
		INSERT INTO workflow_states (id, name, version, status, started_at, metadata)
		VALUES (?, ?, ?, ?, ?, ?)
	`, state.ID, state.Name, state.Version, state.Status, state.StartedAt, string(metadataJSON))

	if err != nil {
		return fmt.Errorf("failed to create workflow state: %w", err)
	}

	// Create initial version
	return sm.createVersion(state)
}

// UpdateWorkflowState updates an existing workflow state
func (sm *StateManager) UpdateWorkflowState(state *WorkflowState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	metadataJSON, err := json.Marshal(state.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	var completedAtSQL interface{}
	if state.CompletedAt != nil {
		completedAtSQL = state.CompletedAt
	}

	_, err = sm.db.Exec(`
		UPDATE workflow_states
		SET status = ?, completed_at = ?, duration_seconds = ?, metadata = ?, error_msg = ?, version = version + 1
		WHERE id = ?
	`, state.Status, completedAtSQL, state.Duration, string(metadataJSON), state.ErrorMsg, state.ID)

	if err != nil {
		return fmt.Errorf("failed to update workflow state: %w", err)
	}

	// Create new version after update
	return sm.createVersion(state)
}

// GetWorkflowState retrieves a workflow state by ID
func (sm *StateManager) GetWorkflowState(id string) (*WorkflowState, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var state WorkflowState
	var metadataJSON string
	var completedAtSQL sql.NullTime

	err := sm.db.QueryRow(`
		SELECT id, name, version, status, started_at, completed_at, duration_seconds, metadata, error_msg, locked_by
		FROM workflow_states
		WHERE id = ?
	`, id).Scan(
		&state.ID, &state.Name, &state.Version, &state.Status, &state.StartedAt,
		&completedAtSQL, &state.Duration, &metadataJSON, &state.ErrorMsg, &state.LockedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workflow state not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow state: %w", err)
	}

	if completedAtSQL.Valid {
		state.CompletedAt = &completedAtSQL.Time
	}

	if err := json.Unmarshal([]byte(metadataJSON), &state.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Load resources
	resources, err := sm.getWorkflowResources(id)
	if err != nil {
		return nil, err
	}
	state.Resources = resources

	// Load outputs
	outputs, err := sm.getWorkflowOutputs(id)
	if err != nil {
		return nil, err
	}
	state.Outputs = outputs

	return &state, nil
}

// GetWorkflowStateByName retrieves the latest workflow state by name
func (sm *StateManager) GetWorkflowStateByName(name string) (*WorkflowState, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var id string
	err := sm.db.QueryRow(`
		SELECT id FROM workflow_states
		WHERE name = ?
		ORDER BY version DESC, started_at DESC
		LIMIT 1
	`, name).Scan(&id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workflow state not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow state by name: %w", err)
	}

	sm.mu.RUnlock()
	return sm.GetWorkflowState(id)
}

// ListWorkflowStates lists all workflow states with optional filters
func (sm *StateManager) ListWorkflowStates(name string, status WorkflowStateStatus) ([]*WorkflowState, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := "SELECT id FROM workflow_states WHERE 1=1"
	args := []interface{}{}

	if name != "" {
		query += " AND name = ?"
		args = append(args, name)
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY started_at DESC"

	rows, err := sm.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflow states: %w", err)
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

// AddResource adds a resource to a workflow state
func (sm *StateManager) AddResource(workflowID string, resource *Resource) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	attributesJSON, err := json.Marshal(resource.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	_, err = sm.db.Exec(`
		INSERT INTO workflow_resources (id, workflow_id, type, name, action, status, attributes)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, resource.ID, workflowID, resource.Type, resource.Name, resource.Action, resource.Status, string(attributesJSON))

	return err
}

// UpdateResource updates a resource in a workflow state
func (sm *StateManager) UpdateResource(resource *Resource) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	attributesJSON, err := json.Marshal(resource.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	_, err = sm.db.Exec(`
		UPDATE workflow_resources
		SET action = ?, status = ?, attributes = ?
		WHERE id = ?
	`, resource.Action, resource.Status, string(attributesJSON), resource.ID)

	return err
}

// getWorkflowResources retrieves all resources for a workflow (internal helper)
func (sm *StateManager) getWorkflowResources(workflowID string) ([]Resource, error) {
	rows, err := sm.db.Query(`
		SELECT id, type, name, action, status, attributes, created_at, updated_at
		FROM workflow_resources
		WHERE workflow_id = ?
		ORDER BY created_at
	`, workflowID)

	if err != nil {
		return nil, fmt.Errorf("failed to get workflow resources: %w", err)
	}
	defer rows.Close()

	var resources []Resource
	for rows.Next() {
		var r Resource
		var attributesJSON string

		if err := rows.Scan(&r.ID, &r.Type, &r.Name, &r.Action, &r.Status, &attributesJSON, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(attributesJSON), &r.Attributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
		}

		resources = append(resources, r)
	}

	return resources, rows.Err()
}

// SetOutput sets an output value for a workflow
func (sm *StateManager) SetOutput(workflowID, key, value string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, err := sm.db.Exec(`
		INSERT INTO workflow_outputs (workflow_id, key, value)
		VALUES (?, ?, ?)
		ON CONFLICT(workflow_id, key) DO UPDATE SET value = excluded.value
	`, workflowID, key, value)

	return err
}

// getWorkflowOutputs retrieves all outputs for a workflow (internal helper)
func (sm *StateManager) getWorkflowOutputs(workflowID string) (map[string]string, error) {
	rows, err := sm.db.Query(`
		SELECT key, value
		FROM workflow_outputs
		WHERE workflow_id = ?
	`, workflowID)

	if err != nil {
		return nil, fmt.Errorf("failed to get workflow outputs: %w", err)
	}
	defer rows.Close()

	outputs := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		outputs[key] = value
	}

	return outputs, rows.Err()
}

// createVersion creates a new version snapshot (internal helper)
func (sm *StateManager) createVersion(state *WorkflowState) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	versionID := fmt.Sprintf("%s-v%d", state.ID, state.Version)

	_, err = sm.db.Exec(`
		INSERT INTO state_versions (id, workflow_id, version, state, created_by, description)
		VALUES (?, ?, ?, ?, ?, ?)
	`, versionID, state.ID, state.Version, string(stateJSON), "system", fmt.Sprintf("Version %d", state.Version))

	return err
}

// GetVersions retrieves all versions of a workflow state
func (sm *StateManager) GetVersions(workflowID string) ([]StateVersion, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	rows, err := sm.db.Query(`
		SELECT id, workflow_id, version, state, created_at, created_by, description
		FROM state_versions
		WHERE workflow_id = ?
		ORDER BY version DESC
	`, workflowID)

	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}
	defer rows.Close()

	var versions []StateVersion
	for rows.Next() {
		var v StateVersion
		if err := rows.Scan(&v.ID, &v.WorkflowID, &v.Version, &v.State, &v.CreatedAt, &v.CreatedBy, &v.Description); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	return versions, rows.Err()
}

// RollbackToVersion rolls back a workflow to a specific version
func (sm *StateManager) RollbackToVersion(workflowID string, version int) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Get the version
	var stateJSON string
	err := sm.db.QueryRow(`
		SELECT state FROM state_versions
		WHERE workflow_id = ? AND version = ?
	`, workflowID, version).Scan(&stateJSON)

	if err == sql.ErrNoRows {
		return fmt.Errorf("version %d not found for workflow %s", version, workflowID)
	}
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	var state WorkflowState
	if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	// Update current state
	state.Status = WorkflowStatusRolledBack
	state.Version++

	return sm.UpdateWorkflowState(&state)
}

// DetectDrift detects drift between expected and actual state
func (sm *StateManager) DetectDrift(workflowID string, resourceID string, expected, actual map[string]interface{}) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Get resource type
	var resourceType string
	err := sm.db.QueryRow(`
		SELECT type FROM workflow_resources WHERE id = ?
	`, resourceID).Scan(&resourceType)

	if err != nil {
		return fmt.Errorf("failed to get resource type: %w", err)
	}

	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)

	drifted := !compareJSON(expected, actual)

	_, err = sm.db.Exec(`
		INSERT INTO drift_detections (workflow_id, resource_id, resource_type, expected, actual, drifted)
		VALUES (?, ?, ?, ?, ?, ?)
	`, workflowID, resourceID, resourceType, string(expectedJSON), string(actualJSON), drifted)

	return err
}

// GetDriftDetections retrieves drift detections for a workflow
func (sm *StateManager) GetDriftDetections(workflowID string) ([]DriftDetection, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	rows, err := sm.db.Query(`
		SELECT workflow_id, resource_id, resource_type, detected_at, expected, actual, drifted
		FROM drift_detections
		WHERE workflow_id = ?
		ORDER BY detected_at DESC
	`, workflowID)

	if err != nil {
		return nil, fmt.Errorf("failed to get drift detections: %w", err)
	}
	defer rows.Close()

	var detections []DriftDetection
	for rows.Next() {
		var d DriftDetection
		var expectedJSON, actualJSON string
		if err := rows.Scan(&d.WorkflowID, &d.ResourceID, &d.ResourceType, &d.DetectedAt, &expectedJSON, &actualJSON, &d.Drifted); err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(expectedJSON), &d.Expected)
		json.Unmarshal([]byte(actualJSON), &d.Actual)

		detections = append(detections, d)
	}

	return detections, rows.Err()
}

// DeleteWorkflowState deletes a workflow state and all related data
func (sm *StateManager) DeleteWorkflowState(id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, err := sm.db.Exec("DELETE FROM workflow_states WHERE id = ?", id)
	return err
}

// compareJSON compares two JSON objects for equality
func compareJSON(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bv, ok := b[k]; !ok || fmt.Sprintf("%v", v) != fmt.Sprintf("%v", bv) {
			return false
		}
	}

	return true
}
