package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderBid struct {
	Id       uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderId  uuid.UUID `gorm:"type:uuid;not null"`
	DriverId uuid.UUID `gorm:"type:uuid;not null"`
	Minutes  int       `gorm:"not null"` // Bid amount
}

func (b *OrderBid) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = uuid.New()
	return
}
