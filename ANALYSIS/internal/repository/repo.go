package repository

import (
	"analysis/internal/models"
	"time"

	"gorm.io/gorm"
)

type AnalysisRepository struct {
	DB *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) *AnalysisRepository {
	var count int64
	db.Model(&models.FeeConfiguration{}).Count(&count)
	if count == 0 {
		db.Create(&models.FeeConfiguration{
			ItemRevenuePercent:     10.0,
			DeliveryRevenuePercent: 5.0,
			AppFeePercent:          6.0,
			AppFeeCap:              250.0,
			SmallOrderThreshold:    599.0,
			SmallOrderFee:          149.0,
			ValidFrom:              time.Now(),
		})
	}
	return &AnalysisRepository{DB: db}
}

func (r *AnalysisRepository) GetActiveConfig() (*models.FeeConfiguration, error) {
	var config models.FeeConfiguration
	err := r.DB.Where("valid_to IS NULL").First(&config).Error
	return &config, err
}

func (r *AnalysisRepository) CreateNewConfig(newConfig models.FeeConfiguration) (*models.FeeConfiguration, error) {
	return &newConfig, r.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Close old config
		if err := tx.Model(&models.FeeConfiguration{}).
			Where("valid_to IS NULL").
			Update("valid_to", time.Now()).Error; err != nil {
			return err
		}
		// 2. Create new config
		newConfig.ValidFrom = time.Now()
		newConfig.ValidTo = nil
		return tx.Create(&newConfig).Error
	})
}

func (r *AnalysisRepository) SaveProfitLog(log *models.OrderProfitLog) error {
	return r.DB.Create(log).Error
}

// 2. Struct for raw SQL results
type AggregatedResult struct {
	Id    string
	Total float64
	Count int
}

// 3. Get Total Revenue (Sum of everything)
func (r *AnalysisRepository) GetTotalRevenue() (float64, error) {
	var total float64
	// COALESCE handles null if table is empty
	err := r.DB.Model(&models.OrderProfitLog{}).Select("COALESCE(SUM(total_profit), 0)").Scan(&total).Error
	return total, err
}

// 4. Get Top Restaurants by Profit
func (r *AnalysisRepository) GetTopRestaurants(limit int) ([]AggregatedResult, error) {
	var results []AggregatedResult
	err := r.DB.Model(&models.OrderProfitLog{}).
		Select("restaurant_id as id, SUM(total_profit) as total, COUNT(*) as count").
		Group("restaurant_id").
		Order("total DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

// 5. Get Top Users by Profit generated
func (r *AnalysisRepository) GetTopUsers(limit int) ([]AggregatedResult, error) {
	var results []AggregatedResult
	err := r.DB.Model(&models.OrderProfitLog{}).
		Select("user_id as id, SUM(total_profit) as total, COUNT(*) as count").
		Group("user_id").
		Order("total DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}
