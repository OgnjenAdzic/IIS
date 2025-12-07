# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

# Postavljamo root projekta unutar kontejnera
WORKDIR /app

# 1. Kopiramo COMMON folder
COPY COMMON ./COMMON

# 2. Kopiramo go.mod i go.sum od AUTH servisa u njegov podfolder
COPY AUTH/go.mod AUTH/go.sum ./AUTH/

# 3. Ulazimo u AUTH folder pre downloada
WORKDIR /app/AUTH

# Downloadujemo zavisnosti (sada radi jer replace gadja ../COMMON)
RUN go mod download

# 4. Kopiramo ostatak koda iz AUTH foldera
COPY AUTH/ .

# Build
RUN go build -o auth cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/AUTH/auth .

EXPOSE 8082
CMD ["./auth"]