# Use the official Golang image with version 1.22.0 as the base image
FROM golang:1.22.0-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the entire project directory into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/SafeTransfer

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Install bash (required for wait-for-it.sh if using bash specific syntax)
RUN apk add --no-cache bash

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the configuration file
COPY configs/config.yaml ./configs/

# Copy the wait-for-it.sh script
COPY wait-for-it.sh /wait-for-it.sh

# After copying your application files
COPY .env ./

# Make the wait-for-it.sh script executable
RUN chmod +x /wait-for-it.sh

# Set the environment variable for JWT secret key
ENV JWT_SECRET_KEY=lojdXnvv3bgUmDSDK+z7i2u7TtlWmu5S0nAH0Ki3mHk=

# Expose the port on which your application listens
EXPOSE 8083

# Start the application using wait-for-it.sh to wait for the db service
CMD ["/wait-for-it.sh", "db:5432", "--", "./main"]
