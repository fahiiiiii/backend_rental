# Use the official Golang image as a build stage
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest  

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]












# # Use official Golang image as base
# FROM golang:1.21-alpine

# # Set working directory
# WORKDIR /app

# # Install required system dependencies
# RUN apk update && apk add --no-cache \
#     git \
#     gcc \
#     musl-dev \
#     bash

# # Install Beego and Bee tool
# RUN go install github.com/beego/bee/v2@latest && \
#     go install github.com/beego/beego/v2@latest

# # Copy go mod files first for better caching
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the rest of the application
# COPY . .

# # Build the application
# RUN go build -o main .

# # Expose the port the app runs on
# EXPOSE 8080

# # Command to run the application using bee (for hot reload during development)
# CMD ["bee", "run"]