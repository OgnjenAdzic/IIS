# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

# Postavljamo root projekta unutar kontejnera
WORKDIR /app

# 1. Kopiramo COMMON folder
COPY COMMON ./COMMON

# 2. Kopiramo go.mod i go.sum od Restaurant servisa u njegov podfolder
COPY SHOPPING_CART/go.mod SHOPPING_CART/go.sum ./SHOPPING_CART/

# 3. Ulazimo u Restaurant folder pre downloada
WORKDIR /app/SHOPPING_CART

# Downloadujemo zavisnosti (sada radi jer replace gadja ../COMMON)
RUN go mod download

# 4. Kopiramo ostatak koda iz Restaurant foldera
COPY SHOPPING_CART/ .

# Build
RUN go build -o shopping_cart cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/SHOPPING_CART/shopping_cart .

EXPOSE 8085
CMD ["./shopping_cart"]