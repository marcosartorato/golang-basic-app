# Start from a small Go image
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the Go app
RUN go build -o myapp ./cmd/myapp

# Final image
FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /app/myapp .

# Expose HTTP port
EXPOSE 8080

# Run app
CMD ["./myapp"]