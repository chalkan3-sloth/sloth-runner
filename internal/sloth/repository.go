package sloth

import (
	"context"
	"errors"
)

var (
	// ErrSlothNotFound is returned when a sloth is not found
	ErrSlothNotFound = errors.New("sloth not found")

	// ErrSlothAlreadyExists is returned when attempting to create a sloth with an existing name
	ErrSlothAlreadyExists = errors.New("sloth with this name already exists")

	// ErrSlothInactive is returned when attempting to use an inactive sloth
	ErrSlothInactive = errors.New("sloth is not active")
)

// Repository defines the interface for sloth persistence
// This follows the Repository Pattern for data access abstraction
type Repository interface {
	// Create adds a new sloth to the repository
	Create(ctx context.Context, sloth *Sloth) error

	// GetByName retrieves a sloth by its name
	GetByName(ctx context.Context, name string) (*Sloth, error)

	// GetByID retrieves a sloth by its ID
	GetByID(ctx context.Context, id string) (*Sloth, error)

	// List returns all sloths, optionally filtered by active status
	List(ctx context.Context, activeOnly bool) ([]*SlothListItem, error)

	// Update updates an existing sloth
	Update(ctx context.Context, sloth *Sloth) error

	// Delete removes a sloth by name
	Delete(ctx context.Context, name string) error

	// SetActive sets the active status of a sloth
	SetActive(ctx context.Context, name string, active bool) error

	// IncrementUsage increments the usage count and updates last used time
	IncrementUsage(ctx context.Context, name string) error

	// Close closes the repository connection
	Close() error
}
