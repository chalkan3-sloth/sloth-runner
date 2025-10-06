package sloth

import (
	"context"
)

// MockRepository is a mock implementation of the Repository interface for testing
type MockRepository struct {
	CreateFunc         func(ctx context.Context, sloth *Sloth) error
	GetByNameFunc      func(ctx context.Context, name string) (*Sloth, error)
	GetByIDFunc        func(ctx context.Context, id string) (*Sloth, error)
	ListFunc           func(ctx context.Context, activeOnly bool) ([]*SlothListItem, error)
	UpdateFunc         func(ctx context.Context, sloth *Sloth) error
	DeleteFunc         func(ctx context.Context, name string) error
	SetActiveFunc      func(ctx context.Context, name string, active bool) error
	IncrementUsageFunc func(ctx context.Context, name string) error
	CloseFunc          func() error
}

func (m *MockRepository) Create(ctx context.Context, sloth *Sloth) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, sloth)
	}
	return nil
}

func (m *MockRepository) GetByName(ctx context.Context, name string) (*Sloth, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, name)
	}
	return nil, ErrSlothNotFound
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*Sloth, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, ErrSlothNotFound
}

func (m *MockRepository) List(ctx context.Context, activeOnly bool) ([]*SlothListItem, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, activeOnly)
	}
	return []*SlothListItem{}, nil
}

func (m *MockRepository) Update(ctx context.Context, sloth *Sloth) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, sloth)
	}
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, name)
	}
	return nil
}

func (m *MockRepository) SetActive(ctx context.Context, name string, active bool) error {
	if m.SetActiveFunc != nil {
		return m.SetActiveFunc(ctx, name, active)
	}
	return nil
}

func (m *MockRepository) IncrementUsage(ctx context.Context, name string) error {
	if m.IncrementUsageFunc != nil {
		return m.IncrementUsageFunc(ctx, name)
	}
	return nil
}

func (m *MockRepository) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}
