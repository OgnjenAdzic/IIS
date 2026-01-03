package repository

import (
	"restaurant/internal/models"
	"time"

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

func (r *RestaurantRepository) UpdateItemPrice(itemId string, newPrice float64) (*models.MenuItem, error) {
	var item models.MenuItem

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&item, "id = ?", itemId).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ItemPriceHistory{}).
			Where("menu_item_id = ? AND valid_to IS NULL", itemId).
			Update("valid_to", time.Now()).Error; err != nil {
			return err
		}

		newHistory := models.ItemPriceHistory{
			MenuItemId: item.Id,
			Price:      newPrice,
			ValidFrom:  time.Now(),
			ValidTo:    nil, // Active
		}
		if err := tx.Create(&newHistory).Error; err != nil {
			return err
		}

		item.Price = newPrice
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		return nil
	})

	return &item, err
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
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var menu models.Menu
		if err := tx.Where("restaurant_id = ?", restaurantId).First(&menu).Error; err != nil {
			return err
		}

		item.MenuId = menu.Id

		// 1. Create the Item (stores current price)
		if err := tx.Create(item).Error; err != nil {
			return err
		}

		// 2. Create the Initial History Record
		history := models.ItemPriceHistory{
			MenuItemId: item.Id,
			Price:      item.Price,
			ValidFrom:  time.Now(),
			ValidTo:    nil, // nil means "Active"
		}

		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *RestaurantRepository) GetByManagerId(managerId string) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.DB.Preload("Menu.Items").Where("manager_id = ?", managerId).Find(&restaurants).Error
	return restaurants, err
}

func (r *RestaurantRepository) GetRestaurantIdByMenuItem(itemId string) (string, error) {
	var item models.MenuItem
	// 1. Get Item to find MenuID
	if err := r.DB.First(&item, "id = ?", itemId).Error; err != nil {
		return "", err
	}

	// 2. Get Menu to find RestaurantID
	var menu models.Menu
	if err := r.DB.First(&menu, "id = ?", item.MenuId).Error; err != nil {
		return "", err
	}

	return menu.RestaurantId.String(), nil
}
