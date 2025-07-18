# ---------------------------
# Stage 1: Build the Go binary
# ---------------------------
FROM golang:1.24.3-alpine AS builder

# Install git for go mod (in caso servano repo remote)
RUN apk add --no-cache git

WORKDIR /app

# Copia i moduli e scarica le dipendenze
COPY go.mod go.sum ./
RUN go mod download

# Copia tutto il codice
COPY . .

# Compila il binario
RUN go build -o api-server fejs.go
# Compiliamo il binario
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./fejs.go

# ---------------------------
# Stage 2: Run the binary
# ---------------------------
FROM alpine:latest

WORKDIR /root/

# Copia il binario compilato dal builder
COPY --from=builder /app/api-server .

# Esponi la porta (modifica se usi una porta diversa)
EXPOSE 8080

# Comando di avvio
CMD ["./api-server"]
