package hooks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Repository manages hook persistence
type Repository struct {
	db         *sql.DB
	EventQueue *EventQueue
}

// NewRepository creates a new hook repository
func NewRepository() (*Repository, error) {
	// Create .sloth-cache directory if it doesn't exist
	cacheDir := filepath.Join(".", ".sloth-cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	dbPath := filepath.Join(cacheDir, "hooks.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &Repository{
		db:         db,
		EventQueue: NewEventQueue(db),
	}
	if err := repo.initialize(); err != nil {
		return nil, err
	}

	// Initialize event queue schema
	if err := repo.EventQueue.InitializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize event queue schema: %w", err)
	}

	return repo, nil
}

// initialize creates the necessary tables
func (r *Repository) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS hooks (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		event_type TEXT NOT NULL,
		file_path TEXT NOT NULL,
		stack TEXT,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		last_run INTEGER,
		run_count INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_hooks_event_type ON hooks(event_type);
	CREATE INDEX IF NOT EXISTS idx_hooks_enabled ON hooks(enabled);
	CREATE INDEX IF NOT EXISTS idx_hooks_stack ON hooks(stack);

	CREATE TABLE IF NOT EXISTS hook_executions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hook_id TEXT NOT NULL,
		success INTEGER NOT NULL,
		output TEXT,
		error TEXT,
		duration_ms INTEGER NOT NULL,
		executed_at INTEGER NOT NULL,
		FOREIGN KEY (hook_id) REFERENCES hooks(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_executions_hook_id ON hook_executions(hook_id);
	CREATE INDEX IF NOT EXISTS idx_executions_executed_at ON hook_executions(executed_at);
	`

	_, err := r.db.Exec(schema)
	return err
}

// Add adds a new hook
func (r *Repository) Add(hook *Hook) error {
	hook.ID = uuid.New().String()
	hook.CreatedAt = time.Now()
	hook.UpdatedAt = time.Now()

	query := `
		INSERT INTO hooks (id, name, description, event_type, file_path, stack, enabled, created_at, updated_at, run_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		hook.ID,
		hook.Name,
		hook.Description,
		hook.EventType,
		hook.FilePath,
		hook.Stack,
		boolToInt(hook.Enabled),
		hook.CreatedAt.Unix(),
		hook.UpdatedAt.Unix(),
		hook.RunCount,
	)

	if err != nil {
		return fmt.Errorf("failed to add hook: %w", err)
	}

	return nil
}

// Get retrieves a hook by ID
func (r *Repository) Get(id string) (*Hook, error) {
	query := `
		SELECT id, name, description, event_type, file_path, stack, enabled,
		       created_at, updated_at, last_run, run_count
		FROM hooks
		WHERE id = ?
	`

	var hook Hook
	var enabled int
	var createdAt, updatedAt int64
	var lastRun sql.NullInt64
	var stack sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&hook.ID,
		&hook.Name,
		&hook.Description,
		&hook.EventType,
		&hook.FilePath,
		&stack,
		&enabled,
		&createdAt,
		&updatedAt,
		&lastRun,
		&hook.RunCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("hook not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get hook: %w", err)
	}

	hook.Enabled = intToBool(enabled)
	hook.CreatedAt = time.Unix(createdAt, 0)
	hook.UpdatedAt = time.Unix(updatedAt, 0)
	if stack.Valid {
		hook.Stack = stack.String
	}
	if lastRun.Valid {
		t := time.Unix(lastRun.Int64, 0)
		hook.LastRun = &t
	}

	return &hook, nil
}

// GetByName retrieves a hook by name
func (r *Repository) GetByName(name string) (*Hook, error) {
	query := `
		SELECT id, name, description, event_type, file_path, stack, enabled,
		       created_at, updated_at, last_run, run_count
		FROM hooks
		WHERE name = ?
	`

	var hook Hook
	var enabled int
	var createdAt, updatedAt int64
	var lastRun sql.NullInt64
	var stack sql.NullString

	err := r.db.QueryRow(query, name).Scan(
		&hook.ID,
		&hook.Name,
		&hook.Description,
		&hook.EventType,
		&hook.FilePath,
		&stack,
		&enabled,
		&createdAt,
		&updatedAt,
		&lastRun,
		&hook.RunCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("hook not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get hook: %w", err)
	}

	hook.Enabled = intToBool(enabled)
	hook.CreatedAt = time.Unix(createdAt, 0)
	hook.UpdatedAt = time.Unix(updatedAt, 0)
	if stack.Valid {
		hook.Stack = stack.String
	}
	if lastRun.Valid {
		t := time.Unix(lastRun.Int64, 0)
		hook.LastRun = &t
	}

	return &hook, nil
}

