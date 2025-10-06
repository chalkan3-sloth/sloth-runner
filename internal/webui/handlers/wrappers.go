package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/chalkan3-sloth/sloth-runner/internal/sloth"
	"github.com/chalkan3-sloth/sloth-runner/internal/ssh"
)

// AgentDBWrapper wraps agent database operations
type AgentDBWrapper struct {
	db *sql.DB
}

// AgentRecord represents an agent
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

// NewAgentDBWrapper creates a new agent DB wrapper
func NewAgentDBWrapper(dbPath string) (*AgentDBWrapper, error) {
	if dbPath == "" {
		dbPath = filepath.Join(".sloth-cache", "agents.db")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open agent database: %w", err)
	}

	return &AgentDBWrapper{db: db}, nil
}

// ListAgents returns all agents
func (w *AgentDBWrapper) ListAgents(ctx context.Context) ([]*AgentRecord, error) {
	query := `SELECT id, name, address, status, last_heartbeat, registered_at, updated_at,
			  last_info_collected, system_info, COALESCE(version, '') as version
			  FROM agents ORDER BY name`

	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*AgentRecord
	for rows.Next() {
		var agent AgentRecord
		if err := rows.Scan(&agent.ID, &agent.Name, &agent.Address, &agent.Status,
			&agent.LastHeartbeat, &agent.RegisteredAt, &agent.UpdatedAt,
			&agent.LastInfoCollected, &agent.SystemInfo, &agent.Version); err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}

	return agents, nil
}

// GetAgent returns an agent by name
func (w *AgentDBWrapper) GetAgent(ctx context.Context, name string) (*AgentRecord, error) {
	query := `SELECT id, name, address, status, last_heartbeat, registered_at, updated_at,
			  last_info_collected, system_info, COALESCE(version, '') as version
			  FROM agents WHERE name = ?`

	var agent AgentRecord
	err := w.db.QueryRowContext(ctx, query, name).Scan(
		&agent.ID, &agent.Name, &agent.Address, &agent.Status,
		&agent.LastHeartbeat, &agent.RegisteredAt, &agent.UpdatedAt,
		&agent.LastInfoCollected, &agent.SystemInfo, &agent.Version,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

// DeleteAgent removes an agent
func (w *AgentDBWrapper) DeleteAgent(ctx context.Context, name string) error {
	_, err := w.db.ExecContext(ctx, "DELETE FROM agents WHERE name = ?", name)
	return err
}

// GetStats returns agent statistics
func (w *AgentDBWrapper) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var total int
	if err := w.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM agents").Scan(&total); err != nil {
		return nil, err
	}
	stats["total"] = total

	return stats, nil
}

// Close closes the database connection
func (w *AgentDBWrapper) Close() error {
	if w.db != nil {
		return w.db.Close()
	}
	return nil
}

// SlothRepoWrapper wraps sloth repository operations
type SlothRepoWrapper struct {
	repo *sloth.SQLiteRepository
}

// NewSlothRepoWrapper creates a new sloth repository wrapper
func NewSlothRepoWrapper(dbPath string) (*SlothRepoWrapper, error) {
	if dbPath == "" {
		dbPath = "/etc/sloth-runner/sloths.db"
	}

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	repo, err := sloth.NewSQLiteRepository(dbPath)
	if err != nil {
		return nil, err
	}

	return &SlothRepoWrapper{repo: repo}, nil
}

// List returns all sloths
func (w *SlothRepoWrapper) List(ctx context.Context, activeOnly bool) ([]*sloth.SlothListItem, error) {
	return w.repo.List(ctx, activeOnly)
}

// Get returns a sloth by name
func (w *SlothRepoWrapper) Get(ctx context.Context, name string) (*sloth.Sloth, error) {
	return w.repo.GetByName(ctx, name)
}

// Create creates a new sloth
func (w *SlothRepoWrapper) Create(ctx context.Context, s *sloth.Sloth) error {
	return w.repo.Create(ctx, s)
}

// Update updates a sloth
func (w *SlothRepoWrapper) Update(ctx context.Context, s *sloth.Sloth) error {
	return w.repo.Update(ctx, s)
}

// Delete deletes a sloth
func (w *SlothRepoWrapper) Delete(ctx context.Context, name string) error {
	return w.repo.Delete(ctx, name)
}

// SetActive sets a sloth's active status
func (w *SlothRepoWrapper) SetActive(ctx context.Context, name string, active bool) error {
	return w.repo.SetActive(ctx, name, active)
}

// Close closes the repository
func (w *SlothRepoWrapper) Close() error {
	if w.repo != nil {
		return w.repo.Close()
	}
	return nil
}

