package fixtures

import (
	"time"

	"github.com/devchuckcamp/gocommerce/cart"
	"github.com/devchuckcamp/gocommerce/money"
)

// Cart fixtures
var (
	// TestUserID is a sample user ID for testing
	TestUserID = "user-001"

	// EmptyCart is an empty cart fixture
	EmptyCart = &cart.Cart{
		ID:        "cart-empty-001",
		UserID:    TestUserID,
		Items:     []cart.CartItem{},
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now(),
	}

	// CartWithItems is a cart with items fixture
	CartWithItems = func() *cart.Cart {
		return &cart.Cart{
			ID:     "cart-001",
			UserID: "user-001",
			Items: []cart.CartItem{
				{
					ID:        "item-001",
					ProductID: "prod-laptop-001",
					Name:      "Professional Laptop",
					SKU:       "LAPTOP-001",
					Quantity:  1,
					Price:     money.Money{Amount: 99999, Currency: "USD"},
					AddedAt:   time.Now().Add(-30 * time.Minute),
				},
				{
					ID:        "item-002",
					ProductID: "prod-phone-001",
					Name:      "Smartphone X",
					SKU:       "PHONE-001",
					Quantity:  2,
					Price:     money.Money{Amount: 79999, Currency: "USD"},
					AddedAt:   time.Now().Add(-20 * time.Minute),
				},
			},
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now(),
		}
	}

	// CartSingleItem is a cart with a single item
	CartSingleItem = func() *cart.Cart {
		return &cart.Cart{
			ID:     "cart-single-001",
			UserID: "user-002",
			Items: []cart.CartItem{
				{
					ID:        "item-single-001",
					ProductID: "prod-tshirt-001",
					Name:      "Classic T-Shirt",
					SKU:       "TSHIRT-001",
					Quantity:  3,
					Price:     money.Money{Amount: 2999, Currency: "USD"},
					AddedAt:   time.Now().Add(-15 * time.Minute),
				},
			},
			CreatedAt: time.Now().Add(-30 * time.Minute),
			UpdatedAt: time.Now(),
		}
	}

	// GuestCart is a cart for guest user (session-based)
	GuestCart = func() *cart.Cart {
		return &cart.Cart{
			ID:        "cart-guest-001",
			SessionID: "session-guest-001",
			Items: []cart.CartItem{
				{
					ID:        "item-guest-001",
					ProductID: "prod-laptop-001",
					Name:      "Professional Laptop",
					SKU:       "LAPTOP-001",
					Quantity:  1,
					Price:     money.Money{Amount: 99999, Currency: "USD"},
					AddedAt:   time.Now().Add(-10 * time.Minute),
				},
			},
			CreatedAt: time.Now().Add(-15 * time.Minute),
			UpdatedAt: time.Now(),
		}
	}
)

// CloneCart creates a deep copy of a cart for test isolation
func CloneCart(c *cart.Cart) *cart.Cart {
	if c == nil {
		return nil
	}

	cloned := &cart.Cart{
		ID:        c.ID,
		UserID:    c.UserID,
		SessionID: c.SessionID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.ExpiresAt != nil {
		exp := *c.ExpiresAt
		cloned.ExpiresAt = &exp
	}

	cloned.Items = make([]cart.CartItem, len(c.Items))
	copy(cloned.Items, c.Items)

	return cloned
}
