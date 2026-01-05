package handlers

import (
	pb "COMMON/analysis/proto"
	"analysis/internal/models"
	"analysis/internal/service"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AnalysisHandler struct {
	pb.UnimplementedAnalysisServiceServer
	service *service.AnalysisService
}

func NewAnalysisHandler(s *service.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{service: s}
}

func (h *AnalysisHandler) CalculateFees(ctx context.Context, req *pb.CalculationRequest) (*pb.FeeResponse, error) {
	res, err := h.service.CalculateFees(req.ItemsTotal, req.DeliveryPrice)
	if err != nil {
		return nil, status.Error(codes.Internal, "Calculation failed")
	}

	return &pb.FeeResponse{
		AppFee:          res.AppFee,
		SmallOrderFee:   res.SmallOrderFee,
		EstimatedProfit: res.EstimatedProfit,
	}, nil
}

func (h *AnalysisHandler) UpdateFeeConfiguration(ctx context.Context, req *pb.UpdateConfigRequest) (*pb.FeeConfigResponse, error) {
	cfg := models.FeeConfiguration{
		ItemRevenuePercent:     req.ItemRevenuePercent,
		DeliveryRevenuePercent: req.DeliveryRevenuePercent,
		AppFeePercent:          req.AppFeePercent,
		AppFeeCap:              req.AppFeeCap,
		SmallOrderThreshold:    req.SmallOrderThreshold,
		SmallOrderFee:          req.SmallOrderFee,
	}

	newCfg, err := h.service.UpdateConfig(cfg)
	if err != nil {
		return nil, status.Error(codes.Internal, "Update failed")
	}

	return mapToProto(newCfg), nil
}

func (h *AnalysisHandler) GetCurrentConfig(ctx context.Context, _ *pb.GetConfigRequest) (*pb.FeeConfigResponse, error) {
	cfg, err := h.service.GetConfig()
	if err != nil {
		return nil, status.Error(codes.Internal, "Fetch failed")
	}
	return mapToProto(cfg), nil
}

func (h *AnalysisHandler) RecordOrderProfit(ctx context.Context, req *pb.RecordProfitRequest) (*pb.RecordProfitResponse, error) {
	// You might want to authorize this so only INTERNAL services can call it
	logEntry := models.OrderProfitLog{
		OrderId:            uuid.MustParse(req.OrderId),
		RestaurantId:       uuid.MustParse(req.RestaurantId),
		UserId:             uuid.MustParse(req.UserId),
		AppFee:             req.AppFee,
		SmallOrderFee:      req.SmallOrderFee,
		ProfitFromItems:    req.ProfitFromItems,
		ProfitFromDelivery: req.ProfitFromDelivery,
		TotalProfit:        req.TotalProfit,
	}

	err := h.service.RecordProfit(logEntry)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to save log")
	}

	return &pb.RecordProfitResponse{Success: true}, nil
}

func (h *AnalysisHandler) GetAnalytics(ctx context.Context, _ *pb.GetAnalyticsRequest) (*pb.AnalyticsResponse, error) {
	if err := h.authorizeAdmin(ctx); err != nil {
		return nil, err
	}

	data, err := h.service.GetAnalytics()
	if err != nil {
		return nil, status.Error(codes.Internal, "Analytics failed")
	}

	// Map Restaurants
	var topRest []*pb.TopEntity
	for _, r := range data.TopRestaurants {
		topRest = append(topRest, &pb.TopEntity{
			Id:     r.Id,
			Amount: r.Total,
			Count:  int32(r.Count),
		})
	}

	// Map Users
	var topUser []*pb.TopEntity
	for _, u := range data.TopUsers {
		topUser = append(topUser, &pb.TopEntity{
			Id:     u.Id,
			Amount: u.Total,
			Count:  int32(u.Count),
		})
	}

	return &pb.AnalyticsResponse{
		TotalRevenue:   data.TotalRevenue,
		TopRestaurants: topRest,
		TopUsers:       topUser,
	}, nil
}

func mapToProto(c *models.FeeConfiguration) *pb.FeeConfigResponse {
	return &pb.FeeConfigResponse{
		Id:                     c.Id.String(),
		ItemRevenuePercent:     c.ItemRevenuePercent,
		DeliveryRevenuePercent: c.DeliveryRevenuePercent,
		AppFeePercent:          c.AppFeePercent,
		AppFeeCap:              c.AppFeeCap,
		SmallOrderThreshold:    c.SmallOrderThreshold,
		SmallOrderFee:          c.SmallOrderFee,
	}
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

func (h *AnalysisHandler) authorizeAdmin(ctx context.Context) error {
	_, role, err := getUserMetadata(ctx)
	if err != nil {
		return err
	}

	if role != "ADMIN" {
		return status.Error(codes.PermissionDenied, "Admin access required")
	}
	return nil
}
