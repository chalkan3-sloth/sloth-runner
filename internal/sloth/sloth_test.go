package sloth

import (
	"encoding/json"
	"testing"
	"time"
)

// Test Sloth struct
func TestSloth_StructCreation(t *testing.T) {
	now := time.Now()
	sloth := &Sloth{
		ID:          "test-id",
		Name:        "test-sloth",
		Description: "Test description",
		FilePath:    "/path/to/file.sloth",
		Content:     "content",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		UsageCount:  5,
		Tags:        "tag1,tag2",
		FileHash:    "abc123",
	}

	if sloth.ID != "test-id" {
		t.Error("Expected ID to be set")
	}

	if sloth.Name != "test-sloth" {
		t.Error("Expected Name to be set")
	}

	if !sloth.IsActive {
		t.Error("Expected IsActive to be true")
	}

	if sloth.UsageCount != 5 {
		t.Error("Expected UsageCount to be 5")
	}
}

func TestSloth_JSONMarshaling(t *testing.T) {
	now := time.Now()
	sloth := &Sloth{
		ID:          "test-id",
		Name:        "test-sloth",
		Description: "Test description",
		FilePath:    "/path/to/file.sloth",
		Content:     "content",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		UsageCount:  5,
		Tags:        "tag1,tag2",
		FileHash:    "abc123",
	}

	data, err := json.Marshal(sloth)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled Sloth
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.ID != sloth.ID {
		t.Errorf("Expected ID %s, got %s", sloth.ID, unmarshaled.ID)
	}

	if unmarshaled.Name != sloth.Name {
		t.Errorf("Expected Name %s, got %s", sloth.Name, unmarshaled.Name)
	}

	if unmarshaled.UsageCount != sloth.UsageCount {
		t.Errorf("Expected UsageCount %d, got %d", sloth.UsageCount, unmarshaled.UsageCount)
	}
}

