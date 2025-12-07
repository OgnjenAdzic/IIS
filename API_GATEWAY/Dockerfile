# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Copy COMMON first
COPY COMMON ./COMMON

# Copy Gateway mod files
COPY API_GATEWAY/go.mod API_GATEWAY/go.sum ./API_GATEWAY/

WORKDIR /app/API_GATEWAY
RUN go mod download

# Copy Gateway source
COPY API_GATEWAY/ .

# Build
RUN go build -o gateway cmd/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/API_GATEWAY/gateway .

EXPOSE 8080
CMD ["./gateway"]