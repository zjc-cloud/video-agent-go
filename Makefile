# Makefile for video-agent-go

.PHONY: help build run test clean docker-build docker-run dev setup deps

# Default target
help:
	@echo "Available commands:"
	@echo "  setup      - Initial project setup"
	@echo "  deps       - Download dependencies"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run - Run with Docker Compose"
	@echo "  dev        - Start development environment"

# Project setup
setup:
	@echo "Setting up project..."
	cp .env.example .env
	mkdir -p uploads/{images,audio,videos,subtitles}
	mkdir -p temp
	mkdir -p logs
	@echo "Setup complete! Please edit .env with your API keys."

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build:
	@echo "Building application..."
	go build -o bin/videoagent cmd/main.go

# Run the application
run: build
	@echo "Starting application..."
	./bin/videoagent

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf temp/*
	rm -rf uploads/*
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t video-agent-go .

# Run with Docker Compose
docker-run:
	@echo "Starting with Docker Compose..."
	docker-compose up --build

# Development environment
dev:
	@echo "Starting development environment..."
	docker-compose up --build -d mysql redis
	@echo "Waiting for database to be ready..."
	sleep 10
	go run cmd/main.go

# Install development tools
dev-tools:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run with hot reload
dev-watch: dev-tools
	air

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Database migration (placeholder)
migrate:
	@echo "Running database migrations..."
	# Add your migration commands here

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	# Add swagger generation here

# Production build
prod-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/videoagent cmd/main.go 