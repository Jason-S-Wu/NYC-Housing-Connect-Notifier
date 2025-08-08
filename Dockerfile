# Build stage
FROM golang:1.24.6-alpine3.22 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Final stage
FROM alpine:3.22

WORKDIR /app

# Install necessary packages
RUN apk update && apk add bash && apk --no-cache add tzdata

# Upgrade packages in a single layer
RUN apk upgrade --no-cache

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Run the binary
CMD ["./main"]