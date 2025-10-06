package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/sloth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlothService_GetSloth(t *testing.T) {
	mockRepo := &sloth.MockRepository{
		GetByNameFunc: func(ctx context.Context, name string) (*sloth.Sloth, error) {
			if name == "test-sloth" {
				return &sloth.Sloth{
					ID:          "test-id",
					Name:        "test-sloth",
					Description: "Test sloth",
					IsActive:    true,
					Content:     "test content",
				}, nil
			}
			return nil, sloth.ErrSlothNotFound
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	result, err := service.GetSloth(context.Background(), "test-sloth")
	require.NoError(t, err)
	assert.Equal(t, "test-sloth", result.Name)
	assert.Equal(t, "test content", result.Content)
}

func TestSlothService_GetActiveSloth(t *testing.T) {
	tests := []struct {
		name        string
		slothName   string
		mockSloth   *sloth.Sloth
		mockError   error
		expectError bool
		expectedErr error
	}{
		{
			name:      "active sloth",
			slothName: "active-sloth",
			mockSloth: &sloth.Sloth{
				Name:     "active-sloth",
				IsActive: true,
				Content:  "active content",
			},
			expectError: false,
		},
		{
			name:      "inactive sloth",
			slothName: "inactive-sloth",
			mockSloth: &sloth.Sloth{
				Name:     "inactive-sloth",
				IsActive: false,
				Content:  "inactive content",
			},
			expectError: true,
			expectedErr: sloth.ErrSlothInactive,
		},
		{
			name:        "sloth not found",
			slothName:   "nonexistent",
			mockError:   sloth.ErrSlothNotFound,
			expectError: true,
			expectedErr: sloth.ErrSlothNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &sloth.MockRepository{
				GetByNameFunc: func(ctx context.Context, name string) (*sloth.Sloth, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockSloth, nil
				},
			}

			service := NewSlothServiceWithRepository(mockRepo)
			defer service.Close()

			result, err := service.GetActiveSloth(context.Background(), tt.slothName)

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, result.IsActive)
			}
		})
	}
}

