package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/devchuckcamp/gocommerce-api/internal/http/handlers"
	"github.com/devchuckcamp/gocommerce-api/internal/services"
	"github.com/devchuckcamp/gocommerce-api/tests/fixtures"
	"github.com/devchuckcamp/gocommerce-api/tests/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupCatalogTestRouter(handler *handlers.CatalogHandler) *gin.Engine {
	router := gin.New()
	router.GET("/catalog/products", handler.ListProducts)
	router.GET("/catalog/products/:id", handler.GetProduct)
	router.GET("/catalog/products/category/:id", handler.GetProductsByCategory)
	router.GET("/catalog/categories", handler.ListCategories)
	router.GET("/catalog/brands", handler.ListBrands)
	return router
}

func TestCatalogHandler_ListProducts(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		setupMock      func(*mocks.MockProductRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "list products successfully",
			queryParams: "",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop
				repo.Products[fixtures.ProductPhone.ID] = fixtures.ProductPhone
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				data, ok := response["data"].([]interface{})
				if !ok {
					t.Fatal("expected data to be an array")
				}
				if len(data) != 2 {
					t.Errorf("expected 2 products, got %d", len(data))
				}
			},
		},
		{
			name:        "list products with pagination",
			queryParams: "?page=1&page_size=10",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				// Check meta exists
				if _, ok := response["meta"]; !ok {
					t.Error("expected meta in response")
				}
			},
		},
		{
			name:        "list products with search keyword",
			queryParams: "?keyword=laptop",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.SearchResults = fixtures.GetActiveProducts()[:1] // Just laptop
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				data, ok := response["data"].([]interface{})
				if !ok {
					t.Fatal("expected data to be an array")
				}
				if len(data) != 1 {
					t.Errorf("expected 1 product, got %d", len(data))
				}
			},
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

			catalogService := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)
			handler := handlers.NewCatalogHandler(catalogService)
			router := setupCatalogTestRouter(handler)

			// Execute
			req, _ := http.NewRequest(http.MethodGet, "/catalog/products"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// Assert
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestCatalogHandler_GetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		setupMock      func(*mocks.MockProductRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "get product successfully",
			productID: "prod-laptop-001",
			setupMock: func(repo *mocks.MockProductRepository) {
				repo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Fatal("expected data to be an object")
				}
				if data["Name"] != "Professional Laptop" {
					t.Errorf("expected name 'Professional Laptop', got %v", data["Name"])
				}
			},
		},
		{
			name:      "product not found",
			productID: "non-existent",
			setupMock: func(repo *mocks.MockProductRepository) {
				// No products
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if _, ok := response["error"]; !ok {
					t.Error("expected error in response")
				}
			},
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

			catalogService := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)
			handler := handlers.NewCatalogHandler(catalogService)
			router := setupCatalogTestRouter(handler)

			// Execute
			req, _ := http.NewRequest(http.MethodGet, "/catalog/products/"+tt.productID, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// Assert
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestCatalogHandler_ListCategories(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockCategoryRepository)
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "list categories successfully",
			setupMock: func(repo *mocks.MockCategoryRepository) {
				repo.Categories[fixtures.CategoryElectronics.ID] = fixtures.CategoryElectronics
				repo.Categories[fixtures.CategoryClothing.ID] = fixtures.CategoryClothing
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "empty categories list",
			setupMock: func(repo *mocks.MockCategoryRepository) {
				// No categories
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
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

			catalogService := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)
			handler := handlers.NewCatalogHandler(catalogService)
			router := setupCatalogTestRouter(handler)

			// Execute
			req, _ := http.NewRequest(http.MethodGet, "/catalog/categories", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// Assert
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}
			data, ok := response["data"].([]interface{})
			if !ok {
				t.Fatal("expected data to be an array")
			}
			if len(data) != tt.expectedCount {
				t.Errorf("expected %d categories, got %d", tt.expectedCount, len(data))
			}
		})
	}
}

func TestCatalogHandler_ListBrands(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockBrandRepository)
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "list brands successfully",
			setupMock: func(repo *mocks.MockBrandRepository) {
				repo.Brands[fixtures.BrandTechCorp.ID] = fixtures.BrandTechCorp
				repo.Brands[fixtures.BrandFashionHub.ID] = fixtures.BrandFashionHub
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "empty brands list",
			setupMock: func(repo *mocks.MockBrandRepository) {
				// No brands
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
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

			catalogService := services.NewCatalogService(productRepo, variantRepo, categoryRepo, brandRepo)
			handler := handlers.NewCatalogHandler(catalogService)
			router := setupCatalogTestRouter(handler)

			// Execute
			req, _ := http.NewRequest(http.MethodGet, "/catalog/brands", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// Assert
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}
			data, ok := response["data"].([]interface{})
			if !ok {
				t.Fatal("expected data to be an array")
			}
			if len(data) != tt.expectedCount {
				t.Errorf("expected %d brands, got %d", tt.expectedCount, len(data))
			}
		})
	}
}
