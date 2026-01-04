package repository

import (
	"pricing/internal/models"
	"time"

	"gorm.io/gorm"
)

type PricingRepository struct {
	DB *gorm.DB
}

func NewPricingRepository(db *gorm.DB) *PricingRepository {
	// Initialize System Status row if not exists
	var count int64
	db.Model(&models.SystemStatus{}).Count(&count)
	if count == 0 {
		db.Create(&models.SystemStatus{Id: 1, IsRushHour: false, IsBadWeather: false})
	}
	// Initialize Default Rule if not exists
	db.Model(&models.PricingRule{}).Count(&count)
	if count == 0 {
		db.Create(&models.PricingRule{
			BasePrice:   150,
			PricePerKm:  50,
			RushHourFee: 150,
			WeatherFee:  125,
			ValidFrom:   time.Now(),
		})
	}
	return &PricingRepository{DB: db}
}

// Get the Rule that is currently active (ValidTo is NULL)
func (r *PricingRepository) GetActiveRule() (*models.PricingRule, error) {
	var rule models.PricingRule
	err := r.DB.Where("valid_to IS NULL").First(&rule).Error
	return &rule, err
}

// Get Toggles
func (r *PricingRepository) GetSystemStatus() (*models.SystemStatus, error) {
	var status models.SystemStatus
	err := r.DB.First(&status, 1).Error
	return &status, err
}

// Update Rules (History Logic)
func (r *PricingRepository) CreateNewRule(newRule models.PricingRule) (*models.PricingRule, error) {
	return &newRule, r.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Close current rule
		if err := tx.Model(&models.PricingRule{}).
			Where("valid_to IS NULL").
			Update("valid_to", time.Now()).Error; err != nil {
			return err
		}

		// 2. Insert new rule
		newRule.ValidFrom = time.Now()
		newRule.ValidTo = nil
		if err := tx.Create(&newRule).Error; err != nil {
			return err
		}
		return nil
	})
}

// Toggle Switches
func (r *PricingRepository) UpdateStatus(rushHour, badWeather bool) (*models.SystemStatus, error) {
	status := models.SystemStatus{Id: 1, IsRushHour: rushHour, IsBadWeather: badWeather}
	err := r.DB.Save(&status).Error
	return &status, err
}
