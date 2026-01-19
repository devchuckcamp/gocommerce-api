package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/devchuckcamp/gocommerce-api/internal/http/response"
	"github.com/devchuckcamp/gocommerce-api/internal/services"
	"github.com/devchuckcamp/gocommerce/catalog"
)

// CatalogHandler handles catalog endpoints
type CatalogHandler struct {
	catalogService *services.CatalogService
}

// NewCatalogHandler creates a new CatalogHandler
func NewCatalogHandler(catalogService *services.CatalogService) *CatalogHandler {
	return &CatalogHandler{
		catalogService: catalogService,
	}
}

// ListProducts lists all products with pagination and search
// GET /products?page=1&page_size=20&keyword=laptop
func (h *CatalogHandler) ListProducts(c *gin.Context) {
	// Get pagination parameters
	params := response.GetPaginationParams(c)

	// Get search keyword
	keyword := c.Query("keyword")

	active := catalog.ProductStatus("active")
	filter := catalog.ProductFilter{
		Status: &active,
		Limit:  params.CalculateLimit(),
		Offset: params.CalculateOffset(),
	}

	// Get products (with search if keyword provided)
	products, err := h.catalogService.SearchProducts(c.Request.Context(), keyword, filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Get total count
	total, err := h.catalogService.CountProducts(c.Request.Context(), filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Build pagination metadata
	meta := response.NewPaginationMeta(params.Page, params.PageSize, total)
	response.SuccessWithPagination(c, products, meta)
}

// GetProduct retrieves a single product by ID
// GET /products/:id
func (h *CatalogHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		response.BadRequest(c, "Product ID is required")
		return
	}

	product, err := h.catalogService.GetProduct(c.Request.Context(), productID)
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}

	response.Success(c, product)
}

// GetProductsByCategory retrieves products by category with pagination
// GET /products/category/:id?page=1&page_size=20
func (h *CatalogHandler) GetProductsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		response.BadRequest(c, "Category ID is required")
		return
	}

	// Get pagination parameters
	params := response.GetPaginationParams(c)

	active := catalog.ProductStatus("active")
	filter := catalog.ProductFilter{
		Status:      &active,
		CategoryIDs: []string{categoryID},
		Limit:       params.CalculateLimit(),
		Offset:      params.CalculateOffset(),
	}

	// Get products
	products, err := h.catalogService.GetProductsByCategory(c.Request.Context(), categoryID, filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Get total count for this category
	total, err := h.catalogService.CountProducts(c.Request.Context(), filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Build pagination metadata
	meta := response.NewPaginationMeta(params.Page, params.PageSize, total)
	response.SuccessWithPagination(c, products, meta)
}

// ListCategories lists all categories with pagination
// GET /categories?page=1&page_size=20
func (h *CatalogHandler) ListCategories(c *gin.Context) {
	// Get pagination parameters
	params := response.GetPaginationParams(c)

	// Get all categories (gocommerce doesn't have pagination for categories yet)
	categories, err := h.catalogService.GetCategories(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Apply pagination manually
	total := int64(len(categories))
	start := params.CalculateOffset()
	end := start + params.CalculateLimit()

	if start > len(categories) {
		start = len(categories)
	}
	if end > len(categories) {
		end = len(categories)
	}

	paginatedCategories := categories[start:end]
	meta := response.NewPaginationMeta(params.Page, params.PageSize, total)
	response.SuccessWithPagination(c, paginatedCategories, meta)
}

// ListBrands lists all brands with pagination
// GET /brands?page=1&page_size=20
func (h *CatalogHandler) ListBrands(c *gin.Context) {
	// Get pagination parameters
	params := response.GetPaginationParams(c)

	// Get all brands (gocommerce doesn't have pagination for brands yet)
	brands, err := h.catalogService.GetBrands(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Apply pagination manually
	total := int64(len(brands))
	start := params.CalculateOffset()
	end := start + params.CalculateLimit()

	if start > len(brands) {
		start = len(brands)
	}
	if end > len(brands) {
		end = len(brands)
	}

	paginatedBrands := brands[start:end]
	meta := response.NewPaginationMeta(params.Page, params.PageSize, total)
	response.SuccessWithPagination(c, paginatedBrands, meta)
}
