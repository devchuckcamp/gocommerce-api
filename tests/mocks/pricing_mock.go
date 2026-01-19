package mocks

import (
	"context"
	"time"

	"github.com/devchuckcamp/gocommerce/money"
	"github.com/devchuckcamp/gocommerce/pricing"
)

// MockSalePriceResolver is a mock implementation of services.SalePriceResolver
type MockSalePriceResolver struct {
	Prices map[string]*pricing.ProductPrice

	FindEffectivePriceError  error
	FindEffectivePricesError error
}

// NewMockSalePriceResolver creates a new mock sale price resolver
func NewMockSalePriceResolver() *MockSalePriceResolver {
	return &MockSalePriceResolver{
		Prices: make(map[string]*pricing.ProductPrice),
	}
}

// FindEffectivePrice returns the effective price for a product
func (m *MockSalePriceResolver) FindEffectivePrice(ctx context.Context, productID string, variantID *string, at time.Time) (*pricing.ProductPrice, error) {
	if m.FindEffectivePriceError != nil {
		return nil, m.FindEffectivePriceError
	}
	if price, ok := m.Prices[productID]; ok {
		return price, nil
	}
	return nil, nil
}

// FindEffectivePrices returns effective prices for multiple products
func (m *MockSalePriceResolver) FindEffectivePrices(ctx context.Context, productIDs []string, at time.Time) (map[string]*pricing.ProductPrice, error) {
	if m.FindEffectivePricesError != nil {
		return nil, m.FindEffectivePricesError
	}
	result := make(map[string]*pricing.ProductPrice)
	for _, id := range productIDs {
		if price, ok := m.Prices[id]; ok {
			result[id] = price
		}
	}
	return result, nil
}

// AddPrice adds a price for a product
func (m *MockSalePriceResolver) AddPrice(productID string, amount int64, currency string) {
	m.Prices[productID] = &pricing.ProductPrice{
		ProductID: productID,
		Price:     money.Money{Amount: amount, Currency: currency},
	}
}

// MockPromotionRepository is a mock implementation of pricing.PromotionRepository
type MockPromotionRepository struct {
	Promotions []*pricing.Promotion

	FindActiveError error
	FindByCodeError error
}

// NewMockPromotionRepository creates a new mock promotion repository
func NewMockPromotionRepository() *MockPromotionRepository {
	return &MockPromotionRepository{
		Promotions: make([]*pricing.Promotion, 0),
	}
}

// FindActive returns active promotions
func (m *MockPromotionRepository) FindActive(ctx context.Context, at time.Time) ([]*pricing.Promotion, error) {
	if m.FindActiveError != nil {
		return nil, m.FindActiveError
	}
	return m.Promotions, nil
}

// FindByCode returns a promotion by code
func (m *MockPromotionRepository) FindByCode(ctx context.Context, code string) (*pricing.Promotion, error) {
	if m.FindByCodeError != nil {
		return nil, m.FindByCodeError
	}
	for _, p := range m.Promotions {
		if p.Code == code {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// FindByID returns a promotion by ID
func (m *MockPromotionRepository) FindByID(ctx context.Context, id string) (*pricing.Promotion, error) {
	for _, p := range m.Promotions {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// Create creates a new promotion
func (m *MockPromotionRepository) Create(ctx context.Context, promotion *pricing.Promotion) error {
	m.Promotions = append(m.Promotions, promotion)
	return nil
}

// Update updates a promotion
func (m *MockPromotionRepository) Update(ctx context.Context, promotion *pricing.Promotion) error {
	for i, p := range m.Promotions {
		if p.ID == promotion.ID {
			m.Promotions[i] = promotion
			return nil
		}
	}
	return ErrNotFound
}

// Delete deletes a promotion
func (m *MockPromotionRepository) Delete(ctx context.Context, id string) error {
	for i, p := range m.Promotions {
		if p.ID == id {
			m.Promotions = append(m.Promotions[:i], m.Promotions[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
