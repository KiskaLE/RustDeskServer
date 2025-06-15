# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go tool templ generate
	@npx @tailwindcss/cli -i cmd/api/webui/view/styles/tailwind.css -o bin/static/styles/styles.css
	@CGO_ENABLED=1 GOOS=linux go build -ldflags="-X 'main.runningInDocker=false'" -o ./bin/api ./cmd/api

# Run the application
dev:
	@echo "Starting development server..."
	@go tool templ generate --watch &
	@npx @tailwindcss/cli -i cmd/api/webui/view/styles/tailwind.css -o cmd/api/webui/view/styles/styles.css --watch &
	@go run cmd/api/main.go