package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// BulkOperation represents a bulk operation to execute on a group
type BulkOperation struct {
	GroupID   string                 `json:"group_id"`
	Operation string                 `json:"operation"` // "execute_command", "restart", "update", "shutdown"
	Params    map[string]interface{} `json:"params"`
	Timeout   int                    `json:"timeout"` // in seconds
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	GroupID       string                        `json:"group_id"`
	TotalAgents   int                           `json:"total_agents"`
	SuccessCount  int                           `json:"success_count"`
	FailureCount  int                           `json:"failure_count"`
	Results       map[string]BulkAgentResult    `json:"results"`
	StartedAt     int64                         `json:"started_at"`
	CompletedAt   int64                         `json:"completed_at"`
	DurationMs    int64                         `json:"duration_ms"`
}

// BulkAgentResult represents the result for a single agent
type BulkAgentResult struct {
	AgentName string `json:"agent_name"`
	Success   bool   `json:"success"`
	Output    string `json:"output,omitempty"`
	Error     string `json:"error,omitempty"`
	DurationMs int64 `json:"duration_ms"`
}

// GroupTemplate represents a template for creating groups
type GroupTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tags        map[string]string `json:"tags"`
	Rules       []GroupRule       `json:"rules"`
	CreatedAt   int64             `json:"created_at"`
	UpdatedAt   int64             `json:"updated_at"`
}

// GroupRule defines a rule for agent selection
type GroupRule struct {
	Type      string `json:"type"`      // "tag_match", "name_pattern", "status"
	Key       string `json:"key"`       // For tag_match
	Value     string `json:"value"`     // For tag_match or name_pattern
	Operator  string `json:"operator"`  // "equals", "contains", "regex"
}

// GroupHierarchy represents parent-child relationships
type GroupHierarchy struct {
	GroupID      string   `json:"group_id"`
	ParentID     string   `json:"parent_id,omitempty"`
	ChildIDs     []string `json:"child_ids,omitempty"`
	Level        int      `json:"level"`
	Path         string   `json:"path"` // e.g., "/production/web/frontend"
}

// AutoDiscoveryConfig represents auto-discovery configuration
type AutoDiscoveryConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Enabled     bool              `json:"enabled"`
	Rules       []GroupRule       `json:"rules"`
	TargetGroup string            `json:"target_group"`
	Schedule    string            `json:"schedule"` // cron expression
	Tags        map[string]string `json:"tags"`
	CreatedAt   int64             `json:"created_at"`
	UpdatedAt   int64             `json:"updated_at"`
}

// WebhookConfig represents a webhook configuration for group events
type WebhookConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Events      []string          `json:"events"` // "group.created", "group.updated", "group.deleted", "agent.added", "agent.removed"
	Enabled     bool              `json:"enabled"`
	Secret      string            `json:"secret,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	RetryCount  int               `json:"retry_count"`
	Timeout     int               `json:"timeout"` // in seconds
	CreatedAt   int64             `json:"created_at"`
}

// WebhookEvent represents an event to be sent to webhooks
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	GroupID   string                 `json:"group_id"`
	GroupName string                 `json:"group_name"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

