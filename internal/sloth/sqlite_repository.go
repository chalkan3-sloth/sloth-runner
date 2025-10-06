package sloth

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepository implements the Repository interface using SQLite
type SQLiteRepository struct {
	db   *sql.DB
	path string
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	if dbPath == "" {
		// Use /etc/sloth-runner/ as the default location
		dbPath = "/etc/sloth-runner/sloths.db"
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create sloth directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &SQLiteRepository{
		db:   db,
		path: dbPath,
	}

	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// initSchema creates the required database tables
func (r *SQLiteRepository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sloths (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		file_path TEXT NOT NULL,
		content TEXT NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_used_at DATETIME,
		usage_count INTEGER DEFAULT 0,
		tags TEXT,
		file_hash TEXT NOT NULL,
		UNIQUE(name)
	);

	CREATE INDEX IF NOT EXISTS idx_sloths_name ON sloths(name);
	CREATE INDEX IF NOT EXISTS idx_sloths_active ON sloths(is_active);
	CREATE INDEX IF NOT EXISTS idx_sloths_updated_at ON sloths(updated_at);
	CREATE INDEX IF NOT EXISTS idx_sloths_hash ON sloths(file_hash);
	`

	_, err := r.db.Exec(schema)
	return err
}

// Create adds a new sloth to the repository
func (r *SQLiteRepository) Create(ctx context.Context, sloth *Sloth) error {
	// Check if sloth with same name already exists
	existing, err := r.GetByName(ctx, sloth.Name)
	if err == nil && existing != nil {
		return ErrSlothAlreadyExists
	}

	query := `
		INSERT INTO sloths (
			id, name, description, file_path, content, is_active,
			created_at, updated_at, tags, file_hash
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		sloth.ID,
		sloth.Name,
		sloth.Description,
		sloth.FilePath,
		sloth.Content,
		sloth.IsActive,
		sloth.CreatedAt,
		sloth.UpdatedAt,
		sloth.Tags,
		sloth.FileHash,
	)

	if err != nil {
		return fmt.Errorf("failed to create sloth: %w", err)
	}

	return nil
}

// GetByName retrieves a sloth by its name
func (r *SQLiteRepository) GetByName(ctx context.Context, name string) (*Sloth, error) {
	query := `
		SELECT id, name, description, file_path, content, is_active,
			   created_at, updated_at, last_used_at, usage_count, tags, file_hash
		FROM sloths
		WHERE name = ?
	`

	sloth := &Sloth{}
	var lastUsedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&sloth.ID,
		&sloth.Name,
		&sloth.Description,
		&sloth.FilePath,
		&sloth.Content,
		&sloth.IsActive,
		&sloth.CreatedAt,
		&sloth.UpdatedAt,
		&lastUsedAt,
		&sloth.UsageCount,
		&sloth.Tags,
		&sloth.FileHash,
	)

	if err == sql.ErrNoRows {
		return nil, ErrSlothNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get sloth: %w", err)
	}

	if lastUsedAt.Valid {
		sloth.LastUsedAt = &lastUsedAt.Time
	}

	return sloth, nil
}

// GetByID retrieves a sloth by its ID
func (r *SQLiteRepository) GetByID(ctx context.Context, id string) (*Sloth, error) {
	query := `
		SELECT id, name, description, file_path, content, is_active,
			   created_at, updated_at, last_used_at, usage_count, tags, file_hash
		FROM sloths
		WHERE id = ?
	`

	sloth := &Sloth{}
	var lastUsedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sloth.ID,
		&sloth.Name,
		&sloth.Description,
		&sloth.FilePath,
		&sloth.Content,
		&sloth.IsActive,
		&sloth.CreatedAt,
		&sloth.UpdatedAt,
		&lastUsedAt,
		&sloth.UsageCount,
		&sloth.Tags,
		&sloth.FileHash,
	)

	if err == sql.ErrNoRows {
		return nil, ErrSlothNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get sloth: %w", err)
	}

	if lastUsedAt.Valid {
		sloth.LastUsedAt = &lastUsedAt.Time
	}

	return sloth, nil
}

// List returns all sloths, optionally filtered by active status
func (r *SQLiteRepository) List(ctx context.Context, activeOnly bool) ([]*SlothListItem, error) {
	query := `
		SELECT name, description, is_active, created_at, last_used_at, usage_count
		FROM sloths
	`

	if activeOnly {
		query += " WHERE is_active = 1"
	}

	query += " ORDER BY updated_at DESC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sloths: %w", err)
	}
	defer rows.Close()

	var sloths []*SlothListItem

	for rows.Next() {
		item := &SlothListItem{}
		var lastUsedAt sql.NullTime

		err := rows.Scan(
			&item.Name,
			&item.Description,
			&item.IsActive,
			&item.CreatedAt,
			&lastUsedAt,
			&item.UsageCount,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan sloth: %w", err)
		}

		if lastUsedAt.Valid {
			item.LastUsedAt = &lastUsedAt.Time
		}

		sloths = append(sloths, item)
	}

	return sloths, rows.Err()
}

// Update updates an existing sloth
func (r *SQLiteRepository) Update(ctx context.Context, sloth *Sloth) error {
	query := `
		UPDATE sloths
		SET description = ?, file_path = ?, content = ?, is_active = ?,
			updated_at = ?, tags = ?, file_hash = ?
		WHERE name = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		sloth.Description,
		sloth.FilePath,
		sloth.Content,
		sloth.IsActive,
		time.Now(),
		sloth.Tags,
		sloth.FileHash,
		sloth.Name,
	)

	if err != nil {
		return fmt.Errorf("failed to update sloth: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrSlothNotFound
	}

	return nil
}

// Delete removes a sloth by name
func (r *SQLiteRepository) Delete(ctx context.Context, name string) error {
	query := `DELETE FROM sloths WHERE name = ?`

	result, err := r.db.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("failed to delete sloth: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrSlothNotFound
	}

	return nil
}

// SetActive sets the active status of a sloth
func (r *SQLiteRepository) SetActive(ctx context.Context, name string, active bool) error {
	query := `
		UPDATE sloths
		SET is_active = ?, updated_at = ?
		WHERE name = ?
	`

	result, err := r.db.ExecContext(ctx, query, active, time.Now(), name)
	if err != nil {
		return fmt.Errorf("failed to set active status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrSlothNotFound
	}

	return nil
}

// IncrementUsage increments the usage count and updates last used time
func (r *SQLiteRepository) IncrementUsage(ctx context.Context, name string) error {
	query := `
		UPDATE sloths
		SET usage_count = usage_count + 1, last_used_at = ?
		WHERE name = ?
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), name)
	if err != nil {
		return fmt.Errorf("failed to increment usage: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrSlothNotFound
	}

	return nil
}

// Close closes the repository connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
