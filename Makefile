# Makefile for Bureau MLM Backend

.PHONY: help build run test clean docker-build docker-run seed-admin generate-gql

# Default target
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  seed-admin     - Seed admin user"
	@echo "  generate-gql   - Generate GraphQL code"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/bureau ./cmd/server

# Run the application
run:
	@echo "Running application..."
	go run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t bureau-mlm-backend .

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build

# Stop Docker Compose services
docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Seed admin user
seed-admin:
	@echo "Seeding admin user..."
	go run ./scripts/seed_admin.go

# Generate GraphQL code
generate-gql:
	@echo "Generating GraphQL code..."
	./scripts/generate_gql.sh

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Install tools
install-tools:
	@echo "Installing tools..."
	go install github.com/99designs/gqlgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

