# Use the official golang image as the base image
FROM golang:1.22.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod ./

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary with race detection enabled
RUN CGO_ENABLED=0 GOOS=linux go build -o app -a -ldflags '-extldflags "-static"' .

# Use a minimal base image to run the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage to the runtime image
COPY --from=builder /app/app .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app"]