# Start from the official Golang image for building
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY ./src .

# Build the Go app
RUN go build -o caddy-gitops .

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from builder
COPY --from=builder /app/caddy-gitops .

# Expose port (change as needed)
EXPOSE 2020

# Run the binary
CMD ["./caddy-gitops"]