package mocks

import (
	"context"

	"github.com/devchuckcamp/gocommerce/cart"
)

// MockCartRepository is a mock implementation of cart.Repository
type MockCartRepository struct {
	Carts map[string]*cart.Cart

	// Error injection
	FindByIDError        error
	FindByUserIDError    error
	FindBySessionIDError error
	SaveError            error
	DeleteError          error
}

// NewMockCartRepository creates a new mock cart repository
func NewMockCartRepository() *MockCartRepository {
	return &MockCartRepository{
		Carts: make(map[string]*cart.Cart),
	}
}

// FindByID returns a cart by ID
func (m *MockCartRepository) FindByID(ctx context.Context, id string) (*cart.Cart, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}
	if c, ok := m.Carts[id]; ok {
		return c, nil
	}
	return nil, cart.ErrCartNotFound
}

// FindByUserID returns a cart by user ID
func (m *MockCartRepository) FindByUserID(ctx context.Context, userID string) (*cart.Cart, error) {
	if m.FindByUserIDError != nil {
		return nil, m.FindByUserIDError
	}
	for _, c := range m.Carts {
		if c.UserID == userID {
			return c, nil
		}
	}
	return nil, cart.ErrCartNotFound
}

// FindBySessionID returns a cart by session ID
func (m *MockCartRepository) FindBySessionID(ctx context.Context, sessionID string) (*cart.Cart, error) {
	if m.FindBySessionIDError != nil {
		return nil, m.FindBySessionIDError
	}
	for _, c := range m.Carts {
		if c.SessionID == sessionID {
			return c, nil
		}
	}
	return nil, cart.ErrCartNotFound
}

// Save saves a cart
func (m *MockCartRepository) Save(ctx context.Context, c *cart.Cart) error {
	if m.SaveError != nil {
		return m.SaveError
	}
	m.Carts[c.ID] = c
	return nil
}

// Delete deletes a cart
func (m *MockCartRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Carts, id)
	return nil
}
