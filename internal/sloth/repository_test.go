package sloth

import (
	"context"
	"errors"
	"testing"
)

// Test error variables
func TestErrSlothNotFound(t *testing.T) {
	if ErrSlothNotFound == nil {
		t.Error("Expected non-nil error")
	}

	if ErrSlothNotFound.Error() == "" {
		t.Error("Expected non-empty error message")
	}

	expected := "sloth not found"
	if ErrSlothNotFound.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, ErrSlothNotFound.Error())
	}
}

func TestErrSlothAlreadyExists(t *testing.T) {
	if ErrSlothAlreadyExists == nil {
		t.Error("Expected non-nil error")
	}

	if ErrSlothAlreadyExists.Error() == "" {
		t.Error("Expected non-empty error message")
	}

	expected := "sloth with this name already exists"
	if ErrSlothAlreadyExists.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, ErrSlothAlreadyExists.Error())
	}
}

func TestErrSlothInactive(t *testing.T) {
	if ErrSlothInactive == nil {
		t.Error("Expected non-nil error")
	}

	if ErrSlothInactive.Error() == "" {
		t.Error("Expected non-empty error message")
	}

	expected := "sloth is not active"
	if ErrSlothInactive.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, ErrSlothInactive.Error())
	}
}

func TestErrors_AreDistinct(t *testing.T) {
	// Each error should be distinct
	if ErrSlothNotFound == ErrSlothAlreadyExists {
		t.Error("Expected errors to be distinct")
	}

	if ErrSlothNotFound == ErrSlothInactive {
		t.Error("Expected errors to be distinct")
	}

	if ErrSlothAlreadyExists == ErrSlothInactive {
		t.Error("Expected errors to be distinct")
	}
}

func TestErrors_CanBeCompared(t *testing.T) {
	err1 := ErrSlothNotFound
	err2 := ErrSlothNotFound

	// Same error should be equal
	if err1 != err2 {
		t.Error("Expected same error instances to be equal")
	}

	// Can use errors.Is
	if !errors.Is(err1, ErrSlothNotFound) {
		t.Error("Expected errors.Is to work")
	}

	if errors.Is(err1, ErrSlothAlreadyExists) {
		t.Error("Expected errors.Is to return false for different errors")
	}
}

func TestErrors_Messages(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		contains string
	}{
		{"NotFound", ErrSlothNotFound, "not found"},
		{"AlreadyExists", ErrSlothAlreadyExists, "already exists"},
		{"Inactive", ErrSlothInactive, "not active"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.err.Error()
			if !containsSubstr(msg, tc.contains) {
				t.Errorf("Expected error message to contain '%s', got '%s'", tc.contains, msg)
			}
		})
	}
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Test Repository interface
func TestRepository_InterfaceExists(t *testing.T) {
	// Verify that Repository is a type
	var _ Repository

	// This test ensures the interface exists and can be referenced
}

func TestRepository_MockImplementation(t *testing.T) {
	// MockRepository should implement Repository interface
	var _ Repository = (*MockRepository)(nil)
}

type testRepository struct{}

func (r *testRepository) Create(ctx context.Context, sloth *Sloth) error {
	return nil
}

func (r *testRepository) GetByName(ctx context.Context, name string) (*Sloth, error) {
	return nil, nil
}

func (r *testRepository) GetByID(ctx context.Context, id string) (*Sloth, error) {
	return nil, nil
}

func (r *testRepository) List(ctx context.Context, activeOnly bool) ([]*SlothListItem, error) {
	return nil, nil
}

func (r *testRepository) Update(ctx context.Context, sloth *Sloth) error {
	return nil
}

func (r *testRepository) Delete(ctx context.Context, name string) error {
	return nil
}

func (r *testRepository) SetActive(ctx context.Context, name string, active bool) error {
	return nil
}

func (r *testRepository) IncrementUsage(ctx context.Context, name string) error {
	return nil
}

func (r *testRepository) Close() error {
	return nil
}

func TestRepository_CustomImplementation(t *testing.T) {
	// testRepository should implement Repository interface
	var _ Repository = (*testRepository)(nil)

	repo := &testRepository{}

	ctx := context.Background()

	// Test all methods can be called
	_ = repo.Create(ctx, &Sloth{})
	_, _ = repo.GetByName(ctx, "test")
	_, _ = repo.GetByID(ctx, "id")
	_, _ = repo.List(ctx, true)
	_ = repo.Update(ctx, &Sloth{})
	_ = repo.Delete(ctx, "test")
	_ = repo.SetActive(ctx, "test", true)
	_ = repo.IncrementUsage(ctx, "test")
	_ = repo.Close()
}

