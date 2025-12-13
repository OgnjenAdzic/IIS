package models

import (
	"github.com/google/uuid"
)

type VehicleType string

const (
	VehicleCar     VehicleType = "CAR"
	VehicleBike    VehicleType = "BIKE"
	VehicleScooter VehicleType = "SCOOTER"
)

// Table: customers
type Customer struct {
	// PK is UserId (from Auth)
	UserId    uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`

	// Address Info
	Address   string  `gorm:"not null"`
	Latitude  float64 `gorm:"type:decimal(10,8);not null"`
	Longitude float64 `gorm:"type:decimal(11,8);not null"`
}

// Table: delivery_persons
type DeliveryPerson struct {
	// PK is UserId (from Auth)
	UserId    uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`

	// Work Info
	Vehicle       VehicleType `gorm:"not null"`
	IsWorking     bool        `gorm:"default:false"`
	DeliveryCount int         `gorm:"default:0"`
}
