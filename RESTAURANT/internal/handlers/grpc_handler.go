package handler

import (
	"context"
	"restaurant/internal/models"
	"restaurant/internal/service"

	pb "COMMON/restaurant/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type RestaurantHandler struct {
	pb.UnimplementedRestaurantServiceServer
	service *service.RestaurantService
}

func NewRestaurantHandler(service *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{service: service}
}

func (h *RestaurantHandler) CreateRestaurant(ctx context.Context, req *pb.CreateRestaurantRequest) (*pb.RestaurantResponse, error) {
	if err := h.authorizeAdmin(ctx); err != nil {
		return nil, err
	}

	managerUUID, err := uuid.Parse(req.ManagerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Manager ID")
	}

	res, err := h.service.CreateRestaurant(req.Name, req.Category, req.Address, req.Latitude, req.Longitude, managerUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create restaurant")
	}
	return mapToProto(res), nil
}

func (h *RestaurantHandler) AddMenuItem(ctx context.Context, req *pb.AddMenuItemRequest) (*pb.MenuItemResponse, error) {
	if err := h.validateManagerAccess(ctx, req.RestaurantId); err != nil {
		return nil, err
	}

	err := h.service.AddMenuItem(req.RestaurantId, req.Name, req.Price)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to add item")
	}

	return &pb.MenuItemResponse{
		Name:  req.Name,
		Price: req.Price,
	}, nil
}

func (h *RestaurantHandler) DeleteMenuItem(ctx context.Context, req *pb.DeleteMenuItemRequest) (*pb.MenuItemResponse, error) {
	restaurantId, err := h.service.GetRestaurantIdByItem(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Item not found")
	}

	if err := h.validateManagerAccess(ctx, restaurantId); err != nil {
		return nil, err
	}

	err = h.service.SoftDeleteMenuItem(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete item")
	}
	return &pb.MenuItemResponse{Id: req.Id}, nil
}

func (h *RestaurantHandler) UpdateItemPrice(ctx context.Context, req *pb.UpdateItemPriceRequest) (*pb.MenuItemResponse, error) {
	restaurantId, err := h.service.GetRestaurantIdByItem(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Item not found")
	}

	if err := h.validateManagerAccess(ctx, restaurantId); err != nil {
		return nil, err
	}

	item, err := h.service.UpdateItemPrice(req.Id, req.Price)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update price")
	}

	return &pb.MenuItemResponse{
		Id:    item.Id.String(),
		Name:  item.Name,
		Price: item.Price,
	}, nil
}

func (h *RestaurantHandler) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.RestaurantResponse, error) {
	if err := h.validateManagerAccess(ctx, req.Id); err != nil {
		return nil, err
	}
	res, err := h.service.UpdateStatus(req.Id, req.IsOpen)
	if err != nil {
		return nil, err
	}
	return mapToProto(res), nil
}

// 4. Get All (Public)
func (h *RestaurantHandler) GetAllRestaurants(ctx context.Context, req *pb.GetAllRestaurantsRequest) (*pb.GetAllRestaurantsResponse, error) {
	userId, role, err := getUserMetadata(ctx)
	if err != nil {
		return nil, err
	}

	var restaurants []models.Restaurant

	if role == "RESTAURANT_WORKER" {
		restaurants, err = h.service.GetByManagerId(userId)
	} else {
		restaurants, err = h.service.GetAll()
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch")
	}

	var protoRestaurants []*pb.RestaurantResponse
	for _, r := range restaurants {
		protoRestaurants = append(protoRestaurants, mapToProto(&r))
	}
	return &pb.GetAllRestaurantsResponse{Restaurants: protoRestaurants}, nil
}

// 5. Get One (Public)
func (h *RestaurantHandler) GetRestaurant(ctx context.Context, req *pb.GetRestaurantRequest) (*pb.RestaurantResponse, error) {
	res, err := h.service.GetById(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Restaurant not found")
	}
	return mapToProto(res), nil
}

func getUserMetadata(ctx context.Context) (string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", status.Error(codes.Unauthenticated, "No metadata found")
	}

	// Extract Role
	roles := md.Get("user-role")
	if len(roles) == 0 {
		return "", "", status.Error(codes.Unauthenticated, "No role found")
	}
	role := roles[0]

	// Extract ID
	ids := md.Get("user-id")
	if len(ids) == 0 {
		return "", "", status.Error(codes.Unauthenticated, "No user ID found")
	}
	userId := ids[0]

	return userId, role, nil
}

// Helper to check role
func (h *RestaurantHandler) authorizeAdmin(ctx context.Context) error {
	_, role, err := getUserMetadata(ctx)
	if err != nil {
		return err
	}

	if role != "ADMIN" {
		return status.Error(codes.PermissionDenied, "Admin access required")
	}
	return nil
}

func (h *RestaurantHandler) validateManagerAccess(ctx context.Context, restaurantId string) error {
	userId, role, err := getUserMetadata(ctx)
	if err != nil {
		return err
	}

	// Rule 1: Admins can manage everything
	if role == "ADMIN" {
		return nil
	}

	// Rule 2: Managers can only manage their own restaurant
	restaurant, err := h.service.GetById(restaurantId)
	if err != nil {
		return status.Error(codes.NotFound, "Restaurant not found")
	}

	if restaurant.ManagerId.String() != userId {
		return status.Error(codes.PermissionDenied, "You are not the manager of this restaurant")
	}

	return nil
}
func mapToProto(r *models.Restaurant) *pb.RestaurantResponse {
	// Map items
	var protoItems []*pb.MenuItem
	for _, item := range r.Menu.Items {
		protoItems = append(protoItems, &pb.MenuItem{
			Id:    item.Id.String(),
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &pb.RestaurantResponse{
		Id:        r.Id.String(),
		Name:      r.Name,
		Category:  r.Category,
		IsOpen:    r.IsOpen,
		Address:   r.Address,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
		ManagerId: r.ManagerId.String(),
		Menu: &pb.Menu{
			Id:    r.Menu.Id.String(),
			Items: protoItems,
		},
	}
}
