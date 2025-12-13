# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

# Postavljamo root projekta unutar kontejnera
WORKDIR /app

# 1. Kopiramo COMMON folder
COPY COMMON ./COMMON

# 2. Kopiramo go.mod i go.sum od STAKEHOLDERS servisa u njegov podfolder
COPY STAKEHOLDERS/go.mod STAKEHOLDERS/go.sum ./STAKEHOLDERS/

# 3. Ulazimo u STAKEHOLDERS folder pre downloada
WORKDIR /app/STAKEHOLDERS

# Downloadujemo zavisnosti (sada radi jer replace gadja ../COMMON)
RUN go mod download

# 4. Kopiramo ostatak koda iz STAKEHOLDERS foldera
COPY STAKEHOLDERS/ .

# Build
RUN go build -o stakeholders cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/STAKEHOLDERS/stakeholders .

EXPOSE 8084
CMD ["./stakeholders"]