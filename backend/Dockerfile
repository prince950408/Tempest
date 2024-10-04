# Stage 1: Build the Go binary
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN go build -o main .

# Stage 2: Use a smaller image to run the Go binary
FROM alpine:3.18

# Set the working directory in the final container
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port on which the Go app will run
EXPOSE 8080

# Run the Go binary
CMD ["./main"]
