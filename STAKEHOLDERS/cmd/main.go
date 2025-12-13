package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"stakeholders/internal/database"
	handler "stakeholders/internal/handlers"
	"stakeholders/internal/repository"
	"stakeholders/internal/service"

	pb "COMMON/stakeholders/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// --- 1. CONFIGURATION ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	// --- 2. DATABASE ---
	db := database.Connect()

	// --- 3. DEPENDENCY INJECTION ---
	// Repo -> Service -> Handler

	// Create Repository
	userRepo := repository.NewStakeholdersRepository(db)

	// Create Services
	userService := service.NewStakeholdersService(userRepo)

	// Create Handler (Controller)
	stakeholdersService := handler.NewStakeholdersHandler(userService)

	// --- 4. NETWORK LISTENER ---
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// --- 5. GRPC SERVER SETUP ---
	grpcServer := grpc.NewServer()

	// Register our implementation with the gRPC server
	pb.RegisterStakeholdersServiceServer(grpcServer, stakeholdersService)

	// Enable Reflection
	// This allows tools like Postman or grpcurl to inspect the schema automatically
	reflection.Register(grpcServer)

	// --- 6. START SERVER ---
	fmt.Printf("stakeholders gRPC Service running on port %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
