package handler

import (
	pb "COMMON/shopping_cart/proto"
	"context"
	"shopping_cart/internal/models"
	"shopping_cart/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type CartHandler struct {
	pb.UnimplementedShoppingCartServiceServer
	service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// Helper to extract User ID from JWT Metadata
func getUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "No metadata found")
	}
	ids := md.Get("user-id")
	if len(ids) == 0 {
		return "", status.Error(codes.Unauthenticated, "User ID missing")
	}
	return ids[0], nil
}

func (h *CartHandler) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.CartResponse, error) {
	userId, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := h.service.GetCart(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "DB Error")
	}
	// If no cart exists yet, return empty
	if cart == nil {
		return &pb.CartResponse{UserId: userId}, nil
	}

	return mapToProto(cart), nil
}

func (h *CartHandler) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.CartResponse, error) {
	userId, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	item := models.CartItem{
		MenuItemId: uuid.MustParse(req.MenuItemId),
		Name:       req.Name,
		Price:      req.Price,
		Quantity:   int(req.Quantity),
	}

	cart, err := h.service.AddItem(userId, req.RestaurantId, item)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to add item")
	}

	return mapToProto(cart), nil
}

// --- IMPLEMENTED METHODS ---

func (h *CartHandler) UpdateItemQuantity(ctx context.Context, req *pb.UpdateQuantityRequest) (*pb.CartResponse, error) {
	userId, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := h.service.UpdateQuantity(userId, req.ItemId, int(req.Quantity))
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update quantity")
	}
	if cart == nil {
		return &pb.CartResponse{UserId: userId}, nil
	}

	return mapToProto(cart), nil
}

func (h *CartHandler) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.CartResponse, error) {
	userId, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := h.service.RemoveItem(userId, req.ItemId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to remove item")
	}
	if cart == nil {
		return &pb.CartResponse{UserId: userId}, nil
	}

	return mapToProto(cart), nil
}

func (h *CartHandler) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.CartResponse, error) {
	userId, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := h.service.ClearCart(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to clear cart")
	}

	// Return empty cart structure
	return mapToProto(cart), nil
}

// --- HELPER MAPPER ---

func mapToProto(c *models.Cart) *pb.CartResponse {
	var total float64
	var items []*pb.CartItem

	for _, i := range c.Items {
		total += i.Price * float64(i.Quantity)
		items = append(items, &pb.CartItem{
			Id:         i.Id.String(),
			MenuItemId: i.MenuItemId.String(),
			Name:       i.Name,
			Price:      i.Price,
			Quantity:   int32(i.Quantity),
		})
	}

	// Handle case where ID might be empty (if it's a temp/cleared cart object)
	cartID := ""
	if c.Id != uuid.Nil {
		cartID = c.Id.String()
	}

	restaurantID := ""
	if c.RestaurantId != uuid.Nil {
		restaurantID = c.RestaurantId.String()
	}

	return &pb.CartResponse{
		Id:           cartID,
		UserId:       c.UserId.String(),
		RestaurantId: restaurantID,
		Items:        items,
		TotalPrice:   total,
	}
}