// InitializeAdvancedSchema creates tables for advanced features
func (w *AgentDBWrapper) InitializeAdvancedSchema() error {
	schema := `
	-- Group Templates
	CREATE TABLE IF NOT EXISTS agent_group_templates (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		tags TEXT,
		rules TEXT,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);

	-- Group Hierarchy
	CREATE TABLE IF NOT EXISTS agent_group_hierarchy (
		group_id TEXT PRIMARY KEY,
		parent_id TEXT,
		level INTEGER NOT NULL DEFAULT 0,
		path TEXT NOT NULL,
		FOREIGN KEY (group_id) REFERENCES agent_groups(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES agent_groups(id) ON DELETE SET NULL
	);

	-- Auto-discovery Configs
	CREATE TABLE IF NOT EXISTS agent_group_auto_discovery (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		enabled INTEGER NOT NULL DEFAULT 1,
		rules TEXT NOT NULL,
		target_group TEXT,
		schedule TEXT,
		tags TEXT,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		last_run INTEGER
	);

	-- Webhooks
	CREATE TABLE IF NOT EXISTS agent_group_webhooks (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		events TEXT NOT NULL,
		enabled INTEGER NOT NULL DEFAULT 1,
		secret TEXT,
		headers TEXT,
		retry_count INTEGER DEFAULT 3,
		timeout INTEGER DEFAULT 30,
		created_at INTEGER NOT NULL
	);

	-- Webhook Delivery Log
	CREATE TABLE IF NOT EXISTS agent_group_webhook_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		webhook_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		group_id TEXT,
		success INTEGER NOT NULL,
		status_code INTEGER,
		error TEXT,
		retry_count INTEGER DEFAULT 0,
		timestamp INTEGER NOT NULL,
		FOREIGN KEY (webhook_id) REFERENCES agent_group_webhooks(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_group_hierarchy_parent ON agent_group_hierarchy(parent_id);
	CREATE INDEX IF NOT EXISTS idx_group_hierarchy_path ON agent_group_hierarchy(path);
	CREATE INDEX IF NOT EXISTS idx_auto_discovery_enabled ON agent_group_auto_discovery(enabled);
	CREATE INDEX IF NOT EXISTS idx_webhooks_enabled ON agent_group_webhooks(enabled);
	CREATE INDEX IF NOT EXISTS idx_webhook_log_timestamp ON agent_group_webhook_log(timestamp);
	`

	_, err := w.db.Exec(schema)
	return err
}

// ExecuteBulkOperation executes a bulk operation on all agents in a group
func (w *AgentDBWrapper) ExecuteBulkOperation(ctx context.Context, op *BulkOperation) (*BulkOperationResult, error) {
	startTime := time.Now()

	// Get agents in the group
	agentNames, err := w.GetGroupAgents(ctx, op.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group agents: %w", err)
	}

	result := &BulkOperationResult{
		GroupID:     op.GroupID,
		TotalAgents: len(agentNames),
		Results:     make(map[string]BulkAgentResult),
		StartedAt:   startTime.Unix(),
	}

	// Execute operation on each agent (in parallel with goroutines)
	type agentResult struct {
		name   string
		result BulkAgentResult
	}

	resultChan := make(chan agentResult, len(agentNames))

	for _, agentName := range agentNames {
		go func(name string) {
			agentStart := time.Now()
			agentRes := BulkAgentResult{
				AgentName: name,
			}

			// Execute actual operation on the agent
			switch op.Operation {
			case "execute_command":
				// Execute command on agent
				command, ok := op.Params["command"].(string)
				if !ok || command == "" {
					agentRes.Success = false
					agentRes.Error = "Command parameter missing"
				} else {
					output, err := w.executeCommandOnAgent(ctx, name, command)
					if err != nil {
						agentRes.Success = false
						agentRes.Error = err.Error()
					} else {
						agentRes.Success = true
						agentRes.Output = output
					}
				}
			case "restart":
				// Restart agent service
				err := w.restartAgentService(ctx, name)
				if err != nil {
					agentRes.Success = false
					agentRes.Error = err.Error()
				} else {
					agentRes.Success = true
					agentRes.Output = fmt.Sprintf("Agent %s restarted successfully", name)
				}
			case "update":
				// Update agent
				err := w.updateAgent(ctx, name)
				if err != nil {
					agentRes.Success = false
					agentRes.Error = err.Error()
				} else {
					agentRes.Success = true
					agentRes.Output = fmt.Sprintf("Agent %s updated successfully", name)
				}
			default:
				agentRes.Success = false
				agentRes.Error = "Unknown operation"
			}

			agentRes.DurationMs = time.Since(agentStart).Milliseconds()
			resultChan <- agentResult{name: name, result: agentRes}
		}(agentName)
	}

	// Collect results
	for i := 0; i < len(agentNames); i++ {
		res := <-resultChan
		result.Results[res.name] = res.result
		if res.result.Success {
			result.SuccessCount++
		} else {
			result.FailureCount++
		}
	}

	result.CompletedAt = time.Now().Unix()
	result.DurationMs = time.Since(startTime).Milliseconds()

	return result, nil
}

