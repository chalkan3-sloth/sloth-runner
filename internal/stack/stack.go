package stack

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// StackState represents the state of a workflow stack
type StackState struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Version         string                 `json:"version"`
	Status          string                 `json:"status"` // created, running, completed, failed
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	WorkflowFile    string                 `json:"workflow_file"`
	TaskResults     map[string]interface{} `json:"task_results"`
	Outputs         map[string]interface{} `json:"outputs"`
	Configuration   map[string]interface{} `json:"configuration"`
	Metadata        map[string]interface{} `json:"metadata"`
	ExecutionCount  int                    `json:"execution_count"`
	LastDuration    time.Duration          `json:"last_duration"`
	LastError       string                 `json:"last_error,omitempty"`
	ResourceVersion string                 `json:"resource_version"`
}

// StackManager manages workflow stacks and their state
type StackManager struct {
	db   *sql.DB
	mu   sync.RWMutex
	path string
}

// NewStackManager creates a new stack manager
func NewStackManager(dbPath string) (*StackManager, error) {
	if dbPath == "" {
		// Use /etc/sloth-runner/ as the default location
		dbPath = "/etc/sloth-runner/stacks.db"
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create stack directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sm := &StackManager{
		db:   db,
		path: dbPath,
	}

	if err := sm.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return sm, nil
}

// initSchema creates the required database tables
func (sm *StackManager) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS stacks (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		version TEXT,
		status TEXT NOT NULL DEFAULT 'created',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		workflow_file TEXT,
		task_results TEXT,
		outputs TEXT,
		configuration TEXT,
		metadata TEXT,
		execution_count INTEGER DEFAULT 0,
		last_duration INTEGER DEFAULT 0,
		last_error TEXT,
		resource_version TEXT NOT NULL DEFAULT '1'
	);

	CREATE INDEX IF NOT EXISTS idx_stacks_name ON stacks(name);
	CREATE INDEX IF NOT EXISTS idx_stacks_status ON stacks(status);
	CREATE INDEX IF NOT EXISTS idx_stacks_updated_at ON stacks(updated_at);

	CREATE TABLE IF NOT EXISTS stack_executions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		stack_id TEXT NOT NULL,
		started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		duration INTEGER,
		status TEXT NOT NULL,
		task_count INTEGER DEFAULT 0,
		success_count INTEGER DEFAULT 0,
		failure_count INTEGER DEFAULT 0,
		outputs TEXT,
		error_message TEXT,
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_executions_stack_id ON stack_executions(stack_id);
	CREATE INDEX IF NOT EXISTS idx_executions_started_at ON stack_executions(started_at);

	CREATE TABLE IF NOT EXISTS resources (
		id TEXT PRIMARY KEY,
		stack_id TEXT NOT NULL,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		module TEXT NOT NULL,
		properties TEXT,
		dependencies TEXT,
		state TEXT NOT NULL DEFAULT 'pending',
		checksum TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_applied DATETIME,
		error_message TEXT,
		metadata TEXT,
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_resources_stack_id ON resources(stack_id);
	CREATE INDEX IF NOT EXISTS idx_resources_type ON resources(type);
	CREATE INDEX IF NOT EXISTS idx_resources_state ON resources(state);
	CREATE INDEX IF NOT EXISTS idx_resources_checksum ON resources(checksum);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_resources_stack_type_name ON resources(stack_id, type, name);
	`

	if _, err := sm.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// CreateStack creates a new stack
func (sm *StackManager) CreateStack(stack *StackState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if stack.ID == "" {
		return fmt.Errorf("stack ID is required")
	}

	if stack.Name == "" {
		return fmt.Errorf("stack name is required")
	}

	// Set defaults
	stack.CreatedAt = time.Now()
	stack.UpdatedAt = time.Now()
	stack.Status = "created"
	stack.ResourceVersion = "1"

	// Serialize JSON fields
	taskResultsJSON, _ := json.Marshal(stack.TaskResults)
	outputsJSON, _ := json.Marshal(stack.Outputs)
	configJSON, _ := json.Marshal(stack.Configuration)
	metadataJSON, _ := json.Marshal(stack.Metadata)

	query := `
		INSERT INTO stacks (
			id, name, description, version, status, created_at, updated_at,
			workflow_file, task_results, outputs, configuration, metadata,
			execution_count, last_duration, resource_version
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := sm.db.Exec(query,
		stack.ID, stack.Name, stack.Description, stack.Version, stack.Status,
		stack.CreatedAt, stack.UpdatedAt, stack.WorkflowFile,
		string(taskResultsJSON), string(outputsJSON), string(configJSON), string(metadataJSON),
		stack.ExecutionCount, int64(stack.LastDuration), stack.ResourceVersion,
	)

	if err != nil {
		return fmt.Errorf("failed to create stack: %w", err)
	}

	return nil
}

// GetStack retrieves a stack by ID
func (sm *StackManager) GetStack(id string) (*StackState, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := `
		SELECT id, name, description, version, status, created_at, updated_at, completed_at,
		       workflow_file, task_results, outputs, configuration, metadata,
		       execution_count, last_duration, COALESCE(last_error, '') as last_error, resource_version
		FROM stacks WHERE id = ?
	`

	row := sm.db.QueryRow(query, id)

	var stack StackState
	var taskResultsJSON, outputsJSON, configJSON, metadataJSON string
	var completedAt sql.NullTime
	var lastDurationNanos int64

	err := row.Scan(
		&stack.ID, &stack.Name, &stack.Description, &stack.Version, &stack.Status,
		&stack.CreatedAt, &stack.UpdatedAt, &completedAt,
		&stack.WorkflowFile, &taskResultsJSON, &outputsJSON, &configJSON, &metadataJSON,
		&stack.ExecutionCount, &lastDurationNanos, &stack.LastError, &stack.ResourceVersion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stack not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get stack: %w", err)
	}

	// Handle nullable completed_at
	if completedAt.Valid {
		stack.CompletedAt = &completedAt.Time
	}

	// Parse duration
	stack.LastDuration = time.Duration(lastDurationNanos)

	// Deserialize JSON fields
	if err := json.Unmarshal([]byte(taskResultsJSON), &stack.TaskResults); err != nil {
		stack.TaskResults = make(map[string]interface{})
	}
	if err := json.Unmarshal([]byte(outputsJSON), &stack.Outputs); err != nil {
		stack.Outputs = make(map[string]interface{})
	}
	if err := json.Unmarshal([]byte(configJSON), &stack.Configuration); err != nil {
		stack.Configuration = make(map[string]interface{})
	}
	if err := json.Unmarshal([]byte(metadataJSON), &stack.Metadata); err != nil {
		stack.Metadata = make(map[string]interface{})
	}

	return &stack, nil
}

// GetStackByName retrieves a stack by name
func (sm *StackManager) GetStackByName(name string) (*StackState, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := `
		SELECT id, name, description, version, status, created_at, updated_at, completed_at,
		       workflow_file, task_results, outputs, configuration, metadata,
		       execution_count, last_duration, COALESCE(last_error, '') as last_error, resource_version
		FROM stacks WHERE name = ? ORDER BY updated_at DESC LIMIT 1
	`

	row := sm.db.QueryRow(query, name)

	var stack StackState
	var taskResultsJSON, outputsJSON, configJSON, metadataJSON string
	var completedAt sql.NullTime
	var lastDurationNanos int64

	err := row.Scan(
		&stack.ID, &stack.Name, &stack.Description, &stack.Version, &stack.Status,
		&stack.CreatedAt, &stack.UpdatedAt, &completedAt,
		&stack.WorkflowFile, &taskResultsJSON, &outputsJSON, &configJSON, &metadataJSON,
		&stack.ExecutionCount, &lastDurationNanos, &stack.LastError, &stack.ResourceVersion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stack not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get stack: %w", err)
	}

	// Handle nullable completed_at
	if completedAt.Valid {
		stack.CompletedAt = &completedAt.Time
	}

	// Parse duration
	stack.LastDuration = time.Duration(lastDurationNanos)

	// Deserialize JSON fields
	if err := json.Unmarshal([]byte(taskResultsJSON), &stack.TaskResults); err != nil {
		stack.TaskResults = make(map[string]interface{})
	}
	if err := json.Unmarshal([]byte(outputsJSON), &stack.Outputs); err != nil {
		stack.Outputs = make(map[string]interface{})
	}
	if err := json.Unmarshal([]byte(configJSON), &stack.Configuration); err != nil {
		stack.Configuration = make(map[string]interface{})
	}
	if err := json.Unmarshal([]byte(metadataJSON), &stack.Metadata); err != nil {
		stack.Metadata = make(map[string]interface{})
	}

	return &stack, nil
}

// UpdateStack updates an existing stack
func (sm *StackManager) UpdateStack(stack *StackState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	stack.UpdatedAt = time.Now()

	// Serialize JSON fields
	taskResultsJSON, _ := json.Marshal(stack.TaskResults)
	outputsJSON, _ := json.Marshal(stack.Outputs)
	configJSON, _ := json.Marshal(stack.Configuration)
	metadataJSON, _ := json.Marshal(stack.Metadata)

	query := `
		UPDATE stacks SET
			name = ?, description = ?, version = ?, status = ?, updated_at = ?, completed_at = ?,
			workflow_file = ?, task_results = ?, outputs = ?, configuration = ?, metadata = ?,
			execution_count = ?, last_duration = ?, last_error = ?, resource_version = ?
		WHERE id = ?
	`

	var completedAt interface{}
	if stack.CompletedAt != nil {
		completedAt = *stack.CompletedAt
	}

	_, err := sm.db.Exec(query,
		stack.Name, stack.Description, stack.Version, stack.Status, stack.UpdatedAt, completedAt,
		stack.WorkflowFile, string(taskResultsJSON), string(outputsJSON), string(configJSON), string(metadataJSON),
		stack.ExecutionCount, int64(stack.LastDuration), stack.LastError, stack.ResourceVersion,
		stack.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update stack: %w", err)
	}

	return nil
}

// ListStacks lists all stacks
func (sm *StackManager) ListStacks() ([]*StackState, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := `
		SELECT id, name, description, version, status, created_at, updated_at, completed_at,
		       workflow_file, task_results, outputs, configuration, metadata,
		       execution_count, last_duration, COALESCE(last_error, '') as last_error, resource_version
		FROM stacks ORDER BY updated_at DESC
	`

	rows, err := sm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list stacks: %w", err)
	}
	defer rows.Close()

	var stacks []*StackState

	for rows.Next() {
		var stack StackState
		var taskResultsJSON, outputsJSON, configJSON, metadataJSON string
		var completedAt sql.NullTime
		var lastDurationNanos int64

		err := rows.Scan(
			&stack.ID, &stack.Name, &stack.Description, &stack.Version, &stack.Status,
			&stack.CreatedAt, &stack.UpdatedAt, &completedAt,
			&stack.WorkflowFile, &taskResultsJSON, &outputsJSON, &configJSON, &metadataJSON,
			&stack.ExecutionCount, &lastDurationNanos, &stack.LastError, &stack.ResourceVersion,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan stack: %w", err)
		}

		// Handle nullable completed_at
		if completedAt.Valid {
			stack.CompletedAt = &completedAt.Time
		}

		// Parse duration
		stack.LastDuration = time.Duration(lastDurationNanos)

		// Deserialize JSON fields
		if err := json.Unmarshal([]byte(taskResultsJSON), &stack.TaskResults); err != nil {
			stack.TaskResults = make(map[string]interface{})
		}
		if err := json.Unmarshal([]byte(outputsJSON), &stack.Outputs); err != nil {
			stack.Outputs = make(map[string]interface{})
		}
		if err := json.Unmarshal([]byte(configJSON), &stack.Configuration); err != nil {
			stack.Configuration = make(map[string]interface{})
		}
		if err := json.Unmarshal([]byte(metadataJSON), &stack.Metadata); err != nil {
			stack.Metadata = make(map[string]interface{})
		}

		stacks = append(stacks, &stack)
	}

	return stacks, nil
}

// DeleteStack deletes a stack
func (sm *StackManager) DeleteStack(id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	query := `DELETE FROM stacks WHERE id = ?`
	result, err := sm.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stack: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stack not found: %s", id)
	}

	return nil
}

