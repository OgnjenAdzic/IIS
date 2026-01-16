package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"order/internal/models"
	"order/internal/repository"

	// Import Protos from OTHER services
	pbAnalysis "COMMON/analysis/proto"
	pbPricing "COMMON/pricing/proto"
	pbRest "COMMON/restaurant/proto"
	pbCart "COMMON/shopping_cart/proto"
	pbStake "COMMON/stakeholders/proto"

	"github.com/google/uuid"

	"google.golang.org/grpc/metadata"
)

type OrderService struct {
	repo *repository.OrderRepository
	// gRPC Clients
	cartClient     pbCart.ShoppingCartServiceClient
	pricingClient  pbPricing.PricingServiceClient
	analysisClient pbAnalysis.AnalysisServiceClient
	restClient     pbRest.RestaurantServiceClient
	stakeClient    pbStake.StakeholdersServiceClient
}

func NewOrderService(repo *repository.OrderRepository, cartC pbCart.ShoppingCartServiceClient, priceC pbPricing.PricingServiceClient, analC pbAnalysis.AnalysisServiceClient, restC pbRest.RestaurantServiceClient, stakeC pbStake.StakeholdersServiceClient) *OrderService {
	return &OrderService{
		repo:           repo,
		cartClient:     cartC,
		pricingClient:  priceC,
		analysisClient: analC,
		restClient:     restC,
		stakeClient:    stakeC,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userId string, customAddr *string, cLat, cLon float64, isPriority bool) (*models.Order, error) {
	md := metadata.Pairs("user-id", userId, "user-role", "CUSTOMER")
	ctx = metadata.NewOutgoingContext(ctx, md)

	fmt.Println("--- 1. Calling Shopping Cart ---")
	cartRes, err := s.cartClient.GetCart(ctx, &pbCart.GetCartRequest{UserId: userId})
	if err != nil {
		fmt.Printf("ERROR calling Cart: %v\n", err)
		return nil, err
	}
	if len(cartRes.Items) == 0 {
		fmt.Println("ERROR: Cart is empty")
		return nil, errors.New("cart is empty")
	}
	fmt.Printf("Cart OK. Items: %d, RestID: %s\n", len(cartRes.Items), cartRes.RestaurantId)

	fmt.Println("--- 2. Calling Restaurant Service ---")
	restRes, err := s.restClient.GetRestaurant(ctx, &pbRest.GetRestaurantRequest{Id: cartRes.RestaurantId})
	if err != nil {
		fmt.Printf("ERROR calling Restaurant: %v\n", err)
		return nil, errors.New("restaurant not found")
	}
	fmt.Printf("Restaurant OK. Lat: %f, Lon: %f\n", restRes.Latitude, restRes.Longitude)

	// 3. Determine Location
	finalLat := cLat
	finalLon := cLon
	finalAddr := ""

	if customAddr != nil && *customAddr != "" {
		finalAddr = *customAddr
		fmt.Println("Using CUSTOM Address")
	} else {
		fmt.Println("--- 3. Calling Stakeholders (Profile) ---")
		profile, err := s.stakeClient.GetCustomer(ctx, &pbStake.GetRequest{UserId: userId})
		if err != nil {
			fmt.Printf("ERROR calling Stakeholders: %v\n", err)
			return nil, errors.New("profile not found")
		}
		finalAddr = profile.Address
		finalLat = profile.Latitude
		finalLon = profile.Longitude
		fmt.Println("Profile OK.")
	}

	fmt.Println("--- 4. Calling Pricing Service ---")
	priceRes, err := s.pricingClient.CalculatePrice(ctx, &pbPricing.CalculatePriceRequest{
		CustomerLat:   finalLat,
		CustomerLon:   finalLon,
		RestaurantLat: restRes.Latitude,
		RestaurantLon: restRes.Longitude,
		IsPriority:    isPriority,
	})
	if err != nil {
		fmt.Printf("ERROR calling Pricing: %v\n", err)
		return nil, errors.New("pricing failed")
	}
	fmt.Printf("Pricing OK. Total: %f\n", priceRes.FinalPrice)

	fmt.Println("--- 5. Calling Analysis Service ---")
	feesRes, err := s.analysisClient.CalculateFees(ctx, &pbAnalysis.CalculationRequest{
		ItemsTotal:    cartRes.TotalPrice,
		DeliveryPrice: priceRes.FinalPrice,
	})
	if err != nil {
		fmt.Printf("ERROR calling Analysis: %v\n", err)
		return nil, errors.New("analysis failed")
	}
	fmt.Println("Analysis OK.")

	// 6. Assemble Order
	totalPrice := cartRes.TotalPrice + priceRes.FinalPrice + feesRes.AppFee + feesRes.SmallOrderFee

	order := &models.Order{
		CustomerId:         uuid.MustParse(userId),
		RestaurantId:       uuid.MustParse(cartRes.RestaurantId),
		Status:             models.StatusPending,
		DeliveryAddress:    finalAddr,
		DeliveryLat:        finalLat,
		DeliveryLon:        finalLon,
		ItemsTotal:         cartRes.TotalPrice,
		DeliveryFee:        priceRes.FinalPrice,
		AppFee:             feesRes.AppFee,
		SmallOrderFee:      feesRes.SmallOrderFee,
		TotalPrice:         totalPrice,
		IsPriority:         isPriority,
		ProfitFromItems:    0,
		ProfitFromDelivery: 0,
		TotalProfit:        feesRes.EstimatedProfit,
		CreatedAt:          time.Now(),
	}

	for _, i := range cartRes.Items {
		order.Items = append(order.Items, models.OrderItem{
			Name:     i.Name,
			Price:    i.Price,
			Quantity: int(i.Quantity),
		})
	}

	fmt.Println("--- 6. Saving to DB ---")
	if err := s.repo.Create(order); err != nil {
		fmt.Printf("ERROR saving to DB: %v\n", err)
		return nil, err
	}
	fmt.Println("DB Save OK.")

	// ... Record Profit & Clear Cart ... (You can log these too, but they usually don't block the response)

	return order, nil
}

func (s *OrderService) GetMyOrders(userId string) ([]models.Order, error) {
	return s.repo.GetByCustomer(userId)
}

func (s *OrderService) GetOrdersByRestaurant(restaurantId string, status models.OrderStatus) ([]models.Order, error) {
	return s.repo.GetByRestaurant(restaurantId, status)
}

// 3. UPDATE STATUS
func (s *OrderService) UpdateStatus(orderId string, status models.OrderStatus, deliveryPersonId *string, prepMinutes int32) (*models.Order, error) {
	err := s.repo.UpdateStatus(orderId, status, deliveryPersonId, int(prepMinutes))
	if err != nil {
		return nil, err
	}

	// 2. Fetch and return the updated order so the frontend sees the change immediately
	updatedOrder, err := s.repo.GetById(orderId)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (s *OrderService) PlaceBid(ctx context.Context, orderId string, driverId string, minutes int) error {
	fmt.Println(">>> [SERVICE] Starting PlaceBid...")

	// Check if driver is busy
	existingJob, err := s.repo.GetActiveJobsForDriver(driverId)

	// Stop if there is a DB error (that isn't "Record Not Found")
	if err != nil {
		fmt.Println(">>> [SERVICE] DB Error checking active job")
		return err
	}

	// Now this check works correctly because repo returns nil
	if existingJob != nil {
		fmt.Println(">>> [SERVICE] Driver is busy!")
		return errors.New("driver already has an active job")
	}

	fmt.Println(">>> [SERVICE] Driver is free. Getting Order...")

	// 1. GET ORDER FRESH STATE
	order, err := s.repo.GetById(orderId)
	if err != nil {
		return err
	}

	// 2. CHECK IF BIDDING IS EXPIRED
	if order.BiddingExpiresAt != nil {
		if time.Now().After(*order.BiddingExpiresAt) {
			return errors.New("bidding window has closed")
		}
	}

	// 4. Check Priority Constraint
	if order.IsPriority {
		deliveryPersonProto, err := s.stakeClient.GetDeliveryPerson(ctx, &pbStake.GetRequest{UserId: driverId})
		if err != nil {
			fmt.Println(">>> [SERVICE] Failed to get Driver Profile")
			return errors.New("delivery person not found")
		}

		fmt.Printf(">>> [SERVICE] Driver Vehicle: %v\n", deliveryPersonProto.Vehicle)
		if deliveryPersonProto.Vehicle != pbStake.VehicleType_CAR {
			fmt.Println(">>> [SERVICE] Priority Blocked: Not a car")
			return errors.New("priority orders require a car")
		}
	}

	// 3. IF FIRST BID -> START TIMER
	if order.BiddingExpiresAt == nil {
		fmt.Printf(">>> [BID] First bid received for %s. Starting 30s timer.\n", orderId)

		// Calculate expiration
		expiresAt := time.Now().Add(45 * time.Second)

		// Save Expiration to DB (So it persists if server crashes)
		err = s.repo.SetBiddingExpiration(orderId, expiresAt)
		if err != nil {
			return err
		}

		// Fire the background goroutine
		go func() {
			time.Sleep(5 * time.Second)
			// Trigger the winner logic
			// We can reuse the background worker logic or call a helper here
			// But relying on the Background Worker (Step 4) is safer!
		}()
	}

	// 4. SAVE BID
	bid := &models.OrderBid{
		OrderId:  uuid.MustParse(orderId),
		DriverId: uuid.MustParse(driverId),
		Minutes:  minutes,
	}
	return s.repo.CreateBid(bid)
}

func (s *OrderService) StartBidProcessingLoop() {
	fmt.Println(">>> [WORKER] Background Bid Processor Started...")

	go func() {
		for {
			// 1. Run every 10 seconds
			time.Sleep(10 * time.Second)

			// 2. Find orders stuck in "PREPARING" without a driver
			orders, err := s.repo.GetOrdersReadyForAssignment()
			if err != nil {
				fmt.Printf(">>> [WORKER] Error checking orders: %v\n", err)
				continue
			}

			// 3. Process each one
			for _, o := range orders {
				fmt.Printf(">>> [WORKER] Picking winner for expired order %s\n", o.Id)

				winner, err := s.repo.PickWinner(o.Id.String())
				if err != nil {
					// No bids yet? Ignore and try again next loop.
					continue
				}

				// Assign
				err = s.repo.AssignDriver(o.Id.String(), winner.DriverId.String(), winner.Minutes)
				if err == nil {
					fmt.Printf(">>> [WORKER] WINNER ASSIGNED: Driver %s (%d minutes) for Order %s\n", winner.DriverId, winner.Minutes, o.Id)
				}
			}
		}
	}()
}
func (s *OrderService) GetAvailableOrders() ([]models.Order, error) {
	return s.repo.GetAvailableOrders()
}

func (s *OrderService) GetActiveJob(driverId string) (*models.Order, error) {
	return s.repo.GetActiveJobsForDriver(driverId)
}
