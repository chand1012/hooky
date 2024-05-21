# Stage 1: Build the Go application
FROM golang:1.22-alpine as builder

# Set the working directory to /app
WORKDIR /app

# Copy the Go source code
COPY . .

# Install dependencies
RUN go get -u -v ./...

# Build the Go application
RUN CGO_ENABLED=0 go build -v -o hooky .

# Stage 2: Create a minimal Docker image for the Hooky app
FROM alpine:latest

# Set the working directory to /app
WORKDIR /app

# Copy the Hooky binary from the builder stage
COPY --from=builder /app/hooky .

ENTRYPOINT [ "/app/hooky" ]
