package database

import (
	"encoding/json"
	"time"

	"github.com/devchuckcamp/gocommerce/money"
)

// Product represents a product in the database
type Product struct {
	ID          string    `gorm:"primaryKey;size:36"`
	SKU         string    `gorm:"uniqueIndex;size:100;not null"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	BasePrice   int64     `gorm:"not null"` // stored as cents
	Currency    string    `gorm:"size:3;not null;default:'USD'"`
	Status      string    `gorm:"size:20;not null;default:'active'"`
	BrandID     string    `gorm:"size:36;index"`
	CategoryID  string    `gorm:"size:36;index"`
	Images      string    `gorm:"type:text"` // JSON array of image URLs
	Metadata    string    `gorm:"type:jsonb"` // JSON metadata (attributes)
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// Variant represents a product variant in the database
type Variant struct {
	ID         string    `gorm:"primaryKey;size:36"`
	ProductID  string    `gorm:"size:36;index;not null"`
	SKU        string    `gorm:"uniqueIndex;size:100;not null"`
	Name       string    `gorm:"size:255;not null"`
	Price      int64     `gorm:"not null"` // stored as cents
	Currency   string    `gorm:"size:3;not null;default:'USD'"`
	Attributes string    `gorm:"type:jsonb"` // JSON attributes like {"color": "red", "size": "L"}
	ImageURL   string    `gorm:"size:500"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

// Category represents a product category in the database
type Category struct {
	ID          string    `gorm:"primaryKey;size:36"`
	Name        string    `gorm:"size:255;not null"`
	Slug        string    `gorm:"uniqueIndex;size:255;not null"`
	Description string    `gorm:"type:text"`
	ParentID    *string   `gorm:"size:36;index"`
	ImageURL    string    `gorm:"size:500"`
	Active      bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// Brand represents a product brand in the database
type Brand struct {
	ID          string    `gorm:"primaryKey;size:36"`
	Name        string    `gorm:"uniqueIndex;size:255;not null"`
	Slug        string    `gorm:"uniqueIndex;size:255;not null"`
	Description string    `gorm:"type:text"`
	LogoURL     string    `gorm:"size:500"`
	Active      bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// Cart represents a shopping cart in the database
type Cart struct {
	ID        string    `gorm:"primaryKey;size:36"`
	UserID    string    `gorm:"size:36;index"`
	SessionID string    `gorm:"size:100;index"`
	Items     string    `gorm:"type:jsonb;not null"` // JSON serialized CartItem array
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// Order represents an order in the database
type Order struct {
	ID              string    `gorm:"primaryKey;size:36"`
	OrderNumber     string    `gorm:"uniqueIndex;size:50;not null"`
	UserID          string    `gorm:"size:36;index;not null"`
	Status          string    `gorm:"size:20;not null;default:'pending'"`
	Items           string    `gorm:"type:jsonb;not null"` // JSON serialized OrderItem array
	ShippingAddress string    `gorm:"type:jsonb;not null"` // JSON serialized Address
	BillingAddress  string    `gorm:"type:jsonb;not null"` // JSON serialized Address
	PaymentMethodID string    `gorm:"size:100"`
	Subtotal        int64     `gorm:"not null"` // stored as cents
	DiscountTotal   int64     `gorm:"not null;default:0"`
	TaxTotal        int64     `gorm:"not null;default:0"`
	ShippingTotal   int64     `gorm:"not null;default:0"`
	Total           int64     `gorm:"not null"`
	Currency        string    `gorm:"size:3;not null;default:'USD'"`
	Notes           string    `gorm:"type:text"`
	IPAddress       string    `gorm:"size:50"`
	UserAgent       string    `gorm:"size:500"`
	CancelledAt     *time.Time
	CancelReason    string    `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

// Promotion represents a discount promotion in the database
type Promotion struct {
	ID                 string     `gorm:"primaryKey;size:36"`
	Code               string     `gorm:"uniqueIndex;size:50;not null"`
	Name               string     `gorm:"size:255;not null"`
	Description        string     `gorm:"type:text"`
	Type               string     `gorm:"size:20;not null"` // percentage, fixed, buy_x_get_y
	DiscountPercentage float64    `gorm:"type:decimal(5,2)"`
	DiscountAmount     int64      // stored as cents
	MinPurchaseAmount  int64      // stored as cents
	MaxDiscountAmount  int64      // stored as cents
	Currency           string     `gorm:"size:3;not null;default:'USD'"`
	StartDate          time.Time  `gorm:"not null"`
	EndDate            time.Time  `gorm:"not null"`
	Active             bool       `gorm:"not null;default:true"`
	UsageLimit         int        `gorm:"default:0"` // 0 = unlimited
	UsageCount         int        `gorm:"default:0"`
	ProductIDs         string     `gorm:"type:jsonb"` // JSON array of product IDs
	CategoryIDs        string     `gorm:"type:jsonb"` // JSON array of category IDs
	CreatedAt          time.Time  `gorm:"not null"`
	UpdatedAt          time.Time  `gorm:"not null"`
}

// Helper functions to convert between domain and database models

// MoneyToInt64 converts money.Money to int64 cents
func MoneyToInt64(m money.Money) int64 {
	return m.Amount
}

// Int64ToMoney converts int64 cents to money.Money
func Int64ToMoney(amount int64, currency string) money.Money {
	return money.Money{
		Amount:   amount,
		Currency: currency,
	}
}

// MarshalJSON marshals any value to JSON string
func MarshalJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// UnmarshalJSON unmarshals JSON string to target
func UnmarshalJSON(data string, target interface{}) error {
	if data == "" {
		return nil
	}
	return json.Unmarshal([]byte(data), target)
}