// RecordExecution records a stack execution
func (sm *StackManager) RecordExecution(stackID string, execution *StackExecution) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	outputsJSON, _ := json.Marshal(execution.Outputs)

	query := `
		INSERT INTO stack_executions (
			stack_id, started_at, completed_at, duration, status,
			task_count, success_count, failure_count, outputs, error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var completedAt interface{}
	if execution.CompletedAt != nil {
		completedAt = *execution.CompletedAt
	}

	_, err := sm.db.Exec(query,
		stackID, execution.StartedAt, completedAt, int64(execution.Duration), execution.Status,
		execution.TaskCount, execution.SuccessCount, execution.FailureCount,
		string(outputsJSON), execution.ErrorMessage,
	)

	if err != nil {
		return fmt.Errorf("failed to record execution: %w", err)
	}

	return nil
}

// StackExecution represents an execution of a stack
type StackExecution struct {
	ID           int64                  `json:"id"`
	StackID      string                 `json:"stack_id"`
	StartedAt    time.Time              `json:"started_at"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
	Duration     time.Duration          `json:"duration"`
	Status       string                 `json:"status"`
	TaskCount    int                    `json:"task_count"`
	SuccessCount int                    `json:"success_count"`
	FailureCount int                    `json:"failure_count"`
	Outputs      map[string]interface{} `json:"outputs"`
	ErrorMessage string                 `json:"error_message,omitempty"`
}

