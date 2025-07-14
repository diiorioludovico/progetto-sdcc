# ---------------------------
# Stage 1: Build the Go binary
# ---------------------------
FROM golang:1.24.3-alpine AS builder

# Install git for go mod (in caso servano repo remote)
RUN apk add --no-cache git

WORKDIR /app

# Copiamo go.mod e go.sum per scaricare solo le dipendenze
COPY go.mod go.sum ./
RUN go mod download

# Copiamo tutto il codice
COPY . .

# Compiliamo il binario
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./server.go

# ---------------------------
# Stage 2: Run the binary
# ---------------------------
FROM alpine:latest

WORKDIR /root/

# Copiamo il binario dal builder
COPY --from=builder /app/app .

# Espone le porte gRPC e HTTP
EXPOSE 50051
EXPOSE 8080

# Comando di avvio
CMD ["./app"]
