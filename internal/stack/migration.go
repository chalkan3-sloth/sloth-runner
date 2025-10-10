//go:build cgo
// +build cgo

package stack

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Migrator handles data migration between old and new systems
type Migrator struct {
	sourceDB *sql.DB // Old workflow_state database
	targetDB *sql.DB // New unified state backend database
}

// NewMigrator creates a new migrator
func NewMigrator(sourceDBPath, targetDBPath string) (*Migrator, error) {
	sourceDB, err := sql.Open("sqlite3", sourceDBPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open source database: %w", err)
	}

	targetDB, err := sql.Open("sqlite3", targetDBPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		sourceDB.Close()
		return nil, fmt.Errorf("failed to open target database: %w", err)
	}

	return &Migrator{
		sourceDB: sourceDB,
		targetDB: targetDB,
	}, nil
}

// Close closes database connections
func (m *Migrator) Close() error {
	if m.sourceDB != nil {
		m.sourceDB.Close()
	}
	if m.targetDB != nil {
		m.targetDB.Close()
	}
	return nil
}

// MigrateWorkflowStates migrates workflow_states to stacks
func (m *Migrator) MigrateWorkflowStates() (int, error) {
	// Check if workflow_states table exists
	var tableExists bool
	err := m.sourceDB.QueryRow(`
		SELECT COUNT(*) > 0 FROM sqlite_master
		WHERE type='table' AND name='workflow_states'
	`).Scan(&tableExists)

	if err != nil || !tableExists {
		return 0, fmt.Errorf("workflow_states table not found in source database")
	}

	// Query all workflow states
	rows, err := m.sourceDB.Query(`
		SELECT id, name, version, status, started_at, completed_at,
		       duration_seconds, metadata, error_msg
		FROM workflow_states
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to query workflow states: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, name, status, metadataJSON, errorMsg string
		var version int
		var durationSeconds int64
		var startedAt time.Time
		var completedAt sql.NullTime

		err := rows.Scan(&id, &name, &version, &status, &startedAt,
			&completedAt, &durationSeconds, &metadataJSON, &errorMsg)
		if err != nil {
			continue
		}

		// Parse metadata
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
			metadata = make(map[string]interface{})
		}

		// Create corresponding stack
		stackID := id // Keep same ID for continuity
		stack := &StackState{
			ID:              stackID,
			Name:            name,
			Description:     fmt.Sprintf("Migrated from workflow_state: %s", name),
			Version:         fmt.Sprintf("v%d", version),
			Status:          status,
			CreatedAt:       startedAt,
			UpdatedAt:       time.Now(),
			WorkflowFile:    "",
			TaskResults:     make(map[string]interface{}),
			Outputs:         make(map[string]interface{}),
			Configuration:   metadata,
			Metadata:        metadata,
			ExecutionCount:  version,
			LastDuration:    time.Duration(durationSeconds) * time.Second,
			LastError:       errorMsg,
			ResourceVersion: fmt.Sprintf("%d", version),
		}

		if completedAt.Valid {
			stack.CompletedAt = &completedAt.Time
		}

		// Insert into target database
		if err := m.createStack(stack); err != nil {
			continue // Skip on error, log in production
		}

		// Migrate resources
		if err := m.migrateWorkflowResources(id, stackID); err != nil {
			continue
		}

		// Migrate outputs
		if err := m.migrateWorkflowOutputs(id, stackID); err != nil {
			continue
		}

		count++
	}

	return count, nil
}

// migrateWorkflowResources migrates resources from workflow_resources to resources
func (m *Migrator) migrateWorkflowResources(workflowID, stackID string) error {
	rows, err := m.sourceDB.Query(`
		SELECT id, type, name, action, status, attributes, created_at, updated_at
		FROM workflow_resources
		WHERE workflow_id = ?
	`, workflowID)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, resType, name, action, status, attributesJSON string
		var createdAt, updatedAt time.Time

		err := rows.Scan(&id, &resType, &name, &action, &status, &attributesJSON, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		// Parse attributes
		var attributes map[string]interface{}
		if err := json.Unmarshal([]byte(attributesJSON), &attributes); err != nil {
			attributes = make(map[string]interface{})
		}

		resource := &Resource{
			ID:           id,
			StackID:      stackID,
			Type:         resType,
			Name:         name,
			Module:       "migrated",
			Properties:   attributes,
			Dependencies: []string{},
			State:        status,
			Checksum:     "",
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
			Metadata:     make(map[string]interface{}),
		}

		resource.Metadata["migrated_action"] = action

		if err := m.createResource(resource); err != nil {
			continue
		}
	}

	return nil
}

// migrateWorkflowOutputs migrates outputs
func (m *Migrator) migrateWorkflowOutputs(workflowID, stackID string) error {
	rows, err := m.sourceDB.Query(`
		SELECT key, value
		FROM workflow_outputs
		WHERE workflow_id = ?
	`, workflowID)

	if err != nil {
		return err
	}
	defer rows.Close()

	// Collect outputs
	outputs := make(map[string]interface{})
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			continue
		}
		outputs[key] = value
	}

	// Update stack with outputs
	if len(outputs) > 0 {
		outputsJSON, _ := json.Marshal(outputs)
		_, err := m.targetDB.Exec(`
			UPDATE stacks SET outputs = ? WHERE id = ?
		`, string(outputsJSON), stackID)
		return err
	}

	return nil
}

// createStack creates a stack in target database (helper)
func (m *Migrator) createStack(stack *StackState) error {
	taskResultsJSON, _ := json.Marshal(stack.TaskResults)
	outputsJSON, _ := json.Marshal(stack.Outputs)
	configJSON, _ := json.Marshal(stack.Configuration)
	metadataJSON, _ := json.Marshal(stack.Metadata)

	var completedAt interface{}
	if stack.CompletedAt != nil {
		completedAt = *stack.CompletedAt
	}

	query := `
		INSERT OR IGNORE INTO stacks (
			id, name, description, version, status, created_at, updated_at, completed_at,
			workflow_file, task_results, outputs, configuration, metadata,
			execution_count, last_duration, last_error, resource_version
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.targetDB.Exec(query,
		stack.ID, stack.Name, stack.Description, stack.Version, stack.Status,
		stack.CreatedAt, stack.UpdatedAt, completedAt,
		stack.WorkflowFile, string(taskResultsJSON), string(outputsJSON),
		string(configJSON), string(metadataJSON),
		stack.ExecutionCount, int64(stack.LastDuration), stack.LastError, stack.ResourceVersion,
	)

	return err
}

// createResource creates a resource in target database (helper)
func (m *Migrator) createResource(resource *Resource) error {
	propertiesJSON, _ := json.Marshal(resource.Properties)
	dependenciesJSON, _ := json.Marshal(resource.Dependencies)
	metadataJSON, _ := json.Marshal(resource.Metadata)

	var lastApplied interface{}
	if resource.LastApplied != nil {
		lastApplied = *resource.LastApplied
	}

	query := `
		INSERT OR IGNORE INTO resources (
			id, stack_id, type, name, module, properties, dependencies,
			state, checksum, created_at, updated_at, last_applied,
			error_message, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.targetDB.Exec(query,
		resource.ID, resource.StackID, resource.Type, resource.Name, resource.Module,
		string(propertiesJSON), string(dependenciesJSON),
		resource.State, resource.Checksum, resource.CreatedAt, resource.UpdatedAt,
		lastApplied, resource.ErrorMessage, string(metadataJSON),
	)

	return err
}

// MigrationReport represents migration results
type MigrationReport struct {
	StacksMigrated     int       `json:"stacks_migrated"`
	ResourcesMigrated  int       `json:"resources_migrated"`
	OutputsMigrated    int       `json:"outputs_migrated"`
	Errors             []string  `json:"errors"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	Duration           string    `json:"duration"`
}

// PerformMigration performs a complete migration and returns a report
func (m *Migrator) PerformMigration() (*MigrationReport, error) {
	report := &MigrationReport{
		StartTime: time.Now(),
		Errors:    []string{},
	}

	// Migrate workflow states
	stackCount, err := m.MigrateWorkflowStates()
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Workflow states migration: %v", err))
	} else {
		report.StacksMigrated = stackCount
	}

	// Count migrated resources
	var resourceCount int
	m.targetDB.QueryRow("SELECT COUNT(*) FROM resources").Scan(&resourceCount)
	report.ResourcesMigrated = resourceCount

	report.EndTime = time.Now()
	report.Duration = report.EndTime.Sub(report.StartTime).String()

	return report, nil
}

// GenerateMigrationScript generates a SQL script for manual migration
func GenerateMigrationScript() string {
	script := `
-- Migration Script: Workflow State to Unified State Backend
-- This script migrates data from the old workflow_states system to the new unified system

-- Step 1: Migrate workflow_states to stacks
INSERT OR IGNORE INTO stacks (
	id, name, description, version, status, created_at, updated_at,
	workflow_file, execution_count, last_duration, last_error, resource_version,
	task_results, outputs, configuration, metadata
)
SELECT
	ws.id,
	ws.name,
	'Migrated from workflow_state: ' || ws.name,
	'v' || ws.version,
	ws.status,
	ws.created_at,
	datetime('now'),
	'',
	ws.version,
	ws.duration_seconds * 1000000000, -- Convert to nanoseconds
	COALESCE(ws.error_msg, ''),
	CAST(ws.version AS TEXT),
	'{}', -- task_results
	'{}', -- outputs (will be populated separately)
	ws.metadata,
	ws.metadata
FROM workflow_states ws;

-- Step 2: Migrate workflow_resources to resources
INSERT OR IGNORE INTO resources (
	id, stack_id, type, name, module, properties, dependencies,
	state, checksum, created_at, updated_at, metadata
)
SELECT
	wr.id,
	wr.workflow_id,
	wr.type,
	wr.name,
	'migrated',
	wr.attributes,
	'[]', -- No dependencies tracked in old system
	wr.status,
	'',
	wr.created_at,
	wr.updated_at,
	json_object('migrated_action', wr.action)
FROM workflow_resources wr;

-- Step 3: Migrate workflow_outputs
-- Note: This requires programmatic handling as outputs need to be aggregated
-- UPDATE stacks SET outputs = (
--   SELECT json_group_object(key, value)
--   FROM workflow_outputs
--   WHERE workflow_id = stacks.id
-- ) WHERE EXISTS (
--   SELECT 1 FROM workflow_outputs WHERE workflow_id = stacks.id
-- );

-- Step 4: Create initial snapshots for migrated stacks
-- (This should be done programmatically after migration)

-- Verification queries:
-- SELECT COUNT(*) as migrated_stacks FROM stacks WHERE description LIKE 'Migrated%';
-- SELECT COUNT(*) as migrated_resources FROM resources WHERE module = 'migrated';
`
	return script
}
