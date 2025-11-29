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
