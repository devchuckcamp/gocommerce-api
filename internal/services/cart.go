package services

import (
	"github.com/devchuckcamp/gocommerce-api/internal/utils"
	"github.com/devchuckcamp/gocommerce/cart"
	"github.com/devchuckcamp/gocommerce/catalog"
	"github.com/devchuckcamp/gocommerce/inventory"
)

// CartService holds the gocommerce cart service
type CartService struct {
	*cart.CartService
}

// NewCartService creates a new CartService using gocommerce domain service
func NewCartService(
	cartRepo cart.Repository,
	productRepo catalog.ProductRepository,
	variantRepo catalog.VariantRepository,
	inventoryService inventory.Service, // can be nil if not using inventory
) *CartService {
	svc := cart.NewCartService(
		cartRepo,
		productRepo,
		variantRepo,
		inventoryService,
		utils.GenerateID,
	)

	return &CartService{
		CartService: svc,
	}
}

// WithPriceResolver attaches an optional price resolver for dynamic pricing
func (s *CartService) WithPriceResolver(resolver cart.PriceResolver) *CartService {
	s.CartService.WithPriceResolver(resolver)
	return s
}
