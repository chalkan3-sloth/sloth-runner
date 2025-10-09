package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// InitializeAgentGroupsSchema creates the necessary tables for agent groups
func (w *AgentDBWrapper) InitializeAgentGroupsSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS agent_groups (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		tags TEXT,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS agent_group_members (
		group_id TEXT NOT NULL,
		agent_name TEXT NOT NULL,
		added_at INTEGER NOT NULL,
		PRIMARY KEY (group_id, agent_name),
		FOREIGN KEY (group_id) REFERENCES agent_groups(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_agent_group_members_group_id ON agent_group_members(group_id);
	CREATE INDEX IF NOT EXISTS idx_agent_group_members_agent_name ON agent_group_members(agent_name);
	`

	_, err := w.db.Exec(schema)
	return err
}

// CreateAgentGroup creates a new agent group
func (w *AgentDBWrapper) CreateAgentGroup(ctx context.Context, group *AgentGroup) error {
	now := time.Now().Unix()

	// Marshal tags to JSON
	tagsJSON, err := json.Marshal(group.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	// Insert group
	query := `INSERT INTO agent_groups (id, name, description, tags, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`

	_, err = w.db.ExecContext(ctx, query, group.ID, group.Name, group.Description, string(tagsJSON), now, now)
	if err != nil {
		return fmt.Errorf("failed to create agent group: %w", err)
	}

	// Add agents to group
	if len(group.AgentNames) > 0 {
		if err := w.AddAgentsToGroup(ctx, group.ID, group.AgentNames); err != nil {
			return fmt.Errorf("failed to add agents to group: %w", err)
		}
	}

	return nil
}

// GetAgentGroup retrieves an agent group by ID
func (w *AgentDBWrapper) GetAgentGroup(ctx context.Context, groupID string) (*AgentGroup, error) {
	query := `SELECT id, name, description, tags, created_at, updated_at FROM agent_groups WHERE id = ?`

	var group AgentGroup
	var tagsJSON string

	err := w.db.QueryRowContext(ctx, query, groupID).Scan(
		&group.ID, &group.Name, &group.Description, &tagsJSON, &group.CreatedAt, &group.CreatedAt, // CreatedAt used twice for created/updated
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found: %s", groupID)
		}
		return nil, err
	}

	// Unmarshal tags
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &group.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}

	// Get agent members
	agentNames, err := w.GetGroupAgents(ctx, groupID)
	if err != nil {
		return nil, err
	}
	group.AgentNames = agentNames
	group.AgentCount = len(agentNames)

	return &group, nil
}

// ListAgentGroups retrieves all agent groups
func (w *AgentDBWrapper) ListAgentGroups(ctx context.Context) ([]*AgentGroup, error) {
	query := `SELECT id, name, description, tags, created_at, updated_at FROM agent_groups ORDER BY name`

	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*AgentGroup
	for rows.Next() {
		var group AgentGroup
		var tagsJSON string

		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &tagsJSON, &group.CreatedAt, &group.CreatedAt); err != nil {
			return nil, err
		}

		// Unmarshal tags
		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &group.Tags); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
			}
		}

		// Get agent count
		agentNames, err := w.GetGroupAgents(ctx, group.ID)
		if err != nil {
			return nil, err
		}
		group.AgentNames = agentNames
		group.AgentCount = len(agentNames)

		groups = append(groups, &group)
	}

	return groups, nil
}

// UpdateAgentGroup updates an agent group
func (w *AgentDBWrapper) UpdateAgentGroup(ctx context.Context, group *AgentGroup) error {
	now := time.Now().Unix()

	// Marshal tags to JSON
	tagsJSON, err := json.Marshal(group.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `UPDATE agent_groups SET name = ?, description = ?, tags = ?, updated_at = ? WHERE id = ?`

	result, err := w.db.ExecContext(ctx, query, group.Name, group.Description, string(tagsJSON), now, group.ID)
	if err != nil {
		return fmt.Errorf("failed to update agent group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group not found: %s", group.ID)
	}

	return nil
}

// DeleteAgentGroup deletes an agent group
func (w *AgentDBWrapper) DeleteAgentGroup(ctx context.Context, groupID string) error {
	// The CASCADE will automatically delete group members
	query := `DELETE FROM agent_groups WHERE id = ?`

	result, err := w.db.ExecContext(ctx, query, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete agent group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group not found: %s", groupID)
	}

	return nil
}

// AddAgentsToGroup adds agents to a group
func (w *AgentDBWrapper) AddAgentsToGroup(ctx context.Context, groupID string, agentNames []string) error {
	now := time.Now().Unix()

	for _, agentName := range agentNames {
		query := `INSERT OR IGNORE INTO agent_group_members (group_id, agent_name, added_at) VALUES (?, ?, ?)`
		_, err := w.db.ExecContext(ctx, query, groupID, agentName, now)
		if err != nil {
			return fmt.Errorf("failed to add agent %s to group: %w", agentName, err)
		}
	}

	return nil
}

// RemoveAgentsFromGroup removes agents from a group
func (w *AgentDBWrapper) RemoveAgentsFromGroup(ctx context.Context, groupID string, agentNames []string) error {
	for _, agentName := range agentNames {
		query := `DELETE FROM agent_group_members WHERE group_id = ? AND agent_name = ?`
		_, err := w.db.ExecContext(ctx, query, groupID, agentName)
		if err != nil {
			return fmt.Errorf("failed to remove agent %s from group: %w", agentName, err)
		}
	}

	return nil
}

// GetGroupAgents retrieves all agents in a group
func (w *AgentDBWrapper) GetGroupAgents(ctx context.Context, groupID string) ([]string, error) {
	query := `SELECT agent_name FROM agent_group_members WHERE group_id = ? ORDER BY agent_name`

	rows, err := w.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agentNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		agentNames = append(agentNames, name)
	}

	return agentNames, nil
}

// GetAgentGroupMembership returns all groups an agent belongs to
func (w *AgentDBWrapper) GetAgentGroupMembership(ctx context.Context, agentName string) ([]string, error) {
	query := `SELECT group_id FROM agent_group_members WHERE agent_name = ?`

	rows, err := w.db.QueryContext(ctx, query, agentName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupIDs []string
	for rows.Next() {
		var groupID string
		if err := rows.Scan(&groupID); err != nil {
			return nil, err
		}
		groupIDs = append(groupIDs, groupID)
	}

	return groupIDs, nil
}

// GetGroupAggregatedMetrics calculates aggregated metrics for all agents in a group
func (w *AgentDBWrapper) GetGroupAggregatedMetrics(ctx context.Context, groupID string) (map[string]interface{}, error) {
	// Get agents in the group
	agentNames, err := w.GetGroupAgents(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if len(agentNames) == 0 {
		return map[string]interface{}{
			"total_agents":       0,
			"healthy_agents":     0,
			"unhealthy_agents":   0,
			"avg_cpu_percent":    0.0,
			"avg_memory_percent": 0.0,
			"avg_disk_percent":   0.0,
		}, nil
	}

	// Build a query to get latest metrics for all agents in the group
	placeholders := make([]interface{}, len(agentNames))
	for i, name := range agentNames {
		placeholders[i] = name
	}

	// Query to get latest metrics for each agent
	// SQLite doesn't support DISTINCT ON, so we use a subquery with MAX
	query := `
		SELECT
			COUNT(*) as total_agents,
			AVG(cpu_percent) as avg_cpu,
			AVG(memory_percent) as avg_memory,
			AVG(disk_percent) as avg_disk
		FROM agent_metrics_history
		WHERE agent_name IN (` + generatePlaceholders(len(agentNames)) + `)
		AND id IN (
			SELECT MAX(id) FROM agent_metrics_history
			WHERE agent_name IN (` + generatePlaceholders(len(agentNames)) + `)
			GROUP BY agent_name
		)
	`

	// Double the placeholders since we use the list twice
	allPlaceholders := append(placeholders, placeholders...)

	var totalAgents int
	var avgCPU, avgMemory, avgDisk sql.NullFloat64

	err = w.db.QueryRowContext(ctx, query, allPlaceholders...).Scan(&totalAgents, &avgCPU, &avgMemory, &avgDisk)
	if err != nil {
		// If no metrics found, return zeros
		return map[string]interface{}{
			"total_agents":       len(agentNames),
			"healthy_agents":     0,
			"unhealthy_agents":   len(agentNames),
			"avg_cpu_percent":    0.0,
			"avg_memory_percent": 0.0,
			"avg_disk_percent":   0.0,
		}, nil
	}

	// Count healthy agents (agents with recent heartbeat)
	healthQuery := `
		SELECT COUNT(*)
		FROM agents
		WHERE name IN (` + generatePlaceholders(len(agentNames)) + `)
		AND status = 'connected'
		AND last_heartbeat > ?
	`

	var healthyAgents int
	thresholdTime := time.Now().Add(-5 * time.Minute).Unix()
	err = w.db.QueryRowContext(ctx, healthQuery, append(placeholders, thresholdTime)...).Scan(&healthyAgents)
	if err != nil {
		healthyAgents = 0
	}

	return map[string]interface{}{
		"total_agents":       len(agentNames),
		"healthy_agents":     healthyAgents,
		"unhealthy_agents":   len(agentNames) - healthyAgents,
		"avg_cpu_percent":    nullFloatToFloat(avgCPU),
		"avg_memory_percent": nullFloatToFloat(avgMemory),
		"avg_disk_percent":   nullFloatToFloat(avgDisk),
	}, nil
}

// Helper function to generate SQL placeholders
func generatePlaceholders(count int) string {
	if count == 0 {
		return ""
	}
	placeholders := "?"
	for i := 1; i < count; i++ {
		placeholders += ", ?"
	}
	return placeholders
}

// Helper function to convert sql.NullFloat64 to float64
func nullFloatToFloat(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0.0
}
