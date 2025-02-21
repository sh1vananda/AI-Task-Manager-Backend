# Stage 1: Build the application
FROM golang:1.21 AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    make \
    libsqlite3-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build with static linking
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags sqlite3 -ldflags="-extldflags=-static" -o main .

# Stage 2: Create final image
FROM gcr.io/distroless/base-debian12

COPY --from=builder /app/main /
COPY *.sql /schema/

EXPOSE 8080
CMD ["/main"]