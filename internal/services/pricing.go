package services

import (
	"github.com/devchuckcamp/gocommerce/pricing"
	"github.com/devchuckcamp/gocommerce/shipping"
	"github.com/devchuckcamp/gocommerce/tax"
)

// PricingService holds the gocommerce pricing service
type PricingService struct {
	pricing.Service
}

// NewPricingService creates a new PricingService using gocommerce domain service
func NewPricingService(
	promotionRepo pricing.PromotionRepository,
	taxCalculator tax.Calculator,
	shippingCalc shipping.RateCalculator, // can be nil if not using shipping
) *PricingService {
	svc := pricing.NewPricingService(
		promotionRepo,
		taxCalculator,
		shippingCalc,
	)

	return &PricingService{
		Service: svc,
	}
}
