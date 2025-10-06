package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/sloth"
	"github.com/google/uuid"
)

// SlothService handles sloth operations
// This implements the Service Layer pattern with business logic
type SlothService struct {
	repo sloth.Repository
}

// NewSlothService creates a new sloth service with default SQLite repository
func NewSlothService() (*SlothService, error) {
	repo, err := sloth.NewSQLiteRepository("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sloth repository: %w", err)
	}
	return &SlothService{repo: repo}, nil
}

// NewSlothServiceWithRepository creates a new sloth service with a custom repository
// This allows dependency injection for testing
func NewSlothServiceWithRepository(repo sloth.Repository) *SlothService {
	return &SlothService{repo: repo}
}

// AddSloth adds a new sloth from a file
func (s *SlothService) AddSloth(ctx context.Context, name, filePath, description string, active bool) error {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate file hash
	hash := sha256.Sum256(content)
	fileHash := fmt.Sprintf("%x", hash)

	// Create sloth object
	newSloth := &sloth.Sloth{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		FilePath:    filePath,
		Content:     string(content),
		IsActive:    active,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UsageCount:  0,
		FileHash:    fileHash,
	}

	// Add to repository
	if err := s.repo.Create(ctx, newSloth); err != nil {
		return fmt.Errorf("failed to add sloth: %w", err)
	}

	return nil
}

// GetSloth retrieves a sloth by name
func (s *SlothService) GetSloth(ctx context.Context, name string) (*sloth.Sloth, error) {
	return s.repo.GetByName(ctx, name)
}

// GetActiveSloth retrieves a sloth by name and checks if it's active
func (s *SlothService) GetActiveSloth(ctx context.Context, name string) (*sloth.Sloth, error) {
	sl, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if !sl.IsActive {
		return nil, sloth.ErrSlothInactive
	}

	return sl, nil
}

// ListSloths returns all sloths, optionally filtered by active status
func (s *SlothService) ListSloths(ctx context.Context, activeOnly bool) ([]*sloth.SlothListItem, error) {
	return s.repo.List(ctx, activeOnly)
}

// RemoveSloth removes a sloth by name
func (s *SlothService) RemoveSloth(ctx context.Context, name string) error {
	return s.repo.Delete(ctx, name)
}

// DeleteSloth is an alias for RemoveSloth (same functionality)
func (s *SlothService) DeleteSloth(ctx context.Context, name string) error {
	return s.RemoveSloth(ctx, name)
}

// ActivateSloth sets a sloth as active
func (s *SlothService) ActivateSloth(ctx context.Context, name string) error {
	return s.repo.SetActive(ctx, name, true)
}

// DeactivateSloth sets a sloth as inactive
func (s *SlothService) DeactivateSloth(ctx context.Context, name string) error {
	return s.repo.SetActive(ctx, name, false)
}

// UseSloth increments usage counter and returns the sloth content
// This is called when a sloth is used in a run command
func (s *SlothService) UseSloth(ctx context.Context, name string) (string, error) {
	// Get sloth and check if active
	sl, err := s.GetActiveSloth(ctx, name)
	if err != nil {
		return "", err
	}

	// Increment usage count
	if err := s.repo.IncrementUsage(ctx, name); err != nil {
		// Log error but don't fail the operation
		// This is not critical for the execution
	}

	return sl.Content, nil
}

// UpdateSloth updates an existing sloth from a file
func (s *SlothService) UpdateSloth(ctx context.Context, name, filePath, description string) error {
	// Get existing sloth
	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return err
	}

	// Read new file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate new file hash
	hash := sha256.Sum256(content)
	fileHash := fmt.Sprintf("%x", hash)

	// Update fields
	existing.FilePath = filePath
	existing.Content = string(content)
	existing.FileHash = fileHash

	if description != "" {
		existing.Description = description
	}

	// Save changes
	if err := s.repo.Update(ctx, existing); err != nil {
		return fmt.Errorf("failed to update sloth: %w", err)
	}

	return nil
}

// WriteContentToFile writes sloth content to a temporary file
// Returns the temporary file path
func (s *SlothService) WriteContentToFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "sloth-*.sloth")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.WriteString(tmpFile, content); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write content: %w", err)
	}

	return tmpFile.Name(), nil
}

// Close closes the service and underlying repository
func (s *SlothService) Close() error {
	return s.repo.Close()
}
