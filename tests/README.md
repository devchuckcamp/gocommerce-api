# GoShop Test Suite

This directory contains the comprehensive test suite for the GoShop e-commerce API.

## Directory Structure

```
tests/
├── unit/                           # Unit tests (no external dependencies)
│   ├── services/                   # Service layer tests
│   │   ├── catalog_service_test.go # CatalogService tests
│   │   └── tax_service_test.go     # SimpleTaxCalculator tests
│   └── handlers/                   # HTTP handler tests
│       └── catalog_handler_test.go # CatalogHandler tests
├── integration/                    # Integration tests (requires database)
│   └── repository/                 # Repository tests against real DB
│       └── product_repository_test.go
├── mocks/                          # Mock implementations
│   ├── catalog_repository.go       # MockProductRepository, MockCategoryRepository, etc.
│   ├── cart_repository.go          # MockCartRepository
│   ├── order_repository.go         # MockOrderRepository
│   └── pricing_mock.go             # MockSalePriceResolver, MockPromotionRepository
├── fixtures/                       # Test data fixtures
│   ├── catalog_fixtures.go         # Product, Category, Brand fixtures
│   ├── cart_fixtures.go            # Cart fixtures
│   └── order_fixtures.go           # Order fixtures
├── helpers/                        # Test utilities and helpers
│   ├── database.go                 # Database test utilities
│   └── http.go                     # HTTP test utilities
└── README.md                       # This file
```

## Running Tests

### Run All Tests

```bash
go test ./tests/... -v
```

### Run Unit Tests Only

```bash
go test ./tests/unit/... -v
```

### Run Integration Tests Only

```bash
go test ./tests/integration/... -v
```

### Run Tests with Race Detection

```bash
go test ./tests/... -race -v
```

### Run Tests with Coverage

```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Run Specific Test

```bash
go test ./tests/unit/services -run TestCatalogService_GetProduct -v
```

## Test Categories

### Unit Tests

Unit tests are isolated tests that don't require external dependencies like databases. They use mock implementations to test business logic.

**Services Tests** (`tests/unit/services/`)
- `TestCatalogService_GetProduct` - Tests product retrieval with sale price resolution
- `TestCatalogService_ListProducts` - Tests product listing
- `TestCatalogService_SearchProducts` - Tests product search functionality
- `TestCatalogService_GetCategories` - Tests category listing
- `TestCatalogService_GetBrands` - Tests brand listing
- `TestCatalogService_GetProductsByCategory` - Tests category filtering
- `TestSimpleTaxCalculator_Calculate` - Tests tax calculation
- `TestSimpleTaxCalculator_GetRatesForAddress` - Tests tax rate lookup

**Handler Tests** (`tests/unit/handlers/`)
- `TestCatalogHandler_ListProducts` - Tests product listing endpoint
- `TestCatalogHandler_GetProduct` - Tests single product endpoint
- `TestCatalogHandler_ListCategories` - Tests category listing endpoint
- `TestCatalogHandler_ListBrands` - Tests brand listing endpoint

### Integration Tests

Integration tests require a running database and test the actual repository implementations.

**Repository Tests** (`tests/integration/repository/`)
- `TestProductRepository_Save` - Tests product persistence
- `TestProductRepository_FindByID` - Tests product lookup by ID
- `TestProductRepository_FindBySKU` - Tests product lookup by SKU
- `TestProductRepository_Search` - Tests product search with real database
- `TestProductRepository_Delete` - Tests product deletion
- `TestCategoryRepository_CRUD` - Tests full CRUD operations for categories
- `TestBrandRepository_CRUD` - Tests full CRUD operations for brands

## Environment Configuration

### Database Configuration

Integration tests use the database connection from environment variables:

1. `TEST_DATABASE_URL` - Preferred for test-specific database
2. `DB_DSN` - Falls back to main application DSN
3. Default: `postgres://commerce:commerce@localhost:5432/commerce?sslmode=disable`

### Running with Test Database

For isolation, use a separate test database:

```bash
export TEST_DATABASE_URL="postgres://commerce:commerce@localhost:5432/commerce_test?sslmode=disable"
go test ./tests/integration/... -v
```

## Test Data Management

### Fixtures

Test fixtures provide consistent test data:

```go
import "github.com/devchuckcamp/gocommerce-api/tests/fixtures"

// Use product fixtures
product := fixtures.ProductLaptop
cart := fixtures.CartWithItems()
order := fixtures.OrderPending()
```

### Test Data Cleanup

Integration tests clean up test data automatically using the `test-` prefix convention:

```go
func cleanupProductTestData(t *testing.T) {
    testDB.Exec("DELETE FROM products WHERE id LIKE 'test-%'")
}
```

## Mocks

Mock implementations allow testing without external dependencies:

```go
import "github.com/devchuckcamp/gocommerce-api/tests/mocks"

// Create mock repository
productRepo := mocks.NewMockProductRepository()

// Setup test data
productRepo.Products[fixtures.ProductLaptop.ID] = fixtures.ProductLaptop

// Inject error for testing error paths
productRepo.FindByIDError = errors.New("database error")
```

## Writing New Tests

### Adding Unit Tests

1. Create test file in appropriate `tests/unit/` subdirectory
2. Use `_test` package suffix (e.g., `services_test`)
3. Use table-driven tests for multiple scenarios
4. Use mocks for all external dependencies

### Adding Integration Tests

1. Create test file in `tests/integration/` subdirectory
2. Use `skipIfNoDatabase(t)` at the start of each test
3. Clean up test data with `defer cleanupTestData(t)`
4. Prefix all test data IDs with `test-`

### Test Naming Convention

- Test functions: `TestStructName_MethodName`
- Subtests: Descriptive lowercase with underscores
- Example: `TestCatalogService_GetProduct/product_not_found`

## CI/CD Integration

For CI/CD pipelines, run tests with:

```bash
# Unit tests (no database required)
go test ./tests/unit/... -race -v

# Integration tests (requires database)
go test ./tests/integration/... -v

# All tests with coverage
go test ./tests/... -race -coverprofile=coverage.out -v
```
