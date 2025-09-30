FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make build-base

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the platform binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o platform ./cmd/platform

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder stage
COPY --from=builder /app/platform .

# Copy any required configuration files
COPY --from=builder /app/configs ./configs

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Use non-root user
USER appuser

# Expose API port
EXPOSE 8081

# Run the platform binary
CMD ["./platform"]
