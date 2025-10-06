package hooks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventQueue manages the persistent event queue
type EventQueue struct {
	db *sql.DB
}

// NewEventQueue creates a new event queue
func NewEventQueue(db *sql.DB) *EventQueue {
	return &EventQueue{db: db}
}

// InitializeSchema creates the events table
func (eq *EventQueue) InitializeSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		data TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		error TEXT,
		timestamp INTEGER NOT NULL,
		created_at INTEGER NOT NULL,
		processed_at INTEGER
	);

	CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);
	CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
	CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);

	CREATE TABLE IF NOT EXISTS file_watchers (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		path TEXT NOT NULL,
		pattern TEXT,
		events TEXT NOT NULL,
		recursive INTEGER DEFAULT 0,
		enabled INTEGER DEFAULT 1,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS event_hook_executions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id TEXT NOT NULL,
		hook_id TEXT NOT NULL,
		hook_name TEXT NOT NULL,
		success INTEGER NOT NULL,
		output TEXT,
		error TEXT,
		duration_ms INTEGER NOT NULL,
		executed_at INTEGER NOT NULL,
		FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
		FOREIGN KEY (hook_id) REFERENCES hooks(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_event_hook_executions_event_id ON event_hook_executions(event_id);
	CREATE INDEX IF NOT EXISTS idx_event_hook_executions_hook_id ON event_hook_executions(hook_id);
	`

	_, err := eq.db.Exec(schema)
	return err
}

// EnqueueEvent adds an event to the queue
func (eq *EventQueue) EnqueueEvent(event *Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Status == "" {
		event.Status = EventStatusPending
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	dataJSON, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	query := `
		INSERT INTO events (id, type, data, status, timestamp, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = eq.db.Exec(query,
		event.ID,
		event.Type,
		string(dataJSON),
		event.Status,
		event.Timestamp.Unix(),
		event.CreatedAt.Unix(),
	)

	return err
}

// GetPendingEvents returns all pending events ordered by creation time
func (eq *EventQueue) GetPendingEvents(limit int) ([]*Event, error) {
	query := `
		SELECT id, type, data, status, error, timestamp, created_at, processed_at
		FROM events
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT ?
	`

	rows, err := eq.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		event, err := eq.scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

// ListEvents returns events with optional filters
func (eq *EventQueue) ListEvents(eventType EventType, status EventStatus, limit int) ([]*Event, error) {
	query := `
		SELECT id, type, data, status, error, timestamp, created_at, processed_at
		FROM events
		WHERE 1=1
	`
	args := []interface{}{}

	if eventType != "" {
		query += " AND type = ?"
		args = append(args, eventType)
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := eq.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		event, err := eq.scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

// GetEvent retrieves a specific event by ID
func (eq *EventQueue) GetEvent(id string) (*Event, error) {
	query := `
		SELECT id, type, data, status, error, timestamp, created_at, processed_at
		FROM events
		WHERE id = ?
	`

	row := eq.db.QueryRow(query, id)
	return eq.scanEvent(row)
}

// UpdateEventStatus updates the status of an event
func (eq *EventQueue) UpdateEventStatus(id string, status EventStatus, errorMsg string) error {
	now := time.Now()
	query := `
		UPDATE events
		SET status = ?, error = ?, processed_at = ?
		WHERE id = ?
	`

	_, err := eq.db.Exec(query, status, errorMsg, now.Unix(), id)
	return err
}

// DeleteEvent removes an event from the queue
func (eq *EventQueue) DeleteEvent(id string) error {
	_, err := eq.db.Exec("DELETE FROM events WHERE id = ?", id)
	return err
}

// CleanupOldEvents removes completed/failed events older than specified duration
func (eq *EventQueue) CleanupOldEvents(olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan).Unix()
	query := `
		DELETE FROM events
		WHERE status IN ('completed', 'failed')
		AND processed_at < ?
	`

	result, err := eq.db.Exec(query, cutoff)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// scanEvent scans a database row into an Event struct
func (eq *EventQueue) scanEvent(scanner interface {
	Scan(dest ...interface{}) error
}) (*Event, error) {
	var event Event
	var dataJSON string
	var timestamp, createdAt int64
	var processedAt sql.NullInt64
	var errorMsg sql.NullString

	err := scanner.Scan(
		&event.ID,
		&event.Type,
		&dataJSON,
		&event.Status,
		&errorMsg,
		&timestamp,
		&createdAt,
		&processedAt,
	)
	if err != nil {
		return nil, err
	}

	// Parse data JSON
	if err := json.Unmarshal([]byte(dataJSON), &event.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	event.Timestamp = time.Unix(timestamp, 0)
	event.CreatedAt = time.Unix(createdAt, 0)

	if processedAt.Valid {
		t := time.Unix(processedAt.Int64, 0)
		event.ProcessedAt = &t
	}

	if errorMsg.Valid {
		event.Error = errorMsg.String
	}

	return &event, nil
}

// AddFileWatcher adds a new file watcher
func (eq *EventQueue) AddFileWatcher(watcher *FileWatcher) error {
	if watcher.ID == "" {
		watcher.ID = uuid.New().String()
	}
	now := time.Now()
	watcher.CreatedAt = now
	watcher.UpdatedAt = now

	eventsJSON, err := json.Marshal(watcher.Events)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO file_watchers (id, name, path, pattern, events, recursive, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = eq.db.Exec(query,
		watcher.ID,
		watcher.Name,
		watcher.Path,
		watcher.Pattern,
		string(eventsJSON),
		watcher.Recursive,
		watcher.Enabled,
		watcher.CreatedAt.Unix(),
		watcher.UpdatedAt.Unix(),
	)

	return err
}

// ListFileWatchers returns all file watchers
func (eq *EventQueue) ListFileWatchers() ([]*FileWatcher, error) {
	query := `
		SELECT id, name, path, pattern, events, recursive, enabled, created_at, updated_at
		FROM file_watchers
		ORDER BY name
	`

	rows, err := eq.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var watchers []*FileWatcher
	for rows.Next() {
		watcher, err := eq.scanFileWatcher(rows)
		if err != nil {
			return nil, err
		}
		watchers = append(watchers, watcher)
	}

	return watchers, rows.Err()
}

// GetFileWatcher retrieves a specific file watcher
func (eq *EventQueue) GetFileWatcher(name string) (*FileWatcher, error) {
	query := `
		SELECT id, name, path, pattern, events, recursive, enabled, created_at, updated_at
		FROM file_watchers
		WHERE name = ?
	`

	row := eq.db.QueryRow(query, name)
	return eq.scanFileWatcher(row)
}

// DeleteFileWatcher removes a file watcher
func (eq *EventQueue) DeleteFileWatcher(name string) error {
	_, err := eq.db.Exec("DELETE FROM file_watchers WHERE name = ?", name)
	return err
}

// EnableFileWatcher enables a file watcher
func (eq *EventQueue) EnableFileWatcher(name string) error {
	_, err := eq.db.Exec("UPDATE file_watchers SET enabled = 1, updated_at = ? WHERE name = ?",
		time.Now().Unix(), name)
	return err
}

// DisableFileWatcher disables a file watcher
func (eq *EventQueue) DisableFileWatcher(name string) error {
	_, err := eq.db.Exec("UPDATE file_watchers SET enabled = 0, updated_at = ? WHERE name = ?",
		time.Now().Unix(), name)
	return err
}

// scanFileWatcher scans a database row into a FileWatcher struct
func (eq *EventQueue) scanFileWatcher(scanner interface {
	Scan(dest ...interface{}) error
}) (*FileWatcher, error) {
	var watcher FileWatcher
	var eventsJSON string
	var createdAt, updatedAt int64
	var recursive, enabled int

	err := scanner.Scan(
		&watcher.ID,
		&watcher.Name,
		&watcher.Path,
		&watcher.Pattern,
		&eventsJSON,
		&recursive,
		&enabled,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(eventsJSON), &watcher.Events); err != nil {
		return nil, err
	}

	watcher.Recursive = recursive == 1
	watcher.Enabled = enabled == 1
	watcher.CreatedAt = time.Unix(createdAt, 0)
	watcher.UpdatedAt = time.Unix(updatedAt, 0)

	return &watcher, nil
}

// RecordEventHookExecution records a hook execution for an event
func (eq *EventQueue) RecordEventHookExecution(eventID string, hookID string, hookName string, result *HookResult) error {
	query := `
		INSERT INTO event_hook_executions (event_id, hook_id, hook_name, success, output, error, duration_ms, executed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	success := 0
	if result.Success {
		success = 1
	}

	_, err := eq.db.Exec(query,
		eventID,
		hookID,
		hookName,
		success,
		result.Output,
		result.Error,
		result.Duration.Milliseconds(),
		result.ExecutedAt.Unix(),
	)

	return err
}

// GetEventHookExecutions retrieves all hook executions for an event
func (eq *EventQueue) GetEventHookExecutions(eventID string) ([]*EventHookExecution, error) {
	query := `
		SELECT id, event_id, hook_id, hook_name, success, output, error, duration_ms, executed_at
		FROM event_hook_executions
		WHERE event_id = ?
		ORDER BY executed_at ASC
	`

	rows, err := eq.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var executions []*EventHookExecution
	for rows.Next() {
		var exec EventHookExecution
		var success int
		var durationMS int64
		var executedAt int64
		var output, errorMsg sql.NullString

		err := rows.Scan(
			&exec.ID,
			&exec.EventID,
			&exec.HookID,
			&exec.HookName,
			&success,
			&output,
			&errorMsg,
			&durationMS,
			&executedAt,
		)
		if err != nil {
			return nil, err
		}

		exec.Success = success == 1
		exec.Duration = time.Duration(durationMS) * time.Millisecond
		exec.ExecutedAt = time.Unix(executedAt, 0)
		if output.Valid {
			exec.Output = output.String
		}
		if errorMsg.Valid {
			exec.Error = errorMsg.String
		}

		executions = append(executions, &exec)
	}

	return executions, rows.Err()
}
