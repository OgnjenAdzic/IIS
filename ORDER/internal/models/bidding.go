package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderBid struct {
	Id       uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderId  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_order_driver"`
	DriverId uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_order_driver"`
	Minutes  int       `gorm:"not null"` // Bid amount
}

func (b *OrderBid) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = uuid.New()
	return
}
