package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"restaurant/internal/database"
	handler "restaurant/internal/handlers"
	"restaurant/internal/repository"
	"restaurant/internal/service"

	pb "COMMON/restaurant/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// --- 1. CONFIGURATION ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // Default port for Auth Service
	}

	// --- 2. DATABASE ---
	db := database.Connect()

	// --- 3. DEPENDENCY INJECTION ---
	// Repo -> Service -> Handler

	// Create Repository
	userRepo := repository.NewRestaurantRepository(db)

	// Create Services
	userService := service.NewRestaurantService(userRepo)

	// Create Handler (Controller)
	restaurantService := handler.NewRestaurantHandler(userService)

	// --- 4. NETWORK LISTENER ---
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// --- 5. GRPC SERVER SETUP ---
	grpcServer := grpc.NewServer()

	// Register our implementation with the gRPC server
	pb.RegisterRestaurantServiceServer(grpcServer, restaurantService)

	// Enable Reflection
	// This allows tools like Postman or grpcurl to inspect the schema automatically
	reflection.Register(grpcServer)

	// --- 6. START SERVER ---
	fmt.Printf("Restaurant gRPC Service running on port %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
