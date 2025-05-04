# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy the catalog service go.mod and go.sum from the repository root’s catalog directory
COPY catalog/go.mod catalog/go.sum ./catalog/

# Copy the metrics package from the repository root (so that when in /app/users, ../metrics exists)
COPY metrics ./metrics

# Copy the rest of the catalog service source code
COPY catalog/ ./catalog/

# Change working directory to catalog so that module replacement ../common is valid
WORKDIR /app/catalog

#RUN sed -i 's/^go 1\.23\.4/go 1.23/' go.mod
# Download dependencies (this will now find ../common)
RUN go mod download

# Build the binary
RUN go build -v -o catalog-service .

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/catalog/catalog-service .

# Expose the port (adjust as needed)
EXPOSE 8080

CMD ["./catalog-service"]