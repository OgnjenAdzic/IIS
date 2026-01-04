# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

# Postavljamo root projekta unutar kontejnera
WORKDIR /app

# 1. Kopiramo COMMON folder
COPY COMMON ./COMMON

# 2. Kopiramo go.mod i go.sum od Restaurant servisa u njegov podfolder
COPY PRICING/go.mod PRICING/go.sum ./PRICING/

# 3. Ulazimo u Restaurant folder pre downloada
WORKDIR /app/PRICING

# Downloadujemo zavisnosti (sada radi jer replace gadja ../COMMON)
RUN go mod download

# 4. Kopiramo ostatak koda iz Restaurant foldera
COPY PRICING/ .

# Build
RUN go build -o pricing cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/PRICING/pricing .

EXPOSE 8086
CMD ["./pricing"]