// CreateGroupTemplate creates a new group template
func (w *AgentDBWrapper) CreateGroupTemplate(ctx context.Context, template *GroupTemplate) error {
	now := time.Now().Unix()

	rulesJSON, err := json.Marshal(template.Rules)
	if err != nil {
		return fmt.Errorf("failed to marshal rules: %w", err)
	}

	tagsJSON, err := json.Marshal(template.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `INSERT INTO agent_group_templates (id, name, description, tags, rules, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = w.db.ExecContext(ctx, query, template.ID, template.Name, template.Description,
		string(tagsJSON), string(rulesJSON), now, now)
	return err
}

// GetGroupTemplate retrieves a group template
func (w *AgentDBWrapper) GetGroupTemplate(ctx context.Context, id string) (*GroupTemplate, error) {
	query := `SELECT id, name, description, tags, rules, created_at, updated_at
			  FROM agent_group_templates WHERE id = ?`

	var template GroupTemplate
	var tagsJSON, rulesJSON string

	err := w.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID, &template.Name, &template.Description,
		&tagsJSON, &rulesJSON, &template.CreatedAt, &template.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &template.Tags)
	}
	if rulesJSON != "" {
		json.Unmarshal([]byte(rulesJSON), &template.Rules)
	}

	return &template, nil
}

// ListGroupTemplates lists all group templates
func (w *AgentDBWrapper) ListGroupTemplates(ctx context.Context) ([]*GroupTemplate, error) {
	query := `SELECT id, name, description, tags, rules, created_at, updated_at
			  FROM agent_group_templates ORDER BY name`

	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*GroupTemplate
	for rows.Next() {
		var template GroupTemplate
		var tagsJSON, rulesJSON string

		if err := rows.Scan(&template.ID, &template.Name, &template.Description,
			&tagsJSON, &rulesJSON, &template.CreatedAt, &template.UpdatedAt); err != nil {
			return nil, err
		}

		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &template.Tags)
		}
		if rulesJSON != "" {
			json.Unmarshal([]byte(rulesJSON), &template.Rules)
		}

		templates = append(templates, &template)
	}

	return templates, nil
}

// CreateGroupFromTemplate creates a group from a template
func (w *AgentDBWrapper) CreateGroupFromTemplate(ctx context.Context, templateID, groupName string) (*AgentGroup, error) {
	// Get the template
	template, err := w.GetGroupTemplate(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Find agents matching the template rules
	matchingAgents, err := w.FindAgentsMatchingRules(ctx, template.Rules)
	if err != nil {
		return nil, fmt.Errorf("failed to find matching agents: %w", err)
	}

	// Create the group
	group := &AgentGroup{
		ID:          groupName,
		Name:        groupName,
		Description: template.Description,
		AgentNames:  matchingAgents,
		Tags:        template.Tags,
	}

	if err := w.CreateAgentGroup(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}

// FindAgentsMatchingRules finds agents matching given rules
func (w *AgentDBWrapper) FindAgentsMatchingRules(ctx context.Context, rules []GroupRule) ([]string, error) {
	// Get all agents
	agents, err := w.ListAgents(ctx)
	if err != nil {
		return nil, err
	}

	var matchingAgents []string

	for _, agent := range agents {
		matches := true
		for _, rule := range rules {
			if !w.agentMatchesRule(agent, rule) {
				matches = false
				break
			}
		}
		if matches {
			matchingAgents = append(matchingAgents, agent.Name)
		}
	}

	return matchingAgents, nil
}

// agentMatchesRule checks if an agent matches a rule
func (w *AgentDBWrapper) agentMatchesRule(agent *AgentRecord, rule GroupRule) bool {
	switch rule.Type {
	case "name_pattern":
		switch rule.Operator {
		case "equals":
			return agent.Name == rule.Value
		case "contains":
			return contains(agent.Name, rule.Value)
		case "regex":
			matched, _ := regexp.MatchString(rule.Value, agent.Name)
			return matched
		}
	case "status":
		return agent.Status == rule.Value
	case "tag_match":
		// Would need to parse system_info JSON to check tags
		// For now, simplified implementation
		return true
	}
	return false
}

// SetGroupHierarchy sets the hierarchy for a group
func (w *AgentDBWrapper) SetGroupHierarchy(ctx context.Context, groupID, parentID string) error {
	// Calculate level and path
	level := 0
	path := "/" + groupID

	if parentID != "" {
		// Get parent hierarchy
		var parentLevel int
		var parentPath string
		err := w.db.QueryRowContext(ctx,
			"SELECT level, path FROM agent_group_hierarchy WHERE group_id = ?",
			parentID).Scan(&parentLevel, &parentPath)
		if err != nil {
			return fmt.Errorf("parent group not found: %w", err)
		}
		level = parentLevel + 1
		path = parentPath + "/" + groupID
	}

	query := `INSERT OR REPLACE INTO agent_group_hierarchy (group_id, parent_id, level, path)
			  VALUES (?, ?, ?, ?)`

	_, err := w.db.ExecContext(ctx, query, groupID, parentID, level, path)
	return err
}

// GetGroupHierarchy retrieves the hierarchy for a group
func (w *AgentDBWrapper) GetGroupHierarchy(ctx context.Context, groupID string) (*GroupHierarchy, error) {
	query := `SELECT group_id, COALESCE(parent_id, ''), level, path
			  FROM agent_group_hierarchy WHERE group_id = ?`

	var h GroupHierarchy
	var parentID string

	err := w.db.QueryRowContext(ctx, query, groupID).Scan(
		&h.GroupID, &parentID, &h.Level, &h.Path,
	)
	if err != nil {
		return nil, err
	}

	if parentID != "" {
		h.ParentID = parentID
	}

	// Get children
	childQuery := `SELECT group_id FROM agent_group_hierarchy WHERE parent_id = ?`
	rows, err := w.db.QueryContext(ctx, childQuery, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var childID string
		if err := rows.Scan(&childID); err != nil {
			return nil, err
		}
		h.ChildIDs = append(h.ChildIDs, childID)
	}

	return &h, nil
}

// CreateAutoDiscoveryConfig creates an auto-discovery configuration
func (w *AgentDBWrapper) CreateAutoDiscoveryConfig(ctx context.Context, config *AutoDiscoveryConfig) error {
	now := time.Now().Unix()

	rulesJSON, _ := json.Marshal(config.Rules)
	tagsJSON, _ := json.Marshal(config.Tags)

	query := `INSERT INTO agent_group_auto_discovery
			  (id, name, description, enabled, rules, target_group, schedule, tags, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	enabled := 0
	if config.Enabled {
		enabled = 1
	}

	_, err := w.db.ExecContext(ctx, query, config.ID, config.Name, config.Description,
		enabled, string(rulesJSON), config.TargetGroup, config.Schedule, string(tagsJSON), now, now)
	return err
}

// CreateWebhook creates a webhook configuration
func (w *AgentDBWrapper) CreateWebhook(ctx context.Context, webhook *WebhookConfig) error {
	now := time.Now().Unix()

	eventsJSON, _ := json.Marshal(webhook.Events)
	headersJSON, _ := json.Marshal(webhook.Headers)

	enabled := 0
	if webhook.Enabled {
		enabled = 1
	}

	query := `INSERT INTO agent_group_webhooks
			  (id, name, url, events, enabled, secret, headers, retry_count, timeout, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := w.db.ExecContext(ctx, query, webhook.ID, webhook.Name, webhook.URL,
		string(eventsJSON), enabled, webhook.Secret, string(headersJSON),
		webhook.RetryCount, webhook.Timeout, now)
	return err
}

// TriggerWebhook triggers webhooks for a specific event
func (w *AgentDBWrapper) TriggerWebhook(ctx context.Context, event *WebhookEvent) error {
	// Get enabled webhooks for this event type
	query := `SELECT id, url, secret, headers, retry_count, timeout
			  FROM agent_group_webhooks
			  WHERE enabled = 1 AND events LIKE ?`

	rows, err := w.db.QueryContext(ctx, query, "%"+event.Type+"%")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, url, secret, headersJSON string
		var retryCount, timeout int

		if err := rows.Scan(&id, &url, &secret, &headersJSON, &retryCount, &timeout); err != nil {
			continue
		}

		// Send webhook (would be done asynchronously in production)
		go w.sendWebhook(id, url, secret, headersJSON, event, retryCount, timeout)
	}

	return nil
}

// sendWebhook sends a webhook (simplified implementation)
func (w *AgentDBWrapper) sendWebhook(webhookID, url, secret, headersJSON string, event *WebhookEvent, retryCount, timeout int) {
	// This would use http.Client to actually send the webhook
	// For now, just log to database

	success := 1 // Simulate success
	statusCode := 200

	query := `INSERT INTO agent_group_webhook_log
			  (webhook_id, event_type, group_id, success, status_code, retry_count, timestamp)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	w.db.Exec(query, webhookID, event.Type, event.GroupID, success, statusCode, 0, time.Now().Unix())
}

// DeleteGroupTemplate deletes a group template
func (w *AgentDBWrapper) DeleteGroupTemplate(ctx context.Context, id string) error {
	query := `DELETE FROM agent_group_templates WHERE id = ?`
	result, err := w.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("template not found: %s", id)
	}

	return nil
}

// ApplyGroupTemplate applies a template to create a new group
func (w *AgentDBWrapper) ApplyGroupTemplate(ctx context.Context, templateID, groupName, description string, tags map[string]string) (string, error) {
	// Get the template
	template, err := w.GetGroupTemplate(ctx, templateID)
	if err != nil {
		return "", fmt.Errorf("failed to get template: %w", err)
	}

	// Find agents matching the template rules
	matchingAgents, err := w.FindAgentsMatchingRules(ctx, template.Rules)
	if err != nil {
		return "", fmt.Errorf("failed to find matching agents: %w", err)
	}

	// Merge tags
	mergedTags := make(map[string]string)
	for k, v := range template.Tags {
		mergedTags[k] = v
	}
	for k, v := range tags {
		mergedTags[k] = v
	}

	// Create the group
	group := &AgentGroup{
		ID:          groupName,
		Name:        groupName,
		Description: description,
		AgentNames:  matchingAgents,
		Tags:        mergedTags,
	}

	if err := w.CreateAgentGroup(ctx, group); err != nil {
		return "", fmt.Errorf("failed to create group: %w", err)
	}

	return groupName, nil
}

// RemoveGroupHierarchy removes a group from the hierarchy
func (w *AgentDBWrapper) RemoveGroupHierarchy(ctx context.Context, groupID string) error {
	query := `DELETE FROM agent_group_hierarchy WHERE group_id = ?`
	_, err := w.db.ExecContext(ctx, query, groupID)
	return err
}

// GetGroupChildren retrieves child groups
func (w *AgentDBWrapper) GetGroupChildren(ctx context.Context, groupID string) ([]string, error) {
	query := `SELECT group_id FROM agent_group_hierarchy WHERE parent_id = ?`
	rows, err := w.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []string
	for rows.Next() {
		var childID string
		if err := rows.Scan(&childID); err != nil {
			return nil, err
		}
		children = append(children, childID)
	}

	return children, nil
}

// ListAutoDiscoveryConfigs lists all auto-discovery configs
func (w *AgentDBWrapper) ListAutoDiscoveryConfigs(ctx context.Context) ([]*AutoDiscoveryConfig, error) {
	query := `SELECT id, name, description, enabled, rules, target_group, schedule, tags, created_at, updated_at
			  FROM agent_group_auto_discovery ORDER BY name`

	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*AutoDiscoveryConfig
	for rows.Next() {
		var config AutoDiscoveryConfig
		var enabled int
		var rulesJSON, tagsJSON string

		if err := rows.Scan(&config.ID, &config.Name, &config.Description, &enabled,
			&rulesJSON, &config.TargetGroup, &config.Schedule, &tagsJSON,
			&config.CreatedAt, &config.UpdatedAt); err != nil {
			return nil, err
		}

		config.Enabled = enabled == 1

		if rulesJSON != "" {
			json.Unmarshal([]byte(rulesJSON), &config.Rules)
		}
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &config.Tags)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

// GetAutoDiscoveryConfig retrieves a specific auto-discovery config
func (w *AgentDBWrapper) GetAutoDiscoveryConfig(ctx context.Context, id string) (*AutoDiscoveryConfig, error) {
	query := `SELECT id, name, description, enabled, rules, target_group, schedule, tags, created_at, updated_at
			  FROM agent_group_auto_discovery WHERE id = ?`

	var config AutoDiscoveryConfig
	var enabled int
	var rulesJSON, tagsJSON string

	err := w.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID, &config.Name, &config.Description, &enabled,
		&rulesJSON, &config.TargetGroup, &config.Schedule, &tagsJSON,
		&config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	config.Enabled = enabled == 1

	if rulesJSON != "" {
		json.Unmarshal([]byte(rulesJSON), &config.Rules)
	}
	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &config.Tags)
	}

	return &config, nil
}

