package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/devchuckcamp/gocommerce-api/internal/database"
	"github.com/devchuckcamp/gocommerce/money"
	"github.com/devchuckcamp/gocommerce/pricing"
)

// ProductPriceRepository implements pricing.ProductPriceRepository using GORM
type ProductPriceRepository struct {
	db *gorm.DB
}

// NewProductPriceRepository creates a new ProductPriceRepository
func NewProductPriceRepository(db *gorm.DB) *ProductPriceRepository {
	return &ProductPriceRepository{db: db}
}

// FindByID finds a product price by ID
func (r *ProductPriceRepository) FindByID(ctx context.Context, id string) (*pricing.ProductPrice, error) {
	var dbPrice database.ProductPrice
	if err := r.db.WithContext(ctx).First(&dbPrice, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product price not found")
		}
		return nil, err
	}
	return r.toDomain(&dbPrice), nil
}

// FindActiveForProduct finds all active prices for a product
func (r *ProductPriceRepository) FindActiveForProduct(ctx context.Context, productID string) ([]*pricing.ProductPrice, error) {
	var dbPrices []database.ProductPrice
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_active = ?", productID, true).
		Order("priority DESC").
		Find(&dbPrices).Error; err != nil {
		return nil, err
	}
	return r.toDomainList(dbPrices), nil
}

// FindActiveForVariant finds all active prices for a variant
func (r *ProductPriceRepository) FindActiveForVariant(ctx context.Context, variantID string) ([]*pricing.ProductPrice, error) {
	var dbPrices []database.ProductPrice
	if err := r.db.WithContext(ctx).
		Where("variant_id = ? AND is_active = ?", variantID, true).
		Order("priority DESC").
		Find(&dbPrices).Error; err != nil {
		return nil, err
	}
	return r.toDomainList(dbPrices), nil
}

// FindEffectivePrice finds the effective price for a product/variant at a given time
func (r *ProductPriceRepository) FindEffectivePrice(ctx context.Context, productID string, variantID *string, at time.Time) (*pricing.ProductPrice, error) {
	var dbPrice database.ProductPrice
	query := r.db.WithContext(ctx).
		Where("product_id = ? AND is_active = ?", productID, true).
		Where("(valid_from IS NULL OR valid_from <= ?)", at).
		Where("(valid_to IS NULL OR valid_to >= ?)", at)

	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	if err := query.Order("priority DESC").First(&dbPrice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No effective price found
		}
		return nil, err
	}
	return r.toDomain(&dbPrice), nil
}

// FindEffectivePrices finds effective prices for multiple products at a given time
func (r *ProductPriceRepository) FindEffectivePrices(ctx context.Context, productIDs []string, at time.Time) (map[string]*pricing.ProductPrice, error) {
	result := make(map[string]*pricing.ProductPrice)
	if len(productIDs) == 0 {
		return result, nil
	}

	var dbPrices []database.ProductPrice
	if err := r.db.WithContext(ctx).
		Where("product_id IN ? AND is_active = ?", productIDs, true).
		Where("variant_id IS NULL"). // Base product prices only
		Where("(valid_from IS NULL OR valid_from <= ?)", at).
		Where("(valid_to IS NULL OR valid_to >= ?)", at).
		Order("priority DESC").
		Find(&dbPrices).Error; err != nil {
		return nil, err
	}

	// Group by product_id, keeping highest priority (first due to ORDER BY)
	for _, dbPrice := range dbPrices {
		if _, exists := result[dbPrice.ProductID]; !exists {
			result[dbPrice.ProductID] = r.toDomain(&dbPrice)
		}
	}
	return result, nil
}

// Save saves a product price
func (r *ProductPriceRepository) Save(ctx context.Context, price *pricing.ProductPrice) error {
	dbPrice := r.toDatabase(price)
	return r.db.WithContext(ctx).Save(dbPrice).Error
}

// Delete deletes a product price by ID
func (r *ProductPriceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.ProductPrice{}, "id = ?", id).Error
}

// DeleteByProductID deletes all prices for a product
func (r *ProductPriceRepository) DeleteByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Delete(&database.ProductPrice{}, "product_id = ?", productID).Error
}

// Helper methods for ProductPriceRepository

func (r *ProductPriceRepository) toDomain(dbPrice *database.ProductPrice) *pricing.ProductPrice {
	return &pricing.ProductPrice{
		ID:        dbPrice.ID,
		ProductID: dbPrice.ProductID,
		VariantID: dbPrice.VariantID,
		Price:     database.Int64ToMoney(dbPrice.PriceAmount, dbPrice.PriceCurrency),
		ValidFrom: dbPrice.ValidFrom,
		ValidTo:   dbPrice.ValidTo,
		Priority:  dbPrice.Priority,
		PriceType: pricing.PriceType(dbPrice.PriceType),
		IsActive:  dbPrice.IsActive,
		CreatedAt: dbPrice.CreatedAt,
		UpdatedAt: dbPrice.UpdatedAt,
	}
}

func (r *ProductPriceRepository) toDomainList(dbPrices []database.ProductPrice) []*pricing.ProductPrice {
	prices := make([]*pricing.ProductPrice, len(dbPrices))
	for i, dbPrice := range dbPrices {
		prices[i] = r.toDomain(&dbPrice)
	}
	return prices
}

func (r *ProductPriceRepository) toDatabase(price *pricing.ProductPrice) *database.ProductPrice {
	return &database.ProductPrice{
		ID:            price.ID,
		ProductID:     price.ProductID,
		VariantID:     price.VariantID,
		PriceAmount:   price.Price.Amount,
		PriceCurrency: price.Price.Currency,
		ValidFrom:     price.ValidFrom,
		ValidTo:       price.ValidTo,
		Priority:      price.Priority,
		PriceType:     string(price.PriceType),
		IsActive:      price.IsActive,
		CreatedAt:     price.CreatedAt,
		UpdatedAt:     price.UpdatedAt,
	}
}

