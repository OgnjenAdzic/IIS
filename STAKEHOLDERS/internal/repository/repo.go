package repository

import (
	"stakeholders/internal/models"

	"gorm.io/gorm"
)

type StakeholdersRepository struct {
	DB *gorm.DB
}

func NewStakeholdersRepository(db *gorm.DB) *StakeholdersRepository {
	return &StakeholdersRepository{DB: db}
}

// --- CUSTOMER ---
func (r *StakeholdersRepository) CreateCustomer(c *models.Customer) error {
	return r.DB.Create(c).Error
}

func (r *StakeholdersRepository) GetCustomer(userId string) (*models.Customer, error) {
	var c models.Customer
	err := r.DB.Where("user_id = ?", userId).First(&c).Error
	return &c, err
}

// --- DELIVERY PERSON ---
func (r *StakeholdersRepository) CreateDeliveryPerson(dp *models.DeliveryPerson) error {
	return r.DB.Create(dp).Error
}

func (r *StakeholdersRepository) GetDeliveryPerson(userId string) (*models.DeliveryPerson, error) {
	var dp models.DeliveryPerson
	err := r.DB.Where("user_id = ?", userId).First(&dp).Error
	return &dp, err
}

func (r *StakeholdersRepository) UpdateWorkingStatus(userId string, isWorking bool) (*models.DeliveryPerson, error) {
	var dp models.DeliveryPerson
	if err := r.DB.First(&dp, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}
	dp.IsWorking = isWorking
	r.DB.Save(&dp)
	return &dp, nil
}
