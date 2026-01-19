package mocks

import (
	"context"
	"errors"

	"github.com/devchuckcamp/gocommerce/catalog"
)

// Common errors for testing
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrDatabase      = errors.New("database error")
)

// MockProductRepository is a mock implementation of catalog.ProductRepository
type MockProductRepository struct {
	Products      map[string]*catalog.Product
	SearchResults []*catalog.Product

	// Error injection
	FindByIDError       error
	FindBySKUError      error
	FindByCategoryError error
	FindByBrandError    error
	SearchError         error
	SaveError           error
	DeleteError         error
}

// NewMockProductRepository creates a new mock product repository
func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		Products: make(map[string]*catalog.Product),
	}
}

// FindByID returns a product by ID
func (m *MockProductRepository) FindByID(ctx context.Context, id string) (*catalog.Product, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}
	if product, ok := m.Products[id]; ok {
		return product, nil
	}
	return nil, ErrNotFound
}

// FindBySKU returns a product by SKU
func (m *MockProductRepository) FindBySKU(ctx context.Context, sku string) (*catalog.Product, error) {
	if m.FindBySKUError != nil {
		return nil, m.FindBySKUError
	}
	for _, p := range m.Products {
		if p.SKU == sku {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// FindByCategory returns products in a category
func (m *MockProductRepository) FindByCategory(ctx context.Context, categoryID string, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	if m.FindByCategoryError != nil {
		return nil, m.FindByCategoryError
	}
	products := make([]*catalog.Product, 0)
	for _, p := range m.Products {
		if p.CategoryID == categoryID {
			products = append(products, p)
		}
	}
	return products, nil
}

// FindByBrand returns products by brand ID
func (m *MockProductRepository) FindByBrand(ctx context.Context, brandID string, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	if m.FindByBrandError != nil {
		return nil, m.FindByBrandError
	}
	products := make([]*catalog.Product, 0)
	for _, p := range m.Products {
		if p.BrandID == brandID {
			products = append(products, p)
		}
	}
	return products, nil
}

// Search returns products matching the search criteria
func (m *MockProductRepository) Search(ctx context.Context, query string, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	if m.SearchError != nil {
		return nil, m.SearchError
	}
	if m.SearchResults != nil {
		return m.SearchResults, nil
	}
	// Return all products if no specific results set
	products := make([]*catalog.Product, 0, len(m.Products))
	for _, p := range m.Products {
		products = append(products, p)
	}
	return products, nil
}

// Save saves a product (create or update)
func (m *MockProductRepository) Save(ctx context.Context, product *catalog.Product) error {
	if m.SaveError != nil {
		return m.SaveError
	}
	m.Products[product.ID] = product
	return nil
}

// Delete deletes a product
func (m *MockProductRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if _, exists := m.Products[id]; !exists {
		return ErrNotFound
	}
	delete(m.Products, id)
	return nil
}

// MockCategoryRepository is a mock implementation of catalog.CategoryRepository
type MockCategoryRepository struct {
	Categories map[string]*catalog.Category

	FindByIDError     error
	FindBySlugError   error
	FindChildrenError error
	FindRootsError    error
	FindAllError      error
	SaveError         error
	DeleteError       error
}

// NewMockCategoryRepository creates a new mock category repository
func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		Categories: make(map[string]*catalog.Category),
	}
}

// FindByID returns a category by ID
func (m *MockCategoryRepository) FindByID(ctx context.Context, id string) (*catalog.Category, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}
	if category, ok := m.Categories[id]; ok {
		return category, nil
	}
	return nil, ErrNotFound
}

