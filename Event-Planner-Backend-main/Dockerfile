# Stage 1: Build Go backend
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

# Install build dependencies (required for CGO with Alpine)
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Stage 2: Run Go backend
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache ca-certificates wget

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Set execute permissions
RUN chmod +x /root/main

# Expose port
EXPOSE 8080

# Health check - using /api/ping endpoint
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/api/ping || exit 1

# Run the application
CMD ["./main"]
