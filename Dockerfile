# Use the official Golang image
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Install necessary libraries (if required)
RUN apk add --no-cache bash

# Copy the binary from the builder stage
COPY --from=builder /app/main /main

# Expose the port
EXPOSE 8080

# Start the application
CMD ["./main"]