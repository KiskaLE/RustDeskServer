# Stage 1: Builder
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO (if necessary for sqlite)
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./cmd/api

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/bin/api .
COPY --from=builder /app/db /app/db
COPY --from=builder /app/.env .

CMD ["/app/api"]
EXPOSE ${PORT:-8080}