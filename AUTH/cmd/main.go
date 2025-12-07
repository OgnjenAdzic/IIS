package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"AUTH/internal/database"
	handler "AUTH/internal/handlers"
	"AUTH/internal/repository"
	"AUTH/internal/service"

	pb "COMMON/auth/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// --- 1. CONFIGURATION ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Default port for Auth Service
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my_super_secret_dev_key" // Fallback for development
	}

	// --- 2. DATABASE ---
	db := database.Connect()

	// --- 3. DEPENDENCY INJECTION ---
	// Repo -> Service -> Handler

	// Create Repository
	userRepo := repository.NewUserRepository(db)

	// Create Services
	userService := service.NewUserService(userRepo)
	jwtService := service.NewJWTService(jwtSecret, "delivery-app-auth")

	// Create Handler (Controller)
	authHandler := handler.NewAuthHandler(userService, jwtService)

	// --- 4. NETWORK LISTENER ---
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// --- 5. GRPC SERVER SETUP ---
	grpcServer := grpc.NewServer()

	// Register our implementation with the gRPC server
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	// Enable Reflection
	// This allows tools like Postman or grpcurl to inspect the schema automatically
	reflection.Register(grpcServer)

	// --- 6. START SERVER ---
	fmt.Printf("Auth gRPC Service running on port %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
