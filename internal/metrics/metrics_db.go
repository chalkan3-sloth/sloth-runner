package metrics

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
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

// MetricsDB handles storage and retrieval of historical metrics with optimizations
type MetricsDB struct {
	db *sql.DB

	// Optimizations: batch writes
	batchMu     sync.Mutex
	batchBuffer []batchMetric
	batchSize   int
	flushTimer  *time.Timer
}

// batchMetric holds a metric waiting to be batch-inserted
type batchMetric struct {
	AgentName string
	Metric    MetricPoint
}

// NewMetricsDB creates a new metrics database with optimizations
func NewMetricsDB(dbPath string) (*MetricsDB, error) {
	// Log database path for debugging
	fmt.Printf("📊 Opening metrics database at: %s\n", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metrics database: %w", err)
	}

	// Optimizations: WAL mode, pragmas
	schema := `
	PRAGMA journal_mode = WAL;
	PRAGMA synchronous = NORMAL;
	PRAGMA cache_size = 10000;
	PRAGMA temp_store = MEMORY;

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

	metricsDB := &MetricsDB{
		db:          db,
		batchBuffer: make([]batchMetric, 0, 10),
		batchSize:   10, // Batch write every 10 metrics
	}

	// Start auto-flush timer (5 seconds)
	metricsDB.startAutoFlush()

	return metricsDB, nil
}

// startAutoFlush starts a timer to flush batch periodically
func (m *MetricsDB) startAutoFlush() {
	m.flushTimer = time.AfterFunc(5*time.Second, func() {
		m.flushBatch()
		m.startAutoFlush() // Restart timer
	})
}

// StoreMetric stores a metric point for an agent (optimized with batching)
func (m *MetricsDB) StoreMetric(ctx context.Context, agentName string, metric MetricPoint) error {
	m.batchMu.Lock()
	defer m.batchMu.Unlock()

	// Add to batch buffer
	m.batchBuffer = append(m.batchBuffer, batchMetric{
		AgentName: agentName,
		Metric:    metric,
	})

	// Flush if batch is full
	if len(m.batchBuffer) >= m.batchSize {
		return m.flushBatchLocked()
	}

	return nil
}

// flushBatch flushes the batch buffer to database
func (m *MetricsDB) flushBatch() error {
	m.batchMu.Lock()
	defer m.batchMu.Unlock()
	return m.flushBatchLocked()
}

// flushBatchLocked flushes with lock already held
func (m *MetricsDB) flushBatchLocked() error {
	if len(m.batchBuffer) == 0 {
		return nil
	}

	// Begin transaction for batch insert
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO agent_metrics (
			agent_name, timestamp, cpu_percent, memory_percent, memory_used_bytes,
			disk_percent, load_avg_1min, load_avg_5min, load_avg_15min, process_count
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert all batched metrics
	for _, bm := range m.batchBuffer {
		_, err := stmt.Exec(
			bm.AgentName,
			bm.Metric.Timestamp,
			bm.Metric.CPUPercent,
			bm.Metric.MemoryPercent,
			bm.Metric.MemoryUsedBytes,
			bm.Metric.DiskPercent,
			bm.Metric.LoadAvg1Min,
			bm.Metric.LoadAvg5Min,
			bm.Metric.LoadAvg15Min,
			bm.Metric.ProcessCount,
		)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// Clear buffer
	m.batchBuffer = m.batchBuffer[:0]
	return nil
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
	// Stop flush timer
	if m.flushTimer != nil {
		m.flushTimer.Stop()
	}

	// Flush any remaining metrics
	m.flushBatch()

	return m.db.Close()
}
