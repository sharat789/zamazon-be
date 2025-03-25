# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy the transactions service go.mod and go.sum from the repository root’s transactions directory
COPY transactions/go.mod transactions/go.sum ./transactions/

# Copy the common package from the repository root (so that when in /app/transactions, ../common exists)
COPY common ./common

# Copy the rest of the transactions service source code
COPY transactions/ ./transactions/

# Change working directory to transactions so that module replacement ../common is valid
WORKDIR /app/transactions

#RUN sed -i 's/^go 1\.23\.4/go 1.23/' go.mod
# Download dependencies (this will now find ../common)
RUN go mod download

# Build the binary
RUN go build -v -o transactions-service .

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/transactions/transactions-service .

# Expose the port (adjust as needed)
EXPOSE 8080

CMD ["./transactions-service"]