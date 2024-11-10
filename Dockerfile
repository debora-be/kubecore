# Build stage
FROM golang:1.22-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o kubecore .

# Final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/kubecore /usr/local/bin/kubecore

# Create non-root user
RUN adduser -D -g '' kubecore
USER kubecore

# Set the entrypoint
ENTRYPOINT ["kubecore"]