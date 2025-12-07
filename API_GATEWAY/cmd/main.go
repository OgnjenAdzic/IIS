package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"API_GATEWAY/middleware"

	// Import definitions from COMMON
	pbAuth "COMMON/auth/proto"

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
