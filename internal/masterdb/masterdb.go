package masterdb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Master represents a master server configuration
type Master struct {
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Description string    `json:"description"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MasterDB handles master server database operations
type MasterDB struct {
	db *sql.DB
}

// NewMasterDB creates a new MasterDB instance
func NewMasterDB(dbPath string) (*MasterDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	mdb := &MasterDB{db: db}
	if err := mdb.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return mdb, nil
}

// initSchema creates the masters table if it doesn't exist
func (m *MasterDB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS masters (
		name TEXT PRIMARY KEY,
		address TEXT NOT NULL,
		description TEXT,
		is_default INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_masters_default ON masters(is_default);
	`

	_, err := m.db.Exec(schema)
	return err
}

// Add adds a new master server configuration
// Returns error if master with the same name already exists
func (m *MasterDB) Add(master *Master) error {
	// Validate inputs
	if master.Name == "" {
		return fmt.Errorf("master name cannot be empty")
	}
	if master.Address == "" {
		return fmt.Errorf("master address cannot be empty")
	}

	// Check if master already exists
	existing, err := m.Get(master.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("master '%s' already exists with address '%s'. Use 'update' to modify it", master.Name, existing.Address)
	}

	// If this is the first master, make it default
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		master.IsDefault = true
	}

	query := `
	INSERT INTO masters (name, address, description, is_default, created_at, updated_at)
	VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`

	_, err = m.db.Exec(query, master.Name, master.Address, master.Description, master.IsDefault)
	if err != nil {
		return fmt.Errorf("failed to add master '%s': %w", master.Name, err)
	}
	return nil
}

// Update updates an existing master server configuration
func (m *MasterDB) Update(master *Master) error {
	// Validate inputs
	if master.Name == "" {
		return fmt.Errorf("master name cannot be empty")
	}
	if master.Address == "" {
		return fmt.Errorf("master address cannot be empty")
	}

	// Check if master exists
	if _, err := m.Get(master.Name); err != nil {
		return fmt.Errorf("master '%s' not found. Use 'add' to create it", master.Name)
	}

	query := `
	UPDATE masters SET
		address = ?,
		description = ?,
		updated_at = datetime('now')
	WHERE name = ?
	`

	result, err := m.db.Exec(query, master.Address, master.Description, master.Name)
	if err != nil {
		return fmt.Errorf("failed to update master '%s': %w", master.Name, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("master '%s' not found", master.Name)
	}

	return nil
}

// Get retrieves a master by name
func (m *MasterDB) Get(name string) (*Master, error) {
	master := &Master{}
	query := `SELECT name, address, description, is_default, created_at, updated_at FROM masters WHERE name = ?`

	err := m.db.QueryRow(query, name).Scan(
		&master.Name,
		&master.Address,
		&master.Description,
		&master.IsDefault,
		&master.CreatedAt,
		&master.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("master '%s' not found", name)
	}
	if err != nil {
		return nil, err
	}

	return master, nil
}

// GetDefault retrieves the default master
func (m *MasterDB) GetDefault() (*Master, error) {
	master := &Master{}
	query := `SELECT name, address, description, is_default, created_at, updated_at FROM masters WHERE is_default = 1 LIMIT 1`

	err := m.db.QueryRow(query).Scan(
		&master.Name,
		&master.Address,
		&master.Description,
		&master.IsDefault,
		&master.CreatedAt,
		&master.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no default master configured")
	}
	if err != nil {
		return nil, err
	}

	return master, nil
}

// List retrieves all masters
func (m *MasterDB) List() ([]*Master, error) {
	query := `SELECT name, address, description, is_default, created_at, updated_at FROM masters ORDER BY is_default DESC, name ASC`

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var masters []*Master
	for rows.Next() {
		master := &Master{}
		err := rows.Scan(
			&master.Name,
			&master.Address,
			&master.Description,
			&master.IsDefault,
			&master.CreatedAt,
			&master.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		masters = append(masters, master)
	}

	return masters, rows.Err()
}

// SetDefault sets a master as the default
func (m *MasterDB) SetDefault(name string) error {
	// First, verify the master exists
	if _, err := m.Get(name); err != nil {
		return err
	}

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Unset all defaults
	if _, err := tx.Exec("UPDATE masters SET is_default = 0"); err != nil {
		return err
	}

	// Set new default
	if _, err := tx.Exec("UPDATE masters SET is_default = 1, updated_at = datetime('now') WHERE name = ?", name); err != nil {
		return err
	}

	return tx.Commit()
}

// Delete removes a master
func (m *MasterDB) Delete(name string) error {
	// Check if it's the default
	master, err := m.Get(name)
	if err != nil {
		return err
	}

	if master.IsDefault {
		// If deleting default, check if there are other masters
		count, err := m.Count()
		if err != nil {
			return err
		}
		if count > 1 {
			return fmt.Errorf("cannot delete default master '%s'. Select a different default first", name)
		}
	}

	_, err = m.db.Exec("DELETE FROM masters WHERE name = ?", name)
	return err
}

// Count returns the total number of masters
func (m *MasterDB) Count() (int, error) {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM masters").Scan(&count)
	return count, err
}

// Close closes the database connection
func (m *MasterDB) Close() error {
	return m.db.Close()
}
