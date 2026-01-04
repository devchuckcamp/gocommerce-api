package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/devchuckcamp/gocommerce-api/internal/database"
	"github.com/devchuckcamp/gocommerce/cart"
	"github.com/devchuckcamp/gocommerce/money"
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

	return r.toDomain(ctx, &dbCart)
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

	return r.toDomain(ctx, &dbCart)
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

	return r.toDomain(ctx, &dbCart)
}

// Save saves a cart
func (r *CartRepository) Save(ctx context.Context, c *cart.Cart) error {
	if c == nil {
		return errors.New("cart cannot be nil")
	}

	// Sync cart header + items in a single transaction.
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dbCart := r.toDatabase(c)
		if err := tx.Save(dbCart).Error; err != nil {
			return err
		}

		// Load existing item IDs.
		var existingIDs []string
		if err := tx.Model(&database.CartItem{}).
			Where("cart_id = ?", c.ID).
			Pluck("id", &existingIDs).Error; err != nil {
			return err
		}
		existing := make(map[string]struct{}, len(existingIDs))
		for _, id := range existingIDs {
			existing[id] = struct{}{}
		}

		desired := make(map[string]struct{}, len(c.Items))
		for _, item := range c.Items {
			desired[item.ID] = struct{}{}
			dbItem := r.toDatabaseItem(c.ID, item)
			if err := tx.Save(dbItem).Error; err != nil {
				return err
			}
		}

		// Delete items that were removed.
		var toDelete []string
		for id := range existing {
			if _, ok := desired[id]; !ok {
				toDelete = append(toDelete, id)
			}
		}
		if len(toDelete) > 0 {
			if err := tx.Where("cart_id = ? AND id IN ?", c.ID, toDelete).Delete(&database.CartItem{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete deletes a cart
func (r *CartRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Cart{}, "id = ?", id).Error
}

// Helper methods

func (r *CartRepository) toDomain(ctx context.Context, dbCart *database.Cart) (*cart.Cart, error) {
	var dbItems []database.CartItem
	if err := r.db.WithContext(ctx).
		Where("cart_id = ?", dbCart.ID).
		Order("added_at ASC").
		Find(&dbItems).Error; err != nil {
		return nil, err
	}

	items := make([]cart.CartItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		item, err := r.toDomainItem(&dbItem)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return &cart.Cart{
		ID:        dbCart.ID,
		UserID:    dbCart.UserID,
		SessionID: dbCart.SessionID,
		Items:     items,
		CreatedAt: dbCart.CreatedAt,
		UpdatedAt: dbCart.UpdatedAt,
		ExpiresAt: dbCart.ExpiresAt,
	}, nil
}

func (r *CartRepository) toDatabase(c *cart.Cart) *database.Cart {
	createdAt := c.CreatedAt
	updatedAt := c.UpdatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	return &database.Cart{
		ID:        c.ID,
		UserID:    c.UserID,
		SessionID: c.SessionID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiresAt: c.ExpiresAt,
	}
}

func (r *CartRepository) toDomainItem(dbItem *database.CartItem) (*cart.CartItem, error) {
	var attrs map[string]string
	if err := database.UnmarshalJSON(dbItem.Attributes, &attrs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart item attributes: %w", err)
	}

	price := money.Money{Amount: dbItem.PriceAmount, Currency: dbItem.PriceCurrency}
	return &cart.CartItem{
		ID:         dbItem.ID,
		ProductID:  dbItem.ProductID,
		VariantID:  dbItem.VariantID,
		SKU:        dbItem.SKU,
		Name:       dbItem.Name,
		Price:      price,
		Quantity:   dbItem.Quantity,
		Attributes: attrs,
		AddedAt:    dbItem.AddedAt,
	}, nil
}

func (r *CartRepository) toDatabaseItem(cartID string, item cart.CartItem) *database.CartItem {
	attrs := "{}"
	if item.Attributes != nil {
		attrs = database.MarshalJSON(item.Attributes)
	}

	addedAt := item.AddedAt
	if addedAt.IsZero() {
		addedAt = time.Now()
	}

	return &database.CartItem{
		ID:            item.ID,
		CartID:        cartID,
		ProductID:     item.ProductID,
		VariantID:     item.VariantID,
		SKU:           item.SKU,
		Name:          item.Name,
		PriceAmount:   item.Price.Amount,
		PriceCurrency: item.Price.Currency,
		Quantity:      item.Quantity,
		Attributes:    attrs,
		AddedAt:       addedAt,
	}
}
