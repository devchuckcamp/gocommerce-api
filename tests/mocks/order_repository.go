package mocks

import (
	"context"

	"github.com/devchuckcamp/gocommerce/orders"
)

// MockOrderRepository is a mock implementation of orders.Repository
type MockOrderRepository struct {
	Orders map[string]*orders.Order

	// Error injection
	FindByIDError          error
	FindByOrderNumberError error
	FindByUserIDError      error
	SaveError              error
	DeleteError            error
}

// NewMockOrderRepository creates a new mock order repository
func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		Orders: make(map[string]*orders.Order),
	}
}

// FindByID returns an order by ID
func (m *MockOrderRepository) FindByID(ctx context.Context, id string) (*orders.Order, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}
	if o, ok := m.Orders[id]; ok {
		return o, nil
	}
	return nil, orders.ErrOrderNotFound
}

// FindByOrderNumber returns an order by order number
func (m *MockOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*orders.Order, error) {
	if m.FindByOrderNumberError != nil {
		return nil, m.FindByOrderNumberError
	}
	for _, o := range m.Orders {
		if o.OrderNumber == orderNumber {
			return o, nil
		}
	}
	return nil, orders.ErrOrderNotFound
}

// FindByUserID returns orders by user ID
func (m *MockOrderRepository) FindByUserID(ctx context.Context, userID string, filter orders.OrderFilter) ([]*orders.Order, error) {
	if m.FindByUserIDError != nil {
		return nil, m.FindByUserIDError
	}
	result := make([]*orders.Order, 0)
	for _, o := range m.Orders {
		if o.UserID == userID {
			result = append(result, o)
		}
	}
	return result, nil
}

// Save saves an order
func (m *MockOrderRepository) Save(ctx context.Context, o *orders.Order) error {
	if m.SaveError != nil {
		return m.SaveError
	}
	m.Orders[o.ID] = o
	return nil
}

// Delete deletes an order
func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Orders, id)
	return nil
}
