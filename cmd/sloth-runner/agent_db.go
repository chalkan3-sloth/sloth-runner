package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// AgentDB manages the SQLite database for agents
type AgentDB struct {
	db *sql.DB
}

// AgentRecord represents an agent record in the database
type AgentRecord struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Address           string `json:"address"`
	Status            string `json:"status"`
	LastHeartbeat     int64  `json:"last_heartbeat"`
	RegisteredAt      int64  `json:"registered_at"`
	UpdatedAt         int64  `json:"updated_at"`
	LastInfoCollected int64  `json:"last_info_collected"`
	SystemInfo        string `json:"system_info"`
	Version           string `json:"version"`
}

// NewAgentDB creates a new AgentDB instance
func NewAgentDB(dbPath string) (*AgentDB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	agentDB := &AgentDB{db: db}

	// Initialize database schema
	if err := agentDB.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return agentDB, nil
}

// initSchema creates the agents table if it doesn't exist
func (adb *AgentDB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		address TEXT NOT NULL,
		status TEXT DEFAULT 'Active',
		last_heartbeat INTEGER DEFAULT 0,
		registered_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		last_info_collected INTEGER DEFAULT 0,
		system_info TEXT DEFAULT ''
	);

	CREATE INDEX IF NOT EXISTS idx_agents_name ON agents(name);
	CREATE INDEX IF NOT EXISTS idx_agents_status ON agents(status);
	CREATE INDEX IF NOT EXISTS idx_agents_last_heartbeat ON agents(last_heartbeat);
	`

	_, err := adb.db.Exec(schema)
	if err != nil {
		return err
	}

	// Create metrics history table
	metricsSchema := `
	CREATE TABLE IF NOT EXISTS agent_metrics_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		agent_name TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		cpu_percent REAL,
		memory_percent REAL,
		disk_percent REAL,
		load_avg_1min REAL,
		load_avg_5min REAL,
		load_avg_15min REAL,
		FOREIGN KEY (agent_name) REFERENCES agents(name)
	);

	CREATE INDEX IF NOT EXISTS idx_metrics_agent_time ON agent_metrics_history(agent_name, timestamp DESC);
	`

	_, err = adb.db.Exec(metricsSchema)
	if err != nil {
		return fmt.Errorf("failed to create metrics history table: %w", err)
	}

	// Migration: Add new columns if they don't exist (for existing databases)
	migrations := []string{
		`ALTER TABLE agents ADD COLUMN last_info_collected INTEGER DEFAULT 0`,
		`ALTER TABLE agents ADD COLUMN system_info TEXT DEFAULT ''`,
		`ALTER TABLE agents ADD COLUMN version TEXT DEFAULT ''`,
	}

	for _, migration := range migrations {
		// Ignore errors if column already exists
		adb.db.Exec(migration)
	}

	return nil
}

// RegisterAgent registers or updates an agent in the database
func (adb *AgentDB) RegisterAgent(name, address string) error {
	now := time.Now().Unix()

	// Use INSERT OR REPLACE to handle both new registrations and updates
	query := `
	INSERT OR REPLACE INTO agents (name, address, status, last_heartbeat, registered_at, updated_at)
	VALUES (?, ?, 'Active', ?, 
		COALESCE((SELECT registered_at FROM agents WHERE name = ?), ?), ?)
	`

	_, err := adb.db.Exec(query, name, address, now, name, now, now)
	if err != nil {
		return fmt.Errorf("failed to register agent: %w", err)
	}

	return nil
}

// UpdateHeartbeat updates the last heartbeat timestamp for an agent
func (adb *AgentDB) UpdateHeartbeat(name string) error {
	query := `UPDATE agents SET last_heartbeat = ?, updated_at = ? WHERE name = ?`

	now := time.Now().Unix()
	result, err := adb.db.Exec(query, now, now, name)
	if err != nil {
		return fmt.Errorf("failed to update heartbeat: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("agent not found: %s", name)
	}

	return nil
}

// UpdateSystemInfo updates the system information for an agent
func (adb *AgentDB) UpdateSystemInfo(name string, systemInfo string) error {
	query := `UPDATE agents SET system_info = ?, last_info_collected = ?, updated_at = ? WHERE name = ?`

	now := time.Now().Unix()
	result, err := adb.db.Exec(query, systemInfo, now, now, name)
	if err != nil {
		return fmt.Errorf("failed to update system info: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("agent not found: %s", name)
	}

	return nil
}

// UpdateVersion updates the version for an agent
func (adb *AgentDB) UpdateVersion(name string, version string) error {
	query := `UPDATE agents SET version = ?, updated_at = ? WHERE name = ?`

	now := time.Now().Unix()
	result, err := adb.db.Exec(query, version, now, name)
	if err != nil {
		return fmt.Errorf("failed to update version: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("agent not found: %s", name)
	}

	return nil
}

// GetAgent retrieves an agent by name
func (adb *AgentDB) GetAgent(name string) (*AgentRecord, error) {
	query := `SELECT id, name, address, status, last_heartbeat, registered_at, updated_at,
			  last_info_collected, system_info, version
			  FROM agents WHERE name = ?`

	var agent AgentRecord
	err := adb.db.QueryRow(query, name).Scan(
		&agent.ID,
		&agent.Name,
		&agent.Address,
		&agent.Status,
		&agent.LastHeartbeat,
		&agent.RegisteredAt,
		&agent.UpdatedAt,
		&agent.LastInfoCollected,
		&agent.SystemInfo,
		&agent.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agent not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}

	return &agent, nil
}

// ListAgents retrieves all agents from the database
func (adb *AgentDB) ListAgents() ([]*AgentRecord, error) {
	query := `SELECT id, name, address, status, last_heartbeat, registered_at, updated_at,
			  last_info_collected, system_info, version
			  FROM agents ORDER BY name`

	rows, err := adb.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query agents: %w", err)
	}
	defer rows.Close()

	var agents []*AgentRecord
	now := time.Now().Unix()

	for rows.Next() {
		var agent AgentRecord
		err := rows.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Address,
			&agent.Status,
			&agent.LastHeartbeat,
			&agent.RegisteredAt,
			&agent.UpdatedAt,
			&agent.LastInfoCollected,
			&agent.SystemInfo,
			&agent.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agent row: %w", err)
		}

		// Update status based on heartbeat (agent is active if heartbeat within last 60 seconds)
		if agent.LastHeartbeat > 0 && now-agent.LastHeartbeat < 60 {
			agent.Status = "Active"
		} else {
			agent.Status = "Inactive"
		}

		agents = append(agents, &agent)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating agent rows: %w", err)
	}

	return agents, nil
}

// GetAgentAddress retrieves the address of an agent by name
func (adb *AgentDB) GetAgentAddress(name string) (string, error) {
	query := `SELECT address FROM agents WHERE name = ? AND last_heartbeat > ?`

	// Only consider agents that have sent a heartbeat in the last 60 seconds
	cutoff := time.Now().Unix() - 60

	var address string
	err := adb.db.QueryRow(query, name, cutoff).Scan(&address)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("active agent not found: %s", name)
		}
		return "", fmt.Errorf("failed to get agent address: %w", err)
	}

	return address, nil
}

// RemoveAgent removes an agent from the database
func (adb *AgentDB) RemoveAgent(name string) error {
	query := `DELETE FROM agents WHERE name = ?`

	result, err := adb.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to remove agent: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("agent not found: %s", name)
	}

	return nil
}

// UnregisterAgent removes an agent from the database (alias for RemoveAgent)
func (adb *AgentDB) UnregisterAgent(name string) error {
	return adb.RemoveAgent(name)
}

// CleanupInactiveAgents removes agents that haven't sent heartbeat for a specified duration
func (adb *AgentDB) CleanupInactiveAgents(maxInactiveHours int) (int, error) {
	cutoff := time.Now().Unix() - int64(maxInactiveHours*3600)

	query := `DELETE FROM agents WHERE last_heartbeat < ? AND last_heartbeat > 0`

	result, err := adb.db.Exec(query, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup inactive agents: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to check affected rows: %w", err)
	}

	return int(rowsAffected), nil
}

// Close closes the database connection
func (adb *AgentDB) Close() error {
	if adb.db != nil {
		return adb.db.Close()
	}
	return nil
}

// GetStats returns database statistics
func (adb *AgentDB) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total agents
	var totalAgents int
	err := adb.db.QueryRow("SELECT COUNT(*) FROM agents").Scan(&totalAgents)
	if err != nil {
		return nil, fmt.Errorf("failed to get total agents: %w", err)
	}
	stats["total_agents"] = totalAgents

	// Active agents (heartbeat within last 60 seconds)
	cutoff := time.Now().Unix() - 60
	var activeAgents int
	err = adb.db.QueryRow("SELECT COUNT(*) FROM agents WHERE last_heartbeat > ?", cutoff).Scan(&activeAgents)
	if err != nil {
		return nil, fmt.Errorf("failed to get active agents: %w", err)
	}
	stats["active_agents"] = activeAgents

	// Inactive agents
	stats["inactive_agents"] = totalAgents - activeAgents

	return stats, nil
}

// SaveMetrics saves agent metrics to history
func (adb *AgentDB) SaveMetrics(agentName string, cpuPercent, memoryPercent, diskPercent, loadAvg1, loadAvg5, loadAvg15 float64) error {
	query := `INSERT INTO agent_metrics_history
		(agent_name, timestamp, cpu_percent, memory_percent, disk_percent, load_avg_1min, load_avg_5min, load_avg_15min)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := adb.db.Exec(query, agentName, time.Now().Unix(), cpuPercent, memoryPercent, diskPercent, loadAvg1, loadAvg5, loadAvg15)
	if err != nil {
		return fmt.Errorf("failed to save metrics: %w", err)
	}

	return nil
}

