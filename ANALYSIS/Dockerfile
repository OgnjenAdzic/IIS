# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

# Postavljamo root projekta unutar kontejnera
WORKDIR /app

# 1. Kopiramo COMMON folder
COPY COMMON ./COMMON

# 2. Kopiramo go.mod i go.sum od Restaurant servisa u njegov podfolder
COPY ANALYSIS/go.mod ANALYSIS/go.sum ./ANALYSIS/

# 3. Ulazimo u Restaurant folder pre downloada
WORKDIR /app/ANALYSIS

# Downloadujemo zavisnosti (sada radi jer replace gadja ../COMMON)
RUN go mod download

# 4. Kopiramo ostatak koda iz Restaurant foldera
COPY ANALYSIS/ .

# Build
RUN go build -o analysis cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ANALYSIS/analysis .

EXPOSE 8087
CMD ["./analysis"]