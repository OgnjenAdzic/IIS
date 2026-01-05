package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeeConfiguration struct {
	Id uuid.UUID `gorm:"type:uuid;primaryKey"`

	// Profit Settings (Internal)
	ItemRevenuePercent     float64 `gorm:"not null"`
	DeliveryRevenuePercent float64 `gorm:"not null"`

	// Customer Fees
	AppFeePercent float64 `gorm:"not null"`
	AppFeeCap     float64 `gorm:"not null"`

	SmallOrderThreshold float64 `gorm:"not null"`
	SmallOrderFee       float64 `gorm:"not null"`

	ValidFrom time.Time `gorm:"not null"`
	ValidTo   *time.Time
}

type OrderProfitLog struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderId      uuid.UUID `gorm:"type:uuid;unique;not null"` // Link to Order Service
	RestaurantId uuid.UUID `gorm:"type:uuid;not null"`        // For grouping
	UserId       uuid.UUID `gorm:"type:uuid;not null"`        // For grouping

	// Financials
	AppFee             float64
	SmallOrderFee      float64
	ProfitFromItems    float64
	ProfitFromDelivery float64
	TotalProfit        float64 `gorm:"not null"` // The main metric

	CreatedAt time.Time
}

func (c *FeeConfiguration) BeforeCreate(tx *gorm.DB) (err error) {
	c.Id = uuid.New()
	return
}

func (o *OrderProfitLog) BeforeCreate(tx *gorm.DB) (err error) {
	o.Id = uuid.New()
	return
}
