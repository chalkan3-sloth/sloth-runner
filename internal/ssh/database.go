package ssh

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Profile represents an SSH connection profile
type Profile struct {
	Name               string    `json:"name"`
	Host               string    `json:"host"`
	User               string    `json:"user"`
	Port               int       `json:"port"`
	KeyPath            string    `json:"key_path,omitempty"`
	Description        string    `json:"description,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	LastUsed           *time.Time `json:"last_used,omitempty"`
	UseCount           int       `json:"use_count"`
	ConnectionTimeout  int       `json:"connection_timeout"`
	KeepaliveInterval  int       `json:"keepalive_interval"`
	StrictHostChecking bool      `json:"strict_host_checking"`
}

// AuditLog represents an SSH operation audit entry
type AuditLog struct {
	ID               int       `json:"id"`
	ProfileName      string    `json:"profile_name"`
	Action           string    `json:"action"`
	Command          string    `json:"command,omitempty"`
	Timestamp        time.Time `json:"timestamp"`
	Success          bool      `json:"success"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	DurationMs       int       `json:"duration_ms,omitempty"`
	BytesTransferred int       `json:"bytes_transferred,omitempty"`
}

// Database manages SSH profiles in SQLite
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new SSH profile database manager
func NewDatabase(dbPath string) (*Database, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set database file permissions (only if file exists)
	if _, err := os.Stat(dbPath); err == nil {
		if err := os.Chmod(dbPath, 0600); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set database permissions: %w", err)
		}
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Database{db: db}, nil
}

