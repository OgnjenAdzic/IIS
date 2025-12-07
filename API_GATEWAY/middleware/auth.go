package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. SKIP PUBLIC ROUTES
		// These endpoints don't require a token
		path := r.URL.Path
		if strings.Contains(path, "/auth/login") ||
			strings.Contains(path, "/auth/register") ||
			strings.Contains(path, "/swagger") {
			next.ServeHTTP(w, r)
			return
		}

		// 2. CHECK HEADER
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		// 3. VALIDATE TOKEN
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Use the same secret as Auth Service
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "my_super_secure_delivery_app_secret_key_2025"
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. EXTRACT CLAIMS
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		// DEBUG: Print what is inside the token
		fmt.Printf("Token Claims: %+v\n", claims)

		// 5. INJECT INTO HEADERS FOR GRPC
		// Note: "sub" is the standard claim for User ID we used in Auth Service
		if userID, ok := claims["sub"].(string); ok {
			r.Header.Set("Grpc-Metadata-User-Id", userID)
		}
		if role, ok := claims["role"].(string); ok {
			r.Header.Set("Grpc-Metadata-User-Role", role)
		}
		if username, ok := claims["username"].(string); ok {
			r.Header.Set("Grpc-Metadata-User-Username", username)
		}

		next.ServeHTTP(w, r)
	})
}
