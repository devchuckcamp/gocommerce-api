package handlers

import (
	"github.com/devchuckcamp/goauthx"
	"github.com/gin-gonic/gin"

	"github.com/devchuckcamp/gocommerce-api/internal/http/middleware"
	"github.com/devchuckcamp/gocommerce-api/internal/http/response"
	"github.com/devchuckcamp/gocommerce-api/internal/services"
	"github.com/devchuckcamp/gocommerce/orders"
)

// OrderHandler handles order endpoints
type OrderHandler struct {
	orderService *services.OrderService
	cartService  *services.CartService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService *services.OrderService, cartService *services.CartService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		cartService:  cartService,
	}
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	ShippingAddress  AddressRequest  `json:"shipping_address" binding:"required"`
	BillingAddress   *AddressRequest `json:"billing_address"`
	PaymentMethodID  string          `json:"payment_method_id"`
	PromotionCodes   []string        `json:"promotion_codes"`
	ShippingMethodID string          `json:"shipping_method_id"`
	Notes            string          `json:"notes"`
}

// AddressRequest represents an address
type AddressRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Company     string `json:"company"`
	Address1    string `json:"address1" binding:"required"`
	Address2    string `json:"address2"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Country     string `json:"country" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

// CreateOrder creates a new order from the user's cart
// POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Get user's cart
	cart, err := h.cartService.GetOrCreateCart(c.Request.Context(), userID, "")
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Check if cart has items
	if len(cart.Items) == 0 {
		response.BadRequest(c, "Cart is empty")
		return
	}

	// Convert addresses
	shippingAddr := orders.Address{
		FirstName:    req.ShippingAddress.FirstName,
		LastName:     req.ShippingAddress.LastName,
		Company:      req.ShippingAddress.Company,
		AddressLine1: req.ShippingAddress.Address1,
		AddressLine2: req.ShippingAddress.Address2,
		City:         req.ShippingAddress.City,
		State:        req.ShippingAddress.State,
		PostalCode:   req.ShippingAddress.PostalCode,
		Country:      req.ShippingAddress.Country,
		Phone:        req.ShippingAddress.PhoneNumber,
	}

	billingAddr := shippingAddr
	if req.BillingAddress != nil {
		billingAddr = orders.Address{
			FirstName:    req.BillingAddress.FirstName,
			LastName:     req.BillingAddress.LastName,
			Company:      req.BillingAddress.Company,
			AddressLine1: req.BillingAddress.Address1,
			AddressLine2: req.BillingAddress.Address2,
			City:         req.BillingAddress.City,
			State:        req.BillingAddress.State,
			PostalCode:   req.BillingAddress.PostalCode,
			Country:      req.BillingAddress.Country,
			Phone:        req.BillingAddress.PhoneNumber,
		}
	}

	// Create order using gocommerce domain service
	createReq := orders.CreateOrderRequest{
		Cart:             cart,
		UserID:           userID,
		ShippingAddress:  shippingAddr,
		BillingAddress:   billingAddr,
		PaymentMethodID:  req.PaymentMethodID,
		PromotionCodes:   req.PromotionCodes,
		ShippingMethodID: req.ShippingMethodID,
		Notes:            req.Notes,
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
	}

	order, err := h.orderService.CreateFromCart(c.Request.Context(), createReq)
	if err != nil {
		if err == orders.ErrEmptyCart {
			response.BadRequest(c, "Cart is empty")
			return
		}
		if err == orders.ErrInvalidAddress {
			response.BadRequest(c, "Invalid address")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, order)
}

// ListOrders lists the current user's orders with pagination
// GET /orders?page=1&page_size=20
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get pagination parameters
	params := response.GetPaginationParams(c)

	filter := orders.OrderFilter{
		Limit:  params.CalculateLimit(),
		Offset: params.CalculateOffset(),
	}

	ordersList, err := h.orderService.GetUserOrders(c.Request.Context(), userID, filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Apply pagination with total count
	// Note: gocommerce doesn't provide count, so we estimate from results
	total := int64(len(ordersList))
	if len(ordersList) == params.CalculateLimit() {
		// If we got a full page, there might be more
		total = int64(params.Page * params.PageSize) // Estimate
	}

	meta := response.NewPaginationMeta(params.Page, params.PageSize, total)
	response.SuccessWithPagination(c, ordersList, meta)
}

// GetOrder retrieves a specific order by ID
// GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		response.BadRequest(c, "Order ID is required")
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		if err == orders.ErrOrderNotFound {
			response.NotFound(c, "Order not found")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	// Check if user owns the order or has support/admin role
	if order.UserID != userID {
		// Allow admin, manager, or customer experience roles to view any order
		if !hasAnyRole(c, string(goauthx.RoleAdmin), string(goauthx.RoleManager), string(goauthx.RoleCustomerExperience)) {
			response.Forbidden(c, "You don't have permission to view this order")
			return
		}
	}

	response.Success(c, order)
}

// hasAnyRole checks if the user has any of the specified roles
func hasAnyRole(c *gin.Context, roles ...string) bool {
	userRoles, ok := middleware.GetUserRoles(c)
	if !ok {
		return false
	}

	for _, userRole := range userRoles {
		for _, requiredRole := range roles {
			if userRole == requiredRole {
				return true
			}
		}
	}
	return false
}