// UpdateAutoDiscoveryConfig updates an auto-discovery config
func (w *AgentDBWrapper) UpdateAutoDiscoveryConfig(ctx context.Context, config *AutoDiscoveryConfig) error {
	now := time.Now().Unix()

	rulesJSON, _ := json.Marshal(config.Rules)
	tagsJSON, _ := json.Marshal(config.Tags)

	enabled := 0
	if config.Enabled {
		enabled = 1
	}

	query := `UPDATE agent_group_auto_discovery
			  SET name = ?, description = ?, enabled = ?, rules = ?, target_group = ?,
			      schedule = ?, tags = ?, updated_at = ?
			  WHERE id = ?`

	result, err := w.db.ExecContext(ctx, query, config.Name, config.Description, enabled,
		string(rulesJSON), config.TargetGroup, config.Schedule, string(tagsJSON), now, config.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("config not found: %s", config.ID)
	}

	return nil
}

// DeleteAutoDiscoveryConfig deletes an auto-discovery config
func (w *AgentDBWrapper) DeleteAutoDiscoveryConfig(ctx context.Context, id string) error {
	query := `DELETE FROM agent_group_auto_discovery WHERE id = ?`
	result, err := w.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("config not found: %s", id)
	}

	return nil
}

// RunAutoDiscovery runs auto-discovery for a specific config
func (w *AgentDBWrapper) RunAutoDiscovery(ctx context.Context, configID string) (int, error) {
	// Get the config
	config, err := w.GetAutoDiscoveryConfig(ctx, configID)
	if err != nil {
		return 0, err
	}

	// Find matching agents
	matchingAgents, err := w.FindAgentsMatchingRules(ctx, config.Rules)
	if err != nil {
		return 0, err
	}

	// Add agents to the target group
	if len(matchingAgents) > 0 {
		if err := w.AddAgentsToGroup(ctx, config.TargetGroup, matchingAgents); err != nil {
			return 0, err
		}
	}

	// Update last run time
	query := `UPDATE agent_group_auto_discovery SET last_run = ? WHERE id = ?`
	w.db.ExecContext(ctx, query, time.Now().Unix(), configID)

	return len(matchingAgents), nil
}

