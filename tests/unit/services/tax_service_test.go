package services_test

import (
	"context"
	"testing"

	"github.com/devchuckcamp/gocommerce/money"
	"github.com/devchuckcamp/gocommerce/tax"

	"github.com/devchuckcamp/gocommerce-api/internal/services"
)

func TestSimpleTaxCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name           string
		taxRate        float64
		lineItems      []tax.TaxableItem
		shippingCost   money.Money
		expectedTotal  int64
		expectedShip   int64
	}{
		{
			name:    "calculate tax for single item",
			taxRate: 0.0875,
			lineItems: []tax.TaxableItem{
				{
					ID:        "item-1",
					Amount:    money.Money{Amount: 10000, Currency: "USD"},
					Quantity:  1,
					IsTaxable: true,
				},
			},
			shippingCost:  money.Money{Amount: 0, Currency: "USD"},
			expectedTotal: 875, // 10000 * 0.0875
			expectedShip:  0,
		},
		{
			name:    "calculate tax for multiple items",
			taxRate: 0.0875,
			lineItems: []tax.TaxableItem{
				{
					ID:        "item-1",
					Amount:    money.Money{Amount: 10000, Currency: "USD"},
					Quantity:  2,
					IsTaxable: true,
				},
				{
					ID:        "item-2",
					Amount:    money.Money{Amount: 5000, Currency: "USD"},
					Quantity:  1,
					IsTaxable: true,
				},
			},
			shippingCost:  money.Money{Amount: 0, Currency: "USD"},
			expectedTotal: 2187, // (20000 + 5000) * 0.0875 = 2187.5 -> 2187
			expectedShip:  0,
		},
		{
			name:    "calculate tax with shipping",
			taxRate: 0.10,
			lineItems: []tax.TaxableItem{
				{
					ID:        "item-1",
					Amount:    money.Money{Amount: 10000, Currency: "USD"},
					Quantity:  1,
					IsTaxable: true,
				},
			},
			shippingCost:  money.Money{Amount: 1000, Currency: "USD"},
			expectedTotal: 1100, // (10000 * 0.10) + (1000 * 0.10)
			expectedShip:  100,
		},
		{
			name:    "skip non-taxable items",
			taxRate: 0.0875,
			lineItems: []tax.TaxableItem{
				{
					ID:        "item-1",
					Amount:    money.Money{Amount: 10000, Currency: "USD"},
					Quantity:  1,
					IsTaxable: true,
				},
				{
					ID:        "item-2",
					Amount:    money.Money{Amount: 5000, Currency: "USD"},
					Quantity:  1,
					IsTaxable: false,
				},
			},
			shippingCost:  money.Money{Amount: 0, Currency: "USD"},
			expectedTotal: 875, // Only first item is taxed
			expectedShip:  0,
		},
		{
			name:    "zero tax rate",
			taxRate: 0.0,
			lineItems: []tax.TaxableItem{
				{
					ID:        "item-1",
					Amount:    money.Money{Amount: 10000, Currency: "USD"},
					Quantity:  1,
					IsTaxable: true,
				},
			},
			shippingCost:  money.Money{Amount: 500, Currency: "USD"},
			expectedTotal: 0,
			expectedShip:  0,
		},
		{
			name:          "empty line items",
			taxRate:       0.0875,
			lineItems:     []tax.TaxableItem{},
			shippingCost:  money.Money{Amount: 1000, Currency: "USD"},
			expectedTotal: 87, // Only shipping tax
			expectedShip:  87,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			calculator := services.NewSimpleTaxCalculator(tt.taxRate)

			req := tax.CalculationRequest{
				LineItems:    tt.lineItems,
				ShippingCost: tt.shippingCost,
				Address: tax.Address{
					State:   "NY",
					Country: "US",
				},
			}

			// Execute
			result, err := calculator.Calculate(context.Background(), req)

			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.TotalTax.Amount != tt.expectedTotal {
				t.Errorf("expected total tax %d, got %d", tt.expectedTotal, result.TotalTax.Amount)
			}

			if result.ShippingTax.Amount != tt.expectedShip {
				t.Errorf("expected shipping tax %d, got %d", tt.expectedShip, result.ShippingTax.Amount)
			}
		})
	}
}

func TestSimpleTaxCalculator_GetRatesForAddress(t *testing.T) {
	tests := []struct {
		name         string
		taxRate      float64
		address      tax.Address
		expectedRate float64
	}{
		{
			name:    "get rate for NY",
			taxRate: 0.0875,
			address: tax.Address{
				State:   "NY",
				Country: "US",
			},
			expectedRate: 0.0875,
		},
		{
			name:    "get rate for CA",
			taxRate: 0.0725,
			address: tax.Address{
				State:   "CA",
				Country: "US",
			},
			expectedRate: 0.0725,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			calculator := services.NewSimpleTaxCalculator(tt.taxRate)

			// Execute
			rates, err := calculator.GetRatesForAddress(context.Background(), tt.address)

			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(rates) != 1 {
				t.Errorf("expected 1 rate, got %d", len(rates))
				return
			}

			if rates[0].Rate != tt.expectedRate {
				t.Errorf("expected rate %f, got %f", tt.expectedRate, rates[0].Rate)
			}

			if rates[0].State != tt.address.State {
				t.Errorf("expected state %q, got %q", tt.address.State, rates[0].State)
			}
		})
	}
}
