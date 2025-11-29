package http

import (
	"github.com/gin-gonic/gin"

	"github.com/devchuckcamp/goauthx"
	"github.com/devchuckcamp/gocommerce-api/internal/http/handlers"
	"github.com/devchuckcamp/gocommerce-api/internal/http/middleware"
	"github.com/devchuckcamp/gocommerce-api/internal/services"
)

// Server holds the HTTP server configuration
type Server struct {
	router *gin.Engine
}

// NewServer creates a new HTTP server
func NewServer(
	authService *goauthx.Service,
	catalogService *services.CatalogService,
	cartService *services.CartService,
	orderService *services.OrderService,
) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Apply global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	catalogHandler := handlers.NewCatalogHandler(catalogService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService, cartService)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Register routes
	setupRoutes(router, authHandler, catalogHandler, cartHandler, orderHandler, authMiddleware)

	return &Server{
		router: router,
	}
}

// setupRoutes sets up all API routes
func setupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	catalogHandler *handlers.CatalogHandler,
	cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 group
	v1 := router.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)

		// Google OAuth routes
		auth.GET("/google", authHandler.GoogleOAuthURL)
		auth.GET("/google/callback", authHandler.GoogleOAuthCallback)

		// Protected auth routes
		authProtected := auth.Group("")
		authProtected.Use(authMiddleware.Authenticate())
		{
			authProtected.GET("/profile", authHandler.Profile)
			authProtected.POST("/logout", authHandler.Logout)
		}
	}

	// Catalog routes (public)
	catalog := v1.Group("/catalog")
	{
		catalog.GET("/products", catalogHandler.ListProducts)
		catalog.GET("/products/:id", catalogHandler.GetProduct)
		catalog.GET("/products/category/:id", catalogHandler.GetProductsByCategory)
		catalog.GET("/categories", catalogHandler.ListCategories)
		catalog.GET("/brands", catalogHandler.ListBrands)
	}

	// Cart routes (protected)
	cart := v1.Group("/cart")
	cart.Use(authMiddleware.Authenticate())
	{
		cart.GET("", cartHandler.GetCart)
		cart.POST("/items", cartHandler.AddItem)
		cart.PATCH("/items/:id", cartHandler.UpdateItemQuantity)
		cart.DELETE("/items/:id", cartHandler.RemoveItem)
		cart.DELETE("", cartHandler.ClearCart)
	}

	// Order routes (protected)
	orders := v1.Group("/orders")
	orders.Use(authMiddleware.Authenticate())
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.ListOrders)
		orders.GET("/:id", orderHandler.GetOrder)
	}
}

// Router returns the Gin router instance
func (s *Server) Router() *gin.Engine {
	return s.router
}
