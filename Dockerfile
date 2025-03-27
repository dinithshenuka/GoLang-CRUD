# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bookapi ./cmd/main.go

# Final stage
FROM alpine:3.19

# Add CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata && \
    mkdir -p /app/data

WORKDIR /app

# Copy binary and data from builder stage
COPY --from=builder /app/bookapi /app/
COPY --from=builder /app/.env /app/.env

# Create non-root user for security
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app
USER appuser

# Expose port
EXPOSE 8080

# Command to run
CMD ["./bookapi"]