package services

import (
	"github.com/devchuckcamp/gocommerce-api/internal/utils"
	"github.com/devchuckcamp/gocommerce/inventory"
	"github.com/devchuckcamp/gocommerce/orders"
	"github.com/devchuckcamp/gocommerce/payments"
	"github.com/devchuckcamp/gocommerce/pricing"
)

// OrderService holds the gocommerce order service
type OrderService struct {
	orders.Service
}

// NewOrderService creates a new OrderService using gocommerce domain service
func NewOrderService(
	orderRepo orders.Repository,
	pricingService pricing.Service,
	inventoryService inventory.Service, // can be nil if not using inventory
	paymentGateway payments.Gateway, // can be nil for now
) *OrderService {
	svc := orders.NewOrderService(
		orderRepo,
		pricingService,
		inventoryService,
		paymentGateway,
		utils.GenerateOrderNumber,
		utils.GenerateID,
	)

	return &OrderService{
		Service: svc,
	}
}
