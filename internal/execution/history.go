package execution

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ExecutionStatus represents the status of an execution
type ExecutionStatus string

const (
	StatusRunning   ExecutionStatus = "running"
	StatusCompleted ExecutionStatus = "completed"
	StatusFailed    ExecutionStatus = "failed"
	StatusCancelled ExecutionStatus = "cancelled"
)

// Execution represents a workflow execution record
type Execution struct {
	ID           string                 `json:"id"`
	WorkflowName string                 `json:"workflow_name"`
	WorkflowFile string                 `json:"workflow_file"`
	GroupName    string                 `json:"group_name,omitempty"`
	Status       ExecutionStatus        `json:"status"`
	StartTime    int64                  `json:"start_time"`
	EndTime      int64                  `json:"end_time,omitempty"`
	Duration     int64                  `json:"duration,omitempty"` // milliseconds
	AgentName    string                 `json:"agent_name,omitempty"`
	User         string                 `json:"user,omitempty"`
	ExitCode     int                    `json:"exit_code"`
	Output       string                 `json:"output,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	TasksTotal   int                    `json:"tasks_total"`
	TasksSuccess int                    `json:"tasks_success"`
	TasksFailed  int                    `json:"tasks_failed"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// TaskExecution represents a single task execution within a workflow
type TaskExecution struct {
	ID          string          `json:"id"`
	ExecutionID string          `json:"execution_id"`
	TaskName    string          `json:"task_name"`
	Status      ExecutionStatus `json:"status"`
	StartTime   int64           `json:"start_time"`
	EndTime     int64           `json:"end_time,omitempty"`
	Duration    int64           `json:"duration,omitempty"`
	Output      string          `json:"output,omitempty"`
	Error       string          `json:"error,omitempty"`
	Changed     bool            `json:"changed"`
}

// HistoryDB manages execution history in SQLite
type HistoryDB struct {
	db *sql.DB
}

// NewHistoryDB creates a new history database connection
func NewHistoryDB(dbPath string) (*HistoryDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	h := &HistoryDB{db: db}
	if err := h.initialize(); err != nil {
		return nil, err
	}

	return h, nil
}

// initialize creates the necessary tables
func (h *HistoryDB) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS executions (
		id TEXT PRIMARY KEY,
		workflow_name TEXT NOT NULL,
		workflow_file TEXT NOT NULL,
		group_name TEXT,
		status TEXT NOT NULL,
		start_time INTEGER NOT NULL,
		end_time INTEGER,
		duration INTEGER,
		agent_name TEXT,
		user TEXT,
		exit_code INTEGER DEFAULT 0,
		output TEXT,
		error_message TEXT,
		tasks_total INTEGER DEFAULT 0,
		tasks_success INTEGER DEFAULT 0,
		tasks_failed INTEGER DEFAULT 0,
		metadata TEXT,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);

	CREATE TABLE IF NOT EXISTS task_executions (
		id TEXT PRIMARY KEY,
		execution_id TEXT NOT NULL,
		task_name TEXT NOT NULL,
		status TEXT NOT NULL,
		start_time INTEGER NOT NULL,
		end_time INTEGER,
		duration INTEGER,
		output TEXT,
		error TEXT,
		changed INTEGER DEFAULT 0,
		created_at INTEGER DEFAULT (strftime('%s', 'now')),
		FOREIGN KEY (execution_id) REFERENCES executions(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_executions_workflow ON executions(workflow_name);
	CREATE INDEX IF NOT EXISTS idx_executions_status ON executions(status);
	CREATE INDEX IF NOT EXISTS idx_executions_start_time ON executions(start_time DESC);
	CREATE INDEX IF NOT EXISTS idx_executions_agent ON executions(agent_name);
	CREATE INDEX IF NOT EXISTS idx_executions_group ON executions(group_name);
	CREATE INDEX IF NOT EXISTS idx_task_executions_execution ON task_executions(execution_id);
	`

	_, err := h.db.Exec(schema)
	return err
}

// CreateExecution creates a new execution record
func (h *HistoryDB) CreateExecution(exec *Execution) error {
	metadataJSON, _ := json.Marshal(exec.Metadata)

	query := `
		INSERT INTO executions (
			id, workflow_name, workflow_file, group_name, status, start_time,
			agent_name, user, tasks_total, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := h.db.Exec(query,
		exec.ID, exec.WorkflowName, exec.WorkflowFile, exec.GroupName,
		exec.Status, exec.StartTime, exec.AgentName, exec.User,
		exec.TasksTotal, string(metadataJSON),
	)

	return err
}

// UpdateExecution updates an existing execution record
func (h *HistoryDB) UpdateExecution(exec *Execution) error {
	query := `
		UPDATE executions SET
			status = ?,
			end_time = ?,
			duration = ?,
			exit_code = ?,
			output = ?,
			error_message = ?,
			tasks_success = ?,
			tasks_failed = ?
		WHERE id = ?
	`

	_, err := h.db.Exec(query,
		exec.Status, exec.EndTime, exec.Duration, exec.ExitCode,
		exec.Output, exec.ErrorMessage, exec.TasksSuccess, exec.TasksFailed,
		exec.ID,
	)

	return err
}

// CreateTaskExecution creates a new task execution record
func (h *HistoryDB) CreateTaskExecution(task *TaskExecution) error {
	query := `
		INSERT INTO task_executions (
			id, execution_id, task_name, status, start_time, changed
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	changed := 0
	if task.Changed {
		changed = 1
	}

	_, err := h.db.Exec(query,
		task.ID, task.ExecutionID, task.TaskName, task.Status,
		task.StartTime, changed,
	)

	return err
}

// UpdateTaskExecution updates a task execution record
func (h *HistoryDB) UpdateTaskExecution(task *TaskExecution) error {
	query := `
		UPDATE task_executions SET
			status = ?,
			end_time = ?,
			duration = ?,
			output = ?,
			error = ?,
			changed = ?
		WHERE id = ?
	`

	changed := 0
	if task.Changed {
		changed = 1
	}

	_, err := h.db.Exec(query,
		task.Status, task.EndTime, task.Duration,
		task.Output, task.Error, changed, task.ID,
	)

	return err
}

// GetExecution retrieves an execution by ID
func (h *HistoryDB) GetExecution(id string) (*Execution, error) {
	query := `
		SELECT id, workflow_name, workflow_file, group_name, status,
			start_time, end_time, duration, agent_name, user, exit_code,
			output, error_message, tasks_total, tasks_success, tasks_failed, metadata
		FROM executions WHERE id = ?
	`

	var exec Execution
	var metadataJSON sql.NullString
	var groupName, agentName, user, output, errorMsg sql.NullString
	var endTime, duration sql.NullInt64

	err := h.db.QueryRow(query, id).Scan(
		&exec.ID, &exec.WorkflowName, &exec.WorkflowFile, &groupName, &exec.Status,
		&exec.StartTime, &endTime, &duration, &agentName, &user, &exec.ExitCode,
		&output, &errorMsg, &exec.TasksTotal, &exec.TasksSuccess, &exec.TasksFailed,
		&metadataJSON,
	)

	if err != nil {
		return nil, err
	}

	if groupName.Valid {
		exec.GroupName = groupName.String
	}
	if agentName.Valid {
		exec.AgentName = agentName.String
	}
	if user.Valid {
		exec.User = user.String
	}
	if output.Valid {
		exec.Output = output.String
	}
	if errorMsg.Valid {
		exec.ErrorMessage = errorMsg.String
	}
	if endTime.Valid {
		exec.EndTime = endTime.Int64
	}
	if duration.Valid {
		exec.Duration = duration.Int64
	}
	if metadataJSON.Valid {
		json.Unmarshal([]byte(metadataJSON.String), &exec.Metadata)
	}

	return &exec, nil
}

// ListExecutions retrieves executions with filters
func (h *HistoryDB) ListExecutions(filters map[string]interface{}, limit, offset int) ([]*Execution, error) {
	query := `
		SELECT id, workflow_name, workflow_file, group_name, status,
			start_time, end_time, duration, agent_name, user, exit_code,
			tasks_total, tasks_success, tasks_failed
		FROM executions
		WHERE 1=1
	`
	args := []interface{}{}

	if workflow, ok := filters["workflow"]; ok {
		query += " AND workflow_name = ?"
		args = append(args, workflow)
	}
	if status, ok := filters["status"]; ok {
		query += " AND status = ?"
		args = append(args, status)
	}
	if agent, ok := filters["agent"]; ok {
		query += " AND agent_name = ?"
		args = append(args, agent)
	}
	if group, ok := filters["group"]; ok {
		query += " AND group_name = ?"
		args = append(args, group)
	}
	if since, ok := filters["since"]; ok {
		query += " AND start_time >= ?"
		args = append(args, since)
	}

	query += " ORDER BY start_time DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var executions []*Execution
	for rows.Next() {
		var exec Execution
		var groupName, agentName, user sql.NullString
		var endTime, duration sql.NullInt64

		err := rows.Scan(
			&exec.ID, &exec.WorkflowName, &exec.WorkflowFile, &groupName, &exec.Status,
			&exec.StartTime, &endTime, &duration, &agentName, &user, &exec.ExitCode,
			&exec.TasksTotal, &exec.TasksSuccess, &exec.TasksFailed,
		)
		if err != nil {
			continue
		}

		if groupName.Valid {
			exec.GroupName = groupName.String
		}
		if agentName.Valid {
			exec.AgentName = agentName.String
		}
		if user.Valid {
			exec.User = user.String
		}
		if endTime.Valid {
			exec.EndTime = endTime.Int64
		}
		if duration.Valid {
			exec.Duration = duration.Int64
		}

		executions = append(executions, &exec)
	}

	return executions, nil
}

// GetTaskExecutions retrieves all task executions for a given execution
func (h *HistoryDB) GetTaskExecutions(executionID string) ([]*TaskExecution, error) {
	query := `
		SELECT id, execution_id, task_name, status, start_time, end_time,
			duration, output, error, changed
		FROM task_executions
		WHERE execution_id = ?
		ORDER BY start_time ASC
	`

	rows, err := h.db.Query(query, executionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*TaskExecution
	for rows.Next() {
		var task TaskExecution
		var endTime, duration sql.NullInt64
		var output, errorStr sql.NullString
		var changed int

		err := rows.Scan(
			&task.ID, &task.ExecutionID, &task.TaskName, &task.Status,
			&task.StartTime, &endTime, &duration, &output, &errorStr, &changed,
		)
		if err != nil {
			continue
		}

		if endTime.Valid {
			task.EndTime = endTime.Int64
		}
		if duration.Valid {
			task.Duration = duration.Int64
		}
		if output.Valid {
			task.Output = output.String
		}
		if errorStr.Valid {
			task.Error = errorStr.String
		}
		task.Changed = changed == 1

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// GetStatistics returns execution statistics
func (h *HistoryDB) GetStatistics(filters map[string]interface{}) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END), 0) as completed,
			COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) as failed,
			COALESCE(SUM(CASE WHEN status = 'running' THEN 1 ELSE 0 END), 0) as running,
			AVG(CASE WHEN duration > 0 THEN duration ELSE NULL END) as avg_duration
		FROM executions
		WHERE 1=1
	`
	args := []interface{}{}

	if workflow, ok := filters["workflow"]; ok {
		query += " AND workflow_name = ?"
		args = append(args, workflow)
	}
	if since, ok := filters["since"]; ok {
		query += " AND start_time >= ?"
		args = append(args, since)
	}

	var total, completed, failed, running int
	var avgDuration sql.NullFloat64

	err := h.db.QueryRow(query, args...).Scan(
		&total, &completed, &failed, &running, &avgDuration,
	)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total":     total,
		"completed": completed,
		"failed":    failed,
		"running":   running,
	}

	if avgDuration.Valid {
		stats["avg_duration"] = int64(avgDuration.Float64)
	} else {
		stats["avg_duration"] = 0
	}

	if total > 0 {
		stats["success_rate"] = float64(completed) / float64(total) * 100
	} else {
		stats["success_rate"] = 0.0
	}

	return stats, nil
}

// DeleteOldExecutions deletes executions older than the specified days
func (h *HistoryDB) DeleteOldExecutions(days int) (int64, error) {
	cutoffTime := time.Now().AddDate(0, 0, -days).Unix()

	result, err := h.db.Exec("DELETE FROM executions WHERE start_time < ?", cutoffTime)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Close closes the database connection
func (h *HistoryDB) Close() error {
	return h.db.Close()
}
