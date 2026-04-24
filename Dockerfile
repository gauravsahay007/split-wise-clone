# -------- Build Stage --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for some dependencies)
RUN apk add --no-cache git

# Copy go mod files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of code
COPY . .

# Build the Go binary
RUN go build -o main .

# -------- Run Stage --------
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port (change if needed)
EXPOSE 8080

# Run the app
CMD ["./main"]