// List retrieves all hooks
func (r *Repository) List() ([]*Hook, error) {
	query := `
		SELECT id, name, description, event_type, file_path, stack, enabled,
		       created_at, updated_at, last_run, run_count
		FROM hooks
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list hooks: %w", err)
	}
	defer rows.Close()

	var hooks []*Hook
	for rows.Next() {
		var hook Hook
		var enabled int
		var createdAt, updatedAt int64
		var lastRun sql.NullInt64
		var stack sql.NullString

		err := rows.Scan(
			&hook.ID,
			&hook.Name,
			&hook.Description,
			&hook.EventType,
			&hook.FilePath,
			&stack,
			&enabled,
			&createdAt,
			&updatedAt,
			&lastRun,
			&hook.RunCount,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan hook: %w", err)
		}

		hook.Enabled = intToBool(enabled)
		hook.CreatedAt = time.Unix(createdAt, 0)
		hook.UpdatedAt = time.Unix(updatedAt, 0)
		if stack.Valid {
			hook.Stack = stack.String
		}
		if lastRun.Valid {
			t := time.Unix(lastRun.Int64, 0)
			hook.LastRun = &t
		}

		hooks = append(hooks, &hook)
	}

	return hooks, nil
}

// ListByEventType retrieves all enabled hooks for a specific event type
func (r *Repository) ListByEventType(eventType EventType) ([]*Hook, error) {
	query := `
		SELECT id, name, description, event_type, file_path, stack, enabled,
		       created_at, updated_at, last_run, run_count
		FROM hooks
		WHERE event_type = ? AND enabled = 1
		ORDER BY name
	`

	rows, err := r.db.Query(query, eventType)
	if err != nil {
		return nil, fmt.Errorf("failed to list hooks by event type: %w", err)
	}
	defer rows.Close()

	var hooks []*Hook
	for rows.Next() {
		var hook Hook
		var enabled int
		var createdAt, updatedAt int64
		var lastRun sql.NullInt64
		var stack sql.NullString

		err := rows.Scan(
			&hook.ID,
			&hook.Name,
			&hook.Description,
			&hook.EventType,
			&hook.FilePath,
			&stack,
			&enabled,
			&createdAt,
			&updatedAt,
			&lastRun,
			&hook.RunCount,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan hook: %w", err)
		}

		hook.Enabled = intToBool(enabled)
		hook.CreatedAt = time.Unix(createdAt, 0)
		hook.UpdatedAt = time.Unix(updatedAt, 0)
		if stack.Valid {
			hook.Stack = stack.String
		}
		if lastRun.Valid {
			t := time.Unix(lastRun.Int64, 0)
			hook.LastRun = &t
		}

		hooks = append(hooks, &hook)
	}

	return hooks, nil
}

// ListByStack retrieves all hooks for a specific stack
func (r *Repository) ListByStack(stack string) ([]*Hook, error) {
	query := `
		SELECT id, name, description, event_type, file_path, stack, enabled,
		       created_at, updated_at, last_run, run_count
		FROM hooks
		WHERE stack = ?
		ORDER BY name
	`

	rows, err := r.db.Query(query, stack)
	if err != nil {
		return nil, fmt.Errorf("failed to list hooks by stack: %w", err)
	}
	defer rows.Close()

	var hooks []*Hook
	for rows.Next() {
		var hook Hook
		var enabled int
		var createdAt, updatedAt int64
		var lastRun sql.NullInt64
		var stackVal sql.NullString

		err := rows.Scan(
			&hook.ID,
			&hook.Name,
			&hook.Description,
			&hook.EventType,
			&hook.FilePath,
			&stackVal,
			&enabled,
			&createdAt,
			&updatedAt,
			&lastRun,
			&hook.RunCount,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan hook: %w", err)
		}

		hook.Enabled = intToBool(enabled)
		hook.CreatedAt = time.Unix(createdAt, 0)
		hook.UpdatedAt = time.Unix(updatedAt, 0)
		if stackVal.Valid {
			hook.Stack = stackVal.String
		}
		if lastRun.Valid {
			t := time.Unix(lastRun.Int64, 0)
			hook.LastRun = &t
		}

		hooks = append(hooks, &hook)
	}

	return hooks, nil
}

// Update updates an existing hook
func (r *Repository) Update(hook *Hook) error {
	hook.UpdatedAt = time.Now()

	query := `
		UPDATE hooks
		SET name = ?, description = ?, event_type = ?, file_path = ?,
		    stack = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		hook.Name,
		hook.Description,
		hook.EventType,
		hook.FilePath,
		hook.Stack,
		boolToInt(hook.Enabled),
		hook.UpdatedAt.Unix(),
		hook.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update hook: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("hook not found: %s", hook.ID)
	}

	return nil
}

