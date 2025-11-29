package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devchuckcamp/goauthx"

	"github.com/devchuckcamp/gocommerce-api/internal/config"
	"github.com/devchuckcamp/gocommerce-api/internal/database"
	httpserver "github.com/devchuckcamp/gocommerce-api/internal/http"
	"github.com/devchuckcamp/gocommerce-api/internal/repository"
	"github.com/devchuckcamp/gocommerce-api/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Starting E-Commerce API...")
	log.Printf("Database: %s", cfg.Database.Driver)

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize goauthx
	authxConfig := cfg.ToGoAuthXConfig()
	authStore, err := goauthx.NewStore(authxConfig.Database)
	if err != nil {
		log.Fatalf("Failed to create auth store: %v", err)
	}

	// Run goauthx migrations
	log.Println("Running goauthx migrations...")
	authMigrator := goauthx.NewMigrator(authStore, authxConfig.Database.Driver)
	if err := authMigrator.Up(context.Background()); err != nil {
		log.Printf("Warning: Auth migrations error: %v", err)
	} else {
		log.Println("✓ goauthx migrations completed successfully")
	}

	// Run gocommerce migrations
	log.Println("Running gocommerce migrations...")
	if err := db.RunCommerceMigrations(context.Background()); err != nil {
		log.Fatalf("Failed to run gocommerce migrations: %v", err)
	}

	// Optionally seed the database (for development)
	if os.Getenv("SEED_DB") == "true" {
		log.Println("Seeding database with sample data...")
		if err := db.SeedCommerce(context.Background()); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	}

	authService, err := goauthx.NewService(authxConfig, authStore)
	if err != nil {
		log.Fatalf("Failed to create auth service: %v", err)
	}

	log.Println("Authentication service initialized")

	// Initialize repositories
	productRepo := repository.NewProductRepository(db.DB)
	variantRepo := repository.NewVariantRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	brandRepo := repository.NewBrandRepository(db.DB)
	cartRepo := repository.NewCartRepository(db.DB)
	orderRepo := repository.NewOrderRepository(db.DB)
	promotionRepo := repository.NewPromotionRepository(db.DB)

	log.Println("Repositories initialized")

	// Initialize services
	// Tax calculator (8.75% tax rate for example)
	taxCalculator := services.NewSimpleTaxCalculator(0.0875)

	// Create catalog service
	catalogService := services.NewCatalogService(
		productRepo,
		variantRepo,
		categoryRepo,
		brandRepo,
	)

	// Create cart service (no inventory service for now)
	cartService := services.NewCartService(
		cartRepo,
		productRepo,
		variantRepo,
		nil, // inventoryService
	)

	// Create pricing service (no shipping calculator for now)
	pricingService := services.NewPricingService(
		promotionRepo,
		taxCalculator,
		nil, // shippingCalculator
	)

	// Create order service (no inventory or payment gateway for now)
	orderService := services.NewOrderService(
		orderRepo,
		pricingService.Service,
		nil, // inventoryService
		nil, // paymentGateway
	)

	log.Println("Domain services initialized")

	// Create HTTP server
	server := httpserver.NewServer(
		authService,
		catalogService,
		cartService,
		orderService,
	)

	// Setup HTTP server
	httpSrv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      server.Router(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("E-Commerce API is running")
	log.Printf("API available at http://localhost:%s/api/v1", cfg.Server.Port)
	log.Printf("Health check: http://localhost:%s/health", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
