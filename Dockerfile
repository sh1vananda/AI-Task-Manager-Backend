# Stage 1: Build the application
FROM golang:1.21 AS builder

# Install necessary tools and libraries for SQLite
RUN apt-get update && apt-get install -y \
    gcc \
    make \
    sqlite3 \
    libsqlite3-dev

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Debug: List files before building
RUN ls -la

# Build the application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

# Debug: Verify the binary was created
RUN ls -la

# Ensure the binary has execution permissions
RUN chmod +x main

# Stage 2: Create the final image
FROM debian:buster-slim

# Install SQLite runtime dependencies
RUN apt-get update && apt-get install -y sqlite3

# Copy the binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/main

# Debug: Verify the binary exists in the final image
RUN ls -la /usr/local/bin

# Expose the port
EXPOSE 8080

# Start the application
CMD ["main"]