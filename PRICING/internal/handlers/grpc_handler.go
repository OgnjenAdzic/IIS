package handlers

import (
	pb "COMMON/pricing/proto"
	"context"
	"pricing/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PricingHandler struct {
	pb.UnimplementedPricingServiceServer
	service *service.PricingService
}

func NewPricingHandler(service *service.PricingService) *PricingHandler {
	return &PricingHandler{service: service}
}

func (h *PricingHandler) CalculatePrice(ctx context.Context, req *pb.CalculatePriceRequest) (*pb.CalculatePriceResponse, error) {
	res, err := h.service.CalculatePrice(
		req.CustomerLat, req.CustomerLon,
		req.RestaurantLat, req.RestaurantLon,
		req.IsPriority,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "Calculation failed")
	}

	return &pb.CalculatePriceResponse{
		FinalPrice:    res.FinalPrice,
		DistanceKm:    res.DistanceKm,
		BasePrice:     res.BasePrice,
		DistancePrice: res.DistancePrice,
		RushHourFee:   res.RushHourFee,
		WeatherFee:    res.WeatherFee,
	}, nil
}

func (h *PricingHandler) UpdatePricingRules(ctx context.Context, req *pb.UpdateRulesRequest) (*pb.PricingRulesResponse, error) {
	if err := h.authorizeAdmin(ctx); err != nil {
		return nil, err
	}
	rule, err := h.service.UpdateRules(req.BasePrice, req.PricePerKm, req.RushHourFee, req.WeatherFee)
	if err != nil {
		return nil, status.Error(codes.Internal, "Update failed")
	}

	return &pb.PricingRulesResponse{
		Id:          rule.Id.String(),
		BasePrice:   rule.BasePrice,
		PricePerKm:  rule.PricePerKm,
		RushHourFee: rule.RushHourFee,
		WeatherFee:  rule.WeatherFee,
	}, nil
}

func (h *PricingHandler) UpdateSystemStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.SystemStatusResponse, error) {
	// TODO: Add Admin Authorization Check here
	st, err := h.service.UpdateStatus(req.IsRushHour, req.IsBadWeather)
	if err != nil {
		return nil, status.Error(codes.Internal, "Update failed")
	}

	return &pb.SystemStatusResponse{
		IsRushHour:   st.IsRushHour,
		IsBadWeather: st.IsBadWeather,
	}, nil
}

func (h *PricingHandler) GetCurrentConfig(ctx context.Context, _ *pb.GetConfigRequest) (*pb.ConfigResponse, error) {
	rule, st, err := h.service.GetConfig()
	if err != nil {
		return nil, status.Error(codes.Internal, "Fetch failed")
	}

	return &pb.ConfigResponse{
		Rules:  &pb.PricingRulesResponse{BasePrice: rule.BasePrice, PricePerKm: rule.PricePerKm, RushHourFee: rule.RushHourFee, WeatherFee: rule.WeatherFee},
		Status: &pb.SystemStatusResponse{IsRushHour: st.IsRushHour, IsBadWeather: st.IsBadWeather},
	}, nil
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

func (h *PricingHandler) authorizeAdmin(ctx context.Context) error {
	_, role, err := getUserMetadata(ctx)
	if err != nil {
		return err
	}

	if role != "ADMIN" {
		return status.Error(codes.PermissionDenied, "Admin access required")
	}
	return nil
}
