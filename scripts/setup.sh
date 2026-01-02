#!/bin/bash

# Setup script for development environment

set -e

echo "🚀 Setting up Distributed Scraping Engine..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose."
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ .env file created. Please update it with your configuration."
fi

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod download

# Generate Swagger docs
echo "📚 Generating Swagger documentation..."
if command -v swag &> /dev/null; then
    swag init -g cmd/main.go -o docs/
else
    echo "⚠️  Swagger CLI not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"
fi

# Start PostgreSQL with Docker Compose
echo "🐘 Starting PostgreSQL database..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 10

# Run database migrations
echo "🗄️  Running database migrations..."
go run cmd/main.go &
SERVER_PID=$!
sleep 5
kill $SERVER_PID

echo "✅ Setup completed successfully!"
echo ""
echo "🎯 Next steps:"
echo "1. Update .env file with your configuration"
echo "2. Run 'docker-compose up' to start all services"
echo "3. Visit http://localhost:8080/swagger/ for API documentation"
echo "4. Visit http://localhost:8080/health for health check"