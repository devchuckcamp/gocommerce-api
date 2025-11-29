package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/devchuckcamp/gocommerce-api/internal/http/middleware"
	"github.com/devchuckcamp/gocommerce-api/internal/http/response"
	"github.com/devchuckcamp/gocommerce-api/internal/services"
	"github.com/devchuckcamp/gocommerce/cart"
)

// CartHandler handles cart endpoints
type CartHandler struct {
	cartService *services.CartService
}

// NewCartHandler creates a new CartHandler
func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// GetCart retrieves the current user's cart
// GET /cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Try to get existing cart or create new one
	cart, err := h.cartService.GetOrCreateCart(c.Request.Context(), userID, "")
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, cart)
}

// AddItemRequest represents the request to add an item to cart
type AddItemRequest struct {
	ProductID  string            `json:"product_id" binding:"required"`
	VariantID  *string           `json:"variant_id"`
	Quantity   int               `json:"quantity" binding:"required,gt=0"`
	Attributes map[string]string `json:"attributes"`
}

// AddItem adds an item to the cart
// POST /cart/items
func (h *CartHandler) AddItem(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Get or create cart
	currentCart, err := h.cartService.GetOrCreateCart(c.Request.Context(), userID, "")
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Add item using gocommerce domain service
	addReq := cart.AddItemRequest{
		ProductID:  req.ProductID,
		VariantID:  req.VariantID,
		Quantity:   req.Quantity,
		Attributes: req.Attributes,
	}

	updatedCart, err := h.cartService.AddItem(c.Request.Context(), currentCart.ID, addReq)
	if err != nil {
		if err == cart.ErrOutOfStock {
			response.BadRequest(c, "Product is out of stock")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, updatedCart)
}

// UpdateItemQuantityRequest represents the request to update item quantity
type UpdateItemQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=0"`
}

// UpdateItemQuantity updates the quantity of an item in the cart
// PATCH /cart/items/:id
func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	itemID := c.Param("id")
	if itemID == "" {
		response.BadRequest(c, "Item ID is required")
		return
	}

	var req UpdateItemQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Get cart
	currentCart, err := h.cartService.GetOrCreateCart(c.Request.Context(), userID, "")
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Update quantity
	updatedCart, err := h.cartService.UpdateItemQuantity(c.Request.Context(), currentCart.ID, itemID, req.Quantity)
	if err != nil {
		if err == cart.ErrItemNotFound {
			response.NotFound(c, "Item not found in cart")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, updatedCart)
}

// RemoveItem removes an item from the cart
// DELETE /cart/items/:id
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	itemID := c.Param("id")
	if itemID == "" {
		response.BadRequest(c, "Item ID is required")
		return
	}

	// Get cart
	currentCart, err := h.cartService.GetOrCreateCart(c.Request.Context(), userID, "")
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Remove item
	updatedCart, err := h.cartService.RemoveItem(c.Request.Context(), currentCart.ID, itemID)
	if err != nil {
		if err == cart.ErrItemNotFound {
			response.NotFound(c, "Item not found in cart")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, updatedCart)
}

// ClearCart clears all items from the cart
// DELETE /cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get cart
	currentCart, err := h.cartService.GetOrCreateCart(c.Request.Context(), userID, "")
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Clear cart
	updatedCart, err := h.cartService.Clear(c.Request.Context(), currentCart.ID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, updatedCart)
}
