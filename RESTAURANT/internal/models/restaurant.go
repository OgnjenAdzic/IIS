package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Restaurant struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Category  string
	IsOpen    bool    `gorm:"default:true"`
	Address   string  `gorm:"not null"`
	Latitude  float64 `gorm:"type:decimal(10,8);not null"`
	Longitude float64 `gorm:"type:decimal(11,8);not null"`
	Menu      Menu    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Menu struct {
	Id           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	RestaurantId uuid.UUID  `gorm:"type:uuid;not null"`
	Items        []MenuItem `gorm:"foreignKey:MenuId"`
}

type MenuItem struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	MenuId    uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	IsDeleted bool      `gorm:"default:false"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	r.Id = uuid.New()
	return
}
func (m *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = uuid.New()
	return
}
func (mi *MenuItem) BeforeCreate(tx *gorm.DB) (err error) {
	mi.Id = uuid.New()
	return
}
