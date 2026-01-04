package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId       uuid.UUID `gorm:"type:uuid;unique;not null"`
	RestaurantId uuid.UUID `gorm:"type:uuid"`

	Items []CartItem `gorm:"foreignKey:CartId;constraint:OnDelete:CASCADE"`
}

type CartItem struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey"`
	CartId     uuid.UUID `gorm:"type:uuid;not null"`
	MenuItemId uuid.UUID `gorm:"type:uuid;not null"`
	Name       string    `gorm:"not null"`
	Price      float64   `gorm:"not null"`
	Quantity   int       `gorm:"not null;default:1"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) (err error) {
	c.Id = uuid.New()
	return
}
func (ci *CartItem) BeforeCreate(tx *gorm.DB) (err error) {
	ci.Id = uuid.New()
	return
}
