package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/devchuckcamp/gocommerce-api/internal/database"
	"github.com/devchuckcamp/gocommerce/cart"
)

// CartRepository implements cart.Repository using GORM
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new CartRepository
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// FindByID finds a cart by ID
func (r *CartRepository) FindByID(ctx context.Context, id string) (*cart.Cart, error) {
	var dbCart database.Cart
	if err := r.db.WithContext(ctx).First(&dbCart, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cart.ErrCartNotFound
		}
		return nil, err
	}

	return r.toDomain(&dbCart)
}

// FindByUserID finds a cart by user ID
func (r *CartRepository) FindByUserID(ctx context.Context, userID string) (*cart.Cart, error) {
	var dbCart database.Cart
	if err := r.db.WithContext(ctx).First(&dbCart, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cart.ErrCartNotFound
		}
		return nil, err
	}

	return r.toDomain(&dbCart)
}

// FindBySessionID finds a cart by session ID
func (r *CartRepository) FindBySessionID(ctx context.Context, sessionID string) (*cart.Cart, error) {
	var dbCart database.Cart
	if err := r.db.WithContext(ctx).First(&dbCart, "session_id = ?", sessionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cart.ErrCartNotFound
		}
		return nil, err
	}

	return r.toDomain(&dbCart)
}

// Save saves a cart
func (r *CartRepository) Save(ctx context.Context, c *cart.Cart) error {
	dbCart := r.toDatabase(c)
	return r.db.WithContext(ctx).Save(dbCart).Error
}

// Delete deletes a cart
func (r *CartRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Cart{}, "id = ?", id).Error
}

// Helper methods

func (r *CartRepository) toDomain(dbCart *database.Cart) (*cart.Cart, error) {
	var items []cart.CartItem
	if err := database.UnmarshalJSON(dbCart.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart items: %w", err)
	}

	return &cart.Cart{
		ID:        dbCart.ID,
		UserID:    dbCart.UserID,
		SessionID: dbCart.SessionID,
		Items:     items,
		CreatedAt: dbCart.CreatedAt,
		UpdatedAt: dbCart.UpdatedAt,
	}, nil
}

func (r *CartRepository) toDatabase(c *cart.Cart) *database.Cart {
	return &database.Cart{
		ID:        c.ID,
		UserID:    c.UserID,
		SessionID: c.SessionID,
		Items:     database.MarshalJSON(c.Items),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
