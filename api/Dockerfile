# Build stage
FROM golang:1.24.2-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./

# Final stage
FROM alpine:latest

# Install ca-certificates and curl for healthcheck
RUN apk --no-cache add ca-certificates curl

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 