// GetMetricsHistory retrieves metrics history for an agent
func (adb *AgentDB) GetMetricsHistory(agentName string, limit int) ([]map[string]interface{}, error) {
	query := `SELECT timestamp, cpu_percent, memory_percent, disk_percent, load_avg_1min, load_avg_5min, load_avg_15min
		FROM agent_metrics_history
		WHERE agent_name = ?
		ORDER BY timestamp DESC
		LIMIT ?`

	rows, err := adb.db.Query(query, agentName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics history: %w", err)
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var timestamp int64
		var cpu, memory, disk, load1, load5, load15 float64

		err := rows.Scan(&timestamp, &cpu, &memory, &disk, &load1, &load5, &load15)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metrics row: %w", err)
		}

		history = append(history, map[string]interface{}{
			"timestamp":      timestamp,
			"cpu_percent":    cpu,
			"memory_percent": memory,
			"disk_percent":   disk,
			"load_avg_1min":  load1,
			"load_avg_5min":  load5,
			"load_avg_15min": load15,
		})
	}

	return history, nil
}

// CleanupOldMetrics removes metrics older than specified days
func (adb *AgentDB) CleanupOldMetrics(daysToKeep int) (int, error) {
	cutoff := time.Now().Unix() - int64(daysToKeep*24*3600)

	query := `DELETE FROM agent_metrics_history WHERE timestamp < ?`

	result, err := adb.db.Exec(query, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old metrics: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to check affected rows: %w", err)
	}

	return int(rowsAffected), nil
}
