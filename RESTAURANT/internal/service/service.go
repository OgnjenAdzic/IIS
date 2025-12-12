package service

import (
	"restaurant/internal/models"
	"restaurant/internal/repository"
)

type RestaurantService struct {
	repo *repository.RestaurantRepository
}

func NewRestaurantService(repo *repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{repo: repo}
}

func (s *RestaurantService) CreateRestaurant(name, category, address string, lat, lon float64) (*models.Restaurant, error) {
	restaurant := &models.Restaurant{
		Name:     name,
		Category: category,
		IsOpen:   true,
		Address:  address,
		Latitude: lat,
	}
	err := s.repo.Create(restaurant)
	return restaurant, err
}

func (s *RestaurantService) GetAll() ([]models.Restaurant, error) {
	return s.repo.GetAll()
}

func (s *RestaurantService) GetById(id string) (*models.Restaurant, error) {
	return s.repo.GetById(id)
}

func (s *RestaurantService) SoftDeleteMenuItem(id string) error {
	return s.repo.SoftDeleteMenuItem(id)
}

func (s *RestaurantService) UpdateItemPrice(id string, price float64) (*models.MenuItem, error) {
	return s.repo.UpdateItemPrice(id, price)
}

func (s *RestaurantService) UpdateStatus(id string, isOpen bool) (*models.Restaurant, error) {
	return s.repo.UpdateStatus(id, isOpen)
}

func (s *RestaurantService) AddMenuItem(restaurantId, name string, price float64) error {
	item := &models.MenuItem{
		Name:  name,
		Price: price,
	}
	return s.repo.AddMenuItem(restaurantId, item)
}