// GetStackExecutions retrieves executions for a stack
func (sm *StackManager) GetStackExecutions(stackID string, limit int) ([]*StackExecution, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := `
		SELECT id, stack_id, started_at, completed_at, duration, status,
		       task_count, success_count, failure_count, outputs, error_message
		FROM stack_executions 
		WHERE stack_id = ? 
		ORDER BY started_at DESC 
		LIMIT ?
	`

	rows, err := sm.db.Query(query, stackID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions: %w", err)
	}
	defer rows.Close()

	var executions []*StackExecution

	for rows.Next() {
		var exec StackExecution
		var completedAt sql.NullTime
		var durationNanos int64
		var outputsJSON string

		err := rows.Scan(
			&exec.ID, &exec.StackID, &exec.StartedAt, &completedAt, &durationNanos, &exec.Status,
			&exec.TaskCount, &exec.SuccessCount, &exec.FailureCount, &outputsJSON, &exec.ErrorMessage,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan execution: %w", err)
		}

		// Handle nullable completed_at
		if completedAt.Valid {
			exec.CompletedAt = &completedAt.Time
		}

		// Parse duration
		exec.Duration = time.Duration(durationNanos)

		// Deserialize outputs
		if err := json.Unmarshal([]byte(outputsJSON), &exec.Outputs); err != nil {
			exec.Outputs = make(map[string]interface{})
		}

		executions = append(executions, &exec)
	}

	return executions, nil
}

