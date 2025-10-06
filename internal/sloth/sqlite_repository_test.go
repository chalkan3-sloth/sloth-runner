package sloth

import (
	"context"
	"testing"
	"time"
)

// TestSQLiteRepository_Create tests creating a new sloth
func TestSQLiteRepository_Create(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	sloth := &Sloth{
		ID:          "test-id-1",
		Name:        "test-sloth",
		Description: "Test description",
		FilePath:    "/path/to/test.sloth",
		Content:     "test content",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FileHash:    "abcd1234",
	}

	err = repo.Create(ctx, sloth)
	if err != nil {
		t.Fatalf("Failed to create sloth: %v", err)
	}

	// Verify it was created
	retrieved, err := repo.GetByName(ctx, "test-sloth")
	if err != nil {
		t.Fatalf("Failed to retrieve sloth: %v", err)
	}

	if retrieved.Name != sloth.Name {
		t.Errorf("Expected name %s, got %s", sloth.Name, retrieved.Name)
	}
}

// TestSQLiteRepository_CreateDuplicate tests creating a duplicate sloth
func TestSQLiteRepository_CreateDuplicate(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	sloth1 := &Sloth{
		ID:        "test-id-1",
		Name:      "duplicate",
		Content:   "content1",
		FilePath:  "/path/to/file1.sloth",
		FileHash:  "hash1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	sloth2 := &Sloth{
		ID:        "test-id-2",
		Name:      "duplicate",
		Content:   "content2",
		FilePath:  "/path/to/file2.sloth",
		FileHash:  "hash2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Create(ctx, sloth1)
	if err != nil {
		t.Fatalf("Failed to create first sloth: %v", err)
	}

	err = repo.Create(ctx, sloth2)
	if err != ErrSlothAlreadyExists {
		t.Errorf("Expected ErrSlothAlreadyExists, got %v", err)
	}
}

// TestSQLiteRepository_GetByName tests retrieving a sloth by name
func TestSQLiteRepository_GetByName(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	sloth := &Sloth{
		ID:        "test-id",
		Name:      "findme",
		Content:   "test content",
		FilePath:  "/path/to/test.sloth",
		FileHash:  "hash123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.Create(ctx, sloth)

	retrieved, err := repo.GetByName(ctx, "findme")
	if err != nil {
		t.Fatalf("Failed to get sloth: %v", err)
	}

	if retrieved.Name != "findme" {
		t.Errorf("Expected name 'findme', got %s", retrieved.Name)
	}

	// Test not found
	_, err = repo.GetByName(ctx, "nonexistent")
	if err != ErrSlothNotFound {
		t.Errorf("Expected ErrSlothNotFound, got %v", err)
	}
}

// TestSQLiteRepository_List tests listing sloths
func TestSQLiteRepository_List(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Create active sloth
	active := &Sloth{
		ID:        "active-id",
		Name:      "active-sloth",
		Content:   "active content",
		FilePath:  "/path/to/active.sloth",
		FileHash:  "activehash",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create inactive sloth
	inactive := &Sloth{
		ID:        "inactive-id",
		Name:      "inactive-sloth",
		Content:   "inactive content",
		FilePath:  "/path/to/inactive.sloth",
		FileHash:  "inactivehash",
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.Create(ctx, active)
	repo.Create(ctx, inactive)

	// List all
	all, err := repo.List(ctx, false)
	if err != nil {
		t.Fatalf("Failed to list all: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("Expected 2 sloths, got %d", len(all))
	}

	// List active only
	activeOnly, err := repo.List(ctx, true)
	if err != nil {
		t.Fatalf("Failed to list active: %v", err)
	}
	if len(activeOnly) != 1 {
		t.Errorf("Expected 1 active sloth, got %d", len(activeOnly))
	}
	if activeOnly[0].Name != "active-sloth" {
		t.Errorf("Expected 'active-sloth', got %s", activeOnly[0].Name)
	}
}

// TestSQLiteRepository_Update tests updating a sloth
func TestSQLiteRepository_Update(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	original := &Sloth{
		ID:          "update-id",
		Name:        "update-test",
		Description: "Original description",
		Content:     "original content",
		FilePath:    "/path/to/original.sloth",
		FileHash:    "originalhash",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.Create(ctx, original)

	// Update
	original.Description = "Updated description"
	original.Content = "updated content"

	err = repo.Update(ctx, original)
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	// Verify
	updated, err := repo.GetByName(ctx, "update-test")
	if err != nil {
		t.Fatalf("Failed to get updated sloth: %v", err)
	}

	if updated.Description != "Updated description" {
		t.Errorf("Expected updated description, got %s", updated.Description)
	}
}

// TestSQLiteRepository_Delete tests deleting a sloth
func TestSQLiteRepository_Delete(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	sloth := &Sloth{
		ID:        "delete-id",
		Name:      "delete-me",
		Content:   "content",
		FilePath:  "/path/to/delete.sloth",
		FileHash:  "deletehash",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.Create(ctx, sloth)

	// Delete
	err = repo.Delete(ctx, "delete-me")
	if err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	// Verify deleted
	_, err = repo.GetByName(ctx, "delete-me")
	if err != ErrSlothNotFound {
		t.Errorf("Expected ErrSlothNotFound after delete, got %v", err)
	}

	// Try deleting non-existent
	err = repo.Delete(ctx, "nonexistent")
	if err != ErrSlothNotFound {
		t.Errorf("Expected ErrSlothNotFound for nonexistent, got %v", err)
	}
}

// TestSQLiteRepository_SetActive tests setting active status
func TestSQLiteRepository_SetActive(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	sloth := &Sloth{
		ID:        "active-test-id",
		Name:      "toggle-active",
		Content:   "content",
		FilePath:  "/path/to/test.sloth",
		FileHash:  "hash",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.Create(ctx, sloth)

	// Deactivate
	err = repo.SetActive(ctx, "toggle-active", false)
	if err != nil {
		t.Fatalf("Failed to deactivate: %v", err)
	}

	updated, _ := repo.GetByName(ctx, "toggle-active")
	if updated.IsActive {
		t.Error("Expected sloth to be inactive")
	}

	// Activate
	err = repo.SetActive(ctx, "toggle-active", true)
	if err != nil {
		t.Fatalf("Failed to activate: %v", err)
	}

	updated, _ = repo.GetByName(ctx, "toggle-active")
	if !updated.IsActive {
		t.Error("Expected sloth to be active")
	}
}

// TestSQLiteRepository_IncrementUsage tests incrementing usage count
func TestSQLiteRepository_IncrementUsage(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	sloth := &Sloth{
		ID:         "usage-id",
		Name:       "usage-test",
		Content:    "content",
		FilePath:   "/path/to/test.sloth",
		FileHash:   "hash",
		UsageCount: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	repo.Create(ctx, sloth)

	// Increment once
	err = repo.IncrementUsage(ctx, "usage-test")
	if err != nil {
		t.Fatalf("Failed to increment: %v", err)
	}

	updated, _ := repo.GetByName(ctx, "usage-test")
	if updated.UsageCount != 1 {
		t.Errorf("Expected usage count 1, got %d", updated.UsageCount)
	}

	if updated.LastUsedAt == nil {
		t.Error("Expected LastUsedAt to be set")
	}

	// Increment again
	err = repo.IncrementUsage(ctx, "usage-test")
	if err != nil {
		t.Fatalf("Failed to increment second time: %v", err)
	}

	updated, _ = repo.GetByName(ctx, "usage-test")
	if updated.UsageCount != 2 {
		t.Errorf("Expected usage count 2, got %d", updated.UsageCount)
	}
}
