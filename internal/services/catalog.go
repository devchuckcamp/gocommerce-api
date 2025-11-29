package services

import (
	"context"

	"github.com/devchuckcamp/gocommerce/catalog"
)

// CatalogService provides additional catalog operations
type CatalogService struct {
	productRepo  catalog.ProductRepository
	variantRepo  catalog.VariantRepository
	categoryRepo catalog.CategoryRepository
	brandRepo    catalog.BrandRepository
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

// GetProduct retrieves a product by ID
func (s *CatalogService) GetProduct(ctx context.Context, id string) (*catalog.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}

// ListProducts lists products with optional filters
func (s *CatalogService) ListProducts(ctx context.Context, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	// For a simple list, we can search with empty query
	return s.productRepo.Search(ctx, "", filter)
}

// SearchProducts searches products by keyword
func (s *CatalogService) SearchProducts(ctx context.Context, keyword string, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	return s.productRepo.Search(ctx, keyword, filter)
}

// GetProductsByCategory retrieves products in a category
func (s *CatalogService) GetProductsByCategory(ctx context.Context, categoryID string, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	return s.productRepo.FindByCategory(ctx, categoryID, filter)
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
