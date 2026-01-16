package handler

import (
	pb "COMMON/order/proto"
	"context"
	"fmt"
	"order/internal/models"
	"order/internal/service"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// Helper: Get User ID from JWT Metadata
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

func (h *OrderHandler) getAurothorizationRole(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "No metadata found")
	}
	roles := md.Get("user-role")
	if len(roles) == 0 {
		return status.Error(codes.Unauthenticated, "User role missing")
	}
	if roles[0] != "RESTAURANT_WORKER" && roles[0] != "DELIVERY_PERSON" {
		return status.Error(codes.PermissionDenied, "Admin access required")
	}
	return nil
}

// 1. CREATE ORDER
func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	fmt.Println(">>> [ORDER] Received CreateOrder Request") // DEBUG

	userId, err := getUserID(ctx)
	if err != nil {
		fmt.Printf(">>> [ORDER] Auth Error: %v\n", err) // DEBUG
		return nil, err
	}
	fmt.Printf(">>> [ORDER] User ID: %s\n", userId) // DEBUG

	var customAddr *string
	if req.CustomAddress != "" {
		customAddr = &req.CustomAddress
	}

	order, err := h.service.CreateOrder(
		ctx,
		userId,
		customAddr,
		req.CustomLat,
		req.CustomLon,
		req.IsPriority,
	)

	if err != nil {
		// THIS IS THE IMPORTANT LOG
		fmt.Printf(">>> [ORDER] Service Failed: %v\n", err)
		return nil, status.Errorf(codes.Internal, "Order creation failed: %v", err)
	}

	fmt.Println(">>> [ORDER] Success!")
	return mapToProto(order), nil
}

// 2. GET MY ORDERS
func (h *OrderHandler) GetMyOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	// If the request userId is empty, try to get from token
	reqId := req.UserId
	if reqId == "" {
		id, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}
		reqId = id
	}

	orders, err := h.service.GetMyOrders(reqId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch orders")
	}

	var protoOrders []*pb.OrderResponse
	for _, o := range orders {
		protoOrders = append(protoOrders, mapToProto(&o))
	}

	return &pb.GetOrdersResponse{Orders: protoOrders}, nil
}

func (h *OrderHandler) GetRestaurantOrders(ctx context.Context, req *pb.GetRestaurantOrdersRequest) (*pb.GetOrdersResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "No metadata found")
	}
	roles := md.Get("user-role")
	if len(roles) == 0 {
		return nil, status.Error(codes.Unauthenticated, "User role missing")
	}
	if roles[0] != "RESTAURANT_WORKER" {
		return nil, status.Error(codes.PermissionDenied, "Restaurant worker access required")
	}

	orders, err := h.service.GetOrdersByRestaurant(req.RestaurantId, (models.OrderStatus)(req.Status))
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch orders")
	}

	var protoOrders []*pb.OrderResponse
	for _, o := range orders {
		protoOrders = append(protoOrders, mapToProto(&o))
	}

	return &pb.GetOrdersResponse{Orders: protoOrders}, nil
}

// 3. UPDATE STATUS
func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.OrderResponse, error) {
	if err := h.getAurothorizationRole(ctx); err != nil {
		return nil, err
	}
	// Convert Proto Enum to String
	statusStr := models.StatusPending
	switch req.Status {
	case pb.OrderStatus_PREPARING:
		statusStr = models.StatusPreparing
	case pb.OrderStatus_READY_FOR_PICKUP:
		statusStr = models.StatusReady
	case pb.OrderStatus_IN_DELIVERY:
		statusStr = models.StatusInDelivery
	case pb.OrderStatus_DELIVERED:
		statusStr = models.StatusDelivered
	case pb.OrderStatus_CANCELLED:
		statusStr = models.StatusCancelled
	}

	var deliveryPersonId *string
	if req.DeliveryPersonId != "" {
		deliveryPersonId = &req.DeliveryPersonId
	}

	fmt.Printf(">>> [DEBUG] Request Minutes: %d\n", req.FoodReadyAt) // DEBUG 1

	order, err := h.service.UpdateStatus(req.OrderId, statusStr, deliveryPersonId, req.FoodReadyAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "Update failed")
	}

	fmt.Printf(">>> [DEBUG] DB Object Time: %v\n", order.FoodReadyAt) // DEBUG 2

	response := mapToProto(order)

	fmt.Printf(">>> [DEBUG] Proto Response Time: %d\n", response.PreparingFoodDeliveryTime) // DEBUG 3

	return response, nil
}

func (h *OrderHandler) BidForOrder(ctx context.Context, req *pb.BidRequest) (*pb.BidResponse, error) {
	driverId, err := getUserID(ctx) // From Token
	if err != nil {
		fmt.Println(">>> [DEBUG Bid] Error :", err)
		return nil, err

	}
	fmt.Println(">>> [DEBUG Bid] Driver ID:", driverId)

	err = h.service.PlaceBid(ctx, req.OrderId, driverId, int(req.Minutes))
	if err != nil {
		return &pb.BidResponse{Success: false, Message: "Failed to bid"}, nil
	}
	return &pb.BidResponse{Success: true, Message: "Bid placed"}, nil
}

func (h *OrderHandler) GetAvailableOrders(ctx context.Context, _ *pb.GetAvailableOrdersRequest) (*pb.GetOrdersResponse, error) {
	orders, err := h.service.GetAvailableOrders()
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed")
	}

	var protoOrders []*pb.OrderResponse
	for _, o := range orders {
		protoOrders = append(protoOrders, mapToProto(&o))
	}
	return &pb.GetOrdersResponse{Orders: protoOrders}, nil
}

func (h *OrderHandler) GetMyActiveJob(ctx context.Context, _ *pb.GetOrdersRequest) (*pb.OrderResponse, error) {
	driverId, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(">>> [DEBUG Active Jobs] Fetching active job for driver ID:", driverId) // DEBUG

	order, err := h.service.GetActiveJob(driverId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get active job")
	}

	if order == nil {
		return nil, status.Error(codes.NotFound, "No active job found")
	}

	return mapToProto(order), nil
}

// MAPPER
func mapToProto(o *models.Order) *pb.OrderResponse {
	var items []*pb.OrderItem
	for _, i := range o.Items {
		items = append(items, &pb.OrderItem{
			Name:     i.Name,
			Price:    i.Price,
			Quantity: int32(i.Quantity),
		})
	}

	dId := ""
	if o.DeliveryPersonId != nil {
		dId = o.DeliveryPersonId.String()
	}

	return &pb.OrderResponse{
		Id:                        o.Id.String(),
		RestaurantId:              o.RestaurantId.String(),
		CustomerId:                o.CustomerId.String(),
		DeliveryPersonId:          dId,
		Status:                    string(o.Status),
		DeliveryAddress:           o.DeliveryAddress,
		ItemsTotal:                o.ItemsTotal,
		DeliveryFee:               o.DeliveryFee,
		AppFee:                    o.AppFee,
		SmallOrderFee:             o.SmallOrderFee,
		TotalPrice:                o.TotalPrice,
		Items:                     items,
		CreatedAt:                 o.CreatedAt.Format(time.RFC3339),
		IsPriority:                o.IsPriority,
		PreparingFoodDeliveryTime: int32(o.FoodReadyAt),
		DeliveryDurationTime:      int32(o.DeliveryDuration),
	}
}
