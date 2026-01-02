.PHONY: build run test clean docker-build docker-run setup deploy

# Variables
APP_NAME=distributed-scraper
DOCKER_IMAGE=$(APP_NAME):latest
GO_VERSION=1.21

# Build the application
build:
	@echo "🔨 Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) cmd/main.go

# Run the application locally
run:
	@echo "🚀 Running $(APP_NAME)..."
	@go run cmd/main.go

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@docker-compose down -v

# Generate Swagger documentation
swagger:
	@echo "📚 Generating Swagger docs..."
	@swag init -g cmd/main.go -o docs/

# Setup development environment
setup:
	@echo "⚙️  Setting up development environment..."
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Run with Docker Compose
docker-run:
	@echo "🐳 Running with Docker Compose..."
	@docker-compose up -d

# Deploy to production
deploy:
	@echo "🚀 Deploying to production..."
	@chmod +x scripts/deploy.sh
	@./scripts/deploy.sh

# Database migrations
migrate:
	@echo "🗄️  Running database migrations..."
	@go run cmd/main.go migrate

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run

# Format code
fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...

# Security scan
security:
	@echo "🔒 Running security scan..."
	@gosec ./...

# Performance benchmark
benchmark:
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  setup         - Setup development environment"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  deploy        - Deploy to production"
	@echo "  migrate       - Run database migrations"
	@echo "  lint          - Lint code"
	@echo "  fmt           - Format code"
	@echo "  security      - Run security scan"
	@echo "  benchmark     - Run performance benchmarks"