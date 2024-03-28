# Use the official Golang image with version 1.22.0 as the base image
FROM golang:1.22.1-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the entire project directory into the container
COPY . .

# Build the Go application with proper LDFLAGS for security
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o main ./cmd/SafeTransfer

# Use a minimal distroless image for the final stage
FROM gcr.io/distroless/static-debian11

# Create a non-root user
USER nonroot:nonroot

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port on which your application listens
EXPOSE 8083

# Start the application
CMD ["./main"]