// PromotionRepository implements pricing.PromotionRepository using GORM
type PromotionRepository struct {
	db *gorm.DB
}

// NewPromotionRepository creates a new PromotionRepository
func NewPromotionRepository(db *gorm.DB) *PromotionRepository {
	return &PromotionRepository{db: db}
}

// FindByCode finds a promotion by code
func (r *PromotionRepository) FindByCode(ctx context.Context, code string) (*pricing.Promotion, error) {
	var dbPromotion database.Promotion
	if err := r.db.WithContext(ctx).First(&dbPromotion, "code = ? AND active = ?", code, true).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("promotion not found")
		}
		return nil, err
	}

	// Check if promotion is valid (within date range)
	now := time.Now()
	if now.Before(dbPromotion.StartDate) || now.After(dbPromotion.EndDate) {
		return nil, fmt.Errorf("promotion not valid")
	}

	return r.toDomain(&dbPromotion)
}

// FindActive finds all active promotions
func (r *PromotionRepository) FindActive(ctx context.Context) ([]*pricing.Promotion, error) {
	now := time.Now()
	var dbPromotions []database.Promotion
	if err := r.db.WithContext(ctx).
		Where("active = ? AND start_date <= ? AND end_date >= ?", true, now, now).
		Find(&dbPromotions).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbPromotions)
}

// Save saves a promotion
func (r *PromotionRepository) Save(ctx context.Context, promotion *pricing.Promotion) error {
	dbPromotion := r.toDatabase(promotion)
	return r.db.WithContext(ctx).Save(dbPromotion).Error
}

// Helper methods

func (r *PromotionRepository) toDomain(dbPromotion *database.Promotion) (*pricing.Promotion, error) {
	var productIDs []string
	if err := database.UnmarshalJSON(dbPromotion.ProductIDs, &productIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product IDs: %w", err)
	}

	var categoryIDs []string
	if err := database.UnmarshalJSON(dbPromotion.CategoryIDs, &categoryIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal category IDs: %w", err)
	}

	var value float64
	if dbPromotion.DiscountPercentage > 0 {
		value = dbPromotion.DiscountPercentage
	} else if dbPromotion.DiscountAmount > 0 {
		value = float64(dbPromotion.DiscountAmount)
	}

	var minPurchase, maxDiscount *money.Money
	if dbPromotion.MinPurchaseAmount > 0 {
		m := database.Int64ToMoney(dbPromotion.MinPurchaseAmount, dbPromotion.Currency)
		minPurchase = &m
	}
	if dbPromotion.MaxDiscountAmount > 0 {
		m := database.Int64ToMoney(dbPromotion.MaxDiscountAmount, dbPromotion.Currency)
		maxDiscount = &m
	}

	return &pricing.Promotion{
		ID:                    dbPromotion.ID,
		Code:                  dbPromotion.Code,
		Name:                  dbPromotion.Name,
		Description:           dbPromotion.Description,
		DiscountType:          pricing.DiscountType(dbPromotion.Type),
		Value:                 value,
		MinPurchase:           minPurchase,
		MaxDiscount:           maxDiscount,
		ValidFrom:             dbPromotion.StartDate,
		ValidTo:               dbPromotion.EndDate,
		IsActive:              dbPromotion.Active,
		UsageLimit:            dbPromotion.UsageLimit,
		UsageCount:            dbPromotion.UsageCount,
		ApplicableProductIDs:  productIDs,
		ApplicableCategoryIDs: categoryIDs,
	}, nil
}

func (r *PromotionRepository) toDomainList(dbPromotions []database.Promotion) ([]*pricing.Promotion, error) {
	promotions := make([]*pricing.Promotion, 0, len(dbPromotions))
	for _, dbPromotion := range dbPromotions {
		promotion, err := r.toDomain(&dbPromotion)
		if err != nil {
			return nil, err
		}
		promotions = append(promotions, promotion)
	}
	return promotions, nil
}

func (r *PromotionRepository) toDatabase(promotion *pricing.Promotion) *database.Promotion {
	var discountPercentage float64
	var discountAmount int64
	var currency string = "USD"

	if promotion.DiscountType == "percentage" {
		discountPercentage = promotion.Value
	} else {
		discountAmount = int64(promotion.Value)
	}

	var minPurchase, maxDiscount int64
	if promotion.MinPurchase != nil {
		minPurchase = database.MoneyToInt64(*promotion.MinPurchase)
		currency = promotion.MinPurchase.Currency
	}
	if promotion.MaxDiscount != nil {
		maxDiscount = database.MoneyToInt64(*promotion.MaxDiscount)
		currency = promotion.MaxDiscount.Currency
	}

	return &database.Promotion{
		ID:                 promotion.ID,
		Code:               promotion.Code,
		Name:               promotion.Name,
		Description:        promotion.Description,
		Type:               string(promotion.DiscountType),
		DiscountPercentage: discountPercentage,
		DiscountAmount:     discountAmount,
		MinPurchaseAmount:  minPurchase,
		MaxDiscountAmount:  maxDiscount,
		Currency:           currency,
		StartDate:          promotion.ValidFrom,
		EndDate:            promotion.ValidTo,
		Active:             promotion.IsActive,
		UsageLimit:         promotion.UsageLimit,
		UsageCount:         promotion.UsageCount,
		ProductIDs:         database.MarshalJSON(promotion.ApplicableProductIDs),
		CategoryIDs:        database.MarshalJSON(promotion.ApplicableCategoryIDs),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
