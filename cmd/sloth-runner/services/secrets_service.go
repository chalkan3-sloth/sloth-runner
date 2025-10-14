//go:build cgo
// +build cgo

package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/crypto"
	_ "github.com/mattn/go-sqlite3"
)

// SecretsService manages encrypted secrets for stacks
type SecretsService struct {
	db *sql.DB
}

// Secret represents an encrypted secret
type Secret struct {
	ID             int
	StackID        string
	Name           string
	EncryptedValue string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewSecretsService creates a new secrets service
func NewSecretsService() (*SecretsService, error) {
	// Always use /etc/sloth-runner/ as per system requirements
	dbPath := "/etc/sloth-runner/secrets.db"
	dbDir := filepath.Dir(dbPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	service := &SecretsService{db: db}

	if err := service.initialize(); err != nil {
		db.Close()
		return nil, err
	}

	return service, nil
}

// initialize creates the secrets table
func (s *SecretsService) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS secrets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		stack_id TEXT NOT NULL,
		name TEXT NOT NULL,
		encrypted_value TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		UNIQUE(stack_id, name)
	);

	CREATE INDEX IF NOT EXISTS idx_secrets_stack_id ON secrets(stack_id);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *SecretsService) Close() error {
	return s.db.Close()
}

// AddSecret adds or updates an encrypted secret for a stack
func (s *SecretsService) AddSecret(ctx context.Context, stackID, name, value, password string, salt []byte) error {
	// Encrypt the value
	encryptedValue, err := crypto.Encrypt(value, password, salt)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	now := time.Now().Unix()

	// Insert or update secret
	query := `
	INSERT INTO secrets (stack_id, name, encrypted_value, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(stack_id, name)
	DO UPDATE SET encrypted_value = excluded.encrypted_value, updated_at = excluded.updated_at
	`

	_, err = s.db.ExecContext(ctx, query, stackID, name, encryptedValue, now, now)
	if err != nil {
		return fmt.Errorf("failed to add secret: %w", err)
	}

	return nil
}

// GetSecret retrieves and decrypts a secret
func (s *SecretsService) GetSecret(ctx context.Context, stackID, name, password string, salt []byte) (string, error) {
	var encryptedValue string

	query := `SELECT encrypted_value FROM secrets WHERE stack_id = ? AND name = ?`
	err := s.db.QueryRowContext(ctx, query, stackID, name).Scan(&encryptedValue)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("secret '%s' not found for stack '%s'", name, stackID)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	// Decrypt the value
	plaintext, err := crypto.Decrypt(encryptedValue, password, salt)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return plaintext, nil
}

// ListSecrets lists all secrets for a stack (names only, values are encrypted)
func (s *SecretsService) ListSecrets(ctx context.Context, stackID string) ([]Secret, error) {
	query := `
	SELECT id, stack_id, name, encrypted_value, created_at, updated_at
	FROM secrets
	WHERE stack_id = ?
	ORDER BY name
	`

	rows, err := s.db.QueryContext(ctx, query, stackID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	defer rows.Close()

	var secrets []Secret
	for rows.Next() {
		var secret Secret
		var createdAt, updatedAt int64

		err := rows.Scan(
			&secret.ID,
			&secret.StackID,
			&secret.Name,
			&secret.EncryptedValue,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}

		secret.CreatedAt = time.Unix(createdAt, 0)
		secret.UpdatedAt = time.Unix(updatedAt, 0)

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

// GetAllSecrets retrieves and decrypts all secrets for a stack
func (s *SecretsService) GetAllSecrets(ctx context.Context, stackID, password string, salt []byte) (map[string]string, error) {
	secrets, err := s.ListSecrets(ctx, stackID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, secret := range secrets {
		plaintext, err := crypto.Decrypt(secret.EncryptedValue, password, salt)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt secret '%s': %w", secret.Name, err)
		}
		result[secret.Name] = plaintext
	}

	return result, nil
}

// RemoveSecret removes a secret from a stack
func (s *SecretsService) RemoveSecret(ctx context.Context, stackID, name string) error {
	query := `DELETE FROM secrets WHERE stack_id = ? AND name = ?`

	result, err := s.db.ExecContext(ctx, query, stackID, name)
	if err != nil {
		return fmt.Errorf("failed to remove secret: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("secret '%s' not found for stack '%s'", name, stackID)
	}

	return nil
}

// RemoveAllSecrets removes all secrets for a stack
func (s *SecretsService) RemoveAllSecrets(ctx context.Context, stackID string) error {
	query := `DELETE FROM secrets WHERE stack_id = ?`
	_, err := s.db.ExecContext(ctx, query, stackID)
	if err != nil {
		return fmt.Errorf("failed to remove secrets: %w", err)
	}
	return nil
}

// HasSecrets checks if a stack has any secrets
func (s *SecretsService) HasSecrets(ctx context.Context, stackID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM secrets WHERE stack_id = ?`
	err := s.db.QueryRowContext(ctx, query, stackID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check secrets: %w", err)
	}
	return count > 0, nil
}

// EncryptionSalt represents a salt for a stack
type EncryptionSalt struct {
	StackID string
	Salt    []byte
}

// GetOrCreateSalt gets or creates encryption salt for a stack
func GetOrCreateSalt(stackService *StackService, stackID string) ([]byte, error) {
	// Try to get existing salt from stack metadata
	stack, err := stackService.GetStack(stackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stack: %w", err)
	}

	// Check if stack has salt in metadata
	if saltB64Val, ok := stack.Metadata["encryption_salt"]; ok {
		// Type assert to string
		saltB64, ok := saltB64Val.(string)
		if !ok {
			return nil, fmt.Errorf("encryption_salt is not a string")
		}
		salt, err := base64.StdEncoding.DecodeString(saltB64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode salt: %w", err)
		}
		return salt, nil
	}

	// Generate new salt
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Store salt in stack metadata
	if stack.Metadata == nil {
		stack.Metadata = make(map[string]interface{})
	}
	stack.Metadata["encryption_salt"] = base64.StdEncoding.EncodeToString(salt)

	// Update stack using manager
	err = stackService.GetManager().UpdateStack(stack)
	if err != nil {
		return nil, fmt.Errorf("failed to update stack: %w", err)
	}

	return salt, nil
}
