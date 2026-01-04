package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"shopping_cart/internal/database"
	handler "shopping_cart/internal/handlers"
	"shopping_cart/internal/repository"
	"shopping_cart/internal/service"

	pb "COMMON/shopping_cart/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// --- 1. CONFIGURATION ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085" // Default port for Shopping Cart Service
	}

	// --- 2. DATABASE ---
	db := database.Connect()

	// --- 3. DEPENDENCY INJECTION ---

	// Create Repository
	cartRepo := repository.NewCartRepository(db)

	// Create Services
	cartService := service.NewCartService(cartRepo)

	// Create Handler (Controller)
	cartHandler := handler.NewCartHandler(cartService)

	// --- 4. NETWORK LISTENER ---
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// --- 5. GRPC SERVER SETUP ---
	grpcServer := grpc.NewServer()

	// Register our implementation with the gRPC server
	pb.RegisterShoppingCartServiceServer(grpcServer, cartHandler)

	// Enable Reflection
	reflection.Register(grpcServer)

	// --- 6. START SERVER ---
	fmt.Printf("Shopping Cart gRPC Service running on port %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
