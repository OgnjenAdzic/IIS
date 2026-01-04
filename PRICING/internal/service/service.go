package service

import (
	"math"
	"pricing/internal/models"
	"pricing/internal/repository"
)

type PricingService struct {
	repo *repository.PricingRepository
}

func NewPricingService(repo *repository.PricingRepository) *PricingService {
	return &PricingService{repo: repo}
}

type PriceResult struct {
	FinalPrice    float64
	DistanceKm    float64
	BasePrice     float64
	DistancePrice float64
	RushHourFee   float64
	WeatherFee    float64
}

// Haversine Formula to calculate distance between two points in KM
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km

	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (s *PricingService) CalculatePrice(custLat, custLon, restLat, restLon float64, isPriority bool) (*PriceResult, error) {
	// 1. Get Config
	rule, err := s.repo.GetActiveRule()
	if err != nil {
		return nil, err
	}
	status, err := s.repo.GetSystemStatus()
	if err != nil {
		return nil, err
	}

	// 2. Calculate Distance
	distanceKm := calculateDistance(custLat, custLon, restLat, restLon)

	// 3. Logic: Distance Price
	// "if distance is less than 1km we will only charge base price"
	distancePrice := 0.0
	if distanceKm >= 1.0 {
		distancePrice = distanceKm * rule.PricePerKm
	}

	// 4. Logic: Add-ons (Rush Hour & Weather)
	rushFee := 0.0
	weatherFee := 0.0

	if status.IsRushHour {
		rushFee = rule.RushHourFee
	}
	if status.IsBadWeather {
		weatherFee = rule.WeatherFee
	}

	// 5. Subtotal
	total := rule.BasePrice + distancePrice + rushFee + weatherFee

	// 6. Logic: Priority (20% increase)
	if isPriority {
		total = total * 1.20
	}

	return &PriceResult{
		FinalPrice:    math.Round(total*100) / 100, // Round to 2 decimals
		DistanceKm:    math.Round(distanceKm*100) / 100,
		BasePrice:     rule.BasePrice,
		DistancePrice: distancePrice,
		RushHourFee:   rushFee,
		WeatherFee:    weatherFee,
	}, nil
}

// Update Rules (Arguments changed to Fee)
func (s *PricingService) UpdateRules(base, perKm, rushFee, weatherFee float64) (*models.PricingRule, error) {
	rule := models.PricingRule{
		BasePrice:   base,
		PricePerKm:  perKm,
		RushHourFee: rushFee,
		WeatherFee:  weatherFee,
	}
	return s.repo.CreateNewRule(rule)
}

func (s *PricingService) UpdateStatus(isRush, isWeather bool) (*models.SystemStatus, error) {
	return s.repo.UpdateStatus(isRush, isWeather)
}

func (s *PricingService) GetConfig() (*models.PricingRule, *models.SystemStatus, error) {
	r, err := s.repo.GetActiveRule()
	if err != nil {
		return nil, nil, err
	}
	st, err := s.repo.GetSystemStatus()
	return r, st, err
}
