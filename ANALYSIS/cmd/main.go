package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"analysis/internal/database"
	handler "analysis/internal/handlers"
	"analysis/internal/repository"
	"analysis/internal/service"

	pb "COMMON/analysis/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// --- 1. CONFIGURATION ---
	// Default to 8086 (or whatever you set in docker-compose for Pricing)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	// --- 2. DATABASE ---
	db := database.Connect()

	// --- 3. DEPENDENCY INJECTION ---

	// Create Repository (This triggers the default rule creation logic)
	analysisRepo := repository.NewAnalysisRepository(db)

	// Create Service
	analysisService := service.NewAnalysisService(analysisRepo)

	// Create Handler
	analysisHandler := handler.NewAnalysisHandler(analysisService)

	// --- 4. NETWORK LISTENER ---
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// --- 5. GRPC SERVER SETUP ---
	grpcServer := grpc.NewServer()

	// Register the Analysis implementation
	pb.RegisterAnalysisServiceServer(grpcServer, analysisHandler)

	// Enable Reflection (useful for debugging tools)
	reflection.Register(grpcServer)

	// --- 6. START SERVER ---
	fmt.Printf("Pricing gRPC Service running on port %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
