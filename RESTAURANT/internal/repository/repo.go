package repository

import (
	"restaurant/internal/models"

	"gorm.io/gorm"
)

type RestaurantRepository struct {
	DB *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{DB: db}
}

func (r *RestaurantRepository) Create(restaurant *models.Restaurant) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Menu").Create(restaurant).Error; err != nil {
			return err
		}

		menu := models.Menu{
			RestaurantId: restaurant.Id,
		}

		if err := tx.Create(&menu).Error; err != nil {
			return err
		}

		restaurant.Menu = menu

		return nil
	})
}

func (r *RestaurantRepository) GetAll() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	// Preload Menu and Items so we see everything
	err := r.DB.Preload("Menu.Items", "is_deleted = ?", false).Find(&restaurants).Error
	return restaurants, err
}

func (r *RestaurantRepository) GetById(id string) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := r.DB.Preload("Menu.Items", "is_deleted = ?", false).Where("id = ?", id).First(&restaurant).Error
	return &restaurant, err
}

func (r *RestaurantRepository) SoftDeleteMenuItem(itemId string) error {
	return r.DB.Model(&models.MenuItem{}).Where("id = ?", itemId).Update("is_deleted", true).Error
}

func (r *RestaurantRepository) UpdateItemPrice(itemId string, price float64) (*models.MenuItem, error) {
	var item models.MenuItem
	if err := r.DB.First(&item, "id = ?", itemId).Error; err != nil {
		return nil, err
	}

	item.Price = price
	if err := r.DB.Save(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RestaurantRepository) UpdateStatus(id string, isOpen bool) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.DB.First(&restaurant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	restaurant.IsOpen = isOpen
	r.DB.Save(&restaurant)
	return &restaurant, nil
}

func (r *RestaurantRepository) AddMenuItem(restaurantId string, item *models.MenuItem) error {
	// 1. Find the menu for this restaurant
	var menu models.Menu
	if err := r.DB.Where("restaurant_id = ?", restaurantId).First(&menu).Error; err != nil {
		return err
	}

	// 2. Assign MenuId to the item
	item.MenuId = menu.Id

	// 3. Save Item
	return r.DB.Create(item).Error
}
