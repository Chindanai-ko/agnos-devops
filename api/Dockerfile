# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go.mod and download dependencies
COPY go.mod ./
# Note: go.sum is missing in the current directory, so we skip copying it for now
# RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 for the API
EXPOSE 8080

# Run the binary
CMD ["./main"]
