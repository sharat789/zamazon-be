# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy the metrics package first (shared dependency)
COPY metrics ./metrics/

# Copy the auth service module files
COPY auth/go.mod auth/go.sum ./auth/

# Copy the auth service source code
COPY auth/ ./auth/

# Change working directory to auth for building
WORKDIR /app/auth

# Download dependencies
RUN go mod download

# Build the binary
RUN go build -v -o auth-service ./cmd

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/auth/auth-service .

# Expose the port
EXPOSE 8080

CMD ["./auth-service"]