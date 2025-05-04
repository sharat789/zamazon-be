# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy the users service go.mod and go.sum from the repository root’s users directory
COPY users/go.mod users/go.sum ./users/

# Copy the metrics package from the repository root (so that when in /app/users, ../metrics exists)
COPY metrics ./metrics

# Copy the rest of the users service source code
COPY users/ ./users/

# Change working directory to users so that module replacement ../common is valid
WORKDIR /app/users

#RUN sed -i 's/^go 1\.23\.4/go 1.23/' go.mod
# Download dependencies (this will now find ../common)
RUN go mod download

# Build the binary
RUN go build -v -o users-service .

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/users/users-service .

# Expose the port (adjust as needed)
EXPOSE 8080

CMD ["./users-service"]