package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/devchuckcamp/gocommerce-api/internal/database"
	"github.com/devchuckcamp/gocommerce/orders"
)

// OrderRepository implements orders.Repository using GORM
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// FindByID finds an order by ID
func (r *OrderRepository) FindByID(ctx context.Context, id string) (*orders.Order, error) {
	var dbOrder database.Order
	if err := r.db.WithContext(ctx).First(&dbOrder, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, orders.ErrOrderNotFound
		}
		return nil, err
	}

	return r.toDomain(&dbOrder)
}

// FindByOrderNumber finds an order by order number
func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*orders.Order, error) {
	var dbOrder database.Order
	if err := r.db.WithContext(ctx).First(&dbOrder, "order_number = ?", orderNumber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, orders.ErrOrderNotFound
		}
		return nil, err
	}

	return r.toDomain(&dbOrder)
}

// FindByUserID finds orders by user ID
func (r *OrderRepository) FindByUserID(ctx context.Context, userID string, filter orders.OrderFilter) ([]*orders.Order, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	query = r.applyFilter(query, filter)

	var dbOrders []database.Order
	if err := query.Find(&dbOrders).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbOrders)
}

// Save saves an order
func (r *OrderRepository) Save(ctx context.Context, order *orders.Order) error {
	dbOrder := r.toDatabase(order)
	return r.db.WithContext(ctx).Save(dbOrder).Error
}

// Delete deletes an order
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Order{}, "id = ?", id).Error
}

// Helper methods

func (r *OrderRepository) applyFilter(query *gorm.DB, filter orders.OrderFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	query = query.Order("created_at DESC")
	return query
}

func (r *OrderRepository) toDomain(dbOrder *database.Order) (*orders.Order, error) {
	var items []orders.OrderItem
	if err := database.UnmarshalJSON(dbOrder.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
	}

	var shippingAddress orders.Address
	if err := database.UnmarshalJSON(dbOrder.ShippingAddress, &shippingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipping address: %w", err)
	}

	var billingAddress orders.Address
	if err := database.UnmarshalJSON(dbOrder.BillingAddress, &billingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal billing address: %w", err)
	}

	return &orders.Order{
		ID:              dbOrder.ID,
		OrderNumber:     dbOrder.OrderNumber,
		UserID:          dbOrder.UserID,
		Status:          orders.OrderStatus(dbOrder.Status),
		Items:           items,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		PaymentMethodID: dbOrder.PaymentMethodID,
		Subtotal:        database.Int64ToMoney(dbOrder.Subtotal, dbOrder.Currency),
		DiscountTotal:   database.Int64ToMoney(dbOrder.DiscountTotal, dbOrder.Currency),
		TaxTotal:        database.Int64ToMoney(dbOrder.TaxTotal, dbOrder.Currency),
		ShippingTotal:   database.Int64ToMoney(dbOrder.ShippingTotal, dbOrder.Currency),
		Total:           database.Int64ToMoney(dbOrder.Total, dbOrder.Currency),
		Notes:           dbOrder.Notes,
		IPAddress:       dbOrder.IPAddress,
		UserAgent:       dbOrder.UserAgent,
		CanceledAt:      dbOrder.CancelledAt,
		CreatedAt:       dbOrder.CreatedAt,
		UpdatedAt:       dbOrder.UpdatedAt,
	}, nil
}

func (r *OrderRepository) toDomainList(dbOrders []database.Order) ([]*orders.Order, error) {
	ordersList := make([]*orders.Order, 0, len(dbOrders))
	for _, dbOrder := range dbOrders {
		order, err := r.toDomain(&dbOrder)
		if err != nil {
			return nil, err
		}
		ordersList = append(ordersList, order)
	}
	return ordersList, nil
}

func (r *OrderRepository) toDatabase(order *orders.Order) *database.Order {
	return &database.Order{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		UserID:          order.UserID,
		Status:          string(order.Status),
		Items:           database.MarshalJSON(order.Items),
		ShippingAddress: database.MarshalJSON(order.ShippingAddress),
		BillingAddress:  database.MarshalJSON(order.BillingAddress),
		PaymentMethodID: order.PaymentMethodID,
		Subtotal:        database.MoneyToInt64(order.Subtotal),
		DiscountTotal:   database.MoneyToInt64(order.DiscountTotal),
		TaxTotal:        database.MoneyToInt64(order.TaxTotal),
		ShippingTotal:   database.MoneyToInt64(order.ShippingTotal),
		Total:           database.MoneyToInt64(order.Total),
		Currency:        order.Total.Currency,
		Notes:           order.Notes,
		IPAddress:       order.IPAddress,
		UserAgent:       order.UserAgent,
		CancelledAt:     order.CanceledAt,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}
