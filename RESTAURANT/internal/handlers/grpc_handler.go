package handler

import (
	"context"
	"restaurant/internal/models"
	"restaurant/internal/service"

	pb "COMMON/restaurant/proto"

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

// Helper to check role
func authorizeAdmin(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "No metadata found")
	}

	roles := md.Get("user-role")
	if len(roles) == 0 {
		return status.Error(codes.Unauthenticated, "No role found")
	}

	if roles[0] != "ADMIN" {
		return status.Error(codes.PermissionDenied, "Admin access required")
	}
	return nil
}

func authorizeRestaurantWorker(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "No metadata found")
	}

	roles := md.Get("user-role")
	if len(roles) == 0 {
		return status.Error(codes.Unauthenticated, "No role found")
	}

	if roles[0] != "RESTAURANT_WORKER" {
		return status.Error(codes.PermissionDenied, "Restaurant worker access required")
	}
	return nil
}

func (h *RestaurantHandler) CreateRestaurant(ctx context.Context, req *pb.CreateRestaurantRequest) (*pb.RestaurantResponse, error) {
	if err := authorizeAdmin(ctx); err != nil {
		return nil, err
	}

	res, err := h.service.CreateRestaurant(req.Name, req.Category, req.Address, req.Latitude, req.Longitude)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create restaurant")
	}
	return mapToProto(res), nil
}

func (h *RestaurantHandler) AddMenuItem(ctx context.Context, req *pb.AddMenuItemRequest) (*pb.MenuItemResponse, error) {
	if err := authorizeRestaurantWorker(ctx); err != nil {
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
	if err := authorizeRestaurantWorker(ctx); err != nil {
		return nil, err
	}

	err := h.service.SoftDeleteMenuItem(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete item")
	}
	return &pb.MenuItemResponse{Id: req.Id}, nil
}

func (h *RestaurantHandler) UpdateItemPrice(ctx context.Context, req *pb.UpdateItemPriceRequest) (*pb.MenuItemResponse, error) {
	if err := authorizeRestaurantWorker(ctx); err != nil {
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
	if err := authorizeRestaurantWorker(ctx); err != nil {
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
	restaurants, err := h.service.GetAll()
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

// Mapper Helper
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
		Menu: &pb.Menu{
			Id:    r.Menu.Id.String(),
			Items: protoItems,
		},
	}
}
