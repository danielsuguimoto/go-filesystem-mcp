# Build stage
FROM golang:1.22.1-alpine3.19 AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-filesystem-mcp .

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/go-filesystem-mcp .

# Create a directory for mounting volumes
RUN mkdir -p /data

# Set the entrypoint
ENTRYPOINT ["./go-filesystem-mcp", "/data"]
