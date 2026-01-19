package fixtures

import (
	"time"

	"github.com/devchuckcamp/gocommerce/money"
	"github.com/devchuckcamp/gocommerce/orders"
)

// Order fixtures
var (
	// TestShippingAddress is a sample shipping address
	TestShippingAddress = orders.Address{
		FirstName:    "John",
		LastName:     "Doe",
		AddressLine1: "123 Main St",
		AddressLine2: "Apt 4",
		City:         "New York",
		State:        "NY",
		PostalCode:   "10001",
		Country:      "US",
		Phone:        "+1-555-555-5555",
	}

	// TestBillingAddress is a sample billing address
	TestBillingAddress = orders.Address{
		FirstName:    "John",
		LastName:     "Doe",
		Company:      "Acme Inc",
		AddressLine1: "456 Business Ave",
		City:         "New York",
		State:        "NY",
		PostalCode:   "10002",
		Country:      "US",
		Phone:        "+1-555-555-5556",
	}

	// OrderPending is a pending order fixture
	OrderPending = func() *orders.Order {
		return &orders.Order{
			ID:              "order-pending-001",
			OrderNumber:     "ORD-2024-00001",
			UserID:          "user-001",
			Status:          orders.OrderStatusPending,
			ShippingAddress: TestShippingAddress,
			BillingAddress:  TestBillingAddress,
			Items: []orders.OrderItem{
				{
					ID:        "oi-001",
					ProductID: "prod-laptop-001",
					Name:      "Professional Laptop",
					SKU:       "LAPTOP-001",
					Quantity:  1,
					UnitPrice: money.Money{Amount: 99999, Currency: "USD"},
					Total:     money.Money{Amount: 99999, Currency: "USD"},
				},
			},
			Subtotal:      money.Money{Amount: 99999, Currency: "USD"},
			TaxTotal:      money.Money{Amount: 8750, Currency: "USD"},
			ShippingTotal: money.Money{Amount: 1000, Currency: "USD"},
			Total:         money.Money{Amount: 109749, Currency: "USD"},
			CreatedAt:     time.Now().Add(-2 * time.Hour),
			UpdatedAt:     time.Now(),
		}
	}

	// OrderProcessing is a processing order fixture
	OrderProcessing = func() *orders.Order {
		return &orders.Order{
			ID:              "order-processing-001",
			OrderNumber:     "ORD-2024-00002",
			UserID:          "user-001",
			Status:          orders.OrderStatusProcessing,
			ShippingAddress: TestShippingAddress,
			BillingAddress:  TestShippingAddress,
			Items: []orders.OrderItem{
				{
					ID:        "oi-002",
					ProductID: "prod-phone-001",
					Name:      "Smartphone X",
					SKU:       "PHONE-001",
					Quantity:  1,
					UnitPrice: money.Money{Amount: 79999, Currency: "USD"},
					Total:     money.Money{Amount: 79999, Currency: "USD"},
				},
			},
			Subtotal:      money.Money{Amount: 79999, Currency: "USD"},
			TaxTotal:      money.Money{Amount: 7000, Currency: "USD"},
			ShippingTotal: money.Money{Amount: 0, Currency: "USD"},
			Total:         money.Money{Amount: 86999, Currency: "USD"},
			CreatedAt:     time.Now().Add(-24 * time.Hour),
			UpdatedAt:     time.Now().Add(-12 * time.Hour),
		}
	}

	// OrderCompleted is a completed order fixture
	OrderCompleted = func() *orders.Order {
		return &orders.Order{
			ID:              "order-completed-001",
			OrderNumber:     "ORD-2024-00003",
			UserID:          "user-002",
			Status:          orders.OrderStatusDelivered,
			ShippingAddress: TestShippingAddress,
			BillingAddress:  TestBillingAddress,
			Items: []orders.OrderItem{
				{
					ID:        "oi-003",
					ProductID: "prod-tshirt-001",
					Name:      "Classic T-Shirt",
					SKU:       "TSHIRT-001",
					Quantity:  3,
					UnitPrice: money.Money{Amount: 2999, Currency: "USD"},
					Total:     money.Money{Amount: 8997, Currency: "USD"},
				},
			},
			Subtotal:      money.Money{Amount: 8997, Currency: "USD"},
			TaxTotal:      money.Money{Amount: 787, Currency: "USD"},
			ShippingTotal: money.Money{Amount: 500, Currency: "USD"},
			Total:         money.Money{Amount: 10284, Currency: "USD"},
			CreatedAt:     time.Now().Add(-72 * time.Hour),
			UpdatedAt:     time.Now().Add(-24 * time.Hour),
		}
	}

	// OrderCancelled is a cancelled order fixture
	OrderCancelled = func() *orders.Order {
		return &orders.Order{
			ID:              "order-cancelled-001",
			OrderNumber:     "ORD-2024-00004",
			UserID:          "user-001",
			Status:          orders.OrderStatusCanceled,
			ShippingAddress: TestShippingAddress,
			BillingAddress:  TestShippingAddress,
			Items: []orders.OrderItem{
				{
					ID:        "oi-004",
					ProductID: "prod-laptop-001",
					Name:      "Professional Laptop",
					SKU:       "LAPTOP-001",
					Quantity:  2,
					UnitPrice: money.Money{Amount: 99999, Currency: "USD"},
					Total:     money.Money{Amount: 199998, Currency: "USD"},
				},
			},
			Subtotal:      money.Money{Amount: 199998, Currency: "USD"},
			TaxTotal:      money.Money{Amount: 17500, Currency: "USD"},
			ShippingTotal: money.Money{Amount: 0, Currency: "USD"},
			Total:         money.Money{Amount: 217498, Currency: "USD"},
			Notes:         "Cancelled by customer",
			CreatedAt:     time.Now().Add(-48 * time.Hour),
			UpdatedAt:     time.Now().Add(-24 * time.Hour),
		}
	}
)

// CloneOrder creates a deep copy of an order for test isolation
func CloneOrder(o *orders.Order) *orders.Order {
	if o == nil {
		return nil
	}

	cloned := &orders.Order{
		ID:              o.ID,
		OrderNumber:     o.OrderNumber,
		UserID:          o.UserID,
		Status:          o.Status,
		ShippingAddress: o.ShippingAddress,
		BillingAddress:  o.BillingAddress,
		PaymentMethodID: o.PaymentMethodID,
		Subtotal:        o.Subtotal,
		DiscountTotal:   o.DiscountTotal,
		TaxTotal:        o.TaxTotal,
		ShippingTotal:   o.ShippingTotal,
		Total:           o.Total,
		Notes:           o.Notes,
		IPAddress:       o.IPAddress,
		UserAgent:       o.UserAgent,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}

	if o.CompletedAt != nil {
		t := *o.CompletedAt
		cloned.CompletedAt = &t
	}

	if o.CanceledAt != nil {
		t := *o.CanceledAt
		cloned.CanceledAt = &t
	}

	cloned.Items = make([]orders.OrderItem, len(o.Items))
	copy(cloned.Items, o.Items)

	return cloned
}

// GetOrdersByUser returns orders for a specific user
func GetOrdersByUser(userID string) []*orders.Order {
	allOrders := []*orders.Order{
		OrderPending(),
		OrderProcessing(),
		OrderCompleted(),
		OrderCancelled(),
	}

	var result []*orders.Order
	for _, o := range allOrders {
		if o.UserID == userID {
			result = append(result, o)
		}
	}
	return result
}
