# Docker Patterns for Agnos DevOps

## Multi-stage Go Dockerfile
Use this template for both API and Worker services.

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Final stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

# Environment variables will be injected at runtime via K8s/Docker Compose
# ENV APP_ENV=PROD

EXPOSE 8080
CMD ["./main"]
```

## Environment Configuration
The assignment specifies DEV, UAT, and PROD environments. 
Application logic should read environment variables (e.g., `APP_ENV`, `DB_URL`) to adjust behavior.
In Docker, these can be set via `--env` or `env_file`.
