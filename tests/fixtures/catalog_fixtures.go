package fixtures

import (
	"time"

	"github.com/devchuckcamp/gocommerce/catalog"
	"github.com/devchuckcamp/gocommerce/money"
)

// Product status constants
const (
	StatusActive   catalog.ProductStatus = "active"
	StatusInactive catalog.ProductStatus = "inactive"
)

// Product fixtures
var (
	ProductLaptop = &catalog.Product{
		ID:          "prod-laptop-001",
		SKU:         "LAPTOP-001",
		Name:        "Professional Laptop",
		Description: "High-performance laptop for professionals",
		BasePrice:   money.Money{Amount: 99999, Currency: "USD"},
		Status:      StatusActive,
		BrandID:     "brand-techcorp",
		CategoryID:  "cat-electronics",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	ProductPhone = &catalog.Product{
		ID:          "prod-phone-001",
		SKU:         "PHONE-001",
		Name:        "Smartphone X",
		Description: "Latest smartphone with advanced features",
		BasePrice:   money.Money{Amount: 79999, Currency: "USD"},
		Status:      StatusActive,
		BrandID:     "brand-techcorp",
		CategoryID:  "cat-electronics",
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	ProductTShirt = &catalog.Product{
		ID:          "prod-tshirt-001",
		SKU:         "TSHIRT-001",
		Name:        "Classic T-Shirt",
		Description: "Comfortable cotton t-shirt",
		BasePrice:   money.Money{Amount: 2999, Currency: "USD"},
		Status:      StatusActive,
		BrandID:     "brand-fashionhub",
		CategoryID:  "cat-clothing",
		CreatedAt:   time.Now().Add(-72 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	ProductInactive = &catalog.Product{
		ID:          "prod-inactive-001",
		SKU:         "INACTIVE-001",
		Name:        "Discontinued Product",
		Description: "This product is no longer available",
		BasePrice:   money.Money{Amount: 4999, Currency: "USD"},
		Status:      StatusInactive,
		BrandID:     "brand-techcorp",
		CategoryID:  "cat-electronics",
		CreatedAt:   time.Now().Add(-96 * time.Hour),
		UpdatedAt:   time.Now(),
	}
)

// Category fixtures
var (
	CategoryElectronics = &catalog.Category{
		ID:          "cat-electronics",
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "Electronic devices and gadgets",
		CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	CategoryClothing = &catalog.Category{
		ID:          "cat-clothing",
		Name:        "Clothing",
		Slug:        "clothing",
		Description: "Apparel and fashion",
		CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	CategoryBooks = &catalog.Category{
		ID:          "cat-books",
		Name:        "Books",
		Slug:        "books",
		Description: "Books and publications",
		CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:   time.Now(),
	}
)

// Brand fixtures
var (
	BrandTechCorp = &catalog.Brand{
		ID:          "brand-techcorp",
		Name:        "TechCorp",
		Slug:        "techcorp",
		Description: "Leading technology manufacturer",
		CreatedAt:   time.Now().Add(-60 * 24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	BrandFashionHub = &catalog.Brand{
		ID:          "brand-fashionhub",
		Name:        "FashionHub",
		Slug:        "fashionhub",
		Description: "Premium fashion brand",
		CreatedAt:   time.Now().Add(-60 * 24 * time.Hour),
		UpdatedAt:   time.Now(),
	}
)

// Variant fixtures
var (
	VariantTShirtSmallRed = &catalog.Variant{
		ID:        "var-tshirt-s-red",
		ProductID: "prod-tshirt-001",
		SKU:       "TSHIRT-001-S-RED",
		Name:      "Classic T-Shirt - Small Red",
		Price:     money.Money{Amount: 2999, Currency: "USD"},
		CreatedAt: time.Now().Add(-72 * time.Hour),
		UpdatedAt: time.Now(),
	}

	VariantTShirtMediumBlue = &catalog.Variant{
		ID:        "var-tshirt-m-blue",
		ProductID: "prod-tshirt-001",
		SKU:       "TSHIRT-001-M-BLUE",
		Name:      "Classic T-Shirt - Medium Blue",
		Price:     money.Money{Amount: 2999, Currency: "USD"},
		CreatedAt: time.Now().Add(-72 * time.Hour),
		UpdatedAt: time.Now(),
	}
)

// GetAllProducts returns all product fixtures
func GetAllProducts() []*catalog.Product {
	return []*catalog.Product{
		ProductLaptop,
		ProductPhone,
		ProductTShirt,
		ProductInactive,
	}
}

// GetActiveProducts returns only active product fixtures
func GetActiveProducts() []*catalog.Product {
	return []*catalog.Product{
		ProductLaptop,
		ProductPhone,
		ProductTShirt,
	}
}

// GetAllCategories returns all category fixtures
func GetAllCategories() []*catalog.Category {
	return []*catalog.Category{
		CategoryElectronics,
		CategoryClothing,
		CategoryBooks,
	}
}

// GetAllBrands returns all brand fixtures
func GetAllBrands() []*catalog.Brand {
	return []*catalog.Brand{
		BrandTechCorp,
		BrandFashionHub,
	}
}

// CloneProduct creates a deep copy of a product for test isolation
func CloneProduct(p *catalog.Product) *catalog.Product {
	return &catalog.Product{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		BasePrice:   p.BasePrice,
		Status:      p.Status,
		BrandID:     p.BrandID,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