func TestSlothService_ListSloths(t *testing.T) {
	mockRepo := &sloth.MockRepository{
		ListFunc: func(ctx context.Context, activeOnly bool) ([]*sloth.SlothListItem, error) {
			items := []*sloth.SlothListItem{
				{Name: "sloth1", IsActive: true, UsageCount: 5},
				{Name: "sloth2", IsActive: false, UsageCount: 2},
			}
			if activeOnly {
				return []*sloth.SlothListItem{items[0]}, nil
			}
			return items, nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	// Test listing all sloths
	all, err := service.ListSloths(context.Background(), false)
	require.NoError(t, err)
	assert.Len(t, all, 2)

	// Test listing active only
	active, err := service.ListSloths(context.Background(), true)
	require.NoError(t, err)
	assert.Len(t, active, 1)
	assert.Equal(t, "sloth1", active[0].Name)
}

func TestSlothService_RemoveSloth(t *testing.T) {
	deleted := ""
	mockRepo := &sloth.MockRepository{
		DeleteFunc: func(ctx context.Context, name string) error {
			deleted = name
			return nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	err := service.RemoveSloth(context.Background(), "test-sloth")
	require.NoError(t, err)
	assert.Equal(t, "test-sloth", deleted)
}

func TestSlothService_ActivateSloth(t *testing.T) {
	var activatedName string
	var activatedStatus bool

	mockRepo := &sloth.MockRepository{
		SetActiveFunc: func(ctx context.Context, name string, active bool) error {
			activatedName = name
			activatedStatus = active
			return nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	err := service.ActivateSloth(context.Background(), "test-sloth")
	require.NoError(t, err)
	assert.Equal(t, "test-sloth", activatedName)
	assert.True(t, activatedStatus)
}

func TestSlothService_DeactivateSloth(t *testing.T) {
	var deactivatedName string
	var deactivatedStatus bool

	mockRepo := &sloth.MockRepository{
		SetActiveFunc: func(ctx context.Context, name string, active bool) error {
			deactivatedName = name
			deactivatedStatus = active
			return nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	err := service.DeactivateSloth(context.Background(), "test-sloth")
	require.NoError(t, err)
	assert.Equal(t, "test-sloth", deactivatedName)
	assert.False(t, deactivatedStatus)
}

func TestSlothService_UseSloth(t *testing.T) {
	usageIncremented := false

	mockRepo := &sloth.MockRepository{
		GetByNameFunc: func(ctx context.Context, name string) (*sloth.Sloth, error) {
			return &sloth.Sloth{
				Name:     "test-sloth",
				IsActive: true,
				Content:  "test content for use",
			}, nil
		},
		IncrementUsageFunc: func(ctx context.Context, name string) error {
			usageIncremented = true
			return nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	content, err := service.UseSloth(context.Background(), "test-sloth")
	require.NoError(t, err)
	assert.Equal(t, "test content for use", content)
	assert.True(t, usageIncremented)
}

func TestSlothService_AddSloth(t *testing.T) {
	// Create a temporary test file
	tmpFile, err := os.CreateTemp("", "test-sloth-*.sloth")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	testContent := "task test { command = \"echo test\" }"
	_, err = tmpFile.WriteString(testContent)
	require.NoError(t, err)
	tmpFile.Close()

	var createdSloth *sloth.Sloth
	mockRepo := &sloth.MockRepository{
		CreateFunc: func(ctx context.Context, s *sloth.Sloth) error {
			createdSloth = s
			return nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	err = service.AddSloth(context.Background(), "test-sloth", tmpFile.Name(), "Test description", true)
	require.NoError(t, err)
	require.NotNil(t, createdSloth)
	assert.Equal(t, "test-sloth", createdSloth.Name)
	assert.Equal(t, "Test description", createdSloth.Description)
	assert.Equal(t, testContent, createdSloth.Content)
	assert.True(t, createdSloth.IsActive)
	assert.NotEmpty(t, createdSloth.FileHash)
}

func TestSlothService_AddSloth_FileNotFound(t *testing.T) {
	mockRepo := &sloth.MockRepository{}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	err := service.AddSloth(context.Background(), "test-sloth", "/nonexistent/file.sloth", "Test", true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestSlothService_UpdateSloth(t *testing.T) {
	// Create a temporary test file with new content
	tmpFile, err := os.CreateTemp("", "test-sloth-*.sloth")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	newContent := "task updated { command = \"echo updated\" }"
	_, err = tmpFile.WriteString(newContent)
	require.NoError(t, err)
	tmpFile.Close()

	existingSloth := &sloth.Sloth{
		ID:          "test-id",
		Name:        "test-sloth",
		Description: "Old description",
		Content:     "old content",
		IsActive:    true,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
		FileHash:    "oldhash",
	}

	var updatedSloth *sloth.Sloth
	mockRepo := &sloth.MockRepository{
		GetByNameFunc: func(ctx context.Context, name string) (*sloth.Sloth, error) {
			return existingSloth, nil
		},
		UpdateFunc: func(ctx context.Context, s *sloth.Sloth) error {
			updatedSloth = s
			return nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	err = service.UpdateSloth(context.Background(), "test-sloth", tmpFile.Name(), "Updated description")
	require.NoError(t, err)
	require.NotNil(t, updatedSloth)
	assert.Equal(t, newContent, updatedSloth.Content)
	assert.Equal(t, "Updated description", updatedSloth.Description)
	assert.NotEqual(t, "oldhash", updatedSloth.FileHash)
	assert.NotEmpty(t, updatedSloth.FileHash)
}

func TestSlothService_WriteContentToFile(t *testing.T) {
	mockRepo := &sloth.MockRepository{}
	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	testContent := "task write_test { command = \"echo write test\" }"

	filePath, err := service.WriteContentToFile(testContent)
	require.NoError(t, err)
	defer os.Remove(filePath)

	// Verify file exists and has correct content
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, testContent, string(content))

	// Verify file has .sloth extension
	assert.Contains(t, filePath, ".sloth")
}

func TestSlothService_DeleteSloth(t *testing.T) {
	deleted := ""
	mockRepo := &sloth.MockRepository{
		DeleteFunc: func(ctx context.Context, name string) error {
			deleted = name
			return nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	defer service.Close()

	// DeleteSloth should be an alias for RemoveSloth
	err := service.DeleteSloth(context.Background(), "test-sloth")
	require.NoError(t, err)
	assert.Equal(t, "test-sloth", deleted)
}

func TestSlothService_Close(t *testing.T) {
	closeCalled := false
	mockRepo := &sloth.MockRepository{
		CloseFunc: func() error {
			closeCalled = true
			return nil
		},
	}

	service := NewSlothServiceWithRepository(mockRepo)
	err := service.Close()

	require.NoError(t, err)
	assert.True(t, closeCalled)
}
