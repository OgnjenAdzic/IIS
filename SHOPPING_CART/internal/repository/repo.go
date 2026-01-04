package repository

import (
	"errors"
	"shopping_cart/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{DB: db}
}

// Get or Create Cart
func (r *CartRepository) GetCartByUserId(userId string) (*models.Cart, error) {
	var cart models.Cart
	err := r.DB.Preload("Items").Where("user_id = ?", userId).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Empty is fine, we handle it in service
	}
	return &cart, err
}

func (r *CartRepository) AddItem(userId, restaurantId string, item models.CartItem) (*models.Cart, error) {
	var cart models.Cart

	// Transaction to ensure safety
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		// 1. Find Cart
		result := tx.Preload("Items").Where("user_id = ?", userId).First(&cart)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Create New Cart
				cart = models.Cart{
					UserId:       uuid.MustParse(userId),
					RestaurantId: uuid.MustParse(restaurantId),
				}
				if err := tx.Create(&cart).Error; err != nil {
					return err
				}
			} else {
				return result.Error
			}
		}

		// 2. Check Restaurant Conflict
		if cart.RestaurantId.String() != restaurantId {
			// Logic: If user adds from new restaurant, clear old items and switch restaurant
			if err := tx.Where("cart_id = ?", cart.Id).Delete(&models.CartItem{}).Error; err != nil {
				return err
			}
			cart.RestaurantId = uuid.MustParse(restaurantId)
			cart.Items = []models.CartItem{} // Clear in memory
			if err := tx.Save(&cart).Error; err != nil {
				return err
			}
		}

		// 3. Check if Item exists, Update or Create
		var existingItem models.CartItem
		// Check in DB or memory loop. DB check is safer here.
		res := tx.Where("cart_id = ? AND menu_item_id = ?", cart.Id, item.MenuItemId).First(&existingItem)

		if res.Error == nil {
			// Update Quantity
			existingItem.Quantity += item.Quantity
			tx.Save(&existingItem)
		} else {
			// Create New
			item.CartId = cart.Id
			tx.Create(&item)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return updated cart
	return r.GetCartByUserId(userId)
}

func (r *CartRepository) UpdateQuantity(itemId string, quantity int) error {
	if quantity <= 0 {
		return r.DB.Delete(&models.CartItem{}, "id = ?", itemId).Error
	}
	return r.DB.Model(&models.CartItem{}).Where("id = ?", itemId).Update("quantity", quantity).Error
}

func (r *CartRepository) RemoveItem(itemId string) error {
	return r.DB.Delete(&models.CartItem{}, "id = ?", itemId).Error
}

func (r *CartRepository) ClearCart(userId string) error {
	return r.DB.Where("user_id = ?", userId).Delete(&models.Cart{}).Error
}