// Delete removes a hook
func (r *Repository) Delete(id string) error {
	query := `DELETE FROM hooks WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete hook: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("hook not found: %s", id)
	}

	return nil
}

// Enable enables a hook
func (r *Repository) Enable(id string) error {
	query := `UPDATE hooks SET enabled = 1, updated_at = ? WHERE id = ?`

	result, err := r.db.Exec(query, time.Now().Unix(), id)
	if err != nil {
		return fmt.Errorf("failed to enable hook: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("hook not found: %s", id)
	}

	return nil
}

// Disable disables a hook
func (r *Repository) Disable(id string) error {
	query := `UPDATE hooks SET enabled = 0, updated_at = ? WHERE id = ?`

	result, err := r.db.Exec(query, time.Now().Unix(), id)
	if err != nil {
		return fmt.Errorf("failed to disable hook: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("hook not found: %s", id)
	}

	return nil
}

// RecordExecution records a hook execution
func (r *Repository) RecordExecution(result *HookResult) error {
	// Update hook's last_run and run_count
	updateQuery := `
		UPDATE hooks
		SET last_run = ?, run_count = run_count + 1, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(updateQuery,
		result.ExecutedAt.Unix(),
		time.Now().Unix(),
		result.HookID,
	)

	if err != nil {
		return fmt.Errorf("failed to update hook stats: %w", err)
	}

	// Insert execution record
	insertQuery := `
		INSERT INTO hook_executions (hook_id, success, output, error, duration_ms, executed_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(insertQuery,
		result.HookID,
		boolToInt(result.Success),
		result.Output,
		result.Error,
		result.Duration.Milliseconds(),
		result.ExecutedAt.Unix(),
	)

	if err != nil {
		return fmt.Errorf("failed to record execution: %w", err)
	}

	return nil
}

// GetExecutionHistory retrieves execution history for a hook
func (r *Repository) GetExecutionHistory(hookID string, limit int) ([]*HookResult, error) {
	query := `
		SELECT hook_id, success, output, error, duration_ms, executed_at
		FROM hook_executions
		WHERE hook_id = ?
		ORDER BY executed_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, hookID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution history: %w", err)
	}
	defer rows.Close()

	var results []*HookResult
	for rows.Next() {
		var result HookResult
		var success int
		var durationMS int64
		var executedAt int64

		err := rows.Scan(
			&result.HookID,
			&success,
			&result.Output,
			&result.Error,
			&durationMS,
			&executedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan execution: %w", err)
		}

		result.Success = intToBool(success)
		result.Duration = time.Duration(durationMS) * time.Millisecond
		result.ExecutedAt = time.Unix(executedAt, 0)

		results = append(results, &result)
	}

	return results, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}

// Helper functions
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	return i == 1
}

// ExportToJSON exports all hooks to JSON
func (r *Repository) ExportToJSON() (string, error) {
	hooks, err := r.List()
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(hooks, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal hooks: %w", err)
	}

	return string(data), nil
}
