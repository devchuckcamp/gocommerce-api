package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/devchuckcamp/gocommerce-api/internal/repository"
	"github.com/devchuckcamp/gocommerce/catalog"
	"github.com/devchuckcamp/gocommerce/money"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// Setup test database connection
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		dsn = "postgres://commerce:commerce@localhost:5432/commerce?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		// Skip integration tests if database is not available
		os.Exit(0)
	}

	testDB = db

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func skipIfNoDatabase(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: database not available")
	}
}

func cleanupProductTestData(t *testing.T) {
	t.Helper()
	if testDB == nil {
		return
	}
	// Clean up any test data with 'test-' prefix
	testDB.Exec("DELETE FROM products WHERE id LIKE 'test-%'")
	testDB.Exec("DELETE FROM categories WHERE id LIKE 'test-%'")
	testDB.Exec("DELETE FROM brands WHERE id LIKE 'test-%'")
}

func TestProductRepository_Save(t *testing.T) {
	skipIfNoDatabase(t)
	defer cleanupProductTestData(t)

	repo := repository.NewProductRepository(testDB)
	ctx := context.Background()

	// First create the test brand and category
	brandRepo := repository.NewBrandRepository(testDB)
	categoryRepo := repository.NewCategoryRepository(testDB)

	testBrand := &catalog.Brand{
		ID:          "test-brand-001",
		Name:        "Test Brand",
		Slug:        "test-brand-001",
		Description: "Test brand for integration tests",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := brandRepo.Save(ctx, testBrand); err != nil {
		t.Fatalf("failed to create test brand: %v", err)
	}

	testCategory := &catalog.Category{
		ID:          "test-category-001",
		Name:        "Test Category",
		Slug:        "test-category-001",
		Description: "Test category for integration tests",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := categoryRepo.Save(ctx, testCategory); err != nil {
		t.Fatalf("failed to create test category: %v", err)
	}

	// Test saving a product
	product := &catalog.Product{
		ID:          "test-product-001",
		SKU:         "TEST-SKU-001",
		Name:        "Test Product",
		Description: "A test product for integration testing",
		BasePrice:   money.Money{Amount: 9999, Currency: "USD"},
		Status:      catalog.ProductStatus("active"),
		BrandID:     "test-brand-001",
		CategoryID:  "test-category-001",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Save(ctx, product)
	if err != nil {
		t.Fatalf("failed to save product: %v", err)
	}

	// Verify product was saved
	saved, err := repo.FindByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("failed to find saved product: %v", err)
	}

	if saved.Name != product.Name {
		t.Errorf("expected name %q, got %q", product.Name, saved.Name)
	}

	if saved.SKU != product.SKU {
		t.Errorf("expected SKU %q, got %q", product.SKU, saved.SKU)
	}

	if saved.BasePrice.Amount != product.BasePrice.Amount {
		t.Errorf("expected price %d, got %d", product.BasePrice.Amount, saved.BasePrice.Amount)
	}
}

func TestProductRepository_FindByID(t *testing.T) {
	skipIfNoDatabase(t)
	defer cleanupProductTestData(t)

	repo := repository.NewProductRepository(testDB)
	brandRepo := repository.NewBrandRepository(testDB)
	categoryRepo := repository.NewCategoryRepository(testDB)
	ctx := context.Background()

	// Setup test data
	testBrand := &catalog.Brand{
		ID:        "test-brand-002",
		Name:      "Test Brand 2",
		Slug:      "test-brand-002",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	brandRepo.Save(ctx, testBrand)

	testCategory := &catalog.Category{
		ID:        "test-category-002",
		Name:      "Test Category 2",
		Slug:      "test-category-002",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	categoryRepo.Save(ctx, testCategory)

	product := &catalog.Product{
		ID:          "test-product-002",
		SKU:         "TEST-SKU-002",
		Name:        "Test Product 2",
		Description: "Another test product",
		BasePrice:   money.Money{Amount: 4999, Currency: "USD"},
		Status:      catalog.ProductStatus("active"),
		BrandID:     "test-brand-002",
		CategoryID:  "test-category-002",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Save(ctx, product)

	tests := []struct {
		name          string
		productID     string
		expectError   bool
		expectedName  string
	}{
		{
			name:         "find existing product",
			productID:    "test-product-002",
			expectError:  false,
			expectedName: "Test Product 2",
		},
		{
			name:        "product not found",
			productID:   "non-existent-id",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByID(ctx, tt.productID)

			if tt.expectError {
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
		})
	}
}

func TestProductRepository_FindBySKU(t *testing.T) {
	skipIfNoDatabase(t)
	defer cleanupProductTestData(t)

	repo := repository.NewProductRepository(testDB)
	brandRepo := repository.NewBrandRepository(testDB)
	categoryRepo := repository.NewCategoryRepository(testDB)
	ctx := context.Background()

	// Setup test data
	testBrand := &catalog.Brand{
		ID:        "test-brand-003",
		Name:      "Test Brand 3",
		Slug:      "test-brand-003",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	brandRepo.Save(ctx, testBrand)

	testCategory := &catalog.Category{
		ID:        "test-category-003",
		Name:      "Test Category 3",
		Slug:      "test-category-003",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	categoryRepo.Save(ctx, testCategory)

	product := &catalog.Product{
		ID:          "test-product-003",
		SKU:         "UNIQUE-SKU-003",
		Name:        "Test Product with Unique SKU",
		Description: "Product with unique SKU for testing",
		BasePrice:   money.Money{Amount: 7999, Currency: "USD"},
		Status:      catalog.ProductStatus("active"),
		BrandID:     "test-brand-003",
		CategoryID:  "test-category-003",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Save(ctx, product)

	tests := []struct {
		name        string
		sku         string
		expectError bool
		expectedID  string
	}{
		{
			name:        "find by existing SKU",
			sku:         "UNIQUE-SKU-003",
			expectError: false,
			expectedID:  "test-product-003",
		},
		{
			name:        "SKU not found",
			sku:         "NON-EXISTENT-SKU",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindBySKU(ctx, tt.sku)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.ID != tt.expectedID {
				t.Errorf("expected ID %q, got %q", tt.expectedID, result.ID)
			}
		})
	}
}

func TestProductRepository_Search(t *testing.T) {
	skipIfNoDatabase(t)
	defer cleanupProductTestData(t)

	repo := repository.NewProductRepository(testDB)
	brandRepo := repository.NewBrandRepository(testDB)
	categoryRepo := repository.NewCategoryRepository(testDB)
	ctx := context.Background()

	// Setup test data
	testBrand := &catalog.Brand{
		ID:        "test-brand-search",
		Name:      "Test Brand Search",
		Slug:      "test-brand-search",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	brandRepo.Save(ctx, testBrand)

	testCategory := &catalog.Category{
		ID:        "test-category-search",
		Name:      "Test Category Search",
		Slug:      "test-category-search",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	categoryRepo.Save(ctx, testCategory)

	products := []*catalog.Product{
		{
			ID:          "test-product-search-001",
			SKU:         "SEARCH-SKU-001",
			Name:        "Gaming Laptop",
			Description: "High-performance gaming laptop",
			BasePrice:   money.Money{Amount: 149999, Currency: "USD"},
			Status:      catalog.ProductStatus("active"),
			BrandID:     "test-brand-search",
			CategoryID:  "test-category-search",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "test-product-search-002",
			SKU:         "SEARCH-SKU-002",
			Name:        "Smartphone Pro",
			Description: "Latest smartphone with laptop-level performance",
			BasePrice:   money.Money{Amount: 99999, Currency: "USD"},
			Status:      catalog.ProductStatus("active"),
			BrandID:     "test-brand-search",
			CategoryID:  "test-category-search",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "test-product-search-003",
			SKU:         "SEARCH-SKU-003",
			Name:        "Wireless Mouse",
			Description: "Ergonomic wireless mouse",
			BasePrice:   money.Money{Amount: 4999, Currency: "USD"},
			Status:      catalog.ProductStatus("active"),
			BrandID:     "test-brand-search",
			CategoryID:  "test-category-search",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, p := range products {
		repo.Save(ctx, p)
	}

	tests := []struct {
		name          string
		query         string
		expectedMin   int
	}{
		{
			name:        "search for laptop",
			query:       "laptop",
			expectedMin: 1, // At least the gaming laptop and description of smartphone
		},
		{
			name:        "search for wireless",
			query:       "wireless",
			expectedMin: 1,
		},
		{
			name:        "search with no results",
			query:       "nonexistentproductxyz",
			expectedMin: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := repo.Search(ctx, tt.query, catalog.ProductFilter{})

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Filter to only count test products
			testResults := 0
			for _, r := range results {
				if len(r.ID) > 5 && r.ID[:5] == "test-" {
					testResults++
				}
			}

			if testResults < tt.expectedMin {
				t.Errorf("expected at least %d results, got %d", tt.expectedMin, testResults)
			}
		})
	}
}

func TestProductRepository_Delete(t *testing.T) {
	skipIfNoDatabase(t)
	defer cleanupProductTestData(t)

	repo := repository.NewProductRepository(testDB)
	brandRepo := repository.NewBrandRepository(testDB)
	categoryRepo := repository.NewCategoryRepository(testDB)
	ctx := context.Background()

	// Setup test data
	testBrand := &catalog.Brand{
		ID:        "test-brand-delete",
		Name:      "Test Brand Delete",
		Slug:      "test-brand-delete",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	brandRepo.Save(ctx, testBrand)

	testCategory := &catalog.Category{
		ID:        "test-category-delete",
		Name:      "Test Category Delete",
		Slug:      "test-category-delete",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	categoryRepo.Save(ctx, testCategory)

	product := &catalog.Product{
		ID:          "test-product-delete",
		SKU:         "DELETE-SKU",
		Name:        "Product to Delete",
		Description: "This product will be deleted",
		BasePrice:   money.Money{Amount: 1999, Currency: "USD"},
		Status:      catalog.ProductStatus("active"),
		BrandID:     "test-brand-delete",
		CategoryID:  "test-category-delete",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Save(ctx, product)

	// Verify product exists
	_, err := repo.FindByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("failed to find product before delete: %v", err)
	}

	// Delete the product
	err = repo.Delete(ctx, product.ID)
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}

	// Verify product no longer exists
	_, err = repo.FindByID(ctx, product.ID)
	if err == nil {
		t.Error("expected error after deleting product, got nil")
	}
}

// TestCategoryRepository_CRUD tests Category repository operations
func TestCategoryRepository_CRUD(t *testing.T) {
	skipIfNoDatabase(t)
	defer cleanupProductTestData(t)

	repo := repository.NewCategoryRepository(testDB)
	ctx := context.Background()

	// Test Save
	category := &catalog.Category{
		ID:          "test-category-crud",
		Name:        "CRUD Test Category",
		Slug:        "test-category-crud",
		Description: "Category for CRUD testing",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Save(ctx, category)
	if err != nil {
		t.Fatalf("failed to save category: %v", err)
	}

	// Test FindByID
	found, err := repo.FindByID(ctx, category.ID)
	if err != nil {
		t.Fatalf("failed to find category by ID: %v", err)
	}
	if found.Name != category.Name {
		t.Errorf("expected name %q, got %q", category.Name, found.Name)
	}

	// Test FindBySlug
	foundBySlug, err := repo.FindBySlug(ctx, category.Slug)
	if err != nil {
		t.Fatalf("failed to find category by slug: %v", err)
	}
	if foundBySlug.ID != category.ID {
		t.Errorf("expected ID %q, got %q", category.ID, foundBySlug.ID)
	}

	// Test Delete
	err = repo.Delete(ctx, category.ID)
	if err != nil {
		t.Fatalf("failed to delete category: %v", err)
	}

	_, err = repo.FindByID(ctx, category.ID)
	if err == nil {
		t.Error("expected error after deleting category, got nil")
	}
}

// TestBrandRepository_CRUD tests Brand repository operations
func TestBrandRepository_CRUD(t *testing.T) {
	skipIfNoDatabase(t)
	defer cleanupProductTestData(t)

	repo := repository.NewBrandRepository(testDB)
	ctx := context.Background()

	// Test Save
	brand := &catalog.Brand{
		ID:          "test-brand-crud",
		Name:        "CRUD Test Brand",
		Slug:        "test-brand-crud",
		Description: "Brand for CRUD testing",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Save(ctx, brand)
	if err != nil {
		t.Fatalf("failed to save brand: %v", err)
	}

	// Test FindByID
	found, err := repo.FindByID(ctx, brand.ID)
	if err != nil {
		t.Fatalf("failed to find brand by ID: %v", err)
	}
	if found.Name != brand.Name {
		t.Errorf("expected name %q, got %q", brand.Name, found.Name)
	}

	// Test FindBySlug
	foundBySlug, err := repo.FindBySlug(ctx, brand.Slug)
	if err != nil {
		t.Fatalf("failed to find brand by slug: %v", err)
	}
	if foundBySlug.ID != brand.ID {
		t.Errorf("expected ID %q, got %q", brand.ID, foundBySlug.ID)
	}

	// Test Delete
	err = repo.Delete(ctx, brand.ID)
	if err != nil {
		t.Fatalf("failed to delete brand: %v", err)
	}

	_, err = repo.FindByID(ctx, brand.ID)
	if err == nil {
		t.Error("expected error after deleting brand, got nil")
	}
}
