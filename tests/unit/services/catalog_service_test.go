package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/devchuckcamp/gocommerce/catalog"

	"github.com/devchuckcamp/gocommerce-api/internal/services"
	"github.com/devchuckcamp/gocommerce-api/tests/fixtures"
	"github.com/devchuckcamp/gocommerce-api/tests/mocks"
)

func TestCatalogService_GetProduct(t *testing.T) {
	tests := []struct {
		name          string
		productID     string
		setupMock     func(*mocks.MockProductRepository, *mocks.MockSalePriceResolver)
		expectedError bool
		expectedName  string
		hasSalePrice  bool
	}{
		{
			name:      "successfully get product without sale price",
			productID: "prod-laptop-001",
			setupMock: func(repo *mocks.MockProductRepository, resolver *mocks.MockSalePriceResolver) {
				repo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop
			},
			expectedError: false,
			expectedName:  "Professional Laptop",
			hasSalePrice:  false,
		},
		{
			name:      "successfully get product with sale price",
			productID: "prod-laptop-001",
			setupMock: func(repo *mocks.MockProductRepository, resolver *mocks.MockSalePriceResolver) {
				repo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop
				resolver.AddPrice("prod-laptop-001", 89999, "USD")
			},
			expectedError: false,
			expectedName:  "Professional Laptop",
			hasSalePrice:  true,
		},
		{
			name:      "product not found",
			productID: "non-existent",
			setupMock: func(repo *mocks.MockProductRepository, resolver *mocks.MockSalePriceResolver) {
				// No products added
			},
			expectedError: true,
		},
		{
			name:      "repository error",
			productID: "prod-laptop-001",
			setupMock: func(repo *mocks.MockProductRepository, resolver *mocks.MockSalePriceResolver) {
				repo.FindByIDError = errors.New("database connection failed")
			},
			expectedError: true,
		},
		{
			name:      "sale price resolver error - still returns product",
			productID: "prod-laptop-001",
			setupMock: func(repo *mocks.MockProductRepository, resolver *mocks.MockSalePriceResolver) {
				repo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop
				resolver.FindEffectivePriceError = errors.New("price service unavailable")
			},
			expectedError: false,
			expectedName:  "Professional Laptop",
			hasSalePrice:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			productRepo := mocks.NewMockProductRepository()
			variantRepo := mocks.NewMockVariantRepository()
			categoryRepo := mocks.NewMockCategoryRepository()
			brandRepo := mocks.NewMockBrandRepository()
			priceResolver := mocks.NewMockSalePriceResolver()

			tt.setupMock(productRepo, priceResolver)

			svc := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo).
				WithSalePriceResolver(priceResolver)

			// Execute
			result, err := svc.GetProduct(context.Background(), tt.productID)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, result.Name)
			}

			if tt.hasSalePrice && result.SalePrice == nil {
				t.Error("expected sale price, got nil")
			}

			if !tt.hasSalePrice && result.SalePrice != nil {
				t.Error("expected no sale price, got one")
			}
		})
	}
}

