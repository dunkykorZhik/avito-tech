# 1. Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/main.go

# 2. Run stage (smaller image)
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main .

# copy configs (add this!)
COPY internal/app/config ./internal/app/config
# Expose your app port (adjust if different)
EXPOSE 8080

CMD ["./main"]
