package service

import (
	"analysis/internal/models"
	"analysis/internal/repository"
	"math"
)

type AnalysisService struct {
	repo *repository.AnalysisRepository
}

func NewAnalysisService(repo *repository.AnalysisRepository) *AnalysisService {
	return &AnalysisService{repo: repo}
}

type FeeResult struct {
	AppFee             float64
	SmallOrderFee      float64
	ProfitFromItems    float64 // Added
	ProfitFromDelivery float64 // Added
	EstimatedProfit    float64
}

func (s *AnalysisService) CalculateFees(itemsTotal, deliveryPrice float64) (*FeeResult, error) {
	config, err := s.repo.GetActiveConfig()
	if err != nil {
		return nil, err
	}

	rawAppFee := itemsTotal * (config.AppFeePercent / 100.0)
	appFee := math.Min(rawAppFee, config.AppFeeCap)

	smallOrderFee := 0.0
	if itemsTotal < config.SmallOrderThreshold {
		smallOrderFee = config.SmallOrderFee
	}

	profitItems := itemsTotal * (config.ItemRevenuePercent / 100.0)
	profitDelivery := deliveryPrice * (config.DeliveryRevenuePercent / 100.0)

	totalProfit := profitItems + profitDelivery + appFee + smallOrderFee

	return &FeeResult{
		AppFee:             math.Round(appFee*100) / 100,
		SmallOrderFee:      smallOrderFee,
		ProfitFromItems:    math.Round(profitItems*100) / 100,    // Return this
		ProfitFromDelivery: math.Round(profitDelivery*100) / 100, // Return this
		EstimatedProfit:    math.Round(totalProfit*100) / 100,
	}, nil
}

func (s *AnalysisService) UpdateConfig(cfg models.FeeConfiguration) (*models.FeeConfiguration, error) {
	return s.repo.CreateNewConfig(cfg)
}

func (s *AnalysisService) GetConfig() (*models.FeeConfiguration, error) {
	return s.repo.GetActiveConfig()
}

type AnalyticsResult struct {
	TotalRevenue   float64
	TopRestaurants []repository.AggregatedResult
	TopUsers       []repository.AggregatedResult
}

func (s *AnalysisService) RecordProfit(req models.OrderProfitLog) error {
	return s.repo.SaveProfitLog(&req)
}

func (s *AnalysisService) GetHistory() ([]models.OrderProfitLog, error) {
	return s.repo.GetProfitHistory()
}

func (s *AnalysisService) GetAnalytics() (*AnalyticsResult, error) {
	total, err := s.repo.GetTotalRevenue()
	if err != nil {
		return nil, err
	}

	restaurants, err := s.repo.GetTopRestaurants(5) // Top 5
	if err != nil {
		return nil, err
	}

	users, err := s.repo.GetTopUsers(5) // Top 5
	if err != nil {
		return nil, err
	}

	return &AnalyticsResult{
		TotalRevenue:   total,
		TopRestaurants: restaurants,
		TopUsers:       users,
	}, nil
}