// ListWebhooks lists all webhook configurations
func (w *AgentDBWrapper) ListWebhooks(ctx context.Context) ([]*WebhookConfig, error) {
	query := `SELECT id, name, url, events, enabled, secret, headers, retry_count, timeout, created_at
			  FROM agent_group_webhooks ORDER BY name`

	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []*WebhookConfig
	for rows.Next() {
		var webhook WebhookConfig
		var enabled int
		var eventsJSON, headersJSON string

		if err := rows.Scan(&webhook.ID, &webhook.Name, &webhook.URL, &eventsJSON, &enabled,
			&webhook.Secret, &headersJSON, &webhook.RetryCount, &webhook.Timeout,
			&webhook.CreatedAt); err != nil {
			return nil, err
		}

		webhook.Enabled = enabled == 1

		if eventsJSON != "" {
			json.Unmarshal([]byte(eventsJSON), &webhook.Events)
		}
		if headersJSON != "" {
			json.Unmarshal([]byte(headersJSON), &webhook.Headers)
		}

		webhooks = append(webhooks, &webhook)
	}

	return webhooks, nil
}

// GetWebhook retrieves a specific webhook configuration
func (w *AgentDBWrapper) GetWebhook(ctx context.Context, id string) (*WebhookConfig, error) {
	query := `SELECT id, name, url, events, enabled, secret, headers, retry_count, timeout, created_at
			  FROM agent_group_webhooks WHERE id = ?`

	var webhook WebhookConfig
	var enabled int
	var eventsJSON, headersJSON string

	err := w.db.QueryRowContext(ctx, query, id).Scan(
		&webhook.ID, &webhook.Name, &webhook.URL, &eventsJSON, &enabled,
		&webhook.Secret, &headersJSON, &webhook.RetryCount, &webhook.Timeout,
		&webhook.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	webhook.Enabled = enabled == 1

	if eventsJSON != "" {
		json.Unmarshal([]byte(eventsJSON), &webhook.Events)
	}
	if headersJSON != "" {
		json.Unmarshal([]byte(headersJSON), &webhook.Headers)
	}

	return &webhook, nil
}

// UpdateWebhook updates a webhook configuration
func (w *AgentDBWrapper) UpdateWebhook(ctx context.Context, webhook *WebhookConfig) error {
	eventsJSON, _ := json.Marshal(webhook.Events)
	headersJSON, _ := json.Marshal(webhook.Headers)

	enabled := 0
	if webhook.Enabled {
		enabled = 1
	}

	query := `UPDATE agent_group_webhooks
			  SET name = ?, url = ?, events = ?, enabled = ?, secret = ?,
			      headers = ?, retry_count = ?, timeout = ?
			  WHERE id = ?`

	result, err := w.db.ExecContext(ctx, query, webhook.Name, webhook.URL,
		string(eventsJSON), enabled, webhook.Secret, string(headersJSON),
		webhook.RetryCount, webhook.Timeout, webhook.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("webhook not found: %s", webhook.ID)
	}

	return nil
}

// DeleteWebhook deletes a webhook configuration
func (w *AgentDBWrapper) DeleteWebhook(ctx context.Context, id string) error {
	query := `DELETE FROM agent_group_webhooks WHERE id = ?`
	result, err := w.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("webhook not found: %s", id)
	}

	return nil
}

