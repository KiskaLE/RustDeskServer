# Stage 1: Builder
FROM golang:1.24.3 AS builder

WORKDIR /app

# Install build dependencies for CGO (if necessary for sqlite)
RUN apt-get update && apt-get install -y gcc libc6-dev
RUN apt update && apt install -y nodejs npm

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go tool templ generate
RUN npm install && npx @tailwindcss/cli -i /app/cmd/api/webui/view/styles/tailwind.css -o /app/bin/static/styles/styles.css

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-X 'main.runningInDocker=true'" -o /app/bin/api /app/cmd/api

# Stage 2: Final Image
FROM debian:stable-slim

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/static ./static

CMD ["/app/api"]
EXPOSE 3000