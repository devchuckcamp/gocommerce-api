package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/devchuckcamp/gocommerce/migrations"
)

// RunCommerceMigrations runs gocommerce migrations using the migrations package
func (db *DB) RunCommerceMigrations(ctx context.Context) error {
	// Get underlying sql.DB for migrations
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Create migration executor
	executor := newGormExecutor(sqlDB)

	// Create migration repository (PostgreSQL)
	repo := migrations.NewPostgreSQLRepository(executor, migrations.TableName)

	// Create migration manager
	manager := migrations.NewManager(repo, executor)

	// Register example migrations (creates tables for catalog, cart, orders, pricing)
	if err := manager.RegisterMultiple(migrations.PostgreSQLExampleMigrations); err != nil {
		return fmt.Errorf("failed to register migrations: %w", err)
	}

	// Run migrations
	if err := manager.Up(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✓ gocommerce migrations completed successfully")
	return nil
}

// SeedCommerce seeds the database with sample e-commerce data
func (db *DB) SeedCommerce(ctx context.Context) error {
	// Get underlying sql.DB for seeding
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Create seeder
	executor := newGormExecutor(sqlDB)
	seeder := migrations.NewSeeder(executor)

	// Register all seeds
	seeder.RegisterMultiple(migrations.AllSeeds)

	// Run all seeds
	if err := seeder.Run(ctx); err != nil {
		return fmt.Errorf("failed to run seeds: %w", err)
	}

	log.Println("✓ Database seeded successfully")
	return nil
}

// gormExecutor implements migrations.Executor interface for GORM
type gormExecutor struct {
	db *sql.DB
	tx *sql.Tx
}

func newGormExecutor(db *sql.DB) *gormExecutor {
	return &gormExecutor{db: db}
}

func (e *gormExecutor) Exec(ctx context.Context, query string, args ...interface{}) error {
	var err error
	if e.tx != nil {
		_, err = e.tx.ExecContext(ctx, query, args...)
	} else {
		_, err = e.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (e *gormExecutor) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error

	if e.tx != nil {
		rows, err = e.tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = e.db.QueryContext(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Prepare result
	var result []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{} to represent each column
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the result into the column pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map to hold the row data
		row := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}

		result = append(result, row)
	}

	return result, rows.Err()
}

func (e *gormExecutor) Begin(ctx context.Context) (migrations.Executor, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &gormExecutor{db: e.db, tx: tx}, nil
}

func (e *gormExecutor) Commit(ctx context.Context) error {
	if e.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}
	return e.tx.Commit()
}

func (e *gormExecutor) Rollback(ctx context.Context) error {
	if e.tx == nil {
		return fmt.Errorf("no transaction to rollback")
	}
	return e.tx.Rollback()
}

// Seed adds sample data to the database (useful for development)
func (db *DB) Seed() error {
	log.Println("Seeding database with sample data...")

	// Check if we already have data
	var count int64
	db.Model(&Product{}).Count(&count)
	if count > 0 {
		log.Println("Database already has data, skipping seed")
		return nil
	}

	// Create sample categories
	categories := []Category{
		{
			ID:          "cat-1",
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "Electronic devices and gadgets",
			Active:      true,
		},
		{
			ID:          "cat-2",
			Name:        "Clothing",
			Slug:        "clothing",
			Description: "Apparel and fashion",
			Active:      true,
		},
		{
			ID:          "cat-3",
			Name:        "Books",
			Slug:        "books",
			Description: "Books and publications",
			Active:      true,
		},
	}

	for _, cat := range categories {
		if err := db.Create(&cat).Error; err != nil {
			return fmt.Errorf("failed to create category: %w", err)
		}
	}

	// Create sample brands
	brands := []Brand{
		{
			ID:          "brand-1",
			Name:        "TechCorp",
			Slug:        "techcorp",
			Description: "Leading technology manufacturer",
			Active:      true,
		},
		{
			ID:          "brand-2",
			Name:        "FashionHub",
			Slug:        "fashionhub",
			Description: "Premium fashion brand",
			Active:      true,
		},
	}

	for _, brand := range brands {
		if err := db.Create(&brand).Error; err != nil {
			return fmt.Errorf("failed to create brand: %w", err)
		}
	}

	// Create sample products
	products := []Product{
		{
			ID:          "prod-1",
			SKU:         "LAPTOP-001",
			Name:        "Professional Laptop",
			Description: "High-performance laptop for professionals",
			BasePrice:   99999, // $999.99
			Currency:    "USD",
			Status:      "active",
			BrandID:     "brand-1",
			CategoryID:  "cat-1",
		},
		{
			ID:          "prod-2",
			SKU:         "PHONE-001",
			Name:        "Smartphone X",
			Description: "Latest smartphone with advanced features",
			BasePrice:   79999, // $799.99
			Currency:    "USD",
			Status:      "active",
			BrandID:     "brand-1",
			CategoryID:  "cat-1",
		},
		{
			ID:          "prod-3",
			SKU:         "TSHIRT-001",
			Name:        "Classic T-Shirt",
			Description: "Comfortable cotton t-shirt",
			BasePrice:   2999, // $29.99
			Currency:    "USD",
			Status:      "active",
			BrandID:     "brand-2",
			CategoryID:  "cat-2",
		},
	}

	for _, prod := range products {
		if err := db.Create(&prod).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}
	}

	// Create sample variants for t-shirt
	variants := []Variant{
		{
			ID:         "var-1",
			ProductID:  "prod-3",
			SKU:        "TSHIRT-001-S-RED",
			Name:       "Classic T-Shirt - Small Red",
			Price:      2999,
			Currency:   "USD",
			Attributes: `{"size": "S", "color": "Red"}`,
		},
		{
			ID:         "var-2",
			ProductID:  "prod-3",
			SKU:        "TSHIRT-001-M-BLUE",
			Name:       "Classic T-Shirt - Medium Blue",
			Price:      2999,
			Currency:   "USD",
			Attributes: `{"size": "M", "color": "Blue"}`,
		},
	}

	for _, variant := range variants {
		if err := db.Create(&variant).Error; err != nil {
			return fmt.Errorf("failed to create variant: %w", err)
		}
	}

	log.Println("Database seeded successfully")
	return nil
}
