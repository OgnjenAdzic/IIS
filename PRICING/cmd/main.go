package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"pricing/internal/database"
	handler "pricing/internal/handlers"
	"pricing/internal/repository"
	"pricing/internal/service"

	pb "COMMON/pricing/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// --- 1. CONFIGURATION ---
	// Default to 8086 (or whatever you set in docker-compose for Pricing)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	// --- 2. DATABASE ---
	db := database.Connect()

	// --- 3. DEPENDENCY INJECTION ---

	// Create Repository (This triggers the default rule creation logic)
	pricingRepo := repository.NewPricingRepository(db)

	// Create Service
	pricingService := service.NewPricingService(pricingRepo)

	// Create Handler
	pricingHandler := handler.NewPricingHandler(pricingService)

	// --- 4. NETWORK LISTENER ---
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// --- 5. GRPC SERVER SETUP ---
	grpcServer := grpc.NewServer()

	// Register the Pricing implementation
	pb.RegisterPricingServiceServer(grpcServer, pricingHandler)

	// Enable Reflection (useful for debugging tools)
	reflection.Register(grpcServer)

	// --- 6. START SERVER ---
	fmt.Printf("Pricing gRPC Service running on port %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