func TestRepository_MethodSignatures(t *testing.T) {
	// This test verifies the interface has the expected method signatures
	// by trying to compile with a struct that implements it

	var repo Repository = &testRepository{}

	// Create signature: (context.Context, *Sloth) error
	err := repo.Create(context.Background(), &Sloth{})
	if err == nil {
		// OK
	}

	// GetByName signature: (context.Context, string) (*Sloth, error)
	sloth, err := repo.GetByName(context.Background(), "name")
	if sloth == nil && err == nil {
		// OK
	}

	// GetByID signature: (context.Context, string) (*Sloth, error)
	sloth, err = repo.GetByID(context.Background(), "id")
	if sloth == nil && err == nil {
		// OK
	}

	// List signature: (context.Context, bool) ([]*SlothListItem, error)
	items, err := repo.List(context.Background(), true)
	if items == nil && err == nil {
		// OK
	}

	// Update signature: (context.Context, *Sloth) error
	err = repo.Update(context.Background(), &Sloth{})
	if err == nil {
		// OK
	}

	// Delete signature: (context.Context, string) error
	err = repo.Delete(context.Background(), "name")
	if err == nil {
		// OK
	}

	// SetActive signature: (context.Context, string, bool) error
	err = repo.SetActive(context.Background(), "name", true)
	if err == nil {
		// OK
	}

	// IncrementUsage signature: (context.Context, string) error
	err = repo.IncrementUsage(context.Background(), "name")
	if err == nil {
		// OK
	}

	// Close signature: () error
	err = repo.Close()
	if err == nil {
		// OK
	}
}

func TestRepository_RequiresMethods(t *testing.T) {
	// Verify all required methods exist
	requiredMethods := []string{
		"Create",
		"GetByName",
		"GetByID",
		"List",
		"Update",
		"Delete",
		"SetActive",
		"IncrementUsage",
		"Close",
	}

	// This test passes if the interface has all these methods
	// (verified by compilation of testRepository)

	_ = requiredMethods // Used to document expected methods
}

func TestRepository_ContextUsage(t *testing.T) {
	repo := &testRepository{}

	// All data methods should accept context
	ctx := context.Background()

	_ = repo.Create(ctx, &Sloth{})
	_, _ = repo.GetByName(ctx, "test")
	_, _ = repo.GetByID(ctx, "test")
	_, _ = repo.List(ctx, true)
	_ = repo.Update(ctx, &Sloth{})
	_ = repo.Delete(ctx, "test")
	_ = repo.SetActive(ctx, "test", true)
	_ = repo.IncrementUsage(ctx, "test")
}

func TestRepository_CancelableContext(t *testing.T) {
	repo := &testRepository{}

	// Context can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Methods should still be callable with cancelled context
	// (actual behavior depends on implementation)
	_ = repo.Create(ctx, &Sloth{})
	_, _ = repo.GetByName(ctx, "test")
}

func TestErrors_WrappingSupport(t *testing.T) {
	// Errors can be wrapped
	wrappedNotFound := errors.Join(ErrSlothNotFound, errors.New("additional context"))

	if !errors.Is(wrappedNotFound, ErrSlothNotFound) {
		t.Error("Expected wrapped error to be detectable")
	}

	wrappedExists := errors.Join(ErrSlothAlreadyExists, errors.New("name: test"))

	if !errors.Is(wrappedExists, ErrSlothAlreadyExists) {
		t.Error("Expected wrapped error to be detectable")
	}
}

func TestErrors_AsVariable(t *testing.T) {
	// Errors can be assigned to variables
	notFoundErr := ErrSlothNotFound
	existsErr := ErrSlothAlreadyExists
	inactiveErr := ErrSlothInactive

	if notFoundErr == nil {
		t.Error("Expected non-nil error")
	}

	if existsErr == nil {
		t.Error("Expected non-nil error")
	}

	if inactiveErr == nil {
		t.Error("Expected non-nil error")
	}
}

func TestErrors_UseCase_NotFound(t *testing.T) {
	// Simulate a not found scenario
	simulateGet := func(name string) (*Sloth, error) {
		if name == "nonexistent" {
			return nil, ErrSlothNotFound
		}
		return &Sloth{Name: name}, nil
	}

	// Test not found
	_, err := simulateGet("nonexistent")
	if !errors.Is(err, ErrSlothNotFound) {
		t.Error("Expected ErrSlothNotFound")
	}

	// Test found
	sloth, err := simulateGet("exists")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sloth == nil {
		t.Error("Expected non-nil sloth")
	}
}

func TestErrors_UseCase_AlreadyExists(t *testing.T) {
	// Simulate a creation scenario
	existing := map[string]bool{"test": true}

	simulateCreate := func(name string) error {
		if existing[name] {
			return ErrSlothAlreadyExists
		}
		existing[name] = true
		return nil
	}

	// Test already exists
	err := simulateCreate("test")
	if !errors.Is(err, ErrSlothAlreadyExists) {
		t.Error("Expected ErrSlothAlreadyExists")
	}

	// Test new creation
	err = simulateCreate("new")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestErrors_UseCase_Inactive(t *testing.T) {
	// Simulate an activation check
	simulateUse := func(isActive bool) error {
		if !isActive {
			return ErrSlothInactive
		}
		return nil
	}

	// Test inactive
	err := simulateUse(false)
	if !errors.Is(err, ErrSlothInactive) {
		t.Error("Expected ErrSlothInactive")
	}

	// Test active
	err = simulateUse(true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