// HookRepoWrapper wraps hook repository operations
type HookRepoWrapper struct {
	repo *hooks.Repository
}

// NewHookRepoWrapper creates a new hook repository wrapper
func NewHookRepoWrapper(dbPath string) (*HookRepoWrapper, error) {
	// The hooks.NewRepository creates its own DB path
	repo, err := hooks.NewRepository()
	if err != nil {
		return nil, err
	}

	return &HookRepoWrapper{repo: repo}, nil
}

// List returns all hooks
func (w *HookRepoWrapper) List() ([]*hooks.Hook, error) {
	return w.repo.List()
}

// Get returns a hook by ID
func (w *HookRepoWrapper) Get(id string) (*hooks.Hook, error) {
	return w.repo.Get(id)
}

// Add adds a new hook
func (w *HookRepoWrapper) Add(hook *hooks.Hook) error {
	return w.repo.Add(hook)
}

// Update updates a hook
func (w *HookRepoWrapper) Update(hook *hooks.Hook) error {
	return w.repo.Update(hook)
}

// Delete deletes a hook
func (w *HookRepoWrapper) Delete(id string) error {
	return w.repo.Delete(id)
}

// Enable enables a hook
func (w *HookRepoWrapper) Enable(id string) error {
	return w.repo.Enable(id)
}

// Disable disables a hook
func (w *HookRepoWrapper) Disable(id string) error {
	return w.repo.Disable(id)
}

// GetExecutionHistory returns execution history
func (w *HookRepoWrapper) GetExecutionHistory(hookID string, limit int) ([]*hooks.HookResult, error) {
	return w.repo.GetExecutionHistory(hookID, limit)
}

// GetEventQueue returns the event queue
func (w *HookRepoWrapper) GetEventQueue() *hooks.EventQueue {
	return w.repo.EventQueue
}

// Close closes the repository
func (w *HookRepoWrapper) Close() error {
	if w.repo != nil {
		return w.repo.Close()
	}
	return nil
}

// SecretsServiceWrapper wraps secrets service operations
type SecretsServiceWrapper struct {
	db *sql.DB
}

// NewSecretsServiceWrapper creates a new secrets service wrapper
func NewSecretsServiceWrapper(dbPath string) (*SecretsServiceWrapper, error) {
	if dbPath == "" {
		homeDir, _ := os.UserHomeDir()
		dbPath = filepath.Join(homeDir, ".sloth-runner", "secrets.db")
	}

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &SecretsServiceWrapper{db: db}, nil
}

// ListSecrets lists secrets for a stack (names only)
func (w *SecretsServiceWrapper) ListSecrets(ctx context.Context, stackID string) ([]string, error) {
	query := "SELECT name FROM secrets WHERE stack_id = ? ORDER BY name"

	rows, err := w.db.QueryContext(ctx, query, stackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		secrets = append(secrets, name)
	}

	return secrets, nil
}

// Close closes the database connection
func (w *SecretsServiceWrapper) Close() error {
	if w.db != nil {
		return w.db.Close()
	}
	return nil
}

// SSHDBWrapper wraps SSH database operations
type SSHDBWrapper struct {
	db *ssh.Database
}

// NewSSHDBWrapper creates a new SSH database wrapper
func NewSSHDBWrapper(dbPath string) (*SSHDBWrapper, error) {
	if dbPath == "" {
		dbPath = ssh.GetDefaultDatabasePath()
	}

	db, err := ssh.NewDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	return &SSHDBWrapper{db: db}, nil
}

// ListProfiles returns all SSH profiles
func (w *SSHDBWrapper) ListProfiles() ([]*ssh.Profile, error) {
	return w.db.ListProfiles()
}

// GetProfile returns an SSH profile by name
func (w *SSHDBWrapper) GetProfile(name string) (*ssh.Profile, error) {
	return w.db.GetProfile(name)
}

// AddProfile adds a new SSH profile
func (w *SSHDBWrapper) AddProfile(profile *ssh.Profile) error {
	return w.db.AddProfile(profile)
}

// UpdateProfile updates an SSH profile
func (w *SSHDBWrapper) UpdateProfile(name string, updates map[string]interface{}) error {
	return w.db.UpdateProfile(name, updates)
}

// RemoveProfile removes an SSH profile
func (w *SSHDBWrapper) RemoveProfile(name string) error {
	return w.db.RemoveProfile(name)
}

// GetAuditLogs returns audit logs for a profile
func (w *SSHDBWrapper) GetAuditLogs(profileName string, limit int) ([]*ssh.AuditLog, error) {
	return w.db.GetAuditLogs(profileName, limit)
}

// Close closes the database connection
func (w *SSHDBWrapper) Close() error {
	if w.db != nil {
		return w.db.Close()
	}
	return nil
}
