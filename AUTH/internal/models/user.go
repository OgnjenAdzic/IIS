package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin            Role = "ADMIN"
	RoleCustomer         Role = "CUSTOMER"
	RoleDeliveryPerson   Role = "DELIVERY_PERSON"
	RoleRestaurantWorker Role = "RESTAURANT_WORKER"
)

type User struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Username string    `json:"username" gorm:"unique;not null"`
	Password string    `json:"password" gorm:"not null"`
	Role     Role      `json:"role" gorm:"not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New()
	return
}