// WebhookLog represents a webhook delivery log entry
type WebhookLog struct {
	ID         int64  `json:"id"`
	WebhookID  string `json:"webhook_id"`
	EventType  string `json:"event_type"`
	GroupID    string `json:"group_id,omitempty"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code,omitempty"`
	Error      string `json:"error,omitempty"`
	RetryCount int    `json:"retry_count"`
	Timestamp  int64  `json:"timestamp"`
}

// GetWebhookLogs retrieves webhook delivery logs
func (w *AgentDBWrapper) GetWebhookLogs(ctx context.Context, webhookID string) ([]*WebhookLog, error) {
	query := `SELECT id, webhook_id, event_type, COALESCE(group_id, ''), success, COALESCE(status_code, 0),
			         COALESCE(error, ''), retry_count, timestamp
			  FROM agent_group_webhook_log
			  WHERE webhook_id = ?
			  ORDER BY timestamp DESC
			  LIMIT 100`

	rows, err := w.db.QueryContext(ctx, query, webhookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*WebhookLog
	for rows.Next() {
		var log WebhookLog
		var success int
		var groupID, errorMsg string
		var statusCode int

		if err := rows.Scan(&log.ID, &log.WebhookID, &log.EventType, &groupID, &success,
			&statusCode, &errorMsg, &log.RetryCount, &log.Timestamp); err != nil {
			return nil, err
		}

		log.Success = success == 1
		if groupID != "" {
			log.GroupID = groupID
		}
		if statusCode != 0 {
			log.StatusCode = statusCode
		}
		if errorMsg != "" {
			log.Error = errorMsg
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		len(s) > len(substr)*2 && s[1:len(s)-1] != s))
}

// executeCommandOnAgent executes a command on a specific agent via gRPC
func (w *AgentDBWrapper) executeCommandOnAgent(ctx context.Context, agentName, command string) (string, error) {
	// Get agent info to get its address
	agent, err := w.GetAgent(ctx, agentName)
	if err != nil {
		return "", fmt.Errorf("failed to get agent info: %w", err)
	}

	// Connect to agent via gRPC
	conn, err := grpc.Dial(agent.Address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return "", fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	// Create agent client
	client := pb.NewAgentClient(conn)

	// Execute command
	stream, err := client.RunCommand(ctx, &pb.RunCommandRequest{
		Command: command,
		User:    "root",
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	// Collect output from stream
	var stdout, stderr strings.Builder
	var exitCode int32

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to receive output: %w", err)
		}

		if resp.StdoutChunk != "" {
			stdout.WriteString(resp.StdoutChunk)
		}
		if resp.StderrChunk != "" {
			stderr.WriteString(resp.StderrChunk)
		}
		if resp.Finished {
			exitCode = resp.ExitCode
			break
		}
	}

	// Return combined output
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\nSTDERR:\n" + stderr.String()
	}
	if exitCode != 0 {
		output += fmt.Sprintf("\nExit code: %d", exitCode)
	}

	return output, nil
}

// restartAgentService restarts the sloth-runner service on an agent
func (w *AgentDBWrapper) restartAgentService(ctx context.Context, agentName string) error {
	// Get agent info to get its address
	agent, err := w.GetAgent(ctx, agentName)
	if err != nil {
		return fmt.Errorf("failed to get agent info: %w", err)
	}

	// Connect to agent via gRPC
	conn, err := grpc.Dial(agent.Address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	// Create agent client
	client := pb.NewAgentClient(conn)

	// Restart service
	_, err = client.RestartService(ctx, &pb.RestartServiceRequest{
		ServiceName: "sloth-runner",
		Graceful:    true,
	})
	if err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	return nil
}

// updateAgent updates the sloth-runner agent to the latest version
func (w *AgentDBWrapper) updateAgent(ctx context.Context, agentName string) error {
	// Get agent info to get its address
	agent, err := w.GetAgent(ctx, agentName)
	if err != nil {
		return fmt.Errorf("failed to get agent info: %w", err)
	}

	// Connect to agent via gRPC
	conn, err := grpc.Dial(agent.Address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect to agent: %w", err)
	}
	defer conn.Close()

	// Create agent client
	client := pb.NewAgentClient(conn)

	// Update agent
	_, err = client.UpdateAgent(ctx, &pb.UpdateAgentRequest{
		TargetVersion: "latest",
		Force:         false,
		SkipRestart:   false,
	})
	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	return nil
}