func TestSloth_LastUsedAtOptional(t *testing.T) {
	// Test without LastUsedAt
	sloth1 := &Sloth{
		ID:   "test-1",
		Name: "test-sloth-1",
	}

	data, err := json.Marshal(sloth1)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// LastUsedAt should be omitted when nil
	if contains(string(data), "last_used_at") {
		t.Error("Expected last_used_at to be omitted when nil")
	}

	// Test with LastUsedAt
	now := time.Now()
	sloth2 := &Sloth{
		ID:         "test-2",
		Name:       "test-sloth-2",
		LastUsedAt: &now,
	}

	data, err = json.Marshal(sloth2)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// LastUsedAt should be included when set
	if !contains(string(data), "last_used_at") {
		t.Error("Expected last_used_at to be included when set")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSloth_DefaultValues(t *testing.T) {
	sloth := &Sloth{}

	if sloth.ID != "" {
		t.Error("Expected default ID to be empty")
	}

	if sloth.IsActive {
		t.Error("Expected default IsActive to be false")
	}

	if sloth.UsageCount != 0 {
		t.Error("Expected default UsageCount to be 0")
	}

	if sloth.LastUsedAt != nil {
		t.Error("Expected default LastUsedAt to be nil")
	}
}

func TestSloth_Tags(t *testing.T) {
	sloth := &Sloth{
		Tags: "production,deployment,web",
	}

	if sloth.Tags != "production,deployment,web" {
		t.Errorf("Expected tags to be preserved: %s", sloth.Tags)
	}

	// Empty tags
	sloth2 := &Sloth{}
	if sloth2.Tags != "" {
		t.Error("Expected empty tags by default")
	}
}

func TestSloth_FileHash(t *testing.T) {
	sloth := &Sloth{
		FileHash: "sha256:abcdef1234567890",
	}

	if sloth.FileHash != "sha256:abcdef1234567890" {
		t.Errorf("Expected FileHash to be preserved: %s", sloth.FileHash)
	}
}

func TestSloth_TimestampFields(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)
	lastUsedAt := updatedAt.Add(1 * time.Hour)

	sloth := &Sloth{
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		LastUsedAt: &lastUsedAt,
	}

	if sloth.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if sloth.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	if sloth.LastUsedAt == nil {
		t.Fatal("Expected LastUsedAt to be set")
	}

	if sloth.LastUsedAt.Before(updatedAt) {
		t.Error("Expected LastUsedAt to be after UpdatedAt")
	}
}

// Test SlothListItem struct
func TestSlothListItem_StructCreation(t *testing.T) {
	now := time.Now()
	item := &SlothListItem{
		Name:        "test-sloth",
		Description: "Test description",
		IsActive:    true,
		CreatedAt:   now,
		UsageCount:  10,
	}

	if item.Name != "test-sloth" {
		t.Error("Expected Name to be set")
	}

	if item.Description != "Test description" {
		t.Error("Expected Description to be set")
	}

	if !item.IsActive {
		t.Error("Expected IsActive to be true")
	}

	if item.UsageCount != 10 {
		t.Error("Expected UsageCount to be 10")
	}
}

func TestSlothListItem_JSONMarshaling(t *testing.T) {
	now := time.Now()
	item := &SlothListItem{
		Name:        "test-sloth",
		Description: "Test description",
		IsActive:    true,
		CreatedAt:   now,
		UsageCount:  10,
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled SlothListItem
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.Name != item.Name {
		t.Errorf("Expected Name %s, got %s", item.Name, unmarshaled.Name)
	}

	if unmarshaled.UsageCount != item.UsageCount {
		t.Errorf("Expected UsageCount %d, got %d", item.UsageCount, unmarshaled.UsageCount)
	}
}

func TestSlothListItem_LastUsedAtOptional(t *testing.T) {
	// Without LastUsedAt
	item1 := &SlothListItem{
		Name: "test-1",
	}

	data, err := json.Marshal(item1)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	if contains(string(data), "last_used_at") {
		t.Error("Expected last_used_at to be omitted when nil")
	}

	// With LastUsedAt
	now := time.Now()
	item2 := &SlothListItem{
		Name:       "test-2",
		LastUsedAt: &now,
	}

	data, err = json.Marshal(item2)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	if !contains(string(data), "last_used_at") {
		t.Error("Expected last_used_at to be included when set")
	}
}

func TestSlothListItem_DefaultValues(t *testing.T) {
	item := &SlothListItem{}

	if item.Name != "" {
		t.Error("Expected default Name to be empty")
	}

	if item.IsActive {
		t.Error("Expected default IsActive to be false")
	}

	if item.UsageCount != 0 {
		t.Error("Expected default UsageCount to be 0")
	}

	if item.LastUsedAt != nil {
		t.Error("Expected default LastUsedAt to be nil")
	}
}

func TestSlothListItem_SimplifiedView(t *testing.T) {
	// SlothListItem should have fewer fields than Sloth
	// It's a simplified view for listing

	fullSloth := &Sloth{
		ID:          "id-123",
		Name:        "test-sloth",
		Description: "description",
		FilePath:    "/path/to/file",
		Content:     "large content here...",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UsageCount:  5,
		Tags:        "tag1,tag2",
		FileHash:    "hash123",
	}

	// Convert to list item (manually for testing)
	listItem := &SlothListItem{
		Name:        fullSloth.Name,
		Description: fullSloth.Description,
		IsActive:    fullSloth.IsActive,
		CreatedAt:   fullSloth.CreatedAt,
		LastUsedAt:  fullSloth.LastUsedAt,
		UsageCount:  fullSloth.UsageCount,
	}

	if listItem.Name != fullSloth.Name {
		t.Error("Expected Name to match")
	}

	if listItem.IsActive != fullSloth.IsActive {
		t.Error("Expected IsActive to match")
	}

	if listItem.UsageCount != fullSloth.UsageCount {
		t.Error("Expected UsageCount to match")
	}
}

func TestSloth_ZeroValue(t *testing.T) {
	var sloth Sloth

	// Verify zero values
	if sloth.ID != "" {
		t.Error("Expected zero ID")
	}

	if sloth.IsActive {
		t.Error("Expected zero IsActive")
	}

	if sloth.UsageCount != 0 {
		t.Error("Expected zero UsageCount")
	}

	if !sloth.CreatedAt.IsZero() {
		t.Error("Expected zero CreatedAt")
	}

	if !sloth.UpdatedAt.IsZero() {
		t.Error("Expected zero UpdatedAt")
	}
}

func TestSlothListItem_ZeroValue(t *testing.T) {
	var item SlothListItem

	// Verify zero values
	if item.Name != "" {
		t.Error("Expected zero Name")
	}

	if item.IsActive {
		t.Error("Expected zero IsActive")
	}

	if item.UsageCount != 0 {
		t.Error("Expected zero UsageCount")
	}

	if !item.CreatedAt.IsZero() {
		t.Error("Expected zero CreatedAt")
	}
}

func TestSloth_UsageTracking(t *testing.T) {
	sloth := &Sloth{
		Name:       "test",
		UsageCount: 0,
	}

	// Simulate usage increment
	sloth.UsageCount++
	now := time.Now()
	sloth.LastUsedAt = &now

	if sloth.UsageCount != 1 {
		t.Errorf("Expected UsageCount 1, got %d", sloth.UsageCount)
	}

	if sloth.LastUsedAt == nil {
		t.Error("Expected LastUsedAt to be set")
	}
}

func TestSloth_ActiveStatus(t *testing.T) {
	sloth := &Sloth{
		Name:     "test",
		IsActive: false,
	}

	// Activate
	sloth.IsActive = true
	if !sloth.IsActive {
		t.Error("Expected sloth to be active")
	}

	// Deactivate
	sloth.IsActive = false
	if sloth.IsActive {
		t.Error("Expected sloth to be inactive")
	}
}

func TestSloth_ContentStorage(t *testing.T) {
	content := `
workflow "example" {
  task "hello" {
    command = "echo 'Hello World'"
  }
}
`

	sloth := &Sloth{
		Name:    "test",
		Content: content,
	}

	if sloth.Content != content {
		t.Error("Expected content to be preserved exactly")
	}

	if len(sloth.Content) == 0 {
		t.Error("Expected content to have length")
	}
}

func TestSloth_FilePath(t *testing.T) {
	paths := []string{
		"/path/to/file.sloth",
		"relative/path.sloth",
		"./file.sloth",
		"../parent/file.sloth",
		"C:\\Windows\\path\\file.sloth",
	}

	for _, path := range paths {
		sloth := &Sloth{
			FilePath: path,
		}

		if sloth.FilePath != path {
			t.Errorf("Expected FilePath to be preserved: %s", path)
		}
	}
}

func TestSloth_MultipleInstances(t *testing.T) {
	sloth1 := &Sloth{ID: "1", Name: "sloth-1"}
	sloth2 := &Sloth{ID: "2", Name: "sloth-2"}

	if sloth1.ID == sloth2.ID {
		t.Error("Expected different IDs")
	}

	if sloth1.Name == sloth2.Name {
		t.Error("Expected different Names")
	}

	// Modifying one shouldn't affect the other
	sloth1.IsActive = true
	if sloth2.IsActive {
		t.Error("Expected sloth2 to be unaffected")
	}
}

func TestSlothListItem_Sorting(t *testing.T) {
	now := time.Now()
	item1 := &SlothListItem{
		Name:      "alpha",
		CreatedAt: now,
	}

	item2 := &SlothListItem{
		Name:      "beta",
		CreatedAt: now.Add(1 * time.Hour),
	}

	// Can be sorted by name
	if item1.Name > item2.Name {
		t.Error("Expected item1 to sort before item2 by name")
	}

	// Can be sorted by creation time
	if item1.CreatedAt.After(item2.CreatedAt) {
		t.Error("Expected item1 to sort before item2 by creation time")
	}
}
