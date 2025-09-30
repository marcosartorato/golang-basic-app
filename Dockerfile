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

# The application is going to listen two HTTP ports:
# - 8080 for the main application
# - 9090 for the metrics endpoint
EXPOSE 8080
EXPOSE 9090

# Run app
CMD ["./myapp"]