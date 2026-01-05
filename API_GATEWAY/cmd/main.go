package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"API_GATEWAY/middleware"

	// Import definitions from COMMON
	pbAnalysis "COMMON/analysis/proto"
	pbAuth "COMMON/auth/proto"
	pbPricing "COMMON/pricing/proto"
	pbRestaurant "COMMON/restaurant/proto"
	pbCart "COMMON/shopping_cart/proto"
	pbStakeholders "COMMON/stakeholders/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. Create gRPC Gateway Mux
	// This translates HTTP JSON requests into gRPC calls
	gwmux := runtime.NewServeMux(
		// Optional: Custom header matcher to pass specific headers to gRPC metadata
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "Grpc-Metadata-User-Id":
				return "user-id", true
			case "Grpc-Metadata-User-Role":
				return "user-role", true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// --- AUTH SERVICE ---
	authAddr := os.Getenv("AUTH_SERVICE_ADDRESS")
	if authAddr == "" {
		authAddr = "auth-service:8082" // Default Docker name
	}

	err := pbAuth.RegisterAuthServiceHandlerFromEndpoint(ctx, gwmux, authAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register Auth Service: %v", err)
	}
	fmt.Println("Auth Service Registered at " + authAddr)

	// --- RESTAURANT SERVICE ---
	restAddr := os.Getenv("RESTAURANT_SERVICE_ADDRESS")
	if restAddr == "" {
		restAddr = "restaurant-service:8083"
	}
	err = pbRestaurant.RegisterRestaurantServiceHandlerFromEndpoint(ctx, gwmux, restAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register Restaurant: %v", err)
	}
	fmt.Println("Restaurant Service Registered at " + restAddr)

	// --- STAKEHOLDERS SERVICE ---
	stakeAddr := os.Getenv("STAKEHOLDERS_SERVICE_ADDRESS")
	if stakeAddr == "" {
		stakeAddr = "stakeholder-service:8084"
	}
	err = pbStakeholders.RegisterStakeholdersServiceHandlerFromEndpoint(ctx, gwmux, stakeAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register stakeholder: %v", err)
	}
	fmt.Println("stakeholder Service Registered at " + stakeAddr)

	// --- CART SERVICE ---
	cartAddr := os.Getenv("CART_SERVICE_ADDRESS")
	if cartAddr == "" {
		cartAddr = "shopping-cart-service:8085"
	}
	err = pbCart.RegisterShoppingCartServiceHandlerFromEndpoint(ctx, gwmux, cartAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register shopping cart: %v", err)
	}
	fmt.Println("Shopping Cart Service Registered at " + cartAddr)

	// --- PRICING SERVICE ---
	pricingAddr := os.Getenv("PRICING_SERVICE_ADDRESS")
	if pricingAddr == "" {
		pricingAddr = "pricing-service:8086"
	}
	err = pbPricing.RegisterPricingServiceHandlerFromEndpoint(ctx, gwmux, pricingAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register pricing service: %v", err)
	}
	fmt.Println("Pricing Service Registered at " + pricingAddr)

	// --- ANALYSIS SERVICE ---
	analysisAddr := os.Getenv("ANALYSIS_SERVICE_ADDRESS")
	if analysisAddr == "" {
		analysisAddr = "analysis-service:8087"
	}
	err = pbAnalysis.RegisterAnalysisServiceHandlerFromEndpoint(ctx, gwmux, analysisAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register analysis service: %v", err)
	}
	fmt.Println("Analysis Service Registered at " + analysisAddr)

	// 2. Start HTTP server (and proxy calls to gRPC server endpoint)
	rootMux := http.NewServeMux()

	rootMux.Handle("/", gwmux)

	authHandler := middleware.AuthMiddleware(rootMux)
	finalHandler := middleware.CorsMiddleware(authHandler)

	port := "8080"
	fmt.Printf("API Gateway running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		log.Fatal(err)
	}
}
