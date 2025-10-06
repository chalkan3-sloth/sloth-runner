package sloth

import (
	"time"
)

// Sloth represents a saved .sloth file in the database
type Sloth struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FilePath    string    `json:"file_path"`
	Content     string    `json:"content"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	UsageCount  int       `json:"usage_count"`
	Tags        string    `json:"tags"` // Comma-separated tags
	FileHash    string    `json:"file_hash"` // SHA256 hash of content
}

// SlothListItem represents a simplified view for listing
type SlothListItem struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	UsageCount  int       `json:"usage_count"`
}
