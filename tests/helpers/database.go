package helpers

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB holds the test database connection
type TestDB struct {
	*gorm.DB
}

// SetupTestDB creates a test database connection
// It reads configuration from .env file in project root
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Try to load .env from project root
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://commerce:commerce@localhost:5432/commerce?sslmode=disable"
	}

	// Use a test-specific database if available
	testDSN := os.Getenv("TEST_DB_DSN")
	if testDSN != "" {
		dsn = testDSN
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return &TestDB{db}
}

// Close closes the database connection
func (db *TestDB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// TruncateTables truncates the specified tables
func (db *TestDB) TruncateTables(t *testing.T, tables ...string) {
	t.Helper()
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Logf("Warning: could not truncate table %s: %v", table, err)
		}
	}
}

// WithTransaction runs a test function within a transaction and rolls back after
func (db *TestDB) WithTransaction(t *testing.T, fn func(tx *gorm.DB)) {
	t.Helper()
	tx := db.Begin()
	defer tx.Rollback()

	fn(tx)
}

// Ping checks if the database connection is alive
func (db *TestDB) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// SkipIfNoDatabase skips the test if database is not available
func SkipIfNoDatabase(t *testing.T) *TestDB {
	t.Helper()

	// Try to load .env
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://commerce:commerce@localhost:5432/commerce?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Skipf("Skipping test: cannot get sql.DB: %v", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("Skipping test: cannot ping database: %v", err)
		return nil
	}

	return &TestDB{db}
}
