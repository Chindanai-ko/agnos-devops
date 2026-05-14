# Build stage
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Final stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port 9090 for metrics
EXPOSE 9090

USER appuser

# Run the binary
CMD ["./main"]
