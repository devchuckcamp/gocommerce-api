package services

import (
	"context"

	"github.com/devchuckcamp/gocommerce/money"
	"github.com/devchuckcamp/gocommerce/tax"
)

// SimpleTaxCalculator implements tax.Calculator with a fixed tax rate
type SimpleTaxCalculator struct {
	rate float64 // e.g., 0.0875 for 8.75%
}

// NewSimpleTaxCalculator creates a new SimpleTaxCalculator
func NewSimpleTaxCalculator(rate float64) *SimpleTaxCalculator {
	return &SimpleTaxCalculator{rate: rate}
}

// Calculate calculates tax for the given request
func (c *SimpleTaxCalculator) Calculate(ctx context.Context, req tax.CalculationRequest) (*tax.CalculationResult, error) {
	// Calculate total from line items
	currency := "USD"
	if len(req.LineItems) > 0 {
		currency = req.LineItems[0].Amount.Currency
	}
	
	var subtotal int64
	lineItemTaxes := make([]tax.LineItemTax, len(req.LineItems))
	
	for i, item := range req.LineItems {
		if !item.IsTaxable {
			continue
		}
		
		itemTotal := item.Amount.Amount * int64(item.Quantity)
		itemTax := int64(float64(itemTotal) * c.rate)
		subtotal += itemTotal
		
		lineItemTaxes[i] = tax.LineItemTax{
			LineItemID: item.ID,
			TaxAmount:  money.Money{Amount: itemTax, Currency: currency},
			TaxRates: []tax.AppliedTaxRate{
				{
					Name:         "Sales Tax",
					Rate:         c.rate,
					Jurisdiction: req.Address.State,
				},
			},
		}
	}

	// Calculate shipping tax
	shippingTax := int64(float64(req.ShippingCost.Amount) * c.rate)
	totalTax := int64(float64(subtotal) * c.rate) + shippingTax

	result := &tax.CalculationResult{
		TotalTax: money.Money{Amount: totalTax, Currency: currency},
		TaxRates: []tax.AppliedTaxRate{
			{
				Name:         "Sales Tax",
				Rate:         c.rate,
				Jurisdiction: req.Address.State,
			},
		},
		LineItemTaxes: lineItemTaxes,
		ShippingTax:   money.Money{Amount: shippingTax, Currency: currency},
	}

	return result, nil
}

// GetRatesForAddress returns the tax rates for a given address
func (c *SimpleTaxCalculator) GetRatesForAddress(ctx context.Context, address tax.Address) ([]tax.TaxRate, error) {
	return []tax.TaxRate{
		{
			ID:       "state-tax",
			Name:     "Sales Tax",
			Rate:     c.rate,
			State:    address.State,
			TaxType:  "sales",
			Priority: 1,
		},
	}, nil
}
