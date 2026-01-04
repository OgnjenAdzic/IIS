package service

import (
	"shopping_cart/internal/models"
	"shopping_cart/internal/repository"
)

type CartService struct {
	repo *repository.CartRepository
}

func NewCartService(repo *repository.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) GetCart(userId string) (*models.Cart, error) {
	return s.repo.GetCartByUserId(userId)
}

func (s *CartService) AddItem(userId, restaurantId string, item models.CartItem) (*models.Cart, error) {
	return s.repo.AddItem(userId, restaurantId, item)
}

func (s *CartService) UpdateQuantity(userId, itemId string, quantity int) (*models.Cart, error) {
	// 1. Perform Update
	err := s.repo.UpdateQuantity(itemId, quantity)
	if err != nil {
		return nil, err
	}
	// 2. Return fresh cart state
	return s.repo.GetCartByUserId(userId)
}

func (s *CartService) RemoveItem(userId, itemId string) (*models.Cart, error) {
	err := s.repo.RemoveItem(itemId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetCartByUserId(userId)
}

func (s *CartService) ClearCart(userId string) (*models.Cart, error) {
	err := s.repo.ClearCart(userId)
	if err != nil {
		return nil, err
	}
	// Return empty cart object (or nil, depending on preference)
	return &models.Cart{Items: []models.CartItem{}}, nil
}