// GetDefaultDatabasePath returns the default database path
func GetDefaultDatabasePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".sloth-runner", "ssh_profiles.db")
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	// Main profiles table
	profilesSQL := `
	CREATE TABLE IF NOT EXISTS ssh_profiles (
		name TEXT PRIMARY KEY,
		host TEXT NOT NULL,
		user TEXT NOT NULL,
		port INTEGER DEFAULT 22,
		key_path TEXT,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_used TIMESTAMP,
		use_count INTEGER DEFAULT 0,
		connection_timeout INTEGER DEFAULT 30,
		keepalive_interval INTEGER DEFAULT 60,
		strict_host_checking BOOLEAN DEFAULT 1,
		UNIQUE(host, user, port),
		CHECK(length(name) > 0 AND length(name) <= 50),
		CHECK(length(host) > 0),
		CHECK(length(user) > 0),
		CHECK(port > 0 AND port <= 65535)
	);`

	// Audit log table
	auditSQL := `
	CREATE TABLE IF NOT EXISTS ssh_audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_name TEXT NOT NULL,
		action TEXT NOT NULL,
		command TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		success BOOLEAN,
		error_message TEXT,
		duration_ms INTEGER,
		bytes_transferred INTEGER,
		FOREIGN KEY(profile_name) REFERENCES ssh_profiles(name) ON DELETE CASCADE
	);`

	// Create indexes
	indexSQL := []string{
		"CREATE INDEX IF NOT EXISTS idx_ssh_profiles_host ON ssh_profiles(host);",
		"CREATE INDEX IF NOT EXISTS idx_ssh_profiles_last_used ON ssh_profiles(last_used);",
		"CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON ssh_audit_log(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_audit_profile ON ssh_audit_log(profile_name);",
	}

	// Execute table creation
	if _, err := db.Exec(profilesSQL); err != nil {
		return fmt.Errorf("failed to create profiles table: %w", err)
	}

	if _, err := db.Exec(auditSQL); err != nil {
		return fmt.Errorf("failed to create audit table: %w", err)
	}

	// Create indexes
	for _, idx := range indexSQL {
		if _, err := db.Exec(idx); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	// Create trigger for updating timestamp
	triggerSQL := `
	CREATE TRIGGER IF NOT EXISTS update_ssh_profile_timestamp
	AFTER UPDATE ON ssh_profiles
	BEGIN
		UPDATE ssh_profiles
		SET updated_at = CURRENT_TIMESTAMP
		WHERE name = NEW.name;
	END;`

	if _, err := db.Exec(triggerSQL); err != nil {
		return fmt.Errorf("failed to create trigger: %w", err)
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// AddProfile adds a new SSH profile
func (d *Database) AddProfile(profile *Profile) error {
	query := `
		INSERT INTO ssh_profiles (
			name, host, user, port, key_path, description,
			connection_timeout, keepalive_interval, strict_host_checking
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.db.Exec(query,
		profile.Name,
		profile.Host,
		profile.User,
		profile.Port,
		profile.KeyPath,
		profile.Description,
		profile.ConnectionTimeout,
		profile.KeepaliveInterval,
		profile.StrictHostChecking,
	)

	if err != nil {
		return fmt.Errorf("failed to add profile: %w", err)
	}

	return nil
}

// GetProfile retrieves a profile by name
func (d *Database) GetProfile(name string) (*Profile, error) {
	query := `
		SELECT name, host, user, port, key_path, description,
			   created_at, updated_at, last_used, use_count,
			   connection_timeout, keepalive_interval, strict_host_checking
		FROM ssh_profiles
		WHERE name = ?`

	var profile Profile
	var lastUsed sql.NullTime

	err := d.db.QueryRow(query, name).Scan(
		&profile.Name,
		&profile.Host,
		&profile.User,
		&profile.Port,
		&profile.KeyPath,
		&profile.Description,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&lastUsed,
		&profile.UseCount,
		&profile.ConnectionTimeout,
		&profile.KeepaliveInterval,
		&profile.StrictHostChecking,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile '%s' not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	if lastUsed.Valid {
		profile.LastUsed = &lastUsed.Time
	}

	return &profile, nil
}

// ListProfiles lists all profiles
func (d *Database) ListProfiles() ([]*Profile, error) {
	query := `
		SELECT name, host, user, port, key_path, description,
			   created_at, updated_at, last_used, use_count,
			   connection_timeout, keepalive_interval, strict_host_checking
		FROM ssh_profiles
		ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}
	defer rows.Close()

	var profiles []*Profile
	for rows.Next() {
		var profile Profile
		var lastUsed sql.NullTime

		err := rows.Scan(
			&profile.Name,
			&profile.Host,
			&profile.User,
			&profile.Port,
			&profile.KeyPath,
			&profile.Description,
			&profile.CreatedAt,
			&profile.UpdatedAt,
			&lastUsed,
			&profile.UseCount,
			&profile.ConnectionTimeout,
			&profile.KeepaliveInterval,
			&profile.StrictHostChecking,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}

		if lastUsed.Valid {
			profile.LastUsed = &lastUsed.Time
		}

		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

// UpdateProfile updates an existing profile
func (d *Database) UpdateProfile(name string, updates map[string]interface{}) error {
	// Build dynamic update query
	var setClause []string
	var args []interface{}

	for field, value := range updates {
		switch field {
		case "host", "user", "key_path", "description":
			setClause = append(setClause, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		case "port", "connection_timeout", "keepalive_interval":
			setClause = append(setClause, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		case "strict_host_checking":
			setClause = append(setClause, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		}
	}

	if len(setClause) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	// Add name for WHERE clause
	args = append(args, name)

	query := fmt.Sprintf("UPDATE ssh_profiles SET %s WHERE name = ?",
		joinStrings(setClause, ", "))

	result, err := d.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("profile '%s' not found", name)
	}

	return nil
}

// RemoveProfile removes a profile
func (d *Database) RemoveProfile(name string) error {
	query := "DELETE FROM ssh_profiles WHERE name = ?"

	result, err := d.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to remove profile: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("profile '%s' not found", name)
	}

	return nil
}

// AddAuditLog adds an audit log entry
func (d *Database) AddAuditLog(log *AuditLog) error {
	query := `
		INSERT INTO ssh_audit_log (
			profile_name, action, command, success,
			error_message, duration_ms, bytes_transferred
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := d.db.Exec(query,
		log.ProfileName,
		log.Action,
		log.Command,
		log.Success,
		log.ErrorMessage,
		log.DurationMs,
		log.BytesTransferred,
	)

	if err != nil {
		return fmt.Errorf("failed to add audit log: %w", err)
	}

	// Update usage statistics if successful execution
	if log.Action == "execute" && log.Success {
		updateQuery := `
			UPDATE ssh_profiles
			SET last_used = CURRENT_TIMESTAMP, use_count = use_count + 1
			WHERE name = ?`

		d.db.Exec(updateQuery, log.ProfileName)
	}

	return nil
}

// GetAuditLogs retrieves audit logs for a profile
func (d *Database) GetAuditLogs(profileName string, limit int) ([]*AuditLog, error) {
	query := `
		SELECT id, profile_name, action, command, timestamp,
			   success, error_message, duration_ms, bytes_transferred
		FROM ssh_audit_log
		WHERE profile_name = ?
		ORDER BY timestamp DESC
		LIMIT ?`

	rows, err := d.db.Query(query, profileName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.ProfileName,
			&log.Action,
			&log.Command,
			&log.Timestamp,
			&log.Success,
			&log.ErrorMessage,
			&log.DurationMs,
			&log.BytesTransferred,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, &log)
	}

	return logs, nil
}

// ValidateKeyFile validates SSH key file permissions
func ValidateKeyFile(keyPath string) error {
	info, err := os.Stat(keyPath)
	if err != nil {
		return fmt.Errorf("cannot access key file: %w", err)
	}

	mode := info.Mode()
	if mode.Perm() != 0600 {
		return fmt.Errorf("private key file has incorrect permissions %v (must be 600)", mode.Perm())
	}

	return nil
}

// joinStrings joins strings with a separator (helper function)
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}