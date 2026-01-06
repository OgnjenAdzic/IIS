package repository

import (
	"order/internal/models"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) GetById(id string) (*models.Order, error) {
	var order models.Order
	err := r.DB.Preload("Items").Where("id = ?", id).First(&order).Error
	return &order, err
}

func (r *OrderRepository) GetByCustomer(userId string) ([]models.Order, error) {
	var orders []models.Order
	err := r.DB.Preload("Items").Where("customer_id = ?", userId).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateStatus(orderId string, status models.OrderStatus, deliveryPersonId *string, readyAt *time.Time) error {
	updates := map[string]interface{}{"status": status}

	if status == models.StatusDelivered {
		updates["delivered_at"] = time.Now()
	}
	if deliveryPersonId != nil {
		updates["delivery_person_id"] = *deliveryPersonId
	}
	if readyAt != nil {
		updates["estimated_ready_at"] = *readyAt
	}

	return r.DB.Model(&models.Order{}).Where("id = ?", orderId).Updates(updates).Error
}

func (r *OrderRepository) GetByRestaurant(restaurantId string, status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	query := r.DB.Preload("Items").Where("restaurant_id = ?", restaurantId)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at desc").Find(&orders).Error
	return orders, err
}
