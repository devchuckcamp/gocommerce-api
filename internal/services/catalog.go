package services

import (
	"context"
	"time"

	"github.com/devchuckcamp/gocommerce/catalog"
	"github.com/devchuckcamp/gocommerce/money"
	"github.com/devchuckcamp/gocommerce/pricing"
)

// SalePriceResolver is an interface for resolving effective sale prices for products
type SalePriceResolver interface {
	FindEffectivePrice(ctx context.Context, productID string, variantID *string, at time.Time) (*pricing.ProductPrice, error)
	FindEffectivePrices(ctx context.Context, productIDs []string, at time.Time) (map[string]*pricing.ProductPrice, error)
}

// ProductResponse wraps catalog.Product with sale price information
type ProductResponse struct {
	*catalog.Product
	SalePrice *money.Money `json:"SalePrice,omitempty"`
}

// CatalogService provides additional catalog operations
type CatalogService struct {
	productRepo       catalog.ProductRepository
	variantRepo       catalog.VariantRepository
	categoryRepo      catalog.CategoryRepository
	brandRepo         catalog.BrandRepository
	salePriceResolver SalePriceResolver
}

// NewCatalogService creates a new CatalogService
func NewCatalogService(
	productRepo catalog.ProductRepository,
	variantRepo catalog.VariantRepository,
	categoryRepo catalog.CategoryRepository,
	brandRepo catalog.BrandRepository,
) *CatalogService {
	return &CatalogService{
		productRepo:  productRepo,
		variantRepo:  variantRepo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
	}
}

// WithSalePriceResolver attaches the sale price resolver for sale price resolution
func (s *CatalogService) WithSalePriceResolver(resolver SalePriceResolver) *CatalogService {
	s.salePriceResolver = resolver
	return s
}

// GetProduct retrieves a product by ID with sale price
func (s *CatalogService) GetProduct(ctx context.Context, id string) (*ProductResponse, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := &ProductResponse{Product: product}

	// Fetch sale price if resolver is available
	if s.salePriceResolver != nil {
		if salePrice, err := s.salePriceResolver.FindEffectivePrice(ctx, id, nil, time.Now()); err == nil && salePrice != nil {
			response.SalePrice = &salePrice.Price
		}
	}

	return response, nil
}

// ListProducts lists products with optional filters including sale prices
func (s *CatalogService) ListProducts(ctx context.Context, filter catalog.ProductFilter) ([]*ProductResponse, error) {
	products, err := s.productRepo.Search(ctx, "", filter)
	if err != nil {
		return nil, err
	}

	return s.enrichWithSalePrices(ctx, products)
}

// SearchProducts searches products by keyword with sale prices
func (s *CatalogService) SearchProducts(ctx context.Context, keyword string, filter catalog.ProductFilter) ([]*ProductResponse, error) {
	products, err := s.productRepo.Search(ctx, keyword, filter)
	if err != nil {
		return nil, err
	}

	return s.enrichWithSalePrices(ctx, products)
}

// GetProductsByCategory retrieves products in a category with sale prices
func (s *CatalogService) GetProductsByCategory(ctx context.Context, categoryID string, filter catalog.ProductFilter) ([]*ProductResponse, error) {
	products, err := s.productRepo.FindByCategory(ctx, categoryID, filter)
	if err != nil {
		return nil, err
	}

	return s.enrichWithSalePrices(ctx, products)
}

// GetCategories retrieves all categories
func (s *CatalogService) GetCategories(ctx context.Context) ([]*catalog.Category, error) {
	return s.categoryRepo.FindAll(ctx)
}

// GetBrands retrieves all brands
func (s *CatalogService) GetBrands(ctx context.Context) ([]*catalog.Brand, error) {
	return s.brandRepo.FindAll(ctx)
}

// CountProducts counts total products matching the filter
func (s *CatalogService) CountProducts(ctx context.Context, filter catalog.ProductFilter) (int64, error) {
	if repo, ok := s.productRepo.(interface {
		CountProducts(ctx context.Context, filter catalog.ProductFilter) (int64, error)
	}); ok {
		return repo.CountProducts(ctx, filter)
	}
	return 0, nil
}

// enrichWithSalePrices batch-fetches sale prices for products and returns ProductResponses
func (s *CatalogService) enrichWithSalePrices(ctx context.Context, products []*catalog.Product) ([]*ProductResponse, error) {
	responses := make([]*ProductResponse, len(products))

	// If no resolver, return products without sale prices
	if s.salePriceResolver == nil {
		for i, product := range products {
			responses[i] = &ProductResponse{Product: product}
		}
		return responses, nil
	}

	// Collect product IDs for batch query
	productIDs := make([]string, len(products))
	for i, product := range products {
		productIDs[i] = product.ID
	}

	// Batch fetch sale prices
	salePrices, err := s.salePriceResolver.FindEffectivePrices(ctx, productIDs, time.Now())
	if err != nil {
		// On error, return products without sale prices
		for i, product := range products {
			responses[i] = &ProductResponse{Product: product}
		}
		return responses, nil
	}

	// Map products to responses with sale prices
	for i, product := range products {
		response := &ProductResponse{Product: product}
		if salePrice, exists := salePrices[product.ID]; exists {
			response.SalePrice = &salePrice.Price
		}
		responses[i] = response
	}

	return responses, nil
}
