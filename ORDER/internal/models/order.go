package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusPreparing  OrderStatus = "PREPARING"
	StatusReady      OrderStatus = "READY_FOR_PICKUP"
	StatusInDelivery OrderStatus = "IN_DELIVERY"
	StatusDelivered  OrderStatus = "DELIVERED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

type Order struct {
	Id               uuid.UUID  `gorm:"type:uuid;primaryKey"`
	CustomerId       uuid.UUID  `gorm:"type:uuid;not null"`
	RestaurantId     uuid.UUID  `gorm:"type:uuid;not null"`
	DeliveryPersonId *uuid.UUID `gorm:"type:uuid"` // Nullable initially

	Status OrderStatus `gorm:"default:'PENDING'"`

	DeliveryAddress string  `gorm:"not null"`
	DeliveryLat     float64 `gorm:"type:decimal(10,8)"`
	DeliveryLon     float64 `gorm:"type:decimal(11,8)"`

	ItemsTotal    float64
	DeliveryFee   float64
	AppFee        float64
	SmallOrderFee float64
	TotalPrice    float64

	IsPriority bool `gorm:"default:false"`

	ProfitFromItems    float64
	ProfitFromDelivery float64
	TotalProfit        float64

	Items []OrderItem `gorm:"foreignKey:OrderId"`

	CreatedAt   time.Time
	DeliveredAt *time.Time

	EstimatedReadyAt *time.Time
}

type OrderItem struct {
	Id       uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderId  uuid.UUID `gorm:"type:uuid;not null"`
	Name     string
	Price    float64
	Quantity int
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.Id = uuid.New()
	return
}
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	oi.Id = uuid.New()
	return
}