func TestCatalogService_ListProducts(t *testing.T) {
	tests := []struct {
		name          string
		filter        catalog.ProductFilter
		setupMock     func(*mocks.MockProductRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "list all products",
			filter: catalog.ProductFilter{},
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop
				repo.Products[fixtures.ProductPhone.ID] = fixtures.ProductPhone
				repo.Products[fixtures.ProductTShirt.ID] = fixtures.ProductTShirt
			},
			expectedCount: 3,
			expectedError: false,
		},
		{
			name:   "empty product list",
			filter: catalog.ProductFilter{},
			setupMock: func(repo *mocks.MockProductRepository) {
				// No products
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:   "repository error",
			filter: catalog.ProductFilter{},
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.SearchError = errors.New("database error")
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			productRepo := mocks.NewMockProductRepository()
			variantRepo := mocks.NewMockVariantRepository()
			categoryRepo := mocks.NewMockCategoryRepository()
			brandRepo := mocks.NewMockBrandRepository()

			tt.setupMock(productRepo)

			svc := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)

			// Execute
			result, err := svc.ListProducts(context.Background(), tt.filter)

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d products, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestCatalogService_SearchProducts(t *testing.T) {
	tests := []struct {
		name          string
		keyword       string
		setupMock     func(*mocks.MockProductRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name:    "search by keyword",
			keyword: "laptop",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.SearchResults = []*catalog.Product{fixtures.ProductLaptop}
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:    "no results found",
			keyword: "nonexistent",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.SearchResults = []*catalog.Product{}
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:    "search error",
			keyword: "test",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.SearchError = errors.New("search failed")
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			productRepo := mocks.NewMockProductRepository()
			variantRepo := mocks.NewMockVariantRepository()
			categoryRepo := mocks.NewMockCategoryRepository()
			brandRepo := mocks.NewMockBrandRepository()

			tt.setupMock(productRepo)

			svc := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)

			// Execute
			result, err := svc.SearchProducts(context.Background(), tt.keyword, catalog.ProductFilter{})

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d products, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestCatalogService_GetCategories(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mocks.MockCategoryRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name: "list all categories",
			setupMock: func(repo *mocks.MockCategoryRepository) {
				repo.Categories[fixtures.CategoryElectronics.ID] = fixtures.CategoryElectronics
				repo.Categories[fixtures.CategoryClothing.ID] = fixtures.CategoryClothing
				repo.Categories[fixtures.CategoryBooks.ID] = fixtures.CategoryBooks
			},
			expectedCount: 3,
			expectedError: false,
		},
		{
			name: "empty categories",
			setupMock: func(repo *mocks.MockCategoryRepository) {
				// No categories
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "repository error",
			setupMock: func(repo *mocks.MockCategoryRepository) {
				repo.FindAllError = errors.New("database error")
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			productRepo := mocks.NewMockProductRepository()
			variantRepo := mocks.NewMockVariantRepository()
			categoryRepo := mocks.NewMockCategoryRepository()
			brandRepo := mocks.NewMockBrandRepository()

			tt.setupMock(categoryRepo)

			svc := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)

			// Execute
			result, err := svc.GetCategories(context.Background())

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d categories, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestCatalogService_GetBrands(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mocks.MockBrandRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name: "list all brands",
			setupMock: func(repo *mocks.MockBrandRepository) {
				repo.Brands[fixtures.BrandTechCorp.ID] = fixtures.BrandTechCorp
				repo.Brands[fixtures.BrandFashionHub.ID] = fixtures.BrandFashionHub
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "empty brands",
			setupMock: func(repo *mocks.MockBrandRepository) {
				// No brands
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "repository error",
			setupMock: func(repo *mocks.MockBrandRepository) {
				repo.FindAllError = errors.New("database error")
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			productRepo := mocks.NewMockProductRepository()
			variantRepo := mocks.NewMockVariantRepository()
			categoryRepo := mocks.NewMockCategoryRepository()
			brandRepo := mocks.NewMockBrandRepository()

			tt.setupMock(brandRepo)

			svc := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)

			// Execute
			result, err := svc.GetBrands(context.Background())

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d brands, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestCatalogService_GetProductsByCategory(t *testing.T) {
	tests := []struct {
		name          string
		categoryID    string
		setupMock     func(*mocks.MockProductRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name:       "get products in category",
			categoryID: "cat-electronics",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop
				repo.Products[fixtures.ProductPhone.ID] = fixtures.ProductPhone
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:       "empty category",
			categoryID: "cat-empty",
			setupMock: func(repo *mocks.MockProductRepository) {
				// No products in this category
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:       "repository error",
			categoryID: "cat-electronics",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.FindByCategoryError = errors.New("database error")
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			productRepo := mocks.NewMockProductRepository()
			variantRepo := mocks.NewMockVariantRepository()
			categoryRepo := mocks.NewMockCategoryRepository()
			brandRepo := mocks.NewMockBrandRepository()

			tt.setupMock(productRepo)

			svc := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)

			// Execute
			result, err := svc.GetProductsByCategory(context.Background(), tt.categoryID, catalog.ProductFilter{})

			// Assert
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d products, got %d", tt.expectedCount, len(result))
			}
		})
	}
}
