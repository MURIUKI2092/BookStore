# Stage 1: Build the Go binary
FROM golang:1.20-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o bookstore .

# Stage 2: Run the Go binary
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bookstore .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./bookstore"]
