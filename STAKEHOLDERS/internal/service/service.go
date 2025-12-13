package service

import (
	"stakeholders/internal/models"
	"stakeholders/internal/repository"

	"github.com/google/uuid"
)

type StakeholdersService struct {
	repo *repository.StakeholdersRepository
}

func NewStakeholdersService(repo *repository.StakeholdersRepository) *StakeholdersService {
	return &StakeholdersService{repo: repo}
}

// --- CUSTOMER ---
func (s *StakeholdersService) CreateCustomer(userIdStr, fname, lname, address string, lat, lon float64) (*models.Customer, error) {
	uid, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil, err
	}

	cust := &models.Customer{
		UserId:    uid,
		FirstName: fname,
		LastName:  lname,
		Address:   address,
		Latitude:  lat,
		Longitude: lon,
	}
	err = s.repo.CreateCustomer(cust)
	return cust, err
}

func (s *StakeholdersService) GetCustomer(userId string) (*models.Customer, error) {
	return s.repo.GetCustomer(userId)
}

// --- DELIVERY PERSON ---
func (s *StakeholdersService) CreateDeliveryPerson(userIdStr, fname, lname, vehicleStr string) (*models.DeliveryPerson, error) {
	uid, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil, err
	}

	dp := &models.DeliveryPerson{
		UserId:    uid,
		FirstName: fname,
		LastName:  lname,
		Vehicle:   models.VehicleType(vehicleStr),
		IsWorking: false, // Default
	}
	err = s.repo.CreateDeliveryPerson(dp)
	return dp, err
}

func (s *StakeholdersService) GetDeliveryPerson(userId string) (*models.DeliveryPerson, error) {
	return s.repo.GetDeliveryPerson(userId)
}

func (s *StakeholdersService) UpdateWorkingStatus(userId string, isWorking bool) (*models.DeliveryPerson, error) {
	return s.repo.UpdateWorkingStatus(userId, isWorking)
}
