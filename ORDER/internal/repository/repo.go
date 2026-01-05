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

func (r *OrderRepository) UpdateStatus(orderId string, status models.OrderStatus, deliveryPersonId *string) error {
	updates := map[string]interface{}{"status": status}

	if status == models.StatusDelivered {
		updates["delivered_at"] = time.Now()
	}
	if deliveryPersonId != nil {
		updates["delivery_person_id"] = *deliveryPersonId
	}

	return r.DB.Model(&models.Order{}).Where("id = ?", orderId).Updates(updates).Error
}
