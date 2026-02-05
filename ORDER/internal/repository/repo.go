package repository

import (
	"errors"
	"fmt"
	"order/internal/models"
	"time"

	"github.com/google/uuid"
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

func (r *OrderRepository) UpdateStatus(orderId string, status models.OrderStatus, deliveryPersonId *string, readyAt int) error {
	updates := map[string]interface{}{"status": status}

	if status == models.StatusDelivered {
		updates["delivered_at"] = time.Now()
	}
	if deliveryPersonId != nil {
		updates["delivery_person_id"] = *deliveryPersonId
	}
	if readyAt > 0 {
		updates["food_ready_at"] = readyAt
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

func (r *OrderRepository) GetAvailableOrders(driverId string) ([]models.Order, error) {
	var orders []models.Order
	err := r.DB.Preload("Items").
		Where("status IN ? AND delivery_person_id IS NULL", []string{string(models.StatusPreparing), string(models.StatusReady)}).
		Order("created_at desc").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	var biddedOrdersIds []uuid.UUID
	r.DB.Model(&models.OrderBid{}).
		Where("driver_id = ?", driverId).
		Pluck("order_id", &biddedOrdersIds)

	bidMap := make(map[uuid.UUID]bool)
	for _, id := range biddedOrdersIds {
		bidMap[id] = true
	}

	for i := range orders {
		if bidMap[orders[i].Id] {
			orders[i].HasCurrentDriverBidded = true
		}
	}
	return orders, err
}

func (r *OrderRepository) CreateBid(bid *models.OrderBid) error {
	fmt.Println(">>> [REPO] Executing INSERT...")
	result := r.DB.Create(bid)
	if result.Error != nil {
		fmt.Printf(">>> [REPO] DB Error: %v\n", result.Error)
		return result.Error
	}
	fmt.Printf(">>> [REPO] Inserted. Rows affected: %d\n", result.RowsAffected)
	return nil
}

func (r *OrderRepository) GetSortedBids(orderId string) ([]models.OrderBid, error) {
	var bids []models.OrderBid
	err := r.DB.Where("order_id = ?", orderId).
		Order("minutes asc, id asc").
		Find(&bids).Error
	return bids, err
}

func (r *OrderRepository) AssignDriver(orderId string, driverId string, minutes int) error {
	result := r.DB.Model(&models.Order{}).
		Where("id = ? AND delivery_person_id IS NULL", orderId).
		Updates(map[string]interface{}{
			"delivery_person_id": driverId,
			"delivery_duration":  minutes,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order already assigned")
	}
	return nil
}

func (r *OrderRepository) GetActiveJobsForDriver(driverId string) (*models.Order, error) {
	var order models.Order
	err := r.DB.Preload("Items").
		Where("delivery_person_id = ? AND status IN ?", driverId, []string{string(models.StatusPreparing), string(models.StatusReady), string(models.StatusInDelivery)}).
		First(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Return nil so Service knows driver is free
	}

	return &order, err
}

func (r *OrderRepository) HasDriverBid(orderId, driverId string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.OrderBid{}).
		Where("order_id = ? AND driver_id = ?", orderId, driverId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OrderRepository) SetBiddingExpiration(orderId string, expiresAt time.Time) error {
	return r.DB.Model(&models.Order{}).Where("id = ?", orderId).Update("bidding_expires_at", expiresAt).Error
}

// Updated logic to find orders ready for assignment
func (r *OrderRepository) GetOrdersReadyForAssignment() ([]models.Order, error) {
	var orders []models.Order

	// Logic:
	// 1. Status is PREPARING or READY
	// 2. Driver is NOT assigned yet
	// 3. BiddingExpiresAt IS NOT NULL (meaning at least one bid happened)
	// 4. BiddingExpiresAt < NOW (meaning time is up)

	now := time.Now().UTC()

	err := r.DB.Where(
		"status IN ? AND delivery_person_id IS NULL AND bidding_expires_at IS NOT NULL AND bidding_expires_at < ?",
		[]string{"PREPARING", "READY_FOR_PICKUP"},
		now,
	).Find(&orders).Error

	return orders, err
}
