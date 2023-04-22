# Use the official Golang image as a base image
FROM golang:1.17-alpine as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files into the workspace
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source files into the workspace
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order_service ./cmd/main.go

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/order_service .

# Install the ca-certificates package for HTTPS support
RUN apk --no-cache add ca-certificates

# Expose the port the service will run on
EXPOSE 8080

# Run the service
CMD ["./order_service"]
