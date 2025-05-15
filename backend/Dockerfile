# Backend Dockerfile for Canny Clone
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final lightweight image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy config files
COPY config.*.json ./

# Expose the port the application runs on
EXPOSE 8080

# Set environment variable for configuration
ENV APP_ENV=develop

# Run the binary
CMD ["./main"]