// FindBySlug returns a category by slug
func (m *MockCategoryRepository) FindBySlug(ctx context.Context, slug string) (*catalog.Category, error) {
	if m.FindBySlugError != nil {
		return nil, m.FindBySlugError
	}
	for _, c := range m.Categories {
		if c.Slug == slug {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

// FindChildren returns child categories
func (m *MockCategoryRepository) FindChildren(ctx context.Context, parentID string) ([]*catalog.Category, error) {
	if m.FindChildrenError != nil {
		return nil, m.FindChildrenError
	}
	categories := make([]*catalog.Category, 0)
	for _, c := range m.Categories {
		if c.ParentID != nil && *c.ParentID == parentID {
			categories = append(categories, c)
		}
	}
	return categories, nil
}

// FindRoots returns root categories (categories without a parent)
func (m *MockCategoryRepository) FindRoots(ctx context.Context) ([]*catalog.Category, error) {
	if m.FindRootsError != nil {
		return nil, m.FindRootsError
	}
	categories := make([]*catalog.Category, 0)
	for _, c := range m.Categories {
		if c.ParentID == nil {
			categories = append(categories, c)
		}
	}
	return categories, nil
}

// FindAll returns all categories
func (m *MockCategoryRepository) FindAll(ctx context.Context) ([]*catalog.Category, error) {
	if m.FindAllError != nil {
		return nil, m.FindAllError
	}
	categories := make([]*catalog.Category, 0, len(m.Categories))
	for _, c := range m.Categories {
		categories = append(categories, c)
	}
	return categories, nil
}

// Save saves a category
func (m *MockCategoryRepository) Save(ctx context.Context, category *catalog.Category) error {
	if m.SaveError != nil {
		return m.SaveError
	}
	m.Categories[category.ID] = category
	return nil
}

// Delete deletes a category
func (m *MockCategoryRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Categories, id)
	return nil
}

// MockBrandRepository is a mock implementation of catalog.BrandRepository
type MockBrandRepository struct {
	Brands map[string]*catalog.Brand

	FindByIDError   error
	FindBySlugError error
	FindAllError    error
	SaveError       error
	DeleteError     error
}

// NewMockBrandRepository creates a new mock brand repository
func NewMockBrandRepository() *MockBrandRepository {
	return &MockBrandRepository{
		Brands: make(map[string]*catalog.Brand),
	}
}

// FindByID returns a brand by ID
func (m *MockBrandRepository) FindByID(ctx context.Context, id string) (*catalog.Brand, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}
	if brand, ok := m.Brands[id]; ok {
		return brand, nil
	}
	return nil, ErrNotFound
}

// FindBySlug returns a brand by slug
func (m *MockBrandRepository) FindBySlug(ctx context.Context, slug string) (*catalog.Brand, error) {
	if m.FindBySlugError != nil {
		return nil, m.FindBySlugError
	}
	for _, b := range m.Brands {
		if b.Slug == slug {
			return b, nil
		}
	}
	return nil, ErrNotFound
}

// FindAll returns all brands
func (m *MockBrandRepository) FindAll(ctx context.Context) ([]*catalog.Brand, error) {
	if m.FindAllError != nil {
		return nil, m.FindAllError
	}
	brands := make([]*catalog.Brand, 0, len(m.Brands))
	for _, b := range m.Brands {
		brands = append(brands, b)
	}
	return brands, nil
}

// Save saves a brand
func (m *MockBrandRepository) Save(ctx context.Context, brand *catalog.Brand) error {
	if m.SaveError != nil {
		return m.SaveError
	}
	m.Brands[brand.ID] = brand
	return nil
}

// Delete deletes a brand
func (m *MockBrandRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Brands, id)
	return nil
}

// MockVariantRepository is a mock implementation of catalog.VariantRepository
type MockVariantRepository struct {
	Variants map[string]*catalog.Variant

	FindByIDError        error
	FindBySKUError       error
	FindByProductIDError error
	SaveError            error
	DeleteError          error
}

// NewMockVariantRepository creates a new mock variant repository
func NewMockVariantRepository() *MockVariantRepository {
	return &MockVariantRepository{
		Variants: make(map[string]*catalog.Variant),
	}
}

// FindByID returns a variant by ID
func (m *MockVariantRepository) FindByID(ctx context.Context, id string) (*catalog.Variant, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}
	if variant, ok := m.Variants[id]; ok {
		return variant, nil
	}
	return nil, ErrNotFound
}

// FindBySKU returns a variant by SKU
func (m *MockVariantRepository) FindBySKU(ctx context.Context, sku string) (*catalog.Variant, error) {
	if m.FindBySKUError != nil {
		return nil, m.FindBySKUError
	}
	for _, v := range m.Variants {
		if v.SKU == sku {
			return v, nil
		}
	}
	return nil, ErrNotFound
}

// FindByProductID returns variants for a product
func (m *MockVariantRepository) FindByProductID(ctx context.Context, productID string) ([]*catalog.Variant, error) {
	if m.FindByProductIDError != nil {
		return nil, m.FindByProductIDError
	}
	variants := make([]*catalog.Variant, 0)
	for _, v := range m.Variants {
		if v.ProductID == productID {
			variants = append(variants, v)
		}
	}
	return variants, nil
}

// Save saves a variant
func (m *MockVariantRepository) Save(ctx context.Context, variant *catalog.Variant) error {
	if m.SaveError != nil {
		return m.SaveError
	}
	m.Variants[variant.ID] = variant
	return nil
}

// Delete deletes a variant
func (m *MockVariantRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Variants, id)
	return nil
}
