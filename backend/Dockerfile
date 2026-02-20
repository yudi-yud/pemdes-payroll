# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files from root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (entire project)
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o payroll-server main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/payroll-server .

# Expose port
EXPOSE 3000

# Run the application
CMD ["./payroll-server"]
