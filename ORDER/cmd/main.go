package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"order/internal/database"
	handler "order/internal/handlers"
	"order/internal/repository"
	"order/internal/service"

	pbOrder "COMMON/order/proto"
	// Import other protos
	pbAnalysis "COMMON/analysis/proto"
	pbPricing "COMMON/pricing/proto"
	pbRest "COMMON/restaurant/proto"
	pbCart "COMMON/shopping_cart/proto"
	pbStake "COMMON/stakeholders/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. CONFIG
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	// 2. DB
	db := database.Connect()

	// 3. gRPC CLIENTS (Connecting to other services)
	// We use "insecure" because we are inside the internal Docker network
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Shopping Cart
	connCart, err := grpc.NewClient("shopping-cart-service:8085", opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Cart: %v", err)
	}
	cartClient := pbCart.NewShoppingCartServiceClient(connCart)

	// Pricing
	connPricing, err := grpc.NewClient("pricing-service:8086", opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Pricing: %v", err)
	}
	pricingClient := pbPricing.NewPricingServiceClient(connPricing)

	// Analysis
	connAnalysis, err := grpc.NewClient("analysis-service:8087", opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Analysis: %v", err)
	}
	analysisClient := pbAnalysis.NewAnalysisServiceClient(connAnalysis)

	// Restaurant
	connRest, err := grpc.NewClient("restaurant-service:8083", opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Restaurant: %v", err)
	}
	restClient := pbRest.NewRestaurantServiceClient(connRest)

	// Stakeholders
	connStake, err := grpc.NewClient("stakeholders-service:8084", opts...) // Note: Check port (50051 or 8084 depending on your compose)
	if err != nil {
		log.Fatalf("Failed to connect to Stakeholders: %v", err)
	}
	stakeClient := pbStake.NewStakeholdersServiceClient(connStake)

	// 4. WIRING
	orderRepo := repository.NewOrderRepository(db)

	// Inject all clients into the service
	orderService := service.NewOrderService(
		orderRepo,
		cartClient,
		pricingClient,
		analysisClient,
		restClient,
		stakeClient,
	)

	orderHandler := handler.NewOrderHandler(orderService)

	// 5. SERVER
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pbOrder.RegisterOrderServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	fmt.Printf("Order Service running on port %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
