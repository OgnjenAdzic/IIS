package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PricingRule struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	BasePrice   float64   `gorm:"not null"`
	PricePerKm  float64   `gorm:"not null"`
	RushHourFee float64   `gorm:"not null;default:0.0"`
	WeatherFee  float64   `gorm:"not null;default:0.0"`
	ValidFrom   time.Time `gorm:"not null"`
	ValidTo     *time.Time
}

type SystemStatus struct {
	Id           int  `gorm:"primaryKey;autoIncrement:false"`
	IsRushHour   bool `gorm:"default:false"`
	IsBadWeather bool `gorm:"default:false"`
}

func (p *PricingRule) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id = uuid.New()
	return
}
