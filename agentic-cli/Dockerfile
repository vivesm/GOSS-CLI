# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
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
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$(date +%Y%m%d)-docker" -o gemini cmd/gemini/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

WORKDIR /home/appuser

# Copy the binary from builder stage
COPY --from=builder /app/gemini .

# Copy configuration example
COPY --from=builder /app/goss_config.json.example .

# Change ownership to appuser
RUN chown -R appuser:appuser /home/appuser

# Switch to non-root user
USER appuser

# Expose port (not used by CLI but for documentation)
EXPOSE 8080

# Default command
ENTRYPOINT ["./gemini"]
CMD ["--help"]