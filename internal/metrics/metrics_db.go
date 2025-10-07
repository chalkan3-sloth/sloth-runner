package metrics

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// MetricPoint represents a single metric data point
type MetricPoint struct {
	Timestamp       int64   `json:"timestamp"`
	CPUPercent      float64 `json:"cpu_percent"`
	MemoryPercent   float64 `json:"memory_percent"`
	MemoryUsedBytes uint64  `json:"memory_used_bytes"`
	DiskPercent     float64 `json:"disk_percent"`
	LoadAvg1Min     float64 `json:"load_avg_1min"`
	LoadAvg5Min     float64 `json:"load_avg_5min"`
	LoadAvg15Min    float64 `json:"load_avg_15min"`
	ProcessCount    int     `json:"process_count"`
}

// MetricsDB handles storage and retrieval of historical metrics
type MetricsDB struct {
	db *sql.DB
}

// NewMetricsDB creates a new metrics database
func NewMetricsDB(dbPath string) (*MetricsDB, error) {
	// Log database path for debugging
	fmt.Printf("📊 Opening metrics database at: %s\n", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metrics database: %w", err)
	}

	// Create schema
	schema := `
	CREATE TABLE IF NOT EXISTS agent_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		agent_name TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		cpu_percent REAL,
		memory_percent REAL,
		memory_used_bytes INTEGER,
		disk_percent REAL,
		load_avg_1min REAL,
		load_avg_5min REAL,
		load_avg_15min REAL,
		process_count INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_agent_timestamp ON agent_metrics(agent_name, timestamp);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON agent_metrics(timestamp);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &MetricsDB{db: db}, nil
}

// StoreMetric stores a metric point for an agent
func (m *MetricsDB) StoreMetric(ctx context.Context, agentName string, metric MetricPoint) error {
	query := `
		INSERT INTO agent_metrics (
			agent_name, timestamp, cpu_percent, memory_percent, memory_used_bytes,
			disk_percent, load_avg_1min, load_avg_5min, load_avg_15min, process_count
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.db.ExecContext(ctx, query,
		agentName,
		metric.Timestamp,
		metric.CPUPercent,
		metric.MemoryPercent,
		metric.MemoryUsedBytes,
		metric.DiskPercent,
		metric.LoadAvg1Min,
		metric.LoadAvg5Min,
		metric.LoadAvg15Min,
		metric.ProcessCount,
	)

	return err
}

// GetMetricsHistory returns metrics for an agent within a time range
func (m *MetricsDB) GetMetricsHistory(ctx context.Context, agentName string, startTime, endTime int64, maxPoints int) ([]MetricPoint, error) {
	query := `
		SELECT timestamp, cpu_percent, memory_percent, memory_used_bytes, disk_percent,
		       load_avg_1min, load_avg_5min, load_avg_15min, process_count
		FROM agent_metrics
		WHERE agent_name = ? AND timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`

	if maxPoints > 0 {
		// Calculate interval to downsample if too many points
		countQuery := `SELECT COUNT(*) FROM agent_metrics WHERE agent_name = ? AND timestamp BETWEEN ? AND ?`
		var totalPoints int
		err := m.db.QueryRowContext(ctx, countQuery, agentName, startTime, endTime).Scan(&totalPoints)
		if err != nil {
			return nil, err
		}

		if totalPoints > maxPoints {
			// Downsample by taking every Nth point
			interval := totalPoints / maxPoints
			query = fmt.Sprintf(`
				SELECT timestamp, cpu_percent, memory_percent, memory_used_bytes, disk_percent,
				       load_avg_1min, load_avg_5min, load_avg_15min, process_count
				FROM (
					SELECT *, ROW_NUMBER() OVER (ORDER BY timestamp) as rn
					FROM agent_metrics
					WHERE agent_name = ? AND timestamp BETWEEN ? AND ?
				)
				WHERE rn %% %d = 0
				ORDER BY timestamp ASC
			`, interval)
		}
	}

	rows, err := m.db.QueryContext(ctx, query, agentName, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []MetricPoint
	for rows.Next() {
		var mp MetricPoint
		err := rows.Scan(
			&mp.Timestamp,
			&mp.CPUPercent,
			&mp.MemoryPercent,
			&mp.MemoryUsedBytes,
			&mp.DiskPercent,
			&mp.LoadAvg1Min,
			&mp.LoadAvg5Min,
			&mp.LoadAvg15Min,
			&mp.ProcessCount,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, mp)
	}

	return metrics, rows.Err()
}

// GetLatestMetric returns the most recent metric for an agent
func (m *MetricsDB) GetLatestMetric(ctx context.Context, agentName string) (*MetricPoint, error) {
	query := `
		SELECT timestamp, cpu_percent, memory_percent, memory_used_bytes, disk_percent,
		       load_avg_1min, load_avg_5min, load_avg_15min, process_count
		FROM agent_metrics
		WHERE agent_name = ?
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var mp MetricPoint
	err := m.db.QueryRowContext(ctx, query, agentName).Scan(
		&mp.Timestamp,
		&mp.CPUPercent,
		&mp.MemoryPercent,
		&mp.MemoryUsedBytes,
		&mp.DiskPercent,
		&mp.LoadAvg1Min,
		&mp.LoadAvg5Min,
		&mp.LoadAvg15Min,
		&mp.ProcessCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &mp, nil
}

// CleanupOldMetrics deletes metrics older than the specified duration
func (m *MetricsDB) CleanupOldMetrics(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan).Unix()
	_, err := m.db.ExecContext(ctx, "DELETE FROM agent_metrics WHERE timestamp < ?", cutoff)
	return err
}

// GetAgentNames returns all agent names that have metrics
func (m *MetricsDB) GetAgentNames(ctx context.Context) ([]string, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT DISTINCT agent_name FROM agent_metrics ORDER BY agent_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, rows.Err()
}

// Close closes the database connection
func (m *MetricsDB) Close() error {
	return m.db.Close()
}