// Resource represents a managed resource in a stack
type Resource struct {
	ID           string                 `json:"id"`
	StackID      string                 `json:"stack_id"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	Module       string                 `json:"module"`
	Properties   map[string]interface{} `json:"properties"`
	Dependencies []string               `json:"dependencies"`
	State        string                 `json:"state"` // pending, applied, failed, drift
	Checksum     string                 `json:"checksum"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	LastApplied  *time.Time             `json:"last_applied,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// CreateResource creates a new resource in the stack
func (sm *StackManager) CreateResource(resource *Resource) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if resource.ID == "" {
		return fmt.Errorf("resource ID is required")
	}

	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()

	propertiesJSON, _ := json.Marshal(resource.Properties)
	dependenciesJSON, _ := json.Marshal(resource.Dependencies)
	metadataJSON, _ := json.Marshal(resource.Metadata)

	query := `
		INSERT INTO resources (
			id, stack_id, type, name, module, properties, dependencies,
			state, checksum, created_at, updated_at, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := sm.db.Exec(query,
		resource.ID, resource.StackID, resource.Type, resource.Name, resource.Module,
		string(propertiesJSON), string(dependenciesJSON),
		resource.State, resource.Checksum, resource.CreatedAt, resource.UpdatedAt,
		string(metadataJSON),
	)

	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	return nil
}

// GetResource retrieves a resource by ID
func (sm *StackManager) GetResource(id string) (*Resource, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := `
		SELECT id, stack_id, type, name, module, properties, dependencies,
		       state, checksum, created_at, updated_at, last_applied,
		       COALESCE(error_message, '') as error_message, metadata
		FROM resources WHERE id = ?
	`

	row := sm.db.QueryRow(query, id)

	var resource Resource
	var propertiesJSON, dependenciesJSON, metadataJSON string
	var lastApplied sql.NullTime

	err := row.Scan(
		&resource.ID, &resource.StackID, &resource.Type, &resource.Name, &resource.Module,
		&propertiesJSON, &dependenciesJSON,
		&resource.State, &resource.Checksum, &resource.CreatedAt, &resource.UpdatedAt,
		&lastApplied, &resource.ErrorMessage, &metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("resource not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	if lastApplied.Valid {
		resource.LastApplied = &lastApplied.Time
	}

	json.Unmarshal([]byte(propertiesJSON), &resource.Properties)
	json.Unmarshal([]byte(dependenciesJSON), &resource.Dependencies)
	json.Unmarshal([]byte(metadataJSON), &resource.Metadata)

	return &resource, nil
}

// UpdateResource updates a resource
func (sm *StackManager) UpdateResource(resource *Resource) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	resource.UpdatedAt = time.Now()

	propertiesJSON, _ := json.Marshal(resource.Properties)
	dependenciesJSON, _ := json.Marshal(resource.Dependencies)
	metadataJSON, _ := json.Marshal(resource.Metadata)

	query := `
		UPDATE resources SET
			type = ?, name = ?, module = ?, properties = ?, dependencies = ?,
			state = ?, checksum = ?, updated_at = ?, last_applied = ?,
			error_message = ?, metadata = ?
		WHERE id = ?
	`

	var lastApplied interface{}
	if resource.LastApplied != nil {
		lastApplied = *resource.LastApplied
	}

	_, err := sm.db.Exec(query,
		resource.Type, resource.Name, resource.Module, string(propertiesJSON), string(dependenciesJSON),
		resource.State, resource.Checksum, resource.UpdatedAt, lastApplied,
		resource.ErrorMessage, string(metadataJSON),
		resource.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	return nil
}

// ListResources lists all resources in a stack
func (sm *StackManager) ListResources(stackID string) ([]*Resource, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := `
		SELECT id, stack_id, type, name, module, properties, dependencies,
		       state, checksum, created_at, updated_at, last_applied,
		       COALESCE(error_message, '') as error_message, metadata
		FROM resources WHERE stack_id = ? ORDER BY created_at ASC
	`

	rows, err := sm.db.Query(query, stackID)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}
	defer rows.Close()

	var resources []*Resource

	for rows.Next() {
		var resource Resource
		var propertiesJSON, dependenciesJSON, metadataJSON string
		var lastApplied sql.NullTime

		err := rows.Scan(
			&resource.ID, &resource.StackID, &resource.Type, &resource.Name, &resource.Module,
			&propertiesJSON, &dependenciesJSON,
			&resource.State, &resource.Checksum, &resource.CreatedAt, &resource.UpdatedAt,
			&lastApplied, &resource.ErrorMessage, &metadataJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan resource: %w", err)
		}

		if lastApplied.Valid {
			resource.LastApplied = &lastApplied.Time
		}

		json.Unmarshal([]byte(propertiesJSON), &resource.Properties)
		json.Unmarshal([]byte(dependenciesJSON), &resource.Dependencies)
		json.Unmarshal([]byte(metadataJSON), &resource.Metadata)

		resources = append(resources, &resource)
	}

	return resources, nil
}

// DeleteResource deletes a resource
func (sm *StackManager) DeleteResource(id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	query := `DELETE FROM resources WHERE id = ?`
	result, err := sm.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("resource not found: %s", id)
	}

	return nil
}

