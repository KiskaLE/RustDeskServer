FROM golang:1.24.3-alpine AS build
RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/api ./cmd/api
CMD ["/app/bin/api"]
EXPOSE 8080