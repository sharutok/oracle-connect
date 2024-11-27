# Use the official Golang image as a builder
FROM golang:1.22 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp

# Use a minimal image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

# Expose the application port
EXPOSE 7777

# Command to run the executable
CMD ["./myapp"]