// GetResourceByStackAndName retrieves a resource by stack ID, type and name
func (sm *StackManager) GetResourceByStackAndName(stackID, resType, name string) (*Resource, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	query := `
		SELECT id, stack_id, type, name, module, properties, dependencies,
		       state, checksum, created_at, updated_at, last_applied,
		       COALESCE(error_message, '') as error_message, metadata
		FROM resources WHERE stack_id = ? AND type = ? AND name = ?
	`

	row := sm.db.QueryRow(query, stackID, resType, name)

	var resource Resource
	var propertiesJSON, dependenciesJSON, metadataJSON string
	var lastApplied sql.NullTime

	err := row.Scan(
		&resource.ID, &resource.StackID, &resource.Type, &resource.Name, &resource.Module,
		&propertiesJSON, &dependenciesJSON,
		&resource.State, &resource.Checksum, &resource.CreatedAt, &resource.UpdatedAt,
		&lastApplied, &resource.ErrorMessage, &metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is not an error
		}
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	if lastApplied.Valid {
		resource.LastApplied = &lastApplied.Time
	}

	json.Unmarshal([]byte(propertiesJSON), &resource.Properties)
	json.Unmarshal([]byte(dependenciesJSON), &resource.Dependencies)
	json.Unmarshal([]byte(metadataJSON), &resource.Metadata)

	return &resource, nil
}

// Close closes the database connection
func (sm *StackManager) Close() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.db.Close()
}